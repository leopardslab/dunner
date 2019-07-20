# Dunner

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b2275e331d2745dc9527d45efbbf2da2)](https://app.codacy.com/app/Leopardslab/dunner?utm_source=github.com&utm_medium=referral&utm_content=leopardslab/dunner&utm_campaign=Badge_Grade_Dashboard)
[![Codecov branch](https://img.shields.io/codecov/c/github/leopardslab/dunner/master.svg)](https://codecov.io/gh/leopardslab/dunner)
[![Build Status](https://travis-ci.org/leopardslab/Dunner.svg?branch=master)](https://travis-ci.org/leopardslab/Dunner)

> The Docker Task Runner

Dunner is a task runner tool like Grunt but uses Docker images like CircleCI do. You can define tasks and steps of the tasks in your `.dunner.yaml` file and then run these steps with `dunner do <taskname>`.


Example `.dunner.yaml`

```yaml
deploy:
- image: 'emeraldsquad/sonar-scanner'
  commands:
  - ['sonar', 'scan']
- image: 'golang'
  commands:
  - ['go', 'install']
- image: 'mesosphere/aws-cli'
  commands:
  - ['aws', 'elasticbeanstalk update-application --application-name myapp']
  envs: 
   - AWS_ACCESS_KEY_ID=`$AWS_KEY`
   - AWS_SECRET_ACCESS_KEY=`$AWS_SECRET`
   - AWS_DEFAULT_REGION=us-east1
- follow: 'status' #This refers to another task and can pass args too
  args: 'prod'
status:
- image: 'mesosphere/aws-cli'
  commands:
  - ['aws', 'elasticbeanstalk describe-events --environment-name $1'] 
  # This uses args passed to the task, `$1` means first arg
  envs: 
   - AWS_ACCESS_KEY_ID=`$AWS_KEY`
   - AWS_SECRET_ACCESS_KEY=`$AWS_SECRET`
   - AWS_DEFAULT_REGION=us-east1
```

Now you can use as,
* `dunner do deploy`
* `dunner do status prod`

# Guides

* [User Guide](https://github.com/leopardslab/dunner/wiki/User-Guide)
* [Installation Guide](https://github.com/leopardslab/dunner/wiki/Installation-Guide)
* [Developer Guide](https://github.com/leopardslab/dunner/wiki/Developer-Guide)

## Development Plan 

### [`v0.1`](https://github.com/leopardslab/dunner/milestone/2)
- [x] Ability to define set of tasks and steps and run the task
- [x] Mount current dir as a volume
- [x] Ability to pass arguments to tasks
### [`v1.0`](https://github.com/leopardslab/dunner/milestone/1) 
- [x] Ability to add ENV variables
- [x] Ability to define the sub-dir that should be mounted to the task containers
- [x] Ability to mount other dirs to the task containers
- [x] Ability to use a task as a step for another task
- [x] Ability to get ENV, param, etc values from host environment variables or `.env` file
- [x] Ability to install as a Snap package

### [`v2.0`](https://github.com/leopardslab/dunner/milestone/3) 
- [x] Ability to Dry Run 
- [x] Ability to verfiy the `.dunner.yaml` file
- [x] Ability to define multiple commands for the same step
- [ ] Ability to install as a Deb package
- [ ] Ability to install as a RPM package
- [ ] Ability to install as a Brew package
