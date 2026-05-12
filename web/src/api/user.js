/* 
    用户管理
*/
import axios from '@/axios/http';
const user = {
    //用户列表
    userList(params) {
        return axios.get('/smart/user/list', params)
    },
    //用户详情
    userInfo(params) {
        return axios.get('/smart/user/info', params)
    },
    //用户删除-批量
    userMultDelete(params) {
        return axios.get('/smart/user/del', params)
    },
    //单个删除
    userDel(params) {
        return axios.get('/smart/user/del', params)
    },
    //新增用户
    addUser(params) {
        return axios.post('/smart/user/save', params)
    },
    //编辑用户
    updateUser(params) {
        return axios.post('/smart/user/save', params)
    },
    //修改密码
    updatePwd(params) {
        return axios.post('/smart/user/updatepw', params)
    },
    //修改密码
    resetPwd(params) {
        return axios.post('/smart/user/resetpw', params)
    },
    //修改有效期
    updateExp(params) {
        return axios.post('/smart/user/updateuserexp', params)
    },
    //设置禁用启用用户
    setIsLoginUser(params) {
        return axios.get('/smart/user/changestatus' , params)
    },
    //youxiaoqi 
    saveYouxiaoqi(params) {
        return axios.post('/smart/user/list/update_account_expire', params)
    },
    // 用户组bufen接口
    usergroupList(params) {
        return axios.get('/smart/usergroup/page', params)
    },
    addUsergroup(params) {
        return axios.post('/smart/usergroup/create', params)
    },
    updateUsergroup(params) {
        return axios.post('/smart/usergroup/update', params)
    },
    getupgroup(params) {//新建弹窗内下拉 上级组
        return axios.get('/smart/usergroup/select', params)
    },
    // 用户组成员预选列表
    getWaitUserList(params) {
        return axios.post('/smart/usergroup/preselectionpage', params)
    },
    // 用户组成员已选列表
    getSelectedUserList(params) {
        return axios.post('/smart/usergroup/alreadypage', params)
    },
    // 添加组成员
    addGroupUser(params) {
        return axios.post('/smart/usergroup/relation', params)
    },
    //单个删除
    groupDel(params) {
        return axios.post('/smart/usergroup/updatestatus', params)
    },
}
export default user;