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

	// La respuesta HTTP (lo que devuelve el servicio) no debe exponer el texto
	// cifrado: los mensajes viajan descifrados.
	if cotizacion.Mensajes[0].Contenido != "Hola, quiero cotizar este vehículo." {
		t.Errorf("el mensaje del cliente debería devolverse descifrado: %q", cotizacion.Mensajes[0].Contenido)
	}
	if cotizacion.Mensajes[1].Contenido != "Respuesta para: Hola, quiero cotizar este vehículo." {
		t.Errorf("la respuesta de la IA debería devolverse descifrada: %q", cotizacion.Mensajes[1].Contenido)
	}

	// El contenido debe guardarse cifrado en la base.
	guardada := repoCotizaciones.porID[cotizacion.ID]
	if guardada.Mensajes[0].Contenido == "Hola, quiero cotizar este vehículo." {
		t.Error("el mensaje del cliente se guardó en claro, debería estar cifrado")
	}
	if guardada.Mensajes[1].Contenido == "Respuesta para: Hola, quiero cotizar este vehículo." {
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
	// Al cerrar, la respuesta no debe exponer el texto cifrado.
	if len(cerrada.Mensajes) != 2 {
		t.Fatalf("se esperaban 2 mensajes, se obtuvieron %d", len(cerrada.Mensajes))
	}
	if cerrada.Mensajes[0].Contenido != "Hola" || cerrada.Mensajes[1].Contenido != "Respuesta para: Hola" {
		t.Error("al cerrar, los mensajes deberían viajar descifrados")
	}

	_, err = servicio.Cerrar(context.Background(), 5, cotizacion.ID)
	if !errors.Is(err, ErrCotizacionYaCerrada) {
		t.Errorf("se esperaba ErrCotizacionYaCerrada, se obtuvo %v", err)
	}
}

