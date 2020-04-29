package tools

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_MaxSizePoolSize       = 10
	_MongoConnectedTimeout = time.Second * 15
)

func GetMongoInstance(host string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _MongoConnectedTimeout)
	defer cancel()

	return mongo.Connect(ctx, options.Client().
		ApplyURI(fmt.Sprintf("mongodb://%s", host)).
		SetMaxPoolSize(_MaxSizePoolSize))
}
