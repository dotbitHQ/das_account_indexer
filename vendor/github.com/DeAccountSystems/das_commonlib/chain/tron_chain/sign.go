package tron_chain

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	tcr "github.com/tron-us/go-common/crypto"
	"strings"
)

/**
 * Copyright (C), 2019-2021
 * FileName: sign
 * Author:   LinGuanHong
 * Date:     2021/7/9 10:36
 * Description:
 */

const TronAddrHexPrefix = "41"

func VerifySignTron(address, signature, signMsg string) error {
	signMsgByt, _ := hex.DecodeString(signMsg)
	messageBytes := append([]byte("\x19TRON Signed Message:\n32"), signMsgByt...)
	messageDigest := crypto.Keccak256(messageBytes)
	signHex := strings.TrimPrefix(signature, "0x")
	if len(signHex) != 130 {
		return fmt.Errorf("invalid signature length: %s - %d", signature, len(signHex))
	}
	sigBytes, _ := hex.DecodeString(signHex)
	r := sigBytes[0:32]
	s := sigBytes[32:64]
	v := sigBytes[64]
	pub, err := crypto.Ecrecover(messageDigest, append(append(r, s...), v-27))
	if err != nil {
		return fmt.Errorf("crypto.Ecrecover err: %s", err.Error())
	}
	pubKey, err := crypto.UnmarshalPubkey(pub)
	if err != nil {
		return fmt.Errorf("crypto.UnmarshalPubkey err: %s", err.Error())
	} else {
		recoveredAddr := crypto.PubkeyToAddress(*pubKey)
		tronAddr := TronAddrHexPrefix + recoveredAddr.String()[2:]
		fmt.Println(recoveredAddr.String(), recoveredAddr.Hash())
		base58address, err := tcr.Encode58Check(&tronAddr)
		if err != nil {
			return fmt.Errorf("crypto_tron.Encode58Check err: %s", err.Error())
		}
		if *base58address == address {
			return nil
		}
		return fmt.Errorf("addresses do not match")
	}
}