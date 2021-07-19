package gotype

import (
	"github.com/DeAccountSystems/das_commonlib/common"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: quote_cell
 * Author:   LinGuanHong
 * Date:     2021/7/16 4:40
 * Description:
 */

type QuoteCell struct {
	Data    []byte
	CellDep types.CellDep
}

func (q *QuoteCell) Quote() (int64, error) {
	return common.BytesToInt64(q.Data[2:]), nil
}
