package handler

import (
	"context"
	"das_account_indexer/types"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/db"
	blockparserTypes "github.com/af913337456/blockparser/types"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
	"testing"
)

/**
 * Copyright (C), 2019-2021
 * FileName: handler_test
 * Author:   LinGuanHong
 * Date:     2021/7/12 1:06
 * Description:
 */

func Test_HandleConfirmProposalTx(t *testing.T) {

	host := ""

	celltype.UseVersion3SystemScriptCodeHash()

	rpcClient, err := rpc.DialWithIndexer(fmt.Sprintf("http://%s:8114", host), fmt.Sprintf("http://%s:8116", host))
	if err != nil {
		panic(fmt.Errorf("init rpcClient err: %s", err.Error()))
	}
	txStatus, err := rpcClient.GetTransaction(context.TODO(), ckbTypes.HexToHash("0x2df31e65a97685107323d9efe3186365c0fce5c1734b06871257d910902a4b9a"))
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
	_ = HandleConfirmProposalTx("", p)
}
