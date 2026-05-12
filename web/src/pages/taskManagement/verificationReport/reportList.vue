<template>
    <div>
        <div class="main-title  ">
            <label for="">验证报告</label>
        </div>
        <div class="tasklist context_box_bg">
            <div class="search-box">
                <div class="operationbutton">
                    <!-- <el-button type="primary"  size="small" @click="btReportTesting">验证报告</el-button>  -->
                    <xzbutton type="primary" @click="btReportTesting" size="small">验证报告</xzbutton>
                    <!-- <el-popover
                    popper-class="delButton_popper"
                    placement="bottom-start"
                    width="170"
                    style="padding-left:8px"
                    trigger="click" 
                    :visible-arrow="false"
                    v-model="alldelvisible" >
                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                    <div style="text-align: right; margin: 0" class="" >
                        <el-button size="mini" class="delCancel" @click="alldelvisible = false" >取消</el-button>
                        <el-button size="mini" type="primary"  @click="btnMultiDeleteTask">确定</el-button>
                    </div>  
                    <el-button type="warning" size="small"  slot="reference" :disabled="!multipleSelection.length">删除</el-button> 
                </el-popover>   -->
                    <delbutton :width="170" @click="btnMultiDeleteTask" :disabled="!multipleSelection.length" style="margin-left: 8px;">
                    </delbutton>
                </div>
                <div class="serach-condition">
                    <!-- <div>
                        <el-select v-model="formData.task_level" style=" width:140px" placeholder="任务风险" size="small">
                            <el-option v-for="(item,index) in taskRiskList" :key="index" :label="item[1]"
                                :value="item[0]">
                            </el-option>
                        </el-select>
                    </div>
                    <div>
                        <el-select v-model="formData.manufacturer" style=" width:140px" placeholder="验证厂商" size="small">
                            <el-option v-for="(item,index) in taskstatus" :key="index" :label="item.label"
                                :value="item.value">
                            </el-option>
                        </el-select>
                    </div> -->
                    <div>
                        <el-date-picker v-model="formData.time" type="daterange" format="yyyy-MM-dd"
                            value-format="yyyy-MM-dd HH:mm:ss" :default-time="['00:00:00','23:59:59']"
                            range-separator="—" start-placeholder="开始日期" end-placeholder="结束日期" size="small" clearable>
                        </el-date-picker>
                    </div>
                    <div class="search-text">
                        <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearch" v-model="formData.search_field" class="input-with-select"
                            size="small" clearable> </el-input>
                        <!-- <el-button type="primary"  size="small" @click="handlesearch">搜索</el-button> -->
                        <xzbutton type="primary" @click="handlesearch" :disabled="false" size="small">搜索</xzbutton>
                    </div>
                    <div>
                        <!-- <el-button type="primary"  size="small" @click="handleReset">重置</el-button> -->
                        <xzbutton type="primary" @click="handleReset" :disabled="false" size="small">重置</xzbutton>
                    </div>
                </div>


            </div>
            <!-- 表格主体 -->
            <el-table :data="tableData" @select="handChecked" @select-all="handSelectAll"  
                @cell-mouse-enter="mouseenter" 
                @cell-mouse-leave="mouseleave" 
                @selection-change="handleSelectionChange" 
                :row-key="rowKeyRe" ref="multipleTable" tooltip-effect="dark"
                style="width: 100%">
                <el-table-column type="selection" width="55" :reserve-selection="true">
                </el-table-column>
                <el-table-column prop="name" label="任务名称">
                    <template slot-scope="scope">
                        <el-link @click="btnTaskinfo(scope.row)">{{scope.row.name}}</el-link>
                    </template>
                </el-table-column>
                <el-table-column prop="executeTypeName" label="执行方式"> 
                </el-table-column>
                <!--  :render-header="icon" -->
                <el-table-column prop="name" label="核实漏洞" show-overflow-tooltip>
                    <template slot="header">
                        核实漏洞
                        <el-tooltip class="item" effect="dark" placement="right">
                            <div slot="content">
                                从左至右依次是“利用成功”、“验证成功”、“验证失败” <br />
                                利用成功 : 通过EXP验证发现漏洞存在 <br />
                                验证成功 : 通过POC验证发现漏洞存在 <br />
                                验证失败 : 通过POC验证发现漏洞不存在<br />
                            </div>
                            <i class="iconfont icontishi" style="position: absolute;top:0;left:66px"></i>
                        </el-tooltip>
                    </template>
                    <template slot-scope="scope">
                        <span class="tag_status tag_danger bug_status" v-if="scope.row.exp >99">99+</span>
                        <span class="tag_status tag_danger bug_status"
                            v-if="scope.row.exp <=99">{{scope.row.exp}}</span>
                        <span class="tag_status tag_warning bug_status" v-if="scope.row.verify >99">99+</span>
                        <span class="tag_status tag_warning bug_status"
                            v-if="scope.row.verify <=99">{{scope.row.verify}}</span>
                        <span class="tag_status tag_primary bug_status" v-if="scope.row.failed >99">99+</span>
                        <span class="tag_status tag_primary bug_status"
                            v-if="scope.row.failed <=99">{{scope.row.failed}}</span>
                    </template>
                </el-table-column>
                <el-table-column prop="name" label="任务风险">
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('任务风险')"
                            :class="(formData.task_level !== '' && formData.task_level !== 0) ? 'active' : ''">任务风险<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.task_level" clearable
                        size="small" ref="loglistRef" @change="handlesearch">
                        <el-option v-for="(item, index) in risklevellist" :key="index" :label="item.label" :value="item.value">
                        </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                        <span :class="[ 
                    {'riskstyle risk_hight': scope.row.risk=='1' } ,
                    {'riskstyle risk_middle': scope.row.risk=='2' },
                    {'riskstyle risk_low':scope.row.risk =='3' },
                    {'riskstyle risk_nofind':scope.row.risk =='4' }]"><i></i>{{scope.row.riskName }}</span>
                    </template>
                    <!-- <template slot-scope="scope">    
                        <div style="    line-height: 14px;">{{scope.row.risk}}</div>
                    </template>   -->
                </el-table-column>
                <el-table-column prop="name" label="验证厂商">
                    <template slot-scope="scope" slot="header"> 
                        <span class="cursorPointer" @click="clickButton('验证厂商')"
                            :class="(formData.manufacturer !== '' && formData.manufacturer !== 0) ? 'active' : ''">验证厂商<i
                            class="iconfont iconshaixuan"></i>
                        </span> 
                        <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.manufacturer" clearable
                            size="small" ref="manufacturer" @change="handlesearch">
                            <el-option v-for="(item, index) in manufacturerlist" :key="index" :label="item.label" :value="item.value">
                            </el-option>
                        </el-select>
                    </template>
                    <template slot-scope="scope">
                        {{ scope.row.producerName }}
                        <!-- <div style="line-height: 14px;" v-if="scope.row.product == 1">启明天镜</div>
                        <div style="line-height: 14px;" v-if="scope.row.product == 2">绿盟nsfocus</div> -->
                    </template>
                </el-table-column>
                <el-table-column prop="name" label="更新时间" :show-overflow-tooltip="true">
                    <template slot-scope="scope">
                        <div style="    line-height: 14px;">{{scope.row.updateTime}}</div>
                    </template>
                </el-table-column>
                <el-table-column prop="name" label="状态">
                    <template slot-scope="scope">
                     <!-- 待触发  1 ；待执行  2；运行中   3；已完成   4 ；暂停中   5  -->
                    <div v-if="showOperateButton && rowId == scope.row.id  ">
                        <!-- 待触发  1 ； -->
                        <div v-if="scope.row.status == 1">
                            <el-link class="link_danger" :underline="false" @click="btnstopcrontask(scope.row.id)">
                                结束</el-link>
                            <el-popover placement="bottom" width="170" :visible-arrow="false"
                                :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel"
                                        @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消
                                    </el-button>
                                    <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                </div>
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0"
                                    slot="reference">删除</el-link>
                            </el-popover>
                        </div>
                        <!--  待执行  2；-->
                        <div v-if="scope.row.status == 2">
                            <el-link class="link_danger" :underline="false" @click="btnstopcrontask(scope.row.id)">
                                结束</el-link>
                            <el-popover placement="bottom" width="170" :visible-arrow="false"
                                :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel"
                                        @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消
                                    </el-button>
                                    <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                </div>
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0"
                                    slot="reference">删除</el-link>
                            </el-popover>
                        </div>
                        <!-- 运行中   3； -->
                        <div v-if="scope.row.status == 3 ">
                            <el-link class="link_danger" :underline="false" @click="btnstopcrontask(scope.row.id)">
                                结束</el-link>
                            <el-popover placement="bottom" width="60" trigger="click" popper-class="learnMore"
                                :visible-arrow="false" style="padding:0">
                                <ul class="operationbox">
                                    <!-- <li @click="btnPauseTask(scope.row.id)">暂停 </li> -->
                                    <li>
                                        <el-popover placement="bottom" width="170" :visible-arrow="false"
                                            :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                                            <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                            <div style="text-align: right; margin: 0">
                                                <el-button size="mini" class="delCancel"
                                                    @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">
                                                    取消</el-button>
                                                <el-button size="mini" type="primary" @click="btnDel(scope)">确定
                                                </el-button>
                                            </div>
                                            <span slot="reference" style="cursor:pointer">删除</span>
                                        </el-popover>
                                    </li>
                                </ul>
                                <el-link :underline="false" class="link_info" slot="reference">更多</el-link>
                            </el-popover>
                        </div>
                        <!-- 已结束 -->
                        <div v-if="scope.row.status == 4 ">
                            <el-popover placement="bottom" width="170" :visible-arrow="false"
                                :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel"
                                        @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消
                                    </el-button>
                                    <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                </div>
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0"
                                    slot="reference">删除</el-link>
                            </el-popover>
                        </div>
                    </div>
                    <div v-else> 
                        <span :class="[ 
                            { 'tag_status tag_danger1': scope.row.status == 1 } ,
                            { 'tag_status tag_warning': scope.row.status ==2 },
                            { 'tag_status tag_primary': scope.row.status == 3 },
                            { 'tag_status tag_success': scope.row.status ==4 },
                        { 'tag_status tag_danger': scope.row.status == 5 }]"><i></i>{{ scope.row.statusName }}</span>
                    </div> 
                </template>
                </el-table-column> 
            </el-table>
            <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange"
                :current-page="currentPage" :page-size="pageSize" layout=" total,  prev, pager, next, sizes,jumper"
                :total="totalpage">
            </el-pagination>
        </div>
    </div>
