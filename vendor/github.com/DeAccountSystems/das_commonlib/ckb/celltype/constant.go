package celltype

import "github.com/DeAccountSystems/das_commonlib/common"

/**
 * Copyright (C), 2019-2020
 * FileName: value
 * Author:   LinGuanHong
 * Date:     2020/12/20 3:12
 * Description:
 */

const (
	witnessDas = "das"
	witnessDasCharLen = 3
	witnessDasTableTypeEndIndex = 7
)
const OneDaySec = uint64(24 * 3600)
const OneYearDays = uint64(365)
const CellVersionByteLen = 4
const MoleculeBytesHeaderSize = 4
const OneCkb = uint64(1e8)
const DasAccountSuffix = ".bit"
const CkbTxMinOutputCKBValue = 61 * OneCkb
const AccountCellDataAccountIdStartIndex = 72
const AccountCellBaseCap = 200 * OneCkb
const IncomeCellBaseCap  = 106 * OneCkb
const OneYearSec = int64(3600 * 24 * 365)
const HashBytesLen = 32
const ETHScriptLockWitnessBytesLen = 65
const MinAccountCharsLen = 2
const DiscountRateBase = 10000
const DasLockArgsMinBytesLen = 1 + 20 + 1 + 20

const (
	DasCellDataVersion1 = uint32(1)
	DasCellDataVersion2 = uint32(2)
)

type DasNetType int
const (
	DasNetType_Testnet2 DasNetType = 2
	DasNetType_Testnet3 DasNetType = 3
	DasNetType_Mainnet  DasNetType = 0
)

func LatestVersion() uint32 {
	return DasCellDataVersion2
}

var (
	NullDasLockManagerArg = make([]byte,DasLockArgsMinBytesLen / 2 -1)
	RootAccountDataAccountByte = make([]byte,29)
	EmptyDataHash  = make([]byte,HashBytesLen)
	EmptyAccountId = DasAccountId{}
)

const (
	PwLockMainNetCodeHash = "0xbf43c3602455798c1a61a596e0d95278864c552fafe231c063b3fabf97a8febc"
	PwLockTestNetCodeHash = "0x58c5f491aba6d61678b7cf7edf4910b1f5e00ec0cde2f42e0abb4fd9aff25a63"
)

// type cell's args
var (
	ContractCodeHash          = "00000000000000000000000000000000000000000000000000545950455f4944"
	DasPwLockCellCodeArgs     = "d5eee5a3ac9d65658535b4bdad25e22a81c032f5bbdf5ace45605a33482eeb45"
	DasLockCellCodeArgs       = "0xc3fd71e4f537b8d77a412b896304abf1a60daaa7f0fab10f83e8649a4f1e9713"
	ConfigCellCodeArgs        = "0x92610ed55bbc6d865ab4d84da3259606951417c537edb5b47c8cd0bc7b7b492e"
	WalletCellCodeArgs        = "9b6d4934ad0366a3a047f24778197000d776c45b2dc68b2738477e730b5b6b91"
	ApplyRegisterCellCodeArgs = "0x43b56d4fa45b57680b4cea21819ea5100c209ebb9434f141a53a05bdee93e4d6"
	RefCellCodeArgs           = "34572aae7e930aa06fdd58cd7b42d3db005f27a2d11333cf08a74188128fc090"
	PreAccountCellCodeArgs    = "0x9d608a334270b7ee7c5b61422bcb5a6021552fa4ec1f2d31acc02b2c4306265e"
	ProposeCellCodeArgs       = "0xe789cf86f36fe1c67c04b2aad300867d1fc2778511365ce0b169d0518e860175"
	AccountCellCodeArgs       = "0x37844013d5230454359d93dea9074d653f94dadc1a36fbe88fc01ac8456cddc7"
	IncomeCellCodeArgs        = "0x54d53b0db02b7ca2ecaf1cf6bbe5a9011c8ae6e1dba6d45444e1f3f79eb13896"
)

