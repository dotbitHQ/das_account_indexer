package celltype

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strconv"
)

/**
 * Copyright (C), 2019-2020
 * FileName: types
 * Author:   LinGuanHong
 * Date:     2020/12/18 3:58
 * Description:
 */

var DasActionWitness = NewDasWitnessData(TableType_Action, []byte{})

type TypeInputCell struct {
	InputIndex uint32          `json:"input_index"`
	Input      types.CellInput `json:"input"`
	LockType   LockScriptType  `json:"lock_type"`
	CellCap    uint64          `json:"cell_cap"`
}

type BuildTransactionRet struct {
	LockType   LockScriptType     `json:"lock_type"`
	Group      []int              `json:"group"`
	WitnessArg *types.WitnessArgs `json:"witness_arg"`
}

type AddDasOutputCallback func(cellCap uint64, outputIndex uint32)

type CellDepWithWitnessFunc func(inputIndex uint32) ([]byte, error)

type CellDepWithWitness struct {
	CellDep        *types.CellDep
	GetWitnessData CellDepWithWitnessFunc
}

type CellWitnessFunc func(inputIndex uint32) ([]byte, error)

// [das, type, table]
type DASWitnessDataObj struct {
	Tag       string    `json:"tag"`
	TableType TableType `json:"table_type"`
	TableBys  []byte    `json:"table_bys"`
}

func NewDasWitnessDataFromSlice(rawData []byte) (*DASWitnessDataObj, error) {
	tempByte := make([]byte,len(rawData))
	copy(tempByte,rawData)
	if size := len(tempByte); size <= 8 { // header's size + min(data)'s size
		return nil, fmt.Errorf("invalid rawData size: %d", size)
	}
	dasStrTag := string(tempByte[:witnessDasCharLen])
	if dasStrTag != witnessDas {
		return nil, fmt.Errorf("invalid dasStrTag, your: %s, want: %s", dasStrTag,witnessDas)
	}
	tableType, err := MoleculeU32ToGo(tempByte[witnessDasCharLen:witnessDasTableTypeEndIndex])
	if err != nil {
		return nil, fmt.Errorf("invalid tableType err: %s", err.Error())
	}
	return &DASWitnessDataObj{
		Tag:       dasStrTag,
		TableType: TableType(tableType),
		TableBys:  tempByte[witnessDasTableTypeEndIndex:],
	}, nil
}

func NewDasWitnessData(tableType TableType, tableBys []byte) *DASWitnessDataObj {
	return &DASWitnessDataObj{
		Tag:       witnessDas,
		TableType: tableType,
		TableBys:  tableBys,
	}
}
func (d *DASWitnessDataObj) ToWitness() []byte {
	if d.TableBys == nil {
		return nil
	}
	bytebuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuf, binary.LittleEndian, d.TableType)
	temp := append([]byte(d.Tag), bytebuf.Bytes()...)
	return append(temp, d.TableBys...)
}

type DASCellBaseInfoDep struct {
	TxHash  types.Hash    `json:"tx_hash"`
	TxIndex uint          `json:"tx_index"`
	DepType types.DepType `json:"dep_type"`
}

func (c DASCellBaseInfoDep) ToDepCell() *types.CellDep {
	return &types.CellDep{
		OutPoint: &types.OutPoint{
			TxHash: c.TxHash,
			Index:  c.TxIndex,
		},
		DepType: c.DepType,
	}
}

type DASCellBaseInfoOut struct {
	CodeHash     types.Hash           `json:"code_hash"`
	CodeHashType types.ScriptHashType `json:"code_hash_type"`
	Args         []byte               `json:"args"`
}

func DASCellBaseInfoOutFromScript(script *types.Script) DASCellBaseInfoOut {
	return DASCellBaseInfoOut{
		CodeHash:     script.CodeHash,
		CodeHashType: script.HashType,
		Args:         script.Args,
	}
}

func (c DASCellBaseInfoOut) Script() *types.Script {
	current := &types.Script{
		CodeHash: c.CodeHash,
		HashType: c.CodeHashType,
		Args:     c.Args,
	}
	return current
}

func (c DASCellBaseInfoOut) SameCodeHash(script *types.Script) bool {
	return c.CodeHash == script.CodeHash
}

func (c DASCellBaseInfoOut) SameScript(script *types.Script) bool {
	return c.Script().Equals(script)
}

type DASCellBaseInfo struct {
	Name string `json:"name"`
	Dep  DASCellBaseInfoDep `json:"dep"`
	Out  DASCellBaseInfoOut `json:"out"`
	ContractTypeScript types.Script `json:"contract_type_script"`
}

