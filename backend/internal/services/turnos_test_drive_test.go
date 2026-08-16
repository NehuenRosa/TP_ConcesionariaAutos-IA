package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"concesionaria/backend/internal/models"
)

func vehiculoDisponible(id uint) *models.Vehiculo {
	return &models.Vehiculo{ID: id, Marca: "Volkswagen", Modelo: "Gol", Estado: models.EstadoDisponible}
}

func fechaDeHoy() string {
	return time.Now().Format("2006-01-02")
}

// fechaDeManana devuelve la fecha de mañana: los tests que crean turnos usan
// esta fecha para no depender de la hora actual (el servicio rechaza franjas
// de hoy que ya comenzaron).
func fechaDeManana() string {
	return time.Now().AddDate(0, 0, 1).Format("2006-01-02")
}

func TestSolicitarExitoso(t *testing.T) {
	vehiculos := nuevoFakeVehiculoRepository()
	vehiculos.porID[1] = vehiculoDisponible(1)
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), vehiculos)

	creado, err := servicio.Solicitar(context.Background(), 7, 1, fechaDeManana(), "10:00")
	if err != nil {
		t.Fatalf("Solicitar devolvió error: %v", err)
	}
	if creado.ID == 0 {
		t.Error("el turno creado no tiene ID asignado")
	}
	if creado.Estado != models.EstadoTurnoSolicitado {
		t.Errorf("estado esperado %q, obtenido %q", models.EstadoTurnoSolicitado, creado.Estado)
	}
	if creado.ClienteID != 7 || creado.VehiculoID != 1 {
		t.Errorf("relaciones incorrectas: cliente %d, vehículo %d", creado.ClienteID, creado.VehiculoID)
	}
	if creado.Franja != "10:00" {
		t.Errorf("franja esperada %q, obtenida %q", "10:00", creado.Franja)
	}
}

func TestSolicitarFranjaInvalida(t *testing.T) {
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Solicitar(context.Background(), 7, 1, fechaDeManana(), "noche")
	if !errors.Is(err, ErrDatosTurnoInvalidos) {
		t.Errorf("se esperaba ErrDatosTurnoInvalidos, se obtuvo %v", err)
	}
}

func TestSolicitarFechaInvalida(t *testing.T) {
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Solicitar(context.Background(), 7, 1, "31/12/2025", "10:00")
	if !errors.Is(err, ErrDatosTurnoInvalidos) {
		t.Errorf("formato inválido: se esperaba ErrDatosTurnoInvalidos, se obtuvo %v", err)
	}

	pasado := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err = servicio.Solicitar(context.Background(), 7, 1, pasado, "10:00")
	if !errors.Is(err, ErrDatosTurnoInvalidos) {
		t.Errorf("fecha pasada: se esperaba ErrDatosTurnoInvalidos, se obtuvo %v", err)
	}
}

func TestSolicitarVehiculoInexistente(t *testing.T) {
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Solicitar(context.Background(), 7, 99, fechaDeManana(), "10:00")
	if !errors.Is(err, ErrVehiculoNoDisponible) {
		t.Errorf("se esperaba ErrVehiculoNoDisponible, se obtuvo %v", err)
	}
}

func TestSolicitarVehiculoNoDisponible(t *testing.T) {
	vehiculos := nuevoFakeVehiculoRepository()
	vehiculos.porID[1] = &models.Vehiculo{ID: 1, Estado: models.EstadoReservado}
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), vehiculos)

	_, err := servicio.Solicitar(context.Background(), 7, 1, fechaDeManana(), "10:00")
	if !errors.Is(err, ErrVehiculoNoDisponible) {
		t.Errorf("se esperaba ErrVehiculoNoDisponible, se obtuvo %v", err)
	}
}

func TestSolicitarSuperpuesto(t *testing.T) {
	vehiculos := nuevoFakeVehiculoRepository()
	vehiculos.porID[1] = vehiculoDisponible(1)
	turnos := nuevoFakeTurnoRepository()
	turnos.superpuesto = true
	servicio := NuevoTurnoTestDriveService(turnos, vehiculos)

	_, err := servicio.Solicitar(context.Background(), 7, 1, fechaDeManana(), "10:00")
	if !errors.Is(err, ErrTurnoSuperpuesto) {
		t.Errorf("se esperaba ErrTurnoSuperpuesto, se obtuvo %v", err)
	}
}

func TestSolicitarFranjaDeHoyYaComenzada(t *testing.T) {
	vehiculos := nuevoFakeVehiculoRepository()
	vehiculos.porID[1] = vehiculoDisponible(1)
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), vehiculos)

	// Se arma una fecha de hoy y una franja a la que no se puede llegar por la
	// hora actual; se espera ErrTurnoEnPasado. Se usa la primera franja del
	// catálogo, que a cualquier hora de la jornada ya comenzó si el test corre
	// después de su inicio (nunca es un horario futuro para hoy al cierre de
	// la jornada). Para evitar falsos positivos se saltea el test si la franja
	// todavía no comenzó (ej. tests que corren muy temprano).
	primera := models.FranjasDisponibles()[0].ID
	ahora := time.Now()
	inicio, _ := time.Parse("15:04", primera)
	comienzo := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), inicio.Hour(), inicio.Minute(), 0, 0, ahora.Location())
	if comienzo.After(ahora) {
		t.Skipf("la franja %q todavía no comenzó hoy", primera)
	}

	_, err := servicio.Solicitar(context.Background(), 7, 1, fechaDeHoy(), primera)
	if !errors.Is(err, ErrTurnoEnPasado) {
		t.Errorf("se esperaba ErrTurnoEnPasado, se obtuvo %v", err)
	}
}

