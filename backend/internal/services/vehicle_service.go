package services

import (
	"errors"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/repositories"
)

type VehicleService struct {
	vehicleRepo *repositories.VehicleRepository
}

func NewVehicleService(vehicleRepo *repositories.VehicleRepository) *VehicleService {
	return &VehicleService{vehicleRepo: vehicleRepo}
}

func (s *VehicleService) Create(vehicle *models.Vehicle) error {
	if vehicle.Status == "" {
		vehicle.Status = models.VehicleAvailable
	}
	return s.vehicleRepo.Create(vehicle)
}

func (s *VehicleService) GetByID(id uint) (*models.Vehicle, error) {
	return s.vehicleRepo.FindByID(id)
}

func (s *VehicleService) Update(vehicle *models.Vehicle) error {
	existing, err := s.vehicleRepo.FindByID(vehicle.ID)
	if err != nil {
		return errors.New("vehiculo no encontrado")
	}
	vehicle.CreatedAt = existing.CreatedAt
	vehicle.Status = existing.Status
	return s.vehicleRepo.Update(vehicle)
}

func (s *VehicleService) Delete(id uint) error {
	_, err := s.vehicleRepo.FindByID(id)
	if err != nil {
		return errors.New("vehiculo no encontrado")
	}
	return s.vehicleRepo.Delete(id)
}

func (s *VehicleService) List(filter models.VehicleFilter) ([]models.Vehicle, int64, error) {
	return s.vehicleRepo.List(filter)
}

func (s *VehicleService) GetBrands() ([]string, error) {
	return s.vehicleRepo.GetDistinctBrands()
}

func (s *VehicleService) GetStats() (map[string]interface{}, error) {
	return s.vehicleRepo.GetStats()
}
