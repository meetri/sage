package clients

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/meetri/ymltree"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"
)

func NewDockerClient(host string, certpath string, ver string, timeout int) (cli *client.Client, err error) {

	var cl *http.Client

	if len(certpath) > 0 {
		dockerCertPath := certpath
		options := tlsconfig.Options{
			CAFile:             filepath.Join(dockerCertPath, "ca.pem"),
			CertFile:           filepath.Join(dockerCertPath, "cert.pem"),
			KeyFile:            filepath.Join(dockerCertPath, "key.pem"),
			InsecureSkipVerify: os.Getenv("DOCKER_TLS_VERIFY") == "",
		}
		tlsc, err := tlsconfig.Client(options)
		if err != nil {
			return nil, err
		}

		to := time.Duration(time.Duration(timeout) * time.Second)

		cl = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsc,
			},
			Timeout: to,
		}

	}

	cli, err = client.NewClient(host, ver, cl, nil)

	return
}

func GetAllContainers(hosts []interface{}, certglobal string, timeout int) []interface{} {

	numcores := runtime.GOMAXPROCS(0)
	if numcores > runtime.NumCPU() {
		numcores = runtime.NumCPU()
	}
	sem := make(chan int, numcores)

	var wg sync.WaitGroup

	l := make([]interface{}, len(hosts))
	for idx, host := range hosts {
		wg.Add(1)
		sem <- 1
		go func(id int, selhost interface{}) {
			getContainerDetails(id, l, selhost, certglobal, timeout, sem)
			wg.Done()
		}(idx, host)
	}
	wg.Wait()

	return l

}

func getContainerDetails(idx int, l []interface{}, hostdata interface{}, certglobal string, timeout int, sem chan int) {

	var hostname string
	var hostalias string
	var certpath string
	var dockerver string
	var timeout_seconds int

	defer func() {
		<-sem
	}()

	if reflect.TypeOf(hostdata).Kind() == reflect.String {
		hostname = hostdata.(string)
		hostalias = hostname
		certpath = certglobal
		timeout_seconds = timeout
	} else if reflect.TypeOf(hostdata).Kind() == reflect.Map {
		hostname = hostdata.(ymltree.Map).FindDefault("host", "")
		hostalias = hostdata.(ymltree.Map).FindDefault("alias", hostname)
		certpath = hostdata.(ymltree.Map).FindDefault("certpath", certglobal)
		timeout_seconds = hostdata.(ymltree.Map).FindDefaultInt("timeout", timeout)
	}

	cli, err := NewDockerClient(hostname, certpath, dockerver, timeout_seconds)
	if err == nil {
		containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true, Size: false})
		if err == nil {
			m := make(ymltree.Map)
			m["container"] = containers
			m["hostalias"] = hostalias
			m["hostname"] = hostname
			m["cli"] = cli
			l[idx] = m
		} else {
			log.Printf("failed to get container details")
		}
	} else {
		log.Printf("failed to create client for %s", hostname)

	}

}
