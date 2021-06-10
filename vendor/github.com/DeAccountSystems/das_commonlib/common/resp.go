package common

import "github.com/DeAccountSystems/das_commonlib/common/dascode"

/**
 * Copyright (C), 2019-2020
 * FileName: types.
 * Author:   LinGuanHong
 * Date:     2020/12/22 3:26
 * Description:
 */

type ReqResp struct {
	ErrNo  dascode.DAS_CODE `json:"errno"`
	ErrMsg string           `json:"errmsg"`
	Data   interface{}      `json:"data"`
}
