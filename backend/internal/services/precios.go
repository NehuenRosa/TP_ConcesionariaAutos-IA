package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Fuente de los valores de referencia.
const fuentePrecios = "Guía Oficial de Precios de la CCA (Cámara del Comercio Automotor)"

// Errores del servicio de precios.
var (
	// ErrPrecioNoEncontrado indica que no hay valor de referencia para el vehículo.
	ErrPrecioNoEncontrado = errors.New("no se encontró un valor de referencia para ese vehículo")
	// ErrPrecioNoDisponible indica que la fuente de precios no respondió.
	ErrPrecioNoDisponible = errors.New("la fuente de precios no está disponible")
)

// ReferenciaPrecio es el valor de mercado real de un vehículo según la fuente oficial.
type ReferenciaPrecio struct {
	Marca     string
	Modelo    string
	Version   string
	Anio      int
	PrecioUSD float64
	PrecioARS float64
	Fuente    string
}

// ServicioPrecios devuelve valores de mercado de vehículos usados según la
// Guía Oficial de la CCA, a través de la API pública de ArgAutos.
type ServicioPrecios interface {
	// Buscar devuelve el valor de referencia para una marca, modelo y año.
	Buscar(ctx context.Context, marca string, modelo string, anio int) (ReferenciaPrecio, error)
}

// servicioPrecios implementa ServicioPrecios consultando la API de ArgAutos
// con una caché en memoria para no superar el límite de pedidos anónimos.
type servicioPrecios struct {
	baseURL  string
	cliente  *http.Client
	mu       sync.Mutex
	cache    map[string]entradaPrecio
	duracion time.Duration
}

type entradaPrecio struct {
	referencia ReferenciaPrecio
	vencimiento time.Time
}

// NuevoServicioPrecios crea el servicio de precios apuntando a la API dada.
func NuevoServicioPrecios(baseURL string) ServicioPrecios {
	if baseURL == "" {
		baseURL = "https://argautos.com/api/v1"
	}
	return &servicioPrecios{
		baseURL: strings.TrimRight(baseURL, "/"),
		cliente: &http.Client{Timeout: 15 * time.Second},
		cache:   make(map[string]entradaPrecio),
		duracion: 24 * time.Hour,
	}
}

// Buscar consulta la API y guarda el resultado en caché.
func (s *servicioPrecios) Buscar(ctx context.Context, marca string, modelo string, anio int) (ReferenciaPrecio, error) {
	marca = strings.TrimSpace(marca)
	modelo = strings.TrimSpace(modelo)
	if marca == "" || modelo == "" {
		return ReferenciaPrecio{}, ErrPrecioNoEncontrado
	}

	clave := strings.ToLower(fmt.Sprintf("%s|%s|%d", marca, modelo, anio))
	s.mu.Lock()
	if entrada, ok := s.cache[clave]; ok && time.Now().Before(entrada.vencimiento) {
		s.mu.Unlock()
		return entrada.referencia, nil
	}
	s.mu.Unlock()

	referencia, err := s.consultar(ctx, marca, modelo, anio)
	if err != nil {
		return ReferenciaPrecio{}, err
	}

	s.mu.Lock()
	s.cache[clave] = entradaPrecio{referencia: referencia, vencimiento: time.Now().Add(s.duracion)}
	s.mu.Unlock()
	return referencia, nil
}

// consultar realiza las llamadas de búsqueda y valuación a la API.
func (s *servicioPrecios) consultar(ctx context.Context, marca string, modelo string, anio int) (ReferenciaPrecio, error) {
	versionID, version, precioUSD, anioReferencia, err := s.buscarVersion(ctx, marca, modelo, anio)
	if err != nil {
		return ReferenciaPrecio{}, err
	}

	precioARS, err := s.obtenerPrecioARS(ctx, versionID, anioReferencia)
	if err != nil {
		return ReferenciaPrecio{}, err
	}

	return ReferenciaPrecio{
		Marca:     strings.ToUpper(marca[:1]) + marca[1:],
		Modelo:    strings.ToUpper(modelo[:1]) + modelo[1:],
		Version:   version,
		Anio:      anioReferencia,
		PrecioUSD: precioUSD,
		PrecioARS: precioARS,
		Fuente:    fuentePrecios,
	}, nil
}

