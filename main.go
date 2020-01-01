package main

import (
	"context"
	"eth_mongodb_plugin/config"
	"eth_mongodb_plugin/data"
	"eth_mongodb_plugin/mongodb"
	"fmt"
	"time"
)


func main() {
    v, err := config.Init()
    if err != nil {
      panic(err)
    }
	//ethIp := fmt.Sprintf("http://%s:%s",v.Get("ETH.host"),v.Get("ETH.port"))
	ethIp := fmt.Sprintf("http://%s:%s",v.Get("localETH.host"),v.Get("localETH.port"))
    mongoIp := fmt.Sprintf("mongodb://%s:%s",v.Get("database.labMongodb.host"),v.Get("database.labMongodb.port"))
	fmt.Println(mongoIp)
    mong, err := mongodb.NewCollection(mongoIp,"eth")
	mong.BlockIndex()
    mong.LogIndex()
    mong.ReceiptIndex()

	mobileCli, _ := data.NewEthMobile(ethIp)

	for {
		blockInfo, receiptsArr, logsArr, err := mobileCli.GetBlock(-1)
		ctx := context.Background()
		res, err := mong.BlockSearch(ctx, blockInfo.Number)
		if err != nil {
			mong.BlockInsert(ctx, blockInfo)
			mong.ReceiptsInsert(ctx, receiptsArr)
			mong.LogsInsert(ctx, logsArr)
			fmt.Println("插入第", blockInfo.Number)
		}else {
			fmt.Println("已是最新块")
		}
		fmt.Println(res)
		time.NewTimer(time.Second * 2)
	}
}