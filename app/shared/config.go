package shared

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

// Config - Application config
type AppConfig struct {
	DB     DBConfig     `json:"db" yaml:"db" mapstructure:"db"`
	Server ServerConfig `json:"server" yaml:"server" mapstructure:"server"`
	Jwk    JwkConfig    `json:"jwk" yaml:"jwk" mapstructure:"jwk"`
	Redis  RedisConfig  `json:"redis" yaml:"redis" mapstructure:"redis"`
}

//DBConfig - Database config
type DBConfig struct {
	Host     string `json:"host" yaml:"host" mapstructure:"host"`
	Port     int    `json:"port" yaml:"port" mapstructure:"port"`
	User     string `json:"user" yaml:"user" mapstructure:"user"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
	DBName   string `json:"dbname" yaml:"dbname" mapstructure:"dbname"`
	SSLMode  string `json:"sslmode" yaml:"sslmode" mapstructure:"sslmode"`
}

type ServerConfig struct {
	Port   string `json:"port" yaml:"port" mapstructure:"port"`
	Domain string `json:"domain" yaml:"domain" mapstructure:"domain"`
}

type JwkConfig struct {
	Keys string `json:"keys"`
}

type RedisConfig struct {
	Address  string `json:"addr" yaml:"addr" mapstructure:"addr"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
}

const IsStatConfigName = "gouthy-config"

// Gets the application configuration directory
func GetAppConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := path.Join(configDir, "gouthy")
	return appConfigDir, nil
}

// GetConfigFilePath - gets a default config file path
func GetConfigFilePath() (string, error) {
	appConfigDir, err := GetAppConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(appConfigDir, IsStatConfigName), nil
}

// Save the config to the specified file
func (config *AppConfig) Save(file string) error {
	content, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, content, 0644); err != nil {
		return err
	}
	return err
}

func (config *AppConfig) Dump() (string, error) {
	content, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// SaveToDefaultLocation - Saves a config to the default location ~/.config/isstat/gouthy-config.yml
func (config *AppConfig) SaveToDefaultLocation() error {
	filePath, err := GetConfigFilePath()
	if err != nil {
		return err
	}
	return config.Save(filePath)
}

// LoadConfig - Loads a config from the configuration
func LoadConfig(cfgFile string) error {
	setDefaults()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find configDir directory.
		appConfigDir, err := GetAppConfigDir()
		if err != nil {
			return err
		}

		viper.AddConfigPath(appConfigDir)
		workingDirectory, err := os.Getwd()
		if err == nil {
			viper.AddConfigPath(workingDirectory)
		}
		viper.AddConfigPath(path.Join(workingDirectory, "config"))
		viper.SetConfigName(IsStatConfigName)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.WithField("file", viper.ConfigFileUsed()).Info("Using config file")
	} else {
		log.WithField("file", viper.ConfigFileUsed()).WithError(err).Debug("Config file not found")
		return nil
	}
	return nil
}

// GetAppConfig - Unmarshal the app configuration using the viper
func GetAppConfig() (AppConfig, error) {
	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		log.WithError(err).WithField("file", viper.ConfigFileUsed()).Error("Unable to parse config")
		return config, err
	}
	return config, nil
}

func setDefaults() {
	viper.SetDefault("server.domain", "localhost")
	viper.SetDefault("server.port", 5000)
	viper.SetDefault("dryrun", false)
}

func SetupLogger(loggingLevel string) {
	if loggingLevel == "" {
		loggingLevel = os.Getenv("LOG_LEVEL")
		if loggingLevel == "" {
			loggingLevel = "warning"
		}
	}

	level, err := log.ParseLevel(loggingLevel)
	if err != nil {
		log.WithError(err).WithField("level", loggingLevel).Warning("Unable to parse the log level")
		level = log.WarnLevel
	}

	log.SetLevel(level)
	log.SetOutput(os.Stderr)
}
