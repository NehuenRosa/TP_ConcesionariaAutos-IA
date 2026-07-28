package repositories

import (
	"time"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"gorm.io/gorm"
)

type TestDriveRepository struct {
	db *gorm.DB
}

func NewTestDriveRepository(db *gorm.DB) *TestDriveRepository {
	return &TestDriveRepository{db: db}
}

func (r *TestDriveRepository) Create(td *models.TestDrive) error {
	return r.db.Create(td).Error
}

func (r *TestDriveRepository) FindByID(id uint) (*models.TestDrive, error) {
	var td models.TestDrive
	err := r.db.Preload("Client").Preload("Vehicle").First(&td, id).Error
	return &td, err
}

func (r *TestDriveRepository) Update(td *models.TestDrive) error {
	return r.db.Save(td).Error
}

func (r *TestDriveRepository) ListByClient(clientID uint) ([]models.TestDrive, error) {
	var tds []models.TestDrive
	err := r.db.Preload("Vehicle").Where("client_id = ?", clientID).Order("scheduled_at desc").Find(&tds).Error
	return tds, err
}

func (r *TestDriveRepository) ListAll() ([]models.TestDrive, error) {
	var tds []models.TestDrive
	err := r.db.Preload("Client").Preload("Vehicle").Order("scheduled_at desc").Find(&tds).Error
	return tds, err
}

func (r *TestDriveRepository) HasOverlap(vehicleID uint, scheduledAt time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.TestDrive{}).
		Where("vehicle_id = ? AND scheduled_at = ? AND status != ?", vehicleID, scheduledAt, models.TDStatusCancelled).
		Count(&count).Error
	return count > 0, err
}

func (r *TestDriveRepository) GetTestDriveCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.TestDrive{}).Count(&count).Error
	return count, err
}

func (r *TestDriveRepository) GetScheduledCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.TestDrive{}).
		Where("status IN ?", []models.TestDriveStatus{models.TDStatusPending, models.TDStatusConfirmed}).
		Count(&count).Error
	return count, err
}
