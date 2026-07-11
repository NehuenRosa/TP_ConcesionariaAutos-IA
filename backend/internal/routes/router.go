package routes

import (
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/config"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/middleware"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.CORS(cfg.FrontendURL))

	auth := middleware.AuthMiddleware(cfg.JWTSecret)

	api := r.Group("/api")

	RegisterAuthRoutes(api, db, cfg)
	RegisterVehicleRoutes(api, db, cfg)
	RegisterConsultationRoutes(api, db, cfg, auth)
	RegisterTestDriveRoutes(api, db, cfg, auth)
	RegisterReservationRoutes(api, db, cfg, auth)
	RegisterAdminRoutes(api, db, cfg, auth)
	RegisterChatbotRoutes(api, db, cfg)

	return r
}

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Vehicle{},
		&models.Consultation{},
		&models.ConsultationResponse{},
		&models.TestDrive{},
		&models.Reservation{},
	)
}
