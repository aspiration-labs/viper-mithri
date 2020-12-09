package cmd

type RootConfig struct {
	ApiUrl string `mapstructure:"api_url"`
	Hostname string `mapstructure:"hostname"`
	Auth RootAuthConfig `mapstructure:",squash"`
}

type RootAuthConfig struct {
	Username string `mapstructure:"auth_username"`
	Password string `mapstructure:"auth_password"`
}

var RootAppConfig RootConfig

var rootDefaults = map[string]interface{}{
	"api_url": "http://localhost/api",
	"auth_username": "zzyzx",
	"auth_password": "12fa",
}

