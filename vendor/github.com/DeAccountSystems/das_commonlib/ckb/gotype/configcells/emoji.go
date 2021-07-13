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

type CfgChatSetEmoji struct {
	Data *ConfigCellChildDataObj
}

func (c *CfgChatSetEmoji) Ready() bool{
	return c.Data != nil
}

func (c *CfgChatSetEmoji) Name() string {
	return "configCellCharsetEmoji:"
}

func (c *CfgChatSetEmoji) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgChatSetEmoji) MocluObj() interface{} {
	return nil
}

func (c *CfgChatSetEmoji) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_CharSetEmoji
}

func (c *CfgChatSetEmoji) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}