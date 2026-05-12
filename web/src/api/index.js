import axios from '@/axios/http'
 

const index = {
    //任务统计
    gettaskinfostat(params){
        return axios.get('/smart/homepage/taskinfostat',params) 
    },
    //目标统计
    gettargetriskstat(){
        return axios.get('/smart/homepage/targetriskstat') 
    },
    // 漏洞类型统计
    getvultypestat(params){
        return axios.get('/smart/homepage/vultypestat',params) 
    },
    // 漏洞发现趋势
    getvulfindtrendstat(params){
        return axios.get('/smart/homepage/vulfindtrendstat',params) 
    },
    //工具统计
    gettoolinfostat(){
        return axios.get('/smart/homepage/toolinfostat') 
    },
    // 漏洞统计
    gettaskvulriskstat(){
        return axios.get('/smart/homepage/taskvulriskstat') 
    },
    //漏洞取证
    getvulevidencestat(){
        return axios.get('/smart/homepage/vulevidencestat') 
    },
    //最新消息
    getmessagestat(){
        return axios.get('/smart/homepage/messagestat') 
    },

}

export  default  index