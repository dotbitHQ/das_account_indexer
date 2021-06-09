package celltype

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DA-Services/das_commonlib/ckb/collector"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
	"golang.org/x/sync/syncmap"
	"strings"
	"time"
)

/**
 * Copyright (C), 2019-2020
 * FileName: cell_info
 * Author:   LinGuanHong
 * Date:     2020/12/22 3:01
 * Description:
 */

var (
	TestNetLockScriptDep = DASCellBaseInfoDep{
		TxHash:  types.HexToHash("0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37"),
		TxIndex: 0,
		DepType: types.DepTypeDepGroup,
	}

	TestNetETHSoScriptDep = DASCellBaseInfoDep{
		TxHash:  types.HexToHash("0x3bffc9beff67d5f93b60b378c68a9910ecc936e5bff0348b3bdf99c4f416213d"),
		TxIndex: 0,
		DepType: types.DepTypeCode,
	}
	// 0xb988070e97c6eda68705e146985bcf2d3b3215cbb619eb61337523bc440d42e0
	TestNetCKBSoScriptDep = DASCellBaseInfoDep{
		TxHash:  types.HexToHash("0xe08b6487bab378df62d1abe58faebecdfefc5dc4297627c1f7240441db69355b"),
		TxIndex: 0,
		DepType: types.DepTypeCode,
	}

	DasETHLockCellInfo = DASCellBaseInfo{
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x57a62003daeab9d54aa29b944fc3b451213a5ebdf2e232216a3cfed0dde61b38"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash(PwLockTestNetCodeHash), // default
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
	}
	DasBTCLockCellInfo = DASCellBaseInfoOut{
		CodeHash:     types.HexToHash(""),
		CodeHashType: types.HashTypeType,
		Args:         emptyHexToArgsBytes(),
	}
	DasAnyOneCanSendCellInfo = DASCellBaseInfo{
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x88462008b19c9ac86fb9fef7150c4f6ef7305d457d6b200c8852852012923bf1"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f"), // default
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
	}

	DasLockCellScript = DASCellBaseInfo{
		Name: "das_lock_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x22b7e4a537b107b32d3e1c5704455b30e04a63f0e97347b32155be49510ae0d0"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd"),
			CodeHashType: types.HashTypeType,
			Args:         dasLockDefaultBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(DasLockCellCodeArgs).Bytes(),
		},
	}

	DasAnyOneCanPayCellInfo = DASCellBaseInfoOut{
		CodeHash:     types.HexToHash(utils.AnyoneCanPayCodeHashOnAggron), // default
		CodeHashType: types.HashTypeType,
		Args:         emptyHexToArgsBytes(),
	}
	DasActionCellScript = DASCellBaseInfo{
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash(""),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash(""),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(""),
			HashType: "",
			Args:     nil,
		},
	}
	// DasWalletCellScript = DASCellBaseInfo{
	// 	Name: "wallet_cell",
	// 	Dep: DASCellBaseInfoDep{
	// 		TxHash:  types.HexToHash("0xaac18fd80a6f9265913518e303fe57d1c93d961ef7badbc1289b9dbe667a8ab42"), //"0x440b323f2821aa808c1bad365c10ffb451058864a11f63b5669a5597ac0e8e0f"
	// 		TxIndex: 0,
	// 		DepType: types.DepTypeCode,
	// 	},
	// 	Out: DASCellBaseInfoOut{
	// 		CodeHash:     types.HexToHash("0x9878b226df9465c215fd3c94dc9f9bf6648d5bea48a24579cf83274fe13801d2"),
	// 		CodeHashType: types.HashTypeType,
	// 		Args:         emptyHexToArgsBytes(),
	// 	},
	// 	ContractTypeScript: types.Script{
	// 		CodeHash: types.HexToHash(ContractCodeHash),
	// 		HashType: types.HashTypeType,
	// 		Args:     types.HexToHash(WalletCellCodeArgs).Bytes(),
	// 	},
	// }
	DasApplyRegisterCellScript = DASCellBaseInfo{
		Name: "apply_register_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0xbc4dec1c2a3b1a9bf76df3a66357c62ec4b543abb595b3ed10fe64e126efc509"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0xa2c3a2b18da897bd24391a921956e45d245b46169d6acc9a0663316d15b51cb1"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(ApplyRegisterCellCodeArgs).Bytes(),
		},
	}
	// DasRefCellScript = DASCellBaseInfo{
	// 	Name: "ref_cell",
	// 	Dep: DASCellBaseInfoDep{
	// 		TxHash:  types.HexToHash("0x86a83fc53d64e0cfbc94ccc003b8eee00617c8aa16a2aa1188d41842ee97dc15"),
	// 		TxIndex: 0,
	// 		DepType: types.DepTypeCode,
	// 	},
	// 	Out: DASCellBaseInfoOut{
	// 		CodeHash:     types.HexToHash("0xe79953f024552e6130220a03d2497dc7c2f784f4297c69ba21d0c423915350e5"),
	// 		CodeHashType: types.HashTypeType,
	// 		Args:         emptyHexToArgsBytes(),
	// 	},
	// 	ContractTypeScript: types.Script{
	// 		CodeHash: types.HexToHash(ContractCodeHash),
	// 		HashType: types.HashTypeType,
	// 		Args:     types.HexToHash(RefCellCodeArgs).Bytes(),
	// 	},
	// }
	DasPreAccountCellScript = DASCellBaseInfo{
		Name: "preAccount_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0xb4353dd3ada2b41b8932edbd853a853e81d50b4c8648c1afd93384b946425d15"), //"0x21b25ab337cbbc7aad691f0f767ec5a852bbb8f6b9ff53dd00e0505f72f1f89a"
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x92d6a9525b9a054222982ab4740be6fe4281e65fff52ab252e7daf9306e12e3f"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(PreAccountCellCodeArgs).Bytes(),
		},
	}
	DasProposeCellScript = DASCellBaseInfo{
		Name: "propose_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0xf3cf92357436e6b6438e33c5d68521cac816baff6ef60e9bfc733453a335a8d4"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x4154b5f9114b8d2dd8323eead5d5e71d0959a2dc73f0672e829ae4dabffdb2d8"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(ProposeCellCodeArgs).Bytes(),
		},
	}
	DasAccountCellScript = DASCellBaseInfo{
		Name: "account_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x9e867e0b7bcbd329b8fe311c8839e10bacac7303280b8124932c66f726c38d8a"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x274775e475c1252b5333c20e1512b7b1296c4c5b52a25aa2ebd6e41f5894c41f"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(AccountCellCodeArgs).Bytes(),
		},
	}
	DasBiddingCellScript = DASCellBaseInfo{
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x123"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x123"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
	}
	DasOnSaleCellScript = DASCellBaseInfo{
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x123"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x123"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
	}
	DasIncomeCellScript = DASCellBaseInfo{
		Name: "income_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0xa411dc40662eaf2c43d165c071947e7440e5ec01193954dbf06670bc6bf221c4"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(IncomeCellCodeArgs).Bytes(),
		},
	}
	// DasQuoteCellScript = DASCellBaseInfo{
	// 	Dep: DASCellBaseInfoDep{
	// 		TxHash:  types.HexToHash(""),
	// 		TxIndex: 0,
	// 		DepType: "",
	// 	},
	// 	Out: DASCellBaseInfoOut{
	// 		CodeHash:     types.HexToHash(""),
	// 		CodeHashType: "",
	// 		Args:         nil,
	// 	},
	// }
	DasConfigCellScript = DASCellBaseInfo{
		Name: "config_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x97cf78ef50809505bba4ac78d8ee7908eccd1119aa08775814202e7801f4895b"),
			TxIndex: 0,
			DepType: "",
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x489ff2195ed41aac9a9265c653d8ca57c825b22db765b9e08d537572ff2cbc1b"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(ConfigCellCodeArgs).Bytes(),
		},
	}
	DasHeightCellScript = DASCellBaseInfo{
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x1bc39fc942746cf961f338c33626bfea999c96eb06334541859426580643fd51"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x7a6db6793ecf341f8f5289bc164d4a417c5adb99ab86a750230d7d14e73768e7"),
			CodeHashType: types.HashTypeType,
			Args:         hexToArgsBytes("0x5123c074feef10b58c061b6d16a70a397b30957024f2a262102206213a808d3700000000"),
		},
	}
	DasTimeCellScript = DASCellBaseInfo{
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0xc0f2b262c8dbd5c8da3376cf81f3d3c69582fefcc3eba36e88f708c1a4d505fe"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0xd78423449320291c41adcce741276c47df1dbb0bca212d0017db66297be88f19"),
			CodeHashType: types.HashTypeType,
			Args:         hexToArgsBytes("0x248f00f2a594ae982501113267d487acd27b343e081e04d1fd0490b3288b38d900000000"),
		},
	}
	SystemCodeScriptMap = syncmap.Map{} // map[types.Hash]*DASCellBaseInfo{}
)

