package contracts_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/nihonge/homo_blockchain/client_utils"
	"github.com/nihonge/homo_blockchain/globals"
	"github.com/nihonge/homo_blockchain/preCompiledContracts"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
	"github.com/tuneinsight/lattigo/v6/schemes/ckks"
	"testing"
)

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
	var err error
	kgen := ckks.NewKeyGenerator(globals.Params)
	sk := kgen.GenSecretKeyNew()
	// 定义go中明文
	values1 := []float64{1.2345678, 2.0, 7.8}
	values2 := []float64{8.7654321, 5.5, 6.6}

	//本地加密成密文
	myenc := client_utils.Encryptor{}
	st1 := myenc.Encrypt(sk, values1)
	st2 := myenc.Encrypt(sk, values2)

	// 密态数据计算
	// 声明加密预编译合约
	compute := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x1})]
	// 声明计算器
	rlk := kgen.GenRelinearizationKeyNew(sk)
	evk := rlwe.NewMemEvaluationKeySet(rlk)
	//编码计算器
	evk_byte, _ := evk.MarshalBinary()
	//编码密文
	st1_byte, _ := st1.MarshalBinary()
	st2_byte, _ := st2.MarshalBinary()

	//测试加法
	ciphertext_bytes := [][]byte{}
	ciphertext_bytes = append(ciphertext_bytes, globals.Addition, evk_byte, st1_byte, st2_byte)
	ans_byte, err := compute.Run(globals.Encode(ciphertext_bytes...))
	if err != nil {
		t.Errorf("密态数据计算出错:%v", err)
	}
	// 解密
	mydec := client_utils.Decryptor{}
	ans := new(rlwe.Ciphertext)
	ans.UnmarshalBinary(ans_byte)
	decode_ans := mydec.Decrypt(sk, ans)
	if err != nil {
		t.Errorf("密文解密出错:%v", err)
	}
	// 验证结果
	fmt.Print("期望:")
	for i := 0; i < len(values1); i++ {
		fmt.Printf("%20.10f  ", values1[i]+values2[i])
	}
	fmt.Println()
	fmt.Print("实际:")
	for i := 0; i < len(values1); i++ {
		fmt.Printf("%20.10f  ", decode_ans[i])
	}
	fmt.Println()

	//测试减法
	ciphertext_bytes = [][]byte{}
	ciphertext_bytes = append(ciphertext_bytes, globals.Subtraction, evk_byte, st1_byte, st2_byte)
	ans_byte, err = compute.Run(globals.Encode(ciphertext_bytes...))
	if err != nil {
		t.Errorf("密态数据计算出错:%v", err)
	}
	// 解密
	mydec = client_utils.Decryptor{}
	ans = new(rlwe.Ciphertext)
	ans.UnmarshalBinary(ans_byte)
	decode_ans = mydec.Decrypt(sk, ans)
	// 验证结果
	fmt.Print("期望:")
	for i := 0; i < len(values1); i++ {
		fmt.Printf("%20.10f  ", values1[i]-values2[i])
	}
	fmt.Println()
	fmt.Print("实际:")
	for i := 0; i < len(values1); i++ {
		fmt.Printf("%20.10f  ", decode_ans[i])
	}
	fmt.Println()

	//测试乘法
	ciphertext_bytes = [][]byte{}
	ciphertext_bytes = append(ciphertext_bytes, globals.Multiplication, evk_byte, st1_byte, st2_byte)
	ans_byte, err = compute.Run(globals.Encode(ciphertext_bytes...))
	if err != nil {
		t.Errorf("密态数据计算出错:%v", err)
	}
	// 解密
	mydec = client_utils.Decryptor{}
	ans = new(rlwe.Ciphertext)
	ans.UnmarshalBinary(ans_byte)
	decode_ans = mydec.Decrypt(sk, ans)
	if err != nil {
		t.Errorf("密文解密出错:%v", err)
	}
	// 验证结果
	fmt.Print("期望:")
	for i := 0; i < len(values1); i++ {
		fmt.Printf("%20.10f  ", values1[i]*values2[i])
	}
	fmt.Println()
	fmt.Print("实际:")
	for i := 0; i < len(values1); i++ {
		fmt.Printf("%20.10f  ", decode_ans[i])
	}
	fmt.Println()
}

// // TestKeyGenerator_Run tests the Run method of the keyGenerator.
// func TestKeyGenerator_Run(t *testing.T) {
// 	// 初始化测试参数
// 	name := "nihonge5201314"
// 	// 创建 keyGenerator 实例
// 	kg := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x4})]

// 	// 调用 Run 方法
// 	output, err := kg.Run([]byte(name))
// 	if err != nil {
// 		t.Fatalf("Run() error = %v", err)
// 	}

// 	// 验证
// 	secret, err := globals.GetUserKey(name)
// 	if err != nil {
// 		t.Fatalf("用户未注册成功")
// 	}
// 	if string(output) != secret {
// 		t.Fatalf("密码存储不正确")
// 	}
// 	fmt.Printf("用户%s注册成功\n", name)
// 	globals.DeleteUser(name)
// }

// // TestKeyGenerator_RequiredGas tests the RequiredGas method of the keyGenerator.
// func TestKeyGenerator_RequiredGas(t *testing.T) {
// 	kg := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x4})]

// 	// 测试 RequiredGas 方法
// 	input := []byte("test input")
// 	expectedGas := uint64(0) // 预期的Gas费用

// 	actualGas := kg.RequiredGas(input)
// 	if actualGas != expectedGas {
// 		t.Errorf("RequiredGas() = %v, want %v", actualGas, expectedGas)
// 	}
// }
