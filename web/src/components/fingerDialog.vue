<template>
  <div>
        <el-dialog
            :title="dialogTitle"
            :visible.sync="visible"
            width="1184px"
            class="buginfobox fingerDialog" 
            :close-on-click-modal="false" 
            :show-close="false" >
            <div class="dialog_b_btn">  
                <el-button size="small" @click="toEdit()" v-if="status == 2 ">编辑</el-button>
                <el-button size="small" @click="saveNewfinger" v-if="status !== 2">保存</el-button>
				<el-button size="small" @click="handleClose">关闭</el-button>
            </div>
            <div class="buginfo_box" > 
                <div class="bugbasicinfo">
                    <el-form :model="newFingerform" :rules="mainRules" label-width="125px"  status-icon  ref="" >
                         <el-form-item
                            class=""
                            label="指纹名称"
                            prop="cnName" > 
                            <el-input 
                                size="small"
                                v-model="newFingerform.cnName"  
                                autocomplete="off"  :disabled="isDisabled"
                                placeholder="请输入指纹名称" style="width: 320px;"></el-input> 
                        </el-form-item>    
                        <el-form-item label="指纹类型" prop="appClass" >
                            <el-select class="typeSelect" v-model="newFingerform.appClass" placeholder="请选择"  :disabled="isDisabled" clearable size="small"  >
                                <el-option
                                    v-for="item in vulobjectlist"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item label="分层" prop="level" >
                            <el-select class="typeSelect" v-model="newFingerform.level" placeholder="请选择"  :disabled="isDisabled" clearable size="small"  >
                                <el-option
                                    v-for="item in levellist"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item
                            class="name"
                            label="系统自带"
                            prop="source" > 
                            <el-select class="typeSelect" v-model="newFingerform.source" placeholder="请选择"  :disabled="isDisabled" clearable size="small"  >
                                <el-option
                                    v-for="item in sourcelist"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value">
                                </el-option>
                            </el-select>
                        </el-form-item> 
                        <el-form-item label="指纹分类" prop="fingerType" >
                            <el-select class="typeSelect" v-model="newFingerform.fingerType" placeholder="请选择"  :disabled="isDisabled" clearable size="small"  >
                                <el-option
                                    v-for="item in fingerTypelist"
                                    :key="item.value"
                                    :label="item.label"
                                    :value="item.value">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item
                            class=""
                            label="版本"
                            prop=" " > 
                            <el-input 
                                size="small"
                                v-model="newFingerform.appVersion"  
                                autocomplete="off"   :disabled="isDisabled"
                                placeholder="请输入版本" style="width: 320px;"></el-input> 
                        </el-form-item> 
                        <el-form-item
                            class=""
                            label="指纹规则"
                            prop="flag" > 
                            <el-input 
                                type="textarea"
                                resize="none"
                                rows='5'  :disabled="isDisabled"
                                v-model="newFingerform.flag"  
                                autocomplete="off"  
                                placeholder="" style="width: 720px;"></el-input> 
                        </el-form-item> 
                        <!-- <el-form-item
                            class="fingerinfo1"
                            label="指纹信息"
                            prop="info" > 
                            <el-input 
                                v-model="newFingerform.info"  
                                autocomplete="off"
                                placeholder="请输入指纹信息" ></el-input> 
                        </el-form-item>  -->
                        <!-- <el-form-item
                            class="width100"
                            label="正则规则"
                            prop="regex" > 
             
                                <el-table
                                     class="diatable"
                                    :data="zhengzeTable"
                                    size='small'  @selection-change="handleSelectionRuleChange"
                                    style="width: 100%" :class="status === 2 ? 'queryTable' : ''">
                                    <el-table-column  
                                        width="55" 
                                        type="selection">
                                    </el-table-column> 
                          
                                    <el-table-column
                                        prop="value"
                                        label="规则" >
                                    </el-table-column>
                         
                                </el-table>
                        </el-form-item>  -->
                        <!-- <el-form-item
                            class="width100 noStar"
                            label="图标hash"
                            prop="info" > 
                        

                                <el-table
                                    class="diatable"
                                    :data="metaTable"
                                    size='small'  @selection-change="handleSelectionMetaChange"
                                    style="width: 100%" :class="status === 2 ? 'queryTable' : ''">
                                   
                                    <el-table-column prop="value"
                                        label="规则" > 
                                    </el-table-column> 
                                  
                                </el-table>
                        </el-form-item> -->
                        <!-- <el-form-item
                        class="width100 noStar"
                        label="证书hash"
                        prop="info" > 
  
                            <el-table
                                class="diatable"
                                :data="headerTable"
                                size='small'  @selection-change="handleSelectionHeaderChange"
                                style="width: 100%" :class="status === 2 ? 'queryTable' : ''">
                               
                           
                                <el-table-column  prop="value"                                
                                    label="规则"  v-if="status === 2"> 
                                </el-table-column>  
                           
                                 
                            </el-table>
                    </el-form-item>   -->
                        
                </el-form>
                </div>   
            </div>
        </el-dialog>
        <el-dialog
            title="正则规则"
            width="1184px"
            :visible.sync="rulevisible"
            class="rulebox" 
            :close-on-click-modal="false" 
            :show-close="false" >
            <div class="dialog_b_btn">  
                <el-button size="small" @click="saveNewRule">确定</el-button>
				<el-button size="small" @click="handleCloseRule">取消</el-button>
            </div>
            <div class="buginfo_box" > 
                <div class="bugbasicinfo">
                    <el-form :model="newRuleForm" label-width="112px"  status-icon  ref="newRuleForm" :rules="ruleRules">
                        <!-- <el-form-item
                            class="rulesmallbox"
                            label="规则类型"
                            prop="value"
                            style="margin-bottom: 20px;"> 
                            <el-input  
                                class="ruleinfo"
                                v-model="newRuleForm.value"
                                autocomplete="off" 
                                placeholder="" ></el-input> 
                        </el-form-item> -->
                        <el-form-item
                            label=""
                            prop="value"  label = "规则类型" > 
                            <el-select v-model="newRuleForm.value" style=" width:320px;" clearable  size="small">  
                                <el-option
                                    v-for="(item,i) in validateRules"
                                    :key="i"
                                    :label="item"
                                    :value="item"> 
                                </el-option>
                            </el-select> 
                        </el-form-item>
                         <el-form-item
                             class="rulesmallbox"
                             prop="version" label="正则表达式"> 
                            <el-input 
                                type="textarea"  
                                class="ruleinfo"
                                 v-model="newRuleForm.version"
                                autocomplete="off" 
                                placeholder="" ></el-input> 
                        </el-form-item>
                </el-form>
                </div>   
            </div>
        </el-dialog>
        <el-dialog
            title="meta规则"
            :visible.sync="metavisible"
            width="1184px"
            class="hashbox" 
            :close-on-click-modal="false" 
            :show-close="false" @close="closeMetaDialog">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="saveMeta">确定</el-button>
				<el-button size="small" @click="handleCloseMeta">取消</el-button>
            </div>
            <div class="hash-mainDiv" > 
                    <el-form :model="generateMetaForm" class="hashEl-form" label-width="112px"  status-icon  ref="metaForm" :rules="hashRules">
                        <el-form-item
                        class="name"
                        label="key"
                        prop="key" > 
                        <el-input 
                            v-model="generateMetaForm.key"  
                            autocomplete="off" 
                            placeholder="请输入key" ></el-input> 
                        </el-form-item> 
                        <el-form-item
                        class="name"
                        label="value"
                        prop="value" > 
                        <el-input 
                            v-model="generateMetaForm.value"  
                            autocomplete="off" 
                            placeholder="请输入value" ></el-input> 
                    </el-form-item> 
                </el-form>
            </div>
        </el-dialog>
        <el-dialog
        title="headers规则"
        :visible.sync="headervisible"
        width="1184px"
        class="hashbox" 
        :close-on-click-modal="false" 
        :show-close="false" @close="closeHeaderDialog">
        <div class="dialog_b_btn">  
            <el-button size="small" @click="saveHeaderForm">确定</el-button>
            <el-button size="small" @click="handleCloseHeader">取消</el-button>
        </div>
        <div class="hash-mainDiv" > 
                <el-form :model="generateHeaderForm" class="hashEl-form" label-width="112px"  status-icon  ref="headerForm" :rules="headerRules">
                    <el-form-item
                    class="name"
                    label="key"
                    prop="key" > 
                    <el-input 
                        v-model="generateHeaderForm.key"  
                        autocomplete="off" 
                        placeholder="请输入key" ></el-input> 
                    </el-form-item> 
                    <el-form-item
                    class="name"
                    label="value"
                    prop="value" > 
                    <el-input 
                        v-model="generateHeaderForm.value"  
                        autocomplete="off" 
                        placeholder="请输入value" ></el-input> 
                </el-form-item> 
            </el-form>
        </div>
    </el-dialog>
  </div>
