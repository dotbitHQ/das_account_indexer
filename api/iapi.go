package api

import (
	"context"
	"github.com/DeAccountSystems/das_commonlib/common"
)

/**
 * Copyright (C), 2019-2021
 * FileName: iapi
 * Author:   LinGuanHong
 * Date:     2021/7/11 5:51
 * Description:
 */

type IApi interface {
	SearchAccount(ctx context.Context, account string) common.ReqResp
	GetAddressAccount(address string) common.ReqResp
	Close()
}
