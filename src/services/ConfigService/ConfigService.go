package ConfigService

import (
	"log"
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Task struct {
	Name   string    `yaml:"name"`
	Image   string    `yaml:"image"`
	Command [] string `yaml:"command"`
}

type Configs struct {
	Tasks map[string][]Task
}

func GetConfigs() Configs {
	// TODO Should get the name of the Dunner file from a Constant
	fileContents, err := ioutil.ReadFile("./.dunner.yaml")
	if err != nil {
		log.Fatal(err)
	}

	var configs Configs
	if err := yaml.Unmarshal(fileContents, &configs.Tasks); err != nil {
		log.Fatal(err)
	}

	return configs
}
