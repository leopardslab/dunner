package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/leopardslab/dunner/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logger.Log

var rootCmd = &cobra.Command{
	Use:   "dunner",
	Short: "Dunner is a Docker based task-runner",
	Long:  `You can define a set of commands and on what Docker images these commands should run as steps. A task has many steps. Then you can run these tasks with 'dunner do nameoftask'`,
	Run: func(cmd *cobra.Command, args []string) {

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			log.Fatal(err)
		}
		cli.NegotiateAPIVersion(context.Background())

		fmt.Println("Dunner running!")
	},
}

func init() {

	// Verbose Mode
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose mode")
	if err := viper.BindPFlag("Verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		log.Fatal(err)
	}

	// Dunner task file
	rootCmd.PersistentFlags().StringP("task-file", "t", ".dunner.yaml", "Task file to be run")
	if err := rootCmd.MarkPersistentFlagFilename("task-file", "yaml", "yml"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("DunnerTaskFile", rootCmd.PersistentFlags().Lookup("task-file")); err != nil {
		log.Fatal(err)
	}

	// Environment file
	rootCmd.PersistentFlags().StringP("env-file", "e", ".env", "Environment file")
	if err := rootCmd.MarkPersistentFlagFilename("env-file", "env"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("DotenvFile", rootCmd.PersistentFlags().Lookup("env-file")); err != nil {
		log.Fatal(err)
	}

	// Working directory
	rootCmd.PersistentFlags().StringP("context", "C", "./", "Working directory")
	if err := rootCmd.MarkPersistentFlagDirname("env-file"); err != nil {
		log.Fatal(err)
	}
	if err := viper.BindPFlag("WorkingDirectory", rootCmd.PersistentFlags().Lookup("context")); err != nil {
		log.Fatal(err)
	}

}

// Execute method executes the 'Run' method of rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
