package config

import (
	"bytes"
	"github.com/docker/docker/api/types"
	"strings"
)

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
	termsc["name"] = terms["name"]
	termsc["command"] = Truncate(terms["command"], 32)
	termsc["image"] = terms["image"]
	termsc["address"] = terms["address"]
	termsc["network"] = terms["network"]
	termsc["orchestration"] = terms["orchestration"]

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
	orch := ""
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

	return terms
}
