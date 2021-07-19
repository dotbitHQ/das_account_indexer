package configcells

import (
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
)

/**
 * Copyright (C), 2019-2021
 * FileName: main
 * Author:   LinGuanHong
 * Date:     2021/5/17 2:51
 * Description:
 */

type CfgMain struct {
	Data *ConfigCellChildDataObj
	MocluData *celltype.ConfigCellMain
}

func (c *CfgMain) Ready() bool{
	return c.Data != nil && c.MocluData != nil && c.MocluData.FieldCount() > 0
}

func (c *CfgMain) Name() string {
	return "configCellMain:"
}

func (c *CfgMain) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	if len(c.Data.MoleculeData) == 0 {
		temp := celltype.ConfigCellMainDefault()
		c.MocluData = &temp
		return nil
	}
	obj, err := celltype.ConfigCellMainFromSlice(c.Data.MoleculeData, false)
	if err != nil {
		return fmt.Errorf("ConfigCellMainFromSlice %s",err.Error())
	}
	c.MocluData = obj
	return nil
}

func (c *CfgMain) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_Main
}

func (c *CfgMain) MocluObj() interface{} {
	return c.MocluData
}

func (c *CfgMain) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}