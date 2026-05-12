<template>
    <div class="commandCount">
        <div class="result">
            <p>时间：{{ currentTime }}</p>
            <div v-for="(item,index) in lists" :key="index">
                <span>{{ item.key }}</span>
                <pre class="showResult">{{ item.result }}</pre>
            </div>
             
        </div>
        <div class="shellCount">
            shell ><el-input class="inputCommand" id="inputCommand" v-model="inputCmd" @keyup.enter.native="onSubmitInput"></el-input>
        </div>
    </div>
</template>
<style lang="less" scoped>
.commandCount{
    margin-top:16px;
    width: 100%;
    height: 240px;
    background: #FFFFFF;
    border-radius: 4px;
    border: 1px solid #E8E8F5;
    .result{
        display:block;
        width:100%;
        height:159px;
        overflow-y:auto ;
        border-bottom:1px solid #E8E8F5;
        box-sizing:border-box;
        padding:18px 24px;
        .showResult{
            padding-top:8px;
            // height:24px;
            // line-height:24px;
        }
        span{
            display:block;
            padding-top:8px;
        }
    }
    .shellCount{
        padding:24px;
        .inputCommand{
            // margin:24px 24px 24px 0;
            border-color:rgba(255,255,255,0);
            box-sizing:border-box;
            display:inline-block;
            width:90%;
        }
        /deep/ .el-input__inner{
            border:none !important; 
            padding-left:0;
            width:90%;
        }
    }
}
</style>
<script>
import { number } from "echarts";

export default({
    name:'command',
    props:{
        resultId:Number,
        cmd:String,
        SelrowID:Number,
    },
    data(){
        return{
            inputCmd:'',
            currentTime:'',
            lists:[
                // {'key':'hsakhf','result':'sddsgds'}
            ], 
        }
    },
    mounted(){
        
        this.getFocus();
        this.getTime();
    },
    methods:{
        getFocus(){
            this.$nextTick(() =>{
               document.getElementById('inputCommand').focus();
            })
        },
        onSubmitInput(){//执行命令  
            this.lists = [];
            
            ///task/risk/command/execute/
            this.$ajax({
                // method:'get',
                method:'post',
                url: '/smart/task/exceshell',
                params: {
                    id:this.SelrowID,
                    cmd:this.inputCmd,
                  
                } 
            })
            .then(res =>{
                let dt = res.data; 
                if(dt.code == 200){    
                    this.lists.push(dt.data)
                    this.inputCmd='';
                }else{
                    this.$message({
                        message:dt.msg,
                        type: 'error'
                    });
                }
            }).catch(err=>{})
        },
        // 其他页面传入
        addResult (data) {
            this.lists.push(data)
        },
        getTime(){
            var date = new Date();
			
			var year = date.getFullYear();        //年 ,从 Date 对象以四位数字返回年份
			var month = date.getMonth() + 1;      //月 ,从 Date 对象返回月份 (0 ~ 11) ,date.getMonth()比实际月份少 1 个月
			var day = date.getDate();             //日 ,从 Date 对象返回一个月中的某一天 (1 ~ 31)
			
			var hours = date.getHours();          //小时 ,返回 Date 对象的小时 (0 ~ 23)
			var minutes = date.getMinutes();      //分钟 ,返回 Date 对象的分钟 (0 ~ 59)
			var seconds = date.getSeconds();      //秒 ,返回 Date 对象的秒数 (0 ~ 59)   
			
			//获取当前系统时间  
			var currentDate = year + "-" + month + "-" + day + " " + hours + ":" + minutes + ":" + seconds;
			// alert(currentDate);
			
			//修改月份格式
			if (month >= 1 && month <= 9) {
		        month = "0" + month;
		        }
			
			//修改日期格式
			if (day >= 0 && day <= 9) {
		        day = "0" + day;
		        }
			
			//修改小时格式
			if (hours >= 0 && hours <= 9) {
		        hours = "0" + hours;
		        }
			
			//修改分钟格式
			if (minutes >= 0 && minutes <= 9) {
		        minutes = "0" + minutes;
		        }
			
			//修改秒格式
			if (seconds >= 0 && seconds <= 9) {
		        seconds = "0" + seconds;
		        }
			
			//获取当前系统时间  格式(yyyy-mm-dd hh:mm:ss)
			this.currentTime = year + "-" + month + "-" + day + " " + hours + ":" + minutes + ":" + seconds;
        }
    },
    watch:{
        cmd(val){
            this.inputCmd = val;
        }
    }
})
</script>