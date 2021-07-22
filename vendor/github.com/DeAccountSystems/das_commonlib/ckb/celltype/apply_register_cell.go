package celltype

import (
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2020
 * FileName: statecell
 * Author:   LinGuanHong
 * Date:     2020/12/18 3:58
 * Description:
 */

var DefaultApplyRegisterCellParam = func(args []byte, account DasAccount, height,timeUnix uint64, senderLockScript *types.Script) *ApplyRegisterCellParam {
	return &ApplyRegisterCellParam{
		Version:         1,
		PubkeyHashBytes: args,
		Account:         account,
		Height:          height,
		TimeUnix:        timeUnix,
		CellCodeInfo:    DasApplyRegisterCellScript,
		SenderLockScriptInfo: DASCellBaseInfo{
			Dep: DasApplyRegisterCellScript.Dep,
			Out: DASCellBaseInfoOutFromScript(senderLockScript),
		},
	}
}

type ApplyRegisterCell struct {
	p *ApplyRegisterCellParam
}

func NewApplyRegisterCell(p *ApplyRegisterCellParam) *ApplyRegisterCell {
	return &ApplyRegisterCell{p: p}
}
func (c *ApplyRegisterCell) SoDeps() []types.CellDep {
	return nil
}
func (c *ApplyRegisterCell) LockDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.SenderLockScriptInfo.Dep.TxHash,
			Index:  c.p.SenderLockScriptInfo.Dep.TxIndex,
		},
		DepType: c.p.SenderLockScriptInfo.Dep.DepType,
	}
}
func (c *ApplyRegisterCell) TypeDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.p.CellCodeInfo.Dep.TxHash,
			Index:  c.p.CellCodeInfo.Dep.TxIndex,
		},
		DepType: c.p.CellCodeInfo.Dep.DepType,
	}
}
func (c *ApplyRegisterCell) LockScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.SenderLockScriptInfo.Out.CodeHash,
		HashType: c.p.SenderLockScriptInfo.Out.CodeHashType,
		Args:     c.p.SenderLockScriptInfo.Out.Args,
	}
}
func (c *ApplyRegisterCell) TypeScript() *types.Script {
	return &types.Script{
		CodeHash: c.p.CellCodeInfo.Out.CodeHash,
		HashType: c.p.CellCodeInfo.Out.CodeHashType,
		Args:     nil,
	}
}

func (c *ApplyRegisterCell) TableType() TableType {
	return 0
}

func (c *ApplyRegisterCell) Data() ([]byte, error) {
	idHash, err := ApplyRegisterDataId(c.p.PubkeyHashBytes, c.p.Account)
	if err != nil {
		return nil, fmt.Errorf("ApplyRegisterDataId err: %s", err.Error())
	}
	temp := append(idHash, GoUint64ToBytes(c.p.Height)...)
	return append(temp, GoUint64ToBytes(c.p.TimeUnix)...), nil
}

func ApplyRegisterDataId(pubKeyHexBytes []byte, account DasAccount) ([]byte, error) {
	if err := account.ValidErr(); err != nil {
		return nil, fmt.Errorf("ApplyRegisterDataId err: %s", err.Error())
	}
	accountBytes := []byte(account)
	targetBytes := append(pubKeyHexBytes, accountBytes...)
	return blake2b.Blake256(targetBytes)
}
