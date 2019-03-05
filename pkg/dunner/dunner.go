package dunner

import (
	"os"
	"strings"
	"sync"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/leopardslab/Dunner/internal/logger"
	"github.com/leopardslab/Dunner/pkg/config"
	"github.com/leopardslab/Dunner/pkg/docker"
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

	var wg sync.WaitGroup
	for _, stepDefinition := range configs.Tasks[args[0]] {
		if async {
			wg.Add(1)
		}
		step := docker.Step{
			Task:    args[0],
			Name:    stepDefinition.Name,
			Image:   stepDefinition.Image,
			Command: stepDefinition.Command,
			Env:     stepDefinition.Envs,
		}

		if err := config.DecodeMount(stepDefinition.Mounts, &step); err != nil {
			log.Fatal(err)
		}

		if async {
			go process(&step, &wg)
		} else {
			process(&step, &wg)
		}
	}

	wg.Wait()
}

func process(s *docker.Step, wg *sync.WaitGroup) {
	var async = viper.GetBool("Async")
	if async {
		defer wg.Done()
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