type WalletCellParam struct {
	AccountId              DasAccountId    `json:"-"`
	CellCodeInfo           DASCellBaseInfo `json:"cell_code_info"`
	AnyoneCanPayScriptInfo DASCellBaseInfo `json:"anyone_can_pay_script_info"`
}

type ApplyRegisterCellParam struct {
	Version              uint32          `json:"version"`
	PubkeyHashBytes      []byte          `json:"pubkey_hash_bytes"`
	Account              DasAccount      `json:"account"`
	Height               uint64          `json:"height"`
	CellCodeInfo         DASCellBaseInfo `json:"cell_code_info"`
	SenderLockScriptInfo DASCellBaseInfo `json:"sender_lock_script_info"`
}

type PreAccountCellTxDataParam struct {
	NewAccountCellData *PreAccountCellData `json:"-"`
}
type PreAccountCellParam struct {
	Version uint32 `json:"version"`
	// Data                      Data                `json:"data"`
	Account                   DasAccount                `json:"account"`
	TxDataParam               PreAccountCellTxDataParam `json:"-"`
	CellCodeInfo              DASCellBaseInfo           `json:"cell_code_info"`
	AlwaysSpendableScriptInfo DASCellBaseInfo           `json:"always_spendable_script_info"`
}

type RefcellParam struct {
	Version        uint32          `json:"version"`
	Data           string          `json:"data"`
	AccountId      DasAccountId    `json:"-"`
	RefType        RefCellType     `json:"ref_type"`
	CellCodeInfo   DASCellBaseInfo `json:"cell_code_info"`
	UserLockScript DASCellBaseInfo `json:"user_lock_script"`
}

type QuoteCellParam struct {
	Price uint64 `json:"price"`
	CellCodeInfo              DASCellBaseInfo         `json:"cell_code_info"`
}

/**
lock: <always_success>
type:
  code_hash: <on_sale_script>
  type: type
  args: [id] // AccountCell 的 ID
data: hash(data: OnSaleCellData)

witness:
  table Data {
    old: table DataEntityOpt {
    	index: Uint32,
    	version: Uint32,
    	entity: OnSaleCellData
    },
    new: table DataEntityOpt {
      index: Uint32,
      version: Uint32,
      entity: OnSaleCellData
    },
  }

======
table OnSaleCellData {
    // the price of account
    price: Uint64,
}
*/
type OnSaleCellParam struct {
	Version uint32 `json:"version"`
	// Data                      Data            `json:"data"`
	OnSaleCellData            OnSaleCellData  `json:"-"`
	Price                     uint64          `json:"price"`
	AccountId                 DasAccountId    `json:"account_id"`
	CellCodeInfo              DASCellBaseInfo `json:"cell_code_info"`
	AlwaysSpendableScriptInfo DASCellBaseInfo `json:"always_spendable_script_info"`
}

type IncomeCellParam struct {
	Version uint32 `json:"version"`
	// Data                      Data            `json:"data"`
	IncomeCellData            IncomeCellData  `json:"-"`
	CellCodeInfo              DASCellBaseInfo `json:"cell_code_info"`
	AlwaysSpendableScriptInfo DASCellBaseInfo `json:"always_spendable_script_info"`
}

type BiddingCellParam struct {
	Version                   uint32          `json:"version"`
	Data                      Data            `json:"data"`
	Price                     uint64          `json:"price"`
	AccountId                 DasAccountId    `json:"account_id"`
	CellCodeInfo              DASCellBaseInfo `json:"cell_code_info"`
	AlwaysSpendableScriptInfo DASCellBaseInfo `json:"always_spendable_script_info"`
}

// type AccountCommonParam struct {
// 	InstanceId string `json:"instance_id"`
// 	// Quantity   uint64 `json:"quantity"`
// 	// TokenLogic string `json:"token_logic"`
// }
//
// func (a AccountCommonParam) ToBytes() []byte {
// 	retBytes := []byte{}
// 	instanceId := GoHexToMoleculeHash(a.InstanceId)
// 	retBytes = append(retBytes, instanceId.RawData()...)
// 	// quantity := GoUint64ToMoleculeU64(a.Quantity)
// 	// retBytes = append(retBytes, quantity.RawData()...)
// 	// tokenLogic := GoHexToMoleculeHash(a.TokenLogic)
// 	// retBytes = append(retBytes, tokenLogic.RawData()...)
// 	return retBytes
// }

// func AccountCommonParamByteLen() int {
// 	return 32 + CellVersionByteLen
// }

type ProposeCellParam struct {
	// AccountCommonParam
	Version                   uint32           `json:"version"`
	Data                      Data             `json:"data"`
	TxDataParam               ProposalCellData `json:"tx_data_param"`
	CellCodeInfo              DASCellBaseInfo  `json:"cell_code_info"`
	AlwaysSpendableScriptInfo DASCellBaseInfo  `json:"always_spendable_script_info"`
}

