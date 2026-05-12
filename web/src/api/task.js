
import http from '@/axios/http'
 

const task = {
    getNodeData(params){ //节点列表
        return http.get('/smart/system/nodeallenable', params)
    },
    taskEnum(){ //枚举
        return http.get('/smart/task/enum')
    },
    taskList(params){ //任务列表
        return http.get('/smart/task/list', params)
    },
    taskCopy(params){
        return http.get('/smart/task/copy', params)
    },
    taskDel(params){
        return http.post('/smart/task/del', params)
    }, 
    taskExport(params){ //导出
        return http.get('/smart/task/taskthreeexport', params)
    },
    taskchangestate(params){  //暂停/结束/开始或重启任务
        return http.get('/smart/task/changestate', params)
    },
    websitelogincheck(params){ //网站登陆凭证校验
        return http.post('/smart/task/websitelogincheck', params)
    },
    taskSave(params){ //保存任务
        return http.post('/smart/task/save', params)
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
    // 测试
    testvultest (params) {
        return http.get('/smart/task/testvultest', params)
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
    // 8.2-------------
    //断开
    riskList(params){
        return http.get('/smart/task/vulevidencelist', params)
    },
    // 远程会话-列表
    riskList2(params){
        return http.get('/smart/task/remotesessionlist', params)
    },
   
    // 取证详情
    riskDetail(params){
        return http.get('/smart/task/vulevidenceinfo', params)
    },
     // 远程会话-详情
    riskDetail2(params){
        return http.get('/smart/task/remotesessioninfo', params)
    },
    riskType(params){
        return http.get('/smart/task/risktypeenum', params)
    },
    // 删除-取证
    delvulevidence(params){
        return http.get('/smart/task/delvulevidence', params)
    },
    // 删除-远程会话
    delvulevidence2(params){
        return http.get('/smart/task/delremotesession', params)
    },
    // 抓取信息下拉
    captureinfoenum(params){
        return http.get('/smart/task/captureinfoenum', params)
    },
    // 抓取信息
    captureinfo(params){
        return http.get('/smart/task/captureinfo', params)
    },
    break(params){
        return http.get('/smart/task/break', params)
    },
    // 文件目录
    filemanagement(params){
        return http.post('/smart/task/filemanagement', params)
    },
    
    //新增目标
    addtarget(params){
        return http.get('/smart/task/addtarget', params)
    },
    //动态添加攻击面
    addattackface(params){
        return http.postJson('/smart/task/addattackface', params)
    },
    //动态添加漏洞
    addvul(params){
        return http.postJson('/smart/task/addvul', params)
    },
    //路径图
    tasktargetmap(params){
        return http.get('/smart/task/tasktargetmap', params)
    },
    //攻击拓扑-路径图
    tasktopologymap(params){
        return http.get('/smart/task/tasktopologymap', params)
    },
    //攻击拓扑-节点详情1
    tasktopologymapnodedetail(params){
        return http.get('/smart/task/tasktopologymapnodedetail', params)
    },

    //已下载--删除
    delfile(params){
        return http.get('/smart/task/delfile', params)
    },
    
    //已下载--下载
    downloadfile(params){
        return http.get('/smart/task/downloadfile', params)
    },
    //任务配置详情
    configinfo(params){
        return http.get('/smart/task/configinfo', params)
    },
    // 新增批量收集信息
    exceshellmany (params) {
        return http.post('/smart/task/exceshellmany', params)
    },
    
};

const task_group = { 
    //任务组列表
    grouplist(params){
        return http.get('/smart/taskgroup/list', params)
    },
    // 新建任务组
    createGroup(params){
        return http.post('/smart/taskgroup/create', params)
    },
    // 编辑任务组
    updateGroup(params){
        return http.post('/smart/taskgroup/groupedit', params)
    },
    //删除任务组
    delGroup(params){
        return http.get('/smart/taskgroup/delete', params)
    },
    //获取任务组内任务的接口
    getTaskbyGroup(params){
        return http.get('/smart/taskgroup/tasklist', params)
    },
    //获取任务统计信息接口
    getTaskGroupOverview(params){
        return http.get('/smart/taskgroup/overview', params)
    },
    //获取任务运行状态接口
    getTaskGroupStatus(params){
        return http.get('/smart/taskgroup/status', params)
    },
    groupBindTask(params){ //任务组与任务绑定
        return http.post('/smart/taskgroup/groupbind', params)
    },
};
const logic={
    logicEnum(){ //枚举接口
        return http.get('/smart/logic/enum')
    },
    getTasklist(params){ //任务列表
        return http.get('/smart/logic/tasklist', params)
    },
    logicTaskStop(params){ //任务结束
        return http.get('/smart/logic/taskstop', params)
    },
    logicTaskcopy(params){
        return http.get('/smart/logic/taskcopy', params)
    },
    getTargetlist(params){ //目标列表
        return http.get('/smart/logic/targetlist', params)
    },
    getVulnlist(params){ //漏洞列表
        return http.get('/smart/logic/vullist', params)
    },
    getLoglist(params){ //日志列表
        return http.get('/smart/logic/loglist', params) 
    },
    vulnDetail(params){
        return http.get('/smart/logic/vulinfo', params) 
    },
    loginfo(params){
        return http.get('/smart/logic/loginfo', params) 
    },
    taskStop(params){
        return http.get('/smart/logic/taskstop', params) 
    },
    taskDel(params){
        return http.get('/smart/logic/taskdel', params) 
    },
    createtask(params){
        return http.post('/smart/logic/taskcreate', params) 
    },
    vultest(params){ //漏洞测试
        return http.post('/smart/logic/vultest', params)
    },
    vuldel(params){
        return http.get('/smart/logic/vuldel', params) 
    },
}
export  {
    task,
    task_group,
    logic,
}
