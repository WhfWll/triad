
/**
 * 社会工程 钓鱼活动
 */

import axios from '@/axios/http';

const gophish = {
    campaignall(param){
        return axios.get('/smart/gophish/campaign/all',param) 
    },
    campaigncreate(param){
        return axios.postJson('/smart/gophish/campaign/create',param) 
    },
    campaignresult(param){
        return axios.get('/smart/gophish/campaign/result',param) 
    },
    campaigndetail(param){
        return axios.get('/smart/gophish/campaign/detail',param) 
    },
    campaigncomplete(param){
        return axios.postJson('/smart/gophish/campaign/complete',param) 
    },
    campaigndelete(param){
        return axios.postJson('/smart/gophish/campaign/delete',param) 
    },
    // 目标组
    groupall(param){
        return axios.get('/smart/gophish/group/all',param) 
    },
    groupcreate(param){
        return axios.postJson('/smart/gophish/group/create',param) 
    },
    groupupdate(param){
        return axios.postJson('/smart/gophish/group/update',param) 
    },
    groupdetail(param){
        return axios.get('/smart/gophish/group/detail',param) 
    },
    groupdelete(param){
        return axios.postJson('/smart/gophish/group/delete',param) 
    },
    //模板
    templateall(param){
        return axios.get('/smart/gophish/template/all',param) 
    },
    templatecreate(param){
        return axios.postJson('/smart/gophish/template/create',param) 
    },
    templateupdate(param){
        return axios.postJson('/smart/gophish/template/update',param) 
    },
    templatedelete(param){
        return axios.postJson('/smart/gophish/template/delete',param) 
    },
    // 钓鱼网站
    pageall(param){
        return axios.get('/smart/gophish/page/all',param) 
    },
    pagecreate(param){
        return axios.postJson('/smart/gophish/page/create',param) 
    },
    pageupdate(param){
        return axios.postJson('/smart/gophish/page/update',param) 
    },
    pagedelete(param){
        return axios.postJson('/smart/gophish/page/delete',param) 
    },
    pageimport_site(param){
        return axios.postJson('/smart/gophish/page/import_site',param) 
    },
    //发件配置
    profileall(param){
        return axios.get('/smart/gophish/profile/all',param) 
    },
    profilecreate(param){
        return axios.postJson('/smart/gophish/profile/create',param) 
    },
    profileupdate(param){
        return axios.postJson('/smart/gophish/profile/update',param) 
    },
    profiledelete(param){
        return axios.postJson('/smart/gophish/profile/delete',param) 
    },
    send_test_email(param){
        return axios.postJson('/smart/gophish/profile/send_test_email',param) 
    },
}

export default gophish;