package services

import (
	"errors"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/repositories"
)

type ConsultationService struct {
	consultationRepo *repositories.ConsultationRepository
	vehicleRepo      *repositories.VehicleRepository
}

func NewConsultationService(consultationRepo *repositories.ConsultationRepository, vehicleRepo *repositories.VehicleRepository) *ConsultationService {
	return &ConsultationService{consultationRepo: consultationRepo, vehicleRepo: vehicleRepo}
}

func (s *ConsultationService) Create(clientID uint, req models.CreateConsultationRequest) (*models.Consultation, error) {
	vehicle, err := s.vehicleRepo.FindByID(req.VehicleID)
	if err != nil {
		return nil, errors.New("vehiculo no encontrado")
	}
	if vehicle.Status != models.VehicleAvailable {
		return nil, errors.New("el vehiculo no esta disponible")
	}

	consultation := &models.Consultation{
		ClientID:  clientID,
		VehicleID: req.VehicleID,
		Message:   req.Message,
		Status:    models.ConsultPending,
	}

	if err := s.consultationRepo.Create(consultation); err != nil {
		return nil, err
	}

	return s.consultationRepo.FindByID(consultation.ID)
}

func (s *ConsultationService) GetByID(id uint) (*models.Consultation, error) {
	return s.consultationRepo.FindByID(id)
}

func (s *ConsultationService) UpdateStatus(id uint, status models.ConsultationStatus, sellerID uint) error {
	consultation, err := s.consultationRepo.FindByID(id)
	if err != nil {
		return errors.New("consulta no encontrada")
	}

	consultation.Status = status
	consultation.AssignedTo = &sellerID
	return s.consultationRepo.Update(consultation)
}

func (s *ConsultationService) ListByClient(clientID uint) ([]models.Consultation, error) {
	return s.consultationRepo.ListByClient(clientID)
}

func (s *ConsultationService) ListAll() ([]models.Consultation, error) {
	return s.consultationRepo.ListAll()
}

func (s *ConsultationService) AddResponse(consultationID, userID uint, message string) (*models.Consultation, error) {
	consultation, err := s.consultationRepo.FindByID(consultationID)
	if err != nil {
		return nil, errors.New("consulta no encontrada")
	}

	response := &models.ConsultationResponse{
		ConsultationID: consultationID,
		UserID:         userID,
		Message:        message,
	}

	if err := s.consultationRepo.CreateResponse(response); err != nil {
		return nil, err
	}

	if consultation.Status == models.ConsultPending {
		consultation.Status = models.ConsultInProgress
		if err := s.consultationRepo.Update(consultation); err != nil {
			return nil, err
		}
	}

	return s.consultationRepo.FindByID(consultationID)
}
