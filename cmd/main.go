package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"sort"
	"time"

	"das_account_indexer/api"
	"das_account_indexer/config"

	"github.com/DA-Services/das_commonlib/cfg"
	"github.com/DA-Services/das_commonlib/ckb/celltype"
	"github.com/DA-Services/das_commonlib/dasrpc"
	"github.com/DA-Services/das_commonlib/sys"
	"github.com/eager7/elog"
	"github.com/fsnotify/fsnotify"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/urfave/cli"
)

/**
 * Copyright (C), 2019-2021
 * FileName: main
 * Author:   LinGuanHong
 * Date:     2021/4/1 3:56 下午
 * Description:
 */

var (
	rpcImpl *dasrpc.JsonrpcServiceImpl
	log     = elog.NewLogger("server", elog.NoticeLevel)
	exit    = make(chan bool)
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

func runServer(ctx *cli.Context) error {
	configFilePath := readConfigFilePath(ctx)
	cfg.InitCfgFromFile(configFilePath, &config.Cfg)
	if err := sys.AddConfigFileWatcher(configFilePath, func(optName fsnotify.Op) {
		cfg.InitCfgFromFile(configFilePath, &config.Cfg)
	}); err != nil {
		return fmt.Errorf("AddConfigFileWatcher err: %s", err.Error())
	}
	sys.ListenSysInterrupt(func(sig os.Signal) {
		log.Warn(fmt.Sprintf("signal [%s] to exit server...., time: %s", sig.String(), time.Now().String()))
		if rpcImpl != nil {
			log.Warn("stop rpc client...")
			rpcImpl.Stop()
		}
		log.Info("exist server success!")
		exit <- true
	})
	celltype.UseVersion2SystemScriptCodeHash()
	rpcClient, err := rpc.DialWithIndexer(config.Cfg.Chain.CKB.NodeUrl, config.Cfg.Chain.CKB.IndexerUrl)
	if err != nil {
		panic(fmt.Errorf("init rpcClient failed: %s", err.Error()))
	}

	// The current service does not need to send system transactions, so it is unnecessary to synchronize cellDeps.
	// systemScripts, err := utils.NewSystemScripts(rpcClient)
	// if err != nil {
	// 	panic(fmt.Errorf("init NewSystemScripts err: %s", err.Error()))
	// }
	// celltype.TimingAsyncSystemCodeScriptOutPoint(rpcClient, &types.Script{
	// 	CodeHash: systemScripts.SecpSingleSigCell.CellHash,
	// 	HashType: types.HashTypeType,
	// }, nil, nil)
	if err = runRpcServer(rpcClient); err != nil {
		return err
	}
	<-exit
	return nil
}

func runRpcServer(client rpc.Client) error {
	methodPrefix := "das"
	publicPort := config.Cfg.Rpc.Port
	rpcHandler := api.NewRpcHandler(client)
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
