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

type CfgUnavailable struct {
	Data *ConfigCellChildDataObj
}

func (c *CfgUnavailable) Ready() bool{
	return c.Data != nil
}

func (c *CfgUnavailable) Name() string {
	return "configCellUnavailable:"
}

func (c *CfgUnavailable) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgUnavailable) MocluObj() interface{} {
	return nil
}

func (c *CfgUnavailable) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_Unavailable
}

func (c *CfgUnavailable) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}