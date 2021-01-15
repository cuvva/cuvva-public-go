package config

import (
	"encoding/json"
)

// Getter takes an environment variable name, and returns its
// value or "" if not set.
type Getter func(string) string

// ConfigEnvironmentVariable is the expected environment variable
// storing JSON configuration data.
const ConfigEnvironmentVariable = "CONFIG"

// FromEnvironment unmarshals JSON configuration from the CONFIG
// environment variable into dst.
func FromEnvironment(get Getter, dst interface{}) error {
	if env := get(ConfigEnvironmentVariable); env != "" {
		err := json.Unmarshal([]byte(env), dst)
		if err != nil {
			return err
		}
	}

	return nil
}

// EnvironmentName returns the name of the current execution environment
// from CONFIG. If no environment is detected, "local" is returned.
func EnvironmentName(get Getter) string {
	cfg := struct {
		Env string `json:"env"`
	}{"local"}

	json.Unmarshal([]byte(get(ConfigEnvironmentVariable)), &cfg)

	return cfg.Env
}
