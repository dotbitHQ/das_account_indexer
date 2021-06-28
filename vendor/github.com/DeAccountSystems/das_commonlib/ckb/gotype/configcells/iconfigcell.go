package configcells

import (
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: config_child_cells
 * Author:   LinGuanHong
 * Date:     2021/5/17 2:21 下午
 * Description:
 */

type ConfigCellChildDataObj struct {
	CellDep      types.CellDep
	WitnessData  []byte
	MoleculeData []byte
}

type IConfigChild interface {
	Tag() celltype.TableType
	Name() string
	Witness() *celltype.CellDepWithWitness
	NotifyData(data *ConfigCellChildDataObj) error
	MocluObj() interface{}
	Ready() bool
}


