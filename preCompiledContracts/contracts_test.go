package contracts_test

import (
	"bytes"
	"fmt"
	"github.com/tuneinsight/lattigo/v5/core/rlwe"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"myproject/globals"
	"myproject/preCompiledContracts"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestKeyGenerator_Run tests the Run method of the keyGenerator.
func TestKeyGenerator_Run(t *testing.T) {
	// 初始化测试参数
	name := "nihonge5201314"
	// 创建 keyGenerator 实例
	kg := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x4})]

	// 调用 Run 方法
	output, err := kg.Run([]byte(name))
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	// 验证
	secret, err := globals.GetUserKey(name)
	if err != nil {
		t.Fatalf("用户未注册成功")
	}
	if string(output) != secret {
		t.Fatalf("密码存储不正确")
	}
	fmt.Printf("用户%s注册成功\n", name)
	globals.DeleteUser(name)
}

// TestKeyGenerator_RequiredGas tests the RequiredGas method of the keyGenerator.
func TestKeyGenerator_RequiredGas(t *testing.T) {
	kg := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x4})]

	// 测试 RequiredGas 方法
	input := []byte("test input")
	expectedGas := uint64(0) // 预期的Gas费用

	actualGas := kg.RequiredGas(input)
	if actualGas != expectedGas {
		t.Errorf("RequiredGas() = %v, want %v", actualGas, expectedGas)
	}
}

func TestComputeRequiredGas(t *testing.T) {
	e := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x1})]

	// 测试输入
	input := []byte{0x01, 0x02, 0x03}

	// 预期的Gas值
	expectedGas := uint64(0)

	// 调用函数
	gas := e.RequiredGas(input)

	// 验证结果
	if gas != expectedGas {
		t.Errorf("RequiredGas() = %d; want %d", gas, expectedGas)
	}
}

// 测试加密，密态数据计算和解密后是否正确
func TestComputeRun(t *testing.T) {
	kgen := rlwe.NewKeyGenerator(globals.Params)
	sk := kgen.GenSecretKeyNew()
	ecd := hefloat.NewEncoder(globals.Params) // 用于把go中切片类型进行编码转换
	sk_bytes, err := sk.MarshalBinary()       //将密钥转化为字节流方便调用预编译合约
	if err != nil {
		t.Errorf("Failed to serialize secret key: %v", err)
	}
	// 定义go中明文
	values := make([]float64, 2)
	values[0] = 100
	values[1] = 100
	pt := hefloat.NewPlaintext(globals.Params, 2) //初始化明文
	// Encodes the vector of plaintext values
	if err = ecd.Encode(values, pt); err != nil {
		t.Errorf("Failed to encode values to plaintext: %v", err)
	}
	pt_byte, err := pt.MarshalBinary()
	if err != nil {
		t.Errorf("明文序列化失败: %v", err)
	}

	encrypt := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x1})]
	decrypt := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x2})]
	compute := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x3})]

	// 加密
	ciphertext, err := encrypt.Run(globals.Encode(sk_bytes, pt_byte))
	if err != nil {
		t.Errorf("加密出现错误:%v", err)
	}

	// 密态数据计算

	// 解密

	// 验证结果
	if !bytes.Equal(output, expectedOutput) {
		t.Errorf("Run() = %v; want %v", output, expectedOutput)
	}
}
