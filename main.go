package main

import (
	"eth-data/config"
	"eth-data/data"
	"fmt"
)

func main() {
    v, err := config.Init()
    if err != nil {
       panic(err)
    }
	//ethIp := fmt.Sprintf("http://%s:%s",v.Get("ETH.host"),v.Get("ETH.port"))
	ethIp := fmt.Sprintf("http://%s:%s",v.Get("localETH.host"),v.Get("localETH.port"))
    //mongoIp := fmt.Sprintf("mongodb://%s:%s",v.Get("database.mongodb.host"),v.Get("database.mongodb.port"))

    ethCli, _ := data.NewEthClient(ethIp)
	ethCli.GetReceiptByTxHash("0x7507187ba1b15ed0271871766d9ab3213f140ed6eb5548c3f4dc41172c3b6e3f")

	//mobileCli, _ := data.NewEthMobile(ethIp)
	//blockInfo, err := mobileCli.GetBlock(9000000)
    //fmt.Println(blockInfo)
    //ethCli.GetReceiptByTxHash("0xa98e2844107c0377fd538455c1cf00a4780eda5547edb06ed55aacc2d6663551")

    //插入mongodb
    //ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    //client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoIp))
    //collection := client.Database("eth").Collection("blocks")
    //ctx, _ = context.WithTimeout(context.Background(), 5*time.Second)
    //res, _ := collection.InsertOne(ctx, &blockInfo)
    //id := res.InsertedID
    //fmt.Println(id)
}