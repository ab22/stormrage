package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ab22/env"
)

// Config struct that contains all of the configuration variables
// that are set up in the environment.
type Config struct {
	Secret            string `env:"SECRET_KEY" envDefault:"SOME-VERY-SECRET-AND-RANDOM-KEY"`
	Port              int    `env:"PORT" envDefault:"1337"`
	Env               string `env:"ENV" envDefault:"DEV"`
	HostURL           string `env:"HOST_URL" envDefault:"http://localhost:1337/"`
	SessionCookieName string
	SessionLifeTime   time.Duration
}

// NewConfig initializes a new Config structure.
func New() (*Config, error) {
	cfg := &Config{
		SessionCookieName: "__session",
		SessionLifeTime:   time.Minute * 30,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if the most important fields are set and are not empty
// values.
func (c *Config) Validate() error {
	var errorMsg = "config: field [%v] was not set!"

	// Config validation.
	if c.Secret == "" {
		return fmt.Errorf(errorMsg, "Secret")
	}

	return nil
}

// Print configuration values to the log. Some user and password fields
// are omitted for security reasons.
func (c *Config) Print() {
	log.Println("----------------------------------")
	log.Println("-Stormrage Project")
	log.Println("         Host URL:", c.HostURL)
	log.Println(" Application Port:", c.Port)
	log.Println("----------------------------------")
}