var (
	ActionParam_Owner   = []byte{0}
	ActionParam_Manager = []byte{1}
)

type PwCoreLockScriptType uint8

const (
	PwCoreLockScriptType_ETH  PwCoreLockScriptType = 1
	PwCoreLockScriptType_EOS  PwCoreLockScriptType = 2
	PwCoreLockScriptType_TRON PwCoreLockScriptType = 3
)

type RefCellType uint8

const (
	RefCellType_Owner   = 0
	RefCellType_Manager = 1
)

type TableType uint32
type AccountCharType uint32
type AccountCellStatus uint8
type DataEntityChangeType uint

func (t TableType) IsConfigType() bool {
	return t >= TableType_ConfigCell_Account
}

/**
const ActionData = 0,
const AccountCellData,
const OnSaleCellData,
const BiddingCellData,
const ProposalCellData,
const PreAccountCellData,
const IncomeCellData,
const ConfigCellAccount = 100,
const ConfigCellApply,
const ConfigCellCharSet,
const ConfigCellIncome,
const ConfigCellMain,
const ConfigCellPrice,
const ConfigCellProposal,
const ConfigCellProfitRate,
const ConfigCellRecordKeyNamespace,
const ConfigCellPreservedAccount00 = 150,
*/
func (t TableType) ValidateType() bool {
	return t <= TableType_IncomeCell ||
		(t >= TableType_ConfigCell_Account && t <= TableType_ConfigCell_RecordNamespace) ||
		(t >= TableType_ConfigCell_PreservedAccount00 && t <= TableType_ConfigCell_PreservedAccount19) ||
		(t >=TableType_ConfigCell_CharSetEmoji && t <= TableType_ConfigCell_CharSetHanT)
}
const (
	TableType_Action       TableType = 0
	TableType_AccountCell  TableType = 1
	TableType_OnSaleCell     TableType = 2
	TableType_BidingCell     TableType = 3
	TableType_ProposeCell    TableType = 4
	TableType_PreAccountCell TableType = 5
	TableType_IncomeCell 	 TableType = 6

	TableType_ConfigCell_Account       TableType = 100
	TableType_ConfigCell_Apply         TableType = 101
	TableType_ConfigCell_CharSet        TableType = 102
	TableType_ConfigCell_Income         TableType = 103

	TableType_ConfigCell_Main         TableType = 104
	TableType_ConfigCell_Price         TableType = 105
	TableType_ConfigCell_Proposal         TableType = 106
	TableType_ConfigCell_ProfitRate         TableType = 107

	TableType_ConfigCell_RecordNamespace       TableType = 108

	TableType_ConfigCell_PreservedAccount00     TableType = 10000
	TableType_ConfigCell_PreservedAccount01     TableType = 10001
	TableType_ConfigCell_PreservedAccount02     TableType = 10002
	TableType_ConfigCell_PreservedAccount03     TableType = 10003
	TableType_ConfigCell_PreservedAccount04     TableType = 10004
	TableType_ConfigCell_PreservedAccount05     TableType = 10005
	TableType_ConfigCell_PreservedAccount06     TableType = 10006
	TableType_ConfigCell_PreservedAccount07     TableType = 10007
	TableType_ConfigCell_PreservedAccount08     TableType = 10008
	TableType_ConfigCell_PreservedAccount09     TableType = 10009
	TableType_ConfigCell_PreservedAccount10     TableType = 10010
	TableType_ConfigCell_PreservedAccount11     TableType = 10011
	TableType_ConfigCell_PreservedAccount12     TableType = 10012
	TableType_ConfigCell_PreservedAccount13     TableType = 10013
	TableType_ConfigCell_PreservedAccount14     TableType = 10014
	TableType_ConfigCell_PreservedAccount15     TableType = 10015
	TableType_ConfigCell_PreservedAccount16     TableType = 10016
	TableType_ConfigCell_PreservedAccount17     TableType = 10017
	TableType_ConfigCell_PreservedAccount18     TableType = 10018
	TableType_ConfigCell_PreservedAccount19     TableType = 10019

	TableType_ConfigCell_CharSetEmoji TableType = 100000
	TableType_ConfigCell_CharSetDigit TableType = 100001
	TableType_ConfigCell_CharSetEn    TableType = 100002
	TableType_ConfigCell_CharSetHanS  TableType = 100003
	TableType_ConfigCell_CharSetHanT  TableType = 100004
	// TableType_ConfigCell_BLOOM_FILTER TableType = 11
)

