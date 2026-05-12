
import http from '@/axios/http'
 

const login = {
    captcha(){
        return http.get('/smart/user/logincaptcha',)
    },
    submitLogin(params){
        return http.post('/smart/user/login', params)
    },
    //免密登录
    loginfreepassa(params){
        return http.post('/smart/user/loginfreepassb',params)    
    }
}

export  default  login

