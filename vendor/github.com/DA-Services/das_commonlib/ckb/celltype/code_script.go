package celltype

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/DA-Services/das_commonlib/common"
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
 * Date:     2020/12/22 3:01 下午
 * Description:
 */


var (
	TestNetLockScriptDep = DASCellBaseInfoDep{
		TxHash:  types.HexToHash("0xf8de3bb47d055cdf460d93a2a6e1b05f7432f9777c8c474abf4eec1d4aee5d37"),
		TxIndex: 0,
		DepType: types.DepTypeDepGroup,
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
	DasWalletCellScript = DASCellBaseInfo{
		Name: "wallet_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0xaac18fd80a6f9265913518e303fe57d1c93d961ef7badbc1289b9dbe667a8ab42"), //"0x440b323f2821aa808c1bad365c10ffb451058864a11f63b5669a5597ac0e8e0f"
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x9878b226df9465c215fd3c94dc9f9bf6648d5bea48a24579cf83274fe13801d2"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(WalletCellCodeArgs).Bytes(),
		},
	}
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
	DasRefCellScript = DASCellBaseInfo{
		Name: "ref_cell",
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0x86a83fc53d64e0cfbc94ccc003b8eee00617c8aa16a2aa1188d41842ee97dc15"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0xe79953f024552e6130220a03d2497dc7c2f784f4297c69ba21d0c423915350e5"),
			CodeHashType: types.HashTypeType,
			Args:         emptyHexToArgsBytes(),
		},
		ContractTypeScript: types.Script{
			CodeHash: types.HexToHash(ContractCodeHash),
			HashType: types.HashTypeType,
			Args:     types.HexToHash(RefCellCodeArgs).Bytes(),
		},
	}
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
			TxHash:  types.HexToHash("0x711bb5cec27b3a5c00da3a6dc0772be8651f7f92fd9bf09d77578b29227c1748"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x5f6a4cc2cd6369dbcf38ddfbc4323cf4695c2e8c20aed572b5db6adc2faf9d50"),
			CodeHashType: types.HashTypeType,
			Args:         hexToArgsBytes("0xe1a958a4c112af95a1220c6fee5f969972a3d8ce13fb7b3211f71abb5db1824102000000"),
		},
	}
	DasTimeCellScript = DASCellBaseInfo{
		Dep: DASCellBaseInfoDep{
			TxHash:  types.HexToHash("0xf3c13ffbaa1d34b8fac6cd848fa04db2e6b4e2c967c3c178295be2e7cdd77164"),
			TxIndex: 0,
			DepType: types.DepTypeCode,
		},
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0xe4fd6f46ab1fd3d5b377df9e2d4ea77e3b52f53ac3319595bb38d097ea051cfd"),
			CodeHashType: types.HashTypeType,
			Args:         hexToArgsBytes("0xd0c1c7156f2e310a12822e2cc336398ec4ef194abc1f96023b743f3249f09e2102000000"),
		},
	}
	SystemCodeScriptMap = syncmap.Map{} // map[types.Hash]*DASCellBaseInfo{}
)

func init() {
	initMap()
}

func initMap()  {
	SystemCodeScriptMap.Store(DasApplyRegisterCellScript.Out.CodeHash,&DasApplyRegisterCellScript)
	SystemCodeScriptMap.Store(DasPreAccountCellScript.Out.CodeHash,&DasPreAccountCellScript)
	SystemCodeScriptMap.Store(DasAccountCellScript.Out.CodeHash,&DasAccountCellScript)
	SystemCodeScriptMap.Store(DasBiddingCellScript.Out.CodeHash,&DasBiddingCellScript)
	SystemCodeScriptMap.Store(DasOnSaleCellScript.Out.CodeHash,&DasOnSaleCellScript)
	SystemCodeScriptMap.Store(DasProposeCellScript.Out.CodeHash,&DasProposeCellScript)
	SystemCodeScriptMap.Store(DasWalletCellScript.Out.CodeHash,&DasWalletCellScript)
	SystemCodeScriptMap.Store(DasRefCellScript.Out.CodeHash,&DasRefCellScript)
}

// testnet version 2
func UseVersion2SystemScriptCodeHash()  {
	DasApplyRegisterCellScript.Out.CodeHash = types.HexToHash("0x0fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd1")
	DasPreAccountCellScript.Out.CodeHash = types.HexToHash("0x6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b7544")
	DasAccountCellScript.Out.CodeHash = types.HexToHash("0x5148d4c832ee9020ef646fb454ee81852d9e28b930eb8c667804e6a51b0a00fc")
	// DasBiddingCellScript.Out.CodeHash = types.HexToHash("0x711bb5cec27b3a5c00da3a6dc0772be8651f7f92fd9bf09d77578b29227c1748")
	// DasOnSaleCellScript.Out.CodeHash = types.HexToHash("0x711bb5cec27b3a5c00da3a6dc0772be8651f7f92fd9bf09d77578b29227c1748")
	DasProposeCellScript.Out.CodeHash = types.HexToHash("0xc432a01b4e0b948e57c6291924914e548a7109028114b97d2815c16d3a06f329")
	DasWalletCellScript.Out.CodeHash = types.HexToHash("0x066a699f5bba9dc4b45bfd7a46f1c5bb1a092dc0eb078810358fad2f07698c37")
	DasRefCellScript.Out.CodeHash = types.HexToHash("0xec5abfd61507cda957d6adc3264ca9bc7120d6db3bf15a50795624e8af54aefa")
	DasConfigCellScript.Out.CodeHash = types.HexToHash("0x79bf0bc0f911c11cb85e51de9ecaf6630ce5bb1cac26ea9c15dd7d08b91c943a")
	DasAnyOneCanSendCellInfo.Out.CodeHash = types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f")
	initMap()
}

func TimingAsyncSystemCodeScriptOutPoint(rpcClient rpc.Client,superLock *types.Script,errHandle func(err error),successHandle func())  {
	sync := func() {
		SystemCodeScriptMap.Range(func(key, value interface{}) bool {
			item := value.(*DASCellBaseInfo)
			if item.ContractTypeScript.Args == nil {
				return true
			}
			searchKey := &indexer.SearchKey{
				Script:     &item.ContractTypeScript,
				ScriptType: indexer.ScriptTypeType,
				Filter: &indexer.CellsFilter{
					Script: superLock,
				},
			}
			liveCells, _, err := common.LoadLiveCells(rpcClient, searchKey, 10000000*OneCkb, true, false, func(cell *indexer.LiveCell) bool {
				return cell.Output.Type != nil
			})
			if err != nil && errHandle != nil {
				errHandle(fmt.Errorf("LoadAllScriptCodeCell err: %s", err.Error()))
				return false
			}
			for _, liveCell := range liveCells {
				scriptCodeOutput := liveCell.Output
				typeId := CalTypeIdFromScript(scriptCodeOutput.Type)
				_ = SetSystemCodeScriptOutPoint(typeId, types.OutPoint{
					TxHash: liveCell.OutPoint.TxHash,
					Index:  liveCell.OutPoint.Index,
				})
			}
			return true
		})
		if successHandle != nil {
			successHandle()
		}
	}
	sync()
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				sync()
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
