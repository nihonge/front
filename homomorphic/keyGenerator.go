package homomorphic

import (
	"fmt"
	"github.com/nihonge/front/globals"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"log"
)

// keyGenerator负责产生同态加密密钥，输入用户名（任意），会安全地产生一个密钥（与用户名无关），用来加密数据并上传
// 密钥很长，生成之后只能先保存在本地，若更换客户端位置登录则需要导出密钥
type KeyGenerator struct{}

func (c *KeyGenerator) GenerateKey(input string) error {
	//检查用户名是否注册过
	_, err := globals.GetUserPrivateKey(input)
	if err == nil {
		return fmt.Errorf("用户已注册")
	}

	kgen := rlwe.NewKeyGenerator(globals.Params)
	sk := kgen.GenSecretKeyNew()
	// pk := kgen.GenPublicKeyNew(sk)

	skCode, err := sk.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to serialize secret key: %v", err)
	}
	globals.RegisterUser(input, skCode)

	// fmt.Printf("Private Key: %v\n", sk)
	// fmt.Printf("Public Key: %v\n", pk)
	return nil
}
