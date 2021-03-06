package handler

import (
	"das_account_indexer/types"
	"encoding/hex"
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
		if err := removeItemFromOwnerList(db, writeBatch, item); err != nil {
			return 0, fmt.Errorf("removeItemFromOwnerList err: %s", err.Error())
		}
	}
	return accountSize, nil
}

func removeItemFromOwnerList(db *gorocksdb.DB, writeBatch *gorocksdb.WriteBatch, item types.AccountReturnObj) error {
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
	sameOwnerMap := map[string]types.AccountReturnObjList{}
	// make group and replace the AccountKey_AccountId data
	for i := 0; i < accountSize; i++ {
		item := accountList[i]
		jsonBys := item.JsonBys()
		writeBatch.Put(AccountKey_AccountId(item.AccountData.AccountId()), jsonBys)
		ownerLockArgsKey := AccountKey_OwnerArgHex(item.AccountData.OwnerLockArgsHex)
		ownerHexKey := hex.EncodeToString(ownerLockArgsKey)
		if preList := sameOwnerMap[ownerHexKey]; len(preList) > 0 {
			preList = append(preList, item)
			sameOwnerMap[ownerHexKey] = preList
		} else {
			newList := types.AccountReturnObjList{}
			newList = append(newList, item)
			sameOwnerMap[ownerHexKey] = newList
		}
	}
	// replace owner data
	for ownerHexKey, ownerItemList := range sameOwnerMap {
		ownerLockArgsKey, _ := hex.DecodeString(ownerHexKey)
		jsonArrBys, err := rocksdb.RocksDbSafeGet(db, ownerLockArgsKey)
		if err != nil {
			return 0, fmt.Errorf("RocksDbSafeGet err: %s", err.Error())
		} else if jsonArrBys == nil {
			writeBatch.Put(ownerLockArgsKey, ownerItemList.JsonBys())
		} else {
			oldList, err := types.AccountReturnObjListFromBys(jsonArrBys)
			if err != nil {
				return 0, fmt.Errorf("AccountReturnObjListFromBys err: %s", err.Error())
			}
			newList := types.AccountReturnObjList{}
			ownerListMap := ownerItemList.ToAccountIdMap()
			oldListSize := len(oldList)
			for i := 0; i < oldListSize; i++ {
				storeItem := oldList[i]
				mapKey := storeItem.AccountData.AccountIdHex
				if newItem := ownerListMap[mapKey]; newItem.AccountData.Account != "" {
					// oldList exist, use the new record
					storeItem = newItem
					delete(ownerListMap, mapKey)
				} else {
					// not exist, keep the old record
				}
				newList = append(newList, storeItem)
			}
			for _, absNewItem := range ownerListMap {
				newList = append(newList, absNewItem)
			}
			if len(newList) == 0 {
				continue
			}
			// print log
			for _, newItem := range newList {
				log.Info(fmt.Sprintf(
					"storeAccountInfoToRocksDb, add new item, account: %s, id: %s, owner: %s",
					newItem.AccountData.Account, newItem.AccountData.AccountIdHex, newItem.AccountData.OwnerLockArgsHex))
			}
			writeBatch.Put(ownerLockArgsKey, newList.JsonBys())
		}
	}
	return accountSize, nil
}
