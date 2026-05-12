<template>
  	<div class="bigDiv"> 
  		<div class="main-title  ">  
             
            <i class="nav_icon"></i>
            <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
            <label >辅助工具</label>
	  	</div> 
	  	<div class="auxiliarytool context_box_bg"> 
	  		<el-tabs v-model="activeName"   @tab-click="handleClick">
                <el-tab-pane  label="httplog" name="httplog" class="httplog" v-if="false">
                    <div class="basic">
                        <div><strong>介绍：</strong>HTTP/TCP Log功能用于辅助无回显漏洞盲测</div>
                        <div><strong>使用帮助：</strong>默认使用系统的HTTP/TCP Log盲测模块，也可以选择用户自定义的盲测模块，使用自定义HTTP/TCP Log时，系统将不能判断盲测漏洞是否存在</div>
                    </div>
                    <div style="margin-top:20px" class="search-box">
                        <xzbutton 
                        type="primary" 
                        @click="clearall" 
                        :disabled="tableData.length == 0"
                        size="small"  >清空http记录</xzbutton>  
                        <div class="serach-condition">
                            <div class="search-text">
                                <el-input placeholder="请输入关键字"  @keydown.enter.native="handlesearch"   v-model="search_field" class="input-with-select"  size="small" clearable > </el-input>
                                <!-- <el-button type="primary"  size="small" @click="handlesearch">搜索</el-button> -->
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
                                size="small"  >重置</xzbutton>  
                            </div>  
					
                        </div>
                    </div>
                    <div>
                        <el-table
                        :data="tableData" 
                        style="width: 100%">
                        <el-table-column
                            prop="type"
                            label="反连类型"
                            >
                        </el-table-column>
                        <el-table-column
                            prop="host"
                            label="连接来源">
                        </el-table-column>
                        <el-table-column
                            prop="token"
                            label="TOKEN"
                            >
                        </el-table-column>
                            <el-table-column
                            prop="response"
                            label="响应"
                            
                            >
                        </el-table-column>
                        <el-table-column
                            prop="create_time"
                            label="请求时间" 
                            >
                        </el-table-column>
                        </el-table>
                        <el-pagination
                            :page-size="httplogpageSize" 
                            background
                            layout="total,  prev, pager, next, sizes,jumper"
                            :total="totalpage"
                            :current-page="currentpage"
                            @current-change = "currentchange"
                             @size-change="handleSizeChangehttplog" 
                            >
                        </el-pagination>
                    </div>
                </el-tab-pane> 
                <el-tab-pane label="反连服务器" name="reverse">
                    <reverseSever ref="reverseSeverRef"></reverseSever>
                </el-tab-pane>
<!-- DNS log............................................................................................................. -->
                <!-- <el-tab-pane label="DNS log" name="tab1">
                    <dnslog v-if="activeName === 'tab1'"/>
                </el-tab-pane> -->
<!--IP域名绑定............................................................................................................. -->
               <!-- <el-tab-pane label="IP域名绑定" name="tab2">
                    <iplog ref="iplog" />
                </el-tab-pane> -->
<!-- 渗透资源............................................................................................................. -->
                <el-tab-pane v-if="false" label="渗透资源" name="Infiltration">
                    <About></About>
                    <Operation 
                        style="margin-top:20px"
                        one="上传" 
                        two="删除"
                        three="搜索"
                        four="重置"
                        :delList="waitDeleteList"
                        @handleOneClick="handleOperationUpdate"
                        @handleTwoClick="handleOperationDel"
                        @handleThreeClick="handleOperationSearch"
                        @handleFourClick="handleOperationReset"
                     ></Operation> 
                      <div>
                        <el-table :data="tableData" style="width: 100%"  @selection-change="handleSelectionChangeInfiltration" @cell-mouse-enter="mouseenter" @cell-mouse-leave="mouseleave" class="tooltable">
                            <el-table-column type="selection" width="55"> </el-table-column>
                            <el-table-column prop="name" label="名称" width="200"></el-table-column>
                            <el-table-column prop="file_suffix" label="文件格式" width="150"></el-table-column>
                            <el-table-column prop="protocol" label="应用协议"  width="150"></el-table-column>
                            <el-table-column prop="filepath" label="路径地址"></el-table-column>
                            <el-table-column prop="description" label="描述" >
                                <template slot-scope="scope">
                                    <div v-if="showEditFileNameButton && rowId == scope.row.id">
                                        <el-button @click="btnCopyTool(scope.row)" type="text" size="small">复制</el-button>
                                        <el-divider direction="vertical"></el-divider>
                                        <el-button @click="downloadBtn(scope.row)" type="text" size="small">下载</el-button>
                                        <el-divider direction="vertical"></el-divider>
                                          <el-popover
                                            placement="bottom"
                                            width="170"   
                                            :visible-arrow="false"
                                            :ref="`popover_id-${scope.row.id}`"
                                            popper-class="delButton_popper" >
                                            <p class="delText"><i class="el-icon-warning"></i>确定删除吗？</p>
                                            <div style="text-align: right; margin: 0">
                                                <el-button size="mini" class="delCancel" @click="scope._self.$refs[`popover_id-${scope.row.id}`].doClose()">取消</el-button>
                                                <el-button size="mini" type="primary" @click="deleteBtn(scope)">确定</el-button>
                                            </div>  
                                            <!-- <el-button :underline="false"  slot="reference" class="link_danger linkafter" type="text" size="small" >删除</el-button> -->
                                            <el-link :underline="false" class="link_danger linkafter" slot="reference"  size="small" >删除</el-link>
                                        </el-popover> 
                                    </div>
                                    <div v-else >
                                        <span>{{scope.row.description}}</span>
                                    </div>
                                </template>
                            </el-table-column>
                            <!-- <el-table-column label="操作"> -->
                               
                            <!-- </el-table-column> -->
                        </el-table>
                        <el-pagination
                            :page-size="page_size" 
                            background
                            layout="total,  prev, pager, next, sizes,jumper"
                            :total="total"
                            :current-page="page"
                            @size-change="handleSizeChange"
                            @current-change="handleCurrentChange"
                            >
                        </el-pagination>
                    </div>
                </el-tab-pane>
