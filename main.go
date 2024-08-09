package main

import (
	"bufio"
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
		case "2":
			fmt.Print("输入明文: ")
			plaintext, _ := reader.ReadString('\n')
			plaintext = strings.TrimSpace(plaintext)
			fmt.Print("输入密钥: ")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)
			ciphertext, err := encrypt.Encrypt(plaintext, key)
			if err != nil {
				fmt.Println("加密错误:", err)
			} else {
				fmt.Println("加密结果:", ciphertext)
			}
		case "3":
			fmt.Print("输入密文: ")
			ciphertext, _ := reader.ReadString('\n')
			ciphertext = strings.TrimSpace(ciphertext)
			fmt.Print("输入密钥: ")
			key, _ := reader.ReadString('\n')
			key = strings.TrimSpace(key)
			plaintext, err := encrypt.Decrypt(ciphertext, key)
			if err != nil {
				fmt.Println("解密错误:", err)
			} else {
				fmt.Println("解密结果:", plaintext)
			}
		case "4":
			a, b := getTwoIntegers(reader)
			result := compute.Subtract(a, b)
			fmt.Printf("结果: %d - %d = %d\n", a, b, result)
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
