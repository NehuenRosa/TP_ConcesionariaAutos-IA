package services

import (
	"context"
	"errors"
	"testing"

	"concesionaria/backend/internal/models"
)

func nuevoServicioConsultas(consultas *fakeConsultaRepository, vehiculos *fakeVehiculoRepository) ConsultaService {
	return NuevoConsultaService(consultas, vehiculos)
}

func TestCrearConsultaExitoso(t *testing.T) {
	vehiculos := nuevoFakeVehiculoRepository()
	vehiculos.porID[1] = vehiculoDisponible(1)
	consultas := nuevoFakeConsultaRepository()
	servicio := nuevoServicioConsultas(consultas, vehiculos)

	creada, err := servicio.Crear(context.Background(), 5, 1, "  ¿Tiene garantía?  ")
	if err != nil {
		t.Fatalf("Crear devolvió error: %v", err)
	}
	if creada.ID == 0 {
		t.Error("la consulta creada no tiene ID asignado")
	}
	if creada.Estado != models.EstadoPendiente {
		t.Errorf("estado esperado %q, obtenido %q", models.EstadoPendiente, creada.Estado)
	}
	if len(creada.Mensajes) != 1 {
		t.Fatalf("se esperaba 1 mensaje inicial, hay %d", len(creada.Mensajes))
	}
	if creada.Mensajes[0].Contenido != "¿Tiene garantía?" {
		t.Errorf("mensaje sin recortar: %q", creada.Mensajes[0].Contenido)
	}
	if creada.Mensajes[0].RemitenteID != 5 {
		t.Errorf("remitente esperado 5, obtenido %d", creada.Mensajes[0].RemitenteID)
	}
	if creada.VendedorID != nil {
		t.Error("una consulta nueva no debería tener vendedor asignado")
	}
}

func TestCrearConsultaMensajeVacio(t *testing.T) {
	servicio := nuevoServicioConsultas(nuevoFakeConsultaRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Crear(context.Background(), 5, 1, "   ")
	if !errors.Is(err, ErrMensajeVacio) {
		t.Errorf("se esperaba ErrMensajeVacio, se obtuvo %v", err)
	}
}

func TestCrearConsultaVehiculoNoDisponible(t *testing.T) {
	vehiculos := nuevoFakeVehiculoRepository()
	vehiculos.porID[1] = &models.Vehiculo{ID: 1, Estado: models.EstadoVendido}
	servicio := nuevoServicioConsultas(nuevoFakeConsultaRepository(), vehiculos)

	_, err := servicio.Crear(context.Background(), 5, 1, "¿Hay stock?")
	if !errors.Is(err, ErrVehiculoNoDisponible) {
		t.Errorf("se esperaba ErrVehiculoNoDisponible, se obtuvo %v", err)
	}
}

func TestCrearConsultaVehiculoInexistente(t *testing.T) {
	servicio := nuevoServicioConsultas(nuevoFakeConsultaRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Crear(context.Background(), 5, 99, "¿Hay stock?")
	if !errors.Is(err, ErrVehiculoNoDisponible) {
		t.Errorf("se esperaba ErrVehiculoNoDisponible, se obtuvo %v", err)
	}
}

func TestTomarExitoso(t *testing.T) {
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, Estado: models.EstadoPendiente}
	consultas.tomable[1] = true
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	tomada, err := servicio.Tomar(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Tomar devolvió error: %v", err)
	}
	if tomada.VendedorID == nil || *tomada.VendedorID != 10 {
		t.Errorf("vendedor esperado 10, obtenido %v", tomada.VendedorID)
	}
	if tomada.Estado != models.EstadoEnConversacion {
		t.Errorf("estado esperado %q, obtenido %q", models.EstadoEnConversacion, tomada.Estado)
	}
}

func TestTomarYaTomadaPorOtroVendedor(t *testing.T) {
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, Estado: models.EstadoEnConversacion}
	consultas.tomable[1] = false
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	_, err := servicio.Tomar(context.Background(), 1, 10)
	if !errors.Is(err, ErrConsultaNoPendiente) {
		t.Errorf("se esperaba ErrConsultaNoPendiente, se obtuvo %v", err)
	}
}

func TestTomarConsultaInexistente(t *testing.T) {
	servicio := nuevoServicioConsultas(nuevoFakeConsultaRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Tomar(context.Background(), 99, 10)
	if !errors.Is(err, ErrConsultaNoEncontrada) {
		t.Errorf("se esperaba ErrConsultaNoEncontrada, se obtuvo %v", err)
	}
}

func TestCerrarExitoso(t *testing.T) {
	vendedor := uint(10)
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, VendedorID: &vendedor, Estado: models.EstadoEnConversacion}
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	cerrada, err := servicio.Cerrar(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("Cerrar devolvió error: %v", err)
	}
	if cerrada.Estado != models.EstadoCerrada {
		t.Errorf("estado esperado %q, obtenido %q", models.EstadoCerrada, cerrada.Estado)
	}
}

