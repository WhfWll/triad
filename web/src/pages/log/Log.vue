<template>
    <div>
        <div class="main-title  ">
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i>
            <label class="taskSceneBtn" >日志审计</label>
        </div>
        <div class="loglist context_box_bg"> 
            <div class="search-box"  >
                <div class="operationbutton" > 
					<el-button type="primary" @click="clearLog"  size="small">清空日志</el-button>
                    <!-- <el-button type="primary" @click="logset" size="small">日志配置</el-button> -->
                </div>
                <div class="serach-condition" >
                    <div> 
						<!-- <el-select v-model="formData.type"  style=" width:140px;" clearable placeholder="请选择日志类型" size="small">  
							<el-option
                            v-for="(item,index) in typelist"
                            :key="index"
                            :label="item.label"
                            :value="item.value">
                            </el-option>
						</el-select>  -->
	                </div>
                    <div > 
                        <el-date-picker
                            v-model="formData.time"
                            type="daterange"
                            format="yyyy-MM-dd"
                            value-format="yyyy-MM-dd HH:mm:ss"
                            :default-time="['00:00:00','23:59:59']"
                            range-separator="—"
                            start-placeholder="开始日期"
                            end-placeholder="结束日期" 
                            clearable
                            size="small"
                            >
                        </el-date-picker>
					</div>
					<div class="search-text">
						<el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch"  v-model="formData.search_field" class="input-with-select"  size="small" clearable > </el-input>
						<el-button type="primary"  size="small" @click="handlesearch">搜索</el-button>
					</div>
					<div >
						<el-button type="primary"  size="small" @click="handleReset">重置</el-button>
					</div> 
                </div>
            </div>
           
            <el-table
                ref="multipleTable"
                :data="tableData" 
                tooltip-effect="dark"  height="calc(100% - 102px)"
                v-model="Loading" class="myTable" v-if="typelist.length > 0">
                <el-table-column
                    label="日志类型">
                    <template slot-scope="scope" slot="header">
                        
                        <span class="cursorPointer" @click="clickButton('日志类型')"
                            :class="(formData.type !== '' && formData.type !== 0) ? 'active' : ''">日志类型<i
                            class="iconfont iconshaixuan"></i>
                        </span>
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.type" clearable
                        size="small" ref="loglistRef" @change="handlesearch">
                            <el-option v-for="(item, index) in typelist" :key="index" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                            {{scope.row.logTypeName}}
                    </template> 
                </el-table-column> 
                <el-table-column
                prop="content"
                label="日志内容"  width="500"> 
                </el-table-column>
                <el-table-column
                prop="updateTime"
                label="更新时间"
                 >
                </el-table-column>
                <el-table-column
                prop="username"
                label="操作用户"
                >
                </el-table-column>  
                <el-table-column
                prop="ip"
                label="登录IP"
                >
                </el-table-column>  
            </el-table>
            <el-pagination 
				background
                @size-change="handleSizeChange"
				layout=" total, prev, pager, next,sizes, jumper"
				:total="total"
                :page-size="pageSize"
				:current-page="currentpage"
				@current-change = "currentchange" >
			</el-pagination>
        </div>
        <el-dialog 
            title="日志配置" 
            :visible.sync="dialogVisible"  
            :before-close="cancelform" 
            width='1184px' 
            :close-on-click-modal="false" 
            :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="submitform">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div style="padding:24px"> 
                <el-form :model="logform" status-icon  ref="rulelog"   label-width="0">
                    <div>
                        <label class="dialog_item_label" style=" vertical-align: top; margin-top: 10px; width: 103px;">留存时长</label>
                        <el-form-item label="" prop="time" style="display: inline-block;"> 
                            <el-select v-model="logform.time"  size="small" style="width:112px"  >
                                <el-option
                                    v-for="(item,index) in timelist"
                                    :key="index"
                                    :label="item[1]"
                                    :value="item[0]">
                                </el-option>
                            </el-select>
                            月 
                        </el-form-item>
                    </div>
                    <div>
                        <label class="dialog_item_label" style=" vertical-align: top; margin-top: 10px; width: 103px;">留存个数</label>
                        <el-form-item label="" prop="count" style="display: inline-block;"> 
                            <el-select v-model="logform.count"  size="small" style="width:112px"  >
                                <el-option
                                    v-for="(item,index) in countlist"
                                    :key="index"
                                    :label="item[1]"
                                    :value="item[0]">
                                </el-option>
                            </el-select> 
                            条
                        </el-form-item>  
                    </div>
                    
                </el-form> 
            </div>
            
        </el-dialog> 
    </div> 
