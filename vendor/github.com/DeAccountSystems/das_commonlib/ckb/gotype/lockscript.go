package gotype

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/chain/tron_chain"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"

	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
)

/**
 * Copyright (C), 2019-2021
 * FileName: lockscript
 * Author:   LinGuanHong
 * Date:     2021/6/9 6:23
 * Description:
 */

type PayTypeLockScriptsRet struct {
	FeeCellScript       *ckbTypes.Script
	UserScript          *ckbTypes.Script
	ScriptType          celltype.LockScriptType
	DasLockArgsParam    celltype.DasLockArgsPairParam
	DasLockArgsParam712 celltype.DasLockArgsPairParam
	Err                 error
}

// ckb argsHexStr, eth address
func GetOwnerArgsFromDasLockArgs(args []byte) (celltype.ChainType, string) {
	if len(args) != celltype.DasLockArgsMinBytesLen {
		return 0, ""
	}
	indexType := args[0]
	ownerArgsBytes := args[1 : celltype.DasLockArgsMinBytesLen/2]
	switch celltype.DasLockCodeHashIndexType(indexType) {
	case celltype.DasLockCodeHashIndexType_CKB_Normal:
		return celltype.ChainType_CKB, hex.EncodeToString(ownerArgsBytes)
	case celltype.DasLockCodeHashIndexType_CKB_AnyOne:
		return celltype.ChainType_CKB, hex.EncodeToString(ownerArgsBytes)
	case celltype.DasLockCodeHashIndexType_ETH_Normal, celltype.DasLockCodeHashIndexType_712_Normal:
		return celltype.ChainType_ETH, "0x" + hex.EncodeToString(ownerArgsBytes)
	case celltype.DasLockCodeHashIndexType_TRON_Normal:
		return celltype.ChainType_TRON, tron_chain.TronAddrHexPrefix + hex.EncodeToString(ownerArgsBytes)
	default:
		return 0, ""
	}
}

func GetDasLockScript(chainType celltype.ChainType, address Address) (*ckbTypes.Script, error) {
	switch chainType {
	case celltype.ChainType_CKB:
		return address.DasLockScript_CKB()
	case celltype.ChainType_ETH:
		return address.DasLockScript(celltype.DasLockCodeHashIndexType_ETH_Normal)
	case celltype.ChainType_BTC:
		return address.DasLockScript(celltype.ScriptType_BTC.ToDasLockCodeHashIndexType())
	case celltype.ChainType_TRON:
		return address.DasLockScript(celltype.DasLockCodeHashIndexType_TRON_Normal)
	}
	return nil, fmt.Errorf("unknow chain type:%d %s", chainType, address)
}

func GetDasLockScript712(chainType celltype.ChainType, address Address) (*ckbTypes.Script, error) {
	switch chainType {
	case celltype.ChainType_CKB:
		return address.DasLockScript_CKB()
	case celltype.ChainType_ETH:
		return address.DasLockScript(celltype.DasLockCodeHashIndexType_712_Normal)
	case celltype.ChainType_BTC:
		return address.DasLockScript(celltype.ScriptType_BTC.ToDasLockCodeHashIndexType())
	case celltype.ChainType_TRON:
		return address.DasLockScript(celltype.DasLockCodeHashIndexType_TRON_Normal)
	}
	return nil, fmt.Errorf("unknow chain type:%d %s", chainType, address)
}