<!--Ping............................................................................................................. -->
                <el-tab-pane label="Ping" name="tab3">
                    <div class="bg_style"  >
                          <Ping v-if="activeName === 'tab3'"/>
                    </div>
                  
                </el-tab-pane>
<!--Traceroute............................................................................................................. -->
                <el-tab-pane label="Traceroute" name="tab4">
                    <div class="bg_style"  >
                        <Traceroute ref = 'traceroute' v-if="activeName === 'tab4'"/>
                    </div>
                </el-tab-pane>

                <!-- <el-tab-pane label="脚本" name="script">
                    <div class="basic" style="margin-bottom:20px">
                        <div><strong>介绍：</strong>脚本模板生成工具，帮助使用者规范脚本</div>
                        <div><strong>使用帮助：</strong>参数填写完成后，将数据带入脚本中作为说明数据，选择需要的输入参数和输出参数后即可下载脚本模板，编辑自己脚本的核心逻辑代码</div>
                    </div>
                    <el-form :model="scriptform" :rules="rules" ref="ruleForm" label-width="100px"  >
                        <div> 
                            <el-form-item label="漏洞名称：" prop="bugname" style="display: inline-block;width:32%">
                                <el-input v-model="scriptform.bugname" maxlength="100" placeholder="漏洞名称" show-word-limit></el-input>
                            </el-form-item>
                            <el-form-item label="漏洞等级：" prop="vullevel" style="display: inline-block ">
                                <el-select v-model="scriptform.vullevel"   style="width:220px">
                                    <el-option label="高" :value="4"></el-option>
                                    <el-option label="中" :value="3"></el-option>
                                    <el-option label="低" :value="2"></el-option>
                                    <el-option label="信息" :value="1"></el-option>
                                </el-select>
                            </el-form-item> 
                            <el-form-item label="官方网站：" style="display: inline-block;width:32%">
                                <el-input v-model="scriptform.appPowerLink" maxlength="100" placeholder="这个组件的官方网站" show-word-limit></el-input>
                            </el-form-item>
                            <el-form-item label="组件名：" style="display: inline-block;width:32%">
                                <el-input v-model="scriptform.appName" maxlength="50" placeholder="组件名" show-word-limit></el-input>
                            </el-form-item>
                            <el-form-item label="漏洞编号：" style="display: inline-block ">
                                <el-input v-model="scriptform.vulID"   style="width:220px" maxlength="20" placeholder="CVE编号或CNVD或CNNVD" show-word-limit></el-input>
                            </el-form-item>
                            <el-form-item label="影响版本：" style="display: inline-block;width:32%">
                                <el-input v-model="scriptform.appVersion" maxlength="20" show-word-limit placeholder="影响版本"></el-input>
                            </el-form-item>
                            <el-form-item label="作者：" prop="author" style="display: inline-block;width:32%">
                                <el-input v-model="scriptform.author" maxlength="50" placeholder="自己的名字" show-word-limit></el-input>
                            </el-form-item>
                            <el-form-item label="漏洞类型：" prop="vulType"  style="display: inline-block "> 
                                <el-select v-model="scriptform.vulType" placeholder="请选择漏洞类型"  style="width:220px" >
                                    <el-option
                                    v-for="(item,index) in inputtypelist"
                                    :key="index"
                                    :label="item"
                                    :value="item" 
                                    >
                                    </el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item label="时间：" style="display: inline-block;width:32%">
                                <el-date-picker
                                    v-model="scriptform.vulDate"
                                    type="datetime"
                                    placeholder="漏洞报出时间" style="width:100%;">
                                </el-date-picker>
                            </el-form-item>
                            <el-form-item label="创建时间：" style="display: inline-block;width:32%">
                                <el-date-picker
                                    v-model="scriptform.createDate"
                                    type="datetime"
                                    placeholder="脚本编写时间" style="width:32%">
                                </el-date-picker>
                            </el-form-item>
                            <el-form-item label="参考链接：" style="display: inline-block;width:calc(32% + 325px)">
                                <el-input  type="textarea" v-model="scriptform.references" placeholder="引用" maxlength="300" show-word-limit></el-input>
                            </el-form-item>
                            <el-form-item label="漏洞描述：" prop="desc">
                                <el-input type="textarea" v-model="scriptform.desc" placeholder="漏洞描述"></el-input>
                            </el-form-item> 
                            <el-form-item label="修复方案：" prop="solution" >
                                <el-input type="textarea" v-model="scriptform.solution" placeholder="修复方案"></el-input>
                            </el-form-item> 
                        </div> 
                        <el-row :gutter="20">
                            <el-col :lg="12" :md="24">
                                <el-row :gutter="20">
                                    <el-col :lg="10" :md="10">
                                        <div class="scriptbox">
                                            <strong>脚本输入参数</strong>
                                            <div style=" height:260px">
                                                <el-select v-model="inputtype" placeholder="请选择类型" style=" width: 100%;"  @change="changeinputtype">
                                                    <el-option
                                                    v-for="(item,index) in inputtypelist"
                                                    :key="index"
                                                    :label="item"
                                                    :value="item"
                                                   
                                                    >
                                                    </el-option>
                                                </el-select>
                                                <el-select v-model="inputtags" filterable placeholder="请选择tags"  style="margin:10px 0; width: 100%;">
                                                    <el-option
                                                    v-for="(item,index) of tagslist"
                                                    :key="index"
                                                    :label="index"
                                                    :value="index" 
                                                    :title="item"
                                                    > 
                                                    </el-option>
                                                </el-select>
                                                <el-select
                                                    v-model="inputparameter"
                                                    multiple
                                                    filterable
                                                    allow-create
                                                    default-first-option
                                                    placeholder="选择参数" style=" width: 100%;">
                                                    <el-option
                                                    v-for="(item,index) in parameter"
                                                    :key="index"
                                                    :label="item "
                                                    :value="item">
                                                    </el-option>
                                                </el-select>
                                            </div>
                                        </div>
                                    </el-col>
                                    <el-col :lg="4"  :md="4">
                                        <div class="Buttonbox" style=" padding-top: 88px;">
                                            <div>
                                                <el-button type="primary" size="small" @click="btnRequired" title="运行该脚本所必须的标签和参数">必选<i class="el-icon-arrow-right"></i></el-button>
                                            </div>
                                            <div>
                                                <el-button type="primary" size="small"  @click="btnOnly" title="运行该脚本所唯一使用的标签和参数">唯一<i class="el-icon-arrow-right"></i></el-button>
                                            </div>
                                            <div>
                                                <el-button type="primary" size="small" @click="btnOptional" title="运行该脚本所使用的可选参数">可选<i class="el-icon-arrow-right"></i></el-button>
                                            </div>
                                            
                                        </div>
                                    </el-col>
                                    <el-col :lg="10"  :md="10" >
                                        <div class="scriptbox">
                                            <strong>已选数据</strong>
                                            <div class="slelecteddata" style="height:260px;overflow-y: auto;">
                                                <div>
                                                    <span>必选 :</span>
                                                    <div >
                                                        <div  v-for="(item,index) in selectedbx" :key="index"> 
                                                            <span>{{item.selectedbxtags}}：</span> 
                                                            <span v-for="(key,i) in item.can" :key="i">
                                                                {{key}} <i  class="el-icon-error" @click="removebx(index,i)"></i>
                                                            </span> 
                                                        </div>
                                                    </div>
                                                </div>
                                                <div>
                                                    <span>唯一：</span>
                                                    <div> 
                                                        <div  v-for="(item,index) in selectedwy" :key="index"> 
                                                            <span>{{item.selectedwytags}}：</span> 
                                                            <span v-for="(key,i) in item.can" :key="i">
                                                                {{key}} <i  class="el-icon-error" @click="removewy(index,i)"></i>
                                                            </span> 
                                                        </div> 
                                                    </div>
                                                </div>
                                                <div>
                                                    <span>可选：</span>
                                                    <div>
                                                        <div  v-for="(item,index) in selectedkx" :key="index"> 
                                                            <span>{{item.selectedkxtags}}：</span> 
                                                            <span v-for="(key,i) in item.can" :key="i">
                                                                {{key}} <i  class="el-icon-error" @click="removekx(index,i)"></i>
                                                            </span> 
                                                        </div> 
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </el-col>
                                </el-row> 
                            </el-col>
                            <el-col :lg="12" :md="24">
                                <el-row :gutter="20">
                                    <el-col :lg="10"  :md="10">
                                        <div class="scriptbox">
                                            <strong>脚本输出参数</strong>
                                            <div style=" height:260px">
                                                <el-select v-model="outputtype" placeholder="请选择类型" style=" width: 100%;"  @change="changeoutputtype">
                                                    <el-option
                                                    v-for="(item,index) in inputtypelist"
                                                    :key="index"
                                                    :label="item"
                                                    :value="item" 
                                                    >
                                                    </el-option>
                                                </el-select>
                                                <el-select v-model="outputtags" filterable placeholder="请选择tags"  style="margin:10px 0; width: 100%;">
                                                    <el-option
                                                    v-for="(item,index) of tagslist"
                                                    :key="index"
                                                    :label="index"
                                                    :value="index" 
                                                    :title="item"
                                                    > 
                                                    </el-option>
                                                </el-select>
                                                <el-select
                                                    v-model="outputparameter"
                                                    multiple
                                                    filterable
                                                    allow-create
                                                    default-first-option
                                                    placeholder="选择参数" style=" width: 100%;">
                                                    <el-option
                                                    v-for="(item,index) in parameter"
                                                    :key="index"
                                                    :label="item"
                                                    :value="item">
                                                    </el-option>
                                                </el-select>
                                            </div>
                                        </div>
                                    </el-col>
                                    <el-col :lg="4"  :md="4">
                                        <div class="Buttonbox" style=" padding-top: 88px;">
                                            <div>
                                                <el-button type="primary" size="small" @click="btnOutput"><i class="el-icon-arrow-right"></i></el-button>
                                            </div>
                                            
                                        </div>
                                    </el-col>
                                    <el-col :lg="10"   :md="10">
                                        <div class="scriptbox">
                                            <strong>已选数据</strong>
                                            <div class="slelecteddata" style="height:260px;overflow-y: auto;">
                                                <div> 
                                                    <div>
                                                        <div  v-for="(item,index) in selectedout" :key="index"> 
                                                            <span>{{item.selectedouttags}}：</span> 
                                                            <span v-for="(key,i) in item.can" :key="i">
                                                                {{key}} <i  class="el-icon-error" @click="removeOutput(index,i)"></i>
                                                            </span> 
                                                        </div> 
                                                    </div>
                                                </div>
                                                
                                            </div>
                                        </div>
                                    </el-col>
                                </el-row> 
                            </el-col> 
                        </el-row> 
                        <div style="margin-top:20px;    text-align: right;">
                            <el-form-item> 
                                <el-button type="primary" @click="preview()">预览</el-button>
                                <el-button @click="Downloadtemplate()">下载模板</el-button>
                            </el-form-item>
                        </div>
                        
                    </el-form>
                </el-tab-pane> -->

                <!-- 空5项目 -->
                <!-- <el-tab-pane label="工具库" name="tab5">
                    <tool ref = 'toolRef' ></tool>
                </el-tab-pane> -->

                <el-tab-pane label="Agent" name="agent" v-if="false">
                    <div class="basic">
                        <div><strong>介绍：</strong>Agent用于内网扫描流量转化，进行内网扫描时，系统流量发送到Agent服务器，Agent服务翻将流量转发给Agent，再通过Agent转发列测试目标</div>
                        <div><strong>使用帮助：</strong>使用Ageni进行内网扫描时，需要开启Agent服务器，并且将透任务的代理配量为Agent服务器，同时在中转机上启动Agent</div>
                    </div>
                    <div style="margin-top: 20px;">
                        <el-row :gutter="20"> 
                            <el-col :span="12">
                                <div class="title_left_line">
                                    <label>Agent服务器</label>
                                    <span>
                                        <el-switch v-model="isagnetopen" class="elSwitch"  >
                                        </el-switch>
                                    </span>
                                </div>
                                <div class="div_block">
                                    <el-form ref="agentserveform" :model="agentserveform" label-width="140px" class="sysform"  >
                                        <el-form-item label="服务器IP：" class="syswarnvalue" prop="cpu_threshold">
                                            <el-input v-model="agentserveform.ip"  
                                                style="width:calc(100% - 190px)">
                                            </el-input>  
                                        </el-form-item>
                                        <el-form-item label="转发端口：" class="syswarnvalue" prop="memory_threshold">
                                            <el-input v-model="agentserveform.port"  
                                                style="width:calc(100% - 190px)">
                                            </el-input>  
                                        </el-form-item>
                                        <el-form-item label="管理端口：" class="syswarnvalue" prop="disk_threshold">
                                            <el-input v-model="agentserveform.port"  
                                                style="width:calc(100% - 190px)">
                                            </el-input>  
                                        </el-form-item> 
                                    </el-form>
                                    <el-button type="primary" class="div_blockbtn" 
                                        @click="btnSaveAgent">保存设置</el-button>
                                </div>
                            </el-col>
                            <el-col :span="12">
                                <div class="title_left_line">
                                    <label>Agent客户端下载</label>
                                    <span>
                                      
                                    </span>
                                </div>
                                <div class="div_block">
                                    <el-form ref="agentserveform" :model="agentserveform" label-width="140px" class="sysform"  >
                                        <el-form-item label="Agent客户端版本：" class="syswarnvalue" prop="cpu_threshold">
                                            <el-input v-model="agentserveform.ip"  
                                                style="width:calc(100% - 190px)">
                                            </el-input>  
                                        </el-form-item>
                                        
                                    </el-form>
                                    <el-button type="primary" class="div_blockbtn" 
                                        @click="btnDownAgent">下载Agent</el-button>
                                </div>
                            </el-col>
                        </el-row>

                    </div> 
                </el-tab-pane>
            </el-tabs> 
            
	  	</div>
	  	<el-dialog title="清空http记录" :visible.sync="dialogFormVisible"   width="640px" :before-close="cancelform" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="submitForm">确定</el-button>
                <el-button size="small" @click="cancelform">取消</el-button>
            </div>
            <div class="dialogtxt" >此操作将清空工具下的所有HTTP记录，确认执行吗？</div>
            
        </el-dialog>

        <el-dialog title="脚本模板预览" :visible.sync="dialogScriptVisible"   width="50%"   :close-on-click-modal="false" :show-close="false">
            <div style="height:350px;overflow-y: auto;">
                <pre>{{scriptcontent}}</pre>
            </div>
            <div slot="footer" class="dialog-footer"> 
                <el-button type="primary" @click="download" >下载</el-button> 
            </div>
        </el-dialog>

        <el-dialog class="upload_dialog shentoudia" title="上传" :visible.sync="createInfiltration"  :before-close="cancelform" width='1184px' :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">  
                <el-button size="small" @click="saveAndSendData">保存</el-button>
                <el-button size="small" @click="handlecancel">取消</el-button>
            </div>
            <Edit
                :exhibitionList="exhibitionList"
                :dictionariesType="dictionariesType"
                :pageType="pageType"
                ref="edit"
            ></Edit>
        </el-dialog>
  	</div>
