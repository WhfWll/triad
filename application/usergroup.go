package application

import (
	"context"
	"errors"
	"smart/api/typespec"
	"smart/services"
	"smart/tools/data"
	"strconv"
	"strings"
)

type UserGroup struct {
}

// UserGroupCreate 创建用户组
func (ug *UserGroup) UserGroupCreate(ctx context.Context, opid int, req *typespec.UserGroupCreateReq, res *typespec.UserGroupCreateRes) (err error) {
	var (
		userService  services.User
		userGroupSrv services.UserGroup
	)
	if err = userService.UserIsAdmin(ctx, opid); err != nil {
		return err
	}
	// 如果开启了范围测试那么必须有值
	if req.RangeOpen != 0 && req.Range == "" {
		return errors.New("范围测试不可为空")
	}
	// 当前最大支持3级 需验证
	if req.Pid != 0 {
		if checkUserGroupLevel(ctx, req.Pid, 1) {
			return errors.New("最大仅支持3级")
		}
	}
	// 添加用户组信息
	if err = userGroupSrv.AddUserGroup(ctx, req.Name, req.Range, req.Pid, opid, req.RangeOpen); err != nil {
		return err
	}
	return
}

// UserGroupSelect 选择用户组
func (ug *UserGroup) UserGroupSelect(ctx context.Context, res *typespec.UserGroupSelectRes) (err error) {
	var userGroupSrv services.UserGroup
	userGroupsList, err := userGroupSrv.GetAllUserGroupList(ctx)
	if err != nil {
		return
	}
	var treeData []data.TreeList
	for _, v := range userGroupsList {
		treeData = append(treeData, data.TreeList{
			Id:    v.ID,
			PId:   v.Pid,
			Value: strconv.Itoa(v.ID),
			Label: v.Name,
		})
	}
	res.List = data.Func.Tree(treeData)
	return
}

// UserGroupEdit 修改用户组
func (ug *UserGroup) UserGroupEdit(ctx context.Context, opid int, req *typespec.UserGroupUpdateReq, res *typespec.UserGroupUpdateRes) (err error) {
	var (
		userService  services.User
		userGroupSrv services.UserGroup
	)
	// 默认组不允许编辑
	if req.Id == 1 {
		return errors.New("默认组禁止编辑")
	}
	// 判断是否是管理员
	if err = userService.UserIsAdmin(ctx, opid); err != nil {
		return err
	}
	// 如果开启了范围测试那么必须有值
	if req.RangeOpen != 0 && req.Range == "" {
		return errors.New("范围测试不可为空")
	}
	// 当前最大支持3级 需验证
	if req.Pid != 0 {
		if checkUserGroupLevel(ctx, req.Pid, 1) {
			return errors.New("最大仅支持3级")
		}
	}
	// 更新用户组信息
	if err = userGroupSrv.UpdateUserGroup(ctx, req.Name, req.Range, req.Pid, req.RangeOpen, req.Id); err != nil {
		return err
	}
	return
}

// UserGroupUpdateStatus 修改 1/多个 用户组状态
func (ug *UserGroup) UserGroupUpdateStatus(ctx context.Context, opid int, req *typespec.UserGroupUpdateStatusReq, res *typespec.UserGroupUpdateStatusRes) (err error) {
	var (
		userService  services.User
		userGroupSrv services.UserGroup
	)
	// 判断是否是管理员
	if err = userService.UserIsAdmin(ctx, opid); err != nil {
		return err
	}
	ids := strings.Split(req.Ids, ",")
	for _, id := range ids {
		if id == "1" {
			return errors.New("默认组禁止编辑")
		}
	}
	// 更新用户组信息
	if err = userGroupSrv.UpdateUserGroupStatus(ctx, ids, req.Status); err != nil {
		return err
	}
	return
}

// checkUserGroupLevel 递归检查用户组的层级是够超过三级
func checkUserGroupLevel(ctx context.Context, pid int, level int) bool {
	if level > 3 {
		return true // 已经超过三级
	}
	var userGroupUsersSer services.UserGroup
	info, _ := userGroupUsersSer.GetUserGroupListByIds(ctx, []int{pid})
	if len(info) == 0 {
		return false
	}
	if info[0].ID == 0 {
		return false
	}
	return checkUserGroupLevel(ctx, info[0].Pid, level+1)
}

// GetUserGroupList 用户组列表
func (ug *UserGroup) GetUserGroupList(ctx context.Context, req *typespec.UserGroupListReq, res *typespec.UserGroupListRes) error {
	var (
		userGroupUserListService services.UserGroupUserList
		userGroupListItems       []typespec.UserGroupListItemRes
		userGroupService         services.UserGroup
	)
	// 获取用户组信息
	userGroupList, count, err := userGroupService.GetUserGroupList(ctx, req.Page, req.Size)
	if err != nil {
		return err
	}
	// 获取有父节点的组
	pids := make([]int, 0)
	for _, item := range userGroupList {
		if item.Pid != 0 {
			pids = append(pids, item.Pid)
		}
	}
	// 构造group name和id的map
	pmap := make(map[int]string)
	if len(pids) > 0 {
		list, err := userGroupService.GetUserGroupListByIds(ctx, pids)
		if err != nil {
			return err
		}
		for _, info := range list {
			pmap[info.ID] = info.Name
		}
	}
	// 记录组内有多少用户
	ids := make([]int, 0)
	for _, item := range userGroupList {
		ids = append(ids, item.Pid)
	}
	// 获取组员相关信息
	groupNumberMap, _ := userGroupUserListService.GetGroupNumber(ctx)
	for _, item := range userGroupList {
		var group typespec.UserGroupListItemRes
		group.Id = item.ID
		group.Name = item.Name
		group.CreateTime = item.CreateTime.String()
		group.Pid = item.Pid
		group.Number = groupNumberMap[item.ID]
		group.RangeOpen = item.IsRangeOpen
		group.Range = item.Range
		if item.Pid == 0 {
			group.PidStr = "无"
		} else {
			group.PidStr = pmap[item.Pid]
			if group.PidStr == "" {
				group.PidStr = "无"
			}
			pidArr := make([]interface{}, 0)
			pidArr = append(pidArr, strconv.Itoa(item.ID))
			pidArr = append(pidArr, strconv.Itoa(item.Pid))
			// 查询父级ID
			pval, _ := userGroupService.GetUserGroupInfoByPid(ctx, item.Pid)
			if pval.ID != 0 {
				pidArr = append(pidArr, strconv.Itoa(pval.Pid))
			}
			//if len(pidArr) > 1 {
			//	garr := garray.NewFrom(pidArr, true)
			//	group.PidArr = append(group.PidArr, garr.Reverse().Slice())
			//} else {
			group.PidArr = append(group.PidArr, pidArr)
			//}
		}
		userGroupListItems = append(userGroupListItems, group)
	}
	*res = typespec.UserGroupListRes{
		Page:  req.Page,
		Size:  req.Size,
		Total: int(count),
		List:  userGroupListItems,
	}
	return nil
}
