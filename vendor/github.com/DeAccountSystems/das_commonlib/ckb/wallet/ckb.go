package wallet

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/bech32"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/secp256k1"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
	"strings"
)

/**
 * Copyright (C), 2019-2020
 * FileName: contract_owner
 * Author:   LinGuanHong
 * Date:     2020/12/21 10:10
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

// payload = type(01) | code hash index(00) | pubkey Blake160
// docs: https://github.com/nervosnetwork/rfcs/blob/master/rfcs/0021-ckb-address-format/0021-ckb-address-format.md
func GetShortAddressFromLockScriptArgs(args string,isTestNet bool) (string,error) {
	prefix := PREFIX_MAINNET
	if isTestNet {
		prefix = PREFIX_TESTNET
	}
	hexStr, err := toHexStrObj(args)
	if err != nil {
		return "", err
	}
	payload := append([]byte{
		uint8(1), // type
		uint8(0)}, // code_hash_index
		hexStr.Bytes()...)
	return encodeAddress(prefix,payload)
}

type NewWalletObj struct {
	PriKeyHex  string
	PubKeyHex  string
	AddressHex string
}

func (n *NewWalletObj) Json() string {
	bys,_ := json.Marshal(n)
	return string(bys)
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

// func VerifySign(msg []byte, sign []byte, ckbPubkeyHex string) (bool, error) {
// 	recoveredPub, err := crypto.Ecrecover(msg, sign)
// 	if err != nil {
// 		return false, err
// 	}
// 	pubKey, err := crypto.UnmarshalPubkey(recoveredPub)
// 	if err != nil {
// 		return false, err
// 	}
// 	return hex.EncodeToString(ethSecp256k1.CompressPubkey(pubKey.X, pubKey.Y)) == ckbPubkeyHex, nil
// }

func byteString(b []byte) (s string) {
	s = ""
	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02x", b[i])
	}
	return s
}

func genBlake160(pubKeyBin []byte) []byte {
	data,_ := blake2b.Blake160(pubKeyBin)
	return data
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

func encodeAddress(prefix string, payload []byte) (string, error) {
	payload, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("bech32.ConvertBits err: %s",err.Error())
	}
	addr, err := bech32.Encode(prefix, payload)
	if err != nil {
		return "", fmt.Errorf("bech32.Encode err: %s",err.Error())
	}
	return addr, nil
}

type HexStrObj struct {
	bytes  []byte
	hexStr string
}

func toHexStrObj(hexStr string) (*HexStrObj, error) {
	HexStrPrefix := "0x"
	if !strings.HasPrefix(hexStr, HexStrPrefix) {
		hexStr = HexStrPrefix + hexStr
	}
	if len(hexStr) == 2 {
		return &HexStrObj{
			bytes:  []byte{},
			hexStr: HexStrPrefix,
		}, nil
	}

	body := hexStr[2:]
	if len(body)%2 == 1 {
		body = "0" + body
	}

	b, err := hex.DecodeString(body)
	if err != nil {
		return nil, fmt.Errorf("DecodeString err: %s",err.Error())
	}
	return &HexStrObj{
		bytes:  b,
		hexStr: HexStrPrefix + body,
	}, nil
}

func (hs *HexStrObj) Bytes() []byte {
	return hs.bytes
}