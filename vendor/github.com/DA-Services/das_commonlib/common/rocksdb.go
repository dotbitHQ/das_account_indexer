package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/tecbot/gorocksdb"
)

/**
 * Copyright (C), 2019-2021
 * FileName: rocksdb_util
 * Author:   LinGuanHong
 * Date:     2021/1/8 2:30 下午
 * Description:
 */

type IHashListBatchWrite interface {
	Obj(i int) interface{}
	Size() int
	Hash(i int) types.Hash
}

func RocksDbSafeGet(rocksdb *gorocksdb.DB, key []byte) ([]byte, error) {
	return RocksDbSafeGetWithReadOpt(rocksdb, gorocksdb.NewDefaultReadOptions(), key)
}

func RocksDbSafeGetWithReadOpt(rocksdb *gorocksdb.DB, rOpts *gorocksdb.ReadOptions, key []byte) ([]byte, error) {
	result, err := rocksdb.Get(rOpts, key)
	if err != nil {
		if result != nil {
			result.Free()
		}
		return nil, err
	}
	if size := len(result.Data()); size == 0 {
		return nil, nil // empty
	} else {
		tmpData := make([]byte, len(result.Data()))
		copy(tmpData, result.Data())
		defer result.Free()
		return tmpData, nil
	}
}

func ExistOrStore(rocksdb *gorocksdb.DB, wb *gorocksdb.WriteBatch, key []byte) (bool, error) {
	data, err := RocksDbSafeGet(rocksdb, key)
	if err != nil {
		return false, fmt.Errorf("read failed: %s", err.Error())
	}
	if len(data) > 0 {
		return true, nil
	}
	wb.Put(key, []byte{1})
	return false, nil
}

func RocksDbBatchWriteWithWB(wb *gorocksdb.WriteBatch, objList IHashListBatchWrite, key func(hash types.Hash) []byte) error {
	size := objList.Size()
	for index := 0; index < size; index++ {
		if bys, err := rlp.EncodeToBytes(objList.Obj(index)); err != nil {
			return fmt.Errorf("rocksDbBatchWriteWithWB rlp.EncodeToBytes failed: %s", err.Error())
		} else {
			wb.Put(key(objList.Hash(index)), bys)
		}
	}
	return nil
}

func RocksDbBatchDelWithWB(wb *gorocksdb.WriteBatch, keyList [][]byte) {
	size := len(keyList)
	for index := 0; index < size; index++ {
		wb.Delete(keyList[index])
	}
}

func RocksDbIteratorLoad(rocksdb *gorocksdb.DB, keyPrefix []byte, handleBytes func(keyBytes, dataBytes []byte) error) error {
	opts := gorocksdb.NewDefaultReadOptions()
	opts.SetFillCache(false)
	defer opts.Destroy()
	reader := rocksdb.NewIterator(opts)
	defer reader.Close()
	for reader.Seek(keyPrefix); ; reader.Next() {
		if valid := reader.ValidForPrefix(keyPrefix); !valid {
			// lsm tree，key 是有顺序的，如果出现一个不是同前缀，那么就直接退出，结束遍历
			break
		}
		if _, err := SafeHandleReaderKV(reader.Key(), reader.Value(), handleBytes); err != nil {
			return err
		}
	}
	return nil
}

func SafeHandleReaderKV(key, value *gorocksdb.Slice, handleBytes func(keyBytes, dataBytes []byte) error) (bool, error) {
	defer func() {
		if key != nil {
			key.Free()
		}
		if value != nil {
			value.Free()
		}
	}()
	if bys := value.Data(); len(bys) > 0 {
		tmpKeyData := make([]byte, len(key.Data()))
		copy(tmpKeyData, key.Data())
		tmpValueData := make([]byte, len(bys))
		copy(tmpValueData, bys)
		if err := handleBytes(tmpKeyData, tmpValueData); err != nil {
			return false, fmt.Errorf("handleBytes err: %s", err.Error())
		} else {
			return false, nil
		}
	} else {
		return true, nil
	}
}

func RocksDbLoadOneByPrefix(rocksdb *gorocksdb.DB, keyPrefix []byte, handleBytes func(keyBytes,dataBytes []byte) (interface{}, error)) (interface{}, error) {
	return loadOneByPrefix(rocksdb,keyPrefix,false,handleBytes)
}

func loadOneByPrefix(rocksdb *gorocksdb.DB, keyPrefix []byte, lastOne bool, handleBytes func(keyBytes,dataBytes []byte) (interface{}, error)) (interface{}, error) {
	opts := gorocksdb.NewDefaultReadOptions()
	opts.SetFillCache(false)
	reader := rocksdb.NewIterator(opts)
	reader.Seek(keyPrefix)

	free := func() {
		opts.Destroy()
		reader.Close()
	}
	if valid := reader.ValidForPrefix(keyPrefix); !valid {
		free()
		return nil, nil
	}
	_key_ := reader.Key()
	value := reader.Value()
	defer func() {
		value.Free()
		free()
	}()
	if bys := value.Data(); len(bys) > 0 {
		keyRawDataBytes := _key_.Data()
		keyTmpData := make([]byte, len(keyRawDataBytes))
		copy(keyTmpData, keyRawDataBytes)
		tmpData := make([]byte, len(bys))
		copy(tmpData, bys)
		if val, err := handleBytes(keyTmpData,tmpData); err != nil {
			return nil, fmt.Errorf("RocksDbLoadOneByPrefix handleBytes failed: %s", err.Error())
		} else {
			return val, nil
		}
	}
	return nil, fmt.Errorf("RocksDbLoadOneByPrefix key exist, but data is empty")
}

func RocksDbNormalStore(rocksDb *gorocksdb.DB, key []byte, bytes func() ([]byte, error)) error {
	if bys, err := bytes(); err != nil {
		return fmt.Errorf("RocksDbNormalStore Bytes err: %s", err.Error())
	} else if err := rocksDb.Put(gorocksdb.NewDefaultWriteOptions(), key, bys); err != nil {
		return fmt.Errorf("RocksDbNormalStore put err: %s", err.Error())
	}
	return nil
}

func RocksDbNormalLoad(rocksDb *gorocksdb.DB, key []byte, handle func(bytes []byte) error) error {
	if valueData, err := RocksDbSafeGet(rocksDb, key); err != nil {
		return fmt.Errorf("RocksDbNormalLoad get key hash value err: %s", err.Error())
	} else if len(valueData) == 0 {
		return nil
	} else {
		if err := handle(valueData); err != nil {
			return fmt.Errorf("RocksDbNormalLoad lp.decode failed: %s", err.Error())
		} else {
			return nil
		}
	}
}
