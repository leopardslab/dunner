package cmd

import (
	"fmt"
	"os"

	"github.com/leopardslab/dunner/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:     "validate",
	Short:   "Validate the dunner task file `.dunner.yaml`",
	Long:    "You can validate task file `.dunner.yaml` with this command to see if there are any parse errors",
	Run:     Validate,
	Args:    cobra.NoArgs,
	Aliases: []string{"v"},
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
		fmt.Println("Validation failed with following errors:")
		for _, err := range errs {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}
	fmt.Println("Validation successful!")
}
