package config

import (
	"github.com/leopardslab/Dunner/internal/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var log = logger.Log

type Task struct {
	Name    string    `yaml:"name"`
	Image   string    `yaml:"image"`
	Command [] string `yaml:"command"`
}

type Configs struct {
	Tasks map[string][]Task
}

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

func (c *Configs) GetAllImages() ([] string, error) {
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
