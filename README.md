# Dunner

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/b2275e331d2745dc9527d45efbbf2da2)](https://app.codacy.com/app/rehrumesh/Dunner?utm_source=github.com&utm_medium=referral&utm_content=leopardslab/Dunner&utm_campaign=Badge_Grade_Settings)

Dunner is a task runner tool like Grunt but uses Docker images like CircleCI do. You can define tasks and steps of the tasks in your `.dunner.yaml` file and then run these steps with `Dunner do taskname`


Example `.dunner.yaml`

```yaml
deploy:
- name: 'emeraldsquad/sonar-scanner'
  command: ['sonar', 'scan']
- name: 'golang'
  command: ['go', 'install']
- name: 'mesosphere/aws-cli'
  command: ['aws', 'elasticbeanstalk update-application --application-name myapp']
  envs: 
   - AWS_ACCESS_KEY_ID=`$AWS_KEY`
   - AWS_SECRET_ACCESS_KEY=`$AWS_SECRET`
   - AWS_DEFAULT_REGION=us-east1
- name: '@status' #This refers to another task and can pass args too
  args: 'prod'
status:
- name: 'mesosphere/aws-cli'
  command: ['aws', 'elasticbeanstalk describe-events --environment-name $1'] 
  # This uses args passed to the task, `$1` means first arg
  envs: 
   - AWS_ACCESS_KEY_ID=`$AWS_KEY`
   - AWS_SECRET_ACCESS_KEY=`$AWS_SECRET`
   - AWS_DEFAULT_REGION=us-east1
```

Now you can use as,
 1. `Dunner do deploy`
 2. `Dunner do status prod`


## NOTE
This work is still in progress. See the development plan.

## Development Plan 

### [`v0.1`](https://github.com/leopardslab/Dunner/milestone/2)
- [x] Ability to define set of tasks and steps and run the task
- [x] Mount current dir as a volume
- [ ] Ability to pass arguments to tasks
### [`v1.0`](https://github.com/leopardslab/Dunner/milestone/1) 
- [ ] Ability to add ENV variables
- [ ] Ability to define the sub-dir that should be mounted to the task containers
- [ ] Ability to mount other dirs to the task containers
- [ ] Ability to use a task as a step for another task
- [ ] Ability to get ENV, param, etc values from host environment variables or `.env` file

# Guides

* [User Guide](https://github.com/leopardslab/Dunner/wiki/User-Guide)
* [Installation Guide](https://github.com/leopardslab/Dunner/wiki/Installation-Guide)
* [Developer Guide](https://github.com/leopardslab/Dunner/wiki/Developer-Guide)
