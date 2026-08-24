// Package googleid verifica ID tokens emitidos por Google Identity Services
// contra los certificados públicos de Google (JWKS), sin depender de
// bibliotecas externas: usa github.com/golang-jwt/jwt/v5 y la stdlib.
package googleid

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Emisores válidos del ID token de Google.
var emisoresValidos = []string{"accounts.google.com", "https://accounts.google.com"}

const (
	// urlCertificados es el JWKS público de Google.
	urlCertificados = "https://www.googleapis.com/oauth2/v3/certs"
	// ttlClaves es cuánto tiempo se cachean los certificados públicos.
	ttlClaves = time.Hour
	// margenExpiracion tolera pequeñas derivas de reloj al validar exp/nbf.
	margenExpiracion = 2 * time.Minute
)

// Errores del paquete.
var (
	// ErrTokenInvalido indica que el ID token no es válido o no confiable.
	ErrTokenInvalido = errors.New("credencial de Google inválida")
	// ErrCertificadosNoDisponibles indica que no se pudieron obtener las
	// claves públicas de Google.
	ErrCertificadosNoDisponibles = errors.New("certificados de Google no disponibles")
)

// Identidad son los datos verificados de la cuenta de Google.
type Identidad struct {
	Sub    string
	Email  string
	Nombre string
}

// reclamosGoogle son los claims relevantes del ID token de Google.
type reclamosGoogle struct {
	Email           string `json:"email"`
	EmailVerificado bool   `json:"email_verified"`
	Nombre          string `json:"name"`
	jwt.RegisteredClaims
}

// jwkGoogle es una clave pública en formato JWK del JWKS de Google.
type jwkGoogle struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// jwks es el documento completo de claves públicas.
type jwks struct {
	Claves []jwkGoogle `json:"keys"`
}

// Verificador valida ID tokens de Google para una audiencia (client ID)
// dada, cacheando el JWKS en memoria.
type Verificador struct {
	clienteID   string
	urlJWKS     string
	clienteHTTP *http.Client

	mu            sync.Mutex
	claves        map[string]*rsa.PublicKey
	validasHasta  time.Time
}

// NuevoVerificador crea un verificador para la audiencia indicada. Un
// clienteID vacío deja al verificador inutilizable (ningún token va a
// coincidir con la audiencia esperada).
func NuevoVerificador(clienteID string) *Verificador {
	return &Verificador{
		clienteID:   clienteID,
		urlJWKS:     urlCertificados,
		clienteHTTP: &http.Client{Timeout: 10 * time.Second},
		claves:      make(map[string]*rsa.PublicKey),
	}
}

// Verificar parsea y valida el ID token: firma RS256 contra los
// certificados de Google, emisor, audiencia, expiración (con margen) y
// email verificado. Devuelve la identidad verificada o un error.
func (v *Verificador) Verificar(ctx context.Context, idToken string) (*Identidad, error) {
	if strings.TrimSpace(idToken) == "" {
		return nil, ErrTokenInvalido
	}

	reclamos := &reclamosGoogle{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithLeeway(margenExpiracion),
	)
	parsed, err := parser.ParseWithClaims(idToken, reclamos, v.clavePorKID(ctx))
	if err != nil || !parsed.Valid {
		return nil, ErrTokenInvalido
	}

	if !emisorValido(reclamos.Issuer) ||
		reclamos.Subject == "" ||
		reclamos.Email == "" ||
		!reclamos.EmailVerificado ||
		!audienciaValida(reclamos.Audience, v.clienteID) {
		return nil, ErrTokenInvalido
	}

	nombre := strings.TrimSpace(reclamos.Nombre)
	if nombre == "" {
		if usuario, _, ok := strings.Cut(reclamos.Email, "@"); ok {
			nombre = usuario
		}
	}
	return &Identidad{Sub: reclamos.Subject, Email: reclamos.Email, Nombre: nombre}, nil
}

// clavePorKID devuelve la clave pública correspondiente al `kid` del token,
// refrescando el caché si es desconocido o está vencido.
func (v *Verificador) clavePorKID(ctx context.Context) jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			return nil, ErrTokenInvalido
		}
		if clave, ok := v.clave(kid); ok {
			return clave, nil
		}
		if err := v.refrescarClaves(ctx); err != nil {
			return nil, err
		}
		clave, ok := v.clave(kid)
		if !ok {
			return nil, ErrTokenInvalido
		}
		return clave, nil
	}
}

// clave consulta el caché de claves bajo mutex.
func (v *Verificador) clave(kid string) (*rsa.PublicKey, bool) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if time.Now().After(v.validasHasta) {
		return nil, false
	}
	clave, ok := v.claves[kid]
	return clave, ok
}

// refrescarClaves descarga el JWKS de Google y reemplaza el caché.
func (v *Verificador) refrescarClaves(ctx context.Context) error {
	pedido, err := http.NewRequestWithContext(ctx, http.MethodGet, v.urlJWKS, nil)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCertificadosNoDisponibles, err)
	}
	respuesta, err := v.clienteHTTP.Do(pedido)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCertificadosNoDisponibles, err)
	}
	defer respuesta.Body.Close()
	if respuesta.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: estado %d", ErrCertificadosNoDisponibles, respuesta.StatusCode)
	}
	cuerpo, err := io.ReadAll(io.LimitReader(respuesta.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("%w: %v", ErrCertificadosNoDisponibles, err)
	}

	var documento jwks
	if err := json.Unmarshal(cuerpo, &documento); err != nil {
		return fmt.Errorf("%w: %v", ErrCertificadosNoDisponibles, err)
	}
	claves := make(map[string]*rsa.PublicKey, len(documento.Claves))
	for _, jwk := range documento.Claves {
		if jwk.Kty != "RSA" || jwk.Kid == "" {
			continue
		}
		clave, err := claveRSA(jwk)
		if err != nil {
			continue
		}
		claves[jwk.Kid] = clave
	}
	if len(claves) == 0 {
		return ErrCertificadosNoDisponibles
	}

	v.mu.Lock()
	v.claves = claves
	v.validasHasta = time.Now().Add(ttlClaves)
	v.mu.Unlock()
	return nil
}

// emisorValido informa si el emisor es uno de los aceptados de Google.
func emisorValido(emisor string) bool {
	for _, valido := range emisoresValidos {
		if emisor == valido {
			return true
		}
	}
	return false
}

// audienciaValida informa si el client ID está entre las audiencias del token.
func audienciaValida(audiencias jwt.ClaimStrings, clienteID string) bool {
	if clienteID == "" {
		return false
	}
	for _, audiencia := range audiencias {
		if audiencia == clienteID {
			return true
		}
	}
	return false
}

// claveRSA convierte un JWK (n, e en base64url) a una clave RSA.
func claveRSA(jwk jwkGoogle) (*rsa.PublicKey, error) {
	modulo, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("decodificar módulo del JWK: %w", err)
	}
	exponenteBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("decodificar exponente del JWK: %w", err)
	}
	exponente := 0
	for _, b := range exponenteBytes {
		exponente = exponente<<8 + int(b)
	}
	return &rsa.PublicKey{N: new(big.Int).SetBytes(modulo), E: exponente}, nil
}
