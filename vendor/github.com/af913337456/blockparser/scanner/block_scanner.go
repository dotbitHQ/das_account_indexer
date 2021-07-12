package scanner

import (
	"errors"
	"fmt"
	"github.com/af913337456/blockparser/types"
	"math/big"
	"strings"
	"sync"
	"time"
)

/**
  author : LinGuanHong
  github : https://github.com/af913337456
  blog   : http://www.cnblogs.com/linguanh
  time   : 14:49
*/

// DES: 区块扫描者。遍历出区块的交易，方便从交易中解析出数据，做自定义操作
// 区块遍历者

type BlockScanner struct {
	chain         IBlockChain     // 接口实现者
	dbEngine      IDatabaseEngine // 数据库
	logger        ILog
	lastBlock     *types.Block // 用来存储每次遍历后，上一次的区块
	lastNumber    *big.Int     // 上一次区块的区块号
	fork          bool         // 区块分叉标记位
	stop          chan bool    // 用来控制是否停止遍历的管道
	exit          chan bool    // 彻底退出
	DelayControl  types.DelayControl
	blockCount    int64
	FrontNumber   int // 扫描的区块前置数，target = current - FrontNumber
	onceDo        sync.Once
	waitToChasing bool
}

type InitBlockScanner struct {
	Chain       IBlockChain
	Db          IDatabaseEngine
	Log         ILog
	Control     types.DelayControl
	FrontNumber int
}

func NewBlockScanner(param InitBlockScanner) *BlockScanner {
	return &BlockScanner{
		chain:        param.Chain,
		dbEngine:     param.Db,
		logger:       param.Log,
		lastBlock:    &types.Block{},
		fork:         false,
		stop:         make(chan bool),
		exit:         make(chan bool),
		onceDo:       sync.Once{},
		DelayControl: param.Control,
		FrontNumber:  param.FrontNumber,
	}
}

// 初始化：内部在开始遍历时赋值 lastBlock
func (scanner *BlockScanner) init() error {
	if scanner.waitToChasing {
		return nil
	}
	setNextBlockNumber := func() {
		// 区块 hash 不为空，证明不是首次启动了。是后续的启动
		scanner.lastNumber = new(big.Int).SetUint64(scanner.lastBlock.BlockNumber)
		// 下面加 1，因为上一次数据库存的是已经遍历完了的，接下来的是它的下一个
		scanner.lastNumber.Add(scanner.lastNumber, new(big.Int).SetInt64(1))
	}
	// 从数据库中寻找出上一次成功遍历的且不是分叉的区块
	dbBlock, err := scanner.dbEngine.GetDbLastBlock()
	if err != nil {
		return err
	}
	scanner.info("last time scan latest block info:", string(dbBlock.ToBytes()))
	if scanner.lastBlock.BlockHash != "" {
		// 被设置了的情况，与 db 的比较，找出最新的
		if dbBlock.BlockNumber > scanner.lastBlock.BlockNumber {
			scanner.lastBlock = dbBlock
		}
		setNextBlockNumber()
		return nil
	}
	scanner.lastBlock = dbBlock
	if scanner.lastBlock.BlockHash == "" {
		// 区块 hash 为空，证明是整个程序的首次启动，那么从节点中获取最新生成的区块
		// GetLatestBlockNumber 获取最新区块的区块号
		latestBlockNumber, err := scanner.chain.GetLatestBlockNumber()
		if err != nil {
			return err
		}
		// GetBlockInfoByNumber 根据区块号获取区块数据
		latestBlock, err := scanner.chain.GetBlockInfoByNumber(latestBlockNumber)
		if err != nil {
			return err
		}
		if latestBlock.BlockNumber == 0 {
			panic(latestBlockNumber.String())
		}
		// 下面是赋值区块遍历者的 lastBlock 变量
		scanner.lastBlock.BlockHash = latestBlock.BlockHash
		scanner.lastBlock.ParentHash = latestBlock.ParentHash
		scanner.lastBlock.BlockNumber = latestBlock.BlockNumber
		scanner.lastBlock.CreateTime = scanner.hexToTen(latestBlock.Timestamp).Int64()
		scanner.lastNumber = latestBlockNumber
	} else {
		setNextBlockNumber()
	}
	return nil
}

// 设置固定的开始高度，起效要求要在开始 scan 之前
func (scanner *BlockScanner) SetStartScannerHeight(height int64) error {
	if height <= 0 {
		return nil
	}
	fmt.Println("start retryGetBlockInfoByNumber...")
	blockInfo, err := scanner.retryGetBlockInfoByNumber(new(big.Int).SetInt64(height))
	if err != nil {
		return err
	}
	if blockInfo.IsEmpty() {
		// 指定的区块高度，还没到，那么轮训等待
		scanner.waitToChasing = true
		scanner.lastNumber = new(big.Int).SetInt64(height)
		return nil
	}
	scanner.lastBlock.BlockHash = blockInfo.BlockHash
	scanner.lastBlock.ParentHash = blockInfo.ParentHash
	scanner.lastBlock.BlockNumber = blockInfo.BlockNumber
	scanner.lastBlock.CreateTime = scanner.hexToTen(blockInfo.Timestamp).Int64()
	return nil
}

