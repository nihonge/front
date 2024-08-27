package client_utils

import (
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"myproject/globals"
)

// 明文加密为同态加密密文，并上链
// 接受密钥和明文作为输入
type Encryptor struct{}

func (c *Encryptor) Encrypt(sk *rlwe.SecretKey, values interface{}) *rlwe.Ciphertext {
	var err error
	// Encryptor
	enc := rlwe.NewEncryptor(globals.Params, sk)
	// Encoder
	pt := ckks.NewPlaintext(globals.Params, globals.Params.MaxLevel())
	// Encodes the vector of plaintext values
	ecd := ckks.NewEncoder(globals.Params)
	if err = ecd.Encode(values, pt); err != nil {
		panic(err)
	}
	var ct *rlwe.Ciphertext
	if ct, err = enc.EncryptNew(pt); err != nil {
		panic(err)
	}
	return ct
}