func (a AccountCellStatus) Str() string {
	switch a {
	case AccountWitnessStatus_Exist:
		return "exist"
	case AccountWitnessStatus_New:
		return "new"
	case AccountWitnessStatus_Proposed:
		return "proposed"
	}
	return "unknown"
}

// type CfgCellType int
// const (
// 	CfgCellType_ConfigCellMain        CfgCellType = 0
// 	CfgCellType_ConfigCellRegister    CfgCellType = 1
// 	CfgCellType_ConfigCellBloomFilter CfgCellType = 2
// 	CfgCellType_ConfigCellMarket      CfgCellType = 3
// )

type ChainType uint

const (
	ChainType_CKB  ChainType = 0
	ChainType_ETH  ChainType = 1
	ChainType_BTC  ChainType = 2
	ChainType_TRON ChainType = 3
	ChainType_WX   ChainType = 4
)

type LockScriptType int

const (
	// use to group inputs when combine tx
	ScriptType_User LockScriptType = 0
	ScriptType_Any  LockScriptType = 1
	ScriptType_ETH  LockScriptType = 2
	ScriptType_BTC  LockScriptType = 3
	ScriptType_DasManager_User  LockScriptType = 4
	ScriptType_DasOwner_User    LockScriptType = 5
	ScriptType_TRON  LockScriptType = 6
)

func (l LockScriptType) ToDasLockCodeHashIndexType() DasLockCodeHashIndexType {
	switch l {
	case ScriptType_User:
		return DasLockCodeHashIndexType_CKB_Normal
	case ScriptType_Any:
		return DasLockCodeHashIndexType_CKB_AnyOne
	case ScriptType_ETH:
		return DasLockCodeHashIndexType_ETH_Normal
	case ScriptType_TRON:
		return DasLockCodeHashIndexType_TRON_Normal
	default:
		return DasLockCodeHashIndexType_CKB_Normal
	}
}

type DasAccountSearchStatus int

const (
	SearchStatus_NotOpenRegister  DasAccountSearchStatus = -1
	SearchStatus_Registerable     DasAccountSearchStatus = 0
	SearchStatus_ReservedAccount  DasAccountSearchStatus = 1 
	SearchStatus_OnePriceSell     DasAccountSearchStatus = 2
	SearchStatus_AuctionSell      DasAccountSearchStatus = 3
	SearchStatus_CandidateAccount DasAccountSearchStatus = 4
	SearchStatus_Registered       DasAccountSearchStatus = 5
)

type MarketType int

// 0x01 for primary，0x02 for secondary
const (
	MarketType_Primary   = 1
	MarketType_Secondary = 2
)

const (
	AccountChar_Emoji  AccountCharType = 0
	AccountChar_Number AccountCharType = 1
	AccountChar_En     AccountCharType = 2
	AccountChar_Zh_Cn  AccountCharType = 3
)

type DasLockCodeHashIndexType uint8

const (
	DasLockCodeHashIndexType_CKB_Normal DasLockCodeHashIndexType = 0
	DasLockCodeHashIndexType_CKB_MultiS DasLockCodeHashIndexType = 1
	DasLockCodeHashIndexType_CKB_AnyOne DasLockCodeHashIndexType = 2
	DasLockCodeHashIndexType_ETH_Normal DasLockCodeHashIndexType = 3
	DasLockCodeHashIndexType_TRON_Normal DasLockCodeHashIndexType = 4
)

