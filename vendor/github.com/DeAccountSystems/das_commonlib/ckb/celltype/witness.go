package celltype

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/common"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
	"math/big"
	"reflect"
	"strings"
)

/**
 * Copyright (C), 2019-2020
 * FileName: util
 * Author:   LinGuanHong
 * Date:     2020/12/18 2:57
 * Description:
 */

func GoBytesToMoleculeHash(bytes []byte) Hash {
	byteArr := [32]Byte{}
	size := len(bytes)
	for i := 0; i < size; i++ {
		byteArr[i] = *ByteFromSliceUnchecked([]byte{bytes[i]})
	}
	return NewHashBuilder().Set(byteArr).Build()
}

func GoHexToMoleculeHash(hexStr string) Hash {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	bytes, _ := hex.DecodeString(hexStr)
	byteArr := [32]Byte{}
	size := len(bytes)
	for i := 0; i < size; i++ {
		byteArr[i] = *ByteFromSliceUnchecked([]byte{bytes[i]})
	}
	return NewHashBuilder().Set(byteArr).Build()
}

func GoUint8ToMoleculeU8(i uint8) Uint8 {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.LittleEndian, i)
	return *Uint8FromSliceUnchecked(bytebuf.Bytes())
}

func GoUint32ToMoleculeU32(i uint32) Uint32 {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.LittleEndian, i)
	return *Uint32FromSliceUnchecked(bytebuf.Bytes())
}

func GoUint32ToBytes(i uint32) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.LittleEndian, i)
	return bytebuf.Bytes()
}

func GoUint64ToBytes(i uint64) []byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.LittleEndian, i)
	return bytebuf.Bytes()
}

func GoUint64ToMoleculeU64(i uint64) Uint64 {
	return *Uint64FromSliceUnchecked(GoUint64ToBytes(i))
}

func GoStrToMoleculeBytes(str string) Bytes {
	if str == "" {
		return BytesDefault()
	}
	strBytes := []byte(str)
	return GoBytesToMoleculeBytes(strBytes)
}

func GoBytesToMoleculeBytes(bys []byte) Bytes {
	_bytesBuilder := NewBytesBuilder()
	for _, bye := range bys {
		_bytesBuilder.Push(*ByteFromSliceUnchecked([]byte{bye}))
	}
	return _bytesBuilder.Build()
}

func GoByteToMoleculeByte(byte byte) Byte {
	return NewByte(byte)
}

func GoTimeUnixToMoleculeBytes(timeSec int64) [8]Byte {
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.LittleEndian, timeSec)
	timestampByteArr := [8]Byte{}
	bytes := bytebuf.Bytes()
	size := len(bytes)
	for i := 0; i < size; i++ {
		timestampByteArr[i] = *ByteFromSliceUnchecked([]byte{bytes[i]})
	}
	return timestampByteArr
}

func GoBytesToMoleculeAccountBytes(bys []byte) [dasAccountIdLen]Byte {
	byteArr := [dasAccountIdLen]Byte{}
	size := len(bys)
	for i := 0; i < size; i++ {
		byteArr[i] = *ByteFromSliceUnchecked([]byte{bys[i]})
	}
	return byteArr
}

func GoCkbScriptToMoleculeScript(script types.Script) Script {
	// data 0x00 ，type 0x01
	ht := 0
	if script.HashType == types.HashTypeType {
		ht = 1
	}
	argBytes := BytesDefault()
	if script.Args != nil {
		argBytes = GoBytesToMoleculeBytes(script.Args)
	}
	return NewScriptBuilder().
		CodeHash(GoHexToMoleculeHash(script.CodeHash.String())).
		HashType(GoByteToMoleculeByte(byte(ht))).
		Args(argBytes).
		Build()
}

func MoleculeScriptToGo(s Script) (*types.Script, error) {
	t, err := MoleculeU8ToGo(s.HashType().AsSlice())
	if err != nil {
		return nil, err
	}
	hashType := types.HashTypeData
	if t == 1 {
		hashType = types.HashTypeType
	}
	return &types.Script{
		CodeHash: types.BytesToHash(s.CodeHash().RawData()),
		HashType: hashType,
		Args:     s.Args().RawData(),
	}, nil
}

