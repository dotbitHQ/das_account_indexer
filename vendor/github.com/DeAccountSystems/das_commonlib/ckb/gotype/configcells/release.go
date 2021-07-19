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

type CfgRelease struct {
	Data *ConfigCellChildDataObj
	MocluData *celltype.ConfigCellRelease
}

func (c *CfgRelease) Ready() bool{
	return c.Data != nil && c.MocluData != nil && c.MocluData.FieldCount() > 0
}

func (c *CfgRelease) Name() string {
	return "configCellRelease:"
}

func (c *CfgRelease) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	if len(c.Data.MoleculeData) == 0 {
		temp := celltype.ConfigCellReleaseDefault()
		c.MocluData = &temp
		return nil
	}
	obj, err := celltype.ConfigCellReleaseFromSlice(c.Data.MoleculeData, false)
	if err != nil {
		return fmt.Errorf("ConfigCellReleaseFromSlice %s",err.Error())
	}
	c.MocluData = obj
	return nil
}

func (c *CfgRelease) MocluObj() interface{} {
	return c.MocluData
}

func (c *CfgRelease) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_Release
}

func (c *CfgRelease) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}