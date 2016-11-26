package config

import (
	"fmt"
	"log"
	"path"
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
	FrontendAppPath   string
	SessionCookieName string
	SessionLifeTime   time.Duration

	DB struct {
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     int    `env:"DB_PORT" envDefault:"5432"`
		User     string `env:"DB_USER" envDefault:"postgres"`
		Password string `env:"DB_PASS" envDefault:"1234"`
		Name     string `env:"DB_NAME" envDefault:"abemar"`
		LogMode  bool   `env:"DB_LOG_MODE" envDefault:"False"`
	}

	PrivateRouter struct {
		Address  string `env:"PR_ADDR"`
		Port     string `env:"PR_PORT" envDefault:"8728"`
		User     string `env:"PR_USER"`
		Password string `env:"PR_PASS"`
	}
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

	var frontendFolder string
	if cfg.Env == "DEV" {
		frontendFolder = "app"
	} else {
		frontendFolder = "dist"
	}

	cfg.FrontendAppPath = path.Join("frontend/abemar-mikrotik", frontendFolder)

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

	// DB validation.
	if c.DB.Host == "" {
		return fmt.Errorf(errorMsg, "DB.Host")
	}

	if c.DB.Port == 0 {
		return fmt.Errorf(errorMsg, "DB.Port")
	}

	if c.DB.User == "" {
		return fmt.Errorf(errorMsg, "DB.User")
	}

	if c.DB.Password == "" {
		return fmt.Errorf(errorMsg, "DB.Password")
	}

	if c.DB.Name == "" {
		return fmt.Errorf(errorMsg, "DB.Name")
	}

	// Private Router config validation.
	if c.PrivateRouter.Address == "" {
		return fmt.Errorf(errorMsg, "PrivateRouter.Address")
	}

	if c.PrivateRouter.Port == "" {
		return fmt.Errorf(errorMsg, "PrivateRouter.Port")
	}

	if c.PrivateRouter.User == "" {
		return fmt.Errorf(errorMsg, "PrivateRouter.User")
	}

	if c.PrivateRouter.Password == "" {
		return fmt.Errorf(errorMsg, "PrivateRouter.Password")
	}

	return nil
}

// Print configuration values to the log. Some user and password fields
// are omitted for security reasons.
func (c *Config) Print() {
	log.Println("----------------------------------")
	log.Println("-Stormrage Project")
	log.Println("            Host URL:", c.HostURL)
	log.Println("    Application Port:", c.Port)
	log.Println("       Database Host:", c.DB.Host)
	log.Println("       Database Port:", c.DB.Port)
	log.Println("       Database Name:", c.DB.Name)
	log.Println("         Db Log mode:", c.DB.LogMode)
	log.Println("       Frontend path:", c.FrontendAppPath)
	log.Println(" Private Router Addr:", c.PrivateRouter.Address)
	log.Println(" Private Router Port:", c.PrivateRouter.Port)
	log.Println("----------------------------------")
}
