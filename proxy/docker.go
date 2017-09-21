package proxy

import (
	//"log"
	"github.com/docker/docker/api/types"
	"os"
	"os/exec"
)

type Docker struct {
	host map[string]string
}

func (doc *Docker) getConnectionArgs() []string {
	var args []string
	args = append(args, "-H", doc.host["hostname"])
	if len(doc.host["certpath"]) > 0 {
		args = append(args, "--tlscacert", doc.host["certpath"]+"/ca.pem")
		args = append(args, "--tlscert", doc.host["certpath"]+"/cert.pem")
		args = append(args, "--tlskey", doc.host["certpath"]+"/key.pem")
		args = append(args, "--tlsverify")
	}
	return args
}

func (doc *Docker) Logs(cid string, args map[string]interface{}) (err error) {

	follow := args["follow"].(bool)
	tail := args["tail"].(string)

	xargs := doc.getConnectionArgs()
	xargs = append(xargs, "logs")

	if follow {
		xargs = append(xargs, "-f")
	}

	if len(tail) > 0 {
		xargs = append(xargs, "--tail", tail)
	}

	xargs = append(xargs, cid)
	docker_cmd := doc.host["binary"]
	cmd := exec.Command(docker_cmd, xargs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	return

}

func (doc *Docker) Proxy(xargs ...string) (err error) {

	args := doc.getConnectionArgs()
	args = append(args, xargs...)

	docker_cmd := doc.host["binary"]
	cmd := exec.Command(docker_cmd, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	return

}

func (doc *Docker) Stats(cid string) (err error) {
	return doc.Proxy("stats", cid)
}

func (doc *Docker) Restart(cid string) (err error) {
	return doc.Proxy("restart", cid)
}

func (doc *Docker) Stop(cid string) (err error) {
	return doc.Proxy("stop", cid)
}

func (doc *Docker) Start(cid string) (err error) {
	return doc.Proxy("start", cid)
}

func (doc *Docker) Remove(cid string) (err error) {
	return doc.Proxy("rm", cid)
}

func (doc *Docker) Destroy(cid string) (err error) {
	err = doc.Proxy("stop", cid)
	if err == nil {
		err = doc.Proxy("rm", cid)
	}
	return
}

func (doc *Docker) Inspect(cid string) (err error) {
	return doc.Proxy("inspect", cid)
}

func (doc *Docker) Connect(cid string) (err error) {

	args := doc.getConnectionArgs()
	args = append(args, "exec", "-it", cid, "/bin/sh")

	cmd := exec.Command(doc.host["binary"], args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	return

}

func (doc *Docker) GetId(container types.Container) string {
	return container.ID
}

func (doc *Docker) GetId2(container types.Container) string {
	return container.ID
}

func NewDockerProxy(conf map[string]string) (OrchestrationProxy, error) {

	cli := new(Docker)
	cli.host = conf

	return cli, nil
}

func CreateDockerCli(config map[string]string) *Docker {
	var cli Docker
	cli.host = config

	return &cli
}
