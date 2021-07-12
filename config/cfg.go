package config

import (
	"encoding/json"
)

/**
 * Copyright (C), 2019-2019
 * FileName: bn_config
 * Author:   LinGuanHong
 * Date:     2019-10-23 17:55
 * Description: yaml
 */

type DasConfig struct {
	Server struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Debug   bool   `json:"debug"`
	} `json:"server"`
	Log struct {
		Detailed bool `json:"detailed"`
	} `json:"log"`
	Rpc struct {
		Port string `json:"port"`
	} `json:"rpc"`
	Chain struct {
		CKB struct {
			NodeUrl      string `json:"node_url"       yaml:"node_url"`
			IndexerUrl   string `json:"indexer_url"    yaml:"indexer_url"`
			LocalStorage struct {
				ParseDelayMs   int64  `json:"parse_delay_ms" yaml:"parse_delay_ms"`
				InfoDbDataPath string `json:"info_db_data_path" yaml:"info_db_data_path"`
				BlockParser    struct {
					RocksDB struct {
						DataPath string `json:"data_path" yaml:"data_path"`
					} `json:"rocks_db" yaml:"rocks_db"`
					FrontNumber int `json:"front_number" yaml:"front_number"`
					ApiReqRetry struct {
						Times   int   `json:"times"`
						DelayMs int64 `json:"delay_ms" yaml:"delay_ms"`
					} `json:"api_req_retry" yaml:"api_req_retry"`
					ScanControl struct {
						RoundDelayMs int64 `json:"round_delay_ms" yaml:"round_delay_ms"`
						CatchDelayMs int64 `json:"catch_delay_ms" yaml:"catch_delay_ms"`
					} `json:"scan_control" yaml:"scan_control"`
					StartHeight int64 `json:"start_height" yaml:"start_height"`
				} `json:"block_parser" yaml:"block_parser"`
			} `json:"local_storage" yaml:"local_storage"`
		} `json:"ckb"`
	} `json:"chain"`
}

func (c DasConfig) ToStr() string {
	bys, _ := json.Marshal(c)
	return string(bys)
}

var Cfg DasConfig
