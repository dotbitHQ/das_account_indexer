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

func buildP(txHash string) (*DASActionHandleFuncParam, *gorocksdb.DB) {
	host := ""

	rpcClient, err := rpc.DialWithIndexer(fmt.Sprintf("http://%s:8114", host), fmt.Sprintf("http://%s:8116", host))
	if err != nil {
		panic(fmt.Errorf("init rpcClient err: %s", err.Error()))
	}
	txStatus, err := rpcClient.GetTransaction(
		context.TODO(),
		ckbTypes.HexToHash(txHash))
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
	return p, infoDb
}

func Test_HandleAccountCellType(t *testing.T) {
	// celltype.UseVersion3SystemScriptCodeHash()
	// celltype.DasAccountCellScript.Out.CodeHash = ckbTypes.HexToHash("0x334540e23ec513f691cdd9490818237cbc9675861e4f19c480e0c520c715fd34")
	celltype.UseVersionReleaseSystemScriptCodeHash()
	p, _ := buildP("0xa0d957793b3f7c0fe7335e9c8aa309c7344d3cafd0bcee257825b1ce1cdda323")
	resp := HandleAccountCellType(celltype.Action_ConfirmProposal, p)
	log.Warn("resp.err:-->", resp.err)
}

func Test_HandleConfirmProposeTx(t *testing.T) {
	celltype.UseVersionReleaseSystemScriptCodeHash()
	fmt.Println(celltype.DasProposeCellScript.Out.CodeHash.String())
	p, infoDb := buildP("0xbc01e3ae6ce550c8d24dc6f7a819331cd525875340844a090b9fb4c50a12d152")
	resp := HandleConfirmProposalTx("", p)
	log.Warn("resp.err:-->", resp.err)
	// ret1 := searchAccount(infoDb, "d55213.bit")
	// fmt.Println(ret1.ErrNo)
	// fmt.Println(ret1.Data)
	//
	ret := getAddressAccount("0x910a7b702388fe5a4a48933327b6b908b674969f", infoDb)
	fmt.Println(ret.ErrMsg)
	bys, _ := json.Marshal(ret.Data)
	fmt.Println(string(bys))
}
func Test_HandleActionTx(t *testing.T) {
	celltype.UseVersion3SystemScriptCodeHash()
	p, infoDb := buildP("")
	resp := HandleTransferAccountTx("", p)
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
	accountList, err := types.AccountReturnObjListFromBys(&jsonArrBys)
	if err != nil {
		return common.ReqResp{ErrNo: dascode.Err_Internal, ErrMsg: fmt.Errorf("AccountReturnObjListFromBys err: %s", err.Error()).Error()}
	}
	return common.ReqResp{ErrNo: dascode.DAS_SUCCESS, Data: accountList}
}
