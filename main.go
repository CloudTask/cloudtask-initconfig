package main

import "github.com/samuel/go-zookeeper/zk"

import (
	"bytes"
	"encoding/json"
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
		WebsiteHost string                 `json:"websitehost"`
		CenterHost  string                 `json:"centerhost"`
		Storage     map[string]interface{} `json:"storage"`
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

	<-event
	defer conn.Close()
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

func readConfiguration() (*Configuration, error) {

	data, err := ioutil.ReadFile("./ServerConfig.json")
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

	conf, err := readConfiguration()
	if err != nil {
		fmt.Errorf("ServerConfig.json invalid, %s", err)
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
		fmt.Errorf("init server config failure, %s", err)
		return
	}
	fmt.Printf("zookeeper path: %s\n", path)
	fmt.Printf("data: %s\n", string(data))
	fmt.Printf("init to zookeeper successed!\n")
}
