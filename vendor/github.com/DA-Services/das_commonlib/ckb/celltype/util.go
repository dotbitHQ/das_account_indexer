package celltype

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
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
 * Date:     2020/12/18 2:57 下午
 * Description:
 */

// int64 4Byte
// func PackCellDataWithVersion(version uint32, cellData []byte) []byte {
// 	bytebuf := bytes.NewBuffer([]byte{})
// 	_ = binary.Write(bytebuf, binary.LittleEndian, version)
// 	return append(bytebuf.Bytes(), cellData...)
// }

// func UnpackCellDataWithVersion(cellData []byte) []byte {
// 	return cellData[CellVersionByteLen:]
// }

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

func GoBytesToMoleculeAccountBytes(bys []byte) [10]Byte {
	byteArr := [10]Byte{}
	size := len(bys)
	for i := 0; i < size; i++ {
		byteArr[i] = *ByteFromSliceUnchecked([]byte{bys[i]})
	}
	return byteArr
}

func GoCkbScriptToMoleculeScript(script types.Script) Script {
	// 这里 data 类型应该就是 0x00 ，type 就是 0x01
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
	/**
	[
		{
			emoji
			[]byte("🌹")
		},
		{
			en
			[]byte("a")
		},
		{
			zh
			[]byte("你")
		}
	]
	*/
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

func registerFee(price, quote, discount uint64) uint64 {
	// CKB 年费 = CKB 年费 - (CKB 年费 * 折扣率 / 10000)
	if discount >= DiscountRateBase {
		discount = DiscountRateBase - 1
	}
	var retVal uint64
	if price < quote {
		retVal = (price * OneCkb) / quote
	} else {
		retVal = (price / quote) * OneCkb
	}
	if discount == 0 {
		return retVal
	}
	retVal = retVal - (retVal*discount)/DiscountRateBase
	return retVal
}

func CalPreAccountCellCap(years uint, price, quote, discountRate uint64, account DasAccount) uint64 {
	// PreAccountCell.capacity >= c + AccountCell 基础成本 + RefCell 基础成本 + Account 字节长度
	registerYearFee := registerFee(price, quote, discountRate) * uint64(years)
	storageFee := AccountCellBaseCap //+ 2*RefCellBaseCap
	accountCharFee := uint64(len([]byte(account))) * OneCkb
	return registerYearFee + storageFee + accountCharFee
}

func CalBuyAccountYearSec(years uint) int64 {
	return OneYearSec * int64(years)
}

func ParseTxWitnessToDasWitnessObj(rawData []byte) (*ParseDasWitnessBysDataObj, error) {
	ret := &ParseDasWitnessBysDataObj{}
	dasWitnessObj, err := NewDasWitnessDataFromSlice(rawData)
	if err != nil {
		return nil, fmt.Errorf("fail to parse dasWitness data: %s", err.Error())
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

func BuildDasCommonMoleculeDataObj(depIndex, oldIndex, newIndex uint32, depMolecule, oldMolecule, newMolecule ICellData) *Data {
	var (
		depData DataEntity
		oldData DataEntity
		newData = NewDataEntityBuilder().
			Index(GoUint32ToMoleculeU32(newIndex)).
			Version(GoUint32ToMoleculeU32(1)).
			Entity(GoBytesToMoleculeBytes(newMolecule.AsSlice())).
			Build()
		dataBuilder = NewDataBuilder().
				New(NewDataEntityOptBuilder().Set(newData).Build())
	)
	if !IsInterfaceNil(depMolecule) {
		depData = NewDataEntityBuilder().
			Index(GoUint32ToMoleculeU32(depIndex)).
			Version(GoUint32ToMoleculeU32(1)).
			Entity(GoBytesToMoleculeBytes(depMolecule.AsSlice())).
			Build()
		dataBuilder.Dep(NewDataEntityOptBuilder().Set(depData).Build())
	} else {
		dataBuilder.Dep(DataEntityOptDefault())
	}
	if !IsInterfaceNil(oldMolecule) {
		oldData = NewDataEntityBuilder().
			Index(GoUint32ToMoleculeU32(oldIndex)).
			Version(GoUint32ToMoleculeU32(1)).
			Entity(GoBytesToMoleculeBytes(oldMolecule.AsSlice())).
			Build()
		dataBuilder.Old(NewDataEntityOptBuilder().Set(oldData).Build())
	} else {
		dataBuilder.Old(DataEntityOptDefault())
	}
	d := dataBuilder.Build()
	return &d
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

/**
singlePrice = ConfigCell.price / quote * 10^8 / 365 * 86400
expiredAt = ((PreAccountCell.capacity - AccountCell.capacity - RefCell.capacity) / singlePrice
*/
func CalAccountCellExpiredAt(param CalAccountCellExpiredAtParam, registerAt int64) (uint64, error) {
	fmt.Println("CalAccountCellExpiredAt Param ====>", param.Json())
	if param.PreAccountCellCap < param.AccountCellCap+param.RefCellCap {
		return 0, fmt.Errorf("CalAccountCellExpiredAt invalid cap, preAccCell: %d, accCell: %d", param.PreAccountCellCap, param.AccountCellCap)
	} else {
		paid := param.PreAccountCellCap - param.AccountCellCap - param.RefCellCap
		registerFee := registerFee(param.PriceConfigNew, param.Quote, param.DiscountRate)
		durationInt := paid * oneYearDays / registerFee * oneDaySec
		fmt.Println("CalAccountCellExpiredAt registerFee ====>", registerFee)
		fmt.Println("CalAccountCellExpiredAt storageFee ====>", paid)
		fmt.Println("CalAccountCellExpiredAt duration   ====>", durationInt)
		return uint64(registerAt) + durationInt, nil // 1648195213
	}
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
