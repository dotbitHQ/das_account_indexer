package celltype

import (
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: income_cell
 * Author:   LinGuanHong
 * Date:     2021/5/20 10:15
 * Description:
 */

var DefaultIncomeCellParam = func(data *IncomeCellData) *IncomeCellParam {
	return &IncomeCellParam{
		Version:        1,
		IncomeCellData: *data,
		CellCodeInfo:              DasIncomeCellScript,
		AlwaysSpendableScriptInfo: DasAnyOneCanSendCellInfo,
	}
}

type IncomeCell struct {
	p *IncomeCellParam
}

func NewIncomeCell(p *IncomeCellParam) *IncomeCell {
	return &IncomeCell{p: p}
}

func (c *IncomeCell) SoDeps() []types.CellDep {
	return nil
}

func (c *IncomeCell) LockDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.AlwaysSpendableScriptInfo.Dep.TxHash,
			Index:  c.p.AlwaysSpendableScriptInfo.Dep.TxIndex,
		},
		DepType: c.p.AlwaysSpendableScriptInfo.Dep.DepType,
	}
}
func (c *IncomeCell) TypeDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.CellCodeInfo.Dep.TxHash,
			Index:  c.p.CellCodeInfo.Dep.TxIndex,
		},
		DepType: c.p.CellCodeInfo.Dep.DepType,
	}
}
func (c *IncomeCell) LockScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.AlwaysSpendableScriptInfo.Out.CodeHash,
		HashType: c.p.AlwaysSpendableScriptInfo.Out.CodeHashType,
		Args:     c.p.AlwaysSpendableScriptInfo.Out.Args,
	}
}
func (c *IncomeCell) TypeScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
		HashType: c.p.CellCodeInfo.Out.CodeHashType,
	}
}

func (c *IncomeCell) Data() ([]byte, error) {
	rawBytes := c.p.IncomeCellData.AsSlice()
	tempByte := make([]byte,len(rawBytes))
	copy(tempByte,rawBytes)
	bys, err := blake2b.Blake256(tempByte)
	if err != nil {
		return nil, err
	}
	// fmt.Println("incomeCell's dataHash:",hex.EncodeToString(bys))
	return bys, nil
}

func (c *IncomeCell) TableType() TableType {
	return TableType_IncomeCell
}



