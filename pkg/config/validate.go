package config

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Validate the configuration, fails if there are any errors
func Validate(_ *cobra.Command, args []string) {
	var dunnerFile = viper.GetString("DunnerTaskFile")

	configs, err := GetConfigs(dunnerFile)
	if err != nil {
		log.Fatal(err)
	}

	errs, ok := configs.Validate()
	for _, err := range errs {
		log.Error(err)
	}
	if !ok {
		os.Exit(1)
	}
}
