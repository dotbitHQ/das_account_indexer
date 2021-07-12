package handler

import (
	"encoding/hex"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
)

/**
 * Copyright (C), 2019-2021
 * FileName: keys
 * Author:   LinGuanHong
 * Date:     2021/7/10 5:51
 * Description:
 */

var (
	AccountKey_AccountId = func(accountId celltype.DasAccountId) []byte {
		return append([]byte("account_id_"), accountId.Bytes()...)
	}

	AccountKey_OwnerArgHex = func(ownerArgHex string) []byte {
		ownerArgBys, _ := hex.DecodeString(ownerArgHex)
		return append([]byte("account_owner_address_"), ownerArgBys...)
	}

	AccountKey_OwnerArgHex_Bys = func(ownerArgBys []byte) []byte {
		return append([]byte("account_owner_address_"), ownerArgBys...)
	}

	// AccountKey_Record = func(address gotype.Address) []byte {
	// 	return append([]byte("account_owner_address_"),[]byte(address.OriginStr())...)
	// }
)
