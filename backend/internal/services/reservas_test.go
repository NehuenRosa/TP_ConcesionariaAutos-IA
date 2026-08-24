package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"concesionaria/backend/internal/models"
)

// nuevoServicioReservas arma el servicio con fakes y un vehículo disponible de
// precio conocido ($10.000.000 → seña $500.000).
func nuevoServicioReservas(t *testing.T) (ReservaService, *fakeReservaRepository, *fakeVehiculoRepository) {
	t.Helper()

	vehiculos := nuevoFakeVehiculoRepository()
	vehiculos.porID[1] = &models.Vehiculo{
		ID:      1,
		Marca:   "Toyota",
		Modelo:  "Corolla",
		Anio:    2022,
		Precio:  10000000,
		Estado:  models.EstadoDisponible,
	}
	reservas := nuevoFakeReservaRepository(vehiculos)
	servicio := NuevoReservaService(reservas, vehiculos, "cbu-de-prueba", "alias-de-prueba")
	return servicio, reservas, vehiculos
}

// guardarReserva inserta directamente una reserva en el fake con el vencimiento
// indicado.
func guardarReserva(repo *fakeReservaRepository, id uint, clienteID uint, vehiculoID uint, vencimiento time.Time) *models.Reserva {
	reserva := &models.Reserva{
		ID:                     id,
		VehiculoID:             vehiculoID,
		ClienteID:              clienteID,
		Estado:                 models.EstadoReservaActiva,
		VencimientoComprobante: vencimiento,
	}
	repo.porID[id] = reserva
	return reserva
}

func TestCrearFijaVencimientoYCalculaMontoSenia(t *testing.T) {
	servicio, _, _ := nuevoServicioReservas(t)

	antes := time.Now()
	reserva, err := servicio.Crear(context.Background(), 7, 1)
	if err != nil {
		t.Fatalf("Crear devolvió error: %v", err)
	}
	despues := time.Now()

	if !reserva.EsActiva() {
		t.Fatalf("la reserva debió quedar activa")
	}
	vencimiento := reserva.VencimientoComprobante
	if vencimiento.Before(antes.Add(PlazoComprobante-time.Second)) || vencimiento.After(despues.Add(PlazoComprobante)) {
		t.Errorf("vencimiento fuera del plazo esperado: %v", vencimiento)
	}
	if monto := CalcularMontoSenia(10000000); monto != 500000 {
		t.Errorf("monto de seña incorrecto: %v", monto)
	}
}

func TestObtenerDatosTransferenciaDevuelveCbuAliasYMonto(t *testing.T) {
	servicio, _, _ := nuevoServicioReservas(t)

	datos, err := servicio.ObtenerDatosTransferencia(context.Background(), 1)
	if err != nil {
		t.Fatalf("ObtenerDatosTransferencia devolvió error: %v", err)
	}
	if datos.CBU != "cbu-de-prueba" || datos.Alias != "alias-de-prueba" {
		t.Errorf("datos bancarios incorrectos: %+v", datos)
	}
	if datos.Monto != 500000 {
		t.Errorf("monto esperado 500000, obtenido %v", datos.Monto)
	}
}

func TestObtenerDatosTransferenciaVehiculoNoDisponible(t *testing.T) {
	servicio, _, vehiculos := nuevoServicioReservas(t)
	vehiculos.porID[1].Estado = models.EstadoReservado

	if _, err := servicio.ObtenerDatosTransferencia(context.Background(), 1); !errors.Is(err, ErrVehiculoNoDisponible) {
		t.Errorf("se esperaba ErrVehiculoNoDisponible, se obtuvo %v", err)
	}
}

func TestSubirComprobanteValido(t *testing.T) {
	servicio, repo, _ := nuevoServicioReservas(t)
	guardarReserva(repo, 3, 7, 1, time.Now().Add(time.Hour))

	png := []byte{0x89, 'P', 'N', 'G'}
	reserva, err := servicio.SubirComprobante(context.Background(), 3, 7, "comprobante.PNG", png)
	if err != nil {
		t.Fatalf("SubirComprobante devolvió error: %v", err)
	}
	if reserva.ComprobanteEnviadoAt == nil {
		t.Fatal("el envío debió registrarse en la reserva")
	}
	comprobante, err := repo.ObtenerComprobantePorReservaID(context.Background(), 3)
	if err != nil {
		t.Fatalf("el comprobante debió guardarse: %v", err)
	}
	if comprobante.MIME != "image/png" || len(comprobante.Datos) == 0 {
		t.Errorf("comprobante mal guardado: mime=%q datos=%d", comprobante.MIME, len(comprobante.Datos))
	}

	// El dueño lo puede ver; un tercero no; el personal sí.
	if _, err := servicio.ObtenerComprobante(context.Background(), 3, 7, false); err != nil {
		t.Errorf("el dueño debió poder ver su comprobante: %v", err)
	}
	if _, err := servicio.ObtenerComprobante(context.Background(), 3, 99, false); !errors.Is(err, ErrReservaProhibida) {
		t.Errorf("se esperaba ErrReservaProhibida para un tercero, se obtuvo %v", err)
	}
	if _, err := servicio.ObtenerComprobante(context.Background(), 3, 99, true); err != nil {
		t.Errorf("el personal debió poder ver el comprobante: %v", err)
	}
}

