package dummyelastic

// Dummy elastic struct
type dummy struct {
	config Config
}

// Config for dummy elastic
type Config struct {
	Mocking map[string]string
}
