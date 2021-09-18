package api

import (
	"bytes"
	"context"
	"das_account_indexer/util"
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
	testNet       bool
	rpcClient     rpc.Client
	systemScripts *utils.SystemScripts
}

func NewRpcHandler(testNet bool, client rpc.Client) *RpcHandler {
	systemScripts, err := utils.NewSystemScripts(client)
	if err != nil {
		panic(fmt.Errorf("init NewSystemScripts err: %s", err.Error()))
	}
	return &RpcHandler{testNet: testNet, rpcClient: client, systemScripts: systemScripts}
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
	return common.ReqResp{ErrNo: dascode.DAS_SUCCESS, Data: accountInfo.ToAccountReturnObj1(r.testNet)}
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
	return common.ReqResp{ErrNo: dascode.DAS_SUCCESS, Data: accountInfo.ToAccountReturnObjList1List(r.testNet)}
}

func (r *RpcHandler) Close() {

}

func (r *RpcHandler) loadOneAccountCellByLockScript(address gotype.Address) (types.AccountReturnObjList, error) {
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
	accountList := []types.AccountReturnObj{}
	liveCellLen := len(liveCells)
	log.Info("total accounts:", liveCellLen)
	for i := 0; i < liveCellLen; i++ {
		account, err := r.parseLiveCellToAccount(&liveCells[i], func(cellData *celltype.AccountCellData, outputIndex uint32) bool {
			return true
		})
		if err != nil {
			log.Error("parseLiveCellToAccount err:", err.Error())
			continue
		}
		if account != nil {
			accountList = append(accountList, *account)
		}
	}
	return accountList, err
}

func (r *RpcHandler) loadOneAccountCellById(targetAccountId celltype.DasAccountId) (*types.AccountReturnObj, error) {
	timeStart := time.Now()
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
	log.Warn("load account time spend:", time.Since(timeStart).String())
	timeStart1 := time.Now()
	obj, err := r.parseLiveCellToAccount(&liveCells[0], func(cellData *celltype.AccountCellData, outputIndex uint32) bool {
		if targetAccountId != celltype.AccountCharsToAccount(*cellData.Account()).AccountId() {
			return false
		}
		return true
	})
	log.Warn("parse account time spend:", time.Since(timeStart1).String())
	return obj, err
}

func (r *RpcHandler) parseLiveCellToAccount(cell *indexer.LiveCell, filter types.AccountFilterFunc) (*types.AccountReturnObj, error) {
	// get witness
	rawTx, err := r.rpcClient.GetTransaction(context.TODO(), cell.OutPoint.TxHash)
	if err != nil {
		return nil, fmt.Errorf("get raw tx err: %s", err.Error())
	}
	if len(rawTx.Transaction.Witnesses) == 0 {
		return nil, fmt.Errorf("invalid accountCell witness data,it empty, txHash: %s", rawTx.Transaction.Hash.String())
	}
	list, err := util.ParseChainAccountToJsonFormat(rawTx.Transaction, filter)
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, errors.New("account not exist")
	}
	return &list[0], nil
}
