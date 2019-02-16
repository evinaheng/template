package internal

// ServerConfig server module
type ServerConfig struct {
	Port  string
	Env   string
	Debug int
}

// APIConfig API module
type APIConfig struct {
	Endpoint string
}

type EnvironmentConfig struct {
	GlobalTimeout int
}
