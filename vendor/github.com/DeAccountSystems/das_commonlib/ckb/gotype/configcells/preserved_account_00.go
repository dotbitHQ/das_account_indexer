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

type CfgPreservedAccount00 struct {
	Data *ConfigCellChildDataObj
}

func (c *CfgPreservedAccount00) Ready() bool{
	return c.Data != nil
}

func (c *CfgPreservedAccount00) Name() string {
	return "configCellPreservedAccount00:"
}

func (c *CfgPreservedAccount00) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgPreservedAccount00) MocluObj() interface{} {
	return nil
}

func (c *CfgPreservedAccount00) Tag() celltype.TableType {
	return celltype.TableType_CONFIG_CELL_PreservedAccount00
}

func (c *CfgPreservedAccount00) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}