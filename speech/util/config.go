// Copyright 2017 The Roobo AI Platform Authors. All rights reserved.

// Package util implements some utility functions.
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

type CfgData map[string]*json.RawMessage

type redisConf struct {
	addr, passwd, db string
}

var (
	cfgData       CfgData
	onceSysConf   sync.Once
	onceRedisConf sync.Once
	rConf         redisConf
)

func GetRedisConf() (addr, passwd, db string, err error) {
	onceRedisConf.Do(func() {
		onceSysConf.Do(func() {
			err = initSysConf("./conf/app.json")
			return
		})
		s1, err1 := GetCfgVal("", "redis", "addr")
		s2, err2 := GetCfgVal("", "redis", "passwd")
		s3, err3 := GetCfgVal("5", "redis", "db")
		if err1 != nil || err2 != nil || err3 != nil {
			log.Fatalf("get redis config error: %v; %v; %v", err1, err2, err3)
		}
		rConf = redisConf{addr: s1.(string), passwd: s2.(string), db: s3.(string)}
	})
	if rConf.addr == "" {
		return "", "", "", errors.New("redis unset")
	}
	return rConf.addr, rConf.passwd, rConf.db, nil
}

// NOTE: not thread-safe
func InitSpecConf(cfgFile string) (CfgData, error) {
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read json file: %s, error: %s",
			cfgFile, err))
	}
	raw := make(CfgData)
	err = json.Unmarshal(cfgBytes, &raw)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to Unmarshal file: %s, error: %s",
			cfgFile, err))
	}
	return raw, nil
}

// TODO: load configure file dynamic by start a timer and goroutine
// NOTE: not thread-safe
func initSysConf(cfgFile string) error {
	raw, err := InitSpecConf(cfgFile)
	if err != nil {
		return err
	}
	cfgData = raw
	return nil
}

func GetSpecCfgVal(cfgData CfgData, def interface{}, keys ...string) (
	v interface{}, err error) {
	v = def
	var m interface{}
	m = cfgData
	for i, k := range keys {
		if m, ok := m.(CfgData); ok {
			if d, ok := m[k]; ok {
				if i == len(keys)-1 {
					if err1 := json.Unmarshal(*d, &v); err1 != nil {
						err = errors.New(fmt.Sprintf("failed to Unmarshal config, key: %v,"+
							" error: %s", keys, err1))
					}
					if _, ok := def.(int); ok {
						v = int(v.(float64))
					}
					return
				}
				if err1 := json.Unmarshal(*d, &m); err1 != nil {
					err = errors.New(fmt.Sprintf("failed to Unmarshal config, key: %v,"+
						" error: %s", keys, err1))
					return
				}
			}
		}
	}
	err = errors.New(fmt.Sprintf("GetCfgVal error: invalid Key: %v", keys))
	return
}

func GetCfgVal(def interface{}, keys ...string) (v interface{}, err error) {
	return GetSpecCfgVal(cfgData, def, keys...)
}
