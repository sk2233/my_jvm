/*
@author: sk
@date: 2024/12/29
*/
package main

import (
	"encoding/binary"
	"io"
)

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadAll(reader io.ReadCloser) []byte {
	bs, err := io.ReadAll(reader)
	HandleErr(err)
	err = reader.Close()
	HandleErr(err)
	return bs
}

func ParseU8(bs []byte, index int) uint8 {
	return bs[index]
}

func ParseU16(bs []byte, index int) uint16 {
	return binary.BigEndian.Uint16(bs[index : index+2])
}

func ParseI16(bs []byte, index int) int16 {
	return int16(ParseU16(bs, index))
}
