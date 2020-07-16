package config

import (
	"encoding/json"

	"github.com/crazy-max/diun/v4/internal/model"
	"github.com/crazy-max/gonfig"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Config holds configuration details
type Config struct {
	Db        *model.Db        `yaml:"db,omitempty" json:"db,omitempty"`
	Watch     *model.Watch     `yaml:"watch,omitempty" json:"watch,omitempty"`
	Notif     *model.Notif     `yaml:"notif,omitempty" json:"notif,omitempty"`
	RegOpts   model.RegOpts    `yaml:"regopts,omitempty" json:"regopts,omitempty" validate:"unique=Name,dive"`
	Providers *model.Providers `yaml:"providers,omitempty" json:"providers,omitempty" validate:"required"`
}

// Load returns Config struct
func Load(cfgfile string) (*Config, error) {
	cfg := Config{
		Db:    (&model.Db{}).GetDefaults(),
		Watch: (&model.Watch{}).GetDefaults(),
	}

	fileLoader := gonfig.NewFileLoader(gonfig.FileLoaderConfig{
		Filename: cfgfile,
	})
	if found, err := fileLoader.Load(&cfg); err != nil {
		return nil, errors.Wrap(err, "Failed to decode configuration from file")
	} else if !found {
		log.Debug().Msg("No configuration file found")
	} else {
		log.Info().Msgf("Configuration loaded from file: %s", fileLoader.GetFilename())
	}

	envLoader := gonfig.NewEnvLoader(gonfig.EnvLoaderConfig{
		Prefix: "DIUN_",
	})
	if found, err := envLoader.Load(&cfg); err != nil {
		return nil, errors.Wrap(err, "Failed to decode configuration from environment variables")
	} else if !found {
		log.Debug().Msg("No DIUN_* environment variables defined")
	} else {
		log.Info().Msgf("Configuration loaded from %d environment variables", len(envLoader.GetVars()))
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// String returns the string representation of configuration
func (cfg *Config) String() string {
	b, _ := json.MarshalIndent(cfg, "", "  ")
	return string(b)
}
