package module

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
)

const (
	_ConfigPath = "./config.json"
)

type Config struct {
	Interval   int64    `json:"interval"`
	RedisHosts []string `json:"redis_hosts"`
	MongoHost  string   `json:"mongo_host"`
}

func (this_ *Config) String() string {
	jstr, _ := json.Marshal(this_)
	return string(jstr)
}

func LoadConfig() *Config {
	f, err := os.OpenFile(_ConfigPath, os.O_RDONLY, 0)
	if err != nil {
		glog.Errorln(err)
		os.Exit(101)
	}

	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		glog.Errorln(err)
		os.Exit(102)
	}

	cfg := &Config{}
	err = json.Unmarshal(data, cfg)
	if err != nil {
		glog.Errorln(err)
		os.Exit(103)
	}

	return cfg
}
