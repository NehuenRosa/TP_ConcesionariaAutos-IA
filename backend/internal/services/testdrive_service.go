package services

import (
	"errors"
	"time"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/repositories"
)

type TestDriveService struct {
	testDriveRepo *repositories.TestDriveRepository
	vehicleRepo   *repositories.VehicleRepository
}

func NewTestDriveService(testDriveRepo *repositories.TestDriveRepository, vehicleRepo *repositories.VehicleRepository) *TestDriveService {
	return &TestDriveService{testDriveRepo: testDriveRepo, vehicleRepo: vehicleRepo}
}

func (s *TestDriveService) Create(clientID uint, req models.CreateTestDriveRequest) (*models.TestDrive, error) {
	vehicle, err := s.vehicleRepo.FindByID(req.VehicleID)
	if err != nil {
		return nil, errors.New("vehiculo no encontrado")
	}
	if vehicle.Status != models.VehicleAvailable {
		return nil, errors.New("el vehiculo no esta disponible")
	}

	scheduledAt, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		return nil, errors.New("formato de fecha invalido. Use ISO 8601")
	}

	if scheduledAt.Before(time.Now()) {
		return nil, errors.New("la fecha debe ser futura")
	}

	overlap, err := s.testDriveRepo.HasOverlap(req.VehicleID, scheduledAt)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("el horario seleccionado ya esta ocupado")
	}

	td := &models.TestDrive{
		ClientID:    clientID,
		VehicleID:   req.VehicleID,
		ScheduledAt: scheduledAt,
		Status:      models.TDStatusPending,
		Notes:       req.Notes,
	}

	if err := s.testDriveRepo.Create(td); err != nil {
		return nil, err
	}

	return s.testDriveRepo.FindByID(td.ID)
}

func (s *TestDriveService) GetByID(id uint) (*models.TestDrive, error) {
	return s.testDriveRepo.FindByID(id)
}

func (s *TestDriveService) UpdateStatus(id uint, status models.TestDriveStatus) error {
	td, err := s.testDriveRepo.FindByID(id)
	if err != nil {
		return errors.New("turno no encontrado")
	}
	td.Status = status
	return s.testDriveRepo.Update(td)
}

func (s *TestDriveService) ListByClient(clientID uint) ([]models.TestDrive, error) {
	return s.testDriveRepo.ListByClient(clientID)
}

func (s *TestDriveService) ListAll() ([]models.TestDrive, error) {
	return s.testDriveRepo.ListAll()
}
