package DockerService

import (
	"os"
	log "github.com/sirupsen/logrus"
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
		log.Fatal(err)
	}

	_, err = cli.ImagePull(ctx, step.Image, types.ImagePullOptions{})
	if err != nil {
		log.Fatal(err)
	}

	path, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
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
		ShowStdout: true, ShowStderr: true,})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Running task '%+v' on '%+v' Docker with command '%+v'", step.Task, step.Image, strings.Join(step.Command, " "))
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

}
