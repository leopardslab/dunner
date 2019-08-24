package internal

// DefaultTaskFilePermission is the default file permission of dunner task file
const DefaultTaskFilePermission = 0644

// DefaultTaskFileContents is the default dunner taskfile contents, used when initialized with dunner
const DefaultTaskFileContents = `# This is an example dunner task file. Please make any required changes.
# (Optional) Set any environment variables to be exported in the container
# for every step of every task (can be overridden)
envs:
  - PERM=775

# (Optional) List of directories that are to be mounted on the container
# for every step of every task (can be overridden)
mounts:
  - /tmp:/tmp:w

# List of all task objects
tasks:
  build:
    # (Optional) Set any environment variables to be exported in the container
    # for every step of 'build' task (can be overridden)
    envs:
      - PERM=775

    # (Optional) List of directories that are to be mounted on the container
    # for every step of 'build' task (can be overridden)
    mounts:
      - /tmp:/tmp:w

    # List of all step objects for 'build' task
    steps:
      - name: setup
        # Image name that has to be pulled from a registry
        image: node:latest
        # List of commands that has to be run inside the container
        commands:
          - ["npm", "--version"]
          - ["npm", "install"]
        # (Optional) List of directories that are to be mounted on the container
        mounts:
          - /tmp:/tmp:w
        # (Optional) Set any environment variables to be exported in the container
        envs:
          - PERM=775
`

// DefaultDunnerTaskFileName is the default dunner task file name
const DefaultDunnerTaskFileName = ".dunner.yaml"
