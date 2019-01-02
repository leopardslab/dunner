package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

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
