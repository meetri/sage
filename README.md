# sage

### Features
- Aggregate all your docker containers from multiple hosts
- Easily connect, view logs, stop, start, and restart containers on any host
- gather important container naming details based on detected orchestration platform used

### Example Host Configuration ~/.sage/hosts.yml
```
name: MySwarm
alias: dcos1
active_cluster: all
certpath: /Users/woiam/.sage/keys/dcos/
# cluster config options coming in next version...
clusters:
    - { name: "cloud1", certpath: "/Users/woiam/.sage/keys/dcos" }
    - { name: "local", certpath: "" }
hosts:
  - { alias: "swarm-node1", cluster: "cloud1", host: "tcp://mesos-node1.mycloud.com:2376" }
  - { alias: "swarm-node2", cluster: "cloud1", host: "tcp://mesos-node6.mycloud.com:2376" }
  - { alias: "swarm-node3", cluster: "cloud1", host: "tcp://mesos-node7.mycloud.com:2376" }
  - { alias: "local", host: "tcp://192.168.59.104:2376" }
  - { alias: "pi", host: "tcp://raspberry.pi:2376","binary":"/Users/whiam/.dvm/bin/docker/1.11.0/docker","certpath":"" }
```

### Example Configuration ~/.sage/config.yml
define which fields get shown during a [sage ps] command
```
fields:
  - id
  - hostalias
  - orchestration
  - name
  - network
  - address
  - image
  - command
  - ports
  - labels
```


### Help `sage --help`
```
NAME:
   sage - sage [OPTIONS] [ARGUMENTS]

USAGE:
   sage [global options] command [command options] [arguments...]

VERSION:
   0.1.2

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
     destroy              stop and remove container
     env                  env -s host1
     stats, stats         get container stats
     help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --host value   run on selected hosts (default: "all")
   --help, -h     show help
   --version, -v  print the version
```
