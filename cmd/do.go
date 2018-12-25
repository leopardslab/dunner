package cmd

import (
	"context"
	"log"
	"io/ioutil"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"os"
	"github.com/docker/docker/pkg/stdcopy"
	"strings"
)

type Config struct {
	Image   string    `yaml:"image"`
	Command [] string `yaml:"command"`
}

func init() {
	rootCmd.AddCommand(doCmd)
}




var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Do whatever you say",
	Long:  `You can run any task defined on the '.dunner.yaml' with this command`,
	Run: Command(cmd,args),
}
