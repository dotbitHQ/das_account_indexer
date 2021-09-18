package gotype

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: account_cell
 * Author:   LinGuanHong
 * Date:     2021/1/24 12:04
 * Description:
 */

type AccountCell struct {
	CellCap       uint64                     `json:"cell_cap"`
	AccountId     celltype.DasAccountId      `json:"account_id"`
	Status        uint8                      `json:"status"`
	Point         types.OutPoint             `json:"point"`
	WitnessStatus celltype.AccountCellStatus `json:"witness_status"`
	Data          []byte                     `json:"-"`
	WitnessData   []byte                     `json:"-"`
	DasLockArgs   []byte                     `json:"-"`
}

type PreCurAccountCell struct {
	Pre *AccountCell `json:"pre"`
	Cur *AccountCell `json:"cur"`
}

type AccountCellList []AccountCell

func (a AccountCellList) Size() int {
	return len(a)
}

func (a AccountCellList) Obj(i int) interface{} {
	return &a[i]
}

func (a AccountCellList) Hash(i int) types.Hash {
	return a[i].Point.TxHash
}

func (a AccountCellList) AccountCellExistMap() map[types.Hash]bool {
	retMap := map[types.Hash]bool{}
	for _, item := range a {
		retMap[item.Point.TxHash] = true
	}
	if len(retMap) > 0 {
		return retMap
	}
	return nil
}

func (a *AccountCell) ExpiredAt() (int64, error) {
	return celltype.ExpiredAtFromOutputData(a.Data)
}

func (a *AccountCell) CellDep() *types.CellDep {
	return &types.CellDep{OutPoint: &a.Point, DepType: types.DepTypeCode}
}

func (a *AccountCell) CellInput() *types.CellInput {
	return &types.CellInput{
		Since:          0,
		PreviousOutput: &a.Point,
	}
}

func (a *AccountCell) GetDepWitness() *celltype.CellDepWithWitness {
	return &celltype.CellDepWithWitness{
		CellDep: a.CellDep(),
		GetWitnessData: func(index uint32) ([]byte, error) {
			newBys, err := celltype.ChangeMoleculeData(celltype.NewToDep, index, a.WitnessData)
			if err != nil {
				return nil, fmt.Errorf(
					"AccountCell CellDepWithWitness err: %s, accountCell witnessStatus: %d, txHash: %s",
					err.Error(), a.WitnessStatus, a.Point.TxHash.String())
			}
			return newBys, nil
		}}
}

func (a *AccountCell) SetNextAccountId(accountId celltype.DasAccountId) {
	a.Data = celltype.SetAccountCellNextAccountId(a.Data, accountId)
}

func (a *AccountCell) NextAccountId() (celltype.DasAccountId, error) {
	return celltype.NextAccountIdFromOutputData(a.Data)
}

func (a *AccountCell) SameTx(target *AccountCell) bool {
	return a.Point.TxHash == target.Point.TxHash
}

func (a *AccountCell) Bytes() ([]byte, error) {
	return rlp.EncodeToBytes(a)
}

func (a *AccountCell) dataObj() (*celltype.Data, error) {
	if a.WitnessData == nil {
		return nil, nil
	} else if a.WitnessStatus == celltype.AccountWitnessStatus_Proposed {
		return nil, nil
	}
	var err error
	witnessObj, err := celltype.NewDasWitnessDataFromSlice(a.WitnessData)
	if err != nil {
		return nil, fmt.Errorf("accountCellData NewDasWitnessDataFromSlice err: %s", err.Error())
	}
	data, err := celltype.DataFromSlice(witnessObj.TableBys, false)
	if err != nil {
		return nil, fmt.Errorf("accountCellData DataFromSlice err: %s", err.Error())
	}
	return data, nil
}

