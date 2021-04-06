package api

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"das_account_indexer/model"

	"github.com/DA-Services/das_commonlib/ckb/celltype"
	"github.com/DA-Services/das_commonlib/common"
	"github.com/DA-Services/das_commonlib/common/dascode"
	"github.com/eager7/elog"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
)

/**
 * Copyright (C), 2019-2021
 * FileName: handler
 * Author:   LinGuanHong
 * Date:     2021/4/1 4:11 下午
 * Description:
 */

var (
	log      = elog.NewLogger("rpc_handler", elog.NoticeLevel)
	emptyErr = errors.New("not exist")
)

type RpcHandler struct {
	rpcClient rpc.Client
}

func NewRpcHandler(client rpc.Client) *RpcHandler {
	return &RpcHandler{rpcClient: client}
}

func (r *RpcHandler) Hello() string {
	return "hi"
}

func (r *RpcHandler) SearchAccount(account string) common.ReqResp {
	log.Info("accept SearchAccount:", account)
	dasAccount := celltype.DasAccountFromStr(account)
	if err := dasAccount.ValidErr(); err != nil {
		return common.ReqResp{ErrNo: dascode.Err_AccountFormatInvalid, ErrMsg: err.Error()}
	}
	accountInfo, err := r.loadOneAccountCell(dasAccount.AccountId())
	if err != nil {
		if err == emptyErr {
			return common.ReqResp{ErrNo: dascode.Err_AccountNotExist, ErrMsg: err.Error()}
		}
		return common.ReqResp{ErrNo: dascode.Err_BaseParamInvalid, ErrMsg: err.Error()}
	}
	return common.ReqResp{ErrNo: dascode.DAS_SUCCESS, Data: accountInfo}
}

func (r *RpcHandler) loadOneAccountCell(targetAccountId celltype.DasAccountId) (*model.AccountReturnObj, error) {
	searchKey := &indexer.SearchKey{
		Script:     celltype.DasAccountCellScript.Out.Script(),
		ScriptType: indexer.ScriptTypeType,
	}
	const accountBytesLen = 10
	liveCells, _, err := common.LoadLiveCells(r.rpcClient, searchKey, 200*celltype.OneCkb, true, false, func(cell *indexer.LiveCell) bool {
		min := celltype.HashBytesLen + accountBytesLen*2
		accountId, err := celltype.AccountIdFromOutputData(cell.OutputData)
		if err != nil {
			return false
		}
		return len(cell.OutputData) > min && accountId.Compare(targetAccountId) == 0
	})
	if len(liveCells) != 1 {
		return nil, emptyErr
	}
	cell := liveCells[0]
	// get witness
	rawTx, err := r.rpcClient.GetTransaction(context.TODO(), cell.OutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("get raw tx err: %s", err.Error())
	}
	if len(rawTx.Transaction.Witnesses) == 0 {
		return nil, fmt.Errorf("invalid accountCell witness data empty, txHash: %s", rawTx.Transaction.Hash.String())
	}
	var (
		thisAccountCellData *celltype.AccountCellData
		witnessData         []byte
	)
	err = celltype.GetTargetCellFromWitness(rawTx.Transaction, func(rawWitnessData []byte, witnessParseObj *celltype.ParseDasWitnessBysDataObj) (bool, error) {
		witnessDataObj := witnessParseObj.WitnessObj
		switch witnessDataObj.TableType {
		case celltype.TableType_ACCOUNT_CELL:
			witnessObj, err := celltype.ParseTxWitnessToDasWitnessObj(rawWitnessData)
			if err != nil {
				return false, fmt.Errorf("ParseTxWitnessToDasWitnessObj err: %s", err.Error())
			}
			accountCellData, err := celltype.AccountCellDataFromSlice(witnessObj.MoleculeNewDataEntity.Entity().RawData(), false)
			if err != nil {
				return false, fmt.Errorf("AccountCellDataFromSlice err: %s", err.Error())
			}
			if targetAccountId != celltype.AccountCharsToAccount(*accountCellData.Account()).AccountId() {
				return false, nil // next one
			}
			thisAccountCellData = accountCellData
			witnessData = rawWitnessData
			return true, nil
		}
		return false, nil
	}, func(err error) {
		log.Warn("GetTargetCellFromWitness [accountCell] skip one item:", err.Error())
	})
	if err != nil {
		return nil, err
	}
	nextAccountId, err := celltype.NextAccountIdFromOutputData(cell.OutputData)
	if err != nil {
		return nil, err
	}
	registerAt, err := celltype.MoleculeU64ToGo(thisAccountCellData.RegisteredAt().RawData())
	if err != nil {
		return nil, fmt.Errorf("parse registerAt err: %s", err.Error())
	}
	expiredAt, err := celltype.ExpiredAtFromOutputData(cell.OutputData)
	if err != nil {
		return nil, fmt.Errorf("parse expiredAt err: %s", err.Error())
	}
	accountStatus, _ := celltype.MoleculeU8ToGo(thisAccountCellData.Status().RawData())
	ownerLock, err := celltype.MoleculeScriptToGo(*thisAccountCellData.OwnerLock())
	if err != nil {
		return nil, fmt.Errorf("parse ownerLock err: %s", err.Error())
	}
	managerLock, err := celltype.MoleculeScriptToGo(*thisAccountCellData.ManagerLock())
	if err != nil {
		return nil, fmt.Errorf("parse managerLock err: %s", err.Error())
	}
	return &model.AccountReturnObj{
		OutPoint:   *cell.OutPoint,
		WitnessHex: hex.EncodeToString(witnessData),
		AccountData: model.AccountData{
			Account:           celltype.AccountCharsToAccount(*thisAccountCellData.Account()).Str(),
			AccountIdHex:      celltype.DasAccountIdFromBytes(thisAccountCellData.Id().RawData()).HexStr(),
			NextAccountIdHex:  nextAccountId.HexStr(),
			CreateAtUnix:      registerAt,
			ExpiredAtUnix:     uint64(expiredAt),
			Status:            celltype.AccountCellStatus(accountStatus),
			OwnerLockScript:   *ownerLock,
			ManagerLockScript: *managerLock,
			Records:           celltype.MoleculeRecordsToGo(*thisAccountCellData.Records()),
		},
	}, err
}
