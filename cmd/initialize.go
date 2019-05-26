package cmd

import (
	"github.com/leopardslab/dunner/pkg/initialize"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Generates a dunner task file `.dunner.yaml`",
	Long:    "You can initialize any project with dunner task file. It generates a default task file `.dunner.yaml`, you can customize it based on needs. You can override the name of task file using -t flag.",
	Run:     initialize.Initialize,
	Args:    cobra.NoArgs,
	Aliases: []string{"i"},
}
