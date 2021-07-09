package celltype

import "strings"

/**
 * Copyright (C), 2019-2021
 * FileName: account_char
 * Author:   LinGuanHong
 * Date:     2021/7/2 5:20
 * Description:
 */

func AccountCharsToAccount(accountChars AccountChars) DasAccount {
	index := uint(0)
	accountRawBytes := []byte{}
	accountCharsSize := accountChars.ItemCount()
	for ; index < accountCharsSize; index++ {
		char := accountChars.Get(index)
		accountRawBytes = append(accountRawBytes, char.Bytes().RawData()...)
	}
	accountStr := string(accountRawBytes)
	if accountStr != "" && !strings.HasSuffix(accountStr, DasAccountSuffix) {
		accountStr = accountStr + DasAccountSuffix
	}
	return DasAccount(accountStr)
}

func AccountCharsToAccountId(accountChars AccountChars) DasAccountId {
	index := uint(0)
	accountCharsSize := accountChars.ItemCount()
	accountRawBytes := []byte{}
	for ; index < accountCharsSize; index++ {
		char := accountChars.Get(index)
		accountRawBytes = append(accountRawBytes, char.Bytes().RawData()...)
	}
	accountStr := string(accountRawBytes)
	if !strings.HasSuffix(accountStr, DasAccountSuffix) {
		accountStr = accountStr + DasAccountSuffix
	}
	return DasAccountFromStr(accountStr).AccountId()
}