package googleid

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// parDeClaves genera un par de claves RSA para firmar tokens de prueba.
func parDeClaves(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	clave, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generar clave RSA: %v", err)
	}
	return clave
}

// servidorJWKS levanta un servidor que expone el JWKS con la clave dada.
func servidorJWKS(t *testing.T, kid string, clave *rsa.PublicKey) *httptest.Server {
	t.Helper()
	documento := jwks{Claves: []jwkGoogle{{
		Kty: "RSA",
		Kid: kid,
		N:   base64.RawURLEncoding.EncodeToString(clave.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(clave.E)).Bytes()),
	}}}
	cuerpo, err := json.Marshal(documento)
	if err != nil {
		t.Fatalf("serializar JWKS: %v", err)
	}
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(cuerpo)
	}))
	t.Cleanup(servidor.Close)
	return servidor
}

// firmar genera un ID token de prueba con los reclamos indicados.
func firmar(t *testing.T, clave *rsa.PrivateKey, kid string, reclamos jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, reclamos)
	token.Header["kid"] = kid
	firmado, err := token.SignedString(clave)
	if err != nil {
		t.Fatalf("firmar token: %v", err)
	}
	return firmado
}

// reclamosValidos son los claims de un ID token de Google correcto.
func reclamosValidos(audiencia string) jwt.MapClaims {
	ahora := time.Now()
	return jwt.MapClaims{
		"iss":            "https://accounts.google.com",
		"aud":            audiencia,
		"sub":            "1234567890",
		"email":          "cliente@ejemplo.com",
		"email_verified": true,
		"name":           "Cliente Ejemplo",
		"exp":            ahora.Add(10 * time.Minute).Unix(),
		"iat":            ahora.Unix(),
	}
}

const clienteIDPrueba = "client-id-prueba.apps.googleusercontent.com"

func TestVerificarTokenValido(t *testing.T) {
	clave := parDeClaves(t)
	servidor := servidorJWKS(t, "kid-1", &clave.PublicKey)
	verificador := NuevoVerificador(clienteIDPrueba)
	verificador.urlJWKS = servidor.URL

	identidad, err := verificador.Verificar(context.Background(), firmar(t, clave, "kid-1", reclamosValidos(clienteIDPrueba)))
	if err != nil {
		t.Fatalf("Verificar() error = %v, esperaba identidad", err)
	}
	if identidad.Email != "cliente@ejemplo.com" || identidad.Sub != "1234567890" || identidad.Nombre != "Cliente Ejemplo" {
		t.Fatalf("identidad incompleta: %+v", identidad)
	}
}

func TestVerificarFirmaInvalida(t *testing.T) {
	claveCorrecta := parDeClaves(t)
	claveAjena := parDeClaves(t)
	servidor := servidorJWKS(t, "kid-1", &claveCorrecta.PublicKey)
	verificador := NuevoVerificador(clienteIDPrueba)
	verificador.urlJWKS = servidor.URL

	_, err := verificador.Verificar(context.Background(), firmar(t, claveAjena, "kid-1", reclamosValidos(clienteIDPrueba)))
	if err == nil {
		t.Fatal("Verificar() esperaba error por firma inválida")
	}
}

func TestVerificarAudienciaIncorrecta(t *testing.T) {
	clave := parDeClaves(t)
	servidor := servidorJWKS(t, "kid-1", &clave.PublicKey)
	verificador := NuevoVerificador(clienteIDPrueba)
	verificador.urlJWKS = servidor.URL

	reclamos := reclamosValidos("otra-app.apps.googleusercontent.com")
	_, err := verificador.Verificar(context.Background(), firmar(t, clave, "kid-1", reclamos))
	if err == nil {
		t.Fatal("Verificar() esperaba error por audiencia incorrecta")
	}
}

func TestVerificarTokenExpirado(t *testing.T) {
	clave := parDeClaves(t)
	servidor := servidorJWKS(t, "kid-1", &clave.PublicKey)
	verificador := NuevoVerificador(clienteIDPrueba)
	verificador.urlJWKS = servidor.URL

	reclamos := reclamosValidos(clienteIDPrueba)
	reclamos["exp"] = time.Now().Add(-time.Hour).Unix()
	_, err := verificador.Verificar(context.Background(), firmar(t, clave, "kid-1", reclamos))
	if err == nil {
		t.Fatal("Verificar() esperaba error por token expirado")
	}
}

