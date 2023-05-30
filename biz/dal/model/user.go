package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string `gorm:"type:varchar(255)"`
	Password   string `gorm:"type:varchar(255)"`
	PhoneNum   string `gorm:"type:varchar(255)"`
	Permission byte
}