</template>
<style lang="less" scoped >
    .loglist{
        padding: 24px; 
        background: #fff;
        height: calc(100% - 39px);
        box-sizing: border-box;
        box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
    }	
  
    .dialog_b_btn{
        position: absolute;
        top: 15px;
        right: 24px;
        font-size: 14px;
        button{
            color: #4C7AE3;
        } 
    }
     @media (max-width: 1440px) {
         
        /deep/ .el-dialog{
            height: calc(100% - 96px);
        }
    }
    @media  (min-width: 1440px) { 
        /deep/ .el-dialog{
            height: calc(100% - 176px);
        }
    }
    .dialog_item_label{
		font-size: 14px;
		border-left: 3px solid #4C7AE3;
		padding-left: 8px; 
        font-weight:500;
        width: 113px;
        display: inline-block;
        // height: 18px;
        line-height: 16px; 
        box-sizing: border-box;
    } 
    /deep/ .myTable{
        thead {
            .cursorPointer{
                cursor: pointer;
                &.active{
                    color:#4C7AE3;
                    i{
                        color:#4C7AE3;
                    }
                }
            }
            .cell{
                line-height: 15px;
                >span{
                    position: absolute;
                }
            }
            .iconfont{
                color:rgba(72,72,102,0.32);
                margin-left:5px;
            }
            .el-select{
                height: 0;
                visibility: hidden;
                .el-input, .el-input__inner{
                    height: 0!important;
                }
            }
        }
    }
</style> 
<script>  
 import log from '@/api/log.js'
