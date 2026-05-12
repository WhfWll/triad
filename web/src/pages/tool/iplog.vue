<!--IP域名绑定............................................................................................................. -->
<template>
  	    <div > 
	  	    <div class="auxiliarytool context_box_bg"> 
            <div class="basic">
                <div><strong>介绍：</strong> IP域名绑定是将IP与域名进行映射，在执行渗透任务时，将直接对域名映射的IP进行测试，不再去查询域名关联的IP。</div>
                <div><strong>使用帮助：</strong>IP域名绑定时，允许单个IP绑定1个或多个域名。</div>
            </div>
            <div style="margin-top:20px" class="search-box">
                <xzbutton 
                    type="primary" 
                    @click="newAdd()"
                    size="small">新增绑定</xzbutton>  
                <delbutton 
                    :width="170"  
                    :disabled="!multipleSelection.length" style="margin-left: 8px;" @click="AllDel()"></delbutton> 
                <div class="serach-condition">
                    <div class="search-text">
                        <el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearch"  v-model="search_field" class="input-with-select"  size="small" clearable > </el-input>
                        <xzbutton 
                        type="primary" 
                        @click="handlesearch" 
                        :disabled="false" 
                        size="small"  >搜索</xzbutton>  
                    </div>
                    <div >
                        <xzbutton 
                        type="primary" 
                        @click="handleReset" 
                        :disabled="false" 
                        size="small">重置</xzbutton>  
                    </div>  
            
                </div>
            </div>
            <div>
                <el-table
                @selection-change="handleSelectionChange"
                :data="tableData" 
                style="width: 100%" @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column type="selection" width="55">
                </el-table-column>
                <el-table-column
                    prop="ip"
                    label="IP">
                </el-table-column>
                <el-table-column
                    prop="hosts"
                    label="域名">
                    <template slot-scope="scope" > 
                        <div v-if="showEditFileNameButton && rowId == scope.row.id">
                            <el-link  :underline="false" class="link_primary"  @click="handleInfo(scope.row)">编辑</el-link>  
                            <!-- v-if="scope.row.source !== 1" -->
                            <el-popover
                                placement="bottom"
                                width="170"   
                                :ref="`popover_id-${scope.row.id}`"
                                popper-class="delButton_popper" >
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                    <el-button size="mini" type="primary" @click="handleDel(scope)">确定</el-button>
                                </div> 
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0" slot="reference">删除</el-link>  
                            </el-popover>
                            <!-- <el-link :underline="false" class="link_danger linkafter" style="padding:0">删除</el-link>  -->
                        </div>
                        <div v-else >
                            <span>{{scope.row.hosts}}</span>
                        </div> 
					</template>
                </el-table-column>
                </el-table>
                <el-pagination
                    :page-size="pageSize" 
                    background
                    layout="total,  prev, pager, next, sizes,jumper"
                    :total="totalpage"
                    :current-page="currentpage"
                    @current-change = "currentchange"
                     @size-change="handleSizeChange" 
                    >
                </el-pagination>
            </div>
	  	</div>
          <el-dialog
          :title="title"
          :visible.sync="dialogFormVisible"
          width="1184px"
          class="fingerValidate" 
          :close-on-click-modal="false" 
          :show-close="false">
          <div class="dialog_b_btn">  
            <el-button size="small" @click="saveAdd()">确定</el-button>
              <el-button size="small" @click="cancaliplogVisible()">关闭</el-button>
          </div>
          <div class="buginfo_box" > 
              <div class="bugbasicinfo">
                  <el-form :model="checkIPform" status-icon  ref="ruleFormaddip" :rules="rules" label-width="80px">
                    <el-form-item prop="ip" label="IP">
                        <!-- <label class="dialog_item_label topline">IP</label> -->
                        <el-input v-model="checkIPform.ip" size="small" style="width:520px" autocomplete="off"
                            placeholder="请输入ip"></el-input>
                    </el-form-item>
                    <el-form-item label="域名" prop="domain">
                        <el-input v-model="checkIPform.domain" size="small" style="width:520px" autocomplete="off"  type="textarea"  
                        :rows="4" placeholder="请输入域名"></el-input>
                    </el-form-item>
              </el-form>
              </div>   
          </div>
      </el-dialog>
  	</div>
</template>

