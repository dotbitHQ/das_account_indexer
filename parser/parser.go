package parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"das_account_indexer/parser/handler"
	accountIndexerTypes "das_account_indexer/types"

	blockparser "github.com/DeAccountSystems/das_commonlib/chain/ckb_rocksdb_parser"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/eager7/elog"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/tecbot/gorocksdb"
)

/**
 * Copyright (C), 2019-2021
 * FileName: parser
 * Author:   LinGuanHong
 * Date:     2021/7/9 4:12
 * Description:
 */

var (
	log = elog.NewLogger("tx_parser", elog.NoticeLevel)
)

type TxParser struct {
	rpcClient          rpc.Client
	rocksdb            *gorocksdb.DB
	context            context.Context
	actionRegister     *handler.ActionRegister
	currentBlockNumber uint64
	latestBlockNumber  uint64
	targetBlockHeight  uint64
}

type InitTxParserParam struct {
	RpcClient         rpc.Client
	Rocksdb           *gorocksdb.DB
	Context           context.Context
	TargetBlockHeight uint64
	FontBlockNumber   uint64
}

func NewParserRpcTx(p *InitTxParserParam) *TxParser {
	txParser := &TxParser{
		rpcClient:         p.RpcClient,
		rocksdb:           p.Rocksdb,
		context:           p.Context,
		targetBlockHeight: p.TargetBlockHeight,
		actionRegister:    handler.NewActionRegister(),
	}
	go txParser.getChainLatestBlockNumber(p.FontBlockNumber)
	return txParser
}

func (p *TxParser) GetCurrentBlockNumber() *uint64 {
	return &p.currentBlockNumber
}

func (p *TxParser) BlockSyncFinish() bool {
	finishSync := p.latestBlockNumber > 0 && p.currentBlockNumber >= p.latestBlockNumber
	log.Warn(fmt.Sprintf("sync blockNumber info, latest: %d, current: %d", p.latestBlockNumber, p.currentBlockNumber))
	return finishSync
}

func (p *TxParser) getChainLatestBlockNumber(blockFontNumber uint64) {
	for {
		if p.context != nil && p.context.Err() != nil {
			return
		}
		if blockNumber, err := p.rpcClient.GetTipBlockNumber(context.TODO()); err != nil {
			log.Error(fmt.Sprintf("getChainLatestBlockNumber err: %s", err.Error()))
		} else {
			p.latestBlockNumber = blockNumber - blockFontNumber
		}
		time.Sleep(time.Second)
	}
}

func (p *TxParser) Handler(data []byte, delayMs *int64) error {
	var msgData = blockparser.TxMsgData{}
	if err := json.Unmarshal(data, &msgData); err != nil {
		return fmt.Errorf("TxParser handle tx failed, cant unmarshal: %s", err.Error())
	}
	return p.Handle1(msgData, delayMs)
}
func (p *TxParser) Handle1(msgData blockparser.TxMsgData, delayMs *int64) error {
	txSize := len(msgData.Txs)
	nowHeight := msgData.BlockBaseInfo.BlockNumber
	log.Info("------------------ start handle one tx ------------------")
	log.Info(fmt.Sprintf("-----block: height: %d hash: %s, tx count: %d-----", nowHeight, msgData.BlockBaseInfo.BlockHash, txSize))
	defer func() {
		if *delayMs > 0 {
			time.Sleep(time.Millisecond * time.Duration(*delayMs))
		}
	}()
	if targetHeight := p.targetBlockHeight; targetHeight != 0 {
		if targetHeight > nowHeight {
			log.Info(fmt.Sprintf("skip, target: %d, now: %d", targetHeight, nowHeight))
			return nil
		}
	}
	for txIndex := 0; txIndex < txSize; txIndex++ {
		tx := msgData.Txs[txIndex]
		if cellSize := len(tx.Outputs); cellSize == 0 {
			log.Info("tx output is zero")
			continue
		}
		// get action
		log.Info(fmt.Sprintf("parse txHash: %s", tx.Hash.String()))
		if actionName, err := celltype.GetActionNameFromWitnessData(tx); err != nil {
			log.Warn("skip this tx:", err.Error())
		} else {
			log.Info("tx action name:", actionName)
			if handleFunc := p.actionRegister.GetTxActionHandleFunc(actionName); handleFunc != nil {
				handleRet := handleFunc(actionName, &handler.DASActionHandleFuncParam{
					Base: &accountIndexerTypes.ParserHandleBaseTxInfo{
						Tx:       *tx,
						TxIndex:  uint8(txIndex),
						ScanInfo: msgData.BlockBaseInfo,
					},
					RpcClient: p.rpcClient,
					Rocksdb:   p.rocksdb,
				})
				if err := handleRet.Error(); err != nil {
					log.Error("handle tx err:", err.Error())
				}
				if handleRet.Data != nil {
					bys, _ := json.Marshal(handleRet)
					log.Info("action handle ret:", string(bys))
				}
				if handleRet.Rollback {
					return errors.New("err happen, dont commit this block")
				}
			} else {
				log.Error(fmt.Sprintf("action handleFunc not found!"))
			}
		}
	}
	p.currentBlockNumber = msgData.BlockBaseInfo.BlockNumber
	return nil
}

func (p *TxParser) Close() {
	if p.rocksdb != nil {
		p.rocksdb.Close()
	}
	if p.rpcClient != nil {
		p.rpcClient.Close()
	}
}
