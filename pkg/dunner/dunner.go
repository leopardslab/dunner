/*
Package dunner consists of the main executing functions for the Dunner application.
*/
package dunner

import (
	"fmt"
	"os"
	os_user "os/user"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/leopardslab/dunner/internal/logger"
	"github.com/leopardslab/dunner/pkg/config"
	"github.com/leopardslab/dunner/pkg/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var log = logger.Log

// Do method is invoked for command-line use
func Do(_ *cobra.Command, args []string) {
	logger.InitColorOutput()

	var async = viper.GetBool("Async")

	if verbose := viper.GetBool("Verbose"); async && verbose {
		log.Warn("Silencing verbose in asynchronous mode")
		viper.Set("Verbose", false)
	}

	var dunnerFile = viper.GetString("DunnerTaskFile")

	configs, err := config.GetConfigs(dunnerFile)
	if err != nil {
		log.Fatal(err)
	}
	errs := configs.Validate()
	if len(errs) != 0 {
		fmt.Println("Validation failed with following errors:")
		for _, err := range errs {
			logger.ErrorOutput(err.Error())
		}
		os.Exit(1)
	}

	if err = ExecTask(configs, args[0], args[1:], nil); err != nil {
		log.Fatal(err)
	}
}

// ExecTask processes the parsed tasks from the dunner task file
func ExecTask(configs *config.Configs, taskName string, args []string, parentStep *config.Step) error {
	var async = viper.GetBool("Async")
	var wg sync.WaitGroup

	if _, exists := configs.Tasks[taskName]; !exists {
		return fmt.Errorf("dunner: task '%s' does not exist", taskName)
	}
	for _, stepDefinition := range configs.Tasks[taskName].Steps {
		err := stepDefinition.ParseStepEnv()
		if err != nil {
			return err
		}
		if async {
			wg.Add(1)
		}
		step := docker.Step{
			Task:     taskName,
			Name:     stepDefinition.Name,
			Image:    stepDefinition.Image,
			Command:  stepDefinition.Command,
			Commands: stepDefinition.Commands,
			Env:      stepDefinition.Envs,
			WorkDir:  stepDefinition.Dir,
			Follow:   stepDefinition.Follow,
			Args:     stepDefinition.Args,
			User:     getDunnerUser(stepDefinition),
		}

		if err := PassGlobals(&step, configs, &stepDefinition, parentStep); err != nil {
			log.Fatal(err)
		}

		if async {
			go Process(configs, &step, &wg, args, &stepDefinition)
		} else {
			Process(configs, &step, &wg, args, &stepDefinition)
		}
	}

	wg.Wait()
	return nil
}

// Process executes a single step of the task.
func Process(configs *config.Configs, s *docker.Step, wg *sync.WaitGroup, args []string, dunnerStep *config.Step) {
	var async = viper.GetBool("Async")
	if async {
		defer wg.Done()
	}

	if s.Follow != "" {
		if async {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				ExecTask(configs, s.Follow, s.Args, dunnerStep)
				wg.Done()
			}(wg)
		} else {
			ExecTask(configs, s.Follow, s.Args, dunnerStep)
		}
		return
	}

	if err := PassArgs(s, &args); err != nil {
		log.Fatal(err)
	}

	if s.Image == "" {
		log.Fatalf(`dunner: image repository name cannot be empty`)
	}

	err := (*s).Exec()
	if err != nil {
		log.Fatal(err)
	}
}

// PassArgs replaces argument variables,of the form '`$d`', where d is a number, with dth argument.
func PassArgs(s *docker.Step, args *[]string) error {
	var gErr error
	var commands [][]string
	if s.Command != nil {
		commands = [][]string{s.Command}
	} else {
		commands = s.Commands
	}
	for i, cmd := range commands {
		for j, subStr := range cmd {
			regex := regexp.MustCompile(`\$[1-9][0-9]*`)
			subStr = regex.ReplaceAllStringFunc(subStr, func(str string) string {
				j, err := strconv.Atoi(strings.Trim(str, "$"))
				if err != nil {
					log.Fatal(err)
				}
				if j > len(*args) {
					gErr = fmt.Errorf(`dunner: insufficient number of arguments passed`)
					return ""
				}
				return (*args)[j-1]
			})
			if gErr != nil {
				return gErr
			}
			if s.Command != nil {
				s.Command[j] = subStr
			} else {
				s.Commands[i][j] = subStr
			}
		}
	}
	return gErr
}

// getDunnerUser returns the user value from step, if empty returns first found value in order:
// UID env variable, current user ID, current user name.
func getDunnerUser(step config.Step) string {
	if step.User != "" {
		return step.User
	}
	dunnerUser := os.Getenv("UID")
	if dunnerUser == "" {
		user, err := os_user.Current()
		if err != nil {
			dunnerUser = os.Getenv("USER")
			log.Debugf("Unable to find current user id: %s. Using `%s` as Docker user.", err.Error(), dunnerUser)
		} else {
			dunnerUser = user.Uid
		}
	}
	return dunnerUser
}

// PassGlobals uses passes the environment variables and directory mounts that
// are present in the upper scopes in dunner file.
//
// In the case of environment variables, if a different value of variable is given
// in a lower scope as compared to an upper scope, the value from the upper scope
// is overridden by the lower scope variable definition.
// While in the case of directory mounts, similar comparision is done when two mounts
// from different scopes have
// the same destination (target) path.
//
// Since both of these parings are independent of each other, they are carried out
// concurrently on two different goroutines to increase the execution speed.
func PassGlobals(step *docker.Step, configs *config.Configs, stepDefinition *config.Step, parentStep *config.Step) error {
	var wg sync.WaitGroup
	wg.Add(2)

	// Parsing environment variable. Environment variable are overridden if
	// same key is present in the lower scopes.
	go func() {
		envKeys := make(map[string]struct{})
		for _, env := range (*step).Env {
			envKeys[strings.Split(env, "=")[0]] = struct{}{}
		}
		var taskEnvs []string
		if parentStep != nil {
			taskEnvs = append(taskEnvs, parentStep.Envs...)
		}
		taskEnvs = append(taskEnvs, (*configs).Tasks[step.Task].Envs...)
		for _, env := range taskEnvs {
			k := strings.Split(env, "=")[0]
			if _, present := envKeys[k]; !present {
				step.Env = append(step.Env, env)
				envKeys[k] = struct{}{}
			}
		}
		for _, env := range (*configs).Envs {
			k := strings.Split(env, "=")[0]
			if _, present := envKeys[k]; !present {
				step.Env = append(step.Env, env)
			}
		}
		wg.Done()
	}()

	// Parsing of directory mounts. Mounts are overridden if same destination is
	// present in the lower scopes.
	go func() {
		targets := make(map[string]struct{})
		allMounts := (*stepDefinition).Mounts
		for _, mount := range (*stepDefinition).Mounts {
			targets[strings.Split(mount, ":")[1]] = struct{}{}
		}
		var taskMounts []string
		if parentStep != nil {
			taskMounts = append(taskMounts, parentStep.Mounts...)
		}
		taskMounts = append(taskMounts, (*configs).Tasks[step.Task].Mounts...)
		for _, mount := range taskMounts {
			k := strings.Split(mount, ":")[1]
			if _, present := targets[k]; !present {
				allMounts = append(allMounts, mount)
				targets[k] = struct{}{}
			}
		}
		for _, mount := range (*configs).Mounts {
			k := strings.Split(mount, ":")[1]
			if _, present := targets[k]; !present {
				allMounts = append(allMounts, mount)
			}
		}
		if err := config.DecodeMount(allMounts, step); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	wg.Wait()
	return nil
}
