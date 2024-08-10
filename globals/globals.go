package globals

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/tuneinsight/lattigo/v5/he/hefloat"
)

// 声明全局变量

var (
	Params            hefloat.Parameters        //实数同态加密参数
	private_key_store = make(map[string]string) // 用于存储用户密钥的全局变量
	address_store     = make(map[string]string) // 用于存储用户地址的全局变量
	mu                sync.Mutex                // 用于确保并发安全
	folderName        string                    // 数据存储的文件夹名 /data
	private_key_file  string
	address_file      string
)

func init() {
	fmt.Println("全局变量初始化……")
	fmt.Println("——·——·——·——·——·——·——")
	var err error

	// 获取当前文件的绝对路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Failed to get current file path")
	}
	// 获取项目根目录路径（移除 'globals' 目录）
	projectRoot := filepath.Dir(filepath.Dir(filename))
	// 构建目标文件的路径
	folderName = filepath.Join(projectRoot, "data")
	private_key_file = filepath.Join(folderName, "user_private_key.json")
	address_file = filepath.Join(folderName, "user_address.json")

	// 128-bit secure parameters enabling depth-7 circuits.
	// LogN:14, LogQP: 431.
	// 加密参数
	if Params, err = hefloat.NewParametersFromLiteral(
		hefloat.ParametersLiteral{
			LogN:            14,                                    // log2(ring degree)
			LogQ:            []int{55, 45, 45, 45, 45, 45, 45, 45}, // log2(primes Q) (ciphertext modulus)
			LogP:            []int{61},                             // log2(primes P) (auxiliary modulus)
			LogDefaultScale: 45,                                    // log2(scale)
		}); err != nil {
		log.Fatalf("Failed to initialize parameters: %v", err)
	}

	// 在程序启动时加载数据
	loadUserData()
}

// loadUserData 从文件加载用户数据
func loadUserData() {
	if _, err := os.Stat(private_key_file); os.IsNotExist(err) {
		// 文件不存在，不需要加载数据
		fmt.Println("密钥文件数据加载完毕")
		fmt.Println("——·——·——·——·——·——·——")
		return
	}

	data, err := os.ReadFile(private_key_file)
	if err != nil {
		log.Fatalf("Failed to read user data from file: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()

	err = json.Unmarshal(data, &private_key_store)
	if err != nil {
		log.Fatalf("Failed to unmarshal user data: %v", err)
	}

	if _, err := os.Stat(address_file); os.IsNotExist(err) {
		// 文件不存在，不需要加载数据
		fmt.Println("地址文件数据加载完毕")
		return
	}

	data, err = os.ReadFile(address_file)
	if err != nil {
		log.Fatalf("Failed to read user data from file: %v", err)
	}

	err = json.Unmarshal(data, &address_store)
	if err != nil {
		log.Fatalf("Failed to unmarshal user data: %v", err)
	}

	fmt.Println("User data loaded from file successfully")
}

// registerUser 注册新用户
func RegisterUser(username, password string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := private_key_store[username]; exists {
		return fmt.Errorf("user %s already exists", username)
	}

	private_key_store[username] = password
	address_store[username] = keyToAddr(password)
	fmt.Printf("User %s registered successfully\n", username)
	saveUserData()
	return nil
}

// saveUserData 将用户数据保存到文件
func saveUserData() {
	// 不再在此处加锁，假设调用此函数时已经加锁
	data, err := json.Marshal(private_key_store)
	if err != nil {
		log.Fatalf("Failed to marshal user data: %v", err)
	}
	err = os.WriteFile(private_key_file, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write user data to file: %v", err)
	}

	data, err = json.Marshal(address_store)
	if err != nil {
		log.Fatalf("Failed to marshal user data: %v", err)
	}
	err = os.WriteFile(address_file, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write user data to file: %v", err)
	}

	fmt.Println("User data saved to file successfully")
}

// deleteUser 删除用户
func DeleteUser(username string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := private_key_store[username]; !exists {
		return fmt.Errorf("user %s does not exist", username)
	}

	delete(private_key_store, username)
	delete(address_store, username)
	fmt.Printf("User %s deleted successfully\n", username)
	saveUserData()
	return nil
}

// getUser 获取用户密码
func GetUserKey(username string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	password, exists := private_key_store[username]
	if !exists {
		return "", fmt.Errorf("user %s does not exist", username)
	}

	return password, nil
}
func GetUserAddr(username string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	password, exists := address_store[username]
	if !exists {
		return "", fmt.Errorf("user %s does not exist", username)
	}

	return password, nil
}

// 展示用户名和地址
func ShowUser() {
	fmt.Printf("\n共%d人 ", len(address_store))
	fmt.Printf("%-10s %-10s\n", "用户名称", "地址")

	for key, val := range address_store {
		fmt.Printf("%-10s %-10s\n", key, val)
	}
	fmt.Println()
}

// 密钥转地址 → 计算 SHA-256 哈希值的字符串
func keyToAddr(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashBytes := hash.Sum(nil)
	// 将哈希值转换为十六进制字符串
	return hex.EncodeToString(hashBytes)
}
