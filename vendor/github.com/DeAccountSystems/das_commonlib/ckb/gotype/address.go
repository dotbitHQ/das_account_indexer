package gotype

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/ckb/wallet"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

/**
 * Copyright (C), 2019-2021
 * FileName: address
 * Author:   LinGuanHong
 * Date:     2021/3/1 10:11
 * Description:
 */

type Address string

func (r Address) Str() string {
	return strings.ToLower(string(r))
}

func (r Address) HexBys(singleSigCellHash types.Hash) ([]byte, error) {
	addrStr := r.Str()
	if len(addrStr) < 3 {
		return nil, errors.New("invalid address")
	}
	switch addrStr[:2] {
	case "ck":
		script, err := r.CKBLockScript(singleSigCellHash)
		if err != nil {
			return nil, err
		}
		return script.Args, nil
	case "0x":
		args, err := hex.DecodeString(addrStr[2:])
		if err != nil {
			return nil, err
		}
		return args, nil
	case "41":
		args, err := hex.DecodeString(addrStr[2:])
		if err != nil {
			return nil, err
		}
		return args, nil
	default:
		return nil, errors.New("unSupport chain address")
	}
}

func (r Address) DasLockScript_CKB() (*types.Script, error) {
	argHex, err := wallet.GetLockScriptArgsFromShortAddress(r.Str())
	if err != nil {
		return nil, fmt.Errorf("GetLockScriptArgsFromShortAddress err:%s", err.Error())
	}
	argBys, _ := hex.DecodeString(argHex)
	indexBytes := celltype.DasLockCodeHashIndexType_CKB_Normal.Bytes()
	bytes := append(indexBytes, argBys...)
	bytes = append(bytes, indexBytes...)
	bytes = append(bytes, argBys...)
	return &types.Script{
		CodeHash: celltype.DasLockCellScript.Out.CodeHash,
		HashType: types.HashTypeType,
		Args:     argBys,
	}, nil
}

func (r Address) CKBLockScript(singleSigCellHash types.Hash) (*types.Script, error) {
	argHex, err := wallet.GetLockScriptArgsFromShortAddress(r.Str())
	if err != nil {
		return nil, fmt.Errorf("GetLockScriptArgsFromShortAddress err:%s", err.Error())
	}
	argBys, _ := hex.DecodeString(argHex)
	return &types.Script{
		CodeHash: singleSigCellHash,
		HashType: types.HashTypeType,
		Args:     argBys,
	}, nil
}

func (r Address) DasLockScript(indexType celltype.DasLockCodeHashIndexType) (*types.Script, error) {
	addrStr := r.Str()
	if strings.HasPrefix(addrStr, "0x") {
		addrStr = addrStr[2:]
	} else if strings.HasPrefix(addrStr, "41") {
		addrStr = addrStr[2:]
	}
	argBys, err := hex.DecodeString(addrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid eth address err:%s", err.Error())
	}
	indexBytes := indexType.Bytes()
	bytes := append(indexBytes, argBys...)
	bytes = append(bytes, indexBytes...)
	bytes = append(bytes, argBys...)
	if len(bytes) != celltype.DasLockArgsMinBytesLen {
		return nil, errors.New("invalid das-lock args")
	}
	return &types.Script{
		CodeHash: celltype.DasLockCellScript.Out.CodeHash,
		HashType: types.HashTypeType,
		Args:     bytes,
	}, nil
}

func (r Address) BTCLockScript() (*types.Script, error) {
	// todo
	return &types.Script{
		CodeHash: celltype.DasBTCLockCellInfo.CodeHash,
		HashType: types.HashTypeType,
		Args:     nil,
	}, nil
}
