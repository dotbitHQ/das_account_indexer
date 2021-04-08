package celltype

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strconv"
)

/**
 * Copyright (C), 2019-2020
 * FileName: types
 * Author:   LinGuanHong
 * Date:     2020/12/18 3:58 下午
 * Description:
 */

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

// type InputWithWitness struct {
// 	CellInput      TypeInputCell          `json:"cell_input"`
// 	GetWitnessData CellDepWithWitnessFunc `json:"-"`
// }

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

/**
- [0:3] 3 个字节固定值为 `0x646173`，这是 `das` 三个字母的 ascii 编码，指明接下来的数据是 DAS 系统数据；
- [4:7] 4 个字节为小端编码的 u32 整形，它是对第 8 字节之后数据类型的标识，具体值详见[Type 常量列表](#Type 常量列表)。首先要通过这个标识判断出具体的数据类型，然后才能用 molecule 编码去解码，下文会解释什么是 molecule 编码；
- [8:] 第 8 字节开始往后的都是 molecule 编码的特殊数据结构，其整体结构如下；
*/
func NewDasWitnessDataFromSlice(rawData []byte) (*DASWitnessDataObj, error) {
	if size := len(rawData); size <= 8 { // header'size + min(data)'size
		return nil, fmt.Errorf("invalid rawData size: %d", size)
	}
	tag := string(rawData[:3])
	if tag != witnessDas {
		return nil, fmt.Errorf("invalid tag: %s", tag)
	}
	tableType, err := MoleculeU32ToGo(rawData[3:7])
	if err != nil {
		return nil, fmt.Errorf("invalid tableType err: %s", err.Error())
	}
	return &DASWitnessDataObj{
		Tag:       tag,
		TableType: TableType(tableType),
		TableBys:  rawData[7:],
	}, nil
}

var DasActionWitness = NewDasWitnessData(TableType_ACTION, []byte{})

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

/**
lock: <always_success>
type:
  code_hash: <bidding_script>
  type: type
  args: [id] // AccountCell 的 ID
data: hash(data: BiddingCellData)

witness:
  table Data {
    old: table DataEntityOpt {
    	index: Uint32,
    	version: Uint32,
    	entity: BiddingCellData
    },
    new: table DataEntityOpt {
      index: Uint32,
      version: Uint32,
      entity: BiddingCellData
    },
  }

======
table BiddingCellData {
    // market type, 0x01 for primary，0x02 for secondary
    market_type: Uint8,
    // starting bidding price
    starting_price: Uint64,
    // current bidding price
    current_price: Uint64,
    // latest bidder's lock script
    current_bidder: ScriptOpt,
}
*/
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
	AccountInfo AccountCellData `json:"-"`
}

type AccountCellDatas struct {
	NewAccountCellData *AccountCellTxDataParam `json:"-"`
}
type AccountCellParam struct {
	TxDataParam               *AccountCellTxDataParam `json:"-"`
	Version                   uint32                  `json:"version"`
	CellCodeInfo              DASCellBaseInfo         `json:"cell_code_info"`
	AlwaysSpendableScriptInfo DASCellBaseInfo         `json:"always_spendable_script_info"`
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

func (p ProposeWitnessSliceDataObjectLL) ToMoleculeProposalCellData(incomeLockScript *types.Script, proposerWalletId []byte) ProposalCellData {
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
		ProposerWallet(GoBytesToMoleculeBytes(proposerWalletId)).
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
		if item.Label == "" || item.Type == "" || item.Value == "" { 
			return nil, errors.New("invalid records, label, value, type cant empty")
		}
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