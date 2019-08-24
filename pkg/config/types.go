package config

// Step defines a single step for a task
type Step struct {
	// Name given as string to identify the task
	Name string `yaml:"name"`

	// Image is the repo name on which Docker containers are built
	Image string `yaml:"image" validate:"required_without=Follow"`

	// SubDir is the primary directory on which task is to be run
	SubDir string `yaml:"dir"`

	// The command which runs on the container and exits
	Command []string `yaml:"command" validate:"omitempty,dive,required"`

	// The list of commands that are to be run in sequence
	Commands [][]string `yaml:"commands" validate:"omitempty,dive,omitempty,dive,required"`

	// The list of environment variables to be exported inside the container
	Envs []string `yaml:"envs"`

	// The directories to be mounted on the container as bind volumes
	Mounts []string `yaml:"mounts" validate:"omitempty,dive,min=1,mountdir,parsedir"`

	// The next task that must be executed if this does go successfully
	Follow string `yaml:"follow" validate:"omitempty,follow_exist"`

	// The list of arguments that are to be passed
	Args []string `yaml:"args"`

	// User that will run the command(s) inside the container, also support user:group
	User string `yaml:"user"`
}

// Task describes a single task composed of multiple steps to be run in a docker container
type Task struct {
	Envs   []string `yaml:"envs"`   // Environment variables common to all steps
	Mounts []string `yaml:"mounts"` // Directory mounts common to all steps
	Steps  []Step   `yaml:"steps"`
}

// Configs describes the parsed information from the dunner file.
// It is a map of task name as keys and the list of tasks associated with it.
type Configs struct {
	Envs   []string        `yaml:"envs"`   // Environment variables common to all tasks
	Mounts []string        `yaml:"mounts"` // Directory mounts common to all tasks
	Tasks  map[string]Task `yaml:"tasks" validate:"dive,keys,required,endkeys,required,min=1,required"`
}
