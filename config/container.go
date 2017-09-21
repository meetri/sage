package config

import (
	"bytes"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/meetri/ymltree"
	"strings"
)

type SearchResults struct {
	matched       int
	containerList types.Container
	hostdetails   map[string]string
}

func (s SearchResults) IsMatched() bool {
	return s.matched > 0
}

func (s SearchResults) Matches() int {
	return s.matched
}

func (s SearchResults) IsMatchedOne() bool {
	return s.matched == 1
}

func (s SearchResults) Containers() types.Container {
	return s.containerList
}

func (s SearchResults) Hosts() map[string]string {
	return s.hostdetails
}

func (s *SearchResults) setMatch(c types.Container, h map[string]string) {
	s.matched++
	s.containerList = c
	s.hostdetails = h
}

func SearchContainersNew(args []string, ac []interface{}) (matchedList []SearchResults, err error) {

	mlen := 1
	if len(args) > 0 {
		mlen = len(args)
	}
	//matchedList := make([]SearchResults, mlen)
	matchedList = make([]SearchResults, mlen)

	total_matches := 0
	for _, cdata := range ac {

		if cdata != nil {

			containers := cdata.(ymltree.Map).Find("container")
			hostdetails := cdata.(ymltree.Map).Find("host").(map[string]string)

			for _, container := range containers.([]types.Container) {

				terms := GetContainerTerms(container)

				if len(args) == 0 {
					matchedList[0].setMatch(container, hostdetails)
					total_matches++
				} else {
					for idx, arg := range args {
						if MatchTerms(arg, terms) {
							matchedList[idx].setMatch(container, hostdetails)
							total_matches++
							break
						}
					}
				}

			}

		}
	}

	if total_matches == 0 {
		err = errors.New("no containers found\n")
	}

	return

}

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
		err = errors.New("no containers found\n")
	} else if matched > 1 {
		err = errors.New("found multiple containers that matched query\n")
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

	ipaddr := ""
	networkname := ""
	for netname, network := range container.NetworkSettings.Networks {
		ipaddr = network.IPAddress
		networkname = netname
		break
	}

	name := container.Names[0][1:]
	orch := "docker"
	for k, v := range container.Labels {
		if k == "MESOS_TASK_ID" {
			name = strings.Split(v, ".")[0]
			orch = "mesos"
			break
		} else if k == "com.docker.swarm.task.name" {
			name = strings.Split(v, ".")[0] + "." + strings.Split(v, ".")[1]
			orch = "swarm"
			break
		}
	}

	terms := make(map[string]string)
	terms["id"] = container.ID
	terms["name"] = name
	terms["command"] = container.Command
	terms["image"] = container.Image
	terms["address"] = ipaddr
	terms["network"] = networkname
	terms["orchestration"] = orch
	terms["status"] = container.Status
	terms["state"] = container.State

	return terms
}
