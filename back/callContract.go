package back

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nihonge/front/globals"
	"github.com/nihonge/front/homomorphic"
	"github.com/nihonge/front/http"
	"github.com/tuneinsight/lattigo/v6/core/rlwe"
)

var (
	sk   = &rlwe.SecretKey{}
	kgen = rlwe.NewKeyGenerator(globals.Params)
)

const (
	UploadDataMethod = "uploadData"
	GetDataMethod    = "getData"
	ComputeMethod    = "compute"
	BusinessContract = "businessContract"
)

// 定义结构体
type RequestData struct {
	MethodName string `json:"method_name"`
	Params     []byte `json:"params"`
}

func Test(a string, b string) string {
	return a + b
}

/*
computeType:
1.加 2.减 3.乘
*/
func Compute(computeType string, values1 []float64, values2 []float64) error {
	kgen_byte, _ := globals.GetUserPrivateKey(globals.CurrentUser)
	sk.UnmarshalBinary(kgen_byte)
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
	fmt.Println()
	data := globals.Encode(ciphertext_bytes...)

	// 准备调用消息
	request_data := RequestData{
		MethodName: ComputeMethod,
		Params:     data,
	}
	// 发送数据至服务器
	response, err := http.PostToServer(request_data)
	if err != nil {
		return err
	}
	// 获取 Base64 编码的数据
	encodedData, ok := response["data"].(string)
	if !ok {
		return errors.New("callContract.go 99:type error")
	}
	// 解码 Base64 数据
	ans_byte, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return err
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

func UploadData(key string, data []byte) error {
	// 编码 uploadData 调用
	var args struct {
		Key  string `json:"key"`
		Data []byte `json:"data"`
	}
	args.Key = "1"
	args.Data = []byte("exampleData")
	args_byte, err := json.Marshal(args)
	if err != nil {
		return err
	}
	request_data := RequestData{
		MethodName: UploadDataMethod,
		Params:     args_byte,
	}
	response, err := http.PostToServer(request_data)
	if err != nil {
		return err
	}
	// 获取 Base64 编码的数据
	if response["status"] == "success" {
		fmt.Println("上传加密数据成功！")
	} else {
		return errors.New("上传加密数据失败！")
	}
	return nil
}
func GetData(key string) error {
	// 编码 uploadData 调用
	var args struct {
		Key string `json:"key"`
	}
	args.Key = "nihonge"
	args_byte, err := json.Marshal(args)
	if err != nil {
		return err
	}
	request_data := RequestData{
		MethodName: GetDataMethod,
		Params:     args_byte,
	}
	response, err := http.PostToServer(request_data)
	if err != nil {
		return err
	}
	encodedData, ok := response["data"].(string)
	if !ok {
		return errors.New("callContract.go 99:type error")
	}
	// 解码 Base64 数据
	ans_byte, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return err
	}
	fmt.Println("获取的数据为:", string(ans_byte))
	return nil
}
