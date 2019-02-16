package v5

import (
	eslib "gopkg.in/olivere/elastic.v5"
)

// ElasticV5 module
type elasticV5 struct {
	config Config
	client *eslib.Client
}

// Config for elastic module
type Config struct {
	Endpoint string
}
