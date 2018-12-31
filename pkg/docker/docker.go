package docker

import (
	"context"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"docker.io/go-docker/api/types/mount"
	log "github.com/sirupsen/logrus"
	"io"
	"path/filepath"
	"strings"
)

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

	_, err = cli.ImagePull(ctx, step.Image, types.ImagePullOptions{})
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
