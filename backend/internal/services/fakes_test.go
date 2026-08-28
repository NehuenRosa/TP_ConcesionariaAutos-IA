package services

import (
	"context"
	"time"

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
	porID        map[uint]*models.TurnoTestDrive
	superpuesto  bool
	clienteDuplic bool
	ocupadas     []string
	siguiente    uint
}

func nuevoFakeTurnoRepository() *fakeTurnoRepository {
	return &fakeTurnoRepository{
		porID:     make(map[uint]*models.TurnoTestDrive),
		siguiente: 1,
	}
}

func (f *fakeTurnoRepository) CrearSiDisponible(_ context.Context, turno *models.TurnoTestDrive) (*models.TurnoTestDrive, error) {
	if f.clienteDuplic {
		return nil, repositories.ErrClienteYaTieneTurno
	}
	if f.superpuesto {
		return nil, repositories.ErrFranjaOcupada
	}
	turno.ID = f.siguiente
	f.siguiente++
	f.porID[turno.ID] = turno
	return turno, nil
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

func (f *fakeTurnoRepository) Ocupadas(_ context.Context, _ uint, _ string) ([]string, error) {
	return f.ocupadas, nil
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

type fakeCotizacionRepository struct {
	porID            map[uint]*models.Cotizacion
	siguiente        uint
	siguienteMensaje uint
}

func nuevoFakeCotizacionRepository() *fakeCotizacionRepository {
	return &fakeCotizacionRepository{
		porID:            make(map[uint]*models.Cotizacion),
		siguiente:        1,
		siguienteMensaje: 1,
	}
}

func (f *fakeCotizacionRepository) Crear(_ context.Context, cotizacion *models.Cotizacion) (*models.Cotizacion, error) {
	cotizacion.ID = f.siguiente
	f.siguiente++
	for i := range cotizacion.Mensajes {
		cotizacion.Mensajes[i].CotizacionID = cotizacion.ID
		cotizacion.Mensajes[i].ID = f.siguienteMensaje
		f.siguienteMensaje++
	}
	// Se guarda una copia: el servicio descifra en el modelo recibido y no debe
	// afectar al valor persistido.
	copia := copiarCotizacion(cotizacion)
	f.porID[cotizacion.ID] = &copia
	return &copia, nil
}

func (f *fakeCotizacionRepository) AgregarMensaje(_ context.Context, mensaje *models.MensajeCotizacion) error {
	cotizacion, ok := f.porID[mensaje.CotizacionID]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	mensaje.ID = f.siguienteMensaje
	f.siguienteMensaje++
	cotizacion.Mensajes = append(cotizacion.Mensajes, *mensaje)
	return nil
}

func (f *fakeCotizacionRepository) ObtenerPorID(_ context.Context, id uint) (*models.Cotizacion, error) {
	if cotizacion, ok := f.porID[id]; ok {
		copia := copiarCotizacion(cotizacion)
		return &copia, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeCotizacionRepository) ListarPorCliente(_ context.Context, clienteID uint) ([]models.Cotizacion, error) {
	var cotizaciones []models.Cotizacion
	for _, cotizacion := range f.porID {
		if cotizacion.ClienteID == clienteID {
			cotizaciones = append(cotizaciones, copiarCotizacion(cotizacion))
		}
	}
	return cotizaciones, nil
}

func (f *fakeCotizacionRepository) ListarBandeja(_ context.Context) ([]models.Cotizacion, error) {
	cotizaciones := make([]models.Cotizacion, 0, len(f.porID))
	for _, cotizacion := range f.porID {
		cotizaciones = append(cotizaciones, copiarCotizacion(cotizacion))
	}
	return cotizaciones, nil
}

func (f *fakeCotizacionRepository) Actualizar(_ context.Context, cotizacion *models.Cotizacion) error {
	copia := copiarCotizacion(cotizacion)
	f.porID[cotizacion.ID] = &copia
	return nil
}

func (f *fakeCotizacionRepository) ContarNoLeidosDeCliente(_ context.Context, clienteID uint) (int64, error) {
	var total int64
	for _, c := range f.porID {
		if c.ClienteID != clienteID {
			continue
		}
		for _, m := range c.Mensajes {
			if m.Remitente == models.RemitenteVendedor && !m.LeidoPorCliente {
				total++
			}
		}
	}
	return total, nil
}

func (f *fakeCotizacionRepository) ContarNoLeidosParaPersonal(_ context.Context, vendedorID uint) (int64, error) {
	var total int64
	for _, c := range f.porID {
		if c.Estado != models.EstadoCotizacionAbierta {
			continue
		}
		if c.VendedorID != nil && *c.VendedorID != vendedorID {
			continue
		}
		for _, m := range c.Mensajes {
			if m.Remitente == models.RemitenteCliente && !m.LeidoPorVendedor {
				total++
			}
		}
	}
	return total, nil
}

func (f *fakeCotizacionRepository) MarcarLeidasParaCliente(_ context.Context, cotizacionID uint) error {
	c := f.porID[cotizacionID]
	if c == nil {
		return nil
	}
	for i := range c.Mensajes {
		if c.Mensajes[i].Remitente != models.RemitenteCliente {
			c.Mensajes[i].LeidoPorCliente = true
		}
	}
	return nil
}

func (f *fakeCotizacionRepository) MarcarLeidasParaPersonal(_ context.Context, cotizacionID uint) error {
	c := f.porID[cotizacionID]
	if c == nil {
		return nil
	}
	for i := range c.Mensajes {
		if c.Mensajes[i].Remitente == models.RemitenteCliente {
			c.Mensajes[i].LeidoPorVendedor = true
		}
	}
	return nil
}

// copiarCotizacion devuelve una copia profunda (mensajes incluidos) para que
// las mutaciones del servicio no alteren el estado persistido del fake.
func copiarCotizacion(cotizacion *models.Cotizacion) models.Cotizacion {
	copia := *cotizacion
	copia.Mensajes = make([]models.MensajeCotizacion, len(cotizacion.Mensajes))
	copy(copia.Mensajes, cotizacion.Mensajes)
	return copia
}

// fakeGeneradorCotizacion devuelve una respuesta fija basada en el mensaje del
// cliente, sin llamar a ningún LLM.
type fakeGeneradorCotizacion struct {
	respuestas map[string]string
}

func (f *fakeGeneradorCotizacion) GenerarCotizacion(_ context.Context, _ models.Vehiculo, _ []TurnoChat, mensaje string) (string, error) {
	if respuesta, ok := f.respuestas[mensaje]; ok {
		return respuesta, nil
	}
	return "Respuesta para: " + mensaje, nil
}

// fakeReservaRepository simula la persistencia de reservas aplicando las mismas
// reglas de expiración que el repositorio real sobre GORM.
type fakeReservaRepository struct {
	porID         map[uint]*models.Reserva
	comprobantes  map[uint]*models.ComprobanteReserva
	vehiculos     *fakeVehiculoRepository
	cantidadLotes int64
	siguiente     uint
}

func nuevoFakeReservaRepository(vehiculos *fakeVehiculoRepository) *fakeReservaRepository {
	return &fakeReservaRepository{
		porID:        make(map[uint]*models.Reserva),
		comprobantes: make(map[uint]*models.ComprobanteReserva),
		vehiculos:    vehiculos,
		siguiente:    1,
	}
}

func (f *fakeReservaRepository) CrearYReservar(_ context.Context, reserva *models.Reserva) (*models.Reserva, error) {
	reserva.ID = f.siguiente
	f.siguiente++
	f.porID[reserva.ID] = reserva
	return reserva, nil
}

func (f *fakeReservaRepository) ObtenerPorID(_ context.Context, id uint) (*models.Reserva, error) {
	if reserva, ok := f.porID[id]; ok {
		return reserva, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeReservaRepository) ListarPorCliente(_ context.Context, _ uint) ([]models.Reserva, error) {
	return nil, nil
}

func (f *fakeReservaRepository) Listar(_ context.Context, _ string) ([]models.Reserva, error) {
	return nil, nil
}

func (f *fakeReservaRepository) ConfirmarVentaYMarcarVendido(_ context.Context, reserva *models.Reserva) (*models.Reserva, error) {
	guardada, ok := f.porID[reserva.ID]
	if !ok || !guardada.EsActiva() {
		return nil, repositories.ErrReservaYaNoActiva
	}
	guardada.Estado = models.EstadoReservaVendida
	f.marcarVehiculo(guardada.VehiculoID, models.EstadoVendido)
	return guardada, nil
}

func (f *fakeReservaRepository) CancelarYLiberar(_ context.Context, reserva *models.Reserva) (*models.Reserva, error) {
	guardada, ok := f.porID[reserva.ID]
	if !ok || !guardada.EsActiva() {
		return nil, repositories.ErrReservaYaNoActiva
	}
	guardada.Estado = models.EstadoReservaCancelada
	f.liberarVehiculo(guardada.VehiculoID)
	return guardada, nil
}

func (f *fakeReservaRepository) GuardarComprobante(_ context.Context, reserva *models.Reserva, comprobante *models.ComprobanteReserva) error {
	comprobante.ID = uint(len(f.comprobantes)) + 1
	f.comprobantes[reserva.ID] = comprobante
	guardada := f.porID[reserva.ID]
	guardada.ComprobanteEnviadoAt = reserva.ComprobanteEnviadoAt
	return nil
}

func (f *fakeReservaRepository) ObtenerComprobantePorReservaID(_ context.Context, reservaID uint) (*models.ComprobanteReserva, error) {
	if comprobante, ok := f.comprobantes[reservaID]; ok {
		return comprobante, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *fakeReservaRepository) ExpirarSiVencida(_ context.Context, reserva *models.Reserva) (bool, error) {
	guardada, ok := f.porID[reserva.ID]
	if !ok || !guardada.ComprobanteVencido(time.Now()) {
		return false, nil
	}
	guardada.Estado = models.EstadoReservaCancelada
	f.liberarVehiculo(guardada.VehiculoID)
	return true, nil
}

func (f *fakeReservaRepository) ExpirarVencidas(_ context.Context) (int64, error) {
	var cantidad int64
	for _, reserva := range f.porID {
		if reserva.ComprobanteVencido(time.Now()) {
			reserva.Estado = models.EstadoReservaCancelada
			f.liberarVehiculo(reserva.VehiculoID)
			cantidad++
		}
	}
	f.cantidadLotes += cantidad
	return cantidad, nil
}

func (f *fakeReservaRepository) liberarVehiculo(vehiculoID uint) {
	f.marcarVehiculo(vehiculoID, models.EstadoDisponible)
}

func (f *fakeReservaRepository) marcarVehiculo(vehiculoID uint, estado string) {
	if vehiculo, ok := f.vehiculos.porID[vehiculoID]; ok {
		vehiculo.Estado = estado
	}
}