// buscarVersion llama a /search, elige la versión cuyo año referencia se acerca
// más al año del vehículo y devuelve el precio en USD de esa versión.
func (s *servicioPrecios) buscarVersion(ctx context.Context, marca string, modelo string, anio int) (int, string, float64, int, error) {
	consulta := url.QueryEscape(strings.TrimSpace(marca) + " " + strings.TrimSpace(modelo))
	destino := fmt.Sprintf("%s/search?q=%s", s.baseURL, consulta)

	var respuesta struct {
		Datos []struct {
			VersionID  int     `json:"version_id"`
			Version    string  `json:"version"`
			Precio     string  `json:"price"`
			PrecioAnio int     `json:"price_year"`
		} `json:"data"`
	}
	if err := s.obtenerJSON(ctx, destino, &respuesta); err != nil {
		return 0, "", 0, 0, err
	}
	if len(respuesta.Datos) == 0 {
		return 0, "", 0, 0, ErrPrecioNoEncontrado
	}

	mejor := -1
	for i := range respuesta.Datos {
		entrada := &respuesta.Datos[i]
		if entrada.VersionID == 0 || entrada.Precio == "" {
			continue
		}
		if mejor == -1 || cercaniaAnio(entrada.PrecioAnio, anio) < cercaniaAnio(respuesta.Datos[mejor].PrecioAnio, anio) {
			mejor = i
		}
	}
	if mejor == -1 {
		return 0, "", 0, 0, ErrPrecioNoEncontrado
	}

	entrada := &respuesta.Datos[mejor]
	precioUSD, err := parseFloat(entrada.Precio)
	if err != nil {
		return 0, "", 0, 0, ErrPrecioNoEncontrado
	}
	return entrada.VersionID, entrada.Version, precioUSD, entrada.PrecioAnio, nil
}

// cercaniaAnio mide qué tan lejos está el año de la versión del año buscado.
// Un año en 0 indica 0km (lo más nuevo posible).
func cercaniaAnio(anioVersion int, anioBuscado int) int {
	if anioBuscado <= 0 {
		return 0
	}
	if anioVersion <= 0 {
		return math.MaxInt
	}
	if anioVersion > anioBuscado {
		return anioVersion - anioBuscado
	}
	return anioBuscado - anioVersion
}

// obtenerPrecioARS llama a /versions/{id}/valuations en pesos y elige el año
// más cercano al año de referencia.
func (s *servicioPrecios) obtenerPrecioARS(ctx context.Context, versionID int, anio int) (float64, error) {
	destino := fmt.Sprintf("%s/versions/%d/valuations?currency=ars", s.baseURL, versionID)

	var respuesta struct {
		Datos []struct {
			Anio  int    `json:"year"`
			Precio string `json:"price"`
		} `json:"data"`
	}
	if err := s.obtenerJSON(ctx, destino, &respuesta); err != nil {
		return 0, err
	}
	if len(respuesta.Datos) == 0 {
		return 0, ErrPrecioNoEncontrado
	}

	mejor := -1
	for i := range respuesta.Datos {
		entrada := &respuesta.Datos[i]
		if entrada.Precio == "" {
			continue
		}
		if mejor == -1 || cercaniaAnio(entrada.Anio, anio) < cercaniaAnio(respuesta.Datos[mejor].Anio, anio) {
			mejor = i
		}
	}
	if mejor == -1 {
		return 0, ErrPrecioNoEncontrado
	}
	return parseFloat(respuesta.Datos[mejor].Precio)
}

// obtenerJSON ejecuta un GET y decodifica la respuesta JSON.
func (s *servicioPrecios) obtenerJSON(ctx context.Context, destino string, destinoJSON any) error {
	peticion, err := http.NewRequestWithContext(ctx, http.MethodGet, destino, nil)
	if err != nil {
		return fmt.Errorf("crear pedido de precios: %w", err)
	}
	peticion.Header.Set("Accept", "application/json")
	peticion.Header.Set("User-Agent", "concesionaria-api/1.0")

	respuesta, err := s.cliente.Do(peticion)
	if err != nil {
		return ErrPrecioNoDisponible
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode == http.StatusNotFound || respuesta.StatusCode == http.StatusUnprocessableEntity {
		return ErrPrecioNoEncontrado
	}
	if respuesta.StatusCode != http.StatusOK {
		return ErrPrecioNoDisponible
	}

	cuerpo, err := io.ReadAll(io.LimitReader(respuesta.Body, 2*1024*1024))
	if err != nil {
		return ErrPrecioNoDisponible
	}
	if err := json.Unmarshal(cuerpo, destinoJSON); err != nil {
		return ErrPrecioNoDisponible
	}
	return nil
}

// parseFloat convierte un precio numérico en formato texto o número.
func parseFloat(valor string) (float64, error) {
	valor = strings.TrimSpace(valor)
	if valor == "" {
		return 0, fmt.Errorf("precio vacío")
	}
	var numero float64
	if err := json.Unmarshal([]byte(valor), &numero); err != nil {
		return 0, err
	}
	return numero, nil
}
