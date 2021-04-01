package common

import (
	"encoding/binary"
	"unsafe"
)

/**
 * Copyright (C), 2019-2021
 * FileName: bytes
 * Author:   LinGuanHong
 * Date:     2021/1/15 2:48 下午
 * Description:
 */

func Uint8ToBytes(i uint8) []byte {
	return []byte{i}
}

func Int64ToBytes(i int64) []byte {
	return Uint64ToBytes(uint64(i))
}

func Uint64ToBytes(i uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BytesToUint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func BytesToInt64_LittleEndian(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf))
}

func SplitKeyName(key []byte) (string, string) {
	k := string(key)
	length := len(key)
	okString := string(k[1 : length-2])
	ttype := string(k[length-1 : length])
	return okString, ttype
}

func Str2bytes(s string) []byte {
	ptr := (*[2]uintptr)(unsafe.Pointer(&s))
	btr := [3]uintptr{ptr[0], ptr[1], ptr[1]}
	return *(*[]byte)(unsafe.Pointer(&btr))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
