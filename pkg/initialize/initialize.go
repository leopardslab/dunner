package initialize

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/leopardslab/dunner/internal"
	"github.com/leopardslab/dunner/internal/logger"
	"github.com/leopardslab/dunner/internal/util"
	"github.com/leopardslab/dunner/pkg/global"

	yaml "gopkg.in/yaml.v2"
)

type recipeMetadata struct {
	Name               string `yaml:"name"`
	Description        string `yaml:"description"`
	Version            string `yaml:"version"`
	PreInstallCmd      string `yaml:"preInstallCmd"`
	PostInstallMessage string `yaml:"postInstallMessage"`
}

// InitProject generates a dunner task file with default template
func InitProject(filename string, args []string) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		return fmt.Errorf("%s already exists", filename)
	}
	if len(args) == 1 && args[0] != "" {
		return InitWithRecipe(filename, args[0])
	}
	logger.Log.Infof("Generating %s file", filename)
	return ioutil.WriteFile(filename, []byte(internal.DefaultTaskFileContents), internal.DefaultTaskFilePermission)
}

// InitWithRecipe initializes the project with given dunner recipe, returns an error if invalid
func InitWithRecipe(filename string, templateName string) error {
	metadataURL := getMetadataURL(templateName)

	fmt.Printf("Downloading metadata of dunner recipe %s...\n", templateName)
	metadata, err := getRecipeMetadata(metadataURL)
	if err != nil {
		return fmt.Errorf("Failed to initialize project with %s recipe. %s", templateName, err.Error())
	}

	if metadata.PreInstallCmd != "" {
		command := strings.Fields(metadata.PreInstallCmd)
		cmd, err := util.ExecuteSystemCommand(command, os.Stdout, os.Stderr)
		if err != nil {
			return fmt.Errorf("Failed to execute pre-install command from recipe: %s", err.Error())
		}
		if err = cmd.Wait(); err != nil {
			return fmt.Errorf("Pre-install command execution of recipe failed")
		}
	}

	fmt.Println("Downloading dunner task file of recipe...")
	dunnerRecipeURL := getDunnerTaskURLOfRecipe(templateName)
	dErr := util.Download(dunnerRecipeURL, filename)
	if dErr != nil {
		return fmt.Errorf("Failed to initialize project with %s recipe. %s", templateName, dErr.Error())
	}
	if metadata.PostInstallMessage != "" {
		fmt.Println(metadata.PostInstallMessage)
	}
	return nil
}

var getMetadataURL = func(templateName string) string {
	return fmt.Sprintf("%s%s/metadata.yml", global.DunnerCookbookRecipesURL, templateName)
}

var getDunnerTaskURLOfRecipe = func(templateName string) string {
	return fmt.Sprintf("%s%s/.dunner.yaml", global.DunnerCookbookRecipesURL, templateName)
}

func getRecipeMetadata(metadataURL string) (*recipeMetadata, error) {
	metadataContents, err := util.GetURLContents(metadataURL)
	if err != nil {
		return nil, fmt.Errorf("Failed to download metadata: %s", err.Error())
	}
	metadata := recipeMetadata{}
	if err := yaml.Unmarshal(metadataContents, &metadata); err != nil {
		logger.Log.Debugf("Failed to unmarshal metadata: %s", err.Error())
		return nil, fmt.Errorf("Failed to unmarshal metadata")
	}
	return &metadata, nil
}
