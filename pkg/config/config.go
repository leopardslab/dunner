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

// GetAllImages extracts a set of images required for all the tasks to run
func (c *Configs) GetAllImages() ([]string, error) {
	var images []string
	var set = make(map[string]bool)

	for _, tasks := range c.Tasks {
		for _, task := range tasks {
			set[task.Image] = true
		}
	}

	for img := range set {
		images = append(images, img)
	}
	return images, nil
}
