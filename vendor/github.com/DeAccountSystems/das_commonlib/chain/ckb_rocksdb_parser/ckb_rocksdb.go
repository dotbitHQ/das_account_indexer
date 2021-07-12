package ckb_rocksdb_parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/db"
	"github.com/af913337456/blockparser/types"
	ckbtype "github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/tecbot/gorocksdb"
)

/**
 * Copyright (C), 2019-2020
 * FileName: ckb_database
 * Author:   LinGuanHong
 * Date:     2020/12/10 1:05
 * Description:
 */

type MsgHandler struct {
	Receive func(info *TxMsgData) error
	Close   func()
}

type CKBRocksDb struct {
	msgHandler MsgHandler
	writer     *gorocksdb.WriteOptions
	db         *gorocksdb.DB
	wb         *gorocksdb.WriteBatch
}

var (
	strToByte = func(str string) []byte {
		return []byte(str)
	}
	blockLatestKey = strToByte("block_latest")
	blockHashKey   = func(blockHash string) []byte {
		return strToByte(fmt.Sprintf("block_%s", blockHash))
	}
	blockNumberKey = func(blockNumber uint64) []byte {
		return strToByte(fmt.Sprintf("block_%d", blockNumber))
	}
)

func NewCKBRocksDb(dataDir string, handler MsgHandler) *CKBRocksDb {
	fmt.Println(fmt.Sprintf("start db [%s] init...", dataDir))
	if _db, err := db.NewDefaultRocksNormalDb(dataDir); err != nil {
		panic(fmt.Errorf("NewCKBRocksDb init db failed: (%s)", err.Error()))
	} else {
		fmt.Println("finish rocksdb init...")
		return &CKBRocksDb{
			db:         _db,
			writer:     gorocksdb.NewDefaultWriteOptions(),
			msgHandler: handler,
		}
	}
}

func (ckb *CKBRocksDb) parseSliceToStructPoint(slice *gorocksdb.Slice, receiver interface{}) error {
	if dataBytes := slice.Data(); len(slice.Data()) == 0 {
		return nil
	} else if err := json.Unmarshal(dataBytes, receiver); err != nil {
		return fmt.Errorf("parseSliceToStructPoint err: %s, data: (%s)", err.Error(), string(dataBytes))
	}
	return nil
}

func (ckb *CKBRocksDb) GetDbLastBlock() (*types.Block, error) {
	block := &types.Block{}
	if slice, err := ckb.db.Get(gorocksdb.NewDefaultReadOptions(), blockLatestKey); err != nil {
		return nil, err
	} else {
		defer slice.Free()
		err := ckb.parseSliceToStructPoint(slice, block)
		return block, err
	}
}
func (ckb *CKBRocksDb) GetDbBlockByHash(blockHash string) (*types.Block, error) {
	block := &types.Block{}
	if slice, err := ckb.db.Get(gorocksdb.NewDefaultReadOptions(), blockHashKey(blockHash)); err != nil {
		return nil, err
	} else {
		defer slice.Free()
		err := ckb.parseSliceToStructPoint(slice, block)
		return block, err
	}
}
func (ckb *CKBRocksDb) RecordBlock(block *types.Block, isUpdate, commitAfterOpt bool) error {
	ckb.wb.Put(blockLatestKey, block.ToBytes())
	ckb.wb.Put(blockHashKey(block.BlockHash), block.ToBytes())
	ckb.wb.Put(blockNumberKey(block.BlockNumber), block.ToBytes())
	if commitAfterOpt {
		return ckb.TxCommit()
	}
	return nil
}
func (ckb *CKBRocksDb) HandleForkEvent(info *types.BlockForkInfo) error {
	return ckb.db.Put(ckb.writer, blockLatestKey, info.BlockEnd.ToBytes())
}

func (ckb *CKBRocksDb) TransactionHandler(block *types.ScannerBlockInfo, dbTx interface{}, blockTxs interface{}) error {
	data := TxMsgData{
		BlockBaseInfo: *block,
		Txs:           (blockTxs).([]*ckbtype.Transaction),
	}
	txSize := len(data.Txs)
	block.TxCount = &txSize
	if ckb.msgHandler.Receive == nil {
		return nil
	}
	return ckb.msgHandler.Receive(&data)
}

func (ckb *CKBRocksDb) TxOpen() (interface{}, error) {
	if ckb.wb == nil {
		ckb.wb = gorocksdb.NewWriteBatch()
	}
	return nil, nil
}
func (ckb *CKBRocksDb) TxCommit() error {
	if ckb.wb == nil {
		return errors.New("writeBatch is nil")
	}
	return ckb.db.Write(ckb.writer, ckb.wb)
}
func (ckb *CKBRocksDb) TxRollBack() error {
	if ckb.wb != nil {
		ckb.wb.Clear()
		_ = ckb.wb.Destroy
	}
	return nil
}
func (ckb *CKBRocksDb) TxClose() {
	if ckb.wb != nil {
		ckb.wb.Clear()
		_ = ckb.wb.Destroy
	}
	ckb.wb = nil
}
func (ckb *CKBRocksDb) Close() error {
	// dont destroy tx here, because it has been call on TxClose()
	if ckb.writer != nil {
		ckb.writer.Destroy()
	}
	if ckb.wb != nil {
		ckb.TxClose()
	}
	if ckb.db != nil {
		ckb.db.Close()
	}
	if ckb.msgHandler.Close != nil {
		ckb.msgHandler.Close() // stop
	}
	return nil
}
