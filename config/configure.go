package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

var (
	once sync.Once
	Conf DbConfig
)

type DbConfig struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Database    string `json:"database"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	AutoMigrate bool   `json:"auto_migrate"`
	LogMode     bool   `json:"log_mode"`
}

func Use(path string) {
	once.Do(func() {
		file, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal("配置文件读取失败: ", err)
		}
		err = json.Unmarshal(file, &Conf)
		if err != nil {
			log.Fatal("配置文件读取失败: ", err)
		}
	})
}