</template>
<style scoped lang="less">
/deep/.el-table thead .cell{
    overflow: initial !important;
}
/deep/ thead {
  .cursorPointer {
    cursor: pointer; 
    position: absolute;
      // top: 6px;
    &.active {
      color: #4C7AE3;

      i {
        color: #4C7AE3;
      }
    }
  }

  .cell {
    line-height: 15px;

    >span {
      position: absolute;
    }
  }

  .iconfont {
    color: rgba(72, 72, 102, 0.32);
    margin-left: 5px;
  }

  .el-select {
    height: 0;
    visibility: hidden;

    .el-input,
    .el-input__inner {
      height: 0 !important;
    }
  }
}
 
.context_box_bg{
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
}

.loophole1{color:orange}
.loophole2{color:red}
.loophole3{color:darkgray}
</style>

<script>
import xzbutton from "@/components/XzButton.vue"; 
import delbutton from "@/components/DelButton.vue";
import { task, target } from '@/api/task.js'
export default({
    name:'verificationReport',
    components: {
    	xzbutton,
		delbutton, 
  	},
    data(){
        var reportName = (rule, value, callback) =>{
            if (!value) {
              return callback(new Error('任务名称不能为空'));
            }else{
                const pattern = /^[a-zA-Z\d\_\u2E80-\u9FFF]{0,8}$/
                if(pattern.test(value)){
                    callback();
                }else{
                    return callback(new Error('任务名称格式有误'));
                }
            }
        }
        return {
            taskRiskList:[
                [1,'高危'],[2,'中危'],[3,'低危'],[4,'安全']
            ],
            manufacturerlist:[
                // {
                //     value: '',
                //     label: '全部'
                // },
                // {
                //     value: '1',
                //     label: '启明天镜'
                // },
                // {
                //     value: '2',
                //     label: '绿盟nsfocus'
                // },
            ],
            formData:{
                task_level:'',
                manufacturer:'',
                time:'',
                starttime:'',
                endtime:'',
                search_field:'',
                page_num:1,
            },
            tableData:[
                // {task_status:2,name:'192.168.0.13...扫描任务',type:'即时执行',high_number:'7',middle_number:'4',low_number:'25',risk_level:'1',risk_level_name:'高危',manufacturer:'绿盟',time:'2020-08-14 10:37',operation_Status:'2',task_status_number:'2'}
            ],
            currentPage:1,
            pageSize:10,
            totalpage:0,
            dialogReportTesting:false,
            typeOptions:[
                {
                    value: '1',
                    label: '立即执行'
                },
                {
                    value: '2',
                    label: '定时执行'
                },
            ],
            importReoprtForm:{
				name:'',
                type:'1',
				reports:'',
			},
			fileList:[],
			reportRule:{
				name:[
					{ required: true, trigger: 'blur',validator:reportName  },
				],
                type:[
                    {required: true, message: '请选择执行方式', trigger: 'blur'}
                ],
				reports:[
					{required: true, message: '请上传扫描报告', trigger: 'blur'}
				]
			},
    		multipleSelection:[],
            alldelvisible:false,
            singleDeleteVisable:false,
            timer: null,
            timer1:null,
            checkedList:[],
            timermillisec:0, //时间间隔
            risklevellist:[],
            showOperateButton:false,
            rowId:'',
        }
    },
    mounted () { 
        this.getEnum();
        // this.getData();
    },
    created:function(){
        // this.timermillisec = this.commonjs.timermillisec;
        this.$store.state.activefirstMenu="/verificationReport";  
        this.pageSize = this.$commonjs.pageSize;
        //  this.timer1 = setInterval(()=>{
            this.getData()
        // },this.timermillisec)
        
    },
    beforeDestroy () {
        // clearInterval(this.timer)
        // clearInterval(this.timer1)
    },
    methods:{
        getEnum() { 
            this.$ajax.get('/smart/reportverify/enum').
            then((res) =>{
                var res = res.data;    
                if(res.code == 200){ 
                    this.risklevellist = res.data.risk;
                    this.risklevellist.unshift( { label: "全部", value: 0 });
                    this.manufacturerlist = res.data.producerType;
                    this.manufacturerlist.unshift( { label: "全部", value: 0 });
                }
				else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }
            }).catch((error) =>{
                console.log(error);
            });
            
        },
        // 图标点击事件
        clickButton(type) {
            switch (type) {
                case '任务风险':
                    this.$refs.loglistRef.toggleMenu();
                    break;
                case '验证厂商':
                    this.$refs.manufacturer.toggleMenu();
                    break;
            }
        },
        rowKeyRe(row){
            return row.id
        },
        handChecked(sel,row){
          this.checkedList = sel
        },
        handSelectAll(ael){
          this.checkedList = ael
        },
        icon(h,{column}){
            const inReview = '从左至右依次是“利用成功”、“验证成功”、“验证失败”； 利用成功，通过EXP验证发现漏洞存在；验证成功，通过POC验证发现漏洞存在；验证失败：通过POC验证发现漏洞不存在；' 
			return h('div', {
					style: { 
							'padding-left':' 0 !important',
							'height': '16px',
							'line-height': '16px',
							'overflow': 'initial',
						}
					},
					[ h('span', column.label),
                    h('el-tooltip', {
                        props: {
                            placement: 'top'
                        }, 
                    }, [
                        h('div', {
                            slot: 'content',
                            style: {
                                // 'width':'100px',
                                whiteSpace: 'normal', 
                            },
                            innerHTML:'123'
                        }, inReview), 
                        h('i', {
                            class: 'iconfont icontishi',
                            style: 'color:rgba(72,72,102,0.32);margin-left:5px;vertical-align: initial;'
                        })
                    ],)
                ],
    　　　　 )
        },
        getData(){
            this.Loading = true; 
			if(!this.formData.time){
				this.formData.starttime = '';
				this.formData.stoptime = '';
			}
			else{
				this.formData.starttime = this.formData.time[0];
				this.formData.endtime = this.formData.time[1];
			} 
            this.$ajax.get('/smart/reportverify/tasklist',{
                params:{
                    risk:this.formData.task_level,
                    producer:this.formData.manufacturer,
                    startTime:this.formData.starttime,
                    endTime:this.formData.endtime,
                    search:this.formData.search_field,
                    size:this.pageSize,
                    page:this.formData.page_num
                }
            }).
            then((res) =>{
                var res = res.data;   
                this.Loading = false; 
                if(res.code == 200){
              
                        this.tableData = res.data.list; 
                        this.totalpage = res.data.total;
                        // if (res.is_refresh === 1) {
                        //     this.timer = setInterval(() => {
                        //         this.getData2()
                        //     }, this.timermillisec)
                        // }
                }
				else{
                    this.$message({
                        message:res.msg,
                        type: 'error'
                    });
                }
            }).catch((error) =>{
                console.log(error);
            });
            
        },
        getData2(){
            this.Loading = true; 
			if(!this.formData.time){
				this.formData.starttime = '';
				this.formData.stoptime = '';
			}
			else{
				this.formData.starttime = this.formData.time[0];
				this.formData.endtime = this.formData.time[1];
			} 
            this.$ajax.get('/smart/reportverify/tasklist',{
                params:{
                    risk:this.formData.task_level,
                    producer:this.formData.task_status,
                    start_date:this.formData.starttime,
                    end_date:this.formData.endtime,
                    search:this.formData.search_field,
                    size:this.pageSize,
                    page:this.formData.page_num
                }
            }).
            then((res) =>{
                var res = res.data;   
                this.Loading = false; 
                if(res.success){ 
                    this.tableData = res.results; 
                    this.totalpage = res.count;
                    if (res.is_refresh !== 1) {
                        clearInterval(this.timer)
                    }
                }
				else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                    clearInterval(this.timer)
                }
            }).catch((error) =>{
                console.log(error);
            });
        },
        handleSizeChange(t){
            this.formData.page_num = 1;
            this.pageSize = t;
            this.getData();
        },
		handleCurrentChange(t){
            this.formData.page_num = t;
            this.getData();
		},
        btnDel(scope){ //单个删除任务 
            this.$ajax.get('/smart/reportverify/taskdelete',{
                params:{
                    taskId: scope.row.id
                }
            })
            .then((dt) => { 
                let res = dt.data;
                if(res.code == 200){
					this.$message({
                        message:'删除任务成功',
                        type: 'success'
					});
					scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
                    this.getData();
                }else{
					this.$message({
                        message:dt.msg,
                        type: 'error'
                    });
				}
            })
            .catch((error) => {
                console.log(error);
            })
		},
        btReportTesting(){
            // this.$ajax.get(('/v2/task/report_test/generate_task_id'),{
            //     params:{}
            // }).then((data) =>{
            //     var res = data.data;  
            //     if(res.code == 200){
            //         this.reportId = res.data.task_id;
                    this.$router.push({
                        path: `uploadReport`,
                        query: {
                            id: this.reportId
                        }
                    })
                // }
                
            // }).catch((error) =>{
            //     console.log(error);
            // })
        },
        handleExceed(files, fileList) {
			const fileType = file.name.substring(file.name.lastIndexOf(".") + 1);
			let isFileType = '' ;
			if(fileType == 'html' || fileType == 'zip'){
				isFileType = fileType;
			}
			const isLt100M = file.size / 1024 / 1024 < 100;
			if (!isFileType) {
				this.$message.error('上传文件只能是html/zip格式!'); 
			}
			if (!isLt100M) {
				this.$message.error('上传文件大小不能超过 100MB!');
			}
			// return isFileType && isLt100M;
        	this.$message.warning(`当前限制选择 20 个文件，本次选择了 ${files.length} 个文件，共选择了 ${files.length + fileList.length} 个文件`);
        },
		onBeforeUpload(file){//上传前验证格式与文件大小
			const fileType = file.name.substring(file.name.lastIndexOf(".") + 1);
			const isType = fileType === "html'||'zip";
			const isLt100M = file.size / 1024 / 1024 < 1;
			if (!isType) {
				this.$message.error('上传文件只能是html/zip格式!');
			}
			if (!isLt100M) {
				this.$message.error('上传文件大小不能超过 1MB!');
			}
			return isType && isLt100M;
		},
		fileChange(file,files){
			// const fileType = file.name.substring(file.name.lastIndexOf(".") + 1);
			// let isFileType = '' ;
			// if(fileType == 'html' || fileType == 'zip'){
			// 	isFileType = fileType;
			// }
			// const isLt100M = file.size / 1024 / 1024 < 1;
			// if (!isFileType) {
			// 	this.$message.error('上传文件只能是html/zip格式!');
			// }
			// if (!isLt100M) {
			// 	this.$message.error('上传文件大小不能超过 1MB!');
			// }
			// return isFileType && isLt100M;
			// console.log(file)
		},
		handlecancelReport(){
			this.dialogReportTesting = false;
			this.$refs.reportsForm.resetFields();
			this.$refs.reportRef.clearFiles();
		},
		btnStartReport(){
			// let formData = new FormData();
			// formData.append('file', this.form.file);
			// this.$ajax({
			// 	url:'',
			// 	method:'post',
			// 	data:''
			// }).then((data) =>{

			// }).catch((error) =>{

			// })
		},
        handlesearch(){
            this.formData.page_num = 1;
    		this.getData();
    		this.currentpage = 1;
        },
        handleReset(){
            this.formData.page_num =1;
			this.formData.search_field='';
			this.formData.time = ''; 
			this.formData.starttime = '';
            this.formData.endtime = '';
			this.formData.task_status ='';
            this.formData.task_level = '';
            this.formData.manufacturer='';
			this.pageSize = 10;
			this.currentpage = 1;
			this.getData();
        },
        btnTaskinfo(row){
            this.$router.push({
                path:`/reportTaskInfo`,
                query: { 
					id: row.id, 
                    risk_level:row.risk,
                    name:row.name,
                    status:row.status
                }
            })
            localStorage.setItem('taskTab', 'tabs1');
			localStorage.setItem('task_id', row.id );
        },
        btnstopcrontask(id){ //结束
            this.$ajax({
                url:'/smart/reportverify/taskstop',
                method:'get',
                params:{
                    taskId:id
                }
            }).
            then((res) =>{
                let dt = res.data;
                if(dt.code == 200){
                    this.$message({
                      message:'任务结束成功',
                      type: 'success'
                    });
                    this.getData();
                }else{
                    this.$message({
                        message:dt.msg,
                        type: 'error'
                    });
                }
            }).
            catch((error) =>{
                console.log(error); 
            })
        },
        handleSelectionChange:function(val){ 
    		this.multipleSelection = val;
    	},
        btnMultiDeleteTask(){ //批量删除
			if(this.multipleSelection.length == 0) return;
    		let _ids = [];
    		for (var i = 0; i < this.multipleSelection.length; i++) {
    			_ids.push(this.multipleSelection[i].id);
			}
            let stringIds = _ids.join(',')
			// this.$ajax({
            //     method:'delete',
            //     url:'/reportverify/v1/tasks/delete/?task_id='+stringIds,
            //     data:{
            //        ids:_ids.join(',')  
            //     } 
            // }) 
            this.$ajax.get('/smart/reportverify/taskdelete',{
                params:{
                    taskId: _ids.join(',')
                }
            })
            .then((dt) => { 
                let res = dt.data;
                if(res.code ==200){
					this.$message({
                        message:'删除任务成功',
                        type: 'success'
					});
					this.alldelvisible = false;
                    this.getData();
                }else{
					this.$message({
                        message:dt.msg,
                        type: 'error'
                    });
				}
            })
            .catch((error) => {
                console.log(error);
            })
		},
        singleDelete(){

        },
        mouseenter(row, colum, cell, event) {
            this.showOperateButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断  
        },
        mouseleave(row, colum, cell, event) { 
            
            if (!this.$refs['popover_id-' + row.id]){
            this.showOperateButton = false;
            this.rowId = "";
            return;
            }else{
            let isShow = this.$refs['popover_id-' + row.id].showPopper;
            if (!isShow) {
                this.showOperateButton = false;
                this.rowId = "";
            }
            } 
             
        },
    }
})
</script>