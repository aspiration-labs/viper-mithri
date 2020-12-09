/*
Copyright © 2020 aspiration.com

*/
package cmd

import (
	"fmt"
	"github.com/aspiration-labs/viper-mithri/mithri"
	"os"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(os.Stderr, "serve called with %#v, %#v\n", RootAppConfig, ServeAppConfig)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	mithri.AddCommand(rootCmd, serveDefaults, &ServeAppConfig,"serve")
}
