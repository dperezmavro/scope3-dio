package config

type Environment struct {
	// DefaultRequestTimeout int
	Name string
}

type Service struct {
	Name    string
	Version int
}

type Config struct {
	Environment    Environment
	Service        Service
	Port           int
	Scope3APIToken string
	WaitForMissing bool
}
