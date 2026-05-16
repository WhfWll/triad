<template>
    <div>
        <div class="main-title  ">
            Agent管理
        </div>
        <div class="list_box">
            <div class="search-box">
                <div class="operationbutton">
                    <el-button type="primary" size="small" style="margin-right:8px" @click="downloadagent">下载agent</el-button> 
                </div>
                <div class="serach-condition">
                    
                    <div class="search-text">
                        <el-input placeholder="搜索关键字" @keydown.enter.native="handlesearch" v-model="search_item.search" class="input-with-select"
                            size="small" clearable> </el-input>
                        <el-button type="primary" @click="handlesearch" size="small">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" @click="handleReset" size="small">重置</el-button>
                    </div>
                </div>
            </div>


            <el-table ref="multipleTable" :data="tableData" tooltip-effect="dark" v-model="Loading" style="width: 100%" class="myTable"
                @selection-change="handleSelectionChange" @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave">
                <el-table-column type="selection" width="55">
                </el-table-column>
                <el-table-column prop="name" label="Agent名称">
                    <template slot-scope="scope"> 
                        <el-link>{{scope.row.name}}</el-link>
                    </template>
                </el-table-column>
                <el-table-column prop="ip" label="IP"></el-table-column> 
                <el-table-column prop="onlineStatusEnum" label="在线状态"></el-table-column>
                <el-table-column prop="statusEnum" label="节点状态">
