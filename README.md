# Dunner [![Release](https://img.shields.io/github/release/leopardslab/dunner.svg)](https://img.shields.io/github/release/leopardslab/dunner.svg)

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b2275e331d2745dc9527d45efbbf2da2)](https://app.codacy.com/app/Leopardslab/dunner?utm_source=github.com&utm_medium=referral&utm_content=leopardslab/dunner&utm_campaign=Badge_Grade_Dashboard)
[![Codecov branch](https://img.shields.io/codecov/c/github/leopardslab/dunner/master.svg)](https://codecov.io/gh/leopardslab/dunner)
[![Build Status](https://travis-ci.org/leopardslab/Dunner.svg?branch=master)](https://travis-ci.org/leopardslab/Dunner)
[![GoDoc](https://godoc.org/github.com/leopardslab/dunner?status.svg)](https://godoc.org/github.com/leopardslab/dunner)
[![GoReport](https://goreportcard.com/badge/github.com/leopardslab/dunner)](https://goreportcard.com/report/github.com/leopardslab/dunner)
[![Join the chat at https://gitter.im/LeaopardLabs/Dunner](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/LeaopardLabs/Dunner)

> The Docker Task Runner

Dunner is a task runner tool based on Docker, simple and flexible. You can define tasks and configure the environment in your `.dunner.yaml` file and then run as `dunner do <taskname>`.

Example `.dunner.yaml`:

```yaml
envs:
  - AWS_ACCESS_KEY_ID=`$AWS_KEY`
  - AWS_SECRET_ACCESS_KEY=`$AWS_SECRET`
  - AWS_DEFAULT_REGION=us-east1
tasks:
  deploy:
    steps:
      - image: 'emeraldsquad/sonar-scanner'
        commands:
          - ['sonar', 'scan']
      - image: 'golang'
        commands:
          - ['go', 'install']
      - image: 'mesosphere/aws-cli'
        commands:
          - ['aws', 'elasticbeanstalk update-application --application-name myapp']
      - follow: 'status' #This refers to another task and can pass args too
        args: 'prod'
  status:
    steps:
      - image: 'mesosphere/aws-cli'
        commands:
          # This uses args passed to the task, `$1` means first arg
          - ['aws', 'elasticbeanstalk describe-events --environment-name $1']
```

Running `dunner do deploy` from command-line executes `deploy` task inside a Docker container. It creates a Docker container using specified image, executes given commands and shows results, all with just simple configuration!


## Features

* [Easy Configuration](https://github.com/leopardslab/dunner/wiki/User-Guide#how-to-write-a-dunner-file) to run tasks inside container
* [Multiple commands](https://github.com/leopardslab/dunner/wiki/User-Guide#multiple-commands)
* [Mount external directories](https://github.com/leopardslab/dunner/wiki/User-Guide#mounting-external-directories)
* [Environment Variables](https://github.com/leopardslab/dunner/wiki/User-Guide#exporting-environment-variables)
* [Pass arguments through CLI](https://github.com/leopardslab/dunner/wiki/User-Guide#passing-arguments-through-cli)
* [Add dependent task](https://github.com/leopardslab/dunner/wiki/User-Guide#use-a-task-as-a-step-for-another-task)
* [Asynchronous mode](https://github.com/leopardslab/dunner/wiki/User-Guide#asynchronous-mode)
* [Dry Run](https://github.com/leopardslab/dunner/wiki/User-Guide#dry-run)

and [more](https://github.com/leopardslab/dunner/wiki/User-Guide)...

# Getting Started

Read more about [Why Dunner](https://github.com/leopardslab/dunner/wiki/Introduction-to-Dunner) and refer our guides for [installation](https://github.com/leopardslab/Dunner/wiki/Installation-Guide) and [usage](https://github.com/leopardslab/dunner/wiki/User-Guide).

| [**User Documentation**](https://github.com/leopardslab/dunner/wiki/User-Guide)     | [**Installation Guide**](https://github.com/leopardslab/dunner/wiki/Installation-Guide)     | [**Dunner Examples**](https://github.com/leopardslab/dunner-cookbook)           | [**Contributing**](https://github.com/leopardslab/dunner/wiki/Developer-Guide)           | [**Dunner GoCD Plugin**](https://github.com/leopardslab/dunner-gocd-plugin#dunner-gocd-plugin)           | 
|:-------------------------------------:|:-------------------------------:|:-----------------------------------:|:---------------------------------------------:| :--------------------------------------:|
| Learn more about using Dunner | Getting started with Dunner | Dunner Cookbook Recipes | How can you contribute to Dunner? | Have a look at Dunner [GoCD](https://www.gocd.org/) Plugin |


## Development Plan

Have a look at [Dunner Milestones](https://github.com/leopardslab/dunner/milestones) for our future plans.


## Contributing

We'd love your help to fix bugs and add features. The maintainers actively manage the issues list, and try to highlight issues suitable for newcomers. The project follows the typical GitHub pull request model. Before starting any work, please either comment on an existing issue, or file a new one. Refer our [Developer Guide](https://github.com/leopardslab/dunner/wiki/Developer-Guide) for more.
#Dunner logo
https://github.com/leopardslab/dunner/issues/190#issue-543767762