func init() {
	initMap()
}

func initMap() {
	SystemCodeScriptMap.Store(DasLockCellScript.Out.CodeHash, &DasLockCellScript)
	SystemCodeScriptMap.Store(DasApplyRegisterCellScript.Out.CodeHash, &DasApplyRegisterCellScript)
	SystemCodeScriptMap.Store(DasPreAccountCellScript.Out.CodeHash, &DasPreAccountCellScript)
	SystemCodeScriptMap.Store(DasAccountCellScript.Out.CodeHash, &DasAccountCellScript)
	SystemCodeScriptMap.Store(DasBiddingCellScript.Out.CodeHash, &DasBiddingCellScript)
	SystemCodeScriptMap.Store(DasOnSaleCellScript.Out.CodeHash, &DasOnSaleCellScript)
	SystemCodeScriptMap.Store(DasProposeCellScript.Out.CodeHash, &DasProposeCellScript)
	SystemCodeScriptMap.Store(DasConfigCellScript.Out.CodeHash, &DasConfigCellScript)
	SystemCodeScriptMap.Store(DasIncomeCellScript.Out.CodeHash, &DasIncomeCellScript)

}

// testnet version 3
func UseVersion3SystemScriptCodeHash() {
	DasApplyRegisterCellScript.Out.CodeHash = types.HexToHash("0xd8e70cbc0d61daee85b8e121fcb6f278c4536ac26cf9cdce36957a2aa289d4d9")
	DasPreAccountCellScript.Out.CodeHash = types.HexToHash("0x9f7ce0892e4484c058d547b648f969266a43d61d8635bb8460252597bc1a7ecd")
	DasAccountCellScript.Out.CodeHash = types.HexToHash("0xf727c4459d3fbd3f2caf59884a5984f66b3891c69501cecc2959104b3f6f39e0")
	// DasBiddingCellScript.Out.CodeHash = types.HexToHash("0x711bb5cec27b3a5c00da3a6dc0772be8651f7f92fd9bf09d77578b29227c1748")
	// DasOnSaleCellScript.Out.CodeHash = types.HexToHash("0x711bb5cec27b3a5c00da3a6dc0772be8651f7f92fd9bf09d77578b29227c1748")
	DasProposeCellScript.Out.CodeHash = types.HexToHash("0xfa14b6fd2fd690295a22e3ae5c2e189b6c5b8144ec16409492004b1700b37db7")
	DasLockCellScript.Out.CodeHash = types.HexToHash("0x31c4408a02d6d5b9fcd1ca8b542c08755c84a6265e0e0129e0580a4e904d418d")
	DasConfigCellScript.Out.CodeHash = types.HexToHash("0x474fea002daafd29d3aa4143571570f0b8304ab5d4261d9f9ed8135341656321")
	DasIncomeCellScript.Out.CodeHash = types.HexToHash("0x4f2afb853a5a161d8d656b90aa94417b63fa43a2dee19a144b4b8d95b873131c")
	DasAnyOneCanSendCellInfo.Out.CodeHash = types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f")

	initMap()
}

