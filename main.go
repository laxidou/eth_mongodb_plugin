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
	config.Execute()
	config.EmpApp.EmpSetting()
	mong, err := mongodb.NewCollection(config.EmpApp.MongoDBIp,config.EmpApp.DatabaseName)
	if err != nil {
		panic(err)
	}
   	if config.EmpApp.CreateIndex == true{
		mong.BlockIndex()
		mong.LogIndex()
		mong.ReceiptIndex()
	}

	mobileCli, _ := data.NewEthMobile(config.EmpApp.EthIp)
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