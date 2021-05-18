package celltype

import (
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: proposalcell
 * Author:   LinGuanHong
 * Date:     2021/1/10 12:22 下午
 * Description:
 */

/**
lock: <always_success>
type:
  code_hash: <proposal_script>,
  type: type,
  args: [],
data:
  // 每个提案以自己的 ProposalCellData 的 hash 为 ID，基于前一个提案发起新提案时，新提案需要携带前一个提案的 hash
  hash_of_first_proposal
  ...
  hash_of_prev_proposal
  hash_of_self(data: ProposalCellData)

witness:
  table Data {
    old: None,
    new: table DataEntityOpt {
      index: Uint32,
      version: Uint32,
      entity: ProposalCellData
    },
  }

======
table ProposalCellData {
    starter_lock: Script,
    slices: SliceList,
}

vector SliceList <SL>;

// SL is used here for "slice" because "slice" may be a keyword in some languages.
vector SL <ProposalItem>;

table ProposalItem {
  account_id: AccountId,
  item_type: Uint8,
  // When account is at the end of the linked list, its next pointer should be None.
  next: AccountIdOpt,
}

====== 举例来说看起来就像下面这样
table ProposalCellData {
  starter_lock: Script,
  slices: [
    [
      { account_id: xxx, item_type: exist, next: xxx },
      { account_id: xxx, item_type: proposed, next: xxx },
      { account_id: xxx, item_type: proposed, next: xxx },
      { account_id: xxx, item_type: new, next: xxx },
    ],
    [
      { account_id: xxx, item_type: exist, next: xxx },
      { account_id: xxx, item_type: proposed, next: xxx },
      { account_id: xxx, item_type: new, next: xxx },
      { account_id: xxx, item_type: proposed, next: xxx },
    ],
    [
      { account_id: xxx, item_type: exist, next: xxx },
      { account_id: xxx, item_type: new, next: xxx },
    ],
    ...
  ]
}
```

- type.args 的结构：当提案只有自己时，那么参数中就只有一个 ID，如果提案基于其他提案发起时，那么参数中前面的部分就是父提案的 ID ，按照依赖关系排列，最后才是自己的 ID ；ID 的计算方式就是整个 ProposalCellData 的 hash ，因为提案发布后不可更改，所以对于一个提案来说 ID 是不变的，但是完全相同的提案 ID 可能相同；
- starter_lock，即提案发起者接收利润分成的的 lock 脚本，即地址；
- slices，当前提案通过后 AccountCell 链表被修改部分的最终状态，其解释详见 `TODO`；

*/

var TestNetProposeCell = func(new *ProposalCellData) *ProposeCellParam {
	acp := &ProposeCellParam{
		Version:     1,
		TxDataParam: *new,
		// Data:                      *buildDasCommonMoleculeDataObj(depIndex, oldIndex, newIndex, dep, old, new),
		CellCodeInfo:              DasProposeCellScript,
		AlwaysSpendableScriptInfo: DasAnyOneCanSendCellInfo,
	}
	return acp
}

type ProposeCell struct {
	p *ProposeCellParam
}

func NewProposeCell(p *ProposeCellParam) *ProposeCell {
	return &ProposeCell{p: p}
}
func (c *ProposeCell) SoDeps() []types.CellDep {
	return nil
}
func (c *ProposeCell) LockDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.AlwaysSpendableScriptInfo.Dep.TxHash,
			Index:  c.p.AlwaysSpendableScriptInfo.Dep.TxIndex,
		},
		DepType: c.p.AlwaysSpendableScriptInfo.Dep.DepType,
	}
}
func (c *ProposeCell) TypeDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.CellCodeInfo.Dep.TxHash,
			Index:  c.p.CellCodeInfo.Dep.TxIndex,
		},
		DepType: c.p.CellCodeInfo.Dep.DepType,
	}
}
func (c *ProposeCell) LockScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.AlwaysSpendableScriptInfo.Out.CodeHash,
		HashType: c.p.AlwaysSpendableScriptInfo.Out.CodeHashType,
		Args:     c.p.AlwaysSpendableScriptInfo.Out.Args,
	}
}
func (c *ProposeCell) TypeScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
		HashType: c.p.CellCodeInfo.Out.CodeHashType,
		Args:     nil,
	}
}

func (c *ProposeCell) Data() ([]byte, error) {
	hashBys, _ := blake2b.Blake256(c.p.TxDataParam.AsSlice())
	return hashBys, nil
}

func (c *ProposeCell) TableType() TableType {
	return TableType_PROPOSE_CELL
}
