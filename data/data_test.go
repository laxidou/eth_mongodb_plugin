package data

import (
	"fmt"
	"testing"
)

const ip = "http://47.244.53.77:8545"
const blockNumber = 9148655

func TestNewEthMobile(t *testing.T)  {
	ethIp := fmt.Sprintf(ip)
	mobileCli, _ := NewEthMobile(ethIp)
	if mobileCli == nil {
		t.Fatal("Couldn't create EthMobile!")
	}
	t.Log("TestNewEthMobile")
}

func TestGetBlock(t *testing.T) {
	ethIp := fmt.Sprintf(ip)
	mobileCli, _ := NewEthMobile(ethIp)
	blockInfo, receiptsArr, logsArr, err := mobileCli.GetBlock(blockNumber)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(blockInfo)
	fmt.Println(receiptsArr)
	fmt.Println(logsArr)
}