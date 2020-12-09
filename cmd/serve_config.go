package cmd

type ServeConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

var ServeAppConfig ServeConfig

var serveDefaults = map[string]interface{}{
	"port": 8080,
	"host": "127.0.0.1",
}

