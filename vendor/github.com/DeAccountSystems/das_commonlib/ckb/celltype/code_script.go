package celltype

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/collector"
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

	ETHSoScriptDep = DASCellBaseInfoDep{
		TxHash:  types.HexToHash("0x3bffc9beff67d5f93b60b378c68a9910ecc936e5bff0348b3bdf99c4f416213d"),
		TxIndex: 0,
		DepType: types.DepTypeCode,
	}

	TRONSoScriptDep = DASCellBaseInfoDep{
		TxHash:  types.HexToHash("0x9f6b5041638b10e9d53498e0b27db51778274c75efaffddceca93f6ab9e2053c"),
		TxIndex: 0,
		DepType: types.DepTypeCode,
	}

	CKBSoScriptDep = DASCellBaseInfoDep{
		TxHash:  types.HexToHash("0xe08b6487bab378df62d1abe58faebecdfefc5dc4297627c1f7240441db69355b"),
		TxIndex: 0,
		DepType: types.DepTypeCode,
	}

	CKBMultiSoScriptDep = DASCellBaseInfoDep{
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
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x96248cdefb09eed910018a847cfb51ad044c2d7db650112931760e3ef34a7e9a"),
			CodeHashType: types.HashTypeType,
			Args:         hexToArgsBytes("0x02"),
		},
	}
	DasTimeCellScript = DASCellBaseInfo{
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x96248cdefb09eed910018a847cfb51ad044c2d7db650112931760e3ef34a7e9a"),
			CodeHashType: types.HashTypeType,
			Args:         hexToArgsBytes("0x01"),
		},
	}
	DasQuoteCellScript = DASCellBaseInfo{
		Out: DASCellBaseInfoOut{
			CodeHash:     types.HexToHash("0x96248cdefb09eed910018a847cfb51ad044c2d7db650112931760e3ef34a7e9a"),
			CodeHashType: types.HashTypeType,
			Args:         hexToArgsBytes("0x00"),
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

// testnet version 2
func UseVersion2SystemScriptCodeHash() {

	CKBSoScriptDep.TxHash = types.HexToHash("0x209b35208da7d20d882f0871f3979c68c53981bcc4caa71274c035449074d082")
	ETHSoScriptDep.TxHash = types.HexToHash("0xb035c200bf759537d3796edf49b5d6a8ec5f5d78326713f987f31ad24d0b0171")
	TRONSoScriptDep.TxHash = types.HexToHash("0x7dc4ae8fe597045fbd7fe78f2bd26435644a69b755de3824a856f681bacb732b")

	DasApplyRegisterCellScript.ContractTypeScript.Args = types.HexToHash("0xc78fa9066af1624e600ccfb21df9546f900b2afe5d7940d91aefc115653f90d9").Bytes()
	DasPreAccountCellScript.ContractTypeScript.Args = types.HexToHash("0xd3f7ad59632a2ebdc2fe9d41aa69708ed1069b074cd8b297b205f835335d3a6b").Bytes()
	DasAccountCellScript.ContractTypeScript.Args = types.HexToHash("0x6f0b8328b703617508d62d1f017b0d91fab2056de320a7b7faed4c777a356b7b").Bytes()
	DasProposeCellScript.ContractTypeScript.Args = types.HexToHash("0x03d0bb128bd10e666975d9a07c148f6abebe811f511e9574048b30600b065b9a").Bytes()
	DasLockCellScript.ContractTypeScript.Args = types.HexToHash("0xeedd10c7d8fee85c119daf2077fea9cf76b9a92ddca546f1f8e0031682e65aee").Bytes()
	DasConfigCellScript.ContractTypeScript.Args = types.HexToHash("0x34363fad2018db0b3b6919c26870f302da74c3c4ef4456e5665b82c4118eda51").Bytes()
	DasIncomeCellScript.ContractTypeScript.Args = types.HexToHash("0xd7b9d8213671aec03f3a3ab95171e0e79481db2c084586b9ea99914c00ff3716").Bytes()

	DasApplyRegisterCellScript.Out.CodeHash = types.HexToHash("0x0fbff871dd05aee1fda2be38786ad21d52a2765c6025d1ef6927d761d51a3cd1")
	DasPreAccountCellScript.Out.CodeHash = types.HexToHash("0x6c8441233f00741955f65e476721a1a5417997c1e4368801c99c7f617f8b7544")
	DasAccountCellScript.Out.CodeHash = types.HexToHash("0x1106d9eaccde0995a7e07e80dd0ce7509f21752538dfdd1ee2526d24574846b1")
	DasProposeCellScript.Out.CodeHash = types.HexToHash("0x67d48c0911e406518de2116bd91c6af37c05f1db23334ca829d2af3042427e44")
	DasLockCellScript.Out.CodeHash = types.HexToHash("0x326df166e3f0a900a0aee043e31a4dea0f01ea3307e6e235f09d1b4220b75fbd")
	DasConfigCellScript.Out.CodeHash = types.HexToHash("0x030ac2acd9c016f9a4ab13d52c244d23aaea636e0cbd386ec660b79974946517")
	DasIncomeCellScript.Out.CodeHash = types.HexToHash("0x08d1cdc6ab92d9cabe0096a2c7642f73d0ef1b24c94c43f21c6c3a32ffe0bb5e")
	DasAnyOneCanSendCellInfo.Out.CodeHash = types.HexToHash("0xf1ef61b6977508d9ec56fe43399a01e576086a76cf0f7c687d1418335e8c401f")

	initMap()
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

func UseVersionReleaseSystemScriptCodeHash() {

	CKBSoScriptDep.TxHash = types.HexToHash("0x1373db89fd2c7ff1617d4fd6740e916169631c5ab6c9995786645071ab19b822")
	ETHSoScriptDep.TxHash = types.HexToHash("0x43ffc1b114a3a2bb3c94f2c4c55a6e666d1c69d8394af4a035858157aebfc7c4")
	TRONSoScriptDep.TxHash = types.HexToHash("0x17d992cd2c7aaa298c5d3f7c365709f4ac7b25c44aa5d618b6245cfc4a0f0352")

	DasApplyRegisterCellScript.ContractTypeScript.Args = types.HexToHash("0xf18c3eab9fd28adbb793c38be9a59864989c1f739c22d2b6dc3f4284f047a69d").Bytes()
	DasPreAccountCellScript.ContractTypeScript.Args = types.HexToHash("0xf6574955079797010689a22cd172ce55b52d2c34d1e9bc6596e97babc2906f7e").Bytes()
	DasAccountCellScript.ContractTypeScript.Args = types.HexToHash("0x96dc231bbbee6aa474076468640f9e0ad27cf13b1343716a7ce04b116ea18ba8").Bytes()
	DasProposeCellScript.ContractTypeScript.Args = types.HexToHash("0xd7b779b1b30f86a77db6b292c9492906f2437b7d88a8a5994e722619bb1d41c8").Bytes()
	DasLockCellScript.ContractTypeScript.Args = types.HexToHash("0xda22fd296682488687a6035b5fc97c269b72d7de128034389bd03041b40309c0").Bytes()
	DasConfigCellScript.ContractTypeScript.Args = types.HexToHash("0x3775c65aabe8b79980c4933dd2f4347fa5ef03611cef64328685618aa7535794").Bytes()
	DasIncomeCellScript.ContractTypeScript.Args = types.HexToHash("0x108fba6a9b9f2898b4cdf11383ba2a6ed3da951b458c48e5f5de0353bbca2d46").Bytes()

	DasApplyRegisterCellScript.Out.CodeHash = types.HexToHash("0xc024b6efde8d49af665b3245223a8aa889e35ede15bc510392a7fea2dec0a758")
	DasPreAccountCellScript.Out.CodeHash = types.HexToHash("0x18ab87147e8e81000ab1b9f319a5784d4c7b6c98a9cec97d738a5c11f69e7254")
	DasAccountCellScript.Out.CodeHash = types.HexToHash("0x4f170a048198408f4f4d36bdbcddcebe7a0ae85244d3ab08fd40a80cbfc70918")
	DasProposeCellScript.Out.CodeHash = types.HexToHash("0x6127a41ad0549e8574a25b4d87a7414f1e20579306c943c53ffe7d03f3859bbe")
	DasLockCellScript.Out.CodeHash = types.HexToHash("0x9376c3b5811942960a846691e16e477cf43d7c7fa654067c9948dfcd09a32137")
	DasConfigCellScript.Out.CodeHash = types.HexToHash("0x903bff0221b72b2f5d549236b631234b294f10f53e6cc7328af07776e32a6640")
	DasIncomeCellScript.Out.CodeHash = types.HexToHash("0x6c1d69a358954fc471a2ffa82a98aed5a4912e6002a5e761524f2304ab53bf39")

	DasAnyOneCanSendCellInfo.Dep = DASCellBaseInfoDep{
		TxHash:  types.HexToHash("0xcfd3350aa2a5a9277cba3cd784262d206646a10244c9ae924fd39cb4005dd653"),
		TxIndex: 0,
		DepType: types.DepTypeCode,
	}
	DasAnyOneCanSendCellInfo.Out.CodeHash = types.HexToHash("0x303ead37be5eebfcf3504847155538cb623a26f237609df24bd296750c123078")

	DasHeightCellScript.Out.CodeHash = types.HexToHash("0x2e0e5b790cfb346bddc0e82a70f785e90d1537bbfdbdd25f6a3617cc760f887b")
	DasTimeCellScript.Out.CodeHash = types.HexToHash("0x2e0e5b790cfb346bddc0e82a70f785e90d1537bbfdbdd25f6a3617cc760f887b")
	DasQuoteCellScript.Out.CodeHash = types.HexToHash("0x2e0e5b790cfb346bddc0e82a70f785e90d1537bbfdbdd25f6a3617cc760f887b")

	initMap()
}

type TimingAsyncSystemCodeScriptParam struct {
	RpcClient            rpc.Client
	SuperLock            *types.Script
	Duration             time.Duration
	Ctx                  context.Context
	ErrHandle            func(err error)
	SuccessHandle        func()
	InitHandle           func() bool
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
				Filter: &indexer.CellsFilter{
					Script: &item.ContractTypeScript,
				},
			}
			c := collector.NewLiveCellCollector(p.RpcClient, searchKey, indexer.SearchOrderDesc, 20, "", false)
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
				liveCells = append(liveCells, *liveCell)
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
