package handler

import (
	"das_account_indexer/types"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/common/rocksdb"
	"github.com/tecbot/gorocksdb"
)

/**
 * Copyright (C), 2019-2021
 * FileName: common
 * Author:   LinGuanHong
 * Date:     2021/7/11 6:15
 * Description:
 */

func storeAccountInfoToRocksDb(db *gorocksdb.DB, writeBatch *gorocksdb.WriteBatch, accountList types.AccountReturnObjList) (int, error) {
	accountSize := len(accountList)
	for i := 0; i < accountSize; i++ {
		item := accountList[i]
		jsonBys := item.AccountData.JsonBys()
		writeBatch.Put(AccountKey_AccountId(item.AccountData.AccountId()), jsonBys)
		fmt.Println("OwnerLockArgsHex:", item.AccountData.OwnerLockArgsHex)
		ownerLockArgsHexKey := AccountKey_OwnerArgHex(item.AccountData.OwnerLockArgsHex)
		jsonArrBys, err := rocksdb.RocksDbSafeGet(db, ownerLockArgsHexKey)
		if err != nil {
			return 0, fmt.Errorf("RocksDbSafeGet err: %s", err.Error())
		} else if jsonArrBys == nil {
			dbList := types.AccountReturnObjList{}
			dbList = append(dbList, item)
			writeBatch.Put(ownerLockArgsHexKey, dbList.JsonBys())
		} else {
			oldList, err := types.AccountReturnObjListFromBys(jsonArrBys)
			if err != nil {
				return 0, fmt.Errorf("AccountReturnObjListFromBys err: %s", err.Error())
			}
			oldList = append(oldList, item)
			writeBatch.Put(ownerLockArgsHexKey, oldList.JsonBys())
		}
	}
	return accountSize, nil
}
