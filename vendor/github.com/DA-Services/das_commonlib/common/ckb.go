package common

import (
	"fmt"
	"github.com/DA-Services/das_commonlib/ckb/collector"
	"github.com/nervosnetwork/ckb-sdk-go/indexer"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
)

/**
 * Copyright (C), 2019-2021
 * FileName: ckb
 * Author:   LinGuanHong
 * Date:     2021/2/25 10:04
 * Description:
 */

func LoadOneLiveCell(client rpc.Client, key *indexer.SearchKey, capLimit uint64, latest, normal bool, filter func(cell *indexer.LiveCell) bool) ([]indexer.LiveCell, uint64, error) {
	return LoadLiveCellsWithSize(client, key, capLimit, 1, latest, normal, filter)
}

func LoadLiveCells(client rpc.Client, key *indexer.SearchKey, capLimit uint64, latest, normal bool, filter func(cell *indexer.LiveCell) bool) ([]indexer.LiveCell, uint64, error) {
	return LoadLiveCellsWithSize(client, key, capLimit, 100, latest, normal, filter)
}

func LoadLiveNormalCells(client rpc.Client, key *indexer.SearchKey, capLimit uint64, filter func(cell *indexer.LiveCell) bool) ([]indexer.LiveCell, uint64, error) {
	return LoadLiveCellsWithSize(client, key, capLimit, 100, false, true, filter)
}

func LoadLiveCellsWithSize(client rpc.Client, key *indexer.SearchKey, capLimit, size uint64, latest, normal bool, filter func(cell *indexer.LiveCell) bool) ([]indexer.LiveCell, uint64, error) {
	order := indexer.SearchOrderAsc
	// note: different args, wont work
	if latest {
		order = indexer.SearchOrderDesc
	}
	c := collector.NewLiveCellCollector(client, key, order, size, "", normal)
	iterator, err := c.Iterator()
	if err != nil {
		return nil, 0, fmt.Errorf("LoadLiveCells Collect failed: %s", err.Error())
	}
	liveCells := []indexer.LiveCell{}
	totalCap := uint64(0)
NextBatch:
	for iterator.HasNext() {
		liveCell, err := iterator.CurrentItem()
		if err != nil {
			return nil, 0, fmt.Errorf("LoadLiveCells, read iterator current err: %s", err.Error())
		}
		if filter != nil && !filter(liveCell) {
			if err = iterator.Next(); err != nil {
				return nil, 0, fmt.Errorf("LoadLiveCells, read iterator next err: %s", err.Error())
			}
			continue
		}
		totalCap = totalCap + liveCell.Output.Capacity
		liveCells = append(liveCells, *liveCell)
		if totalCap >= capLimit { // enough
			break
		}
		if err = iterator.Next(); err != nil {
			return nil, 0, fmt.Errorf("LoadLiveCells, read iterator next err: %s", err.Error())
		}
	}
	if totalCap < capLimit {
		if err = iterator.Next(); err != nil {
			return nil, 0, fmt.Errorf("LoadLiveCells, read iterator next err: %s", err.Error())
		} else if iterator.HasNext() {
			goto NextBatch
		}
	}
	return liveCells, totalCap, nil
}
