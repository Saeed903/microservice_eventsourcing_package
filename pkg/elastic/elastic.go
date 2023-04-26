package elastic

import (
	"net/http"
	"os"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Addresses []string `mapstructure:"addresses" validate:"required"`
	Username  string   `mapstructure:"username"`
	Password  string   `mapstructure:"password"`

	APIKey        string      `mapstructre:"apiKey"`
	Header        http.Header // Global HTTP request header.
	EnableLogging bool        `mapstructure:"enableLogging"`
}

func NewElasticSearchClient(cfg Config) (*elasticsearch.Client, error) {
	config := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		APIKey:    cfg.APIKey,
		Header:    cfg.Header,
	}

	if cfg.EnableLogging {
		config.Logger = &elastictransport.ColorLogger{Output: os.Stdout, EnableRequestBody: true, EnableResponseBody: true}
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
