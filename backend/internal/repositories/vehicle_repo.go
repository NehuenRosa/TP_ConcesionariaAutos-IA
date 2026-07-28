package repositories

import (
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"gorm.io/gorm"
)

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) Create(vehicle *models.Vehicle) error {
	return r.db.Create(vehicle).Error
}

func (r *VehicleRepository) FindByID(id uint) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	err := r.db.First(&vehicle, id).Error
	return &vehicle, err
}

func (r *VehicleRepository) Update(vehicle *models.Vehicle) error {
	return r.db.Save(vehicle).Error
}

func (r *VehicleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Vehicle{}, id).Error
}

func (r *VehicleRepository) List(filter models.VehicleFilter) ([]models.Vehicle, int64, error) {
	query := r.db.Model(&models.Vehicle{})

	if filter.Search != "" {
		like := "%" + filter.Search + "%"
		query = query.Where("brand ILIKE ? OR model ILIKE ? OR description ILIKE ?", like, like, like)
	}
	if filter.Brand != "" {
		query = query.Where("brand = ?", filter.Brand)
	}
	if filter.Model != "" {
		query = query.Where("model = ?", filter.Model)
	}
	if filter.YearFrom != nil {
		query = query.Where("year >= ?", *filter.YearFrom)
	}
	if filter.YearTo != nil {
		query = query.Where("year <= ?", *filter.YearTo)
	}
	if filter.PriceFrom != nil {
		query = query.Where("price >= ?", *filter.PriceFrom)
	}
	if filter.PriceTo != nil {
		query = query.Where("price <= ?", *filter.PriceTo)
	}
	if filter.Fuel != "" {
		query = query.Where("fuel = ?", filter.Fuel)
	}
	if filter.Condition != "" {
		query = query.Where("condition = ?", filter.Condition)
	}
	if filter.VehicleType != "" {
		query = query.Where("vehicle_type = ?", filter.VehicleType)
	}

	query = query.Where("status = ?", models.VehicleAvailable)

	var total int64
	query.Count(&total)

	sortBy := "created_at"
	sortOrder := "desc"
	if filter.SortBy == "price" || filter.SortBy == "year" {
		sortBy = filter.SortBy
	}
	if filter.SortOrder == "asc" || filter.SortOrder == "desc" {
		sortOrder = filter.SortOrder
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 12
	}
	offset := (filter.Page - 1) * filter.PageSize

	var vehicles []models.Vehicle
	err := query.Order(sortBy + " " + sortOrder).Offset(offset).Limit(filter.PageSize).Find(&vehicles).Error
	return vehicles, total, err
}

func (r *VehicleRepository) ListAll() ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	err := r.db.Find(&vehicles).Error
	return vehicles, err
}

func (r *VehicleRepository) GetDistinctBrands() ([]string, error) {
	var brands []string
	err := r.db.Model(&models.Vehicle{}).Where("status = ?", models.VehicleAvailable).Distinct().Pluck("brand", &brands).Error
	return brands, err
}

func (r *VehicleRepository) GetStats() (map[string]interface{}, error) {
	var total, available, reserved, sold int64
	r.db.Model(&models.Vehicle{}).Count(&total)
	r.db.Model(&models.Vehicle{}).Where("status = ?", models.VehicleAvailable).Count(&available)
	r.db.Model(&models.Vehicle{}).Where("status = ?", models.VehicleReserved).Count(&reserved)
	r.db.Model(&models.Vehicle{}).Where("status = ?", models.VehicleSold).Count(&sold)

	return map[string]interface{}{
		"total":      total,
		"disponible": available,
		"reservado":  reserved,
		"vendido":    sold,
	}, nil
}
