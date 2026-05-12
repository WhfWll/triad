package application

import (
	"context"
	"reflect"
	"smart/api/typespec"
	"smart/services"
	"strconv"
)

type MapSet struct {
}

// CreateSystemSecurity 系统管理 - 系统配置 - 安全设置
func (m *MapSet) CreateSystemSecurity(ctx context.Context, req *typespec.MapSetSystemSecurityReq, resp *typespec.MapSetSystemSecurityRes) (err error) {
	var mapSetService services.MapSet
	// 反射取值
	t := reflect.TypeOf(*req)
	v := reflect.ValueOf(*req)
	for i := 0; i < t.NumField(); i++ {
		content := t.Field(i).Tag.Get("dc")
		objKey := t.Field(i).Tag.Get("json")
		objValue := v.Field(i).Int()
		err = mapSetService.Create(ctx, objKey, strconv.Itoa(int(objValue)), content)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetSystemSecurityInfo 系统管理 - 系统配置 - 安全设置信息
func (m *MapSet) GetSystemSecurityInfo(ctx context.Context, req *typespec.SystemSecurityInfoReq, resp *typespec.SystemSecurityInfoRes) (err error) {
	var mapSetService services.MapSet
	returnData, err := mapSetService.GetsSystemSecurity(ctx, true)
	if err != nil {
		return err
	}
	*resp = typespec.SystemSecurityInfoRes{
		Success: true,
		Security: typespec.SecurityInfo{
			PasswordRank:  returnData.SystemSecurityPasswordRank,
			PasswordCycle: returnData.SystemSecurityPasswordCycle,
			LoginTimeout:  returnData.SystemSecurityLoginTimeout,
			UserLimit:     returnData.SystemSecurityUserLimit,
			AdminLimit:    returnData.SystemSecurityAdminLimit,
			BanTime:       returnData.SystemSecurityBanTime,
			PasswordValid: returnData.SystemSecurityPasswordValid,
			ExpireUnused:  returnData.SystemSecurityExpireUnused,
		},
	}
	return nil
}
