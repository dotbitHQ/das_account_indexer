package celltype

import (
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: on_sale_cell
 * Author:   LinGuanHong
 * Date:     2021/2/22 11:05
 * Description:
 */

var DefaultOnSaleCellParam = func(newIndex uint32, price uint64, accountId DasAccountId) *OnSaleCellParam {
	return &OnSaleCellParam{
		Version:        1,
		Price:          price,
		OnSaleCellData: NewOnSaleCellDataBuilder().Price(GoUint64ToMoleculeU64(price)).Build(),
		// Data: *buildDasCommonMoleculeDataObj(0, 0, newIndex, nil, nil, &onSaleMoleData),
		AccountId:                 accountId,
		CellCodeInfo:              DasOnSaleCellScript,
		AlwaysSpendableScriptInfo: DasAnyOneCanSendCellInfo,
	}
}

type OnSaleCell struct {
	p *OnSaleCellParam
}

func NewOnSaleCell(p *OnSaleCellParam) *OnSaleCell {
	return &OnSaleCell{p: p}
}

func (c *OnSaleCell) SoDeps() []types.CellDep {
	return nil
}

func (c *OnSaleCell) LockDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.AlwaysSpendableScriptInfo.Dep.TxHash,
			Index:  c.p.AlwaysSpendableScriptInfo.Dep.TxIndex,
		},
		DepType: c.p.AlwaysSpendableScriptInfo.Dep.DepType,
	}
}
func (c *OnSaleCell) TypeDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.CellCodeInfo.Dep.TxHash,
			Index:  c.p.CellCodeInfo.Dep.TxIndex,
		},
		DepType: c.p.CellCodeInfo.Dep.DepType,
	}
}
func (c *OnSaleCell) LockScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.AlwaysSpendableScriptInfo.Out.CodeHash,
		HashType: c.p.AlwaysSpendableScriptInfo.Out.CodeHashType,
		Args:     c.p.AlwaysSpendableScriptInfo.Out.Args,
	}
}
func (c *OnSaleCell) TypeScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
		HashType: c.p.CellCodeInfo.Out.CodeHashType,
		Args:     c.p.AccountId.Bytes(),
	}
}

func (c *OnSaleCell) Data() ([]byte, error) {
	bys, err := blake2b.Blake256(c.p.OnSaleCellData.AsSlice())
	if err != nil {
		return nil, err
	}
	return bys, nil
}

func (c *OnSaleCell) TableType() TableType {
	return TableType_OnSaleCell
}
