package celltype

import (
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: quote_cell
 * Author:   LinGuanHong
 * Date:     2021/4/19 10:45 上午
 * Description:
 */

type QuoteCell struct {
	p *QuoteCellParam
}

func NewQuoteCell(p *QuoteCellParam) *QuoteCell {
	return &QuoteCell{p: p}
}
func (c *QuoteCell) SoDeps() []types.CellDep {
	return nil
}
func (c *QuoteCell) LockDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.CellCodeInfo.Dep.TxHash,
			Index:  c.p.CellCodeInfo.Dep.TxIndex,
		},
		DepType: c.p.CellCodeInfo.Dep.DepType,
	}
}
func (c *QuoteCell) TypeDepCell() *types.CellDep {
	return nil
}
func (c *QuoteCell) LockScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
		HashType: c.p.CellCodeInfo.Out.CodeHashType,
		Args:     c.p.CellCodeInfo.Out.Args,
	}
}
func (c *QuoteCell) TypeScript() *types.Script {
	return nil
}

func (c *QuoteCell) Data() ([]byte, error) {
	return GoUint64ToBytes(c.p.Price), nil
}

func (c *QuoteCell) TableType() TableType {
	return 0
}
