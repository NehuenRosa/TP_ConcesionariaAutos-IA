package services

import (
	"context"

	"concesionaria/backend/internal/models"
)

// VehiculoService define el contrato de la lógica de negocio de vehículos.
// TODO: implementar junto con el caso de uso CU-03.
type VehiculoService interface {
	ListarDisponibles(ctx context.Context, pagina int, tamano int) ([]models.Vehiculo, int64, error)
	ObtenerPorID(ctx context.Context, id uint) (*models.Vehiculo, error)
}
