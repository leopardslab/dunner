from src.docker import Docker
from src.configs import Configs

class Do:

    def do_commad(self, params):

        steps = Configs().get(params[0])

        for step in steps:

            docker = Docker()
            docker.set_image(step['image'])
            docker.set_cmd(step['command'])

            if 'envs' in step:
                docker.set_envs(step['envs'])

            docker.do()
