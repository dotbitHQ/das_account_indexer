package gotype

import (
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
)

/**
 * Copyright (C), 2019-2021
 * FileName: account_register
 * Author:   LinGuanHong
 * Date:     2021/6/24 4:03
 * Description:
 */

func registerFee(price, quote, discount uint64) uint64 {
	// CKB 年费 = CKB 年费 - (CKB 年费 * 折扣率 / 10000)
	if discount >= celltype.DiscountRateBase {
		discount = celltype.DiscountRateBase - 1
	}
	var retVal uint64
	if price < quote {
		retVal = (price * celltype.OneCkb) / quote
	} else {
		retVal = (price / quote) * celltype.OneCkb
	}
	if discount == 0 {
		return retVal
	}
	retVal = retVal - (retVal*discount)/celltype.DiscountRateBase
	return retVal
}

func CalPreAccountCellCap(years uint, price, quote, discountRate uint64, configCell *ConfigCell, account celltype.DasAccount, isRenew bool) (uint64,error) {
	registerYearFee := registerFee(price, quote, discountRate) * uint64(years)
	if isRenew {
		return registerYearFee, nil
	}
	accountRegisterFee,err := AccountCellRegisterCap(configCell,account)
	return registerYearFee + accountRegisterFee, err
}

func AccountCellRegisterCap(configCell *ConfigCell,account celltype.DasAccount) (uint64,error) {
	accountCellBaseCap,err := configCell.AccountCellBaseCap()
	if err != nil {
		return 0, fmt.Errorf("AccountCellBaseCap err: %s", err.Error())
	}
	accountCellPrepareCap,err := configCell.AccountCellPrepareCap()
	if err != nil {
		return 0, fmt.Errorf("AccountCellPrepareCap err: %s", err.Error())
	}
	return accountCellBaseCap + uint64(len(account.Bytes())) * celltype.OneCkb + accountCellPrepareCap, nil
}

func CalAccountCellExpiredAt(param celltype.CalAccountCellExpiredAtParam, registerAt int64) (uint64, error) {
	// fmt.Println("CalAccountCellExpiredAt Param ====>", param.Json())
	if param.PreAccountCellCap < param.AccountCellCap+param.RefCellCap {
		return 0, fmt.Errorf("CalAccountCellExpiredAt invalid cap, preAccCell: %d, accCell: %d", param.PreAccountCellCap, param.AccountCellCap)
	} else {
		paid := param.PreAccountCellCap - param.AccountCellCap - param.RefCellCap
		registerFee := registerFee(param.PriceConfigNew, param.Quote, param.DiscountRate)
		durationInt := paid * celltype.OneYearDays / registerFee * celltype.OneDaySec
		// fmt.Println("CalAccountCellExpiredAt registerFee ====>", registerFee)
		// fmt.Println("CalAccountCellExpiredAt storageFee ====>", paid)
		// fmt.Println("CalAccountCellExpiredAt duration   ====>", durationInt)
		return uint64(registerAt) + durationInt, nil // 1648195213
	}
}

func CalBuyAccountYearSec(years uint) int64 {
	return celltype.OneYearSec * int64(years)
}
