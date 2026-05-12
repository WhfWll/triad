package services

import (
	"context"
	"encoding/json"
	"gitlabee.4dogs.cn/common/redis"
	"smart/models/mysqls"
	"smart/tools/data"
	"smart/tools/enums"
	"strconv"
	"time"
)

// 枚举服务

type MapSet struct {
}

// GetProductAuthState 获取授权状态
func (m *MapSet) GetProductAuthState(ctx context.Context) (isAuth bool) {
	isAuth = false
	var mapset mysqls.MapSet
	mapSetRes, err := mapset.GetsByObjKey(ctx, "productAuthState")
	if err != nil {
		isAuth = false
	}
	// 如果返回1 代表已授权
	if mapSetRes.ObjValue == "1" {
		isAuth = true
	}
	return
}

// GetMapValue 获取map_set信息
func (m *MapSet) GetMapValue(ctx context.Context, objKey string) (string, error) {
	var mapset mysqls.MapSet
	mapSetRes, err := mapset.GetsByObjKey(ctx, objKey)
	if err != nil {
		return "", err
	}
	return mapSetRes.ObjValue, nil
}

// UpdateMapValue 更新map_set信息
func (m *MapSet) UpdateMapValue(ctx context.Context, objKey, objValue string) error {
	var mapset mysqls.MapSet
	err := mapset.UpdateObjValueByObjKey(ctx, objKey, objValue)
	if err != nil {
		return err
	}
	return nil
}

// 创建
func (m *MapSet) Create(ctx context.Context, objKey, objValue, content string) (err error) {
	var mapSetModel mysqls.MapSet
	mapSetData, err := mapSetModel.GetsByObjKey(ctx, objKey)
	if err != nil {
		return err
	}
	if mapSetData.Estate == "" {
		mapSetData.Estate = enums.MapSetEstateValid
	}
	mapSetData.ObjKey = objKey
	mapSetData.ObjValue = objValue
	mapSetData.Content = content

	if mapSetData.ID == 0 {
		// 创建
		err = mapSetData.AddMapSet(ctx)
	} else {
		// 更新
		m.DelSystemSecurity(ctx)
		err = mapSetData.UpdateMapSet(ctx)
	}
	return
}

// GetsSystemSecurity 安全策略 - 获取
func (m *MapSet) GetsSystemSecurity(ctx context.Context, isResetCache bool) (returnData enums.SystemsSystemSecurity, err error) {
	// 优先获取缓存数据
	cacheClient, err := redis.NewClient()
	systemsecurity := cacheClient.Get(ctx, "systemsecurity")
	if isResetCache == false && systemsecurity.Val() != "" {
		err = json.Unmarshal([]byte(systemsecurity.Val()), &returnData)
		if returnData.SystemSecurityLoginTimeout == 0 {
			returnData.SystemSecurityLoginTimeout = 60000
		}
		if returnData.SystemSecurityUserLimit == 0 {
			returnData.SystemSecurityUserLimit = 5
		}
		if returnData.SystemSecurityAdminLimit == 0 {
			returnData.SystemSecurityAdminLimit = 5
		}
		return
	}
	// 缓存不存在，从数据库获取
	var mysetModel mysqls.MapSet
	data := mysetModel.GetsLikeLeftObjKey(ctx, "systemSecurity")
	mapData := make(map[string]int)
	for _, item := range data {
		mapData[item.ObjKey], _ = strconv.Atoi(item.ObjValue)
	}
	mapDataByte, err := json.Marshal(mapData)
	if err != nil {
		return
	}
	err = json.Unmarshal(mapDataByte, &returnData)

	// 设置默认值
	if returnData.SystemSecurityLoginTimeout == 0 {
		returnData.SystemSecurityLoginTimeout = 60000
	}
	if returnData.SystemSecurityUserLimit == 0 {
		returnData.SystemSecurityUserLimit = 5
	}
	if returnData.SystemSecurityAdminLimit == 0 {
		returnData.SystemSecurityAdminLimit = 5
	}
	if returnData.SystemSecurityBanTime == 0 {
		returnData.SystemSecurityBanTime = 30
	}

	// 设置缓存
	GlobalSystemsecurityByte, err := json.Marshal(returnData)
	if err != nil {
		return
	}
	cacheClient.Set(ctx, "systemsecurity", string(GlobalSystemsecurityByte), 10*time.Hour)
	return
}

// TargetIpWhiteBlack 测试目标黑白名单对测试目标进行过滤
func (m *MapSet) TargetIpWhiteBlack(ctx context.Context, targetList []string) ([]string, error) {
	var mapset mysqls.MapSet
	mapSetRes, err := mapset.GetsByObjKey(ctx, enums.TargetIpWhiteBlackMapSetObjKey)
	if err != nil || len(mapSetRes.ObjValue) == 0 {
		return targetList, err
	}
	var objValue TargetIpWhiteBlackMapSet
	err = json.Unmarshal([]byte(mapSetRes.ObjValue), &objValue)
	if err != nil || objValue.IsOpen == 0 {
		return targetList, err
	}
	if objValue.IsOpen == enums.TargetIpWhiteBlackIsOpenOff { //如果黑白名单关闭，不做过滤， 直接返回
		return targetList, nil
	}
	if objValue.Type == enums.TargetIpWhiteBlackTypeWhite { //开启白名单
		targetList = data.StringArrayWhiteList(targetList, objValue.IpListArray)
	} else if objValue.Type == enums.TargetIpWhiteBlackTypeBlack { //开启黑名单
		targetList = data.StringArrayBlackList(targetList, objValue.IpListArray)
	}
	return targetList, nil
}

// SystemConfigListByMapKeys 通过objKey数组获取系统配置列表
func (m *MapSet) SystemConfigListByMapKeys(ctx context.Context, objKeys any) []mysqls.MapSet {
	var mapSetModel mysqls.MapSet
	return mapSetModel.ListByObjKeys(ctx, objKeys)
}

// DelSystemSecurity 删除安全策略
func (m *MapSet) DelSystemSecurity(ctx context.Context) {
	cacheClient, _ := redis.NewClient()
	cacheClient.Del(ctx, "systemsecurity")
}
