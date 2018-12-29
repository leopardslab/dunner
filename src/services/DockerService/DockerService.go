package DockerService

import (
	"os"
	"log"
	"context"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker"
	"docker.io/go-docker/api/types/container"
	"strings"
	"github.com/docker/docker/pkg/stdcopy"
	"docker.io/go-docker/api/types/mount"
	"path/filepath"
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

func (step Step) Do() {

	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ImagePull(ctx, step.Image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	path, err := filepath.Abs("./")
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: step.Image,
		Cmd:   step.Command,
		WorkingDir: "/dunner",
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: path,
				Target: "/dunner",
			},
		},
	}, nil, "")
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

	log.Printf("Running task '%+v' on '%+v' Docker with command '%+v'", step.Task, step.Image, strings.Join(step.Command, " "))
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

}
