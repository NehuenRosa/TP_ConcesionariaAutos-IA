package services

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"

	"concesionaria/backend/internal/cifrado"
	"concesionaria/backend/internal/models"
)

// nuevoCifradorPrueba crea un cifrador real con clave fija para las pruebas del
// flujo de cotizaciones.
func nuevoCifradorPrueba(t *testing.T) cifrado.Cifrador {
	t.Helper()
	cifrador, err := cifrado.NuevoCifrador("clave-de-prueba-para-cotizaciones")
	if err != nil {
		t.Fatalf("no se pudo crear el cifrador: %v", err)
	}
	return cifrador
}

// vehiculoDisponiblePrueba devuelve un vehículo disponible de prueba.
func vehiculoDisponiblePrueba(id uint) *models.Vehiculo {
	return &models.Vehiculo{
		ID:     id,
		Marca:  "Toyota",
		Modelo: "Corolla",
		Anio:   2022,
		Precio: 25000,
		Estado: models.EstadoDisponible,
	}
}

func nuevoServicioCotizacionesPrueba(t *testing.T, repoCotizaciones *fakeCotizacionRepository, repoVehiculos *fakeVehiculoRepository, generador *fakeGeneradorCotizacion) CotizacionService {
	t.Helper()
	return NuevoCotizacionService(repoCotizaciones, repoVehiculos, nuevoCifradorPrueba(t), generador)
}

func TestCotizacionCrear(t *testing.T) {
	repoCotizaciones := nuevoFakeCotizacionRepository()
	repoVehiculos := nuevoFakeVehiculoRepository()
	vehiculo := vehiculoDisponiblePrueba(7)
	repoVehiculos.porID[7] = vehiculo
	servicio := nuevoServicioCotizacionesPrueba(t, repoCotizaciones, repoVehiculos, &fakeGeneradorCotizacion{})

	cotizacion, err := servicio.Crear(context.Background(), 5, 7, "")
	if err != nil {
		t.Fatalf("se esperaba crear la cotización, se obtuvo %v", err)
	}

	if cotizacion.VehiculoID != 7 {
		t.Errorf("vehiculoID esperado 7, obtenido %d", cotizacion.VehiculoID)
	}
	if cotizacion.ClienteID != 5 {
		t.Errorf("clienteID esperado 5, obtenido %d", cotizacion.ClienteID)
	}
	if cotizacion.Estado != models.EstadoCotizacionAbierta {
		t.Errorf("estado esperado abierta, obtenido %s", cotizacion.Estado)
	}
	if len(cotizacion.Mensajes) != 2 {
		t.Fatalf("se esperaban 2 mensajes iniciales, se obtuvieron %d", len(cotizacion.Mensajes))
	}
	if cotizacion.Mensajes[0].Remitente != models.RemitenteCliente || cotizacion.Mensajes[1].Remitente != models.RemitenteIA {
		t.Error("el orden de remitentes debería ser cliente y luego IA")
	}

	// El contenido debe guardarse cifrado en la base.
	if cotizacion.Mensajes[0].Contenido == "Hola, quiero cotizar este vehículo." {
		t.Error("el mensaje del cliente se guardó en claro, debería estar cifrado")
	}
	if cotizacion.Mensajes[1].Contenido == "Respuesta para: Hola, quiero cotizar este vehículo." {
		t.Error("la respuesta de la IA se guardó en claro, debería estar cifrada")
	}

	// Y descifrarse al leer.
	obtenida, err := servicio.ObtenerPorCliente(context.Background(), 5, cotizacion.ID)
	if err != nil {
		t.Fatalf("no se pudo obtener la cotización: %v", err)
	}
	if obtenida.Mensajes[0].Contenido != "Hola, quiero cotizar este vehículo." {
		t.Errorf("el mensaje descifrado no coincide: %q", obtenida.Mensajes[0].Contenido)
	}
	if obtenida.Mensajes[1].Contenido != "Respuesta para: Hola, quiero cotizar este vehículo." {
		t.Errorf("la respuesta descifrada no coincide: %q", obtenida.Mensajes[1].Contenido)
	}
}

func TestCotizacionCrearVehiculoNoDisponible(t *testing.T) {
	repoCotizaciones := nuevoFakeCotizacionRepository()
	repoVehiculos := nuevoFakeVehiculoRepository()
	vehiculo := vehiculoDisponiblePrueba(7)
	vehiculo.Estado = models.EstadoReservado
	repoVehiculos.porID[7] = vehiculo
	servicio := nuevoServicioCotizacionesPrueba(t, repoCotizaciones, repoVehiculos, &fakeGeneradorCotizacion{})

	_, err := servicio.Crear(context.Background(), 5, 7, "")
	if !errors.Is(err, ErrVehiculoNoDisponible) {
		t.Errorf("se esperaba ErrVehiculoNoDisponible, se obtuvo %v", err)
	}
}

