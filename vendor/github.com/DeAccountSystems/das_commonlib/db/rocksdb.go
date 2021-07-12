package db

import (
	"fmt"
	"github.com/tecbot/gorocksdb"
)

/**
 * Copyright (C), 2019-2020
 * FileName: rocksdb
 * Author:   LinGuanHong
 * Date:     2020/12/25 10:53
 * Description:
 */

func defaultRateLimiter() *gorocksdb.RateLimiter {
	return gorocksdb.NewRateLimiter(1024, 100*1000, 10)
}

func NewDefaultRocksTxDb(dataDir string) (*gorocksdb.TransactionDB, error) {
	rateLimiter := defaultRateLimiter()
	txOpts := gorocksdb.NewDefaultOptions()
	txOpts.SetRateLimiter(rateLimiter)
	txOpts.SetCreateIfMissing(true)
	txDb, err := gorocksdb.OpenTransactionDb(txOpts, gorocksdb.NewDefaultTransactionDBOptions(), dataDir)
	if err != nil {
		return nil, fmt.Errorf("NewDefaultRocksTxDb failed: (%s)", err.Error())
	}
	return txDb, nil
}

func NewDefaultRocksNormalDb(dataDir string) (*gorocksdb.DB, error) {
	opts := gorocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, dataDir)
	if err != nil {
		return nil, fmt.Errorf("NewDefaultRocksNormalDb failed: (%s)", err.Error())
	}
	return db, nil
}

func NewDefaultReadOnlyRocksDb(dataDir string) (*gorocksdb.DB, error) {
	rateLimiter := defaultRateLimiter()
	opts := gorocksdb.NewDefaultOptions()
	opts.SetRateLimiter(rateLimiter)
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDbForReadOnly(opts, dataDir, false)
	if err != nil {
		return nil, fmt.Errorf("NewDefaultReadOnlyRocksDb failed: (%s)", err.Error())
	}
	return db, nil
}
