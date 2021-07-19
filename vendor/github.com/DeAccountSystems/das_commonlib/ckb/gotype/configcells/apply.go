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

type CfgApply struct {
	Data *ConfigCellChildDataObj
	MocluData *celltype.ConfigCellApply
}

func (c *CfgApply) Ready() bool{
	return c.Data != nil && c.MocluData != nil && c.MocluData.FieldCount() > 0
}

func (c *CfgApply) Name() string {
	return "configCellApply:"
}

func (c *CfgApply) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	if len(c.Data.MoleculeData) == 0 {
		temp := celltype.ConfigCellApplyDefault()
		c.MocluData = &temp
		return nil
	}
	obj, err := celltype.ConfigCellApplyFromSlice(c.Data.MoleculeData, false)
	if err != nil {
		return fmt.Errorf("ConfigCellApplyFromSlice %s",err.Error())
	}
	c.MocluData = obj
	return nil
}

func (c *CfgApply) MocluObj() interface{} {
	return c.MocluData
}

func (c *CfgApply) Tag() celltype.TableType {
	return celltype.TableType_ConfigCell_Apply
}

func (c *CfgApply) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}