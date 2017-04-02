package proxy

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
)

type OrchestrationProxyFactory func(conf map[string]string) (OrchestrationProxy, error)

type OrchestrationProxy interface {
	Logs(cid string, args map[string]interface{}) (err error)
	Restart(cid string) (err error)
	Stop(cid string) (err error)
	Start(cid string) (err error)
	Remove(cid string) (err error)
	GetId(container types.Container) string
	GetId2(container types.Container) string
}

var (
	orchestrationProxyList = make(map[string]OrchestrationProxyFactory)
)

func init() {
	Register("docker", NewDockerProxy)
	Register("mesos", NewDCOSProxy)
}

func Register(name string, prox OrchestrationProxyFactory) {

	if _, ok := orchestrationProxyList[name]; !ok {
		orchestrationProxyList[name] = prox
	}
}

func Create(name string, conf map[string]string) (OrchestrationProxy, error) {
	oProxy, ok := orchestrationProxyList[name]
	if ok {
		return oProxy(conf)
	}
	return nil, errors.New(fmt.Sprintf("%s proxy has not been registered", name))
}
