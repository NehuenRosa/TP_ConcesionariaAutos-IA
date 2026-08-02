package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Reclamos contiene los datos del usuario incluidos en el token JWT.
type Reclamos struct {
	UsuarioID uint   `json:"usuario_id"`
	Rol       string `json:"rol"`
	jwt.RegisteredClaims
}

// Generar firma un token JWT para el usuario indicado con una duración dada.
func Generar(usuarioID uint, rol string, secreto string, duracion time.Duration) (string, error) {
	reclamos := Reclamos{
		UsuarioID: usuarioID,
		Rol:       rol,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", usuarioID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracion)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, reclamos)
	firmado, err := tokenJWT.SignedString([]byte(secreto))
	if err != nil {
		return "", fmt.Errorf("firmar token JWT: %w", err)
	}
	return firmado, nil
}

// Validar verifica la firma y la expiración de un token y devuelve sus reclamos.
func Validar(tokenString string, secreto string) (*Reclamos, error) {
	reclamos := &Reclamos{}
	parsed, err := jwt.ParseWithClaims(tokenString, reclamos, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método de firma inesperado")
		}
		return []byte(secreto), nil
	})
	if err != nil {
		return nil, fmt.Errorf("validar token JWT: %w", err)
	}
	if !parsed.Valid {
		return nil, errors.New("token JWT inválido")
	}
	return reclamos, nil
}
