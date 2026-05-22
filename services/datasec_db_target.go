package services

import (
	"encoding/base64"

	aesEncryption "smart/tools/encryption"
)

type DatasecDBTargetCrypto struct {
	aesEcb aesEncryption.AesEcb
}

func NewDatasecDBTargetCrypto() *DatasecDBTargetCrypto {
	return &DatasecDBTargetCrypto{}
}

func (c *DatasecDBTargetCrypto) EncryptPassword(password string) string {
	if password == "" {
		return ""
	}
	encrypted := c.aesEcb.AesEncryptECB([]byte(password), aesKey)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func (c *DatasecDBTargetCrypto) DecryptPassword(encryptedBase64 string) string {
	if encryptedBase64 == "" {
		return ""
	}
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return ""
	}
	decryptedBytes := c.aesEcb.AesDecryptECB(encryptedBytes, aesKey)
	return string(decryptedBytes)
}
