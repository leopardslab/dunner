package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dunner",
	Long:  `All software has versions. This is Dunners's`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("v0.1")
	},
}
