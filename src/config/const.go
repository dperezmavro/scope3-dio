package config

// environment variable names
const (
	// the environment's name (stage/prod etc)
	envVarEnvName = "ENV"

	// the port where the service listens
	envVarPort = "PORT"

	// the serice's name
	envVarServiceName = "SERVICE"
)

// defaults
const (
	// name of default environment
	defaultEnvName = "local"

	// default port
	defaultPort = 3000

	// default service name
	defaultServiceName = "localService"
)
