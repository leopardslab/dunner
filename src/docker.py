

class Docker:
    def __init__(self):
        self.image = ''
        self.cmd = []
        self.envs = []
        self.workDir = ''

    def set_image(self, image):
        if image is not '':
            self.image = image

    def set_cmd(self, cmd):
        if cmd is not []:
            self.cmd = cmd

    def set_envs(self, envs):
        if envs is not []:
            self.envs = envs

    def set_work_dir(self, work_dir):
        if work_dir is not '':
            self.workDir = work_dir

    def do(self):
        if self.image is not '' and self.cmd is not []:
            print(
                "Creating a container with Docker '"+self.image+
                "' with command '"+" ".join(self.cmd)+"'"
            )
        if self.envs is not []:
            print(
                "using envs '"+",".join(self.envs)+"'"
            )
