package core

import (
	"github.com/meetri/ymltree"
	"os"
	"sync"
)

type ConfigData struct {
	data  ymltree.Map
	hosts ymltree.Map
}

var (
	config     ConfigData
	configOnce sync.Once
)

func Config() *ConfigData {

	configOnce.Do(func() {
		home := os.Getenv("HOME")
		config.hosts, _ = ymltree.Load(home + "/.sage/hosts.yml")
		config.data, _ = ymltree.Load(home + "/.sage/config.yml")
	})

	return &config

}

func (this ConfigData) Hosts() ymltree.ConfigMap {
	return this.hosts
}

func (this ConfigData) Data() ymltree.ConfigMap {
	return this.data
}
