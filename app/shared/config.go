package shared

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const AppName = "gouthy"

var disabledFeature = FeatureConfig{Enabled: false}

type FeaturesConfig map[string]FeatureConfig

type FeatureParams map[string]interface{}

// Config - Application config
type AppConfig struct {
	DB       DBConfig       `json:"db" yaml:"db" mapstructure:"db"`
	Server   ServerConfig   `json:"server" yaml:"server" mapstructure:"server"`
	Jwk      JwkConfig      `json:"jwk" yaml:"jwk" mapstructure:"jwk"`
	Redis    RedisConfig    `json:"redis" yaml:"redis" mapstructure:"redis"`
	Features FeaturesConfig `json:"features" yaml:"features" mapstructure:"features"`
}

//DBConfig - Database config
type DBConfig struct {
	Default     string                   `json:"default" yaml:"default" mapstructure:"default"`
	DataSources map[string]DBEntryConfig `json:"datasources" yaml:"datasources" mapstructure:"datasources"`
	AutoMigrate bool                     `json:"automigrate" yaml:"automigrate" mapstructure:"automigrate"`
}

func (c *DBConfig) GetDefault() *DBEntryConfig {
	entry, ok := c.DataSources[c.Default]
	if ok {
		return &entry
	}
	return &DBEntryConfig{
		Uri:    "file:memdb1?mode=memory&cache=shared",
		DBType: "sqlite",
	}
}

type DBEntryConfig struct {
	Uri         string   `json:"uri" yaml:"uri" mapstructure:"uri"`
	DBType      string   `json:"db_type" yaml:"db_type" mapstructure:"db_type"`
	AutoMigrate bool     `json:"automigrate" yaml:"automigrate" mapstructure:"automigrate"`
	DataImport  []string `json:"data_import" yaml:"data_import" mapstructure:"data_import"`
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

type FeatureConfig struct {
	Enabled bool          `json:"enabled" yaml:"enabled" mapstructure:"enabled"`
	Params  FeatureParams `json:"params" yaml:"params" mapstructure:"params"`
}

const IsStatConfigName = "gouthy-config"

// Gets the application configuration directory
func getAppConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appConfigDir := path.Join(configDir, AppName)
	return appConfigDir, nil
}

// getConfigFilePath - gets a default config file path
func getConfigFilePath() (string, error) {
	appConfigDir, err := getAppConfigDir()
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

// SaveToDefaultLocation - Saves a config to the default location ~/.config/gouthy/gouthy-config.yml
func (config *AppConfig) SaveToDefaultLocation() error {
	filePath, err := getConfigFilePath()
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
		appConfigDir, err := getAppConfigDir()
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

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.MergeInConfig(); err == nil {
		log.WithField("file", viper.ConfigFileUsed()).Info("Using config file")
	} else {
		log.WithField("file", viper.ConfigFileUsed()).WithError(err).Debug("Config file not found")
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
	viper.SetDefault("database.inmemory", false)
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

// Feature config

func (f *FeatureConfig) GetParamString(name string, d string) string {
	str, ok := f.Params[name]
	if !ok {
		return d
	}
	return str.(string)
}

func (f *FeatureConfig) GetParamInt(name string, d int) int {
	str, ok := f.Params[name]
	if !ok {
		return d
	}
	return str.(int)
}

func (f *FeatureConfig) GetParamsBool(name string, d bool) bool {
	str, ok := f.Params[name]
	if !ok {
		return d
	}
	return str.(bool)
}

func (c *FeaturesConfig) PasswordPolicy() *FeatureConfig {
	return c.GetFeature("password_policy")
}

func (c *FeaturesConfig) GetFeature(name string) *FeatureConfig {
	if feature, ok := (*c)[name]; ok {
		return &feature
	}

	return &disabledFeature
}
