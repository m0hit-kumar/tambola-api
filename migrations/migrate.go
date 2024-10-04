package migrations

import (
	"github.com/m0hit-kumar/tambola/models"
	"gorm.io/gorm"
)

var tables = []interface{}{&models.Users{}, &models.Books{}, &models.TicketDesign{}}

func MigrateTables(db *gorm.DB) error {
	err := db.AutoMigrate(tables...)
 	return err
}
