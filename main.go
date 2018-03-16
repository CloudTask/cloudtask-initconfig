package main

import "github.com/samuel/go-zookeeper/zk"

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

//Configuration is exported
//zookeeper cluster parameters.
type Configuration struct {
	Zookeeper struct {
		Hosts string `json:"hosts"`
		Root  string `json:"root"`
	} `json:"zookeeper"`
	ServerConfig struct {
		WebsiteHost   string                 `json:"websitehost"`
		CenterHost    string                 `json:"centerhost"`
		StorageDriver map[string]interface{} `json:"storagedriver"`
	} `json:"serverconfig"`
}

var (
	zkFlags       = int32(0)
	zkACL         = zk.WorldACL(zk.PermAll)
	zkConnTimeout = time.Second * 15
)

func initServerConfigData(conf *Configuration) (string, []byte, error) {

	hosts := strings.Split(conf.Zookeeper.Hosts, ",")
	conn, event, err := zk.Connect(hosts, zkConnTimeout)
	if err != nil {
		return "", nil, err
	}

	defer conn.Close()
	<-event
	if ret, _, _ := conn.Exists(conf.Zookeeper.Root); !ret {
		if _, err := conn.Create(conf.Zookeeper.Root, []byte{}, zkFlags, zkACL); err != nil {
			return "", nil, fmt.Errorf("zookeeper root: %s failure", conf.Zookeeper.Root)
		}
	}

	serverConfigPath := conf.Zookeeper.Root + "/ServerConfig"
	ret, _, err := conn.Exists(serverConfigPath)
	if err != nil {
		return "", nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	if err = json.NewEncoder(buf).Encode(conf.ServerConfig); err != nil {
		return "", nil, err
	}

	data := buf.Bytes()
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

func readConfiguration(configFile string) (*Configuration, error) {

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	conf := &Configuration{}
	err = json.NewDecoder(bytes.NewBuffer(data)).Decode(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func main() {

	var configFile string
	flag.StringVar(&configFile, "f", "./ServerConfig.json", "server config file path.")
	flag.Parse()

	conf, err := readConfiguration(configFile)
	if err != nil {
		fmt.Printf("server config file invalid, %s", err)
		return
	}

	if ret := strings.HasPrefix(conf.Zookeeper.Root, "/"); !ret {
		conf.Zookeeper.Root = "/" + conf.Zookeeper.Root
	}

	if ret := strings.HasSuffix(conf.Zookeeper.Root, "/"); ret {
		conf.Zookeeper.Root = strings.TrimSuffix(conf.Zookeeper.Root, "/")
	}

	path, data, err := initServerConfigData(conf)
	if err != nil {
		fmt.Printf("init server config failure, %s", err)
		return
	}
	fmt.Printf("zookeeper path: %s\n", path)
	fmt.Printf("serverconfig: %s\n", string(data))
	fmt.Printf("initconfig to zookeeper successed!\n")
}
