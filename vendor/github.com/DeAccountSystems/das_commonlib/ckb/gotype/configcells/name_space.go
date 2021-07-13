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

type CfgNameSpace struct {
	Data *ConfigCellChildDataObj
}

func (c *CfgNameSpace) Ready() bool{
	return c.Data != nil
}

func (c *CfgNameSpace) Name() string {
	return "configCellNameSpace:"
}

func (c *CfgNameSpace) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgNameSpace) MocluObj() interface{} {
	return nil
}

func (c *CfgNameSpace) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_RecordNamespace
}

func (c *CfgNameSpace) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}