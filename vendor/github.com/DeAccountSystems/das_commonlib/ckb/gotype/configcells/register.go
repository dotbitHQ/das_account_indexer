package configcells
//
// import (
// 	"fmt"
// 	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
// )
//
// /**
//  * Copyright (C), 2019-2021
//  * FileName: main
//  * Author:   LinGuanHong
//  * Date:     2021/5/17 2:51
//  * Description:
//  */
//
// type CfgRegister struct {
// 	Data *ConfigCellChildDataObj
// 	MocluData *celltype.ConfigCellRegister
// }
//
// func (c *CfgRegister) Ready() bool{
// 	return c.Data != nil && c.MocluData.FieldCount() > 0
// }
//
// func (c *CfgRegister) Name() string {
// 	return "configCellRegister:"
// }
//
// func (c *CfgRegister) NotifyData(Data *ConfigCellChildDataObj) error {
// 	c.Data = Data
// 	if len(c.Data.MoleculeData) == 0 {
// 		temp := celltype.ConfigCellRegisterDefault()
// 		c.MocluData = &temp
// 		return nil
// 	}
// 	obj, err := celltype.ConfigCellRegisterFromSlice(c.Data.MoleculeData, false)
// 	if err != nil {
// 		return fmt.Errorf("ConfigCellRegisterFromSlice %s",err.Error())
// 	}
// 	c.MocluData = obj
// 	return nil
// }
//
// func (c *CfgRegister) MocluObj() interface{} {
// 	return c.MocluData
// }
//
// func (c *CfgRegister) Tag() celltype.CfgCellType {
// 	return celltype.CfgCellType_ConfigCellRegister
// }
//
// func (c *CfgRegister) Witness() *celltype.CellDepWithWitness {
// 	return &celltype.CellDepWithWitness{
// 		CellDep: &c.Data.CellDep,
// 		GetWitnessData: func(index uint32) ([]byte, error) {
// 			return c.Data.WitnessData, nil
// 		}}
// }