type TimingAsyncSystemCodeScriptParam struct {
	RpcClient     rpc.Client
	SuperLock     *types.Script
	Duration      time.Duration
	Ctx           context.Context
	ErrHandle     func(err error)
	SuccessHandle func()
	InitHandle    func() bool
	FirstSuccessCallBack func()
}

func TimingAsyncSystemCodeScriptOutPoint(p *TimingAsyncSystemCodeScriptParam) {
	if p.SuperLock == nil {
		if p.ErrHandle != nil {
			p.ErrHandle(errors.New("superLock cant be null"))
		}
		return
	}
	isNeedSync := true
	if p.InitHandle != nil {
		isNeedSync = p.InitHandle()
	}
	sync := func(callback bool) {
		liveCells := []indexer.LiveCell{}
		SystemCodeScriptMap.Range(func(key, value interface{}) bool {
			item := value.(*DASCellBaseInfo)
			if item.ContractTypeScript.Args == nil {
				return true
			}
			searchKey := &indexer.SearchKey{
				Script:     p.SuperLock,
				ScriptType: indexer.ScriptTypeLock,
				Filter:     &indexer.CellsFilter{
					Script: &item.ContractTypeScript,
				},
			}
			c := collector.NewLiveCellCollector(p.RpcClient, searchKey, indexer.SearchOrderDesc, 20, "",false)
			iterator, err := c.Iterator()
			if err != nil {
				p.ErrHandle(fmt.Errorf("LoadLiveCells Collect failed: %s", err.Error()))
				return false
			}
			for iterator.HasNext() {
				liveCell, err := iterator.CurrentItem()
				if err != nil {
					p.ErrHandle(fmt.Errorf("LoadLiveCells, read iterator current err: %s", err.Error()))
					return false
				}
				liveCells = append(liveCells,*liveCell)
				if err = iterator.Next(); err != nil {
					p.ErrHandle(fmt.Errorf("LoadLiveCells, read iterator next err: %s", err.Error()))
					return false
				}
			}
			return true
		})
		for _, liveCell := range liveCells {
			scriptCodeOutput := liveCell.Output
			typeId := CalTypeIdFromScript(scriptCodeOutput.Type)
			_ = SetSystemCodeScriptOutPoint(typeId, types.OutPoint{
				TxHash: liveCell.OutPoint.TxHash,
				Index:  liveCell.OutPoint.Index,
			})
		}
		if p.SuccessHandle != nil {
			p.SuccessHandle()
		}
		if callback && p.FirstSuccessCallBack != nil {
			p.FirstSuccessCallBack()
		}
	}
	if isNeedSync {
		sync(true)
	}
	go func() {
		ticker := time.NewTicker(p.Duration)
		defer ticker.Stop()
		if p.Ctx == nil {
			p.Ctx = context.TODO()
		}
		for {
			select {
			case <-p.Ctx.Done():
				return
			case <-ticker.C:
				sync(false)
			}
		}
	}()
}

