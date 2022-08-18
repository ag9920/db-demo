package dal

import (
	"fmt"
	"sync"

	"gorm.io/gorm"

	"gorm.io/driver/sqlite"

	"github.com/ag9920/db-demo/gendemo/dal/model"
)

var DB *gorm.DB
var once sync.Once

func init() {
	once.Do(func() {
		DB = ConnectDB().Debug()
		_ = DB.AutoMigrate(&model.User{}, &model.Passport{})
	})
}

func ConnectDB() (conn *gorm.DB) {
	conn, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
	return conn
}
