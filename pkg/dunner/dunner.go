package dunner

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/docker/docker/pkg/stdcopy"
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
			Task:    taskName,
			Name:    stepDefinition.Name,
			Image:   stepDefinition.Image,
			Command: stepDefinition.Command,
			Env:     stepDefinition.Envs,
			WorkDir: stepDefinition.SubDir,
			Args:    stepDefinition.Args,
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

	if newTask := regexp.MustCompile(`^@\w+$`).FindString(s.Name); newTask != "" {
		newTask = strings.Trim(newTask, "@")
		if async {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				execTask(configs, newTask, s.Args)
				wg.Done()
			}(wg)
		} else {
			execTask(configs, newTask, s.Args)
		}
		return
	}

	if err := passArgs(s, &args); err != nil {
		log.Fatal(err)
	}

	if s.Image == "" {
		log.Fatalf(`dunner: image repository name cannot be empty`)
	}

	pout, err := (*s).Exec()
	if err != nil {
		log.Fatal(err)
	}

	log.Infof(
		"Running task '%+v' on '%+v' Docker with command '%+v'",
		s.Task,
		s.Image,
		strings.Join(s.Command, " "),
	)

	if _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, *pout); err != nil {
		log.Fatal(err)
	}

	if err = (*pout).Close(); err != nil {
		log.Fatal(err)
	}
}

func passArgs(s *docker.Step, args *[]string) error {
	for i, subStr := range s.Command {
		regex := regexp.MustCompile(`\$[1-9][0-9]*`)
		subStr = regex.ReplaceAllStringFunc(subStr, func(str string) string {
			j, err := strconv.Atoi(strings.Trim(str, "$"))
			if err != nil {
				log.Fatal(err)
			}
			return (*args)[j-1]
		})
		s.Command[i] = subStr
	}
	return nil
}
