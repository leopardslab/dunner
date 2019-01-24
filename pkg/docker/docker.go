package docker

import (
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/leopardslab/Dunner/internal/logger"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var log = logger.Log

type Step struct {
	Task    string
	Name    string
	Image   string
	Command [] string
	Env     map[string]string
	WorkDir string
	Volumes map[string]string
}

func (step Step) Do() (*io.ReadCloser, error) {

	var (
		hostMountFilepath   = "./"
		containerWorkingDir = "/dunner"
		hostMountTarget     = "/dunner"
	)

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	path, err := filepath.Abs(hostMountFilepath)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:      step.Image,
			Cmd:        step.Command,
			WorkingDir: containerWorkingDir,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: path,
					Target: hostMountTarget,
				},
			},
		},
		nil, "")
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Running task '%+v' on '%+v' Docker with command '%+v'", step.Task, step.Image, strings.Join(step.Command, " "))
	return &out, nil

}

func PullImages(images *[]string) error {

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, image := range *images {
		wg.Add(1)

		go func(image string) {

			defer wg.Done()
			log.Infof("Pulling image '%s'", image)

			out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
			if err != nil {
				log.Fatal(err)
			}

			termFd, isTerm := term.GetFdInfo(os.Stdout)
			if err = jsonmessage.DisplayJSONMessagesStream(out, ioutil.Discard, termFd, isTerm, nil); err != nil {
				log.Fatal(err)
			}

			if err = out.Close(); err != nil {
				log.Fatal(err)
			}
		}(image)
	}
	wg.Wait()
	log.Info("Pull complete\n\n")

	return nil
}
