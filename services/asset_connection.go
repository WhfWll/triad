package services

import (
	"context"
	"encoding/base64"
	"smart/models/mysqls"
	aesEncryption "smart/tools/encryption"
	"time"
)

type AssetConnectionService struct {
	aesEcb aesEncryption.AesEcb
}

// CreateAssetConnection 新增连接（自动加密密码）
func (svc *AssetConnectionService) CreateAssetConnection(ctx context.Context, assetID, port, protocol int, ip, username, password string) error {
	assetConMysql := mysqls.AssetConnection{
		AssetID:    assetID,
		IP:         ip,
		Port:       port,
		Protocol:   protocol,
		Username:   username,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
	if len(password) > 0 {
		encrypted := svc.aesEcb.AesEncryptECB([]byte(password), aesKey)
		assetConMysql.Password = base64.StdEncoding.EncodeToString(encrypted)
	}
	return assetConMysql.CreateAssetConnection(ctx)
}

// UpdateAssetConnection 修改连接（自动加密密码）
func (svc *AssetConnectionService) UpdateAssetConnection(ctx context.Context, ip string, assetID, port, protocol int, username, password string) error {
	assetConMysql := mysqls.AssetConnection{
		AssetID:  assetID,
		IP:       ip,
		Port:     port,
		Protocol: protocol,
		Username: username,
	}
	if len(password) > 0 {
		encrypted := svc.aesEcb.AesEncryptECB([]byte(password), aesKey)
		assetConMysql.Password = base64.StdEncoding.EncodeToString(encrypted)
	}
	return assetConMysql.UpsertByIP(ctx)
}

// NOTE:一个资产可能有多个连接方式（SSH、RDP、Telnet、MySQL 等）
// GetConnectionsByIP 通过ip获取某资产的所有连接信息
func (svc *AssetConnectionService) GetConnectionsByIP(ctx context.Context, ip string) ([]mysqls.AssetConnection, error) {
	var assetConn mysqls.AssetConnection
	return assetConn.GetByAssetIP(ctx, ip)
}

// GetConnectionsByProtocol 通过连接方式获取连接信息，返回 map[ip]连接信息切片
func (svc *AssetConnectionService) GetConnectionsByProtocol(ctx context.Context, protocol int) (map[string][]mysqls.AssetConnection, error) {
	var assetConn mysqls.AssetConnection
	list, err := assetConn.GetByProtocol(ctx, protocol)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]mysqls.AssetConnection)
	for _, conn := range list {
		result[conn.IP] = append(result[conn.IP], conn)
	}
	return result, nil
}

// GetAllConnections 按 IP 分组返回所有连接信息（不区分协议）
func (svc *AssetConnectionService) GetAllConnections(ctx context.Context) (map[string][]mysqls.AssetConnection, error) {
	var assetConn mysqls.AssetConnection
	list, err := assetConn.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]mysqls.AssetConnection)
	for _, conn := range list {
		result[conn.IP] = append(result[conn.IP], conn)
	}
	return result, nil
}

// DeleteAssetConnection 删除连接信息
func (svc *AssetConnectionService) DeleteAssetConnection(ctx context.Context, id int) error {
	var conn mysqls.AssetConnection
	return conn.DeleteAssetConnection(ctx, id)
}

// GetConnectionsByAssetID 获取某资产的所有连接信息
func (svc *AssetConnectionService) GetConnectionsByAssetID(ctx context.Context, assetID int) ([]mysqls.AssetConnection, error) {
	var conn mysqls.AssetConnection
	return conn.GetByAssetID(ctx, assetID)
}

// GetConnectionByID 获取单个连接信息
func (svc *AssetConnectionService) GetConnectionByID(ctx context.Context, id int) (*mysqls.AssetConnection, error) {
	var conn mysqls.AssetConnection
	return conn.GetByID(ctx, id)
}

func (svc *AssetConnectionService) DecryptPassword(encryptedBase64 string) string {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return encryptedBase64
	}
	decryptedBytes := svc.aesEcb.AesDecryptECB(encryptedBytes, aesKey)
	return string(decryptedBytes)
}

// NOTE:一个资产可能有多个连接方式（SSH、RDP、Telnet、MySQL 等）当前默认先按照一种处理
// GetConnectionsByIPWithDecryptedPassword 通过ip检索连接信息 返回解密后的密码信息
func (svc *AssetConnectionService) GetConnectionsByIPWithDecryptedPassword(ctx context.Context, ip string) ([]mysqls.AssetConnection, error) {
	var assetConn mysqls.AssetConnection
	assetConnList, err := assetConn.GetByAssetIP(ctx, ip)
	if err != nil {
		return nil, err
	}
	for i := range assetConnList {
		if len(assetConnList[i].Password) > 0 {
			decoded, _ := base64.StdEncoding.DecodeString(assetConnList[i].Password)
			decrypted := svc.aesEcb.AesDecryptECB(decoded, aesKey)
			assetConnList[i].Password = string(decrypted)
		}
	}
	return assetConnList, nil
}

// GetAssetConnectionsList 获取资产连接列表
func (svc *AssetConnectionService) GetAssetConnectionsList(ctx context.Context, ip string, port, protocol, page, size int) ([]mysqls.AssetConnection, int64, error) {
	var assetConn mysqls.AssetConnection
	assetConnList, count, err := assetConn.GetAssetConnList(ctx, ip, port, protocol, page, size)
	if err != nil {
		return nil, 0, err
	}
	for i := range assetConnList {
		if len(assetConnList[i].Password) > 0 {
			decoded, _ := base64.StdEncoding.DecodeString(assetConnList[i].Password)
			decrypted := svc.aesEcb.AesDecryptECB(decoded, aesKey)
			assetConnList[i].Password = string(decrypted)
		}
	}
	return assetConnList, count, nil
}
