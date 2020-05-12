package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

func (a *AllCollection)BlockStateUpdate(ctx context.Context,blockNumber int64, state int) (int64, error) {
	res, err := a.blockState.UpdateOne(ctx, &BlockState{BlockNumber:blockNumber}, bson.M{"$set": bson.M{"blockstate": state}})
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}
