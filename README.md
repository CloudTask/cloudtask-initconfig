# Cloudtask InitConfig

The cloudtask platform initialize configuration tool.

### Dependencies

- [Zookeeper cluster 3.4.6+](https://zookeeper.apache.org)   

- [Mongodb cluster 3.0.12+](https://www.mongodb.com)  
  

### Usage

``` bash
$ ./cloudtask-initconfig -f ./ServerConfig.json
```

### ServerConfig.json

``` json
{
    "zookeeper": {
        "hosts": "192.168.2.80:2181,192.168.2.81:2181,192.168.2.82:2181",
        "root": "/cloudtask"
    },
    "serverconfig": {
        "websitehost": "192.168.2.80:8091",
        "centerhost": "192.168.2.80:8985",
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