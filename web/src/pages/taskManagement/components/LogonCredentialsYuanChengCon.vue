<template>
    <div>
        <el-tabs  v-model="activeName" type="card" @tab-click="handleClick">
            <el-tab-pane label="风险详情" name="first" >
                <template>
                    <div style="position:absolute;top:0;right:0">
                        <el-dropdown trigger="click">
                            <el-button type="primary">
                                抓取信息<i class="el-icon-arrow-down el-icon--right"></i>
                            </el-button>
                            <el-dropdown-menu slot="dropdown">
                                <el-dropdown-item @click.native="captureSelect(item)" v-for="(item,index) in selectForZHuaquArr" :key="index">{{item.label}}</el-dropdown-item>
                            </el-dropdown-menu>
                      </el-dropdown>
                       <el-button size="mini" type="primary" @click.native="collectInfo" style="margin-left:10px">批量收集信息</el-button>
                      <el-button size="mini" type="primary" @click.native="duanKai">断开</el-button>
                    </div>
                        <!-- type  敏感数据  -->
                  <div  style="">
                    <div v-for="(item,index) in detailSelect.detail" :key="index" style="margin-left:30px;margin-top:20px">
                        <label for="" class="lbtxt">{{ item.title }}：
                             <span style="margin-left:20px">
                                <span>{{ item.value }}</span>
                            </span>
                        </label>
                    </div>
                  </div>
                <!-- type  敏感数据  -->
                 
                    <!--登录凭证1，数据库2，远程控制3，信息泄露4，文件泄露5-->
                    <!-- <SQLInjection :tableData="sqlTableData"></SQLInjection> -->
                    <!-- <sensitiveTable  :tableData="sqlTableData" ></sensitiveTable> -->
                    <command :cmd="cmdForZhuaQu" ref="commandEvent" :resultId="check_result_id" :SelrowID="SelrowID"></command>
                    <!-- <command ref="commandEvent" v-if="details.riskType == '远程控制'" :resultId="check_result_id"></command> -->
                </template>
            </el-tab-pane>
            <el-tab-pane label="文件目录" name="second"  >

                   <template>
                    <div style="padding-left:50px">
                        <el-button type="primary" size="small"  style="margin-bottom:10px" @click="btnReturnprev">返回上一级</el-button>
                        <!-- <el-button style="margin-left:20px;margin-bottom:10px"  type="primary"  >上传文件</el-button> -->
                    </div>
                        <el-table
                            ref="multipleTable"
                            :data="downloadTableData2"
                            tooltip-effect="dark"
                            style="width: 100%;"
                            class="context_box_bg "
                        >
                            <el-table-column
                                prop="name"
                                label="文件名">
                                <template slot-scope="scope" > 
                                    <div v-if="scope.row.type == 'directory'">
                                        <i class="el-icon-folder-opened"></i>
                                        <el-link :underline="false" @click="btnClickfolder(scope.row)">  {{ scope.row.name }}</el-link> 
                                    </div>
                                    <div v-else>
                                        <i class="el-icon-document"></i>
                                        {{ scope.row.name }}
                                    </div>
                                </template>
                            </el-table-column>
                            <el-table-column
                                prop="size"
                                label="文件大小">
                            </el-table-column>
                            <el-table-column
                                prop="lastModified"
                                label="更新时间">
                            </el-table-column>
                            <!-- <el-table-column
                                prop=""
                                label="操作">
                                <template slot-scope="scope" > 
                                    <div >
                                        <el-link :underline="false"  @click="downfile2(scope.row.loacl)">下载</el-link>
                                        <el-link 
                                            class="link_danger" 
                                            :underline="false"   >
                                                <el-popover
                                                    placement="bottom"
                                                    width="170"   
                                                    :visible-arrow="false"
                                                    :ref="`popover_id-${scope.$index}`"
                                                    popper-class="delButton_popper" >
                                                    <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                                    <div style="text-align: right; margin: 0">
                                                        <el-button size="mini" class="delCancel" @click="closeDelPopover(scope._self.$refs[`popover_id-${scope.$index}`])">取消</el-button>
                                                        <el-button size="mini" type="primary" @click="targetDelete2(scope,'yes')">确定</el-button>
                                                    </div>  
                                                    <span  slot="reference"  >删除</span>
                                                </el-popover>
                                        </el-link>
                                    </div>
                                </template>
                            </el-table-column> -->
                        </el-table>
                        <!-- <el-pagination
                            :page-size="downLoadPageSize" 
                            background
                            layout=" total, prev, pager, next, sizes, jumper"
                            :total="downLoadTotalpage"
                            :current-page="downLoadCurrentpage"
                            @current-change = "currentchange"
                            @size-change="handleSizeChange">
                        </el-pagination>  -->
                    </template>


            </el-tab-pane>
            <el-tab-pane label="已下载文件" name="dowmload" >
                <template>
                    <el-table
                        ref="multipleTable"
                        :data="detailSelect.downloadedFiles"
                        tooltip-effect="dark"
                        style="width: 100%;margin-left:30px"
                        class="context_box_bg "
                    >
                        <el-table-column
                            prop="fileName"
                            label="名称">
                        </el-table-column>
                        <el-table-column
                            prop="fileSize"
                            label="文件大小">
                        </el-table-column>
                        <!-- <el-table-column
                            prop="update_time"
                            label="更新时间">
                        </el-table-column> -->
                        <el-table-column
                            prop=""
                            label="操作">
                            <template slot-scope="scope" > 
                                <div >
                                    <el-link :underline="false"  @click="downfile(scope.row)">下载</el-link>
                                    <el-link 
                                        class="link_danger" 
                                        :underline="false"   >
                                            <!-- <el-popover
                                                placement="bottom"
                                                width="170"   
                                                :visible-arrow="false"
                                                :ref="`popover_id-${scope.$index}`"
                                                popper-class="delButton_popper" >
                                                <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                                <div style="text-align: right; margin: 0">
                                                    <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.$index}`].doClose()">取消</el-button>
                                                    <el-button size="mini" type="primary" @click="targetDelete(scope.row,'yes')">确定</el-button>
                                                </div>  
                                                <span  slot="reference"  >删除</span>
                                            </el-popover> -->
                                            <span  @click="targetDelete(scope.row,'yes')"  >删除</span>
                                    </el-link>
                                </div>
                            </template>
                        </el-table-column>
                    </el-table>
                    <el-pagination
                        :page-size="downLoadPageSize" 
                        background
                        layout=" total, prev, pager, next, sizes, jumper"
                        :total="downLoadTotalpage"
                        :current-page="downLoadCurrentpage"
                        @current-change = "currentchange"
                        @size-change="handleSizeChange">
                    </el-pagination> 
                </template>
            </el-tab-pane>
         
        </el-tabs>
        <!-- 批量收集信息弹窗 -->
            <el-dialog
            title="批量收集信息"
            :visible.sync="dialogFormVisible"
            width="1184px"
            class="fingerValidate" 
            :close-on-click-modal="false" 
            :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="saveAdd()">立即收集</el-button>
                <el-button size="small" @click="cancaliplogVisible()">关闭</el-button>
            </div>
            <div class="buginfo_box" style="padding: 20px;box-sizing:border-box;"> 
                <div class="bugbasicinfo">
                    <el-form :model="checkIPform" status-icon  ref="ruleFormaddip" label-width="80px">
                        <el-form-item label="" prop="name">
                            <label class="dialog_item_label">抓取信息</label>
                            <div>
                                <span style="display: inline-block;margin-left:10px">信息类型</span>      
                                <el-checkbox-group v-model="captureInfoIds" style="display: inline-block;margin-left:50px">
                                    <el-checkbox  v-for="item in selectForZHuaquArr" :key="item.value"  :label="item.value">{{item.label}}</el-checkbox>
                                    <!-- <el-checkbox label="网卡信息"></el-checkbox>
                                    <el-checkbox label="系统信息"></el-checkbox>
                                    <el-checkbox label="进程信息"></el-checkbox>
                                     <el-checkbox label="端口信息"></el-checkbox>
                                    <el-checkbox label="环境变量"></el-checkbox>
                                    <el-checkbox label="用户信息"></el-checkbox> -->
                                    <!-- <el-checkbox label="截屏"></el-checkbox> -->
                                </el-checkbox-group>
                            </div>
                            
                        </el-form-item>
                        <el-form-item prop="" label="">
                            <label class="dialog_item_label">文件收集</label>
                            <div style="display:flex;align-items:top">
                                <span style="display: inline-block;margin:0 20px 0 10px">文件名关键字</span>      
                                <el-input  type="textarea" rows="4" v-model="checkIPform.fileName" size="small" style="width:520px" autocomplete="off"
                                placeholder=""></el-input>
                            </div>                           
                        </el-form-item>
                        <el-form-item prop="" label="">
                             <!-- <label class="dialog_item_label">文件搜索路径</label> -->
                            <div style="display:flex;align-items:top">
                                <span style="display: inline-block;margin:0 20px 0 10px">文件搜索路径</span>      
                                <el-input  type="textarea" rows="4" v-model="checkIPform.filePath" size="small" style="width:520px" autocomplete="off"
                                placeholder=""></el-input>
                            </div>        
                        </el-form-item>
                        <div style="margin:20px 0 0 90px;color:#999;">
                            <i class="el-icon-warning"></i><span style="display: inline-block;margin-left:10px;">批量收集的信息将在已下载的文件中展示</span>
                        </div>
                </el-form>
                </div>   
            </div>
        </el-dialog>
    </div>
</template>
<style lang="less" scoped>
    .lbtxt {
      display: block;
      font-size: 14px;
      margin: 10px 0;
      border-left: 3px solid #4c7ae3;
      padding-left: 10px;
    }
/deep/ .el-table--enable-row-hover .el-table__body tr:hover>td {
    background: transparent !important;
}
/deep/ .el-tabs__header {
        margin: 0 0 24px !important;
        box-shadow:none !important;
    }
/deep/.el-tabs__header{
    padding-top:12px;
    background: transparent !important;
}
/deep/.el-tabs--card>.el-tabs__header .el-tabs__nav{
    border:none;
}
/deep/.el-tabs--card>.el-tabs__header{
    border-bottom:1px solid transparent;
}
/deep/.el-tabs--card>.el-tabs__header .el-tabs__item{
    border:1px solid #E4E7ED;
    border-radius: 2px 0px 0px 2px;
    width:88px;
    height:32px;
    line-height:32px;
    padding:0!important;
    text-align:center;
}
/deep/.el-tabs__item.is-active{
    border-color:rgba(76, 122, 227, 1)!important;
}
/deep/.el-tabs__content{
    padding:0 24px 24px 24px;
}
ul{
    margin:0;
    padding:0;
    height: 300px;
    overflow-y: auto;
}
.itemLi{
    list-style-type:none;
    margin-top:16px;
    overflow: hidden;
    // :not(:first-child){
    //     margin-top:16px;
    // }
    .circleBlue{
        display:inline-block;
        background: #4C7AE3;
        width: 8px;
        height: 8px;
        border-radius: 8px;
        margin-right:32px;
    }
    label{
        // display:block;
        // float:left;
        font-size:13px;
        font-weight: 500;
        color: rgba(72, 72, 102, 0.87);
        line-height: 14px;
        margin-right:23px;
        // display:inline-block;
        display:block;
        float:left;
        width:65px;
    }
    span{
        // display:inline-block;
        display:block;
        float:left;
        width:85%;
        // float:left;
        font-weight: 400;
        color: rgba(72, 72, 102, 0.64);
        line-height: 14px;
    }
}
#svgCanvas{
    width:100%;
    height:400px;
}
.svgCanvas /deep/ g.node rect {
    padding: 5px;
    fill: #dbe4f9;
    stroke: #E8E8F5;
    stroke-width: 1px;
}
.svgCanvas  /deep/ g.node text {
    fill: rgba(72, 72, 102, 0.87);
    // color: #4c7ae3;
    font: 14px sans-serif;
    /*font-weight:700;*/
    // font-size: 12px;
    cursor: pointer;
}

.svgCanvas  /deep/  g.type-freeze>rect {
    fill: #F7F7FB;
}
.svgCanvas  /deep/  g>g.label {
    width:100%;
    text-align:center;
}

.svgCanvas  /deep/  .node text {
    font-weight: 300;
    font-family: "Helvetica Neue", Helvetica, Arial, sans-serf;
    font-size: 12px;
    pointer-events: none;
    text-anchor: middle;
    fill: #000;
  }

.svgCanvas  /deep/  .label g text tspan:last-child {
    font-size: 10px;
    margin-top: 5px;
    color: rgba(72, 72, 102, 0.64);
    // dy: 1.5em;
  }

.svgCanvas  /deep/  .label g {
    transform: translate(0, -7px);
  }

.svgCanvas  /deep/  .node rect {
    fill: white;
    stroke-width: 0px;
    color: white;
  }

.svgCanvas  /deep/  .edgePath path {
    stroke: #c4c4ce;
    stroke-width: 1px;
    // stroke-dasharray:5,2;
}
 #myMenu{ 
    position: fixed; 
    padding: 14px 14px;
    border-radius: 4px;
    background: #fff; 
    box-shadow:0px 4px 8px 0px rgba(72,72,102,0.32);
    max-height: 300px;
    max-width: 500px;
    min-width: 300px; 
    word-wrap: break-word; 
    overflow-y: auto;
    z-index: 99999;
}
#myMenu >div{
        font-size: 14px;
        margin-bottom: 10px;
}
#myMenu >div:last-child{
    margin-bottom:0;
}
#myMenu >div >label{
    display: inline-block;
    // width: 100px;
    max-width:400px;
    font-size: 13px;   
    border-left: 3px solid #4C7AE3;
    padding-left: 10px;
    height: 14px;
    line-height: 14px;
    color:rgba(72,72,102,0.64);
    vertical-align: sub;
}
#myMenu >div >span{
    display: inline-block;
    width: calc(100% - 130px);
    color:rgba(72,72,102,0.64);
    vertical-align: top;
}
</style>
<script>
import sensitiveTable from "./SensitiveTable"//文件泄露
import SQLInjection from "./SQLInjection.vue"//数据库
import command from "./Command.vue"//远程控制命令执行
import * as d3 from 'd3' 
import dagreD3 from 'dagre-d3'
import {
    task
} from '@/api/task.js'
import $ from 'jquery' 
export default({
    name:'logonCredentials',
    components:{
        sensitiveTable,
        SQLInjection,
        command
    },
    props:{
        SelrowID:Number,
        expandData:Array,
        details:Object,
        detailSelect:Object,
        sqlTable:Array
    },
    data(){
        return {
            selectForZHuaqu:'',//抓取信息下拉
            selectForZHuaquArr:'',//抓取信息下拉
            activeName: 'first',
            check_result_id:'',
            state:[
                // { id: 0, label: 'wwww.baidu.com', class: 'type-freeze',cont:'hdskfhdskhkf' },
                // { id: 1, label: 'V1\n数据同步同同同', class: 'type-freeze' ,},
                // { id: 2, label: 'V2\nhive-sql', class: 'type-freeze' },
                // { id: 3, label: 'V3\nspark-sql', class: 'type-freeze' },
                // { id: 4, label: 'V4\nshell', class: 'type-freeze' },
                // { id: 5, label: 'V5\npython', class: 'type-freeze' },
                // { id: 6, label: 'V6\n虚节点同同', class: 'type-freeze' },
                // { id: 7, label: 'V7\nspark-sql', class: 'type-freeze' },
                // { id: 8, label: 'V8\nshell', class: 'type-freeze' },
                // { id: 9, label: 'V9\n数据同步同', class: 'type-freeze' },
                // { id: 10, label: 'V10\nshell', class: 'type-freeze' },
                // { id: 11, label: 'V11\nspark-sql', class: 'type-freeze' },
                // { id: 12, label: 'V12\nspark-sql', class: 'type-freeze' },
                // { id: 13, label: 'V13\n虚节点同同同同同', class: 'type-freeze' },
                // { id: 14, label: 'V14\n数据同步同', class: 'type-freeze' },
            ],
            edg:[
                // { start: 0, end: 1, option: {} },
                // { start: 1, end: 4, option: {} },
                // { start: 1, end: 3, option: {} },
                // { start: 1, end: 2, option: {} },
                // { start: 6, end: 7, option: {} },
                // { start: 5, end: 6, option: {} },
                // { start: 9, end: 10, option: { color: 'red'} },
                // { start: 8, end: 9, option: {} },
                // { start: 11, end: 12, option: {} },
                // { start: 8, end: 11, option: {} },
                // { start: 5, end: 8, option: {} },
                // { start: 1, end: 5, option: {} },
                // { start: 13, end: 14, option: {} },
                // { start: 1, end: 13, option: {} },
            ],
            diagGraph:{},
            downLoadPageSize:10,
            downLoadTotalpage:0,
            downLoadCurrentpage:1,
            downloadTableData:[ ],
            downloadTableData2:[ ],
            menuflag:false,
            menustyle:'',
            target_id:'',
            sqlTableData:[],
            cmdForZhuaQu:'',
            dialogFormVisible:false, // 批量收集信息
            checkIPform:{
                filePath: '',
            	fileName:'',
            },
            captureInfoIds:[],
            typeValue: '', // 远控类型对应的数值
            farControlTypeList: [
                {
                    label: '远控',
                    value: 1
                },
                {
                    label: '反弹shell',
                    value: 2
                }
            ],
            path:'/'
        }
    },
    mounted(){
        // this.getGraphData();
        // this.initdagreD3();
        this.getCaptureinfoenum();
        console.log(this.SelrowID)
    },
    created(){
        this.sqlTableData = this.sqlTable;
        this.target_id = this.details.target_id;
        this.check_result_id = this.details.checkResultID;
        console.log(3,  this.detailSelect.detail)
        let typeObj = this.detailSelect.detail.find(item => item.title === '远控类型')
        let typeName = ''
        if (typeObj) {
            typeName = typeObj.value
        }
        let obj = this.farControlTypeList.find(item => item.label === typeName)
        if (obj) {
            this.typeValue = obj.value
        }
        console.log(101, this.detailSelect)
    },
    methods: {
        closePOP(aaa){
            // this.$refs.aaa.doClose();
        },
        async saveAdd(){
            // this.$refs['ruleFormaddip'].validate(async (valid) => {
            //     if(valid){
                    let params = {
                        remoteSessionId:this.SelrowID,
                        // captureType: '',
                        captureType: this.typeValue,//int,系统，1-远控,2-反弹shell，默认2
                        captureInfoIds: this.captureInfoIds.join(','), //string,信息类型,多个用英文逗号隔开
                        fileName: this.checkIPform.fileName,//string,文件名关键字,captureInfoIds与fileName字段至少有一个不为空
                        filePath: this.checkIPform.filePath
                    }
                    const res = await task.exceshellmany(params)
                    if(res.code == 200){
                        this.dialogFormVisible = false
                        this.$refs.commandEvent.addResult({
                            key: '',
                            result: res.data.result
                        })
                        this.$message({
                            message:res.msg||'收集成功',
                            type: 'success'
                        });
                    }else{
                        this.$message({
                            message:res.msg,
                            type: 'error'
                        });
                    } 
            //     }
            // })
        }, 
         cancaliplogVisible(){
            this.dialogFormVisible = false;
            this.checkIPform.name ='';
            this.checkIPform.ip = '';
            this.checkIPform.id = '';
            this.checkIPform.port = '';
            // this.$refs.ruleFormaddip.resetFields(); 
        }, 
        // 删除
        targetDelete2(scope,flag){
            if(flag == 'yes'){
                this.$ajax({
                    method:'post',
                    url: '/smart/task/deletefile',
                    data: {
                        target_id:this.target_id,
                        file_name:scope.row.file_name,
                        file_path:scope.row.file_path,
                        file_size:scope.row.file_size,
                        file_type:scope.row.file_type,
                        file_time:scope.row.file_time,
                    } 
                })
                .then(res =>{
                    let dt = res.data; 
                    if(dt.code == 200){   
                        this.$message({
                            message:dt.msg,
                            type: 'success'
                        });
                        this.getDownLoadData();
                    }else{
                        this.$message({
                            message:dt.msg,
                            type: 'error'
                        });
                    }
                }).catch(err=>{})
            }
        },
        targetDelete(scope,flag){
            console.log(scope,'scopescopescopescopescope',flag);
                 
            if(flag == 'yes'){
                this.$ajax({
                    method:'get',
                    url: '/smart/task/delfile',
                    params: {
                        remoteSessionID:this.SelrowID,
                        fileName:scope.fileName,
                    } 
                })
                .then(res =>{
                    let dt = res.data; 
                    if(dt.code == 200){   
                        this.$message({
                            message:dt.msg,
                            type: 'success'
                        });
                       this.getDownLoadData();
                    }else{
                        this.$message({
                            message:dt.msg,
                            type: 'error'
                        });
                    }
                }).catch(err=>{})
            }
        },
        // 下载 -已下载文件
        downfile(scope){
            window.open('/api/smart/task/downloadfile?remoteSessionID='+this.SelrowID+'&fileName='+scope.fileName);
            // var urldd = window.location.host
            // console.log(urldd,'00');
            //  this.$ajax({
            //         method:'get',
            //         url: `/smart/task/downloadfile`,
            //         params: {
            //             remoteSessionID:this.SelrowID,
            //             fileName:scope.fileName,
            //         } 
            //     })
                // var link = document.createElement('a');
                // link.href = `${urldd}/smart/task/downloadfile?remoteSessionID=${this.SelrowID}&fileName=${scope.fileName}`;
                // link.download = 'myfile.zip';
                // link.click();
                
                // .then(res =>{
                //     let dt = res.data; 
                //     if(dt.code == 200){   
                //         this.$message({
                //             message:dt.msg,
                //             type: 'success'
                //         });
                //         this.getDownLoadData();
                //     }else{
                //         this.$message({
                //             message:dt.msg,
                //             type: 'error'
                //         });
                //     }
                // }).catch(err=>{})
        },
         // 下载 -文件目录
        downfile2(url){

            window.open(url);
        },
        handleClick(tab, event, noNeedRequest) { 
            if(tab.name == 'second'){
                // var statePoint = 0; // 当前选中的点   
                // this.getGraphData(statePoint); 
                this.getfilelist(); 

            }if(tab.name == 'dowmload'){ 
               
            }if(tab.name == 'first'){
                if(this.details.riskType == '远程控制'){
                    this.$refs.commandEvent.getFocus();
                }
            }
        },
        // 文件目录
        async getfilelist(){ 
            const res = await task.filemanagement({
                id:this.SelrowID,
                path:this.path,
            });
            if(res.code == 200){
                this.downloadTableData2 = res.data.list;
            }
        },
        //文件夹点击
        async btnClickfolder(row){
            this.path +=row.name+'/';
            const res = await task.filemanagement({
                id:this.SelrowID,
                path: this.path
            });
            if(res.code == 200){
                this.downloadTableData2 = res.data.list;
            }
        },
        //返回上一级
        async btnReturnprev(){ 
            // let lastSlashIndex = this.path.lastIndexOf('/',this.path.length - 2); // 从倒数第二个字符开始查找
            // let result = this.path.substring(0, lastSlashIndex);
            let result = this.path.replace(/(\/)[^\/]*\/?$/, '$1');
            this.path = result;
            console.log( this.path)
            // return;
            const res = await task.filemanagement({
                id:this.SelrowID,
                path:  this.path
            });
            if(res.code == 200){
                this.downloadTableData2 = res.data.list;
            }
        },
        captureSelect(item){
            console.log("你真的奇奇怪怪",item);
            this.captureinfo(Number(item.value))
        },
        // 批量收集信息
        collectInfo () {
            this.dialogFormVisible = true;
        },
        // 断开
        async duanKai(){
            const dt = await task.break({
            });
            if (dt.code == 200) {
                this.cmdForZhuaQu = dt.data.cmd;
            }
        },
        //抓取信息
        async captureinfo(id){
            let _this = this;
            const dt = await task.captureinfo({
                id: id,
                captureType: this.typeValue
            });
            if (dt.code == 200) {
                this.cmdForZhuaQu = dt.data.cmd;
                setTimeout(()=>{
                    _this.$refs.commandEvent.onSubmitInput();
                },500)
                
                // this.selectForZHuaqu = dt.data;
                // this.state = dt.result.state;
                // this.edg = dt.result.edge;

                // this.diagGraph.init(statePoint, this.state, this.edg); //创建关系图 
            }
          
        },
        //抓取信息枚举--下拉
        async getCaptureinfoenum(){

            const dt = await task.captureinfoenum();
            if (dt.code == 200) {
                this.selectForZHuaqu = dt.data;
                this.selectForZHuaquArr = []
                 for (const n in this.selectForZHuaqu) {
                    this.selectForZHuaquArr.push({
                        value: n,
                        label: this.selectForZHuaqu[n]
                    })
                 }
                // this.state = dt.result.state;
                // this.edg = dt.result.edge;

                // this.diagGraph.init(statePoint, this.state, this.edg); //创建关系图 
            }
          
        },
        async getGraphData(statePoint){
            this.check_result_id = this.details.checkResultID;

            const dt = await task.riskAttackPath({
                check_result_id: this.details.checkResultID,
                risk_type: this.details.riskTypeNum,
                vul_name: this.details.vulName,
                target_url: this.details.targetUrl,
            });
            if (dt.code == 200) {
                this.state = dt.result.state;
                this.edg = dt.result.edge;

                this.diagGraph.init(statePoint, this.state, this.edg); //创建关系图 
            }
          
        },
        initdagreD3(){
            let _that = this;
            this.diagGraph = { //diag图数据操作 
                state:[],
                edg:[],
                statePoint: '',
                g: '',
                init: function (statePoint ,state, edg) {
                    this.statePoint = statePoint
                    this.state = state
                    this.edg = edg
                    this.createG();
                    this.renderG(); 
                },
                drawNode: function () { 
                    for (let i in this.state) { //画点
                        let el = this.state[i] 
                        let style = 'padding:20px';
                        this.g.setNode(el.id, {
                            id: el.id,
                            label: el.label,
                            // labelType:"html",
                            // label: el.label+'<i class=\"iconfont iconshouye\"></i> ',
                            class: el.class,
                            // style: style,
                            node_content: el.cont,
                        }); 
                    }
                    this.g.nodes().forEach((v) => { //画圆角
                        var node = this.g.node(v); 
                        node.rx = node.ry = 5;
                          
                    });
                },
                addNode:function(nodelist){},
                drawEdg: function () {
                    for (let i in this.edg) { // 画连线
                        let el = this.edg[i]
                        if (el.start === this.statePoint || el.end === this.statePoint) {
                            this.g.setEdge(el.start, el.end, {
                              style: "stroke: #4C7AE3; fill: none;",
                              arrowheadStyle: "fill: #c4c4ce;stroke: #c4c4ce;",
                              arrowhead: 'undirected',
                              curve: d3.curveBasis,//控制线类型，弧线
                            });
                        } else {
                            this.g.setEdge(el.start, el.end, {
                                arrowhead: 'undirected',
                                curve: d3.curveBasis,//控制线类型，弧线
                            });
                        }
                    }
                },
                createG: function () {
                    this.g = new dagreD3.graphlib.Graph()
                      .setGraph({
                        rankdir: 'TB', //设置方向
                        ranksep:32,
                      })
                      .setDefaultEdgeLabel(function () { return {}; });
                },
                renderG: function () {  
                    var render = new dagreD3.render();
                    var svg = d3.select("#svgCanvas"); //声明节点
                    
                    svg.select("g").remove(); //删除以前的节点，清空画面
                    var svgGroup = svg.append("g");
                    var inner = svg.select("g");
                    var zoom = d3.zoom().on("zoom", function () { //添加鼠标滚轮放大缩小事件
                        inner.attr("transform", d3.event.transform);
                        _that.menuflag=false;
                    });
                    svg.call(zoom);
                    this.drawNode();//画点
                    this.drawEdg();// 画连线
                    render(d3.select("svg g"), this.g); //渲染节点 
                    
                    var max = svg._groups[0][0].clientWidth > svg._groups[0][0].clientHeight ? svg._groups[0][0].clientWidth : svg._groups[0][0].clientHeight;
                   
                    var initialScale = 0.8;
                    var tWidth = (svg._groups[0][0].clientWidth - this.g.graph().width * initialScale) / 2;
                    var tHeight = (svg._groups[0][0].clientHeight - this.g.graph().height * initialScale) / 2;
                    var trans = d3.zoomIdentity.translate(tWidth, 0).scale(1.5);  
                    svg.call(zoom.transform,trans); //元素居中

                   
                }, 
            }
        },
   
        handleScroll(e) { //滚动，清空右键东西
            // this.menuflag=false; 
            if(e.target.id && e.target.id == 'myMenu'){
               
            }else{
                 this.menuflag = false;
            }
        },
        // 获取下载文件列表
        getDownLoadData(scope){
                    //  刷新页面
              this.$router.go(0);
                
                    
        //    scope._self.$refs[`popover_id-${scope.$index}`].doClose()

        },
        currentchange(t){
            this.formData.page_num = t; 
            this.getDownLoadData();
            this.downLoadCurrentpage = t;
        },
        handleSizeChange(t){
            this.formData.page_num = 1;
            this.downLoadPageSize = t;
            this.getDownLoadData();
        }
    }
})
</script>