func (a *AccountCell) GetOldAccountCellData() (*celltype.VersionAccountCell, error) {
	data, err := a.dataObj()
	if err != nil {
		return nil, fmt.Errorf("GetOldAccountCellData DataFromSlice err: %s", err.Error())
	}
	cellData, err := data.New().IntoDataEntity() // outside use it to fill input position, means old data
	if err != nil {
		// old no data, get from dep
		cellData, err = data.Dep().IntoDataEntity()
		if err != nil {
			return nil, fmt.Errorf("GetOldAccountCellData dep.IntoDataEntity err: %s", err.Error())
		}
	}
	return compatibleParse(cellData)
}

func (a *AccountCell) GetNewAccountCellData() (*celltype.VersionAccountCell, error) {
	data, err := a.dataObj()
	if err != nil {
		return nil, fmt.Errorf("GetNewAccountCellData DataFromSlice err: %s", err.Error())
	} else if data == nil {
		return nil, nil
	}
	cellData, err := data.New().IntoDataEntity()
	if err != nil {
		return nil, fmt.Errorf("GetNewAccountCellData new.IntoDataEntity err: %s", err.Error())
	}
	return compatibleParse(cellData)
}

func compatibleParse(cellData *celltype.DataEntity) (*celltype.VersionAccountCell, error) {
	var accountCellData *celltype.VersionAccountCell
	version, err := celltype.MoleculeU32ToGo(cellData.Version().RawData())
	if err != nil {
		return nil, fmt.Errorf("compatibleParse MoleculeU32ToGo index err: %s", err.Error())
	}
	switch version {
	case celltype.DasCellDataVersion1:
		if accountCellDataV1, err := celltype.AccountCellDataV1FromSlice(cellData.Entity().RawData(), false); err != nil {
			return nil, fmt.Errorf("compatibleParse AccountCellDataV1FromSlice err: %s", err.Error())
		} else {
			newCellData := celltype.NewAccountCellDataBuilder().
				Records(*accountCellDataV1.Records()).
				Id(*accountCellDataV1.Id()).
				Status(*accountCellDataV1.Status()).
				Account(*accountCellDataV1.Account()).
				RegisteredAt(*accountCellDataV1.RegisteredAt()).
				LastTransferAccountAt(celltype.TimestampDefault()).
				LastEditRecordsAt(celltype.TimestampDefault()).
				LastEditManagerAt(celltype.TimestampDefault()).
				Build()
			accountCellData = &celltype.VersionAccountCell{
				Version:     celltype.DasCellDataVersion1,
				OriginSlice: cellData.Entity().RawData(),
				CellData:    &newCellData,
			}
		}
		break
	case celltype.DasCellDataVersion2:
		newCellData, err := celltype.AccountCellDataFromSlice(cellData.Entity().RawData(), false)
		if err != nil {
			return nil, fmt.Errorf("compatibleParse AccountCellDataFromSlice err: %s", err.Error())
		}
		accountCellData = &celltype.VersionAccountCell{
			Version:     celltype.DasCellDataVersion2,
			OriginSlice: cellData.Entity().RawData(),
			CellData:    newCellData,
		}
		break
	default:
		return nil, fmt.Errorf("compatibleParse invalid version: %d", version)
	}
	return accountCellData, nil
}

func (a *AccountCell) ParseDasLockArgsIndexType() (celltype.DasLockCodeHashIndexType, celltype.DasLockCodeHashIndexType) {
	tempBytes := make([]byte, len(a.DasLockArgs))
	copy(tempBytes, a.DasLockArgs)
	ownerType := celltype.DasLockCodeHashIndexType(tempBytes[0])
	managerTypeStartIndex := celltype.DasLockArgsMinBytesLen / 2
	managerType := celltype.DasLockCodeHashIndexType(tempBytes[managerTypeStartIndex : managerTypeStartIndex+1][0])
	return ownerType, managerType
}

func (a *AccountCell) DasLockOwnerBytes() []byte {
	if len(a.DasLockArgs) < celltype.DasLockArgsMinBytesLen {
		return nil
	}
	return a.DasLockArgs[:celltype.DasLockArgsMinBytesLen/2]
}

func (a *AccountCell) DasLockManagerBytes() []byte {
	if len(a.DasLockArgs) < celltype.DasLockArgsMinBytesLen {
		return nil
	}
	return a.DasLockArgs[celltype.DasLockArgsMinBytesLen/2:]
}

