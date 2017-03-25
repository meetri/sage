package proxy

import (
	"log"
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

func (doc *Docker) Connect(cid string) (err error) {

	args := doc.getConnectionArgs()
	args = append(args, "exec", "-it", cid, "/bin/sh")

	cmd := exec.Command("docker", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Printf("starting...")
	err = cmd.Start()
	if err == nil {
		log.Printf("waiting...")
		err = cmd.Wait()
	}

	log.Printf("done...")

	return

}

func CreateDockerCli(config map[string]string) *Docker {
	var cli Docker
	cli.host = config

	return &cli
}
