package blockchain

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nihonge/homo_blockchain/globals"
	"github.com/nihonge/homo_blockchain/homomorphic"
	"github.com/nihonge/homo_blockchain/http"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
)

var (
	contracts = map[string]interface{}{
		"0x21": test,
		"0x22": compute,
	}
	sk     = &rlwe.SecretKey{}
	kgen   = rlwe.NewKeyGenerator(globals.Params)
	client = &ethclient.Client{}
	err    error
)

func test(a string, b string) string {
	return a + b
}

/*
computeType:
1.加 2.减 3.乘
*/
func compute(computeType string, values1 []float64, values2 []float64) error {
	//本地加密成密文
	myenc := &homomorphic.Encryptor{}
	st1 := myenc.Encrypt(sk, values1)
	st2 := myenc.Encrypt(sk, values2)

	// 声明计算器
	rlk := kgen.GenRelinearizationKeyNew(sk)
	evk := rlwe.NewMemEvaluationKeySet(rlk)
	//编码计算器
	evk_byte, _ := evk.MarshalBinary()
	//编码密文
	st1_byte, _ := st1.MarshalBinary()
	st2_byte, _ := st2.MarshalBinary()
	ciphertext_bytes := [][]byte{}

	switch computeType {
	case "1":
		ciphertext_bytes = append(ciphertext_bytes, globals.Addition, evk_byte, st1_byte, st2_byte)
		fmt.Print("期望结果:")
		for i := 0; i < len(values1); i++ {
			fmt.Printf("%20.10f  ", values1[i]+values2[i])
		}
	case "2":
		ciphertext_bytes = append(ciphertext_bytes, globals.Subtraction, evk_byte, st1_byte, st2_byte)
		fmt.Print("期望结果:")
		for i := 0; i < len(values1); i++ {
			fmt.Printf("%20.10f  ", values1[i]-values2[i])
		}
	case "3":
		ciphertext_bytes = append(ciphertext_bytes, globals.Multiplication, evk_byte, st1_byte, st2_byte)
		fmt.Print("期望结果:")
		for i := 0; i < len(values1); i++ {
			fmt.Printf("%20.10f  ", values1[i]*values2[i])
		}
	default:
		return errors.New("不存在该运算")
	}
	data := globals.Encode(ciphertext_bytes...)
	// 预编译合约地址 (compute 地址为 0x22)
	contractAddress := common.HexToAddress("0x22")

	// 准备调用消息
	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}
	// 调用预编译合约
	// ans_byte, err := client.CallContract(context.Background(), msg, nil)
	// if err != nil {
	// 	log.Fatalf("Failed to call contract: %v", err)
	// }

	// 发送数据至服务器
	ret_json, err := http.PostToServer(msg)
	if err != nil {
		return err
	}
	ans_byte, ok := ret_json["data"].([]byte)
	if !ok {
		return errors.New("callContract.go 99:type error")
	}

	// 解密
	mydec := &homomorphic.Decryptor{}
	ans := new(rlwe.Ciphertext)
	ans.UnmarshalBinary(ans_byte)
	decode_ans := mydec.Decrypt(sk, ans)
	// 验证结果
	fmt.Println()
	fmt.Print("实际结果:")
	for i := 0; i < len(values1); i++ {
		fmt.Printf("%20.10f  ", decode_ans[i])
	}
	fmt.Println()
	return nil
}

// 输入要调用的合约地址，以及参数
func Call(addr string, params ...interface{}) error {
	kgen_byte, _ := globals.GetUserPrivateKey(globals.CurrentUser)
	sk.UnmarshalBinary(kgen_byte)
	// 连接到本地 Geth 节点
	client, err = ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	// 检查预编译合约地址是否存在（后期调用业务合约可能不需在客户端检查）
	f, exists := contracts[addr]
	if !exists {
		return errors.New("函数不存在")
	}
	// 通过反射获取函数的反射值和类型
	funcValue := reflect.ValueOf(f)
	funcType := funcValue.Type()
	// 检查参数的数量是否正确
	if len(params) != funcType.NumIn() {
		return errors.New("参数数量不正确")
	}
	// 准备参数
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		if reflect.TypeOf(param) != funcType.In(i) {
			fmt.Println("get:", reflect.TypeOf(param), " want", funcType.In(i))
			return fmt.Errorf("参数 %d 类型不匹配", i+1)
		}
		in[i] = reflect.ValueOf(param)
	}
	// 调用函数
	result := funcValue.Call(in)
	if result[0].IsNil() {
		return nil
	}
	return result[0].Interface().(error)
}

// func GetBlocknumber() {
// 	client, err := ethclient.Dial("http://localhost:8545")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// 获取区块信息
// 	header, err := client.HeaderByNumber(context.Background(), nil)
// 	if err != nil {
// 		log.Fatalf("获取最新区块头失败: %v", err)
// 	}
// 	log.Printf("最新区块号: %v", header.Number.String())
// }
// func GetBanlance() {
// 	client, err := ethclient.Dial("http://localhost:8545")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// 获取账户余额
// 	account := common.HexToAddress("0x0e58872222579801cc0ef86933a456bad552fa42")
// 	balance, err := client.BalanceAt(context.Background(), account, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to retrieve account balance: %v", err)
// 	}

// 	fmt.Printf("Balance: %s\n", balance.String())
// }
