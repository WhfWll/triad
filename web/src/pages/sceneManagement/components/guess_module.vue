<template>
    <el-form :model="webGuess" label-width="0" status-icon ref="commandFormRef"  >
        <el-form-item label="" v-show="ishttp != 1">
            <label class="dialog_item_label_m"  >爆破服务</label>

            <div class="serve_list frame_width" style="display: inline-block;" :class="flag == 2 ? 'frame_width_m' : ''"  > 

                <el-checkbox-group v-model="webGuess.services"  style=" display: inline-block;line-height: 30px;" :disabled="!isUpdate"
                    @change="handleChange()" >
                    <el-checkbox  v-for="(item, i) in serve_list" :key="i"  :label="item.value" :value="item.value">{{item.label}}</el-checkbox>
                    
                </el-checkbox-group>
            </div> 
        </el-form-item>
        <!-- <el-form-item prop = ' ' label=" " v-if="isHttp" > 
            <label class="dialog_item_label_m">验证码识别</label>
            <el-radio-group v-model="webGuess.verification_code_detect" class="radiobox" :disabled="!isUpdate" >
                <el-radio :label="1">数字与字母</el-radio>
                <el-radio :label="2">简单计算</el-radio>
            </el-radio-group>
        </el-form-item>   -->
         <el-form-item label=""  v-show="isImgMa">
            <label class="dialog_item_label_m">图形验证码</label> 
            <el-radio-group v-model="webGuess.captchaMode" :disabled="!isUpdate" > 
                <el-radio v-for="(item,i) in [
                {label:'字母数字',value:'common_alphanumeric'},
                {label:'简单算术',value:'common_arithmetic'},
                ]" :key="i" :label="item.value" :value="item.value">{{ item.label }}</el-radio>
            </el-radio-group>
        </el-form-item>
         <el-form-item label=""  >
            <label class="dialog_item_label_m">字典来源：</label> 
            <el-radio-group v-model="webGuess.dictType" :disabled="!isUpdate" > 
                <el-radio v-for="(item,i) in type_list" :key="i" :label="item.value" :value="item.value">{{ item.label }}</el-radio>
            </el-radio-group>
        </el-form-item>
        <div v-if="webGuess.dictType == 2">
            <el-form-item   label="" prop="scan_dict" >
                <label class="dialog_item_label_m">账号字典</label>
                <el-select v-model="webGuess.commonUserDict"  size="small" placeholder="请选择" class="frame_width"
                    :class="flag == 2 ? 'frame_width_m' : ''"
                        :disabled="!isUpdate">
                    <el-option v-for="(item, index) in commonUserDict" :key="index" :label="item.label" :value="item.value">
                    </el-option>
                </el-select>
            </el-form-item>
            <el-form-item   label="" prop="scan_dict" >
                <label class="dialog_item_label_m">口令字典</label>
                <el-select v-model="webGuess.commonPassDict"   size="small" placeholder="请选择" class="frame_width"
                    :class="flag == 2 ? 'frame_width_m' : ''"
                        :disabled="!isUpdate">
                    <el-option v-for="(item, index) in commonPassDict" :key="index" :label="item.label" :value="item.value">
                    </el-option>
                </el-select>
            </el-form-item>
        </div>
        <div v-if="webGuess.dictType == 3">
            <el-form-item label="" >
                <label class="dialog_item_label_m">补充账号</label>
                <el-input type="textarea"  rows="4" :disabled="!isUpdate" placeholder="可选择输入补充账号"
                    v-model="webGuess.addAccount" class="frame_width" resize="none" :class="flag == 2? 'frame_width_m' :''"
                    style=" vertical-align: text-top;">
                </el-input>
            </el-form-item> 
            <el-form-item label="" >
                <label class="dialog_item_label_m">补充口令</label>
                <el-input type="textarea"  rows="4" :disabled="!isUpdate" placeholder="可选择输入补充账号"
                    v-model="webGuess.addPass" class="frame_width" resize="none" :class="flag == 2 ? 'frame_width_m' : ''"
                    style=" vertical-align: text-top;">
                </el-input>
            </el-form-item>
            <el-form-item label="" v-show="ishttp != 1" >
                    <label class="dialog_item_label_m"> </label>
                    <el-checkbox v-model="webGuess.onlyUseAdd" :disabled="!isUpdate">只用补充凭证</el-checkbox>
                </el-form-item>
        </div>
        
       <el-form-item label=""  v-show="createtask !=1">
            <label class="dialog_item_label_m">猜测次数</label>
            <el-select v-model="webGuess.guessNum" size="small" placeholder="请选择"   class="frame_width"
                :class="flag == 2 ? 'frame_width_m' : ''"
                :disabled="!isUpdate">
                <el-option v-for="(item, index) in guessCount_list" :key="index" :label="item.label" :value="item.value"></el-option>
            </el-select>
        </el-form-item>
          <el-form-item label=""  v-show="createtask != 1">
            <label class="dialog_item_label_m">猜测时间</label>
            <el-select v-model="webGuess.guessTimeout" size="small" placeholder="请选择"   class="frame_width"
                :class="flag == 2 ? 'frame_width_m' : ''"
                :disabled="!isUpdate">
                <el-option v-for="(item, index) in guessTime_list" :key="index" :label="item.label" :value="item.value"></el-option>
            </el-select>
        </el-form-item>
        <el-form-item prop="pass_dict" label=""  v-show="createtask != 1">
            <label class="dialog_item_label_m">猜测速率</label>
            <el-select v-model="webGuess.guessRate" size="small" placeholder="请选择"   class="frame_width"
                :class="flag == 2 ? 'frame_width_m' : ''"
                :disabled="!isUpdate">
                <el-option v-for="(item, index) in guessRate_list" :key="index" :label="item.label" :value="item.value"></el-option>
            </el-select>
        </el-form-item>
       
    </el-form>
