package common

import "time"

/**
 * Copyright (C), 2019-2021
 * FileName: req
 * Author:   LinGuanHong
 * Description:
 */

func RetryReq(retryTime int, delayTime time.Duration, reqFunc func() (interface{}, error)) (interface{}, error) {
	var (
		ret interface{}
		err error
	)
	for i := 0; i <= retryTime; i++ {
		if ret, err = reqFunc(); err != nil {
			time.Sleep(delayTime)
			continue
		}
		break
	}
	return ret, nil
}