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

func deleteAccountInfoToRocksDb(db *gorocksdb.DB, writeBatch *gorocksdb.WriteBatch, accountList types.AccountReturnObjList) (int, error) {
	accountSize := len(accountList)
	for i := 0; i < accountSize; i++ {
		item := accountList[i]
		writeBatch.Delete(AccountKey_AccountId(item.AccountData.AccountId()))
		if err := removeItemFromOwnerList(db, writeBatch, &item); err != nil {
			return 0, fmt.Errorf("removeItemFromOwnerList err: %s", err.Error())
		}
	}
	return accountSize, nil
}

func removeItemFromOwnerList(db *gorocksdb.DB, writeBatch *gorocksdb.WriteBatch, item *types.AccountReturnObj) error {
	ownerLockArgsHexKey := AccountKey_OwnerArgHex(item.AccountData.OwnerLockArgsHex)
	jsonArrBys, err := rocksdb.RocksDbSafeGet(db, ownerLockArgsHexKey)
	if err != nil {
		return fmt.Errorf("RocksDbSafeGet err: %s", err.Error())
	} else if jsonArrBys != nil {
		oldList, err := types.AccountReturnObjListFromBys(jsonArrBys)
		if err != nil {
			return fmt.Errorf("AccountReturnObjListFromBys err: %s", err.Error())
		}
		oldListSize := len(oldList)
		newList := types.AccountReturnObjList{}
		for i := 0; i < oldListSize; i++ {
			oldItem := oldList[i]
			if oldItem.AccountData.AccountIdHex == item.AccountData.AccountIdHex {
				continue
			}
			newList = append(newList, oldItem)
		}
		writeBatch.Put(ownerLockArgsHexKey, newList.JsonBys())
	}
	return nil
}

func storeAccountInfoToRocksDb(db *gorocksdb.DB, writeBatch *gorocksdb.WriteBatch, accountList types.AccountReturnObjList) (int, error) {
	accountSize := len(accountList)
	for i := 0; i < accountSize; i++ {
		item := accountList[i]
		jsonBys := item.JsonBys()
		writeBatch.Put(AccountKey_AccountId(item.AccountData.AccountId()), jsonBys)
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
			oldListSize := len(oldList)
			newList := types.AccountReturnObjList{}
			for i := 0; i < oldListSize; i++ {
				if oldList[i].AccountData.AccountIdHex != item.AccountData.AccountIdHex {
					newList = append(newList, oldList[i])
				}
			}
			log.Info(fmt.Sprintf(
				"storeAccountInfoToRocksDb, add new item, account: %s, id: %s",
				item.AccountData.Account, item.AccountData.AccountIdHex))
			newList = append(newList, item)
			writeBatch.Put(ownerLockArgsHexKey, newList.JsonBys())
		}
	}
	return accountSize, nil
}
