package cmd

import (
	"github.com/leopardslab/dunner/internal/logger"
	"github.com/leopardslab/dunner/pkg/initialize"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listRecipesCmd)
}

var listRecipesCmd = &cobra.Command{
	Use:   "list-recipes",
	Short: "Lists all available recipes from dunner cookbook",
	Long:  "This lists all the available recipes from Dunner cookbook with which you can initialize any project as `dunner init <recipe_name>`.",
	Run:   ListRecipes,
	Args:  cobra.NoArgs,
}

// ListRecipes command invoked from command line lists all available dunner recipes in cookbook
func ListRecipes(_ *cobra.Command, args []string) {
	if err := initialize.ListRecipes(); err != nil {
		logger.Log.Fatalf("Failed to list dunner recipes: %s", err.Error())
	}
}
