package types

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/ckb/gotype"
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

// origin
type AccountData struct {
	Account           string                      `json:"account"`
	AccountIdHex      string                      `json:"account_id_hex"`
	NextAccountIdHex  string                      `json:"next_account_id_hex"`
	CreateAtUnix      uint64                      `json:"create_at_unix"`
	ExpiredAtUnix     uint64                      `json:"expired_at_unix"`
	Status            celltype.AccountCellStatus  `json:"status"`
	RawDasLockArgsHex string                      `json:"raw_das_lock_args_hex"`
	OwnerLockArgsHex  string                      `json:"owner_lock_args_hex"`
	ManagerLockArgHex string                      `json:"manager_lock_arg_hex"`
	Records           celltype.EditRecordItemList `json:"records"`
}

// repair
type AccountData1 struct {
	Account             string                     `json:"account"`
	AccountIdHex        string                     `json:"account_id_hex"`
	NextAccountIdHex    string                     `json:"next_account_id_hex"`
	CreateAtUnix        uint64                     `json:"create_at_unix"`
	ExpiredAtUnix       uint64                     `json:"expired_at_unix"`
	Status              celltype.AccountCellStatus `json:"status"`
	OwnerLockChainType  string                     `json:"owner_lock_chain_type"`
	OwnerLockArgsHex    string                     `json:"owner_lock_args_hex"`
	OwnerAddress        string                     `json:"owner_address"`
	ManageLockChainType string                     `json:"manage_lock_chain_type"`
	ManagerAddress      string                     `json:"manager_address"`
	ManagerLockArgsHex  string                     `json:"manager_lock_args_hex"`
	Records             []SimpleRecordItem         `json:"records"`
}

func (a AccountData) AccountId() celltype.DasAccountId {
	return celltype.DasAccountFromStr(a.Account).AccountId()
}
func (a AccountData) JsonBys() []byte {
	bys, _ := json.Marshal(a)
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
	var (
		rawDasLockArgsHex   = a.AccountData.RawDasLockArgsHex
		rawDasLockArgsBytes []byte
	)
	if strings.HasPrefix(rawDasLockArgsHex, "0x") {
		rawDasLockArgsHex = rawDasLockArgsHex[2:]
	}
	rawDasLockArgsBytes, _ = hex.DecodeString(rawDasLockArgsHex)
	if len(rawDasLockArgsBytes) < celltype.DasLockArgsMinBytesLen {
		rawDasLockArgsBytes = bytes.Repeat([]byte{0}, celltype.DasLockArgsMinBytesLen)
	}
	ownerChainType := celltype.ChainType(rawDasLockArgsBytes[0])
	managerChainType := celltype.ChainType(rawDasLockArgsBytes[celltype.DasLockArgsMinBytesLen/2])
	return AccountReturnObj1{
		OutPoint:   a.OutPoint,
		WitnessHex: a.WitnessHex,
		AccountData: AccountData1{
			Account:             a.AccountData.Account,
			AccountIdHex:        a.AccountData.AccountIdHex,
			NextAccountIdHex:    a.AccountData.NextAccountIdHex,
			CreateAtUnix:        a.AccountData.CreateAtUnix,
			ExpiredAtUnix:       a.AccountData.ExpiredAtUnix,
			Status:              a.AccountData.Status,
			OwnerAddress:        gotype.PubkeyHashToAddress(ownerChainType, removeOx(a.AccountData.OwnerLockArgsHex)).OriginStr(),
			OwnerLockArgsHex:    appendOx(a.AccountData.OwnerLockArgsHex),
			OwnerLockChainType:  ownerChainType.String(),
			ManagerAddress:      gotype.PubkeyHashToAddress(managerChainType, removeOx(a.AccountData.ManagerLockArgHex)).OriginStr(),
			ManageLockChainType: managerChainType.String(),
			ManagerLockArgsHex:  appendOx(a.AccountData.ManagerLockArgHex),
			Records:             originRecordsToNewRecords(a.AccountData.Records),
		},
	}
}

type AccountReturnObjList []AccountReturnObj

func (a AccountReturnObjList) JsonBys() []byte {
	bys, _ := json.Marshal(a)
	return bys
}

func removeOx(hex string) string {
	if strings.HasPrefix(hex, "0x") {
		return hex[2:]
	}
	return hex
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

func (a AccountReturnObjList) ToAccountReturnObjList1List() []AccountReturnObj1 {
	retList := make([]AccountReturnObj1, 0, len(a))
	for i := 0; i < len(a); i++ {
		retList = append(retList, a[i].ToAccountReturnObj1())
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
