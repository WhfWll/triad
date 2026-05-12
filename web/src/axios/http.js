/**
* http.js
* 封装axios， 
*/
import axios from './axios.js'
import QS from 'qs';
 

export default {
    /**
     * get方法，对应get请求
     * @param {String} url [请求的url地址]
     * @param {Object} params [请求时携带的参数]
     */
    get(url, params) {
        return new Promise((resolve, reject) => {
            axios.get(url, {
                params: params
            }).then(res => {
                resolve(res.data)
            }).catch(err => {
                reject(err)
            })
        })
    },
    getfile(url, type, params) {
        return new Promise((resolve, reject) => {
            axios({
                url,
                method: 'get',
                responseType: type,
                params: params
            }).then(res => {
                resolve(res.data)
            }).catch(err => {
                reject(err)
            });
        })
    }, 

    /**
     * post方法，对应post请求
     * @param {String} url [请求的url地址]
     * @param {Object} params [请求时携带的参数]
     */
    post(url, params) {
        return new Promise((resolve, reject) => { 
        axios.post(url, QS.stringify(params),)
                .then(res => {
                    resolve(res.data)
                })
                .catch(err => {
                    reject(err)
                })
        }) 
    },
    postJson(url, params) {
        return new Promise((resolve, reject) => { 
        axios.post(url, params ,{
            headers: {
                'content-type': 'application/json;charset=utf-8'
            }
        })
        .then(res => {
            resolve(res.data)
        })
        .catch(err => {
            reject(err)
        })
        }) 
    },
    
    /**
    * patch方法，对应patch请求
    * @param {String} url [请求的url地址]
    * @param {Object} params [请求时携带的参数]
    */
    patch(url, params) {
        return new Promise((resolve, reject) => {
            axios.patch(url, QS.stringify(params))
                .then(response => {
                    resolve(response.data);
                }, err => {
                    reject(err)
                })
        })
    },
    /**
    * put方法，对应put请求
    * @param {String} url [请求的url地址]
    * @param {Object} params [请求时携带的参数]
    */
    put(url, params) {
        return new Promise((resolve, reject) => {
            axios.put(url, QS.stringify(params))
                .then(response => {
                    resolve(response.data);
                }, err => {
                    reject(err)
                })
        })
    },
    /**
    * delete方法，对应delete请求
    * @param {String} url [请求的url地址]
    * @param {Object} params [请求时携带的参数]
    */
    delete(url, params) {
        return new Promise((resolve, reject) => {
            axios.delete(url, {
                    data: params
                })
                .then(response => {
                    resolve(response.data);
                }, err => {
                    reject(err)
                })

            // axios({
            //     method: 'delete',
            //     url: url,
            //     data: params 
            // }).then(response => {
            //     resolve(response.data); 
            // }, err => {
            //     reject(err)
            // })

        })
    },
    /**
     * postFormData方法，对应post请求，用来提交文件+数据
     * @param {String} url [请求的url地址]
     * @param {Object} params [请求时携带的参数],不用转换成formdata
     */
    postFormData(url, params) {
        return new Promise((resolve, reject) => {
            axios({
                headers: {
                    'Content-Type': 'multipart/form-data'// ;boundary=----WebKitFormBoundaryQ6d2Qh69dv9wad2u
                },
                transformRequest: [function (data) { // 在请求之前对data传参进行格式转换
                    const formData = new FormData()
                    Object.keys(data).forEach(key => {
                        formData.append(key, data[key])
                    })
                    return formData
                }],
                url,
                method: 'post',
                data: params
            }).then(res => {
                resolve(res.data)
            }).catch(err => {
                reject(err)
            });

        })


        
    },
    /* 
    * postBlob方法，对应post请求，传入responseType:blob
    * 
    */
    postBlob(url, params){
        return new Promise((resolve, reject) => {
            axios({ 
                url,
                method: 'post',
                responseType: 'blob', 
                data: params
            }).then(res => {
                resolve(res.data)
            }).catch(err => {
                reject(err)
            });



            // axios.post(url, QS.stringify(params))
            //     .then(res => {
            //         resolve(res.data)
            //     })
            //     .catch(err => {
            //         reject(err)
            //     })
            //  }) 
        })
    }
}