package configcells

import (
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
)

/**
 * Copyright (C), 2019-2021
 * FileName: main
 * Author:   LinGuanHong
 * Date:     2021/5/17 2:51
 * Description:
 */

type CfgChatSetHans struct {
	Data *ConfigCellChildDataObj
}

func (c *CfgChatSetHans) Ready() bool {
	return c.Data != nil
}

func (c *CfgChatSetHans) Name() string {
	return "configCellCharsetHanSimple:"
}

func (c *CfgChatSetHans) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgChatSetHans) MocluObj() interface{} {
	return nil
}

func (c *CfgChatSetHans) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_CharSetHanS
}

func (c *CfgChatSetHans) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}
