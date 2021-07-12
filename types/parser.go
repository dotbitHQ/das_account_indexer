package types

import (
	blockparserTypes "github.com/af913337456/blockparser/types"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: parser
 * Author:   LinGuanHong
 * Date:     2021/7/10 4:53
 * Description:
 */

type ParserHandleBaseTxInfo struct {
	ScanInfo blockparserTypes.ScannerBlockInfo
	Tx       types.Transaction
	TxIndex  uint8
}
