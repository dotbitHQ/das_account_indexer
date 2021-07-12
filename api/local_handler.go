package api

import (
	"context"
	"das_account_indexer/parser/handler"
	"das_account_indexer/types"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/gotype"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
	"time"

	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/common"
	"github.com/DeAccountSystems/das_commonlib/common/dascode"
	rocksdbUtil "github.com/DeAccountSystems/das_commonlib/common/rocksdb"
	"github.com/tecbot/gorocksdb"
)

/**
 * Copyright (C), 2019-2021
 * FileName: local_handler
 * Author:   LinGuanHong
 * Date:     2021/7/11 5:50
 * Description:
 */

type RpcLocalHandler struct {
	rocksDb       *gorocksdb.DB
	systemScripts *utils.SystemScripts
}

func NewRpcLocalHandler(rpcClient rpc.Client, rocksDb *gorocksdb.DB) *RpcLocalHandler {
	systemScripts, err := utils.NewSystemScripts(rpcClient)
	if err != nil {
		panic(fmt.Errorf("init NewSystemScripts err: %s", err.Error()))
	}
	return &RpcLocalHandler{rocksDb: rocksDb, systemScripts: systemScripts}
}

func (r *RpcLocalHandler) SearchAccount(ctx context.Context, account string) common.ReqResp {
	log.Info("accept SearchAccount:", account)
	timeStart := time.Now()
	dasAccount := celltype.DasAccountFromStr(account)
	if err := dasAccount.ValidErr(); err != nil {
		return common.ReqResp{ErrNo: dascode.Err_AccountFormatInvalid, ErrMsg: err.Error()}
	}
	jsonBys, err := rocksdbUtil.RocksDbSafeGet(r.rocksDb, handler.AccountKey_AccountId(dasAccount.AccountId()))
	if err != nil {
		return common.ReqResp{ErrNo: dascode.Err_Internal, ErrMsg: fmt.Errorf("RocksDbSafeGet err: %s", err.Error()).Error()}
	} else if jsonBys == nil {
		return common.ReqResp{ErrNo: dascode.Err_AccountNotExist, ErrMsg: "account not exist, it may not be stored in the local database yet"}
	}
	returnRet := &types.AccountReturnObj{}
	if err = json.Unmarshal(jsonBys, returnRet); err != nil {
		return common.ReqResp{ErrNo: dascode.Err_Internal, ErrMsg: fmt.Errorf("unmarshal err: %s", err.Error()).Error()}
	}
	log.Info("time spend:", time.Since(timeStart).String())
	return common.ReqResp{ErrNo: dascode.DAS_SUCCESS, Data: returnRet}
}

func (r *RpcLocalHandler) GetAddressAccount(address string) common.ReqResp {
	log.Info("accept GetAddressAccount:", address)
	addrLockScriptOwnerArgs, err := gotype.Address(address).HexBys(r.systemScripts.SecpSingleSigCell.CellHash)
	if err != nil {
		return common.ReqResp{ErrNo: dascode.Err_Internal, ErrMsg: fmt.Errorf("parse address to lockArgs err: %s", err.Error()).Error()}
	}
	jsonArrBys, err := rocksdbUtil.RocksDbSafeGet(r.rocksDb, handler.AccountKey_OwnerArgHex_Bys(addrLockScriptOwnerArgs))
	if err != nil {
		return common.ReqResp{ErrNo: dascode.Err_Internal, ErrMsg: fmt.Errorf("RocksDbSafeGet err: %s", err.Error()).Error()}
	} else if jsonArrBys == nil {
		return common.ReqResp{ErrNo: dascode.Err_AccountNotExist, ErrMsg: "account not exist, it may not be stored in the local database yet"}
	}
	accountList, err := types.AccountReturnObjListFromBys(jsonArrBys)
	if err != nil {
		return common.ReqResp{ErrNo: dascode.Err_Internal, ErrMsg: fmt.Errorf("AccountReturnObjListFromBys err: %s", err.Error()).Error()}
	}
	return common.ReqResp{ErrNo: dascode.DAS_SUCCESS, Data: accountList}
}

func (r *RpcLocalHandler) Close() {

}
