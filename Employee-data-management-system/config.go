package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Database string `yaml:"database"`
		Charset  string `yaml:"charset"`
	} `yaml:"database"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"JWT"`
}

var DB *gorm.DB
var jwtSecret []byte

func InitConfig() {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal("failed to read config file")
	}

	var config Config
	err = yaml.Unmarshal(f, &config)
	if err != nil {
		log.Fatal("failed to unmarshal config file")
	}
	if config.Database.User == "" ||
		config.Database.Password == "" ||
		config.Database.Host == "" ||
		config.Database.Port == 0 ||
		config.Database.Charset == "" ||
		config.JWT.Secret == "" {
		log.Fatal("database configuration is incomplete")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Database.User, config.Database.Password,
		config.Database.Host, config.Database.Port,
		config.Database.Database, config.Database.Charset)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	jwtSecret = []byte(config.JWT.Secret)

	DB.AutoMigrate(&Employee{}, &User{})
}
