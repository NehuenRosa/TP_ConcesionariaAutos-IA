package routes

import (
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/config"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/handlers"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/middleware"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/repositories"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	userRepo := repositories.NewUserRepository(db)
	authSvc := services.NewAuthService(userRepo, cfg)
	handler := handlers.NewAuthHandler(authSvc)
	auth := middleware.AuthMiddleware(cfg.JWTSecret)

	api.POST("/auth/register", handler.Register)
	api.POST("/auth/login", handler.Login)
	api.GET("/auth/me", auth, handler.Me)
}

func RegisterVehicleRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	vehicleRepo := repositories.NewVehicleRepository(db)
	vehicleSvc := services.NewVehicleService(vehicleRepo)
	handler := handlers.NewVehicleHandler(vehicleSvc)
	auth := middleware.AuthMiddleware(cfg.JWTSecret)
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin))

	api.GET("/vehicles", handler.List)
	api.GET("/vehicles/brands", handler.GetBrands)
	api.GET("/vehicles/:id", handler.GetByID)
	api.POST("/vehicles", auth, adminOnly, handler.Create)
	api.PUT("/vehicles/:id", auth, adminOnly, handler.Update)
	api.DELETE("/vehicles/:id", auth, adminOnly, handler.Delete)
}

func RegisterConsultationRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config, auth gin.HandlerFunc) {
	consultationRepo := repositories.NewConsultationRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	consultationSvc := services.NewConsultationService(consultationRepo, vehicleRepo)
	handler := handlers.NewConsultationHandler(consultationSvc)
	sellerOnly := middleware.RoleMiddleware(string(models.RoleSeller), string(models.RoleAdmin))

	api.POST("/consultations", auth, handler.Create)
	api.GET("/consultations/mine", auth, handler.ListMy)
	api.GET("/consultations", auth, sellerOnly, handler.ListAll)
	api.GET("/consultations/pending/count", auth, sellerOnly, handler.GetPendingCount)
	api.GET("/consultations/notifications/count", auth, handler.GetNotificationCounts)
	api.GET("/consultations/:id", auth, handler.GetByID)
	api.PATCH("/consultations/:id/status", auth, sellerOnly, handler.UpdateStatus)
	api.POST("/consultations/:id/responses", auth, handler.AddResponse)
	api.DELETE("/consultations/:id", auth, handler.Delete)
}

func RegisterTestDriveRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config, auth gin.HandlerFunc) {
	testDriveRepo := repositories.NewTestDriveRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	testDriveSvc := services.NewTestDriveService(testDriveRepo, vehicleRepo)
	handler := handlers.NewTestDriveHandler(testDriveSvc)
	sellerOnly := middleware.RoleMiddleware(string(models.RoleSeller), string(models.RoleAdmin))

	api.POST("/test-drives", auth, handler.Create)
	api.GET("/test-drives/mine", auth, handler.ListMy)
	api.GET("/test-drives", auth, sellerOnly, handler.ListAll)
	api.GET("/test-drives/:id", auth, handler.GetByID)
	api.PATCH("/test-drives/:id/status", auth, sellerOnly, handler.UpdateStatus)
}

func RegisterReservationRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config, auth gin.HandlerFunc) {
	reservationRepo := repositories.NewReservationRepository(db)
	vehicleRepo := repositories.NewVehicleRepository(db)
	reservationSvc := services.NewReservationService(reservationRepo, vehicleRepo)
	handler := handlers.NewReservationHandler(reservationSvc)
	sellerOnly := middleware.RoleMiddleware(string(models.RoleSeller), string(models.RoleAdmin))

	api.POST("/reservations", auth, handler.Create)
	api.GET("/reservations/mine", auth, handler.ListMy)
	api.GET("/reservations", auth, sellerOnly, handler.ListAll)
	api.GET("/reservations/:id", auth, handler.GetByID)
	api.POST("/reservations/:id/confirm", auth, sellerOnly, handler.Confirm)
	api.POST("/reservations/:id/cancel", auth, sellerOnly, handler.Cancel)
}

func RegisterAdminRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config, auth gin.HandlerFunc) {
	vehicleRepo := repositories.NewVehicleRepository(db)
	consultationRepo := repositories.NewConsultationRepository(db)
	testDriveRepo := repositories.NewTestDriveRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)
	adminSvc := services.NewAdminService(vehicleRepo, consultationRepo, testDriveRepo, reservationRepo)
	handler := handlers.NewAdminHandler(adminSvc)
	adminOnly := middleware.RoleMiddleware(string(models.RoleAdmin))

	api.GET("/admin/dashboard", auth, adminOnly, handler.Dashboard)
}

func RegisterChatbotRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	vehicleRepo := repositories.NewVehicleRepository(db)
	chatbotSvc := services.NewChatbotService(vehicleRepo, cfg)
	handler := handlers.NewChatbotHandler(chatbotSvc)

	api.GET("/chatbot/status", handler.Status)
	api.POST("/chatbot/ask", handler.Ask)
}
