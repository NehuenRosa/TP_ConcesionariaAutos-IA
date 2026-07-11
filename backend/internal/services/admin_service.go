package services

import (
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/repositories"
)

type AdminService struct {
	vehicleRepo      *repositories.VehicleRepository
	consultationRepo *repositories.ConsultationRepository
	testDriveRepo    *repositories.TestDriveRepository
	reservationRepo  *repositories.ReservationRepository
}

func NewAdminService(
	vehicleRepo *repositories.VehicleRepository,
	consultationRepo *repositories.ConsultationRepository,
	testDriveRepo *repositories.TestDriveRepository,
	reservationRepo *repositories.ReservationRepository,
) *AdminService {
	return &AdminService{
		vehicleRepo:      vehicleRepo,
		consultationRepo: consultationRepo,
		testDriveRepo:    testDriveRepo,
		reservationRepo:  reservationRepo,
	}
}

func (s *AdminService) GetDashboard() (map[string]interface{}, error) {
	vehicleStats, err := s.vehicleRepo.GetStats()
	if err != nil {
		return nil, err
	}

	consultationStats, err := s.consultationRepo.GetConsultationCountByPeriod()
	if err != nil {
		return nil, err
	}

	testDriveCount, err := s.testDriveRepo.GetTestDriveCount()
	if err != nil {
		return nil, err
	}

	scheduledTestDrives, err := s.testDriveRepo.GetScheduledCount()
	if err != nil {
		return nil, err
	}

	activeReservations, err := s.reservationRepo.GetActiveReservationCount()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"vehiculos":            vehicleStats,
		"consultas":            consultationStats,
		"test_drives_totales":  testDriveCount,
		"test_drives_pendientes": scheduledTestDrives,
		"reservas_activas":     activeReservations,
	}, nil
}
