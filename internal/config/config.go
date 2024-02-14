package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordCommandPrefix string
	DiscordBotToken      string
	RestEnabled          bool
	RestGinRelease       bool
	RestHostname         string
}

// NewConfig creates a new Config object and returns it along with any error encountered.
//
// No parameters.
// Returns *Config and error.
func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if err := validateMandatoryConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		DiscordCommandPrefix: os.Getenv("DISCORD_COMMAND_PREFIX"),
		DiscordBotToken:      os.Getenv("DISCORD_BOT_TOKEN"),
		RestEnabled:          getenvAsBool("REST_ENABLED"),
		RestGinRelease:       getenvAsBool("REST_GIN_RELEASE"),
		RestHostname:         os.Getenv("REST_HOSTNAME"),
	}

	return config, nil
}

// String returns the JSON representation of the Config struct.
//
// No parameters.
// Returns a string.
func (c *Config) String() string {
	configMap := map[string]interface{}{
		"DiscordCommandPrefix": c.DiscordCommandPrefix,
		"DiscordBotToken":      c.DiscordBotToken,
		"RestEnabled":          c.RestEnabled,
		"RestGinRelease":       c.RestGinRelease,
		"RestHostname":         c.RestHostname,
	}

	jsonString, err := json.MarshalIndent(configMap, "", "    ")
	if err != nil {
		return ""
	}

	return string(jsonString)
}

// validateMandatoryConfig checks for the presence of mandatory configuration keys in the environment variables and returns an error if any are missing.
//
// No parameters.
// Returns an error.
func validateMandatoryConfig() error {
	mandatoryKeys := []string{
		"DISCORD_COMMAND_PREFIX", "DISCORD_BOT_TOKEN", "REST_ENABLED",
	}

	for _, key := range mandatoryKeys {
		if os.Getenv(key) == "" {
			return errors.New("Missing mandatory configuration: " + key)
		}
	}

	return nil
}

// getenvAsBool returns the boolean value of the environment variable specified by the key.
//
// It takes a string key as a parameter and returns a boolean value.
func getenvAsBool(key string) bool {
	val := os.Getenv(key)

	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}

	return boolValue
}
