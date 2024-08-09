package globals

import (
	"encoding/json"
	"fmt"
	"github.com/tuneinsight/lattigo/v5/he/hefloat"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

// 声明全局变量

var (
	Params    hefloat.Parameters        //实数同态加密参数
	userStore = make(map[string]string) // 用于存储用户数据的全局变量
	mu        sync.Mutex                // 用于确保并发安全
	fileName  string                    // 数据存储的文件名
)

func init() {
	fmt.Println("全局变量初始化……")
	var err error

	// 获取当前文件的绝对路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Failed to get current file path")
	}
	// 获取项目根目录路径（移除 'globals' 目录）
	projectRoot := filepath.Dir(filepath.Dir(filename))
	// 构建目标文件的路径
	fileName = filepath.Join(projectRoot, "data", "user_keys.json")

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
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		// 文件不存在，不需要加载数据
		fmt.Println("密钥文件数据加载完毕")
		return
	}

	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read user data from file: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()

	err = json.Unmarshal(data, &userStore)
	if err != nil {
		log.Fatalf("Failed to unmarshal user data: %v", err)
	}

	fmt.Println("User data loaded from file successfully")
}

// registerUser 注册新用户
func RegisterUser(username, password string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := userStore[username]; exists {
		return fmt.Errorf("user %s already exists", username)
	}

	userStore[username] = password
	fmt.Printf("User %s registered successfully\n", username)
	saveUserData()
	return nil
}

// saveUserData 将用户数据保存到文件
func saveUserData() {
	// 不再在此处加锁，假设调用此函数时已经加锁
	data, err := json.Marshal(userStore)
	if err != nil {
		log.Fatalf("Failed to marshal user data: %v", err)
	}

	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write user data to file: %v", err)
	}

	fmt.Println("User data saved to file successfully")
}

// deleteUser 删除用户
func DeleteUser(username string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := userStore[username]; !exists {
		return fmt.Errorf("user %s does not exist", username)
	}

	delete(userStore, username)
	fmt.Printf("User %s deleted successfully\n", username)
	saveUserData()
	return nil
}

// getUser 获取用户密码
func GetUser(username string) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	password, exists := userStore[username]
	if !exists {
		return "", fmt.Errorf("user %s does not exist", username)
	}

	return password, nil
}
