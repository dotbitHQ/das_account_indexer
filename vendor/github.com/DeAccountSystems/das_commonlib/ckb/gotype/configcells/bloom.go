package configcells
//
// import (
// 	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
// )
//
// /**
//  * Copyright (C), 2019-2021
//  * FileName: main
//  * Author:   LinGuanHong
//  * Date:     2021/5/17 2:51 下午
//  * Description:
//  */
//
// type CfgBloom struct {
// 	Data *ConfigCellChildDataObj
// }
//
// func (c *CfgBloom) Ready() bool{
// 	return c.Data != nil
// }
//
// func (c *CfgBloom) Name() string {
// 	return "configCellBloom:"
// }
//
// func (c *CfgBloom) NotifyData(Data *ConfigCellChildDataObj) error {
// 	c.Data = Data
// 	return nil
// }
//
// func (c *CfgBloom) Tag() celltype.TableType {
// 	return celltype.CfgCellType_ConfigCellBloomFilter
// }
//
// func (c *CfgBloom) MocluObj() interface{} {
// 	return nil
// }
//
// func (c *CfgBloom) Witness() *celltype.CellDepWithWitness {
// 	return &celltype.CellDepWithWitness{
// 		CellDep: &c.Data.CellDep,
// 		GetWitnessData: func(index uint32) ([]byte, error) {
// 			return c.Data.WitnessData, nil
// 		}}
// }