func TestSubirComprobanteInvalido(t *testing.T) {
	servicio, repo, _ := nuevoServicioReservas(t)
	guardarReserva(repo, 3, 7, 1, time.Now().Add(time.Hour))

	if _, err := servicio.SubirComprobante(context.Background(), 3, 7, "comprobante.pdf", []byte("pdf")); !errors.Is(err, ErrComprobanteInvalido) {
		t.Errorf("se esperaba ErrComprobanteInvalido por extensión, se obtuvo %v", err)
	}
	grande := make([]byte, MaximoPesoComprobanteBytes+1)
	if _, err := servicio.SubirComprobante(context.Background(), 3, 7, "foto.jpg", grande); !errors.Is(err, ErrComprobanteInvalido) {
		t.Errorf("se esperaba ErrComprobanteInvalido por peso, se obtuvo %v", err)
	}
	if _, err := servicio.SubirComprobante(context.Background(), 3, 99, "foto.jpg", []byte("jpg")); !errors.Is(err, ErrReservaNoEncontrada) {
		t.Errorf("una reserva ajena debe tratarse como inexistente, se obtuvo %v", err)
	}
}

func TestSubirComprobanteFueraDePlazoAnulaLaReserva(t *testing.T) {
	servicio, repo, vehiculos := nuevoServicioReservas(t)
	guardarReserva(repo, 3, 7, 1, time.Now().Add(-time.Minute))

	if _, err := servicio.SubirComprobante(context.Background(), 3, 7, "foto.jpg", []byte("jpg")); !errors.Is(err, ErrComprobanteFueraDePlazo) {
		t.Fatalf("se esperaba ErrComprobanteFueraDePlazo, se obtuvo %v", err)
	}
	guardada := repo.porID[3]
	if guardada.Estado != models.EstadoReservaCancelada {
		t.Errorf("la reserva debió anularse, quedó %q", guardada.Estado)
	}
	if vehiculos.porID[1].Estado != models.EstadoDisponible {
		t.Errorf("el vehículo debió liberarse, quedó %q", vehiculos.porID[1].Estado)
	}
}

func TestConfirmarVentaSobreVencidaAplicaExpiracionPrimero(t *testing.T) {
	servicio, repo, vehiculos := nuevoServicioReservas(t)
	guardarReserva(repo, 3, 7, 1, time.Now().Add(-time.Hour))

	if _, err := servicio.ConfirmarVenta(context.Background(), 3); !errors.Is(err, ErrReservaEstadoInvalido) {
		t.Fatalf("se esperaba ErrReservaEstadoInvalido, se obtuvo %v", err)
	}
	if repo.porID[3].Estado != models.EstadoReservaCancelada {
		t.Errorf("la reserva debió expirar a cancelada, quedó %q", repo.porID[3].Estado)
	}
	if vehiculos.porID[1].Estado != models.EstadoDisponible {
		t.Errorf("el vehículo debió liberarse, quedó %q", vehiculos.porID[1].Estado)
	}
}

func TestCancelarSobreVencidaAplicaExpiracionPrimero(t *testing.T) {
	servicio, repo, vehiculos := nuevoServicioReservas(t)
	guardarReserva(repo, 3, 7, 1, time.Now().Add(-time.Hour))

	if _, err := servicio.Cancelar(context.Background(), 3, 7); !errors.Is(err, ErrReservaEstadoInvalido) {
		t.Fatalf("se esperaba ErrReservaEstadoInvalido, se obtuvo %v", err)
	}
	if repo.porID[3].Estado != models.EstadoReservaCancelada {
		t.Errorf("la reserva debió expirar a cancelada, quedó %q", repo.porID[3].Estado)
	}
	if vehiculos.porID[1].Estado != models.EstadoDisponible {
		t.Errorf("el vehículo debió liberarse, quedó %q", vehiculos.porID[1].Estado)
	}
}

