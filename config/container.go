package config

import (
	"bytes"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/meetri/ymltree"
	"strings"
)

func SearchContainers(args []string, ac []interface{}) (mh map[string]string, mc types.Container, err error) {

	matched := 0
	for _, cdata := range ac {

		if cdata != nil {

			containers := cdata.(ymltree.Map).Find("container")
			hostdetails := cdata.(ymltree.Map).Find("host").(map[string]string)

			for _, container := range containers.([]types.Container) {

				terms := GetContainerTerms(container)

				match := false

				if len(args) == 0 {
					match = true
				} else {
					for _, arg := range args {
						if MatchTerms(arg, terms) {
							match = true
							break
						}
					}
				}

				if match {
					matched++
					mc = container
					mh = hostdetails
				}

			}

		}
	}
	if matched == 0 {
		err = errors.New("no containers found")
	}

	return

}

func Truncate(str string, strlen int) string {

	if len(str) > strlen {
		str = str[0:strlen]
	}

	return str
}

func MatchTerms(term string, termlist map[string]string) bool {

	for _, matchterm := range termlist {
		if strings.Contains(matchterm, term) {
			return true
		}
	}

	return false

}

func WriteTerms(terms map[string]string, termlist []string) string {

	var buffer bytes.Buffer
	for _, term := range termlist {
		buffer.WriteString(terms[term] + "\t")
	}

	buffer.WriteString("\n")

	return buffer.String()
}

func TruncateContainerTerms(terms map[string]string) map[string]string {

	termsc := make(map[string]string)
	for k, v := range terms {
		termsc[k] = v
	}

	termsc["id"] = Truncate(terms["id"], 16)
	termsc["command"] = Truncate(terms["command"], 32)

	return termsc

}

func GetContainerTerms(container types.Container) map[string]string {

	networkMode := container.HostConfig.NetworkMode
	ipaddr := ""
	if container.NetworkSettings.Networks[networkMode] != nil {
		ipaddr = container.NetworkSettings.Networks[networkMode].IPAddress
	} else if container.NetworkSettings.Networks["bridge"] != nil {
		ipaddr = container.NetworkSettings.Networks["bridge"].IPAddress
	}

	name := container.Names[0][1:]
	orch := "docker"
	for k, v := range container.Labels {
		if k == "MESOS_TASK_ID" {
			name = strings.Split(v, ".")[0]
			orch = "mesos"
			break
		}
	}

	terms := make(map[string]string)
	terms["id"] = container.ID
	terms["name"] = name
	terms["command"] = container.Command
	terms["image"] = container.Image
	terms["address"] = ipaddr
	terms["network"] = container.HostConfig.NetworkMode
	terms["orchestration"] = orch
	terms["status"] = container.Status
	terms["state"] = container.State

	return terms
}
