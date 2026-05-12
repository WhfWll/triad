
import axios from 'axios';
import router from '../router';
import { Message } from 'element-ui';

// import { encryptCBC, decryptCBC } from '../commonFunction/des.js'

let URL_PREFIX = window.location.protocol + "//" + window.location.host;
// let URL_PREFIX = window.location.protocol + "//" + window.location.hostname + ":8011";
// let URL_PREFIX  = '';
if (process.env.NODE_ENV === 'development') {
    //开发环境下的代理地址，解决本地跨域跨域，配置在config目录下的index.js dev.proxyTable中
    URL_PREFIX = "/api"
    // URL_PREFIX ='http://192.168.0.61:80';
} else {
    //生产环境下的地址
    URL_PREFIX += '/api';
}

const $ajax = axios.create({
    baseURL: URL_PREFIX
});

$ajax.interceptors.request.use(
    config => {
        //判断 token 
        //是否登录  
        // var user = decryptCBC(localStorage.getItem('user'), '4dogs.cn');
        // if (!user) {
        //     // this.$router.push('/login');
        //     router.push('/login');
        // } 
        let token = localStorage.getItem('token')
        if (token) {
            config.headers['Authorization'] = token
        } else {
            router.push('/login');
        }
        return config
    },
    err => {
        return Promise.reject(err)
    }
)
$ajax.interceptors.response.use(function (response) {
    if(response.data.code == 401){
        // Message.error('登录超时，请重新登录');
        router.push('/login');
    }
    if(response.data.code == 403){
        //  没权限
        router.push('/404');
    }
    let token = response?.headers['authorization']
    if (!token) {
        token = response?.headers['Authorization']
    }
    // console.log('token', token)
    if (token) {
        // console.log('取到go的token了')
        localStorage.setItem('token', token)
    } else {
        // console.log('token没取到')
    }
    if (response.headers && (
        response.headers['content-disposition'] === 'application/octet-stream' ||
        response.headers['content-disposition'] === 'application/octet-stream;charset=utf-8' ||
        response.headers['content-disposition'] === 'multipart/form-data' ||
        (response.headers['content-disposition'] && response.headers['content-disposition'].includes('attachment;filename=')) ||
        (response.headers['content-disposition'] && response.headers['content-disposition'].includes('attachment; filename=')) ||
        (response.headers['Content-Disposition'] && response.headers['Content-Disposition'].includes('attachment;filename=')) ||
        (response.headers['Content-Disposition'] && response.headers['Content-Disposition'].includes('attachment; filename='))
        
    )
    ) {
        console.log('response', response)
        let filename = response?.headers['content-disposition']
        console.log('filename', filename)
        if (!filename) {
            filename = response?.headers['Content-Disposition']
        }
        sessionStorage.filename = filename.split('filename=')[1]
    }

    return response;
}, function (error) {
    console.log('error', error)
    console.log('error2', error.response)
    console.log('error3', error.response.status)
    if (error && error.response) {
        switch (error.response.status) {
            case 400:
                // Message.error('错误请求');
                Message.error(error.response.data.error);
                break;
            case 401:
                // Message.error('访问被拒绝');
                router.push('/404');
                break;
            case 403:
                // Message.error('禁止访问');
                // Message.error(error.response.data.msg);
                router.push('/login');
                break;
            case 404:
                Message.error('请求错误，未找到资源');
                router.push('/404');
                break;
            case 408:
                Message.error('请求超时');
                break;
            case 500:
                //Message.error('服务器端出错');
                // router.push('/500');
                break;
            case 502:
                Message.error('网络错误');
                break;
            case 503:
                Message.error('服务不可用');
                break;
            case 504:
                Message.error('网关超时');
                break;
            default:
                Message.error('连接错误' + error.response.status);

        }
    } else {
        Message.error('网络出现问题，请稍后再试');
    }


    return Promise.reject(error);
});

export default $ajax