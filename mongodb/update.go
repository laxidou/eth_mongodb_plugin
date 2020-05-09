package mongodb

import (
	"context"
)

func (a *AllCollection)BlockStateUpdate(ctx context.Context,blockNumber int64, state int) (interface{}, error) {
	res, err := a.blockState.UpdateOne(ctx, map[string]int64{"blocknumber": blockNumber}, map[string]int{"blockstate": state})
	if err != nil {
		return nil, err
	}
	return res, nil
}
