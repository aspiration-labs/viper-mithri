package mithri

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	cmdName string
	viper *viper.Viper
	defaults map[string]interface{}
	appConfig interface{}
	// Cobra PersistentFlags
	inputCfgFile string
	viperConfigType string
}

var useEnv bool
var envPrefix string

// Compile from defaults and optional cfgFile into new viper instance and a config struct.
//func ReadConfig(cfgFile string, defaults map[string]interface{}, config interface{}, useEnv bool, envPrefix string) (*viper.Viper, error) {

var mithriConfigs map[string]*Config

func configDecoder(decoderConfig *mapstructure.DecoderConfig) {
	decoderConfig.ErrorUnused = true
}

func InitConfig(config *Config) (*viper.Viper, error) {
	config.viper = viper.New()
	v := config.viper

	for key, value := range config.defaults {
		v.SetDefault(key, value)
	}

	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
		v.AutomaticEnv()
	} else if useEnv {
		v.AutomaticEnv()
	}

	if config.inputCfgFile != "" {
		v.SetConfigFile(config.inputCfgFile)
	}
	if fileErr := v.ReadInConfig(); fileErr != nil {
		if v.ConfigFileUsed() != "" {
			fmt.Fprintf(os.Stderr, "%s: %v", v.ConfigFileUsed(), fileErr)
		}
	}

	err := v.Unmarshal(config.appConfig, configDecoder)

	return v, err
}
