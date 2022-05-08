package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v2"
)

var globalConf atomic.Value

// Conf 配置
type Conf struct {
	AppID       string      `yaml:"appid"`
	Secret      string      `yaml:"secret"`
	Addr        string      `yaml:"addr"`
	AccessToken AccessToken `yaml:"access_token"`
	Token       string      `yaml:"token"` // 配置服务器使用的 token
}

type AccessToken struct {
	DurationInSec int32  `yaml:"duration_in_sec"`
	TimeoutInSec  int32  `yaml:"timeout_in_sec"`
	URL           string `yaml:"url"`
	Retries       int32  `yaml:"retries"`
}

// InitConf 初始化配置, name 是配置文件名
func InitConf(ctx context.Context, name string) error {
	log.Printf("start init conf:%v\n", name)
	conf, err := getConf(name)
	if err != nil {
		return err
	}
	globalConf.Store(conf)
	go updateConf(ctx, name)
	log.Printf("end init confname:%v, conf:%+v\n", name, GetConf())
	return nil
}

func getConf(name string) (*Conf, error) {
	body, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("read all conf file:%v err:%v", name, err)
	}
	conf := &Conf{}
	if err = yaml.Unmarshal(body, conf); err != nil {
		return nil, fmt.Errorf("unmarshal conf err:%v, body:%s", err, body)
	}
	return conf, nil
}

func updateConf(ctx context.Context, name string) {
	log.Printf("start updateConf:%v\n", name)
	tick := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Println("updateConf ctx done")
		case <-tick.C:
			conf, err := getConf(name)
			if err != nil {
				log.Printf("ERROR get conf err:%v, name:%v\n", err, name)
				continue
			}
			if !reflect.DeepEqual(conf, GetConf()) {
				log.Printf("get new conf:%+v\n", conf)
			}
			globalConf.Store(conf)
		}
	}
}

// GetConf 获取配置
func GetConf() *Conf {
	return globalConf.Load().(*Conf)
}
