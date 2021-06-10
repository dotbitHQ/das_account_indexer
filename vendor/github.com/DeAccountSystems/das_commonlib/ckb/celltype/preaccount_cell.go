package celltype

import (
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2020
 * FileName: statecell
 * Author:   LinGuanHong
 * Date:     2020/12/18 3:58
 * Description:
 */

var DefaultPreAccountCellParam = func(account DasAccount, new *PreAccountCellData) *PreAccountCellParam {
	return &PreAccountCellParam{
		Version: 1,
		Account: account,
		// Data:         *buildDasCommonMoleculeDataObj(depIndex, oldIndex, newIndex, dep, old, new),
		CellCodeInfo: DasPreAccountCellScript,
		TxDataParam: PreAccountCellTxDataParam{
			NewAccountCellData: new,
		},
		AlwaysSpendableScriptInfo: DasAnyOneCanSendCellInfo,
	}
}

type PreAccountCell struct {
	p *PreAccountCellParam
}

func NewPreAccountCell(p *PreAccountCellParam) *PreAccountCell {
	return &PreAccountCell{p: p}
}
func (c *PreAccountCell) SoDeps() []types.CellDep {
	return nil
}
func (c *PreAccountCell) LockDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.AlwaysSpendableScriptInfo.Dep.TxHash,
			Index:  c.p.AlwaysSpendableScriptInfo.Dep.TxIndex,
		},
		DepType: c.p.AlwaysSpendableScriptInfo.Dep.DepType,
	}
}
func (c *PreAccountCell) TypeDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.CellCodeInfo.Dep.TxHash,
			Index:  c.p.CellCodeInfo.Dep.TxIndex,
		},
		DepType: c.p.CellCodeInfo.Dep.DepType,
	}
}
func (c *PreAccountCell) LockScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.AlwaysSpendableScriptInfo.Out.CodeHash,
		HashType: c.p.AlwaysSpendableScriptInfo.Out.CodeHashType,
		Args:     c.p.AlwaysSpendableScriptInfo.Out.Args,
	}
}
func (c *PreAccountCell) TypeScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
		HashType: c.p.CellCodeInfo.Out.CodeHashType,
		Args:     c.p.CellCodeInfo.Out.Args,
	}
}

func (c *PreAccountCell) TableType() TableType {
	return TableType_PRE_ACCOUNT_CELL
}

func (c *PreAccountCell) Data() ([]byte, error) {
	dataHash, err := blake2b.Blake256(c.p.TxDataParam.NewAccountCellData.AsSlice())
	if err != nil {
		return nil, err
	}
	accountId := c.p.Account.AccountId()
	return append(dataHash, accountId.Bytes()...), nil
}
