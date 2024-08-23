package client_utils

import (
	"fmt"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"myproject/globals"
)

// 明文加密为同态加密密文，并上链
// 接受密钥和明文作为输入
type encrypt struct{}

func (c *encrypt) RequiredGas(input []byte) uint64 {
	//gas需要参考现有的预编译合约，目前只能根据计算量大致定义
	return 0
}
func (c *encrypt) Run(input []byte) ([]byte, error) {
	decode, err := globals.Decode(input)
	if err != nil {
		return []byte{}, fmt.Errorf("解码错误:%v", err)
	}
	fmt.Println("encrypt:字节数组个数:", len(decode))
	var sk rlwe.SecretKey
	sk.UnmarshalBinary(decode[0])
	// Encryptor
	var pt rlwe.Plaintext
	enc := rlwe.NewEncryptor(globals.Params, &sk)
	var ct *rlwe.Ciphertext //密文
	var encrypted_data [][]byte
	for i := 1; i < len(decode); i++ {
		pt.UnmarshalBinary(decode[i])
		if ct, err = enc.EncryptNew(&pt); err != nil {
			return []byte{}, fmt.Errorf("加密错误:%v", err)
		}
		ct_byte, err := ct.MarshalBinary()
		if err != nil {
			return []byte{}, fmt.Errorf("密文序列化错误:%v", err)
		}
		encrypted_data = append(encrypted_data, ct_byte)
	}
	return globals.Encode(encrypted_data...), nil
}
