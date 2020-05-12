package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *AllCollection)BlockSearch(ctx context.Context, blockNumber int64) (bson.Raw, error) {
	res := a.blocks.FindOne(ctx, map[string]int64{"number": blockNumber})
	blockInfo, err := res.DecodeBytes()
	if err != nil {
		return nil, err
	}
	return blockInfo, nil
}

func (a *AllCollection)BlockStateSearch(ctx context.Context, blockNumber int64) (bson.Raw, error) {
	res := a.blockState.FindOne(ctx, map[string]int64{"blocknumber": blockNumber})
	StateInfo, err := res.DecodeBytes()
	if err != nil {
		return nil, err
	}
	return StateInfo, nil
}