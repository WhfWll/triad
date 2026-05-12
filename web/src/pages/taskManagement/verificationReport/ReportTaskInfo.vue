<template>
    <div class="outContainer"> 
        <div class="portMenu">

        </div>
      <div class="main-title  ">   
          <router-link :underline="false" class="classA" :to="{ path: '/verificationReport' }" >验证报告</router-link>  
          <label class="currentpagetitle"> 
               <el-tooltip class="item" effect="dark" :content="task_name" placement="bottom">
                  <span> {{task_name}}</span>
              </el-tooltip>
          </label>
        </div> 
      <div class="taskinfolist context_box_bg">
          <el-tabs v-model="activeTableName" @tab-click="handleClick">
              <el-tab-pane label="任务参数" name="tabs1" >
                  <div class="part1 mt24"> 
                      <el-row :gutter="24">
                          <el-col :span="18">
                              <div class="info_box info_box_height150 info_box_bg1">
                                <el-row style="z-index:9">
                                      <el-col :span="8">
                                          <div class="info_block">
                                              <i class="micon iconfont iconrenwumingcheng"></i>
                                              <label for="" class="lbname">任务名称</label>
                                              <span class="spvalue">
                                                  <el-tooltip class="item" effect="dark" :content="taskinfo.taskname" placement="top-start">
                                                      <span> {{taskinfo.taskname}}</span>
                                                  </el-tooltip>
                                              </span>
                                          </div>
                                      </el-col>
                                      <el-col :span="8">
                                          <div class="info_block">
                                              <i class="micon iconfont iconchuangjianshijian"></i>
                                              <label for="" class="lbname">创建时间</label>
                                              <span class="spvalue">{{taskinfo.create_time}}</span>
                                          </div>
                                      </el-col>
                                      <el-col :span="8">
                                          <div class="info_block">
                                              <i class="micon iconfont iconzhihangshijian"></i>
                                              <label for="" class="lbname">执行时间 </label>
                                              <span class="spvalue">{{taskinfo.execute_time}}</span>
                                          </div>
                                      </el-col>
                                      <el-col :span="8">
                                          <div class="info_block">
                                              <i class="micon iconfont iconzhihangfangshi"></i>
                                              <label for="" class="lbname">执行方式</label>
                                              <span class="spvalue">{{taskinfo.priority}}</span>
                                          </div>
                                      </el-col>
                                      
                                      <el-col :span="8">
                                          <div class="info_block">
                                              <i class="micon iconfont iconbaogaochangshang"></i>
                                              <label for="" class="lbname">报告厂商</label>
                                              <span class="spvalue">{{taskinfo.manufacturer}}</span>
                                          </div>
                                      </el-col> 
                                      <el-col :span="8">
                                          <div class="info_block">
                                              <i class="micon iconfont icontijiaozhe"></i>
                                              <label for="" class="lbname">提交者 </label>
                                              <span class="spvalue">{{taskinfo.user}}</span>
                                          </div>
                                      </el-col> 
                                  </el-row>
                                  <span  class="iconcircular circular1"></span>
                                  <span  class="iconcircular circular2"></span>
                                  <span  class="iconcircular circular3"></span>
                                  <span  class="iconcircular circular4"></span>
                              </div>
                          </el-col>
                          <el-col :span="6"> 
                              <div  :class="{'info_box info_box_height150 info_box_bg2-high': risk_level=='1',
                                          'info_box info_box_height150 info_box_bg2-mind': risk_level=='2',
                                          'info_box info_box_height150 info_box_bg2-low': risk_level=='3',
                                          'info_box info_box_height150 info_box_bg2-nofind': risk_level=='4',
                                          }" >
                                  <div style="position: relative;z-index:9"> 
                                      <div class="info_block" >
                                          <i class="micon iconfont iconfengxiandengji"></i>
                                          <label for="" class="lbname">任务风险</label>
                                          <span class="spvalue">{{taskinfo.name}}</span>
                                      </div>
                                      <div>
                                          <i class="micon iconfont iconzhuangtai"></i>
                                          <label for="" class="lbname">任务状态 </label>
                                          <span class="spvalue">{{taskinfo.task_status}}</span>
                                      </div>
                                  </div> 
                                  <span  class="iconcircular circular5"></span>
                                  <span  class="iconcircular circular6"></span>
                                  <span  class="iconcircular circular7"></span>
                                  <span  class="iconcircular circular8"></span>
                              </div>
                          </el-col> 
                      </el-row>  
                      <div>  
                          <div  >
                              <div class="part_title">目标统计</div>
                          </div>
                          <el-table
                              ref="multipleTable"
                              :data="targetTableData"
                              style="width: 100%;"
                              class="context_box_bg "
                              >
                              <el-table-column
                                  prop="fileName"
                                  label="文件名">
                              </el-table-column>
                              <el-table-column
                                  prop="fileSize"
                                  label="文件大小">
                              </el-table-column>
                              <el-table-column
                                  prop="number"
                                  label="目标数量">
                              </el-table-column>
                          </el-table>
                      </div>
                  </div>
              </el-tab-pane>
              <el-tab-pane label="测试目标" name="tabs2" >
                  <div class="part1 bgColor2" >
                      <div class="outCont">
                          <div class="targetLabeling">
                              <div style="overflow: hidden;"> 
                                  <div class="tagging_block">
                                          <i class="labelColor colorPurple"></i>
                                          <label for="" class="labelName">导入总目标</label>
                                          <span class="spvalue">{{ tabs2TargetData.all_target_number }}</span>
                                      </div> 
                                  <div class="tagging_block">
                                      <i class="labelColor colorGreen"></i>
                                      <label for="" class="labelName">存活目标</label>
                                      <span class="spvalue">{{ tabs2TargetData.alive_number }}</span>
                                  </div> 
                                      
                              </div>
                              <div style="overflow: hidden;">
                                  <div class="tagging_block">
                                      <i class="labelColor colorWaterRed"></i>
                                      <label for="" class="labelName">高危漏洞目标</label>
                                      <span class="spvalue">{{ tabs2TargetData.high_risk_number }}</span>
                                  </div>
                                  <div class="tagging_block">
                                      <i class="labelColor colorPaleYellow"></i>
                                      <label for="" class="labelName">中危漏洞目标</label>
                                      <span class="spvalue">{{ tabs2TargetData.middle_risk_number }}</span>
                                  </div>
                                  <div class="tagging_block">
                                      <i class="labelColor colorBlue"></i>
                                      <label for="" class="labelName">低危漏洞目标</label>
                                      <span class="spvalue">{{ tabs2TargetData.low_risk_number }}</span>
                                  </div>
                                  <div class="tagging_block">
                                      <i class="labelColor colorBlue"></i>
                                      <label for="" class="labelName">安全目标</label>
                                      <span class="spvalue"> {{ tabs2TargetData.safe_risk_number }}</span>
                                  </div>
                              </div>
                          </div>
                      </div>
                  <div class=""> 
                      <div class="search-box">
                          <div class="operationbutton"  > 
                              <delbutton 
                              :width="170"  
                              @click="btnMultiDeleteTarget"   
                              :disabled="!multipleSelectionTarget.length"></delbutton>
                          </div> 
                          <div class="serach-condition" > 
                              <div class="search-text">
                                  <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearchTarget"  v-model="formData.search" class="input-with-select"  size="small" clearable  > </el-input>
                              
                                  <xzbutton 
                                  type="primary" 
                                  @click="handlesearchTarget" 
                                  :disabled="false" 
                                  size="small"  >搜索</xzbutton> 
                              </div>
                              <div > 
                                  <xzbutton 
                                  type="primary" 
                                  @click="handleResetTarget" 
                                  :disabled="false" 
                                  size="small"  >重置</xzbutton> 
                              </div>  
                          </div>
                      </div>
                      <div class="targetCont">
                          <!--   @select="handChecked2"
                              @select-all="handSelectAll2"
                               -->
                          <el-table  
                            :row-key="rowKey2"
                              ref="multipleTable"
                              :data="tabs2TargetTableData" 
                              tooltip-effect="dark"
                              style="width:100%"
                              v-loading="targetcloading"
                              @selection-change="handleSelectionChangeTarget"
                              @cell-mouse-enter="mouseentertarget" 
                            @cell-mouse-leave="mouseleavetarget" 
                              >
                              <el-table-column
                                  type="selection"
                                  width="55"
                                   >
                              </el-table-column> 
                              <el-table-column
                                  prop="name"
                                  :show-overflow-tooltip="true"
                                  label="任务目标" > 
                                  <template 
                                      slot-scope="scope"  > 
                                      <el-link   @click="btnTasktargetinfo(scope.row.name,'info')">{{scope.row.name}}</el-link> 
                                  </template>
                              </el-table-column> 
                              <el-table-column
                                  prop="os"
                                  label="系统">
                              </el-table-column>
                              <el-table-column
                                  prop="target_risk"
                                  label="风险" > 
                                   <template slot="header"> 
                                      <span class="cursorPointer" @click="clickButton('风险')"
                                          :class="(formData.risk_level !== '' && formData.risk_level !== 0) ? 'active' : ''">风险<i
                                          class="iconfont iconshaixuan"></i>
                                      </span> 
                                      <el-select popper-class="thSelect" style=" width:150px;" v-model="formData.risk_level" clearable
                                      size="small" ref="risk" @change="handlesearchTarget">
                                      <el-option v-for="(item, index) in level" :key="index" :label="item.label" :value="item.value">
                                      </el-option>
                                      </el-select>
                                  </template>
                                  <template slot-scope="scope" >   
                                      <span class="riskstyle risk_hight" v-if="scope.row.risk == 1"><i></i>{{scope.row.riskName}}</span>
                                      <span class="riskstyle risk_middle" v-if="scope.row.risk == 2"><i></i>{{scope.row.riskName}}</span>
                                      <span class="riskstyle risk_low" v-if="scope.row.risk == 3"><i></i>{{scope.row.riskName}}</span>
                                      <span class="riskstyle risk_nofind" v-if="scope.row.risk == 4"><i></i>{{scope.row.riskName}}</span> 
                                  </template>
                              </el-table-column>
                              <el-table-column  label="漏洞" 
                                  width="300px"
                                  :render-header="icons">
                                      <template slot-scope="scope" >   
                                          <span class="tag_status tag_danger bug_status" v-if="scope.row.exp >99" >99+</span>
                                          <span class="tag_status tag_danger bug_status" v-if="scope.row.exp <=99" >{{scope.row.exp}}</span>
                                          <span class="tag_status tag_warning bug_status" v-if="scope.row.verify >99" >99+</span>
                                          <span class="tag_status tag_warning bug_status" v-if="scope.row.verify <=99"  >{{scope.row.verify}}</span>
                                          <span class="tag_status tag_primary bug_status" v-if="scope.row.failed >99"  >99+</span>
                                          <span class="tag_status tag_primary bug_status" v-if="scope.row.failed <=99" >{{scope.row.failed}}</span>
                                          <span class="tag_status tag_success bug_status" v-if="scope.row.unVerify >99"  >99+</span>
                                          <span class="tag_status tag_success bug_status" v-if="scope.row.unVerify <=99" >{{scope.row.unVerify}}</span>
                                          <!-- <span class="tag_status tag_total bug_status" v-if="scope.row.all_vulns_number <=99" >{{scope.row.all_vulns_number}}</span> -->
                                      </template> 
                              </el-table-column> 
                              <!-- <el-table-column
                                  prop="target_status"
                                  label="状态"  > 
                                  <template slot-scope="scope" >      
                                      <span class="tag_status tag_warning" v-if="scope.row.target_status == 1">待开始</span>
                                      <span class="tag_status tag_primary" v-if="scope.row.target_status == 2">运行中</span>
                                      <span class="tag_status tag_success" v-if="scope.row.target_status == 3">已结束</span>
                                      <span class="tag_status tag_danger" v-if="scope.row.target_status == 4">暂停中</span>
                                  </template>
                              </el-table-column>  -->
                              <el-table-column label="状态"  > 
                                  <template slot-scope="scope" > 
                                      <div v-if="showOperateButton && rowId == scope.row.id ">
                                          <el-link 
                                              class="link_danger" 
                                              :underline="false"   >
                                                  <el-popover
                                                      placement="bottom"
                                                      width="170"   
                                                      :visible-arrow="false"
                                                      :ref="`popover_id-${scope.row.id}`"
                                                      popper-class="delButton_popper" >
                                                      <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                                      <div style="text-align: right; margin: 0">
                                                          <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                                          <el-button size="mini" type="primary" @click="targetDelete(scope,'yes')">确定</el-button>
                                                      </div>  
                                                      <span  slot="reference"  >删除</span>
                                                  </el-popover>
                                          </el-link>
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
                          <el-pagination
                              :page-size="targetPageSize" 
                              background
                              layout=" total, prev, pager, next, sizes, jumper"
                              :total="targetTotalpage"
                              :current-page="targetCurrentpage"
                              @current-change = "currentchange"
                              @size-change="handleSizeChangetarget"  >
                          </el-pagination> 
                      </div>
                      
                  </div>
                  </div>
              </el-tab-pane>
              <el-tab-pane label="开放端口" name="tabs3" >
                  <div class="part1">
                      <div class="search-box">
                          <div class="serach-condition" > 
                              <!-- <div>
                                   <el-select v-model="portFormData.target_ids" placeholder="目标类型" size="small" class="selectwidth">
                                      <el-option
                                      v-for="(item,index) in targetlist"
                                      :key="index"
                                      :label="item[1]"
                                      :value="item[0]">
                                      </el-option>
                                  </el-select>
                              </div> -->
                              <!-- <div class="search-text">
                                  <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearchPort"  handlesearchTarget v-model="portFormData.search_fields" class="input-with-select"  size="small"  > </el-input>
                         
                                  <xzbutton 
                                  type="primary" 
                                  @click="handlesearchPort" 
                                  :disabled="false" 
                                  size="small"  >搜索</xzbutton> 
                                  
                              </div>
                              <div > 
                                  <xzbutton 
                                  type="primary" 
                                  @click="handleResetReport" 
                                  :disabled="false" 
                                  size="small"  >重置</xzbutton>
                              </div>   -->
                          </div>
                      </div>
                      <el-table
                          :data="portListtableData" 
                          tooltip-effect="dark"
                          style="width: 100%"
                          v-loading="portLoading"
                          >
                          <el-table-column
                              prop="target"
                              label="任务目标">
                          </el-table-column>
                          <el-table-column
                              prop="port"
                              label="端口">
                          </el-table-column>
                          <el-table-column
                              prop="scheme"
                              label="协议">
                          </el-table-column>
                          <el-table-column
                              prop="service"
                              label="服务">
                          </el-table-column>
                          <el-table-column
                              prop="component"
                              :show-overflow-tooltip="true"
                              label="应用/组件">
                          </el-table-column> 
                      </el-table>
                      <el-pagination
                          :page-size="portPageSize" 
                          background
                          layout=" total, prev, pager, next, sizes, jumper"
                          :total="portTotalpage"
                          :current-page="portCurrentpage"
                          @current-change = "portCurrentchange"
                          @size-change="handleSizeChangePort"  >
                      </el-pagination>
                  </div>
              </el-tab-pane>
              <el-tab-pane label="漏洞信息" name="tabs4" >
                  <div class="part1 bgColor2" >
                      <div class="outCont">
                          <div class="targetLabeling">
                              <div class="tagging_block">
                                  <i class="labelColor loopOldRed"></i>
                                  <label for="" class="labelName">导入总漏洞</label>
                                  <span class="spvalue">{{ loop.all_vuln_number }}</span>
                              </div>
                              <div class="tagging_block">
                                  <i class="labelColor loopDeepBlue"></i>
                                  <label for="" class="labelName">验证成功漏洞</label>
                                  <span class="spvalue">{{ loop.prove_success_number }}</span>
                              </div>
                              <div class="tagging_block">
                                  <i class="labelColor loopGreen"></i>
                                  <label for="" class="labelName">利用成功漏洞</label>
                                  <span class="spvalue">{{ loop.use_success_number }}</span>
                              </div>
                              <div class="tagging_block">
                                  <i class="labelColor loopLightYellow"></i>
                                  <label for="" class="labelName">验证失败漏洞</label>
                                  <span class="spvalue">{{ loop.prove_faile_number }}</span>
                              </div>
                              <div class="tagging_block">
                                  <i class="labelColor loopWaterRed"></i>
                                  <label for="" class="labelName">未能验证漏洞</label>
                                  <span class="spvalue">{{ loop.no_prove_number }}</span>
                              </div>  
                              
                              <div class="tagging_block">
                                  <i class="labelColor loopOrange"></i>
                                  <label for="" class="labelName">高危漏洞</label>
                                  <span class="spvalue">{{ loop.high_vuln_number }}</span>
                              </div>
                              <div class="tagging_block">
                                  <i class="labelColor loopYellow"></i>
                                  <label for="" class="labelName">中危漏洞</label>
                                  <span class="spvalue">{{ loop.middle_vuln_number }}</span>
                              </div>
                              <div class="tagging_block">
                                  <i class="labelColor loopBlue"></i>
                                  <label for="" class="labelName">低危漏洞</label>
                                  <span class="spvalue">{{ loop.low_vuln_number }}</span>
                              </div>
                          </div>
                      </div>
                      <div class="search-box"> 
                          <div class="operationbutton"  >   
                              <delbutton 
                              :width="170"  
                              @click="btnMultiDeleteLoop"   
                              :disabled="!multipleSelectLoop.length" style="margin-left: 8px;"></delbutton>
                          </div> 
                          <div class="serach-condition" > 
                              <!-- <div > 
                                  <el-select v-model="formDataLoophole.vuln_risk"   placeholder="漏洞风险" class="selectwidth" size="small">  
                                      <el-option
                                          v-for="(item,index) in buglevel"
                                          :key="index"
                                          :label="item[1]"
                                          :value="item[0]">
                                          </el-option>
                                  </el-select> 
                              </div>  
                              <div > 
                                  <el-select v-model="formDataLoophole.vuln_status"   placeholder="漏洞状态" class="selectwidth" size="small">  
                                      <el-option
                                          v-for="(item,index) in verificationStatus"
                                          :key="index"
                                          :label="item[1]"
                                          :value="item[0]">
                                          </el-option>
                                  </el-select> 
                              </div>  -->
                              <div class="search-text">
                                  <el-input placeholder="请输入关键字"  v-model="formDataLoophole.search_fields" class="input-with-select"  size="small" clearable  > </el-input>
                                 
                                  <xzbutton 
                                  type="primary" 
                                  @click="loopHandlesearchbug" 
                                  :disabled="false" 
                                  size="small"  >搜索</xzbutton>  
                              </div>
                              <div > 
                                  <xzbutton 
                                  type="primary" 
                                  @click="loopHandleResetbug" 
                                  :disabled="false" 
                                  size="small"  >重置</xzbutton> 
                              </div>  
                          </div>
                      </div>
                      <!--    :row-key="rowKey" @select="handChecked4"
                          @select-all="handSelectAll4" -->
                      <el-table
                          ref="multipleTableBug"
                          :data="buglisttableData" 
                          tooltip-effect="dark"
                          v-loading="loopLoading"
                          :row-key="rowKey"
                          style="width: 100%"  
                          @selection-change="handleSelectionChangeLoop">
                          <el-table-column
                              type="selection"
                              width="55"
                              :reserve-selection="true">
                          </el-table-column> 
                          <el-table-column
                              prop="name"
                              label="漏洞名称" 
                              :show-overflow-tooltip="true"> 
                              <template slot-scope="scope"> 
                                  <el-link @click="btnbuginfo(scope.row)">{{ scope.row.name }}</el-link> 
                              </template>
                          </el-table-column>
                          <el-table-column
                              prop="vuln_risk"
                              label="漏洞风险" >
                          <template slot="header"> 
                              <span class="cursorPointer" @click="clickButton('漏洞风险')"
                                  :class="(formDataLoophole.vuln_risk !== '' && formDataLoophole.vuln_risk !== 0) ? 'active' : ''">漏洞风险<i
                                  class="iconfont iconshaixuan"></i>
                              </span> 
                              <el-select popper-class="thSelect" style=" width:150px;" v-model="formDataLoophole.vuln_risk" clearable
                              size="small" ref="bugrisk" @change="loopHandlesearchbug">
                              <el-option v-for="(item, index) in buglevel" :key="index" :label="item.label" :value="item.value">
                              </el-option>
                              </el-select>
                          </template>
                              <template slot-scope="scope" >   
                                  <span class="riskstyle risk_hight" v-if="scope.row.risk == 1"><i></i>{{scope.row.riskName}}</span>
                                  <span class="riskstyle risk_middle" v-if="scope.row.risk == 2"><i></i>{{scope.row.riskName}}</span>
                                  <span class="riskstyle risk_low" v-if="scope.row.risk == 3"><i></i>{{scope.row.riskName}}</span>
                                  <span class="riskstyle risk_nofind" v-if="scope.row.risk == 4"><i></i>{{scope.row.riskName}}</span>
                              </template>
                          </el-table-column>
                          <el-table-column
                              prop="vuln_status" :show-overflow-tooltip="true"
                               
                              label="漏洞状态" >
                              <template slot="header"> 
                                  <span class="cursorPointer" @click="clickButton('漏洞状态')"
                                      :class="(formDataLoophole.vuln_status !== '' && formDataLoophole.vuln_status !== 0) ? 'active' : ''">漏洞状态<i
                                      class="iconfont iconshaixuan"></i>
                                  </span> 
                                  <el-select popper-class="thSelect" style=" width:150px;" v-model="formDataLoophole.vuln_status" clearable
                                  size="small" ref="vuln_status" @change="loopHandlesearchbug">
                                  <el-option v-for="(item, index) in verificationStatus" :key="index" :label="item.label" :value="item.value">
                                  </el-option>
                                  </el-select>
                              </template>
                              <template slot-scope="scope" >   
                                  <span class="tag_status tag_waitting" v-if="scope.row.status == 1"><i></i>{{scope.row.statusName}}</span>
                                  <span class="tag_status tag_fail" v-if="scope.row.status == 2"><i></i>{{scope.row.statusName}}</span>
                                  <span class="tag_status tag_success" v-if="scope.row.status == 3"><i></i>{{scope.row.statusName}}</span>
                                  <span class="tag_status tag_warning" v-if="scope.row.status == 4"><i></i>{{scope.row.statusName}}</span>  
                                  <span class="tag_status tag_primary" v-if="scope.row.status == 5"><i></i>{{scope.row.statusName}}</span>   
                              </template>
                          </el-table-column>
                          <el-table-column
                              prop="location" :show-overflow-tooltip="true"
                              label="漏洞位置" >
                          </el-table-column> 
                          <el-table-column label="操作"  >
                              <template slot-scope="scope">  
                                  <!-- <el-link :underline="false" class="link_primary"  @click="btnbuguse(scope.row)" style="padding-left: 0px;" v-if="scope.row.can_used">利用</el-link> -->
                                  <!-- <el-link :underline="false" class="link_primary" 
                                      @click="btnverification(scope.row)" 
                                      style="padding-left: 0px;" 
                                      v-if="scope.row.status!=2">
                                          <span >验证</span>
                                  </el-link>
                                  <el-link :underline="false" class="link_default" style="cursor: default;padding-left: 0px;"  v-else>- -</el-link> -->
                                  <el-popover
                                      placement="bottom"
                                      width="200"   
                                      :visible-arrow="false"
                                      :ref="`popover-${scope.row.id}`"
                                      popper-class="delButton_popper" >
                                      <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                      <div style="text-align: right; margin: 0">
                                          <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover-${scope.row.id}`].doClose()">取消</el-button>
                                          <el-button size="mini" type="primary" @click="btnLoopDel(scope,'yes')">确定</el-button>
                                      </div> 
                                      <el-link :underline="false" class="link_danger " style="padding-right:10px" slot="reference" >删除</el-link>  
                                  </el-popover> 
                                  
                              </template>
                          </el-table-column>
                      </el-table>
                      <el-pagination 
                          background
                          @size-change="handleSizeloop"
                          @current-change="handleCurrentloop"
                          :current-page="loopCurrentPage" 
                          :page-size="loopPageSize"
                          layout="total,  prev, pager, next, sizes, jumper"
                          :total="loopTotal">
                      </el-pagination>
                  </div>
              </el-tab-pane>
              <!-- <el-tab-pane label="验证日志" name="tabs5" >
                  <div class="part1 ">
                      <div class="search-box">  
                          <div class="operationbutton"  > 
                              <el-button type="warning"  size="small"  slot="reference" @click="emptyLog">清空</el-button> 
                          </div> 
                          <div class="serach-condition" > 
                              <div class="search-text">
                                  <el-input placeholder="请输入关键字" @keydown.enter.native="handlesearchlog"  v-model="formDataJournal.search" class="input-with-select"  size="small"> </el-input>
                                  <xzbutton 
                                  type="primary" 
                                  @click="handlesearchlog" 
                                  :disabled="false" 
                                  size="small"  >搜索</xzbutton>  
                              </div>
                              <div > 
                                  <xzbutton 
                                  type="primary" 
                                  @click="handleResetlog" 
                                  :disabled="false" 
                                  size="small"  >重置</xzbutton> 
                              </div>  
                          </div>
                      </div>
                      <el-table
                          ref="multipleTableLog"
                          :data="logtableData" 
                          tooltip-effect="dark"
                          style="width: 100%" 
                          v-loading="logLoading"
                          > 
                          
                          <el-table-column
                              prop="pocname"
                              label="pocname" > 

                          </el-table-column>
                          <el-table-column
                                  prop="type_enum"
                                  label="类型" >
                          </el-table-column>
                           <el-table-column
                                  prop="content"
                                  label="内容" >
                          </el-table-column>
                          <el-table-column
                                  prop="create_time"
                                  label="创建时间" >
                          </el-table-column>
                      </el-table>
                      <el-pagination 
                          background
                          @size-change="handleSizeChangelog"
                          @current-change="handleCurrentChangelog"
                          :current-page="currentPageLog" 
                          :page-size="pageSizeLog"
                          layout="total,  prev, pager, next, sizes, jumper"
                          :total="totalLog">
                      </el-pagination>
                  </div>
              </el-tab-pane> -->
          </el-tabs>
      </div>

      <el-dialog title="变更状态" 
          :visible.sync="upadtestatusdialogVisible"    
          width="1184px"
          :close-on-click-modal="false" 
          :validate-on-rule-change="false" 
          :show-close="false"
          class="updatestatus">
          <div class="dialog_b_btn">   
              <el-button size="small" @click="changeStatus">确定</el-button>
              <el-button size="small" @click="upadtestatusdialogVisible = false">关闭</el-button>
          </div>
          <div style="padding:24px">
              <div  >
                  <el-form
                      ref="statusform"
                      :model="statusform"
                      label-width="0"
                      class="clearfix"
                      style="text-align: center;" >
                      <el-form-item label=""   style="margin-bottom:0">
                          <el-select v-model="statusform.status" placeholder="请选择" class="selstatus" style="width:360px;height:40px">
                              <el-option
                                  v-for="(item,index) in statuslist"
                                  :key="index"
                                  :label="item[1]"
                                  :value="item[0]">
                              </el-option>
                          </el-select>
                      </el-form-item> 
                  </el-form>
              </div> 
          </div> 
      </el-dialog>

      <!-- 验证 -->
      <el-dialog :title="taskTitle+'验证'" 
          :visible.sync="TestdialogVisible"    
          width="1184px"
          :close-on-click-modal="false" 
          :validate-on-rule-change="false" 
          :show-close="false">
          <div class="dialog_b_btn">   
              <el-button size="small" @click="startTest"  v-loading="yzloading">开始验证</el-button>
              <el-button size="small" @click="closeTest()">关闭</el-button>
          </div>
          <div style="padding:24px">
              <div  style="padding-bottom:32px;color:rgba(72, 72, 102, 0.64);border-bottom:1px solid  #E8E8F5">
                  {{target_result}}
              </div> 
              <div style="margin-top:26px">
                  <div class="controlbox" v-for="(item,index) in verify_result" :key = index >
                      <div v-if="item.str">
                           {{item.str}}
                      </div>
                      <el-row :gutter="20" v-else>
                          <el-col :span="12">
                              <div class="message" >
                                  <label class="title_bg title_bg1">请求报文</label> 
                                  <div><pre>{{item.back[0]}}</pre></div>
                              </div> 
                          </el-col> 
                          <el-col :span="12"> 
                              <div class="message" >
                                  <label class="title_bg title_bg2">响应报文</label> 
                                  <div><pre>{{item.back[1]}}</pre></div>
                              </div>
                          </el-col> 
                      </el-row>
                  </div>
              </div>
          </div> 
      </el-dialog>

      <!-- 漏洞详情 -->
      <el-dialog
          :title="bugform.vuln_name"
          :visible.sync="bugdialogVisible"
          width="1184px"
          class="buginfobox" 
          :close-on-click-modal="false" 
          :show-close="false" >
          <div class="dialog_b_btn">   
              <el-button size="small" @click="cancalbugdialogVisible">关闭</el-button>
          </div>
          <div class="buginfo_box" > 
              <div class="bugbasicinfo">
                  <el-table
                      :data="bugbasicinfo"
                      size='small'
                      style="width: 100%">
                      <!-- <el-table-column
                          prop="vuln_type"
                          label="漏洞类型" >
                          <template slot-scope="scope">
                              <span>{{scope.row.vuln_type}}</span> 
                          </template>
                      </el-table-column> -->
                      <el-table-column
                          prop="vuln_risk"
                          label="漏洞风险" >
                          <template slot-scope="scope">
                              <span v-if="!is_bugUpdate" > 
                                  <span class="riskstyle risk_hight" v-if="scope.row.vuln_risk == 1"><i></i>{{scope.row.vuln_risk_name}}</span>
                                  <span class="riskstyle risk_middle" v-if="scope.row.vuln_risk == 2"><i></i>{{scope.row.vuln_risk_name}}</span>
                                  <span class="riskstyle risk_low" v-if="scope.row.vuln_risk == 3"><i></i>{{scope.row.vuln_risk_name}}</span>
                                  <span class="riskstyle risk_nofind" v-if="scope.row.vuln_risk == 4"><i></i>{{scope.row.vuln_risk_name}}</span> 
                              </span>
                              <el-select v-model="updateinfo.vuln_rsik"   size="mini" v-if="is_bugUpdate" >  
                                  <el-option
                                      v-for="(item,index) in buglevel"
                                      :key="index"
                                      :label="item.label"
                                      :value="item.value">
                                      </el-option>
                              </el-select> 
                          </template>
                      </el-table-column>
                      <el-table-column
                          prop="cvss"
                          label="CVSS评分">
                      </el-table-column>
                      <!-- <el-table-column
                          prop="vuln_object"
                          label="漏洞对象">
                      </el-table-column>
                      <el-table-column
                          prop="vul_id"
                          label="VUL_ID">
                          <template slot-scope="scope">
                              <span  >{{scope.row.vul_id}}</span> 
                          </template>
                      </el-table-column> -->
                      <!-- <el-table-column
                          prop="vuln_status"
                          label="漏洞状态">
                          <template slot-scope="scope"> 
                            
                                  <span class="tag_status tag_waitting" v-if="scope.row.vuln_status == 1"><i></i>待验证</span>
                                  <span class="tag_status tag_fail" v-if="scope.row.vuln_status == 2"><i></i>未能验证</span>
                                  <span class="tag_status tag_success" v-if="scope.row.vuln_status == 3"><i></i>验证成功</span>
                                  <span class="tag_status tag_warning" v-if="scope.row.vuln_status == 4"><i></i>验证失败</span>  
                                  <span class="tag_status tag_primary" v-if="scope.row.vuln_status == 5"><i></i>利用成功</span>
                                  <i class="iconfont iconduocijiance" v-if="scope.row.is_repeat"></i> 
                           
                              <el-select v-model="updateinfo.status" size="mini" v-if="is_bugUpdate">
                                  <el-option
                                      v-for="(item,index) in statuslist"
                                      :key="index"
                                      :label="item[1]"
                                      :value="item[0]">
                                  </el-option>
                              </el-select>

                          </template>
                      </el-table-column> -->
                      <!-- <el-table-column
                          prop="vuln_time"
                          label="公布时间">
                      </el-table-column> -->
                  </el-table>
              </div>
              <div class="bugotherinfo">
                  <div class="part_title">漏洞描述</div>
                  <div class="content" v-if="!is_bugUpdate"> {{bugform.detail}}  </div>
                  <el-input class="textarea" type="textarea" v-model="updateinfo.detail"  size="mini" :row="3" v-if="is_bugUpdate" ></el-input>
              </div>
              <div class="bugotherinfo">
                  <div class="part_title">修复建议</div>
                  <div class="content" v-if="!is_bugUpdate">{{bugform.fix_suggest}}</div>
                  <el-input class="textarea" type="textarea" v-model="updateinfo.fix_suggest"  size="mini"  :row="3" v-if="is_bugUpdate" ></el-input>
              </div>
               <!-- <div class="bugotherinfo">
                  <div class="part_title">漏洞结果</div>
                  <div class="content" v-if="!is_bugUpdate">{{bugform.vuln_result}}</div>
                  <el-input class="textarea" type="textarea" v-model="updateinfo.vuln_result"  size="mini"  :row="3" v-if="is_bugUpdate" ></el-input>
              </div> -->
              <div class="bugotherinfo">
                  <div class="part_title">cve编号</div>
                  <div class="content" >{{ bugform.cve }}</div> 
              </div>
               <div class="bugotherinfo">
                  <div class="part_title">cnnvd编号</div>
                  <div class="content" >{{ bugform.cnnvd }}</div> 
              </div>
              <!-- <div class="bugotherinfo">
                  <div class="part_title">参考链接</div>
                  <div class="content" v-if="!is_bugUpdate">{{bugform.ref_url}}</div>
                  <el-input class="textarea" type="textarea" v-model="updateinfo.ref_url"  size="mini"  :row="3" v-if="is_bugUpdate" ></el-input>
              </div>
              <div class="bugotherinfo">
                  <div class="part_title">漏洞地址</div>
                  <div class="content" v-if="!is_bugUpdate">{{bugform.vuln_location}}</div>
                   <el-input class="textarea" type="textarea" v-model="updateinfo.vuln_location"  size="mini"  :row="3" v-if="is_bugUpdate" ></el-input>
              </div>
              <div class="bugotherinfo">
                  <div class="part_title">漏洞验证报文</div>
                  <div class="bugbasicinfo">  
                      <el-table
                          :data="bugmessage"
                          size='small'
                          row-key="id"
                          :expand-row-keys="expands"
                          @expand-change="handleExpandChange"
                          @row-click="rowClick"
                          style="width: 100%">
                          <el-table-column type="expand">
                              <template slot-scope="scope">
                                  <el-row :gutter="20">
                                      <el-col :span="12">
                                          <div class="message" >
                                              <label class="title_bg title_bg1">请求报文</label> 
                                              <div><pre>{{scope.row.request_pack}}</pre></div>
                                          </div>
                                          <div class="message" >
                                              <label class="title_bg title_bg1">响应报文</label> 
                                              <div><pre>{{scope.row.response_pack}}</pre></div>
                                          </div>
                                      </el-col>
                                      <el-col :span="12">
                                          <div class="message requestpack" >
                                              <label class="title_bg title_bg2">请求报文</label>
                                              <div>
                                                  <div class="packheight " >
                                                      <div class="packinput" > 
                                                          <el-input
                                                              class="packtxt"
                                                              type="textarea"
                                                              :rows="7" 
                                                              v-model="requestpack">
                                                          </el-input> 
                                                      </div>
                                                      <div class="packbtn"> 
                                                          <el-button  style="padding:9px 24px; float: right;" size="mini" type="primary"  @click="btnpacket(scope)">验证</el-button>
                                                      </div>
                                                  </div>
                                              </div>
                                          </div>
                                          <div class="message" >
                                              <label class="title_bg title_bg2">响应报文</label>
                                              <div><pre>{{responsepack}}</pre></div>
                                          </div>
                                      </el-col> 
                                  </el-row>
                              </template>
                          </el-table-column>
                          <el-table-column
                              type="index" 
                              label="步骤" >
                          </el-table-column>
                          <el-table-column
                              prop="method"
                              label="HTTP方法" > 
                          </el-table-column>
                          <el-table-column
                              prop="url"
                              label="URL">
                          </el-table-column>
                          <el-table-column
                              prop="status"
                              label="状态码">
                          </el-table-column> 
                      </el-table> 
                  </div> 
              </div>
               -->
          </div>
      </el-dialog>

      <!-- <reportTargetList 
          v-if="rightMenu"
          ref="reportList"
          :activeName="activeTableName" 
          :rangs="Listrangs"
          :taskId='task_id'
          @tabsEvent="tabsEventFunction(arguments)">
      </reportTargetList> -->
  </div>
</template>
<style lang="less" scoped>
 @import "../css/reportTaskInfo.less";
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
</style>
<script>
// import reportTargetList from '@/components/ReportTargetList.vue'
// import basicsinfo from "@/components/Basicsinfo.vue";
import $ from 'jquery'
import xzbutton from "@/components/XzButton.vue"; 
import delbutton from "@/components/DelButton.vue";
export default({
  name:'',
  components: {
    //   reportTargetList, 
    //   basicsinfo,
      xzbutton,
      delbutton, 
  },
  data(){
      return{
          timermillisec:0 ,//时间间隔
          task_id:this.$route.query.id,
          task_name:this.$route.query.name,
          risk_level:this.$route.query.risk_level,
          risk_status:this.$route.query.status,
          // task_name:'',
          activeTableName:'tabs1',
          activeName:'tabs1', 
          taskinfo:{
              taskname:'',
              create_time:'',
              execute_time:'',
              priority:'',
              manufacturer:'',
              user:'',
              risk_level:'',
              task_status:''
          },
          targetTableData:[
              // {fileName:'sdjk',fileSize:'89',number:'78'}
          ],
          targetlist:[[0,'全部'],[1,'存活目标'],[2,'不存活目标']],
          level:[ ],
          formData:{target_type:0,risk_level:'',search:'',page_num:1},
          portFormData:{
              target_ids:1,
              search_fields:'',
              page_num:1
          },
          multipleSelectionTarget:[],//测试目标多选框组
          checkedList2:[],//存测试目标
          portListtableData:[],
          VerifiedLoop:[[1,'已验证漏洞'],[2,'未能验证漏洞']],
          buglevel:[ ],
          verificationStatus:[],
          formDataLoophole:{
              already_prove:1,
              vuln_risk:'',
              vuln_status:'',
              search_fields:'',
              page_num:1
          },
          buglisttableData:[],
          multipleSelectionbug:[],
          upadtestatusdialogVisible:false,
          statuslist:[
              [1,'待验证'],[2,'未能验证'],[3,'验证成功'],[4,'验证失败'],[5,'利用成功']
          ],
          statusform:{
              status:1,
          },
          formDataJournal:{
              search:'',
              page_num:1
          },
          logtableData:[],
          targetTotalpage:0,
          targetCurrentpage:1,
          targetPageSize:10,
          rightMenu:false,
          tabs2TargetData:{
              risk_target_number:'',
              no_risk_number:'',
              alive_number:'',
              no_alive_number:'',
              all_target_number:'',
              high_risk_number:'',
              middle_risk_number:'',
              low_risk_number:'',
              safe_risk_number:'',
          },
          tabs2TargetTableData:[
              {
                  target_url:'',
                  opearting_system:'',
                  target_risk:'',
                  use_success_number:'',
                  verify_success_number:'',
                  verify_failed_number:'',
                  cannot_verify_number:'',
                  all_vulns_number:'',
                  target_status:''
              }
          ],
          portPageSize:10,
          portTotalpage:0,
          portCurrentpage:1,
          loop:{
              prove_success_number:'',
              use_success_number:'',
              prove_faile_number:'',
              no_prove_number:'',
              wait_prove_number:'',
              all_vuln_number:'',
              high_vuln_number:'',
              middle_vuln_number:'',
              low_vuln_number:'',
              // page_num:1
          },
          loopCurrentPage:1,
          loopTotal:0,
          loopPageSize:10,
          multipleSelectLoop:[],
          checkedList4:[],//漏洞信息多选框组
          target_ids:'',
          rangs:'',
          currentPageLog:1,
          pageSizeLog:10,
          totalLog:0,
          targetcloading:false,
          portLoading:false,
          loopLoading:false,
          logLoading:false,
          checkedAll:true,
          Listrangs:[],
          taskTitle:'',
          TestdialogVisible:false,
          yzloading:false,
          verify_result:[],
          target_result:'',
          timer:null,
          request_pack:'',
          response_pack:'',
          bugform:{
             vuln_name:'',
             detail:'',
             fix_suggest:'',
             vuln_result:'',
             ref_url:'',
             vuln_location:'',
             vul_id:''         
          },
          bugmessage:[],
          bugdialogVisible:false,
          is_bugUpdate:false,
          updateinfo:{
              check_vuln_id:'',
              name:'',
              object:'',
              type:'',
              risk_lever:'',
              cve:'',  
              cnvd:'',    
              cnnvd:'',    
              detail:'',   
              fix_suggest:'',
              ref_url:'',    
              vuln_location:'',
              risk_lever_number:'',
              status:[],
          },
          bugbasicinfo:[],
          requestpack:'',
          responsepack:'',
          expands: [],  // 要展开的行，数值的元素是row的key值
          pocname:'',
          vuln_id:'',
          target_id:'',
          alldelvisible:false,
          bugalldelvisible:false,
          basicsinfo:[],
          risklevel:[],
          difModel:'验证报告',
          timerTable : null,
          timeTest:null,
          timeOpen:null,
          timeBug:null,
          timeLog:null,
          showOperateButton:false,
          rowId:''
      }
  },
  created:function(){
    //   this.timermillisec = this.commonjs.timermillisec;
      this.$store.state.activefirstMenu="/verificationReport"; 
    //   // this.pageSize = this.commonjs.pageSize; 
    //   this.targetPageSize = 10; 
    //   this.portPageSize = 10; 
    //   this.loopPageSize = 10; 
    //   this.pageSizeLog = 10; 
    //   this.task_id = localStorage.getItem('task_id');
    //   // 测试目标定时器
    //   this.timeBug = setInterval(()=>{
    //       this.tabs4TableData(false);
    //   },this.timermillisec)
    //       clearTimeout(this.timerTable)
    //   clearTimeout(this.timeTest)
    //   clearTimeout(this.timeOpen)
    //   clearTimeout(this.timeLog)
    this.getEnum();
  },
  mounted:function(){ 
    //   this.taskParameters();
      var _tab =  !(localStorage.getItem('taskTab')) ? 'tabs1' :localStorage.getItem('taskTab') ;
      this.activeTableName = _tab;
      if(  this.activeTableName =='tabs5'){
          this.rightMenu = false; 

          //列表
        //   this.verificationLogData(false); 
      }else{
          this.rightMenu = false;
      }

      if(this.activeTableName == 'tabs1'){
          this.taskParameters();
      }
      if(this.activeTableName == 'tabs2'){
          this.testObjectives();  
          this.testTargetData(true);
      }
      if (this.activeTableName == 'tabs3') {
          this.openReport();
      }
      if (this.activeTableName == 'tabs4') { 
        this.testObjectives();  
          this.tabs4TableData();
      }

  },
  beforeDestroy () {
      clearInterval(this.timerTable)
      clearInterval(this.timeTest)
      clearInterval(this.timeOpen)
      clearInterval(this.timeBug)
      clearInterval(this.timeLog)
  }, 
  methods:{
        getEnum() { 
            this.$ajax.get('/smart/reportverify/enum').
            then((res) =>{
                var res = res.data;    
                if(res.code == 200){ 
                    this.level = res.data.risk;
                    this.level.unshift( { label: "全部", value: 0 }); 
                    this.buglevel = res.data.vulRisk;
                    this.buglevel.unshift( { label: "全部", value: 0 }); 
                    this.verificationStatus = res.data.vulStatus;
                    this.verificationStatus.unshift( { label: "全部", value: 0 }); 
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
              case '风险':
                  this.$refs.risk.toggleMenu();
                  break;
              case '状态':
                  this.$refs.status.toggleMenu();
                  break;
              case '漏洞风险':
                  this.$refs.bugrisk.toggleMenu();
                  break;
              case '漏洞状态':
                  this.$refs.vuln_status.toggleMenu();
                  break;
          }
      },
      rowKey2(row){
          return row.id
      },
      rowKey(row){
        return row.vuln_id
      },
      handChecked2(sel,row){
        this.checkedList2 = sel
      },
      handSelectAll2(ael){
        this.checkedList2 = ael
      },
       handChecked4(sel,row){
        this.checkedList4 = sel
      },
      handSelectAll4(ael){
        this.checkedList4 = ael
      },
      handleClick(val){
          if(  val.name =='tabs5'){
            //   this.rightMenu = false;
              // this.activeTableName = val.name;
              // this.$nextTick(() =>{
              //     this.$refs.reportList.getData();
              // })
              // $('.menu').removeClass('openAnimate').addClass('putAway');
              // this.btnStyle = false;
              // this.checkedAll = true;
              // this.Listrangs = [];
              // // if(val.name == 'tabs3'){
              // //     this.openReport();
              // // }
              // if(val.name == 'tabs4'){ 
              //     // this.tabs4TableData();
              // }
            //    this.verificationLogData(false);
          }else{
              this.rightMenu = false;
          }
          if(val.name == 'tabs1'){
              this.taskParameters();
          }
          if(val.name == 'tabs2'){
              this.testObjectives();
              this.testTargetData(true);
          }
          if(val.name == 'tabs3'){
             this.openReport();
          }
          if (val.name == 'tabs4') { 
            this.testObjectives();
              this.tabs4TableData();
          }
          this.activeTableName = val.name;
          localStorage.setItem('taskTab', val.name);
      },
      icons(h,{column}){
          const inReview = '从左至右依次是“利用成功”、“验证成功”、“验证失败”、“未能验证”' 
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
                      }
                  }, [
                      h('div', {
                          slot: 'content',
                          style: {
                              // 'width':'100px',
                              whiteSpace: 'normal', 
                          }
                      }, inReview), 
                      h('i', {
                          class: 'iconfont icontishi',
                          style: 'color:rgba(72,72,102,0.32);margin-left:5px;vertical-align: initial;'
                      })
                  ],)
              ],
  　　　　 )
      },
      bugIcons(h,{column}){
          const inReview = `利用成功，通过EXP验证发现漏洞存在；验证成功，通过POC验证发现漏洞存在；验证失败：通过POC验证发现漏洞不存在；未能验证，未对漏洞进行验证`
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
                      }
                  }, [
                      h('div', {
                          slot: 'content',
                          tt:'ee',
                          style: {
                              'width':'250px',
                              whiteSpace: 'normal', 
                          }
                      }, inReview), 
                      h('i', {
                          class: 'iconfont icontishi',
                          style: 'color:rgba(72,72,102,0.32);margin-left:5px;vertical-align: initial;'
                      })
                  ],)
              ],
  　　　　 )
      },
      updateStatus(){
          this.upadtestatusdialogVisible = true;
      },
      taskParameters(){
        this.targetTableData = [];
          this.$ajax.get('/smart/reportverify/taskdetail',{
              params:{
                  id:this.task_id
              }
          }).
          then((res) =>{ 
              if(res.data.code == 200){
                  let dt = res.data.data;
                  this.taskinfo.taskname = dt.name;
                  this.taskinfo.create_time = dt.createTime;
                  this.taskinfo.execute_time = dt.updateTime;
                  this.taskinfo.priority = dt.executeTypeName;
                  this.taskinfo.manufacturer = dt.producerName;
                //   this.targetTableData = dt.file_list;
                  this.taskinfo.user = dt.user;
                  this.taskinfo.name = dt.riskName;
                  this.taskinfo.task_status = dt.statusName;

                //   this.getSubdata(this.taskinfo);// 传入子组件的数据
               
                let filelist = {
                    fileName:dt.fileName,fileSize:dt.fileSize,number:dt.targetNumber
                }
                this.targetTableData.push(filelist)

              }else{
                  this.$message({
                      message: res.data.msg,
                      type: "error"
                  });
              }
          }).
          catch((error) =>{
              console.log(error)
          })
      },
      getSubdata(data){ // 传入子组件的数据
          this.basicsinfo=[];
          this.risklevel = [];
          this.basicsinfo.push({ name:'任务名称', value:data.taskname, icon:'iconrenwumingcheng1' });
          this.basicsinfo.push({ name:'创建时间', value:data.create_time, icon:'iconchuangjianshijian' });
          this.basicsinfo.push({ name:'执行时间', value:data.execute_time, icon:'iconzhihangshijian' });
          this.basicsinfo.push({ name:'执行方式', value:data.priority, icon:'iconzhihangfangshi' });
          this.basicsinfo.push({ name:'报告厂商', value:data.manufacturer, icon:'iconceshifangan' });
          this.basicsinfo.push({ name:'提交者', value:data.user, icon:'icontijiaozhe' });
          // 右边
          this.risklevel.push({ name:'任务风险', value:data.name, icon:'iconfengxiandengji' })
          this.risklevel.push({ name:'任务状态', value:data.task_status, icon:'iconzhuangtai' })
      },
        //统计   
      testObjectives(){
          this.$ajax.get('/smart/reportverify/stats',{
              params:{
                taskId:this.task_id
              }
          }).
          then((res) =>{
              if(res.data.code ==200){
                  let dt = res.data.data;
                //   this.tabs2TargetData.risk_target_number = dt.risk_target_number;
                //   this.tabs2TargetData.no_risk_number = dt.no_risk_number;
                  this.tabs2TargetData.alive_number = dt.aliveTarget; //存活
                //   this.tabs2TargetData.no_alive_number = dt.no_alive_number;
                  this.tabs2TargetData.all_target_number = dt.allTarget; //导入总目标
                  this.tabs2TargetData.high_risk_number = dt.highTarget;
                  this.tabs2TargetData.middle_risk_number = dt.middleTarget;
                  this.tabs2TargetData.low_risk_number = dt.lowTarget;
                  this.tabs2TargetData.safe_risk_number = dt.safeTarget;
                //   -----------

                    this.loop.prove_success_number = dt.verify;
                  this.loop.use_success_number = dt.exp
                  this.loop.prove_faile_number = dt.failed;
                  this.loop.no_prove_number = dt.unVerify;
                //   this.loop.wait_prove_number = dt.wait_prove_number;
                  this.loop.all_vuln_number = dt.allVul;
                  this.loop.high_vuln_number = dt.highVul;
                  this.loop.middle_vuln_number = dt.middleVul;
                  this.loop.low_vuln_number = dt.lowVul;

              }else{
                  this.$message({
                      message: res.data.msg,
                      type: "error"
                  });
              }
          }).
          catch((error) =>{
              console.log(error)
          })
      },
      testTargetData(status){
          this.targetcloading = status;
          this.$ajax.get('/smart/reportverify/targetlist',{
              params:{
                //   alive_number:this.formData.target_type,
                  risk:this.formData.risk_level,
                //   search_fields:this.formData.search,
                  page:this.formData.page_num,
                  size:this.targetPageSize,
                  taskId:this.task_id
              }
          }).
          then((res) =>{
              let dt = res.data;
              if(dt.code == 200){
                  this.tabs2TargetTableData = dt.data.list || [];
                  this.targetTotalpage = dt.data.total;
                  this.targetcloading = false;
              }else{
                  this.$message({
                      message: dt.msg,
                      type: "error"
                  });
              }
          }).
          catch((error) =>{
              console.log(error)
          })
      },
      handleSelectionChangeTarget(val){
          this.multipleSelectionTarget = val;
      },
      targetDelete(scope,single){
          var ids = '';
          if(single == 'yes'){
              ids = scope.row['id']
          }else{
              ids = scope
          }
          this.$ajax({
              url:'/smart/reportverify/targetdelete',
              method:'get',
              params:{
                targetId:ids
              }
          }).
          then((res) =>{
              let dt = res.data;
              if(dt.code == 200){
                  this.$message({
                      message:'目标删除成功',
                      type: 'success'
                  });
                  if(single == 'yes'){
                      scope._self.$refs[`popover_id-${scope.row.id}`].doClose();
                  }else{
                      this.alldelvisible = false;
                  }
                  this.testTargetData(true);
              }
          }).
          catch((error) =>{
              console.log(error);
          })
      },
      btnMultiDeleteTarget(){
          let _arr = [];
          this.multipleSelectionTarget.forEach((item,index) =>{
              _arr.push(item.id)
          })
          let stringArr = _arr.join(',');
          this.targetDelete(stringArr,'no')
      },
      currentchange(t){
          this.formData.page_num = t; 
          this.testTargetData(true);
          this.targetCurrentpage = t;
      },
      handleSizeChangetarget(t){
          this.formData.page_num = 1;
          this.targetPageSize = t;
          this.testTargetData(true);
      },
      handlesearchTarget(){
          this.formData.page_num = 1;
          this.testTargetData(true);
          this.targetCurrentpage = 1;
      },
      handleResetTarget(){
          this.formData.page_num = 1;
          this.formData.search='';
          this.targetPageSize = 10;
          this.targetCurrentpage = 1;
          this.formData.target_type = 0;
          this.formData.risk_level = '';
          this.testTargetData(true);
      },
      mouseentertarget(row, colum, cell, event) {
        this.showOperateButton = true;
        this.rowId = row.id   //赋值行id，便于页面判断  
        },
        mouseleavetarget(row, colum, cell, event) { 
            
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


      // tab3
      openReport(status){
          this.portLoading = status;
          this.$ajax({
              url:'/smart/reportverify/portlist',
              method:'get',
              params: {
                  taskId:this.task_id,
                  search:this.portFormData.search_fields,
                  page:this.portFormData.page_num,
                  size:this.portPageSize
              } 
          }). 
          then((res) =>{
              let dt = res.data;
              if(dt.code ==200){
                  this.portListtableData = dt.data.list;
                  this.portTotalpage = dt.data.total;
                  this.portLoading = false;
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
      handlesearchPort(){
          this.openReport(true);
          this.portFormData.page_num = 1;
          this.portCurrentpage = 1;
      },
      handleResetReport(){
          this.portFormData.target_ids = 1;
          this.portFormData.search_fields = '';
          this.portPageSize = 10;
          this.portFormData.page_num = 1;
          this.openReport(true);
      },
      portCurrentchange(t){
          this.portFormData.page_num = t;
          this.openReport(true);
          this.portCurrentpage = t;
      },
      handleSizeChangePort(t){
          this.portPageSize = t;
          this.openReport(true);
          this.portFormData.page_num = 1;
      },
     
      // aaa
      tabs4TableData(status){
          this.loopLoading = status;
          this.$ajax({
              url:'/smart/reportverify/vullist',
              method:'get',
              params:{
                  taskId:this.task_id,
                //   already_prove:0,
                risk:this.formDataLoophole.vuln_risk,
                status:this.formDataLoophole.vuln_status,
                  search:this.formDataLoophole.search_fields,
                  page:this.formDataLoophole.page_num,
                  size:this.loopPageSize
              }
          }).
          then((res) =>{
              let dt = res.data;
              if(dt.code == 200){
                  this.buglisttableData = dt.data.list || [];
                  this.loopTotal = dt.data.total;
                  this.loopLoading = false;
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
      loopHandlesearchbug(){
          this.tabs4TableData(true);
          this.loop.page_num = 1;
          this.loopCurrentPage = 1;
      },
      loopHandleResetbug(){
          this.formDataLoophole.type = 1;
          this.formDataLoophole.vuln_risk = '';
          
          this.formDataLoophole.vuln_status = '';
          this.formDataLoophole.search_fields = '';
          this.formDataLoophole.page_num = 1;
          this.tabs4TableData(true);
      },
      handleSizeloop(t){
          this.formDataLoophole.page_num = 1;
          this.loopPageSize = t;
          this.tabs4TableData(true);
      },
      handleCurrentloop(t){
          this.formDataLoophole.page_num = t;
          this.tabs4TableData(true);
          this.loopCurrentPage = t;
      },
      btnLoopDel(scope,single){
          var id = '';
          if(single == 'yes'){
              id = scope.row.id;
          }else{
              id = scope;
          }
          this.$ajax({
              url:'/smart/reportverify/vuldelete' ,
              method:'get',
              params:{
                vulId:id
              }
          }).
          then((res) =>{
              let dt = res.data;
              if(dt.code == 200){
                  this.$message({
                      message:'删除成功',
                      type: 'success'
                  });
                  if(single == 'yes'){
                      scope._self.$refs[`popover-${scope.row.id}`].doClose()
                  }else{
                      this.bugalldelvisible = false;
                  }
                  this.tabs4TableData(false);
              }
          }).
          catch((error) =>{
              console.log(error)
          })
      },
      handleSelectionChangeLoop(val){
            this.multipleSelectLoop = val;
      },
      btnMultiDeleteLoop(){
          let arr = []
          this.multipleSelectLoop.forEach((item,index) =>{
              arr.push(item.id);
          })
          let stringArr = arr.join(',');
          this.btnLoopDel(stringArr,'no');
      },
      changeStatus(){
          let arr = []
          this.multipleSelectLoop.forEach((item,index) =>{
              arr.push(item.vuln_id);
          })
          let stringArr = arr.join(',');
          this.$ajax({
              url:'/reportverify/v1/vulns/status/',
              method:'post',
              data:this.qs.stringify({
                  vuln_ids:stringArr,
                  vuln_status:this.statusform.status
              })
          }).
          then((res) =>{
              let dt = res.data;
              if(dt.success){
                  this.upadtestatusdialogVisible = false;
                  this.tabs4TableData(true);
              }
          }).
          catch((error) =>{
              console.log(error)
          })
      },
      btnverification(val){
          this.taskTitle = val.vuln_name;
          this.verify_result.length = 0;
          this.TestdialogVisible = true;
          this.target_id = val.target_id;
          this.vuln_id = val.vuln_id;
          this.pocname = val.pocname;
          this.$ajax.get('/reportverify/v1/vulns/verify/result/',{
              params:{
                  vuln_id:val.vuln_id,
                  target_id:val.target_id
              }
          }).
          then((res) =>{
              let dt = res.data;
              if(dt.success){
                  this.target_result = res.data.results.target_result;
                  for(var i=0; i< res.data.results.verify_result.length;i++){
                      var item =  res.data.results.verify_result[i];
                      if(this.isJSON(item)){ //报文
                          var back = JSON.parse(item);  
                          this.verify_result.push({'str':'success：'+back.success});
                          for( var key in back ){ 
                              if(key != 'request' && key != 'response' && key != 'success'){ 
                                  if(back[key] instanceof Array && back[key].length == 0){
                                      this.verify_result.push({'str':key+':'+'[]'})
                                  }else{
                                      this.verify_result.push({'str':key+':'+back[key]})
                                  }
                              } 
                          }
                          // if(back.request.length != 0 && back.response.length != 0){
                              this.verify_result.push({'back':[back.request[0],back.response[0]]})
                          // }
                          // this.verify_result.push({'back':[back.request[0],back.response[0]]})
                      }else{
                          this.verify_result.push({'str':item})
                      }
                  }
                 // this.tabs4TableData();
              }else{
                  this.$message({
                      message:res.error,
                      type: 'error'
                  });
              }
          }).
          catch((error) =>{
              console.log(error)
          })
      },
      isJSON(str) {
          if (typeof str == 'string') {
              try {
                  var obj=JSON.parse(str);
                  if(typeof obj == 'object' && obj ){
                      return true;
                  }else{
                      return false;
                  }

              } catch(e) { 
                  return false;
              }
          } 
      },
      startTest(){ //开始验证
          this.yzloading = true;
          this.$ajax({
              method:'post',
              url:'/reportverify/v1/vulns/verify/',
              data: {
                  target_id:this.target_id ,
                  vuln_id:this.vuln_id,
                  pocname:this.pocname
              } 
          })
          .then(dt =>{
              let res = dt.data;
              if(res.success){   
                  this.yzloading = false;
                  this.$message({
                      message:'开始验证',
                      type: 'success'
                  });
                  this.fnResult();
              }else{
                  this.$message({
                      message:res.error,
                      type: 'error'
                  });
              }
          }).catch(err=>{})
      },
      closeTest(){
          this.yzloading = false;
          this.TestdialogVisible = false; 
          clearInterval(this.timer);
      },
      fnResult(){
          let that = this;
          that.verify_result.length = 0;
          this.timer = setInterval(function(){
              that.$ajax.get('/reportverify/v1/vulns/verify/result/',{
                  params: { 
                      target_id:that.target_id,
                      vuln_id:that.vuln_id,
                  }
              }).then(dt=>{
                  let res = dt.data;
                  if(res.success){  
                      that.target_result = res.results.target_result;
                      // that.verify_result = res.data.verify_result;
                      for(var i=0; i< res.results.verify_result.length;i++){
                      var item =  res.results.verify_result[i];
                      if(that.isJSON(item)){ //报文
                          var back = JSON.parse(item);  
                          that.verify_result.push({'str':'success：'+back.success});
                          for( var key in back ){ 
                              if(key != 'request' && key != 'response' && key != 'success'){ 
                                  if(back[key] instanceof Array && back[key].length == 0){
                                      that.verify_result.push({'str':key+':'+'[]'})
                                  }else{
                                      that.verify_result.push({'str':key+':'+back[key]})
                                  }
                                  // that.verify_result.push({'str':key+'：url：'+back[key].url})
                              } 
                          }
                          // if(back.request.length != 0 && back.response.length != 0){
                              that.verify_result.push({'back':[back.request[0],back.response[0]]})
                          // }
                          // that.verify_result.push({'back':[back.request[0],back.response[0]]})
                      }else{
                          that.verify_result.push({'str':item})
                      }
                      if(item == '验证结束'){
                          clearInterval(that.timer);
                      }
                  }
                  }else{
                      that.$message({
                          message:res.error,
                          type: 'error'
                      });
                  }
                  
              }).catch(err=>{})
          },5000);
      },
      btnbuginfo(row){   

        this.$ajax({
              url:'/smart/reportverify/vuldetail',
              method:'get',
              params: {
                vulId:row.id
              } 
          }).
          then((res) =>{
              let dt = res.data;
              if(dt.code == 200){
                this.bugdialogVisible = true; 
                this.bugbasicinfo = [];
                var bugjosn= { 
                    vuln_risk: row.risk,
                    vuln_risk_name:row.riskName,
                    cvss: dt.data.cvss, 
                };
                
                this.bugbasicinfo.push(bugjosn); 
                this.bugform.vul_id =  dt.datavulId;
                this.bugform.vuln_name =  dt.data.name;
                this.bugform.detail =  dt.data.desc;
                this.bugform.cve = dt.data.cve;
                this.bugform.cnnvd = dt.data.cnnvd;
                this.bugform.fix_suggest =  dt.data.fix; 
              }else{
                  this.$message({
                      message:dt.msg,
                      type: 'error'
                  });
              }
          }).
          catch((error) =>{
              console.log(error)
          })


        //   this.bugdialogVisible = true; 
        //   this.bugbasicinfo = [];
        //   var bugjosn= { 
        //       vuln_risk: row.vuln_risk,
        //       cvss: row.cvss, 
        //   };
          
        //   this.bugbasicinfo.push(bugjosn); 
        //   this.bugform.vul_id =  row.vul_id;
        //   this.bugform.vuln_name =  row.vuln_name;
        //   this.bugform.detail =  row.desc;
        //   this.bugform.cve = row.cve;
        //   this.bugform.cnnvd = row.cnnvd;
        //   this.bugform.fix_suggest =  row.fix; 

          
      },
      cancalbugdialogVisible(){
          this.bugdialogVisible = false;
          this.is_bugUpdate = false;
          this.responsepack = '';
          // this.updateBugtxt = '编辑';
          this.expands = [];
      },
    //   rowClick(row, event, column){ ////在<table>里，我们已经设置row的key值设置为每行数据id：row-key="id"
    //       Array.prototype.remove = function (val) {
    //           let index = this.indexOf(val);
    //           if (index > -1) {
    //               this.splice(index, 1);
    //           }
    //       }; 
    //       if (this.expands.indexOf(row.id) < 0) {
    //           this.expands.push(row.id);
    //       } else {
    //           this.expands.remove(row.id);
    //       }
    //       if(this.expands[0]){
    //           console.log(this.expands[0]);
    //       }
    //   },
      handleExpandChange(row,expandedRows){
          var that = this;
          that.requestpack  = [];
          if(expandedRows.length > 0){ 
              that.expands = []
              if (row) { 
                  that.expands.push(row.id); // 每次push进去的是每行的ID  
                  that.requestpack = row.request_pack;
                  this.responsepack = '';
              } 
          } 
          else {
              that.expands = []; // 默认不展开
          }
      },
      btnpacket(scope){
          let index = scope.$index;
          this.$ajax({
              method:'post',
              url:'/tools/packet/',
              data: {
                  req_info: window.btoa(this.requestpack),
              } 
          })
          .then(dt=>{ 
              let res = dt.data;
              if(res.success){
                  this.$message({
                      message:'验证报文成功',
                      type: 'success'
                  });
                  this.responsepack = window.atob(res.resp_info)
              }else{
                   this.$message({
                      message:dt.error,
                      type: 'error'
                  });
              }
          }).catch(err=>{})
      },
      // 验证日志
      verificationLogData(status){
          this.logLoading = status;
          this.$ajax({
              url:'/v2/task/report_test/detail_log_list',
              method:'get',
              params: {
                  task_id:this.task_id,
                  page:this.formDataJournal.page_num,
                  size:this.pageSizeLog,
                  search:this.formDataJournal.search
              } 
          }).
          then((res) =>{
              let dt = res.data;
              if(dt.code == 200){
                  this.logtableData = dt.data.list;
                  this.totalLog = dt.data.total;
                  this.logLoading = false;
              }else{
                  this.$message({
                      message:dt.msg,
                      type: 'error'
                  });
              }
          }).
          catch((error) =>{
              console.log(error)
          })
      },
      handleSizeChangelog(t){
          this.pageSizeLog = t;
          this.verificationLogData(true);
      },
      handleCurrentChangelog(t){
          this.currentPageLog = t;
          this.verificationLogData(true);
          this.formDataJournal.page_num = t;
      },
      handlesearchlog(){
          this.verificationLogData(true);
      },
      handleResetlog(){
          this.formDataJournal.search = '';
          this.formDataJournal.page_num = 1;
          this.verificationLogData(true);
      },
      emptyLog(){ //清空日志
          this.$ajax({
              url:'/v2/task/report_test/detail_log_del' ,
              method:'post',
              data:{
                  task_id: this.task_id
              }
          }).then((res) =>{
              let dt = res.data;
              if(dt.code ==200){
                  this.$message({
                      message: '清空日志成功',
                      type: 'success'
                  });
                  this.verificationLogData(true);
              }else{
                  this.$message({
                      message: dt.msg,
                      type: 'error'
                  });
              }
          }).catch((error) =>{
              console.log(error)
          })
      },
      tabsEventFunction(arr){
          this.target_ids = arr[1];
          this.rangs = arr[0];
          this.checkedAll = arr[3]
          if(arr[2] == 'tabs3'){
              this.openReport(true);
          }
          if(arr[2] == 'tabs4'){ 
              this.tabs4TableData();
          }
          if(arr[2] == 'tabs5'){
              this.verificationLogData(true);
          }
      }
  }
})
</script>