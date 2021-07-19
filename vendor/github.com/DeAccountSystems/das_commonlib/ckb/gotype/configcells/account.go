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

type CfgAccount struct {
	Data *ConfigCellChildDataObj
	MocluData *celltype.ConfigCellAccount
}

func (c *CfgAccount) Ready() bool{
	return c.Data != nil && c.MocluData != nil && c.MocluData.FieldCount() > 0
}

func (c *CfgAccount) Name() string {
	return "configCellAccount:"
}

func (c *CfgAccount) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	if len(c.Data.MoleculeData) == 0 {
		temp := celltype.ConfigCellAccountDefault()
		c.MocluData = &temp
		return nil
	}
	obj, err := celltype.ConfigCellAccountFromSlice(c.Data.MoleculeData, false)
	if err != nil {
		return fmt.Errorf("ConfigCellAccountFromSlice %s",err.Error())
	}
	c.MocluData = obj
	return nil
}

func (c *CfgAccount) MocluObj() interface{} {
	return c.MocluData
}

func (c *CfgAccount) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_Account
}

func (c *CfgAccount) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}