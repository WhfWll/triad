package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"smart/models/redises"
	"smart/services"
	"smart/tools/enums"
	"smart/tools/errno"
	"smart/tools/network"
	"smart/tools/utils"
	"strconv"
	"strings"
)

// AccessWhite 验证访问ip是否在系统ip白名单中
func AccessWhite(c *gin.Context) {
	var mapSetService services.MapSet
	ctx := context.Background()
	ipWhiteObjValue, err := mapSetService.GetMapValue(ctx, enums.SystemAccessIpWhiteMapSetObjKey)

	// 添加管理员账户可以不受系统ip白名单限制逻辑
	uid, ok := c.Get("uid")
	if ok {
		var cacheClient redises.RedisStr
		cacheAdminId, _ := cacheClient.Get(ctx, enums.AdminCacheKey)
		if cacheAdminId == "" {
			var srv services.User
			user, err := srv.GetForUsername(ctx, "admin")
			if err != nil || user.ID == 0 {
				log.Println("GetForUsername error: ", err.Error())
				cacheAdminId = "1"
			} else {
				cacheAdminId = strconv.Itoa(user.ID)
			}
			err = cacheClient.SetString1(ctx, enums.AdminCacheKey, cacheAdminId, 24*60*60)
		}
		if cacheAdminId != "" {
			if cacheAdminId == strconv.Itoa(uid.(int)) {
				c.Next()
				return
			}
		}
	}
	if err != nil {
		fail(c, errno.LoginAuthErr, err)
		return
	}
	if ipWhiteObjValue == "" {
		c.Next()
		return
	}
	var ipWhiteMapSet services.SystemSettingIpWhiteMapSet
	if err = json.Unmarshal([]byte(ipWhiteObjValue), &ipWhiteMapSet); err != nil {
		fail(c, errno.LoginAuthErr, err)
		return
	}
	if ipWhiteMapSet.IsOpen == enums.ConfigClose {
		c.Next()
		return
	}

	// 但前请求地址
	uri := c.Request.URL.String()
	// 是否在白名单
	ok, err = isWhiteApi(uri)
	if err != nil {
		fail(c, errno.LoginAuthErr, err)
		return
	}
	if ok {
		c.Next()
		return
	}

	accessIp := utils.GetClientIp(c)
	fmt.Println("access ip: ", accessIp)
	ipList := []string{"127.0.0.1", "localhost"} //本地ip永远在白名单中
	for _, item := range strings.Split(ipWhiteMapSet.Ip, ",") {
		if network.IpSegmentTools.VerifyNetmaskIpSegment(item) {
			netmaskIpList, err := network.IpSegmentTools.HandleNetmaskIpSegment(item)
			if err != nil {
				fail(c, errno.LoginAuthErr, err)
				return
			}
			ipList = append(ipList, netmaskIpList...)
		} else if network.IpSegmentTools.VerifyCrossbarIpSegment(item) {
			crossbarIpList, err := network.IpSegmentTools.HandleCrossbarIpSegment(item)
			if err != nil {
				fail(c, errno.LoginAuthErr, err)
			}
			ipList = append(ipList, crossbarIpList...)
		} else if network.IpSegmentTools.VerifyIp(item) {
			ipList = append(ipList, item)
		} else {
			fail(c, errno.LoginAuthErr, errors.New("白名单数据异常"))
			return
		}
	}
	if !utils.UrlHandleLogic.In(accessIp, ipList) {
		fail(c, errno.LoginAuthErr, errors.New("当前主机IP不允许访问系统"))
		return
	}

	// 允许继续
	c.Next()
}
