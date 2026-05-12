<template>
    <el-form :model="webPathForm" label-width="0" status-icon ref="webPathFormRef"   >
        <el-form-item class="formItem1" label="" > 
                <div class="infobox2">
                    <el-form-item label="" prop="guess_rate">
                        <label class="dialog_item_label_m">猜测速率</label>
                        <el-select v-model="webPathForm.guessRate" size="small" 
                            placeholder="请选择" class="frame_width" 
                            :class="flag == 2 ? 'frame_width_m' : ''" :disabled="!isUpdate">
                            <el-option v-for="(item, index) in guessRate_list" 
                                :key="index" :label="item.label" 
                                :value="item.value"></el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="" prop="guess_rate">
                            <label class="dialog_item_label_m">猜测时长</label>
                            <el-select v-model="webPathForm.guessTimeout" size="small" placeholder="请选择" 
                            class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''"
                            :disabled="!isUpdate">
                                <el-option v-for="(item, index) in guessTime_list" :key="index" :label="item.label" :value="item.value"></el-option>
                            </el-select>
                        </el-form-item>
                    <el-form-item   label="" prop="scan_dict" >
                        <label class="dialog_item_label_m">路径字典</label>
                        <el-select v-model="webPathForm.scanDict" filterable placeholder="请选择" 
                        class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''"
                        multiple :disabled="!isUpdate">
                            <el-option v-for="(item, index) in pathDict_list" :key="index" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </el-form-item>
                    <el-form-item label="" prop="title_black" >
                        <label class="dialog_item_label_m">排除标题黑名单</label>
                        <el-input v-model="webPathForm.titleBlack" size="small" 
                        class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''"
                        autocomplete="off" :disabled="!isUpdate"
                            placeholder="请输入排除标题黑名单"></el-input>
                    </el-form-item>
                    <el-form-item label="" prop="title_black" >
                        <label class="dialog_item_label_m">打开动态语言识别</label>
                        <el-switch v-model="webPathForm.isIntelligent" size="small" 
                        class="frame_width" :class="flag == 2 ? 'frame_width_m' : ''"
                        autocomplete="off" :disabled="!isUpdate"
                            placeholder="请输入排除标题黑名单"></el-switch>
                    </el-form-item>
                </div>
            </el-form-item>
    </el-form>
</template>
<style lang="less" scoped>
.frame_width{
        width: 720px;
    }
.frame_width_m{
    width: 400px;
}
.dialog_item_label_m {
    display: inline-block;
    min-width: 100px;
    font-size: 14px;
    font-weight: 500;
    color: rgba(72, 72, 102, 0.87);
    line-height: 16px;
    margin-right: 16px;
    margin-left: 10px;
}

/deep/ .el-form-item {
    margin-bottom: 8px;
}
</style>
<script>
export default {
    name: '',
    data() {
        return {
            webPathForm:{
                guessRate: 0, // 猜测速率
                guessTimeout:0,
                scanDict: [], // 路径字典
                titleBlack: '', // 排除标题黑名单
                isIntelligent: false, // 打开动态语言识别
            },
            guessTime_list: [ ], // 猜测时间下拉列表
            guessRate_list: [ ], // 猜测速率下拉列表
            pathDict_list: [ ], // 路径字典下拉列表
            isUpdate:false,
            webPathScan:{},
            defaultscanDict:[],
        }
    },
    props: { 
        flag:{}, 
    },
    created() { 
    },
    methods: {
        getEnum(webPathScan,flag){
            this.webPathScan = webPathScan;
            this.guessTime_list = this.webPathScan.times;
            this.pathDict_list = this.webPathScan.scanDict;
            this.defaultscanDict = this.webPathScan.scanDictDefault;
            this.guessRate_list = this.webPathScan.speed; 
            if(flag == 1){
               
                this.guessTime_list.forEach(item=>{
                    if( item.isDefault==true){
                        this.webPathForm.guessTimeout = item.value
                    }
                })
                this.guessRate_list.forEach(item=>{
                    if(item.isDefault==true){
                        this.webPathForm.guessRate = item.value
                    }
                })
                this.pathDict_list.forEach(item => { 
                    if( item.isDefault==true){ 
                        this.webPathForm.scanDict.push(item.value)
                    }
                    
                }) 
                this.webPathForm.titleBlack = this.webPathScan.titleBlack;
            }
        },
        getIsUpdate(flag) {
            this.isUpdate = flag;
           
        },
        getConifg(_config) {
            this.webPathForm = _config;
       
        },
        getAllData() {
            return this.webPathForm;
        },
    },
}
</script>