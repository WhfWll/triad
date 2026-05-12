<template>
    <div>
        <div  class="tasktarget_box">
            <div class="search-box"> 
                <div class="operationbutton"> 
                    <el-button type="primary" size="small" @click="btnNewTask" :disabled='!isShowShouQuan'>新建</el-button>
                    <el-popover popper-class="delButton_popper" placement="bottom-start" width="170"
                        trigger="click" :visible-arrow="false" v-model="alldelvisible">
                        <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                        <div style="text-align: right; margin: 0" class="">
                            <el-button size="mini" class="delCancel" @click="alldelvisible = false">取消
                            </el-button>
                            <el-button size="mini" type="primary" @click="btnMultiDelete">确定
                            </el-button>
                        </div>
                        <el-button  type="warning" size="small" slot="reference" :disabled="!multipleSelection.length">删除
                        </el-button>
                    </el-popover> 
                    <el-button :disabled="!multipleSelection.length" style="margin-left:10px" type="primary" size="small" @click="btnReport(multipleSelection)">生成报告</el-button>
                    
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
        <el-table ref="targetTable" :data="tableData" tooltip-effect="dark" style="width: 100%"
            @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave"  
            @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55">
            </el-table-column>
              <el-table-column prop="taskName" label="任务名称"> 
                <template slot-scope="scope">
                    <el-link @click="btnTasktargetinfo(scope.row)"> 
                        {{scope.row.taskName}}
                    </el-link>
                </template>
            </el-table-column> 
            <el-table-column prop="executeTypeName" label="任务类型">
          <template slot-scope="scope">
            {{ scope.row.executeTypeName}}
          </template>
        </el-table-column>
        <el-table-column   label="任务风险" show-overflow-tooltip>
          <template slot="header"> 
              <span class="cursorPointer" @click="clickButton('任务风险')"
                :class="(formData.risk_level !== '' && formData.risk_level !== 0) ? 'active' : ''">任务风险<i
                  class="iconfont iconshaixuan"></i>
              </span> 
            <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.risk_level" clearable
              size="small" ref="loglistRef" @change="handlesearch">
              <el-option v-for="(item, index) in risklevellist" :key="index" :label="item.label" :value="item.value">
              </el-option>
            </el-select>
          </template>
          <template slot-scope="scope">
            <span :class="[
              { 'riskstyle risk_hight': scope.row.riskLevel == '1' },
              { 'riskstyle risk_middle': scope.row.riskLevel == '2' },
              { 'riskstyle risk_low': scope.row.riskLevel == '3' },
              { 'riskstyle risk_nofind': scope.row.riskLevel == '4' }
            ]"><i></i>{{ scope.row.riskLevelName }} </span> 
          </template>
        </el-table-column>
        <el-table-column prop="" label="目标风险">
          <template slot="header">
            目标风险
            <el-tooltip class="item" effect="dark" placement="right">
              <div slot="content">
                从左往右依次为高危目标、中危目标、低危目标、安全目标
              </div>
              <i class="iconfont icontishi" style="position: absolute;top:0;left:66px"></i>
            </el-tooltip>
          </template>
          <template slot-scope="scope">
            <span class="tag_status tag_danger bug_status">{{scope.row.targetRisk[0]}}</span>
            <span class="tag_status tag_warning bug_status">{{scope.row.targetRisk[1]}}</span>
            <span class="tag_status tag_primary bug_status">{{scope.row.targetRisk[2]}}</span>
            <span class="tag_status tag_success bug_status">{{scope.row.targetRisk[3]}}</span>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="更新时间"> </el-table-column>
        <el-table-column prop="statusName" label="状态">
            
        </el-table-column>

        </el-table>
    </div>
</template>
<style scoped>

</style>
<script>
export default {
    name:'',
    data(){
        return{
            alldelvisible:false,
            multipleSelection:[],
            formData:{
                search:'',
                page:1
            },
            tableData:[],
        }
    },
    created(){

    },
    methods:{
        handlesearch(){},
        handleReset(){},
        btnMultiDelete(){

        },
        handleSelectionChange(){

        },
        mouseenter(){},
        mouseleave(){},
    }
}
</script>
