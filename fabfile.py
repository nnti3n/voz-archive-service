from fabric.api import *

env.hosts = ['aws']
env.use_ssh_config = True

def exists(path):
    with settings(warn_only=True):
        return run('test -e %s' % path)

def deploy():
    code_dir = '$GOPATH/src/github.com/nnti3n/voz-archive-service'
    with settings(warn_only=True):
        if exists(code_dir).failed:
            run("go get github.com/nnti3n/voz-archive-service")
    with cd(code_dir):
        run("git pull")
        run("dep ensure")
        run("go build -o voz-worker serviceWorker/main.go")
        run("go build -o voz-interface interface/main.go")
        run("sudo systemctl restart voz-worker")
        run("sudo systemctl restart voz-interface")
        
