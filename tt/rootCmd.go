package tt

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os"
)

var debug bool
var cfgFile string
var flagAll bool
var paramClient string
var paramProject string
var paramTask string

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "time tracking on comand line",
	Long:  `time tracking on comand line`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(ErrorString, CharError, err)
	}
}

func init() {

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/tt/tt.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&debug, FlagDebug, "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag(FlagDebug, rootCmd.PersistentFlags().Lookup(FlagDebug))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath("$XDG_CONFIG_HOME/tt")
		viper.AddConfigPath(home + "/.config/tt")
		viper.SetConfigType("yaml")
		viper.SetConfigName("tt")
	}

	viper.SetEnvPrefix("tt")
	viper.BindEnv("db")
	viper.GetString("db")
	if err := viper.ReadInConfig(); err != nil {
		// Set default values for parameters
		viper.Set("debug", false)
	}

	if viper.GetBool("debug") {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		fmt.Fprintln(os.Stderr, "Using Database file:", viper.GetString("db"))
	}

	initDB()
}
