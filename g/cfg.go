package g

import (
	"encoding/json"
	"github.com/toolkits/file"
	"log"
	"sync"
)


type GlobalConfig struct {
	AccessKeyId     string  `json:"accessKeyId"`
	AccessKeySecret string  `json:"accessKeySecret"`
	Appid           string  `json:"appid"`
	Appsecret       string   `json:"appsecret"`
	Callback        string `json:"callback"`
	ServerPort      string `json:"server_port"`
	MysqlCon        string `json:"mysql_connect"`
}

var (
	config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

func ParseConfig(cfg string) {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		log.Fatalln("config file:", cfg, "is not existent")
	}


	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		log.Fatalln("read config file:", cfg, "fail:", err)
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
}