</template>
<style lang="less" scoped>
.frame_width {
    width: 720px;
} 
.frame_width_m{
    width: 400px;
}
.serve_list {
    display: inline-block;  
    vertical-align: middle;
}
.dialog_item_label_m {
    display: inline-block;
    min-width: 100px;
    font-size: 14px;
    font-weight: 500;
    color: rgba(72, 72, 102, 0.87);
    line-height: 16px;
    margin-right: 6px;
    margin-left: 10px;
}
</style>
<script>
export default {
    data(){
        return{
            isImgMa:false,
            webGuess:{
                services:[],
                guessNum:0,
                guessTimeout:300,
                guessRate:3,
                onlyUseAdd:false,
                addAccount:'',
                addPass:'',
                dictType:1,//字典来源
                commonUserDict:'',
                commonPassDict:'',
                captchaMode:'common_arithmetic',
                // verification_code_detect:''
            },
            guessCount_list: [], // 猜测次数下拉列表
            guessTime_list: [ ], // 猜测时间下拉列表
            guessRate_list: [], // 猜测速率下拉列表 
            isUpdate:false,
            checkList:[],
            serve_list:[],
            checkedData:[], 
            type_list:[],
            commonUserDict:[],
            commonPassDict:[],
            ishttpvalue:'',
            isHttp:false,
        }
    },
    props: {
        guess_config: {}, 
        flag:{},
        createtask:{},
        ishttp:{}
    },
    created() {
        // this.webGuess = this.guess_config; 
    },
    methods: {
        getEnum(guess,flag){
            this.serve_list = guess.services;
            this.type_list = guess.type;
            this.webGuess.onlyUseAdd = guess.onlyUseAdd;
            this.commonUserDict=guess.commonUserDict;
            this.commonPassDict = guess.commonPassDict;
            this.guessCount_list = guess.guessNum;
            this.guessTime_list = guess.guessTimeout;
            this.guessRate_list = guess.guessRate;
            // this.serve_list.forEach(item =>{
            //     if( item.label == 'http'){
            //         this.ishttpvalue=item.value;
            //     }
            // })
            if(flag == 1){
                this.serve_list.forEach(item =>{
                    if( item.isDefault == true){
                        this.webGuess.services.push(item.value);
                    }
                })
                this.type_list.forEach(item =>{
                    if( item.isDefault==true){
                        this.webGuess.dictType = item.value
                    }
                })
                this.guessCount_list.forEach(item =>{
                    if( item.isDefault==true){
                        this.webGuess.guessNum = item.value
                    }
                })
                this.guessTime_list.forEach(item =>{
                    if( item.isDefault==true){
                        this.webGuess.guessTimeout = item.value
                    }
                })
                this.guessRate_list.forEach(item =>{
                    if( item.isDefault==true){
                        this.webGuess.guessRate = item.value
                    }
                })
            }
        },
        getIsUpdate(_isUpdate) {
            this.isUpdate = _isUpdate;
        },
        getConifg(_config) {  
            if(_config.commonUserDict == 0){
                _config.commonUserDict =''
            }
            if(_config.commonPassDict == 0){ 
                _config.commonPassDict =''; 
            } 
            this.webGuess = _config;  
            if(!this.webGuess.services){
                this.webGuess.services=[]
            }
        },
        getAllData() {
            // let _infos = [];
            // this.serve_list.forEach(item=>{
            //     if(item.ischecked){
            //         item.user_dict = item.user_checked;
            //         item.pass_dict = item.pass_checked;
            //         _infos.push(item);
            //     }
            //     if(this.ishttp == 1){ //是http
            //         if (item.ischecked) {
            //             _infos.push(item);
            //         }
            //     }
            // }); 
            // this.webGuess.infos = _infos; 
            
            if(this.webGuess.commonUserDict =='') {
                this.webGuess.commonUserDict=0;
            }
            if(this.webGuess.commonPassDict =='') {
                this.webGuess.commonPassDict=0;
            }
            if(this.webGuess.dictType == 2){
                //通用
                this.webGuess.addAccount = '';
                this.webGuess.addPass = '';
                this.webGuess.onlyUseAdd = false;
            }
            if(this.webGuess.dictType == 3){
                //补充
                this.webGuess.commonUserDict = 0;
                this.webGuess.commonPassDict = 0;
            }


            return this.webGuess;
        },
        fnWebGetData() {
            //获取口令猜测数据
            this.$ajax
                .get("/task/task/get/user_password/dictionary/", {})
                .then((dt) => {
                    let res = dt.data;
                    if (res.success) { 
                        res.result.forEach(ele  => {
                            ele.ischecked = false;
                        })
                        this.serve_list = res.result; 
                    
                        // this.wordguesstable.pop();

                        //设置选中状态
                        this.$nextTick(() => {
                            let _that = this;
                            let _checkedData = this.checkedData;
                            let _data = this.serve_list;
                            for (let i = 0; i < _data.length; i++) {
                                for (let j = 0; j < _checkedData.length; j++) {
                                    if (_data[i].service == _checkedData[j].service) { 
                                        _data[i].user_dict = _checkedData.user_checked;
                                        _data[i].pass_dict = _checkedData.pass_checked;
                                        _data[i].ischecked = true;

                                    }
                                }
                            }

                        })
                    }
                })
                .catch((err) => { });
        },
        handleChange(){ //选择 服务
            console.log(this.webGuess.services)
            if(this.webGuess.services.indexOf(this.ishttpvalue) !=-1){
                this.isHttp = true;
            }
            this.webGuess.services.forEach(item =>{
                if(item == '22'){
                    this.isImgMa = true;
                }else{
                this.isImgMa = false;
            }
            })
            
            // if(serve.pass_checked =='http'){
            //     this.webGuess.http_weak_password_blasting.is_open = serve.ischecked;
            // }
        }, 
    },
    watch: {
       //监听数组得变化：webGuess.services
        'webGuess.services': function (newVal, oldVal) {
            console.log(newVal,oldVal)
            if(newVal.indexOf(this.ishttpvalue) !=-1){
                this.isHttp = true;
            }else{
                this.isHttp = false;
            }
            this.webGuess.services.forEach(item =>{
                if(item == '22'){
                    this.isImgMa = true;
                }else{
                this.isImgMa = false;
            }
            })
    },
}}
</script>