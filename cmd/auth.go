package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"smart/services"
	"smart/tools/enums"
	"strconv"
	"time"
)

func main() {
	var (
		days         string
		serialNumber string
		generate     bool
	)
	flag.StringVar(&days, "d", "31", "auth days")
	flag.StringVar(&serialNumber, "s", "", "the system serial number")
	flag.BoolVar(&generate, "g", false, "the system serial number")
	flag.Parse()
	if serialNumber == "" && generate == false {
		fmt.Println("请输入序列号,或者选择生成序列号")
		os.Exit(1)
	}
	if generate {
		generateSerialNumber()
	} else {
		generateAuthCode(serialNumber, days)
	}
}

func generateAuthCode(serialNumber string, days string) string {
	_, err := strconv.Atoi(days)
	if err != nil {
		fmt.Println("请输入正确的天数")
		os.Exit(1)
	}
	tempData, err := buildAuthData(serialNumber, days)
	if err != nil {
		panic(err)
	}
	var auth services.Auth
	ciphertext := auth.RsaEncrypt(tempData, []byte(enums.SWPubKey))
	fmt.Println("公钥加密后的数据：", hex.EncodeToString(ciphertext))
	sourceData := auth.RsaDecryptString(ciphertext, []byte(enums.SWPrvKey))
	fmt.Println("私钥解密后的数据：", string(sourceData))
	return hex.EncodeToString(ciphertext)
}

func generateSerialNumber() {
	ctx := context.Background()
	var auth services.Auth
	//err := config.NewConfig("../config.json") //加载配置文件
	//if err != nil {
	//	return
	//}
	//mysql.Setup()
	//productID, _ := auth.GetProductID(ctx)
	//fmt.Println(productID)
	serialNumber, err := auth.GenerateSystemSerialNumber(ctx)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println(serialNumber)
	}
}

func buildAuthData(serialNumber string, days string) ([]byte, error) {
	tempMap := make(map[string]string, 0)
	tempMap["productID"] = serialNumber
	tempMap["generateTime"] = time.Now().Format("2006-01-02")
	tempMap["authDays"] = days
	tempData, err := json.Marshal(tempMap)
	if err != nil {
		return nil, nil
	}
	return tempData, nil
}

//RSA公钥私钥产生
func GenRsaKey() (prvkey, pubkey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	prvkey = pem.EncodeToMemory(block)
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkey = pem.EncodeToMemory(block)
	return
}

////签名
//func RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
//	h := sha256.New()
//	h.Write(data)
//	hashed := h.Sum(nil)
//	block, _ := pem.Decode(keyBytes)
//	if block == nil {
//		panic(errors.New("private key error"))
//	}
//	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		fmt.Println("ParsePKCS8PrivateKey err", err)
//		panic(err)
//	}
//
//	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
//	if err != nil {
//		fmt.Printf("Error from signing: %s\n", err)
//		panic(err)
//	}
//
//	return signature
//}

////验证
//func RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
//	block, _ := pem.Decode(keyBytes)
//	if block == nil {
//		panic(errors.New("public key error"))
//	}
//	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		panic(err)
//	}
//
//	hashed := sha256.Sum256(data)
//	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
//	if err != nil {
//		panic(err)
//	}
//	return true
//}

//
//// 公钥加密
//func RsaEncrypt(data, keyBytes []byte) []byte {
//	//解密pem格式的公钥
//	block, _ := pem.Decode(keyBytes)
//	if block == nil {
//		panic(errors.New("public key error"))
//	}
//	// 解析公钥
//	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
//	if err != nil {
//		panic(err)
//	}
//	// 类型断言
//	pub := pubInterface.(*rsa.PublicKey)
//	//加密
//	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
//	if err != nil {
//		panic(err)
//	}
//	return ciphertext
//}
//
//// 私钥解密
//func RsaDecrypt(ciphertext, keyBytes []byte) []byte {
//	//获取私钥
//	block, _ := pem.Decode(keyBytes)
//	if block == nil {
//		panic(errors.New("private key error!"))
//	}
//	//解析PKCS1格式的私钥
//	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
//	if err != nil {
//		panic(err)
//	}
//	// 解密
//	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
//	if err != nil {
//		panic(err)
//	}
//	return data
//}
