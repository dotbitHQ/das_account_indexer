package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"time"

	"das_account_indexer/cmd/base"
	"das_account_indexer/config"
	"das_account_indexer/parser"
	"das_account_indexer/parser/handler"
	"github.com/urfave/cli"

	"github.com/DeAccountSystems/das_commonlib/cfg"
	blockparser "github.com/DeAccountSystems/das_commonlib/chain/ckb_rocksdb_parser"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/db"
	"github.com/af913337456/blockparser/types"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: main
 * Author:   LinGuanHong
 * Date:     2021/8/11 2:51
 * Description:
 */

var (
	rpcClient rpc.Client
)

var blockNumbersArgName = "block_numbers"

func main() {
	app := func() *cli.App {
		return cli.NewApp()
	}()
	app.Name = "indexer_cli"
	app.HideVersion = true
	globalFlags := []cli.Flag{}
	globalFlags = append(globalFlags, base.CmdBaseFlag...)
	app.Flags = append(app.Flags, globalFlags...)
	app.Commands = commands()
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Before = func(ctx *cli.Context) error {
		debug.FreeOSMemory()
		minCore := runtime.NumCPU() // below go version 1.5,returns 1
		if minCore < 4 {
			minCore = 4
		}
		runtime.GOMAXPROCS(minCore)
		return nil
	}
	app.After = func(ctx *cli.Context) error {
		return nil
	}
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func commands() []cli.Command {
	var (
		ctaa = "consume_transfer_account_action"
	)
	commands := []cli.Command{
		{
			Name:        ctaa,
			ShortName:   ctaa,
			Description: "consume the transfer account action by block number, demo: ./indexer_cli --config=\"\" consume_transfer_account_action --" + blockNumbersArgName + "={1,2,3}",
			ArgsUsage:   "block number array, before run this, you need to stop indexer first",
			Flags: []cli.Flag{
				cli.Int64SliceFlag{
					Name:  blockNumbersArgName,
					Usage: blockNumbersArgName + "={1,2,3}",
				},
			},
			Before:   commonBefore,
			Action:   consumeTransferAccountAction,
			HelpName: ctaa,
		}}
	return commands
}

func commonBefore(ctx *cli.Context) error {
	cfg.InitCfgFromFile(base.ReadConfigFilePath(ctx), &config.Cfg)
	var err error
	rpcClient, err = rpc.DialWithIndexer(config.Cfg.Chain.CKB.NodeUrl, config.Cfg.Chain.CKB.IndexerUrl)
	if err != nil {
		return fmt.Errorf("init rpcClient failed: %s", err.Error())
	}
	if err = base.SetSystemCodeHash(base.ReadNetType(ctx)); err != nil {
		return fmt.Errorf("SetSystemCodeHash failed: %s", err.Error())
	}
	return nil
}

func consumeTransferAccountAction(ctx *cli.Context) error {
	rocksDb, err := db.NewDefaultRocksNormalDb(config.Cfg.Chain.CKB.LocalStorage.InfoDbDataPath)
	if err != nil {
		return fmt.Errorf("NewDefaultReadOnlyRocksDb err: %s", err.Error())
	}
	blockNumbers := ctx.Int64Slice(blockNumbersArgName)
	if blockNumbers == nil || len(blockNumbers) == 0 {
		return fmt.Errorf("param %s can't empty", blockNumbersArgName)
	}
	fmt.Println(fmt.Sprintf("accept %s param: %v", blockNumbersArgName, blockNumbers))
	cmdRegister := handler.NewActionRegisterWithoutListenCmd()
	cmdRegister.Register(celltype.Action_TransferAccount, handler.HandleTransferAccountTx)
	txParser := parser.NewParserRpcTxWithCmdRegister(&parser.InitTxParserParam{
		RpcClient:         rpcClient,
		Rocksdb:           rocksDb,
		Context:           context.TODO(),
		TargetBlockHeight: uint64(config.Cfg.Chain.CKB.LocalStorage.BlockParser.StartHeight),
		FontBlockNumber:   uint64(config.Cfg.Chain.CKB.LocalStorage.BlockParser.FrontNumber),
	}, cmdRegister)
	requester := blockparser.NewCKBBlockChainWithRpcClient(context.TODO(), rpcClient, &types.RetryConfig{
		RetryTime: 2,
		DelayTime: time.Second * 2,
	})
	defer func() {
		txParser.Close() // rocks close inside
		requester.Close()
	}()
	for _, number := range blockNumbers {
		blockData, err := requester.GetBlockInfoByNumber(new(big.Int).SetInt64(number))
		if err != nil {
			return fmt.Errorf("GetBlockInfoByNumber err: %s", err.Error())
		}
		msgData := blockparser.TxMsgData{
			BlockBaseInfo: *blockData,
			Txs:           (blockData.Txs).([]*ckbTypes.Transaction),
		}
		delayMs := int64(200)
		if err = txParser.Handle1(msgData, &delayMs); err != nil {
			return fmt.Errorf("txParser.Handle1 err: %s", err.Error())
		}
	}
	return nil
}
