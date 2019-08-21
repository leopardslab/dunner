package dunner

import (
	"fmt"

	"github.com/leopardslab/dunner/pkg/config"
	"github.com/spf13/viper"
)

// ListTasks lists all the available dunner tasks, if there are errors, it returns `error`
func ListTasks() error {
	var dunnerFile = viper.GetString("DunnerTaskFile")

	configs, err := config.GetConfigs(dunnerFile)
	if err != nil {
		return err
	}

	if len(configs.Tasks) == 0 {
		fmt.Println("No dunner tasks found")
	} else {
		fmt.Println("Available Dunner tasks:")
		for taskName := range configs.Tasks {
			fmt.Println(taskName)
		}
		fmt.Println("Run `dunner do <task_name>` to run a dunner task.")
	}
	return nil
}
