package main

import (
	"fmt"
	"log"
	"time"

	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/config"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/routes"
	"github.com/NehuenRosa/TP_ConcesionariaAutos-IA/backend/internal/seed"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Waiting for database... attempt %d/30: %v", i+1, err)
		time.Sleep(time.Second)
	}
	if err != nil {
		log.Fatalf("Error connecting to database after retries: %v", err)
	}
	fmt.Println("Connected to database")

	if err := routes.RunMigrations(db); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	seed.Run(db)

	router := routes.Setup(db, cfg)

	log.Printf("Server running on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
