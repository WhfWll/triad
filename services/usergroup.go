package services

import (
	"context"
	"smart/models/mysqls"
	"smart/tools/enums"
	"strconv"
	"time"
)

// 用户管理 - 用户组管理

type UserGroup struct {
}

// AddUserGroup 添加用户组
func (ug *UserGroup) AddUserGroup(ctx context.Context, name, rangeInfo string, pid, userID, isRangeOpen int) error {
	var userGroupModel = mysqls.UserGroups{
		Name:        name,
		Pid:         pid,
		Range:       rangeInfo,
		IsRangeOpen: isRangeOpen,
		UserID:      userID,
		Status:      enums.UserGroupStatusNormal,
		UpdateTime:  time.Now(),
		CreateTime:  time.Now(),
	}
	return userGroupModel.AddUserGroups(ctx)
}

// UpdateUserGroup 更新用户组
func (ug *UserGroup) UpdateUserGroup(ctx context.Context, name, rangeInfo string, pid, isRangeOpen, id int) error {
	var userGroupModel = mysqls.UserGroups{
		ID:          id,
		Name:        name,
		Pid:         pid,
		Range:       rangeInfo,
		IsRangeOpen: isRangeOpen,
		UpdateTime:  time.Now(),
	}
	return userGroupModel.UpdateUserGroups(ctx)
}

// UpdateUserGroupStatus 更新多个用户组状态
func (ug *UserGroup) UpdateUserGroupStatus(ctx context.Context, ids []string, status int) error {
	var userGroupModel = mysqls.UserGroups{
		Status:     status,
		UpdateTime: time.Now(),
	}
	return userGroupModel.UpdateUserGroupsStatus(ctx, ids)
}

// GetUserGroupList 获取用户组列表
func (ug *UserGroup) GetUserGroupList(ctx context.Context, page, limit int) ([]mysqls.UserGroups, int64, error) {
	var userGroup mysqls.UserGroups
	userGroupList, count, err := userGroup.GetUserGroupsList(ctx, page, limit)
	if err != nil {
		return []mysqls.UserGroups{}, 0, err
	}
	return userGroupList, count, nil
}

// GetAllUserGroupList 获取所有用户组列表
func (ug *UserGroup) GetAllUserGroupList(ctx context.Context) ([]mysqls.UserGroups, error) {
	var userGroup mysqls.UserGroups
	userGroupList, err := userGroup.GetAllUserGroupsList(ctx)
	if err != nil {
		return []mysqls.UserGroups{}, err
	}
	return userGroupList, nil
}

// GetUserGroupListByIds 通过ids获取用户组列表
func (ug *UserGroup) GetUserGroupListByIds(ctx context.Context, ids []int) ([]mysqls.UserGroups, error) {
	var userGroup mysqls.UserGroups
	userGroupList, err := userGroup.GetUserGroupsListByIDs(ctx, ids)
	if err != nil {
		return []mysqls.UserGroups{}, err
	}
	return userGroupList, nil
}

// GetUserGroupInfoByPid 通过pid获取用户信息
func (ug *UserGroup) GetUserGroupInfoByPid(ctx context.Context, pid int) (mysqls.UserGroups, error) {
	var userGroup mysqls.UserGroups
	userGroupInfo, err := userGroup.GetUserGroupsInfoByID(ctx, pid)
	if err != nil {
		return mysqls.UserGroups{}, err
	}
	return userGroupInfo, nil
}

// GetUserGroupListByQuery 通过query获取用户组列表
func (ug *UserGroup) GetUserGroupListByQuery(ctx context.Context, query string) ([]mysqls.UserGroups, error) {
	var userGroup mysqls.UserGroups
	userGroupList, err := userGroup.GetUserGroupsListByQuery(ctx, query)
	if err != nil {
		return []mysqls.UserGroups{}, err
	}
	return userGroupList, nil
}

// AllMaps 所有组
func (ug *UserGroup) AllMaps(ctx context.Context) map[int]mysqls.UserGroups {
	var userGroup mysqls.UserGroups
	allGroups := userGroup.All(ctx)
	// 所有组转map
	allGroupMaps := make(map[int]mysqls.UserGroups)
	for _, item := range allGroups {
		allGroupMaps[item.ID] = item
	}
	return allGroupMaps
}

// 用户所属组
// @return [用户ID：[组1，组2，...]]
type UserGroupItem struct {
	GroupStr string
	GroupArr [][]string
}

func (ug *UserGroup) GetUsersGroups(ctx context.Context, userIds []int) map[int]UserGroupItem {
	// 所属组
	var userGroupUsersModel mysqls.UserGroupsUsers
	userGroupUsers := userGroupUsersModel.GetByUserIds(ctx, userIds)

	// 所有组
	var userGroupService UserGroup
	userGroups := userGroupService.AllMaps(ctx)

	// 关联数据 [user_id:组1、组2]
	groupData := make(map[int][]mysqls.UserGroups)
	for _, item := range userGroupUsers {
		// 获取组名称
		if group, ok := userGroups[item.UserGroupsID]; ok {
			if _, ok := groupData[item.UserID]; ok {
				groupData[item.UserID] = append(groupData[item.UserID], group)
			} else {
				gS := make([]mysqls.UserGroups, 0)
				gS = append(gS, group)
				groupData[item.UserID] = gS
			}
		}
	}

	// 依据pid排序，分对应的组
	// [user_id:[group_str:"",group_arr:[[],[],[]]]]

	returnData := make(map[int]UserGroupItem)
	for userId, item := range groupData {
		treeGroups := make([]map[string]string, 0)
		ug.userGroupTree(item, 0, &treeGroups)

		index := 0
		for k, tree := range treeGroups {
			if k == 0 && tree == nil {
				continue
			}

			userItem, ok := returnData[userId]
			if !ok {
				userItem = UserGroupItem{}
				returnData[userId] = userItem
			}

			//returnItem[userId]
			if tree == nil {
				// 下一组
				index++
			}

			if tree["name"] != "" {
				userItem.GroupStr += tree["name"]
				userItem.GroupArr[index] = append(userItem.GroupArr[index], tree["id"])
				returnData[userId] = userItem
			}
		}
	}

	return returnData
}

func (ug *UserGroup) userGroupTree(resources []mysqls.UserGroups, pid int, aaa *[]map[string]string) {
	for _, item := range resources {
		if pid == 0 {
			*aaa = append(*aaa, nil)
		}
		if item.Pid == pid {
			*aaa = append(*aaa, map[string]string{
				"id":   strconv.Itoa(item.ID),
				"name": item.Name,
			})
			ug.userGroupTree(resources, item.ID, aaa)
		}
	}
}
