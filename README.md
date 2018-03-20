# Cloudtask InitConfig

The cloudtask platform initialize configuration tool.   

It sets the system common configuration parameters to the zookeeper path.

cloudtask zookeeper configuration path is `/cloudtask/ServerConfig`

### Dependencies

- [Zookeeper cluster 3.4.6+](https://zookeeper.apache.org)   

- [Mongodb cluster 3.0.12+](https://www.mongodb.com)  
  

### Usage

> binary
``` bash
$ ./cloudtask-initconfig -f ./ServerConfig.json
```

> docker image
``` bash
$ docker run -it --rm \
  -v /opt/cloudtask/ServerConfig.json:/cloudtask-initconfig/ServerConfig.json \
  --name=cloudtask-initconfig \
  cloudtask/cloudtask-initconfig:1.0.0
```


### Output Successed
``` bash
2018/03/16 15:07:29 Connected to 192.168.2.80:2181
2018/03/16 15:07:29 Authenticated: id=99692315792834560, timeout=15000
2018/03/16 15:07:29 Re-submitting `0` credentials after reconnect
zookeeper path: /cloudtask/ServerConfig
2018/03/16 15:07:30 Recv loop terminated: err=EOF
serverconfig: {"websitehost":"192.168.2.80:8091","centerhost":"192.168.2.80:8985","storagedriver":{"mongo":{"auth":{"password":"ds4dev","user":"datastoreAdmin"},"database":"cloudtask","hosts":"192.168.2.80:27017,192.168.2.81:27017,192.168.2.82:27017","options":["maxPoolSize=20","replicaSet=mgoCluster","authSource=admin"]}}}

2018/03/16 15:07:30 Send loop terminated: err=<nil>
initconfig to zookeeper successed!
```

### Checking Zookeeper OK?
``` bash
$ ./zkCli.sh -server 192.168.2.80
[zk: 192.168.2.80(CONNECTED) 0] get /cloudtask/ServerConfig
{"websitehost":"192.168.2.80:8091","centerhost":"192.168.2.80:8985","storagedriver":{"mongo":{"auth":{"password":"ds4dev","user":"datastoreAdmin"},"database":"cloudtask","hosts":"192.168.2.80:27017,192.168.2.81:27017,192.168.2.82:27017","options":["maxPoolSize=20","replicaSet=mgoCluster","authSource=admin"]}}}
```

### ServerConfig.json

``` json
{
    "zookeeper": {
        "hosts": "192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181",
        "root": "/cloudtask"
    },
    "serverconfig": {
        "websitehost": "http://192.168.2.80:8091",
        "centerhost": "http://192.168.2.80:8985",
        "storagedriver": {
            "mongo": {
                "hosts": "192.168.2.80:27017,192.168.2.81:27017,192.168.2.82:27017",
                "database": "cloudtask",
                "auth": {
                    "user": "datastoreAdmin",
                    "password": "ds4dev"
                },
                "options": [
                    "maxPoolSize=20",
                    "replicaSet=mgoCluster",
                    "authSource=admin"
                ]
            }
        }
    }
}
```

- zookeeper   
  `hosts`: set the zookeeper cluster hosts address to ensure that the state is running.   
  `root`: cloudtask zookeeper root path.   
- serverconfig   
  `websitehost`: cloudtask-web http address.   
  `centerhost`: cloudtask-center scheduler http address.      
  `storagedriver`: mongodb drvier cluster configs, currently only supports mongodb database, `mongo` key is `Required`.    
  `mongo.hosts`: set the mongodb cluster hosts address.   
  `mongo.database`: cloudtask database name.   
  `mongo.auth`: mongodb database safety certificate, if no security certificate, please ignorable 'auth' key.   
  `mongo.options`: mongodb cluster more k/v pair options, please see https://docs.mongodb.com/manual/reference/connection-string/#connections-connection-options 


### Example Mongo Driver No Security Setting

``` json
{
    "zookeeper": {
        "hosts": "192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181",
        "root": "/cloudtask"
    },
    "serverconfig": {
        "websitehost": "http://192.168.2.80:8091",
        "centerhost": "http://192.168.2.80:8985",
        "storagedriver": {
            "mongo": {
                "hosts": "192.168.2.80:27017,192.168.2.81:27017,192.168.2.82:27017",
                "database": "cloudtask",
                "options": [
                    "maxPoolSize=20",
                    "replicaSet=mgoCluster",
                    "authSource=admin"
                ]
            }
        }
    }
}
```
### License
cloudtask source code is licensed under the [Apache Licence 2.0](http://www.apache.org/licenses/LICENSE-2.0.html). 