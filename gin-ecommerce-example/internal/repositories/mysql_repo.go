package repositories

import (
	"gin-ecommerce-example/internal/models"
	"gin-ecommerce-example/pkg/database"
)

func CreateProduct(product *models.Product) error {
	return database.DB.Create(product).Error
}

func GetProduct(id uint) (*models.Product, error) {
	var product models.Product
	err := database.DB.First(&product, id).Error
	return &product, err
}

func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := database.DB.Find(&products).Error
	return products, err
}

func UpdateProduct(product *models.Product) error {
	return database.DB.Save(product).Error
}

func DeleteProduct(id uint) error {
	return database.DB.Delete(&models.Product{}, id).Error
}
