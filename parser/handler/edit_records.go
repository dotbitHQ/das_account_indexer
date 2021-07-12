package handler

import (
	"das_account_indexer/util"
	"fmt"
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

func HandleEditRecordsTx(actionName string, p *DASActionHandleFuncParam) DASActionHandleFuncResp {
	var (
		tx         = p.Base.Tx
		resp       = NewDASActionHandleFuncResp("HandleEditRecordsTx")
		writeOpt   = gorocksdb.NewDefaultWriteOptions()
		writeBatch = gorocksdb.NewWriteBatch()
	)
	defer func() {
		writeOpt.Destroy()
		writeBatch.Destroy()
	}()
	if gotype.IsEditRecordsTx(tx) {
		return resp.SetErr(fmt.Errorf("IsEditRecordsTx err: invalid editManagerTx"))
	}
	accountList, err := util.ParseChainAccountToJsonFormat(&tx, nil)
	if err != nil {
		return resp.SetErr(fmt.Errorf("ParseChainAccountToJsonFormat err: %s", err.Error()))
	}
	accountSize, err := storeAccountInfoToRocksDb(p.Rocksdb, writeBatch, accountList)
	if err != nil {
		return resp.SetErr(fmt.Errorf("storeAccountInfoToRocksDb err: %s", err.Error()))
	}
	if err = p.Rocksdb.Write(writeOpt, writeBatch); err != nil {
		resp.Rollback = true
		return resp.SetErr(fmt.Errorf("rocksdb write data err: %s", err.Error()))
	} else {
		log.Info("HandleEditRecordsTx storage success, total number:", accountSize)
	}
	return resp
}
