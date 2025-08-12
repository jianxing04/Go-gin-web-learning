package repositories

import (
	"encoding/json"
	"gin-ecommerce-example/internal/models"
	"gin-ecommerce-example/pkg/database"
	"time"
)

const cacheKey = "products"

func GetCachedProducts() ([]models.Product, error) {
	val, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}
	var products []models.Product
	json.Unmarshal([]byte(val), &products)
	return products, nil
}

func SetCachedProducts(products []models.Product) error {
	data, _ := json.Marshal(products)
	return database.RedisClient.Set(database.Ctx, cacheKey, data, 10*time.Minute).Err()
}

func InvalidateCache() {
	database.RedisClient.Del(database.Ctx, cacheKey)
}
