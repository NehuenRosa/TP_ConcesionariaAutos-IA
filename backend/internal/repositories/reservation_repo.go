package repositories

import (
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"gorm.io/gorm"
)

type ReservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(res *models.Reservation) error {
	return r.db.Create(res).Error
}

func (r *ReservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var res models.Reservation
	err := r.db.Preload("Client").Preload("Vehicle").First(&res, id).Error
	return &res, err
}

func (r *ReservationRepository) Update(res *models.Reservation) error {
	return r.db.Save(res).Error
}

func (r *ReservationRepository) ListByClient(clientID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Where("client_id = ?", clientID).Preload("Vehicle").Find(&reservations).Error
	return reservations, err
}

func (r *ReservationRepository) ListAll() ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Preload("Client").Preload("Vehicle").Find(&reservations).Error
	return reservations, err
}

func (r *ReservationRepository) GetActiveReservationCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).Where("status = ?", models.ReservationActive).Count(&count).Error
	return count, err
}
