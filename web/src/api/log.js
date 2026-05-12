/* 
   日志 api 
*/
import axios from '@/axios/http';
const log = {
    //保存设置
    saveConfiginfo(params){
        return axios.post('/smart/logs/logbackupconfig', params)
    },
    // 查询日志过期时间
    getexptime(){
        return axios.get('/smart/logs/getexptime')
    },
    // 保存日志过期时间
    setexptime(params){
        return axios.post('/smart/logs/setexptime', params)
    },
    saveCopyinfo(params){
        return axios.post('/smart/logs/logbackupnow', params)
    },
    //恢复
    // liveCopy(params){
    //     return axios.post('/logs/logBackup/recover/', params)
    // },
    //获得列表
    getlogs(params){
        return axios.get('/smart/logs/logbackuplist', params) 
    },
    getselectinfo(){
        return axios.get('/logs/logBackup/params/') 
    },
    getconfiginfo(){
        return axios.get('/smart/logs/logbackupconfiginfo') 
    },
    handleDelinfo(params){
        return axios.post('/smart/logs/logbackupdelete', params)
    },
    // 下载
    download (params) {
        return axios.get('/smart/logs/logbackupdownload', params) 
    },
    logEmnu(){
        return axios.get('/smart/logs/enum')
    },
}
export default log;