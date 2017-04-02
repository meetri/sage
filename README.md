# sage

### Features
- Aggregate all your docker containers from multiple hosts
- Detect orchestration platform per container
- perform proxy operations on containers based on detected orchestration platform
  utilizes `dcos` for DCOS, ( `kubectl` for Kubernetes ), and `docker` for docker 
- Easily connect, view logs, stop, start, and restart containers on any host

### Example Host Configuration ~/.sage/hosts.yml
```
name: DCOS Cluster 1
alias: dcos1
certpath: /projects/docker/certs
hosts:
  - { alias: "mesos-master", host: "tcp://node1.mydockerhost.com:2376" }
  - { alias: "mesos-node2", host: "tcp://node2.mydockerhost.com:2376" }
  - { alias: "mesos-node3", host: "tcp://node3.mydockerhost.com:2376" }
  - { alias: "mesos-node4", host: "tcp://node4.mydockerhost.com:2376" }
  - { alias: "local", host: "tcp://192.168.59.104:2376", certpath: "" }
```


### Help `sage --help`
```
NAME:
   sage - sage [OPTIONS] [ARGUMENTS]

USAGE:
   sage [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     list, ls, ps         list running containers
     connect, co, hijack  connect to a running container
     logs, log, l         view container logs
     inspect, inspect, i  inspect container
     docker, d, doc       docker proxy
     restart              restart container
     stop                 stop container
     start                start container
     remove, rm           remove container
     help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value   run on selected hosts (default: "all")
   --help, -h     show help
   --version, -v  print the version
```
