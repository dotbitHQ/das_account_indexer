package dascode

/**
 * Copyright (C), 2019-2021
 * FileName: errcode
 * Author:   LinGuanHong
 * Date:     2021/2/4 1:02
 * Description:
 */

type DAS_CODE int

const DAS_SUCCESS DAS_CODE = 0

const (
	Err_CallIndexer            DAS_CODE = 20000
	Err_Internal               DAS_CODE = 20001
	Err_AccountExpired         DAS_CODE = 20002
	Err_AccountFrozen          DAS_CODE = 20003
	Err_AccountAlreadyRegister DAS_CODE = 20004
	Err_AccountRecordsInvalid  DAS_CODE = 20005
	Err_AccountFormatInvalid   DAS_CODE = 20006
	Err_AccountNotExist        DAS_CODE = 20007
	Err_PubkeyHexFormatInvalid DAS_CODE = 20008
	Err_BaseParamInvalid       DAS_CODE = 20009
	Err_CostPriceNotEqual      DAS_CODE = 20010
	Err_PreRegisterTooEarly    DAS_CODE = 20011
	Err_PreRegisterTooLate     DAS_CODE = 20012
	Err_AccountCellUpdateRateLimit     DAS_CODE = 20013
)
