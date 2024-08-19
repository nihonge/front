package contracts

import (
	"bytes"
	"fmt"
	"log"
	"myproject/globals"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
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

type decrypt struct{}

func (c *decrypt) RequiredGas(input []byte) uint64 {
	return 0
}
func (c *decrypt) Run(input []byte) ([]byte, error) {
	return []byte{1, 2, 3}, nil
}

// 计算链上密态数据，目前为求和
type compute struct{}

func (c *compute) RequiredGas(input []byte) uint64 {
	return 0
}
func (c *compute) Run(input []byte) ([]byte, error) {
	// 第一个字节是运算方法，后面是密文序列
	ciphertext_bytes, err := globals.Decode(input)
	if err != nil {
		return []byte{}, fmt.Errorf("%v", err)
	}
	fmt.Println("compute:字节数组个数:", len(ciphertext_bytes))
	if bytes.Equal(ciphertext_bytes[0], globals.Addition) {
		pt1_byte := ciphertext_bytes[1]
		pt2_byte := ciphertext_bytes[2]
		eva := rlwe.NewEvaluator()
		fmt.Println("加法")
	} else if bytes.Equal(ciphertext_bytes[0], globals.Subtraction) {
		fmt.Println("减法")
	} else if bytes.Equal(ciphertext_bytes[0], globals.Multiplication) {
		fmt.Println("乘法")
	} else if bytes.Equal(ciphertext_bytes[0], globals.Division) {
		fmt.Println("除法")
	}
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
