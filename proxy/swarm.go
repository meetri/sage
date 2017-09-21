package proxy

import (
//"log"
//"github.com/docker/docker/api/types"
//"os"
//"os/exec"
)

type Swarm struct {
	*Docker
	host map[string]string
}

func NewSwarmProxy(conf map[string]string) (OrchestrationProxy, error) {

	cli := new(Swarm)
	cli.host = conf

	return cli, nil
}

func CreateSwarmCli(config map[string]string) *Swarm {
	var cli Swarm
	cli.host = config

	return &cli
}
