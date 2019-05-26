package internal

// DefaultTaskFileContents is the default dunner taskfile contents, used when initialized with dunner
const DefaultTaskFileContents = `# This is an example dunner task file. Please make any required changes.
build:
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

// DefaultTaskFilePermission is the default file permission of dunner task file
const DefaultTaskFilePermission = 0644
