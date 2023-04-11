package main

import (
	"github.com/go-ini/ini"
	"log"
)

var (
	configFileName = "config.conf" // 配置文件名
	config         = &Config{Host: "127.0.0.1", Port: 22, Username: "root"}
)

// Config 配置结构体
type Config struct {
	Host     string `ini:"host" comment:"地址"`
	Port     int    `ini:"port" comment:"端口"`
	Username string `ini:"username" comment:"用户名"`
	Password string `ini:"password" comment:"密码"`
	Cmd      string `ini:"cmd" comment:"命令"`
}

// 初始化配置文件
func initConfig() {
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, config)
	if err != nil {
		log.Fatalln(err)
	}
	err = cfg.SaveTo(configFileName)
	if err != nil {
		log.Fatalln(err)
	}
}

// 加载配置文件
func loadConfig() {
	cfg, err := ini.Load(configFileName)
	if err != nil {
		//log.Printf("读取配置文件失败: %v", err)
		return
	}
	err = cfg.MapTo(config)
	if err != nil {
		log.Println(err)
	}
}