<script>  
import About from "@/components/About.vue";
import Operation from "@/components/Operation";
import xzbutton from "@/components/XzButton.vue";
import delbutton from "@/components/DelButton.vue";
import { auxiliarytool } from '@/api/tool.js'
export default ({
    name:'iplog',
    components:{
        About,
        Operation,
        delbutton,
        xzbutton
    },
    data(){  
        var validatePass2 = (rule, value, callback) => { 
            if (!value) {
              callback(new Error('请输入IP地址'));
            } else {
                const re =
                    /^(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|[0-9])\.((1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)\.){2}(1\d{2}|2[0-4]\d|25[0-5]|[1-9]\d|\d)$/
                if (re.test(value)) {
                    callback();
                } else {
                    return callback(new Error('请输入正确的IP地址'));
                }
            }
        };
    	return{
            dialogFormVisible:false,
            checkIPform:{
                id: '',
            	ip:'',
                domain: '',
            },
            showEditFileNameButton:false,
            rowId:'',
            multipleSelection: [],
            page:1,
            total:0,
            pageSize:10,
            loading:false,  
    		activeName:'httplog',
            tableData:[],
            search_field:'',
            totalpage:0,
            currentpage:1,
            createInfiltration:false,
            pageType:'auxi',
            rules:{
                ip: [
                    { required: true, message: 'IP不能为空', trigger: 'blur' }, 
                    { validator: validatePass2, trigger: 'blur' },
                ],
                domain: [
                    { required: true, message: '域名不能为空', trigger: 'blur' }, 
                ],
            },
            title: '新增绑定'
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/auxiliarytool"; 
    },
    mounted:function(){   
        this.getIpData(); 
    },
    methods:{  
        // 漏洞页列表数据显示
        async getIpData(){
            let params = {
                search:this.search_field,
                page:this.page,
                size: this.pageSize,
            }
            const res = await auxiliarytool.getiplogData(params)
            if(res.code==200){ 
                this.tableData = res.data.list;
                this.totalpage = res.data.total;
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            } 
        },
        cancaliplogVisible(){
            this.dialogFormVisible = false;
            this.checkIPform.domain ='';
            this.checkIPform.ip = '';
            this.checkIPform.id = '';
            this.$refs.ruleFormaddip.resetFields(); 
        }, 
        async saveAdd(){
            if (this.checkIPform.id) {
                let params = {
                    id: this.checkIPform.id,
                    ip:this.checkIPform.ip,
                    hosts:this.checkIPform.domain,
                }
                const res = await auxiliarytool.updateIplog(params)
                if(res.code == 200){
                    let obj = this.tableData.find(item => item.id === this.checkIPform.id)
                    if (obj) {
                        obj.ip = this.checkIPform.ip
                        obj.domain = this.checkIPform.domain
                    }
                    this.cancaliplogVisible()
                    // this.page = 1
                    this.getIpData();
                    this.$message({
                        message:res.msg || '编辑成功',
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:res.msg || '编辑失败',
                        type: 'error'
                    });
                //关闭弹窗
                this.cancaliplogVisible()
                }  
            } else {
                let params = {
                    ip:this.checkIPform.ip,
                    hosts:this.checkIPform.domain,
                }
                const res = await auxiliarytool.addIplog(params)
                if(res.code == 200){
                    this.cancaliplogVisible()
                    this.page = 1
                    this.getIpData();
                    this.$message({
                        message:res.msg||'新增成功',
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                } 
            }
            

        },
         // 删除规则
         async handleDel (scope){
            let params = {
                id:[scope.row.id].join(',')
            }
            const res = await auxiliarytool.deleteipLog(params)
            if(res.code==200){ 
                this.$message({
                    message:res.msg || '删除成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                this.getIpData();
                
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        },
        async AllDel(){ //批量删除
			if(this.multipleSelection.length == 0) return;
    		var ids = [];
    		for (var i = 0; i < this.multipleSelection.length; i++) {
    			ids.push(this.multipleSelection[i].id);
    		}
            let params = {
                id:ids.join(',')
            }
            const res = await auxiliarytool.deleteipLog(params)
            if(res.code==200){ 
                this.$message({
                    message:res.msg||'删除成功',
                    type: 'success'
                });
                // this.alldelvisible = false;
                this.getIpData();
            }else{
                this.$message({
                    message:res.msg||'删除失败',
                    type: 'error'
                });
            }  
		},
         handleSelectionChange(val){
            this.multipleSelection = val;
        },
        newAdd(){
            this.title = '新增绑定'
            this.dialogFormVisible = true;
        },
        // 编辑
        handleInfo (row) {
            this.checkIPform.domain = row.hosts;
            this.checkIPform.ip = row.ip;
            this.checkIPform.id = row.id;
            this.title = '编辑绑定'
            this.dialogFormVisible = true;
        },
        handleReset(){
            this.search_field = '';
            this.page = 1;
            this.getIpData();
        },
        handlesearch(){
            this.page = 1; 
            this.getIpData();
            this.currentpage = 1;
        },
        currentchange(t){
            this.page = t; 
            this.currentpage = t;
            this.getIpData();
        },
        handleSizeChange(t) {
            this.page = 1;
            this.pageSize = t;
            this.getIpData();
        },
        mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
            if (this.$refs['popover_id-' + row.id]) {
                let t = this.$refs['popover_id-' + row.id].showPopper;
                if(!t){
                    this.showEditFileNameButton = false;
                this.rowId = "";
                }
            }
        },
        //  提交表单
        // saveAndSendData(){
        //     let valiDatas = this.$refs.edit.handleEdit();
        //     if(valiDatas === null)
        //         return;
        //     const data = this.handleFormDate(valiDatas);
        //     let config = {
        //         headers: {
        //             'Content-Type': 'multipart/form-data'
        //         }
        //     }
        //     this.$ajax({
        //         url:"/tools/assists/penetration_resources/create/",
        //         method:"POST",
        //         data
        //     })
        //     .then((data) => {
        //         let dt = data.data
        //         if (dt.success) {
        //             this.$message.success(dt.msg);
        //             this.createInfiltration = false;
        //             this.getPenetrationResources();
        //             this.$refs.edit.handleClearFiles();
        //             this.$refs.edit.form.remarks = ''

        //         } else {
        //             this.$message.error(dt.msg)
        //         }
        //     })
        //     .catch((data)=>{
        //         console.log(data); //错误信息
        //     });
        // },
    }
})
 
</script>

<style lang="less" scoped>    
    /deep/ .el-table{
        tr{
            height: 55px;
        }
         .cell {
         line-height: 14px;
        }
    }
    .bugbasicinfo{
        padding: 24px;
        .el-form-item{
            margin-bottom: 20px;
        }
        /deep/ .el-form-item__label{
            text-align: left;
            border-left: 3px solid #4C7AE3;
            line-height: 16px;
            padding-left: 8px;
            margin-top: 12px;
            position: relative;
            &:before{
                position: absolute;
                margin-left: 54px;
                margin-right: 0;
            }
        }
        .flexBet{
            position: relative;
            span:nth-child(2){
                position: absolute;
                right: 6px;
                top: 3px;
                color: #F56C6C;
            }
        }
    }
        .tooltable{
            .el-button--small{
                padding: 0;
            }
        }
       /deep/ .el-dialog{
            height: 192px !important;
        }
        .upload_dialog /deep/ .el-dialog{
            height: calc(100% - 96px) !important;
        }
       
        .dialogtxt{
            text-align: center;
            margin-top: 55px;
        }
        /deep/ .el-table td:not(.el-table-column--selection):first-child .cell, 
        /deep/ .el-table th:not(.el-table-column--selection):first-child .cell{
            padding-left: 32px !important;
        }
        /deep/ .el-tabs__item{
            height: 48px;
            line-height: 48px;
            padding: 0 24px;
        }
        /deep/ .el-tabs__item.is-active{
            color: #4C7AE3;
            font-weight: 500;
        }
        /deep/ .el-tabs__nav-wrap{
            padding: 0 24px; 
        }
        /deep/ .el-tabs__nav-wrap::after{
            background: #E8E8F5;
            height: 1px;
        }
        /deep/ .el-tabs__header{
            margin: 0 0 24px;
        }
        .auxiliarytool{ 
           background: #fff;
            min-height: calc(100% - 39px);
            box-sizing: border-box;
            // padding: 24px;
            /* box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12); */
            /* border-radius: 4px; */
             /deep/ .el-tabs__header{
                margin:0;
            }
            /deep/ .el-tabs__content{
                padding: 24px ; 
                background: #fff;
                box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
                border-radius: 4px;
                // min-height: 680px;
            }
        }	
        .auxiliarytool .el-tabs__item.is-top.is-active{
            background: #fff;
        }
        .auxiliarytool .el-tabs__header{
            margin:0;
        }
        .auxiliarytool .el-tabs__content{
            padding:20px 30px; 
            background: #fff;
            min-height: 500px;
        }
        .auxiliarytool .tabsbox{
            background: #fff;
        }
        .auxiliarytool .basic{
            background: #F7F7FB !important;;
            border:1px solid #e2e5ed;
            border-left: 2px solid #4c7ae3;
            padding: 5px 10px;
        }
        .auxiliarytool .basic > div{
            color: rgba(72, 72, 102, 0.64);
            font-weight: 500;
            margin:10px 8px;
            font-size: 13px;
        }
        .auxiliarytool .basic > div strong{
            display: inline-block;
            width: 80px;
            color: #4c7ae3;
            
        }
        .auxiliarytool .el-tabs__content{
            min-height: 650px;
        } 
        .auxiliarytool .el-date-editor.el-input, 
        .auxiliarytool .el-date-editor.el-input__inner{
            width: 100% !important;
        }
        .scriptbox{ 
            border:1px solid #ebebeb;
            
        }
        .scriptbox >strong{
            display: inline-block;
            width: 100%;
            background: #f2f3f9;
            padding: 10px 20px;
            box-sizing: border-box;
            font-size: 14px;
        }
        .scriptbox >div{
            width: 100%;
            padding: 10px 20px;
            box-sizing: border-box;
            overflow-y: auto;
            overflow-x: hidden;
        }
        .Buttonbox{
            text-align: center;
        }
        .Buttonbox >div{ 
            display: inline-block;
            margin: 5px 0;
        }
        .el-icon-error{
            vertical-align: middle;
            font-size: 14px;
            cursor: pointer;
        }
        /deep/ .shentoudia .el-form-item__label{
            text-align: left!important;;
        }
        /deep/ .shentoudia .el-form-item__label:after{
            left:-6px;
        }
        /deep/ .fingerValidate{
        .el-dialog{
            height: 400px!important;
        }
    }
    </style>
