package tron_chain

import (
	"fmt"
	tcr "github.com/tron-us/go-common/crypto"
)

/**
 * Copyright (C), 2019-2021
 * FileName: addr
 * Author:   LinGuanHong
 * Date:     2021/7/9 10:36
 * Description:
 */

func PubkeyHexToBase58(address string) (string, error) {
	tAddr, err := tcr.Encode58Check(&address)
	if err != nil {
		return "", fmt.Errorf("encode 58check:%v", err)
	}
	return *tAddr, nil
}

func PubkeyHexFromBase58(address string) (string, error) {
	addr, err := tcr.Decode58Check(&address)
	if err != nil {
		return "", fmt.Errorf("decode base58:%v", err)
	}
	return *addr, nil
}
