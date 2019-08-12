/*
Package docker is the interface of dunner to communicate with the Docker Engine through
methods wrapping over Docker client library.
*/
package docker

import (
	"bytes"
	"context"
	"flag"
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
	User      string            // User that will run the command(s) inside the container, also support user:group
}

// Result stores the output of commands run using `docker exec`
type Result struct {
	Output string
	Error  string
}

// Exec method is used to execute the task described in the corresponding step. It returns an object of the
// struct `Result` with the corresponding output and/or error.
//
// Note: A working internet connection is mandatory for the Docker container to contact Docker Hub to find the image and/or
// corresponding updates.
func (step Step) Exec() error {
	var (
		hostMountFilepath          = viper.GetString("WorkingDirectory")
		containerDefaultWorkingDir = "/dunner"
		hostMountTarget            = "/dunner"
		defaultCommand             = []string{"tail", "-f", "/dev/null"}
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

	done := make(chan bool)
	go func() {
		ticker := time.Tick(time.Second / 2)
		busyChars := []string{`-`, `\`, `|`, `/`}
		x := 0
	loop:
		for true {
			select {
			case stop := <-done:
				if stop {
					break loop
				}
			default:
				x %= 4
				<-ticker
				if flag.Lookup("test.v") == nil {
					fmt.Printf("\rPulling image: '%s'... %s", step.Image, busyChars[x])
				}
				x++
			}
		}
		fmt.Print("\r")
		log.Infof("Pulled image: '%s'", step.Image)
	}()

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

	done <- true

	if err = out.Close(); err != nil {
		log.Fatal(err)
	}

	var containerWorkingDir = containerDefaultWorkingDir
	if step.WorkDir != "" {
		if step.WorkDir[0] == '/' {
			containerWorkingDir = step.WorkDir
		} else {
			containerWorkingDir = filepath.Join(hostMountTarget, step.WorkDir)
		}
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:      step.Image,
			Cmd:        defaultCommand,
			Env:        step.Env,
			WorkingDir: containerWorkingDir,
			User:       step.User,
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

	commands := step.Commands
	if len(commands) == 0 {
		commands = append(commands, step.Command)
	}

	dryRun := viper.GetBool("Dry-run")
	for _, cmd := range commands {
		done := make(chan bool)
		show := make(chan bool)
		go func() {
			ticker := time.Tick(time.Second / 2)
			busyChars := []string{`-`, `\`, `|`, `/`}
			x := 0
		loop:
			for true {
				select {
				case stop := <-done:
					if stop {
						break loop
					}
				default:
					if flag.Lookup("test.v") == nil {
						x %= 4
						<-ticker
						fmt.Printf("\rRunning command '%s' of '%s' task on a container of '%s' image... %s",
							strings.Join(cmd, " "),
							step.Task,
							step.Image,
							busyChars[x],
						)
						x++
					}
				}
			}
			fmt.Print("\r")
			log.Infof("Finished running command '%+s' on '%+s' docker",
				strings.Join(cmd, " "),
				step.Image,
			)
			show <- true
			return
		}()

		if dryRun {
			continue
		}
		r, err := runCmd(ctx, cli, resp.ID, cmd)
		done <- true
		if err != nil {
			log.Fatal(err)
		}
		if <-show {
			if r != nil && r.Output != "" {
				fmt.Printf(`OUT: %s`, r.Output)
			}
			if r != nil && r.Error != "" {
				fmt.Printf(`ERR: %s`, r.Error)
			}
		}
	}
	return nil
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
		Output: out.String(),
		Error:  errOut.String(),
	}
	return &result
}
