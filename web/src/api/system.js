/* 
    系统设置 api 
*/
import axios from '@/axios/http';
const system = {
    //获得版本信息
    getversion(params){
        return axios.get('/smart/system/authinfo', params)
        // return axios.get('/systems/information/version/', params)
    },
    //获得2023 11 18 add
    getPassWordCheck(params){
        return axios.get('/smart/user/passwordcheck', params)
    },
    //可利用评分配置信息
    getusescoreinfo(params){
        return axios.get('/smart/system/usescoreinfo', params)
    },
    //可利用评分配置保存
    getusescoresaveSave(params){
        return axios.post('/smart/system/usescoresave', params)
    },
    //测试范围校验信息
    gettestscopeinfo(params){
        return axios.get('/smart/system/testscopeinfo', params)
    },
    //测试范围校验保存
    gettestscopesave(params){
        return axios.post('/smart/system/testscopesave', params)
    },
    //添加系统授权
    getSystemAuthSave(params){
        return axios.post('/smart/system/authsave', params)
        // return axios.get('/systems/information/version/', params)
    },
    //获得基本信息
    getBasicInfo(params){
        return axios.get('/systems/information/get_basic_info/', params)
    },
    //获得系统uid
    getSystemUid(){
        return axios.get('/systems/interfaces/system_uid/')
    },
    //关机
    systemShutdown(){
        return axios.post('/systems/interfaces/shutdown/')
    },
     //获得设置信息（告警、安全）
     getSysteminfo(){
        return axios.get('/smart/system/initialinfo') 
    },

    // //获得告警设置信息
    // getWarninfo(){
    //     return axios.get('/systems/securities/initial/settings/') 
    // },
    // 保存告警设置
    saveSystemWarn(params){
        return axios.post('/systems/securities/warning/', params)
    },
    // //获得安全设置信息 旧的
    // getSecurityinfo(){
    //     return axios.get('/systems/securities/initial/settings/') 
    // },
    //获得安全设置信息 新的
    getSafeinfo(){
        return axios.get('/system/safe/page') 
    },
    //保存安全设置 新的
    saveSecurity(params){
        return axios.post('/smart/system/securities', params)
    },
    //获得API密钥
    getTokens(params){
        return axios.get('/smart/system/tokenlist', params) 
    },
    tokenDel(params){ //秘钥删除
        return axios.post('/smart/system/tokendel',params)
    },
    //获得用户列表
    getusernames(params){
        return axios.get('/smart/user/list',params) 
    },
    //生成密钥
    submitgenerateToken(params){
        // return axios.post('/systems/auth/generate/token/', params) 
        return axios.post('/smart/system/generatetoken', params) 
    },
    //收集授权信息
    getAuthorinfo(){
        return axios.postBlob('/systems/auth/download/')
    },
    current_version_info(){
        return axios.get('/systems/interfaces/current_version_info/') 
    },
    //系统监控 折线图接口
    getDatabasetrend(){
        return axios.get('/smart/system/cpuinfo')
    },
    //系统监控 折线图接口2
    getDataneicun(){
        return axios.get('/smart/system/memoryinfo')
    },
     //系统监控 饼图接口
    getpiedata(){
        return axios.get('/smart/system/diskinfo')
    },
    //系统离线更新
    getSystemoffupgrade(params){
        return axios.postJson('/smart/system/confirmupgrade', params)
    },
    //手动回滚
    manualRollback(params){
        return axios.postJson('/smart/system/manualrollback', params)
    },

    //获取升级进度
    getUpgradeStatus(){
        return axios.get('/smart/system/upgradestatus')
    },
    getAIenum(){
        return axios.get('/smart/llmmodel/enums')
    },
    //获取大模型列表
    getAiModels(params){
        return axios.get('/smart/llmmodel/list', params)
    },
    //新增/编辑大模型
    addOrUpdateAiModel(params){
        return axios.post('/smart/llmmodel/save', params)
    },
    modelcheck(params){
        return axios.post('/smart/llmmodel/check', params)
    },
    //删除大模型
    deleteAiModels(params){
        return axios.post('/smart/llmmodel/delete', params)
    },
    enhancedetail(){
        return axios.get('/smart/llmmodel/enhance/detail')
    },
    enhanceedit(params){
         return axios.post('/smart/llmmodel/enhance/edit',params)
    },
    //设为默认大模型
    setDefaultAiModel(params){
        return axios.post('/smart/llmmodel/setdefault', params)
    },
    //获取场景配置
    getScenarioConfig(params){
        return axios.get('/smart/llmmodel/scenarioconfig', params)
    },
    //保存场景配置
    saveScenarioConfig(params){
        return axios.post('/smart/llmmodel/scenarioconfig/save', params)
    },
    //获取AI场景应用配置列表
    getAiScenarioList(params){
        return axios.get('/smart/aiscenario/list', params)
    },
    //AI场景测试接口
    testAiScenario(params){
        return axios.post('/smart/aiscenario/test', params)
    },
    //应用场景配置接口
    configAiScenario(params){
        return axios.post('/smart/aiscenario/config', params)
    },
}
const otherset = {
    //获得设置信息（告警、安全）
    getSysteminfo(){
        return axios.get('/smart/system/monitorwarninfo') 
    },
    // 保存告警设置
    saveSystemWarn(params){
    return axios.post('/smart/system/monitorwarnsave', params)
    },
    //获取 系统访问白名单
    getWhite(){
        return axios.get('/smart/system/ipwhiteinfo') 
    },
    // 保存系统访问白名单
    saveWhite(params){
        return axios.post('/smart/system/ipwhitesave', params)
    },
    //获取syslog获取接口
    getSyslog(){
        return axios.get('/smart/system/sysloginfo') 
    },
    //syslog
    saveSyslog(params){
        return axios.post('/smart/system/syslogsave', params)
    },
    //邮箱配置获取接口
    getEmail(){
        // return axios.get('/systems/mail/info/') 
        return axios.get('/smart/system/mailinfo') 
    },
    // 保存邮箱配置
    saveEmail(params){
        // return axios.post('/systems/mail/save/', params)
        return axios.post('/smart/system/mailsave', params)
    },
    //网络配置获取接口
    getNetwork(){
        return axios.get('/smart/system/networkconfiginfo') 
    },
    // 保存网络配置
    saveNetwork(params){
        return axios.post('/smart/system/networkconfigsave', params)
    },
    //路由获取接口
      getLuyou(){
        return axios.get('/smart/system/routelist') 
    },
    // 新增路由
    addnewLuyou(params){
        return axios.post('/smart/system/routeadd', params)
    },
    // 删除路由
   deleteLuyou(params){
        return axios.post('/smart/system/routedelete', params)
    },

}
const node = {
    // 节点列表
    getData(params){
        return axios.get('/smart/system/nodelist', params)
    },
    // 新增节点
    addNode(params){
        return axios.post('/smart/system/nodeadd', params);
    },
    // 删除
    nodeDel(params) {
        return axios.post('/smart/system/nodedel', params)
    },
    // 启用、禁用
    nodeauthorize(params){
        return axios.post('/smart/system/nodedisorenable', params);
    },
    // 获取所有可用节点
    getEnableNode(params){
        return axios.get('/smart/system/nodeallenable', params)
    },
    // 获取是否开启分布式引擎
    getStatus(params){
        return axios.get('/smart/system/nodegetdistribute', params)
    },
     // 设置是否开启分布式引擎
     setStatus(params){
        return axios.post('/smart/system/nodesetdistribute', params)
    },
    // getCPUChart(params){ //CPU详情
    //     return axios.get('/node/cpu/', params)
    // },
    // getRAMChart(params){ //内存详情
    //     return axios.get('/node/memory/', params)
    // },
    // 节点详情
    getNodeInfo(params){
        return axios.get('/smart/system/nodeinfo', params)
    },
    // 编辑节点
    editNode(params){
        return axios.post('/smart/system/nodeedit', params)
    },
    // saveUpdateNode(params) { //编辑保存
    //     return axios.post('/node/edit/', params)
    // },
    // 下载
    download (params) {
        return axios.get('/smart/system/nodedownload', params) 
    }
}

