package seed

import (
	"log"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/models"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/repositories"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) {
	userRepo := repositories.NewUserRepository(db)

	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	sellerHash, _ := bcrypt.GenerateFromPassword([]byte("vendedor123"), bcrypt.DefaultCost)
	clientHash, _ := bcrypt.GenerateFromPassword([]byte("cliente123"), bcrypt.DefaultCost)

	users := []models.User{
		{Name: "Admin", Email: "admin@concesionaria.com", Password: string(adminHash), Role: models.RoleAdmin},
		{Name: "Vendedor Juan", Email: "vendedor@concesionaria.com", Password: string(sellerHash), Role: models.RoleSeller},
		{Name: "Cliente Pedro", Email: "cliente@test.com", Password: string(clientHash), Role: models.RoleClient},
	}
	for _, u := range users {
		existing, _ := userRepo.FindByEmail(u.Email)
		if existing == nil || existing.ID == 0 {
			userRepo.Create(&u)
		}
	}

	vehicleRepo := repositories.NewVehicleRepository(db)
	vehicles := []models.Vehicle{
		{Brand: "Toyota", Model: "Corolla", Year: 2023, Price: 25000000, Mileage: 0, Fuel: models.FuelGasoline, Transmission: models.TransmissionAutomatic, Condition: models.ConditionNew, Color: "Blanco", Description: "Sedan compacto, 4 puertas, full equipo.", Images: pq.StringArray{"https://placehold.co/600x400?text=Toyota+Corolla"}, Status: models.VehicleAvailable, VehicleType: "sedan"},
		{Brand: "Ford", Model: "Ranger", Year: 2022, Price: 35000000, Mileage: 15000, Fuel: models.FuelDiesel, Transmission: models.TransmissionManual, Condition: models.ConditionUsed, Color: "Gris", Description: "Pickup mediana 4x4, ideal para trabajo y familia.", Images: pq.StringArray{"https://placehold.co/600x400?text=Ford+Ranger"}, Status: models.VehicleAvailable, VehicleType: "pickup"},
		{Brand: "Volkswagen", Model: "Golf", Year: 2024, Price: 28000000, Mileage: 0, Fuel: models.FuelGasoline, Transmission: models.TransmissionAutomatic, Condition: models.ConditionNew, Color: "Negro", Description: "Hatchback premium con tecnologia de punta.", Images: pq.StringArray{"https://placehold.co/600x400?text=VW+Golf"}, Status: models.VehicleAvailable, VehicleType: "hatchback"},
		{Brand: "Chevrolet", Model: "Tracker", Year: 2021, Price: 18500000, Mileage: 30000, Fuel: models.FuelGasoline, Transmission: models.TransmissionAutomatic, Condition: models.ConditionUsed, Color: "Rojo", Description: "SUV compacta, excelente relacion precio-calidad.", Images: pq.StringArray{"https://placehold.co/600x400?text=Chevrolet+Tracker"}, Status: models.VehicleAvailable, VehicleType: "suv"},
		{Brand: "Tesla", Model: "Model 3", Year: 2024, Price: 55000000, Mileage: 0, Fuel: models.FuelElectric, Transmission: models.TransmissionAutomatic, Condition: models.ConditionNew, Color: "Azul", Description: "Sedan electrico de alto rendimiento.", Images: pq.StringArray{"https://placehold.co/600x400?text=Tesla+Model3"}, Status: models.VehicleAvailable, VehicleType: "sedan"},
		{Brand: "Toyota", Model: "Hilux", Year: 2023, Price: 40000000, Mileage: 5000, Fuel: models.FuelDiesel, Transmission: models.TransmissionAutomatic, Condition: models.ConditionUsed, Color: "Plata", Description: "Pickup full 4x4 con cubierta y proteccion de caja.", Images: pq.StringArray{"https://placehold.co/600x400?text=Toyota+Hilux"}, Status: models.VehicleAvailable, VehicleType: "pickup"},
	}
	var count int64
	db.Model(&models.Vehicle{}).Count(&count)
	if count == 0 {
		for _, v := range vehicles {
			vehicleRepo.Create(&v)
		}
	}

	log.Println("Seed data loaded successfully")
}
