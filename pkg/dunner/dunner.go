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
)

var log = logger.Log

// Do method is invoked for command-line use
func Do(_ *cobra.Command, args []string) {
	const async = false

	// TODO Should get the name of the Dunner file
	// from a constant or an environment variable or config file
	var dunnerFile = ".dunner.yaml"

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
	const async = false
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
