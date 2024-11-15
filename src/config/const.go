package config

// environment variable names
const (
	// the environment's name (stage/prod etc)
	envVarEnvName = "ENV"

	// the port where the service listens
	envVarPort = "PORT"

	// the serice's name
	envVarServiceName = "SERVICE"

	// scope3 api token
	envVarApiToken = "SCOPE3_API_TOKEN"
)

// defaults
const (
	// name of default environment
	defaultEnvName = "local"

	// default port
	defaultPort = 3000

	// default service name
	defaultServiceName = "localService"

	// default scope3 token
	defaultApiToken = "unusable_token"
)