func TestVerificarEmailSinVerificar(t *testing.T) {
	clave := parDeClaves(t)
	servidor := servidorJWKS(t, "kid-1", &clave.PublicKey)
	verificador := NuevoVerificador(clienteIDPrueba)
	verificador.urlJWKS = servidor.URL

	reclamos := reclamosValidos(clienteIDPrueba)
	reclamos["email_verified"] = false
	_, err := verificador.Verificar(context.Background(), firmar(t, clave, "kid-1", reclamos))
	if err == nil {
		t.Fatal("Verificar() esperaba error por email sin verificar")
	}
}

func TestVerificarEmisorDesconocido(t *testing.T) {
	clave := parDeClaves(t)
	servidor := servidorJWKS(t, "kid-1", &clave.PublicKey)
	verificador := NuevoVerificador(clienteIDPrueba)
	verificador.urlJWKS = servidor.URL

	reclamos := reclamosValidos(clienteIDPrueba)
	reclamos["iss"] = "https://malicioso.ejemplo"
	_, err := verificador.Verificar(context.Background(), firmar(t, clave, "kid-1", reclamos))
	if err == nil {
		t.Fatal("Verificar() esperaba error por emisor desconocido")
	}
}

func TestVerificarAlgoritmoInesperado(t *testing.T) {
	clave := parDeClaves(t)
	servidor := servidorJWKS(t, "kid-1", &clave.PublicKey)
	verificador := NuevoVerificador(clienteIDPrueba)
	verificador.urlJWKS = servidor.URL

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, reclamosValidos(clienteIDPrueba))
	token.Header["kid"] = "kid-1"
	firmado, err := token.SignedString([]byte("secreto"))
	if err != nil {
		t.Fatalf("firmar token HS256: %v", err)
	}
	if _, err := verificador.Verificar(context.Background(), firmado); err == nil {
		t.Fatal("Verificar() esperaba error por algoritmo no permitido")
	}
}

func TestVerificarCredencialVacia(t *testing.T) {
	verificador := NuevoVerificador(clienteIDPrueba)
	if _, err := verificador.Verificar(context.Background(), ""); err == nil {
		t.Fatal("Verificar() esperaba error con credencial vacía")
	}
}

func TestVerificarKidDesconocidoRefrescaCache(t *testing.T) {
	// El primer JWKS solo conoce kid-viejo; tras refrescar aparece kid-nuevo.
	claveVieja := parDeClaves(t)
	claveNueva := parDeClaves(t)
	pedidos := 0
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pedidos++
		kidActual := "kid-viejo"
		claveActual := &claveVieja.PublicKey
		if pedidos > 1 {
			kidActual = "kid-nuevo"
			claveActual = &claveNueva.PublicKey
		}
		documento := jwks{Claves: []jwkGoogle{{
			Kty: "RSA",
			Kid: kidActual,
			N:   base64.RawURLEncoding.EncodeToString(claveActual.N.Bytes()),
			E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(claveActual.E)).Bytes()),
		}}}
		cuerpo, _ := json.Marshal(documento)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(cuerpo)
	}))
	defer servidor.Close()

	verificador := NuevoVerificador(clienteIDPrueba)
	verificador.urlJWKS = servidor.URL

	// Primero se valida un token con el kid viejo (carga inicial).
	if _, err := verificador.Verificar(context.Background(), firmar(t, claveVieja, "kid-viejo", reclamosValidos(clienteIDPrueba))); err != nil {
		t.Fatalf("primera validación: %v", err)
	}
	// Un token con un kid nuevo fuerza el refresco del caché.
	if _, err := verificador.Verificar(context.Background(), firmar(t, claveNueva, "kid-nuevo", reclamosValidos(clienteIDPrueba))); err != nil {
		t.Fatalf("validación tras refresco de caché: %v", err)
	}
}
