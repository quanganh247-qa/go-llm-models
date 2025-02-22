package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(databaseURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// // Auto-migrate the schemas
	// err = db.AutoMigrate(
	// 	&models.Patient{},
	// 	&models.Allergy{},
	// 	// &models.Note{},
	// 	// &models.Breed{},
	// )
	// if err != nil {
	// 	log.Fatalf("Failed to migrate database: %v", err)
	// }

	fmt.Println("Database connected and migrated successfully")
	return db
}
