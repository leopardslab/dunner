package cmd

import (
	"github.com/leopardslab/dunner/internal/logger"
	"github.com/leopardslab/dunner/pkg/initialize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Generates a dunner task file `.dunner.yaml`",
	Long:    "You can initialize any project with dunner task file. It generates a default task file `.dunner.yaml`, you can customize it based on needs. You can override the name of task file using -t flag.",
	Run:     Initialize,
	Args:    cobra.NoArgs,
	Aliases: []string{"i"},
}

// Initialize command invoked from command line generates a dunner task file with default template
func Initialize(_ *cobra.Command, args []string) {
	var dunnerFile = viper.GetString("DunnerTaskFile")
	if err := initialize.InitProject(dunnerFile); err != nil {
		logger.Log.Fatalf("Failed to initialize project: %s", err.Error())
	}
	logger.Log.Infof("Dunner task file `%s` created. Please make any required changes.", dunnerFile)
}
