package ckb_rocksdb_parser

import (
	blockparserTypes "github.com/af913337456/blockparser/types"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: types
 * Author:   LinGuanHong
 * Date:     2021/7/9 5:08
 * Description:
 */

type TxMsgData struct {
	BlockBaseInfo blockparserTypes.ScannerBlockInfo `json:"block_base_info"`
	Txs           []*types.Transaction              `json:"block_tx_list"`
}