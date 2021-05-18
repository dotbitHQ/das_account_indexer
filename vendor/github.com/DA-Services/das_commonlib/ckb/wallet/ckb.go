package wallet

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	ethSecp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/minio/blake2b-simd"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/bech32"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/secp256k1"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
)

/**
 * Copyright (C), 2019-2020
 * FileName: contract_owner
 * Author:   LinGuanHong
 * Date:     2020/12/21 10:10 下午
 * Description:
 */

const (
	PREFIX_MAINNET = "ckb"
	PREFIX_TESTNET = "ckt"
)

type CkbWalletObj struct {
	SystemScripts *utils.SystemScripts
	Secp256k1Key  *secp256k1.Secp256k1Key
	LockScript    *types.Script
}

func InitCkbWallet(privateKeyHex string, systemScript *utils.SystemScripts) (*CkbWalletObj, error) {
	key, err := secp256k1.HexToKey(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("InitCkbWallet HexToKey err: %s", err.Error())
	}
	lockScript, err := key.Script(systemScript)
	if err != nil {
		return nil, fmt.Errorf("InitCkbWallet LockScript err: %s", err.Error())
	}
	return &CkbWalletObj{
		Secp256k1Key:  key,
		LockScript:    lockScript,
		SystemScripts: systemScript,
	}, nil
}

func GetShortAddressFromLockScriptArgs(args string) {

}

type NewWalletObj struct {
	PriKeyHex  string
	PubKeyHex  string
	AddressHex string
}
func CreateCKBWallet(isTestNet bool) (*NewWalletObj,error) {
	seed := rand.Reader
	keyPair, err := GenerateKey(seed)
	if err != nil {
		return nil, fmt.Errorf("GenerateKey err:%s",err.Error())
	}
	rawPubKey := keyPair.PublicKey
	privBytes := keyPair.ToBytes()
	privKey := byteString(privBytes)
	compressionPubKey := rawPubKey.ToBytes()
	pubKey := byteString(compressionPubKey)
	blake160 := genBlake160(compressionPubKey)
	if isTestNet {
		addr,err := genCkbAddr(PREFIX_TESTNET, blake160)
		return &NewWalletObj{
			PriKeyHex:  privKey,
			PubKeyHex:  pubKey,
			AddressHex: addr,
		}, err
	} else {
		addr,err := genCkbAddr(PREFIX_MAINNET, blake160)
		return &NewWalletObj{
			PriKeyHex:  privKey,
			PubKeyHex:  pubKey,
			AddressHex: addr,
		}, err
	}
}

func GetLockScriptArgsFromShortAddress(address string) (string, error) {
	_, bys, err := bech32.Decode(address)
	if err != nil {
		return "", fmt.Errorf("bech32.Decode err: %s", err.Error())
	}
	converted, err := bech32.ConvertBits(bys, 5, 8, false)
	if err != nil {
		return "", fmt.Errorf("bech32.ConvertBits err: %s", err.Error())
	}
	ret := hex.EncodeToString(converted)[4:]
	const bysSize = 40
	if size := len(ret); size != bysSize {
		return "", fmt.Errorf("invalid args bytes len, want: %d, your: %d",bysSize, size)
	}
	return ret, nil
}

func VerifySign(msg []byte, sign []byte, ckbPubkeyHex string) (bool, error) {
	recoveredPub, err := crypto.Ecrecover(msg, sign)
	if err != nil {
		return false, err
	}
	pubKey, err := crypto.UnmarshalPubkey(recoveredPub)
	if err != nil {
		return false, err
	}
	return hex.EncodeToString(ethSecp256k1.CompressPubkey(pubKey.X, pubKey.Y)) == ckbPubkeyHex, nil
}

func byteString(b []byte) (s string) {
	s = ""
	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02x", b[i])
	}
	return s
}

func genBlake160(pubKeyBin []byte) []byte {
	sum := blake2b.Sum256(pubKeyBin)
	return sum[:20]
}

func genCkbAddr(prefix string, blake160Addr []byte) (string,error) {
	pType, _ := hex.DecodeString("01")
	flag, _ := hex.DecodeString("00")
	payload := append(pType, flag...)
	payload = append(payload, blake160Addr...)

	converted, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("ConvertBits err:%s",err.Error())
	}
	if addr, err := bech32.Encode(prefix, converted); err != nil {
		return "", fmt.Errorf("bech32.Encode err:%s",err.Error())
	} else {
		return addr, nil
	}
}