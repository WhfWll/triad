<template>
    <div>
        <div class="main-title  ">
            系统授权1
        </div>
        <div class="systembox context_box_bg">
           
            <el-tabs v-model="activeName" @tab-click="handleClick">
                <!-- <el-tab-pane label="基础参数" name="tab1"> 
                    <div class="tabsbox  systeminfolist">
                        <div>
                            <label>IP地址</label>
                            <span>{{formBasics.ip}}</span>
                        </div>
                        <div>
                            <label>运行端口</label>
                            <span>{{ formBasics.port }}</span>
                        </div>
                        <div>
                            <label>CPU</label>
                            <span>{{formBasics.cpu}}</span>
                        </div>
                        <div>
                            <label>内存</label>
                            <span>{{formBasics.memories}}</span>
                        </div> 
                    </div>
                </el-tab-pane>  -->
                <!-- <el-tab-pane label="检测设置" name="tab3">
                    <div class="tabsbox">
                        <el-form ref="form" :model="formcheck" label-width="100px"> 
                            <el-form-item label="扫描端口：">
                                <el-input v-model="formcheck.port" class="iptw" placeholder="请输入扫描端口"></el-input>
                            </el-form-item>
                            <el-form-item label="user-agent：">
                                <el-input v-model="formcheck.useragent" class="iptw" placeholder="请输入user-agent"></el-input>
                            </el-form-item>
                            <el-form-item label="扫描线程：">
                                <el-input v-model="formcheck.thread" class="iptw" placeholder="请输入扫描线程"></el-input>
                            </el-form-item>
                            <el-form-item label="超时时长 ：">
                                <el-input v-model="formcheck.times" class="iptw" placeholder="请输入超时时长"></el-input>
                            </el-form-item>
                            <el-form-item label="cookie ：">
                                <el-input v-model="formcheck.cookie" class="iptw" placeholder="请输入cookie"></el-input>
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" @click="handleSubmitCheck">保存</el-button>
                                <el-button @click="handleRecoveryDefault">恢复默认设置</el-button>
                            </el-form-item>
                        </el-form>
                    </div>
                </el-tab-pane> -->
                <!-- <el-tab-pane label="消息设置" name="tab4">
                    <div class="tabsbox">
                        <el-form ref="formMessage" :model="formMessage" label-width="100px">
                            <el-form-item label="提醒方式：">
                                <el-select v-model="formMessage.remindtype" placeholder="请选择提醒方式" class="iptw" @change="changeremindtype">
                                    <el-option label="未设置提醒" value="0"></el-option>
                                    <el-option label="邮件提醒" value="1"></el-option> 
                                </el-select>
                            </el-form-item>
                            <el-form-item label="邮箱地址：" prop="email"  > 
                                <el-input v-model="formMessage.email" class="iptw"  placeholder="请输入邮箱地址" ></el-input>
                            </el-form-item>
                            <el-form-item label="提醒阈值："  > 
                                <el-select v-model="formMessage.threshold" placeholder="请选择提醒阈值" class="iptw">
                                    <el-option label="检测到高危进行提醒" value="1"></el-option>
                                    <el-option label="检测到中危及以上进行提醒" value="2"></el-option>
                                    <el-option label="检测到低危及以上进行提醒" value="3"></el-option>
                                </el-select>
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" @click="handleSubmitMessage('formMessage')">保存</el-button> 
                            </el-form-item>
                        </el-form>
                    </div> 
                </el-tab-pane>  -->
                <el-tab-pane label="产品授权 " name="tab6" v-if="true">
                     <div class="tabsbox productinfo systeminfolist">
                        <div style="border-top: 1px solid #E8E8F5;">
                            <label>产品名称</label>
                            <span>{{productinfo.name}}</span>
                        </div>
                        <div>
                            <label>系1统版本</label>
                            <span>{{productinfo.software_version}}
                            </span>
                        </div>
                        <!-- <div>
                            <label>漏洞版本</label>
                            <span>{{productinfo.engine_version}} 
                            </span>
                        </div> -->
                        <!-- <div>
                            <label>软件类型</label>
                            <span>{{productinfo.version_type}}</span>
                        </div> -->
                        <!-- <div>
                            <label>使用期限</label>
                            <span>{{productinfo.time}}</span>
                        </div> -->
                        <!-- <div>
                            <label>登录地址</label>
                            <span>{{productinfo.login_url}}</span>
                        </div> -->
                        <div>
                            <label>使用期限</label>
                            <div class="system">
                                <div class="system_data">
                                     授权日期：<span style="margin:0 10px">{{productinfo.authTime}}</span>
                                   授权时长： <span style="margin:0 10px">{{productinfo.authDays}}</span>
                                   剩余时长： <span style="margin:0 10px">{{productinfo.leftDays}}</span>
                                </div>
                                <div style="float:right">
                                    <el-button type="primary" size="small" style="display:none">在线授权</el-button>
                                    <!-- <el-button type="primary" size="small" @click="fileDialog = !fileDialog">文件授权</el-button> -->
                                    <xzbutton type="primary" @click="fileDialog = !fileDialog" size="small">系统授权22
                                    </xzbutton>
                                    <!-- <el-button type="primary" size="small" @click="handleGetData">收集授权信息</el-button> -->
                                    <xzbutton v-if="false" type="primary" @click="handleGetData" size="small" style="margin-left: 8px;">收集授权信息</xzbutton>
                                </div>
                            </div>
                        </div>
                    </div>
                </el-tab-pane>
                <!-- <el-tab-pane label="关机/重启" name="tab7">
                    <div class="tabsbox productinfo systeminfolist">
                        <div>
                            <label>系统关机</label>
                            <span>关机后需要手动打开设备电源才能启动设备
                                <div style="float:right"> 
                                    <el-checkbox v-model="shutdown" @change="chkshutdown"  style="color:#4C7AE3">立即关机</el-checkbox>
                                </div>
                            </span>
                        </div>
                        <div>
                            <label>快速重启</label>
                            <span>快速重启需要1~3分钟
                                <div style="float:right"> 
                                    <el-checkbox v-model="fastrestart" @change="chkfastrestart"  style="color:#4C7AE3">立即重启</el-checkbox>
                                </div>
                            </span>
                        </div>
                        <div>
                            <label>系统重启</label>
                            <span>系统重启需要5~10分钟
                                <div style="float:right"> 
                                    <el-checkbox v-model="systemrestart" @change="chksystemrestart"  style="color:#4C7AE3">立即重启</el-checkbox>
                                </div>
                            </span>
                        </div> 
                    </div>
                    <p style="font-size:12px;margin-top:20px; color:rgba(72, 72, 102, 0.32)"><i class="iconfont icontishi" style="margin-right:10px;vertical-align: middle;"></i>快速重启是小智服务重启，系统重启是整机重启</p>
                    <el-button type="primary"
                        class="btnloading" style="opacity:0"
                        v-loading.fullscreen.lock="fullscreenLoading"> 
                    </el-button>
                </el-tab-pane> -->
                <!-- <el-tab-pane label="告警设置" name="tab8" class="nopaddingbox">
                    <div class="tabsbox ">
                        <el-form ref="warningform" :model="warningform" label-width="114px" class="sysform"
                            :rules="setUp">
                            <el-form-item label="CPU告警阈值：" class="syswarnvalue" prop="cpu_threshold">
                                <el-input v-model="warningform.cpu_threshold" style="width:200px"></el-input> %
                                <el-tooltip placement="right">
                                    <div slot="content">告警范围1%~100%</div>
                                    <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                                </el-tooltip>
                            </el-form-item>
                            <el-form-item label="内存告警阈值：" class="syswarnvalue" prop="memory_threshold">
                                <el-input v-model="warningform.memory_threshold" style="width:200px"></el-input> %
                                <el-tooltip placement="right">
                                    <div slot="content">告警范围1%~100%</div>
                                    <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                                </el-tooltip>
                            </el-form-item>
                            <el-form-item label="硬盘告警阈值：" class="syswarnvalue" prop="disk_threshold">
                                <el-input v-model="warningform.disk_threshold" style="width:200px"></el-input> %
                                <el-tooltip placement="right">
                                    <div slot="content">告警范围1%~100%</div>
                                    <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                                </el-tooltip>
                            </el-form-item>
                            <el-form-item label="流量告警阈值：" class="syswarnvalue" prop="flow_threshold">
                                <el-input v-model="warningform.flow_threshold" style="width:200px"></el-input> Mbps
                                <el-tooltip placement="right">
                                    <div slot="content">告警范围1~1000M</div>
                                    <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                                </el-tooltip>
                            </el-form-item>
                            <el-form-item class="savebox">
                                <xzbutton type="primary" @click="btnSaveWarning" size="small">保存</xzbutton>
                            </el-form-item>
                        </el-form>
                    </div>
                </el-tab-pane> -->
                <el-tab-pane label="安全设置" name="tab9" v-if="true">
                    <safetyset v-if="activeName === 'tab9'"/>
                </el-tab-pane>
                <!-- <el-tab-pane label="安全设置" name="tab9">
                    <div class="tabsbox">
                        <el-form ref="securityform" :model="securityform" label-width="150px" class="securityform"
                            :rules="safety">
                            <el-form-item label="普通用户登录限制：" prop="user_limit">
                                <el-input v-model="securityform.user_limit"
                                    style="width:200px;margin-right:10px;margin-left:18px"></el-input>次
                                <el-tooltip placement="bottom">
                                    <div slot="content">普通用户一天内连续登录错误超过限制次数时，<br />系统在限定时间内禁止普通用户登录</div>
                                    <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                                </el-tooltip>
                            </el-form-item>
                            <div class="clearfix">
                                <el-form-item label="管理/审计员登录限制：" label-width="170px" prop="admin_limit"
                                    style="float:left;width:450px">
                                    <el-input v-model="securityform.admin_limit" style="width:200px;margin-right:10px">
                                    </el-input>次
                                    <el-tooltip placement="bottom">
                                        <div slot="content">管理员/审计员一天内连续登录错误超过限制次数时</div>
                                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                                    </el-tooltip>
                                </el-form-item>
                                <el-form-item label="禁止登录：" style="float:left;width:500px">
                                    <el-input v-model="securityform.ban_time" style="width:200px;margin-right:10px">
                                    </el-input>分钟
                                    <el-tooltip placement="bottom">
                                        <div slot="content">系统在限定时间内禁止管理员/审计员登录</div>
                                        <i class="iconfont icontishi icontsstyle" style="vertical-align: middle; "></i>
                                    </el-tooltip>
                                </el-form-item>
                            </div>
                            <el-form-item class="savebox">
                                <xzbutton type="primary" @click="benSavesecurity" size="small">保存</xzbutton>
                            </el-form-item>
                        </el-form>
                    </div>
                </el-tab-pane> -->
                <el-tab-pane label="系统升级" name="tab10" v-if="false">
                    <div class="tabsbox productinfo sysupgradebox systeminfolist">
                        <div class="firstline">
                            <label>当前系统版本</label>
                            <span class="blueword"  >{{systemupdate.current_sys_version}}</span>
                            <!-- <span class="blueword">{{systemupdate.current_sys_version}}</span> -->
                            <span class="updatetime">更新时间 <span class="timeword">{{systemupdate.update_time}}</span>
                            </span>

                        </div>
                        <div class="secondline" style="">
                            <label>更新文件</label>
                            <span class="spversion"
                                style="width: calc(100% - 160px);display:flex;justify-content: space-between;">
                                <label v-if="!newVersion">{{uploadFileName}}</label>
                                <!-- <span class="greytip"><span class="icon-name">下载中</span><i class="el-icon-loading"></i></span> -->
                                <span class="orangetip" v-if="updateStatus == 'fail'"><span
                                        class="icon-name">上传失败！</span></span>
                                <span class="greytip" v-if="updateStatus == 'running'"><span
                                        class="icon-name">上传中</span><i class="el-icon-loading"></i></span>
                                <span class="bluetip" v-if="updateStatus == 'finish'"><span
                                        class="icon-name">已上传！</span></span>
                            </span>
                        </div>
                        <div style="border-bottom:none">
                            <div style="margin-top:24px">
                                <xzbutton type="primary" @click="btnoffupdate" size="small">离线更新</xzbutton>
                                <input type="file" class="btnUploadID2" ref="upload" @change="changeuploaID2($event)"
                                    style="display:none" id="input-file-ID2" />
                            </div>
                            <p class="systemtip">
                                提示：1 系统更新时，系统将停止所有运行的任务。 2 系统更新过程持续时间10分钟左右，更新过程中不能访问系统。 3 离线更新包括用户上传更新文件和系统加载更新文件两步。
                            </p>
                        </div>
                        <!-- 2022年三月初新设计 -->

                    </div>
                    <!-- <div class=" tabsbox productinfo sysupgradebox systeminfolist"> -->

                       
                    <!-- </div> -->
                    <!-- <div class="tabsbox productinfo sysupgradebox systeminfolist"> -->
                    <!-- <div>
                            <label>在线自动更新</label>
                             <div class="planTime">
                                <el-select 
                                    v-model="onlineUpdate.plan_time1"
                                    size="small"
                                    placeholder="请选择" 
                                    style="width:126px;margin-right:16px;margin-left:20px" >
                                    <el-option value="1" label="周"></el-option>
                                    <el-option value="2" label="月"></el-option>
                                </el-select>
                             <el-select 
                                    v-model="onlineUpdate.plan_time2"
                                    size="small"
                                    placeholder="请选择"
                                    style="width:126px;margin-right:16px"  v-if="onlineUpdate.plan_time1 === '2'">
                                    <el-option   
                                        v-for="i in 31" 
                                        :key="i" 
                                        :label="i+'号'"
                                        :value="i"> </el-option> 
                                </el-select>
                            <el-select 
                                    
                                    v-model="onlineUpdate.plan_time3"
                                    size="small"
                                    placeholder="请选择"
                                    style="width:126px;margin-right:16px" v-if="onlineUpdate.plan_time1 === '1'">
                                    <el-option value="1" label="周一"></el-option>
                                    <el-option value="2" label="周二"></el-option>
                                    <el-option value="3" label="周三"></el-option>
                                    <el-option value="4" label="周四"></el-option>
                                    <el-option value="5" label="周五"></el-option>
                                    <el-option value="6" label="周六"></el-option>
                                    <el-option value="7" label="周日"></el-option>
                                </el-select>
                                <el-time-picker v-model="onlineUpdate.plan_time4" placeholder="选择时间" value-format="HH:mm:ss">
                                </el-time-picker>
                                <el-switch
                                    class="switchbtn"
                                    active-color="#4c7ae3"
                                    v-model="onlineUpdate.isopen"
                                    inactive-color="#E8E8F5"
                                    @change="fnReptileSwitch">
                                </el-switch>
                           </div>
                        </div> -->


                    <!-- </div> -->


                </el-tab-pane>

                <el-tab-pane label="配置备份" name="tab13" v-if="false">
                    <sysconfigcopy v-if="activeName === 'tab13'"/>
                </el-tab-pane>
                <el-tab-pane label="API秘钥" name="tab11" v-if="false">
                    <!-- <el-button type="primary"  @click="generateKey" size="small" style="margin-bottom:24px;">生成秘钥</el-button> -->
                    <div>
                        <xzbutton type="primary" @click="generateKey" size="small">生成秘钥</xzbutton>
                        <el-tooltip placement="right" style="margin-bottom:24px;display:inline-block;">
                            <div slot="content">通过API接口可以调用系统功能模块，API密钥是通过API接口登录的凭证，每个用户允许生成唯一密钥</div>
                            <i class="iconfont icontishi icontsstyle" style="vertical-align: initial; "></i>
                        </el-tooltip>
                    </div>

                    <div class="totalbox">
                        <el-table class="apitable tablePadding" :data="keyTableData" style="width: 100%;">
                            <el-table-column prop="username" label="用户名">
                            </el-table-column>
                            <el-table-column prop="token" label="API认证秘钥">
                            </el-table-column>
                            <el-table-column prop="add_time" width="300" label="生成日期">
                            </el-table-column>
                        </el-table>
                    </div>

                </el-tab-pane>
                    <el-tab-pane label="系统监控" name="tab12" v-if="false">
                    <sysmonitoring v-if="activeName === 'tab12'"/>
                </el-tab-pane>
                <el-tab-pane label="其他设置" name="tab4" v-if="false">
                    <othersysset v-if="activeName === 'tab4'" />
                </el-tab-pane>
                <el-tab-pane label="业务设置" name="tab2" v-if="false">
                    <businessset v-if="activeName === 'tab2'" />
                </el-tab-pane>
            </el-tabs>

        </div>
        <el-dialog title="设置权限" :visible.sync="dialogFormVisible" width="30%" :before-close="cancelform"
            :close-on-click-modal="false" :show-close="false">
            <el-form :model="formJurisdiction" label-width="100px" status-icon>
                <el-form-item label="激活用户：">
                    <el-switch v-model="formJurisdiction.user"></el-switch>
                </el-form-item>
            </el-form>
            <div slot="footer" class="dialog-footer">
                <el-button @click="cancelform()">取 消</el-button>
                <el-button type="primary" @click="submitForm()">确 定</el-button>
            </div>
        </el-dialog>
        <el-dialog title="系统版本升级" :visible.sync="systemupgradeDialogVisible" width="40%"
            :before-close="systemupgradehandleClose" :close-on-click-modal="false" :show-close="false">
            <div :class="[isfirst ? 'first systemupgrade upgradeactive' : 'first systemupgrade']">
                <p>当前软件版本：V2.0190</p>
                <div style="margin-top:20px;margin-bottom:20px;">
                    <label><i style="color:red">*</i> 升级文件：</label>
                    <el-link type="primary" class="btnUpload" @click="clickupload()"><i
                            class="iconfont iconfujian"></i>上传升级文件</el-link>
                    <!-- <el-button type="primary"  plain  class="btnUpload" @click="clickupload()">上传文件</el-button>     -->
                    <input type="file" name="" class="btnUploadID" @change="changeuploaID($event)" style="display:none"
                        id="input-file-ID" accept="*.zip">
                    <label>{{systemfileName}}</label>
                </div>
                <p>为保障升级成功，请选择原厂提供的目标升级文件</p>
            </div>
            <div :class="[!isfirst ? 'second systemupgrade upgradeactive' : 'second systemupgrade']">
                <p style="margin-bottom:20px">系统升级后需要重起才能生效</p>
                <div>
                    <el-radio v-model="restart" label="0">稍后重启</el-radio>
                    <el-radio v-model="restart" label="1">立即重启</el-radio>
                </div>
            </div>
            <div slot="footer" class="dialog-footer">

                <el-button @click="systemupgradehandleClose()">取 消</el-button>
                <el-button v-if="!isfirst" @click="systemupgradeprestep">上一步</el-button>
                <el-button v-if="isfirst" @click="systemupgradenext">下一步</el-button>
                <el-button v-if="!isfirst" type="primary" @click="systemupgradesubmitForm()">确 定</el-button>
            </div>
        </el-dialog>
        <el-dialog title="知识库版本升级" :visible.sync="knowledgeDialogVisible" width="40%"
            :before-close="knowledgehandleClose" :close-on-click-modal="false" :show-close="false">
            <div :class="[knowledgeisfirst ? 'first knowledgeupgrade upgradeactive' : 'first knowledgeupgrade']">
                <p>当前软件版本：V2.0190</p>
                <div style="margin-top:20px;margin-bottom:20px;">
                    <label><i style="color:red">*</i> 升级文件：</label>
                    <el-link type="primary" class="btnUpload1" @click="clickupload1()"><i
                            class="iconfont iconfujian"></i>上传升级文件</el-link>
                    <!-- <el-button type="primary"  plain  class="btnUpload" @click="clickupload()">上传文件</el-button>     -->
                    <input type="file" name="" class="btnUploadID1" @change="changeuploaID1($event)"
                        style="display:none" id="input-file-ID" accept="*.zip">
                    <label>{{systemfileName}}</label>
                </div>
                <p>为保障升级成功，请选择原厂提供的目标升级文件</p>
            </div>
            <div :class="[!knowledgeisfirst ? 'second knowledgeupgrade upgradeactive' : 'second knowledgeupgrade']">
                <p style="margin-bottom:20px">知识库版本升级后需要重起才能生效</p>
                <div>
                    <el-radio v-model="knowledgerestart" label="0">稍后重启</el-radio>
                    <el-radio v-model="knowledgerestart" label="1">立即重启</el-radio>
                </div>
            </div>
            <div slot="footer" class="dialog-footer">

                <el-button @click="knowledgehandleClose()">取 消</el-button>
                <el-button v-if="!knowledgeisfirst" @click="knowledgeprestep">上一步</el-button>
                <el-button v-if="knowledgeisfirst" @click="knowledgenext">下一步</el-button>
                <el-button v-if="!knowledgeisfirst" type="primary" @click="knowledgesubmitForm()">确 定</el-button>
            </div>
        </el-dialog>



        <el-dialog title="系统关机" :visible.sync="shutdownDialogVisible" width="40%" :before-close="shutdownhandleClose"
            :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="shutdownsubmitForm">确定</el-button>
                <el-button size="small" @click="shutdownhandleClose">取消</el-button>
            </div>
            <div class="dialogtxt">关机后需要手动打开设备电源才能起动设备，确定关机？</div>
            <!-- <div slot="footer" class="dialog-footer">
                <el-button @click="shutdownhandleClose()">取 消</el-button>
                <el-button type="primary"  @click="shutdownsubmitForm()">确 定</el-button>
            </div> -->
        </el-dialog>
        <el-dialog title="快速重启" :visible.sync="fastrestartDialogVisible" width="40%"
            :before-close="fastrestarthandleClose" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="fastrestartsubmitForm">确定</el-button>
                <el-button size="small" @click="fastrestarthandleClose">取消</el-button>
            </div>
            <div class="dialogtxt">确定重起小智服务？</div>
            <!-- <div slot="footer" class="dialog-footer">
                <el-button @click="fastrestarthandleClose()">取 消</el-button>
                <el-button type="primary"  @click="fastrestartsubmitForm()">确 定</el-button>
            </div> -->
        </el-dialog>
        <el-dialog title="系统重启" :visible.sync="systemrestartDialogVisible" width="40%"
            :before-close="systemrestarthandleClose" :close-on-click-modal="false" :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="systemrestartsubmitForm">确定</el-button>
                <el-button size="small" @click="systemrestarthandleClose">取消</el-button>
            </div>
            <div class="dialogtxt">确定重起小智服务？</div>
            <!-- <div slot="footer" class="dialog-footer">
                <el-button @click="systemrestarthandleClose()">取 消</el-button>
                <el-button type="primary"  @click="systemrestartsubmitForm()">确 定</el-button>
            </div> -->
        </el-dialog>
        <el-dialog title="系统授权" width='1184px' :visible.sync="fileDialog" :close-on-click-modal="false"
            :show-close="false" class="only_dialog">
            <div class="dialog_b_btn only_dialog">
               <el-button  size="small" @click="handleFileUpload">立即授权</el-button>
                <el-button size="small" @click="fileDialog = !fileDialog">取消</el-button>
            </div>
            <div style="padding:24px" class="">
                <div class="">
                    <div style="border-left:4px solid #4c7ae3;padding-left:20px;margin-bottom:40px">
                        <span class="formbox_upload_file">系统授权码</span>
                        <!-- <span an class="formbox_upload_filename" style="fontSize:15px;opacity:.6;">{{filename}}</span> -->
                    </div>
                    <div>
                        <el-input type="textarea"  :rows="6" v-model="systemSTR"  style="" placeholder="请输入系统授权码"
                            maxlength="50000"></el-input>
                    </div>
                    <div class="formbox_upload_btn" v-if="true">
                        <el-upload  v-show="false" :limit="1" :show-file-list="false" :on-change="handleUploadChange"
                            :auto-upload="false" action="https://jsonplaceholder.typicode.com/posts/"
                            class="formbox_upload_btn_upload" ref="uploadRef">
                            <el-button type="primary" size="small">导入</el-button>
                        </el-upload>
                        <!-- <el-button type="primary" class="" size="small" @click="handleFileUpload">授权</el-button> -->
                    </div>
                </div>
            </div>
        </el-dialog>


        <el-dialog title="生成秘钥" :visible.sync="dialogKeyVisible" :before-close="cancelKey" :close-on-click-modal="false"
            :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitKey">生成秘钥</el-button>
                <el-button size="small" @click="cancelKey">取消</el-button>
            </div>
            <div style="padding:24px">
                <label class="dialog_item_label">用户名</label>
                <el-select v-model="userKey" placeholder="请选择">
                    <el-option v-for="(item,index) in userNames" :key="index" :label="item" :value="item">
                    </el-option>
                </el-select>
            </div>
        </el-dialog>
        <!-- 升级还原确定弹窗 -->
        <el-dialog title="离线更新" :visible.sync="dialogSysVisible" :before-close="cancelKey" :close-on-click-modal="false"
            :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="sureUpload">确定</el-button>
                <el-button size="small" @click="dialogSysVisible = false">取消</el-button>
            </div>
            <div style="padding:24px">
                确定要上传离线更新文件包吗？
            </div>
        </el-dialog>
    </div>
