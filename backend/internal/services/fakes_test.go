package services

import (
	"context"

	"concesionaria/backend/internal/models"
	"concesionaria/backend/internal/repositories"

	"gorm.io/gorm"
)

type fakeUsuarioRepository struct {
	porEmail  map[string]*models.Usuario
	porID     map[uint]*models.Usuario
	siguiente uint
}

func nuevoFakeUsuarioRepository() *fakeUsuarioRepository {
	return &fakeUsuarioRepository{
		porEmail:  make(map[string]*models.Usuario),
		porID:     make(map[uint]*models.Usuario),
		siguiente: 1,
	}
}

func (f *fakeUsuarioRepository) Crear(_ context.Context, usuario *models.Usuario) (*models.Usuario, error) {
	if _, existe := f.porEmail[usuario.Email]; existe {
		return nil, gorm.ErrDuplicatedKey
	}
	usuario.ID = f.siguiente
	f.siguiente++
	f.porEmail[usuario.Email] = usuario
	f.porID[usuario.ID] = usuario
	return usuario, nil
}

func (f *fakeUsuarioRepository) ObtenerPorEmail(_ context.Context, email string) (*models.Usuario, error) {
	if usuario, ok := f.porEmail[email]; ok {
		return usuario, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeUsuarioRepository) ObtenerPorID(_ context.Context, id uint) (*models.Usuario, error) {
	if usuario, ok := f.porID[id]; ok {
		return usuario, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeUsuarioRepository) Listar(_ context.Context) ([]models.Usuario, error) {
	usuarios := make([]models.Usuario, 0, len(f.porID))
	for _, usuario := range f.porID {
		usuarios = append(usuarios, *usuario)
	}
	return usuarios, nil
}

func (f *fakeUsuarioRepository) Actualizar(_ context.Context, usuario *models.Usuario) error {
	f.porID[usuario.ID] = usuario
	f.porEmail[usuario.Email] = usuario
	return nil
}

func (f *fakeUsuarioRepository) Eliminar(_ context.Context, id uint) error {
	if usuario, ok := f.porID[id]; ok {
		delete(f.porID, id)
		delete(f.porEmail, usuario.Email)
	}
	return nil
}

type fakeVehiculoRepository struct {
	porID map[uint]*models.Vehiculo
}

func nuevoFakeVehiculoRepository() *fakeVehiculoRepository {
	return &fakeVehiculoRepository{porID: make(map[uint]*models.Vehiculo)}
}

func (f *fakeVehiculoRepository) ObtenerPorID(_ context.Context, id uint) (*models.Vehiculo, error) {
	if vehiculo, ok := f.porID[id]; ok {
		return vehiculo, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeVehiculoRepository) Listar(_ context.Context, _ string, _ repositories.FiltrosBusqueda, _ int, _ int) ([]models.Vehiculo, int64, error) {
	return nil, 0, nil
}

func (f *fakeVehiculoRepository) ListarParaGestion(_ context.Context, _ string, _ int, _ int) ([]models.Vehiculo, int64, error) {
	return nil, 0, nil
}

func (f *fakeVehiculoRepository) Crear(_ context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error) {
	return vehiculo, nil
}

func (f *fakeVehiculoRepository) Actualizar(_ context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error) {
	return vehiculo, nil
}

func (f *fakeVehiculoRepository) DarDeBaja(_ context.Context, _ uint) error {
	return nil
}

type fakeTurnoRepository struct {
	porID       map[uint]*models.TurnoTestDrive
	superpuesto bool
	siguiente   uint
}

func nuevoFakeTurnoRepository() *fakeTurnoRepository {
	return &fakeTurnoRepository{
		porID:     make(map[uint]*models.TurnoTestDrive),
		siguiente: 1,
	}
}

func (f *fakeTurnoRepository) CrearSiSinSuperposicion(_ context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, bool, error) {
	if f.superpuesto {
		return nil, false, nil
	}
	turno.ID = f.siguiente
	f.siguiente++
	f.porID[turno.ID] = turno
	return turno, true, nil
}

func (f *fakeTurnoRepository) ObtenerPorID(_ context.Context, id uint) (*models.TurnoTestDrive, error) {
	if turno, ok := f.porID[id]; ok {
		return turno, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeTurnoRepository) ListarPorCliente(_ context.Context, _ uint) ([]models.TurnoTestDrive, error) {
	return nil, nil
}

func (f *fakeTurnoRepository) Listar(_ context.Context, _ string) ([]models.TurnoTestDrive, error) {
	return nil, nil
}

func (f *fakeTurnoRepository) Actualizar(_ context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error) {
	f.porID[turno.ID] = turno
	return turno, nil
}

type fakeConsultaRepository struct {
	porID     map[uint]*models.Consulta
	tomable   map[uint]bool
	siguiente uint
}

func nuevoFakeConsultaRepository() *fakeConsultaRepository {
	return &fakeConsultaRepository{
		porID:     make(map[uint]*models.Consulta),
		tomable:   make(map[uint]bool),
		siguiente: 1,
	}
}

func (f *fakeConsultaRepository) Crear(_ context.Context, consulta *models.Consulta) (*models.Consulta, error) {
	consulta.ID = f.siguiente
	f.siguiente++
	for i := range consulta.Mensajes {
		consulta.Mensajes[i].ConsultaID = consulta.ID
	}
	f.porID[consulta.ID] = consulta
	return consulta, nil
}

func (f *fakeConsultaRepository) ObtenerPorID(_ context.Context, id uint) (*models.Consulta, error) {
	if consulta, ok := f.porID[id]; ok {
		return consulta, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeConsultaRepository) ListarPorCliente(_ context.Context, _ uint) ([]models.Consulta, error) {
	return nil, nil
}

func (f *fakeConsultaRepository) ListarPorVendedor(_ context.Context, _ uint) ([]models.Consulta, error) {
	return nil, nil
}

func (f *fakeConsultaRepository) ListarPorUsuario(_ context.Context, _ uint) ([]uint, error) {
	return nil, nil
}

func (f *fakeConsultaRepository) TomarSiPendiente(_ context.Context, consultaID uint, vendedorID uint) (bool, error) {
	if !f.tomable[consultaID] {
		return false, nil
	}
	consulta, ok := f.porID[consultaID]
	if !ok {
		return false, nil
	}
	vendedor := vendedorID
	consulta.VendedorID = &vendedor
	consulta.Estado = models.EstadoEnConversacion
	return true, nil
}

func (f *fakeConsultaRepository) Actualizar(_ context.Context, consulta *models.Consulta) (*models.Consulta, error) {
	f.porID[consulta.ID] = consulta
	return consulta, nil
}

func (f *fakeConsultaRepository) Eliminar(_ context.Context, id uint) error {
	if _, ok := f.porID[id]; ok {
		delete(f.porID, id)
		return nil
	}
	return gorm.ErrRecordNotFound
}
