package services

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"smart/api/typespec"
	"smart/models/mysqls"
	"smart/models/redises"
	"smart/tools/auth"
	"smart/tools/encryption"
	"smart/tools/enums"
	"smart/tools/utils"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gitlabee.4dogs.cn/common/config"
)

// 用户管理 - 用户管理
type User struct {
}

// UpdateIsAlive 退出登录修改在线状态
func (user *User) UpdateIsAlive(ctx context.Context, uid, isAlive int, token string) error {
	userModel := mysqls.User{
		ID:              uid,
		IsAlive:         isAlive,
		LastOperateTime: time.Now(),
		LastLoginTime:   time.Now(),
		OperatorID:      uid,
		Token:           token,
		TokenCreateTime: time.Now(),
	}
	if err := userModel.UpdateUser(ctx); err != nil {
		return err
	}
	var redisHash redises.RedisHash
	if err := redisHash.HDel(ctx, enums.RedisUserAlive, strconv.Itoa(uid)); err != nil {
		return err
	}
	return nil
}

// UpdateUserStatus 修改用户状态
func (user *User) UpdateUserStatus(ctx context.Context, opUid, uid, status int) error {
	userModel := mysqls.User{
		ID:         uid,
		Status:     status,
		UpdateTime: time.Now(),
		OperatorID: opUid,
	}
	if err := userModel.UpdateUser(ctx); err != nil {
		return err
	}
	return nil
}

// UpdateAccountExpireTime 修改账号过期时间
func (user *User) UpdateAccountExpireTime(ctx context.Context, opUid, uid int, accountExpireTime string) error {
	expireTime, _ := time.Parse(enums.TimeYMDBarLayout, accountExpireTime)
	userModel := mysqls.User{
		ID:                uid,
		AccountExpireTime: expireTime,
		UpdateTime:        time.Now(),
		LastOperateTime:   time.Now(),
		Status:            enums.UserStatusSuccess,
		OperatorID:        opUid,
	}
	if err := userModel.UpdateUser(ctx); err != nil {
		return err
	}
	return nil
}

// UpdatePassWord 修改用户密码
func (user *User) UpdatePassWord(ctx context.Context, opUid, uid int, password, oldPassword string) error {
	userModel := mysqls.User{
		ID:                  uid,
		Password:            password,
		UpdateTime:          time.Now(),
		OperatorID:          opUid,
		PasswordChangeTime:  time.Now(),
		PasswordAlreadyUsed: oldPassword,
	}
	if err := userModel.UpdateUser(ctx); err != nil {
		return err
	}
	return nil
}

// DelUser 删除用户
func (user *User) DelUser(ctx context.Context, opUid int, uids string) error {
	ids := strings.Split(uids, ",")
	userModel := mysqls.User{
		Estate:     "deleted",
		UpdateTime: time.Now(),
		OperatorID: opUid,
	}
	if err := userModel.DeleteUser(ctx, ids); err != nil {
		return err
	}
	return nil
}

// GetUserDetail 获取用户信息
func (user *User) GetUserDetail(ctx context.Context, userID int) (userInfo mysqls.User, err error) {
	userModel := mysqls.User{
		ID: userID,
	}
	userInfo, err = userModel.GetUser(ctx)
	if err != nil {
		return
	}
	return
}

// GetUserAllValidList 获取所有有效的用户列表
func (user *User) GetUserAllValidList(ctx context.Context) (userList []mysqls.User, count int64, err error) {
	userModel := mysqls.User{
		Estate: "valid",
	}
	userList, count, err = userModel.GetAllUserList(ctx)
	if err != nil {
		return
	}
	return
}

// GetUserList 获取用户列表
func (user *User) GetUserList(ctx context.Context, search string, page, size int) (userList []mysqls.User, count int64, err error) {
	userModel := mysqls.User{
		Estate: "valid",
	}
	userList, count, err = userModel.GetUserList(ctx, page, size, search)
	if err != nil {
		return
	}
	return
}

// GenerateToken 生成token
func (u *User) GenerateToken(uid int) (string, error) {
	var mapSet MapSet
	security, err := mapSet.GetsSystemSecurity(context.Background(), false)
	if err != nil {
		return "", err
	}
	expireSeconds := (time.Duration(security.SystemSecurityLoginTimeout) * time.Minute).Seconds()
	tokenPayload := jwt.MapClaims{
		"uid": uid,
		"exp": int(time.Now().Unix() + int64(expireSeconds)),
	}
	var jwt auth.JwtFactory
	return jwt.Encode(tokenPayload)
}

