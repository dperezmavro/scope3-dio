package config

import "fmt"

// returns a new Config object, or error
func New() (*Config, error) {
	envName := fromEnvStringDefault(envVarEnvName, defaultEnvName)
	serviceName := fromEnvStringDefault(envVarServiceName, defaultServiceName)
	waitForMissing := fromEnvBoolDefault(envVarWaitForMissing, defaultWaitForMissing)

	apiToken, err := fromEnvString(envVarAPIToken)
	if err != nil {
		return nil, fmt.Errorf("error populating config: %+v", err)
	}

	port, err := fromEnvIntDefault(envVarPort, defaultPort)
	if err != nil {
		return nil, fmt.Errorf("error populating config: %+v", err)
	}

	version, err := fromEnvIntDefault(envVarServiceVersion, defaultServiceVersion)
	if err != nil {
		return nil, fmt.Errorf("error populating config: %+v", err)
	}

	c := &Config{
		Service: Service{
			Name:    serviceName,
			Version: version,
		},
		Environment: Environment{
			Name: envName,
		},
		Port:           port,
		Scope3APIToken: apiToken,
		WaitForMissing: waitForMissing,
	}

	return c, nil
}