// type UpdateAccountCellObj struct {
// 	OldData *AccountCellData
// 	NewData *AccountCellTxDataParam
// }
//
// func (a *UpdateAccountCellObj) ToAccountCell() *AccountCell {
// 	return NewAccountCell(TestNetAccountCell(a.NewData))
// }

type AccountCellTxDataParam struct {
	NextAccountId DasAccountId `json:"next_account_id"`
	// RegisteredAt  uint64          `json:"registered_at"`
	ExpiredAt   uint64          `json:"expired_at"`
	AccountInfo VersionAccountCell `json:"-"`
}

/**
args: [
    owner_code_hash_index,
    owner_pubkey_hash,
    manager_code_hash_index,
    manager_pubkey_hash,
  ]
*/

type DasLockArgsPairParam struct {
	HashIndexType DasLockCodeHashIndexType
	Script types.Script
}

func (d DasLockArgsPairParam) Bytes() []byte {
	if len(d.Script.Args) == DasLockArgsMinBytesLen {
		return d.Script.Args
	}
	return append(d.HashIndexType.Bytes(),d.Script.Args...)
}

type DasLockParam struct {
	OwnerCodeHashIndexByte []byte
	OwnerPubkeyHashByte  []byte
	ManagerCodeHashIndex []byte
	ManagerPubkeyHash []byte
}

func (d *DasLockParam) Bytes() []byte {
	ownerBytes := append(d.OwnerCodeHashIndexByte,d.OwnerPubkeyHashByte...)
	return append(append(ownerBytes,d.ManagerCodeHashIndex...),d.ManagerPubkeyHash...)
}

type AccountCellDatas struct {
	NewAccountCellData *AccountCellTxDataParam `json:"-"`
}
type AccountCellParam struct {
	TestNet                   bool `json:"test_net"`
	TxDataParam               *AccountCellTxDataParam `json:"-"`
	DataBytes                 []byte `json:"data_bytes"`
	Version                   uint32                  `json:"version"`
	CellCodeInfo              DASCellBaseInfo         `json:"cell_code_info"`
	DasLock                   DASCellBaseInfo `json:"das_lock"`
	DasLockParam *DasLockParam `json:"das_lock_param"`
}

type ParseDasWitnessBysDataObj struct {
	WitnessObj            *DASWitnessDataObj
	MoleculeData          *Data
	MoleculeDepDataEntity *DataEntity
	MoleculeOldDataEntity *DataEntity
	MoleculeNewDataEntity *DataEntity
}

func (p ParseDasWitnessBysDataObj) NewEntity() (*DataEntity, uint32, error) {
	if p.MoleculeNewDataEntity != nil && len(p.MoleculeNewDataEntity.inner) > 0 {
		index, err := MoleculeU32ToGo(p.MoleculeNewDataEntity.Index().RawData())
		if err != nil {
			return nil, 0, err
		}
		return p.MoleculeNewDataEntity, index, nil
	}
	return nil, 0, nil
}

func (p ParseDasWitnessBysDataObj) DepEntity() (*DataEntity, uint32, error) {
	if p.MoleculeDepDataEntity != nil && len(p.MoleculeDepDataEntity.inner) > 0 {
		index, err := MoleculeU32ToGo(p.MoleculeDepDataEntity.Index().RawData())
		if err != nil {
			return nil, 0, err
		}
		return p.MoleculeDepDataEntity, index, nil
	}
	return nil, 0, nil
}

func (p ParseDasWitnessBysDataObj) OldEntity() (*DataEntity, uint32, error) {
	if p.MoleculeOldDataEntity != nil && len(p.MoleculeOldDataEntity.inner) > 0 {
		index, err := MoleculeU32ToGo(p.MoleculeOldDataEntity.Index().RawData())
		if err != nil {
			return nil, 0, err
		}
		return p.MoleculeOldDataEntity, index, nil
	}
	return nil, 0, nil
}

type ProposeWitnessSliceDataObject struct {
	AccountId DasAccountId      `json:"account_id"`
	ItemType  AccountCellStatus `json:"item_type"`
	Next      DasAccountId      `json:"next"`
}

func (p ProposeWitnessSliceDataObject) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(
		`{"account_id":"%s","item_type":"%s","next":"%s"}`,
		p.AccountId.HexStr(),
		p.ItemType.Str(), p.Next.HexStr())), nil
}

type ProposeWitnessSliceDataObjectList []ProposeWitnessSliceDataObject

