/* 工具管理
*   fingerprint：指纹库
*/
import axios from '@/axios/http';
import $ajax from '@/axios/axios';
const fingerprint = {
    // 删除
    deleteFingerprint(params) {
        return axios.get('/smart/tools/delfinger', params)
    },
    //漏洞对象统计
    getObjectStatistics(params) {
        return axios.get('/tools/fingerprint/type/statistics/', params)
    },
    //指纹类型下拉
    getVulObjectlist(params) {
        return axios.get('/smart/tools/fingerenum', params)
    },
    //指纹列表
    getObjectData(params) {
        return axios.get('/smart/tools/fingerlist', params)
    },
    //获取详情
    getInfo(params) {
        return axios.get('/smart/tools/fingerinfo', params)
    },
    //验证里的验证方式下拉
    getvulrisk(params) {
        return axios.get('/tools/fingerprint/verify/mode/select/', params)
    },
    //验证指纹
    checkFinger(params) {
        return axios.post('/tools/fingerprint/verify/', params)
    },
    //开始验证
    checkFinger(params) {
        return axios.get('/tools/fingerprint/verify/result/', params)
    },
    //编辑
    toEdit(params) {
        return axios.get('/tools/fingerprint/detail/', params)
    },
    //正则指纹类型下拉
    openrule(params) {
        return axios.get('/tools/fingerprint/re/rule/type/select/', params)
    },
    // 创建新指纹
    createNewfinger(params) {
        return axios.post('/smart/tools/addfinger', params)
    },
    // 更新新指纹
    updateNewfinger(params) {
        return axios.post('/smart/tools/editfinger', params)
    },
    // 指纹测试
    testfinger(params){
        return axios.postJson('/smart/tools/testfinger', params)
    },
    // 指纹测试日志
    testfingerlog(params){ 
        return axios.get('/smart/tools/testfingerresult', params)
    }
}
const vulnerability = {
    //漏洞对象统计
    getObjectStatistics(params) {
        return axios.get('/tools/vul/vul_class_statistics/', params)
    }, 
    //所有漏洞下拉
    getVulObjectlist(params) {
        return axios.get('/smart/tools/vulenum', params)
    }, 
    //漏洞页列表数据显示
    getObjectData(params) {
        return axios.get('/smart/tools/vullist', params)
    }, 
    //单个漏洞详情
    handleInfo(params) {
        return axios.post('/smart/tools/vuledit', params)
    },
    //详情漏洞风险下拉
    getvulrisk(params) {
        return axios.get('/smart/tools/vulenum', params)
    },  
    //单个漏洞详情编辑保存功能
    saveEditvul(params) {
        return axios.post('/smart/tools/vuledit', params)
    }, 
    // 删除...........................新增................................
    deletevul(params) {
        return axios.delete('/tools/vul/delete/', params)
    },
    //启用
    openVul(params) {
        return axios.post('/smart/tools/vuleditstatus', params)
    },
    //验证
    checkVulP(params) {
        return axios.post('/tools/vul/verify/', params)
    },
    //验证日志接口
    verifyLog(params) {
        return axios.get('/tools/vul/verifyLog/', params)
    },
    //导入VulKit漏洞脚本（zip压缩包）
    importVulnVulKit(data) {
        // 直接传FormData，不要手动设置Content-Type，让浏览器自动处理boundary
        return $ajax.post('/smart/tools/importvulnvulkit', data)
    }
  
}
const dictionary = {
    // 目标删除
    handleDel(params) {
        // return axios.delete('/dictionary/keyvalue/multiple_delete/', params)
        return axios.get('/smart/tools/dictdel', params)
    },
    //单个漏洞详情编辑保存功能
    handleEdit(params) {
        return axios.post('/tools/vul/detail/edit/', params)
    }, 
    //新增保存
    handleCreateSave(params) {
        // return axios.post('/dictionary/keyvalue/', params)
        return axios.post('/smart/tools/dictaddoredit', params)
        
    }, 
    //获取展示列表数据
    getMapList(params) {
        // return axios.get('/dictionary/keyvalue/', params)
        return axios.get('/smart/tools/dictlist', params)
    },
    //获取需要编辑的字典数据
    getDetailData(params) {
        // return axios.get('/dictionary/keyvalue/detail/', params)
        return axios.get('/smart/tools/dictinfo', params)
    },
    // 获取下拉数据
    getServiceSelect(params) {
        // return axios.get('/dictionary/keyvalue/params/', params)
        return axios.get('/smart/tools/enum', params)
    },
    //编辑保存
    handleEditSave(params) {
        // return axios.post('/dictionary/keyvalue/edit/dictionary/', params)
        return axios.post('/smart/tools/dictionaryaddoredit', params)
    }, 
    // 默认字典
    defaultData(params){
        // return axios.post('/tools/dictionary/set_default', params)
        return axios.get('/smart/tools/dictsetdefault', params)
    }
}
const auxiliarytool = {
    //gettags
    gettags(params) {
        return axios.get('/tools/interfaces/tag/list/', params)
    }, 
    //getHttplogToken
    getHttplogToken(params) {
        return axios.get('/tools/assists/httplog_token/', params)
    }, 
    //httplog列表页
    getHttploglist(params) {
        return axios.get('/smart/tools/supportiphostlist', params)
    },
    // 删除按钮事件
    handleOperationDel(params) {
        return axios.delete('/tools/assists/penetration_resources/delete/', params)
    },
    //渗透资源
    getPenetrationResources(params) {
        return axios.get('/tools/assists/penetration_resources/get/', params)
    },
    // 确认清空
    submitForm() {
        return axios.post('/tools/http_log/truncate')
    },
    //保存
    savescript(params) {
        return axios.post('/tools/assists/script/', params)
    },
    // iplog..........................................................
    //iplog列表数据显示
    getiplogData(params) {
        return axios.get('/smart/tools/supportiphostlist', params)
    }, 
    // iplog删除
    deleteipLog(params) {
        return axios.post('/smart/tools/supportiphostdel', params)
    },
    addIplog(params) {
        return axios.post('/smart/tools/supportiphostcreate', params)
    },
    updateIplog(params) {
        return axios.post('/smart/tools/supportiphostcreate', params)
    },
    // Ping..........................................................
     // 停止ping
    stopPing(params) {
        return axios.post('/smart/tools/supportpingstop', params)
    }, 
    // 开始ping
    openPing(params) {
        return axios.post('/smart/tools/supportpingcreate', params)
    },
    // ping页面
    logPing(params) {
        return axios.get('/smart/tools/supportpingresult', params)
    },
    // Traceroute.........................................................
     // 停止Traceroute
     stopTraceroute(params) {
        return axios.post('/smart/tools/supporttraceroutestop', params)
    }, 
    // 开始Traceroute
    openTraceroute(params) {
        return axios.post('/smart/tools/supporttraceroutecreate', params)
    },
    // Traceroute页面
    logTraceroute(params) {
        return axios.get('/smart/tools/supporttracerouteresult', params)
    },

     // ---------------------工具库-----------------------------------
     toollist(params){
        return axios.get('/smart/tools/toolfilelist', params)
    },
    // 反连服务器
    //获得状态
    reversestatus(){
        return axios.get('/smart/reverse/status')
    },
    //开启
    reversestart(params){
         return axios.post('/smart/reverse/start',params)
    },
    reversestop(){
        return axios.post('/smart/reverse/stop')
    },
    reversemessage(params){
        return axios.get('/smart/reverse/message',params)
    },
    reverseclear(){
        return axios.post('/smart/reverse/clear')
    },

}
export {
    fingerprint,
    vulnerability,
    dictionary,
    auxiliarytool
};