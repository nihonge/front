package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"myproject/contracts"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

func init() {
}

func main() {
	keyGenAddress := common.BytesToAddress([]byte{0x1})
	keyGenContract := contracts.PrecompiledContractsMap[keyGenAddress]
	sk := keyGenContract.Run
	fmt.Println(sk)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("请选择操作：")
		fmt.Println("1. 生成密钥")
		fmt.Println("2. 加密")
		fmt.Println("3. 解密")
		fmt.Println("4. 密态数据计算")
		fmt.Println("5. 退出")
		fmt.Print("输入操作编号: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Print("请输入用户名:")
			username, _ := reader.ReadString('\n')
			// 创建 keyGenerator 实例
			kg := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x4})]

			// 调用 Run 方法
			output, err := kg.Run([]byte(username))
			switch err {
			case fmt.Errorf("用户已注册"):
				fmt.Println("用户已注册")
			case nil:
			default:
				fmt.Println("密钥生成失败")
			}
			//将密钥转化为16进制可读字符串
			hexString := hex.EncodeToString(output)
			fmt.Printf("生成密钥为%s\n", hexString[:40])
		case "2":
			fmt.Print("输入明文: ")
			plaintext, _ := reader.ReadString('\n')
			plaintext = strings.TrimSpace(plaintext)
			fmt.Print("输入密钥: ")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)
		case "3":
			fmt.Print("输入密文: ")
			ciphertext, _ := reader.ReadString('\n')
			ciphertext = strings.TrimSpace(ciphertext)
			fmt.Print("输入密钥: ")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)
		case "4":
		case "5":
			fmt.Println("退出程序")
			// 在程序退出时保存数据
			return
		default:
			fmt.Println("无效的操作编号，请重新输入")
		}
	}
}

func getTwoIntegers(reader *bufio.Reader) (int, int) {
	fmt.Print("输入第一个整数: ")
	input1, _ := reader.ReadString('\n')
	input1 = strings.TrimSpace(input1)
	a, _ := strconv.Atoi(input1)

	fmt.Print("输入第二个整数: ")
	input2, _ := reader.ReadString('\n')
	input2 = strings.TrimSpace(input2)
	b, _ := strconv.Atoi(input2)

	return a, b
}
