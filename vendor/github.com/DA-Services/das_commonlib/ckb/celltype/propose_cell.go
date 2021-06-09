package celltype

import (
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: proposalcell
 * Author:   LinGuanHong
 * Date:     2021/1/10 12:22
 * Description:
 */

var DefaultProposeCellParam = func(new *ProposalCellData) *ProposeCellParam {
	acp := &ProposeCellParam{
		Version:     1,
		TxDataParam: *new,
		CellCodeInfo:              DasProposeCellScript,
		AlwaysSpendableScriptInfo: DasAnyOneCanSendCellInfo,
	}
	return acp
}

type ProposeCell struct {
	p *ProposeCellParam
}

func NewProposeCell(p *ProposeCellParam) *ProposeCell {
	return &ProposeCell{p: p}
}
func (c *ProposeCell) SoDeps() []types.CellDep {
	return nil
}
func (c *ProposeCell) LockDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.AlwaysSpendableScriptInfo.Dep.TxHash,
			Index:  c.p.AlwaysSpendableScriptInfo.Dep.TxIndex,
		},
		DepType: c.p.AlwaysSpendableScriptInfo.Dep.DepType,
	}
}
func (c *ProposeCell) TypeDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.CellCodeInfo.Dep.TxHash,
			Index:  c.p.CellCodeInfo.Dep.TxIndex,
		},
		DepType: c.p.CellCodeInfo.Dep.DepType,
	}
}
func (c *ProposeCell) LockScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.AlwaysSpendableScriptInfo.Out.CodeHash,
		HashType: c.p.AlwaysSpendableScriptInfo.Out.CodeHashType,
		Args:     c.p.AlwaysSpendableScriptInfo.Out.Args,
	}
}
func (c *ProposeCell) TypeScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
		HashType: c.p.CellCodeInfo.Out.CodeHashType,
		Args:     nil,
	}
}

func (c *ProposeCell) Data() ([]byte, error) {
	hashBys, _ := blake2b.Blake256(c.p.TxDataParam.AsSlice())
	return hashBys, nil
}

func (c *ProposeCell) TableType() TableType {
	return TableType_PROPOSE_CELL
}