func MoleculeRecordsToGo(records Records) EditRecordItemList {
	index := uint(0)
	recordSize := records.ItemCount()
	retList := make([]EditRecordItem, 0, recordSize)
	for ; index < recordSize; index++ {
		record := records.Get(index)
		ttlU32, _ := MoleculeU32ToGo(record.RecordTtl().RawData())
		retList = append(retList, EditRecordItem{
			Key:   string(record.RecordKey().RawData()),
			Type:  string(record.RecordType().RawData()),
			Label: string(record.RecordLabel().RawData()),
			Value: string(record.RecordValue().RawData()),
			TTL:   fmt.Sprintf("%d", ttlU32),
		})
	}
	return retList
}

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

func MoleculeU8ToGo(bys []byte) (uint8, error) {
	var t uint8
	bytesBuffer := bytes.NewBuffer(bys)
	if err := binary.Read(bytesBuffer, binary.LittleEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func MoleculeU32ToGo(bys []byte) (uint32, error) {
	var t uint32
	bytesBuffer := bytes.NewBuffer(bys)
	if err := binary.Read(bytesBuffer, binary.LittleEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func MoleculeU64ToGo(bys []byte) (uint64, error) {
	var t uint64
	bytesBuffer := bytes.NewBuffer(bys)
	if err := binary.Read(bytesBuffer, binary.LittleEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func MoleculeU64ToGo_BigEndian(bys []byte) (uint64, error) {
	var t uint64
	bytesBuffer := bytes.NewBuffer(bys)
	if err := binary.Read(bytesBuffer, binary.BigEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func MoleculeU32ToGoPercentage(bys []byte) (float64, error) {
	v, e := MoleculeU32ToGo(bys)
	if e != nil {
		return 0, e
	}
	a := new(big.Rat).SetFloat64(float64(v))
	b := new(big.Rat).SetInt64(10000)
	r, _ := new(big.Rat).Quo(a, b).Float64()
	return r, nil
}

func CalDasAwardCap(cap uint64, rate float64) (uint64, error) {
	a := new(big.Rat).SetFloat64(float64(cap))
	b := new(big.Rat).SetFloat64(rate)
	r, _ := new(big.Rat).Mul(a, b).Float64()
	return uint64(r), nil
}

func CalAccountSpend(account DasAccount) uint64 {
	return uint64(len([]byte(account))) * OneCkb
}

func ParseTxWitnessToDasWitnessObj(rawData []byte) (*ParseDasWitnessBysDataObj, error) {
	ret := &ParseDasWitnessBysDataObj{}
	dasWitnessObj, err := NewDasWitnessDataFromSlice(rawData)
	if err != nil {
		return nil, fmt.Errorf("fail to parse dasWitness data: %s", err.Error())
	}
	if tableType := dasWitnessObj.TableType; !tableType.ValidateType() {
		return nil, fmt.Errorf("invalid tableType, your: %d", tableType)
	}
	if dasWitnessObj.TableType == TableType_ACTION {
		ret.WitnessObj = DasActionWitness
		return ret, nil
	}
	ret.WitnessObj = dasWitnessObj
	if dasWitnessObj.TableType.IsConfigType() {
		newDataEntity := NewDataEntityBuilder().Entity(GoBytesToMoleculeBytes(dasWitnessObj.TableBys)).Build()
		newOpt := NewDataEntityOptBuilder().Set(newDataEntity).Build()
		data := NewDataBuilder().Dep(DataEntityOptDefault()).Old(DataEntityOptDefault()).New(newOpt).Build()
		ret.MoleculeNewDataEntity = &newDataEntity
		ret.MoleculeData = &data
		return ret, nil
	}
	data := DataFromSliceUnchecked(dasWitnessObj.TableBys)
	ret.MoleculeData = data
	if data.Dep().IsNone() {
		ret.MoleculeDepDataEntity = nil
	} else {
		ret.MoleculeDepDataEntity = DataEntityFromSliceUnchecked(data.Dep().AsSlice())
	}
	if data.Old().IsNone() {
		ret.MoleculeOldDataEntity = nil
	} else {
		ret.MoleculeOldDataEntity = DataEntityFromSliceUnchecked(data.Old().AsSlice())
	}
	ret.MoleculeNewDataEntity = DataEntityFromSliceUnchecked(data.New().AsSlice())
	return ret, nil
}

var (
	accountCellType = reflect.TypeOf(&AccountCellData{})
	versionAccountCellType = reflect.TypeOf(&VersionAccountCell{})
	accountCellVersion1FieldCount = uint(5)
)
type VersionAccountCell struct {
	Version     uint32
	OriginSlice []byte
	CellData    *AccountCellData
}
func (v *VersionAccountCell) AsSlice() []byte {
	versionBys := GoUint32ToBytes(v.Version)
	tempByte := make([]byte,len(versionBys) +len(v.OriginSlice))
	srcByte := append(versionBys,v.OriginSlice...)
	copy(tempByte,srcByte)
	return tempByte
}

func IsVersion2AccountCell(cellData *AccountCellData) bool {
	if cellData.Len() == accountCellVersion1FieldCount {return true}
	empty_TS := TimestampDefault()
	emptyEMA := bytes.Compare(cellData.LastEditManagerAt().RawData(),empty_TS.RawData()) == 0
	emptyERA := bytes.Compare(cellData.LastEditRecordsAt().RawData(),empty_TS.RawData()) == 0
	emptyETA := bytes.Compare(cellData.LastTransferAccountAt().RawData(),empty_TS.RawData()) == 0
	return emptyERA && emptyEMA && emptyETA
}
func versionAndSlice(molecule ICellData) (*Uint32,[]byte,error) {
	version := GoUint32ToMoleculeU32(DasCellDataVersion1) // default is version 1
	if IsInterfaceNil(molecule) {
		return &version, nil, nil
	}
	switch reflect.TypeOf(molecule) {
	case versionAccountCellType:
		sliceBytes := molecule.AsSlice()
		versionUint32 := common.BytesToUint32_LittleEndian(sliceBytes[:4])
		version = GoUint32ToMoleculeU32(versionUint32)
		return &version, sliceBytes[4:],nil
	case accountCellType:
		if accountCellData,err := AccountCellDataFromSlice(molecule.AsSlice(),false); err != nil {
			return nil, nil, fmt.Errorf("AccountCellDataFromSlice err: %s",err.Error())
		} else if IsVersion2AccountCell(accountCellData) { // version 2
			version = GoUint32ToMoleculeU32(DasCellDataVersion2)
		}
		break
	}
	return &version, molecule.AsSlice(),nil
}
func BuildDasCommonMoleculeDataObj(depIndex, oldIndex, newIndex uint32, depMolecule, oldMolecule, newMolecule ICellData) (*Data,error) {
	var (
		versionDep,sliceBytesDep,depErr = versionAndSlice(depMolecule)
		versionOld,sliceBytesOld,oldErr = versionAndSlice(oldMolecule)
		versionNew,sliceBytesNew,newErr = versionAndSlice(newMolecule)
	)
	if depErr != nil {
		return nil, fmt.Errorf("parse version depErr: %s",depErr.Error())
	}
	if oldErr != nil {
		return nil, fmt.Errorf("parse version oldErr: %s",oldErr.Error())
	}
	if newErr != nil {
		return nil, fmt.Errorf("parse version newErr: %s",newErr.Error())
	}
	var (
		depData DataEntity
		oldData DataEntity
		newData = NewDataEntityBuilder().
			Index(GoUint32ToMoleculeU32(newIndex)).
			Version(*versionNew).
			Entity(GoBytesToMoleculeBytes(sliceBytesNew)).
			Build()
		dataBuilder = NewDataBuilder().
				New(NewDataEntityOptBuilder().Set(newData).Build())
	)
	if !IsInterfaceNil(depMolecule) {
		depData = NewDataEntityBuilder().
			Index(GoUint32ToMoleculeU32(depIndex)).
			Version(*versionDep).
			Entity(GoBytesToMoleculeBytes(sliceBytesDep)).
			Build()
		dataBuilder.Dep(NewDataEntityOptBuilder().Set(depData).Build())
	} else {
		dataBuilder.Dep(DataEntityOptDefault())
	}
	if !IsInterfaceNil(oldMolecule) {
		oldData = NewDataEntityBuilder().
			Index(GoUint32ToMoleculeU32(oldIndex)).
			Version(*versionOld).
			Entity(GoBytesToMoleculeBytes(sliceBytesOld)).
			Build()
		dataBuilder.Old(NewDataEntityOptBuilder().Set(oldData).Build())
	} else {
		dataBuilder.Old(DataEntityOptDefault())
	}
	d := dataBuilder.Build()
	return &d, nil
}

type ReqFindTargetTypeScriptParam struct {
	Ctx       context.Context
	RpcClient rpc.Client
	InputList []*types.CellInput
	IsLock    bool
	CodeHash  types.Hash
}
type FindTargetTypeScriptRet struct {
	Output *types.CellOutput
	Data   []byte
	Tx     *types.Transaction
}

func FindTargetTypeScriptByInputList(p *ReqFindTargetTypeScriptParam) (*FindTargetTypeScriptRet, error) {
	codeHash := p.CodeHash
	for _, item := range p.InputList {
		tx, err := p.RpcClient.GetTransaction(p.Ctx, item.PreviousOutput.TxHash)
		if err != nil {
			return nil, fmt.Errorf("FindSenderLockScriptByInputList err: %s", err.Error())
		}
		size := len(tx.Transaction.Outputs)
		for i := 0; i < size; i++ {
			output := tx.Transaction.Outputs[i]
			if p.IsLock {
				if output.Lock != nil && output.Lock.CodeHash == codeHash &&
					output.Lock.HashType == types.HashTypeType && item.PreviousOutput.Index == uint(i) {
					return &FindTargetTypeScriptRet{
						Output: output,
						Data:   tx.Transaction.OutputsData[i],
						Tx:     tx.Transaction,
					}, nil
				}
			} else {
				if output.Type != nil &&
					output.Type.CodeHash == codeHash &&
					output.Type.HashType == types.HashTypeType &&
					item.PreviousOutput.Index == uint(i) {
					return &FindTargetTypeScriptRet{
						Output: output,
						Data:   tx.Transaction.OutputsData[i],
						Tx:     tx.Transaction,
					}, nil
				}
			}
		}
	}
	return nil, errors.New("FindSenderLockScriptByInputList not found")
}

// const sameIndexMark = 999999
// func ChangeMoleculeDataSameIndex(changeType DataEntityChangeType, originWitnessData []byte) ([]byte, error) {
// 	return ChangeMoleculeData(changeType,sameIndexMark, originWitnessData)
// }

func ChangeMoleculeData(changeType DataEntityChangeType, index uint32, originWitnessData []byte) ([]byte, error) {
	witnessObj, err := NewDasWitnessDataFromSlice(originWitnessData)
	if err != nil {
		return nil, fmt.Errorf("ChangeMoleculeData NewDasWitnessDataFromSlice err: %s", err.Error())
	}
	oldData, err := DataFromSlice(witnessObj.TableBys, false)
	if err != nil {
		return nil, fmt.Errorf("ChangeMoleculeData DataFromSlice err: %s", err.Error())
	}
	// bys := data.New().AsSlice()
	// dataNewBys := make([]byte, 0, len(bys))
	newData := Data{}
	depToX := func(changeType DataEntityChangeType) error {
		if entityOpt := oldData.Dep(); !entityOpt.IsNone() {
			entity, _ := entityOpt.IntoDataEntity()
			dataEntity := NewDataEntityBuilder().
				Version(*entity.Version()).
				Index(GoUint32ToMoleculeU32(index)). // reset index
				Entity(*entity.Entity()).
				Build()
			dataEntityOpt := NewDataEntityOptBuilder().Set(dataEntity).Build()
			if changeType == DepToInput {
				newData = NewDataBuilder().New(DataEntityOptDefault()).Old(dataEntityOpt).Dep(DataEntityOptDefault()).Build()
			} else if changeType == depToDep {
				newData = NewDataBuilder().New(DataEntityOptDefault()).Old(DataEntityOptDefault()).Dep(dataEntityOpt).Build()
			}
		} else {
			return errors.New("ChangeMoleculeData both new ans dep are empty data")
		}
		return nil
	}
	switch changeType {
	case NewToDep:
		oldNewDataEntity, err := oldData.New().IntoDataEntity()
		if err != nil {
			// no data
			if err := depToX(depToDep); err != nil {
				return nil, err
			}
		} else {
			depDataEntity := NewDataEntityBuilder().
				Version(*oldNewDataEntity.Version()).
				Index(GoUint32ToMoleculeU32(index)).
				Entity(*oldNewDataEntity.Entity()).
				Build()
			depDataEntityOpt := NewDataEntityOptBuilder().Set(depDataEntity).Build()
			newData = NewDataBuilder().New(DataEntityOptDefault()).Old(DataEntityOptDefault()).Dep(depDataEntityOpt).Build()
		}
		break
	case NewToInput:
		oldNewDataEntity, err := oldData.New().IntoDataEntity()
		if err != nil {
			// no data
			if err := depToX(DepToInput); err != nil {
				return nil, err
			}
		} else {
			oldDataEntity := NewDataEntityBuilder().
				Version(*oldNewDataEntity.Version()).
				Index(GoUint32ToMoleculeU32(index)).
				Entity(*oldNewDataEntity.Entity()).
				Build()
			oldDataEntityOpt := NewDataEntityOptBuilder().Set(oldDataEntity).Build()
			newData = NewDataBuilder().New(DataEntityOptDefault()).Old(oldDataEntityOpt).Dep(DataEntityOptDefault()).Build()
		}
		break
	case DepToInput:
		if err := depToX(DepToInput); err != nil {
			return nil, err
		}
		break
	default:
		return nil, errors.New("unSupport changeType")
	}
	newDataBytes := (&newData).AsSlice()
	newWitnessData := NewDasWitnessData(witnessObj.TableType, newDataBytes)
	return newWitnessData.ToWitness(), nil
}

func GetScriptTypeFromLockScript(ckbSysScript *utils.SystemScripts, lockScript *types.Script) (LockScriptType, error) {
	lockCodeHash := lockScript.CodeHash
	switch lockCodeHash {
	case ckbSysScript.SecpSingleSigCell.CellHash:
		return ScriptType_User, nil
	case DasAnyOneCanSendCellInfo.Out.CodeHash:
		return ScriptType_Any, nil
	case DasETHLockCellInfo.Out.CodeHash:
		return ScriptType_ETH, nil
	case DasBTCLockCellInfo.CodeHash:
		return ScriptType_BTC, nil
	default:
		return -1, errors.New("invalid lockScript")
	}
}

func IsValidETHLockScriptSignature(signBytes []byte) error {
	if len(signBytes) != ETHScriptLockWitnessBytesLen {
		return fmt.Errorf("invalid signed bys, signed bytes len: %d", ETHScriptLockWitnessBytesLen)
	}
	if signBytes[0] != byte(PwCoreLockScriptType_ETH) {
		return fmt.Errorf("invalid signed bys, first byte must 1, %d", signBytes[0])
	}
	return nil
}

func CalTypeIdFromScript(script *types.Script) types.Hash {
	bys, _ := script.Serialize()
	bysRet, _ := blake2b.Blake256(bys)
	return types.BytesToHash(bysRet)
}

type SkipHandle func(err error)
type ValidHandle func(rawWitnessData []byte, witnessParseObj *ParseDasWitnessBysDataObj) (bool, error)

func GetTargetCellFromWitness(tx *types.Transaction, handle ValidHandle, skipHandle SkipHandle) error {
	inputSize := len(tx.Inputs)
	witnessSize := len(tx.Witnesses)
	for i := inputSize + 1; i < witnessSize; i++ { // (inputSize + 1) skip action cell
		rawWitnessBytes := tx.Witnesses[i]
		if dasObj, err := ParseTxWitnessToDasWitnessObj(rawWitnessBytes); err != nil {
			skipHandle(fmt.Errorf("GetTargetCellFromTx ParseTxWitnessToDasWitnessObj err: %s, skip this one", err.Error()))
		} else {
			if stop, resp := handle(rawWitnessBytes, dasObj); resp != nil {
				return resp
			} else if stop {
				break
			}
		}
	}
	return nil
}

func IsInterfaceNil(i interface{}) bool {
	ret := i == nil
	if !ret {
		defer func() {
			recover()
		}()
		ret = reflect.ValueOf(i).IsNil()
	}
	return ret
}
