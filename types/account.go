package types

import (
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: account
 * Author:   LinGuanHong
 * Date:     2021/4/1 4:46
 * Description:
 */

type AccountFilterFunc func(cellData *celltype.AccountCellData) bool

type AccountData struct {
	Account           string                      `json:"account"`
	AccountIdHex      string                      `json:"account_id_hex"`
	NextAccountIdHex  string                      `json:"next_account_id_hex"`
	CreateAtUnix      uint64                      `json:"create_at_unix"`
	ExpiredAtUnix     uint64                      `json:"expired_at_unix"`
	Status            celltype.AccountCellStatus  `json:"status"`
	OwnerLockArgsHex  string                      `json:"owner_lock_args_hex"`
	ManagerLockArgHex string                      `json:"manager_lock_arg_hex"`
	Records           celltype.EditRecordItemList `json:"records"`
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

func (a AccountReturnObj) JsonBys() []byte {
	bys, _ := json.Marshal(a)
	return bys
}

type AccountReturnObjList []AccountReturnObj

func (a AccountReturnObjList) JsonBys() []byte {
	bys, _ := json.Marshal(a)
	return bys
}

func AccountReturnObjListFromBys(listBys []byte) (AccountReturnObjList, error) {
	list := &AccountReturnObjList{}
	if err := json.Unmarshal(listBys, list); err != nil {
		return nil, err
	}
	return *list, nil
}
