package application

import (
	"context"
	"errors"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/enums"
	"strings"
)

type UserGroupUserList struct {
}

// GroupUserPreselectionList 组内用户预选列表
func (ugul *UserGroupUserList) GroupUserPreselectionList(ctx context.Context, req *typespec.UserGroupUserPreselectionReq, res *typespec.UserGroupUserPreselectionRes) (err error) {
	var (
		userService          services.User
		userGroupSrv         services.UserGroup
		userGroupUserListSrv services.UserGroupUserList
	)
	// 1 验证是否存在
	userGroupInfo, err := userGroupSrv.GetUserGroupListByIds(ctx, []int{req.Id})
	if err != nil {
		return err
	}
	if len(userGroupInfo) > 0 && userGroupInfo[0].ID == 0 {
		return errors.New("未知的组，请重新选择组")
	}
	// 2 获取用户列表
	userList, count, err := userService.GetUserListByType(ctx, 1, req.Page, req.Size, req.Keyword)
	if err != nil {
		return
	}
	// 3 组内用户列表
	groupUsersMap, err := userGroupUserListSrv.GetUserGroupListByIdsMap(ctx, []int{req.Id})
	if err != nil {
		return
	}
	// 4 获取用户ID组
	ids := make([]int, 0)
	for _, item := range userList {
		ids = append(ids, item.ID)
	}
	// 5 获取用户ID和用户组名称的关系
	userGroupNameMap, err := userGroupUserListSrv.GetUserAndUserGroupName(ctx, ids)
	if err != nil {
		return
	}
	// 6 处理返回数据
	for _, item := range userList {
		var user typespec.UserGroupUserPreselectionItemRes
		user.Id = item.ID
		user.Name = item.Username
		user.Role = item.Type
		user.RoleStr = enums.RoleChoice[item.Type]
		user.Status = item.Status
		user.StatusStr = enums.StatusMap[item.Status]
		user.GroupList = userGroupNameMap[item.ID]
		// 是否已在组内
		if _, ok := groupUsersMap[item.ID]; ok {
			user.Selected = true
		} else {
			user.Selected = false
		}
		res.List = append(res.List, user)
	}
	res.Total = int(count)
	res.Page = req.Page
	res.Size = req.Size
	return
}

// GroupUserAlreadyList 用户组 - 组内成员 - 已选列表
func (ugul *UserGroupUserList) GroupUserAlreadyList(ctx context.Context, req *typespec.UserGroupUserAlreadyReq, res *typespec.UserGroupUserAlreadyRes) (err error) {
	var (
		userGroupUserListSrv services.UserGroupUserList
		userSrv              services.User
	)
	groupUsers, err := userGroupUserListSrv.GetUserGroupListByIds(ctx, []int{req.Id})
	if err != nil {
		return
	}
	if len(groupUsers) > 0 {
		userIds := make([]int, 0)
		for _, item := range groupUsers {
			userIds = append(userIds, item.UserID)
		}
		users, _ := userSrv.AllForIds(ctx, userIds)
		for _, item := range users {
			res.List = append(res.List, typespec.UserGroupUserAlreadyItemRes{
				Id:   item.ID,
				Name: item.Username,
			})
		}
	}
	return
}

// GroupUserRelation 用户组 - 组内成员 - 已选列表
func (ugul *UserGroupUserList) GroupUserRelation(ctx context.Context, opID int, req *typespec.UserGroupUserRelationReq, res *typespec.UserGroupUserRelationRes) (err error) {
	var (
		userService          services.User
		userGroupUserListSrv services.UserGroupUserList
	)
	if err = userService.UserIsAdmin(ctx, opID); err != nil {
		return err
	}
	// 1 验证用户是否存在
	userIds := strings.Split(req.UserIds, ",")
	users, _ := userService.GetUserListByIds(ctx, userIds)
	if len(users) != len(userIds) {
		return errors.New("已选列表存在未知用户")
	}
	// 2 判断是否为普通用户
	for _, item := range users {
		if item.Type != enums.UserRoleOrdinary {
			return errors.New(item.Username + " 非普通用户，无法进组")
		}
	}
	// 3根据用户组 对应的组内成员
	if err = userGroupUserListSrv.DelUserGroupListByGroupID(ctx, req.GroupId); err != nil {
		return
	}
	// 4 添加新的用户组成员
	if err = userGroupUserListSrv.BathAddUserGroupInfo(ctx, userIds, req.GroupId, opID); err != nil {
		return
	}
	return
}
