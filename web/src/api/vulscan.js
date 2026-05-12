
import http from '@/axios/http'
 

const vulscan = {
    taskList(params){ //任务列表
        return http.get('/smart/vulscan/tasklist', params)
    },
    taskSave(params){ //任务保存
        return http.postJson('/smart/vulscan/tasksave', params)
    },
    targetList(params){ //目标列表
        return http.get('/smart/vulscan/targetlist', params)
    },
    taskDelete(params){ //任务删除
        return http.get('/smart/vulscan/taskdelete', params)
    },
    taskStop(params){ //任务结束
        return http.get('/smart/vulscan/taskstop', params)
    },
    taskOverview(params){ //cve列表
        return http.get('/smart/vulscan/taskoverview', params)
    },
    taskState(params){ //cve列表
        return http.get('/smart/vulscan/taskstate', params)
    },
    vulList(params){ //漏洞列表
        return http.get('/smart/vulscan/vullist', params)
    },
    vulDetail(params){ //漏洞详情
        return http.get('/smart/vulscan/vuldetail', params)
    },
    cveList(params){ //cve列表
        return http.get('/smart/vulscan/cvelist', params)
    },
    cveDetail(params){ //cve列表
        return http.get('/smart/vulscan/cvedetail', params)
    },
};


export  {
    vulscan,
}