func (a *AccountCell) SameOwner(indexHashType celltype.DasLockCodeHashIndexType, args []byte) error {
	ownerType, _ := a.ParseDasLockArgsIndexType()
	if len(a.DasLockArgs) < celltype.DasLockArgsMinBytesLen {
		return errors.New("invalid dasLockArgs len")
	}
	end := celltype.DasLockArgsMinBytesLen / 2
	if len(args) == celltype.DasLockArgsMinBytesLen {
		// das-lock
		if indexHashType == celltype.DasLockCodeHashIndexType_712_Normal {
			if ownerType == indexHashType && bytes.Compare(a.DasLockArgs[1:end], args[1:end]) == 0 {
				return nil
			} else if ownerType == celltype.DasLockCodeHashIndexType_ETH_Normal && bytes.Compare(a.DasLockArgs[1:end], args[1:end]) == 0 {
				return nil
			}
		} else {
			if ownerType == indexHashType && bytes.Compare(a.DasLockArgs[1:end], args[1:end]) == 0 {
				return nil
			}
		}
	} else {
		// other, such as pw-lock
		if ownerType == indexHashType && bytes.Compare(a.DasLockArgs[1:end], args) == 0 {
			return nil
		}
	}
	return errors.New("invalid account owner")
}

func (a *AccountCell) SameManager(indexHashType celltype.DasLockCodeHashIndexType, args []byte) error {
	_, managerType := a.ParseDasLockArgsIndexType()
	if len(a.DasLockArgs) < celltype.DasLockArgsMinBytesLen {
		return errors.New("invalid dasLockArgs len")
	}
	start := celltype.DasLockArgsMinBytesLen/2 + 1
	if len(args) == celltype.DasLockArgsMinBytesLen {
		// das-lock
		if indexHashType == celltype.DasLockCodeHashIndexType_712_Normal {
			if managerType == indexHashType && bytes.Compare(a.DasLockArgs[start:], args[start:]) == 0 {
				return nil
			} else if managerType == celltype.DasLockCodeHashIndexType_ETH_Normal && bytes.Compare(a.DasLockArgs[start:], args[start:]) == 0 {
				return nil
			}
		} else {
			if managerType == indexHashType && bytes.Compare(a.DasLockArgs[start:], args[start:]) == 0 {
				return nil
			}
		}
	} else {
		// other, such as pw-lock
		if managerType == indexHashType && bytes.Compare(a.DasLockArgs[start:], args) == 0 {
			return nil
		}
	}
	return errors.New("invalid account manager")
}

func (a *AccountCell) TypeInputCellInConfirmPropose() *celltype.TypeInputCell {
	input := &celltype.TypeInputCell{
		InputIndex: 0,
		Input:      *a.CellInput(),
		LockType:   celltype.ScriptType_Any, // keep same type as preAccountCell, so that they can in the same inputGroup
		CellCap:    a.CellCap,
	}
	return input
}

func (a *AccountCell) TypeInputCell(checkOwnerSign bool) *celltype.TypeInputCell {
	ownerType, managerType := a.ParseDasLockArgsIndexType()
	var targetLockType = celltype.ScriptType_Any
	if checkOwnerSign {
		targetLockType = ownerType.ToScriptType(true)
	} else {
		targetLockType = managerType.ToScriptType(false)
	}
	input := &celltype.TypeInputCell{
		InputIndex: 0,
		Input:      *a.CellInput(),
		LockType:   targetLockType,
		CellCap:    a.CellCap,
	}
	return input
}

type UpdateAccountCellInfo struct {
	OldData           *celltype.VersionAccountCell
	OutputAccountCell *celltype.AccountCell
	NewData           *celltype.VersionAccountCell
}

