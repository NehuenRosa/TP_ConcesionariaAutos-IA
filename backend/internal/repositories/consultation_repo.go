package repositories

import (
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"gorm.io/gorm"
)

type ConsultationRepository struct {
	db *gorm.DB
}

func NewConsultationRepository(db *gorm.DB) *ConsultationRepository {
	return &ConsultationRepository{db: db}
}

func (r *ConsultationRepository) Create(c *models.Consultation) error {
	return r.db.Create(c).Error
}

func (r *ConsultationRepository) FindByID(id uint) (*models.Consultation, error) {
	var c models.Consultation
	err := r.db.Preload("Client").Preload("Vehicle").Preload("Responses.User").First(&c, id).Error
	return &c, err
}

func (r *ConsultationRepository) Update(c *models.Consultation) error {
	return r.db.Save(c).Error
}

func (r *ConsultationRepository) ListByClient(clientID uint) ([]models.Consultation, error) {
	var consultations []models.Consultation
	err := r.db.Where("client_id = ?", clientID).Preload("Vehicle").Preload("Responses.User").Find(&consultations).Error
	return consultations, err
}

func (r *ConsultationRepository) ListAll() ([]models.Consultation, error) {
	var consultations []models.Consultation
	err := r.db.Preload("Client").Preload("Vehicle").Preload("Responses.User").Find(&consultations).Error
	return consultations, err
}

func (r *ConsultationRepository) ListByStatus(status models.ConsultationStatus) ([]models.Consultation, error) {
	var consultations []models.Consultation
	err := r.db.Where("status = ?", status).Preload("Client").Preload("Vehicle").Preload("Responses.User").Find(&consultations).Error
	return consultations, err
}

func (r *ConsultationRepository) CreateResponse(response *models.ConsultationResponse) error {
	return r.db.Create(response).Error
}

func (r *ConsultationRepository) GetConsultationCountByPeriod() (map[string]int64, error) {
	var total int64
	r.db.Model(&models.Consultation{}).Count(&total)
	return map[string]int64{"total": total}, nil
}
