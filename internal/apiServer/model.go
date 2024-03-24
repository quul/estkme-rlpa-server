package apiServer

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string
	PasswordHash string
}
