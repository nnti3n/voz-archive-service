from fabric.api import *

env.hosts = ['ec2-user@ec2-54-169-119-62.ap-southeast-1.compute.amazonaws.com']

def deploy():
    code_dir = '/home/ec2-user/projects/src/github.com/nnti3n/voz-archive-plus'
    with settings(warn_only=True):
        if run("test -d %s" % code_dir).failed:
            run("go get github.com/nnti3n/voz-archive-plus")
    with cd(code_dir):
        run("git pull")
        run("dep ensure")