func TestCotizacionTomaYCierrePersonalConMensajesDescifrados(t *testing.T) {
	repoCotizaciones := nuevoFakeCotizacionRepository()
	repoVehiculos := nuevoFakeVehiculoRepository()
	repoVehiculos.porID[7] = vehiculoDisponiblePrueba(7)
	servicio := nuevoServicioCotizacionesPrueba(t, repoCotizaciones, repoVehiculos, &fakeGeneradorCotizacion{})

	cotizacion, err := servicio.Crear(context.Background(), 5, 7, "Hola")
	if err != nil {
		t.Fatalf("no se pudo crear la cotización: %v", err)
	}

	tomada, err := servicio.Tomar(context.Background(), 9, cotizacion.ID)
	if err != nil {
		t.Fatalf("no se pudo tomar la cotización: %v", err)
	}
	if tomada.Mensajes[0].Contenido != "Hola" {
		t.Errorf("al tomar, el mensaje debería viajar descifrado: %q", tomada.Mensajes[0].Contenido)
	}

	cerrada, err := servicio.CerrarPersonal(context.Background(), cotizacion.ID)
	if err != nil {
		t.Fatalf("no se pudo cerrar la cotización: %v", err)
	}
	if cerrada.Estado != models.EstadoCotizacionCerrada {
		t.Errorf("estado esperado cerrada, obtenido %s", cerrada.Estado)
	}
	if cerrada.Mensajes[0].Contenido != "Hola" || cerrada.Mensajes[1].Contenido != "Respuesta para: Hola" {
		t.Error("al cerrar desde el personal, los mensajes deberían viajar descifrados")
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

// sembrarCotizacionPrueba inserta una cotización directamente en el fake con
// los mensajes ya en claro: alcanza para probar conteos y marcas de lectura.
func sembrarCotizacionPrueba(t *testing.T, repo *fakeCotizacionRepository, clienteID uint, vendedorID *uint, estado string, mensajes []models.MensajeCotizacion) *models.Cotizacion {
	t.Helper()
	cotizacion, err := repo.Crear(context.Background(), &models.Cotizacion{
		VehiculoID: 7,
		ClienteID:  clienteID,
		VendedorID: vendedorID,
		Estado:     estado,
		Mensajes:   mensajes,
	})
	if err != nil {
		t.Fatalf("no se pudo sembrar la cotización: %v", err)
	}
	return cotizacion
}

func TestCotizacionObtenerMensajesDesde(t *testing.T) {
	repoCotizaciones := nuevoFakeCotizacionRepository()
	repoVehiculos := nuevoFakeVehiculoRepository()
	repoVehiculos.porID[7] = vehiculoDisponiblePrueba(7)
	servicio := nuevoServicioCotizacionesPrueba(t, repoCotizaciones, repoVehiculos, &fakeGeneradorCotizacion{})
	ctx := context.Background()

	cotizacion, err := servicio.Crear(ctx, 5, 7, "Hola")
	if err != nil {
		t.Fatalf("no se pudo crear la cotización: %v", err)
	}
	if _, err := servicio.EnviarMensaje(ctx, 5, cotizacion.ID, "¿Tienen financiación?"); err != nil {
		t.Fatalf("no se pudo enviar el mensaje: %v", err)
	}

	// El delta a partir del segundo mensaje devuelve solo los dos posteriores,
	// descifrados.
	desdeID := cotizacion.Mensajes[1].ID
	estado, err := servicio.ObtenerMensajesDesde(ctx, 5, models.RolCliente, cotizacion.ID, desdeID)
	if err != nil {
		t.Fatalf("no se pudo obtener el delta: %v", err)
	}
	if estado.Total != 4 {
		t.Errorf("total esperado 4, obtenido %d", estado.Total)
	}
	if len(estado.Mensajes) != 2 {
		t.Fatalf("se esperaban 2 mensajes nuevos, se obtuvieron %d", len(estado.Mensajes))
	}
	if estado.Mensajes[0].Remitente != models.RemitenteCliente || estado.Mensajes[0].Contenido != "¿Tienen financiación?" {
		t.Errorf("el primer mensaje nuevo no es el esperado: %+v", estado.Mensajes[0])
	}
	if estado.Mensajes[1].Remitente != models.RemitenteIA || estado.Mensajes[1].Contenido != "Respuesta para: ¿Tienen financiación?" {
		t.Errorf("el segundo mensaje nuevo no es el esperado: %+v", estado.Mensajes[1])
	}
	if estado.Estado != models.EstadoCotizacionAbierta {
		t.Errorf("estado esperado abierta, obtenido %s", estado.Estado)
	}

	// Con desdeID en 0 devuelve todo el historial descifrado.
	desdeCero, err := servicio.ObtenerMensajesDesde(ctx, 5, models.RolCliente, cotizacion.ID, 0)
	if err != nil {
		t.Fatalf("no se pudo obtener el historial completo: %v", err)
	}
	if len(desdeCero.Mensajes) != 4 || desdeCero.Mensajes[0].Contenido != "Hola" {
		t.Errorf("con desdeId 0 debería devolver el historial completo descifrado")
	}

	// Un cliente ajeno no puede pedir el delta de la cotización de otro.
	if _, err := servicio.ObtenerMensajesDesde(ctx, 99, models.RolCliente, cotizacion.ID, 0); !errors.Is(err, ErrCotizacionNoPertenece) {
		t.Errorf("se esperaba ErrCotizacionNoPertenece, se obtuvo %v", err)
	}

	// El personal puede pedirlo aunque no esté asignado a la conversación.
	if _, err := servicio.ObtenerMensajesDesde(ctx, 9, models.RolVendedor, cotizacion.ID, 0); err != nil {
		t.Errorf("el personal debería poder pedir el delta: %v", err)
	}
}

func TestCotizacionObtenerMensajesDesdeNoEncontrada(t *testing.T) {
	repo := nuevoFakeCotizacionRepository()
	servicio := nuevoServicioCotizacionesPrueba(t, repo, nuevoFakeVehiculoRepository(), &fakeGeneradorCotizacion{})

	if _, err := servicio.ObtenerMensajesDesde(context.Background(), 5, models.RolCliente, 999, 0); !errors.Is(err, ErrCotizacionNoEncontrada) {
		t.Errorf("se esperaba ErrCotizacionNoEncontrada, se obtuvo %v", err)
	}
}

func TestCotizacionContarNoLeidosPorLado(t *testing.T) {
	repo := nuevoFakeCotizacionRepository()
	servicio := nuevoServicioCotizacionesPrueba(t, repo, nuevoFakeVehiculoRepository(), &fakeGeneradorCotizacion{})
	ctx := context.Background()

	// Abierta y sin asignar: mensaje de cliente, respuesta de IA y mensaje del
	// vendedor sin leer.
	sembrarCotizacionPrueba(t, repo, 5, nil, models.EstadoCotizacionAbierta, []models.MensajeCotizacion{
		{Remitente: models.RemitenteCliente, Contenido: "hola"},
		{Remitente: models.RemitenteIA, Contenido: "respuesta"},
		{Remitente: models.RemitenteVendedor, Contenido: "te paso el valor"},
	})

	noLeidosCliente, err := servicio.ContarNoLeidos(ctx, 5, models.RolCliente)
	if err != nil {
		t.Fatalf("error contando para el cliente: %v", err)
	}
	if noLeidosCliente != 2 {
		t.Errorf("el cliente debería tener 2 no leídos (IA y vendedor), obtuvo %d", noLeidosCliente)
	}

	// Para el personal cuenta solo el mensaje del cliente; la IA nunca suma.
	noLeidosVendedor, err := servicio.ContarNoLeidos(ctx, 9, models.RolVendedor)
	if err != nil {
		t.Fatalf("error contando para el personal: %v", err)
	}
	if noLeidosVendedor != 1 {
		t.Errorf("el vendedor debería tener 1 no leído (cliente), obtuvo %d", noLeidosVendedor)
	}

	// Asignada a otro vendedor: sus mensajes de cliente no cuentan para el
	// vendedor 9 pero sí para el asignado.
	asignadaOtro := uint(3)
	sembrarCotizacionPrueba(t, repo, 6, &asignadaOtro, models.EstadoCotizacionAbierta, []models.MensajeCotizacion{
		{Remitente: models.RemitenteCliente, Contenido: "consulta"},
	})
	if n, _ := servicio.ContarNoLeidos(ctx, 9, models.RolVendedor); n != 1 {
		t.Errorf("el vendedor 9 no debería ver la cotización ajena, obtuvo %d", n)
	}
	// El vendedor 3 ve su asignada más la sin asignar de c1.
	if n, _ := servicio.ContarNoLeidos(ctx, 3, models.RolVendedor); n != 2 {
		t.Errorf("el vendedor asignado debería tener 2 no leídos (propia + sin asignar), obtuvo %d", n)
	}

	// Cerrada con mensaje de cliente sin leer: no avisa al personal.
	sembrarCotizacionPrueba(t, repo, 5, nil, models.EstadoCotizacionCerrada, []models.MensajeCotizacion{
		{Remitente: models.RemitenteCliente, Contenido: "última consulta"},
	})
	if n, _ := servicio.ContarNoLeidos(ctx, 9, models.RolVendedor); n != 1 {
		t.Errorf("las cerradas no deberían contar para el personal, obtuvo %d", n)
	}
}

func TestCotizacionMarcarLeidas(t *testing.T) {
	repo := nuevoFakeCotizacionRepository()
	servicio := nuevoServicioCotizacionesPrueba(t, repo, nuevoFakeVehiculoRepository(), &fakeGeneradorCotizacion{})
	ctx := context.Background()

	c1 := sembrarCotizacionPrueba(t, repo, 5, nil, models.EstadoCotizacionAbierta, []models.MensajeCotizacion{
		{Remitente: models.RemitenteCliente, Contenido: "hola"},
		{Remitente: models.RemitenteIA, Contenido: "respuesta"},
		{Remitente: models.RemitenteVendedor, Contenido: "te paso el valor"},
	})
	asignada := uint(9)
	c2 := sembrarCotizacionPrueba(t, repo, 6, &asignada, models.EstadoCotizacionAbierta, []models.MensajeCotizacion{
		{Remitente: models.RemitenteCliente, Contenido: "consulta"},
	})

	// Un cliente ajeno no puede marcar la cotización de otro.
	if err := servicio.MarcarLeidas(ctx, 6, c1.ID, LadoCliente); !errors.Is(err, ErrCotizacionNoPertenece) {
		t.Errorf("se esperaba ErrCotizacionNoPertenece, obtuvo %v", err)
	}

	// El dueño marca su lado: se limpian los mensajes de ia/vendedor.
	if err := servicio.MarcarLeidas(ctx, 5, c1.ID, LadoCliente); err != nil {
		t.Fatalf("no se pudo marcar el lado cliente: %v", err)
	}
	if n, _ := servicio.ContarNoLeidos(ctx, 5, models.RolCliente); n != 0 {
		t.Errorf("después de abrir el hilo el cliente no debería tener no leídos, obtuvo %d", n)
	}
	// El mensaje del cliente sigue pendiente para el personal.
	if n, _ := servicio.ContarNoLeidos(ctx, 9, models.RolVendedor); n != 2 {
		t.Errorf("marcar el lado cliente no debería tocar el lado personal, obtuvo %d", n)
	}

	// Otro vendedor intenta marcar una cotización que no es suya: no hace nada.
	if err := servicio.MarcarLeidas(ctx, 3, c2.ID, LadoPersonal); err != nil {
		t.Fatalf("marcar con un vendedor ajeno no debería fallar: %v", err)
	}
	if n, _ := servicio.ContarNoLeidos(ctx, 9, models.RolVendedor); n != 2 {
		t.Errorf("un vendedor ajeno no debería limpiar pendientes ajenos, obtuvo %d", n)
	}

	// El vendedor asignado abre el hilo y limpia su pendiente.
	if err := servicio.MarcarLeidas(ctx, 9, c2.ID, LadoPersonal); err != nil {
		t.Fatalf("no se pudo marcar el lado personal: %v", err)
	}
	if n, _ := servicio.ContarNoLeidos(ctx, 9, models.RolVendedor); n != 1 {
		t.Errorf("quedaría 1 pendiente (mensaje de cliente en c1), obtuvo %d", n)
	}
}
