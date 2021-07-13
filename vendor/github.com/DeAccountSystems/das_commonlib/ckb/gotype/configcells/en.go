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

type CfgChatSetEn struct {
	Data *ConfigCellChildDataObj
}

func (c *CfgChatSetEn) Ready() bool{
	return c.Data != nil
}

func (c *CfgChatSetEn) Name() string {
	return "configCellCharsetEn:"
}

func (c *CfgChatSetEn) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgChatSetEn) MocluObj() interface{} {
	return nil
}

func (c *CfgChatSetEn) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_CharSetEn
}

func (c *CfgChatSetEn) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}