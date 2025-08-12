package database

import (
	"gin-ecommerce-example/pkg/config"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

var ESClient *elasticsearch.Client

func ConnectES() error {
	var err error
	ESClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: strings.Split(config.AppConfig.ESAddresses, ","),
	})
	return err
}
