package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "mason",
	Short: `mason is a "Proof Of Concept" of client-side builder`,
	Long:  `mason is a "Proof Of Concept" of client-side builder`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mason.yaml)")
	RootCmd.PersistentFlags().BoolP("debug", "d", false, "Set loglevel to DEBUG")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
