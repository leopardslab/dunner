import yaml

DUNNER_FILE = '.dunner.yaml'

class Configs:
    def __init__(self):
        with open("./"+DUNNER_FILE, "r") as config_file:
            config_contents = config_file.read()
            self.configs = yaml.load(config_contents)

    def get(self, task):
        return self.configs[task]