</template>

<script> 
import xzbutton from "@/components/XzButton.vue"; 
import delbutton from "@/components/DelButton.vue";
import { fingerprint } from '@/api/tool.js'
export default {
  name: 'fingerDialog',
  components: {
        xzbutton,
		delbutton,
    },
  data () {
        var validateReg =  (rule, value, callback) => {
            let isReg
            try {
                isReg = eval(value) instanceof RegExp
                } catch (e) {
                isReg = false
            }
            if (isReg) {
                callback();
            } else {
                return callback(new Error('格式不正确,请填写正则表达式'));
            }
        };
        var validateInfo = (rule, value, callback) => {
            if (!['a', 'h', 'o'].includes(value.substr(0, 1))) {
                return callback(new Error('格式不正确'));
            } else {
                var regex = new RegExp(':', 'g'); // 使用g表示整个字符串都要匹配
                var result = value.match(regex);          //match方法可在字符串内检索指定的值，或找到一个或多个正则表达式的匹配。
                var count=!result ? 0 : result.length;
                if (count !== 2) {
                    return callback(new Error('格式不正确'));
                } else {
                    callback();
                }
            }
        };
    return {
        rulevisible:false,
        metavisible:false,
        headervisible:false,
        zhengzeTable:[],
        metaTable:[],
        headerTable:[],
        validateRules:[],
        newFingerform:{
            id: '',
            cnName:'',
            type:'',
            re_rule:'',
            regex:'',
            ico:'',
            description:'',
            rule:'',
            source:'',//系统自带
            fingerType:'',//指纹分类
            appClass:'',//分类
            desc:'',//
            level:'',//分层
            flag:'',//规则
            appVersion:'',
        },
        newRuleForm: {
            value: '',
            version: ''
        },
        generateMetaForm: {
            key: '',
            value: ''
        },
        generateHeaderForm: {
            key: '',
            value: ''
        },
        mainRules:{ 
            cnName:[
                { required: true, message: '指纹名称不能为空', trigger: 'blur' }, 
                { max: 50, message: '指纹名称最大长度不能超过50', trigger: 'blur' } 
            ],  
            appClass:[
                { required: true, message: '指纹类型不能为空', trigger: 'blur' },
            ],
            level:[
                { required: true, message: '请选择分层', trigger: 'blur' },
            ],
            source:[
                { required: true, message: '请选择系统自带', trigger: 'blur' }, 
            ],
            fingerType:[
                { required: true, message: '请选择指纹分类', trigger: 'blur' }, 
            ],
            flag:[
                {required: true, message: '指纹规则不能为空', trigger: 'blur'}
            ]
        },
        ruleRules:{ 
            value:[
                { required: true, message: '规则类型不能为空', trigger: 'blur' },
                // { validator: validateReg, trigger: 'blur' } 
            ],  
            version:[
                { required: true, message: '正则表达式不能为空', trigger: 'blur' },
                // { validator: validateReg, trigger: 'blur' } 
            ],
        },
        hashRules:{ 
            key:[
                { required: true, message: '请输入key', trigger: 'blur' }
            ],  
            value:[
                { required: true, message: '请输入value', trigger: 'blur' },
            ],
        },
        headerRules:{ 
            key:[
                { required: true, message: '请输入key', trigger: 'blur' }
            ],  
            value:[
                { required: true, message: '请输入value', trigger: 'blur' },
            ],
        },
        alldelMetaVisible:false,
        alldelHeaderVisible:false,
        alldelRuleVisible:false,
        multipleRuleSelection: [],
        multipleMetaSelection: [],
        multipleHeaderSelection: [],
        status: 0, // 0 新增  1是编辑 2是详情
        dialogTitle2: '',
        vulobjectlist: [],
        levellist:[],
        sourcelist:[
            { value: 1, label: '系统自带' },
            { value: 2, label: '用户添加' }
        ],
        fingerTypelist:[
            { value: 1, label: 'web指纹' },
            { value: 2, label: '系统指纹' }
        ],
        isDisabled:false,
    }
  },
  props: {
    visible: {
      type: Boolean
    },
    dialogTitle: {
      type: String
    },
    editFingerInfo: {
      type: Object
    },
    // vulobjectlist: {
    //     type: Array
    // }
  },
  watch: {
      editFingerInfo: {
        handler (n, o) {
            this.getDefaultValue()
        },
        deep: true // 深度监听父组件传过来对象变化
          
      }
  },
  mounted () {
     
      this.dialogTitle2 = this.dialogTitle
      this.getFingerType()
      switch (this.dialogTitle2) {
        //0新建。1编辑 2详情
        case '新建指纹':
            this.status = 0
            break
        case '编辑指纹':
            this.status = 1
            break
        default:
            this.status = 2;
            this.isDisabled = true;
            break
      }
      if (this.editFingerInfo) {
        // this.getDefaultValue()
        // console.log(this.editFingerInfo);
        this.getInfo(); //获得详情
      } else {
        this.newFingerform.id = ''
        this.newFingerform.name = ''
        this.newFingerform.type = ''
        this.newFingerform.re_rule = ''
        this.zhengzeTable = []
        this.metaTable = []
        this.headerTable = []
      }
      
  },
  methods: {
    async getInfo(){
        const res = await fingerprint.getInfo({
            id:this.editFingerInfo.id
        });
        if(res.code == 200){
            console.log(res.data);
            this.newFingerform.cnName = res.data.fingerName;
            this.newFingerform.appClass = parseInt(res.data.fingerClass);
            this.newFingerform.source = parseInt(res.data.source);
            this.newFingerform.level = parseInt(res.data.levelID);
            this.newFingerform.appVersion = res.data.appVersion;
            this.newFingerform.fingerType = parseInt(res.data.fingerType);
            this.newFingerform.flag = res.data.flag;

        }

    },
    getDefaultValue () {   
        this.newFingerform.id = this.editFingerInfo.id
        this.newFingerform.cnName = this.editFingerInfo.fingerName
        this.newFingerform.appClass = this.editFingerInfo.fingerClass
        this.newFingerform.flag = this.editFingerInfo.rule
    },
    handleClose () {
      this.$emit('handleClose')
    },
    handleCloseRule () {
        this.rulevisible = false
    },
    saveNewRule () {
        this.$refs['newRuleForm'].validate((valid) => {
            if(valid){
                this.zhengzeTable.push({
                    name: this.newRuleForm.value,
                    value: this.newRuleForm.version
                })
                this.handleCloseRule()
                this.newRuleForm.value = ''
                this.newRuleForm.version = ''
            }
        })
        
    },
    saveMeta () {
        this.$refs['metaForm'].validate((valid) => {
            if(valid){
                this.metaTable.push({
                    name: this.generateMetaForm.key,
                    value: this.generateMetaForm.value
                })
                this.metavisible = false
                
            }
        })
    },
    saveHeaderForm () {
        this.$refs['headerForm'].validate((valid) => {
            if(valid){
                this.headerTable.push({
                    name: this.generateHeaderForm.key,
                    value: this.generateHeaderForm.value
                })
                this.headervisible = false
                this.generateHeaderForm.key = ''
                this.generateHeaderForm.value = ''
            }
        })
    },
    closeMetaDialog () {
        this.generateMetaForm.key = ''
        this.generateMetaForm.value = ''
    },
    closeHeaderDialog () {
        this.generateHeaderForm.key = ''
        this.generateHeaderForm.value = ''
    },
    handleCloseMeta () {
        this.metavisible = false
    },
    handleCloseHeader () {
        this.headervisible = false
    },
    handleSure () {
      
    },
    //正则指纹类型下拉
    async openrule(){
        this.newRuleForm.value = ''
        this.newRuleForm.version = ''
        let params = {}
        const res = await fingerprint.openrule(params)
        if(res.success){  
            this.validateRules = res.results;
            this.rulevisible =true;
        }else{
            // this.$message({
            //     message:res.msg,
            //     type: 'error'
            // });
        }
    },
    openMeta(){
        this.metavisible =true;
    },
    openheader(){
        this.headervisible =true;
    },
    async saveNewfinger(){ //保存
        var params={  
            appName: this.newFingerform.cnName,
            appClass:this.newFingerform.appClass,
            appVersion: this.newFingerform.appVersion ,
            level:this.newFingerform.level ,
            source:this.newFingerform.source,
            fingerType:this.newFingerform.fingerType,
            flag:this.newFingerform.flag
        } 
        let res = '';
        
        if (this.editFingerInfo && this.editFingerInfo.id) {
            params.id = this.editFingerInfo.id;
            res = await fingerprint.updateNewfinger(params)
        } else {
            res = await fingerprint.createNewfinger(params)
        }
        if(res.code ==200){ 
            this.$message({
                message:'保存成功',
                type: 'success'
            }); 
            this.visible = false;
            this.newFingerform.id= '';
            this.newFingerform.cnName= '';
            this.newFingerform.appClass =''; 
            this.newFingerform.appVersion = ''; 
            this.newFingerform.level = ''; 
            this.newFingerform.source = ''; 
            this.newFingerform.fingerType = ''; 
            this.newFingerform.flag = ''; 
            this.$emit('refreshData', true)
        }else{
            this.$message({
                message:res.msg,
                type: 'error'
            });
        }
    },
    async getFingerType () {
          //新建指纹类型下拉
        //   let url = '/tools/fingerprint/type/select/'
        //   if ( this.editFingerInfo && this.editFingerInfo.source === 1) {
        //     url = '/tools/fingerprint/type/select/'
        //   } else {
        //     url = '/tools/fingerprint/type/select/'
        //   }
       
        let params = {}
        const res = await fingerprint.getVulObjectlist(params)
        if(res.code == 200){  
            console.log(1111);
            this.vulobjectlist = res.data.class;
            this.levellist = res.data.level;
            this.devlist = res.data.softOrHard
        }else{
            // this.$message({
            //     message:res.msg,
            //     type: 'error'
            // });
        } 
    },
    // 删除规则
    handleDelRule:function(scope){ //删除 
        // this.$ajax({
        //     method:'delete',
        //     url:'/systems/information/delete/',
        //     data: {
        //         id:scope.row.id+''
        //     } 
        // }).then(data => { 
        //     var dt = data.data;  
        //     if(dt.success){ 
        //         this.$message({
        //             message:dt.msg,
        //             type: 'success'
        //         });
        this.zhengzeTable.splice(scope.$index, 1)
        scope._self.$refs[`popover_id-${scope.$index}`].doClose()
        //         this.getData();
        //     }else{
        //         this.$message({
        //             message:dt.error,
        //             type: 'error'
        //         });
        //     }  
            
        // })
        // .catch(data=>{
        //     console.log(data); //错误信息
        // }); 
    },
    // 删除meta
    handleDelMeta:function(scope){ //删除 
        this.metaTable.splice(scope.$index, 1)
        scope._self.$refs[`popover_id-${scope.$index}`].doClose()
    },
    // 删除header
    handleDelHeader:function(scope){ //删除 
        this.headerTable.splice(scope.$index, 1)
        scope._self.$refs[`popover_id-${scope.$index}`].doClose()
    },
    handleSelectionMetaChange(val){
        this.multipleMetaSelection = val
    },
    handleSelectionHeaderChange(val){
        this.multipleHeaderSelection = val
    },
    AllDelMeta () {
        this.multipleMetaSelection.forEach(item => {
            this.metaTable = this.metaTable.filter(item2 => {
                return item2.name !== item.name && item2.value !== item.value
            })
        })
        this.alldelMetaVisible = false
    },
    AllDelHeader () {
        this.multipleHeaderSelection.forEach(item => {
            this.headerTable = this.headerTable.filter(item2 => {
                return item2.name !== item.name && item2.value !== item.value
            })
        })
        this.alldelHeaderVisible = false
    },
    handleSelectionRuleChange(val){
        console.log(val)
        this.multipleRuleSelection = val
    },
    AllDelRule () {
        this.multipleRuleSelection.forEach(item => {
            this.zhengzeTable = this.zhengzeTable.filter(item2 => {
                return item2.name !== item.name
            })
        })
        this.alldelRuleVisible = false
    },
    // 编辑
    toEdit () {
        this.dialogTitle2 = '编辑指纹'
        this.status = 1;
        this.isDisabled = false;
    }
  }
}
</script>

