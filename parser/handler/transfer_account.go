package handler

import (
	"context"
	"das_account_indexer/util"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/ckb/gotype"
	"github.com/tecbot/gorocksdb"
)

/**
 * Copyright (C), 2019-2021
 * FileName: confirm_propose
 * Author:   LinGuanHong
 * Date:     2021/7/10 4:50
 * Description:
 */

func HandleTransferAccountTx(actionName string, p *DASActionHandleFuncParam) DASActionHandleFuncResp {
	var (
		tx         = p.Base.Tx
		resp       = NewDASActionHandleFuncResp("HandleTransferAccountTx")
		writeOpt   = gorocksdb.NewDefaultWriteOptions()
		writeBatch = gorocksdb.NewWriteBatch()
	)
	defer func() {
		writeOpt.Destroy()
		writeBatch.Destroy()
	}()
	if !gotype.IsTransferAccountTx(tx) {
		return resp.SetErr(fmt.Errorf("IsTransferAccountTx err: invalid editManagerTx"))
	}

	// delete old account
	param := &gotype.ReqFindTargetTypeScriptParam{
		Ctx:       context.TODO(),
		RpcClient: p.RpcClient,
		InputList: tx.Inputs[:],
		IsLock:    false,
		CodeHash:  celltype.DasAccountCellScript.Out.CodeHash,
	}
	ret, err := gotype.FindTargetTypeScriptByInputList(param)
	if err != nil {
		return resp.SetErr(fmt.Errorf("IsTransferAccountTx err: invalid isRecycleAccountTx"))
	}
	if len(ret.Tx.Witnesses) == 0 {
		return resp.SetErr(fmt.Errorf("invalid transferAccount, witness data is empty"))
	}
	accountListOld, err := util.ParseChainAccountToJsonFormat(ret.Tx, func(cellData *celltype.AccountCellData, outputIndex uint32) bool {
		return uint32(ret.PreviousIndex) == outputIndex
	})
	if err != nil {
		return resp.SetErr(fmt.Errorf("ParseChainAccountToJsonFormat err: %s", err.Error()))
	}

	deleteSize, err := deleteAccountInfoToRocksDb(p.Rocksdb, writeBatch, accountListOld)
	if err != nil {
		return resp.SetErr(fmt.Errorf("deleteAccountInfoToRocksDb err: %s", err.Error()))
	}

	// storage new
	accountListNew, err := util.ParseChainAccountToJsonFormat(&tx, nil)
	if err != nil {
		return resp.SetErr(fmt.Errorf("ParseChainAccountToJsonFormat err: %s", err.Error()))
	}
	accountSizeNew, err := storeAccountInfoToRocksDb(p.Rocksdb, writeBatch, accountListNew)
	if err != nil {
		return resp.SetErr(fmt.Errorf("storeAccountInfoToRocksDb err: %s", err.Error()))
	}

	if deleteSize != accountSizeNew {
		return resp.SetErr(fmt.Errorf("transferAccount err: account number not equal, old: %d, new: %d", deleteSize, accountSizeNew))
	}
	log.Info(fmt.Sprintf(
		"transfer, account: %s, from: %s to: %s",
		accountListOld[0].AccountData.Account, accountListOld[0].AccountData.OwnerLockArgsHex, accountListNew[0].AccountData.OwnerLockArgsHex))
	if err = p.Rocksdb.Write(writeOpt, writeBatch); err != nil {
		resp.Rollback = true
		return resp.SetErr(fmt.Errorf("rocksdb write data err: %s", err.Error()))
	} else {
		log.Info("HandleTransferAccountTx storage success, total number:", accountSizeNew)
	}
	return resp
}
