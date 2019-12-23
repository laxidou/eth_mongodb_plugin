package main

import (
	"eth_mongodb_plugin/config"
	"eth_mongodb_plugin/data"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"context"
)

func main() {
    v, err := config.Init()
    if err != nil {
       panic(err)
    }
	//ethIp := fmt.Sprintf("http://%s:%s",v.Get("ETH.host"),v.Get("ETH.port"))
	ethIp := fmt.Sprintf("http://%s:%s",v.Get("localETH.host"),v.Get("localETH.port"))
    mongoIp := fmt.Sprintf("mongodb://%s:%s",v.Get("database.mongodb.host"),v.Get("database.mongodb.port"))
	fmt.Println(mongoIp)
	mobileCli, _ := data.NewEthMobile(ethIp)
	blockInfo, err := mobileCli.GetBlock(-1)

    //插入mongodb
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoIp))
    db := client.Database("ethT")
	blocksDb := db.Collection("blocks")
    ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
    res, err := blocksDb.InsertOne(ctx, &blockInfo)
    if err != nil {
    	fmt.Println(err)
	}
    id := res.InsertedID
    fmt.Println(id)
}