package celltype

import "github.com/DA-Services/das_commonlib/common"

/**
 * Copyright (C), 2019-2020
 * FileName: value
 * Author:   LinGuanHong
 * Date:     2020/12/20 3:12 下午
 * Description:
 */

const (
	witnessDas = "das"
	witnessDasCharLen = 3
	witnessDasTableTypeEndIndex = 7
)
const oneDaySec = uint64(24 * 3600)
const oneYearDays = uint64(365)
const CellVersionByteLen = 4
const MoleculeBytesHeaderSize = 4
const OneCkb = uint64(1e8)
const DasAccountSuffix = ".bit"
const CkbTxMinOutputCKBValue = 61 * OneCkb
const AccountCellDataAccountIdStartIndex = 72
const RefCellBaseCap = 105 * OneCkb
const AccountCellBaseCap = 200 * OneCkb
const IncomeCellBaseCap  = 106 * OneCkb
const WalletCellBaseCap = 84 * OneCkb
const OneYearSec = int64(3600 * 24 * 365)
const HashBytesLen = 32
const ETHScriptLockWitnessBytesLen = 65
const MinAccountCharsLen = 2
const DiscountRateBase = 10000
const DasLockArgsMinBytesLen = 1 + 20 + 1 + 20

const (
	PwLockMainNetCodeHash = "0xbf43c3602455798c1a61a596e0d95278864c552fafe231c063b3fabf97a8febc"
	PwLockTestNetCodeHash = "0x58c5f491aba6d61678b7cf7edf4910b1f5e00ec0cde2f42e0abb4fd9aff25a63"
)

// type cell's args
const (
	ContractCodeHash          = "00000000000000000000000000000000000000000000000000545950455f4944"
	DasPwLockCellCodeArgs     = "d5eee5a3ac9d65658535b4bdad25e22a81c032f5bbdf5ace45605a33482eeb45"
	DasLockCellCodeArgs       = "eedd10c7d8fee85c119daf2077fea9cf76b9a92ddca546f1f8e0031682e65aee"
	ConfigCellCodeArgs        = "34363fad2018db0b3b6919c26870f302da74c3c4ef4456e5665b82c4118eda51"
	WalletCellCodeArgs        = "9b6d4934ad0366a3a047f24778197000d776c45b2dc68b2738477e730b5b6b91"
	ApplyRegisterCellCodeArgs = "c78fa9066af1624e600ccfb21df9546f900b2afe5d7940d91aefc115653f90d9"
	RefCellCodeArgs           = "34572aae7e930aa06fdd58cd7b42d3db005f27a2d11333cf08a74188128fc090"
	PreAccountCellCodeArgs    = "d3f7ad59632a2ebdc2fe9d41aa69708ed1069b074cd8b297b205f835335d3a6b"
	ProposeCellCodeArgs       = "03d0bb128bd10e666975d9a07c148f6abebe811f511e9574048b30600b065b9a"
	AccountCellCodeArgs       = "589c8e33ffde5bd3a6cda1c391f172247a44f826d3752d866050bdd20fa4d34c"
	IncomeCellCodeArgs       = "589c8e33ffde5bd3a6cda1c391f172247a44f826d3752d866050bdd20fa4d34c"
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

var EmptyDataHash = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var EmptyAccountId = DasAccountId{}

type TableType uint32
type AccountCharType uint32
type AccountCellStatus uint8
type DataEntityChangeType uint

func (t TableType) IsConfigType() bool {
	return t >= TableTyte_CONFIG_CELL_ACCOUNT
}

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
	ScriptType_User LockScriptType = 0
	ScriptType_Any  LockScriptType = 1
	ScriptType_ETH  LockScriptType = 2
	ScriptType_BTC  LockScriptType = 3
)

func (l LockScriptType) ToDasLockCodeHashIndexType() DasLockCodeHashIndexType {
	switch l {
	case ScriptType_User:
		return DasLockCodeHashIndexType_CKB_Normal
	case ScriptType_Any:
		return DasLockCodeHashIndexType_CKB_AnyOne
	case ScriptType_ETH:
		return DasLockCodeHashIndexType_ETH_Normal
	default:
		return DasLockCodeHashIndexType_CKB_Normal
	}
}

type DasAccountSearchStatus int

const (
	SearchStatus_NotOpenRegister  DasAccountSearchStatus = -1
	SearchStatus_Registerable     DasAccountSearchStatus = 0
	SearchStatus_ReservedAccount  DasAccountSearchStatus = 1 // 候选
	SearchStatus_OnePriceSell     DasAccountSearchStatus = 2
	SearchStatus_AuctionSell      DasAccountSearchStatus = 3 // 竞拍，候选 -> 竞拍
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
const (
	TableType_ACTION       TableType = 0
	TableType_ACCOUNT_CELL TableType = 1
	// TableType_REGISTER_CELL TableType = 3
	TableType_ON_SALE_CELL     TableType = 2
	TableType_BIDDING_CELL     TableType = 3
	TableType_PROPOSE_CELL     TableType = 4
	TableType_PRE_ACCOUNT_CELL TableType = 5
	TableType_INCOME_CELL 	   TableType = 6

	TableTyte_CONFIG_CELL_ACCOUNT       TableType = 100
	TableTyte_CONFIG_CELL_APPLY         TableType = 101
	TableTyte_CONFIG_CELL_CHARSET         TableType = 102
	TableTyte_CONFIG_CELL_INCOME         TableType = 103

	TableTyte_CONFIG_CELL_MAIN         TableType = 104
	TableTyte_CONFIG_CELL_PRICE         TableType = 105
	TableTyte_CONFIG_CELL_PROPOSAL         TableType = 106
	TableTyte_CONFIG_CELL_PROFITRATE         TableType = 107

	TableTyte_CONFIG_CELL_RECORD_NAMESPACE       TableType = 108
	TableTyte_CONFIG_CELL_PreservedAccount00     TableType = 150
	// TableTyte_CONFIG_CELL_BLOOM_FILTER TableType = 11
)

type DasLockCodeHashIndexType uint8

const (
	DasLockCodeHashIndexType_CKB_Normal DasLockCodeHashIndexType = 0
	DasLockCodeHashIndexType_CKB_MultiS DasLockCodeHashIndexType = 1
	DasLockCodeHashIndexType_CKB_AnyOne DasLockCodeHashIndexType = 2
	DasLockCodeHashIndexType_ETH_Normal DasLockCodeHashIndexType = 3
)

func (t DasLockCodeHashIndexType) Bytes() []byte {
	return common.Uint8ToBytes(uint8(t))
}

func (t DasLockCodeHashIndexType) ToScriptType() LockScriptType {
	switch t {
	case DasLockCodeHashIndexType_CKB_Normal:
		return ScriptType_User
	case DasLockCodeHashIndexType_CKB_AnyOne:
		return ScriptType_Any
	case DasLockCodeHashIndexType_ETH_Normal:
		return ScriptType_ETH
	default:
		return ScriptType_User
	}
}

const (
	/**
	- status ，状态字段：
	  - 0 ，正常；
	  - 1 ，出售中；
	  - 2 ，拍卖中；
	*/
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
)
