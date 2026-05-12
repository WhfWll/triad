import http from '@/axios/http'

const scene = { 
  // 信息收集列表查询及筛选(完成)
  
    getData(params){
        return http.get('/smart/task/taskresultlist', params)
    },
    getSceneinfo(params){ //场景详情
      return http.get('/smart/task/taskresultdetail', params)
    }, 
    //单个/批量删除
    delScene(params) {
        return http.get('/smart/task/taskresultdel', params)
    },
    //默认场景


}

export  default  scene