package ckb_rocksdb_parser

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/DeAccountSystems/das_commonlib/common"
	"github.com/af913337456/blockparser/types"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	ckbtype "github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2020
 * FileName: ckb_normal_blockchain
 * Author:   LinGuanHong
 * Date:     2020/12/9 5:31
 * Description:
 */

type CKBBlockChain struct {
	Ctx       context.Context
	ReqTryCfg *types.RetryConfig
	RpcClient rpc.Client
}

func NewCKBBlockChainWithRpcClient(ctx context.Context, rpcClient rpc.Client, reqTryCfg *types.RetryConfig) *CKBBlockChain {
	if reqTryCfg == nil {
		reqTryCfg = &types.RetryConfig{
			RetryTime: 3,
			DelayTime: time.Second * 5,
		}
	}
	return &CKBBlockChain{
		Ctx:       ctx,
		ReqTryCfg: reqTryCfg,
		RpcClient: rpcClient,
	}
}

func NewCKBBlockChain(ctx context.Context, nodeUrl string, reqTryCfg *types.RetryConfig) *CKBBlockChain {
	if client, err := rpc.Dial(nodeUrl); err != nil {
		panic(fmt.Errorf("NewCKBBlockChain init rpcClient failed: (%s), nodeUrl: %s", err.Error(), nodeUrl))
	} else {
		fmt.Println("finish rpc init...")
		return NewCKBBlockChainWithRpcClient(ctx, client, reqTryCfg)
	}
}

func (ckb *CKBBlockChain) GetParentHash(childHash string) (string, error) {
	ret, err := common.RetryReq(ckb.ReqTryCfg.RetryTime, ckb.ReqTryCfg.DelayTime, func() (interface{}, error) {
		if block, err := ckb.RpcClient.GetBlock(ckb.Ctx, ckbtype.HexToHash(childHash)); err != nil {
			return "", err
		} else {
			return block.Header.ParentHash.String(), nil
		}
	})
	return ret.(string), err
}
func (ckb *CKBBlockChain) GetLatestBlockNumber() (*big.Int, error) {
	ret, err := common.RetryReq(ckb.ReqTryCfg.RetryTime, ckb.ReqTryCfg.DelayTime, func() (interface{}, error) {
		if number, err := ckb.RpcClient.GetTipBlockNumber(ckb.Ctx); err != nil {
			return &big.Int{}, err
		} else {
			return new(big.Int).SetUint64(number), nil
		}
	})
	return ret.(*big.Int), err
}
func (ckb *CKBBlockChain) GetBlockInfoByNumber(blockNumber *big.Int) (*types.ScannerBlockInfo, error) {
	ret, err := common.RetryReq(ckb.ReqTryCfg.RetryTime, ckb.ReqTryCfg.DelayTime, func() (interface{}, error) {
		if block, err := ckb.RpcClient.GetBlockByNumber(ckb.Ctx, blockNumber.Uint64()); err != nil {
			return &types.ScannerBlockInfo{}, err
		} else {
			txCount := len(block.Transactions)
			return &types.ScannerBlockInfo{
				BlockHash:   block.Header.Hash.String(),
				ParentHash:  block.Header.ParentHash.String(),
				Version:     int32(block.Header.Version),
				BlockNumber: block.Header.Number,
				Timestamp:   strconv.FormatInt(int64(block.Header.Timestamp), 10),
				Txs:         block.Transactions,
				TxCount:     &txCount,
			}, nil
		}
	})
	return ret.(*types.ScannerBlockInfo), err
}

func (ckb *CKBBlockChain) Close() {
	if ckb.RpcClient != nil {
		ckb.RpcClient.Close()
	}
}
