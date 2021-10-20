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

const genesisAccountIdHex = "0x0000000000000000000000000000000000000000"

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
		oldList, err := types.AccountReturnObjListFromBys(&jsonArrBys)
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

// func storeAccountInfoToRocksDb(db *gorocksdb.DB, writeBatch *gorocksdb.WriteBatch, accountList types.AccountReturnObjList) (int, error) {
// 	accountSize := len(accountList)
// 	sameOwnerMap := map[string]types.AccountReturnObjList{}
// 	accountIdOwnerMap := map[string]string{}
// 	// make group and replace the AccountKey_AccountId data
// 	for i := 0; i < accountSize; i++ {
// 		item := accountList[i]
// 		jsonBys := item.JsonBys()
// 		accountId := item.AccountData.AccountId()
// 		writeBatch.Put(AccountKey_AccountId(accountId), jsonBys)
// 		ownerLockArgsKey := AccountKey_OwnerArgHex(item.AccountData.OwnerLockArgsHex)
// 		ownerHexKey := hex.EncodeToString(ownerLockArgsKey)
// 		accountIdOwnerMap[accountId.Str()] = ownerHexKey
// 		if preList := sameOwnerMap[ownerHexKey]; len(preList) > 0 {
// 			preList = append(preList, item)
// 			sameOwnerMap[ownerHexKey] = preList
// 			for _, o := range preList {
// 				fmt.Println("---->",o.AccountData.OwnerLockArgsHex,string(o.JsonBys()))
// 			}
// 		} else {
// 			newList := types.AccountReturnObjList{}
// 			newList = append(newList, item)
// 			sameOwnerMap[ownerHexKey] = newList
// 		}
// 	}
// 	// replace owner data
// 	for ownerHexKey, ownerItemList := range sameOwnerMap {
// 		ownerLockArgsKey, _ := hex.DecodeString(ownerHexKey)
// 		jsonArrBys, err := rocksdb.RocksDbSafeGet(db, ownerLockArgsKey)
// 		if err != nil {
// 			return 0, fmt.Errorf("RocksDbSafeGet err: %s", err.Error())
// 		} else if jsonArrBys == nil {
// 			writeBatch.Put(ownerLockArgsKey, ownerItemList.JsonBys())
// 		} else {
// 			oldList, err := types.AccountReturnObjListFromBys(jsonArrBys)
// 			if err != nil {
// 				return 0, fmt.Errorf("AccountReturnObjListFromBys err: %s", err.Error())
// 			}
// 			newList := types.AccountReturnObjList{}
// 			ownerListMap := ownerItemList.ToAccountIdMap()
// 			oldListSize := len(oldList)
// 			for i := 0; i < oldListSize; i++ {
// 				storeItem := oldList[i]
// 				mapKey := storeItem.AccountData.AccountIdHex
// 				if newItem := ownerListMap[mapKey]; newItem.AccountData.Account != "" {
// 					// oldList exist, use the new record
// 					storeItem = newItem
// 					delete(ownerListMap, mapKey)
// 				} else if newItem.AccountData.AccountIdHex == genesisAccountIdHex {
// 					// genesis account storage
// 					storeItem = newItem
// 					delete(ownerListMap, mapKey)
// 				} else if accountIdOwnerMap[storeItem.AccountData.AccountId().Str()] == ownerHexKey {
// 					// not exist, keep the old record
// 				} else {
// 					continue // this account has change owner
// 				}
// 				newList = append(newList, storeItem) // here
// 			}
// 			for _, absNewItem := range ownerListMap {
// 				newList = append(newList, absNewItem)
// 			}
// 			if len(newList) == 0 {
// 				continue
// 			}
// 			// print log
// 			for _, newItem := range newList {
// 				log.Info(fmt.Sprintf(
// 					"storeAccountInfoToRocksDb, add new item, account: %s, id: %s, owner: %s",
// 					newItem.AccountData.Account, newItem.AccountData.AccountIdHex, newItem.AccountData.OwnerLockArgsHex))
// 			}
// 			writeBatch.Put(ownerLockArgsKey, newList.JsonBys())
// 		}
// 	}
// 	return accountSize, nil
// }

func storeAccountInfoToRocksDb(db *gorocksdb.DB, writeBatch *gorocksdb.WriteBatch, accountList types.AccountReturnObjList) (int, error) {
	accountSize := len(accountList)
	sameOwnerMap := map[string]*types.AccountReturnObjList{}
	putsItem := func(ownerLockArgsHexKey []byte, currentList *types.AccountReturnObjList) {
		ownerHexKey := hex.EncodeToString(ownerLockArgsHexKey)
		if preList := sameOwnerMap[ownerHexKey]; preList != nil {
			*currentList = append(*currentList, *preList...)
		}
		writeBatch.Put(ownerLockArgsHexKey, (*currentList).JsonBys())
		sameOwnerMap[ownerHexKey] = currentList
	}
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
			putsItem(ownerLockArgsHexKey, &dbList)
		} else {
			oldList, err := types.AccountReturnObjListFromBys(&jsonArrBys)
			if err != nil {
				return 0, fmt.Errorf("AccountReturnObjListFromBys err: %s", err.Error())
			}
			oldListSize := len(oldList)
			newList := types.AccountReturnObjList{}
			for i := 0; i < oldListSize; i++ {
				if oldList[i].AccountData.AccountIdHex != item.AccountData.AccountIdHex { // skip old record
					// newList = append(newList, oldList[i])
				}
			}
			// newList = append(newList, item)
			log.Info(fmt.Sprintf(
				"storeAccountInfoToRocksDb, add new item, account: %s, id: %s, owner: %s",
				item.AccountData.Account, item.AccountData.AccountIdHex, item.AccountData.OwnerLockArgsHex))
			putsItem(ownerLockArgsHexKey, &newList)
			jsonArrBys = nil
		}
	}
	return accountSize, nil
}
