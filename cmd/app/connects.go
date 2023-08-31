package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/config"
)

func NewConnectToMongo(ctx context.Context, cfg *config.YAML) (*mongo.Database, error) {
	val := cfg.Get("mongodb")
	var (
		db       = val.Get("nameDB")
		user     = val.Get("user")
		password = val.Get("password")
		host     = val.Get("host")
		port     = val.Get("port")
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s", user, password, host, port)))
	if err != nil {
		return nil, fmt.Errorf("connect to mongodb %s:%s: %w", host, port, err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("failed attempt ping: %w", err)
	}
	return client.Database(db.String()), nil
}
