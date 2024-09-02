package homomorphic

import (
	"github.com/nihonge/homo_blockchain/globals"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
)

// 明文加密为同态加密密文，并上链
// 接受密钥和明文作为输入
type Decryptor struct{}

func (c *Decryptor) Decrypt(sk *rlwe.SecretKey, ct *rlwe.Ciphertext) []float64 {
	var err error
	// decryptor
	dec := rlwe.NewDecryptor(globals.Params, sk)
	pt := dec.DecryptNew(ct)
	ecd := ckks.NewEncoder(globals.Params)
	decodePt := make([]float64, globals.Params.MaxSlots())
	if err = ecd.Decode(pt, decodePt); err != nil {
		panic(err)
	}
	return decodePt
}
