package types

import "time"

/**
 * Copyright (C), 2019-2020
 * FileName: request
 * Author:   LinGuanHong
 * Date:     2020/12/11 11:23 上午
 * Description:
 */

type RetryConfig struct {
	RetryTime   int
	DelayTime time.Duration
}