// 整个区块扫描的启动函数
func (scanner *BlockScanner) Start() error {
	if err := scanner.init(); err != nil {
		return err
	}
	scan := func() { // scan 函数，就是区块扫描函数
		if err := scanner.scan(); nil != err {
			scanner.error("scanner err :", err.Error())
			return
		}
		time.Sleep(scanner.DelayControl.RoundDelay) // 延迟一秒开始下一轮
	}
	go func() {
		scanner.info("starting scanner...")
		for {
			select {
			case <-scanner.stop:
				scanner.info("going to stop scanner ...")
				scanner.info("first stop scan client, stop data provider...")
				if scanner.chain != nil {
					scanner.chain.Close()
				}
				scanner.info("then stop database engine, stop data store...")
				if scanner.dbEngine != nil {
					_ = scanner.dbEngine.Close()
				}
				scanner.info("exit scanner!")
				scanner.exit <- true
				return
			default:
				if !scanner.fork {
					scan()
				} else {
					if err := scanner.init(); err != nil {
						scanner.error(err.Error())
						return
					}
					scanner.fork = false
				}
			}
		}
	}()
	return nil
}

func (scanner *BlockScanner) Wait() {
	<-scanner.exit
}

// 公有函数，可以共外部调用，来控制停止区块遍历
func (scanner *BlockScanner) Stop() {
	scanner.stop <- true
}

func (scanner *BlockScanner) info(args ...interface{}) {
	scanner.logger.Info(args...)
}
func (scanner *BlockScanner) error(args ...interface{}) {
	scanner.logger.Error(args...) // need panic, stop the world
}
func (scanner *BlockScanner) warn(args ...interface{}) {
	scanner.logger.Warn(args...)
}

// 是否分叉，返回 true 是分叉
func (scanner *BlockScanner) isFork(currentBlock *types.Block) bool {
	if currentBlock.BlockNumber == 0 {
		panic("invalid block")
	}
	if scanner.lastBlock.IsEmpty() || scanner.lastBlock.BlockHash == currentBlock.BlockHash || scanner.lastBlock.BlockHash == currentBlock.ParentHash {
		scanner.lastBlock = currentBlock // 没有发生分叉，则更新上一区块为当前被检测的
		return false
	}
	return true
}

func (scanner *BlockScanner) forkCheck(currentBlock *types.Block) *types.BlockForkInfo {
	if !scanner.isFork(currentBlock) {
		return nil
	}
	// 获取出最初开始分叉的那个区块
	forkBlock, err := scanner.getStartForkBlock(currentBlock.ParentHash)
	if err != nil {
		panic(err)
	}
	scanner.lastBlock = forkBlock // 更新。从这个区块开始分叉的
	return &types.BlockForkInfo{
		currentBlock,
		forkBlock,
	}
}

// 获取分叉点区块
func (scanner *BlockScanner) getStartForkBlock(parentHash string) (*types.Block, error) {
	// 获取当前区块的父区块，分叉从父区块开始
	scanner.info(fmt.Sprintf("fork event ===> search one pre block, hash: %s", parentHash))
	if parent, err := scanner.dbEngine.GetDbBlockByHash(parentHash); err == nil && !parent.IsEmpty() && parentHash == parent.BlockHash {
		return parent, nil // 本地存在，直接返回分叉点区块
	}
	// 数据库没有父区块记录，准备从链上接口获取
	fatherHash, err := scanner.retryGetBlockInfoByHash(parentHash)
	if err != nil {
		return nil, fmt.Errorf("分叉严重错误，父区块亦无法从链上获取，需要重启区块扫描 %s", err.Error())
	}
	// 继续递归往上查询，直到在数据库中有它的记录
	return scanner.getStartForkBlock(fatherHash)
}

func (scanner *BlockScanner) hexToTen(hex string) *big.Int {
	if !strings.HasPrefix(hex, "0x") {
		ten, _ := new(big.Int).SetString(hex, 10) // 本身就是 10 进制字符串，直接设置
		return ten
	}
	ten, _ := new(big.Int).SetString(hex[2:], 16)
	return ten
}

// 区块号存在，信息获取为空，可能是链上网络延时问题，重试策略函数
func (scanner *BlockScanner) retryGetBlockInfoByNumber(targetNumber *big.Int) (*types.ScannerBlockInfo, error) {
	// 下面调用我们请求者 client 的 GetBlockInfoByNumber 函数
	fullBlock, err := scanner.chain.GetBlockInfoByNumber(targetNumber)
	if err != nil {
		return nil, err
	}
	return fullBlock, nil
}

// 区块 hash 存在，信息获取为空，可能是链上网络或节点问题，重试策略函数
func (scanner *BlockScanner) retryGetBlockInfoByHash(hash string) (string, error) {
	// 下面调用我们请求者 client 的 GetBlockInfoByHash 函数
	parentHash, err := scanner.chain.GetParentHash(hash)
	if err != nil {
		return "", err
	}
	return parentHash, nil
}

