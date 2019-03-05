package docker

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/term"
	"github.com/leopardslab/Dunner/internal/logger"
	"github.com/spf13/viper"
)

var log = logger.Log

// Step describes the information required to run one task in docker container
type Step struct {
	Task      string
	Name      string
	Image     string
	Command   []string
	Env       []string
	WorkDir   string
	Volumes   map[string]string
	ExtMounts []mount.Mount
}

// Exec method is used to execute the task described in the corresponding step
func (step Step) Exec() (*io.ReadCloser, error) {

	var (
		hostMountFilepath   = "./"
		containerWorkingDir = "/dunner"
		hostMountTarget     = "/dunner"
	)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion(viper.GetString("DockerAPIVersion")),
	)
	if err != nil {
		log.Fatal(err)
	}

	path, err := filepath.Abs(hostMountFilepath)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("Pulling an image: '%s'", step.Image)

	out, err := cli.ImagePull(ctx, step.Image, types.ImagePullOptions{})
	if err != nil {
		log.Fatal(err)
	}

	termFd, isTerm := term.GetFdInfo(os.Stdout)
	var verbose = viper.GetBool("Verbose")
	if verbose {
		if err = jsonmessage.DisplayJSONMessagesStream(out, os.Stdout, termFd, isTerm, nil); err != nil {
			log.Fatal(err)
		}
	} else {
		if err = jsonmessage.DisplayJSONMessagesStream(out, ioutil.Discard, termFd, isTerm, nil); err != nil {
			log.Fatal(err)
		}
	}

	if err = out.Close(); err != nil {
		log.Fatal(err)
	}
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:      step.Image,
			Cmd:        step.Command,
			Env:        step.Env,
			WorkingDir: containerWorkingDir,
		},
		&container.HostConfig{
			Mounts: append(step.ExtMounts, mount.Mount{
				Type:   mount.TypeBind,
				Source: path,
				Target: hostMountTarget,
			}),
		},
		nil, "")
	if err != nil {
		log.Fatal(err)
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			log.Fatal(err)
		}
	case <-statusCh:
	}

	out, err = cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &out, nil

}