func (a *AccountCell) setOwner(indexType celltype.DasLockCodeHashIndexType, args []byte) {
	rawBytes := make([]byte, len(args))
	copy(rawBytes, args)
	halfEnd := celltype.DasLockArgsMinBytesLen / 2
	var appendBytes = make([]byte, halfEnd-1)
	if len(rawBytes) == celltype.DasLockArgsMinBytesLen {
		copy(appendBytes, rawBytes[1:halfEnd]) // das-lock
	} else {
		copy(appendBytes, rawBytes) // other, such as pw-lock
	}
	tempBytes := make([]byte, 0, celltype.DasLockArgsMinBytesLen)
	tempBytes = append(tempBytes, indexType.Bytes()...)
	tempBytes = append(tempBytes, appendBytes...)
	tempBytes = append(tempBytes, a.DasLockArgs[halfEnd:]...)
	a.DasLockArgs = tempBytes
}

func (a *AccountCell) setManager(indexType celltype.DasLockCodeHashIndexType, args []byte) {
	rawBytes := make([]byte, len(args))
	copy(rawBytes, args)
	halfEnd := celltype.DasLockArgsMinBytesLen / 2
	var appendBytes = make([]byte, halfEnd-1)
	if len(rawBytes) == celltype.DasLockArgsMinBytesLen {
		copy(appendBytes, rawBytes[halfEnd+1:]) // das-lock
	} else {
		copy(appendBytes, rawBytes) // other, such as pw-lock
	}
	tempBytes := make([]byte, 0, celltype.DasLockArgsMinBytesLen)
	tempBytes = append(tempBytes, a.DasLockArgs[0:halfEnd]...)
	tempBytes = append(tempBytes, indexType.Bytes()...)
	tempBytes = append(tempBytes, appendBytes...)
	a.DasLockArgs = tempBytes
}

func (a *AccountCell) ToDasLockArgParam() *celltype.DasLockParam {
	if size := len(a.DasLockArgs); size < celltype.DasLockArgsMinBytesLen {
		return nil
	}
	endLen := celltype.DasLockArgsMinBytesLen / 2
	return &celltype.DasLockParam{
		OwnerCodeHashIndexByte: a.DasLockArgs[:1],
		OwnerPubkeyHashByte:    a.DasLockArgs[1:endLen],
		ManagerCodeHashIndex:   a.DasLockArgs[endLen : endLen+1],
		ManagerPubkeyHash:      a.DasLockArgs[endLen+1:],
	}
}

type InOutputWitnessCallbackParam struct {
	OldData    *celltype.VersionAccountCell
	NewBuilder *celltype.AccountCellDataBuilder
	ExpiredAt  *int64
}

// callback method return 'true' means use the oldData to newData
func (a *AccountCell) UpdateAccountCellInfos(
	testNet bool,
	callback func(param *InOutputWitnessCallbackParam) bool,
	setOwner func() (indexType celltype.DasLockCodeHashIndexType, args []byte),
	setManager func() (indexType celltype.DasLockCodeHashIndexType, args []byte)) (*UpdateAccountCellInfo, error) {

	oldData, err := a.GetOldAccountCellData()
	if err != nil {
		return nil, fmt.Errorf("GetOldAccountCellData err: %s", err.Error())
	}
	accountId := celltype.NewAccountIdBuilder().Set(celltype.GoBytesToMoleculeAccountBytes(a.AccountId.Bytes())).Build()
	builder := celltype.NewAccountCellDataBuilder().Id(accountId)
	expiredAt, err := a.ExpiredAt()
	if err != nil {
		return nil, fmt.Errorf("parse expiredAt err: %s", err.Error())
	}
	cbp := InOutputWitnessCallbackParam{
		OldData:    oldData,
		NewBuilder: builder,
		ExpiredAt:  &expiredAt,
	}
	var (
		newData celltype.AccountCellData
		useOld  = false
	)
	if callback != nil {
		cbp.NewBuilder.
			Id(*cbp.OldData.CellData.Id()).
			Account(*cbp.OldData.CellData.Account()).
			Status(*cbp.OldData.CellData.Status()).
			RegisteredAt(*cbp.OldData.CellData.RegisteredAt()).
			Records(*cbp.OldData.CellData.Records()).
			LastTransferAccountAt(*cbp.OldData.CellData.LastTransferAccountAt()).
			LastEditRecordsAt(*cbp.OldData.CellData.LastEditRecordsAt()).
			LastEditManagerAt(*cbp.OldData.CellData.LastEditManagerAt())
		if useOld = callback(&cbp); useOld {
			newData = *oldData.CellData
		}
	}
	if !useOld {
		newData = builder.Build()
	}
	nextId, err := a.NextAccountId()
	if err != nil {
		return nil, fmt.Errorf("parse NextAccountId err: %s", err.Error())
	}
	newFullData := &celltype.AccountCellTxDataParam{
		NextAccountId: nextId,
		ExpiredAt:     uint64(*cbp.ExpiredAt),
		AccountInfo: celltype.VersionAccountCell{
			Version:     celltype.LatestVersion(),
			OriginSlice: newData.AsSlice(),
			CellData:    &newData,
		},
	}
	if setOwner != nil {
		ownerIndexType, args := setOwner()
		a.setOwner(ownerIndexType, args)
	}
	if setManager != nil {
		managerIndexType, args := setManager()
		if args == nil {
			args = celltype.NullDasLockManagerArg
		}
		a.setManager(managerIndexType, args)
	}
	return &UpdateAccountCellInfo{
		OldData:           oldData,
		OutputAccountCell: celltype.NewAccountCell(celltype.DefaultAccountCellParam(testNet, newFullData, a.ToDasLockArgParam(), nil)),
		NewData: &celltype.VersionAccountCell{
			Version:     celltype.LatestVersion(),
			OriginSlice: newData.AsSlice(),
			CellData:    &newData,
		},
	}, nil
}