func TestCerrarYaCerrada(t *testing.T) {
	vendedor := uint(10)
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, VendedorID: &vendedor, Estado: models.EstadoCerrada}
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	_, err := servicio.Cerrar(context.Background(), 1, 10)
	if !errors.Is(err, ErrConsultaYaCerrada) {
		t.Errorf("se esperaba ErrConsultaYaCerrada, se obtuvo %v", err)
	}
}

func TestCerrarNoVendedorAsignado(t *testing.T) {
	vendedor := uint(11)
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, VendedorID: &vendedor, Estado: models.EstadoEnConversacion}
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	_, err := servicio.Cerrar(context.Background(), 1, 10)
	if !errors.Is(err, ErrNoEsVendedorAsignado) {
		t.Errorf("se esperaba ErrNoEsVendedorAsignado, se obtuvo %v", err)
	}
}

func TestCerrarConsultaInexistente(t *testing.T) {
	servicio := nuevoServicioConsultas(nuevoFakeConsultaRepository(), nuevoFakeVehiculoRepository())
	_, err := servicio.Cerrar(context.Background(), 99, 10)
	if !errors.Is(err, ErrConsultaNoEncontrada) {
		t.Errorf("se esperaba ErrConsultaNoEncontrada, se obtuvo %v", err)
	}
}

func TestEliminarExitoso(t *testing.T) {
	vendedor := uint(10)
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, VendedorID: &vendedor, Estado: models.EstadoCerrada}
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	if err := servicio.Eliminar(context.Background(), 1, 10); err != nil {
		t.Fatalf("Eliminar devolvió error: %v", err)
	}
	if _, ok := consultas.porID[1]; ok {
		t.Error("la consulta no se eliminó del repositorio")
	}
}

func TestEliminarNoCerrada(t *testing.T) {
	vendedor := uint(10)
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, VendedorID: &vendedor, Estado: models.EstadoEnConversacion}
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	err := servicio.Eliminar(context.Background(), 1, 10)
	if !errors.Is(err, ErrConsultaNoCerrada) {
		t.Errorf("se esperaba ErrConsultaNoCerrada, se obtuvo %v", err)
	}
}

func TestEliminarNoVendedorAsignado(t *testing.T) {
	vendedor := uint(11)
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, VendedorID: &vendedor, Estado: models.EstadoCerrada}
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	err := servicio.Eliminar(context.Background(), 1, 10)
	if !errors.Is(err, ErrNoEsVendedorAsignado) {
		t.Errorf("se esperaba ErrNoEsVendedorAsignado, se obtuvo %v", err)
	}
}

func TestEsParticipante(t *testing.T) {
	vendedor := uint(10)
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, VendedorID: &vendedor, Estado: models.EstadoEnConversacion}
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	cliente, err := servicio.EsParticipante(context.Background(), 1, 5)
	if err != nil || !cliente {
		t.Errorf("el cliente 5 debería participar: %v, %v", cliente, err)
	}

	vendedorOk, err := servicio.EsParticipante(context.Background(), 1, 10)
	if err != nil || !vendedorOk {
		t.Errorf("el vendedor 10 debería participar: %v, %v", vendedorOk, err)
	}

	ajeno, err := servicio.EsParticipante(context.Background(), 1, 99)
	if err != nil || ajeno {
		t.Errorf("el usuario 99 no debería participar: %v, %v", ajeno, err)
	}
}

func TestEsParticipanteConsultaInexistente(t *testing.T) {
	servicio := nuevoServicioConsultas(nuevoFakeConsultaRepository(), nuevoFakeVehiculoRepository())
	participa, err := servicio.EsParticipante(context.Background(), 99, 5)
	if err != nil || participa {
		t.Errorf("una consulta inexistente no debería reportar participación: %v, %v", participa, err)
	}
}

func TestEsVendedorAsignado(t *testing.T) {
	vendedor := uint(10)
	consultas := nuevoFakeConsultaRepository()
	consultas.porID[1] = &models.Consulta{ID: 1, VehiculoID: 1, ClienteID: 5, VendedorID: &vendedor, Estado: models.EstadoEnConversacion}
	servicio := nuevoServicioConsultas(consultas, nuevoFakeVehiculoRepository())

	es, err := servicio.EsVendedorAsignado(context.Background(), 1, 10)
	if err != nil || !es {
		t.Errorf("el vendedor 10 debería estar asignado: %v, %v", es, err)
	}

	noEs, err := servicio.EsVendedorAsignado(context.Background(), 1, 11)
	if err != nil || noEs {
		t.Errorf("el vendedor 11 no debería estar asignado: %v, %v", noEs, err)
	}
}
