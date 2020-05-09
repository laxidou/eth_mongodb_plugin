package mongodb

import (
	"context"
	"eth_mongodb_plugin/data"
)

type BlockState struct {
	BlockNumber int64
	BlockState int
}

func (a *AllCollection)BlockInsert(ctx context.Context,blockInfo *data.BlockInfo) (interface{}, error) {
	res, err := a.blocks.InsertOne(ctx, &blockInfo)
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}

func (a *AllCollection)ReceiptsInsert(ctx context.Context,receiptsArr *[]interface{}) ([]interface{}, error) {
	res, err := a.receipts.InsertMany(ctx, *receiptsArr)
	if err != nil {
		return nil, err
	}
	ids := res.InsertedIDs
	return ids, nil
}

func (a *AllCollection)LogsInsert(ctx context.Context,logsArr *[]interface{}) ([]interface{}, error) {
	res, err := a.logs.InsertMany(ctx, *logsArr)
	if err != nil {
		return nil, err
	}
	ids := res.InsertedIDs
	return ids, nil
}

func (a *AllCollection)BlockStateInsert(ctx context.Context,blockNumber int64) (interface{}, error) {
	res, err := a.blockState.InsertOne(ctx, &BlockState{blockNumber,0})
	if err != nil {
		return nil, err
	}
	id := res.InsertedID
	return id, nil
}