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
	//获取最新区块号
	blockInfo, _, _, _ := mobileCli.GetBlock(-1)
	//从不可逆区块开始拉取
	blockNumber := blockInfo.Number - 8
	blocks := make(chan int64, 10)
	ctx := context.Background()
	go checkBlock(mong, blockNumber - 1, blocks)
	go func() {
		for {
			pullFromChannel(ctx, mong, mobileCli, blockNumber, blocks)
		}
	}()
	reversePull(mong, mobileCli, blockNumber)
}

func pullFromChannel(ctx context.Context, mong *mongodb.AllCollection, mobileCli *data.MobileClient, blockNumber int64, blocks chan int64) {
	getNumber := <- blocks
	//fmt.Println("chennel拉块:",getNumber)
	insertBlock(ctx, mong, mobileCli, getNumber)
}

func reversePull(mong *mongodb.AllCollection, mobileCli *data.MobileClient, blockNumber int64) {
	ctx := context.Background()
	insertRes := insertBlock(ctx, mong, mobileCli, blockNumber - 8)

	time.Sleep(time.Second)
	if !insertRes {
		fmt.Println("已插入最新块", blockNumber)
	}
	blockInfo, _, _, _ := mobileCli.GetBlock(-1)
	reversePull(mong, mobileCli, blockInfo.Number - 8)
}

func insertBlock(ctx context.Context, mong *mongodb.AllCollection, mobileCli *data.MobileClient, blockNumber int64) bool {
	blockInfo, receiptsArr, logsArr, err := mobileCli.GetBlock(blockNumber)
	res, err := mong.BlockStateSearch(ctx, blockNumber)
	info := mongodb.BlockState{}
	bson.Unmarshal(res, &info)
	if err != nil {
		mong.BlockStateInsert(ctx, blockNumber)
		//mong.BlockStateUpdate(ctx, blockNumber, 1)
	}else if info.BlockState == 0 {
		//mong.BlockStateUpdate(ctx, blockNumber, 1)
	}else if info.BlockState == 1 {
		//mong.DeleteBlock(ctx, blockNumber)
	}else if info.BlockState == 2 {
		return false
	}
	mong.BlockInsert(ctx, blockInfo)
	mong.ReceiptsInsert(ctx, receiptsArr)
	mong.LogsInsert(ctx, logsArr)
	fmt.Println("插入第", blockNumber)
	mong.BlockStateUpdate(ctx, blockNumber, 2)
	return true
}

func checkBlock(mong *mongodb.AllCollection, blockNumber int64, blocks chan int64){
	for {
		ctx := context.Background()
		fmt.Println("channel:",len(blocks))
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