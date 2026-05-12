package utils

import (
	"errors"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"math"
	"reflect"
)

type OtherPayloadParam struct {
	Exp      int64
	Iat      int64
	Iss      string
	Jti      string
	RealName string
	Sub      string
	Ticket   string
}

func CheckOtherLoginPayloadParam(data jwtv5.MapClaims) (payload OtherPayloadParam, err error) {

	exp, expOk := data["exp"] // 过期时间
	if !expOk {
		return payload, errors.New("token中未检测到 exp")
	}
	if reflect.TypeOf(exp).String() != "float64" {
		return payload, errors.New("token中检测到 exp 类型非 float64 请确认约定格式是否正确")
	}
	payload.Exp = int64(math.Floor(exp.(float64)))

	iat, iatOk := data["iat"] // 颁发时间
	if !iatOk {
		return payload, errors.New("token中未检测到 iat")
	}
	if reflect.TypeOf(iat).String() != "float64" {
		return payload, errors.New("token中检测到 iat 类型非 float64 请确认约定格式是否正确")
	}
	payload.Iat = int64(math.Floor(iat.(float64)))

	iss, issOk := data["iss"] // 颁发时间
	if !issOk {
		return payload, errors.New("token中未检测到 iss")
	}
	if reflect.TypeOf(iss).String() != "string" {
		return payload, errors.New("token中检测到 iss 类型非 string 请确认约定格式是否正确")
	}
	payload.Iss = iss.(string)

	jti, jtiOk := data["jti"] // id
	if !jtiOk {
		return payload, errors.New("token中未检测到 jti")
	}
	if reflect.TypeOf(jti).String() != "string" {
		return payload, errors.New("token中检测到 jti 类型非 string 请确认约定格式是否正确")
	}
	payload.Jti = jti.(string)

	realname, realnameOk := data["realname"] // realname
	if !realnameOk {
		return payload, errors.New("token中未检测到 realname")
	}
	if reflect.TypeOf(realname).String() != "string" {
		return payload, errors.New("token中检测到 realname 类型非 string 请确认约定格式是否正确")
	}
	payload.RealName = realname.(string)

	sub, subOk := data["sub"] // realname
	if !subOk {
		return payload, errors.New("token中未检测到 sub")
	}
	if reflect.TypeOf(sub).String() != "string" {
		return payload, errors.New("token中检测到 sub 类型非 string 请确认约定格式是否正确")
	}
	payload.Sub = sub.(string)

	ticket, ticketOk := data["ticket"] // ticket
	if !ticketOk {
		return payload, errors.New("token中未检测到 ticket")
	}
	if reflect.TypeOf(ticket).String() != "string" {
		return payload, errors.New("token中检测到 ticket 类型非 string 请确认约定格式是否正确")
	}
	payload.Ticket = ticket.(string)
	return payload, nil
}
