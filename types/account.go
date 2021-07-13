package types

import (
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

/**
 * Copyright (C), 2019-2021
 * FileName: account
 * Author:   LinGuanHong
 * Date:     2021/4/1 4:46
 * Description:
 */

type AccountFilterFunc func(cellData *celltype.AccountCellData) bool
type SimpleRecordItem struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}
type AccountData struct {
	Account           string                      `json:"account"`
	AccountIdHex      string                      `json:"account_id_hex"`
	NextAccountIdHex  string                      `json:"next_account_id_hex"`
	CreateAtUnix      uint64                      `json:"create_at_unix"`
	ExpiredAtUnix     uint64                      `json:"expired_at_unix"`
	Status            celltype.AccountCellStatus  `json:"status"`
	OwnerLockArgsHex  string                      `json:"owner_lock_args_hex"`
	ManagerLockArgHex string                      `json:"manager_lock_arg_hex"`
	Records           celltype.EditRecordItemList `json:"-"`
}
type AccountData1 struct {
	Account           string                     `json:"account"`
	AccountIdHex      string                     `json:"account_id_hex"`
	NextAccountIdHex  string                     `json:"next_account_id_hex"`
	CreateAtUnix      uint64                     `json:"create_at_unix"`
	ExpiredAtUnix     uint64                     `json:"expired_at_unix"`
	Status            celltype.AccountCellStatus `json:"status"`
	OwnerLockArgsHex  string                     `json:"owner_lock_args_hex"`
	ManagerLockArgHex string                     `json:"manager_lock_arg_hex"`
	Records           []SimpleRecordItem         `json:"records"`
}

func (a AccountData) AccountId() celltype.DasAccountId {
	return celltype.DasAccountFromStr(a.Account).AccountId()
}
func (a AccountData) JsonBys() []byte {
	fmt.Println("===>", a)
	bys, err := json.Marshal(a)
	if err != nil {
		fmt.Println(err.Error())
	}
	return bys
}

type AccountReturnObj struct {
	OutPoint    types.OutPoint `json:"out_point"`
	WitnessHex  string         `json:"-"`
	AccountData AccountData    `json:"account_data"`
}

type AccountReturnObj1 struct {
	OutPoint    types.OutPoint `json:"out_point"`
	WitnessHex  string         `json:"-"`
	AccountData AccountData1   `json:"account_data"`
}

func (a AccountReturnObj) JsonBys() []byte {
	bys, _ := json.Marshal(a)
	return bys
}

func (a AccountReturnObj) ToAccountReturnObj1() AccountReturnObj1 {
	return AccountReturnObj1{
		OutPoint:   a.OutPoint,
		WitnessHex: a.WitnessHex,
		AccountData: AccountData1{
			Account:           a.AccountData.Account,
			AccountIdHex:      a.AccountData.AccountIdHex,
			NextAccountIdHex:  a.AccountData.NextAccountIdHex,
			CreateAtUnix:      a.AccountData.CreateAtUnix,
			ExpiredAtUnix:     a.AccountData.ExpiredAtUnix,
			Status:            a.AccountData.Status,
			OwnerLockArgsHex:  appendOx(a.AccountData.OwnerLockArgsHex),
			ManagerLockArgHex: appendOx(a.AccountData.ManagerLockArgHex),
			Records:           originRecordsToNewRecords(a.AccountData.Records),
		},
	}
}

type AccountReturnObjList []AccountReturnObj

func (a AccountReturnObjList) JsonBys() []byte {
	bys, _ := json.Marshal(a)
	return bys
}

func appendOx(hex string) string {
	if !strings.HasPrefix(hex, "0x") {
		return "0x" + hex
	}
	return hex
}

func originRecordsToNewRecords(records celltype.EditRecordItemList) []SimpleRecordItem {
	recordSize := len(records)
	recordList := make([]SimpleRecordItem, 0, recordSize)
	for j := 0; j < recordSize; j++ {
		recordList = append(recordList, SimpleRecordItem{
			Key:   records[j].Type + "." + records[j].Key,
			Label: records[j].Label,
			Value: records[j].Value,
			TTL:   records[j].TTL,
		})
	}
	return recordList
}

func (a AccountReturnObjList) ToAccountReturnObjList1List() []AccountData1 {
	retList := make([]AccountData1, 0, len(a))
	for i := 0; i < len(a); i++ {
		records := a[i].AccountData.Records
		recordSize := len(records)
		recordList := make([]SimpleRecordItem, 0, recordSize)
		for j := 0; j < recordSize; j++ {
			recordList = append(recordList, SimpleRecordItem{
				Key:   records[j].Type + "." + records[j].Key,
				Label: records[j].Label,
				Value: records[j].Value,
				TTL:   records[j].TTL,
			})
		}
		retList = append(retList, a[i].ToAccountReturnObj1().AccountData)
	}
	return retList
}

func AccountReturnObjListFromBys(listBys []byte) (AccountReturnObjList, error) {
	list := &AccountReturnObjList{}
	if err := json.Unmarshal(listBys, list); err != nil {
		return nil, err
	}
	return *list, nil
}
