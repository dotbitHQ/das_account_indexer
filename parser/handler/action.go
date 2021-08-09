package handler

import (
	"das_account_indexer/types"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/eager7/elog"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/tecbot/gorocksdb"
	"sync"
)

var (
	once = sync.Once{}
	log  = elog.NewLogger("parse_handler", elog.NoticeLevel)
)

type DASActionHandleFuncParam struct {
	Base      *types.ParserHandleBaseTxInfo
	RpcClient rpc.Client
	Rocksdb   *gorocksdb.DB
}
type DASActionHandleFuncResp struct {
	Tag      string
	Data     interface{}
	Rollback bool
	err      error
}

func NewDASActionHandleFuncResp(funcName string) DASActionHandleFuncResp {
	return DASActionHandleFuncResp{Tag: funcName}
}
func (resp DASActionHandleFuncResp) SetErr(err error) DASActionHandleFuncResp {
	resp.err = err
	return resp
}

func (resp DASActionHandleFuncResp) Error() error {
	if resp.err == nil {
		return nil
	}
	return fmt.Errorf("funcName:[%s]--> %s", resp.Tag, resp.err.Error())
}

// ==============================================

type DASActionHandleFunc func(actionName string, p *DASActionHandleFuncParam) DASActionHandleFuncResp

type ActionRegister struct {
	handlerMap map[string]DASActionHandleFunc
}

func NewActionRegister() *ActionRegister {
	register := ActionRegister{
		handlerMap: make(map[string]DASActionHandleFunc),
	}
	register.RegisterTxActionHandler()
	return &register
}

func (a *ActionRegister) RegisterTxActionHandler() {
	once.Do(func() {
		a.handlerMap[celltype.Action_ConfirmProposal] = HandleConfirmProposalTx
		a.handlerMap[celltype.Action_EditManager] = HandleEditManagerTx
		a.handlerMap[celltype.Action_EditRecords] = HandleEditRecordsTx
		a.handlerMap[celltype.Action_TransferAccount] = HandleTransferAccountTx
		a.handlerMap[celltype.Action_RecycleExpiredAccount] = HandleExpiredRecycleAccountTx
	})
}

func (a *ActionRegister) GetTxActionHandleFunc(actionName string) DASActionHandleFunc {
	return a.handlerMap[actionName]
}
