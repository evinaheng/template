package internal

// DatabaseConfig database module
type DatabaseConfig struct {
	Driver string
	Master string
	Slave  []string
}

// RedisConfig redis module
type RedisConfig struct {
	Endpoint string
	Timeout  int
	MaxIdle  int
}

// ElasticConfig elastic module
type ElasticConfig struct {
	Endpoint string
}
