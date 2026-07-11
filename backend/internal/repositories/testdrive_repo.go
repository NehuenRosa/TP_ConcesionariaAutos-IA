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
	err := r.db.Where("client_id = ?", clientID).Preload("Vehicle").Find(&tds).Error
	return tds, err
}

func (r *TestDriveRepository) ListAll() ([]models.TestDrive, error) {
	var tds []models.TestDrive
	err := r.db.Preload("Client").Preload("Vehicle").Find(&tds).Error
	return tds, err
}

func (r *TestDriveRepository) HasOverlap(vehicleID uint, scheduledAt time.Time) (bool, error) {
	startWindow := scheduledAt.Add(-1 * time.Hour)
	endWindow := scheduledAt.Add(1 * time.Hour)

	var count int64
	err := r.db.Model(&models.TestDrive{}).
		Where("vehicle_id = ? AND status NOT IN ? AND scheduled_at BETWEEN ? AND ?",
			vehicleID, []string{"cancelado", "completado"}, startWindow, endWindow).
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
	err := r.db.Model(&models.TestDrive{}).Where("status = ?", models.TDStatusConfirmed).Count(&count).Error
	return count, err
}
