package main

import (
	"context"
	"das_account_indexer/parser"
	"fmt"
	blockparser "github.com/DeAccountSystems/das_commonlib/chain/ckb_rocksdb_parser"
	"github.com/DeAccountSystems/das_commonlib/db"
	"github.com/af913337456/blockparser/scanner"
	blockparserTypes "github.com/af913337456/blockparser/types"
	"github.com/tecbot/gorocksdb"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"time"

	"das_account_indexer/api"
	"das_account_indexer/config"

	"github.com/DeAccountSystems/das_commonlib/cfg"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/dasrpc"
	"github.com/DeAccountSystems/das_commonlib/sys"
	"github.com/eager7/elog"
	"github.com/fsnotify/fsnotify"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/urfave/cli"
)

/**
 * Copyright (C), 2019-2021
 * FileName: main
 * Author:   LinGuanHong
 * Date:     2021/4/1 3:56
 * Description:
 */

var (
	rpcImpl  *dasrpc.JsonrpcServiceImpl
	txParser *parser.TxParser
	_scanner *scanner.BlockScanner
	log      = elog.NewLogger("server", elog.NoticeLevel)
	rpcWait  = make(chan bool)
	exit     = make(chan bool)
)

func main() {
	app := func() *cli.App {
		return cli.NewApp()
	}()
	app.Action = runServer
	app.HideVersion = true
	globalFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "config file abs path",
		},
	}
	app.Flags = append(app.Flags, globalFlags...)
	app.Commands = []cli.Command{}
	sort.Sort(cli.CommandsByName(app.Commands))
	app.Before = func(ctx *cli.Context) error {
		debug.FreeOSMemory()
		minCore := runtime.NumCPU()
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

func listenAndHandleInterrupt() {
	sys.ListenSysInterrupt(func(sig os.Signal) {
		log.Warn(fmt.Sprintf("signal [%s] to exit server...., time: %s", sig.String(), time.Now().String()))
		if rpcWait != nil {
			rpcWait <- true
		}
		if rpcImpl != nil {
			log.Warn("stop rpc client...")
			rpcImpl.Stop()
		}
		if _scanner != nil {
			_scanner.Stop()
		}
		if txParser != nil {
			txParser.Close()
		}
		log.Info("exist server success!")
		exit <- true
	})
}

func runServer(ctx *cli.Context) error {
	configFilePath := readConfigFilePath(ctx)
	cfg.InitCfgFromFile(configFilePath, &config.Cfg)
	if err := sys.AddConfigFileWatcher(configFilePath, func(optName fsnotify.Op) {
		cfg.InitCfgFromFile(configFilePath, &config.Cfg)
	}); err != nil {
		return fmt.Errorf("AddConfigFileWatcher err: %s", err.Error())
	}
	log.Info("run at config:\n", config.Cfg.ToStr())

	listenAndHandleInterrupt()

	celltype.UseVersion3SystemScriptCodeHash()
	rpcClient, err := rpc.DialWithIndexer(config.Cfg.Chain.CKB.NodeUrl, config.Cfg.Chain.CKB.IndexerUrl)
	if err != nil {
		panic(fmt.Errorf("init rpcClient failed: %s", err.Error()))
	}

	dataPath := config.Cfg.Chain.CKB.LocalStorage.BlockParser.RocksDB.DataPath
	infoPath := config.Cfg.Chain.CKB.LocalStorage.InfoDbDataPath
	if dataPath != "" && infoPath != "" {
		infoDb, err := db.NewDefaultRocksNormalDb(infoPath)
		if err != nil {
			return fmt.Errorf("NewDefaultRocksNormalDb err: %s", err.Error())
		} else {
			txParser = parser.NewParserRpcTx(&parser.InitTxParserParam{
				RpcClient:         rpcClient,
				Rocksdb:           infoDb,
				TargetBlockHeight: uint64(config.Cfg.Chain.CKB.LocalStorage.BlockParser.StartHeight),
				FontBlockNumber:   uint64(config.Cfg.Chain.CKB.LocalStorage.BlockParser.FrontNumber),
			})
			if err = runChainBlockParser(dataPath, rpcClient, txParser, context.TODO()); err != nil {
				return fmt.Errorf("runChainBlockParser err: %s", err.Error())
			}
		}
		log.Info("rpc server need to wait for block info finish sync, stopping ...")
		go func() {
			for {
				select {
				case <-rpcWait:
					return
				default:
					if txParser.BlockSyncFinish() {
						close(rpcWait)
						rpcWait = nil
						if err = runRpcServer(rpcClient, infoDb); err != nil {
							log.Error("runRpcServer err: %s", err.Error())
						}
					}
				}
				time.Sleep(time.Second * 2)
			}
		}()
	} else {
		if err = runRpcServer(rpcClient, nil); err != nil {
			return fmt.Errorf("runRpcServer err: %s", err.Error())
		}
	}
	<-exit
	return nil
}

func runRpcServer(client rpc.Client, accountDb *gorocksdb.DB) error {
	methodPrefix := "das"
	publicPort := config.Cfg.Rpc.Port
	var rpcHandler api.IApi
	if accountDb == nil {
		rpcHandler = api.NewRpcHandler(client)
	} else {
		rpcHandler = api.NewRpcLocalHandler(client, accountDb)
	}
	rpcImpl = dasrpc.NewJsonrpcService(publicPort, &dasrpc.RpcServiceDelegate{
		Name:    methodPrefix,
		Element: rpcHandler,
	})
	log.Info("rpc serve at:", publicPort)
	if err := rpcImpl.Start(func(w http.ResponseWriter, r *http.Request) {
		// append value to request' ctx
		// newCtx := context.WithValue(r.Context(), "X-Real-IP", "")
		// *r = *r.WithContext(newCtx)
	}); err != nil {
		return err
	}
	return nil
}

func readConfigFilePath(ctx *cli.Context) string {
	if configFileAbsPath := ctx.GlobalString("config"); configFileAbsPath != "" {
		return configFileAbsPath
	} else {
		defaultCfgFilePath := "conf/local_server.yaml"
		return defaultCfgFilePath
	}
}

func runChainBlockParser(dataPath string, rpcClient rpc.Client, txParser *parser.TxParser, ctx context.Context) error {
	var (
		reqRetry     *blockparserTypes.RetryConfig = nil
		retryTimes                                 = config.Cfg.Chain.CKB.LocalStorage.BlockParser.ApiReqRetry.Times
		retryDelayMs                               = config.Cfg.Chain.CKB.LocalStorage.BlockParser.ApiReqRetry.DelayMs
		catchDelayMs                               = config.Cfg.Chain.CKB.LocalStorage.BlockParser.ScanControl.CatchDelayMs
		roundDelayMs                               = config.Cfg.Chain.CKB.LocalStorage.BlockParser.ScanControl.RoundDelayMs
	)
	if retryTimes > 0 && retryDelayMs > 0 {
		reqRetry = &blockparserTypes.RetryConfig{
			RetryTime: retryTimes,
			DelayTime: time.Millisecond * time.Duration(retryDelayMs),
		}
	}
	ckbChain := blockparser.NewCKBBlockChainWithRpcClient(ctx, rpcClient, reqRetry)
	_rocksDb := blockparser.NewCKBRocksDb(dataPath, blockparser.MsgHandler{
		Receive: func(info *blockparser.TxMsgData) error {
			if err := txParser.Handle1(*info, &config.Cfg.Chain.CKB.LocalStorage.ParseDelayMs); err != nil {
				log.Error("handle tx err:", err.Error())
				return err
			}
			return nil
		},
		Close: nil,
	})
	if catchDelayMs <= 0 {
		catchDelayMs = int64(time.Millisecond)
	}
	if roundDelayMs <= 0 {
		roundDelayMs = int64(time.Millisecond)
	}
	_scanner = scanner.NewBlockScanner(scanner.InitBlockScanner{
		Chain: ckbChain,
		Db:    _rocksDb,
		Log:   new(myBlockParserLogger),
		Control: blockparserTypes.DelayControl{
			RoundDelay: time.Millisecond * time.Duration(roundDelayMs),
			CatchDelay: time.Millisecond * time.Duration(catchDelayMs),
		},
		FrontNumber: config.Cfg.Chain.CKB.LocalStorage.BlockParser.FrontNumber,
	})
	go func() {
		if startHeight := config.Cfg.Chain.CKB.LocalStorage.BlockParser.StartHeight; startHeight != 0 {
			_ = _scanner.SetStartScannerHeight(startHeight)
		}
		log.Info("start block parser, wait for message coming...")
		if err := _scanner.Start(); err != nil {
			log.Error("block scanner start err: %s", err.Error())
		}
	}()
	return nil
}

type myBlockParserLogger struct{}

func (l *myBlockParserLogger) Info(args ...interface{}) {
	if config.Cfg.Log.Detailed {
		log.Info(args...)
	}
}
func (l *myBlockParserLogger) Error(args ...interface{}) { log.Error(args...) }
func (l *myBlockParserLogger) Warn(args ...interface{})  { log.Warn(args...) }
