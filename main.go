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
	for _, worker := range conf.Workers {
		go controller.StartWorker(conf.CronSpec, worker.Url, worker.Name, worker.HashFile)
	}

	select {}
}
