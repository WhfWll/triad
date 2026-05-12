import http from '@/axios/http';
const report = {
    reportEnum(){
        return http.get('/smart/report/enum')
    },
    reportList(params){
        return http.get('/smart/report/list', params)
    },
    createReport(params){
        return http.postJson('/smart/report/save', params)
    },
    reportDel(params){
        return http.get('/smart/report/del', params)
    },
    downLoadfile(params){
        return http.get('/smart/report/download', params)
    },
    // 临时接口
    downLoadreporttemp(params){
        return http.get('/smart/report/generate', params)
    },

}

export  {
    report,
}
