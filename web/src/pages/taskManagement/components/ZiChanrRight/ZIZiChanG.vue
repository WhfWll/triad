<template>
    <div class="loglist">
        <div  class="log_box">
            <div class="search-box">
                <div class="operationbutton"> 
                    <!-- <el-popover popper-class="delButton_popper" placement="bottom-start" width="170"
                        trigger="click" :visible-arrow="false" v-model="alldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消
                            </el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDelete">确定
                            </el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除
                        </el-button>
                    </el-popover>  -->
                </div>
                <div class="serach-condition">
                    <div class="search-text">
                        <el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearch"  v-model="formData.search" class="input-with-select" size="small"
                            clearable> </el-input>
                        <el-button type="primary" size="small" @click="handlesearch">搜索</el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleReset">重置</el-button>
                    </div>
                </div>
            </div>
        </div> 
        <el-table border ref="logTable" :data="logtableData" tooltip-effect="dark" style="width: 100%"
            @selection-change="handleSelectionChange" 
            @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" >
           
            <el-table-column prop="subName" label="子资产组">
            </el-table-column> 
            
             <el-table-column prop="levelName" label="级别">
                </el-table-column>
            <el-table-column prop="levelRiskName" label="资产组风险">
              <template slot-scope="scope">
                <span :style="getRiskLevelColor(scope.row.levelRiskName)">
                  {{ scope.row.levelRiskName }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="assetStatics" label="资产统计">
              <template slot-scope="scope">
                    <div style="display:flex;justify-content: space-between;">
                        <div>
                          <span style="color:red;margin-right:3px">{{scope.row.assetStatics.highRiskAsset}}</span>
                          <span style="color:#f9ca84;margin-right:3px">{{scope.row.assetStatics.middleRiskAsset}}</span>
                          <span style="color:#1223ff;margin-right:3px">{{scope.row.assetStatics.lowRiskAsset}}</span>
                          <span style="color:#8ac23a">{{scope.row.assetStatics.safeAsset}}</span>
                        </div>
                          <span style="color:#00">{{scope.row.assetStatics.totalAsset}}</span>
                        
                    </div>
                </template>
            </el-table-column>
     
            <el-table-column prop="vulStatics" label="漏洞统计">
                <template slot-scope="scope">
                    <div v-if="showOperateButton && rowId == scope.row.targetId">
                        <el-link class="link_primary" :underline="false" @click="btnShow(scope.row)">详情</el-link>
                    </div>
                    <div v-else>
                      <div style="display:flex;justify-content: space-between;">
                          <div>
                            <span style="color:red;margin-right:3px">{{scope.row.vulStatics.deadlyVul}}</span>
                            <span style="color:#f9ca84;margin-right:3px">{{scope.row.vulStatics.highRiskVul}}</span>
                            <span style="color:#1223ff;margin-right:3px">{{scope.row.vulStatics.mediumRiskVul}}</span>
                            <span style="color:#8ac23a">{{scope.row.vulStatics.lowRiskVul}}</span>
                          </div>
                            <span style="color:#00">{{scope.row.vulStatics.totalRiskVul}}</span>
                          
                      </div>
                    </div>
                </template>
            </el-table-column>
        </el-table>
        <el-pagination background @size-change="handleSizeChangelog" @current-change="handleCurrentChangelog"
            :current-page="currentPage" :page-size="pageSize" layout="total,  prev, pager, next, sizes, jumper"
            :total="total">
        </el-pagination>

        <el-dialog :title="dialogtitle" :visible.sync="dialogVisible" width="1184px" :close-on-click-modal="false"
            :validate-on-rule-change="false" :show-close="false" >
            <div class="dialog_b_btn"> 
                <el-button size="small" @click="dialogVisible = false">关闭</el-button>
            </div>
            <div style="padding:24px">
                <div>
                    <ul class="loginfolist">
                        <li v-for="(item,index) in loginfolist" :key=index>
                            <span>[*]</span>
                            <span>[{{ item.createTime }}]</span>
                            <span>{{ item.pocname }}：</span>
                            <span>{{ item.result }}</span>
                        </li>
                    </ul>
                </div>
            </div>
        </el-dialog>



    </div>
</template>
<style lang="less" scoped> 
    .loglist{
        background: #fff;
        padding: 24px 24px;
        box-sizing: border-box; 
        box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12);
        border-radius: 4px;
    }
    .loginfolist {
        list-style: none;
        margin-top: 24px;
        li {
            color: rgba(72, 72, 102, 0.6400);
            font-size: 13px;
        }
    }
</style>
<script>
import { traffic } from '@/api/assetManagement'; 
export default{
    name:'loglist3',
    props:{ 
        task_id:{},
    },
    data(){
        return{
            LoginfoID:null,
            multipleSelection:[],
            alldelvisible:false,
            formData:{
                search:'',
                page:1,
            },
            showOperateButton:false,
            rowId:'',
            currentPage:1,
            total:0,
            pageSize:10,
            logtableData:[],
            dialogtitle:'',
            dialogVisible:false,
            loginfolist:[],
        }
    },
    created(){
    //   this.getData()
    },
    mounted(){

    },
    methods:{
      getRiskLevelColor(levelRiskName) {
        const riskColors = {
          '高危': 'red',
          '中危': 'orange',
          '低危': 'yellow',
          '安全': 'green'
        };
        return { color: riskColors[levelRiskName] || 'defaultColor' }; // 如果没有匹配的风险等级，返回默认颜色
      },
        async getData(){
            const res  = await traffic.trafficVulnInfo({
                "groupID": this.$store.state.groupID[0]||-1, 
                "search":this.formData.search, 
                "page": this.formData.page, 
                "size": this.pageSize
            })
            if(res.code== 200){
              
                this.logtableData = res.data.subAssetGroupInfo;
                console.log(this.logtableData,'console.log(this.logtableData);');
                this.total = res.data.count;
            }else{
                this.$message({
                    message: res.msg,
                    type: "error"
                });
            }

        },
        handlesearch(){
            this.formData.page = 1;
            this.currentPage = 1;
            this.getData();
        },
        handleReset(){
            this.formData.page = 1;
            this.formData.search = ''; 
            this.currentPage = 1;
            this.getData();
        },
        handleSizeChangelog(t){
            this.formData.page = 1;
            this.currentPage = 1;
            this.pageSize = t;
            this.getData();
        },
        handleCurrentChangelog(t){ 
            this.formData.page = t;
            this.currentPage = t;
            this.getData();
        },
        btnMultiDelete(){

        },
        btnShow(row){
           console.log(row,'row111111111');
          // 跳转到概览
          this.$emit('gotoFirstTab',[row.groupID,1])
        },
        async getLoginfo(_id){
            const res = await task.loginfo({
                taskLogId:_id
            })
            if(res.code == 200){
                this.loginfolist = res.data.list
            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }
        },
        mouseenter(row, colum, cell, event) { 
            this.showOperateButton = true; 
            this.rowId = row.targetId ;  //赋值行id，便于页面判断 
        },
        mouseleave(row, colum, cell, event) { 
            this.showOperateButton = false; 
            this.rowId='';
        },
        handleSelectionChange(val){
            this.multipleSelection = val;
        },
    },
    watch: {
        dialogVisible(newValue, oldValue) {
            if(newValue == false){
                clearInterval(this.LoginfoID);
            }
        }
    },
}
</script>