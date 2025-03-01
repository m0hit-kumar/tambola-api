package models

import "gorm.io/gorm"

type TicketDesign struct {
	gorm.Model
	HostName   string `json:"hostName" gorm:"not null"`
	Background string `json:"background" gorm:"not null"`
	Border     string `json:"border" gorm:"not null"`
	Text       string `json:"text" gorm:"not null"`
	RoomId   string `json:"roomId" gorm:"not null;unique"`
 }
