<template>
    <el-dialog :title="title"   :visible.sync="comprehensiveTestVisible"  width='1184px' :close-on-click-modal="false" :show-close="false">
        <el-dialog
            width="1184px"
            title="登录检测"
            :visible.sync="logininnerVisible"
            class="subclass"
            :before-close="handleCloseLoginform"
            append-to-body
            :close-on-click-modal="false"
            :show-close="false" >
            <div class="dialog_b_btn">
                <!-- <el-button  size="small" class="btn_login_test" @click="login_validate('login_test')">登录检测</el-button> -->
                <el-button  
                    size="small"  
                    @click="addLoginconfig()" >确定</el-button>
                <el-button  size="small" @click="resetloginForm('loginform')">重置</el-button>
                <el-button  size="small" @click="handleCloseLoginform">关闭</el-button>
            </div>
            <el-form
                ref="loginform"
                :model="loginform"
                label-width="0"
                class="clearfix"
                :rules="loginformrules"
                v-loading.fullscreen.lock="loginloading" >
                <div class="clearfix">
                    <el-form-item
                        label=""
                        prop="ip" class="sub_dialog_label">
                        <label class="dialog_item_label">IP/域名<i class="is-required" style="float:right;">*</i></label>
                        <el-input 
                            :class="loginformip"
                            v-model="loginform.ip"
                            size="small"
                            style="width:320px"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item> 
                    <el-form-item
                        label=""
                        prop="port" class="sub_dialog_label">
                        <label class="dialog_item_label label_mr">端口<i class="is-required" style="float:right;">*</i></label>
                        <el-input
                            :class="loginformport"
                            v-model="loginform.port"
                            size="small"
                            style="width:320px"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item> 
                </div>
                <div class="clearfix">
                    <el-form-item
                        label=""
                        prop="agreement"
                        class="sub_dialog_label" >
                        <label class="dialog_item_label">类型</label>
                        <el-select
                            v-model="loginform.agreement"
                            size="small"
                            placeholder="请选择"
                            style=" width: 320px;" >
                            <el-option
                                v-for="(item,index) in loginagreementlist"
                                :key="index"
                                :label="item"
                                :value="item"
                            ></el-option>
                        </el-select>  
                    </el-form-item>
                    <el-form-item
                        label=""
                        prop="path" class="sub_dialog_label">
                        <label class="dialog_item_label label_mr">路径</label>
                        <el-input
                            v-model="loginform.path"
                            size="small"
                            style="width:320px"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item> 
                </div>
                <div class="clearfix">
                    <el-form-item
                        label=""
                        prop="username" class="sub_dialog_label">
                        <label class="dialog_item_label">用户名</label>
                        <el-input
                            v-model="loginform.username"
                            size="small"
                            style="width:320px"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item> 
                    <el-form-item
                        label=""
                        prop="password" class="sub_dialog_label">
                        <label class="dialog_item_label label_mr">密码</label>
                        <el-input
                            v-model="loginform.password"
                            size="small"
                            style="width:320px"
                            maxlength="50" 
                        ></el-input>
                    </el-form-item> 
                </div>
                <div>
                    <el-form-item
                        label=""
                        prop="cookie" >
                        <label class="dialog_item_label" style="vertical-align: top; margin-top: 10px;">Cookie</label>
                        <el-input
                            class="cookietextarea"
                            v-model="loginform.cookie"
                            type="textarea" 
                            :rows="4" 
                            size="small"
                            style="width:777px" 
                        ></el-input>
                    </el-form-item> 
                </div>  
            </el-form>
        </el-dialog> 
            <el-dialog
                width="1184px"
                title="新增模板"
                :visible.sync="templateVisible" 
                class="dialog_add_template"
                append-to-body >
                <div class="dialog_b_btn"> 
                    <el-button size="small"  @click="submitTemplate()">确定</el-button>
                    <el-button size="small" @click="handleClosetemplateform">取消</el-button>
                </div>
                <el-form
                    ref="templateform"
                    :model="templateform"
                    label-width="0"
                    class="clearfix" >
                    <el-form-item label="" prop="template_name" style="margin-bottom:0">
                        <label class="is-required">*</label>
                        <el-input 
                            :class="templateError"
                            v-model="templateform.template_name" 
                            style="width:calc(100% - 20px);margin-left:8px"></el-input>
                    </el-form-item> 
                </el-form>
            </el-dialog>    
            <div class="dialog_b_btn">
                <el-button size="small" class="btn_Add_Template"  @click="btnUpdateTemplate" v-if="flag == 5">{{updateTxt}}</el-button>
                <el-button size="small" class="btn_Add_Template"  @click="templateVisible=true" v-if="flag == 1">新建模板</el-button>
                <el-button size="small" @click="submithandle" v-if="flag != 5">开始</el-button>
                <el-button size="small" @click="saveTemplate" v-if="flag == 5">保存</el-button>   <!-- 编辑模板，保存-->
				<el-button size="small" @click="cancelTask">取消</el-button>
            </div>
            <div class="dialog_tabs_top" v-if="flag != 5">
                <el-link type="primary" :underline="false"  @click="Totemplate()">模板管理 <i class="el-icon-caret-right"></i></el-link>
            </div>
            <el-form :model="comprehensiveTestform" status-icon  ref="form"  :rules="rules"  style="height:100%" >  
				<el-tabs v-model="activeName" @tab-click="handleClick"> 
					<el-tab-pane label="基础参数" name="tabs1" class="tabs1"> 
                        <el-form-item
                            label=""
                            prop="template"   v-if="flag != 5"> 
                            <label class="dialog_item_label">任务模板</label>
                            <el-select
                                v-model="comprehensiveTestform.template"
                                size="small" 
                                placeholder="请选择"
                                style=" width: 320px;"
                                @change="gettemplateconfig(comprehensiveTestform.template)" >
                                <el-option
                                    v-for="(item,index) in templatelist"
                                    :key="index"
                                    :label="item.template_name"
                                    :value="item.template_id"
                                ></el-option>
                            </el-select>  
                            <el-tooltip placement="right-start">
                                <div slot="content">任务模板给新建任务的功能模块提供默认参数，<br/>模板管理可以管理已有模板。</div> 
                                <i class="iconfont icontishi icontsstyle"></i> 
                            </el-tooltip>
                        </el-form-item> 
                        <el-form-item
                            label=""
                            prop="template"  v-if="flag == 5"> 
                            <label class="dialog_item_label">模板名称</label>
                            <el-input 
                                v-model="comprehensiveTestform.template_name" 
                                size="small"
                                :disabled="!noUpdte"
                                style="width:320px"
                                maxlength="50" value="222"></el-input>
                        </el-form-item> 
                        <div v-if="flag != 5" style="position: relative;">
                            <label class="dialog_item_label" style=" vertical-align: top; margin-top: 10px; width: 103px;">任务目标</label>
                            <el-form-item
                                label=" "  prop="target"  label-width="10px" style="display: inline-block;"> 
                                <el-input
                                    type="textarea" 
                                    :rows="4" 
                                    @input="targetinput"
                                    v-model="comprehensiveTestform.target"
                                    :disabled="flag == 2"
                                    autocomplete="off"
                                    placeholder="任务目标不能为空"
                                    resize="none"
                                    style="width:320px;margin-bottom:10px; margin-top: -32px;"
                                ></el-input>
                                <div v-if="flag != 2">
                                    <el-button type="primary"  size="small" style="vertical-align: top;margin-right:27px" @click="clickupload()">导入</el-button>
                                    <span style="color: rgba(72, 72, 102, 0.32);">只能上传.txt或.xls格式文件</span>
                                    <input
                                        type="file"
                                        class="btnUploadID"
                                        ref="upload"
                                        @change="changeuploaID($event)"
                                        style="display:none"
                                        id="input-file-ID" /> 
                                </div> 
                            </el-form-item>
                            <el-tooltip placement="right-start">
                                <div slot="content">任务目标支持IP、IP段、域名、URL，多个不同目标用“换行”隔开；<br/> 示例：<br/>“192.168.0.127”、“192.168.0.10-127”、<br />“4dogs.cn”、“www.4dogs.cn”、“http://www.4dogs.cn/aqjc/”、<br />“192.168.0.127:8000”</div> 
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; top: 18px;"></i> 
                            </el-tooltip>
                        </div> 
                        <el-form-item
                            label=""
                            prop="taskname" 
                            v-if="flag != 5"> 
                            <label class="dialog_item_label">任务名称<i class="is-required" style="float:right;">*</i></label> 
                            <el-input 
                                v-model="comprehensiveTestform.taskname"
                                :class="tasknameError"
                                size="small"
                                style="width:320px"
                                placeholder="请输入任务名称"
                                maxlength="50" 
                            ></el-input>
                        </el-form-item>  
                        <el-form-item
                            label=""
                            prop="execute_type"  >
                            <label class="dialog_item_label">执行方式</label> 
                            <el-select
                                v-model="comprehensiveTestform.execute_type "
                                size="small"
                                :disabled="!noUpdte"
                                placeholder="请选择"
                                style="width:320px;margin-bottom: 8px;" >
                                <el-option
                                    v-for="item in rantypelist"
                                    :key="item.id"
                                    :label="item.name"
                                    :value="item.id"
                                ></el-option>
                            </el-select>
                            <div class="margin_left_width" v-if="comprehensiveTestform.execute_type == 2">
                                <label class="execute_type_label">计划时间</label>
                                <el-date-picker
                                    :clearable="false"
                                    v-model="comprehensiveTestform.time"
                                    type="datetime"
                                    size="small"
                                    format="yyyy-MM-dd HH:mm:ss"
                                    value-format="yyyy-MM-dd HH:mm:ss"
                                    placeholder="选择日期时间">
                                </el-date-picker>
                            </div>
                            <div class="margin_left_width" v-if="comprehensiveTestform.execute_type == 3" style="margin-bottom: 8px;">
                                <div>
                                    <label class="execute_type_label">计划时间</label>
                                    <el-select 
                                        v-model="comprehensiveTestform.execute_cycletype"
                                        size="small"
                                        placeholder="请选择" 
                                        style="width:126px;margin-right:16px" >
                                        <el-option value="1" label="每月一次"></el-option>
                                        <el-option value="2" label="每周一次"></el-option>
                                    </el-select>
                                    <el-select 
                                        v-if="comprehensiveTestform.execute_cycletype == '1'"
                                        v-model="comprehensiveTestform.execute_cycletype_day"
                                        size="small"
                                        placeholder="请选择"
                                        style="width:126px;margin-right:16px" >
                                        <el-option   
                                            v-for="i in 28" 
                                            :key="i" 
                                            :label="i+'号'"
                                            :value="i"> </el-option> 
                                    </el-select>
                                    <el-select 
                                        v-if="comprehensiveTestform.execute_cycletype == '2'"
                                        v-model="comprehensiveTestform.execute_cycletype_week"
                                        size="small"
                                        placeholder="请选择"
                                        style="width:126px;margin-right:16px" >
                                        <el-option value="1" label="周一"></el-option>
                                        <el-option value="2" label="周二"></el-option>
                                        <el-option value="3" label="周三"></el-option>
                                        <el-option value="4" label="周四"></el-option>
                                        <el-option value="5" label="周五"></el-option>
                                        <el-option value="6" label="周六"></el-option>
                                        <el-option value="7" label="周日"></el-option>
                                    </el-select>
                                    <el-time-select
                                        v-model="comprehensiveTestform.execute_cycletype_starttime"
                                        size="small" 
                                        :clearable="false"
                                        :picker-options="{
                                            start: '00:00',
                                            step: '00:01',
                                            end: '23:59'
                                        }"
                                        style="width:120px!important;"
                                        placeholder="选择时间">
                                    </el-time-select>
                                </div>
                                <div>
                                    <label class="execute_type_label">终止时间</label>
                                    <el-date-picker 
                                        v-model="comprehensiveTestform.endtime"
                                        size="small"
                                        format="yyyy-MM-dd HH:mm:ss"
                                        value-format="yyyy-MM-dd HH:mm:ss"
                                        type="datetime"
                                        placeholder="选择日期时间"
                                        :clearable="false">
                                    </el-date-picker> 
                                </div>
                            </div>
                        </el-form-item> 
					    <el-form-item
                                label=""
                                prop="scheme"  >
								<label class="dialog_item_label">测试方案</label>
                                <el-select
                                    v-model="comprehensiveTestform.scheme"
                                    size="small"
                                    :disabled="!noUpdte"
                                    placeholder="请选择"
                                    style="width:320px" >
                                    <el-option
                                        v-for="item in schemelist"
                                        :key="item.id"
                                        :label="item.name"
                                        :value="item.id"
                                    ></el-option>
                                </el-select>
                                <el-tooltip placement="right-start">
                                    <div slot="content">测试方案用于控制渗透任务的测试范围</div> 
                                    <i class="iconfont icontishi icontsstyle"></i> 
                                </el-tooltip>
                            </el-form-item> 
							<!-- <el-form-item
                                label=""
                                prop="strength" >
								<label class="dialog_item_label">测试强度</label>
                                <el-select
                                    v-model="comprehensiveTestform.strength"
                                    size="small"
                                    :disabled="!noUpdte"
                                    placeholder="请选择"
                                    style="width:320px"  >
                                    <el-option
                                        v-for="item in strengthlist"
                                        :key="item.id"
                                        :label="item.name"
                                        :value="item.id"
                                    ></el-option>
                                </el-select>
                            </el-form-item>    -->
							<!-- <el-form-item
                                label=""
                                prop="rantime" >
								<label class="dialog_item_label">运行时段</label>
                                <el-input
                                    v-model="comprehensiveTestform.rantime"
                                    size="small"
                                    :disabled="!noUpdte"
                                    style="width:320px"
                                    maxlength="50" 
                                    placeholder="时间点之间用“-”隔开，多个时间段用“,”隔开"
                                ></el-input>
                            </el-form-item>  -->
							<el-form-item
                                label=""
                                prop="priority"  >
								<label class="dialog_item_label">任务优先级</label>
                                <el-select
                                    v-model="comprehensiveTestform.priority"
                                    size="small"
                                    :disabled="!noUpdte"
                                    placeholder="请选择"
                                    style="width:320px" >
                                    <el-option
                                        v-for="(item,index) in prioritylist"
                                        :key="index"
                                        :label="item[1]"
                                        :value="item[0]"
                                    ></el-option>
                                </el-select>
                                <el-tooltip placement="right-start">
                                    <div slot="content">多个任务等待执行时，<br />优先级高的任务先执行</div> 
                                    <i class="iconfont icontishi icontsstyle"></i> 
                                </el-tooltip>
                            </el-form-item>  
					</el-tab-pane>
					<el-tab-pane label="高级参数" name="tabs2" class="tabs2" >
                        <div class="high_param_item">
                            <div class="dialog_item_label">端口扫描</div>
                            <el-form-item prop="scan_type" style="margin-bottom:0">
                                <label class="dialog_item_label_m">扫描模式</label> 
                                <div class="dialog_form_value"> 
                                    <el-radio-group v-model="comprehensiveTestform.scan_type" :disabled="!noUpdte" style="margin-bottom:8px"> 
                                        <el-radio :label="1" value="1">智能端口
                                            <el-tooltip placement="top">
                                                <div slot="content">检测到目标存活的基础上，<br />常规端口扫描未发现开放端口，<br />将对剩余端口进行扫描</div> 
                                                <i class="iconfont icontishi icontsstyle" style="margin-left:4px"></i> 
                                            </el-tooltip></el-radio>
                                        <el-radio :label="2" value="2">常用端口</el-radio>
                                        <el-radio :label="3" value="3">指定端口</el-radio>
                                        <el-radio :label="4" value="4">全部端口</el-radio>
                                    </el-radio-group> 
                                    <div style="width: 720px; "> 
                                        <i class="is-required" v-if="comprehensiveTestform.scan_type != 4">*</i> 
                                        <el-input
                                            type="text"
                                            size="small" 
                                            style="width:calc(100% - 10px)"
                                            :disabled="!noUpdte"
                                            v-if="comprehensiveTestform.scan_type == 1"
                                            v-model="comprehensiveTestform.port1"
                                            placeholder="请输入智能端口" ></el-input>
                                        <el-input
                                            type="text"
                                            size="small"
                                            style="width:calc(100% - 10px)"
                                            :disabled="!noUpdte"
                                            v-if="comprehensiveTestform.scan_type == 2"
                                            v-model="comprehensiveTestform.port2"
                                            placeholder="请输入常用端口" ></el-input>
                                        <el-input
                                            type="text"
                                            size="small"
                                            style="width:calc(100% - 10px)"
                                            :disabled="!noUpdte"
                                            v-if="comprehensiveTestform.scan_type == 3"
                                            v-model="comprehensiveTestform.port"
                                            placeholder="请输入指定端口" ></el-input>
                                        <div v-if="comprehensiveTestform.scan_type == 4" style="height: 32px; line-height: 32px;color: rgba(72, 72, 102, 0.64)" >从0~65535进行扫描</div> 
                                    </div> 
                                </div> 
                            </el-form-item>
                            <el-form-item prop="tcp_scan_type" style="margin-bottom:0">
                                <label class="dialog_item_label_m">TCP扫描</label> 
                                <el-radio-group v-model="comprehensiveTestform.tcp_scan_type" :disabled="!noUpdte">
                                    <el-radio :label="1" value="1">TCP-Connect</el-radio>
                                    <el-radio :label="2" value="2">TCP SYN</el-radio>
                                </el-radio-group> 
                            </el-form-item>
                            <el-form-item prop="udp_scan" style="margin-bottom:0" class="udpsan">
                                <label class="dialog_item_label_m">UDP扫描</label> 
                                <el-checkbox v-model="comprehensiveTestform.udp_scan" :disabled="!noUpdte">
                                    <el-tooltip placement="right">
                                        <div slot="content">UDP扫描的端口范围是常用端口</div> 
                                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle;margin-left: 6px;"  ></i> 
                                    </el-tooltip>
                                </el-checkbox> 
                            </el-form-item>
                        </div>
                        <div class="high_param_item">
                            <div>
                                <div class="dialog_item_label">Web爬虫</div>
                                <el-switch v-model="comprehensiveTestform.iswebreptile" :disabled="!noUpdte" class="elSwitch"></el-switch>
                            </div> 
                            <div class="clearfix" style="margin-left: 114px;">
                                <el-form-item
                                    prop="crawler_depth"
                                    label="爬取深度"
                                    style="float:left;width:33%"
                                    label-width="65px" > 
                                    <el-select
                                        v-model="comprehensiveTestform.crawler_depth"
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        size="small"
                                        style="width:74%"
                                        placeholder="请选择" >
                                        <el-option
                                            v-for="(item,index) in crawler_depth"
                                            :key="index"
                                            :label="item"
                                            :value="item"
                                        ></el-option>
                                    </el-select>
                                    <el-tooltip placement="top">
                                        <div slot="content">爬取被测目标的网站目录深度</div> 
                                        <i class="iconfont icontishi icontsstyle"  ></i> 
                                    </el-tooltip>
                                </el-form-item>
                                <!-- <el-form-item
                                    prop="crawler_scope"
                                    label="爬取宽度"
                                    style="float:left;width:33%"
                                    label-width="65px" > 
                                    <el-select
                                        v-model="comprehensiveTestform.crawler_scope"
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        size="small"
                                        style="width:74%"
                                        placeholder="请选择" >
                                        <el-option
                                            v-for="(item,index) in crawler_scope"
                                            :key="index"
                                            :label="item"
                                            :value="item"
                                        ></el-option>
                                    </el-select>
                                    <el-tooltip placement="top">
                                        <div slot="content">爬取被测目标的网站目录深度</div> 
                                        <i class="iconfont icontishi icontsstyle"  ></i> 
                                    </el-tooltip>
                                </el-form-item> -->
                                <el-form-item
                                    label="爬取范围"
                                    style="float:left;width:33%"
                                    label-width="65px" >
                                    <el-select
                                        v-model="comprehensiveTestform.crawl_range"
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        size="small"
                                        placeholder="请选择" >
                                        <el-option
                                            v-for="(item,index) in crawler_rangelist"
                                            :key="index"
                                            :label="item[1]"
                                            :value="item[0]"
                                        ></el-option>
                                    </el-select>
                                </el-form-item>
                                <el-form-item
                                    prop="single_link_timeout"
                                    label="单链接超时"
                                    style="float:left;width:33%"
                                    label-width="78px" > 
                                    <el-select
                                        v-model="comprehensiveTestform.single_link_timeout"
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        size="small"
                                        style="width:72.5%"
                                        placeholder="请选择" >
                                        <el-option
                                            v-for="(item,index) in single_link_timeout"
                                            :key="index"
                                            :label="item"
                                            :value="item"
                                        ></el-option>
                                    </el-select>
                                     <el-tooltip placement="top">
                                        <div slot="content">爬取单个链接等待响应的最长时间</div> 
                                        <i class="iconfont icontishi icontsstyle"  ></i> 
                                    </el-tooltip>
                                </el-form-item>
                            </div>
                            <!-- <div class="clearfix bottomNoe" style="margin-left: 114px;"> 
                                <el-form-item
                                    prop="single_link_timeout"
                                    label="单链接超时"
                                    style="float:left;width:33%"
                                    label-width="78px" > 
                                    <el-select
                                        v-model="comprehensiveTestform.single_link_timeout"
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        size="small"
                                        style="width:72.5%"
                                        placeholder="请选择" >
                                        <el-option
                                            v-for="(item,index) in single_link_timeout"
                                            :key="index"
                                            :label="item"
                                            :value="item"
                                        ></el-option>
                                    </el-select>
                                     <el-tooltip placement="top">
                                        <div slot="content">爬取单个链接等待响应的最长时间</div> 
                                        <i class="iconfont icontishi icontsstyle"  ></i> 
                                    </el-tooltip>
                                </el-form-item>
                                <el-form-item
                                    label="爬取速度"
                                    style="float:left;width:33%"
                                    label-width="65px" >
                                    <el-select
                                        v-model="comprehensiveTestform.crawl_speed"
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        size="small"
                                        placeholder="请选择" >
                                        <el-option
                                            v-for="(item,index) in crawl_speedlist"
                                            :key="index"
                                            :label="item"
                                            :value="item"
                                        ></el-option>
                                    </el-select>
                                </el-form-item>
                                <el-form-item
                                    label="URL去重"
                                    style="float:left;width:33%"
                                    label-width="65px" >
                                    <el-select
                                        v-model="comprehensiveTestform.url_removal"
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        size="small"
                                        placeholder="请选择" >
                                        <el-option  label="是" value="1" ></el-option>
                                        <el-option  label="否" value="2" ></el-option>
                                    </el-select>
                                </el-form-item>
                            </div> -->
                            <!-- <div  style="margin-left: 114px;">
                                <label for="" class="el-form-item__label">分析文件类型</label>
                                <div style="margin: 8px 0;width: 720px;">
                                    <el-input  
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        v-model="comprehensiveTestform.analyze_file_type" 
                                        size="small"  ></el-input> 
                                </div>
                            </div>
                            <div  style="margin-left: 114px;">
                                <label for="" class="el-form-item__label">下载文件类型</label>
                                <div style="margin: 8px 0;width: 720px;">
                                    <el-input  
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        v-model="comprehensiveTestform.download_file_type" 
                                        size="small"  ></el-input> 
                                </div>
                            </div> -->
                            <div  style="margin-left: 114px;">
                                <label for="" class="el-form-item__label">关键字白名单</label>
                                <div style="margin: 8px 0;width: 720px;">
                                    <el-input  
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        v-model="comprehensiveTestform.black_key" 
                                        size="small"  ></el-input> 
                                </div>
                            </div>
                            <div  style="margin-left: 114px;">
                                <label for="" class="el-form-item__label">URL白名单</label>
                                <div style="margin: 8px 0;width: 720px;">
                                    <el-input  
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        v-model="comprehensiveTestform.black_url" 
                                        size="small"  ></el-input> 
                                </div>
                            </div>
                            <div style="margin-left: 114px;">
                                <el-radio-group
                                    v-model="comprehensiveTestform.tcpport"
                                    :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                    style="margin-bottom:10px" >
                                    <el-radio :label="2">默认 header</el-radio>
                                    <el-radio :label="3">自定义 header</el-radio>
                                </el-radio-group>
                                <div  style="margin-bottom: 8px;width: 720px;">
                                    <el-input 
                                        v-if="comprehensiveTestform.tcpport == 2"
                                        v-model="comprehensiveTestform.crawler_header"
                                        disabled
                                        size="small"  ></el-input> 
                                    <el-input
                                        :class="headerError"
                                         v-if="comprehensiveTestform.tcpport == 3"
                                        v-model="comprehensiveTestform.crawler_header1"
                                        :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"
                                        size="small"  ></el-input> 
                                </div> 
                            </div>
                             <!-- <div class="checkstyle" style="margin-left: 114px;"  > 
                                <el-checkbox v-model="comprehensiveTestform.explain_flash"  
                                    :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"> 
                                    解释flash
                                </el-checkbox> 
                                <el-checkbox v-model="comprehensiveTestform.is_trigger_js"  
                                    :disabled="!comprehensiveTestform.iswebreptile || !noUpdte"> 
                                    javascript是否触发
                                </el-checkbox> 
                            </div> -->
                        </div>
                        <div class="high_param_item">
                            <div class="dialog_item_label">Web路径猜测</div>
                            <el-switch v-model="comprehensiveTestform.iswebroute" :disabled="!noUpdte" class="elSwitch"></el-switch>  
                        </div>
                        <!-- <div class="high_param_item">
                            <div>
                                <div class="dialog_item_label">口令猜测</div>
                                <el-switch v-model="comprehensiveTestform.ispassword" :disabled="!noUpdte" class="elSwitch"></el-switch>  
                            </div> 
                            <div class="clearfix" style="margin-left: 114px;">
                                <el-form-item
                                    label="用户字典"
                                    label-width="65px"
                                    style="float:left;width:33%" >
                                    <el-select
                                        v-model="comprehensiveTestform.user_dict"
                                        :disabled="!comprehensiveTestform.ispassword || !noUpdte"
                                        size="small"
                                        placeholder="请选择"
                                        style="width:85%" >
                                        <el-option
                                            v-for="item in userdictlist"
                                            :key="item.id"
                                            :label="item.name"
                                            :value="item.id"
                                        ></el-option>
                                    </el-select>
                                </el-form-item>
                                <el-form-item
                                    prop="pass_dict"
                                    label="密码字典"
                                    label-width="65px"
                                    style="float:left;width:33%" >
                                    <el-select
                                        v-model="comprehensiveTestform.pass_dict"
                                        :disabled="!comprehensiveTestform.ispassword || !noUpdte"
                                        size="small"
                                        placeholder="请选择"
                                        style="width:85%" >
                                        <el-option
                                            v-for="item in pwddictlist"
                                            :key="item.id"
                                            :label="item.name"
                                            :value="item.id"
                                        ></el-option>
                                    </el-select>
                                </el-form-item>
                            </div>
                        </div> -->
                        <div class="high_param_item"  v-if="flag != 5">
                            <div class="clearfix">
                                <label class="dialog_item_label">登录测试</label>
                                <!-- <el-switch v-model="comprehensiveTestform.islogin " :disabled="!noUpdte" class="elSwitch"></el-switch>   -->
                                <div style="float:right;margin-bottom:16px" >
                                    <el-button size="small" type="primary" 
                                        :disabled="!comprehensiveTestform.islogin" 
                                        style="margin-right: 10px;"
                                        @click="AddLogin" >新增登录 </el-button>
                                    <!-- <el-button size="small" type="primary" plain  >批量测试 </el-button> -->
                                </div>
                                
                                <!-- <el-link :underline="false" 
                                    :disabled="!comprehensiveTestform.islogin" 
                                    type="primary" 
                                    style="float:right;margin-right: 10px;" 
                                    @click.stop="logininnerVisible = true"  >
                                    <i class="el-icon-plus"></i>新增登录
                                </el-link> -->
                            </div> 
                            <div class="hiddenbox "  >
                                <el-table
                                    :data="comprehensiveTestform.login_conf"
                                    size='small'
                                    style="width: 100%">
                                    <el-table-column
                                        prop="ip"
                                        label="IP地址" >
                                    </el-table-column>
                                    <el-table-column
                                        prop="scheme"
                                        :render-header="icons"
                                        label="类型" >
                                    </el-table-column>
                                    <el-table-column
                                        prop="port"
                                        label="登录端口">
                                    </el-table-column>
                                    <el-table-column 
                                        label="操作" >
                                        <template slot-scope="scope"> 
                                            <!-- <el-link :underline="false"  class="link_primary"  > 登录测试
                                            </el-link>  -->
                                            <el-link :underline="false"  class="link_primary" 
                                                @click.native.prevent="updatelogin(scope.$index,comprehensiveTestform.login_conf)"> 编辑
                                            </el-link> 
                                            <el-link :underline="false"  class="link_danger"
                                                @click.native.prevent="dellogin(scope.$index,comprehensiveTestform.login_conf)" > 删除
                                            </el-link> 
                                        </template>
                                        </el-table-column>
                                    </el-table>
                            </div>
                        </div>
                        <div  class="high_param_item">
                            <div>
                                <label for="" class="dialog_item_label">远程控制</label>
                                <el-switch v-model="comprehensiveTestform.isRemotecontrol" :disabled="!noUpdte" class="elSwitch"></el-switch>
                                <el-tooltip placement="top">
                                    <div slot="content">上传远控或反弹shell</div> 
                                    <i class="iconfont icontishi icontsstyle" style="vertical-align: text-top;"  ></i> 
                                </el-tooltip>
                            </div>
                        </div> 
					</el-tab-pane>
				</el-tabs> 
            </el-form> 
        </el-dialog> 
