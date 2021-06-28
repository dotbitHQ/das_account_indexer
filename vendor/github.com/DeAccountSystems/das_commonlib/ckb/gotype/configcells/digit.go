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

type CfgChatSetDigit struct {
	Data *ConfigCellChildDataObj
}

func (c *CfgChatSetDigit) Ready() bool{
	return c.Data != nil
}

func (c *CfgChatSetDigit) Name() string {
	return "configCellCharsetDigit:"
}

func (c *CfgChatSetDigit) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgChatSetDigit) MocluObj() interface{} {
	return nil
}

func (c *CfgChatSetDigit) Tag() celltype.TableType {
	return celltype.TableType_CONFIG_CELL_CharSetDigit
}

func (c *CfgChatSetDigit) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}