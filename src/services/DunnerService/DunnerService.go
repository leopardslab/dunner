package DunnerService

import (
	"github.com/spf13/cobra"
	"github.com/leopardslab/Dunner/src/services/DockerService"
	"github.com/leopardslab/Dunner/src/services/ConfigService"
)

func Do(cmd *cobra.Command, args []string) {

	var configs = ConfigService.GetConfigs()
	for _, stepDefinition := range configs.Tasks[args[0]] {
		step := DockerService.Step {
			Task: args[0],
			Name: stepDefinition.Name,
			Image: stepDefinition.Image,
			Command: stepDefinition.Command,
		}
		step.Do()
	}
}