</template>
<style  scoped  lang="less">
    /deep/  .el-dialog  .el-table td:not(.el-table-column--selection):first-child .cell, 
    /deep/  .el-dialog  .el-table th:not(.el-table-column--selection):first-child .cell{
        padding-left: 32px !important;
    }
    .bottomNoe /deep/ .el-form-item__content{
        margin-bottom: 0 !important;
    }
    .checkstyle /deep/ .is-checked .el-checkbox__label{
        font-weight: 500;
        color: #4C7AE3;
        font-size: 13px;
    }
    .checkstyle /deep/ .el-checkbox__label{
        font-weight: 500;
        color: rgba(72, 72, 102, 0.64);
        font-size: 13px;
    }
    .is-required{
        margin-right:4px;
        color:#F56C6C;
        font-size: 12px;
        // vertical-align: sub;
    }
    .is-error /deep/ input{
        border-color: rgb(245, 108, 108) !important;
    }
    .cookietextarea textarea{
        resize: none;
    }
    /deep/ .el-progress__text{
		font-size: 13px !important;
		color:rgba(72,72,102,0.64);
    }
    /deep/ .el-tabs{
        height: 100%;
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
    @media (max-width: 1440px) {
        
        .subclass /deep/ .el-dialog{
            height: calc(100% - 96px);
        }
        /deep/ .el-dialog{
            height: calc(100% - 96px);
        }
    }
    @media  (min-width: 1440px) {
        
        .subclass /deep/ .el-dialog{
            height: calc(100% - 176px);
        }
        /deep/ .el-dialog{
            height: calc(100% - 176px);
        }
    }
    
    /deep/ .el-tabs__header{
        margin: 0 ;
    }
    /deep/ .el-dialog .el-dialog__body{
        background: #F7F7FB;
        height: calc(100% - 62px);
    }
    /deep/ .el-tabs__content{
        padding:  24px;
        height: calc(100% - 96px);
        overflow: auto;
    }
    .tabs1 /deepp/ .el-form-item__content{
        line-height:32;
    }
    .tabs2 /deep/ .el-form-item{
        margin-bottom: 0;
    }
    .tabs2 /deep/ .el-form-item__content{
        line-height: 0;
        margin-bottom: 16px;
        
    }
    .high_param_item /deep/ .el-form-item__label{
        line-height: 33px !important;
    }
    /deep/ .el-form-item__label{
        line-height: 44px;
        font-weight:500;
        color:rgba(72,72,102,0.64);
        font-size: 13px;
    }
    /deep/ .el-radio{
        color:rgba(72,72,102,0.64);
    }
    .dialog_add_template /deep/ .el-dialog__body{
        padding: 72px 152px;
    }
    .dialog_add_template /deep/ .el-dialog{
            height: auto;
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
    .dialog_item_label_m{
        display: inline-block;
        font-size: 14px;
        font-weight:500;
        color:rgba(72,72,102,0.87);
        padding-left: 11px;
        width: 104px;
        // height: 18px;
        line-height: 16px;
    }
    .dialog_form_value{
        display: inline-block; 
        vertical-align: top;
    }
    .high_param_item{
        margin-bottom: 16px;
        .dialog_item_label{
            margin-bottom: 16px;
            width: 114px;
        }
    }
    .dialog_b_btn{
        position: absolute;
        top: 15px;
        right: 24px;
        font-size: 0;
        button{
            color: #4C7AE3;
            font-size: 14px;
        }
        .btn_login_test,
        .btn_Add_Template{
            background:rgba(255,255,255,0.26);
            border-radius:2px;
            border:1px solid rgba(255,255,255,0.5);
            color:#fff; 
        }
    }
    .dialog_tabs_top{
        position: absolute;
        top: 75px;
        right: 24px;
        font-size: 14px;
        color: #4C7AE3;
        z-index: 2;
    }
    .subclass {
        .dialog_item_label{
            width: 100px;
            font-size: 14px;
            display: inline-block;
        } 
        .label_mr{
            width:80px;
        }
    }  
    .sub_dialog_label{
        float: left;
        width: 480px;
    }
    .subclass /deep/ .el-dialog__body{
        height: 100%;
    }
    .hiddenbox{
        background: #fff;
        padding: 24px;
        border: 1px solid rgba(232,232,245,1);
    }
    .uploadfile{
        display: inline-block;
        height: 32px;
        line-height: 32px; 
        i{
            display: inline-block;
            height: 32px;
            line-height: 17px;
            vertical-align: middle;
        }
        span{
            display: inline-block;
            border-bottom:2px solid #4C7AE3;
            height:22px;
            line-height:18px;
            color:rgba(72,72,102,0.64);
            vertical-align: top;
        }
    }
    .margin_left_width{
        margin-left: 114px;
    }
    .execute_type_label{
        color:rgba(72,72,102,0.64);
        font-size: 13px;
        margin-right: 16px;
    }
</style>
<style>
    .el-dialog__header{
        padding: 15px 0;
        background: #4C7AE3; 
    }
    .el-dialog__title{
        display: inline-block;
        font-size: 14px;
        /* width: 104px; */
        padding: 0 24px;
        text-align: center;
        border-left: 2px solid #fff;
        height: 32px;
        line-height: 32px;
        background:rgba(255,255,255,0.12);
        color: #fff; 
    }
   
</style>
<script> 
export default {
    props: {
        value: {}, // 注意此处获取的value对应的就是组件标签中的v-model
        title: {},
        type: {},
        dialogtype: {}, //1:综合检测、2:计算机检测，3:WiFi检测，4:云平台，5:视频监控检测，6:办公自动化
        isInfo: {}, //是否详情
        isAddtarget: {}, //新增检测目标
        flag: {}, //1新建，2重测,3复制任务,5:编辑模板
        taskid: {}, //任务id
        targetid: {}, //目标id
        templateid: {}, //模板id
        templatename: {} //模板名称
    },
    data() {
        let showDate =new Date();
            let seperator ='-';
            let year = showDate.getFullYear();
            let month = showDate.getMonth() + 1;
            let day =showDate.getDate();
            let hour =showDate.getHours();
            let mins  = showDate.getMinutes();
            let sec = showDate.getSeconds();
            var strDate = showDate.getDate();
            if (month >= 1 && month <= 9) {
                month = "0" + month;
            }
            if (strDate >= 0 && strDate <= 9) {
                strDate = "0" + strDate;
            }
        let currentdate = year + seperator + month + seperator + strDate +' '+ hour +':'+mins+':'+sec;
        return {
            sendVal: false,
            activeName:'tabs1',
            noUpdte:false,
            comprehensiveTestVisible: false,
            comprehensiveTestform: {
                template_name:'',
                template: "", 
                taskname: "",
                time:currentdate,
                rantime:currentdate,
                execute_type :1,
                execute_cycletype:'1',
                execute_cycletype_day:1,
                execute_cycletype_week:'1',
                endtime:currentdate,
                scheme: "", 
                target: "",
                priority: "",
                isport: false,
                port: "",
                port1: "",
                port2: "",
                scan_type: 1,
                tcp_scan_type: "",
                iswebreptile: false,
                iswebroute: false,
                ispassword: false,
                strength: "", 
                login_conf: [],
                user_pass_choice: 1,
                user_dict: "",
                Retesttarget: "",
                guess_dict: "",
                password_guess_open: true,
                imitate_ap_open: true,
                imitate_ap_time: "",
                code_embedding_open: true,
                web_shell_open: false,
                bounce_shell_open: false,
                trojan_horse_open: false,
                isRemotecontrol: false,
                udp_scan:false,
                crawler_depth:'',
                crawl_range:'',
                crawler_scope:'1',
                crawl_speed:'5',
                url_removal:'1',
                single_link_timeout:'',
                black_key:'',
                black_url:'',
                analyze_file_type:'',
                download_file_type:'',
                crawler_header:'',
                crawler_header1:'',
                explain_flash:true,
                is_trigger_js:true,
                tcpport:false,
                islogin:false,
                execute_cycletype_starttime:'00:00'
            },
            templateError:'',
            tasknameError:'',
            headerError:'',
            templateVisible:false,
            templateform: {
                template_name: ""
            },
            isport: false,
            iswebreptile: false,
            iswebroute: false,
            ispassword: false, 
            islogin: false, 
            isRemotecontrol: false, 
            loginagreementlist:[],
            templatelist:[],
            rules:{ 
				target:[
					{ required: true, message: '', trigger: 'blur' },  
				],  
            },
            loginformrules:{},
            loginform: {
                ip:'',
                port:'',
                agreement:'',
                username:'',
                password:'',
                path:'',
                cookie:'',
            },
            loginloading:false,
            loginformip:'',
            loginformport:'',
            logininnerVisible:false,
            pwddictlist:[],
            userdictlist:[],
            isLoginSuccess:false,
            crawler_rangelist: [
                [1, " 爬取全域名"],
                [2, " 爬取当前url与子目录"],
                [3, "爬取目标url"]
            ],
            prioritylist: [
                [1, "低"],
                [2, "中"],
                [3, "高"]
            ],
            strengthlist: [
                { id: 1, name: "绿色可控渗透" },
                { id: 2, name: "全量渗透威胁" }
            ],
            rantypelist:[
                { id : 1, name:'立即执行'},
                { id: 2, name: '定时执行' },
                { id: 3, name: '周期执行' },
            ],
            crawler_depth:['1','2','3','4','5','6','7','8','9','10','12','15','20',],
            crawler_scope:['1','2','3','4','5',],
            single_link_timeout:['1','2','3','5','10','20','30','60','120'],
            crawl_speedlist:['1','2','3','5','10','15','20'],
            schemelist:[],
            updateTxt:'编辑模板',
            updateLogin:false,
            updateIndex:-1,
        }
    },
    created: function() { 
    },
    watch: {
        value(newVal, oldVal) {
            // 监测value的变化，并赋值。
            this.comprehensiveTestVisible = newVal; 
            this.activeName = "tabs1";
            this.noUpdte = true;
            if (this.comprehensiveTestVisible) {
                this.comprehensiveTestform.template = "";
                this.gettemplete();
                this.getAllscheme();
                this.getTaskschemes();
                // this.getdictlist(); 
                
                if (this.flag == "2") {//重测
                    this.getretestconfig();
                }
                if(this.flag == "3" || this.flag == "6"){ //复制任务
                    this.getretestconfig();
                }
                if (this.flag == "5") { //编辑模板
                    this.noUpdte = false; 
                    this.updateTxt = '编辑模板';
                    this.gettemplateconfig(this.templateid);
                    this.comprehensiveTestform.template_name = this.title;
                }

            }else{
                this.$refs["form"].resetFields();
            }
            console.log(this.flag);
        },
        comprehensiveTestVisible(val) {
            this.$emit("input", val); // 此处监测showMask目的为关闭弹窗时，重新更换value值，注意emit的事件一定要为input。
        }
    },
    mounted: function() {
        this.comprehensiveTestVisible = this.value; // 在生命周期中，把获取的value值获取给comprehensiveTestVisible
    },
    methods: {
        handleClick(){

        },
        icons(h,{column}){
            const inReview = '登录目标的类型，如XX数据库、XX服务' 
			return h('div', {
					style: { 
							'padding-left':' 0 !important',
							'height': '23px',
							'line-height': '23px',
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
                                'width':'250px',
                                whiteSpace: 'normal', 
                            }
                        }, inReview), 
                        h('i', {
                            class: 'iconfont icontishi',
                            style: 'color:rgba(72,72,102,0.32);margin-left:5px;vertical-align: bottom;'
                        })
                    ],)
                ],
    　　　　 )
        },
        submithandle() { 
            //保存
            if (this.flag == "1") {
                //新建
                this.saveCreateTask();
            } else if (this.flag == "2") {
                //重测
                this.saveRetest();
            } 
            else if (this.flag == "3") {
                //复制
                this.saveCreateTask();
            } 
            else if (this.flag == "6") {
                //测试验证
                this.saveCreateTask();
            } 
            // else if (this.flag == "4") {
            //     //重测目标
            //     this.saveRetestTarget();
            // else if (this.flag == "5") {
            //     //编辑模板
            // }
        },
        saveCreateTask(){
            
            //新建保存
            if(!this.comprehensiveTestform.taskname){
                this.tasknameError='is-error';
                this.$message({
                    message:'任务名称不能为空',
                    type: 'error'
                });
                return ;
            }else{
                this.tasknameError='';
            }
            if(this.comprehensiveTestform.tcpport == 3){ //3 自定义header
                if(!this.comprehensiveTestform.crawler_header1){
                    this.headerError='is-error';
                    this.$message({
                        message:'自定义header不能为空',
                        type: 'error'
                    });
                    return ;
                }else{
                    this.headerError='';
                }
            }else{  
                this.headerError='';
            }
            var json = this.saveparameter(); 
            if(this.flag == "6"){
                json.relation_task_id = this.taskid;
            }
            this.$ajax({
                method:'post',
                url:'/task/task/',
                data:this.qs.stringify(json)
            }).then(dt => {
                let res = dt.data;
                if (res.success == true) {
                    this.dangerBtn();
                } else {
                    this.$message({
                        message: res.error,
                        type: "error"
                    });
                }
            })
            .catch(err => {});
        },
        saveRetest() {  //重测保存 /task/task/restart/
            if(!this.comprehensiveTestform.taskname){
                this.tasknameError='is-error';
                this.$message({
                    message:'任务名称不能为空',
                    type: 'error'
                });
                return ;
            }else{
                this.tasknameError='';
            }
            if(this.comprehensiveTestform.tcpport == 3){ //3 自定义header
                if(!this.comprehensiveTestform.crawler_header1){
                    this.headerError='is-error';
                    this.$message({
                        message:'自定义header不能为空',
                        type: 'error'
                    });
                    return ;
                }else{
                    this.headerError='';
                }
            }else{  
                this.headerError='';
            }
            
            var json = this.saveparameter();
            json.task_id = this.taskid;
            this.$ajax({
                method:'post',
                url:'/task/task/restart/',
                data:this.qs.stringify(json)
            }).then(dt => {
                let res = dt.data;
                if (res.success == true) {
                    this.dangerBtn();
                } else {
                    this.$message({
                        message: res.error,
                        type: "error"
                    });
                }
            })
            .catch(err => {});
        },
        getretestconfig(){ //获得重测任务的配置 /task/task/config/
            this.$ajax.get('/task/task/config/',{
                params: {
                    task_id:this.taskid,
                }
            }).then(dt=>{
                let res = dt.data; 
                this.comprehensiveTestform.strength = res.strength;
                this.comprehensiveTestform.priority = res.priority;
                this.comprehensiveTestform.taskname = res.task_name;
                this.comprehensiveTestform.scheme = res.check_scheme;
                this.comprehensiveTestform.template = res.task_template; 
                this.comprehensiveTestform.target = res.check_target;  
                this.comprehensiveTestform.execute_type = parseInt(res.execute_type) ;
                if( this.comprehensiveTestform.execute_type  == 2){
                    this.comprehensiveTestform.time = res.plan_time;
                }
                if(this.comprehensiveTestform.execute_type  == 3){
                    let arr = res.plan_time.split(',');
                    this.comprehensiveTestform.execute_cycletype = arr[0];
                    if(arr[1] == 1 ){
                        this.comprehensiveTestform.execute_cycletype_day = arr[1];
                    }
                    if(arr[1] == 2 ){
                        this.comprehensiveTestform.execute_cycletype_week = arr[1];
                    }
                    this.comprehensiveTestform.execute_cycletype_starttime = arr[2];
                    this.comprehensiveTestform.endtime = res.finish_time;
                }
                this.comprehensiveTestform.islogin = res.login_check.is_open;
                this.comprehensiveTestform.login_conf =  res.login_check.login_conf;
                this.comprehensiveTestform.isport = res.port_scan.is_open;
                this.comprehensiveTestform.scan_type = res.port_scan.port_scan_type;
                 if(this.comprehensiveTestform.scan_type == 1) {
                    this.comprehensiveTestform.port1 = res.port_scan.scan_port;
                    this.comprehensiveTestform.port2 = res.port_scan.scan_port;
                }
                if(this.comprehensiveTestform.scan_type == 2) {
                    this.comprehensiveTestform.port1 = res.port_scan.scan_port;
                    this.comprehensiveTestform.port2 = res.port_scan.scan_port;
                }
                if(this.comprehensiveTestform.scan_type == 3) {
                    this.comprehensiveTestform.port = res.port_scan.scan_port;
                }

               
                this.comprehensiveTestform.tcp_scan_type =res.port_scan.tcp_scan_type;
                this.comprehensiveTestform.udp_scan = res.port_scan.udp_scan;
                this.comprehensiveTestform.isescape = res.vm_escape.is_open; 

                this.comprehensiveTestform.iswebreptile = res.web_crawler.is_open; 
                this.comprehensiveTestform.crawl_range = res.web_crawler.crawl_range;
                this.comprehensiveTestform.crawler_depth = res.web_crawler.crawler_depth+'';

                this.comprehensiveTestform.crawler_scope = res.web_crawler.crawl_scope+'';
                this.comprehensiveTestform.crawl_speed = res.web_crawler.crawl_speed; 
                res.web_crawler.url_removal ? this.comprehensiveTestform.url_removal  = '1' :this.comprehensiveTestform.url_removal  = '0'
                this.comprehensiveTestform.black_key = res.web_crawler.black_key;
                this.comprehensiveTestform.black_url = res.web_crawler.black_url;
                this.comprehensiveTestform.analyze_file_type = res.web_crawler.analyze_file_type;
                this.comprehensiveTestform.download_file_type = res.web_crawler.download_file_type;
                this.comprehensiveTestform.explain_flash = res.web_crawler.explain_flash;
                this.comprehensiveTestform.is_trigger_js = res.web_crawler.is_trigger_js;
                
                this.comprehensiveTestform.single_link_timeout = res.web_crawler.single_link_timeout+'';
                this.comprehensiveTestform.iswebroute = res.web_crawler.web_path_scan;
                // this.comprehensiveTestform.tcpport = res.web_crawler.self_defined_header ;
                if(!res.web_crawler.self_defined_header ){
                     this.comprehensiveTestform.tcpport = 2;
                    this.comprehensiveTestform.crawler_header = res.web_crawler.crawler_header;
                }else{
                    this.comprehensiveTestform.tcpport = 3;
                    this.comprehensiveTestform.crawler_header1 = res.web_crawler.headers;
                }
                this.comprehensiveTestform.ispassword = res.word_guess.is_open;
                this.comprehensiveTestform.guess_number = res.word_guess.guess_number;
                this.comprehensiveTestform.guess_rate = res.word_guess.guess_rate;
                this.comprehensiveTestform.guess_timeout = res.word_guess.guess_timeout; 
                
                this.comprehensiveTestform.isRemotecontrol = res.remote_control.is_open;
               
                
            }).catch(err=>{})
        },
        dangerBtn() {
            this.$emit("danger");
            this.closeMask();
        },
        closeMask() {
            this.comprehensiveTestVisible = false;
            this.activeName = "tabs1";
        },
        saveparameter(){ //新增任务、重测任务 ----保存参数
            var json = {};  
             json.check_target_type = this.comprehensiveTestform.targetType;
            json.check_target = this.comprehensiveTestform.target.split('\n').join(',');  
            json.task_template = this.comprehensiveTestform.template;
            json.task_name = this.comprehensiveTestform.taskname; 
            json.check_scheme = this.comprehensiveTestform.scheme; 
            json.priority = this.comprehensiveTestform.priority;
            json.execute_type  = this.comprehensiveTestform.execute_type;
            if(json.execute_type == 2){ //定时  
                json.plan_time = this.comprehensiveTestform.time;
            }
            if(json.execute_type == 3){ //周期
                let t = this.comprehensiveTestform.execute_cycletype+','; 
                if(this.comprehensiveTestform.execute_cycletype == 1){ //月
                    t += this.comprehensiveTestform.execute_cycletype_day+',';
                }
                if(this.comprehensiveTestform.execute_cycletype == 2){ //日
                    t += this.comprehensiveTestform.execute_cycletype_week+',';
                } 
                t += this.comprehensiveTestform.execute_cycletype_starttime;
                json.plan_time = t;
                json.finish_time =  this.comprehensiveTestform.endtime; 
            }
            var portscan =  {
                is_open:this.comprehensiveTestform.isport,
                port_scan_type:this.comprehensiveTestform.scan_type,
                tcp_scan_type:this.comprehensiveTestform.tcp_scan_type,
                udp_scan:this.comprehensiveTestform.udp_scan, 
                scan_rate:'',
                
            };  
            if(portscan.port_scan_type == 1){
                portscan.port = this.comprehensiveTestform.port1;
            }
            if(portscan.port_scan_type == 2){
                portscan.port = this.comprehensiveTestform.port2;
            }
            if(portscan.port_scan_type == 3){
                portscan.port = this.comprehensiveTestform.port;
            } 
            json.port_scan =  JSON.stringify(portscan); 
            let web_crawler = {
                is_open:this.comprehensiveTestform.iswebreptile,
                crawler_depth:this.comprehensiveTestform.crawler_depth,
                crawl_scope:this.comprehensiveTestform.crawler_scope,
                crawl_speed:this.comprehensiveTestform.crawl_speed,
                single_link_timeout:this.comprehensiveTestform.single_link_timeout,
                crawl_range:this.comprehensiveTestform.crawl_range,  
                url_removal: this.comprehensiveTestform.url_removal=='1'?true:false, 
            };  
            web_crawler.black_key = this.comprehensiveTestform.black_key;
            web_crawler.black_url = this.comprehensiveTestform.black_url;
            web_crawler.analyze_file_type = this.comprehensiveTestform.analyze_file_type;
            web_crawler.download_file_type = this.comprehensiveTestform.download_file_type;
            web_crawler.explain_flash = this.comprehensiveTestform.explain_flash;
            web_crawler.is_trigger_js = this.comprehensiveTestform.is_trigger_js; 
            if(this.comprehensiveTestform.tcpport == 2 ){
                web_crawler.self_defined_header = false;
                web_crawler.crawler_header = this.comprehensiveTestform.crawler_header
            }else if(this.comprehensiveTestform.tcpport == 3){
                 web_crawler.self_defined_header = true;
                 web_crawler.crawler_header = this.comprehensiveTestform.crawler_header1;
            } 
            json.web_path_scan = JSON.stringify({
                "is_open": this.comprehensiveTestform.iswebroute
            });
            json.web_crawler = JSON.stringify(web_crawler);
             
            json.word_guess = JSON.stringify({
                "is_open": this.comprehensiveTestform.ispassword, 
                "pass_dict":  this.comprehensiveTestform.pass_dict,
                "guess_number": this.comprehensiveTestform.guess_number,
                "guess_rate":  this.comprehensiveTestform.guess_rate,
                "guess_timeout":  this.comprehensiveTestform.guess_timeout,
            }) 

            var logincheck = {
                "is_open":this.comprehensiveTestform.islogin,
                "login_conf": this.comprehensiveTestform.login_conf,
            };
            json.login_check = JSON.stringify(logincheck)  
            //远控参数，除了3wifi检测，都有
            json.remote_control=JSON.stringify({
                is_open:this.comprehensiveTestform.isRemotecontrol, 
            });   
            return json;
        },
        gettemplete() {  
            this.$ajax.get('/task/template/t_list/',{
                params: {
                    task_type: this.dialogtype
                }
            }).then(dt => {
                let res = dt.data;
                this.templatelist = res;
                for (var i = 0; i < this.templatelist.length; i++) {
                    if (this.templatelist[i].is_default) {
                        this.comprehensiveTestform.template = this.templatelist[ i ].template_id;
                        if(this.flag == "1"){ //新建
                            this.gettemplateconfig(this.comprehensiveTestform.template);
                        } 
                    }
                }
            })
            .catch(res => {});
        },
        getAllscheme() {  
            this.$ajax.get('/schemes/interfaces/name/',{
                params: {}
            }).then(dt => {
                let res = dt.data;
                this.schemelist = res;
            }).catch(res => {});
        },
        getTaskschemes() { 
            //获得协议
            this.$ajax.get('/task/task/schemes/',{
                params: { }
            }).then(dt => {
                let res = dt.data;
                if (res.success) {
                    this.loginagreementlist = res.login_check_scheme;
                    // this.escapeformagreementlist = res.vm_escape_scheme;
                }
            })
            .catch(err => {});
        },
        getdictlist() { 
            this.$ajax.get('/dictionary/keyvalue/dict_list/',{
                params: {
                    dict_type: 1
                }
            }).then(dt => {
                let res = dt.data;
                if (res.success) {
                    this.userdictlist = res.result;
                }
            })
            .catch(err => {});
            this.$ajax.get('/dictionary/keyvalue/dict_list/',{
                params: {
                    dict_type: 2
                }
            }).then(dt => {
                 let res = dt.data;
                    if (res.success) {
                        this.pwddictlist = res.result;
                    }
                })
                .catch(err => {});
            this.$ajax.get('/dictionary/keyvalue/dict_list/',{
                params: {
                    dict_type: 3
                }
            }).then(dt => {
                 let res = dt.data;
                    if (res.success) {
                        this.wifidictlist = res.result;
                    }
                })
                .catch(err => {});
        },
        templateconfig() {
            this.$router.push({
                path: "/templatemanage",
                query: {}
            });
        },
        cancelTask() {
            this.comprehensiveTestVisible = false;
        },
        gettemplateconfig(_id){
            //获得模板配置 
            this.$ajax.get('/task/template/config/',{
                params: {
                    id: _id
                }
            }).then(dt => {
                let res = dt.data;
                if (res.success) {
                    this.comprehensiveTestform.priority = res.priority;
                    this.comprehensiveTestform.execute_type  = res.execute_type;
                    this.comprehensiveTestform.scheme = parseInt( res.check_scheme );
                    this.comprehensiveTestform.strength = res.strength;
                    this.comprehensiveTestform.islogin = res.login_check.is_open;
                    this.comprehensiveTestform.login_conf = res.login_check.login_conf;
                    this.comprehensiveTestform.isport = res.port_scan.is_open;
                    this.comprehensiveTestform.scan_type = res.port_scan.port_scan_type;
                    this.comprehensiveTestform.tcp_scan_type = res.port_scan.tcp_scan_type;
                    this.comprehensiveTestform.udp_scan = res.port_scan.udp_scan;
                    if (this.comprehensiveTestform.scan_type == 1) {
                        this.comprehensiveTestform.port1 = res.port_scan.scan_port;
                        this.comprehensiveTestform.port2 = res.port_scan.scan_port;
                    }
                    if (this.comprehensiveTestform.scan_type == 2) {
                         this.comprehensiveTestform.port1 = res.port_scan.scan_port;
                        this.comprehensiveTestform.port2 = res.port_scan.scan_port;
                    }
                    if (this.comprehensiveTestform.scan_type == 3) {
                        this.comprehensiveTestform.port = res.port_scan.scan_port;
                    }
                     if(!res.web_crawler.self_defined_header ){
                        this.comprehensiveTestform.tcpport = 2;
                        this.comprehensiveTestform.crawler_header =
                        res.web_crawler.headers;
                    }else{
                        this.comprehensiveTestform.tcpport = 3;
                        this.comprehensiveTestform.crawler_header1 =
                        res.web_crawler.crawler_header;
                    }
                    this.comprehensiveTestform.iswebreptile =
                        res.web_crawler.is_open;
                    this.comprehensiveTestform.iswebroute =
                        res.web_crawler.web_path_scan;
                    this.comprehensiveTestform.crawl_range =
                        res.web_crawler.crawl_range;
                    this.comprehensiveTestform.crawler_depth =
                        res.web_crawler.crawler_depth+'';
                    
                    this.comprehensiveTestform.crawler_scope = res.web_crawler.crawl_scope+'';
                    this.comprehensiveTestform.crawl_speed = res.web_crawler.crawl_speed; 
                    res.web_crawler.url_removal ? this.comprehensiveTestform.url_removal  = '1' :this.comprehensiveTestform.url_removal  = '0'
                    this.comprehensiveTestform.black_key = res.web_crawler.black_key;
                    this.comprehensiveTestform.black_url = res.web_crawler.black_url;
                    this.comprehensiveTestform.analyze_file_type = res.web_crawler.analyze_file_type;
                    this.comprehensiveTestform.download_file_type = res.web_crawler.download_file_type;
                    this.comprehensiveTestform.explain_flash = res.web_crawler.explain_flash;
                    this.comprehensiveTestform.is_trigger_js = res.web_crawler.is_trigger_js;

            

                    this.comprehensiveTestform.single_link_timeout =
                        res.web_crawler.single_link_timeout+'';
                    this.comprehensiveTestform.ispassword =
                        res.word_guess.is_open;
                    this.comprehensiveTestform.pass_dict =
                        res.word_guess.pass_dict;
                    this.comprehensiveTestform.user_dict =
                        res.word_guess.user_dict;
                    this.comprehensiveTestform.password_guess_open =
                        res.wifi_config.is_crack;
                    this.comprehensiveTestform.guess_dict =
                        res.wifi_config.passwd_dict;
                    this.comprehensiveTestform.imitate_ap_open =
                        res.wifi_config.is_simulate;
                    this.comprehensiveTestform.imitate_ap_time =
                        res.wifi_config.simulate_duration;
                    this.comprehensiveTestform.code_embedding_open =
                        res.wifi_config.code_embedding_open;
                    this.loginAPtableData = res.wifi_config.ap_list;
                    this.comprehensiveTestform.isRemotecontrol =
                        res.remote_control.is_open; 
                    }
                })
                .catch(err => {});
        },
        resetloginForm(formName){
            this.$refs[formName].resetFields();
        },
        login_validate(str){ //登录测试 
            //  this.$refs['loginform'].validate((valid) => {
            //     if(valid){
                    //登录检测验证
                    this.loginloading = true; //等待效果 
                    this.$ajax.get('/task/task/login_validate/',{
                        params: {
                            validate_type: str,
                            ip: this.loginform.ip,
                            port: this.loginform.port,
                            scheme: this.loginform.agreement,
                            db_name: this.loginform.path,
                            username: this.loginform.username,
                            password: this.loginform.password,
                            cookie:this.loginform.cookie
                        }
                    }).then(dt => {
                        let res = dt.data;
                        if (res.success) {
                            this.loginloading = false;
                            this.isLoginSuccess = true;
                        } else {
                            this.loginloading = false;
                            this.$message({
                                message: "登录验证失败",
                                type: "error"
                            });
                        }
                    })
                    .catch(err => {});
            //     }
            
            // })
            
        },
        addLoginconfig(){   //保存登录检测-新增登录
            if(this.updateLogin){
                this.comprehensiveTestform.login_conf.splice(this.updateIndex, 1);
            } 
            if(!this.loginform.ip){
                this.loginformip = 'is-error';
                this.$message({
                    message: "IP/域名不能为空",
                    type: "error"
                });
                return;
            }else{
                this.loginformip = '';
            }
            if(!this.loginform.port){
                this.loginformport = 'is-error';
                this.$message({
                    message: "端口不能为空",
                    type: "error"
                });
                return;
            }else{
                this.loginformport = '';
            }
            var _j = {
                ip: this.loginform.ip,
                port: this.loginform.port,
                scheme: this.loginform.agreement,
                username: this.loginform.username,
                password: this.loginform.password,
                path: this.loginform.path,
                cookie:this.loginform.cookie,
            };
            this.comprehensiveTestform.login_conf.push(_j);
            this.logininnerVisible = false;
            this.loginform.ip = "";
            this.loginform.port = "";
            this.loginform.agreement = "";
            this.loginform.username = "";
            this.loginform.password = "";
            this.loginform.path = "";
            this.loginform.cookie = '';
            this.isLoginSuccess = false;
            this.$refs["loginform"].resetFields();
        },
        handleCloseLoginform() {
            this.loginloading = false;
            this.logininnerVisible = false;
        },
        handleClosetemplateform() {
            // 添加模板名称取消后恢复
            this.$refs['templateform'].resetFields()
            this.templateVisible = false;
        },  
        submitTemplate() {//保存模板
            //  校验成功后修改数据
            if(!this.templateform.template_name){
                this.templateError = 'is-error';
                this.$message({
                    message: "模板名称不能为空",
                    type: "error"
                });
                return ;
            }else{
                this.templateError = '';
            }
            this.$refs["templateform"].validate(valid => {
                if(!valid) return
                var json = {};  
                var portscan =  {
                    is_open:this.comprehensiveTestform.isport,
                    port_scan_type:this.comprehensiveTestform.scan_type,
                    tcp_scan_type:this.comprehensiveTestform.tcp_scan_type,
                    udp_scan:this.comprehensiveTestform.udp_scan, 
                    scan_rate:'',
                }; 
                if(portscan.port_scan_type == 1){
                    portscan.port = this.comprehensiveTestform.port1;
                }
                if(portscan.port_scan_type == 2){
                    portscan.port = this.comprehensiveTestform.port2;
                }
                if(portscan.port_scan_type == 3){
                    portscan.port = this.comprehensiveTestform.port;
                }
                json.port_scan =  JSON.stringify(portscan);  
                json.check_scheme = this.comprehensiveTestform.scheme; 
                json.priority = this.comprehensiveTestform.priority;  
                let webcrawler = {
                    is_open:this.comprehensiveTestform.iswebreptile,
                    crawler_depth:this.comprehensiveTestform.crawler_depth,
                    crawl_scope:this.comprehensiveTestform.crawler_scope ,
                    crawl_speed:this.comprehensiveTestform.crawl_speed,
                    single_link_timeout:this.comprehensiveTestform.single_link_timeout,
                    crawl_range:this.comprehensiveTestform.crawl_range, 
                    url_removal: this.comprehensiveTestform.url_removal  = '1'?true:false,
                } 
                webcrawler.black_key = this.comprehensiveTestform.black_key;
                webcrawler.black_url = this.comprehensiveTestform.black_url;
                webcrawler.analyze_file_type = this.comprehensiveTestform.analyze_file_type;
                webcrawler.download_file_type = this.comprehensiveTestform.download_file_type;
                webcrawler.explain_flash = this.comprehensiveTestform.explain_flash;
                webcrawler.is_trigger_js = this.comprehensiveTestform.is_trigger_js;



                if(this.comprehensiveTestform.tcpport == 2 ){
                    webcrawler.self_defined_header = false;
                    webcrawler.crawler_header = this.comprehensiveTestform.crawler_header
                }else if(this.comprehensiveTestform.tcpport == 3){
                    webcrawler.self_defined_header = true;
                    webcrawler.crawler_header = this.comprehensiveTestform.crawler_header1;
                } 
                json.web_path_scan = JSON.stringify({
                    "is_open": this.comprehensiveTestform.iswebroute
                });
                json.web_crawler = JSON.stringify(webcrawler);   
                var wordguess = JSON.stringify({
                        "is_open": this.comprehensiveTestform.ispassword, 
                        "pass_dict":  this.comprehensiveTestform.pass_dict,
                        "guess_number": this.comprehensiveTestform.guess_number,
                        "guess_rate":  this.comprehensiveTestform.guess_rate,
                        "guess_timeout":  this.comprehensiveTestform.guess_timeout,
                    })
                json.word_guess = wordguess;
              
                // var logincheck = {
                //     "is_open":this.comprehensiveTestform.islogin,
                //     "login_conf": this.comprehensiveTestform.login_conf,
                // };
                // json.login_check = JSON.stringify(logincheck)  
                //远控参数，除了3wifi检测，都有
                json.remote_control=JSON.stringify({
                    is_open:this.comprehensiveTestform.isRemotecontrol, 
                });
                json.execute_type  = this.comprehensiveTestform.execute_type;
                if(json.execute_type == 2){ //定时
                json.plan_time = this.comprehensiveTestform.time;
                }
                if(json.execute_type == 3){ //周期
                    let t = this.comprehensiveTestform.execute_cycletype+',';
                    
                    if(this.comprehensiveTestform.execute_cycletype == 1){ //月
                        t += this.comprehensiveTestform.execute_cycletype_day+',';
                    }
                    if(this.comprehensiveTestform.execute_cycletype == 2){ //日
                        t += this.comprehensiveTestform.execute_cycletype_week+',';
                    }

                    t += this.comprehensiveTestform.execute_cycletype_starttime;
                    json.plan_time = t;
                    json.finish_time =  this.comprehensiveTestform.endtime;
                } 
                json.template_name = this.templateform.template_name; 
                this.$ajax({
                    method:'post',
                    url:'/task/template/',
                    data:this.qs.stringify(json)
                })
               .then(dt => {
                    let res = dt.data;
                    if(res.success){
                        this.$message({
                            message: "保存模板成功",
                            type: "success"
                        });
                        this.templateVisible= false;
                        this.templateform.template_name = '';
                        this.activeName = 'tabs1';
                        this.comprehensiveTestform.template = parseInt(res.template_id); 
                        // this.gettemplateconfig(res.template_id);
                        
                        this.savetemplateShow();  
                        
                    }else{
                        this.$message({
                            message: "保存模板失败",
                            type: "error"
                        });
                    }
                }).catch(err=>{})
            });
        },
        savetemplateShow(){
            this.$ajax.get('/task/template/t_list/',{
                params: {
                    task_type: this.dialogtype
                }
            }).then(dt => {
                let res = dt.data;
                this.templatelist = res;
                this.gettemplateconfig(this.comprehensiveTestform.template);    
            })
            .catch(res => {});
        },
        Totemplate(){ //跳转到模板管理
            this.$router.push({
                path: `/templatemanage`,
                query: { }
            });
        },
        targetinput(e) {
            this.comprehensiveTestform.taskname =
                this.comprehensiveTestform.target.substr(0, 20) +
                "_" +
                this.commonjs.nowtime();
        }, 
        clickupload() {
            document.querySelector(".btnUploadID").click();
        },
        changeuploaID(e) {
            var that = this;
            var f = e.target.files[0];
            if (!f) return;
            // that.taskfile.name = f.name;
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1);
            if (fileSuffix.indexOf("xls") != -1) {
                var reader = new FileReader();
                reader.onload = function(e) {
                    var data = e.target.result;
                    if (that.rABS) {
                        that.wb = XLSX.read(btoa(fixdata(data)), {
                            type: "base64"
                        });
                    } else {
                        that.wb = XLSX.read(data, {
                            type: "binary"
                        });
                    } 
                    let carData = XLSX.utils.sheet_to_json(
                        that.wb.Sheets[that.wb.SheetNames[0]]
                    ); 
                    let arr = [];
                    for (var key in carData) { 
                        for (var k in carData[key]) {
                            if (arr.indexOf(k) === -1) { 
                                arr.push(k);
                            }
                            if (arr.indexOf(carData[key][k]) === -1) { 
                                arr.push(carData[key][k]);
                            }
                        }
                    } 
                    that.comprehensiveTestform.target = arr.join("\n");
                };
                if (that.rABS) {
                    reader.readAsArrayBuffer(f);
                } else {
                    reader.readAsBinaryString(f);
                }
            } else if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function() {
                    if (reader.result) { 
                        that.comprehensiveTestform.target = reader.result;
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        btnUpdateTemplate(){ //编辑模板
            this.noUpdte = true;
            this.updateTxt = '编辑中';
        },
        saveTemplate(){ //编辑模板--保存 
            var json = {};  
            json.id = this.templateid;
            json.template_name = this.comprehensiveTestform.template_name;
            json.execute_type  = this.comprehensiveTestform.execute_type;
            if(json.execute_type == 2){ //定时
            json.plan_time = this.comprehensiveTestform.time;
            }
            if(json.execute_type == 3){ //周期
                let t = this.comprehensiveTestform.execute_cycletype+',';
                
                if(this.comprehensiveTestform.execute_cycletype == 1){ //月
                    t += this.comprehensiveTestform.execute_cycletype_day+',';
                }
                if(this.comprehensiveTestform.execute_cycletype == 2){ //日
                    t += this.comprehensiveTestform.execute_cycletype_week+',';
                }

                t += this.comprehensiveTestform.execute_cycletype_starttime;
                json.plan_time = t;
                json.finish_time =  this.comprehensiveTestform.endtime;
            }
            var portscan =  {
                is_open:this.comprehensiveTestform.isport,
                port_scan_type:this.comprehensiveTestform.scan_type,
                tcp_scan_type:this.comprehensiveTestform.tcp_scan_type,
                udp_scan:this.comprehensiveTestform.udp_scan, 
                scan_rate:'',
            }; 
            if(portscan.port_scan_type == 1){
                portscan.port = this.comprehensiveTestform.port1;
            }
            if(portscan.port_scan_type == 2){
                portscan.port = this.comprehensiveTestform.port2;
            }
            if(portscan.port_scan_type == 3){
                portscan.port = this.comprehensiveTestform.port;
            }
            json.port_scan =  JSON.stringify(portscan);  
            json.check_scheme = this.comprehensiveTestform.scheme, 
            json.priority = this.comprehensiveTestform.priority;
            // let web_crawler = {
            //     is_open:this.comprehensiveTestform.iswebreptile,
            //     crawler_depth:this.comprehensiveTestform.crawler_depth,
            //     single_link_timeout:this.comprehensiveTestform.single_link_timeout,
            //     crawl_range:this.comprehensiveTestform.crawl_range, 
                  
            // };
            let webcrawler = {
                is_open:this.comprehensiveTestform.iswebreptile,
                crawler_depth:this.comprehensiveTestform.crawler_depth,
                crawl_scope:this.comprehensiveTestform.crawler_scope ,
                crawl_speed:this.comprehensiveTestform.crawl_speed,
                single_link_timeout:this.comprehensiveTestform.single_link_timeout,
                crawl_range:this.comprehensiveTestform.crawl_range, 
                url_removal: this.comprehensiveTestform.url_removal  = '1'?true:false,
            } 
            webcrawler.black_key = this.comprehensiveTestform.black_key;
            webcrawler.black_url = this.comprehensiveTestform.black_url;
            webcrawler.analyze_file_type = this.comprehensiveTestform.analyze_file_type;
            webcrawler.download_file_type = this.comprehensiveTestform.download_file_type;
            webcrawler.explain_flash = this.comprehensiveTestform.explain_flash;
            webcrawler.is_trigger_js = this.comprehensiveTestform.is_trigger_js;

 
            
            json.web_path_scan = JSON.stringify({
                    "is_open": this.comprehensiveTestform.iswebroute
                });
            if(this.comprehensiveTestform.tcpport == 2 ){
                web_crawler.self_defined_header = false;
                web_crawler.crawler_header = this.comprehensiveTestform.crawler_header
            }else if(this.comprehensiveTestform.tcpport == 3){
                web_crawler.self_defined_header = true;
                web_crawler.crawler_header = this.comprehensiveTestform.crawler_header1;
            } 
            json.web_crawler = JSON.stringify(web_crawler);   
            
            // var wordguess = JSON.stringify({
            //         "is_open": this.comprehensiveTestform.ispassword, 
            //         "pass_dict":  this.comprehensiveTestform.pass_dict,
            //         "guess_number": this.comprehensiveTestform.guess_number,
            //         "guess_rate":  this.comprehensiveTestform.guess_rate,
            //         "guess_timeout":  this.comprehensiveTestform.guess_timeout,
            //     })
            // json.word_guess = wordguess;
            
            var logincheck = {
                "is_open":this.comprehensiveTestform.islogin,
                // "login_conf": this.comprehensiveTestform.login_conf,
            };
            json.login_check = JSON.stringify(logincheck)  
            //远控参数，除了3wifi检测，都有
            json.remote_control=JSON.stringify({
                is_open:this.comprehensiveTestform.isRemotecontrol, 
            });
                 
            this.$ajax({
                method:'put',
                url:'/task/template/'+this.templateid+'/',
                data:this.qs.stringify(json)
            })
            .then(dt => {
                let res = dt.data;
                if(res.success){  
                    this.dangerBtn();
                }else{
                    this.$message({
                        message: res.error,
                        type: "error"
                    });
                }
            }).catch(err=>{})
        },
        dellogin(index, rows){ //移除添加的登录信息
            rows.splice(index, 1);
        },
        updatelogin(index, rows){ //编辑
            this.logininnerVisible = true;  
            this.updateLogin = true;
            this.updateIndex = index;
            this.loginform.ip = rows[0].ip;
            this.loginform.port = rows[0].port;
            this.loginform.agreement = rows[0].scheme;
            this.loginform.username = rows[0].username;
            this.loginform.password = rows[0].password;
            this.loginform.path = rows[0].path;
            this.loginform.cookie = rows[0].cookie; 
        },
        AddLogin(){
            this.logininnerVisible = true;  
            this.updateLogin = false;
            this.updateIndex =-1;
        },
    }
}
</script>
