package handler

import (
	"context"
	"das_account_indexer/types"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/ckb/gotype"
	"github.com/DeAccountSystems/das_commonlib/common"
	"github.com/DeAccountSystems/das_commonlib/common/dascode"
	rocksdbUtil "github.com/DeAccountSystems/das_commonlib/common/rocksdb"
	"github.com/DeAccountSystems/das_commonlib/db"
	blockparserTypes "github.com/af913337456/blockparser/types"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/tecbot/gorocksdb"
	"testing"
	"time"
)

/**
 * Copyright (C), 2019-2021
 * FileName: handler_test
 * Author:   LinGuanHong
 * Date:     2021/7/12 1:06
 * Description:
 */

func Test_HandleActionTx(t *testing.T) {

	host := ""

	celltype.UseVersionReleaseSystemScriptCodeHash()

	rpcClient, err := rpc.DialWithIndexer(fmt.Sprintf("http://%s:8114", host), fmt.Sprintf("http://%s:8116", host))
	if err != nil {
		panic(fmt.Errorf("init rpcClient err: %s", err.Error()))
	}
	txStatus, err := rpcClient.GetTransaction(context.TODO(), ckbTypes.HexToHash(""))
	if err != nil {
		panic(fmt.Errorf("GetTransaction err: %s", err.Error()))
	}
	infoDb, err := db.NewDefaultRocksNormalDb("../../rocksdb-data/account2")
	if err != nil {
		panic(fmt.Errorf("NewDefaultRocksNormalDb err: %s", err.Error()))
	}
	p := &DASActionHandleFuncParam{
		Base: &types.ParserHandleBaseTxInfo{
			ScanInfo: blockparserTypes.ScannerBlockInfo{},
			Tx:       *txStatus.Transaction,
			TxIndex:  0,
		},
		RpcClient: rpcClient,
		Rocksdb:   infoDb,
	}
	resp := HandleEditRecordsTx("", p)
	log.Warn("resp.err:-->", resp.err)
	ret1 := searchAccount(infoDb, "d55213.bit")
	fmt.Println(ret1.ErrNo)
	fmt.Println(ret1.Data)

	ret := getAddressAccount("", infoDb)
	fmt.Println(ret.ErrMsg)
	fmt.Println(ret.Data)
}

func searchAccount(rocksdb *gorocksdb.DB, account string) common.ReqResp {
	log.Info("accept SearchAccount:", account)
	timeStart := time.Now()
	dasAccount := celltype.DasAccountFromStr(account)
	if err := dasAccount.ValidErr(); err != nil {
		return common.ReqResp{ErrNo: dascode.Err_AccountFormatInvalid, ErrMsg: err.Error()}
	}
	jsonBys, err := rocksdbUtil.RocksDbSafeGet(rocksdb, AccountKey_AccountId(dasAccount.AccountId()))
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

func getAddressAccount(address string, rocksdb *gorocksdb.DB) common.ReqResp {
	log.Info("accept GetAddressAccount:", address)
	addrLockScriptOwnerArgs, err := gotype.Address(address).HexBys(ckbTypes.HexToHash(""))
	if err != nil {
		return common.ReqResp{ErrNo: dascode.Err_Internal, ErrMsg: fmt.Errorf("parse address to lockArgs err: %s", err.Error()).Error()}
	}
	jsonArrBys, err := rocksdbUtil.RocksDbSafeGet(rocksdb, AccountKey_OwnerArgHex_Bys(addrLockScriptOwnerArgs))
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
