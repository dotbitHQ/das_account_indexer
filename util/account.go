package util

import (
	"das_account_indexer/types"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/eager7/elog"

	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/DeAccountSystems/das_commonlib/ckb/gotype"
	ckbTypes "github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: account
 * Author:   LinGuanHong
 * Date:     2021/7/10 5:15
 * Description:
 */

var log = elog.NewLogger("util", elog.NoticeLevel)

type parseAccountPackDataObj struct {
	AccountCellData *celltype.AccountCellData
	WitnessData     []byte
	OutputIndex     uint32
}

func ParseChainAccountToJsonFormat(tx *ckbTypes.Transaction, filter types.AccountFilterFunc) ([]types.AccountReturnObj, error) {
	list := []parseAccountPackDataObj{}
	err := celltype.GetTargetCellFromWitness(tx, func(rawWitnessData []byte, witnessParseObj *celltype.ParseDasWitnessBysDataObj) (bool, error) {
		witnessDataObj := witnessParseObj.WitnessObj
		switch witnessDataObj.TableType {
		case celltype.TableType_AccountCell:
			entity, index, err := witnessParseObj.NewEntity()
			if err != nil {
				return false, fmt.Errorf("witnessParseObj.NewEntity err: %s", err.Error())
			}
			if entity == nil {
				return false, fmt.Errorf("accountCell'new entity is nil, skip this tx")
			}
			versionAccount, err := gotype.VersionCompatibleAccountCellDataFromSlice(entity)
			if err != nil {
				return false, fmt.Errorf("VersionCompatibleAccountCellDataFromSlice err: %s", err.Error())
			}
			if filter != nil && !filter(versionAccount.CellData, index) {
				return false, nil // next one
			}
			list = append(list, parseAccountPackDataObj{
				AccountCellData: versionAccount.CellData,
				WitnessData:     rawWitnessData,
				OutputIndex:     index,
			})
			return false, nil
		}
		return false, nil
	}, func(err error) {
		log.Warn("GetTargetCellFromWitness [accountCell] skip one item:", err.Error())
	})
	if err != nil {
		return nil, err
	}
	accountReturnList := make([]types.AccountReturnObj, 0, len(list))
	for _, item := range list {
		if item.AccountCellData == nil {
			return nil, errors.New("skip this account, not match filter")
		}
		nextAccountId, err := celltype.NextAccountIdFromOutputData(tx.OutputsData[item.OutputIndex])
		if err != nil {
			return nil, err
		}
		registerAt, err := celltype.MoleculeU64ToGo(item.AccountCellData.RegisteredAt().RawData())
		if err != nil {
			return nil, fmt.Errorf("parse registerAt err: %s", err.Error())
		}
		expiredAt, err := celltype.ExpiredAtFromOutputData(tx.OutputsData[item.OutputIndex])
		if err != nil {
			return nil, fmt.Errorf("parse expiredAt err: %s", err.Error())
		}
		accountStatus, _ := celltype.MoleculeU8ToGo(item.AccountCellData.Status().RawData())
		cellAccount := gotype.AccountCell{DasLockArgs: tx.Outputs[item.OutputIndex].Lock.Args}
		ownerBys := cellAccount.DasLockOwnerBytes()
		managerBys := cellAccount.DasLockManagerBytes()
		if ownerBys == nil || managerBys == nil {
			return nil, fmt.Errorf("invalid accountCell, owner or manager empty")
		}
		accountReturnList = append(accountReturnList, types.AccountReturnObj{
			OutPoint: ckbTypes.OutPoint{
				TxHash: tx.Hash,
				Index:  uint(item.OutputIndex),
			},
			WitnessHex: hex.EncodeToString(item.WitnessData),
			AccountData: types.AccountData{
				Account:           celltype.AccountCharsToAccount(*item.AccountCellData.Account()).Str(),
				AccountIdHex:      celltype.DasAccountIdFromBytes(item.AccountCellData.Id().RawData()).HexStr(),
				NextAccountIdHex:  nextAccountId.HexStr(),
				RawDasLockArgsHex: hex.EncodeToString(tx.Outputs[item.OutputIndex].Lock.Args),
				OwnerLockArgsHex:  hex.EncodeToString(ownerBys[1:]),
				ManagerLockArgHex: hex.EncodeToString(managerBys[1:]),
				CreateAtUnix:      registerAt,
				ExpiredAtUnix:     uint64(expiredAt),
				Status:            celltype.AccountCellStatus(accountStatus),
				Records:           celltype.MoleculeRecordsToGo(*item.AccountCellData.Records()),
			},
		})
	}
	return accountReturnList, err
}
