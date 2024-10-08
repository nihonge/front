package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nihonge/front/back"
	"github.com/nihonge/front/globals"
	"github.com/nihonge/front/homomorphic"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	Guest = iota // 0 - 未登录游客
	User         // 1 - 已登录用户
	// Admin        // 2 - 管理员
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
			fmt.Println("3. 通过密钥文件登录（请将密钥文件放在当前文件夹下）")
			fmt.Println("4. 退出")
			fmt.Print("输入操作编号: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			switch input {
			case "1":
				fmt.Print("请输入新注册用户名:")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				// 创建 keyGenerator 实例
				kg := &homomorphic.KeyGenerator{}
				err := kg.GenerateKey(username)
				if err != nil {
					fmt.Printf("%v\n", err)
				} else {
					globals.CurrentUser = username
					state = User
					addr, _ := globals.GetUserAddr(username)
					data, err := json.Marshal(addr)
					if err != nil {
						log.Fatalf("Failed to marshal user data: %v", err)
					}
					fmt.Println("生成密钥为", string(data))
				}
			case "2":
				fmt.Print("请输入用户名:")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				fmt.Print("请输入密码:")
				password, _ := reader.ReadString('\n')
				password = strings.TrimSpace(password)
				addr_byte, err := globals.GetUserAddr(username)
				if err != nil {
					fmt.Println("用户名不存在!")
					continue
				}
				data, err := json.Marshal(addr_byte)
				if err != nil {
					log.Fatalf("Failed to marshal user data: %v", err)
				}
				password = fmt.Sprintf("\"%s\"", password)
				if string(data) != password {
					fmt.Println("密码错误!")
					continue
				}
				fmt.Println("登录成功")
				globals.CurrentUser = username
				state = User
				// //管理员登录
				// if username == "1" {
				// 	const (
				// 		Reset   = "\033[0m"
				// 		Red     = "\033[31m"
				// 		Green   = "\033[32m"
				// 		Yellow  = "\033[33m"
				// 		Blue    = "\033[34m"
				// 		Magenta = "\033[35m"
				// 		Cyan    = "\033[36m"
				// 		White   = "\033[37m"
				// 	)
				// 	fmt.Println(Red + "  _          _ _               _ _                            " + Reset)
				// 	fmt.Println(Green + " | |        | | |             (_) |                           " + Reset)
				// 	fmt.Println(Yellow + " | |__   ___| | | ___    _ __  _| |__   ___  _ __   __ _  ___ " + Reset)
				// 	fmt.Println(Blue + " | '_ \\ / _ \\ | |/ _ \\  | '_ \\| | '_ \\ / _ \\| '_ \\ / _` |/ _ \\" + Reset)
				// 	fmt.Println(Magenta + " | | | |  __/ | | (_) | | | | | | | | | (_) | | | | (_| |  __/" + Reset)
				// 	fmt.Println(Cyan + " |_| |_|\\___|_|_|\\___/  |_| |_|_|_| |_|\\___/|_| |_|\\__, |\\___|" + Reset)
				// 	fmt.Println(White + "                                                    __/ |     " + Reset)
				// 	fmt.Println(Red + "                                                   |___/      " + Reset)
				// 	state = Admin
				// } else {
				// 	fmt.Println("登录成功!")
				// 	state = User
				// }
			case "3":
				fmt.Print("请输入用户名:")
				username, _ := reader.ReadString('\n')
				username = strings.TrimSpace(username)
				key_name := fmt.Sprintf("%v_key.txt", username)
				key, err := os.ReadFile(key_name)
				if err != nil {
					fmt.Println("密钥导入错误:", err)
				}
				globals.RegisterUser(username, key)
				fmt.Println("密钥导入成功！")
				state = User
				globals.CurrentUser = username
			case "4":
				fmt.Println("退出程序")
				return
			default:
				fmt.Println("无效的操作编号，请重新输入")
			}
		case User:
			fmt.Println("请选择操作：")
			fmt.Println("1. 密态数据计算")
			fmt.Println("2. 上传数据")
			fmt.Println("3. 获取数据")
			fmt.Println("4. 导出密钥")
			fmt.Println("5. 退出登录")
			fmt.Print("输入操作编号: ")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				fmt.Println("请选择运算种类：")
				fmt.Println("1. 加法")
				fmt.Println("2. 减法")
				fmt.Println("3. 乘法")
				fmt.Print("输入操作编号：")
				computeType, _ := reader.ReadString('\n')
				computeType = strings.TrimSpace(computeType)
				if computeType != "1" && computeType != "2" && computeType != "3" {
					fmt.Println("无效的操作编号，请重新输入")
					continue
				}
				// 从命令行读取第一个向量
				fmt.Println("请输入第一个向量（以空格分隔的数字）：")
				vector1 := readVector()
				// 从命令行读取第二个向量
				fmt.Println("请输入第二个向量（以空格分隔的数字）：")
				vector2 := readVector()
				// 检查向量长度是否相等
				if len(vector1) != len(vector2) {
					fmt.Println("错误：向量长度不相等")
					continue
				}
				err := back.Compute(computeType, vector1, vector2)
				if err != nil {
					log.Println("调用合约出错:", err)
				}
			case "2":
				err := back.UploadData("1", []byte("helloworld"))
				if err != nil {
					log.Println(err)
				}
			case "3":
				err := back.GetData("nihonge")
				if err != nil {
					log.Println(err)
				}
			case "4":
				data, _ := globals.GetUserPrivateKey(globals.CurrentUser)
				// 将字节数组写入文件
				key_name := fmt.Sprintf("%v_key.txt", globals.CurrentUser)
				err := os.WriteFile(key_name, data, 0644)
				if err != nil {
					log.Println("密钥导出错误:", err)
				}
				fmt.Printf("密钥文件导出成功:%v\n", key_name)
			case "5":
				// 在程序退出时保存数据
				state = Guest
			default:
				fmt.Println("无效的操作编号，请重新输入")
			}
		// case Admin:
		// 	fmt.Println("请选择操作：")
		// 	fmt.Println("1. 查看用户列表")
		// 	fmt.Println("2. 删除用户")
		// 	fmt.Println("3. 退出登录")
		// 	fmt.Print("输入操作编号: ")

		// 	input, _ := reader.ReadString('\n')
		// 	input = strings.TrimSpace(input)

		// 	switch input {
		// 	case "1":
		// 		globals.ShowUser()
		// 	case "2":
		// 		fmt.Print("请输入要删除的用户名:")
		// 		username, _ := reader.ReadString('\n')
		// 		username = strings.TrimSpace(username)
		// 		globals.DeleteUser(username)
		// 	case "3":
		// 		state = Guest
		// 	default:
		// 		fmt.Println("无效的操作编号，请重新输入")
		// 	}
		default:
		}
	}
}

// 读取输入的向量
func readVector() []float64 {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// 分割输入的字符串
	stringValues := strings.Fields(input)
	vector := make([]float64, len(stringValues))

	// 将字符串转换为浮点数
	for i, str := range stringValues {
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			fmt.Println("错误：无法解析输入的数字", str)
			os.Exit(1)
		}
		vector[i] = value
	}

	return vector
}
