package types

import (
	"encoding/json"
	"time"
)

/**
 * Copyright (C), 2019-2020
 * FileName: types
 * Author:   LinGuanHong
 * Date:     2020/12/9 5:59 下午
 * Description:
 */

type DelayControl struct {
	RoundDelay time.Duration // 扫描的间隔时间
	CatchDelay time.Duration // 追块的间隔时间
}

type ScannerBlockInfo struct {
	BlockHash   string      `json:"block_hash"`
	ParentHash  string      `json:"parent_hash"`
	Version     int32       `json:"version"`
	BlockNumber uint64      `json:"block_number"`
	Timestamp   string      `json:"timestamp"`
	Txs         interface{} `json:"txs"`
	TxCount     *int        `json:"tx_count"`
}

func (s *ScannerBlockInfo) IsEmpty() bool {
	return s.BlockHash == ""
}

type BlockForkInfo struct {
	BlockFrom *Block `json:"block_from"`
	BlockEnd  *Block `json:"block_end"`
}

func (b *BlockForkInfo) ToBytes() []byte {
	bys,_ := json.Marshal(b)
	return bys
}

type Block struct {
	Id          int64  `json:"id,omitempty"`
	BlockNumber uint64 `json:"block_number"`   // 区块号
	BlockHash   string `json:"block_hash"`  // 区块 hash
	ParentHash  string `json:"parent_hash"` // 父区块 hash
	Version     int32  `json:"version"`
	CreateTime  int64  `json:"create_time"` // 区块的生成时间
	Fork        bool   `json:"fork"`        // 是否是分叉区块
}

func (s *Block) IsEmpty() bool {
	return s.BlockHash == ""
}

func (b *Block) ToBytes() []byte {
	bys,_ := json.Marshal(b)
	return bys
}














