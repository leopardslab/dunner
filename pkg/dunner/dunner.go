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
			fmt.Println(err.Error())
		}
		os.Exit(1)
	}

	ExecTask(configs, args[0], args[1:])
}

// ExecTask processes the parsed tasks from the dunner task file
func ExecTask(configs *config.Configs, taskName string, args []string) {
	var async = viper.GetBool("Async")
	var wg sync.WaitGroup
	for _, stepDefinition := range configs.Tasks[taskName] {
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
			WorkDir:  stepDefinition.SubDir,
			Follow:   stepDefinition.Follow,
			Args:     stepDefinition.Args,
			User:     getDunnerUser(stepDefinition),
		}

		if err := config.DecodeMount(stepDefinition.Mounts, &step); err != nil {
			log.Fatal(err)
		}

		if async {
			go Process(configs, &step, &wg, args)
		} else {
			Process(configs, &step, &wg, args)
		}
	}

	wg.Wait()
}

// Process executes a single step of the task.
func Process(configs *config.Configs, s *docker.Step, wg *sync.WaitGroup, args []string) {
	var async = viper.GetBool("Async")
	if async {
		defer wg.Done()
	}

	if s.Follow != "" {
		if async {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				ExecTask(configs, s.Follow, s.Args)
				wg.Done()
			}(wg)
		} else {
			ExecTask(configs, s.Follow, s.Args)
		}
		return
	}

	if err := PassArgs(s, &args); err != nil {
		log.Fatal(err)
	}

	if s.Image == "" {
		log.Fatalf(`dunner: image repository name cannot be empty`)
	}

	results, err := (*s).Exec()
	if err != nil {
		log.Fatal(err)
	}

	if results == nil {
		return
	}

	for _, res := range *results {
		log.Infof(
			"Running task '%+v' on '%+v' Docker with command '%+v'",
			s.Task,
			s.Image,
			res.Command,
		)
		if res.Output != "" {
			fmt.Printf(`OUT: %s`, res.Output)
		}
		if res.Error != "" {
			fmt.Printf(`ERR: %s`, res.Error)
		}
	}
}

// PassArgs replaces argument variables,of the form '`$d`', where d is a number, with dth argument.
func PassArgs(s *docker.Step, args *[]string) error {
	for i, cmd := range s.Commands {
		for j, subStr := range cmd {
			regex := regexp.MustCompile(`\$[1-9][0-9]*`)
			subStr = regex.ReplaceAllStringFunc(subStr, func(str string) string {
				j, err := strconv.Atoi(strings.Trim(str, "$"))
				if err != nil {
					log.Fatal(err)
				}
				return (*args)[j-1]
			})
			s.Commands[i][j] = subStr
		}
	}
	return nil
}

// getDunnerUser returns the user value from step, if empty returns first found value in order:
// UID env variable, current user ID, current user name.
func getDunnerUser(task config.Task) string {
	if task.User != "" {
		return task.User
	}
	dunnerUser := os.Getenv("UID")
	if dunnerUser == "" {
		user, err := os_user.Current()
		if err != nil {
			dunnerUser = os.Getenv("USER")
			log.Debugf("Unable to find current user id: %s. Using `%s` as docker user.", err.Error(), dunnerUser)
		} else {
			dunnerUser = user.Uid
		}
	}
	return dunnerUser
}
