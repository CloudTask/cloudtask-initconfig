package main

import "github.com/samuel/go-zookeeper/zk"
import yaml "gopkg.in/yaml.v2"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

//Configuration is exported
//zookeeper cluster parameters.
type Configuration struct {
	Zookeeper struct {
		Hosts string `yaml:"hosts" json:"hosts"`
		Root  string `yaml:"root" json:"root"`
	} `yaml:"zookeeper" json:"zookeeper"`
}

var (
	zkFlags       = int32(0)
	zkACL         = zk.WorldACL(zk.PermAll)
	zkConnTimeout = time.Second * 15
)

func readConfiguration(file string) (*Configuration, error) {

	fd, err := os.OpenFile(file, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	buf, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	conf := &Configuration{}
	if err := yaml.Unmarshal(buf, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func initServerConfigData(hosts string, root string) (string, []byte, error) {

	servers := strings.Split(hosts, ",")
	conn, event, err := zk.Connect(servers, zkConnTimeout)
	if err != nil {
		return "", nil, err
	}

	<-event
	defer conn.Close()
	serverConfigPath := root + "/ServerConfig"
	ret, _, err := conn.Exists(serverConfigPath)
	if err != nil {
		return "", nil, err
	}

	data, err := ioutil.ReadFile("./ServerConfig.json")
	if err != nil {
		return "", nil, fmt.Errorf("ServerConfig.json read failure, %s", err)
	}

	if !ret {
		if _, err := conn.Create(serverConfigPath, data, zkFlags, zkACL); err != nil {
			return "", nil, err
		}
	} else {
		if _, err := conn.Set(serverConfigPath, data, -1); err != nil {
			return "", nil, err
		}
	}
	return serverConfigPath, data, nil
}

func main() {

	var (
		configFile string
		zkHosts    string
		zkRoot     string
	)

	flag.StringVar(&configFile, "f", "./config.yaml", "coudtask initconfig etc.")
	flag.StringVar(&zkHosts, "hosts", "127.0.0.1:2181", "zookeeper hosts.")
	flag.StringVar(&zkRoot, "root", "/cloudtask", "zookeeper root path.")
	flag.Parse()

	if configFile != "" {
		conf, err := readConfiguration(configFile)
		if err != nil {
			fmt.Errorf("config file invalid, %s", err)
			return
		}
		zkHosts = conf.Zookeeper.Hosts
		zkRoot = conf.Zookeeper.Root
	}

	if ret := strings.HasPrefix(zkRoot, "/"); !ret {
		zkRoot = "/" + zkRoot
	}

	if ret := strings.HasSuffix(zkRoot, "/"); ret {
		zkRoot = strings.TrimSuffix(zkRoot, "/")
	}

	path, data, err := initServerConfigData(zkHosts, zkRoot)
	if err != nil {
		fmt.Errorf("init server config failure, %s", err)
		return
	}
	fmt.Printf("zookeeper path: %s\n", path)
	fmt.Printf("data: %s\n", string(data))
	fmt.Printf("init to zookeeper successed!\n")
}