func TestEsHoraDisponibleEn(t *testing.T) {
	utc := time.UTC
	base := time.Date(2026, 8, 15, 11, 30, 0, 0, utc)

	casos := []struct {
		nombre  string
		fecha   string
		franja  string
		ahora   time.Time
		esperar bool
	}{
		{"franja futura del mismo día", "2026-08-15", "14:00", base, true},
		{"franja pasada del mismo día", "2026-08-15", "10:00", base, false},
		{"franja que comienza en el límite", "2026-08-15", "11:30", base, false},
		{"fecha distinta siempre disponible", "2026-08-16", "09:00", base, true},
		{"franja inexistente no disponible", "2026-08-15", "00:00", base, false},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			if resultado := esHoraDisponibleEn(c.fecha, c.franja, c.ahora); resultado != c.esperar {
				t.Errorf("esHoraDisponibleEn(%q, %q, %v) = %v, se esperaba %v",
					c.fecha, c.franja, c.ahora, resultado, c.esperar)
			}
		})
	}
}

func TestCancelarTurnoPropio(t *testing.T) {
	turnos := nuevoFakeTurnoRepository()
	turnos.porID[1] = &models.TurnoTestDrive{
		ID: 1, VehiculoID: 1, ClienteID: 7, Fecha: fechaDeManana(), Franja: "10:00",
		Estado: models.EstadoTurnoSolicitado,
	}
	servicio := NuevoTurnoTestDriveService(turnos, nuevoFakeVehiculoRepository())

	cancelado, err := servicio.Cancelar(context.Background(), 1, 7)
	if err != nil {
		t.Fatalf("Cancelar devolvió error: %v", err)
	}
	if cancelado.Estado != models.EstadoTurnoCancelado {
		t.Errorf("estado esperado %q, obtenido %q", models.EstadoTurnoCancelado, cancelado.Estado)
	}
}

func TestCancelarTurnoAjenoSeTrataComoInexistente(t *testing.T) {
	turnos := nuevoFakeTurnoRepository()
	turnos.porID[1] = &models.TurnoTestDrive{
		ID: 1, VehiculoID: 1, ClienteID: 7, Fecha: fechaDeManana(), Franja: "10:00",
		Estado: models.EstadoTurnoSolicitado,
	}
	servicio := NuevoTurnoTestDriveService(turnos, nuevoFakeVehiculoRepository())

	_, err := servicio.Cancelar(context.Background(), 1, 8)
	if !errors.Is(err, ErrTurnoNoEncontrado) {
		t.Errorf("se esperaba ErrTurnoNoEncontrado, se obtuvo %v", err)
	}
}

func TestCancelarTurnoNoActivo(t *testing.T) {
	turnos := nuevoFakeTurnoRepository()
	turnos.porID[1] = &models.TurnoTestDrive{
		ID: 1, VehiculoID: 1, ClienteID: 7, Fecha: fechaDeManana(), Franja: "10:00",
		Estado: models.EstadoTurnoCompletado,
	}
	servicio := NuevoTurnoTestDriveService(turnos, nuevoFakeVehiculoRepository())

	_, err := servicio.Cancelar(context.Background(), 1, 7)
	if !errors.Is(err, ErrTurnoEstadoInvalido) {
		t.Errorf("se esperaba ErrTurnoEstadoInvalido, se obtuvo %v", err)
	}
}

func TestCancelarTurnoInexistente(t *testing.T) {
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Cancelar(context.Background(), 99, 7)
	if !errors.Is(err, ErrTurnoNoEncontrado) {
		t.Errorf("se esperaba ErrTurnoNoEncontrado, se obtuvo %v", err)
	}
}

func TestConfirmarExitoso(t *testing.T) {
	turnos := nuevoFakeTurnoRepository()
	turnos.porID[1] = &models.TurnoTestDrive{ID: 1, Estado: models.EstadoTurnoSolicitado}
	servicio := NuevoTurnoTestDriveService(turnos, nuevoFakeVehiculoRepository())

	confirmado, err := servicio.Confirmar(context.Background(), 1)
	if err != nil {
		t.Fatalf("Confirmar devolvió error: %v", err)
	}
	if confirmado.Estado != models.EstadoTurnoConfirmado {
		t.Errorf("estado esperado %q, obtenido %q", models.EstadoTurnoConfirmado, confirmado.Estado)
	}
}

