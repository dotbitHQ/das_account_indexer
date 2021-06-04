package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/DA-Services/das_commonlib/ckb/celltype"
	"github.com/DA-Services/das_commonlib/ckb/wallet"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

/**
 * Copyright (C), 2019-2021
 * FileName: address
 * Author:   LinGuanHong
 * Date:     2021/5/6 10:45
 * Description:
 */
type Address string

func (r Address) Str() string {
	return strings.ToLower(string(r))
}

func (r Address) LockScript(singleSigCellHash types.Hash) (*types.Script, error) {
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
		return script, nil
	case "0x":
		script, err := r.ETHLockScript()
		if err != nil {
			return nil, err
		}
		return script, nil
	default:
		return nil, errors.New("unSupport chain address")
	}
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

func (r Address) ETHLockScript() (*types.Script, error) {
	addrStr := r.Str()
	if strings.HasPrefix(addrStr, "0x") {
		addrStr = addrStr[2:]
	}
	argBys, err := hex.DecodeString(addrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid eth address err:%s", err.Error())
	}
	return &types.Script{
		CodeHash: celltype.DasETHLockCellInfo.Out.CodeHash,
		HashType: types.HashTypeType,
		Args:     argBys,
	}, nil
}
