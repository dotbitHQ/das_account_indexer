package configcells

import (
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
)

/**
 * Copyright (C), 2019-2021
 * FileName: main
 * Author:   LinGuanHong
 * Date:     2021/5/17 2:51 下午
 * Description:
 */

type CfgIncome struct {
	Data *ConfigCellChildDataObj
	MocluData *celltype.ConfigCellIncome
}

func (c *CfgIncome) Ready() bool{
	return c.Data != nil && c.MocluData != nil && c.MocluData.FieldCount() > 0
}

func (c *CfgIncome) Name() string {
	return "configCellIncome:"
}

func (c *CfgIncome) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	if len(c.Data.MoleculeData) == 0 {
		temp := celltype.ConfigCellIncomeDefault()
		c.MocluData = &temp
		return nil
	}
	obj, err := celltype.ConfigCellIncomeFromSlice(c.Data.MoleculeData, false)
	if err != nil {
		return fmt.Errorf("ConfigCellIncomeFromSlice %s",err.Error())
	}
	c.MocluData = obj
	return nil
}

func (c *CfgIncome) MocluObj() interface{} {
	return c.MocluData
}

func (c *CfgIncome) Tag() celltype.TableType {
	return celltype.TableType_CONFIG_CELL_INCOME
}

func (c *CfgIncome) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}