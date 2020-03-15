package mongodb

import (
	"fmt"
	"testing"
)

func TestNewCollection(t *testing.T) {
	mong, err := NewCollection("mongodb://127.0.0.1:27017","ethT")
	if err != nil {
		fmt.Println(err)
	}
	mong.BlockIndex()
	mong.LogIndex()
	mong.ReceiptIndex()
}
