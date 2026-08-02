package repositories

import (
	"context"

	"concesionaria/backend/internal/models"
)

// VehiculoRepository define el acceso a datos de vehículos sobre GORM.
// TODO: implementar junto con el caso de uso CU-03.
type VehiculoRepository interface {
	ListarDisponibles(ctx context.Context, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error)
}
