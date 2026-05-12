<template>
<!-- 漏洞利用 -->
    <el-dialog :title="useinfo.title" :visible.sync="bugusedialogVisible" width="1184px" :before-close="handleCloseUse"
        :close-on-click-modal="false" :validate-on-rule-change="false" :show-close="false">
        <div class="dialog_b_btn">
            <el-button size="small" @click="handleCloseUse">关闭</el-button>
        </div>
        <div style="padding:24px">
            <div class="targetaddress">目标地址：{{ useinfo.target }}</div>
            <div class="time">时间：{{ useinfo.time }}</div>
            <div class="usebox">
                <div v-for="(item, index) in useresult" :key="index">
                    smart x > <span>{{item.usecmd}}</span>
                    <div class="cmdresult">{{item.useresult}}</div>
                </div>
                <div>
                    smart x > <el-input :autofocus="true" id="useinput" class="useinput" v-model="usecmd"
                        @keyup.enter.native="onSubmitUse" size="mini"></el-input>
                </div>
            </div>
        </div>
    </el-dialog>
</template>
<style lang="less" scoped>

 .useinput {
     width: 90% !important;
     box-sizing: border-box;

     /deep/ .el-input__inner {
         border: none !important;
         padding-left: 0;
     }
 }
</style>
<script>
export default {
    name: '',
    props:{
        value: {},
        useinfo:{},
        task_id:{},
        target_id:{},
    },
    data(){ 
        return{
            bugusedialogVisible:false, 
            useresult: [],
            usecmd: '',
        }
    },
    watch: {
        value(newVal, oldVal) {
            // 监测value的变化，并赋值。
            if (newVal) { 
                this.bugusedialogVisible = newVal; 
            }

        },
        bugusedialogVisible(val) {
            this.$emit("input", val); // 此处监测showMask目的为关闭弹窗时，重新更换value值，注意emit的事件一定要为input。
        }
    },
    mounted() {
        this.bugusedialogVisible = this.value; // 在生命周期中，把获取的value值获取给bugusedialogVisible

    },
    methods: {
        handleCloseUse() {
            this.bugusedialogVisible = false;
            this.useresult = [];
            this.useinfo.pocname = '';
            this.usecmd = ''
        },
        onSubmitUse() { //漏洞利用
            let _fromdata={
                cmd: this.usecmd,
                pocname: this.useinfo.pocname
            }
            let _url = '';
            if(this.task_id){
                _fromdata.task_id = this.task_id;
                _url = '/v2/task/vul/use/';  
            }
            else{
                _fromdata.target_id = this.target_id;
                _url = '/task/target/vul/use/'
            }
            this.$ajax({
                method: 'post',
                url: _url,
                data: _fromdata
            })
            .then(dt => {
                let res = dt.data;
                if (res.success) {
                    this.useresult.push({ usecmd: this.usecmd, useresult: res.result })
                    this.usecmd = '';
                } else {
                    this.$message({
                        message: res.error,
                        type: 'error'
                    });
                }
            }).catch(err => { })
        },
    },
}
</script>