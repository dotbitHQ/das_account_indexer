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

type CfgPrice struct {
	Data *ConfigCellChildDataObj
	MocluData *celltype.ConfigCellPrice
}

func (c *CfgPrice) Ready() bool{
	return c.Data != nil && c.MocluData != nil && c.MocluData.FieldCount() > 0
}

func (c *CfgPrice) Name() string {
	return "configCellPrice:"
}

func (c *CfgPrice) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	if len(c.Data.MoleculeData) == 0 {
		temp := celltype.ConfigCellPriceDefault()
		c.MocluData = &temp
		return nil
	}
	obj, err := celltype.ConfigCellPriceFromSlice(c.Data.MoleculeData, false)
	if err != nil {
		return fmt.Errorf("ConfigCellPriceFromSlice %s",err.Error())
	}
	c.MocluData = obj
	return nil
}

func (c *CfgPrice) MocluObj() interface{} {
	return c.MocluData
}

func (c *CfgPrice) Tag() celltype.TableType {
	return celltype.TableType_CONFIG_CELL_PRICE
}

func (c *CfgPrice) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}