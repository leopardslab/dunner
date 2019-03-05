package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
	"github.com/leopardslab/Dunner/internal/logger"
	"github.com/spf13/viper"
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

	if err := parseEnv(&configs); err != nil {
		log.Fatal(err)
	}

	return &configs, nil
}

func parseEnv(configs *Configs) error {
	file := viper.GetString("DotenvFile")
	envs, err := godotenv.Read(file)
	if err != nil {
		log.Warn(err)
	}

	for k, tasks := range (*configs).Tasks {
		for j, task := range tasks {
			for i, envVar := range task.Envs {
				var str = strings.Split(envVar, "=")
				if len(str) != 2 {
					return fmt.Errorf(
						`config: invalid format of environment variable: %v`,
						envVar,
					)
				}
				var pattern = "^`\\$.+`$"
				check, err := regexp.MatchString(pattern, str[1])
				if err != nil {
					log.Fatal(err)
				}
				if check {
					var key = strings.Replace(
						strings.Replace(
							str[1],
							"`",
							"",
							-1,
						),
						"$",
						"",
						1,
					)
					var val string
					if v, isSet := os.LookupEnv(key); isSet {
						val = v
					}
					if v, isSet := envs[key]; isSet {
						val = v
					}
					if val == "" {
						return fmt.Errorf(
							`config: could not find environment variable '%v' in %s file or among host environment variables`,
							key,
							file,
						)
					}
					var newEnv = str[0] + "=" + val
					(*configs).Tasks[k][j].Envs[i] = newEnv
				}
			}
		}
	}

	return nil
}
