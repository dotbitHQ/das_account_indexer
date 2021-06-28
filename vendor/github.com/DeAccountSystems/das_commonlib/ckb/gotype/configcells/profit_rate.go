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

type CfgProfitRate struct {
	Data *ConfigCellChildDataObj
	MocluData *celltype.ConfigCellProfitRate
}

func (c *CfgProfitRate) Ready() bool{
	return c.Data != nil && c.MocluData != nil && c.MocluData.FieldCount() > 0
}

func (c *CfgProfitRate) Name() string {
	return "configCellProfitRate:"
}

func (c *CfgProfitRate) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	if len(c.Data.MoleculeData) == 0 {
		temp := celltype.ConfigCellProfitRateDefault()
		c.MocluData = &temp
		return nil
	}
	obj, err := celltype.ConfigCellProfitRateFromSlice(c.Data.MoleculeData, false)
	if err != nil {
		return fmt.Errorf("ConfigCellProfitRateFromSlice %s",err.Error())
	}
	c.MocluData = obj
	return nil
}

func (c *CfgProfitRate) MocluObj() interface{} {
	return c.MocluData
}

func (c *CfgProfitRate) Tag() celltype.TableType {
	return celltype.TableType_CONFIG_CELL_PROFITRATE
}

func (c *CfgProfitRate) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}