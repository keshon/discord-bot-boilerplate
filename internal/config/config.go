package config

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"github.com/gookit/slog"
	"github.com/joho/godotenv"
)

type Config struct {
	DiscordCommandPrefix string
	DiscordBotToken      string
	RestEnabled          bool
	RestGinRelease       bool
	RestHostname         string
}

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

func (c *Config) String() string {
	// Create a map for key-value pairs
	configMap := map[string]interface{}{
		"DiscordCommandPrefix": c.DiscordCommandPrefix,
		"DiscordBotToken":      c.DiscordBotToken,
		"RestEnabled":          c.RestEnabled,
		"RestGinRelease":       c.RestGinRelease,
		"RestHostname":         c.RestHostname,
	}

	// Convert the map to a JSON string
	jsonString, err := json.MarshalIndent(configMap, "", "    ")
	if err != nil {
		slog.Error("Error formatting configuration as JSON")
		return ""
	}

	return string(jsonString)
}

func validateMandatoryConfig() error {

	// Define a list of mandatory environment variable keys
	// Extra overlay to ensure we have all necessary values even if they have default values (e.g. ffmpeg has own)
	// ignore:
	// - REST_GIN_RELEASE
	// - REST_HOSTNAME

	mandatoryKeys := []string{
		"DISCORD_COMMAND_PREFIX", "DISCORD_BOT_TOKEN", "REST_ENABLED",
	}

	// Check if any mandatory key is missing
	for _, key := range mandatoryKeys {
		if os.Getenv(key) == "" {
			return errors.New("Missing mandatory configuration: " + key)
		}
	}

	return nil
}

func getenvAsInt(key string) int {
	val := os.Getenv(key)

	intValue, err := strconv.Atoi(val)
	if err != nil {
		slog.Error("Error parsing integer value from env variable")
		return 0
	}

	return intValue
}

func getenvAsBool(key string) bool {
	val := os.Getenv(key)

	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		slog.Error("Error parsing bool value from env variable")
		return false
	}

	return boolValue
}

func getenvBoolAsInt(key string) int {
	val := os.Getenv(key)

	boolValue, err := strconv.ParseBool(val)
	if err != nil {
		slog.Error("Error parsing bool value from env variable")
		return 0
	}

	if boolValue {
		return 1
	}

	return 0
}
