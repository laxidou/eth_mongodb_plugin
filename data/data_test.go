package data

import (
	"testing"
)

func TestGetBlock(t *testing.T) {
	mobileCli, _ := NewEthMobile("192.168.1.100:8080")
	mobileCli.GetBlock(9148655)
	//blockInfo, receiptsArr, logsArr, err := mobileCli.GetBlock(9148655)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(blockInfo)
	//fmt.Println(receiptsArr)
	//fmt.Println(logsArr)
}