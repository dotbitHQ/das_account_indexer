package scanner

import (
	"github.com/af913337456/blockparser/types"
	"math/big"
)

/**
  author : LinGuanHong
  github : https://github.com/af913337456
  blog   : http://www.cnblogs.com/linguanh
  time   : 10:38
*/

// 作为兼容多条链的区块扫描，独立出公共接口

type IDatabaseEngine interface {
	GetDbLastBlock() (*types.Block, error)
	GetDbBlockByHash(blockHash string) (*types.Block, error)
	RecordBlock(block *types.Block, updateModel, commitAfterOpt bool) error
	HandleForkEvent(info *types.BlockForkInfo) error
	TransactionHandler(block *types.ScannerBlockInfo, dbTx interface{}, blockTxs interface{}) error
	TxOpen() (interface{}, error)
	TxCommit() error
	TxRollBack() error
	TxClose()
	Close() error
}

type IBlockChain interface {
	GetParentHash(childHash string) (string, error)
	GetLatestBlockNumber() (*big.Int, error)
	GetBlockInfoByNumber(blockNumber *big.Int) (*types.ScannerBlockInfo, error)
	Close()
}

type ILog interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}
