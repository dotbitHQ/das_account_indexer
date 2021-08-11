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

func HandleExpiredRecycleAccountTx(actionName string, p *DASActionHandleFuncParam) DASActionHandleFuncResp {
	var (
		tx         = p.Base.Tx
		resp       = NewDASActionHandleFuncResp("HandleRecycleAccountTx")
		writeOpt   = gorocksdb.NewDefaultWriteOptions()
		writeBatch = gorocksdb.NewWriteBatch()
	)
	defer func() {
		writeOpt.Destroy()
		writeBatch.Destroy()
	}()
	param := &gotype.ReqFindTargetTypeScriptParam{
		Ctx:       context.TODO(),
		RpcClient: p.RpcClient,
		InputList: tx.Inputs[:],
		IsLock:    false,
		CodeHash:  celltype.DasAccountCellScript.Out.CodeHash,
	}
	ret, err := gotype.FindTargetTypeScriptByInputList(param)
	if err != nil {
		return resp.SetErr(fmt.Errorf("isRecycleAccountTx err: invalid isRecycleAccountTx"))
	}
	if len(ret.Tx.Witnesses) == 0 {
		return resp.SetErr(fmt.Errorf("invalid recycleAccount, witness data is empty"))
	}
	accountList, err := util.ParseChainAccountToJsonFormat(ret.Tx, func(cellData *celltype.AccountCellData, outputIndex uint32) bool {
		return outputIndex == uint32(tx.Inputs[0].PreviousOutput.Index)
	})
	if err != nil {
		return resp.SetErr(fmt.Errorf("ParseChainAccountToJsonFormat err: %s", err.Error()))
	}
	accountSize, err := deleteAccountInfoToRocksDb(p.Rocksdb, writeBatch, accountList)
	if err != nil {
		return resp.SetErr(fmt.Errorf("deleteAccountInfoToRocksDb err: %s", err.Error()))
	}
	if err = p.Rocksdb.Write(writeOpt, writeBatch); err != nil {
		resp.Rollback = true
		return resp.SetErr(fmt.Errorf("rocksdb write data err: %s", err.Error()))
	} else {
		log.Info("HandleRecycleAccountTx storage success, total number:", accountSize)
	}
	return resp
}
