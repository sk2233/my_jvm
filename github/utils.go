/*
@author: sk
@date: 2024/12/28
*/
package main

import (
	"encoding/binary"
	"io"
	"os"
)

func OpenFile(path string) *os.File {
	file, err := os.Open(path)
	HandleErr(err)
	return file
}

func ReadU32(reader io.Reader) uint32 {
	bs := ReadBytes(reader, 4)
	return binary.BigEndian.Uint32(bs)
}

func ReadU16(reader io.Reader) uint16 {
	bs := ReadBytes(reader, 2)
	return binary.BigEndian.Uint16(bs)
}

func ReadU8(reader io.Reader) uint8 {
	return ReadBytes(reader, 1)[0]
}

func ReadBytes(reader io.Reader, count int) []byte {
	bs := make([]byte, count)
	_, err := reader.Read(bs)
	HandleErr(err)
	return bs
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}
