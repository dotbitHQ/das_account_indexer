package celltype

//
// import "github.com/nervosnetwork/ckb-sdk-go/types"
//
// /**
//  * Copyright (C), 2019-2021
//  * FileName: bidding_cell
//  * Author:   LinGuanHong
//  * Date:     2021/2/22 2:24 下午
//  * Description:
//  */
//
// var TestNetBiddingCell = func(newIndex uint32, price uint64, accountId DasAccountId) *BiddingCellParam {
// 	BiddingMoleData := NewBiddingCellDataBuilder().p(GoUint64ToMoleculeU64(price)).Build()
// 	return &BiddingCellParam{
// 		Version: 1,
// 		Price:   price,
// 		Data: *buildDasCommonMoleculeDataObj(
// 			0, 0, newIndex, nil, nil, &BiddingMoleData),
// 		AccountId:    accountId,
// 		CellCodeInfo: DasBiddingCellScript,
// 		AlwaysSpendableScriptInfo: DASCellBaseInfo{
// 			Dep: DASCellBaseInfoDep{
// 				TxHash:  types.HexToHash("0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37"),
// 				TxIndex: 0,
// 				DepType: types.DepTypeDepGroup,
// 			},
// 			Out: DasAnyOneCanSendCellInfo,
// 		},
// 	}
// }
//
// type BiddingCell struct {
// 	p *BiddingCellParam
// }
//
// func NewBiddingCell(p *BiddingCellParam) *BiddingCell {
// 	return &BiddingCell{p: p}
// }
//
// func (c *BiddingCell) LockDepCell() *types.CellDep {
// 	return &types.CellDep{
// 		OutPoint: &types.OutPoint{
// 			TxHash: c.p.AlwaysSpendableScriptInfo.Dep.TxHash,
// 			Index:  c.p.AlwaysSpendableScriptInfo.Dep.TxIndex,
// 		},
// 		DepType: c.p.AlwaysSpendableScriptInfo.Dep.DepType,
// 	}
// }
// func (c *BiddingCell) TypeDepCell() *types.CellDep {
// 	return &types.CellDep{
// 		OutPoint: &types.OutPoint{
// 			TxHash: c.p.CellCodeInfo.Dep.TxHash,
// 			Index:  c.p.CellCodeInfo.Dep.TxIndex,
// 		},
// 		DepType: c.p.CellCodeInfo.Dep.DepType,
// 	}
// }
// func (c *BiddingCell) LockScript() *types.Script {
// 	return &types.Script{
// 		CodeHash: c.p.AlwaysSpendableScriptInfo.Out.CodeHash,
// 		HashType: c.p.AlwaysSpendableScriptInfo.Out.CodeHashType,
// 		Args:     c.p.AlwaysSpendableScriptInfo.Out.Args,
// 	}
// }
// func (c *BiddingCell) TypeScript() *types.Script {
// 	return &types.Script{
// 		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
// 		HashType: c.p.CellCodeInfo.Out.CodeHashType,
// 		Args:     c.p.AccountId,
// 	}
// }
//
// func (c *BiddingCell) Data() ([]byte, error) {
// 	bys, err := blake2b.Blake256(c.TableData())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return bys, nil
// }
//
// func (c *BiddingCell) TableType() TableType {
// 	return TableType_ON_SALE_CELL
// }
//
// func (c *BiddingCell) TableData() []byte {
// 	return c.p.Data.AsSlice()
// }
//
//
//
