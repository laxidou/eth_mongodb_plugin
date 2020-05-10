package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type AllCollection struct {
	blocks 		*mongo.Collection
	receipts 	*mongo.Collection
	logs 		*mongo.Collection
	blockState 	*mongo.Collection
}

func NewCollection(mongoIp string, databaseName string) (a *AllCollection, _ error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoIp))
	if err != nil {
		return nil, err
	}
	db := client.Database(databaseName)
	blocksCollection := db.Collection("blocks")
	receiptsCollection := db.Collection("receipts")
	logsCollection := db.Collection("logs")
	blockStateCollection := db.Collection("blockState")
	return &AllCollection{blocksCollection,receiptsCollection,logsCollection, blockStateCollection}, nil
}

//block建立索引
func (a *AllCollection)BlockIndex() ([]string,error) {
	opt := options.IndexOptions{}
	opt.SetUnique(true)
	newIndexs := []mongo.IndexModel{
		{Keys: map[string]int{"number": 1}},
		{Keys: map[string]int{"hash": 1}},
		{Options: &opt},
	}
	index := a.blocks.Indexes()
	return createIndexs(&index, &newIndexs)
}

//receipt建立索引
func (a *AllCollection)ReceiptIndex() ([]string,error) {
	opt := options.IndexOptions{}
	opt.SetUnique(true)
	newIndexs := []mongo.IndexModel{
		{Keys: map[string]int{"txhash": -1}},
		{Keys: map[string]int{"blocknumber": -1}},
		{Options: &opt},
	}
	index := a.blocks.Indexes()
	return createIndexs(&index, &newIndexs)
}

//logs建立索引
func (a *AllCollection)LogIndex() ([]string,error) {
	opt := options.IndexOptions{}
	opt.SetUnique(true)
	newIndexs := []mongo.IndexModel{
		{Keys: map[string]int{"address": -1}},
		{Keys: map[string]int{"blocknumber": -1}},
		{Keys: map[string]int{"blockhash": -1}},
		{Keys: map[string]int{"txhash": -1}},
		{Options: &opt},
	}
	index := a.logs.Indexes()
	return createIndexs(&index, &newIndexs)
}

//BlockState建立索引
func (a *AllCollection)BlockStateIndex() ([]string,error) {
	opt := options.IndexOptions{}
	opt.SetUnique(true)
	newIndexs := []mongo.IndexModel{
		{Keys: map[string]int{"blocknumber": -1}},
		{Options: &opt},
	}
	index := a.blockState.Indexes()
	return createIndexs(&index, &newIndexs)
}

func createIndexs(index *mongo.IndexView, newIndexs *[]mongo.IndexModel) ([]string,error) {
	ctx := context.Background()
	res, err := index.CreateMany(ctx, *newIndexs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(res)
	return res, nil
}