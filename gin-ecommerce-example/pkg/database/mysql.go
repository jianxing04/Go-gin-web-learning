package database

import (
	"gin-ecommerce-example/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectMySQL() error {
	var err error
	DB, err = gorm.Open(mysql.Open(config.AppConfig.MySQLDSN), &gorm.Config{})
	return err
}