func (p *ProposeWitnessSliceDataObjectList) Add(accountId DasAccountId, nextId DasAccountId, status AccountCellStatus) {
	*p = append(*p, ProposeWitnessSliceDataObject{AccountId: accountId, Next: nextId, ItemType: status})
}

func ProposeWitnessSliceDataObjectListFromBytes(bys []byte) ([]ProposeWitnessSliceDataObjectList, error) {
	proposeCellData, err := ProposalCellDataFromSlice(bys, false)
	if err != nil {
		return nil, err
	}
	retList := []ProposeWitnessSliceDataObjectList{}
	sliceList := proposeCellData.Slices()
	sliceListSize := sliceList.ItemCount()
	index := uint(0)
	for ; index < sliceListSize; index++ {
		sl := sliceList.Get(index)
		slSize := sl.ItemCount()
		proposeItemIndex := uint(0)
		list := []ProposeWitnessSliceDataObject{}
		for ; proposeItemIndex < slSize; proposeItemIndex++ {
			propose := sl.Get(proposeItemIndex)
			itemTypeUint8, err := MoleculeU8ToGo(propose.ItemType().inner)
			if err != nil {
				return nil, err
			}
			list = append(list, ProposeWitnessSliceDataObject{
				AccountId: DasAccountIdFromBytes(propose.AccountId().RawData()),
				ItemType:  AccountCellStatus(itemTypeUint8),
				Next:      DasAccountIdFromBytes(propose.Next().RawData()),
			})
		}
		retList = append(retList, list)
	}
	return retList, nil
}

type ProposeWitnessSliceDataObjectLL []ProposeWitnessSliceDataObjectList

func (p ProposeWitnessSliceDataObjectLL) ToMoleculeProposalCellData(incomeLockScript *types.Script) ProposalCellData {
	sliceList := make([]SL, 0, len(p))
	for _, slice := range p {
		proposeItemList := make([]ProposalItem, 0, len(slice))
		for _, item := range slice {
			accountId := NewAccountIdBuilder().Set(GoBytesToMoleculeAccountBytes(item.AccountId.Bytes())).Build()
			nextAccountId := NewAccountIdBuilder().Set(GoBytesToMoleculeAccountBytes(item.Next.Bytes())).Build()
			proposeItem := NewProposalItemBuilder().
				AccountId(accountId).
				Next(nextAccountId).
				ItemType(GoUint8ToMoleculeU8(uint8(item.ItemType))).
				Build()
			proposeItemList = append(proposeItemList, proposeItem)
		}
		sliceList = append(sliceList, NewSLBuilder().Set(proposeItemList).Build())
	}
	proposalCellData := NewProposalCellDataBuilder().
		ProposerLock(GoCkbScriptToMoleculeScript(*incomeLockScript)).
		// ProposerWallet(GoBytesToMoleculeBytes(proposerWalletId)).
		Slices(NewSliceListBuilder().Set(sliceList).Build()).
		Build()
	return proposalCellData
}

type CalAccountCellExpiredAtParam struct {
	Quote             uint64 `json:"quote"`
	AccountCellCap    uint64 `json:"account_cell_cap"`
	PriceConfigNew    uint64 `json:"price_config_new"`
	AccountBytesLen   uint32 `json:"account_bytes_len"`
	PreAccountCellCap uint64 `json:"pre_account_cell_cap"`
	RefCellCap        uint64 `json:"ref_cell_cap"`
	DiscountRate      uint64 `json:"discount_rate"`
}

func (c CalAccountCellExpiredAtParam) Json() string {
	bys, _ := json.Marshal(c)
	return string(bys)
}

type EditRecordItem struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}
type EditRecordItemList []EditRecordItem

func (l EditRecordItemList) ToMoleculeRecords() (*Records, error) {
	if len(l) == 0 {
		return nil, nil
	}
	records := NewRecordsBuilder()
	for _, item := range l {
		// if item.Label == "" || item.Type == "" || item.Value == "" {
		// 	return nil, errors.New("invalid records, label, value, type cant empty")
		// }
		ttl, _ := strconv.ParseInt(item.TTL, 10, 64)

		record := NewRecordBuilder().
			RecordKey(GoStrToMoleculeBytes(item.Key)).
			RecordValue(GoStrToMoleculeBytes(item.Value)).
			RecordLabel(GoStrToMoleculeBytes(item.Label)).
			RecordTtl(GoUint32ToMoleculeU32(uint32(ttl))).
			RecordType(GoStrToMoleculeBytes(item.Type)).
			Build()
		records.Push(record)
	}
	recordsBuild := records.Build()
	return &recordsBuild, nil
}

type ReqSendEditRecordsTx struct {
	Records EditRecordItemList `json:"records"`
}