/*
    验证报告
*/
import axios from '@/axios/http';
const verReport = {
    //报错验证报告
    saveUploadreport(params){
        return axios.postFormData('/smart/reportverify/upload', params)
    },
   
}
export default verReport;