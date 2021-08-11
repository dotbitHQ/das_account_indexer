package base

import (
	"errors"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/urfave/cli"
)

/**
 * Copyright (C), 2019-2021
 * FileName: base
 * Author:   LinGuanHong
 * Date:     2021/8/11 2:55
 * Description:
 */

var (
	CmdArgNetType = "net_type"
	CmdBaseFlag   = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Usage: "config file abs path",
		},
		cli.IntFlag{
			Name:  CmdArgNetType,
			Usage: "spec indexer's net type. 1 means mainnet,2 means das-test2, 3 means das-test3",
		},
	}
)

func ReadNetType(ctx *cli.Context) celltype.DasNetType {
	return celltype.DasNetType(ctx.GlobalInt(CmdArgNetType))
}

func ReadConfigFilePath(ctx *cli.Context) string {
	if configFileAbsPath := ctx.GlobalString("config"); configFileAbsPath != "" {
		return configFileAbsPath
	} else {
		defaultCfgFilePath := "conf/local_server.yaml"
		return defaultCfgFilePath
	}
}

func SetSystemCodeHash(neyType celltype.DasNetType) error {
	switch neyType {
	case celltype.DasNetType_Mainnet:
		celltype.UseVersionReleaseSystemScriptCodeHash()
		break
	case celltype.DasNetType_Testnet2:
		celltype.UseVersion2SystemScriptCodeHash()
		break
	case celltype.DasNetType_Testnet3:
		celltype.UseVersion3SystemScriptCodeHash()
		break
	default:
		return errors.New("unSupport DasNetType")
	}
	return nil
}