<!--                     
                </el-table-column>
                <el-table-column label="操作"> -->
                    <template slot-scope="scope">
                        <div v-if="showEditFileNameButton && rowId == scope.row.id">
                            <el-link :underline="false" class="link_primary" @click="btnNodeAuthorize(scope,1)"
                                v-if="scope.row.status === 2">启用</el-link>
                            <el-link :underline="false" class="link_primary" @click="btnNodeAuthorize(scope,2)"
                                v-if="scope.row.status === 1">禁用</el-link>

                            <!-- <el-popover placement="bottom" width="170" :visible-arrow="false"
                                :ref="`popover_id-${scope.row.id}`" popper-class="delButton_popper">
                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                <div style="text-align: right; margin: 0">
                                    <el-button size="mini" class="delCancel"
                                        @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                    <el-button size="mini" type="primary" @click="btnDel(scope)">确定</el-button>
                                </div>
                                <el-link :underline="false" class="link_danger linkafter" style="padding:0"
                                    slot="reference">删除</el-link>
                            </el-popover> -->
                        </div>
                        <div v-else >
							{{scope.row.statusEnum}}
						</div>
                    </template>
                </el-table-column>
            </el-table>
            <el-pagination background @size-change="handleSizeChange" @current-change="handleCurrentChange"
                :current-page="currentPage" :page-size="pageSize" layout=" total,  prev, pager, next, sizes,jumper"
                :total="total">
            </el-pagination>


        </div>
        <el-dialog
          title="添加节点"
          :visible.sync="dialogFormVisible"
          width="1184px"
          class="fingerValidate" 
          :close-on-click-modal="false" 
          :show-close="false">
          <div class="dialog_b_btn">  
            <el-button size="small" @click="saveAdd()">确定</el-button>
              <el-button size="small" @click="cancaliplogVisible()">关闭</el-button>
          </div>
          <div class="buginfo_box" style="padding: 20px;box-sizing:border-box;"> 
              <div class="bugbasicinfo">
                  <el-form :model="checkIPform" status-icon  ref="ruleFormaddip" :rules="rules" label-width="80px">
                    <el-form-item label="节点名称" prop="name">
                        <el-input v-model="checkIPform.name" size="small" style="width:520px" autocomplete="off" placeholder="请输入节点名称"></el-input>
                    </el-form-item>
                    <el-form-item prop="ip" label="节点IP">
                    <el-input v-model="checkIPform.ip" size="small" style="width:520px" autocomplete="off"
                            placeholder="请输入节点ip"></el-input>
                    </el-form-item>
                    <el-form-item prop="port" label="节点端口">
                    <el-input v-model="checkIPform.port" size="small" style="width:520px" autocomplete="off"
                            placeholder="请输入节点端口"></el-input>
                    </el-form-item>
              </el-form>
              </div>   
          </div>
      </el-dialog>

        <el-dialog :title="basicinfo1.name" :visible.sync="dialogVisible" width="1184px" class="buginfobox"
            :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">


                <el-button size="small" @click="btnUpdate" v-if="!is_Update">{{updatetxt}}</el-button>
                <el-button size="small" @click="saveUpdate" v-if="is_Update">保存</el-button>
                <el-button size="small" @click="cancaldialogVisible">关闭</el-button>
                <!-- <el-button size="small" @click="dialogVisible = false">关闭</el-button> -->
            </div>
            <div class="info_box">
                <div class="basicinfo">
                    <el-table :data="basicinfo" ref="myTable" size='small' style="width: 100%">
                        <el-table-column prop="name" width="300" label="节点名称">
                            <template slot-scope="scope">
                                <span v-if="!is_Update">{{ scope.row.name }}</span>
                                <span style="display:inline-block" v-else>
                                    <el-input v-model="nodeform.name" size="mini" maxlength="20"></el-input>
                                </span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="id" label="节点ID">

                        </el-table-column>
                        <el-table-column prop="ip" width="200" label="节点IP">
                        </el-table-column>
                        <el-table-column prop="tasks" label="运行任务数">
                        </el-table-column>
                        <el-table-column prop="online" label="在线状态">
                            <template slot-scope="scope">
                                <span class="tag_status tag_warning" v-if="!scope.row.online">离线</span>

                                <span class="tag_status tag_primary" v-if="scope.row.online">在线</span>
                            </template>
                        </el-table-column>
                        <el-table-column prop="status" label="节点状态">
                            <template slot-scope="scope">
                                <span class="tag_status tag_success" v-if="scope.row.status">已启用</span>

                                <span class="tag_status tag_danger" v-else>已禁用</span>
                            </template>
                        </el-table-column>
                    </el-table>
                </div>
                <!-- 
                    经讨论，被动流量暂时注释掉，by wdh 2022-7-21
                <div class="otherinfo">
                    <div class="part_title">
                        <label>被动流量扫描</label>

                        <span v-if="!is_Update"><span v-if="basicinfo1.agent">开启</span><span v-else>关闭</span></span>
                        <span v-else>
                            <el-switch v-model="nodeform.agent" class="elSwitch">
                            </el-switch>
                        </span>
                    </div>
                    <div>
                        <label>代理网卡：</label>
                        <span v-if="!is_Update">{{ basicinfo1.agentCard }}</span>
                        <span v-else>
                            <el-input v-model="nodeform.agentCard" size="mini" maxlength="20"></el-input>
                        </span>
                    </div>

                    <div>
                        <label>代理端口：</label>
                        <span v-if="!is_Update">{{ basicinfo1.agentPort }}</span>
                        <span v-else>
                            <el-input v-model="nodeform.agentPort" size="mini" maxlength="20"></el-input>
                        </span>

                    </div>

                    <div>
                        <label for="" class="" style="display:inline-block;margin-bottom: 8px;">客户访客白名单</label>
                        <div class="bgcolorfff minH" v-if="!is_Update">{{ basicinfo1.whitelist }}</div>
                        <el-input class="textarea" type="textarea" v-model="nodeform.whitelist" size="mini" :row="3"
                            v-if="is_Update"></el-input>
                    </div>

                </div> -->
                <!-- <div class="otherinfo">
                    <div class="part_title">
                        <label>限定目标测试</label>
                        <span v-if="!is_Update"><span v-if="basicinfo1.limit">开启</span><span v-else>关闭</span></span>
                        <span v-else>
                            <el-switch v-model="nodeform.limit" class="elSwitch">
                            </el-switch>
                        </span>

                    </div>
                    <div class="bgcolorfff minH" v-if="!is_Update">{{ basicinfo1.targets }}</div>
                    <el-input class="textarea" type="textarea" v-model="nodeform.targets" size="mini" :row="3" v-else>
                    </el-input>
                </div> -->
                <div class="otherinfo otherinfo_1">
                  
                    <div>
                        <label>代理渗透：</label>
                        <span v-if="basicinfo1.node_status == ''">未开启</span> 
                        <span v-else-if="basicinfo1.node_status == 'proxy'">开启</span>
                        <span v-else>未开启</span>
                        <span v-if="basicinfo1.node_status == 'proxy'"></span>
                    </div>
                    <div>
                        <label>OpenVPN：</label>
                        <span v-if="basicinfo1.node_status == ''">未开启</span>
                        <span v-else-if="basicinfo1.node_status == 'openvpn'">开启</span>
                        <span v-else>未开启</span>
                    </div>
                    <div>
                        <label>被动流量扫描：</label>
                        <span v-if="basicinfo1.node_status == ''">未开启</span>
                        <span v-else-if="basicinfo1.node_status == 'passive_traffic'">开启</span>
                        <span v-else>未开启</span>
                    </div>
                    
                    <!-- <div >
                        <div v-if="basicinfo1.node_status == 'proxy'">
                            <label>代理渗透：</label>
                            <span>开启</span>
                            <span>{{basicinfo1.proxy_parten}}</span>
                        </div>
                        <div v-if="basicinfo1.node_status == 'openvpn'">
                            <label>OpenVPN：</label>
                            <span>开启</span>
                        </div>
                        <div v-if="basicinfo1.node_status == 'passive_traffic'">
                            <label>被动流量扫描：</label>
                            <span>开启</span>
                        </div>
                    </div> -->
                   
                    <div>
                        <label>加密算法：</label>
                        <span>DES</span>
                    </div>
                </div>
                <div class="otherinfo">
                    <div class="part_title">
                        <label>系统监控</label>
                    </div>
                    <div>
                        <el-row :gutter="20">
                            <el-col :span="12">
                                <div class="chartdiv">
                                    <label for="">最近一小时CPU使用率</label>
                                    <div id="CPU"></div>
                                </div>
                            </el-col>
                            <el-col :span="12">
                                <div class="chartdiv">
                                    <label for="">最近一小时内存使用率</label>
                                    <div id="RAM"></div>
                                </div>
                            </el-col>
                        </el-row>
                    </div>
                </div>
            </div>
        </el-dialog>

        <el-dialog title="下载Agent" :visible.sync="agentdialogVisible">
            <div class="dialog_b_btn"> 
                <el-button size="small" @click="download" >下载</el-button>
                <el-button size="small" @click="cancalAgent">关闭</el-button> 
            </div>
            <div class="info_box">
                <div>
            <el-form ref="reportform" :model="reportform" label-width="150px" class="clearfix" :rules="rules">
                <el-form-item label="选择Agent平台：" style="margin-bottom:0"> 
                    <el-radio-group v-model="format"  > 
                        <el-radio v-for="(item,i) in formatlist" 
                        :key="i" :label="item.value" value="item.value"  
                        >{{item.label}}</el-radio>
                    </el-radio-group>
                </el-form-item>
               
            </el-form>
        </div>
            </div>
        </el-dialog>

    </div>
