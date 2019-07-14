package initialize

import (
	"fmt"

	"github.com/leopardslab/dunner/internal/util"
	"github.com/leopardslab/dunner/pkg/global"

	yaml "gopkg.in/yaml.v2"
)

var recipeListURL = global.DunnerCookbookListURL

type cookbook struct {
	Recipes []string `yaml:"recipes"`
}

// ListRecipes lists all available dunner recipes from dunner cookbook
func ListRecipes() error {
	contents, err := util.GetURLContents(recipeListURL)
	if err != nil {
		return fmt.Errorf("Failed to download list of recipes from dunner cookbook: %s", err.Error())
	}
	cookbook := cookbook{}
	if err := yaml.Unmarshal(contents, &cookbook); err != nil {
		return fmt.Errorf("Failed to read list of dunner recipes: %s", err.Error())	
	}
	if len(cookbook.Recipes) == 0 {
		fmt.Println("No dunner recipes found")
	} else {
		fmt.Println("Available Dunner recipes are as follows:")
		for _, recipeName := range cookbook.Recipes {
			fmt.Println(recipeName)
		}
		fmt.Println("Run `dunner init <recipe_name>` to initialize a project with dunner recipe.")
	}
	return nil
}
