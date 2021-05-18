package celltype

import "github.com/nervosnetwork/ckb-sdk-go/types"

/**
 * Copyright (C), 2019-2020
 * FileName: icelltype
 * Author:   LinGuanHong
 * Date:     2020/12/20 3:25 下午
 * Description:
 */

type ICellType interface {
	SoDeps() []types.CellDep
	LockDepCell() *types.CellDep
	TypeDepCell() *types.CellDep
	LockScript() *types.Script
	TypeScript() *types.Script
	TableType() TableType
	Data() ([]byte, error)
	// TableData() []byte
}

type ICellData interface {
	AsSlice() []byte
}
