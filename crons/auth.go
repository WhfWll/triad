package crons

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"
	"os/exec"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type AuthProvider interface {
	GetAuthInfo(ctx context.Context) (map[string]string, error)
	GenerateSystemSerialNumber(ctx context.Context) (string, error)
	RsaDecrypt(ctx context.Context, authCode string) (map[string]string, error)
	UpdateAuthState(ctx context.Context, authState string) error
	UpdateAuthInfo(ctx context.Context, authInfoMap map[string]string) error
	RsaEncrypt(data, keyBytes []byte) []byte
	CleanAuthRecord(ctx context.Context) error
}

var authProvider AuthProvider = &services.Auth{}

// CheckSystemAuth 检查系统授权
func CheckSystemAuth() {
	auth := authProvider
	log.Println("begin CheckSystemAuth")
	ctx := context.Background()
	//第一步 获取系统授权信息
	authInfoMap, err := auth.GetAuthInfo(ctx)
	if err != nil {
		serialNumber, serr := auth.GenerateSystemSerialNumber(ctx)
		if serr != nil {
			log.Println("CheckSystemAuth GetAuthInfo Err:" + err.Error())
			auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
			return
		}
		def := map[string]string{
			"productID":   serialNumber,
			"productName": "自动化渗透测试系统",
			"authCode":    "",
			"authTime":    "未授权",
			"authExpTime": "",
			"authDays":    "0天",
			"leftDays":    "0",
		}
		_ = auth.UpdateAuthInfo(ctx, def)
		auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
		return
	}
	//第二步 比较产品序列号，判断是否被改动过
	productID := authInfoMap["productID"]
	serialNumber, err := auth.GenerateSystemSerialNumber(ctx)
	if err != nil {
		log.Println("CheckSystemAuth GenerateSystemSerialNumber Err:" + err.Error())
		auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
		return
	}
	if productID != serialNumber {
		log.Println("CheckSystemAuth Err: 产品序列号不一致 serialNumber is " + serialNumber + " productID is " + productID)

		// 兼容旧版本序列号逻辑：
		// 如果不一致，尝试用“旧算法”（带空格换行）再算一次。
		// 如果旧算法算出的结果 == 数据库里的 productID，说明这是升级上来的老用户。
		// 此时不应该废除授权，而应该静默更新数据库里的 productID 为新格式，并保持授权状态。

		// 1. 获取原始 dmidecode 输出（为了重现旧算法，这里需要重新调用一次命令，或者修改 GenerateSystemSerialNumber 返回更多信息）
		// 由于 GenerateSystemSerialNumber 已经修改，我们无法直接复用。
		// 这里我们做一个简化的假设：如果数据库里的 productID 对应的是 "旧算法" 的结果，
		// 那么我们现在生成的 serialNumber (新算法) 应该是 productID 去除格式后的版本？
		// 不完全是，因为 md5(a+b) != md5(trim(a)+trim(b))。
		// 所以我们需要手动模拟旧算法：即不做 TrimSpace，而是 Trim(..., " ")。

		ctx := context.Background()
		oldStyleSerialNumber, err := generateOldSystemSerialNumber(ctx)
		if err == nil && oldStyleSerialNumber == productID {
			log.Println("检测到旧版本序列号格式，尝试自动迁移到新格式...")

			// 自动迁移逻辑：
			// 1. 获取当前的 authCode
			authCode := authInfoMap["authCode"]
			if authCode != "" {
				// 2. 解密 authCode 获取原始授权信息
				decryptMap, decryptErr := auth.RsaDecrypt(ctx, authCode)
				if decryptErr == nil {
					// 3. 构造新的 payload，使用新序列号 (serialNumber) 替换旧 productID
					// 保持 authDays 和 generateTime 不变，确保授权有效期延续
					tempMap := map[string]string{
						"productID":    serialNumber, // 使用新生成的干净序列号
						"generateTime": decryptMap["generateTime"],
						"authDays":     decryptMap["authDays"],
					}

					// 4. 重新加密生成新的 authCode
					tempData, _ := json.Marshal(tempMap)
					ciphertext := auth.RsaEncrypt(tempData, []byte(enums.SWPubKey))
					newAuthCode := hex.EncodeToString(ciphertext)

					// 5. 更新内存中的 authInfoMap，以便后续流程使用新数据
					authInfoMap["productID"] = serialNumber
					authInfoMap["authCode"] = newAuthCode

					// 6. 更新数据库 (UpdateAuthInfo 会在后面统一调用，或者这里提前调用)
					// 注意：CheckSystemAuth 的最后一步会调用 UpdateAuthInfo 保存 authInfoMap
					// 所以我们只需要修改 authInfoMap 即可。
					// 但为了保险起见，同时也为了避免后续逻辑用旧数据判断，我们这里显式更新一下 productID 变量
					productID = serialNumber // 更新本地变量，欺骗后续逻辑认为 ID 是一致的

					log.Println("旧版本序列号自动迁移成功：内存已更新，新ID:", serialNumber)
				} else {
					log.Println("自动迁移失败：无法解密旧 authCode，保持兼容模式。Err:", decryptErr)
				}
			} else {
				log.Println("自动迁移跳过：无有效 authCode，保持兼容模式。")
			}
		} else {
			// 只有当既不等于新序列号，也不等于旧序列号时，才认为是真正的硬件变更

			// 序列号变更时，自动初始化
			tempMap := map[string]string{
				"productID":    serialNumber,
				"generateTime": time.Now().Format("2006-01-02"),
				"authDays":     "0",
			}
			tempData, _ := json.Marshal(tempMap)
			ciphertext := auth.RsaEncrypt(tempData, []byte(enums.SWPubKey))
			_ = auth.UpdateAuthInfo(ctx, map[string]string{
				"productID":   serialNumber,
				"productName": "自动化渗透测试系统",
				"authCode":    hex.EncodeToString(ciphertext),
				"authTime":    "未授权",
				"authExpTime": "",
				"authDays":    "0天",
				"leftDays":    "0",
			})
			_ = auth.CleanAuthRecord(ctx)
			auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
			return
		}
	}
	//第三步 解码授权码，并比较 产品id和授权码中的产品id
	authCode := authInfoMap["authCode"]
	if authCode == "" {
		auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
		return
	}
	decryptMap, err := auth.RsaDecrypt(ctx, authCode)
	if err != nil {
		log.Println("CheckSystemAuth RsaDecrypt Err: " + err.Error())
		auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
		return
	}

	// 兼容逻辑：
	// 数据库里的 productID 可能是旧格式，也可能是新格式。
	// decryptMap["productID"] 是当初生成授权码时绑定的 ID。
	// 正常情况下，decryptMap["productID"] == productID (数据库里的)。
	// 如果它们相等，说明授权码是匹配这个“登记在册”的 ID 的。
	// 至于这个 ID 是新的还是旧的，在上面“第二步”已经验证过它属于当前机器（或者是当前机器的旧格式）了。

	if decryptMap["productID"] != productID {
		log.Println("CheckSystemAuth Err: 授权信息有误，授权信息中产品ID与产品ID不相同 授权信息序列号ID：" + decryptMap["productID"] + " 产品ID： " + productID)
		auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
		return
	}
	//第四步 计算授权时间
	authTimeString := authInfoMap["authTime"]
	authDays, err := strconv.Atoi(decryptMap["authDays"])
	if err != nil {
		log.Println("CheckSystemAuth decryptMap Atoi Err: " + err.Error())
		auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
		return
	}
	authTime, err := time.Parse(enums.ResTimeDayLayout, authTimeString)
	if err != nil {
		// 尝试从 decryptMap 中恢复 generateTime
		// 如果 authTime 是 "未授权" 或其他非法格式，但 decryptMap 是有效的，说明可以信任 decryptMap 里的时间
		if genTime, ok := decryptMap["generateTime"]; ok {
			parsedTime, parseErr := time.Parse(enums.ResTimeDayLayout, genTime)
			if parseErr == nil {
				log.Println("CheckSystemAuth: 自动修复 authTime, 使用 generateTime:", genTime)
				authTime = parsedTime
				authInfoMap["authTime"] = genTime // 更新内存中的值，以便后续保存
				err = nil
			}
		}
	}

	if err != nil {
		log.Println("CheckSystemAuth ParseTime Err: " + err.Error())
		auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
		return
	}
	authExpireTime := authTime.AddDate(0, 0, authDays)

	//第五步 更新授权信息
	authInfoMap["authExpTime"] = authExpireTime.Format(enums.ResTimeDayLayout)
	authInfoMap["authDays"] = decryptMap["authDays"] + "天"
	leftDays := int(math.Ceil(authExpireTime.Sub(time.Now()).Hours() / 24))
	if leftDays > 0 {
		authInfoMap["leftDays"] = strconv.Itoa(leftDays)
	} else {
		authInfoMap["leftDays"] = "-- "
	}
	err = auth.UpdateAuthInfo(ctx, authInfoMap)

	//第六步 判断授权时间
	if time.Now().After(authExpireTime) { // 如果当前时间超过了允许的授权时间，那么将不允许使用
		log.Println("CheckSystemAuth Err: 产品已过期")
		auth.UpdateAuthState(ctx, enums.ProductAuthStateFailed)
		return
	}

	auth.UpdateAuthState(ctx, enums.ProductAuthStateSuccess)

	if err != nil {
		log.Println("CheckSystemAuth Err: 更新授权信息失败")
	}
	log.Println("授权信息更新完成")
}

// generateOldSystemSerialNumber 模拟旧版本的序列号生成逻辑，用于兼容性检查
func generateOldSystemSerialNumber(ctx context.Context) (string, error) {
	// 旧版本逻辑：dmidecode 且使用 strings.Trim(..., " ") 而不是 TrimSpace
	// 注意：这里我们只模拟最核心的 dmidecode 场景，因为只有这个场景会产生不一致

	productNameCmd := []string{"-s", "system-product-name"}
	cmd := exec.CommandContext(ctx, "dmidecode", productNameCmd...)
	productNameOutput, err1 := cmd.CombinedOutput()

	serialNumberCmd := []string{"-s", "system-serial-number"}
	cmd = exec.CommandContext(ctx, "dmidecode", serialNumberCmd...)
	serialNumberOutput, err2 := cmd.CombinedOutput()

	if err1 == nil && err2 == nil {
		// 重点：使用旧的 Trim 逻辑
		return utils.Md5V(strings.Trim(string(productNameOutput), " ") + strings.Trim(string(serialNumberOutput), " ")), nil
	}
	return "", errors.New("old generation method failed")
}
