import axios from '@/axios/axios'
// 登出
export function logout(params){
  return axios({
    url:'/smart/user/logout',
    method:'get',
    params
    // url:'/accounts/logout/'
  })
}

// 任务uid
export function getUid(){
  return axios({
    url:'/task/task/uid/'
  })
}


// 登录检测认证状态
export function loginCheckStatus(data){
  // console.log('执行');
  return axios({
    url:'/task/task/login_check/status/',
    method:'post',
    data
  })
}
// 登录检测下拉菜单
export function loginSelectList(){
  return axios({
    url:'task/task/login_check/select/'
  })
}
// 执行任务
export function implementTask(data){
  return axios({
    url:'/task/task/',
    method:'post',
    ...data
  })
}
// 上传文件
export function updateFile(data){
  return axios({
    // url:'/systems/interfaces/off_update/',
    url:'/smart/system/uploadupgradefile',
    method:'post',
    data
  })
}