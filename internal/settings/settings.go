package settings

import (
	"github.com/spf13/viper"
)

// Init function initializes the default settings for dunner
// These settings can tweaked using appropriate environment variables, or
// defining the configuration in conf present in the appropriate config files
func Init() {
	// Settings file
	viper.SetConfigName("settings")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.dunner/")
	viper.AddConfigPath("/etc/dunner/")
	// TODO Add standard configuration paths for Windows OS

	// Automatic binding of environment variables
	viper.SetEnvPrefix("dunner")
	viper.AutomaticEnv()

	// Files
	viper.SetDefault("DunnerTaskFile", ".dunner.yaml")
	viper.SetDefault("DotenvFile", ".env")
	viper.SetDefault("GlobalLogFile", "/var/log/dunner/logs/")
	viper.SetDefault("LocalLogFile", nil)

	// Modes
	viper.SetDefault("Async", false)
	viper.SetDefault("Verbose", false)

	// Constants
	viper.SetDefault("DockerAPIVersion", "1.39")
}