const proxy = {
    getData(params) { //代理列表
        return axios.get('/node/proxy/', params)
    },
    addProxy(params){ //新增
        return axios.post('/node/proxy/', params)
    },
    MultDelete(params){ //删除
        return axios.delete('/node/proxy/delete/', params)
    },
    editProxy(params) { 
        return axios.post('/node/proxy/edit/', params)
    }
}

//业务工具
const businessset = {
    getTCP(){
        return axios.get('/smart/system/tcpblindtestinfo',)
    },
    saveTCP(params){
        return axios.post('/smart/system/tcpblindtestsave', params)
    },
    getHTTP() {
        return axios.get('/smart/system/httpblindtestinfo',)
    },
    saveHTTP(params) {
        return axios.post('/smart/system/httpblindtestsave', params)
    },
    getDNS() {
        return axios.get('/smart/system/dnsblindtestinfo',)     
    },
    saveDNS(params) {
        return axios.post('/smart/system/dnsblindtestsave', params)
    },
    getconcurrency(){
        return axios.get('/smart/system/curtasksinfo',)
    },
    //并发
    saveconcurrency(params){
        return axios.get('/smart/system/curtaskssave', params)
    },
    getreverse(){
        return axios.get('/smart/system/getreverseiphost',)
    },
    savereverse(params){
        return axios.post('/smart/system/reverseiphostsave', params)
    },
    getTargetIp(){
        return axios.get('/smart/system/targetiplist',)
    },
    saveTargetIp(params){
        return axios.postJson('/smart/system/targetipsave', params)
    },
}
const config = {
    //保存设置
    saveConfiginfo(params){
        return axios.post('/smart/system/configbackupconfigsave', params)
    },
    saveCopyinfo(params){
        return axios.post('/smart/system/configbackupnow', params)
    },
    //恢复
    liveCopy(params){
        return axios.post('/smart/system/configbackuprestore', params)
    },
    //获得列表
    getlogs(params){
        return axios.get('/smart/system/configbackuplist', params) 
    },
    getselectinfo(){
        return axios.get('/systems/backup/params/') 
    },
    getconfiginfo(){
        return axios.get('/smart/system/configbackupconfiginfo') 
    },
    handleDelinfo(params){
        return axios.post('/smart/system/configbackupdelete', params)
    },
    // 下载
    download (params) {
        return axios.get('/smart/system/configbackupdownload', params) 
    }
}

export {
    system,
    node,
    proxy,
    otherset,
    businessset,
    config
} ;


//
export function current_version_info() {
    return axios.get('/systems/interfaces/current_version_info/')
}
