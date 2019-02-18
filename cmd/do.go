package cmd

import (
	"github.com/leopardslab/Dunner/pkg/dunner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(doCmd)

	// Async Mode
	doCmd.Flags().BoolP("async", "A", false, "Async mode")
	if err := viper.BindPFlag("Async", doCmd.Flags().Lookup("async")); err != nil {
		log.Fatal(err)
	}

}

var doCmd = &cobra.Command{
	Use:   "do [taskName]",
	Short: "Do whatever you say",
	Long:  `You can run any task defined on the '.dunner.yaml' with this command`,
	Run:   dunner.Do,
	Args:  cobra.ExactArgs(1),
}
