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
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Should get the name of the Dunner file from a config or ENV
		b, err := ioutil.ReadFile("./.dunner.yaml")
		if err != nil {
			log.Fatal(err)
		}

		var cfg map[string][]Config
		if err := yaml.Unmarshal(b, &cfg); err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		cli, err := docker.NewEnvClient()
		if err != nil {
			panic(err)
		}

		for _, step := range cfg[args[0]] {
			_, err = cli.ImagePull(ctx, step.Image, types.ImagePullOptions{})
			if err != nil {
				panic(err)
			}

			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Image: step.Image,
				Cmd:   step.Command,
			}, nil, nil, "")
			if err != nil {
				panic(err)
			}

			if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				panic(err)
			}

			statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
			select {
			case err := <-errCh:
				if err != nil {
					panic(err)
				}
			case <-statusCh:
			}

			out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				panic(err)
			}

			log.Printf("Running task '%+v' on '%+v' Docker with command '%+v'", args[0], step.Image, strings.Join(step.Command, " "))
			stdcopy.StdCopy(os.Stdout, os.Stderr, out)

		}

	},
}
