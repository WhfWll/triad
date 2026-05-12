import http from '@/axios/http'

const asset = {
    // 1资产枚举接口(完成)
    trafficEnum(){
        return http.get('/smart/asset/enum'); 
    },
    getStatus(){ //2、资产统计(完成)
        return http.get('/smart/asset/stat')
    },
    getData(params){ //列表
        return http.get('/smart/task/flowtasklist', params)
    },
    getAllDataIP(params){ //
        return http.get('/smart/asset/selectall', params)
    },
    addTraffic(params){ //3、最新资产风险变化列表（完成）
        return http.get('/smart/asset/uptodate', params)
    },
    delTraffic(params){ //4、最新资产风险变化忽略（完成）
        return http.get('/smart/asset/uptodateignore', params)
    },
    trafficStatus(params){ //5、资产渗透测试趋势（完成）
        return http.get('/smart/asset/taskstat', params)
    },
    trafficbase(params){//6、漏洞发现趋势(完成)
        return http.get('/smart/asset/vulstat', params)
    },
    downloadhttps(params){ //7、漏洞类型统计（完成）
        return http.get('/smart/asset/vultypestat', params)
    },
    trafficlist(params){ //8、服务类型统计（完成）
        return http.get('/smart/asset/servicetypestat', params)
    },
    trafficlistinfodel(params){//9、资产树整体结构高级搜索（完成）
        return http.get('/smart/assettree/list', params)
    },
    gettrafficHeader(params){//10、获取完整资产组结构（完成）
        return http.get('/smart/asset/group', params)
    },
    trafficVuln(params){ //11资产组概览
        return http.get('/smart/assetgroup/overview', params)
    },
    trafficVulnInfo(params) { //12子资产组列表

        return http.get('/smart/assetgroup/sublist', params)
    },
    trafficvulndel(params){ //13资产组列表
        return http.get('/smart/asset/list', params)
    },
    addApi(params){ //13资产树新增
        return http.post('/smart/assetgroup/add', params)
    },
    editApi(params){ //13资产树新增
        return http.post('/smart/assetgroup/edit', params)
    },
    deleteApi(params){ //13资产树新增
        return http.get('/smart/assetgroup/del', params)
    },
    enumApi(){ //13资产树新增
        return http.get('/smart/assetgroup/enums')
    },
    addAsset(params){ //13资产树新增
        return http.post('/smart/asset/add', params)
    },
    editAsset(params){ //13资产树新增
        return http.post('/smart/asset/edit', params)
    },
    deleteAsset(params){ //13资产树新增
        return http.get('/smart/asset/del', params)
    },
    assetinfo(query){ //资产详情
        return http.get('/smart/assetgroup/info', query)
    },
    detailAsset(params){ //13资产树新增
        return http.get('/smart/asset/detail', params)
    },
    importAsset(data){
        return http.postFormData('/smart/asset/import', data)
    },
    exportAsset(params){ //13资产树新增
        return http.post('/smart/asset/export', params)
    },
    updateManyAssetGroups(data){ //13资产树新增
        return http.post('/smart/asset/asset/updatemanyassetgroups', data)
    },
    statisticsVal(){ //资产综述
        return http.get('/smart/asset/statistics')
    },
    collectApi(params){ //备战api
        return http.get('/smart/asset/collect', params)
    },
    collectDetailApi(params){ //备战api
        return http.get('/smart/activity/event/detail', params)
    },
    overviewApi(params){ //迎战阶段 - 迎战综述 - 概述（完成  ee）
        return http.get('/smart/activity/confront/overview', params)
    },


    trafficlog(params){ //14资产组列表删除
        return http.get('/smart/assetgroup/assetdel', params)
    },
    clearTrafficlog2(params){ //15 资产漏洞列表
        return http.get('/smart/assetgroup/vullist', params)
    },
    clearTrafficlogDel(params){ //16 资产漏洞列表-删除
        return http.get('/smart/assetgroup/vuldel', params)
    },
    danDuZIChanGL(params){ //17 单独资产资产概览
        return http.get('/smart/asset/overview', params)
    },
    danDuZIChanmanageinfo(params){ //18 资产概览-管理信息详情查询
        return http.get('/smart/asset/manageinfo', params)
    },
    danDuZIChanportlist(params){ //19 资产端口列表222
        return http.get('/smart/asset/portlist', params)
    },
    danDuZIChanmanageinfoupdate(params){ //20 更新资产详情信息
        return http.post('/smart/asset/manageinfoupdate', params)
    },
    danDuZIChanvulinfo(params){ //21 漏洞详情：
        return http.get('/api/smart/task/vulinfo', params)
    },
    danDuZIChanvulenums(params){ //22漏洞状态枚举接口
        return http.get('/smart/asset/vulenums', params)
    },
    danDuZIChanvulinfoeditstatus(params){ //23改漏洞状态：
        return http.get('/smart/asset/vulinfoeditstatus', params)
    },
    danDuZIChanvultestinfodel(params){ //24删除漏洞
        return http.get('/smart/asset/vultestinfodel', params)
    },
    danDuZIChanvultestienum(params){ //25其他枚举接口：
        return http.get('/api/smart/asset/enum', params)
        },
    clearTrafficlog(params){ //清空log
        return http.get('/smart/task/flowlogdel', params)
    },

    assetfinddiff(){ //发现资产-数据展示(完成)
        return http.get('/smart/asset/assetfinddiff')
    },
    syncassetfinddiff(params){ //同步发现的资产(完成)
        return http.post('/smart/asset/assetfindsync', params)
    },
    assetpenetrationsync(params){ //资产渗透
        return http.post('/smart/asset/assetpenetrationsync', params)
    },
    assetLDTestList(params){ //漏洞测试列表
        return http.get('/smart/asset/vultestlist', params)
    },
    downloadasset(params){ //资产目标下载
        return http.get('/smart/asset/templatedownload',params)
    },
    importasset(params){ //资产导入
        return http.post('/smart/asset/import',params)
    },
    addasset(params,type){ //添加资产 ,type:0 编辑，1 添加
        if(type ==1){
            return http.post('/smart/asset/add', params)
        }
        else {
            return http.post('/smart/asset/edit', params)
        }
    },
    addassetgroup(params,type){ //添加资产组 type:0 编辑，1 添加
        if(type ==1){
            return http.post('/smart/assetgroup/add', params)
        }else{
            return http.post('/smart/assetgroup/edit', params)
        }
    },

    assetdelete(params){ //删除 资产-资产组
        return http.post('/smart/asset/assetdelete', params)
    },
    assetGroupTree(){ //只有资产组树形结构
        return http.get('/smart/asset/group')
    },
    // assetDetail(params){ //资产详情
    //     return http.get('/smart/asset/detail',params)
    // },
    assetDetail(params) {
        const queryString = new URLSearchParams(params).toString();
        return http.get(`/smart/asset/detail?${queryString}`);
      },
    assetGroupDetail(params){ //资产组详情
        return http.get('/smart/asset/group/detail',params)
    },
    assetDengJiEnum(){ //备案等级枚举
        return http.get('/smart/asset/enum')
    },
}
 
export default asset;
 
