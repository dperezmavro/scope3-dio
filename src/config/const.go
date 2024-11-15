package config

// environment variable names
const (
	// the environment's name (stage/prod etc)
	envVarEnvName = "ENV"

	// the port where the service listens
	envVarPort = "PORT"

	// the serice's name
	envVarServiceName = "SERVICE"

	// the serice's version
	envVarServiceVersion = "VERSION"

	// scope3 api token
	envVarAPIToken = "SCOPE3_API_TOKEN"
)

// defaults
const (
	// name of default environment
	defaultEnvName = "local"

	// default port
	defaultPort = 3000

	// default service name
	defaultServiceName = "localService"

	// default service version
	defaultServiceVersion = 0
)
