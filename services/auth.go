package services

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"smart/models/mysqls"
	"smart/tools/enums"
	"smart/tools/utils"
	"strings"
)

type Auth struct {
}

// GetAuthInfo 获取系统授权信息
func (a *Auth) GetAuthInfo(ctx context.Context) (map[string]string, error) {
	var mapSet mysqls.MapSet
	mapSetRes, err := mapSet.GetsByObjKey(ctx, enums.ProductAuthInfoMapSetKey)
	if err != nil {
		return nil, err
	}
	var authInfoMap map[string]string
	err = json.Unmarshal([]byte(mapSetRes.ObjValue), &authInfoMap)
	if err != nil {
		return nil, err
	}
	return authInfoMap, nil
}

// GenerateSystemSerialNumber 生成系统序列号（仅 Linux：绑定 DMI / machine-id，不使用跨平台 MAC 降级）
func (a *Auth) GenerateSystemSerialNumber(ctx context.Context) (string, error) {
	if runtime.GOOS != "linux" {
		return "", errors.New("系统授权特征码仅支持在 Linux 环境生成")
	}

	// 1. 优先尝试 dmidecode (硬件级唯一，最稳定，但需要权限)
	productNameCmd := []string{"-s", "system-product-name"}
	cmd := exec.CommandContext(ctx, "dmidecode", productNameCmd...)
	productNameOutput, err1 := cmd.CombinedOutput()

	serialNumberCmd := []string{"-s", "system-serial-number"}
	cmd = exec.CommandContext(ctx, "dmidecode", serialNumberCmd...)
	serialNumberOutput, err2 := cmd.CombinedOutput()

	// 只有两个命令都成功才认为是有效的 DMI 信息
	if err1 == nil && err2 == nil {
		pName := strings.TrimSpace(string(productNameOutput))
		sNum := strings.TrimSpace(string(serialNumberOutput))

		// 定义常见的无效/通用 OEM 标识 (忽略大小写)
		invalidValues := []string{
			"To be filled by O.E.M.",
			"To be filled by O.E.M",
			"Default String",
			"System Product Name",
			"System Serial Number",
			"Not Specified",
			"None",
			"12345678",
			"00000000",
			"",
		}

		isInvalid := false
		for _, v := range invalidValues {
			if strings.EqualFold(pName, v) || strings.EqualFold(sNum, v) {
				isInvalid = true
				break
			}
		}

		if !isInvalid {
			fmt.Println("productNameOutput:", pName)
			fmt.Println("serialNumberOutput:", sNum)
			return utils.Md5V(pName + sNum), nil
		}

		fmt.Printf("dmidecode returned generic/invalid values (Product: %s, Serial: %s), trying fallback...\n", pName, sNum)
	}

	fmt.Printf("dmidecode failed (err1: %v, err2: %v), trying fallback...\n", err1, err2)

	// 2. 降级尝试读取 /etc/machine-id (Linux 系统唯一标识，无需 root 权限通常可读)
	machineId, err := os.ReadFile("/etc/machine-id")
	if err == nil {
		fmt.Println("using /etc/machine-id")
		return utils.Md5V(strings.TrimSpace(string(machineId))), nil
	}

	return "", errors.New("failed to generate system serial number: dmidecode and machine-id unavailable")
}

// 公钥加密
func (a *Auth) RsaEncrypt(data, keyBytes []byte) []byte {
	//解密pem格式的公钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

// RsaDecrypt 解密授权码
func (a *Auth) RsaDecrypt(ctx context.Context, authCode string) (map[string]string, error) {
	//获取私钥
	block, _ := pem.Decode([]byte(enums.SWPrvKey))
	if block == nil {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	authCodeByte, err := hex.DecodeString(authCode)
	if err != nil {
		return nil, err
	}
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, authCodeByte)
	if err != nil {
		return nil, err
	}
	var authMap map[string]string
	err = json.Unmarshal(data, &authMap)
	if err != nil {
		return nil, err
	}
	return authMap, nil
}

// 私钥解密
func (a *Auth) RsaDecryptString(ciphertext, keyBytes []byte) []byte {
	//获取私钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error!"))
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 解密
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		panic(err)
	}
	return data
}

// 验证
func (a *Auth) RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

// 签名
func (a *Auth) RsaSignWithSha256(data []byte, keyBytes []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}
	return signature
}

