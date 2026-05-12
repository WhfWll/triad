import axios from '@/axios/http';

const bas = {
    basEnum(){
        return axios.get('/smart/bas/enum') 
    },
    //枚举
    basruleenum(){
        return axios.get('/smart/bas/basruleenum') 
    },
    getBasrule(params){
        return axios.get('/smart/bas/basruleget', params) 
    },
    // 规则详情
    basruleinfo(params){
        return axios.get('/smart/bas/basruleinfo', params) 
    },
    getbastemplate(params){
        return axios.get('/smart/bas/bastemplateget', params) 
    },
    // 编辑规则
    basruleedit(params){
        return axios.post('/smart/bas/basruleedit', params) 
    },
    // 评估方案-删除--剧本集删除
    bastemplatedel(params){
        return axios.post('/smart/bas/bastemplatedel', params) 
    }, 
    //默认
    bastemplatesetdefault(params){
        return axios.post('/smart/bas/bastemplatesetdefault', params) 
    },
    saveScene(params){
        return axios.postJson('/smart/bas/bastemplatecreate', params) 
    },
    bastemplatebyid(params){
        return axios.get('/smart/bas/bastemplatebyid', params) 
    },
     // 是否在线
     basagentisonline(params){
        return axios.postJson('/smart/bas/basagentisonline', params) 
     },
    //bas任务列表
    basTasklist(params){ 
        return axios.get('/smart/bas/bastaskget',params)
    },
    // 删除
    basTaskdel(params){
        return axios.postJson('/smart/bas/bastaskdel',params)
    },
    // 结束
    basTaskstop(params){
        return axios.postJson('/smart/bas/bastaskend',params)
    },
    // 创建任务
    basTaskcreate(params){
        return axios.postJson('/smart/bas/taskcreate',params)
    },
    getSelectNodelist(params){
        return axios.get('/smart/bas/basagentlive',params)
    },
    //目标列表
    basTargetlist(params){
        return axios.get('/smart/bas/bastasktargetpage',params)
    },
    //目标删除
    basTargetdel(params){
        return axios.postJson('/smart/bas/bastasktargetdel',params)
    },
    // 漏洞测试统计
    basvulstat(params){
        return axios.get('/smart/bas/basvulstat',params)
    },
    // 漏洞测试列表
    basvullist(params){
        return axios.get('/smart/bas/basvullist',params)
    },
    // 漏洞测试删除
    basvuldel(params){
        return axios.get('/smart/bas/basvuldel',params)
    },
    //日志列表
    basTargetlog(params){
        return axios.get('/smart/bas/bastasktargetlog',params)
    },
    // agent列表
    agentlist(params){
        return axios.get('/smart/bas/basagentlist',params)
    },
    // agent节点状态变更
    basagentstatusedit(params){
        return axios.postJson('/smart/bas/basagentstatusedit',params)
    },
    //下载
    downagent(params){
        return axios.get('/smart/bas/basagentdownload',params)
    },
}
export default bas;