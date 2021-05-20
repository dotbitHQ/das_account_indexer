package celltype
//
// import (
// 	"github.com/nervosnetwork/ckb-sdk-go/types"
// )
//
// /**
//  * Copyright (C), 2019-2021
//  * FileName: wallet_cell
//  * Author:   LinGuanHong
//  * Date:     2021/2/17 12:30 下午
//  * Description:
//  */
//
// var TestNetWalletCell = func(accountId DasAccountId) *WalletCellParam {
// 	return &WalletCellParam{
// 		AccountId:              accountId,
// 		CellCodeInfo:           DasWalletCellScript,
// 		AnyoneCanPayScriptInfo: DasAnyOneCanSendCellInfo,
// 	}
// }
//
// /**
// lock: <always_success>
// type:
//     code_hash: <wallet-cell-type>
//     type: type
//     args: <account_id>
// data:
// */
//
// type WalletCell struct {
// 	p *WalletCellParam
// }
//
// func NewWalletCell(p *WalletCellParam) *WalletCell {
// 	return &WalletCell{p: p}
// }
// func (c *WalletCell) SoDeps() []types.CellDep {
// 	return nil
// }
// func (c *WalletCell) LockDepCell() *types.CellDep {
// 	return &types.CellDep{
// 		OutPoint: &types.OutPoint{
// 			TxHash: c.p.AnyoneCanPayScriptInfo.Dep.TxHash,
// 			Index:  c.p.AnyoneCanPayScriptInfo.Dep.TxIndex,
// 		},
// 		DepType: c.p.AnyoneCanPayScriptInfo.Dep.DepType,
// 	}
// }
// func (c *WalletCell) TypeDepCell() *types.CellDep {
// 	return &types.CellDep{
// 		OutPoint: &types.OutPoint{
// 			TxHash: c.p.CellCodeInfo.Dep.TxHash,
// 			Index:  c.p.CellCodeInfo.Dep.TxIndex,
// 		},
// 		DepType: c.p.CellCodeInfo.Dep.DepType,
// 	}
// }
// func (c *WalletCell) LockScript() *types.Script {
// 	return &types.Script{
// 		CodeHash: c.p.AnyoneCanPayScriptInfo.Out.CodeHash,
// 		HashType: c.p.AnyoneCanPayScriptInfo.Out.CodeHashType,
// 		Args:     c.p.AnyoneCanPayScriptInfo.Out.Args,
// 	}
// }
// func (c *WalletCell) TypeScript() *types.Script {
// 	return &types.Script{
// 		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
// 		HashType: c.p.CellCodeInfo.Out.CodeHashType,
// 		Args:     c.p.CellCodeInfo.Out.Args,
// 	}
// }
//
// func (c *WalletCell) TableType() TableType {
// 	return 0
// }
//
// func (c *WalletCell) Data() ([]byte, error) {
// 	return c.p.AccountId.Bytes(), nil
// }
//
// func (c *WalletCell) TableData() []byte {
// 	return nil
// }