// UpdateAuthState 更新系统授权状态
func (a *Auth) UpdateAuthState(ctx context.Context, authState string) error {
	var mapSet mysqls.MapSet
	row, err := mapSet.GetsByObjKey(ctx, enums.ProductAuthStateMapSetKey)
	if err != nil {
		return err
	}
	if row.ID == 0 {
		row.Estate = "valid"
		row.ObjKey = enums.ProductAuthStateMapSetKey
		row.ObjValue = authState
		row.Content = "产品授权状态"
		return row.AddMapSet(ctx)
	}
	return mapSet.UpdateObjValueByObjKey(ctx, enums.ProductAuthStateMapSetKey, authState)
}

// UpdateAuthRecord 更新系统授权记录
func (a *Auth) UpdateAuthRecord(ctx context.Context, authRecord string) error {
	var mapSet mysqls.MapSet
	// 1 获取过去的数据并拼接
	mapSetRes, err := mapSet.GetsByObjKey(ctx, enums.ProductAuthRecordMapSetKey)
	if err != nil {
		return err
	}
	if mapSetRes.ID == 0 {
		mapSetRes.Estate = "valid"
		mapSetRes.ObjKey = enums.ProductAuthRecordMapSetKey
		mapSetRes.ObjValue = authRecord
		mapSetRes.Content = "产品授权记录"
		return mapSetRes.AddMapSet(ctx)
	}
	if mapSetRes.ObjValue != "" {
		authRecord = mapSetRes.ObjValue + "," + authRecord
	}
	return mapSet.UpdateObjValueByObjKey(ctx, enums.ProductAuthRecordMapSetKey, authRecord)
}

// CleanAuthRecord 清空授权记录
func (a *Auth) CleanAuthRecord(ctx context.Context) error {
	var mapSet mysqls.MapSet
	row, err := mapSet.GetsByObjKey(ctx, enums.ProductAuthRecordMapSetKey)
	if err != nil {
		return err
	}
	if row.ID == 0 {
		row.Estate = "valid"
		row.ObjKey = enums.ProductAuthRecordMapSetKey
		row.ObjValue = ","
		row.Content = "产品授权记录"
		return row.AddMapSet(ctx)
	}
	return mapSet.UpdateObjValueByObjKey(ctx, enums.ProductAuthRecordMapSetKey, ",")
}

// CheckAuthRecord 查验系统授权记录
func (a *Auth) CheckAuthRecord(ctx context.Context, authRecord string) bool {
	var mapSet mysqls.MapSet
	// 1 获取过去的数据并拼接
	mapSetRes, _ := mapSet.GetsByObjKey(ctx, enums.ProductAuthRecordMapSetKey)
	if mapSetRes.ObjValue == "" {
		return false
	}
	authList := strings.Split(mapSetRes.ObjValue, ",")
	for _, v := range authList {
		if v == authRecord {
			return true
		}
	}
	return false
}

// UpdateAuthInfo 更新授权信息
func (a *Auth) UpdateAuthInfo(ctx context.Context, authInfoMap map[string]string) error {
	var mapSet mysqls.MapSet
	authInfoByte, err := json.Marshal(authInfoMap)
	if err != nil {
		return err
	}
	row, err := mapSet.GetsByObjKey(ctx, enums.ProductAuthInfoMapSetKey)
	if err != nil {
		return err
	}
	if row.ID == 0 {
		row.Estate = "valid"
		row.ObjKey = enums.ProductAuthInfoMapSetKey
		row.ObjValue = string(authInfoByte)
		row.Content = "产品授权信息"
		return row.AddMapSet(ctx)
	}
	return mapSet.UpdateObjValueByObjKey(ctx, enums.ProductAuthInfoMapSetKey, string(authInfoByte))
}

// AddProductID 添加ProductID
func (a *Auth) AddProductID(ctx context.Context, productID string) error {
	var mapSet mysqls.MapSet
	err := mapSet.UpdateObjValueByObjKey(ctx, enums.ProductAuthInfoMapSetKey, productID)
	if err != nil {
		return err
	}
	return nil
}

// GetProductID 获取产品ProductID
func (a *Auth) GetProductID(ctx context.Context) (string, error) {
	var mapSet mysqls.MapSet
	mapSetRes, err := mapSet.GetsByObjKey(ctx, enums.ProductAuthInfoMapSetKey)
	if err != nil {
		return "", err
	}
	var authInfoMap map[string]string
	err = json.Unmarshal([]byte(mapSetRes.ObjValue), &authInfoMap)
	if err != nil {
		return "", err
	}
	if productID, ok := authInfoMap["productID"]; ok && productID != "" {
		return productID, nil
	}
	return "", errors.New("productID not found")
}
