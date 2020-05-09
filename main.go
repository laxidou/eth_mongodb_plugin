package main

import (
	"context"
	"eth_mongodb_plugin/config"
	"eth_mongodb_plugin/data"
	"eth_mongodb_plugin/mongodb"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)


func main() {
	config.Execute()
	config.EmpApp.EmpSetting()
	mong, err := mongodb.NewCollection(config.EmpApp.MongoDBIp,config.EmpApp.DatabaseName)
	if err != nil {
		panic(err)
	}
   	if config.EmpApp.CreateIndex == true {
		mong.BlockIndex()
		mong.LogIndex()
		mong.ReceiptIndex()
		mong.BlockStateIndex()
	}

	mobileCli, _ := data.NewEthMobile(config.EmpApp.EthIp)
	blockInfo, _, _, _ := mobileCli.GetBlock(-1)
	blocks := make(chan int64, 10)
	blockNumber := blockInfo.Number - 8
	ctx := context.Background()
	go checkBlock(mong, blockNumber, blocks)
	go func() {
		for {
			getNumber := <- blocks
			go insertBlock(ctx, mong, mobileCli, getNumber)
		}
	}()
	reversePull(mong, mobileCli, blockNumber, true)
}

func reversePull(mong *mongodb.AllCollection, mobileCli *data.MobileClient, blockNumber int64, startFlag bool) {
	var insertRes bool
	ctx := context.Background()
	if startFlag == true {
		insertRes = insertBlock(ctx, mong, mobileCli, blockNumber)
	}else{
		insertRes = insertBlock(ctx, mong, mobileCli, blockNumber - 8)
	}
	if insertRes {
		fmt.Println("已插入最新块", blockNumber)
	}
	blockInfo, _, _, _ := mobileCli.GetBlock(-1)
	time.Sleep(time.Second)
	reversePull(mong, mobileCli, blockInfo.Number, false)
}

func insertBlock(ctx context.Context, mong *mongodb.AllCollection, mobileCli *data.MobileClient, blockNumber int64) bool {
	blockInfo, receiptsArr, logsArr, err := mobileCli.GetBlock(blockNumber)
	res, err := mong.BlockStateSearch(ctx, blockNumber)
	info := mongodb.BlockState{}
	bson.Unmarshal(res, &info)
	if err != nil {
		mong.BlockStateInsert(ctx, blockNumber)
		mong.BlockStateUpdate(ctx, blockNumber, 1)
	}else if info.BlockState == 0 {
		mong.BlockStateUpdate(ctx, blockNumber, 1)
	}else if info.BlockState == 1 {
		mong.DeleteBlock(ctx, blockNumber)
	}else if info.BlockState == 2 {
		return false
	}
	mong.BlockInsert(ctx, blockInfo)
	mong.ReceiptsInsert(ctx, receiptsArr)
	mong.LogsInsert(ctx, logsArr)
	fmt.Println("插入第", blockNumber)
	fmt.Println(res)
	mong.BlockStateUpdate(ctx, blockNumber, 2)
	return true
}

func checkBlock(mong *mongodb.AllCollection, blockNumber int64, blocks chan int64){
	for {
		ctx := context.Background()
		if blockNumber == 0 {
			close(blocks)
		} else {
			if len(blocks) < 10 {
				res, err := mong.BlockStateSearch(ctx, blockNumber)
				if err != nil {
					mong.BlockStateInsert(ctx, blockNumber)
				} else {
					info := mongodb.BlockState{}
					bson.Unmarshal(res, &info)
					if info.BlockState == 2 {
						blockNumber--
						continue
					} else if info.BlockState == 1 {
						deleteRes, deleteErr := mong.DeleteBlock(ctx, blockNumber)
						if deleteErr != nil {
							fmt.Println("删除脏数据错误", deleteRes)
						} else {
							fmt.Println("删除脏数据", deleteRes)
						}
						mong.BlockStateUpdate(ctx, blockNumber, 0)
					}
				}
				blocks <- blockNumber
			}
		}
		blockNumber--
	}
}