// AllForIds 通过IDs获取所有用户
func (a *User) AllForIds(ctx context.Context, userIds []int) (mapData map[int]mysqls.User, err error) {
	var userModel mysqls.User
	list, err := userModel.GetByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}
	mapData = make(map[int]mysqls.User)
	for _, item := range list {
		mapData[item.ID] = item
	}
	return
}

// GetForId 通过ID获取用户
func (a *User) GetForId(ctx context.Context, userId int) (mysqls.User, error) {
	var userModel mysqls.User
	userModel.ID = userId
	return userModel.GetUser(ctx)
}

// GetForUsername 通过username获取用户
func (a *User) GetForUsername(ctx context.Context, username string) (mysqls.User, error) {
	var userModel mysqls.User
	return userModel.GetByUsername(ctx, username)
}

// UpdateStatusById 依据用户ID设置用户状态
func (a *User) UpdateStatusById(ctx context.Context, id, status int) error {
	var userModel mysqls.User
	return userModel.UpdateStatusById(ctx, id, status)
}

// UpdateStatusById 依据用户ID设置用户最后操作时间
func (a *User) UpdateLastOperateById(ctx context.Context, id int) error {
	var userModel mysqls.User
	return userModel.UpdateLastOperateById(ctx, id)
}

// CheckAccountLongNoLogin 长时间未登录，设置账号禁用（仅普通账号）
func (a *User) CheckAccountLongNoLogin(ctx context.Context) error {
	// 获取账号有多产时间未登录 禁用 单位 月
	var mapSet MapSet
	systemsecurity, err := mapSet.GetsSystemSecurity(ctx, false)
	if err != nil {
		return err
	}

	accountNoLoginDisSecond := systemsecurity.SystemSecurityAccountNoLoginDisable * 30 * 24 * 3600

	// 账号多长时间未有操作
	expireTime := time.Now().Unix() - int64(accountNoLoginDisSecond)
	var usersModel mysqls.User
	users := usersModel.GetsByRoleAndLastOperateTime(ctx, enums.UserRoleOrdinary, expireTime, enums.UserStatusSuccess)

	for _, item := range users {
		fmt.Println("修改用户状态为：账号长时间未登录禁用；用户ID为：" + strconv.Itoa(int(item.ID)))
		item.Status = enums.UserStatusNoLoginDisable
		err = item.UpdateUser(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckAccountExpire 账号有效期，禁用
func (a *User) CheckAccountExpire(ctx context.Context) (err error) {
	var usersModel mysqls.User
	expireTime := time.Now().Unix()
	users := usersModel.GetsByRoleAndAccountTime(ctx, enums.UserRoleOrdinary, expireTime, enums.UserStatusSuccess)

	for _, item := range users {
		item.Status = enums.UserStatusAccountExpireDisable
		err = item.UpdateUser(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// CheckPasswordExpire 密码有效期校验，过期前如未重新修改，设置账号禁用（仅普通账号）
func (a *User) CheckPasswordExpire(ctx context.Context) (err error) {
	expireTime := time.Now().Unix()
	var usersModel mysqls.User
	users := usersModel.GetsByRoleAndPasswordTime(ctx, enums.UserRoleOrdinary, expireTime, enums.UserStatusSuccess)

	for _, item := range users {
		item.Status = enums.UserStatusPassExpireDisable
		err = item.UpdateUser(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (user *User) GetUserForId(ctx context.Context, userId int) (data mysqls.User, err error) {
	var model mysqls.User
	model.ID = userId
	data, err = model.GetUser(ctx)
	if err != nil {
		return
	}

	return
}

func (user *User) CheckUserIsAlready(ctx context.Context, userId int, username, email string) (data bool, err error) {
	var model mysqls.User
	data, err = model.CheckUserIsAlready(ctx, userId, username, email)
	return
}

func (user *User) GetUserForToken(ctx context.Context, token string) (data mysqls.User, err error) {
	var model mysqls.User
	data, err = model.GetUserByToken(ctx, token)
	if err != nil {
		return
	}

	return
}

// 列表 - 分页
func (user *User) Page(ctx context.Context, page, size int) ([]mysqls.User, int64, error) {
	var userModel mysqls.User
	return userModel.GetUserList(ctx, page, size, "")
}

// 是否管理员
func (user *User) UserIsAdmin(ctx context.Context, uid int) error {
	var myUser mysqls.User
	myUser.ID = uid
	data, err := myUser.GetUser(ctx)
	if err != nil {
		return err
	}
	if data.ID == 0 {
		return errors.New("未查询到您的信息，请尝试重新登录")
	}
	if data.Type != enums.UserRoleSuperAdmin {
		return errors.New("您非超级管理员，无法操作")
	}
	return nil
}

// 密码强度
func (user *User) CheckPassRank(ctx context.Context, password string) error {
	if len(password) < 8 {
		return errors.New("密码长度至少8位")
	}

	ok, err := regexp.MatchString(`[a-zA-Z]+`, password)
	if err != nil {
		return err
	}
	need := 0
	if ok {
		need++
	}

	ok, err = regexp.MatchString(`[0-9]+`, password)
	if err != nil {
		return err
	}
	if ok {
		need++
	}

	ok, err = regexp.MatchString(`[^a-zA-Z0-9]+`, password)
	if err != nil {
		return err
	}
	if ok {
		need++
	}

	var mapSet MapSet
	security, err := mapSet.GetsSystemSecurity(ctx, false)
	if err != nil {
		return err
	}
	switch security.SystemSecurityPasswordRank {
	case enums.SystemSecurityPassRankMiddle:
		if need < 2 {
			return errors.New("密码至少包含字母/数字/字符其中2种")
		}
	case enums.SystemSecurityPassRankHigh:
		if need < 3 {
			return errors.New("密码至少包含字母和数字和字符")
		}
	}
	return nil
}

// CreateUser 创建用户
func (user *User) CreateUser(ctx context.Context, username, password, accountExpire, email, groupIds, department, remark string, createUid, role int) error {
	var (
		userModel mysqls.User
		mapSet    MapSet
	)
	// 检查邮箱是否已存在
	userData, err := userModel.GetByUsernameOrEmail(ctx, username, email)
	if userData.ID != 0 {
		return errors.New("邮箱/用户名重复")
	}
	security, err := mapSet.GetsSystemSecurity(ctx, false)
	if err != nil {
		return err
	}
	accountExpireTime, _ := time.Parse("2006-01-02", accountExpire)
	userModel = mysqls.User{
		Username:          username,
		Password:          password,
		Email:             email,
		AccountExpireTime: accountExpireTime,
		IsAlive:           enums.UserOffline,
		Department:        department,
		OperatorID:        createUid,
		Remark:            remark,
		Type:              role,
		Status:            enums.UserStatusSuccess,
		//LastOperateTime:    time.Now(),
		//LastLoginTime:      time.Now(),
		CreateTime:         time.Now(),
		PasswordExpireTime: time.Unix(time.Now().Unix()+int64(security.SystemSecurityPasswordCycle*(3600*24*30)), 0),
		PasswordChangeTime: time.Now(),
	}
	if err = userModel.AddUser(ctx); err != nil {
		return err
	}
	// 添加用户组
	if groupIds != "" {
		var userGroupUsers mysqls.UserGroupsUsers
		if err := userGroupUsers.AddUserGroupsUsersByString(ctx, userModel.ID, groupIds); err != nil {
			return err
		}
	}
	return nil
}

// UpdateUser 更新用户
func (user *User) UpdateUser(ctx context.Context, createUid, uid, role int, username, email, department, remark, groupIds string) error {
	var userModel mysqls.User
	userModel.ID = uid
	userData, _ := userModel.GetByUsernameOrEmail(ctx, username, email)
	updateInfo := map[string]interface{}{
		"username":          username,
		"email":             email,
		"operator_id":       createUid,
		"department":        department,
		"remark":            remark,
		"type":              role,
		"status":            enums.UserStatusSuccess,
		"last_operate_time": time.Now(),
		"update_time":       time.Now(),
	}
	if err := userData.EditUserInfo(ctx, updateInfo); err != nil {
		return err
	}
	// 更新用户组
	var userGroupUsers mysqls.UserGroupsUsers
	if err := userGroupUsers.UpdateUserGroupsUsersByString(ctx, uid, groupIds); err != nil {
		return err
	}
	return nil
}

// CheckUserAddParam 添加用户操作动作
func (user *User) CheckUserAddParam(ctx context.Context, req typespec.UserManageOpReq, uid int) string {
	// 1 验证创建人身份超级管理员可以编辑管理员与普通用户 -> 管理员可以编辑普通用户信息 -> 普通用户只允许编辑自己信息
	opUserInfo, _ := user.GetUserDetail(ctx, uid)
	if opUserInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	if opUserInfo.Type == enums.UserRoleSuperAdmin {
	} else if opUserInfo.Type == enums.UserRoleAdministrator {
		if req.Role == enums.UserRoleSuperAdmin {
			return "管理员不能操作超级管理员"
		}
	} else if opUserInfo.Type == enums.UserRoleOrdinary {
		return "普通用户无权创建用户"
	}
	// 添加用户 验证密码准确性
	if req.Password != req.Repassword {
		return "两次密码输入不一致，请确认"
	}
	// 验证邮箱
	if !utils.IsEmail(req.Email) {
		return "电子邮箱格式错误，请确认"
	}
	// 检查密码强度
	if err := user.CheckPassRank(ctx, req.Password); err != nil {
		return err.Error()
	}
	// 检查密码是否包含完整用户名
	if err := user.CheckPassContainsUser(ctx, req.Password, req.Username); err != nil {
		return err.Error()
	}
	// 检查是否包含键盘排序密码
	if err := user.CheckPassOrderByKeyBoard(ctx, req.Password); err != nil {
		return err.Error()
	}
	return ""
}

// CheckUserEditParam 修改用户操作动作
func (user *User) CheckUserEditParam(ctx context.Context, req typespec.UserManageOpReq, uid int) string {
	// 1 检查用户是否存在
	userInfo, err := user.GetUserForId(ctx, req.ID)
	if err != nil || userInfo.ID == 0 {
		return "未知的用户"
	}
	// 2 验证创建人身份超级管理员可以编辑管理员与普通用户 -> 管理员可以编辑普通用户信息 -> 普通用户只允许编辑自己信息
	opUserInfo, _ := user.GetUserDetail(ctx, uid)
	if opUserInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	if opUserInfo.Type != enums.UserRoleSuperAdmin && req.Role != userInfo.Type {
		return "只有超级管理员可以修改身份"
	}
	if opUserInfo.Type == enums.UserRoleSuperAdmin {
	} else if opUserInfo.Type == enums.UserRoleAdministrator {
		if req.Role == enums.UserRoleSuperAdmin {
			return "管理员不能操作超级管理员"
		}
	} else if opUserInfo.Type == enums.UserRoleOrdinary {
		if req.ID != uid {
			return "普通用户只能修改自己的信息"
		}
	}
	// 3 验证邮箱
	if !utils.IsEmail(req.Email) {
		return "电子邮箱格式错误，请确认"
	}
	// 4 检查邮箱是否已存在
	ok, err := user.CheckUserIsAlready(ctx, req.ID, req.Username, req.Email)
	if err != nil || ok {
		return "用户名或邮箱已存在，请确认"
	}
	return ""
}

// PasswordEncryption 密码加密
func PasswordEncryption(password string) string {
	var iterations int
	err := config.Load("pbkdf2_sha256_iterations", &iterations)
	if err != nil {
		return password
	}
	var pbkdf2 encryption.Pbkdf2Sha256
	encryptionPass, err := pbkdf2.PasswordEncode(password, "", iterations)
	if err != nil {
		return password
	}
	return encryptionPass
}

// PasswordDecodeString 密码解密
func PasswordDecodeString(password string) string {
	var (
		aesCbc  encryption.AesCbc
		certVal string
	)
	if password == "" {
		return ""
	}
	if err := config.Load("aesCbcUserCert", &certVal); err != nil {
		return ""
	}
	if certVal == "" {
		return ""
	}
	hexDePass, err := hex.DecodeString(password)
	if err != nil {
		return ""
	}
	return string(aesCbc.AesDecryptCBC(hexDePass, []byte(certVal)))
}

// CheckUserEditPassParam 处理用户密码修改
func (user *User) CheckUserEditPassParam(ctx context.Context, req typespec.UserPassPWReq, uid int) string {
	// 1 修改密码校验[修改密码只允许用户修改自己的密码]
	if req.UserID != 0 && uid != req.UserID {
		if req.Password != "" {
			return "只能修改自己的密码"
		}
	}
	// 2 验证密码是否重复
	if req.Password != req.Repassword {
		return "两次密码输入不一致，请确认"
	}
	// 3 验证密码准确性
	var pbkdf encryption.Pbkdf2Sha256
	userInfo, err := user.GetUserForId(ctx, req.UserID)
	if err != nil {
		return "无法判断用户身份"
	}
	if ok, _ := pbkdf.PasswordVerify(PasswordDecodeString(req.OldPassword), userInfo.Password); !ok {
		return "原始密码错误"
	} else {
		if req.OldPassword == req.Password {
			return "新密码与旧密码一致：请重设密码"
		}
	}
	// 检查密码强度
	if err := user.CheckPassRank(ctx, req.Password); err != nil {
		return err.Error()
	}
	if err := user.CheckPassRecent(ctx, req.Password, userInfo); err != nil {
		return err.Error()
	}
	return ""
}

// CheckUserPassParam 重置密码
func (user *User) CheckUserPassParam(ctx context.Context, opid, uid int, password string) string {
	opUserInfo, _ := user.GetUserDetail(ctx, opid)
	if opUserInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	userInfo, _ := user.GetUserDetail(ctx, uid)
	if uid == opid {
		return "不能操作自己"
	}
	if userInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	if opUserInfo.Type == enums.UserRoleSuperAdmin {
	} else if opUserInfo.Type == enums.UserRoleAdministrator {
		if userInfo.Type == enums.UserRoleSuperAdmin {
			return "管理员不能操作超级管理员"
		}
	} else if opUserInfo.Type == enums.UserRoleOrdinary {
		return "普通用户无权修改用户状态"
	}
	// 3 验证密码准确性
	var pbkdf encryption.Pbkdf2Sha256
	userInfo, err := user.GetUserForId(ctx, opid)
	if err != nil {
		return "无法判断用户身份"
	}
	if ok, _ := pbkdf.PasswordVerify(password, userInfo.Password); !ok {
		return "原始密码错误"
	}
	return ""
}

// CheckUserStatusParam 修改状态和其他
func (user *User) CheckUserStatusParam(ctx context.Context, opid, uid, editStatus int) string {
	opUserInfo, _ := user.GetUserDetail(ctx, opid)
	if opUserInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	userInfo, _ := user.GetUserDetail(ctx, uid)
	if uid == opid {
		return "不能操作自己"
	}
	if userInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	if editStatus == enums.UserStatusSuccess && userInfo.Status != enums.UserStatusSetDisable {
		return "只能修改手动禁用的账户状态"
	}
	if opUserInfo.Type == enums.UserRoleSuperAdmin {
	} else if opUserInfo.Type == enums.UserRoleAdministrator {
		if userInfo.Type == enums.UserRoleSuperAdmin {
			return "管理员不能操作超级管理员"
		}
	} else if opUserInfo.Type == enums.UserRoleOrdinary {
		if userInfo.Type == enums.UserRoleSuperAdmin || userInfo.Type == enums.UserRoleAdministrator {
			return "普通用户不能操作管理员和超级管理员"
		}
	}
	return ""
}

// CheckUserExpireTimeParam 修改有效期
func (user *User) CheckUserExpireTimeParam(ctx context.Context, opid, uid int, accountExpireTime string) string {
	opUserInfo, _ := user.GetUserDetail(ctx, opid)
	if opUserInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	userInfo, _ := user.GetUserDetail(ctx, uid)
	if uid == opid {
		return "不能操作自己"
	}
	if userInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	if userInfo.Status == enums.UserStatusAccountExpireDisable && accountExpireTime < time.Now().Format(enums.ResTimeDayLayout) {
		return "用户已过期 修改时间错误 无效"
	}
	if opUserInfo.Type == enums.UserRoleSuperAdmin {
	} else if opUserInfo.Type == enums.UserRoleAdministrator {
		if userInfo.Type == enums.UserRoleSuperAdmin {
			return "管理员不能操作超级管理员"
		}
	} else if opUserInfo.Type == enums.UserRoleOrdinary {
		if userInfo.Type == enums.UserRoleSuperAdmin || userInfo.Type == enums.UserRoleAdministrator {
			return "普通用户不能操作管理员和超级管理员"
		}
	}
	return ""
}

// CheckUserDelParam 删除校验
func (user *User) CheckUserDelParam(ctx context.Context, opid int, uids string) string {
	opUserInfo, _ := user.GetUserDetail(ctx, opid)
	if opUserInfo.ID == 0 {
		return "未查询到您的信息，请尝试重新登录"
	}
	ids := strings.Split(uids, ",")
	for _, uidStr := range ids {
		uid, _ := strconv.Atoi(uidStr)
		if uid == opid {
			return "不能操作自己"
		}
		userInfo, _ := user.GetUserDetail(ctx, uid)
		if userInfo.ID == 0 {
			return "未查询到您的信息，请尝试重新登录"
		}
		if userInfo.Type == enums.UserRoleAuditor {
			return "审计员不能被删除 请检查是否选择错误"
		}
		if opUserInfo.Type == enums.UserRoleSuperAdmin {
		} else if opUserInfo.Type == enums.UserRoleAdministrator {
			if userInfo.Type == enums.UserRoleSuperAdmin {
				return "管理员不能操作超级管理员"
			}
		} else if opUserInfo.Type == enums.UserRoleOrdinary {
			return "普通用户无权删除用户"
		}
	}
	return ""
}

// GeneratePlatformToken 生成平台token
func (u *User) GeneratePlatformToken(ctx context.Context, username string) (string, error) {
	originData := "4dogs.cn" + username + time.Now().Format(enums.TimeLayout)
	platformToken := utils.Md5V(originData)
	var userModel mysqls.User
	err := userModel.UpdatePlatformToken(ctx, username, platformToken)
	if err != nil {
		return "", err
	}
	return platformToken, nil
}

// GetUserListByToken 获取用户列表
func (user *User) GetUserListByToken(ctx context.Context, search string, page, size int) (userList []mysqls.User, count int64, err error) {
	userModel := mysqls.User{
		Estate: "valid",
	}
	userList, count, err = userModel.GetUserListByToken(ctx, search, page, size)
	if err != nil {
		return
	}
	return
}

// GetUserListByIds 通过IDs获取所有用户
func (user *User) GetUserListByIds(ctx context.Context, userIds []string) (mapData map[int]mysqls.User, err error) {
	var userModel mysqls.User
	list, err := userModel.GetUserListByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}
	mapData = make(map[int]mysqls.User)
	for _, item := range list {
		mapData[item.ID] = item
	}
	return
}

// GetUserListByType 获取对应角色的用户列表
func (user *User) GetUserListByType(ctx context.Context, role, page, size int, search string) (userList []mysqls.User, count int64, err error) {
	userModel := mysqls.User{
		Estate: "valid",
		Type:   role,
	}
	userList, count, err = userModel.GetAllUserListByType(ctx, page, size, search)
	if err != nil {
		return
	}
	return
}

// UpdatePasswordChangeTime 修改密码更新时间
func (user *User) UpdatePasswordChangeTime(ctx context.Context, uid int, changeTime time.Time) error {
	var userModel mysqls.User
	if err := userModel.UpdatePasswordChangeTime(ctx, uid, changeTime); err != nil {
		return err
	}
	return nil
}

// CheckPassContainsUser 检测密码是否包含用户名
func (user *User) CheckPassContainsUser(ctx context.Context, password, username string) error {
	if strings.Contains(strings.ToLower(password), strings.ToLower(username)) {
		return errors.New("密码不能包含完整用户名字符")
	}
	return nil
}

// CheckPassOrderByKeyBoard 检测密码是否按键盘排序
func (user *User) CheckPassOrderByKeyBoard(ctx context.Context, password string) error {
	if utils.IsKeyboardSorted(password) {
		return errors.New("请避免键盘排序密码")
	}
	return nil
}

// CheckPassRecent 检查是否使用最近的密码
func (user *User) CheckPassRecent(ctx context.Context, password string, userInfo mysqls.User) error {
	if PasswordDecodeString(userInfo.PasswordAlreadyUsed) == password {
		return errors.New("不能使用上次使用的密码")
	}
	return nil
}

func (user *User) AesEncryptPassword(password string) (string, error) {
	var certVal string
	if err := config.Load("aesCbcUserCert", &certVal); err != nil {
		return "", err
	}
	key := []byte(certVal) // 加密的密钥
	aesEcb := encryption.AesEcb{}
	encrypted := aesEcb.AesEncryptECB([]byte(password), key)
	fmt.Println("密文：", hex.EncodeToString(encrypted))
	return hex.EncodeToString(encrypted), nil
}

// TokenDelete 秘钥删除接口
func (user *User) TokenDelete(ctx context.Context, username, token string) (err error) {
	var userModel mysqls.User
	err = userModel.TokenDelete(ctx, username, token)
	if err != nil {
		return
	}
	return
}
