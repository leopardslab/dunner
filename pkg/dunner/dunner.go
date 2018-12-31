package dunner

import (
	"github.com/spf13/cobra"
	"github.com/leopardslab/Dunner/pkg/docker"
	"github.com/leopardslab/Dunner/pkg/config"
)

func Do(cmd *cobra.Command, args []string) {

	var configs = config.GetConfigs()
	for _, stepDefinition := range configs.Tasks[args[0]] {
		step := docker.Step {
			Task: args[0],
			Name: stepDefinition.Name,
			Image: stepDefinition.Image,
			Command: stepDefinition.Command,
		}
		step.Do()
	}
}
