package config

import (
	"log"
	"os"
)

// Configuracion agrupa los valores de configuración cargados del entorno.
type Configuracion struct {
	Host        string
	Puerto      string
	BDHost      string
	BDPuerto    string
	BDUsuario   string
	BDPassword  string
	BDNombre    string
	BDSSL       string
	JWTSecreto  string
	OrigenesCORS string
}

// Cargar lee las variables de entorno y devuelve la configuración de la API.
func Cargar() Configuracion {
	return Configuracion{
		Host:         obtener("HOST_API", "0.0.0.0"),
		Puerto:       obtener("PUERTO_API", "8080"),
		BDHost:       obtener("BD_HOST", "localhost"),
		BDPuerto:     obtener("BD_PUERTO", "5432"),
		BDUsuario:    obtener("BD_USUARIO", "concesionaria"),
		BDPassword:   obtener("BD_PASSWORD", "concesionaria"),
		BDNombre:     obtener("BD_NOMBRE", "concesionaria"),
		BDSSL:        obtener("BD_SSL", "disable"),
		JWTSecreto:   obtener("JWT_SECRETO", "cambiar-en-produccion"),
		OrigenesCORS: obtener("CORS_ORIGENES", "*"),
	}
}

func obtener(clave string, valorPorDefecto string) string {
	valor, existe := os.LookupEnv(clave)
	if !existe || valor == "" {
		if valorPorDefecto == "" {
			log.Printf("Variable de entorno %s no definida", clave)
		}
		return valorPorDefecto
	}
	return valor
}