// 获取要扫描的区块号
func (scanner *BlockScanner) getScannerBlockNumber() (*big.Int, error) {
	// 调用请求者 client 获取公链上最新生成的区块的区块号
	latestNumber, err := scanner.chain.GetLatestBlockNumber()
	if err != nil {
		return nil, fmt.Errorf("GetLatestBlockNumber: %s", err.Error())
	}
	fix := func(number *big.Int) *big.Int {
		return new(big.Int).Sub(number, new(big.Int).SetInt64(int64(scanner.FrontNumber)))
	}
	// 下面使用 new 的形式初始化并设置值，不要直接赋值，
	// 否则会和 lastNumber 的内存地址一样，影响后面的获取区块信息
	targetNumber := new(big.Int).Set(scanner.lastNumber)
	// 比较区块号大小
	if fix(latestNumber).Cmp(targetNumber) < 0 {
		// 最新的区块高度比设置的要小，则等待新区块高度 >= 设置的
	Next:
		for {
			select {
			case <-time.After(scanner.DelayControl.CatchDelay): // 延时x秒重新获取
				if newBlockNumber, err := scanner.chain.GetLatestBlockNumber(); err == nil && fix(newBlockNumber).Cmp(targetNumber) >= 0 {
					scanner.waitToChasing = false
					break Next // 跳出循环
				} else {
					scanner.info(fmt.Sprintf("wait for chasing block, current: %s, target: %s, fix number: %d", newBlockNumber.String(), targetNumber.String(), scanner.FrontNumber))
				}
			}
		}
	}
	return targetNumber, nil // 返回目标区块高度
}

// 扫描区块
func (scanner *BlockScanner) scan() error {
	// 获取要进行扫描的区块号
	targetNumber, err := scanner.getScannerBlockNumber()
	if err != nil {
		return fmt.Errorf("getScannerBlockNumber: %s", err.Error())
	}
	// 使用具有重试策略的函数获取区块信息
	info, err := scanner.retryGetBlockInfoByNumber(targetNumber)
	if err != nil {
		return fmt.Errorf("retryGetBlockInfoByNumber: %s", err.Error())
	}
	// 区块号自增 1，在下次扫描的时候，指向下一个高度的区块
	scanner.lastNumber.Add(scanner.lastNumber, new(big.Int).SetInt64(1))
	// 因为涉及到多表的更新，我们需要采用数据库事务处理
	tx, err := scanner.dbEngine.TxOpen() // 开启事务
	defer scanner.dbEngine.TxClose()
	if err != nil {
		return fmt.Errorf("db open session failed: %s", err.Error())
	}
	// 准备保存区块信息，先判断当前区块记录是否已经存在
	block := &types.Block{}
	if block, err = scanner.dbEngine.GetDbBlockByHash(info.BlockHash); err == nil {
		updateModel := !block.IsEmpty()
		block.BlockNumber = info.BlockNumber
		block.ParentHash = info.ParentHash
		block.Version = info.Version
		if timeStamp := scanner.hexToTen(info.Timestamp); timeStamp != nil {
			block.CreateTime = scanner.hexToTen(info.Timestamp).Int64()
		} else {
			return fmt.Errorf("parse timestamp err, invalid timestamp: [%s]", info.Timestamp)
		}
		block.BlockHash = info.BlockHash
		block.Fork = false
		if err = scanner.dbEngine.RecordBlock(block, updateModel, false); err != nil {
			return fmt.Errorf("tx.Insert: %s", err.Error())
		}
	}
	// 检查区块是否分叉
	if forkInfo := scanner.forkCheck(block); forkInfo != nil {
		scanner.info("fork！", string(forkInfo.ToBytes()))
		if err = scanner.dbEngine.TxRollBack(); err != nil {
			return fmt.Errorf("fork happen, commit failed: %s", err.Error())
		}
		if err := scanner.dbEngine.HandleForkEvent(forkInfo); err != nil {
			return fmt.Errorf("handle fork block failed %s", err.Error())
		}
		scanner.fork = true
		return errors.New("fork happen") // 返回错误，让上层处理并重启区块扫描
	}
	// 解析区块内数据，读取内部的 “transactions” 交易信息
	blockNumber := info.BlockNumber
	zero := 0
	info.TxCount = &zero
	scanner.info("scan block start ==>", "number:", blockNumber, "hash:", info.BlockHash)
	if err := scanner.dbEngine.TransactionHandler(info, tx, info.Txs); err != nil {
		_ = scanner.dbEngine.TxRollBack()
		return fmt.Errorf("TransactionHandler err: %s", err.Error())
	} else {
		scanner.info("TransactionHandler success.")
	}
	scanner.blockCount++
	if block != nil {
		scanner.info(
			"scan block finish, hash:",
			block.BlockHash,
			",block count:", scanner.blockCount, ",this block tx's count:", *info.TxCount, "\n================================")
	}
	return scanner.dbEngine.TxCommit() // 提交事务
}
