package services

import (
	"context"
	"smart/models/mysqls"
	"strconv"
	"strings"
	"time"
)

type UserGroupUserList struct {
}

// GetGroupNumber 获取用户组 各个组的数量
func (ug *UserGroupUserList) GetGroupNumber(ctx context.Context) (map[int]int, error) {
	var userGroupUsers mysqls.UserGroupsUsers
	userGroupList, err := userGroupUsers.GetUserGroupsAllUserList(ctx)
	if err != nil {
		return map[int]int{}, err
	}
	groupNumber := make(map[int]int)
	for _, v := range userGroupList {
		if _, ok := groupNumber[v.UserGroupsID]; ok {
			groupNumber[v.UserGroupsID]++
		} else {
			groupNumber[v.UserGroupsID] = 1
		}
	}
	return groupNumber, nil
}

// GetUserGroupListByQuery 通过query获取用户组列表
func (ug *UserGroupUserList) GetUserGroupListByQuery(ctx context.Context, query string) ([]mysqls.UserGroupsUsers, error) {
	var userGroup mysqls.UserGroupsUsers
	userGroupList, err := userGroup.GetUserGroupsUserListByQuery(ctx, query)
	if err != nil {
		return []mysqls.UserGroupsUsers{}, err
	}
	return userGroupList, nil
}

// GetUserGroupListByIds 通过ids获取用户组列表
func (ug *UserGroupUserList) GetUserGroupListByIds(ctx context.Context, ids []int) ([]mysqls.UserGroupsUsers, error) {
	var userGroup mysqls.UserGroupsUsers
	userGroupList, err := userGroup.GetUserGroupsListByUserGroupIDs(ctx, ids)
	if err != nil {
		return []mysqls.UserGroupsUsers{}, err
	}
	return userGroupList, nil
}

// GetUserAndUserGroupName 获取用户对应的用户组名称
func (ug *UserGroupUserList) GetUserAndUserGroupName(ctx context.Context, ids []int) (map[int]string, error) {
	var (
		usergroupSrv   UserGroup
		userGroupUsers mysqls.UserGroupsUsers
	)
	userGroupList, err := userGroupUsers.GetUserGroupsListByUserIDs(ctx, ids)
	if err != nil {
		return map[int]string{}, nil
	}
	groupMap := usergroupSrv.AllMaps(ctx)
	relationMap := make(map[int]string)
	for _, v := range userGroupList {
		relationMap[v.UserID] = groupMap[v.UserGroupsID].Name
	}
	return relationMap, nil
}

// GetUserGroupNamesMap 获取用户对应的用户组名称(支持多组)
func (ug *UserGroupUserList) GetUserGroupNamesMap(ctx context.Context, ids []int) (map[int]string, error) {
	var (
		usergroupSrv   UserGroup
		userGroupUsers mysqls.UserGroupsUsers
	)
	// 获取用户与组的关联
	userGroupList, err := userGroupUsers.GetUserGroupsListByUserIDs(ctx, ids)
	if err != nil {
		return map[int]string{}, nil
	}
	// 获取所有组信息
	groupMap := usergroupSrv.AllMaps(ctx)

	// 拼接
	relationMap := make(map[int][]string)
	for _, v := range userGroupList {
		if group, ok := groupMap[v.UserGroupsID]; ok {
			relationMap[v.UserID] = append(relationMap[v.UserID], group.Name)
		}
	}

	// 转为字符串
	resultMap := make(map[int]string)
	for uid, names := range relationMap {
		resultMap[uid] = strings.Join(names, ",")
	}
	return resultMap, nil
}

// GetUserGroupIdsMap 获取用户对应的用户组ID列表
func (ug *UserGroupUserList) GetUserGroupIdsMap(ctx context.Context, ids []int) (map[int][]int, error) {
	var (
		userGroupUsers mysqls.UserGroupsUsers
	)
	// 获取用户与组的关联
	userGroupList, err := userGroupUsers.GetUserGroupsListByUserIDs(ctx, ids)
	if err != nil {
		return map[int][]int{}, nil
	}

	// 拼接
	relationMap := make(map[int][]int)
	for _, v := range userGroupList {
		relationMap[v.UserID] = append(relationMap[v.UserID], v.UserGroupsID)
	}
	return relationMap, nil
}

// GetUserGroupListByIdsMap 通过ids获取用户组列表map
func (ug *UserGroupUserList) GetUserGroupListByIdsMap(ctx context.Context, ids []int) (map[int]int, error) {
	var userGroup mysqls.UserGroupsUsers
	userGroupList, err := userGroup.GetUserGroupsListByUserGroupIDs(ctx, ids)
	if err != nil {
		return map[int]int{}, err
	}
	// 构造组内用户map
	groupUsersMap := make(map[int]int)
	for _, v := range userGroupList {
		groupUsersMap[v.UserID] = 1
	}
	return groupUsersMap, nil
}

// DelUserGroupListByGroupID 通过用户组ID删除用户组
func (ug *UserGroupUserList) DelUserGroupListByGroupID(ctx context.Context, groupID int) error {
	var userGroup mysqls.UserGroupsUsers
	if err := userGroup.DeleteUserGroupsUsersByGroupID(ctx, groupID); err != nil {
		return err
	}
	return nil
}

// BathAddUserGroupInfo 批量入库
func (ug *UserGroupUserList) BathAddUserGroupInfo(ctx context.Context, userIds []string, groupID, opID int) error {
	var bathData []mysqls.UserGroupsUsers
	for _, v := range userIds {
		UserID, _ := strconv.Atoi(v)
		bathData = append(bathData, mysqls.UserGroupsUsers{
			UserGroupsID: groupID,
			UserID:       UserID,
			SubmitUserID: opID,
			CreateTime:   time.Now(),
		})
	}
	var userGroup mysqls.UserGroupsUsers
	if err := userGroup.AddsUserGroupsUsers(ctx, bathData); err != nil {
		return err
	}
	return nil
}

// GetUserGroupsByIds 通过id获取用户组信息
func (ug *UserGroupUserList) GetUserGroupsByIds(ctx context.Context, ids []int) ([]mysqls.UserGroups, error) {
	var (
		userGroupUserModel mysqls.UserGroupsUsers
		userGroupModel     mysqls.UserGroups
	)
	userGroupList := make([]mysqls.UserGroups, 0)
	userGroupUserList, err := userGroupUserModel.GetUserGroupsListByUserIDs(ctx, ids)
	if err != nil {
		return userGroupList, err
	}
	tempIdList := make([]int, 0)
	for _, item := range userGroupUserList {
		tempIdList = append(tempIdList, item.UserGroupsID)
	}
	return userGroupModel.GetUserGroupsListByIDs(ctx, tempIdList)
}
