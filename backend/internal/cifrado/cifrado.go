package cifrado

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// Cifrado cotiza mensajes sensibles (cotizaciones) en reposo con AES-256-GCM.
// El contenido se guarda cifrado en la base y se descifra al leerlo.
type Cifrador interface {
	// Cifrar devuelve el texto cifrado codificado en base64 (nonce + dato).
	Cifrar(texto string) (string, error)
	// Descifrar recupera el texto original a partir del valor cifrado.
	Descifrar(cifrado string) (string, error)
}

// cifradorAES implementa Cifrador con AES-256-GCM.
type cifradorAES struct {
	gcm cipher.AEAD
}

// NuevoCifrador crea un cifrador a partir de una clave secreta. La clave se
// normaliza a 32 bytes (SHA-256) para usar AES-256.
func NuevoCifrador(claveSecreta string) (Cifrador, error) {
	resumen := sha256.Sum256([]byte(claveSecreta))
	bloque, err := aes.NewCipher(resumen[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(bloque)
	if err != nil {
		return nil, err
	}
	return &cifradorAES{gcm: gcm}, nil
}

// Cifrar cifra el texto con un nonce nuevo y devuelve nonce + cifrado en base64.
func (c *cifradorAES) Cifrar(texto string) (string, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	datos := c.gcm.Seal(nonce, nonce, []byte(texto), nil)
	return base64.StdEncoding.EncodeToString(datos), nil
}

// Descifrar invierte Cifrar recuperando el texto original.
func (c *cifradorAES) Descifrar(cifrado string) (string, error) {
	if cifrado == "" {
		return "", nil
	}
	datos, err := base64.StdEncoding.DecodeString(cifrado)
	if err != nil {
		return "", err
	}
	nonce, dato := datos[:c.gcm.NonceSize()], datos[c.gcm.NonceSize():]
	texto, err := c.gcm.Open(nil, nonce, dato, nil)
	if err != nil {
		return "", errors.New("no se pudo descifrar el contenido")
	}
	return string(texto), nil
}