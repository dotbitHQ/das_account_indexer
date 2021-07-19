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

type CfgPreservedAccountXX struct {
	name string
	TableType celltype.TableType
	Data *ConfigCellChildDataObj
}

func NewCfgPreservedAccount(tableType celltype.TableType,name string) *CfgPreservedAccountXX {
	return &CfgPreservedAccountXX{
		name:      name,
		TableType: tableType,
	}
}

func (c *CfgPreservedAccountXX) Ready() bool{
	return c.Data != nil
}

func (c *CfgPreservedAccountXX) Name() string {
	return c.name
}

func (c *CfgPreservedAccountXX) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	return nil
}

func (c *CfgPreservedAccountXX) MocluObj() interface{} {
	return nil
}

func (c *CfgPreservedAccountXX) Tag() celltype.TableType {
	return c.TableType
}

func (c *CfgPreservedAccountXX) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}