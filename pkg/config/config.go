package config

import (
	"io/ioutil"

	"github.com/leopardslab/Dunner/internal/logger"
	yaml "gopkg.in/yaml.v2"
)

var log = logger.Log

// Task describes a single task to be run in a docker container
type Task struct {
	Name    string   `yaml:"name"`
	Image   string   `yaml:"image"`
	Command []string `yaml:"command"`
	Envs    []string `yaml:"envs"`
}

// Configs describes the parsed information from the dunner file
type Configs struct {
	Tasks map[string][]Task
}

// GetConfigs reads and parses tasks from the dunner file
func GetConfigs(filename string) (*Configs, error) {
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var configs Configs
	if err := yaml.Unmarshal(fileContents, &configs.Tasks); err != nil {
		log.Fatal(err)
	}

	return &configs, nil
}
