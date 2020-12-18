package mithri

import (
	"github.com/spf13/cobra"
)



// Adds a cobra.Command for reading and writing configs.
// Sets PersistentFlags on parentCmd (usually rootCmd) to config options.
func AddCommand(parentCmd *cobra.Command, defaults map[string]interface{}, appConfig interface{}, cmdName string) {
	config := Config{
		cmdName: cmdName,
		defaults: defaults,
		appConfig: appConfig,
	}

	var cmdPrefix string
	if cmdName != "" {
		cmdPrefix = cmdName + "-"
	}
	cobraCommand := cobra.Command{
		Use: cmdPrefix + "config-tool",
		Short: "Read, unmarshal, then write config to file or '-' for stdout",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runWithConfig(cmd, args, &config)
		},
	}
	parentCmd.AddCommand(&cobraCommand)
	parentCmd.PersistentFlags().StringVar(&config.inputCfgFile, cmdPrefix + "config-file", "", "config file to read")
	parentCmd.PersistentFlags().StringVar(&config.viperConfigType, cmdPrefix + "config-type", "yaml", "format of config (default is yaml)")
	if parentCmd.PersistentFlags().Lookup("config-env") == nil {
		parentCmd.PersistentFlags().StringVar(&envPrefix, "config-env", "", "use env with prefix")
	}
	if parentCmd.PersistentFlags().Lookup("use-env") == nil {
		parentCmd.PersistentFlags().BoolVar(&useEnv, "use-env", false, "use env")
	}
	cobra.OnInitialize(func() {
		InitConfig(&config)
	})
}

func runWithConfig(cmd *cobra.Command, args []string, m *Config) {
	v := m.viper
	outputCfgFile := "-"
	if len(args) > 0 {
		outputCfgFile = args[0]
	}

	if outputCfgFile == "-" {
		outputCfgFile = "/dev/stdout"
		v.SetConfigType(m.viperConfigType)
	}
	v.WriteConfigAs(outputCfgFile)
}
