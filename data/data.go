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
	TotalUncles int         "bson:`totalUncle`"
	Uncles []UncleBlock		"json:`group,omitempty` bson:`totalUncle`"
	Transactions []TxData	"json:`group,omitempty` bson:`transactions`"
	Receipts []ReceiptInfo  "json:`group,omitempty` bson:`receipts`"
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

type TxData struct {
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
	ip string
}

// NewEthereumClient connects a client to the given URL.
func NewEthMobile(ethIp string) (c *MobileClient, _ error) {
	cli, err := geth.NewEthereumClient(ethIp)
	return &MobileClient{cli,ethIp}, err
}

//获取区块信息
func (m *MobileClient)GetBlock(blockNumber int64) (block *BlockInfo, err error) {
	ethCtx := geth.NewContext()
	b, err := m.cli.GetBlockByNumber(ethCtx, blockNumber)
	if err != nil {
		fmt.Println(err)
	}
	//获取区块头信息
	h := b.GetHeader()
	blockInfo, err := getHeader(h)
	//获取叔块
	uncle := b.GetUncles()
	uncles, totalUncle, err :=getUncleBlock(uncle)
	//获取交易记录
	txs := b.GetTransactions()
	transactions, totalTxs, receipts, err:= getTransactions(txs)

	////完善区块数据
	blockInfo.TotalTxs = totalTxs
	blockInfo.Transactions = *transactions
	blockInfo.TotalUncles = totalUncle
	blockInfo.Uncles = *uncles
	ethCli, _ := newEthClient(m.ip)
	var receiptsArr = make([]ReceiptInfo, 0)
	for _, repHash := range *receipts{
		receiptsInfo := ethCli.GetReceiptByTxHash(repHash)
		receiptsArr = append(receiptsArr, *receiptsInfo)
	}
	blockInfo.Receipts = receiptsArr
	return blockInfo, nil
}

//获取叔块
func getUncleBlock(uncle *geth.Headers) (unclesInfo *[]UncleBlock, totalNucle int, err error) {
	totalNucle =uncle.Size()
	var uncles =make([]UncleBlock, 0)
	for i := 0; i < totalNucle; i++{
		un, _ := uncle.Get(i)
		unRes, _ :=un.EncodeJSON()
		var uncleHeader map[string]interface{}
		if err := json.Unmarshal([]byte(unRes), &uncleHeader); err != nil{
			return nil, totalNucle, err
		}
		uncleHeader["gasLimit"], _ = math.ParseUint64(uncleHeader["gasLimit"].(string))
		uncleHeader["gasUsed"], _ = math.ParseUint64(uncleHeader["gasUsed"].(string))
		uncleHeader["number"], _ = math.ParseUint64(uncleHeader["number"].(string))
		var uncleInfo UncleBlock
		if err := mapstructure.Decode(uncleHeader, &uncleInfo); err != nil{
			return nil, totalNucle, err
		}
		uncles = append(uncles, uncleInfo)
	}
	return &uncles, totalNucle, nil
}

//获取交易记录
func getTransactions(txs *geth.Transactions) (txsInfo *[]TxData, totalTxs int, receiapts *[]string, err error) {
	fmt.Println("total transactions:",txs.Size())
	totalTxs = txs.Size()
	var receiptArr []string
	var transactions = make([]TxData, 0)
	for i := 0; i < totalTxs; i++{
		var txData TxData
		tx, _ := txs.Get(i)
		txRes, _ :=tx.EncodeJSON()
		if err := json.Unmarshal([]byte(txRes), &txData); err != nil{
			return nil, totalTxs, nil,err
		}
		fmt.Println(txData)
		transactions = append(transactions, txData)
		receiptArr = append(receiptArr, tx.GetHash().GetHex())
	}
	return &transactions, totalTxs, &receiptArr,nil
}

func getHeader(h *geth.Header) (block *BlockInfo, err error) {
	headerJson, _ := h.EncodeJSON()
	//类型转换
	var header map[string]interface{}
	if err := json.Unmarshal([]byte(headerJson), &header); err != nil{
		return nil, err
	}
	header["gasLimit"], _ = math.ParseUint64(header["gasLimit"].(string))
	header["gasUsed"], _ = math.ParseUint64(header["gasUsed"].(string))
	header["number"], _ = math.ParseUint64(header["number"].(string))
	var blockInfo BlockInfo
	if err := mapstructure.Decode(header, &blockInfo); err != nil{
		return nil, err
	}
	return &blockInfo, nil
}