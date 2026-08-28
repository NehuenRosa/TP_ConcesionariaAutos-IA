package repositories

import (
	"context"
	"fmt"
	"strings"

	"concesionaria/backend/internal/models"

	"gorm.io/gorm"
)

// FiltrosBusqueda agrupa los criterios opcionales de búsqueda y filtrado del
// catálogo público. Los rangos usan punteros para distinguir "sin filtro" de
// un valor cero.
type FiltrosBusqueda struct {
	Busqueda      string
	Marca         string
	Modelo        string
	AnioMin       *int
	AnioMax       *int
	PrecioMin     *float64
	PrecioMax     *float64
	Tipo          string
	Combustible   string
	Condicion     string
	OrdenPor      string
	OrdenDireccion string
}

// VehiculoRepository define el acceso a datos de vehículos sobre GORM.
type VehiculoRepository interface {
	// Listar devuelve los vehículos que cumplen el estado y los criterios de
	// búsqueda indicados, paginados, junto con el total de registros que
	// cumplen las condiciones.
	Listar(ctx context.Context, estado string, filtros FiltrosBusqueda, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	// ListarParaGestion devuelve los vehículos paginados para la administración.
	// Si estado está vacío incluye todos los estados; si no, filtra por el estado.
	ListarParaGestion(ctx context.Context, estado string, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	// ObtenerPorID devuelve un vehículo con su galería de imágenes.
	ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error)
	// Crear persiste un vehículo nuevo con sus imágenes.
	Crear(ctx context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error)
	// Actualizar actualiza la ficha y el estado de un vehículo, reemplazando su
	// galería de imágenes por la lista recibida.
	Actualizar(ctx context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error)
	// DarDeBaja cambia el estado del vehículo a dado_de_baja.
	DarDeBaja(ctx context.Context, id uint) error
}

// vehiculoRepository implementa VehiculoRepository sobre GORM.
type vehiculoRepository struct {
	base *gorm.DB
}

// NuevoVehiculoRepository crea un repositorio de vehículos.
func NuevoVehiculoRepository(base *gorm.DB) VehiculoRepository {
	return &vehiculoRepository{base: base}
}

// Listar cuenta y devuelve la página solicitada de vehículos que cumplen el
// estado y los criterios de búsqueda indicados.
func (r *vehiculoRepository) Listar(ctx context.Context, estado string, filtros FiltrosBusqueda, pagina int, tamano int) ([]models.Vehiculo, int64, error) {
	consulta := construirConsultaBusqueda(r.base.WithContext(ctx), filtros).
		Where("estado = ?", estado)

	var total int64
	if err := consulta.Model(&models.Vehiculo{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("contar vehículos: %w", err)
	}

	columnaOrden, direccion := ordenPorDefecto(filtros)
	consulta = consulta.Order(columnaOrden + " " + direccion)

	var vehiculos []models.Vehiculo
	err := consulta.
		Model(&models.Vehiculo{}).
		Preload("Imagenes").
		Offset((pagina - 1) * tamano).
		Limit(tamano).
		Find(&vehiculos).Error
	if err != nil {
		return nil, 0, fmt.Errorf("listar vehículos: %w", err)
	}

	return vehiculos, total, nil
}

// construirConsultaBusqueda aplica a la consulta los criterios opcionales de
// búsqueda y filtrado del catálogo.
func construirConsultaBusqueda(consulta *gorm.DB, filtros FiltrosBusqueda) *gorm.DB {
	if texto := strings.TrimSpace(filtros.Busqueda); texto != "" {
		patron := "%" + escaparComodines(texto) + "%"
		consulta = consulta.Where("marca ILIKE ? OR modelo ILIKE ?", patron, patron)
	}
	if filtros.Marca != "" {
		consulta = consulta.Where("marca = ?", filtros.Marca)
	}
	if filtros.Modelo != "" {
		consulta = consulta.Where("modelo = ?", filtros.Modelo)
	}
	if filtros.AnioMin != nil {
		consulta = consulta.Where("anio >= ?", *filtros.AnioMin)
	}
	if filtros.AnioMax != nil {
		consulta = consulta.Where("anio <= ?", *filtros.AnioMax)
	}
	if filtros.PrecioMin != nil {
		consulta = consulta.Where("precio >= ?", *filtros.PrecioMin)
	}
	if filtros.PrecioMax != nil {
		consulta = consulta.Where("precio <= ?", *filtros.PrecioMax)
	}
	if filtros.Tipo != "" {
		consulta = consulta.Where("LOWER(tipo) = LOWER(?)", filtros.Tipo)
	}
	if filtros.Combustible != "" {
		consulta = consulta.Where("combustible = ?", filtros.Combustible)
	}
	if filtros.Condicion != "" {
		consulta = consulta.Where("condicion = ?", filtros.Condicion)
	}
	return consulta
}

// ordenPorDefecto resuelve la columna y la dirección de ordenamiento con sus
// valores por defecto (año descendente) y una lista blanca para evitar
// inyección SQL por el parámetro de orden.
func ordenPorDefecto(filtros FiltrosBusqueda) (string, string) {
	columna := "anio"
	switch filtros.OrdenPor {
	case "precio":
		columna = "precio"
	case "anio":
		columna = "anio"
	}

	direccion := "DESC"
	switch strings.ToLower(filtros.OrdenDireccion) {
	case "asc":
		direccion = "ASC"
	case "desc":
		direccion = "DESC"
	}
	return columna, direccion
}

// escaparComodines neutraliza los caracteres comodín de ILIKE para que el
// usuario busque literales.
func escaparComodines(texto string) string {
	texto = strings.ReplaceAll(texto, `\`, `\\`)
	texto = strings.ReplaceAll(texto, `%`, `\%`)
	texto = strings.ReplaceAll(texto, `_`, `\_`)
	return texto
}

// ListarParaGestion cuenta y devuelve la página solicitada de vehículos para la
// administración. Con estado vacío incluye todos los estados del stock.
func (r *vehiculoRepository) ListarParaGestion(ctx context.Context, estado string, pagina int, tamano int) ([]models.Vehiculo, int64, error) {
	consulta := r.base.WithContext(ctx).Model(&models.Vehiculo{})
	if estado != "" {
		consulta = consulta.Where("estado = ?", estado)
	}

	var total int64
	if err := consulta.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("contar vehículos para gestión: %w", err)
	}

	var vehiculos []models.Vehiculo
	err := consulta.
		Preload("Imagenes").
		Offset((pagina - 1) * tamano).
		Limit(tamano).
		Find(&vehiculos).Error
	if err != nil {
		return nil, 0, fmt.Errorf("listar vehículos para gestión: %w", err)
	}

	return vehiculos, total, nil
}

// ObtenerPorID devuelve un vehículo con sus imágenes o un error de GORM
// (ErrRecordNotFound si no existe).
func (r *vehiculoRepository) ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error) {
	var vehiculo models.Vehiculo
	err := r.base.WithContext(ctx).
		Preload("Imagenes").
		First(&vehiculo, id).Error
	if err != nil {
		return nil, err
	}
	return &vehiculo, nil
}

// Crear persiste el vehículo y su galería de imágenes, y devuelve el registro
// completo con los IDs asignados.
func (r *vehiculoRepository) Crear(ctx context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(vehiculo).Error; err != nil {
			return err
		}
		for i := range vehiculo.Imagenes {
			vehiculo.Imagenes[i].VehiculoID = vehiculo.ID
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("crear vehículo: %w", err)
	}
	return r.ObtenerPorID(ctx, vehiculo.ID)
}

// Actualizar reemplaza las imágenes existentes por la lista recibida, guarda la
// ficha técnica y el estado, y devuelve el registro actualizado.
func (r *vehiculoRepository) Actualizar(ctx context.Context, vehiculo *models.Vehiculo) (*models.Vehiculo, error) {
	err := r.base.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("vehiculo_id = ?", vehiculo.ID).Delete(&models.Imagen{}).Error; err != nil {
			return err
		}

		imagenes := vehiculo.Imagenes
		vehiculo.Imagenes = nil
		if err := tx.Save(vehiculo).Error; err != nil {
			return err
		}

		for i := range imagenes {
			imagenes[i].ID = 0
			imagenes[i].VehiculoID = vehiculo.ID
		}
		if len(imagenes) > 0 {
			if err := tx.Create(&imagenes).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("actualizar vehículo: %w", err)
	}
	return r.ObtenerPorID(ctx, vehiculo.ID)
}

// DarDeBaja actualiza el estado del vehículo a dado_de_baja. Devuelve
// gorm.ErrRecordNotFound si el vehículo no existe.
func (r *vehiculoRepository) DarDeBaja(ctx context.Context, id uint) error {
	resultado := r.base.WithContext(ctx).
		Model(&models.Vehiculo{}).
		Where("id = ?", id).
		Update("estado", models.EstadoDadoDeBaja)
	if resultado.Error != nil {
		return fmt.Errorf("dar de baja vehículo: %w", resultado.Error)
	}
	if resultado.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