</template>
<style lang="less" scoped>
    .bg_style{
        height: calc(100% - 39px);
       
    }
    .pingDiv,.TracerouteDiv{
        height: calc(100% - 54px);
        padding: 24px; 
        box-sizing: border-box;
        background-color: #fff;
        box-shadow: 0px 2px 4px 1px rgba(76, 122, 227, 0.11);
        border-radius: 4px;
    }
        .title_left_line{
            font-size: 14px;
            margin-bottom: 16px;
            font-weight: 500;
            border-left: 3px solid #4C7AE3;
            padding-left: 10px;
            height: 14px;
            line-height: 14px;
            color: rgba(72, 72, 102, 0.87);
            >span{
                padding-left: 32px;
                font-size: 13px;
                color: rgba(72, 72, 102, 0.63);
            }
            margin-bottom: 16px;
        }
.div_block{
    height: 298px;
    background-color: #fff;
    box-shadow: 0px 2px 4px 1px rgba(76, 122, 227, 0.11);
    border-radius: 4px;
    margin-bottom: 24px;
    padding: 24px;
    box-sizing: border-box;
    /deep/ .el-form-item__label{
        height:32px;
        line-height: 32px;
    }
    /deep/ .el-form-item__content{
        height:32px;
        line-height: 32px;
    }
    .div_blockbtn{
        margin-left: 0;
        width:82px;
        height:32px;
        line-height: 7px;
        padding-right: 40px;
        text-align: center;
        padding-left: 14px;
        margin-left: 10px;
    }
}
/deep/ .el-tabs__content{
    overflow: inherit;
    height: 100%;
    >div{
        height: 100%;
        >div{
            height: 100%;
        }
    }
}
    /deep/ .tooltable{
        tr{
            height: 57px;
        }
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
        background-color: #FFF;
    }
    /deep/ .el-tabs__nav-wrap::after{
        background: #E8E8F5;
        height: 1px;
    }
    /deep/ .el-tabs__header{
        margin: 0 0 24px;
    }
    .auxiliarytool{
        background: inherit;
        height: calc(100% - 39px);
        box-sizing: border-box;
        // padding: 24px;
       
        /deep/ .el-tabs{
            height: 100%;
            // display: flex;
            // flex-direction: column;
        }
         /deep/ .el-tabs__header{
            margin:0;
            margin-bottom: 15px;
            box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.1200);
            border-radius: 4px;
        }
        /deep/ .el-tabs__content{
            flex: 1;
            // padding: 24px ; 
            // background: #fff;
            // box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.1200);
            border-radius: 4px;
            // margin-top: 24px;
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
    .slelecteddata {
        padding: 10px 5px;
    }
    .slelecteddata > div{
        background: #f0f2f5;
        border-radius: 4px;
        padding:5px 10px;
        font-size: 12px;
        color: #606266;
        margin: 5px 0;
    }
    .slelecteddata > div > span{
        display: inline-block;
        color: #4c7ae3;
        
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
    .bigDiv{
        height: 100%;
    }
</style>
<script>  
import About from "@/components/About.vue";
import Operation from "@/components/Operation";
import Edit from "@/components/edit";
import xzbutton from "@/components/XzButton.vue";
import delbutton from "@/components/DelButton.vue";
import dnslog from "@/pages/tool/dnslog.vue";
import iplog from "@/pages/tool/iplog.vue";
import Ping from "@/pages/tool/Ping.vue";
import Traceroute from "@/pages/tool/Traceroute.vue";
// import tool from '@/pages/tool/toolLibrary.vue'
import reverseSever from "./reverseSever.vue";
import { auxiliarytool } from '@/api/tool.js'
export default ({
    name:'auxiliarytool',
    components:{
        About,
        Operation,
        Edit,
        delbutton,
        xzbutton,
        dnslog,
        iplog,
        Ping,
        Traceroute,
        reverseSever,
        // tool
    },
    data(){  
    	return{
            showEditFileNameButton:false,
            rowId:'',
            multipleSelection: [],
            page:1,
            total:0,
            page_size:10,
            loading:false,  
    		activeName:'reverse',
            token:'',
            httpapi:'',
            search_field:'',
            tableData:[],
            httplogpageSize:10,
            totalpage:0,
            currentpage:1,
            dialogFormVisible:false,
            page_num:1, 
            scriptform:{
                bugname:'',
                vullevel:1,
                appPowerLink:'',
                appName:'',
                appVersion:'',
                author:'',
                vulType:'',
                vulDate:'',
                createDate:'',
                references:'',
                solution:'',
                desc:'',
                required:[],
                choice:[],
                options:[],
                output:[],
            },
            rules:{
                bugname:[
                    { required: true, message: '漏洞名称不能为空', trigger: 'blur' }, 
                ],
                vullevel:[
                     {required: true, message: '请选择漏洞等级', trigger: 'change' }, 
                ],
                author:[
                     {required: true, message: '作者不能为空', trigger: 'blur' }, 
                ],
                vulType:[
                    {required: true, message: '请选择漏洞类型', trigger: 'change' }, 
                ],
                desc:[
                     { required: true, message: '漏洞描述不能为空', trigger: 'blur' }, 
                ],
                solution:[
                     { required: true, message: '修复方案不能为空', trigger: 'blur' }, 
                ]
            },
            alltagslist:[
                {
                    type:'后台登录',
                    tags:{
                        'root_url':'一般的 url 参数',
                        'http_login':'http 登录包标签'
                    }   
                },
                {
                    type:'后台登录多因素认证',
                    tags:{
                        
                    }
                }
            ],
            
            inputtypelist:[ ],//'后台登录','后台登录多因素认证','命令执行类','webshell类','csrf','sql_inj','XSS','任意文件下载','任意文件下载利用','敏感文件'
            tagslist:{},
            inputtype:'',
            inputtags:'',
            inputparameter:'',
            selectedbxtags:'',
            selectedbx:[],
            selectedwytags:'',
            selectedwy:[],
            selectedkxtags:'',
            selectedkx:[],
            outputtype:'',
            outputtags:'',
            outputparameter:'',
            selectedouttags:'',
            selectedout:[],
            parameter:[
                  'url',     
                'data',    
                'headers',     
                'cookie',     
                'username',    
                'password',    
                'user_dict',   
                'passwd_dict',   
                'db_type',   
                'db_title',    
                'db_version',   
                'current_db',  
                'tables',   
                'columns', 
                'sql_inj_user_passwd',  
                'sensitive_table_names',  
                'sensitive_column_names',   
                'db_list',  
                'script_import_content',   
                'script_export_content',   
            ],
            downloadfile:'',
            dialogScriptVisible:false,
            scriptcontent:'',
            exhibitionList:{
                name:"名称",
                namePlaceholder:"名称",
                typeOne:"应用协议",
                file:"资源文件",
                uploadName:"添加资源文件",
                // uploadType:"application/octet-stream",
                uploadType:"",
                remarks:"备注"
            },
            dictionariesType: [
              {label: 'http',value: '1'},
              {label: 'rmi',value: '2'},
              {label: 'ldap',value: '3'}
            ],
            pageType:'auxi',
            createInfiltration:false,
            waitDeleteList: [],
            alldelvisible:false,
            iptableData: [],
            isagnetopen:false,
            agentserveform:{

            }
    	}
    }, 
    created:function(){
        this.$store.state.activefirstMenu="/auxiliarytool"; 
    },
    mounted:function(){   
        // this.getHttplogToken();
        // this.getHttploglist();
        // this.gettags()
    },
    methods:{   
        btnCopyTool(row){ //复制任务
            var strs0=row.filepath;
            var textBox = document.createElement('input');
            textBox.value = strs0;
            document.body.appendChild(textBox);
            textBox.select();                             // 选择对象
            document.execCommand("Copy");                 // 执行浏览器复制命令
            textBox.className = 'textBox';
            textBox.style.display='none';
		},
        handleClick(tab, event){
            if(tab.name == 'tab2'){  
                // this.getHttplogToken(); 
                this.$refs.iplog.getIpData();
            } else if (tab.name == 'Infiltration'){
                //  this.gettags()
                this.getPenetrationResources()
            } else if (tab.name == 'tab5'){
                //  this.gettags()
                this.$refs.toolRef.getData();
            }
        },
        // async gettags(){
        //     //this.tagslist
        //     let params =  {}
        //     const res = await auxiliarytool.gettags(params)
        //     debugger
        //     this.alltagslist = res;
        //     for(var i=0;i<this.alltagslist.length;i++){
        //         this.inputtypelist.push(this.alltagslist[i].type)
        //         for(var o in this.alltagslist[i].tags){  
        //             this.tagslist[o] = this.alltagslist[i].tags[o];
        //         }
        //         console.log(this.tagslist);
        //         // this.tagslist.push(this.alltagslist[i].tags) 
        //     }
        // },
        async getHttplogToken(){ 
            let params =  {}
            const res = await auxiliarytool.getHttplogToken(params)
            if(res.success){ 
                this.token = res.token;
                this.httpapi = res.api;
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            } 
        },
          //httplog列表页
        async  getHttploglist(){   
            let params = {
                search:this.search_field,
                page : this.page_num,
                size: this.httplogpageSize
            }
            const res = await auxiliarytool.getHttploglist(params)
            if(res.code == 200){ 
                this.tableData = res.data.list || [];
                this.totalpage = res.data.total;
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }  
        },
        handleSizeChangehttplog(t){
            this.httplogpageSize = t;
            this.page_num = 1;
            this.getHttploglist();
        },
        // .........................................................................................................................................................
        // 渗透资源
        async getPenetrationResources(name){
            let params =  {
                name:name ? name : '',
                page : this.page,
                page_size : this.page_size
            }
            const res = await auxiliarytool.getPenetrationResources(params)
            if(res.success){ 
                this.tableData = res.result
                    this.tableData.forEach((item, index) => {
                        item.visible = false
                        this.$set(this.tableData, index, item)
                    })
                this.total = res.count;
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }  
        },
        handleSizeChange(page_size){
            this.page_size = page_size;
            this.page = 1;
            this.getPenetrationResources();
        },
        handleCurrentChange(page){
            this.page = page;
            this.getPenetrationResources();
        },
        // 删除路由
        deleteBtn (scope) {
            this.waitDeleteList.push(scope.row)
            scope._self.$refs[`popover_id-${scope.row.id}`].doClose()
            this.handleOperationDel()
            
        },
        downloadBtn (item) { 
            var url="/tools/assists/penetration_resources/download/?name=";
            let baseurl = '';
            if(process.env.NODE_ENV === 'development'){
                baseurl = '/api'
                // baseurl = process.env.VUE_APP_API_URL; 
            }else{
                baseurl = '';
                // baseurl = window.location.host;
            } 
            var href = baseurl+url+item.name; 
            // console.log(href);
            const a = document.createElement('a'); // 创建a标签
            a.setAttribute('target', '_blank');// target属性
            a.setAttribute('download', '');// download属性
            a.setAttribute('href',  href);// href链接
            a.click();// 自执行点击事件
            
        },
        // .................................................................................................................................
        handleReset(){
            this.search_field = '';
            this.page_num = 1;
            this.getHttploglist();
        },
        handlesearch(){
            this.page_num = 1; 
            this.getHttploglist();
            this.currentpage = 1;
        },
    	currentchange(t){
            this.page_num = t; 
            this.getHttploglist();
            this.currentpage = t;
        },
        clearall(){
            this.dialogFormVisible = true;
        },
        cancelform(){
            this.dialogFormVisible = false;
        },
        async submitForm(){ //确认清空 
            const res = await auxiliarytool.submitForm()
            if(res.code == 200){  
                this.$message({
                    message:'清空所有http记录成功',
                    type: 'success'
                });
                this.dialogFormVisible = false;
                this.tableData = [];
                this.currentpage = 1;
                this.totalpage = 0;

            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        }, 
        MergeArray(arr1,arr2){
            var _arr = new Array();
            for(var i=0;i<arr1.length;i++){
                _arr.push(arr1[i]);
            }
            for(var i=0;i<arr2.length;i++){
                var flag = true;
                for(var j=0;j<arr1.length;j++){
                    if(arr2[i]==arr1[j]){
                        flag=false;
                        break;
                    }
                }
                if(flag){
                    _arr.push(arr2[i]);
                }
            }
            return _arr;
        },
        changeinputtype(){
            for (let i = 0; i < this.alltagslist.length; i++) {
                if(this.inputtype ==  this.alltagslist[i].type){
                    this.tagslist = this.alltagslist[i].tags;
                    break;
                }
            } 
            this.inputtags = '';
        },
        btnRequired(){ //必选
            if(this.inputparameter.length == 0  || !this.inputtags || !this.inputtype){
                this.$message({
                    message:'请先选择输入参数，再进行操作 ',
                    type: 'error'
                });
                return;
            }
            if(this.selectedbx.length>0){
                let flag = false;
                for(var i=0;i< this.selectedbx.length;i++){
                    if(this.inputtags == this.selectedbx[i].selectedbxtags){
                        this.selectedbx[i].can = this.MergeArray(this.selectedbx[i].can,this.inputparameter); 
                        flag = true;
                        break;
                    }
                }
                if(!flag){
                    this.selectedbx.push({selectedbxtags: this.inputtags,can:this.inputparameter})
                }
            }else{
                this.selectedbx.push({selectedbxtags: this.inputtags,can:this.inputparameter})
            }
              
            this.inputparameter = [];
        },
        removebx(index,i){ //移除必选

            this.selectedbx[index].can.splice(i,1);
            if(this.selectedbx[index].can.length == 0){
                this.selectedbx.splice(index,1)
            }
 
        },
        btnOnly(){ //唯一
            if(this.inputparameter.length == 0  || !this.inputtags || !this.inputtype){
                this.$message({
                    message:'请先选择输入参数，再进行操作 ',
                    type: 'error'
                });
                return;
            }
            if(this.selectedwy.length>0){
                let flag = false;
                for(var i=0;i< this.selectedwy.length;i++){
                    if(this.inputtags == this.selectedwy[i].selectedwytags){
                        this.selectedwy[i].can = this.MergeArray(this.selectedwy[i].can,this.inputparameter); 
                        flag = true;
                        break;
                    }
                }
                if(!flag){
                    this.selectedwy.push({selectedwytags: this.inputtags,can:this.inputparameter})
                }
            }else{
                this.selectedwy.push({selectedwytags: this.inputtags,can:this.inputparameter})
            }
            this.inputparameter = [];
        },
        removewy(index,i){ //移除唯一

            this.selectedwy[index].can.splice(i,1);
            if(this.selectedwy[index].can.length == 0){
                this.selectedwy.splice(index,1)
            }
            // for (let i = 0; i < this.selectedwy.length; i++) {
            //     if(index == i){
            //         this.selectedwy.splice(i,1);
            //         break;
            //     }
            // }
        },
        btnOptional(){ //可选
            if(this.inputparameter.length == 0 || !this.inputtags  || !this.inputtype){
                this.$message({
                    message:'请先选择输入参数，再进行操作 ',
                    type: 'error'
                });
                return;
            }
            if(this.selectedkx.length>0){
                let flag = false;
                for(var i=0;i< this.selectedkx.length;i++){
                    if(this.inputtags == this.selectedkx[i].selectedkxtags){
                        this.selectedkx[i].can = this.MergeArray(this.selectedkx[i].can,this.inputparameter); 
                        flag = true;
                        break;
                    }
                }
                if(!flag){
                    this.selectedkx.push({selectedkxtags: this.inputtags,can:this.inputparameter})
                }
            }else{
                this.selectedkx.push({selectedkxtags: this.inputtags,can:this.inputparameter})
            }
            this.inputparameter = [];
        },
        removekx(index,i){//移除可选
            this.selectedkx[index].can.splice(i,1);
            if(this.selectedkx[index].can.length == 0){
                this.selectedkx.splice(index,1)
            } 
        },
        changeoutputtype(){
            for (let i = 0; i < this.alltagslist.length; i++) {
                if(this.outputtype ==  this.alltagslist[i].type){
                    this.tagslist = this.alltagslist[i].tags;
                    break;
                }
            }
            this.outputtags = '';
        },
        btnOutput(){  
            if(this.outputparameter.length == 0 || !this.outputtags  || !this.outputtype){
                this.$message({
                    message:'请先选择输出参数，再进行操作 ',
                    type: 'error'
                });
                return;
            }
            if(this.selectedout.length>0){
                let flag = false;
                for(var i=0;i< this.selectedout.length;i++){
                    if(this.outputtags == this.selectedout[i].selectedouttags){
                        this.selectedout[i].can = this.MergeArray(this.selectedout[i].can,this.outputparameter); 
                        flag = true;
                        break;
                    }
                }
                if(!flag){
                    this.selectedout.push({selectedouttags: this.outputtags,can:this.outputparameter})
                }
            }else{
                this.selectedout.push({selectedouttags: this.outputtags,can:this.outputparameter})
            }
            this.outputparameter = [];
        },
        removeOutput(index,i){
            //  for (let i = 0; i < this.selectedout.length; i++) {
            //     if(index == i){
            //         this.selectedout.splice(i,1);
            //         break;
            //     }
            // }
            this.selectedout[index].can.splice(i,1);
            if(this.selectedout[index].can.length == 0){
                this.selectedout.splice(index,1)
            }
        },
        savescript(flag){
            this.downloadfile='';
            for(var i=0;i<this.selectedbx.length;i++){
                var tags = this.selectedbx[i].selectedbxtags;
                var tagsvalue = this.selectedbx[i].can;
                var _j = {};
                _j[tags] = tagsvalue;
                this.scriptform.required.push(_j)
            }
            for(var i=0;i<this.selectedwy.length;i++){
                var tags = this.selectedwy[i].selectedwytags;
                var tagsvalue = this.selectedwy[i].can;
                var _j = {};
                _j[tags] = tagsvalue;
                this.scriptform.choice.push(_j)
            }
            for(var i=0;i<this.selectedkx.length;i++){
                var tags = this.selectedkx[i].selectedkxtags;
                var tagsvalue = this.selectedkx[i].can;
                var _j = {};
                _j[tags] = tagsvalue;
                this.scriptform.options.push(_j)
            }
            for(var i=0;i<this.selectedout.length;i++){
                var tags = this.selectedout[i].selectedouttags;
                var tagsvalue = this.selectedout[i].can;
                var _j = {};
                _j[tags] = tagsvalue;
                this.scriptform.output.push(_j)
            }   
            this.$refs.ruleForm.validate(async (valid) => {
			  	if (valid) {
                    const params = {
                        vulID:this.scriptform.vulID,
                        name:this.scriptform.bugname,
                        vullevel:this.scriptform.vullevel,
                        appPowerLink:this.scriptform.appPowerLink,
                        appName:this.scriptform.appName,
                        appVersion:this.scriptform.appVersion,
                        author:this.scriptform.author,
                        vulType:this.scriptform.vulType,
                        vulDate:this.scriptform.vulDate,
                        createDate:this.scriptform.createDate,
                        references:this.scriptform.references,
                        solution:this.scriptform.solution,
                        desc:this.scriptform.desc,
                        required:JSON.stringify(this.scriptform.required),
                        choice:JSON.stringify(this.scriptform.choice),
                        options:JSON.stringify(this.scriptform.options),
                        output:JSON.stringify(this.scriptform.output),
                    }
                    const res = await auxiliarytool.savescript(params)
                    if(res.success){ 
                        this.downloadfile = 'http://192.168.0.79:8000/'+res.file_path;
                        if(flag == 1){
                            this.dialogScriptVisible = true;
                            this.scriptcontent = res.context;
                        }else if(flag == 2){
                            // window.open(this.downloadfile, '_blank');

                            // 创建隐藏的可下载链接
                            var eleLink = document.createElement('a');
                            eleLink.download = '模板.py';
                            eleLink.style.display = 'none';
                            // 字符内容转变成blob地址
                            var blob = new Blob([res.context]);
                            eleLink.href = URL.createObjectURL(blob);
                            // 触发点击
                            document.body.appendChild(eleLink);
                            eleLink.click();
                            // 然后移除
                            document.body.removeChild(eleLink);
                        } 
                        
                    }else{
                        this.$message({
                            message:res.error,
                            type: 'error'
                        });
                    } 
                }
            }); 
        },
        download(){
            // 创建隐藏的可下载链接
            var eleLink = document.createElement('a');
            eleLink.download = '模板.py';
            eleLink.style.display = 'none';
            // 字符内容转变成blob地址
            var blob = new Blob([ this.scriptcontent ]);
            eleLink.href = URL.createObjectURL(blob);
            // 触发点击
            document.body.appendChild(eleLink);
            eleLink.click();
            // 然后移除
            document.body.removeChild(eleLink);
        },
        preview(){
            this.savescript(1);
        },
        Downloadtemplate(){
            this.savescript(2);
        },
        //  渗透资源列表多选 需要删除的列表 
        //  创建一个 delList数组 传给 Operation组件
        //  参考 src\pages\dictionaryLibrary\dictionary.vue
        handleSelectionChangeInfiltration(list){
            console.log(list,'需要删除的列表');
            this.waitDeleteList = list
        },
        //  上传按钮时间
        handleOperationUpdate(){
            this.createInfiltration = true;
        },
        //  删除按钮事件
        async handleOperationDel(value){
            let _ids = [];
            if (this.waitDeleteList.length === 0) {
                return false
            }
            for (var i = 0; i < this.waitDeleteList.length; i++) {
    			_ids.push(this.waitDeleteList[i].id);
			}

            let params = {
                ids:_ids.join(',')  
            }
            const res = await auxiliarytool.handleOperationDel(params)
            // this.$ajax({
            //     method:'delete',
            //     url:'/tools/assists/penetration_resources/delete/',
            //     data:{
            //        ids:_ids.join(',')  
            //     } 
            // }) 
            // .then((dt) => { 
            //     let res = dt.data;
                if(res.success){
					this.$message({
                        message:'删除任务成功',
                        type: 'success'
					});
                    this.waitDeleteList = []
                    // scope._self.$refs[`popover_id-${scope.$index}`].doClose()
					this.getPenetrationResources()

                }else{
					this.$message({
                        message:res.error,
                        type: 'error'
                    });
				}
            // })
            // .catch((error) => {
            //     console.log(error);
            // })
        },
        //  搜索按钮事件
        handleOperationSearch(value){
            console.log('搜索按钮事件值为：',value)
            if (this.activeName === 'Infiltration') {
                this.page_num = 1
                this.getPenetrationResources(value)
            }
        },
        //  重置按钮事件
        handleOperationReset(value){
            console.log('重置按钮事件')
            this.page_num = 1
            this.getPenetrationResources()
        },
        handleFormDate(data){
            const formDate = new FormData();
            formDate.append('name',data.name);
            formDate.append('protocol',data.dictionariesValue);
            // formDate.append('file',data?.files[0] || '');
            if (data.files.length > 0) {
                formDate.append('file',data.files[0].raw || '');
            }
            console.log(formDate.get('file'))
            // formDate.append('service',data.serviceValue);
            formDate.append('description',data.remarks);
            return formDate;
        },
        //  提交表单
        saveAndSendData(){
            let valiDatas = this.$refs.edit.handleEdit();
            if(valiDatas === null)
                return;
            const data = this.handleFormDate(valiDatas);
            // TODO: 接口调用
            // console.log(validatas);
            let config = {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            }
            this.$ajax({
                url:"/tools/assists/penetration_resources/create/",
                method:"POST",
                data
            })
            .then((data) => {
                let dt = data.data
                if (dt.success) {
                    this.$message.success(dt.msg);
                    this.createInfiltration = false;
                    this.getPenetrationResources();
                    this.$refs.edit.handleClearFiles();
                    this.$refs.edit.form.remarks = ''

                } else {
                    this.$message.error(dt.msg)
                }
            })
            .catch((data)=>{
                console.log(data); //错误信息
            });
            // TODO: 调用完成 重置form表单
            // this.$refs.edit.$refs.editRef.resetFields();
        },
        //  取消移除校验
        handlecancel(){
            this.$refs.edit.handleClearFiles();
            this.$refs.edit.form.remarks = ''
            this.createInfiltration = false;    
        },
         handleSelectionChange(val){
            this.multipleSelection = val;
        },
        mouseenter(row,colum,cell,event){ 
            // this.$refs['popover' + row.user_id].showPopper = true;
            this.showEditFileNameButton = true;
            this.rowId = row.id   //赋值行id，便于页面判断
        },
        mouseleave(row, colum, cell, event){
            let t = this.$refs['popover_id-' + row.id].showPopper;
            if(!t){
                  this.showEditFileNameButton = false;
            this.rowId = "";
            }
          
        },
        // handleCurrentChange(t){
        //     this.search_item.page = t;
        //     this.getData();
        // },
        // handleSizeChange(t){
        //     this.search_item.page = 1;
        //     this.pageSize = t;
        //     this.getData();
        // },
        btnSaveAgent(){

        },
        btnDownAgent(){

        },
    }
})
 
</script>
