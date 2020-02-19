package main

import (
	"go-crapy/config"
	"go-crapy/controller"
)

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		panic(err)
	}
	go controller.StartWorker(conf.CronSpec, "http://81rc.81.cn/index.htm", "军队人才网首页", "persistence/hash1")
	//go StartWorker(conf.CronSpec, "http://81rc.81.cn/Civilianpost/index.htm", "军队人才网文职人员")

	select {}
}