func SetSystemCodeScriptOutPoint(typeId types.Hash, point types.OutPoint) *DASCellBaseInfo {
	if item, ok := SystemCodeScriptMap.Load(typeId); !ok {
		return nil
	} else {
		obj := item.(*DASCellBaseInfo)
		obj.Dep.TxHash = point.TxHash
		obj.Dep.TxIndex = point.Index
		// SystemCodeScriptMap.Store(typeId,obj)
		return obj
	}
}

func emptyHexToArgsBytes() []byte {
	return []byte{}
}

func dasLockDefaultBytes() []byte {
	return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
}

func hexToArgsBytes(hexStr string) []byte {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	bys, _ := hex.DecodeString(hexStr)
	return bys
}

func IsSystemCodeScriptReady() bool {
	ready := true
	SystemCodeScriptMap.Range(func(key, value interface{}) bool {
		item := value.(*DASCellBaseInfo)
		if item.Out.CodeHash.Hex() == "0x" {
			ready = false
			return false
		}
		return true
	})
	return ready
}

func SystemCodeScriptBytes() ([]byte, error) {
	return json.Marshal(SystemCodeScriptMap)
}

func SystemCodeScriptFromBytes(bys []byte) error {
	if err := json.Unmarshal(bys, &SystemCodeScriptMap); err != nil {
		return err
	}
	return nil
}

//
// func ParseDasCellsScript(data *ConfigCellMain) map[types.Hash]string {
// 	applyRegisterCodeHash := types.BytesToHash(data.TypeIdTable().ApplyRegisterCell().RawData())
// 	preAccountCellCodeHash := types.BytesToHash(data.TypeIdTable().PreAccountCell().RawData())
// 	biddingCellCodeHash := types.BytesToHash(data.TypeIdTable().BiddingCell().RawData())
// 	accountCellCodeHash := types.BytesToHash(data.TypeIdTable().AccountCell().RawData())
// 	proposeCellCodeHash := types.BytesToHash(data.TypeIdTable().ProposalCell().RawData())
// 	onSaleCellCodeHash := types.BytesToHash(data.TypeIdTable().OnSaleCell().RawData())
// 	walletCellCodeHash := types.BytesToHash(data.TypeIdTable().WalletCell().RawData())
// 	refCellCodeHash := types.BytesToHash(data.TypeIdTable().RefCell().RawData())
//
// 	retMap := map[types.Hash]string{}
// 	retMap[applyRegisterCodeHash] = SystemScript_ApplyRegisterCell
// 	retMap[preAccountCellCodeHash] = SystemScript_PreAccoutnCell
// 	retMap[biddingCellCodeHash] = SystemScript_BiddingCell
// 	retMap[accountCellCodeHash] = SystemScript_AccoutnCell
// 	retMap[proposeCellCodeHash] = SystemScript_ProposeCell
// 	retMap[onSaleCellCodeHash] = SystemScript_OnSaleCell
// 	retMap[walletCellCodeHash] = SystemScript_WalletCell
// 	retMap[refCellCodeHash] = SystemScript_RefCell
// 	return retMap
// }
//
// func SetSystemScript(scriptName string, dasCellBaseInfo *DASCellBaseInfo) error {
// 	if v := SystemCodeScriptMap[scriptName]; v != nil {
// 		*SystemCodeScriptMap[scriptName] = *dasCellBaseInfo
// 		return nil
// 	}
// 	return errors.New("unSupport script")
// }
