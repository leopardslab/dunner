package cmd

import (
	"github.com/spf13/cobra"
	"github.com/leopardslab/Dunner/src/services/DunnerService"
)

type Config struct {
	Image   string    `yaml:"image"`
	Command [] string `yaml:"command"`
}

func init() {
	rootCmd.AddCommand(doCmd)
}

var s DunnerService.Service

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Do whatever you say",
	Long:  `You can run any task defined on the '.dunner.yaml' with this command`,
	Run: s.Do,
}
