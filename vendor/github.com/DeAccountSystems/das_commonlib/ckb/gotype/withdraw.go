package gotype

import (
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: withdraw
 * Author:   LinGuanHong
 * Date:     2021/6/24 11:06
 * Description:
 */

type WithdrawDasLockCell struct {
	OutPoint       *types.OutPoint
	LockScriptArgs []byte // das-lock'args
	CellCap        uint64
}

func (w WithdrawDasLockCell) LockType() celltype.LockScriptType {
	return celltype.DasLockCodeHashIndexType(w.LockScriptArgs[0]).ToScriptType(true)
}
