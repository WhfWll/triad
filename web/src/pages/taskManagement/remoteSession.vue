/* 
  远程会话
 */
<template>
    <div class="riskdata bgColor2">
    
        <div class="risktable">
            <div class="search-box">
                <div v-if="false" class="operationbutton">
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170" trigger="click"
                        :visible-arrow="false" v-model="riskVisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="riskVisible = false">取消
                            </el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDeleteRisk">确定
                            </el-button>
                        </div>
                        <el-button type="warning" size="small" slot="reference"
                            :disabled="!multipleSelectionRisk.length">删除</el-button>
                    </el-popover>
                </div>
                <div class="serach-condition">
                    <!-- <div>
                        <el-select v-model="riskFormData.risk_type" placeholder="风险类型" size="small" class="selectwidth">
                            <el-option v-for="(item, index) in risklist" :key="index" :label="item.name_zh"
                                :value="item.name">
                            </el-option>
                        </el-select>
                    </div> -->
                    <div class="search-text">
                        <el-input style="width:400px" placeholder="搜索测试目标、文件名、信息" @keydown.enter.native="handlesearchRisk" v-model="riskFormData.search" class="input-with-select"
                            size="small" clearable> </el-input>
                        <el-button type="primary" size="small" @click="handlesearchRisk">搜索
                        </el-button>
                    </div>
                    <div>
                        <el-button type="primary" size="small" @click="handleResetRisk">重置
                        </el-button>
                    </div>
                </div>
            </div>
        
            <el-table ref="riskTable" :data="riskTableData" tooltip-effect="dark" style="width: 100%"  height="calc(100% - 102px)"
                @selection-change="handleSelectionChangeRisk" @expand-change="handleExpandChange"
                :expand-row-keys="expands" :row-key="getRowKeys">
                <el-table-column type="expand">
                    <template slot-scope="props">
                        
                        <LogonCredentialsYuanChengCon  :SelrowID="SelrowID" :detailSelect="detailSelect" :expandData="expandData" :details="props.row" :sqlTable="sqlTable">
                        </LogonCredentialsYuanChengCon>

                    </template>
                </el-table-column>
                
                <el-table-column label="会话编号" prop="sessionNum">
                </el-table-column>
                <el-table-column label="测试目标" prop="targetUrl">
                </el-table-column>
                <el-table-column label="路由" prop="route">
                </el-table-column>
                <el-table-column label="远程控制" prop="remoteControl">
                </el-table-column>
                <el-table-column label="状态" prop="status">
                </el-table-column>
            </el-table>
            <el-pagination :page-size="pageSize" background layout=" total, prev, pager, next, sizes, jumper"
                :total="riskTotalpage" :current-page="currentpageRisk" @current-change="currentchangeRisk"
                @size-change="handleSizeChangeRisk">
            </el-pagination>
           
        </div>
    </div>
