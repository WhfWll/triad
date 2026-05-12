/* 

  资产管理
*/
import axios from '@/axios/http';
const traffic = {
    // 1资产枚举接口(完成)
    trafficEnum(){
        return axios.get('/smart/asset/enum',)
    },
    getStatus(){ //2、资产统计(完成)
        return axios.get('/smart/asset/stat')
    },
    getData(params){ //列表
        return axios.get('/smart/task/flowtasklist', params)
    },
    addTraffic(params){ //3、最新资产风险变化列表（完成）
        return axios.get('/smart/asset/uptodate', params) 
    },
    delTraffic(params){ //4、最新资产风险变化忽略（完成）
        return axios.get('/smart/asset/uptodateignore', params)
    },
    trafficStatus(params){ //5、资产渗透测试趋势（完成）
        return axios.get('/smart/asset/taskstat', params)
    }, 
    trafficbase(params){//6、漏洞发现趋势(完成)
        return axios.get('/smart/asset/vulstat', params)
    },
    downloadhttps(params){ //7、漏洞类型统计（完成）
        return axios.get('/smart/asset/vultypestat', params)
    },
    trafficlist(params){ //8、服务类型统计（完成）
        return axios.get('/smart/asset/servicetypestat', params)
    },
    trafficlistinfodel(params){//9、资产树整体结构高级搜索（完成）
        return axios.post('/smart/asset/tree', params)
    },
    gettrafficHeader(params){//10、获取完整资产组结构（完成）
        return axios.get('/smart/asset/group', params)
    },
    trafficVuln(params){ //11资产组概览
        return axios.get('/smart/assetgroup/overview', params)
    },
    trafficVulnInfo(params) { //12子资产组列表

        return axios.get('/smart/assetgroup/sublist', params)
    },
    trafficvulndel(params){ //13资产组列表
        return axios.get('/smart/assetgroup/list', params) 
    },
    trafficlog(params){ //14资产组列表删除
        return axios.get('/smart/assetgroup/assetdel', params)
    },
    clearTrafficlog2(params){ //15 资产漏洞列表
        return axios.get('/smart/assetgroup/vullist', params)
    },
    clearTrafficlogDel(params){ //16 资产漏洞列表-删除
        return axios.get('/smart/assetgroup/vuldel', params)
    },
    danDuZIChanGL(params){ //17 单独资产资产概览
        return axios.get('/smart/asset/overview', params)
    },
    danDuZIChanmanageinfo(params){ //18 资产概览-管理信息详情查询
        return axios.get('/smart/asset/manageinfo', params)
    },
    danDuZIChanportlist(params){ //19 资产端口列表222
        return axios.get('/smart/asset/portlist', params)
    },
    danDuZIChanmanageinfoupdate(params){ //20 更新资产详情信息
        return axios.post('/smart/asset/manageinfoupdate', params)
    },
    danDuZIChanvulinfo(params){ //21 漏洞详情：
        return axios.get('/api/smart/task/vulinfo', params)
    },
    danDuZIChanvulenums(params){ //22漏洞状态枚举接口
        return axios.get('/smart/asset/vulenums', params)
    },
    danDuZIChanvulinfoeditstatus(params){ //23改漏洞状态：
        return axios.get('/smart/asset/vulinfoeditstatus', params)
    },
    danDuZIChanvultestinfodel(params){ //24删除漏洞
        return axios.get('/smart/asset/vultestinfodel', params)
    },
    danDuZIChanvultestienum(params){ //25其他枚举接口：
        return axios.get('/api/smart/asset/enum', params)
        },
    clearTrafficlog(params){ //清空log
        return axios.get('/smart/task/flowlogdel', params)
    },
    
    assetfinddiff(){ //发现资产-数据展示(完成)
        return axios.get('/smart/asset/assetfinddiff')
    },
    syncassetfinddiff(params){ //同步发现的资产(完成)
        return axios.post('/smart/asset/assetfindsync', params)
    },
    assetpenetrationsync(params){ //资产渗透
        return axios.post('/smart/asset/assetpenetrationsync', params)
    },
    assetLDTestList(params){ //漏洞测试列表
        return axios.get('/smart/asset/vultestlist', params)
    },
    downloadasset(params){ //资产目标下载
        return axios.get('/smart/asset/templatedownload',params)
    },
    importasset(params){ //资产导入
        return axios.post('/smart/asset/import',params)
    },
    assetDetect(params){ //开始资产扫描
        return axios.post('/smart/asset/detect',params)
    },
    addasset(params,type){ //添加资产 ,type:0 编辑，1 添加
        if(type ==1){
            return axios.post('/smart/asset/add', params)
        }
        else {
            return axios.post('/smart/asset/edit', params)
        }
    },
    addassetgroup(params,type){ //添加资产组 type:0 编辑，1 添加
        if(type ==1){
            return axios.post('/smart/assetgroup/add', params)
        }else{
            return axios.post('/smart/assetgroup/edit', params)
        }
    },
      
    assetdelete(params){ //删除 资产-资产组
        return axios.post('/smart/asset/assetdelete', params)
    },
    assetGroupTree(){ //只有资产组树形结构
        return axios.get('/smart/asset/group')
    }, 
    assetDetail(params){ //资产详情
        return axios.get('/smart/asset/detail',params)
    },
    assetGroupDetail(params){ //资产组详情
        return axios.get('/smart/assetgroup/detail',params)
    },
    assetDengJiEnum(){ //备案等级枚举
        return axios.get('/smart/asset/enum') 
    },
    assetDetectProgress(){ //获取资产扫描进度
        return axios.get('/smart/asset/progress') 
    },
}

// export default assetManagement;
export {
    traffic
}