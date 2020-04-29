package model

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	_MongoDatabase   = "gohallserver"
	_MongoCollection = "OnlineUserRecords"
	_RedisKey        = "OnlineUser"
)

type OnlineUser struct {
	ID        string `json:"id"         bson:"_id"`
	Username  string `json:"username"  	bson:"username"`
	UserIdx   int64  `json:"useridx"    bson:"useridx"`
	Platform  int64  `json:"platform"   bson:"platform"`
	GameID    int64  `json:"gameid"     bson:"gameid"`
	LoginTime int64  `json:"login_time" bson:"login_time"`
	OrderID   string `json:"orderid"    bson:"orderid"`
	AgentCode string `json:"agent_code" bson:"agent_code"`
}

func (this_ *OnlineUser) String() string {
	data, _ := json.Marshal(this_)
	return string(data)
}

func _Init(mc *mongo.Client) error {
	ctx := context.TODO()
	indexView := mc.Database(_MongoDatabase).Collection(_MongoCollection).Indexes()
	cur, err := indexView.List(ctx)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	defer cur.Close(ctx)
	result := []bson.M{}
	err = cur.All(ctx, &result)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	if len(result) == 0 {
		_, err := indexView.CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.M{"login_time": 1},
			},
		})
		if err != nil {
			glog.Errorln(err)
			return err
		}
	}

	return nil
}

func CheckOnlineUser(mc *mongo.Client) error {
	return _Init(mc)
}

func GetOnlineUsersFromRedis(rc *redis.Client, interval int64) ([]interface{}, error) {
	dataMap, err := rc.HGetAll(_RedisKey).Result()
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	dataList := []interface{}{}
	for _, v := range dataMap {
		item := &OnlineUser{}

		err = json.Unmarshal([]byte(v), item)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}

		item.ID = primitive.NewObjectID().Hex()
		item.LoginTime = item.LoginTime - item.LoginTime%(interval*60)
		dataList = append(dataList, item)
	}

	return dataList, nil
}

func SetOnlineUsersToMongo(mc *mongo.Client, dataList []interface{}) error {
	col := mc.Database(_MongoDatabase).Collection(_MongoCollection)
	_, err := col.InsertMany(context.TODO(), dataList)
	if err != nil {
		glog.Errorln(err)
		return err
	}

	return nil
}
