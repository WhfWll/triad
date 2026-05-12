/* 

  流量分析
*/
import axios from '@/axios/http';
const traffic = {
    // 枚举
    trafficEnum(){
        return axios.get('/smart/task/flowtaskenum',)
    },
    getStatus(params){ //获得状态
        return axios.get('/smart/task/flowtaskstatus',params)
    },
    getData(params){ //列表
        return axios.get('/smart/task/flowtasklist', params)
    },
    addTraffic(params){ //新建
        return axios.postJson('/smart/task/flowtaskadd', params) 
    },
    delTraffic(params){ //删除
        return axios.get('/smart/task/flowtaskdel', params)
    },
    trafficStatus(params){ //结束
        return axios.get('/smart/task/changeflowtaskstatus', params)
    }, 
    trafficbase(params){//基本详情
        return axios.get('/smart/task/flowtaskinfo', params)
    },
    downloadhttps(params){ //下载证书
        return axios.get('/smart/task/httpscert', params)
    },
    trafficlist(params){ //被动流量列表
        return axios.get('/smart/task/flowbaselist', params)
    },
    trafficlistinfodel(params){
        return axios.post('/smart/task/flowbasedel', params)
    },
    gettrafficHeader(params){
        return axios.get('/smart/task/flowbaseinfo', params)
    },
    trafficVuln(params){ //漏洞
        return axios.get('/smart/task/flowrisklist', params)
    },
    trafficVulnInfo(params) { //漏洞详情
        return axios.get('/smart/task/flowriskinfo', params)
    },
    trafficvulndel(params){
        return axios.post('/smart/task/flowriskdel', params) 
    },
    trafficlog(params){
        return axios.get('/smart/task/flowloginfo', params)
    },
    clearTrafficlog(params){ //清空log
        return axios.get('/smart/task/flowlogdel', params)
    },
    flowtaskedit(params){
        return axios.post('/smart/task/flowtaskedit', params) 
    }
}

// export default traffic;
export {
    traffic
}