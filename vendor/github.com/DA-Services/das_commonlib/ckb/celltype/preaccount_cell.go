package celltype

import (
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2020
 * FileName: statecell
 * Author:   LinGuanHong
 * Date:     2020/12/18 3:58 下午
 * Description:
 */

var TestNetPreAccountCell = func(account DasAccount, new *PreAccountCellData) *PreAccountCellParam {
	return &PreAccountCellParam{
		Version: 1,
		Account: account,
		// Data:         *buildDasCommonMoleculeDataObj(depIndex, oldIndex, newIndex, dep, old, new),
		CellCodeInfo: DasPreAccountCellScript,
		TxDataParam: PreAccountCellTxDataParam{
			NewAccountCellData: new,
		},
		AlwaysSpendableScriptInfo: DasAnyOneCanSendCellInfo,
	}
}

/**
lock: <lock_script>
type: <apply_register_script>
data:
  hash(pubkey_hash + account)
  Timestamp // cell 创建时 TimeCell 的时间
*/

type PreAccountCell struct {
	p *PreAccountCellParam
}

func NewPreAccountCell(p *PreAccountCellParam) *PreAccountCell {
	return &PreAccountCell{p: p}
}

func (c *PreAccountCell) LockDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.AlwaysSpendableScriptInfo.Dep.TxHash,
			Index:  c.p.AlwaysSpendableScriptInfo.Dep.TxIndex,
		},
		DepType: c.p.AlwaysSpendableScriptInfo.Dep.DepType,
	}
}
func (c *PreAccountCell) TypeDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.CellCodeInfo.Dep.TxHash,
			Index:  c.p.CellCodeInfo.Dep.TxIndex,
		},
		DepType: c.p.CellCodeInfo.Dep.DepType,
	}
}
func (c *PreAccountCell) LockScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.AlwaysSpendableScriptInfo.Out.CodeHash,
		HashType: c.p.AlwaysSpendableScriptInfo.Out.CodeHashType,
		Args:     c.p.AlwaysSpendableScriptInfo.Out.Args,
	}
}
func (c *PreAccountCell) TypeScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
		HashType: c.p.CellCodeInfo.Out.CodeHashType,
		Args:     c.p.CellCodeInfo.Out.Args,
	}
}

func (c *PreAccountCell) TableType() TableType {
	return TableType_PRE_ACCOUNT_CELL
}

/**
lock: <always_success>
type: <pre_account_script>
data:
  hash(data: PreAccountCellData)
  id // account ID，生成算法为 hash(account)，然后取前 10 bytes
*/
func (c *PreAccountCell) Data() ([]byte, error) {
	dataHash, err := blake2b.Blake256(c.p.TxDataParam.NewAccountCellData.AsSlice())
	if err != nil {
		return nil, err
	}
	accountId := c.p.Account.AccountId()
	return append(dataHash, accountId.Bytes()...), nil
}
