package main

import (
	"bufio"
	"fmt"
	"myproject/globals"
	"myproject/preCompiledContracts"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const (
	Guest = iota // 0 - 未登录游客
	User         // 1 - 已登录用户
	Admin        // 2 - 管理员
)

func init() {
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	state := Guest
	for {
		switch state {
		case Guest:
			fmt.Println("请选择操作：")
			fmt.Println("1. 注册生成密钥")
			fmt.Println("2. 登录")
			fmt.Println("3. 退出")
			fmt.Print("输入操作编号: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				fmt.Print("请输入新注册用户名:")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				// 创建 keyGenerator 实例
				kg := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x4})]

				// 调用 Run 方法
				_, err := kg.Run([]byte(username))
				switch err {
				case fmt.Errorf("用户已注册"):
					fmt.Println("用户已注册")
				case nil:
				default:
					fmt.Println("密钥生成失败")
				}
				addr, _ := globals.GetUserAddr(username)
				fmt.Printf("生成密钥为%s\n", addr)
			case "2":
				fmt.Print("请输入用户名:")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				fmt.Print("请输入密码:")
				password, _ := reader.ReadString('\n')
				password = strings.TrimSpace(password)
				addr, err := globals.GetUserAddr(username)
				if err != nil {
					fmt.Println("用户名不存在!")
					continue
				}
				if addr != password {
					fmt.Println("密码错误!")
					continue
				}
				//管理员登录
				if username == "1" {
					const (
						Reset   = "\033[0m"
						Red     = "\033[31m"
						Green   = "\033[32m"
						Yellow  = "\033[33m"
						Blue    = "\033[34m"
						Magenta = "\033[35m"
						Cyan    = "\033[36m"
						White   = "\033[37m"
					)
					fmt.Println(Red + "  _          _ _               _ _                            " + Reset)
					fmt.Println(Green + " | |        | | |             (_) |                           " + Reset)
					fmt.Println(Yellow + " | |__   ___| | | ___    _ __  _| |__   ___  _ __   __ _  ___ " + Reset)
					fmt.Println(Blue + " | '_ \\ / _ \\ | |/ _ \\  | '_ \\| | '_ \\ / _ \\| '_ \\ / _` |/ _ \\" + Reset)
					fmt.Println(Magenta + " | | | |  __/ | | (_) | | | | | | | | | (_) | | | | (_| |  __/" + Reset)
					fmt.Println(Cyan + " |_| |_|\\___|_|_|\\___/  |_| |_|_|_| |_|\\___/|_| |_|\\__, |\\___|" + Reset)
					fmt.Println(White + "                                                    __/ |     " + Reset)
					fmt.Println(Red + "                                                   |___/      " + Reset)
					state = Admin
				} else {
					fmt.Println("登录成功!")
					state = User
				}
			case "3":
				fmt.Println("退出程序")
				return
			default:
				fmt.Println("无效的操作编号，请重新输入")
			}
		case User:
			fmt.Println("请选择操作：")
			fmt.Println("1. 加密")
			fmt.Println("2. 解密")
			fmt.Println("3. 密态数据计算")
			fmt.Println("4. 退出登录")
			fmt.Print("输入操作编号: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				fmt.Print("请输入用户名:")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				// 创建 keyGenerator 实例
				kg := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x4})]

				// 调用 Run 方法
				_, err := kg.Run([]byte(username))
				switch err {
				case fmt.Errorf("用户已注册"):
					fmt.Println("用户已注册")
				case nil:
				default:
					fmt.Println("密钥生成失败")
				}
				addr, _ := globals.GetUserAddr(username)
				fmt.Printf("生成密钥为%s\n", addr)
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
				state = Guest
			case "5":
				globals.ShowUser()
			case "6":
				fmt.Println("退出程序")
				// 在程序退出时保存数据
				return
			default:
				fmt.Println("无效的操作编号，请重新输入")
			}
		case Admin:
			fmt.Println("请选择操作：")
			fmt.Println("1. 加密")
			fmt.Println("2. 解密")
			fmt.Println("3. 密态数据计算")
			fmt.Println("4. 查看用户列表")
			fmt.Println("5. 删除用户")
			fmt.Println("6. 退出登录")
			fmt.Print("输入操作编号: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				fmt.Print("请输入用户名:")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				// 创建 keyGenerator 实例
				kg := contracts.PrecompiledContractsMap[common.BytesToAddress([]byte{0x4})]

				// 调用 Run 方法
				_, err := kg.Run([]byte(username))
				switch err {
				case fmt.Errorf("用户已注册"):
					fmt.Println("用户已注册")
				case nil:
				default:
					fmt.Println("密钥生成失败")
				}
				addr, _ := globals.GetUserAddr(username)
				fmt.Printf("生成密钥为%s\n", addr)
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
				globals.ShowUser()
			case "5":
				fmt.Print("请输入要删除的用户名:")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				globals.DeleteUser(username)
			case "6":
				state = Guest
			default:
				fmt.Println("无效的操作编号，请重新输入")
			}
		default:
		}
	}
}
