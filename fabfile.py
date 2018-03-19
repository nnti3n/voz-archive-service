from fabric.api import *

env.hosts = ['aws']
env.use_ssh_config = True

def deploy():
    code_dir = '/home/ec2-user/projects/src/github.com/nnti3n/voz-archive-service'
    with settings(warn_only=True):
        if run("test -d %s" % code_dir).failed:
            run("go get github.com/nnti3n/voz-archive-service")
    with cd(code_dir):
        run("git pull")
        run("dep ensure")
        run("go build -o voz-worker serviceWorker/main.go")
        run("go build -o voz-interface interface/main.go")
        run("sudo systemctl start voz-worker")
        run("sudo systemctl start voz-interface")
        