export default({
    name:'log',
    data(){ 
        return{
            formData:{
                type:0,
                content:'',
                time:'',
                start_time:'',
                end_time:'',
                page:1,
            },
            total:0,
            currentpage:1,
            typelist:[],
            
            typelist2:[],
            tableData:[],
            Loading:false,
            dialogVisible:false,
            logform:{
                time:'',
                count:'',
            },
            timelist:[],
            countlist:[],
            pageSize:10,
            timer1: null,
      		timermillisec:0 //时间间隔

        }
    },
    created() {
        this.$store.state.activefirstMenu="/log"; 
        this.pageSize = this.commonjs.pageSize;
		this.timermillisec = this.commonjs.timermillisec;
    },
    mounted() {
        this.getData();
        this.getType1();  
        
        this.timer1 = setInterval(() => {
			this.getData(true);
		}, this.timermillisec );
    },
    beforeDestroy() {
    //页面销毁清除定时器
    this.timer1 ? clearInterval(this.timer1) : "";
    clearTimeout(this.commonjs.timeer);
    this.commonjs.timeer = null;
  },
    methods: {
        icons(h, { column }) {
            let that = this
            return h(
              "div",

              [
                h("span", column.label),
                h(
                  "i",
                  {
                    slot: "reference",
                    class: "iconfont iconshaixuan",
                    style:"color:rgba(72,72,102,0.32);margin-left:5px;vertical-align:initial",
                    on: {
                        click: function() {
                            that.clickButton(column);
                        }
                    }
                  },
                  ""
                )
              ]
            );
          },
        // 图标点击事件
        clickButton(type) {
            switch (type) {
                case '日志类型':
                this.$refs.loglistRef.toggleMenu();
                break;
            }
        },
        async getType1(){ 
             const res = await log.logEmnu();
              if(res.code == 200){
                this.typelist = res.data.type;
                this.typelist.unshift(
                    { label: "全部", value: 0 }
                ) 
                 
              }else{

              }

            // this.$ajax.get('/smart/logs/enum',{
            //     params: {}
            // })
            // .then(dt => { 
            //     let res = dt.data;
            //     if(res.code === 200){
            //         // console.log(res.data,'res.data');
            //         this.typelist = res.data.type;
            //         console.log(this.typelist,'this.typelist');
                   
            //         this.typelist.unshift(
            //         { label: "全部", value: 0 }
            //         ) 
            //         console.log(this.typelist);
                   
            //         setTimeout(() => {
            //             this.typelist2 = this.typelist
            //         }, 1000);
                   
            //     }
			// 	else{
            //         this.$message({
            //             message:res.msg,
            //             type: 'error'
            //         });
            //     }
                 
            // })
            // .catch(data=>{ });
        },
        getData(){  
            this.Loading = true;  
            if(!this.formData.time){
				this.formData.start_time = '';
				this.formData.end_time = '';
			}
			else{
				this.formData.start_time = this.formData.time[0];
				this.formData.end_time = this.formData.time[1];
			} 
            let logtype = this.formData.type
            if (logtype === 0) {
                logtype = ''
            }
            this.$ajax.get('/smart/logs/logauditlist',{
                params: { 
                    page: this.formData.page, 
  					search:this.formData.search_field,
  					startTime:this.formData.start_time,
  					endTime:this.formData.end_time,
                    logType:logtype, 
                    size:this.pageSize,
                }
            })
            .then(dt => { 
                let res = dt.data;
                if(res.code === 200){
                    this.Loading = false; 
                    this.tableData = res.data.list; 
                    this.total = res.data.total;
                }
				else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
                 
            })
            .catch(data=>{ });
        },
        getconfig(){
            this.$ajax.get('/logs/auditLogs/get/settings/',{
                params: {}
            })
            .then(dt => { 
                let res = dt.data;
                if(res.success){
                    this.timelist = res.param.time;
                    this.countlist = res.param.count;
                    this.logform.time = res.data.save_time;
                    this.logform.count = res.data.save_count;
                    
                }
				else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
                 
            })
            .catch(data=>{ });
        },
        currentchange(t){
            this.formData.page = t; 
            this.getData();
            this.currentpage = t;
        },
        handleSizeChange(t){
            this.formData.page = 1;
            this.pageSize = t;
            this.getData();
        },
        handlesearch(){
            this.formData.page = 1;
            this.currentpage = 1;
            this.getData();
        },
        handleReset(){
            this.formData.search_field = '';
            this.formData.type = 0;
            this.formData.time = '';
            this.formData.page = 1; 
            this.currentpage = 1;
            this.pageSize = 10;
            this.getData();
        },
        clearLog(){
            this.$confirm('确定要清空日志？', '清空日志', {
                distinguishCancelAndClose: true,
                confirmButtonText: '确定',
                cancelButtonText: '取消',  
            }).then(() => {  
				this.$ajax({
	                method:'post',
	                url:'/smart/logs/logauditempty',
	                data:''
	            }) 
                .then(dt => { 
                    let res = dt.data;
                    if(!res.error){
                        this.$message({
                            message:'清空日志成功!',
                            type: 'success'
                        }); 
                        this.formData.page = 1;
                        this.formData.start_time='';
                        this.formData.end_time='';
                        this.formData.type= 0;
                        this.formData.content='';
                        this.getData();
                        this.dialogVisible = false;
                        this.$refs.rulelog.resetFields(); 
                    }
                    else{
                        this.$message({
                            message:res.error,
                            type: 'error'
                        });
                    }
                    
                })
                .catch(data=>{ });
	        }).catch(action => {
	                  
            }); 
        },
        logset(){ //日志配置
            this.dialogVisible = true;
            this.getconfig();
        },
        cancelform(){
            this.dialogVisible = false;
            this.$refs.rulelog.resetFields(); 
        },
        submitform(){ ///
            this.$ajax({
                method:'post',
                url:'/logs/auditLogs/settings/',
                data:this.qs.stringify({
                    save_time:this.logform.time,
                    save_count:this.logform.count
                })
            }) 
            .then(dt =>{
                let res = dt.data;
                if(res.success){
                    this.$message({
                        message:res.msg,
                        type: 'success'
                    });
                    this.dialogVisible = false;
                    this.$refs.rulelog.resetFields(); 
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            }).catch(err=>{})
        },
    },
})
</script>