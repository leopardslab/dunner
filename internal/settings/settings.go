package settings

import (
	"github.com/spf13/viper"
)

// Initialize sets default configuration for the project
func Initialize() {
	// Settings file
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("$HOME/.dunner/")
	viper.AddConfigPath("/etc/dunner/")
	// TODO Add standard configuration paths for Windows OS

	// Automatic binding of environment variables
	viper.SetEnvPrefix("dunner")
	viper.AutomaticEnv()

	// Files
	viper.SetDefault("DunnerTaskFile", ".dunner.yaml")
	viper.SetDefault("GlobalLogFile", "/var/log/dunner/logs/")
	viper.SetDefault("LocalLogFile", nil)

	// Modes
	viper.SetDefault("Async", false)
	viper.SetDefault("Verbose", false)
}
