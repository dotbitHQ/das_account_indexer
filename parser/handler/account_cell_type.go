package handler

import (
	"context"
	"das_account_indexer/types"
	"das_account_indexer/util"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/ckb/gotype"
	"github.com/tecbot/gorocksdb"
	"strings"
)

/**
 * Copyright (C), 2019-2021
 * FileName: account_cell
 * Author:   LinGuanHong
 * Date:     2021/9/8 2:25
 * Description:
 */

func isInputOldAccountCellNotExist(errMsg string) bool {
	return strings.Contains(errMsg, "not found")
}

func HandleAccountCellType(actionName string, p *DASActionHandleFuncParam) DASActionHandleFuncResp {
	var (
		tx         = p.Base.Tx
		resp       = NewDASActionHandleFuncResp("HandleAccountCellType")
		writeOpt   = gorocksdb.NewDefaultWriteOptions()
		writeBatch = gorocksdb.NewWriteBatch()
	)
	defer func() {
		writeOpt.Destroy()
		writeBatch.Destroy()
		log.Info(fmt.Sprintf("---------- End HandleAccountCellType: %s ----------", actionName))
	}()
	log.Info("HandleAccountCellType find action:", actionName)
	// find old accountCell from input
	var accountListOldNumber int
	var accountListOld []types.AccountReturnObj
	if actionName != celltype.Action_ConfirmProposal {
		param := &gotype.ReqFindTargetTypeScriptParam{
			Ctx:       context.TODO(),
			RpcClient: p.RpcClient,
			InputList: tx.Inputs[:],
			IsLock:    false,
			CodeHash:  celltype.DasAccountCellScript.Out.CodeHash,
		}
		inputAccountCellRet, err := gotype.FindTargetTypeScriptByInputList(param)
		if err != nil {
			if !isInputOldAccountCellNotExist(err.Error()) {
				return resp.SetErr(fmt.Errorf("FindTargetTypeScriptByInputList err: %s", err))
			} else {
				// this tx has no input accountCell
			}
		} else {
			if len(inputAccountCellRet.Tx.Witnesses) == 0 {
				return resp.SetErr(fmt.Errorf("invalid HandleAccountCellType, witness data is empty"))
			}
			accountListOld, err = util.ParseChainAccountToJsonFormat(inputAccountCellRet.Tx, func(cellData *celltype.AccountCellData, outputIndex uint32) bool {
				return uint32(inputAccountCellRet.PreviousIndex) == outputIndex
			})
			if err != nil {
				return resp.SetErr(fmt.Errorf("ParseChainAccountToJsonFormat err: %s", err.Error()))
			}
		}
	}
	// try storage new
	accountListNew, err := util.ParseChainAccountToJsonFormat(&tx, nil)
	if err != nil {
		return resp.SetErr(fmt.Errorf("ParseChainAccountToJsonFormat err: %s", err.Error()))
	}
	accountSizeNew := len(accountListNew)
	// try update owner info
	if accountListOldNumber = len(accountListOld); accountListOldNumber > 0 {
		shouldExecuteDelete := false
		if accountSizeNew == 0 {
			// output does not have a new accountCell, this tx maybe is a recycle type tx.
			log.Info("this maybe some kind of recycle account type tx")
			shouldExecuteDelete = true
		} else {
			if accountListOldNumber != accountSizeNew { // for now, DAS accountCell edit type tx's input account number must equal output account number
				return resp.SetErr(fmt.Errorf("accountCell number not equal, old: %d, new: %d", accountListOldNumber, accountSizeNew))
			}
			// The judgment of start with outputs is whether the accountId is equal, but the owner is different
			for i := 0; i < accountListOldNumber; i++ {
				newAccountData := accountListNew[i].AccountData
				oldAccountData := accountListOld[i].AccountData
				if newAccountData.AccountIdHex == oldAccountData.AccountIdHex {
					shouldExecuteDelete = newAccountData.OwnerLockArgsHex != oldAccountData.OwnerLockArgsHex
					break
				}
			}
		}
		if shouldExecuteDelete {
			if _, err = deleteAccountInfoToRocksDb(p.Rocksdb, writeBatch, accountListOld); err != nil {
				return resp.SetErr(fmt.Errorf("deleteAccountInfoToRocksDb err: %s", err.Error()))
			}
			if accountSizeNew > 0 {
				log.Info(fmt.Sprintf(
					"owner transfer happened, account: %s, from: %s to: %s",
					accountListOld[0].AccountData.Account, accountListOld[0].AccountData.OwnerLockArgsHex,
					accountListNew[0].AccountData.OwnerLockArgsHex))

			}
		} else {
			log.Info("no need to delete accountCell info, this maybe some kind of editRecords or renewAccount tx")
		}
	}
	if accountListOldNumber > 0 || accountSizeNew > 0 {
		if _, err := storeAccountInfoToRocksDb(p.Rocksdb, writeBatch, accountListNew); err != nil {
			return resp.SetErr(fmt.Errorf("storeAccountInfoToRocksDb err: %s", err.Error()))
		}
		if err = p.Rocksdb.Write(writeOpt, writeBatch); err != nil {
			resp.Rollback = true
			return resp.SetErr(fmt.Errorf("rocksdb write data err: %s", err.Error()))
		} else {
			log.Info(fmt.Sprintf(
				"HandleAccountCellType success, total number, old: %d, new: %d",
				accountListOldNumber, accountSizeNew))
		}
	}
	return resp
}
