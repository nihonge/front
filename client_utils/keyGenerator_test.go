package client_utils

import (
	"bytes"
	"fmt"
	"github.com/nihonge/homo_blockchain/globals"
	"testing"
)

func TestKeyGenerator(t *testing.T) {
	for i := 1; i < 100; i++ {
		k := &keyGenerator{}
		sk, err := k.GenerateKey("nihonge")
		if err != nil {
			t.Errorf("错误:%v\n", err)
		}
		compressedSk, _ := Compress(sk)
		fmt.Printf("数据: 原始大小 %d bytes, 压缩后大小 %d bytes\n", len(sk), len(compressedSk))
		decompressedSk, err := Decompress(compressedSk)
		if err != nil {
			t.Errorf("解压缩错误:%v\n", err)
		}
		fmt.Printf("解压缩后大小 %d bytes\n", len(decompressedSk))
		if bytes.Equal(sk, decompressedSk) {
			fmt.Println("解压缩后相等")
		} else {
			t.Errorf("压缩出错\n")
		}
		globals.DeleteUser("nihonge")
	}
}
