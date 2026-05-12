package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"gitlabee.4dogs.cn/common/config"
	"reflect"
	"time"
)

type JwtFactory struct {
}

func (jwtFactory *JwtFactory) Encode(mapClaims jwt.MapClaims) (string, error) {
	// 秘钥
	var jwtSecret string
	if err := config.Load("jwt_secret", &jwtSecret); err != nil {
		return "", err
	}
	if jwtSecret == "" {
		return "", errors.New("jwt secret Unknown")
	}
	// 创建Token结构体
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	// 调用加密方法，发挥Token字符串
	signingString, err := claims.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return signingString, nil
}

func (jwtFactory *JwtFactory) Auth(token string) (userId int, err error) {
	// 秘钥
	var jwtSecret string
	if err := config.Load("jwt_secret", &jwtSecret); err != nil {
		return 0, err
	}
	if jwtSecret == "" {
		return 0, errors.New("jwt secret Unknown")
	}
	// 创建Token结构体
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}
	data, ok := tokenObj.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("token 解析失败")
	}

	userIdInter, ok := data["uid"]
	if !ok {
		return 0, errors.New("token 未检测到必要信息")
	}

	if reflect.TypeOf(userIdInter).String() != "float64" {
		return 0, errors.New("token 数据类型错误")
	}

	exp, ok := data["exp"]
	if !ok {
		return 0, errors.New("token 未检测到必要信息")
	}

	if int64(exp.(float64)) <= time.Now().Unix() {
		return 0, errors.New("token 已过期")
	}
	return int(userIdInter.(float64)), nil
}
