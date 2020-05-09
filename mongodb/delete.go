package mongodb

import (
	"context"
)

func (a *AllCollection)DeleteBlock(ctx context.Context, blockNumber int64) (int64, error) {
	res, err := a.blocks.DeleteOne(ctx, map[string]int64{"number": blockNumber})
	if err != nil {
		return -2, err
	}
	return res.DeletedCount, err
}