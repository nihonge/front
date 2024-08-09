package contracts_test

import (
	"fmt"
	"myproject/contracts"
	"myproject/globals"
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
	fmt.Println("before")
	output, err := kg.Run([]byte(name))
	fmt.Println("after")
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	// 验证
	secret, err := globals.GetUser(name)
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
