
import http from '@/axios/http'

export default{
     
    riskvulstatistics(){ //风险统计
        return http.get('/smart/riskvul/statistics');
    },
    riskvulnums(){ 
        return http.get('/smart/riskvul/enums');
    },
    riskvullist(data){
        return http.get('/smart/riskvul/list',data);
    },
    vulinfo(data){
        return http.get('/smart/riskvul/detail',data);
    },
    vuldelbyid(data){
        return http.get('/smart/riskvul/delete',data);
    },
    riskvulnStatus(data){
        return http.get('/smart/riskvul/updatestatus',data);
    },
 
}
