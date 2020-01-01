package mongodb

import (
	"fmt"
	"testing"
)

func TestNewCollection(t *testing.T) {
	mong, err := NewCollection("192.168.1.100:27017","eth")
	if err != nil {
		fmt.Println(err)
	}
	mong.BlockIndex()
	mong.LogIndex()
	mong.ReceiptIndex()
}
