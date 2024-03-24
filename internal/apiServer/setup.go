package apiServer

import (
	"github.com/damonto/estkme-rlpa-server/internal/rlpa"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func connectDB() {
	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&User{})
	if err != nil {
		return
	}

	DB = database
}

var RlpaServerManager rlpa.Manager

func initRlpaServerManager() {
	RlpaServerManager = rlpa.NewManager()
}

func setupServer() {
	connectDB()
	initRlpaServerManager()
}
