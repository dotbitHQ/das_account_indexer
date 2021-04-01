package celltype

//
// import (
// 	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
// 	"github.com/nervosnetwork/ckb-sdk-go/types"
// )
//
// /**
//  * Copyright (C), 2019-2020
//  * FileName: statecell
//  * Author:   LinGuanHong
//  * Date:     2020/12/18 3:58 下午
//  * Description:
//  */
//
// var TestNetStateCell = func(dataBuilder *StateCellData) *StateCellParam {
// 	return &StateCellParam{
// 		Version:      1,
// 		Data:         dataBuilder,
// 		CellCodeInfo: DasStateCellScript,
// 		AlwaysSpendableScriptInfo: DASCellBaseInfo{
// 			Dep: DASCellBaseInfoDep{
// 				TxHash:  "0xec26b0f85ed839ece5f11c4c4e837ec359f5adc4420410f6453b1f6b60fb96a6",
// 				TxIndex: 0,
// 				DepType: types.DepTypeDepGroup,
// 			},
// 			Out: DASCellBaseInfoOut{
// 				CodeHash:     "0x3419a1c09eb2567f6552ee7a8ecffd64155cffe0f1796e6e61ec088d740c1356",
// 				CodeHashType: types.HashTypeType,
// 				Args:         nil,
// 			},
// 		},
// 	}
// }
//
// /**
// lock: <always_spendable_script>
// type: <state_script>
// data:
//   [version: u32]
//   table StateCellData {
//       reserved_account_root: Hash,
//       price_list: PriceList,
//       char_set: CharSet,
//       min_ttl: Uint32,
//       closing_limit_of_primary_market_auction: Uint32,
//       closing_limit_of_secondary_market_auction: Uint32,
//   }
// */
//
// type StateCell struct {
// 	p *StateCellParam
// }
//
// func NewStateCell(p *StateCellParam) *StateCell {
// 	return &StateCell{p: p}
// }
//
// func (c *StateCell) LockDepCell() *types.CellDep {
// 	return &types.CellDep{
// 		OutPoint: &types.OutPoint{
// 			TxHash: types.HexToHash(c.p.AlwaysSpendableScriptInfo.Dep.TxHash),
// 			Index:  c.p.AlwaysSpendableScriptInfo.Dep.TxIndex,
// 		},
// 		DepType: c.p.AlwaysSpendableScriptInfo.Dep.DepType,
// 	}
// }
// func (c *StateCell) TypeDepCell() *types.CellDep {
// 	return &types.CellDep{ // state_cell
// 		OutPoint: &types.OutPoint{
// 			TxHash: types.HexToHash(c.p.CellCodeInfo.Dep.TxHash),
// 			Index:  c.p.CellCodeInfo.Dep.TxIndex, // state_script_tx_index
// 		},
// 		DepType: c.p.CellCodeInfo.Dep.DepType,
// 	}
// }
// func (c *StateCell) LockScript() *types.Script {
// 	return &types.Script{
// 		CodeHash: types.HexToHash(c.p.AlwaysSpendableScriptInfo.Out.CodeHash),
// 		HashType: c.p.AlwaysSpendableScriptInfo.Out.CodeHashType,
// 		Args:     c.p.AlwaysSpendableScriptInfo.Out.Args,
// 	}
// }
// func (c *StateCell) TypeScript() *types.Script {
// 	return &types.Script{
// 		CodeHash: types.HexToHash(c.p.CellCodeInfo.Out.CodeHash),
// 		HashType: c.p.CellCodeInfo.Out.CodeHashType,
// 		Args:     c.p.CellCodeInfo.Out.Args,
// 	}
// }
//
// func (c *StateCell) Data() ([]byte, error) {
// 	return blake2b.Blake256(c.TableData())
// }
//
// func (c *StateCell) TableType() TableType {
// 	return TableType_STATE_CELL
// }
//
// func (c *StateCell) TableData() []byte {
// 	return c.p.Data.AsSlice()
// }
