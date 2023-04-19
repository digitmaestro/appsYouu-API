package prisma

import (
	"context"
	"log"
	"strings"

	"github.com/caarlos0/env"
)

type Config struct {
	Host   string `env:"PRISMA_HOST" envDefault:"localhost:4466"`
	Stage  string `env:"APP_NAMESPACE" envDefault:"keskin-dev"`
	Secret string `env:"PRISMA_SECRET"`
}

func NewConfig() (*Options, error) {
	config := &Config{}
	err := env.Parse(config)
	if err != nil {
		return nil, err
	}

	return &Options{
		Endpoint: createEndpoint(config),
		Secret:   config.Secret,
	}, nil
}

func NewClient(config *Options) (*Client, error) {
	client := New(config)

	// run a simple query at application start to fail if unsuccessful
	// language=GraphQL
	_, err := client.Client.GraphQL(context.Background(), `query { __schema { queryType { name } } }`, map[string]interface{}{})

	if err != nil {
		return nil, err
	}

	return client, nil
}

func createEndpoint(config *Config) string {
	log.Printf("namespace %s", config.Stage)
	split := strings.Split(config.Stage, "-")
	// 'keskin-dev' or 'dev' is allowed
	last := split[len(split)-1]
	log.Printf("using stage %s", last)
	return "http://" + config.Host + "/keskin/" + last
}
