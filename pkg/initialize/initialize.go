package initialize

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/leopardslab/dunner/internal"
	"github.com/leopardslab/dunner/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Initialize command invoked from command line generates a dunner task file with default template
func Initialize(_ *cobra.Command, args []string) {
	var dunnerFile = viper.GetString("DunnerTaskFile")
	if err := initProject(dunnerFile); err != nil {
		logger.Log.Fatalf("Failed to initialize project: %s", err.Error())
	}
	logger.Log.Infof("Dunner task file `%s` created. Please make any required changes.", dunnerFile)
}

func initProject(filename string) error {
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		if err != nil {
			return err
		}
		return fmt.Errorf("%s already exists", filename)
	}
	logger.Log.Infof("Generating %s file", filename)
	return ioutil.WriteFile(filename, []byte(internal.DefaultTaskFileContents), internal.DefaultTaskFilePermission)
}
