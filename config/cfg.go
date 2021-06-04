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
	} `json:"log"`
	Rpc struct {
		Port string `json:"port"`
	} `json:"rpc"`
	Chain struct {
		CKB struct {
			NodeUrl    string `json:"node_url"       yaml:"node_url"`
			IndexerUrl string `json:"indexer_url"    yaml:"indexer_url"`
		} `json:"ckb"`
	} `json:"chain"`
}

func (c DasConfig) ToStr() string {
	bys, _ := json.Marshal(c)
	return string(bys)
}

var Cfg DasConfig
