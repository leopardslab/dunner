package cmd

import (
	"github.com/leopardslab/dunner/internal/logger"
	"github.com/leopardslab/dunner/pkg/dunner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listTasksCmd)
}

var listTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Lists all available tasks in dunner task file",
	Long:  "This lists all the available tasks in dunner task file, `.dunner.yaml` file by default or file passed to `-t` flag",
	Run:   ListTasks,
	Args:  cobra.NoArgs,
}

// ListTasks command invoked from command line lists all available dunner tasks
func ListTasks(_ *cobra.Command, args []string) {
	if err := dunner.ListTasks(); err != nil {
		logger.Log.Fatalf("Failed to list dunner tasks: %s", err.Error())
	}
}
