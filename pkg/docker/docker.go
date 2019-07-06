/*
Package docker is the interface of dunner to communicate with the Docker Engine through
methods wrapping over Docker client library.
*/
package docker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/pkg/term"
	"github.com/leopardslab/dunner/internal/logger"
	"github.com/spf13/viper"
)

var log = logger.Log

// Step describes the information required to run one task in docker container. It is very similar to the concept
// of docker build of a 'Dockerfile' and then a sequence of commands to be executed in `docker run`.
type Step struct {
	Task      string            // The name of the task that the step corresponds to
	Name      string            // Name given to this step for identification purpose
	Image     string            // Image is the repo name on which Docker containers are built
	Command   []string          // The command which runs on the container and exits
	Commands  [][]string        // The list of commands that are to be run in sequence
	Env       []string          // The list of environment variables to be exported inside the container
	WorkDir   string            // The primary directory on which task is to be run
	Volumes   map[string]string // Volumes that are to be attached to the container
	ExtMounts []mount.Mount     // The directories to be mounted on the container as bind volumes
	Follow    string            // The next task that must be executed if this does go successfully
	Args      []string          // The list of arguments that are to be passed
}

// Result stores the output of commands run using `docker exec`
type Result struct {
	Command string
	Output  string
	Error   string
}

// Exec method is used to execute the task described in the corresponding step. It returns an object of the
// struct `Result` with the corresponding output and/or error.
//
// Note: A working internet connection is mandatory for the Docker container to contact Docker Hub to find the image and/or
// corresponding updates.
func (step Step) Exec() (*[]Result, error) {

	var (
		hostMountFilepath          = "./"
		containerDefaultWorkingDir = "/dunner"
		hostMountTarget            = "/dunner"
		defaultCommand             = []string{"tail", "-f", "/dev/null"}
		multipleCommands           = false
	)

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	cli.NegotiateAPIVersion(ctx)

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

	var containerWorkingDir = containerDefaultWorkingDir
	if step.WorkDir != "" {
		containerWorkingDir = filepath.Join(hostMountTarget, step.WorkDir)
	}

	multipleCommands = len(step.Commands) > 0
	if !multipleCommands {
		defaultCommand = step.Command
	}
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:      step.Image,
			Cmd:        defaultCommand,
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

	if len(resp.Warnings) > 0 {
		for warning := range resp.Warnings {
			log.Warn(warning)
		}
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	defer func() {
		dur, err := time.ParseDuration("-1ns") // Negative duration means no force termination
		if err != nil {
			log.Fatal(err)
		}
		if err = cli.ContainerStop(ctx, resp.ID, &dur); err != nil {
			log.Fatal(err)
		}
	}()

	var results []Result
	if dryRun := viper.GetBool("Dry-run"); !dryRun {
		if multipleCommands {
			for _, cmd := range step.Commands {
				r, err := runCmd(ctx, cli, resp.ID, cmd)
				if err != nil {
					log.Fatal(err)
				}
				results = append(results, *r)
			}
		} else {
			statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
			select {
			case err = <-errCh:
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

			results = []Result{*ExtractResult(out, step.Command)}
		}
		return &results, nil
	}
	return nil, nil
}

func runCmd(ctx context.Context, cli *client.Client, containerID string, command []string) (*Result, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf(`config: Command cannot be empty`)
	}

	exec, err := cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          command,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := cli.ContainerExecAttach(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Close()

	return ExtractResult(resp.Reader, command), nil
}

// ExtractResult can parse output and/or error corresponding to the command passed as an argument,
// from an io.Reader and convert to an object of strings.
func ExtractResult(reader io.Reader, command []string) *Result {

	var out, errOut bytes.Buffer
	if _, err := stdcopy.StdCopy(&out, &errOut, reader); err != nil {
		log.Fatal(err)
	}
	var result = Result{
		Command: strings.Join(command, " "),
		Output:  out.String(),
		Error:   errOut.String(),
	}
	return &result
}