func TestExpirarVencidasCancelaSoloLasPendientesVencidas(t *testing.T) {
	servicio, repo, vehiculos := nuevoServicioReservas(t)

	ahora := time.Now()
	vehiculos.porID[2] = &models.Vehiculo{ID: 2, Estado: models.EstadoReservado}
	vehiculos.porID[3] = &models.Vehiculo{ID: 3, Estado: models.EstadoReservado}
	guardarReserva(repo, 10, 7, 1, ahora.Add(-time.Minute))            // vencida sin comprobante: expira
	guardarReserva(repo, 11, 8, 2, ahora.Add(time.Hour))               // activa dentro del plazo: no toca
	enviada := guardarReserva(repo, 12, 9, 3, ahora.Add(-time.Minute)) // con comprobante: no expira
	horaEnvio := ahora.Add(-2 * time.Hour)
	enviada.ComprobanteEnviadoAt = &horaEnvio

	cantidad, err := repo.ExpirarVencidas(context.Background())
	if err != nil {
		t.Fatalf("ExpirarVencidas devolvió error: %v", err)
	}
	if cantidad != 1 {
		t.Errorf("solo una reserva debía expirar, expiraron %d", cantidad)
	}
	if repo.porID[10].Estado != models.EstadoReservaCancelada ||
		repo.porID[11].Estado != models.EstadoReservaActiva ||
		repo.porID[12].Estado != models.EstadoReservaActiva {
		t.Errorf("estados finales incorrectos: %q %q %q",
			repo.porID[10].Estado, repo.porID[11].Estado, repo.porID[12].Estado)
	}
	if vehiculos.porID[1].Estado != models.EstadoDisponible ||
		vehiculos.porID[2].Estado != models.EstadoReservado ||
		vehiculos.porID[3].Estado != models.EstadoReservado {
		t.Errorf("estados de vehículos incorrectos: %q %q %q",
			vehiculos.porID[1].Estado, vehiculos.porID[2].Estado, vehiculos.porID[3].Estado)
	}

	if err := servicio.ExpirarVencidas(context.Background()); err != nil {
		t.Errorf("ExpirarVencidas del servicio devolvió error: %v", err)
	}
}

func TestValidarComprobanteFormatosAdmitidos(t *testing.T) {
	formatos := map[string]string{
		"a.jpg": "image/jpeg", "b.JPEG": "image/jpeg", "c.png": "image/png", "d.webp": "image/webp",
	}
	for nombre, mimeEsperado := range formatos {
		mime, err := validarComprobante(nombre, []byte("datos"))
		if err != nil {
			t.Errorf("%s debió ser válido: %v", nombre, err)
			continue
		}
		if mime != mimeEsperado {
			t.Errorf("%s: MIME esperado %q, obtenido %q", nombre, mimeEsperado, mime)
		}
	}
	if _, err := validarComprobante("sin-extension", []byte("datos")); err == nil {
		t.Error("un archivo sin extensión admitida no debió pasar")
	}
	if _, err := validarComprobante("vacio.png", nil); err == nil {
		t.Error("un archivo vacío no debió pasar")
	}
}

func TestCancelarComoVendedorExigeMotivoYLoGuarda(t *testing.T) {
	servicio, repo, vehiculos := nuevoServicioReservas(t)
	guardarReserva(repo, 5, 7, 1, time.Now().Add(time.Hour))

	// Sin motivo (vacío o solo espacios) no cancela.
	if _, err := servicio.CancelarComoVendedor(context.Background(), 5, "   "); !errors.Is(err, ErrMotivoRequerido) {
		t.Errorf("se esperaba ErrMotivoRequerido, se obtuvo %v", err)
	}
	if !repo.porID[5].EsActiva() {
		t.Error("sin motivo la reserva debía seguir activa")
	}

	motivo := "El comprobante es ilegible y el monto no coincide con la seña"
	reserva, err := servicio.CancelarComoVendedor(context.Background(), 5, motivo)
	if err != nil {
		t.Fatalf("CancelarComoVendedor devolvió error: %v", err)
	}
	if reserva.Estado != models.EstadoReservaCancelada || reserva.MotivoCancelacion != motivo {
		t.Errorf("reserva mal cancelada: estado=%q motivo=%q", reserva.Estado, reserva.MotivoCancelacion)
	}
	if guardada := repo.porID[5]; guardada.MotivoCancelacion != motivo {
		t.Errorf("el motivo debió persistir en el repositorio, quedó %q", guardada.MotivoCancelacion)
	}
	if vehiculos.porID[1].Estado != models.EstadoDisponible {
		t.Errorf("el vehículo debió liberarse, quedó %q", vehiculos.porID[1].Estado)
	}
}

func TestCancelarPropiaNoRegistraMotivo(t *testing.T) {
	servicio, repo, _ := nuevoServicioReservas(t)
	reserva := guardarReserva(repo, 6, 7, 1, time.Now().Add(time.Hour))
	reserva.MotivoCancelacion = "" // sin motivo previo

	if _, err := servicio.Cancelar(context.Background(), 6, 7); err != nil {
		t.Fatalf("Cancelar devolvió error: %v", err)
	}
	if guardada := repo.porID[6]; guardada.Estado != models.EstadoReservaCancelada || guardada.MotivoCancelacion != "" {
		t.Errorf("la cancelación propia debía quedar sin motivo: estado=%q motivo=%q",
			guardada.Estado, guardada.MotivoCancelacion)
	}
}