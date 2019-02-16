package gonsq

import nsq "github.com/nsqio/go-nsq"

type gonsq struct {
	prod *nsq.Producer
}

// Config for GoNSQ
type Config struct {
	EndpointNSQLookupD string
	Endpoint           string
}
