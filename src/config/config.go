package config

import "fmt"

var Default = Config{
	Environment: Environment{
		Name: defaultEnvName,
	},
	Port:        defaultPort,
	ServiceName: defaultServiceName,
}

// returns a new Config object, or error
func New() (*Config, error) {

	envName, err := fromEnvString(envVarEnvName)
	if err != nil {
		return nil, fmt.Errorf("error populating config: %+v", err)
	}

	serviceName, err := fromEnvString(envVarServiceName)
	if err != nil {
		return nil, fmt.Errorf("error populating config: %+v", err)
	}

	apiToken, err := fromEnvString(envVarApiToken)
	if err != nil {
		return nil, fmt.Errorf("error populating config: %+v", err)
	}

	port, err := fromEnvInt(envVarPort)
	if err != nil {
		return nil, fmt.Errorf("error populating config: %+v", err)
	}

	c := &Config{
		ServiceName: serviceName,
		Environment: Environment{
			Name: envName,
		},
		Port:           port,
		Scope3APIToken: apiToken,
	}

	return c, nil
}

type Environment struct {
	// DefaultRequestTimeout int
	Name string
}

type Config struct {
	Environment    Environment
	ServiceName    string
	Port           int
	Scope3APIToken string
}
