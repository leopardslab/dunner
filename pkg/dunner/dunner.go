package dunner

import (
	"fmt"
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

	execTask(configs, args[0], args[1:])
}

func execTask(configs *config.Configs, taskName string, args []string) {
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
		}

		if err := config.DecodeMount(stepDefinition.Mounts, &step); err != nil {
			log.Fatal(err)
		}

		if async {
			go process(configs, &step, &wg, args)
		} else {
			process(configs, &step, &wg, args)
		}
	}

	wg.Wait()
}

func process(configs *config.Configs, s *docker.Step, wg *sync.WaitGroup, args []string) {
	var async = viper.GetBool("Async")
	if async {
		defer wg.Done()
	}

	if s.Follow != "" {
		if async {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				execTask(configs, s.Follow, s.Args)
				wg.Done()
			}(wg)
		} else {
			execTask(configs, s.Follow, s.Args)
		}
		return
	}

	if err := passArgs(s, &args); err != nil {
		log.Fatal(err)
	}

	if s.Image == "" {
		log.Fatalf(`dunner: image repository name cannot be empty`)
	}

	results, err := (*s).Exec()
	if err != nil {
		log.Fatal(err)
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

func passArgs(s *docker.Step, args *[]string) error {
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