func TestConfirmarEstadoInvalido(t *testing.T) {
	turnos := nuevoFakeTurnoRepository()
	turnos.porID[1] = &models.TurnoTestDrive{ID: 1, Estado: models.EstadoTurnoCancelado}
	servicio := NuevoTurnoTestDriveService(turnos, nuevoFakeVehiculoRepository())

	_, err := servicio.Confirmar(context.Background(), 1)
	if !errors.Is(err, ErrTurnoEstadoInvalido) {
		t.Errorf("se esperaba ErrTurnoEstadoInvalido, se obtuvo %v", err)
	}
}

func TestCompletarExitoso(t *testing.T) {
	turnos := nuevoFakeTurnoRepository()
	turnos.porID[1] = &models.TurnoTestDrive{ID: 1, Estado: models.EstadoTurnoConfirmado}
	servicio := NuevoTurnoTestDriveService(turnos, nuevoFakeVehiculoRepository())

	completado, err := servicio.Completar(context.Background(), 1)
	if err != nil {
		t.Fatalf("Completar devolvió error: %v", err)
	}
	if completado.Estado != models.EstadoTurnoCompletado {
		t.Errorf("estado esperado %q, obtenido %q", models.EstadoTurnoCompletado, completado.Estado)
	}
}

func TestCompletarEstadoInvalido(t *testing.T) {
	turnos := nuevoFakeTurnoRepository()
	turnos.porID[1] = &models.TurnoTestDrive{ID: 1, Estado: models.EstadoTurnoSolicitado}
	servicio := NuevoTurnoTestDriveService(turnos, nuevoFakeVehiculoRepository())

	_, err := servicio.Completar(context.Background(), 1)
	if !errors.Is(err, ErrTurnoEstadoInvalido) {
		t.Errorf("se esperaba ErrTurnoEstadoInvalido, se obtuvo %v", err)
	}
}

func TestCancelarComoVendedorExitoso(t *testing.T) {
	turnos := nuevoFakeTurnoRepository()
	turnos.porID[1] = &models.TurnoTestDrive{ID: 1, Estado: models.EstadoTurnoConfirmado}
	servicio := NuevoTurnoTestDriveService(turnos, nuevoFakeVehiculoRepository())

	cancelado, err := servicio.CancelarComoVendedor(context.Background(), 1)
	if err != nil {
		t.Fatalf("CancelarComoVendedor devolvió error: %v", err)
	}
	if cancelado.Estado != models.EstadoTurnoCancelado {
		t.Errorf("estado esperado %q, obtenido %q", models.EstadoTurnoCancelado, cancelado.Estado)
	}
}

func TestListarFiltroInvalido(t *testing.T) {
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Listar(context.Background(), "estado-raro")
	if !errors.Is(err, ErrFiltroEstadoTurnoInvalido) {
		t.Errorf("se esperaba ErrFiltroEstadoTurnoInvalido, se obtuvo %v", err)
	}
}

func TestListarSinFiltroYConFiltroValido(t *testing.T) {
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), nuevoFakeVehiculoRepository())
	if _, err := servicio.Listar(context.Background(), ""); err != nil {
		t.Errorf("sin filtro no debería fallar: %v", err)
	}
	if _, err := servicio.Listar(context.Background(), models.EstadoTurnoSolicitado); err != nil {
		t.Errorf("filtro válido no debería fallar: %v", err)
	}
}

func TestFranjas(t *testing.T) {
	servicio := NuevoTurnoTestDriveService(nuevoFakeTurnoRepository(), nuevoFakeVehiculoRepository())
	franjas := servicio.Franjas()
	esperadas := 7
	if len(franjas) != esperadas {
		t.Fatalf("se esperaban %d franjas, hay %d", esperadas, len(franjas))
	}
	for _, franja := range franjas {
		if !models.FranjaValida(franja.ID) {
			t.Errorf("la franja %q no figura en el catálogo", franja.ID)
		}
		inicio, err := time.Parse("15:04", franja.Inicio)
		if err != nil {
			t.Errorf("la franja %q tiene inicio inválido %q", franja.ID, franja.Inicio)
			continue
		}
		if fin := inicio.Add(time.Hour).Format("15:04"); fin != franja.Fin {
			t.Errorf("la franja %q debería terminar a las %s, termina a las %s", franja.ID, fin, franja.Fin)
		}
	}
	if models.FranjaValida("manana") {
		t.Error("'manana' ya no debería ser una franja válida")
	}
	if !models.FranjaValida("10:00") {
		t.Error("'10:00' debería ser una franja válida")
	}
}

func TestEsFechaValida(t *testing.T) {
	if !esFechaValida(fechaDeHoy()) {
		t.Error("hoy debería ser una fecha válida")
	}
	futuro := time.Now().AddDate(0, 0, 7).Format("2006-01-02")
	if !esFechaValida(futuro) {
		t.Error("una fecha futura debería ser válida")
	}
	pasado := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	if esFechaValida(pasado) {
		t.Error("una fecha pasada no debería ser válida")
	}
	if esFechaValida("31/12/2025") {
		t.Error("un formato inválido no debería ser aceptado")
	}
}
