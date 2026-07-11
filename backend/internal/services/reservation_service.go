package services

import (
	"errors"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/repositories"
)

type ReservationService struct {
	reservationRepo *repositories.ReservationRepository
	vehicleRepo     *repositories.VehicleRepository
}

func NewReservationService(reservationRepo *repositories.ReservationRepository, vehicleRepo *repositories.VehicleRepository) *ReservationService {
	return &ReservationService{reservationRepo: reservationRepo, vehicleRepo: vehicleRepo}
}

func (s *ReservationService) Create(clientID uint, req models.CreateReservationRequest) (*models.Reservation, error) {
	vehicle, err := s.vehicleRepo.FindByID(req.VehicleID)
	if err != nil {
		return nil, errors.New("vehiculo no encontrado")
	}
	if vehicle.Status != models.VehicleAvailable {
		return nil, errors.New("el vehiculo no esta disponible")
	}

	vehicle.Status = models.VehicleReserved
	if err := s.vehicleRepo.Update(vehicle); err != nil {
		return nil, err
	}

	reservation := &models.Reservation{
		ClientID:  clientID,
		VehicleID: req.VehicleID,
		Status:    models.ReservationActive,
		Notes:     req.Notes,
	}

	if err := s.reservationRepo.Create(reservation); err != nil {
		vehicle.Status = models.VehicleAvailable
		s.vehicleRepo.Update(vehicle)
		return nil, err
	}

	return s.reservationRepo.FindByID(reservation.ID)
}

func (s *ReservationService) Confirm(id uint) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("reserva no encontrada")
	}

	reservation.Status = models.ReservationConfirmed
	if err := s.reservationRepo.Update(reservation); err != nil {
		return nil, err
	}

	vehicle, _ := s.vehicleRepo.FindByID(reservation.VehicleID)
	vehicle.Status = models.VehicleSold
	s.vehicleRepo.Update(vehicle)

	return s.reservationRepo.FindByID(id)
}

func (s *ReservationService) Cancel(id uint) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("reserva no encontrada")
	}

	reservation.Status = models.ReservationCancelled
	if err := s.reservationRepo.Update(reservation); err != nil {
		return nil, err
	}

	vehicle, _ := s.vehicleRepo.FindByID(reservation.VehicleID)
	vehicle.Status = models.VehicleAvailable
	s.vehicleRepo.Update(vehicle)

	return s.reservationRepo.FindByID(id)
}

func (s *ReservationService) GetByID(id uint) (*models.Reservation, error) {
	return s.reservationRepo.FindByID(id)
}

func (s *ReservationService) ListByClient(clientID uint) ([]models.Reservation, error) {
	return s.reservationRepo.ListByClient(clientID)
}

func (s *ReservationService) ListAll() ([]models.Reservation, error) {
	return s.reservationRepo.ListAll()
}