</template>
<style scoped lang="less">
    .tabsbox{
        background: #fff;
        padding: 24px;
        border-radius: 4px;
        box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.1200);
        margin-bottom: 10px;
    }

    .totalbox{
        padding: 25px;
        background: #fff;
        box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.1200);
        border-radius: 4px 4px 4px 4px;
        margin-bottom: 10px;
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
    .updatetime{
        float: right!important;
        border-left: none!important;
    }
    /deep/ .el-icon-loading{
        display: inline-block;
        font-size: 18px;
        color: #4C7AE3;
        margin-left: 6px;
    }
    .systemtip{
        margin-top: 15px;
        font-family: Source Han Sans CN-Regular, Source Han Sans CN;
        font-size: 13px;
        font-weight: 400;
        color: #F35F28;
    }
    .greytip{
        float: right;
        font-size: 13px;
        font-weight: 400;
        color: rgba(72, 72, 102, 0.6399999856948853);
    }
    .bluetip{
        float: right;
        font-size: 13px;
        font-weight: 400;
        color: #4C7AE3;
    }
    .orangetip{
        float: right;
        font-size: 13px;
        font-weight: 400;
        color: #F35F28;

    }
    /deep/ .el-checkbox__label{
        font-size: 13px;
    }
    /* /deep/ .el-form-item__label{
        color: rgba(72, 72, 102, 0.64);
        font-weight: 500;
        font-size: 13px;
    } */
    /deep/ .el-form-item__content{
        color: rgba(72, 72, 102, 0.64);
        font-weight: 500;
        font-size: 13px;
    }
    /deep/ .el-dialog{
        height: 192px !important;
    }
    .dialogtxt{
        text-align: center;
        margin-top: 55px;
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
    /deep/ .el-tabs__header{
        margin: 0 0 24px;
    }
    /deep/ .el-table .cell, .el-table th div, .el-table--border td:first-child .cell, .el-table--border th:first-child .cell{
        padding-left:32px;
    }
    
     
.upgradeactive{
    display: block !important;
}
.systemupgrade{
    display: none;
}
.knowledgeupgrade{
    display: none;
}
.sysupgradebox{
    /* margin-bottom: 24px; */
    .spversion{
        label{
            display: inline-block;
            width: 150px;
        }
        span{
            display: inline-block;
            // width: 200px;
        }
    }
}
.context_box_bg{
    background: none;
}
/deep/ .el-tabs__header {
    margin: 0 0 15px;
    background: #fff;
    border-radius: 4px;
    box-shadow: 0px 2px 4px 0px rgb(76 122 227 / 12%);
    border: none;
}
.systembox{
    //
	// background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    // padding: 24px;
    // box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
    /deep/ .el-tabs{
        // min-height: 700px;
    } 
    /deep/ .el-tabs__item.is-top.is-active{
        // background: #fff;
    }
    /deep/ .el-tabs__header{
        margin: 0 0 15px;
        background: #fff;
        border-radius: 4px;
        box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
        border: none;
    }

} 
/deep/ .el-switch.is-checked .el-switch__core{
    border-color: #4C7AE3;
    background-color: #4C7AE3;
}
.system{ 
    div {
        float:left;
        text-align: center;
    }
}
.iptw{
    width: 500px !important;
}
.systeminfolist{
    /* width: 88%; */
    .firstline{
        border-top: 1px solid #E8E8F5;
        margin-bottom:20px;
    }
    .planTime{
        height: 50px;
        line-height:50px;
        
    }
    .switchbtn{
        float: right;
        margin-top: 15px;
        margin-right: 20px;
    }
    .blueword{
        color:#4C7AE3;
        cursor: pointer;
    }
    .timeword{
        margin-left: 10px;
    }
    .system{
        // display: flex;
        // align-items: center;
        // justify-content: space-between;
        font-weight: 500;
        line-height: 30px;
        padding: 10px 20px;
        padding: 10px 20px;
    line-height: 30px;
        width: calc(100% - 160px);
    float: left;
    box-sizing: border-box;
        .system_data{
            background: #fff;
            color: rgba(72, 72, 102, 0.64);
            font-weight: 500;

        }
    }
} 
.systeminfolist>div {
    /* border-top: 1px solid #E8E8F5; */
    border-bottom: 1px solid #E8E8F5;
    overflow: hidden;
    // background: #f5f8fb;
    font-size: 13px;
}
.systeminfolist >div >label,
.systeminfolist >div >span{
    padding: 10px 20px;
    line-height: 30px;
    float: left;
    box-sizing: border-box;
 
}
.systeminfolist>div>label {
    display: inline-block;
    width: 160px;
    text-align: center;
    border-right: 1px solid #E8E8F5;
    background: #F7F7FB ;
    color: rgba(72, 72, 102, 0.87);
    font-weight: 500;
}
.systeminfolist>div>span {
    display: block;
    /* width: calc(100% - 160px); */
    border-left: 1px solid #E8E8F5;
    margin-left: -1px;
    background: #fff;
    min-height: 50px;
    word-wrap: break-word;
    word-break: break-all;
    overflow: hidden;
    color: rgba(72, 72, 102, 0.64);
    font-weight: 500;
}
.only_dialog /deep/ .el-dialog{
    height: 87vh !important;
}
.formbox_upload{
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    font-size: 16px;
    margin-top: 20px;
    .formbox_upload_file{
        margin-left: 40px;
        margin-right: 10px;
    }
    .formbox_upload_filename{
        
    }
    .formbox_upload_btn_upload{
        margin-right: 15px;
    }
    .formbox_upload_btn{
        margin-right: 30px;
        display: flex;
    }
}
.secondline{
        border-top:1px solid #E8E8F5;
        border-bottom:1px solid #E8E8F5;
        border-bottom:none;
    }
/deep/ .savebox .el-form-item__content{
    margin-left: 10px!important;
}
 /deep/.nopaddingbox .el-tabs__content{
   padding-left: 0!important;
}
/deep/ .el-table th.el-table__cell>.cell {
    display: inline-block;
    box-sizing: border-box;
    position: relative;
    vertical-align: middle;
    padding-left: 10px;
    padding-right: 10px;
    width: 100%;
    padding-left: 32px!important;
}
/deep/ .el-input--small .el-input__icon {
    line-height: 52px;
}
/deep/ .el-checkbox__label {
    font-size: 13px;
    color: rgba(72, 72, 102, 0.64);
}
/deep/ .el-form-item__label{
 text-align: left;
}
</style>
<script> 
import xzbutton from "../../components/XzButton.vue"; 
import {encryptCBC,decryptCBC} from '@/commonFunction/des.js' 
import { updateFile} from '@/lib'
import { system } from '@/api/system.js'
import othersysset from '@/pages/system/othersysset'//其他设置
import businessset from '@/pages/system/businessset'//其他设置
import sysmonitoring from '@/pages/system/sysmonitoring'//系统监控
import sysconfigcopy from '@/pages/system/sysconfigcopy'//配置备份
import safetyset from '@/pages/system/safetyset'//安全设置
import $ from 'jquery'
export default({
    name:'systemsetting',
    components: {
    	xzbutton,
        othersysset,//其他设置
        businessset,
        sysmonitoring,//系统监控
        sysconfigcopy,//配置备份
        safetyset,//安全设置
  	},
    data(){ 

        let threshold = (rule, value, callback)=>{
            if (value < 0 || value > 100) return callback(new Error(`只能在0-100之间`));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()

        }
        let memory = (rule, value, callback)=>{
            if (value < 0 || value > 100) return callback(new Error(`只能在0-100之间`));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()
        }
        let disk = (rule, value, callback)=>{
            if (value < 0 || value > 100) return callback(new Error(`只能在0-100之间`));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()

        }
        let flow = (rule, value, callback)=>{
            if (value < 0 || value > 10000) return callback(new Error(`只能在0-10000之间`));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()

        }
 
    	return{
            systemupdate:{
                latest_sys_version:'',
                find_time:'',
                current_sys_version:'',
                update_time:''
            },
            system_uid:'',
            dialogSysVisible:false,
            role:0,  
    		activeName:'tab6', 
            formBasics:{
                ip:'',
                port:'',
                cpu:'',
                memories:'',
            },
            formagent:{
                agreement:'',
                ip:'',
                port:'',
            },
            formcheck:{
                port:'',
                useragent:'',
                thread:'',
                times:'',
                cookie:'',
            },
            formMessage:{
                remindtype:'0',
                email:'',
                threshold:'1',
            },
            tableData:[],
            totalpage:0,
            currentpage:1,
            dialogFormVisible:false,
            formJurisdiction:{
                id:0,
                user:true,
            },
            page_num:1,
            myKey: '4dogs.cn',
            productinfo:{
                name:'',
                engine_version:'',
                software_version:'',
                version_type:'',
                time:'',
                login_url:'',
            },
            upgrade1:false,
            upgrade2:false,
            shutdown:false,
            fastrestart:false,
            systemrestart:false,
            shutdownDialogVisible:false,
            fastrestartDialogVisible:false,
            systemrestartDialogVisible:false,
            systemfileName:'', //系统上传文件名
            systemfile:null,
            systemupgradeDialogVisible:false,
            restart:'0',
            isfirst:true,
            knowledgeDialogVisible:false,
            knowledgeisfirst:true,
            knowledgerestart:0,
            fullscreenLoading:false,
            systemSTR:'',
            // warningform:{
            //     cpu_threshold: 10,
            //     disk_threshold: 0,
            //     flow_threshold:0,
            //     memory_threshold: 0,
            // },
            setUp:{
                // CPU告警阈值
                cpu_threshold:[
                    { required: true, message: 'CPU告警阈值不能为空', trigger: ['blur','change'] }, 
                    { validator:threshold, trigger: ['blur','change'] }
                ],
                // 内存告警阈值
                memory_threshold:[
                    { required: true, message: '内存告警阈值不能为空', trigger: ['blur','change'] }, 
                    { validator:memory, trigger: ['blur','change'] }
                ],
                //  硬盘报警
                disk_threshold:[
                    { required: true, message: '硬盘报警阈值不能为空', trigger: ['blur','change'] }, 
                    { validator:disk, trigger: ['blur','change'] }
                ],
                //  流量告警
                flow_threshold:[
                    { required: true, message: '流量告警阈值不能为空', trigger: ['blur','change'] }, 
                    { validator:flow, trigger: ['blur','change'] }
                ],
                
            },
            keyTableData:[//API秘钥
                // {
                //     user: 'admin',
                //     authentication: 'shdkfdljgosdj',
                //     date: '2021-4-26'
                // }, 
            ],
            dialogKeyVisible:false,
            userKey:'',
            userNames:['test1','Tom'],

            loading:false,
            sys_version:'',
            backupstableData:[],
            alldelvisible:false,
            systemrestoreloading:false,
            offupdateloading:false,
            nonlineupdateeloading:false,
            thisVersion:'',
            newVersion:'',
            fileDialog:false,   // 文件授权dialog
            uploadFile:null,
            filename:"请导入授权文件，再点击授权按钮授权",
            systemLoading:false,
            currentAxiosSource: null, // 当前请求的axios
            onlineUpdate: {
                plan_time1: '1',
                plan_time2: 1,
                plan_time3: '1',
                plan_time4: '10:00:00',
                isopen:false
            },
            newWin: null, // 新窗口对象
            updateStatus: '',
            uploadFileName: '', // 上传文件名
            needUpdate: false, // 是否需要更新
            timer2:'',
            // updateStatusText: '',
            otherconfig:{
                system_threshold:false,
                syslog:false,
                whitelist:false,
            },
    	}
    }, 
    created (){
        this.$store.commit('setPercentage', 0)
        this.$store.state.activefirstMenu="/systemsetting"; 
        this.role = decryptCBC(localStorage.getItem('role'),this.myKey);
    },
    mounted () {
        this.get_product_info();
        // this.get_email_config();
    },
    beforeDestroy(){
       clearTimeout(this.timer2) 
    },
    methods:{
        handleClick(tab, event){
            // console.log(tab.name);
            if(tab.name == 'tab1'){ //基本
                this.get_basic_info()
            }
            else if(tab.name == 'tab2'){  //业务设置
                
            }
            // else if(tab.name == 'tab3'){ //检测
            //     this.get_check_config();
            // }
            // else if(tab.name == 'tab4'){ //提醒消息
            //     this.get_notify_settings();
            // }
            else if(tab.name == 'tab5'){
                this.get_user_list();
            }
            else if(tab.name == 'tab6'){
                this.get_product_info();            
                 
            }
            else if(tab.name == 'tab7'){ //关机重启
                
            }
            else if(tab.name == 'tab8'){ //告警设置
                // this.getwarninginfo();
            }
            else if(tab.name == 'tab9'){  //安全设置
                this.getsecurityinfo();
            }
            else if(tab.name == 'tab10'){  //系统备份与还原 
                this.getVersion2();
            }
            else if(tab.name == 'tab11'){  //API秘钥
                this.fnShowTable();
            }
            // else if(tab.name == 'tab12'){  //系统 监控
            //     this.getvul_database_trend();
            //     this.echartline2();
            // }
        },
        opentip(){
            this.dialogSysVisible = true;
        } ,
        fnReptileSwitch(val){//web爬虫配置
            this.onlineUpdate.isopen = val;
        },  
        get_basic_info(){ //获得基本信息
            // this.$ajax.get('/systems/information/get_basic_info/',{
            //     params: {}
            // })
            system.getBasicInfo()
            .then((data) => { 
                var dt = data;  
                if(dt.success){ 
                    this.formBasics.ip = dt.basic_info.ip_address;
                    this.formBasics.port = dt.basic_info.port;
                    this.formBasics.cpu = dt.basic_info.cpu;
                    this.formBasics.memories = dt.basic_info.memory;
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
            })
            .catch((error) => {
                console.log(error);
            })
         

        },
        
        handleSubmitAgent(){ //保存代理设置

        },
        get_check_config(){ //获得检测设置
            this.$ajax.get('/sysconfig/get_check_config/',{
                params: {}
            })
            .then((data) => { 
                var dt = data.data;  
                if(dt.success){  
                    this.formcheck.port = dt.advance_settings.scan_port;
                    this.formcheck.useragent = dt.advance_settings.user_agent;
                    this.formcheck.thread = dt.advance_settings.thread_number;
                    this.formcheck.times = dt.advance_settings.timeout;
                    this.formcheck.cookie = dt.advance_settings.cookie;
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
            })
            .catch((error) => {
                console.log(error);
            })
        },
        handleSubmitCheck(){ //保存 检测设置
            this.$ajax({
                method:'post',
                url:'/sysconfig/save_check_config/',
                data:this.qs.stringify({
                    scan_port: this.formcheck.port,  
                    user_agent:this.formcheck.useragent,
                    thread_number:this.formcheck.thread,
                    timeout:this.formcheck.times,
                    cookie:this.formcheck.cookie,
                })
            })
            .then(data => { 
                var dt = data.data;  
                if(dt.success){ 
                    this.$message({
                        message:'保存检测设置成功！',
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
                
            })
            .catch(data=>{
                console.log(data); //错误信息
            });
        },
        handleRecoveryDefault(){ //恢复默认检测设置
            this.$ajax.get('/sysconfig/recovery_default_config/',{
                params: {}
            })
            .then((data) => { 
                var dt = data.data;  
                if(dt.success){  
                    this.formcheck.port = dt.advance_settings.scan_port;
                    this.formcheck.useragent = dt.advance_settings.user_agent;
                    this.formcheck.thread = dt.advance_settings.thread_number;
                    this.formcheck.times = dt.advance_settings.timeout;
                    this.formcheck.cookie = dt.advance_settings.cookie;
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
            })
            .catch((error) => {
                console.log(error);
            })
        },
        get_notify_settings(){ //获取提醒消息设置
            this.$ajax.get('/sysconfig/get_notify_settings/',{
                params: {}
            })
            .then((data) => { 
                var dt = data.data;  
                // console.log(dt);
                if(dt.success){  
                    this.formMessage.remindtype = dt.notify_way+'';
                    this.formMessage.email = dt.notify_target;
                    this.formMessage.threshold = dt.notify_threshold;
                    
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
            })
            .catch((error) => {
                console.log(error);
            })
        },
        changeremindtype(val){
            
                // this.formMessage.threshold = '';
                // this.formMessage.email = '';
            
        },
        handleSubmitMessage(formName){ //保存 消息设置
            // this.$refs[formName].validate((valid) => {
            //     if (valid) {
                   this.$ajax({
                        method:'post',
                        url:'/sysconfig/save_notify_settings/',
                        data:this.qs.stringify({
                            notify_way:this.formMessage.remindtype,
                            notify_threshold:this.formMessage.remindtype == 1 ? this.formMessage.threshold :'',
                            notify_target:this.formMessage.remindtype == 1 ? this.formMessage.email :'',
                        })
                    })
                    .then(res => { 
                        var dt = res.data;  
                        if(dt.success){ 
                            this.$message({
                                message:'保存消息设置成功！',
                                type: 'success'
                            });
                        }else{
                            this.$message({
                                message:dt.error,
                                type: 'error'
                            });
                        }  
                        
                    })
                    .catch(error=>{
                        console.log(error); //错误信息
                    });
                // } else {
                //     console.log('error submit!!');
                //     return false;
                // }
        // });
            
        },
        get_user_list(){ //获得用户列表
            this.$ajax.get('/sysconfig/get_user_list/',{
                params: {
                    page_num:this.page_num,   
                }
            })
            .then((data) => { 
                var dt = data.data;  
                // console.log(dt);
                if(dt.success){  
                    this.tableData = dt.user_list;
                    this.totalpage = dt.all_article_num ;
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
            })
            .catch((error) => {
                console.log(error);
            })
        },
        currentchange(t){
            this.page_num = t; 
            this.get_user_list();
            this.currentpage = t;
        },
        fnset(id,is_active){
            this.formJurisdiction.id = id;
            this.formJurisdiction.user = is_active;
            this.dialogFormVisible = true;
        },
        submitForm(){ //用户设置权限
            this.$ajax({
                method:'post',
                url:'/sysconfig/activate_user/',
                data:this.qs.stringify({
                    is_active:this.formJurisdiction.user ? 1:0,
                    user_id : this.formJurisdiction.id 
                })
            })
            .then(res => { 
                var dt = res.data;  
                if(dt.success){ 
                    this.dialogFormVisible = false;
                    this.$message({
                        message:'用户权限设置成功！',
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
                
            })
            .catch(error=>{
                console.log(error); //错误信息
            });
        },
        cancelform(){
            this.dialogFormVisible = false;
        },
        cancelKey(){
            this.dialogKeyVisible = false;
        },
        async get_product_info(){
            
            const dt = await system.getversion();
            console.log('system', dt)
            if(dt.code == 200){ 
                this.productinfo.name = dt.data.productName;
                // this.productinfo.engine_version = dt.data.productID;
                this.productinfo.software_version = dt.data.productID;
                // this.productinfo.version_type = dt.data.type;
                this.productinfo.time = dt.data.use_time;
                // this.productinfo.login_url =  dt.data.login_url;
                   this.productinfo.authTime = dt.data.authTime;
                    this.productinfo.authDays = dt.data.authDays;
                    this.productinfo.leftDays = dt.data.leftDays;
            }else{
                // this.$message({
                //     message:dt.error,
                //     type: 'error'
                // });
            }  
        },
        // 升级还原获取system_uid
        async getsystem_uid(){ 
            const dt = await system.getSystemUid();
            if(dt.success){ 
                this.system_uid = dt.system_uid;
                this.$store.state.system_uid = dt.system_uid;
            }else{
                this.$message({
                    message:dt.error,
                    type: 'error'
                });
            }  
        },
        shutdownhandleClose(){
            this.shutdownDialogVisible = false;
            this.shutdown =false;
        },
        shutdownsubmitForm(){ //系统关机
            this.$ajax({
                method:'post',
                url:'/systems/interfaces/shutdown/',
                data:this.qs.stringify({})
            })
            .then(res => { 
                var dt = res.data;  
                if(dt.success){ 
                    this.shutdownDialogVisible = false;
                    this.$message({
                        message:'系统关机成功！ 系统将在5秒后关机',
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
                
            })
            .catch(error=>{
                console.log(error); //错误信息
            });
        },
        chkshutdown(){
            if(this.shutdown){
                this.shutdownDialogVisible = true;
            }
            
        }, 
        chkfastrestart(){ 
            if(this.fastrestart){
                this.fastrestartDialogVisible = true;
            }
        },
        fastrestarthandleClose(){
            this.fastrestartDialogVisible = false;
            this.fastrestart = false;
        },
        fastrestartsubmitForm(){ //快速重启
            this.$ajax({
                method:'post',
                url:'/systems/interfaces/restart/',
                data:this.qs.stringify({})
            })
            .then(res => { 
                var dt = res.data;  
                if(dt.success){ 
                    this.fastrestartDialogVisible = false;
                    this.$message({
                        message:'快速重启成功！系统将在5秒后重启',
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
                
            })
            .catch(error=>{
                console.log(error); //错误信息
            });
        },
        systemrestarthandleClose(){
            this.systemrestartDialogVisible  = false;
            this.systemrestart = false;
        },
        systemrestartsubmitForm(){ //系统重启
            this.$ajax({
                method:'post',
                url:'/systems/interfaces/reboot/',
                data:this.qs.stringify({})
            })
            .then(res => { 
                var dt = res.data;  
                if(dt.success){ 
                    this.systemrestartDialogVisible = false;
                    this.$message({
                        message:'系统重启成功！系统将在5秒后重启',
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
                
            })
            .catch(error=>{
                console.log(error); //错误信息
            });
        },
        chksystemrestart(){
            if(this.systemrestart){
                this.systemrestartDialogVisible  = true;
            }
           
        },
        clickupload(){ //上传文件
             document.querySelector('.btnUploadID').click();
        },
        changeuploaID(e){
            let deviceFile = e.target.files; 
            for(let i=0;i<deviceFile.length;i++){ 
                this.systemfileName = deviceFile[i].name;
                this.systemfile =  deviceFile[i]; 
            } 
        },
        systemupgradehandleClose(){
            this.systemupgradeDialogVisible = false;
            this.upgrade1 = false;
            this.isfirst = true;
            this.systemfile = null;
            this.restart = '0';
            this.systemfileName = '';
        },
        systemupgradesubmitForm(){ //系统升级
            let formData = new FormData(); 
            formData.append('file',this.systemfile);
            formData.append('reboot',this.restart); //0/1 是否重启
            let config = {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            } 
            this.$ajax.post('/systems/interfaces/system_upgrade/',
                formData, config
            ).then((res)=>{   
                var dt = res.data;  
                if(dt.success){ 
                    this.systemupgradeDialogVisible = false;
                    this.upgrade1 = false;
                    this.isfirst = true;
                    this.systemfile = null;
                    this.restart = '0';
                    this.systemfileName = '';
                    this.$message({
                        message:'升级系统成功！系统将在5秒后升级',
                        type: 'success'
                    });
                    this.fullscreenLoading = true;
                    let timer = setInterval(() => {
                        this.$ajax.get('/systems/interfaces/system_upgrade_status/',{
                            params: {}
                        })
                        .then((data) => { 
                            var dt = data.data;  
                            if(dt.success){ 
                                this.fullscreenLoading = false;
                                clearInterval(timer)
                            } 
                        })
                        .catch((error) => {
                            console.log(error);
                        })
                    }, 3000);
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
            },()=>{ 
                console.log('error'); //错误信息
            }); 
            
        },
        systemupgrade(){
            if(this.upgrade1){
                this.systemupgradeDialogVisible = true;
            }
            // if(this.upgrade2){
            //     this.knowledgeDialogVisible = true;
            // }
        },
        systemupgradeprestep(){ //上一步
            this.isfirst = true;
        },
        systemupgradenext(){ //下一步
            this.isfirst = false;
        },
        knowledgeupgrade(){ //知识库升级
            this.knowledgeDialogVisible = true;
        },
        clickupload1(){ //上传文件
             document.querySelector('.btnUploadID1').click();
        },
        changeuploaID1(e){
            let deviceFile = e.target.files; 
            for(let i=0;i<deviceFile.length;i++){ 
                this.systemfileName = deviceFile[i].name;
                this.systemfile =  deviceFile[i]; 
            } 
        },
        knowledgeprestep(){
            this.knowledgeisfirst = true;
        },
        knowledgenext(){
            this.knowledgeisfirst = false;
        },
        knowledgesubmitForm(){ //确定知识库升级
            let formData = new FormData(); 
            formData.append('file',this.systemfile);
            formData.append('reboot',this.knowledgerestart); //0/1 是否重启
            let config = {
                headers: {
                    'Content-Type': 'multipart/form-data'
                }
            } 
            this.$ajax.post('/systems/interfaces/knowledge_upgrade/',
                formData, config
            ).then((res)=>{   
                var dt = res.data;  
                if(dt.success){ 
                    this.knowledgeDialogVisible = false;
                    this.upgrade2 = false;
                    this.knowledgeisfirst = true;
                    this.systemfile = null;
                    this.knowledgerestart = '0';
                    this.systemfileName = '';
                    this.$message({
                        message:'升级知识库成功！系统将在5秒后升级',
                        type: 'success'
                    });
                    this.fullscreenLoading = true;
                    let timer = setInterval(() => {
                        this.$ajax.get('/systems/interfaces/knowledge_upgrade_status/',{
                            params: {}
                        })
                        .then((data) => { 
                            var dt = data.data;  
                            if(dt.success){ 
                                this.fullscreenLoading = false;
                                clearInterval(timer)
                            } 
                        })
                        .catch((error) => {
                            console.log(error);
                        })
                    }, 3000);
                }else{
                    this.$message({
                        message:dt.error,
                        type: 'error'
                    });
                }  
            },()=>{ 
                console.log('error'); //错误信息
            }); 
        },
        knowledgehandleClose(){
            this.knowledgeDialogVisible = false;
            this.upgrade2 = false;
            this.knowledgeisfirst = true;
            this.systemfile = null;
            this.knowledgerestart = '0';
            this.systemfileName = '';
        },
        // btnSaveWarning(){ //保存告警设置 
        //     this.$refs.warningform.validate( async valid=>{
        //         if(!valid) return ;
           
        //         const res = await system.saveSystemWarn({
        //                     cpu_threshold: this.warningform.cpu_threshold,
        //                     disk_threshold: this.warningform.disk_threshold,
        //                     flow_threshold:parseInt(this.warningform.flow_threshold)*1024,
        //                     memory_threshold: this.warningform.memory_threshold,
        //                 });
        //         if(res.success){
        //             this.$message({
        //                 message:res.msg,
        //                 type: 'success'
        //             });
        //         }else{
        //             this.$message({
        //                 message:res.msg,
        //                 type: 'error'
        //             });
        //         }
        //     })
        // },

        // async getwarninginfo(){ //获取 告警设置 信息
            
        //     const res = await system.getSysteminfo();
        //     if(res.success){  
        //         this.warningform.cpu_threshold = res.warning.cpu_threshold; 
        //         this.warningform.disk_threshold =res.warning.disk_threshold; 
        //         this.warningform.flow_threshold =res.warning.flow_threshold; 
        //         this.warningform.memory_threshold =res.warning.memory_threshold; 
        //     }else{
        //         this.$message({
        //             message:res.error,
        //             type: 'error'
        //         });
        //     }  

        // },
        async getsecurityinfo(){ //获取安全设置信息
            // this.$ajax.get('/systems/securities/initial/settings/',{
            //     params: {}
            // })
            // .then((data) => { 
            //     var res = data.data;   
            //     if(res.success){ 
            //        this.securityform = res.security;
            //     }else{
            //         this.$message({
            //             message:res.msg,
            //             type: 'error'
            //         });
            //     }  
            // })
            // .catch((error) => { })

            const res = await system.getSysteminfo();
            if(res.success){ 
                this.securityform = res.security;
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }  
        },
        getVersion(){ //获得系统版本 
            this.$ajax.get('/systems/information/version/',{
                params: {}
            }).then(dt=>{
                let res = dt.data;
                if(res.success){ 
                    this.systemupdate.sys_version = res.sys_version; 
                    this.systemupdate.vuln_version = res.vuln_version;
                    this.systemupdate.starttime = res.use_time;
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            }).catch(err=>{})
        },  
        getVersion2(){ //获得s升级还原 
            this.$ajax.get('/systems/interfaces/version/',{
                params: {}
            }).then(dt=>{
                let res = dt.data;
                if(res.success){ 
                    this.systemupdate = res.result; 
                    // this.systemupdate.vuln_version = res.vuln_version;
                    // this.systemupdate.starttime = res.use_time;
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            }).catch(err=>{})
        }, 
        getsystemreduction(){ //系统还原
            this.systemrestoreloading = true;
            this.alldelvisible = false;
            this.$ajax({
                headers: {
                    'Retry-After': '503'
                },
                method:'post',
                url:'/systems/interfaces/restore/settings/',
                data:this.qs.stringify({})
            })
            .then(dt=>{
                let res = dt.data;
                if(res.success){
                    this.systemrestoreloading = false;
                    this.$message({
                        message:res.msg,
                        type: 'success'
                    });
                    this.$router.push({ path: '/login' });  
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
            }).catch(err=>{})
        }, 
        btnoffupdate(){ //离线升级
            this.dialogSysVisible = true;
            this.getsystem_uid()
        },
        // 确定上传离线更新包
        sureUpload () {
            document.querySelector(".btnUploadID2").click();
        },
        getoffuploadStatus(){
            this.$ajax.get('/systems/interfaces/off_upload/status/',{
                params: {
                    system_uid:this.system_uid
                }
            }).then(res=>{
                let dt = res.data;
                this.updateStatus = dt.status
                console.log(dt.success,1827);
                if(dt.success){ 
                    switch (dt.status) {
                        case 'running':
                            this.getoffuploadStatus() // runing时，持续请求该路由，直到finish结束
                            break;
                        case 'finish':
                            setTimeout(() => {
                                if (this.needUpdate) {
                                    this.$router.push({
                                        path: "/systemtip"
                                    });
                                }
                            }, 4000)
                            // this.getOfflineStatus(); // 去访问更新状态路由
                            break;
                        case 'fail':
                            this.$message.error('更新失败，请联系管理员')
                            // this.newWin.close()
                            break;
                        }
                }else{
                    // this.newWin.close() // 异常关闭
                    this.$message({
                        message:res.data.error,
                        type: 'error'
                    });
                }
            }).catch(err=>{
                this.$message.error(err.data.err)
                // this.newWin.close() // 异常关闭
            })
        },
        getOfflineStatus (time) {
            // this.systemLoading = true
            this.$ajax.get('/systems/interfaces/off_update/status/')
                .then((res)=>{   
                    var dt = res.data;
                    if(dt.success){
                     
                         switch (dt.status) {

                            case 'running':
                                this.getOfflineStatus() // 持续请求该路由，直到‘finish'结束请求
                                break
                            case 'finish':

                                // if (time === 'again') {
                                //     this.systemLoading = false
                                //     this.$router.push({
                                //         path: "/login",
                                //         query: {}
                                //     });
                                // } else {
                                //     setTimeout(() => {
                                //         this.getOfflineStatus('again')
                                //     }, 10000)
                                // }
                                // 休息20s,直接访问登录页
                                setTimeout(() => {
                                    this.$router.push({
                                        path: "/login",
                                        query: {}
                                    });
                                }, 30000)
                                breaktoNowsystem
                            case 'fail':
                                this.$message.error('更新失败，请联系管理员')
                                this.newWin.close() // 直接结束更新
                                break
                        }
                    }
                })
                .catch(err=>{
                    setTimeout(() => {
                        this.getOfflineStatus()
                    }, 5000)
                })
        },
        // changeuploaID2(e){
        async changeuploaID2(e){
            this.needUpdate = true
            let deviceFile = e.target.files;  
            let systemfileName = deviceFile[0].name;
            let systemfile =  deviceFile[0]; 
            this.uploadFileName = systemfileName
            let formData = new FormData(); 
            formData.append('file',systemfile);
            formData.append('system_uid',this.system_uid);
            this.dialogSysVisible = false;
            this.updateStatus = 'running'
            const { data } = await updateFile(formData)
            $('#input-file-ID2').val('')
           
             if (data.success) {
              this.timer2 = setTimeout(()=>{
              this.getoffuploadStatus();
            },1000)
             }else{
              this.$message.error(data.error)
             }
        },
        btnonlineupdate(){ //zaixian
            this.nonlineupdateeloading = true;
            this.$ajax({ 
                method:'post',
                url:'/systems/interfaces/online/update/',
                data:this.qs.stringify({ })
            })
            .then(dt=>{
                let res = dt.data;
                if(res.success){ 
                    
                    this.$message({
                        message:'立即更新成功,5秒后重启系统',
                        type: 'success'
                    }); 
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
                this.nonlineupdateeloading = false;
            }).catch(err=>{
                this.nonlineupdateeloading = false;
            })
        },
        testupdate(){
            this.offupdateloading = true;
            this.$ajax.get('/systems/interfaces/version/detect/',{
                params: {}
            }).then(dt=>{
                let res = dt.data;
                if(res.success){  
                    
                    this.newVersion = res.msg;
                    this.$message({
                        message:res.msg,
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:res.error,
                        type: 'error'
                    });
                }
                this.offupdateloading = false;
            }).catch(err=>{
                this.offupdateloading = false;
            })
        },
        handleUploadChange(file){
            this.$refs.uploadRef.clearFiles();
            this.uploadFile = file;
            this.filename = file.name;
        },
        async handleGetData(){ 
            const res = await  system.getAuthorinfo();
            const blob = new Blob([res],{
                type:"application/octet-stream"
            });
            let filename = sessionStorage.filename || 'feature.pem';
            const url =  window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.download = filename;
            a.href = url;
            a.click();
            sessionStorage.filename = ''

        },
        //  文件上传
      async handleFileUpload(){
          const res =   await system.getSystemAuthSave({
                authCode:this.systemSTR
            });
            if(res.code == 200){
                this.$message.success('授权成功');
                this.fileDialog = false;
                this.get_product_info();
                this.systemSTR = ''
            }else{
                this.$message.error(res.msg);
                this.fileDialog = false;
                this.systemSTR = ''
            }
            // if(!this.uploadFile) return this.$message.warning('请先上传文件!')
            // const formData = new FormData();
            // formData.append('cipher_text',this.uploadFile.raw)
            // this.$ajax({
            //     method:"POST",
            //     url:"/systems/auth/upload/",
            //     data:formData
            // })
            // .then(res =>{
            //     if(!res.data.success){
            //         this.$message.error(res.data.msg)
            //     }else{
            //         this.$message.success(res.data.msg);
            //         this.fileDialog = false;
            //         this.uploadFile = null;
            //         this.filename = '请选择要上传的文件';
            //         this.get_product_info();
            //     }
            // })
            // .catch(err =>{})

        },

        // API秘钥
        async fnShowTable(){//展示列表 
            const res = await system.getTokens();
            if(res.success){  
                this.keyTableData = res.token_list;
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        },
        async generateKey(){//获取用户名
            this.dialogKeyVisible = true;

            const res = await system.getusernames();
            if(res.success){  
                this.userNames = res.usernames;
                if(this.userNames){
                    this.userKey = this.userNames[0]
                }else{
                    this.userKey = '';
                }
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        },
        async submitKey(){//生成秘钥
            let _userName = this.userKey; 
            const res = await system.submitgenerateToken({ username:_userName });
            if(res.success){  
                this.dialogKeyVisible = false;
                this.$message({
                    message:res.msg,
                    type: 'success'
                });
                this.fnShowTable();
                
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }

        },
        toNowsystem(){
            const routeData = this.$router.resolve({
                    name: 'nowsystip',
                })
               window.open(routeData.href, '_blank')
            // window.open('/nowsystip')

        }
    }
})
 
</script>
