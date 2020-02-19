package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Config struct {
	CronSpec     string
	AppID        int
	AppKey       string
	RegisteTemp  int
	ExpireMinute string
	Phone        string
	Workers      []Worker
}

type Worker struct {
	Url      string
	Name     string
	HashFile string
}

var (
	c      *Config
	locker sync.Mutex
)

func GetConfig() (*Config, error) {
	locker.Lock()
	defer locker.Unlock()
	if c != nil {
		log.Println("use old config")
		return c, nil
	}
	var err error
	v := viper.New()

	v.SetConfigFile("config/config.yaml")
	err = v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return c, err
	}
	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal(err)
		return c, err
	}
	log.Println(c)
	return c, nil
}
