package models

import "gorm.io/gorm"

var tables = []interface{}{&Users{}, &Books{}}

func MigrateTables(db *gorm.DB) error {
	err := db.AutoMigrate(tables...)
	return err
}