func (a *AccountCell) UpdateAccountCellNextId(testNet bool, nextAccountId celltype.DasAccountId) (*UpdateAccountCellInfo, error) {
	oldData, err := a.GetOldAccountCellData()
	if err != nil {
		return nil, fmt.Errorf("GetOldAccountCellData err: %s", err.Error())
	}
	expired, err := a.ExpiredAt()
	if err != nil {
		return nil, fmt.Errorf("parse ExpiredAt err: %s", err.Error())
	}
	newData := &celltype.AccountCellTxDataParam{
		NextAccountId: nextAccountId,
		AccountInfo:   *oldData,
		ExpiredAt:     uint64(expired),
	}
	return &UpdateAccountCellInfo{
		OldData:           oldData,
		OutputAccountCell: celltype.NewAccountCell(celltype.DefaultAccountCellParam(testNet, newData, a.ToDasLockArgParam(), a.Data)),
		NewData:           oldData,
	}, nil
}

func (a *AccountCell) JudgeExpireStatus(expiredCheck, frozenCheck bool, frozenRange int64) (bool, bool, error) {
	var (
		frozen  bool
		expired bool
		err     error
	)
	if frozenCheck {
		if frozen, err = celltype.IsAccountFrozen(a.Data, time.Now().Unix(), int64(frozenRange)); err != nil {
			return false, false, fmt.Errorf("IsAccountFrozen err: %s", err.Error())
		}
	}
	if expiredCheck {
		if expired, err = celltype.IsAccountExpired(a.Data, time.Now().Unix()); err != nil {
			return false, false, fmt.Errorf("IsAccountExpired err: %s", err.Error())
		}
	}
	return expired, frozen, nil
}

func VersionCompatibleAccountCellDataFromSlice(cellData *celltype.DataEntity) (*celltype.VersionAccountCell, error) {
	return compatibleParse(cellData)
}

func BytesToAccountCellTxValue(bys []byte) (*AccountCell, error) {
	ret := &AccountCell{}
	if err := rlp.DecodeBytes(bys, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func CalDasAwardCap(cap uint64, rate float64) (uint64, error) {
	a := new(big.Rat).SetFloat64(float64(cap))
	b := new(big.Rat).SetFloat64(rate)
	r, _ := new(big.Rat).Mul(a, b).Float64()
	return uint64(r), nil
}

func CalAccountSpend(account celltype.DasAccount) uint64 {
	return uint64(len([]byte(account))) * celltype.OneCkb
}
