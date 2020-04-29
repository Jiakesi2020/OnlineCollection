package main

import (
	"context"
	"flag"
	"oc/module"
	"oc/module/model"
	"oc/tools"
	"time"

	"github.com/golang/glog"
)

func main() {

	flag.Set("stderrthreshold", "0")
	flag.Parse()

	for {
		cfg := module.LoadConfig()

		err := _TimerHandler(cfg)
		if err != nil {
			glog.Errorln(err)
			continue
		}

		time.Sleep(time.Duration(cfg.Interval) * time.Minute)
	}
}

func _TimerHandler(cfg *module.Config) error {
	mc, err := tools.GetMongoInstance(cfg.MongoHost)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	defer mc.Disconnect(context.TODO())

	rc, err := tools.GetRedisInstance(cfg.RedisHosts)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	defer rc.Close()

	err = model.CheckOnlineUser(mc)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	dataList, err := model.GetOnlineUsersFromRedis(rc, cfg.Interval)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	err = model.SetOnlineUsersToMongo(mc, dataList)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	return nil
}