func (t DasLockCodeHashIndexType) Bytes() []byte {
	return common.Uint8ToBytes(uint8(t))
}

func (t DasLockCodeHashIndexType) ToScriptType(fromOwner bool) LockScriptType {
	switch t {
	case DasLockCodeHashIndexType_CKB_Normal:
		if fromOwner {
			return ScriptType_DasOwner_User
		} else {
			return ScriptType_DasManager_User
		}
	case DasLockCodeHashIndexType_CKB_AnyOne:
		return ScriptType_Any
	case DasLockCodeHashIndexType_ETH_Normal:
		return ScriptType_ETH
	case DasLockCodeHashIndexType_TRON_Normal:
		return ScriptType_TRON
	default:
		return ScriptType_User
	}
}

const (
	AccountCellStatus_Normal AccountCellStatus = 0
	AccountCellStatus_OnSale AccountCellStatus = 1
	AccountCellStatus_OnBid  AccountCellStatus = 2
)

const (
	AccountWitnessStatus_Exist    AccountCellStatus = 0
	AccountWitnessStatus_Proposed AccountCellStatus = 1
	AccountWitnessStatus_New      AccountCellStatus = 2
)

const (
	NewToDep   DataEntityChangeType = 0
	NewToInput DataEntityChangeType = 1
	DepToInput DataEntityChangeType = 2
	depToDep   DataEntityChangeType = 3
)

const (
	CkbSize_AccountCell = 147 * OneCkb
)

const (
	SystemScript_ApplyRegisterCell = "apply_register_cell"
	SystemScript_PreAccoutnCell    = "preaccount_cell"
	SystemScript_AccoutnCell       = "account_cell"
	SystemScript_BiddingCell       = "bidding_cell"
	SystemScript_OnSaleCell        = "on_sale_cell"
	SystemScript_ProposeCell       = "propose_cell"
	SystemScript_WalletCell        = "wallet_cell"
	SystemScript_RefCell           = "ref_cell"
)

const (
	Action_Deploy                = "deploy"
	Action_Config                = "config"
	Action_AccountChain          = "init_account_chain"
	Action_ApplyRegister         = "apply_register"
	Action_RefundApply           = "apply_apply"
	Action_PreRegister           = "pre_register"
	Action_CreateWallet          = "create_wallet"
	Action_DeleteWallet          = "delete_wallet"
	Action_Propose               = "propose"
	Action_TransferAccount       = "transfer_account"
	Action_RenewAccount          = "renew_account"
	Action_ExtendPropose         = "extend_proposal"
	Action_ConfirmProposal       = "confirm_proposal"
	Action_RecyclePropose        = "recycle_proposal"
	Action_WithdrawFromWallet    = "withdraw_from_wallet"
	Action_Register              = "register"
	Action_VoteBiddingList       = "vote_bidding_list"
	Action_PublishAccount        = "publish_account"
	Action_RejectRegister        = "reject_register"
	Action_PublishBiddingList    = "publish_bidding_list"
	Action_BidAccount            = "bid_account"
	Action_EditManager           = "edit_manager"
	Action_EditRecords           = "edit_records"
	Action_CancelBidding         = "cancel_bidding"
	Action_CloseBidding          = "close_bidding"
	Action_RefundRegister        = "refund_register"
	Action_QuotePriceForCkb      = "quote_price_for_ckb"
	Action_StartAccountSale      = "start_account_sale"
	Action_CancelAccountSale     = "cancel_account_sale"
	Action_StartAccountAuction   = "start_account_auction"
	Action_CancelAccountAuction  = "cancel_account_auction"
	Action_AccuseAccountRepeat   = "accuse_account_repeat"
	Action_AccuseAccountIllegal  = "accuse_account_illegal"
	Action_RecycleExpiredAccount = "recycle_expired_account_by_keeper"
	Action_CancelSaleByKeeper    = "cancel_sale_by_keeper"
	Action_CreateIncome          = "create_income"
	Action_ConsolidateIncome     = "consolidate_income"
)
