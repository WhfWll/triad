<!--Ping............................................................................................................. -->
<template>
    <div class="TracerouteDiv context_box_bg"> 
        <div  >
            <div class="btnbox">
               
                <el-button size="mini" type="primary" :disabled="working || !form.search_field"   @click="startTraceroute()">开始</el-button>
                <el-button size="mini" type="warning"  :disabled="!working"  @click="cancelTraceroute()">停止</el-button>
            </div>
            <div class="search-text">
                <el-form ref="form" :model="form" :rules="rule" label-width="20px">
                    <el-form-item label="" prop="search_field">
                      <el-input placeholder="请输入请输入服务器地址/域名地址"  size="small" clearable class="input" v-model="form.search_field"></el-input>
                    </el-form-item>
                </el-form>
                <!-- <i class="iconfont icontishi icontsstyle"></i> -->
            </div>
        </div>
        <div class="pinginfo">
            <ul class="loglist">
                <li v-for="(item,index) in TracerouteData" :key="index">
                    <span>
                       <pre>{{item}}</pre>
                    </span>
                </li>
            </ul>
        </div>
    </div>
</template>
<style lang="less" scoped>
    .loglist li{
        list-style: none;
    }
    .loglist li span {
        color:rgba(72, 72, 102, 0.64);
        font-size: 13px;    
}
  .btnbox{
      display: inline-block;
      font-size: 14px;
  }
  .search-text{
      margin-left: 10px;
      display: inline-block;
      width:80%;
  }
  .icontsstyle{
      display:inline-block;
      margin-right:10px;
  }
  .input{
      width:90%;
      display:inline-block
  }
  /deep/ .el-button{
      font-size:14px!important;
  }
  .pinginfo{
    height: calc(100vh - 300px);
    overflow-y: auto;
  }
</style>
<script>  
import About from "@/components/About.vue";
import Operation from "@/components/Operation";
import { auxiliarytool } from '@/api/tool.js'
export default ({
    name:'Traceroute',
    components:{
        About,
        Operation
    },
    data(){ 
        var validatePass3 = (rule, value, callback) => { 
            if (!value) {
              callback(new Error('请输入服务器地址'));
            } else {
                const re = /^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])\.((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$/
                const re2 = /^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,6}$/
                if (re.test(value) || re2.test(value)) {
                    callback();
                } else {
                    return callback(new Error('请输入正确的服务器地址/域名地址'));
                }
            }
        }; 
    	return{ 

            TracerouteData:[],
            taskId: '',
            form: {
                search_field:'',
            },
            rule: {
                search_field:[
                     { required: true, message: '服务器地址/域名地址不能为空', trigger: 'blur' },
                     { validator: validatePass3, trigger: 'blur' },
                ], 
            },
            timer: null,
            working: false
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/auxiliarytool"; 
    },
    mounted:function(){   

    },
    beforeDestroy () {
        clearTimeout(this.timer)
    },
    methods:{ 
        startTraceroute(){
            this.$refs.form.validate(async valid=>{
                if (valid) {
                    let params = {
                        ip:this.form.search_field,
                    }
                    const res = await auxiliarytool.openTraceroute(params)
                    if(res.code == 200){
                        this.working = true
                        // this.taskId = res.data.task_id
                        this.taskId = res.data.token
                        this.getTraceroutelog()
                        this.$message({
                            message:res.msg||'任务开始成功',
                            type: 'success'
                        });
                    }else{
                        this.$message({
                            message:res.error,
                            type: 'error'
                        });
                    } 
                }
            })
            
        },
        async getTraceroutelog(){
            let that = this
            let params = {
                token: that.taskId|| localStorage.getItem('token')
            }
            const res = await auxiliarytool.logTraceroute(params)
            if(res.code == 200){
                // this.pingData.push(res.data)
                // this.TracerouteData = res.data.unreadMsg
             if(res.data&&res.data.unreadMsg.length>0){
                     res.data.unreadMsg.forEach(element => {
                        this.TracerouteData.push(element)
                     });
                }
                // if (this.TracerouteData) {
                    this.timer = setTimeout(() => {
                        this.getTraceroutelog()
                    }, 3000)
                // }else {
                //     this.working = false
                // }

            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            } 
        }, 
        async cancelTraceroute(){
            let that = this
            let params = {
                token: that.taskId|| localStorage.getItem('token')
            }
            const res = await auxiliarytool.stopTraceroute(params)
            if(res.code == 200){ 
                this.$message({
                    message:res.msg||'任务已关闭',
                    type: 'success'
                });
                clearTimeout(this.timer)
                this.working = false
                // 停止清空值
                this.taskId = ''
                // this.form.search_field = ''
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            } 
        }, 
    }
})
 
</script>
