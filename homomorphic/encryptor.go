package homomorphic

import (
	"github.com/nihonge/front/globals"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

// 明文加密为同态加密密文，并上链
// 接受密钥和明文作为输入
type Encryptor struct{}

func (c *Encryptor) Encrypt(sk *rlwe.SecretKey, values interface{}) *rlwe.Ciphertext {
	var err error
	// Encryptor
	enc := rlwe.NewEncryptor(globals.Params, sk)
	// Encoder
	//这里直接使用MaxLevel可能导致数据过大
	//可在geth，http.go中修改最大限制
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
	// ct_byte, _ := ct.MarshalBinary()
	// fmt.Println("level=", pt.Level())
	// fmt.Println("密文长度:", len(ct_byte))
	return ct
}
