package celltype
//
// import (
// 	"github.com/nervosnetwork/ckb-sdk-go/types"
// )
//
// /**
//  * Copyright (C), 2019-2020
//  * FileName: refcell
//  * Author:   LinGuanHong
//  * Date:     2020/12/27 11:17 上午
//  * Description:
//  */
//
// var TestNetRefCell = func(lockScript *types.Script, accountId DasAccountId, refType RefCellType) *RefcellParam {
// 	return &RefcellParam{
// 		Version:      1,
// 		AccountId:    accountId,
// 		RefType:      refType,
// 		CellCodeInfo: DasRefCellScript,
// 		UserLockScript: DASCellBaseInfo{
// 			Dep: TestNetLockScriptDep,
// 			Out: DASCellBaseInfoOut{
// 				CodeHash:     lockScript.CodeHash,
// 				CodeHashType: lockScript.HashType,
// 				Args:         lockScript.Args,
// 			},
// 		},
// 	}
// }
//
// type Refcell struct {
// 	p *RefcellParam
// }
//
// func NewRefcell(p *RefcellParam) *Refcell {
// 	return &Refcell{p: p}
// }
// func (c *Refcell) SoDeps() []types.CellDep {
// 	return nil
// }
// func (c *Refcell) LockDepCell() *types.CellDep {
// 	return &types.CellDep{
// 		OutPoint: &types.OutPoint{
// 			TxHash: c.p.UserLockScript.Dep.TxHash,
// 			Index:  c.p.UserLockScript.Dep.TxIndex,
// 		},
// 		DepType: c.p.UserLockScript.Dep.DepType,
// 	}
// }
// func (c *Refcell) TypeDepCell() *types.CellDep {
// 	return &types.CellDep{ // state_cell
// 		OutPoint: &types.OutPoint{
// 			TxHash: c.p.CellCodeInfo.Dep.TxHash,
// 			Index:  c.p.CellCodeInfo.Dep.TxIndex, // state_script_tx_index
// 		},
// 		DepType: c.p.CellCodeInfo.Dep.DepType,
// 	}
// }
// func (c *Refcell) LockScript() *types.Script {
// 	return &types.Script{
// 		CodeHash: c.p.UserLockScript.Out.CodeHash,
// 		HashType: c.p.UserLockScript.Out.CodeHashType,
// 		Args:     c.p.UserLockScript.Out.Args,
// 	}
// }
// func (c *Refcell) TypeScript() *types.Script {
// 	return &types.Script{
// 		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
// 		HashType: c.p.CellCodeInfo.Out.CodeHashType,
// 		Args:     c.p.CellCodeInfo.Out.Args,
// 	}
// }
//
// /**
// data:
//   id // 10 Bytes 的 Account ID
//   role // 1 Bytes 的身份区分符
// */
// func (c *Refcell) Data() ([]byte, error) {
// 	return append(c.p.AccountId.Bytes(), []byte{uint8(c.p.RefType)}...), nil
// }
//
// func (c *Refcell) TableType() TableType {
// 	return 0
// }
