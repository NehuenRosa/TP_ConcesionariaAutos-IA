package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInicialMayuscula(t *testing.T) {
	casos := []struct {
		entrada string
		esperado string
	}{
		{entrada: "volkswagen", esperado: "Volkswagen"},
		{entrada: " toyota ", esperado: "Toyota"},
		{entrada: "gol", esperado: "Gol"},
		{entrada: "škoda", esperado: "Škoda"},
		{entrada: "ŠKODA", esperado: "ŠKODA"},
		{entrada: "", esperado: ""},
	}
	for _, caso := range casos {
		if resultado := inicialMayuscula(caso.entrada); resultado != caso.esperado {
			t.Errorf("inicialMayuscula(%q) = %q, se esperaba %q", caso.entrada, resultado, caso.esperado)
		}
	}
}

func TestCercaniaAnio(t *testing.T) {
	casos := []struct {
		anioVersion int
		anioBuscado int
		esperado    int
	}{
		{anioVersion: 2020, anioBuscado: 2020, esperado: 0},
		{anioVersion: 2018, anioBuscado: 2020, esperado: 2},
		{anioVersion: 2022, anioBuscado: 2020, esperado: 2},
		{anioVersion: 0, anioBuscado: 2020, esperado: math.MaxInt},
		{anioVersion: 0, anioBuscado: 0, esperado: 0},
		{anioVersion: 2020, anioBuscado: 0, esperado: 0},
	}
	for _, caso := range casos {
		if resultado := cercaniaAnio(caso.anioVersion, caso.anioBuscado); resultado != caso.esperado {
			t.Errorf("cercaniaAnio(%d, %d) = %d, se esperaba %d", caso.anioVersion, caso.anioBuscado, resultado, caso.esperado)
		}
	}
}

func TestParseFloat(t *testing.T) {
	if numero, err := parseFloat("25000"); err != nil || numero != 25000 {
		t.Errorf("parseFloat(\"25000\") = %v, %v; se esperaba 25000", numero, err)
	}
	if numero, err := parseFloat("12345.67"); err != nil || numero != 12345.67 {
		t.Errorf("parseFloat(\"12345.67\") = %v, %v; se esperaba 12345.67", numero, err)
	}
	if _, err := parseFloat(""); err == nil {
		t.Error("parseFloat(\"\") debería fallar")
	}
	if _, err := parseFloat("no-es-un-precio"); err == nil {
		t.Error("parseFloat(\"no-es-un-precio\") debería fallar")
	}
}

func TestServicioPreciosBuscarConCache(t *testing.T) {
	pedidosSearch := 0
	pedidosValuaciones := 0

	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/search":
			pedidosSearch++
			fmt.Fprint(w, `{"data":[
				{"version_id":1,"version":"2.0 TSI","price":"25000","price_year":2022},
				{"version_id":2,"version":"2.5","price":"20000","price_year":2018}
			]}`)
		case "/versions/1/valuations":
			pedidosValuaciones++
			fmt.Fprint(w, `{"data":[
				{"year":2018,"price":"5000000"},
				{"year":2022,"price":"8000000"}
			]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer servidor.Close()

	servicio := NuevoServicioPrecios(servidor.URL)
	referencia, err := servicio.Buscar(context.Background(), "volkswagen", "gol", 2020)
	if err != nil {
		t.Fatalf("Buscar devolvió error: %v", err)
	}
	if referencia.Marca != "Volkswagen" {
		t.Errorf("marca esperada %q, obtenida %q", "Volkswagen", referencia.Marca)
	}
	if referencia.Modelo != "Gol" {
		t.Errorf("modelo esperado %q, obtenido %q", "Gol", referencia.Modelo)
	}
	if referencia.PrecioUSD != 25000 {
		t.Errorf("precio USD esperado 25000, obtenido %v", referencia.PrecioUSD)
	}
	if referencia.PrecioARS != 8000000 {
		t.Errorf("precio ARS esperado 8000000, obtenido %v", referencia.PrecioARS)
	}
	if referencia.Anio != 2022 {
		t.Errorf("año esperado 2022, obtenido %d", referencia.Anio)
	}
	if referencia.Fuente != fuentePrecios {
		t.Errorf("fuente esperada %q, obtenida %q", fuentePrecios, referencia.Fuente)
	}

	if _, err := servicio.Buscar(context.Background(), "volkswagen", "gol", 2020); err != nil {
		t.Fatalf("segunda búsqueda (caché) devolvió error: %v", err)
	}
	if pedidosSearch != 1 {
		t.Errorf("la caché no funcionó: se esperaba 1 pedido a /search, hubo %d", pedidosSearch)
	}
	if pedidosValuaciones != 1 {
		t.Errorf("la caché no funcionó: se esperaba 1 pedido de valuaciones, hubo %d", pedidosValuaciones)
	}
}

func TestServicioPreciosEligeVersionCercana(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/search" {
			// Un 0km (price_year 0) no debe ganarle a una versión con año real.
			fmt.Fprint(w, `{"data":[
				{"version_id":1,"version":"0km","price":"30000","price_year":0},
				{"version_id":2,"version":"Usado","price":"18000","price_year":2021}
			]}`)
			return
		}
		fmt.Fprint(w, `{"data":[{"year":2021,"price":"4500000"}]}`)
	}))
	defer servidor.Close()

	servicio := NuevoServicioPrecios(servidor.URL)
	referencia, err := servicio.Buscar(context.Background(), "ford", "focus", 2020)
	if err != nil {
		t.Fatalf("Buscar devolvió error: %v", err)
	}
	if referencia.Version != "Usado" {
		t.Errorf("se esperaba elegir la versión \"Usado\", se eligió %q", referencia.Version)
	}
}

func TestServicioPreciosSinResultados(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[]}`)
	}))
	defer servidor.Close()

	servicio := NuevoServicioPrecios(servidor.URL)
	_, err := servicio.Buscar(context.Background(), "marca", "modelo", 2020)
	if !errors.Is(err, ErrPrecioNoEncontrado) {
		t.Errorf("se esperaba ErrPrecioNoEncontrado, se obtuvo %v", err)
	}
}

func TestServicioPreciosDatosVacios(t *testing.T) {
	servicio := NuevoServicioPrecios("")
	_, err := servicio.Buscar(context.Background(), "  ", "gol", 2020)
	if !errors.Is(err, ErrPrecioNoEncontrado) {
		t.Errorf("se esperaba ErrPrecioNoEncontrado con marca vacía, se obtuvo %v", err)
	}
}