func TestCotizacionEnviarMensaje(t *testing.T) {
	repoCotizaciones := nuevoFakeCotizacionRepository()
	repoVehiculos := nuevoFakeVehiculoRepository()
	repoVehiculos.porID[7] = vehiculoDisponiblePrueba(7)
	servicio := nuevoServicioCotizacionesPrueba(t, repoCotizaciones, repoVehiculos, &fakeGeneradorCotizacion{})

	cotizacion, err := servicio.Crear(context.Background(), 5, 7, "Hola")
	if err != nil {
		t.Fatalf("no se pudo crear la cotización: %v", err)
	}

	respuesta, err := servicio.EnviarMensaje(context.Background(), 5, cotizacion.ID, "¿Tienen financiación?")
	if err != nil {
		t.Fatalf("no se pudo enviar el mensaje: %v", err)
	}
	if respuesta != "Respuesta para: ¿Tienen financiación?" {
		t.Errorf("respuesta inesperada: %q", respuesta)
	}

	guardada, err := servicio.ObtenerPorCliente(context.Background(), 5, cotizacion.ID)
	if err != nil {
		t.Fatalf("no se pudo obtener la cotización: %v", err)
	}
	if len(guardada.Mensajes) != 4 {
		t.Fatalf("se esperaban 4 mensajes, se obtuvieron %d", len(guardada.Mensajes))
	}
	ultimo := guardada.Mensajes[3]
	if ultimo.Remitente != models.RemitenteIA || ultimo.Contenido != "Respuesta para: ¿Tienen financiación?" {
		t.Errorf("el último mensaje debería ser la respuesta de la IA descifrada")
	}
}

func TestCotizacionEnviarMensajeNoPertenece(t *testing.T) {
	repoCotizaciones := nuevoFakeCotizacionRepository()
	repoVehiculos := nuevoFakeVehiculoRepository()
	repoVehiculos.porID[7] = vehiculoDisponiblePrueba(7)
	servicio := nuevoServicioCotizacionesPrueba(t, repoCotizaciones, repoVehiculos, &fakeGeneradorCotizacion{})

	cotizacion, err := servicio.Crear(context.Background(), 5, 7, "Hola")
	if err != nil {
		t.Fatalf("no se pudo crear la cotización: %v", err)
	}

	_, err = servicio.EnviarMensaje(context.Background(), 99, cotizacion.ID, "Hola")
	if !errors.Is(err, ErrCotizacionNoPertenece) {
		t.Errorf("se esperaba ErrCotizacionNoPertenece, se obtuvo %v", err)
	}
}

func TestCotizacionEnviarMensajeCerrada(t *testing.T) {
	repoCotizaciones := nuevoFakeCotizacionRepository()
	repoVehiculos := nuevoFakeVehiculoRepository()
	repoVehiculos.porID[7] = vehiculoDisponiblePrueba(7)
	servicio := nuevoServicioCotizacionesPrueba(t, repoCotizaciones, repoVehiculos, &fakeGeneradorCotizacion{})

	cotizacion, err := servicio.Crear(context.Background(), 5, 7, "Hola")
	if err != nil {
		t.Fatalf("no se pudo crear la cotización: %v", err)
	}
	if _, err := servicio.Cerrar(context.Background(), 5, cotizacion.ID); err != nil {
		t.Fatalf("no se pudo cerrar la cotización: %v", err)
	}

	_, err = servicio.EnviarMensaje(context.Background(), 5, cotizacion.ID, "Hola")
	if !errors.Is(err, ErrCotizacionCerradaMensajes) {
		t.Errorf("se esperaba ErrCotizacionCerradaMensajes, se obtuvo %v", err)
	}
}

func TestCotizacionCerrar(t *testing.T) {
	repoCotizaciones := nuevoFakeCotizacionRepository()
	repoVehiculos := nuevoFakeVehiculoRepository()
	repoVehiculos.porID[7] = vehiculoDisponiblePrueba(7)
	servicio := nuevoServicioCotizacionesPrueba(t, repoCotizaciones, repoVehiculos, &fakeGeneradorCotizacion{})

	cotizacion, err := servicio.Crear(context.Background(), 5, 7, "Hola")
	if err != nil {
		t.Fatalf("no se pudo crear la cotización: %v", err)
	}

	cerrada, err := servicio.Cerrar(context.Background(), 5, cotizacion.ID)
	if err != nil {
		t.Fatalf("no se pudo cerrar la cotización: %v", err)
	}
	if cerrada.Estado != models.EstadoCotizacionCerrada {
		t.Errorf("estado esperado cerrada, obtenido %s", cerrada.Estado)
	}

	_, err = servicio.Cerrar(context.Background(), 5, cotizacion.ID)
	if !errors.Is(err, ErrCotizacionYaCerrada) {
		t.Errorf("se esperaba ErrCotizacionYaCerrada, se obtuvo %v", err)
	}
}

func TestTextoCondicionesComerciales(t *testing.T) {
	texto := textoCondicionesComerciales(1000000)

	casosEsperados := []string{
		"Contado / transferencia bancaria",
		"Tarjeta de crédito",
		"Plan 30/70",
		"24 cuotas",
		"36 cuotas",
		"Plan de ahorro (prenda)",
		"$",
	}
	for _, esperado := range casosEsperados {
		if !strings.Contains(texto, esperado) {
			t.Errorf("se esperaba que el texto contenga %q", esperado)
		}
	}

	if strings.Contains(texto, "NaN") || strings.Contains(texto, "Inf") {
		t.Errorf("el texto no debe contener valores no numéricos: %q", texto)
	}

	if got, quiero := cuotaEstimada(1000000, planFinanciacion{adelantoMin: 0.30, cuotas: 6, tasaMensual: 0}), 116666.67; math.Abs(got-quiero) > 0.01 {
		t.Errorf("cuota del Plan 30/70 esperada %.2f, obtenida %.2f", quiero, got)
	}
}
