package configcells

import (
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
)

/**
 * Copyright (C), 2019-2021
 * FileName: main
 * Author:   LinGuanHong
 * Date:     2021/5/17 2:51 下午
 * Description:
 */

type CfgChatSetHant struct {
	Data *ConfigCellChildDataObj
}

func (c *CfgChatSetHant) Ready() bool {
	return c.Data != nil
}

func (c *CfgChatSetHant) Name() string {
	return "configCellCharsetHanT:"
}

func (c *CfgChatSetHant) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgChatSetHant) MocluObj() interface{} {
	return nil
}

func (c *CfgChatSetHant) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_CharSetHanT
}

func (c *CfgChatSetHant) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}
