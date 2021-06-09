package cfg

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"os"
)

/**
 * Copyright (C), 2019-2020
 * FileName: config
 * Author:   LinGuanHong
 * Date:     2020/12/21 11:26
 * Description:
 */

func InitCfgFromFile(filepath string, receiver interface{}) {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		panic(fmt.Sprintf("open config file failed: %s", err.Error()))
	}
	bys, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Sprintf("read config file failed: %s", err.Error()))
	}
	if err = yaml.Unmarshal(bys, receiver); err != nil {
		panic(fmt.Sprintf("yaml config file formal invalid: %s", err.Error()))
	}
}
