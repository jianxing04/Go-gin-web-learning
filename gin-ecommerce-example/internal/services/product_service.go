package services

import (
	"errors"
	"gin-ecommerce-example/internal/models"
	"gin-ecommerce-example/internal/repositories"
)

func CreateProduct(product *models.Product) error {
	if err := repositories.CreateProduct(product); err != nil {
		return err
	}
	repositories.IndexProduct(product)
	repositories.InvalidateCache()
	return nil
}

func GetProduct(id uint) (*models.Product, error) {
	return repositories.GetProduct(id)
}

func GetAllProducts() ([]models.Product, error) {
	products, err := repositories.GetCachedProducts()
	if err == nil {
		return products, nil
	}
	products, err = repositories.GetAllProducts()
	if err != nil {
		return nil, err
	}
	repositories.SetCachedProducts(products)
	return products, nil
}

func SearchProducts(query string) ([]models.Product, error) {
	if query == "" {
		return nil, errors.New("query required")
	}
	return repositories.SearchProducts(query)
}

func UpdateProduct(product *models.Product) error {
	if err := repositories.UpdateProduct(product); err != nil {
		return err
	}
	repositories.UpdateProductInES(product)
	repositories.InvalidateCache()
	return nil
}

func DeleteProduct(id uint) error {
	if err := repositories.DeleteProduct(id); err != nil {
		return err
	}
	repositories.DeleteProductFromES(id)
	repositories.InvalidateCache()
	return nil
}