<style lang="less" scoped>
// .rulesmallbox {
//      display:flex;
//      justify-content:center;
//      align-items:center;
// }
.diatable{
    // margin-top: 20px;
     border:1px solid #E8E8F5;
     background-color: #FFFFFF;
     padding: 20px;
    /deep/ th{
        padding: 7px 0;
    }
    &.queryTable{
        /deep/ .el-table-column--selection{
            .cell{
                display: none;
            }
            
        }
    }
}
.fingername .el-input{
     width: 320px;
}
.typeSelect{
    width: 320px;
}
    .bugbasicinfo{
        padding: 24px;
        background: #F7F7FB;
    }
 .buginfo_box{
       background: #F7F7FB; 
       height:100%;
    }
    .bugbasicinfo{
        padding: 24px;
        // background: #fff;
        .el-form-item{
            margin-bottom: 14px;
        }
        /deep/ .el-form-item__label{
            text-align: left;
            border-left: 3px solid #4C7AE3;
            line-height: 16px;
            padding-left: 8px;
            margin-top: 12px;
            position: relative;
        }
        /deep/ .el-form-item.is-required:not(.is-no-asterisk) .el-form-item__label-wrap>.el-form-item__label:before, /deep/ .el-form-item.is-required:not(.is-no-asterisk)>.el-form-item__label:before{
            position: absolute;
            right: 0
        }
        .width100{
            position: relative;
            /deep/ .el-form-item__content{
                margin-left: 0!important;
                margin-top: 60px;
            }
            .optionBtns{
                position: absolute;
                right: 0;
                top: -53px;
                .btnMarginR{
                    margin-right: 10px;
                }
            }
        }
        .noStar{
            /deep/ .el-form-item__label:before{
                display: none;
            }
        }
    }
    .bugotherinfo{
        margin-top: 32px;
        .part_title{
            margin-bottom: 8px;
        }
        .content{
            background:rgba(255,255,255,1);
            border-radius:2px;
            border:1px solid rgba(232,232,245,1);
            padding: 12px 16px;
            color:rgba(72,72,102,0.64);
            font-size: 13px;
        }
    }
    .part_title{
        font-size: 14px; 
        margin-bottom:16px;
        font-weight: 800;
        border-left: 3px solid #4C7AE3;
        padding-left: 10px;
        height: 14px;
        line-height: 14px;
        color:rgba(72,72,102,0.87);
    }
    .errormsg{
        font-size: 12px;
        color: #F56C6C;
        padding-top: 4px;
    }
    .nolabelleft{
        border-left: none;
    }
    .infobox /deep/ .el-form-item__content{
        display: flex;
        justify-content: center;
        margin-left: 0px;
    }
    .ruleinfo /deep/ .el-textarea__inner{
        width: 100%;
        height: 88px!important;
        background: #FFFFFF;
        border-radius: 2px;
        border: 1px solid #E8E8F5;
    }
    .rulebox{
        position: fixed;
        /deep/ .el-dialog{
            height:328px;
        }
        /* /deep/ .el-form .el-form-item__content{
            margin-left: 0px;
        } */
    }
    .hashbox{
        .el-input {
            position: relative;
            font-size: 13px;
            display: inline-block;
            width: 60%;
        }
        position: fixed;
        /deep/ .el-dialog{
            height:328px;
            .el-dialog__body{
                /* display: flex; */
                justify-content: center;
                align-items: center;
                padding:50px 30px!important;
                .hash-mainDiv{
                    width: 100%;
                    
                    .hashEl-form{
                        .el-form-item__label {
                            text-align: left;
                            border-left: 3px solid #4C7AE3;
                            line-height: 16px;
                            padding-left: 8px;
                            margin-top: 12px;
                            position: relative;
                        }
                        .rulesmallbox{
                            .el-form-item__content{
                                display: flex;
                                .hash_btns{
                                    width: 248px;
                                    display: flex;
                                    justify-content: space-between;
                                    margin-left: 38px;
                                     .el-button{
                                         height: 32px;
                                         line-height: 32px;
                                         padding: 0 11px;
                                         
                                     }
                                    .uploadBtn{
                                        color: #4C7AE3;
                                        border: 1px solid #4C7AE3;
                                    }
                                    .uploadDiv{
                                        height: 32px;
                                        line-height: 32px;
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
    .isError /deep/ .el-textarea__inner{
        border: 1px solid  #F56C6C  !important;
    }
    /deep/ .el-table tr th:nth-child(1) .cell{
        padding-left: 32px;
    }
    /deep/ .el-table tr td:nth-child(1) .cell{
        padding-left: 32px;
    }
    /deep/ .el-upload-list {
        display:none
    }
    .fingerdelbtn{
        margin-left: 8px;
    }
</style>