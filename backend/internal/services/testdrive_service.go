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
	scheduledAt, err := time.Parse("2006-01-02T15:04", req.ScheduledAt)
	if err != nil {
		return nil, errors.New("formato de fecha invalido (yyyy-mm-ddThh:mm)")
	}

	if scheduledAt.Before(time.Now()) {
		return nil, errors.New("la fecha debe ser futura")
	}

	vehicle, err := s.vehicleRepo.FindByID(req.VehicleID)
	if err != nil {
		return nil, errors.New("vehiculo no encontrado")
	}
	if vehicle.Status != models.VehicleAvailable {
		return nil, errors.New("el vehiculo no esta disponible")
	}

	overlap, err := s.testDriveRepo.HasOverlap(req.VehicleID, scheduledAt)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("ya existe un turno para ese horario")
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

	validTransitions := map[models.TestDriveStatus][]models.TestDriveStatus{
		models.TDStatusPending:   {models.TDStatusConfirmed, models.TDStatusCancelled},
		models.TDStatusConfirmed: {models.TDStatusCompleted, models.TDStatusCancelled},
		models.TDStatusCancelled: {},
		models.TDStatusCompleted: {},
	}

	allowed, ok := validTransitions[td.Status]
	if !ok {
		return errors.New("estado invalido")
	}

	valid := false
	for _, s := range allowed {
		if s == status {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("transicion de estado no permitida")
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
