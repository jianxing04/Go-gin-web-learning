package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gin-ecommerce-example/internal/models"
	"gin-ecommerce-example/pkg/database"

	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const indexName = "products"

func IndexProduct(product *models.Product) error {
	body, _ := json.Marshal(product)
	req := esapi.IndexRequest{
		Index:      indexName,
		DocumentID: fmt.Sprintf("%d", product.ID),
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), database.ESClient)
	if err != nil || res.IsError() {
		return fmt.Errorf("index error: %v", err)
	}
	res.Body.Close()
	return nil
}

func SearchProducts(query string) ([]models.Product, error) {
	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name", "description"},
			},
		},
	}
	json.NewEncoder(&buf).Encode(searchQuery)

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  &buf,
	}
	res, err := req.Do(context.Background(), database.ESClient)
	if err != nil || res.IsError() {
		return nil, fmt.Errorf("search error: %v", err)
	}
	defer res.Body.Close()

	var result struct {
		Hits struct {
			Hits []struct {
				Source models.Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	json.NewDecoder(res.Body).Decode(&result)

	var products []models.Product
	for _, hit := range result.Hits.Hits {
		products = append(products, hit.Source)
	}
	return products, nil
}

func UpdateProductInES(product *models.Product) error {
	body, _ := json.Marshal(product)
	req := esapi.UpdateRequest{
		Index:      indexName,
		DocumentID: fmt.Sprintf("%d", product.ID),
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, body))),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), database.ESClient)
	if err != nil || res.IsError() {
		return fmt.Errorf("update error: %v", err)
	}
	res.Body.Close()
	return nil
}

func DeleteProductFromES(id uint) error {
	req := esapi.DeleteRequest{
		Index:      indexName,
		DocumentID: fmt.Sprintf("%d", id),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), database.ESClient)
	if err != nil || res.IsError() {
		return fmt.Errorf("delete error: %v", err)
	}
	res.Body.Close()
	return nil
}
