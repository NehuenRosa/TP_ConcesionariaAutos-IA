package services

import (
	"context"
	"testing"
	"time"

	"concesionaria/backend/internal/repositories"
)

// fakeRetencionRepository registra los cortes recibidos y devuelve los conteos
// configurados, sin tocar una base real.
type fakeRetencionRepository struct {
	cortesCotizaciones []time.Time
	cortesConsultas    []time.Time
}

func (f *fakeRetencionRepository) PurgarCotizacionesCerradas(_ context.Context, corte time.Time) (int64, int64, error) {
	f.cortesCotizaciones = append(f.cortesCotizaciones, corte)
	return 12, 3, nil
}

func (f *fakeRetencionRepository) PurgarConsultasCerradas(_ context.Context, corte time.Time) (int64, int64, error) {
	f.cortesConsultas = append(f.cortesConsultas, corte)
	return 4, 1, nil
}

var _ repositories.RetencionRepository = (*fakeRetencionRepository)(nil)

func TestRetencionEjecutarPurgaCerradasViejas(t *testing.T) {
	repositorio := &fakeRetencionRepository{}
	servicio := NuevoRetencionService(repositorio)

	resultado, err := servicio.Ejecutar(context.Background(), 180)
	if err != nil {
		t.Fatalf("Ejecutar devolvió error inesperado: %v", err)
	}

	if resultado.MensajesCotizaciones != 12 || resultado.Cotizaciones != 3 {
		t.Fatalf("conteos de cotizaciones incorrectos: mensajes=%d cotizaciones=%d", resultado.MensajesCotizaciones, resultado.Cotizaciones)
	}
	if resultado.MensajesConsultas != 4 || resultado.Consultas != 1 {
		t.Fatalf("conteos de consultas incorrectos: mensajes=%d consultas=%d", resultado.MensajesConsultas, resultado.Consultas)
	}

	if len(repositorio.cortesCotizaciones) != 1 || len(repositorio.cortesConsultas) != 1 {
		t.Fatalf("se esperaba un corte por tabla, got cotizaciones=%d consultas=%d", len(repositorio.cortesCotizaciones), len(repositorio.cortesConsultas))
	}
	esperado := time.Now().AddDate(0, 0, -180)
	for _, corte := range append(repositorio.cortesCotizaciones, repositorio.cortesConsultas...) {
		// Tolerancia de unos segundos para no fallar por diferencias de clock.
		if corte.Sub(esperado) > 5*time.Second || esperado.Sub(corte) > 5*time.Second {
			t.Fatalf("corte recibido %v no coincide con lo esperado %v", corte, esperado)
		}
		despues := corte.Add(179 * 24 * time.Hour)
		if despues.After(time.Now()) {
			t.Fatalf("el corte %v no respeta el plazo de 180 días", corte)
		}
	}
}

func TestRetencionEjecutarRechazaPlazoInvalido(t *testing.T) {
	servicio := NuevoRetencionService(&fakeRetencionRepository{})

	if _, err := servicio.Ejecutar(context.Background(), 0); err == nil {
		t.Fatal("se esperaba error con plazo 0")
	}
}