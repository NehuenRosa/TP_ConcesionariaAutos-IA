package config

import (
	"log"
	"os"
)

// Configuracion agrupa los valores de configuración cargados del entorno.
type Configuracion struct {
	Host         string
	Puerto       string
	BDURL        string
	BDHost       string
	BDPuerto     string
	BDUsuario    string
	BDPassword   string
	BDNombre     string
	BDSSL        string
	JWTSecreto   string
	OrigenesCORS string
	// ProveedorLLM indica el proveedor del LLM: "googleai" (nube, por defecto
	// si hay GOOGLE_API_KEY) u "ollama" (local).
	ProveedorLLM string
	GoogleAIKey  string
	OllamaURL    string
	ModeloChatbot string
	ModeloVision string
	ArgAutosURL  string
}

// Cargar lee las variables de entorno y devuelve la configuración de la API.
func Cargar() Configuracion {
	proveedor := obtener("PROVEEDOR_LLM", "")
	claveGoogleAI := obtener("GOOGLE_API_KEY", "")
	// Si no se indica proveedor, se elige googleai cuando hay GOOGLE_API_KEY y
	// ollama como respaldo local.
	if proveedor == "" {
		if claveGoogleAI != "" {
			proveedor = "googleai"
		} else {
			proveedor = "ollama"
		}
	}
	// Los modelos por defecto dependen del proveedor: Gemini soporta texto +
	// visión; Ollama usa llama3 para el chat y minicpm-v para la tasación.
	// gemini-flash-lite-latest es el alias estable del Gemini gratis más
	// liviano (1M de contexto, texto + visión, disponible para cuentas nuevas).
	modeloChatbotPorDefecto := "llama3"
	modeloVisionPorDefecto := "minicpm-v"
	if proveedor == "googleai" {
		modeloChatbotPorDefecto = "gemini-flash-lite-latest"
		modeloVisionPorDefecto = "gemini-flash-lite-latest"
	}
	return Configuracion{
		Host:         obtener("HOST_API", "0.0.0.0"),
		// PORT lo inyecta Render en la nube; localmente sigue PUERTO_API.
		Puerto:       obtener("PUERTO_API", obtener("PORT", "8080")),
		BDURL:        obtener("BD_URL", ""),
		BDHost:       obtener("BD_HOST", "localhost"),
		BDPuerto:     obtener("BD_PUERTO", "5432"),
		BDUsuario:    obtener("BD_USUARIO", "concesionaria"),
		BDPassword:   obtener("BD_PASSWORD", "concesionaria"),
		BDNombre:     obtener("BD_NOMBRE", "concesionaria"),
		BDSSL:        obtener("BD_SSL", "disable"),
		JWTSecreto:    obtener("JWT_SECRETO", "cambiar-en-produccion"),
		OrigenesCORS:  obtener("CORS_ORIGENES", "*"),
		ProveedorLLM:  proveedor,
		GoogleAIKey:   claveGoogleAI,
		OllamaURL:     obtener("OLLAMA_URL", "http://localhost:11434"),
		ModeloChatbot: obtener("MODELO_CHATBOT", modeloChatbotPorDefecto),
		ModeloVision:  obtener("MODELO_VISION", modeloVisionPorDefecto),
		ArgAutosURL:   obtener("ARGAUTOS_URL", "https://argautos.com/api/v1"),
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
