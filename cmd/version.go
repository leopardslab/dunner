package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	G "github.com/leopardslab/Dunner/pkg/global"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dunner",
	Long:  `All software has versions. This is Dunners's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(G.VERSION)
	},
}
