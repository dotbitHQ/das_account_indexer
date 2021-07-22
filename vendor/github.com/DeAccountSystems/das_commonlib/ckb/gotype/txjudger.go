package gotype

import (
	"context"
	"fmt"
	"github.com/DeAccountSystems/das_commonlib/ckb/celltype"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

/**
 * Copyright (C), 2019-2021
 * FileName: txjudger
 * Author:   LinGuanHong
 * Date:     2021/5/12 4:37
 * Description:
 */

func IsEditManagerTx(tx types.Transaction) bool {
	foundAccountCell := false
	for i := 0; i < len(tx.Outputs); i++ {
		output := tx.Outputs[i]
		if output.Type == nil {
			continue
		}
		if celltype.DasAccountCellScript.Out.SameScript(output.Type) {
			foundAccountCell = true
			break
		}
	}
	return foundAccountCell
}

func IsTransferAccountTx(tx types.Transaction) bool  {
	foundAccountCell := false
	for _, output := range tx.Outputs {
		if output.Type == nil {
			continue
		}
		if celltype.DasAccountCellScript.Out.SameScript(output.Type) {
			foundAccountCell = true
			break
		}
	}
	return foundAccountCell
}

func IsEditRecordsTx(tx types.Transaction) bool {
	foundAccountCell := false
	for i := 0; i < len(tx.Outputs); i++ {
		output := tx.Outputs[i]
		if output.Type == nil {
			continue
		}
		if celltype.DasAccountCellScript.Out.SameScript(output.Type) {
			foundAccountCell = true
			break
		} else if dataBytes := tx.OutputsData[i]; len(dataBytes) == 0 {
			continue
		}
	}
	return foundAccountCell
}

func IsRenewAccountTx(tx types.Transaction) bool {
	foundAccountCell := false
	for i := 0; i < len(tx.Outputs); i++ {
		output := tx.Outputs[i]
		if output.Type == nil {
			continue
		}
		if celltype.DasAccountCellScript.Out.SameScript(output.Type) {
			foundAccountCell = true
			break
		} else if dataBytes := tx.OutputsData[i]; len(dataBytes) == 0 {
			continue
		}
	}
	return foundAccountCell
}

func IsConfirmProposeTx(ctx context.Context,rpcClient rpc.Client,inputs []*types.CellInput) error {
	param := &ReqFindTargetTypeScriptParam{
		Ctx:       ctx,
		RpcClient: rpcClient,
		InputList: inputs,
		IsLock:    false,
		CodeHash:  celltype.DasProposeCellScript.Out.CodeHash,
	}
	if _, err := FindTargetTypeScriptByInputList(param); err != nil {
		return fmt.Errorf(
			"FindSenderLockScriptByInputList proposeCell err: %s, "+
				"invalid confirmPropose Tx, first input not a proposeCell. proposeTx: %s", err.Error(),celltype.DasProposeCellScript.Out.CodeHash.String())
	}
	return nil
}

func IsStartAccountAuctionTx(tx types.Transaction) bool {
	var (
		foundAccountCell = false
		foundBiddingCell = false
	)
	for _, output := range tx.Outputs {
		if output.Type == nil {
			continue
		}
		if celltype.DasAccountCellScript.Out.SameScript(output.Type) {
			foundAccountCell = true
		} else if celltype.DasBiddingCellScript.Out.SameScript(output.Type) {
			foundBiddingCell = true
		}
	}
	return foundAccountCell && foundBiddingCell
}

func IsCancelAccountAuctionTx(tx types.Transaction) bool {
	foundAccountCell := false
	for _, output := range tx.Outputs {
		if output.Type == nil {
			continue
		}
		if celltype.DasAccountCellScript.Out.SameScript(output.Type) {
			foundAccountCell = true
			break
		}
	}
	return foundAccountCell
}

func IsStartAccountSaleTx(tx types.Transaction) (bool,int) {
	var (
		foundAccountCell = false
		foundOnSaleCell  = false
		onSaleIndex      = 0
	)
	for index, output := range tx.Outputs {
		if output.Type == nil {
			continue
		}
		if celltype.DasAccountCellScript.Out.SameScript(output.Type) {
			foundAccountCell = true
		} else if celltype.DasOnSaleCellScript.Out.SameScript(output.Type) {
			foundOnSaleCell = true
			onSaleIndex = index
		}
	}
	return foundAccountCell && foundOnSaleCell, onSaleIndex
}

func IsCancelAccountSaleTx(tx types.Transaction) bool {
	foundAccountCell := false
	for _, output := range tx.Outputs {
		if output.Type == nil {
			continue
		}
		if celltype.DasAccountCellScript.Out.SameScript(output.Type) {
			foundAccountCell = true
			break
		}
	}
	return foundAccountCell
}
