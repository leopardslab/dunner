package cmd

import (
	"os"

	"github.com/leopardslab/dunner/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the dunner task file `.dunner.yaml`",
	Long:  "You can validate task file `.dunner.yaml` with this command to see if there are any parse errors",
	Run:   Validate,
	Args:  cobra.MinimumNArgs(0),
}

// Validate command invoked from command line, validates the dunner task file. If there are errors, it fails with non-zero exit code.
func Validate(_ *cobra.Command, args []string) {
	var dunnerFile = viper.GetString("DunnerTaskFile")

	configs, err := config.GetConfigs(dunnerFile)
	if err != nil {
		log.Fatal(err)
	}

	errs := configs.Validate()
	if len(errs) != 0 {
		for _, err := range errs {
			log.Error(err)
		}
		os.Exit(1)
	}
}