</template>
<style lang="less" scoped>



    /deep/ thead {
        .cursorPointer {
            cursor: pointer;
            position: absolute;
            top: 6px;

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
 

    .risktable {
        // margin-top: 16px;
        height: 100%;
        padding: 24px;
        background: #FFF;
        box-sizing: border-box;
        border-radius: 4px;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    }
</style>
<script>
var echarts = require('echarts');
import LogonCredentialsYuanChengCon from './components/LogonCredentialsYuanChengCon.vue';
import { 
    task
} from '@/api/task.js'
export default {
    name:'riskdata',
    components:{
        LogonCredentialsYuanChengCon
    },
    props: { 
        task_id: {},
    },
    data(){
        return{
            SelrowID: '',
            detailSelect: {},
            pageSize:10,
            riskFormData: {
                risk_type: 0,
                search: '', 
            },
            riskVisible: false,
            multipleSelectionRisk: [],
            risklist: [
                { "name": 0, "name_zh": "全部" },
                { "name": 1, "name_zh": "远程控制" }, 
                { "name": 2, "name_zh": "登录凭证" }, 
                { "name": 3, "name_zh": "敏感数据" }, 
                { "name": 4, "name_zh": "敏感文件" },
                { "name": 5, "name_zh": "数据库" } 
                ],
            riskTotalpage: 0,
            currentpageRisk: 1,
            expands: [],
            riskTableData: [],
            expandData: [],
            sqlTable: [],
            risk_types: '',
            getRowKeys(row) {
                return row.id
            },
        }
    },
    created: function () { 
        this.pageSize = this.commonjs.pageSize;  
        this.getTableData();
 
     
    },
    methods: {
        clickButton(type) {
            switch (type) {
                case '会话编号':
                    this.$refs.list.toggleMenu();
                    break;
            }
        },    

        async getSelectData() {
            // const dt = await task.riskType();
            // //TODO: 临时处理
            // if (dt.code==200) {
            //      let dataV = Object.keys(dt.data).map(item => {
            //         return {
            //             name: Number(item),
            //             name_zh: dt.data[item]
            //         }
            //     })
            //     this.risklist = dataV
            //     console.log(this.risklist,'this.risklist');
            // }
            // else {
            //     this.$message({
            //         message: dt.error,
            //         type: 'error'
            //     });
            // }
           
        }, 
        async getTableData(target_ids) {
            let multipleSelection = []
            // if (notloading) {
            //     multipleSelection = this.multipleSelectionRisk
            // }
            // console.log(this.riskFormData.risk_type);
            let _target_ids = '';
            if (target_ids === undefined || target_ids === '') {
                _target_ids = '';
            } else {
                _target_ids = target_ids;
            }
            const dt = await task.riskList2({
                search:this.riskFormData.search,
                taskId:  this.task_id,
                // taskId:  1,
                page: this.currentpageRisk,
                size: this.pageSize,
            })

            if (dt.code == 200) {
                this.riskTotalpage = dt.data.count;
                this.riskTableData = dt.data.result;
                // if (notloading) {
                //     let ids = []
                //     multipleSelection.forEach(item => {
                //         ids.push(item.check_result_id)
                //     })
                //     this.$nextTick(() => {
                //         this.riskTableData.forEach(item => {
                //             if (ids.includes(item.check_result_id)) {
                //                 this.$refs.riskTable.toggleRowSelection(item, true)
                //             }
                //         })
                //     })
                // }
            } else {
                this.$message({
                    message: dt.error,
                    type: 'error'
                });
            }
        },
        handleSelectionChangeRisk(val) {
            this.multipleSelectionRisk = val;
        },
        // 查看风险类型展开
        handleExpandChange(row, expandedRows) {
            this.SelrowID = row.id
            // console.log(row,'row======================',expandedRows);
            var that = this;
            if (expandedRows.length > 0) {
                that.expands = []
                if (row) {
                    that.check_result_id = row.id;
                    that.risk_types = row.risk_type_num
                    // if(row.risk_type_num == '1'){
                    that.fnexpandajax();
                    // }
                }
            }
            else {
                that.expands = []; // 默认不展开
            }
        },
        // 加载展示内容
        async fnexpandajax() {
             // check_result_id: this.check_result_id,
                // risk_type: this.risk_types
            const dt = await task.riskDetail2({
                id:this.check_result_id,
               
            })
            if (dt.code == 200) {
                this.expandData = dt.data.detail;
                this.detailSelect = dt.data
                // this.sqlTable = dt.table; //wdh 11.10
                this.sqlTable = dt.data.downloadedFiles;
                this.expands.push(this.check_result_id);


            } else {
                this.$message({
                    message: dt.msg,
                    type: 'error'
                });
            }
        },
        // 搜索
        handlesearchRisk() {
            this.getTableData();
        },
        // 重置
        handleResetRisk() {

            this.currentpageRisk = 1;
            this.riskTotalpage =0;
            this.riskFormData.risk_type = '',
            this.riskFormData.search = '',
            this.getTableData();
        },
        currentchangeRisk(t) {
            this.currentpageRisk = t;
            this.getTableData();
        },
        handleSizeChangeRisk(t) {
            this.currentpageRisk = 1;
            this.pageSize = t;
            this.getTableData();
        },
      async  btnMultiDeleteRisk() { 
                  
            if (this.multipleSelectionRisk.length == 0) return; 
              let _ids = this.multipleSelectionRisk.map(item => item.id); 
              this.riskVisible = false;
              const res = await	task.delvulevidence2({
                ids: _ids.join(",")
              }); 
              if (res.code == 200) {
                this.$message({
                  message: "删除任务成功",
                  type: "success"
                }); 
                this.currentpageRisk = 1;
               
                this.getTableData();
              } else {
                this.$message({
                  message: res.msg,
                  type: "error"
                });
              }
        }, 
    },
}
</script>