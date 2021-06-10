package model

import (
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

type AccountData struct {
	Account          string                     `json:"account"`
	AccountIdHex     string                     `json:"account_id_hex"`
	NextAccountIdHex string                     `json:"next_account_id_hex"`
	CreateAtUnix     uint64                     `json:"create_at_unix"`
	ExpiredAtUnix    uint64                     `json:"expired_at_unix"`
	Status           celltype.AccountCellStatus `json:"status"`
	// OwnerLockScript types.Script `json:"owner_lock_script"`
	// ManagerLockScript types.Script `json:"manager_lock_script"`
	Records celltype.EditRecordItemList `json:"records"`
}
type AccountReturnObj struct {
	OutPoint    types.OutPoint `json:"out_point"`
	WitnessHex  string         `json:"-"`
	AccountData AccountData    `json:"account_data"`
}
