package data

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/mobile"
	"github.com/mitchellh/mapstructure"
	"fmt"
)

type BlockInfo struct {
	ParentHash string       "bson:`parentHash`"
	Sha3Uncles string       "bson:`sha3Uncles`"
	Miner string            "bson:`miner`"
	StateRoot string        "bson:`stateRoot`"
	TransactionsRoot string "bson:`transactionsRoot`"
	ReceiptsRoot string     "bson:`receiptsRoot`"
	LogsBloom string        "bson:`logsBloom`"
	Difficulty string       "bson:`difficulty`"
	Number int64            "bson:`number`"
	GasLimit int64          "bson:`gasLimit`"
	GasUsed int64           "bson:`gasUsed`"
	Timestamp string        "bson:`timestamp`"
	ExtraData string        "bson:`extraData`"
	MixHash string          "bson:`mixHash`"
	Nonce string            "bson:`nonce`"
	Hash string             "bson:`hash`"
	TotalTxs int			"bson:`totalTxs`"
	TotalNucles int          "bson:`totalNucle`"
	Uncles []UncleBlock
	Transactions []Txdata
}

type UncleBlock struct {
	ParentHash string       "bson:`parentHash`"
	Sha3Uncles string       "bson:`sha3Uncles`"
	Miner string            "bson:`miner`"
	StateRoot string        "bson:`stateRoot`"
	TransactionsRoot string "bson:`transactionsRoot`"
	ReceiptsRoot string     "bson:`receiptsRoot`"
	LogsBloom string        "bson:`logsBloom`"
	Difficulty string       "bson:`difficulty`"
	Number int64            "bson:`number`"
	GasLimit int64          "bson:`gasLimit`"
	GasUsed int64           "bson:`gasUsed`"
	Timestamp string        "bson:`timestamp`"
	ExtraData string        "bson:`extraData`"
	MixHash string          "bson:`mixHash`"
	Nonce string            "bson:`nonce`"
	Hash string             "bson:`hash`"
}

type Txdata struct {
	Nonce string			"bson:`nonce`"
	GasPrice string			"bson:`gasPrice`"
	Gas string				"bson:`gas`"
	To string				"bson:`to`"
	Value string			"bson:`value`"
	Input string			"bson:`input`"
	V string				"bson:`v`"
	R string				"bson:`r`"
	S string				"bson:`s`"
	Hash string				"bson:`hash`"
}



type MobileClient struct {
	cli *geth.EthereumClient
}

// NewEthereumClient connects a client to the given URL.
func NewEthMobile(ethIp string) (c *MobileClient, _ error) {
	cli, err := geth.NewEthereumClient(ethIp)
	return &MobileClient{cli}, err
	//client, err := ethclient.Dial(ethIp)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return &MobileClient{client}, err
}


//获取区块信息
func (m *MobileClient)GetBlock(blockNumber int64) (block *BlockInfo, err error) {
	eth_ctx := geth.NewContext()
	b, err := m.cli.GetBlockByNumber(eth_ctx, blockNumber)
	if err != nil {
		fmt.Println(err)
	}

	//获取区块头信息
	h := b.GetHeader()
	headerJson, _ := h.EncodeJSON()
	//获取叔块
	uncle := b.GetUncles()
	totalNucle :=uncle.Size()
	var uncles []UncleBlock
	for i := 0; i < totalNucle; i++{
		un, _ := uncle.Get(i)
		unRes, _ :=un.EncodeJSON()
		var uncleHeader map[string]interface{}
		if err := json.Unmarshal([]byte(unRes), &uncleHeader); err != nil{
			return nil, err
		}
		uncleHeader["gasLimit"], _ = math.ParseUint64(uncleHeader["gasLimit"].(string))
		uncleHeader["gasUsed"], _ = math.ParseUint64(uncleHeader["gasUsed"].(string))
		uncleHeader["number"], _ = math.ParseUint64(uncleHeader["number"].(string))
		var uncleInfo UncleBlock
		if err := mapstructure.Decode(uncleHeader, &uncleInfo); err != nil{
			return nil, err
		}
		uncles = append(uncles, uncleInfo)
	}

	//获取交易记录
	txs := b.GetTransactions()
	fmt.Println("total transactions:",txs.Size())
	totalTxs := txs.Size()
	var transactions []Txdata
	for i := 0; i < totalTxs; i++{
		var txdata Txdata
		tx, _ := txs.Get(i)


		tx_res, _ :=tx.EncodeJSON()
		//fmt.Println(tx_res)
		if err := json.Unmarshal([]byte(tx_res), &txdata); err != nil{
			return nil, err
		}
		transactions = append(transactions, txdata)
	}

	//类型转换
	var header map[string]interface{}
	if err := json.Unmarshal([]byte(headerJson), &header); err != nil{
		return nil, err
	}
	header["gasLimit"], _ = math.ParseUint64(header["gasLimit"].(string))
	header["gasUsed"], _ = math.ParseUint64(header["gasUsed"].(string))
	header["number"], _ = math.ParseUint64(header["number"].(string))

	////完善区块数据
	var blockInfo BlockInfo
	if err := mapstructure.Decode(header, &blockInfo); err != nil{
		return nil, err
	}
	blockInfo.TotalTxs = totalTxs
	blockInfo.Transactions = transactions
	blockInfo.TotalNucles = totalNucle
	blockInfo.Uncles = uncles
	return &blockInfo, nil
}

