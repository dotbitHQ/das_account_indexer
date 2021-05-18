package celltype

import (
	"bytes"
	"fmt"
	"github.com/Andrew-M-C/go.emoji/official"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"strings"
)

/**
 * Copyright (C), 2019-2021
 * FileName: account_type
 * Author:   LinGuanHong
 * Date:     2021/2/1 2:21 下午
 * Description:
 */

type DasAccount string

func DasAccountFromStr(account string) DasAccount {
	return DasAccount(account)
}

func (d DasAccount) Bytes() []byte {
	if d == "" {
		return []byte{}
	}
	return []byte(d)
}

func (d DasAccount) Format() string {
	temp := string(d)
	if strings.HasSuffix(temp, DasAccountSuffix) {
		temp = strings.Split(temp, ".")[0]
	}
	return temp
}

func (d DasAccount) ValidErr() error {
	if d == "" ||
		!strings.HasSuffix(string(d), DasAccountSuffix) ||
		strings.Contains(string(d), " ") || strings.Contains(string(d), "_") {
		return fmt.Errorf("invalid account:[%s], demo: helloWorld.bit", d)
	}

	if size := d.AccountValidateLen(); size < MinAccountCharsLen {
		return fmt.Errorf("account's char number min is: %d", MinAccountCharsLen)
	}
	return nil
}

func (d DasAccount) AccountValidateLen() int {
	var (
		emojiCharCounter = 0
		unEmojiCharCounter = 0
	)
	d.parseEmojiStr(func(emoji string) string {
		emojiCharCounter ++
		return ""
	}, func() {
		unEmojiCharCounter ++
	})
	return unEmojiCharCounter + emojiCharCounter
}

func (d DasAccount) parseEmojiStr(isEmoji func(emoji string) string,notEmoji func()) string {
	buff := bytes.Buffer{}
	nextIndex := 0
	s := d.Format()
	for i, r := range s {
		if i < nextIndex {
			continue
		}
		match, length := official.AllSequences.HasEmojiPrefix(s[i:])
		if false == match {
			if notEmoji != nil {
				notEmoji()
			}
			buff.WriteRune(r)
			continue
		}
		nextIndex = i + length
		if isEmoji != nil {
			buff.WriteString(isEmoji(s[i : i+length]))
		}
	}
	return buff.String()
}

func (d DasAccount) Str() string {
	return string(d)
}

func (d DasAccount) AccountId() DasAccountId {
	if len(d) == 0 {
		return EmptyAccountId
	}
	bys, _ := blake2b.Blake160([]byte(d))
	id := &DasAccountId{}
	id.SetBytes(bys)
	return *id
}

const dasAccountIdLen = 10

type DasAccountId [dasAccountIdLen]byte

func BytesToDasAccountId(b []byte) DasAccountId {
	var h DasAccountId
	h.SetBytes(b)
	return h
}

func HexToHash(s string) DasAccountId {
	return BytesToDasAccountId(common.FromHex(s))
}

func (d *DasAccountId) SetBytes(b []byte) {
	bLen := len(b)
	if bLen > len(d) {
		b = b[:dasAccountIdLen]
	}
	copy(d[dasAccountIdLen-len(b):], b)
}

func (d DasAccountId) Point() *DasAccountId {
	return &d
}

func (d DasAccountId) Compare(b DasAccountId) int {
	return bytes.Compare(d.Bytes(), b.Bytes())
}

func DasAccountIdFromBytes(accountRawBytes []byte) DasAccountId {
	id := &DasAccountId{}
	id.SetBytes(accountRawBytes)
	return *id
}

func (d DasAccountId) HexStr() string {
	return hexutil.Encode(d[:])
}

func (d DasAccountId) Str() string {
	return d.HexStr()
}

func (d DasAccountId) Bytes() []byte {
	return d[:]
}
