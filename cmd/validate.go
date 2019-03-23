package cmd

import (
	"github.com/leopardslab/Dunner/pkg/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the dunner task file `.dunner.yaml`",
	Long:  "You can validate task file `.dunner.yaml` with this command to see if there are any parse errors",
	Run:   config.Validate,
	Args:  cobra.MinimumNArgs(0),
}
