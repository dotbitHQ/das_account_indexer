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

type CfgProposal struct {
	Data *ConfigCellChildDataObj
	MocluData *celltype.ConfigCellProposal
}

func (c *CfgProposal) Ready() bool{
	return c.Data != nil && c.MocluData != nil && c.MocluData.FieldCount() > 0
}

func (c *CfgProposal) Name() string {
	return "configCellProposal:"
}

func (c *CfgProposal) NotifyData(Data *ConfigCellChildDataObj) error {
	c.Data = Data
	if len(c.Data.MoleculeData) == 0 {
		temp := celltype.ConfigCellProposalDefault()
		c.MocluData = &temp
		return nil
	}
	obj, err := celltype.ConfigCellProposalFromSlice(c.Data.MoleculeData, false)
	if err != nil {
		return fmt.Errorf("ConfigCellProposalFromSlice %s",err.Error())
	}
	c.MocluData = obj
	return nil
}

func (c *CfgProposal) MocluObj() interface{} {
	return c.MocluData
}

func (c *CfgProposal) Tag() celltype.TableType {
	return celltype.TableType_CONFIG_CELL_PROPOSAL
}

func (c *CfgProposal) Witness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: &c.Data.CellDep,
		GetWitnessData: func(index uint32) ([]byte, error) {
			return c.Data.WitnessData, nil
		}}
}