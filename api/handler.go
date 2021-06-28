package api

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/gotype"
	"time"

	"das_account_indexer/types"

	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/common"
	"github.com/DeAccountSystems/das_commonlib/common/dascode"
	"github.com/eager7/elog"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
)

/**
 * Copyright (C), 2019-2021
 * FileName: handler
 * Author:   LinGuanHong
 * Date:     2021/4/1 4:11
 * Description:
 */

var (
	log      = elog.NewLogger("rpc_handler", elog.NoticeLevel)
	emptyErr = errors.New("not exist")
)

const maxAccountNumber = 100000

type RpcHandler struct {
	rpcClient     rpc.Client
	systemScripts *utils.SystemScripts
}

func NewRpcHandler(client rpc.Client) *RpcHandler {
	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		panic(fmt.Errorf("init NewSystemScripts err: %s", err.Error()))
	}
	return &RpcHandler{rpcClient: client, systemScripts: systemScripts}
}

func (r *RpcHandler) Hello() string {
	return "hi"
}

func (r *RpcHandler) SearchAccount(ctx context.Context, account string) common.ReqResp {
	log.Info("accept SearchAccount:", account)
	timeStart := time.Now()
	dasAccount := celltype.DasAccountFromStr(account)
	if err := dasAccount.ValidErr(); err != nil {
		return common.ReqResp{ErrNo: dascode.Err_AccountFormatInvalid, ErrMsg: err.Error()}
	}
	accountInfo, err := r.loadOneAccountCellById(dasAccount.AccountId())
	if err != nil {
		if err == emptyErr {
			return common.ReqResp{ErrNo: dascode.Err_AccountNotExist, ErrMsg: err.Error()}
		}
		return common.ReqResp{ErrNo: dascode.Err_BaseParamInvalid, ErrMsg: fmt.Sprintf("loadOneAccountCellById err: %s", err.Error())}
	}
	log.Info("time spend:", time.Since(timeStart).String())
	return common.ReqResp{ErrNo: dascode.DAS_SUCCESS, Data: accountInfo}
}

func (r *RpcHandler) GetAddressAccount(address string) common.ReqResp {
	log.Info("accept GetAddressAccount:", address)
	accountInfo, err := r.loadOneAccountCellByLockScript(gotype.Address(address))
	if err != nil {
		if err == emptyErr {
			return common.ReqResp{ErrNo: dascode.Err_AccountNotExist, ErrMsg: err.Error()}
		}
		return common.ReqResp{ErrNo: dascode.Err_BaseParamInvalid, ErrMsg: fmt.Sprintf("loadOneAccountCellByLockScript err: %s", err.Error())}
	}
	return common.ReqResp{ErrNo: dascode.DAS_SUCCESS, Data: accountInfo}
}

func (r *RpcHandler) loadOneAccountCellByLockScript(address gotype.Address) ([]*types.AccountReturnObj, error) {
	addrLockScriptOwnerArgs, err := address.HexBys(r.systemScripts.SecpSingleSigCell.CellHash)
	if err != nil {
		return nil, fmt.Errorf("LockScript err: %s", err.Error())
	}
	searchKey := &indexer.SearchKey{
		Script:     celltype.DasAccountCellScript.Out.Script(),
		ScriptType: indexer.ScriptTypeType,
	}
	liveCells, _, err := common.LoadLiveCellsWithSize(r.rpcClient, searchKey, maxAccountNumber*celltype.AccountCellBaseCap, maxAccountNumber, true, false, func(cell *indexer.LiveCell) bool {
		ownerBytes := cell.Output.Lock.Args[1 : celltype.DasLockArgsMinBytesLen/2]
		return bytes.Compare(ownerBytes, addrLockScriptOwnerArgs) == 0
	})
	if len(liveCells) == 0 {
		return nil, emptyErr
	}
	accountList := []*types.AccountReturnObj{}
	liveCellLen := len(liveCells)
	log.Info("total accounts:", liveCellLen)
	for i := 0; i < liveCellLen; i++ {
		account, err := r.parseLiveCellToAccount(&liveCells[i], func(cellData *celltype.AccountCellData) bool {
			return true
		})
		if err != nil {
			log.Error("parseLiveCellToAccount err:", err.Error())
			continue
		}
		if account != nil {
			accountList = append(accountList, account)
		}
	}
	return accountList, err
}

func (r *RpcHandler) loadOneAccountCellById(targetAccountId celltype.DasAccountId) (*types.AccountReturnObj, error) {
	searchKey := &indexer.SearchKey{
		Script:     celltype.DasAccountCellScript.Out.Script(),
		ScriptType: indexer.ScriptTypeType,
	}
	liveCells, _, err := common.LoadLiveCells(r.rpcClient, searchKey, celltype.AccountCellBaseCap*2, true, false, func(cell *indexer.LiveCell) bool {
		min := celltype.HashBytesLen + len(celltype.EmptyAccountId)*2
		accountId, err := celltype.AccountIdFromOutputData(cell.OutputData)
		if err != nil {
			return false
		}
		nextAccountId, err := celltype.NextAccountIdFromOutputData(cell.OutputData)
		if err != nil {
			return false
		}
		return len(cell.OutputData) > min && accountId.Compare(targetAccountId) == 0 && accountId.Compare(nextAccountId) != 0
	})
	if err != nil {
		log.Error("LoadLiveCells err:", err.Error())
		return nil, err
	}
	if len(liveCells) == 0 {
		return nil, emptyErr
	}
	return r.parseLiveCellToAccount(&liveCells[0], func(cellData *celltype.AccountCellData) bool {
		if targetAccountId != celltype.AccountCharsToAccount(*cellData.Account()).AccountId() {
			return false
		}
		return true
	})
}

func (r *RpcHandler) parseLiveCellToAccount(cell *indexer.LiveCell, filter func(cellData *celltype.AccountCellData) bool) (*types.AccountReturnObj, error) {
	// get witness
	rawTx, err := r.rpcClient.GetTransaction(context.TODO(), cell.OutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("get raw tx err: %s", err.Error())
	}
	if len(rawTx.Transaction.Witnesses) == 0 {
		return nil, fmt.Errorf("invalid accountCell witness data,it empty, txHash: %s", rawTx.Transaction.Hash.String())
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
			accountCellData, err := gotype.VersionCompatibleAccountCellDataFromSlice(witnessObj.MoleculeNewDataEntity)
			if err != nil {
				return false, fmt.Errorf("VersionCompatibleAccountCellDataFromSlice err: %s", err.Error())
			}
			if filter != nil && !filter(accountCellData) {
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
	if thisAccountCellData == nil {
		return nil, errors.New("skip this account, not match filter")
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
	return &types.AccountReturnObj{
		OutPoint:   *cell.OutPoint,
		WitnessHex: hex.EncodeToString(witnessData),
		AccountData: types.AccountData{
			Account:          celltype.AccountCharsToAccount(*thisAccountCellData.Account()).Str(),
			AccountIdHex:     celltype.DasAccountIdFromBytes(thisAccountCellData.Id().RawData()).HexStr(),
			NextAccountIdHex: nextAccountId.HexStr(),
			CreateAtUnix:     registerAt,
			ExpiredAtUnix:    uint64(expiredAt),
			Status:           celltype.AccountCellStatus(accountStatus),
			Records:          celltype.MoleculeRecordsToGo(*thisAccountCellData.Records()),
		},
	}, err
}