</template>
<style scoped lang="less">
@import './css/bas-list-page.less';

.info_box {
    padding: 24px;
}

.basicinfo {
    padding: 24px;
    background: #fff;
    border: 1px solid rgba(232, 232, 245, 1);
    margin-bottom: 24px;
}

.otherinfo {
    margin-top: 32px;
    >div{
        margin-bottom: 16px;
        label{
            margin-left: 6px;
            color: rgba(72, 72, 102, 0.87);
        }
        span{
            display: inline-block;
            color: rgba(72, 72, 102, 0.63);
            font-size: 13px;

        }
    }
    .bgcolorfff{
        background-color: #fff;
        border: 1px solid rgba(232, 232, 245, 1);
    }
    .part_title {
        margin-bottom: 16px;
        >label{
            margin-left: 0;
        }
        >span{
            margin-left: 32px;
        }
    }

    .content {
        // background: rgba(255, 255, 255, 1);
        border-radius: 2px;
        // border: 1px solid rgba(232, 232, 245, 1);
        padding: 12px 16px;
        color: rgba(72, 72, 102, 0.64);
        font-size: 13px;
    }
}
.otherinfo_1  {
    >div{
        display: inline-block;
        width: 20%;
    }
}
.minH{
    height: 50px;
}
.part_title {
    font-size: 14px;
    margin-bottom: 16px;
    font-weight: 500;
    border-left: 3px solid #4C7AE3;
    padding-left: 10px;
    height: 14px;
    line-height: 14px;
    color: rgba(72, 72, 102, 0.87);
}
.chartdiv{
    background-color: #fff;
    height: 300px;
    box-shadow: 0px 2px 4px 0px rgb(76 122 227 / 12%);
    >label{
        width: 100%;
        height: 60px;
        line-height: 60px;
        padding-left: 24px;
        font-weight: bold;
    }
    >div{
        height: calc(100% - 80px);
    }
}
/deep/ .el-dialog__body{
    background-color: rgb(247, 247, 251);
}
/deep/ .textarea textarea{
    resize: none;
}
</style>
<script>
import DelButton from "@/components/DelButton.vue";
// import {node} from '@/api/system.js';
import bas from '@/api/bas.js'
import jsFileDownload from 'js-file-download'
var echarts = require('echarts');
export default {
    name:"node",
    components:{
        DelButton
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
            format:'linux',
            formatlist:[
                {
                    label:'linux',
                    value:'linux'
                },
                {
                    label:'windows',
                    value:'windows'
                },
            ],
            reportform:{},
            agentdialogVisible:false,
             dialogFormVisible:false,
             checkIPform:{
                id: '',
                name: '',
            	ip:'',
                port: '',
            },
            showEditFileNameButton:false,
            rowId:'',
            Loading:false,
            alldelvisible:false,
            search_item: {
                page:1,
            },
            multipleSelection: [],
            tableData:[],
            currentPage:1,
            pageSize:10,
            total:0,
            updatetxt:'编辑',
            is_Update:false,
            dialogVisible:false,
            nodeform:{
                name:'',
                agent: false,
                agentCard: '',
                agentUsers: '',
                agentPort:'', 
                whitelist:'',
                targets:'',
                limit:false,
            },
            basicinfo:[],
            basicinfo1:{
                agent:false,
                agentCard:'',
                agentUsers:'',
                agentPort:'',
                limit:false,
                whitelist:'',
                node_status: '',
            },
             rules:{
                ip: [
                    { required: true, message: '节点IP不能为空', trigger: 'blur' }, 
                    { validator: validatePass2, trigger: 'blur' },
                ],
                name: [
                    { required: true, message: '节点名称不能为空', trigger: 'blur' }, 
                ],
                port: [
                    { required: true, message: '节点端口不能为空', trigger: 'blur' }, 
                ],
            },
           
            timer1: null,
            timer2: null,
      		timermillisec:0 //时间间隔
        }
        
    },
    created: function () {
        this.$store.state.activefirstMenu = "/agent"; 
        this.pageSize = this.commonjs.pageSize;
		this.timermillisec = this.commonjs.timermillisec;
        

    },
    mounted: function () {
        this.getData();
        // // this.timer1 = setInterval(() => {
		// // 	this.getData(true);
		// // }, this.timermillisec );
     
    },
    beforeDestroy() {
    //页面销毁清除定时器
        this.timer1 ? clearInterval(this.timer1) : "";
        this.timer2 ? clearInterval(this.timer2) : "";
        clearTimeout(this.commonjs.timeer);
        this.commonjs.timeer = null;
    },
    methods: {
        // 设置开启 不开启
        async changeStatus (val) {
            let res = await node.setStatus({
                status: val ? 1 : 0
            })
            if (res.code === 200) {
                let text = val ? '开启成功' : '关闭成功'
                this.$message.success(text)
            } else {
                this.$message.error(res.msg)
            }
        },
        newAdd(){
            // this.title = '新增绑定'
            this.dialogFormVisible = true;
        },
        async getData(notloading){
            let multipleSelection = [];
			if (notloading) {
				multipleSelection = this.multipleSelection;
			}
            this.Loading = true;  ///task/template/
            let params = {
                search: this.search_item.search,
                page: this.search_item.page,
                size: this.pageSize,
            }
            const res = await bas.agentlist(params)
            if (res.code === 200) {
                this.Loading = false;
                this.tableData = res.data.list;
                this.total = res.count;
                if (notloading) {//处理定时刷新多选
                    let ids = [];
                    multipleSelection.forEach(item => {
                        ids.push(item.id);
                    });
                    this.$nextTick(() => {
                        this.tableData.forEach(item => {
                        if (ids.includes(item.id)) {
                            console.log(this.$refs)
                            this.$refs.multipleTable.toggleRowSelection(item, true);
                        }
                        });
                    });
                }
            } else {
                this.$message.error(res.msg)
            }
            
        },
        saveAdd(){
            this.$refs['ruleFormaddip'].validate(async (valid) => {
                if(valid){
                    let params = {
                        name:this.checkIPform.name,
                        ip:this.checkIPform.ip,
                        port:this.checkIPform.port
                    }
                    const res = await node.addNode(params)
                    if(res.code == 200){
                        this.cancaliplogVisible()
                        this.page = 1
                        this.getData();
                        this.$message({
                            message:res.msg||'新增成功',
                            type: 'success'
                        });
                    }else{
                        this.$message({
                            message:res.msg,
                            type: 'error'
                        });
                    } 
                }
            })
        },
        cancaliplogVisible(){
            this.dialogFormVisible = false;
            this.checkIPform.name ='';
            this.checkIPform.ip = '';
            this.checkIPform.id = '';
            this.checkIPform.port = '';
            // this.$refs.ruleFormaddip.resetFields(); 
        }, 
        handlesearch(){
            this.search_item.page = 1;
            this.currentPage = 1;
            this.getData();
        },
        handleReset(){
            this.search_item.page = 1;
            this.search_item.search = '';
            this.pageSize = 10;
            this.currentPage = 1;
            this.getData();
        },
        async btnNodeAuthorize(scope,flag){ //启用禁用
            const res = await bas.basagentstatusedit({
                id: scope.row.id,
                status:flag
            });

            if(res.code === 200){
                this.$message({
                    message: flag === 1 ? '启用成功' : '禁用成功',
                    type: 'success'
                });
                this.getData();
            }else{
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async btnMultiDeleteTemplate(){ //批量删除
            if (this.multipleSelection.length == 0) return;
            var _ids = this.multipleSelection.map(item => item.id);
            const res = await node.nodedel({
                id: _ids.join(',')
            })
            if (res.code === 200) {
                this.$message({
                    message:'删除成功',
                    type: 'success'
                }); 
                this.getData();
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        async btnDel(scope) { //单个删除
            const res = await node.nodeDel({
                id: scope.row.id
            })
            if (res.code === 200) {
                this.$message({
                    message: '删除节点成功',
                    type: 'success'
                });
                scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
                this.getData(); 
            } else {
                this.$message({
                    message: res.msg,
                    type: 'error'
                });
            }
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
        handleCurrentChange(t){
            this.search_item.page = t;
            this.getData();
        },
        handleSizeChange(t){
            this.search_item.page = 1;
            this.pageSize = t;
            this.getData();
        },
        async btninfo(row){//弹出详情窗口
            this.basicinfo = [];
            const res = await node.getNodeinfo({
                id:row.id
            });
            this.basicinfo.push(res.data);  //表格内容
            // console.log(this.basicinfo);
            this.basicinfo1.id = res.data.id;
            this.basicinfo1.name = res.data.name;
            this.basicinfo1.agent = res.data.agent;
            this.basicinfo1.agentCard = res.data.agentCard;
            this.basicinfo1.agentPort = res.data.agentPort;
            this.basicinfo1.whitelist = res.data.whitelist;
            this.basicinfo1.limit = res.data.limit;
            this.basicinfo1.targets = res.data.targets;
            this.basicinfo1.node_status = res.data.node_status
            this.dialogVisible=true;
            this.is_Update = false;
            this.$nextTick(async ()=>{
                // let dtx = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
                // let dty = [820, 932, 901, 934, 1290, 1330, 1320];
                
                let cpudt = await node.getCPUChart({
                    id: row.id
                });
                if (cpudt.success) {//加判断，存在后台返回节点不存在的情况, wm
                    this.echartCPU(cpudt.data.timestamp, cpudt.data.ratio);
                }

                let ramdt = await node.getRAMChart({
                    id: row.id
                });
                if (ramdt.success) {
                    this.echartRAM(ramdt.data.timestamp, ramdt.data.ratio);
                }
                this.timer2 = setInterval(async () => {
                    // this.basicinfo = [];
                    const res = await node.getNodeinfo({
                        id:row.id
                    });
                    // 罗 22-8-10讲要保持弹窗内 表格实时刷新 编辑时停止弹窗内的刷新 wm 
                    this.basicinfo[0] = res.data //表格内容 ，若清空数组，再填进去 会导致表格跳动闪烁 wm
                    this.basicinfo1.id = res.data.id;
                    this.basicinfo1.name = res.data.name;
                    this.basicinfo1.agent = res.data.agent;
                    this.basicinfo1.agentCard = res.data.agentCard;
                    this.basicinfo1.agentPort = res.data.agentPort;
                    this.basicinfo1.whitelist = res.data.whitelist;
                    this.basicinfo1.limit = res.data.limit;
                    this.basicinfo1.targets = res.data.targets;
                    let cpudt = await node.getCPUChart({
                        id: row.id
                    });
                    if (cpudt.success) {
                        // let dty = [820, 932, 901, 934, 1290, 1330, 1320];
                        this.echartCPU(cpudt.data.timestamp, cpudt.data.ratio);
                    }

                    let ramdt = await node.getRAMChart({
                        id: row.id
                    });
                    if (ramdt.success) {
                        this.echartRAM(ramdt.data.timestamp, ramdt.data.ratio);
                    }
                }, this.timermillisec );
            })
        },
        
        echartCPU(dtx,dty){ 
            var myChart = echarts.init(document.getElementById('CPU')); 
            myChart.setOption({
                grid: {
                    left: '20',
                    top: '5%',
                    bottom: '24',
                    right: '5%',
                    containLabel: true
                },
                tooltip: {
                    // formatter: 'CPU使用率 <br />{b0}<br /> {c0}<br /> ',
                    // backgroundColor: '#000000',
                    // textStyle: {
                    //     color: '#fff'
                    // }
                    trigger: 'axis'
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    axisLine: {
                        lineStyle: {
                            color: '#d9e1e4'
                        },
                    },
                    axisLabel: {
                        color: '#4e5b5f'
                    },
                    axisTick: {
                        show: false,
                    },
                    data:dtx
                },
                yAxis: {
                    type: 'value',
                    max: 100,
                    min: 0,
                    interval:20,
                    axisLine:{
                        show: false,
                        lineStyle:{
                            color:'#d9e1e4'
                        },
                    },
                    axisLabel: {
                        color: '#4e5b5f',
                        formatter: function (val) {//2022-8-16 罗要求将y轴改为固定的0-100% wm
                            return val + '%';
                        }
                    },
                    axisTick: {
                        show: false,
                    },
                    splitLine:{
                        lineStyle:{
                            type:'dashed'
                        }
                    }
                },
                series: [
                    {
                        data: dty,
                        type: 'line',
                        smooth: true,
                        lineStyle: {
                            color: '#4c7ae3'
                        },
                        itemStyle: {
                            borderWidth: 1,
                            color: '#4c7ae3'
                        }
                    }
                ]
            });
        },
        echartRAM(dtx,dty){
            var myChart = echarts.init(document.getElementById('RAM')); 
            myChart.setOption({
                grid: {
                    left: '20',
                    top: '5%',
                    bottom: '24',
                    right: '5%',
                    containLabel: true
                },
                tooltip: {
                    // formatter: '内存使用率 <br />{b0}<br /> {c0}<br /> ',
                    // backgroundColor: '#000000',
                    // textStyle: {
                    //     color: '#fff'
                    // }
                    trigger: 'axis'
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    axisLine: {
                        lineStyle: {
                            color: '#d9e1e4'
                        },
                    },
                    axisLabel: {
                        color: '#4e5b5f'
                    },
                    axisTick: {
                        show: false,
                    },
                    data: dtx
                },
                yAxis: {
                    type: 'value',
                    axisLine:{
                        show: false,
                        lineStyle:{
                            color:'#d9e1e4'
                        },
                    },
                    axisLabel: {
                        color: '#4e5b5f',
                        formatter: function (val) {//2022-8-16 罗要求将y轴改为固定的0-100% wm
                            return val + '%';
                        }
                    },
                    axisTick: {
                        show: false,
                    },
                    splitLine:{
                        lineStyle:{
                            type:'dashed'
                        }
                    }
                },
                series: [
                    {
                        data: dty,
                        type: 'line',
                        smooth: true,
                        lineStyle: {
                            color: '#4c7ae3'
                        },
                        itemStyle: {
                            borderWidth: 1,
                            color: '#4c7ae3'
                        }
                    }
                ]
            });
        },
        btnUpdate(){
            this.timer2 ? clearInterval(this.timer2) : "";
            this.is_Update=true;
            this.nodeform.id = this.basicinfo1.id;
            this.nodeform.name = this.basicinfo1.name;
            this.nodeform.agent = this.basicinfo1.agent;
            this.nodeform.agentCard = this.basicinfo1.agentCard;
            this.nodeform.whitelist = this.basicinfo1.whitelist;
            this.nodeform.limit = this.basicinfo1.limit;
            this.nodeform.targets = this.basicinfo1.targets;
            this.nodeform.agentPort = this.basicinfo1.agentPort;
        },
        async saveUpdate(){ 
            let params = {
                id: this.nodeform.id,
                name: this.nodeform.name,
                agent: this.nodeform.agent,
                agent_Card: this.nodeform.agentCard,
                agent_port: this.nodeform.agentPort,
                whitelist: this.nodeform.whitelist,
                limit: this.nodeform.limit,
                targets: this.nodeform.targets
            }

            if (!params.name){
                this.$message({
                    message: '节点名称不能为空',
                    type: 'error'
                });
                return;
            }
            if (params.limit){
                if(!params.targets){
                    this.$message({
                        message: '限定目标开启，限定目标不能为空',
                        type: 'error'
                    });
                    return;
                }
            }
            const res = await node.saveUpdateNode(params);
            if(res.success){
                this.$message({
                    message: '编辑节点成功',
                    type: 'success'
                });
                this.dialogVisible = false;
                this.is_Update = false;

                this.getData();
            }else{
                this.$message({
                    message: res.message,
                    type: 'error'
                });
            }

        },
        cancaldialogVisible(){
            this.timer2 ? clearInterval(this.timer2) : "";
            this.is_Update = false;
            this.dialogVisible=false;
        },
        mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
			
            // let t = this.$refs['popover_id-' + row.id].showPopper;
            // if(!t){
                this.showEditFileNameButton = false;
                this.rowId = "";
            // }
          
        },
        downloadagent(){
            this.agentdialogVisible = true;
        },
        //关闭
        cancalAgent(){
            this.agentdialogVisible = false;
        },
        async download(){

            const res = await bas.downagent({
                getTempToken:true
            })
            if(res.code ==200){
              
            // const res = await bas.downagent({
            //     platform:this.format
            // }); 
            // if(res.code){  // 有code码返回，认为是失败
            //     this.$message({
            //         message:res.msg,
            //         type: 'error'
            //     });
            // }else{
            //     let filename = sessionStorage.filename
            //     const blob = new Blob([res],{
            //         type:"application/json"
            //     });
            //     const url =  window.URL.createObjectURL(blob);
            //     const a = document.createElement('a');
            //     // a.download = row.name + '.csv';
            //     a.download = filename
            //     a.href = url;
            //     a.click();

            //     // jsFileDownload(res, filename);
                

            //     sessionStorage.filename = ''
            // }

                let linkhref = '';
                if(process.env.NODE_ENV == 'development'){
                    linkhref += 'http://'+process.env.VUE_APP_API_URL; 
                }else{
                    linkhref += window.location.protocol + "//" + window.location.hostname  ;
                }
                linkhref+='/api/smart/bas/basagentdownload?temp_token='+res.data.tempToken+'&platform='+this.format;

                if(res.data.tempToken){
                    window.open(linkhref, '_blank')
                }
                else{
                    this.$message({
                        message: '下载内容有误',
                        type: 'error'
                    });
                }
            }
        },
       
    },
}
</script>