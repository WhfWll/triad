package services

import (
	"encoding/base64"
	"gitlabee.4dogs.cn/common/config"
	"smart/api/typespec"
	"smart/tools/encryption"
	"sort"
)

type Common struct {
}

func toolsSort(data map[int]string, defaultIndex ...int) (res []typespec.GlobalOptionsItemRes) {
	sortKeys := make([]int, 0)
	for value, _ := range data {
		sortKeys = append(sortKeys, value)
	}

	sort.Ints(sortKeys)

	for _, key := range sortKeys {
		isDefault := false
		// 是否为默认
		if len(defaultIndex) == 1 && defaultIndex[0] == key {
			isDefault = true
		}
		res = append(res, typespec.GlobalOptionsItemRes{
			Value:     key,
			Label:     data[key],
			IsDefault: isDefault,
		})
	}
	return
}

// CbcPasswordEncryption AES cbc模式的加密
func CbcPasswordEncryption(rawPassword string) string {
	var (
		aesCbc  encryption.AesCbc
		certVal string
	)
	if rawPassword == "" {
		return rawPassword
	}
	err := config.Load("aesCbcUserCert", &certVal)
	if err != nil {
		return rawPassword
	}
	if certVal == "" {
		return rawPassword
	}
	encryptedPass := aesCbc.AesEncryptCBC([]byte(rawPassword), []byte(certVal))
	return base64.StdEncoding.EncodeToString(encryptedPass)
}

// CbcPasswordDecodeString AES cbc模式的解密
func CbcPasswordDecodeString(encryptedPassBase64 string) string {
	var (
		aesCbc  encryption.AesCbc
		certVal string
	)

	if encryptedPassBase64 == "" {
		return ""
	}
	if err := config.Load("aesCbcUserCert", &certVal); err != nil {
		return ""
	}
	if certVal == "" {
		return ""
	}
	encryptedPass, err := base64.StdEncoding.DecodeString(encryptedPassBase64)
	if err != nil {
		return ""
	}

	return string(aesCbc.AesDecryptCBC(encryptedPass, []byte(certVal)))
}
