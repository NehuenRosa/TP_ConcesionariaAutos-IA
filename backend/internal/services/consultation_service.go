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

func (s *ConsultationService) GetPendingCount() (int64, error) {
	return s.consultationRepo.CountPending()
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

	if consultation.Status == models.ConsultInProgress && consultation.AssignedTo != nil {
		if userID == consultation.ClientID {
			consultation.HasUnreadMessages = true
			s.consultationRepo.Update(consultation)
		} else if userID == *consultation.AssignedTo {
			consultation.HasUnreadForClient = true
			s.consultationRepo.Update(consultation)
		}
	}

	return s.consultationRepo.FindByID(consultationID)
}

func (s *ConsultationService) Delete(id, userID uint, role string) error {
	consultation, err := s.consultationRepo.FindByID(id)
	if err != nil {
		return errors.New("consulta no encontrada")
	}

	if role != string(models.RoleAdmin) && role != string(models.RoleSeller) && consultation.ClientID != userID {
		return errors.New("no tienes permiso para eliminar esta consulta")
	}

	if err := s.consultationRepo.DeleteResponseByConsultation(id); err != nil {
		return err
	}

	return s.consultationRepo.Delete(id)
}

func (s *ConsultationService) MarkAsRead(id uint) error {
	return s.consultationRepo.MarkAsRead(id)
}

func (s *ConsultationService) MarkAsReadForClient(id uint) error {
	return s.consultationRepo.MarkAsReadForClient(id)
}

func (s *ConsultationService) GetNotificationCounts(role string, userID uint) (map[string]int64, error) {
	result := make(map[string]int64)
	if role == string(models.RoleSeller) || role == string(models.RoleAdmin) {
		pending, err := s.consultationRepo.CountPending()
		if err != nil {
			return nil, err
		}
		unread, err := s.consultationRepo.CountUnreadForSeller()
		if err != nil {
			return nil, err
		}
		result["pending"] = pending
		result["unread"] = unread
		result["total"] = pending + unread
	} else {
		unread, err := s.consultationRepo.CountUnreadForClient(userID)
		if err != nil {
			return nil, err
		}
		result["unread"] = unread
		result["total"] = unread
	}
	return result, nil
}
