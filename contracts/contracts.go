package contracts

import (
	"fmt"
	"log"
	"myproject/globals"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
)

type PrecompiledContract interface {
	RequiredGas(input []byte) uint64  // RequiredPrice calculates the contract gas use
	Run(input []byte) ([]byte, error) // Run runs the precompiled contract
}

var PrecompiledContractsMap = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{0x1}): &encrypt{},
	common.BytesToAddress([]byte{0x2}): &decrypt{},
	common.BytesToAddress([]byte{0x3}): &compute{},
	common.BytesToAddress([]byte{0x4}): &keyGenerator{},
}

// 明文加密为同态加密密文
type encrypt struct{}

func (c *encrypt) RequiredGas(input []byte) uint64 {
	//gas需要参考现有的预编译合约，目前只能根据计算量大致定义
	return 0
}
func (c *encrypt) Run(input []byte) ([]byte, error) {
	//使用lattigo的算法
	return []byte{1, 2, 3}, nil
}

type decrypt struct{}

func (c *decrypt) RequiredGas(input []byte) uint64 {
	return 0
}
func (c *decrypt) Run(input []byte) ([]byte, error) {
	return []byte{1, 2, 3}, nil
}

type compute struct{}

func (c *compute) RequiredGas(input []byte) uint64 {
	return 0
}
func (c *compute) Run(input []byte) ([]byte, error) {
	return []byte{1, 2, 3}, nil
}

type keyGenerator struct{}

func (c *keyGenerator) RequiredGas(input []byte) uint64 {
	return 0
}
func (c *keyGenerator) Run(input []byte) ([]byte, error) {
	//检查用户名是否注册过
	_, err := globals.GetUserKey(string(input))
	if err == nil {
		return []byte{}, fmt.Errorf("用户已注册")
	}

	kgen := rlwe.NewKeyGenerator(globals.Params)
	sk := kgen.GenSecretKeyNew()
	// pk := kgen.GenPublicKeyNew(sk)

	skCode, err := sk.MarshalBinary()
	if err != nil {
		log.Fatalf("Failed to serialize secret key: %v", err)
	}
	globals.RegisterUser(string(input), string(skCode))

	// fmt.Printf("Private Key: %v\n", sk)
	// fmt.Printf("Public Key: %v\n", pk)
	return skCode, nil
}
