package proxy

import (
	"github.com/docker/docker/api/types"
	"os"
	"os/exec"
	"strings"
)

type DCOS struct {
	host map[string]string
}

func (doc *DCOS) Logs(cid string, args map[string]interface{}) (err error) {
	return doc.Proxy("marathon", "app", "show", cid)
}

func (doc *DCOS) Restart(cid string) (err error) {
	return doc.Proxy("marathon", "app", "restart", cid)
}

func (doc *DCOS) Stop(cid string) (err error) {
	return doc.Proxy("marathon", "app", "stop", cid)
}

func (doc *DCOS) Remove(cid string) (err error) {
	return nil
}

func (doc *DCOS) Destroy(cid string) (err error) {
	return doc.Proxy("marathon", "app", "kill", cid)
}

func (doc *DCOS) Start(cid string) (err error) {
	return doc.Proxy("marathon", "app", "start", cid)
}

func (doc *DCOS) Proxy(args ...string) (err error) {

	cmd := exec.Command("dcos", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err == nil {
		err = cmd.Wait()
	}

	return

}

func NewDCOSProxy(conf map[string]string) (OrchestrationProxy, error) {

	cli := new(DCOS)
	cli.host = conf
	return cli, nil
}

func (doc *DCOS) GetId2(container types.Container) string {
	for k, v := range container.Labels {
		if k == "MESOS_TASK_ID" {
			return v
		}
	}
	return ""
}

func (doc *DCOS) GetId(container types.Container) string {
	name := ""
	for k, v := range container.Labels {
		if k == "MESOS_TASK_ID" {
			name = strings.Split(v, ".")[0]
			break
		}
	}
	return strings.Replace(name, "_", "/", -1)
}