func PayTypeLockScripts(sysWallet *ckbTypes.Script, sysScripts *utils.SystemScripts, payType celltype.ChainType, address Address) PayTypeLockScriptsRet {
	var (
		feeCellProviderScript *ckbTypes.Script
		userScript            *ckbTypes.Script
		scriptType            celltype.LockScriptType
		dasLockParam          celltype.DasLockArgsPairParam
		dasLockParam712       celltype.DasLockArgsPairParam
		err                   error
	)
	if payType == celltype.ChainType_CKB {
		feeCellProviderScript, err = address.CKBLockScript(sysScripts.SecpSingleSigCell.CellHash)
	} else {
		feeCellProviderScript = sysWallet // other coins
	}
	switch payType {
	case celltype.ChainType_CKB:
		scriptType = celltype.ScriptType_User
		indexType := scriptType.ToDasLockCodeHashIndexType()
		userScript, err = address.DasLockScript(indexType)
		dasLockParam = celltype.DasLockArgsPairParam{HashIndexType: scriptType.ToDasLockCodeHashIndexType(), Script: *userScript}
		dasLockParam712 = dasLockParam
		break
	case celltype.ChainType_ETH:
		scriptType = celltype.ScriptType_ETH
		indexType := scriptType.ToDasLockCodeHashIndexType()
		userScript, err = address.DasLockScript(indexType)
		dasLockParam = celltype.DasLockArgsPairParam{HashIndexType: indexType, Script: *userScript}

		indexType712 := scriptType.ToDasLockCodeHashIndexType712()
		userScript712, _ := address.DasLockScript(indexType712)
		dasLockParam712 = celltype.DasLockArgsPairParam{
			HashIndexType: indexType712,
			Script:        *userScript712,
		}
		break
	case celltype.ChainType_TRON:
		scriptType = celltype.ScriptType_TRON
		indexType := scriptType.ToDasLockCodeHashIndexType()
		userScript, err = address.DasLockScript(indexType)
		dasLockParam = celltype.DasLockArgsPairParam{HashIndexType: indexType, Script: *userScript}
		dasLockParam712 = dasLockParam
		break
	case celltype.ChainType_BTC:
		scriptType = celltype.ScriptType_BTC
		indexType := scriptType.ToDasLockCodeHashIndexType()
		userScript, err = address.DasLockScript(indexType)
		dasLockParam = celltype.DasLockArgsPairParam{HashIndexType: indexType, Script: *userScript}
		dasLockParam712 = dasLockParam
		break
	default:
		return PayTypeLockScriptsRet{nil, nil, -1, dasLockParam, dasLockParam712, errors.New("invalid payType")}
	}
	return PayTypeLockScriptsRet{feeCellProviderScript, userScript, scriptType, dasLockParam, dasLockParam712, err}
}

func GetScriptTypeFromLockScript(ckbSysScript *utils.SystemScripts, lockScript *ckbTypes.Script) (celltype.LockScriptType, error) {
	lockCodeHash := lockScript.CodeHash
	switch lockCodeHash {
	case ckbSysScript.SecpSingleSigCell.CellHash:
		return celltype.ScriptType_User, nil
	case celltype.DasAnyOneCanSendCellInfo.Out.CodeHash:
		return celltype.ScriptType_Any, nil
	case celltype.DasETHLockCellInfo.Out.CodeHash:
		return celltype.ScriptType_ETH, nil
	case celltype.DasBTCLockCellInfo.CodeHash:
		return celltype.ScriptType_BTC, nil
	default:
		return -1, errors.New("invalid lockScript")
	}
}

type ReqFindTargetTypeScriptParam struct {
	Ctx       context.Context
	RpcClient rpc.Client
	InputList []*ckbTypes.CellInput
	IsLock    bool
	CodeHash  ckbTypes.Hash
}
type FindTargetTypeScriptRet struct {
	Output        *ckbTypes.CellOutput
	Data          []byte
	Tx            *ckbTypes.Transaction
	PreviousIndex uint
}

func FindTargetTypeScriptByInputList(p *ReqFindTargetTypeScriptParam) (*FindTargetTypeScriptRet, error) {
	codeHash := p.CodeHash
	for _, item := range p.InputList {
		tx, err := p.RpcClient.GetTransaction(p.Ctx, item.PreviousOutput.TxHash)
		if err != nil {
			return nil, fmt.Errorf("FindSenderLockScriptByInputList err: %s", err.Error())
		}
		size := len(tx.Transaction.Outputs)
		for i := 0; i < size; i++ {
			output := tx.Transaction.Outputs[i]
			if p.IsLock {
				if output.Lock != nil && output.Lock.CodeHash == codeHash &&
					output.Lock.HashType == ckbTypes.HashTypeType && item.PreviousOutput.Index == uint(i) {
					return &FindTargetTypeScriptRet{
						Output:        output,
						Data:          tx.Transaction.OutputsData[i],
						Tx:            tx.Transaction,
						PreviousIndex: item.PreviousOutput.Index,
					}, nil
				}
			} else {
				if output.Type != nil &&
					output.Type.CodeHash == codeHash &&
					output.Type.HashType == ckbTypes.HashTypeType &&
					item.PreviousOutput.Index == uint(i) {
					return &FindTargetTypeScriptRet{
						Output:        output,
						Data:          tx.Transaction.OutputsData[i],
						Tx:            tx.Transaction,
						PreviousIndex: item.PreviousOutput.Index,
					}, nil
				}
			}
		}
	}
	return nil, errors.New("FindSenderLockScriptByInputList not found")
}
