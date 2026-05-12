<template>
    <div>
        <el-dialog
            :title="dialogtitle"
            :visible.sync="visible"
            width="1184px"
            class="buginfobox fingerDialog" 
            :close-on-click-modal="false" 
            :show-close="false" >
            <div class="dialog_b_btn">   
                <el-button size="small" @click="btnStart" >开始测试</el-button>
				<el-button size="small" @click="handleClose">关闭</el-button>
            </div>
            <div class="form_box">
                <el-form :model="form" label-width="80px"  status-icon  ref="" >
                    <el-form-item
                        class=""
                        label="指纹名称"
                        prop="name" > 
                        <el-input 
                            v-model="form.name"  
                            autocomplete="off" 
                            disabled
                            placeholder="请输入指纹名称" ></el-input> 
                    </el-form-item>  
                    <el-form-item
                        class=""
                        label="指纹规则"
                        prop="rule" > 
                        <el-input 
                            type="textarea"
                            resize="none"
                            rows='3'
                            v-model="form.rule"  
                            autocomplete="off" 
                            disabled
                            placeholder="请输入指纹规则" ></el-input> 
                    </el-form-item>  
                    <el-form-item
                        class=""
                        label="测试地址"
                        prop="url" > 
                        <el-input 
                            v-model="form.url"  
                            autocomplete="off"  
                            placeholder="" ></el-input> 
                    </el-form-item> 
                    <el-form-item
                        class=""
                        label="测试日志"
                        prop="log" > 
                        <el-input 
                            type="textarea"
                            resize="none"
                            rows='5'
                            v-model="form.log"  
                            autocomplete="off"  
                            placeholder="" ></el-input> 
                    </el-form-item>  
                </el-form>
            </div>
        </el-dialog>
    </div>
</template>
<style scoped  lang="less">
.form_box{
    padding: 24px; 
}
</style>
<script>
import { fingerprint } from '@/api/tool.js'
export default {
    data(){
        return{
            form:{
                name:'',
                rule:'',
                url:'',
                log:'',
            },
            dialogtitle:'测试指纹', 
            resultID:'',
            tiemr :null,
        }
    },
    props: {
        visible: {
            type: Boolean
        },
        fingerInfo:{
            type: Object
        }
    },
    watch: {
        fingerInfo: {
            handler (n, o) { 
                this.form.name = this.fingerInfo.fingerName;
                this.form.rule = this.fingerInfo.rule;
            },
            deep: true // 深度监听父组件传过来对象变化
            
        }
    },
    created(){
        // console.log(this.fingerInfo)
    },
    mounted(){
        // console.log(this.fingerInfo)
        if (this.fingerInfo) { 
            this.form.name = this.fingerInfo.fingerName
            console.log(this.fingerInfo)
        } else {
            // this.newFingerform.id = ''
            // this.newFingerform.name = ''
            // this.newFingerform.type = ''
            // this.newFingerform.re_rule = ''
            // this.zhengzeTable = []
            // this.metaTable = []
            // this.headerTable = []
        }
    },
    methods:{
        async btnStart(){ //开始测试
            let param = [{
                key:'root_url',
                value:this.form.url,
            }];
            const res = await fingerprint.testfinger({
                fingerID:this.fingerInfo.id,
                tool_name: this.form.name ,
                param: param  
            });
            if(res.code == 200){
                console.log(123)
                this.resultID = res.data.token;
                this.getTestlog();

            }
        },
        async getTestlog(){
            let _this = this;
            const res = await fingerprint.testfingerlog({
                callId:this.resultID
            });
            if(res.code == 200){
                if(res.data.result){
                     _this.form.log += res.data.result;
                } else if (!res.data.end) {
                     _this.form.log += '正在测试中...\n';
                }

                if(res.data.end){ //true,结束
                    clearInterval(this.tiemr);
                    _this.form.log += '\n测试结束\n';
                }else{
                    //继续循环 
                    clearInterval(this.tiemr);
                    this.tiemr = setTimeout(function(){
                        _this.getTestlog(); 
                    },3000);
                    
                }
            }else{
                clearInterval(this.tiemr);
                this.$message({
                    message:'测试地址有误',
                    type: 'error'
                });
            }

        },
        handleClose(){
            clearInterval(this.tiemr);
            this.form.url;
            this.form.log='';
            this.$emit('handleClose')
        },
    },
}
</script>