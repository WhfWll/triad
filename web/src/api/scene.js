import http from '@/axios/http'

const scene = { 
    getData(params){
        return http.get('/smart/scene/list', params)
    },
    //单个/批量删除
    delScene(params) {
        return http.post('/smart/scene/del', params)
    },
    //默认场景
    defaultScene(params){ 
        return http.post('/smart/scene/setdefault', params)
    },
    //复制场景
    copyScene(params){ 
        return http.post('/smart/scene/copy', params)
    },
    // 新建场景页面
    getSceneEnum(){
        return http.get('/smart/scene/enum')
    }, 
    // 获得漏洞列表
    getVulnlist(params){
        return http.post('/smart/scene/vullist', params)
    },
    getSceneinfo(params){ //场景详情
        return http.get('/smart/scene/info', params)
    }, 
    checkVuln(params){ //检测关联漏洞
        return http.post('/smart/scene/tasktemplateorphan', params)
    }, 
    // 创建或编辑场景
    saveScene(params){
        return http.post('/smart/scene/save', params)
    },
    // 创建渗透任务，选择场景下拉框
    sceneoptions(params){
        return http.get('/smart/scene/sceneoptions', params)
    },
    // 漏洞枚举
    vulnEnum(){
        return http.get('/smart/tools/vulenum')
    },

}

export  default  scene