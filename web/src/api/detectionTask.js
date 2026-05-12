
import http from '@/axios/http'
 

const task = {
    taskEnum(){ //枚举
        return http.get('/smart/task/enum')
    },

    taskList(params){ //WiFi任务列表
        return http.get('/smart/tripartite/wifilist', params)
    },
    wifiList(){ //WiFi任务列表
        return http.get('/smart/tripartite/wifiaplist')
    },

    taskListDetail(params){ //详情列表
        return http.get('/smart/tripartite/xraydetaillist', params)
    },
    taskListDetail2(params){ //burp详情列表
        return http.get('/smart/tripartite/burpsuitedetaillist', params)
    },
    taskCopy(params){
        return http.get('/smart/task/copy', params)
    },
    //删除/批量删除
    taskDel(params){
        return http.post('/smart/tripartite/wifidel', params)
    }, 
    //创建wifi任务--保存
    createWIFITask(params){
        return http.post('/smart/tripartite/wificreate', params)
    }, 
    //上传
    taskUpload(params){
        return http.post('/smart/tripartite/xrayupload', params)
    }, 
    //burp删除/批量删除
    taskDel2(params){
        return http.post('/smart/tripartite/burpsuitedel', params)
    }, 
    taskchangestate(params){  //暂停/结束/开始或重启任务
        return http.get('/smart/task/changestate', params)
    },
    websitelogincheck(params){ //网站登陆凭证校验
        return http.post('/smart/task/websitelogincheck', params)
    },
    taskSave(params){ //保存任务 --创建
        return http.post('/smart/tripartite/xraysave', params)
    },
    taskSave2(params){ //保存任务 --burp创建
        return http.post('/smart/tripartite/burpsuitesave', params)
    },

    //任务详情--任务概览
    getTaskOverview(params){ 
        return http.get('/smart/task/overview', params)
    },
    // 任务状态
    getTaskStatus(params){
        return http.get('/smart/task/getstate', params)
    },
    // 任务详情-测试目标
    getTargetlist(params){
        return http.get('/smart/task/targetlist', params)
    },
    targetdel(params){
        return http.get('/smart/task/targetdel', params)
    },
    targetStop(params){
        return http.get('/smart/task/targetchangestate', params)
    },
    taskVulList (params) {
        return http.get('/smart/task/vullist', params)
    },
    taskVulinfo (params) {
        return http.get('/smart/task/vulinfo', params)
    },
    //删除/批量删除
    taskVulDelete (params) {
        return http.get('/smart/task/vuldel', params)
    },
    //验证
    vulValidate (params) {
        return http.post('/smart/task/vulverify', params)
    },
    //验证
    vulTest (params) {
        return http.post('/smart/task/vultest', params)
    },
    //变更状态
    updateVulStatus (params) {
        return http.post('/task/vul/status_update/', params)
    },
    // 测试日志
    getLoglist(params){
        return http.get('/smart/task/loglist', params)
    },
    loginfo(params){
        return http.get('/smart/task/loginfo', params)
    },
};
 
export  {
    task,
}
