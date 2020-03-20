package config

import (
	"github.com/spf13/viper"
	"os"
	"sync"
)

// How to differentiate configuration file per environment?
// Ans: (read the following reference)
// https://github.com/spf13/viper/issues/80#issuecomment-350049720
// https://medium.com/@felipedutratine/manage-config-in-golang-to-get-variables-from-file-and-env-variables-33d876887152

var (
	conf *Config
	once sync.Once
)

const (
	ProductionEnv  = "PRODUCTION"
	StagingEnv     = "STAGING"
	DevelopmentEnv = "DEVELOPMENT"
)

const (
	AppEnvironmentKey = "APP_ENV"
)

func New() (c *Config, err error) {
	once.Do(func() {
		conf, err = loadConfig()
	})
	if err != nil {
		return
	}
	return conf, nil
}

func loadConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == ProductionEnv {
		viper.SetConfigName("config")
	} else if env == StagingEnv {
		viper.SetConfigName("staging-config")
	} else {
		viper.SetConfigName("devel-config")
	}
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/.pocker")
	// In docker, the executable is located in /pocker/. See the Dockerfile.
	viper.AddConfigPath("/pocker")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	c := new(Config)
	if err := viper.Unmarshal(c); err != nil {
		return nil, err
	}
	return c, nil
}

// GetConfig will return a copy of global setted config instance.
// It will return nil when the config haven't yet initialized
// by calling the New method.
func GetConfig() *Config {
	if conf == nil {
		return conf
	}
	return &(*conf)
}

// Config is the configuration setting
// use within the pocker.
type Config struct {
	Debug    bool `json:"debug"`
	Database struct {
		Host string `json:"host"`
		Port string `json:"port"`
		User string `json:"user"`
		Pass string `json:"pass"`
		Name string `json:"name"`
	} `json:"database"`
	Elasticsearch struct {
		Servers []string
	}
	Server struct {
		Address string `json:"address"`
	} `json:"server"`
	JWTToken string `json:"jwttoken"`
}
