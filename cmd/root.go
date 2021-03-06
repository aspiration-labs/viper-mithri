/*
Copyright © 2020 aspiration.com

*/
package cmd

import (
	"fmt"
	"github.com/aspiration-labs/viper-mithri/mithri"
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "viper-mithri",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "root called with %#v, %#v\n", RootAppConfig, ServeAppConfig)
	},

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	mithri.AddCommand(rootCmd, rootDefaults, &RootAppConfig, "")
}
