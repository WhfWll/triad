/* 
    此页是系统配置
 */
<template>
    <div>
        <div class="systembox context_box_bg">
            <el-tabs v-model="activeName">
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
                <el-tab-pane label="系统授权" name="tabAuth">
                    <div v-if="activeName === 'tabAuth'" class="tabsbox auth-tabs">
                        <div class="section-header">
                            <i class="el-icon-key"></i>
                            <span class="section-title">授权信息</span>
                        </div>
                        <div class="info-row">
                            <span class="info-label">产品名称</span>
                            <span class="info-value">{{ productinfo.name || '未获取' }}</span>
                        </div>
                        <div class="info-row">
                            <span class="info-label">特征码</span>
                            <span class="info-value">{{ productinfo.software_version || '未获取' }}</span>
                        </div>
                        <div class="info-row">
                            <span class="info-label">软件版本</span>
                            <span class="info-value">{{ productinfo.software_display_version || '未获取' }}</span>
                        </div>
                        <div class="info-row">
                            <span class="info-label">授权时间</span>
                            <span class="info-value">{{ productinfo.authTime || '未授权' }}</span>
                        </div>
                        <div class="info-row">
                            <span class="info-label">授权天数</span>
                            <span class="info-value">{{ formatAuthDays(productinfo.authDays) }}</span>
                        </div>
                        <div class="info-row">
                            <span class="info-label">剩余天数</span>
                            <span class="info-value" :class="{'danger-text': productinfo.leftDays <= 30}">
                                {{ formatAuthDays(productinfo.leftDays) }}
                            </span>
                        </div>

                        <div class="auth-divider"></div>

                        <div class="section-header">
                            <i class="el-icon-lock"></i>
                            <span class="section-title">授权操作</span>
                        </div>
                        <el-form label-width="120px">
                            <el-form-item label="系统授权码">
                                <el-input 
                                    type="textarea" 
                                    :rows="8" 
                                    v-model="systemSTR" 
                                    placeholder="请输入系统授权码"
                                    maxlength="50000"
                                    class="auth-code-input"
                                ></el-input>
                            </el-form-item>
                            <el-form-item>
                                <el-button 
                                    type="primary" 
                                    @click="handleFileUpload" 
                                    :loading="authLoading"
                                    size="medium"
                                >
                                    立即授权
                                </el-button>
                                <el-button 
                                    @click="systemSTR = ''" 
                                    size="medium"
                                >
                                    清空
                                </el-button>
                            </el-form-item>
                        </el-form>


                    </div>
                </el-tab-pane>
                <el-tab-pane label="系统监控" name="tab12">
                    <sysmonitoring v-if="activeName === 'tab12'"/>
                </el-tab-pane>
                <el-tab-pane label="安全检查配置" name="tabSec">
                    <security-check-config v-if="activeName === 'tabSec'" />
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
        <el-dialog title="系统授权" width='1184px' top="5vh" :visible.sync="fileDialog" :close-on-click-modal="false"
            :show-close="false" custom-class="auth_dialog_custom only_dialog">
            <div class="dialog_b_btn">
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
        <!-- <el-dialog title="文件授权" width='1184px' :visible.sync="fileDialog" :close-on-click-modal="false"
            :show-close="false" class="only_dialog">
            <div class="dialog_b_btn only_dialog">
                <el-button size="small" @click="fileDialog = !fileDialog">取消</el-button>
            </div>
            <div style="padding:24px" class="formbox">
                <div class="formbox_upload">
                    <div>
                        <span class="formbox_upload_file">授权文件：</span>
                        <span an class="formbox_upload_filename" style="fontSize:15px;opacity:.6;">{{filename}}</span>
                    </div>
                    <div class="formbox_upload_btn">
                        <el-upload :limit="1" :show-file-list="false" :on-change="handleUploadChange"
                            :auto-upload="false" action="https://jsonplaceholder.typicode.com/posts/"
                            class="formbox_upload_btn_upload" ref="uploadRef">
                            <el-button type="primary" size="small">导入</el-button>
                        </el-upload>
                        <el-button type="primary" class="" size="small" @click="handleFileUpload">授权</el-button>
                    </div>
                </div>
            </div>
        </el-dialog> -->


        <el-dialog title="生成秘钥" :visible.sync="dialogKeyVisible" :before-close="cancelKey" :close-on-click-modal="false"
            :show-close="false">
            <div class="dialog_b_btn">
                <el-button size="small" @click="submitKey">生成秘钥</el-button>
                <el-button size="small" @click="cancelKey">取消</el-button>
            </div>
            <div style="padding:24px">
                <label class="dialog_item_label">用户名</label>
                <el-select v-model="userKey" placeholder="请选择">
                    <el-option v-for="(item,index) in userNames" :key="index" :label="item.username" :value="item.username">
                    </el-option>
                </el-select>
            </div>
        </el-dialog>
        <!-- 升级还原确定弹窗 -->
        <el-dialog title="离线更新" :visible.sync="dialogSysVisible" :before-close="cancelKey" :close-on-click-modal="false"
            :show-close="false" custom-class="auto-height-dialog">
            <div class="dialog_b_btn">
                <el-button size="small" @click="sureUpload">确定</el-button>
                <el-button size="small" @click="dialogSysVisible = false">取消</el-button>
            </div>
            <div style="padding:24px">
                确定要上传离线更新文件包吗？
            </div>
        </el-dialog>
        
        <!-- 升级确认弹窗 -->
        <el-dialog title="确认升级" :visible.sync="upgradeConfirmDialogVisible" :close-on-click-modal="false" :show-close="false" custom-class="auto-height-dialog">
             <div style="padding: 24px;">
                <div style="margin-bottom: 20px; font-weight: bold;">检测到升级包，请确认详细信息：</div>
                <div style="line-height: 30px;">
                    <div><label style="width: 100px; display: inline-block;">版本号：</label>{{ upgradeInfo.version }}</div>
                    <div><label style="width: 100px; display: inline-block;">发布时间：</label>{{ upgradeInfo.releaseTime }}</div>
                    <div><label style="width: 100px; display: inline-block;">文件大小：</label>{{ upgradeInfo.size }}</div>
                    <div><label style="width: 100px; display: inline-block;">漏洞类型：</label>{{ upgradeInfo.typeDesc || upgradeInfo.type }}</div>
                    <div><label style="width: 100px; display: inline-block;">升级范围：</label>{{ upgradeInfo.scopeDesc || upgradeInfo.vulnScope }}</div>
                    <div><label style="width: 100px; display: inline-block; vertical-align: top;">更新描述：</label>
                        <div style="display: inline-block; width: calc(100% - 110px);">{{ upgradeInfo.description }}</div>
                    </div>
                </div>
                <div style="margin-top: 20px; color: #F56C6C; font-size: 12px;">
                    <i class="el-icon-warning"></i> 升级过程中系统将暂时不可用，请确认是否立即执行升级？
                </div>
            </div>
            <div slot="footer" class="dialog-footer">
                <el-button @click="upgradeConfirmDialogVisible = false">取消</el-button>
                <el-button type="primary" @click="executeUpgrade" :loading="offupdateloading">确认执行</el-button>
            </div>
        </el-dialog>

      
       
    </div>
</template>
<style scoped lang="less">
@import '../security/css/appsec-tokens.less';

    .tabsbox{
        background: #fff;
        padding: 20px;
        box-sizing: border-box;
        border-radius: 4px;
        box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.1200);
   
    }

    .tabsbox.auth-tabs {
        background: var(--appsec-bg-surface, #1a1a2e);
        border: 1px solid var(--appsec-bg-surface-inset, #16213e);
    }

    .totalbox{
        padding: 25px;
        box-sizing: border-box;
        background: #fff;
        box-shadow: 0px 2px 4px 1px rgba(76,122,227,0.1200);
        border-radius: 4px 4px 4px 4px;
  
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
    /deep/ .el-dialog:not(.ai-model-dialog .el-dialog){
        min-height: 192px !important;
    }
    /deep/ .el-dialog.auto-height-dialog{
        height: auto !important;
    }
    .dialogtxt{
        text-align: center;
        margin-top: 55px;
    }
    /deep/ .el-tabs__item{
        height: 48px;
        line-height: 48px;
        padding: 0 24px;
        color: @appsec-text-body;
        transition: color 0.2s;
        &:hover {
            color: @appsec-text-strong;
        }
    }
    /deep/ .el-tabs__item.is-active{
        color: @appsec-accent;
        font-weight: 600;
    }
    /deep/ .el-tabs__active-bar {
        background-color: @appsec-accent;
    }
    /deep/ .el-tabs__nav-wrap{
        padding: 0 24px; 
    }
    /deep/ .el-tabs__nav-wrap::after{
        background: @appsec-border-subtle;
        height: 1px;
    }
    /deep/ .el-tabs__header{
        margin: 0 0 24px;
    }
    // /deep/ .el-table .cell, .el-table th div, .el-table--border td:first-child .cell, .el-table--border th:first-child .cell{
    //     padding-left:32px;
    // }
    
     
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
    background: @appsec-bg-surface;
    border-radius: @appsec-radius-md;
    box-shadow: @appsec-shadow-card;
    border: 1px solid @appsec-border-default;
}
.systembox{
    //
	// background: #fff;
    height: calc(100% - 39px);
    box-sizing: border-box;
    // padding: 24px;
    // box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
    /deep/ .el-tabs{
        // min-height: 700px;
        height: 100%;
        .el-tabs__content{
            height: calc(100% - 59px);
            >div{
                height: 100%;
                >div{
                    height: 100%;
                }
            }
        }
    } 
    /deep/ .el-tabs__item.is-top.is-active{
        // background: #fff;
    }
    /deep/ .el-tabs__header{
        margin: 0 0 15px;
        background: @appsec-bg-surface;
        border-radius: @appsec-radius-md;
        box-shadow: @appsec-shadow-card;
        border: 1px solid @appsec-border-default;
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
        margin-bottom:0;
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
.systeminfolist >div .restore-btn{
    float: right;
    margin: 10px 20px;
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
    height: auto !important;
}
/deep/ .auth_dialog_custom{
    height: auto !important;
    margin-bottom: 50px;
}
/deep/ .auth_dialog_custom .el-dialog__body{
    height: auto !important;
    overflow: visible !important;
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
    // padding-left: 32px!important;
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
import sysmonitoring from '@/pages/system/sysmonitoring'//系统监控
import SecurityCheckConfig from '@/pages/system/SecurityCheckConfig.vue';
import $ from 'jquery'
export default {
    name:'systemsetting',
    components: {
    	xzbutton,
        sysmonitoring,//系统监控
        SecurityCheckConfig,
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
        //  安全设置相关
        let cycle = (rule, value, callback)=>{
            if(value < 0 || value > 36) callback(new Error('只能输入0-36范围之内'));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()
        }
        let timeout = (rule, value, callback)=>{
            if(value < 0 || value > 120) callback(new Error('只能输入0-120范围之内'));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()
        }
        let user = (rule, value, callback)=>{
            if(value < 0 || value > 10) callback(new Error('只能输入0-10范围之内'));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()
        }
        let admin = (rule, value, callback)=>{
            if(value < 0 || value > 10) callback(new Error('只能输入0-10范围之内'));
            else if (!(/(^[0-9]\d*$)/.test(value))) return callback(new Error('只能输入整数类型'))
            else callback()
        } 
    	return{
            actNum:1,
            inActNum:0,

            systemupdate:{
               latestVersion:'',
                findTime:'',
                currentVersion:'',
                updateTime:'',
                vulVersion:'',
                vulUpdateTime:'',
            },
            system_uid:'',
            dialogSysVisible:false,
            role:0,  
    		activeName:'tabAuth', 
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
                software_display_version:'',
                version_type:'',
                time:'',
                login_url:'',
                authTime:'',
                authDays:0,
                leftDays:0
            },
            authLoading: false,
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
            securityform:{
                password_rank:'',
                password_cycle:'',
                login_timeout:'',
                user_limit:'',
                admin_limit:'',
                ban_time:'',
                password_valid :"",
                expire_unused :''
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
           
            isEditingModel: false,
            editingModelId: null,
           
            // 应用场景配置
            sceneConfig: {
                captchaModel: '', // 验证码识别模型ID
                webStructureModel: '' // 网页结构识别模型ID
            },
            savingSceneConfig: false, // 保存场景配置的loading状态
            safety:{
                // 密码修改周期
                password_cycle:[
                    { required: true, message: '密码修改周期不能为空', trigger: ['blur','change'] }, 
                    { validator:cycle, trigger: ['blur','change'] }
                ],
                // 系统登录超时
                login_timeout:[
                    { required: true, message: '系统登录超时不能为空', trigger: ['blur','change'] },
                    { validator:timeout, trigger: ['blur','change'] }
                ],
                user_limit:[
                    { required: true, message: '普通用户登录限制不能为空', trigger: ['blur','change'] }, 
                    { validator:user, trigger: ['blur','change'] }
                ],
                admin_limit:[
                    { required: true, message: '管理/审计员登录限制不能为空', trigger: ['blur','change'] }, 
                    { validator:admin, trigger: ['blur','change'] }
                ],
            },
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
            systemSTR:'',
            pageSize:10,
            totalpage:0,
			currentpage:1,
            
            
            
            // 场景配置相关数据
            sceneConfig: {
                captchaModel: '',
                webStructureModel: ''
            },
            savingSceneConfig: false,
            
            // 新版离线升级相关
            upgradeConfirmDialogVisible: false,
            upgradeInfo: {
                version: '',
                releaseTime: '',
                size: '',
                description: ''
            }
    	}
    },
    computed: {
        // 获取可用的模型选项（用于下拉选择）
        availableModels() {
            return this.aiModelTableData.map(model => ({
                id: model.id,
                name: model.name,
                modelId: model.modelId,
                platform: model.platform
            }));
        }
    },
    created (){
        this.$store.commit('setPercentage', 0)
        this.$store.state.activefirstMenu="/systemsetting"; 
        this.role = decryptCBC(localStorage.getItem('role'),this.myKey);
    },
    mounted () {
        this.get_product_info();
    },
    watch: {
        activeName: {
            handler(newVal) {
                if (newVal === 'tabAuth') {
                    this.get_product_info();
                }
            },
            immediate: false
        }
    },
    beforeDestroy(){
    //    clearTimeout(this.timer2) 
    },
    methods:{
           formatAuthDays(value) {
            const text = value === undefined || value === null ? '' : String(value).trim()
            if (!text || text === '--') {
                return '--天'
            }
            if (text.includes('过期')) {
                return text
            }
            if (text.endsWith('天')) {
                return text
            }
            if (/^\d+$/.test(text)) {
                return `${text}天`
            }
            return text
            },
           inputChange(val) {
            this.securityform.expire_unused = this.securityform.expire_unused.replace(/[^0-9.]/g, '')
            },
           inputChange2(val) {
            this.securityform.password_valid = this.securityform.password_valid.replace(/[^0-9.]/g, '')
            },
           inputChange3(val) {
            this.securityform.ban_time = this.securityform.ban_time.replace(/[^0-9.]/g, '')
            },
        // 回滚到上次系统版本
        restoreSystem() {
            this.$confirm('确认回滚到上次系统版本？回滚过程中系统将暂时不可用。', '提示', {
                type: 'warning',
                confirmButtonText: '确定',
                cancelButtonText: '取消'
            }).then(() => {
                system.manualRollback({ type: 'SYSTEM' }).then(res => {
                    if (res.code === 200) {
                        this.$message.success('回滚任务已启动');
                        this.startRollbackLoading();
                        this.pollRollbackStatus();
                    } else {
                        this.$message.error(res.msg || '回滚启动失败');
                    }
                }).catch(err => {
                    this.$message.error('请求回滚接口失败');
                    console.error(err);
                });
            }).catch(() => {});
        },
        // 回滚到上次工具库版本
        restoreVul() {
            this.$confirm('确认回滚到上次工具库版本？回滚过程中系统将暂时不可用。', '提示', {
                type: 'warning',
                confirmButtonText: '确定',
                cancelButtonText: '取消'
            }).then(() => {
                system.manualRollback({ type: 'VULN' }).then(res => {
                    if (res.code === 200) {
                        this.$message.success('回滚任务已启动');
                        this.startRollbackLoading();
                        this.pollRollbackStatus();
                    } else {
                        this.$message.error(res.msg || '回滚启动失败');
                    }
                }).catch(err => {
                    this.$message.error('请求回滚接口失败');
                    console.error(err);
                });
            }).catch(() => {});
        },
        startRollbackLoading() {
            if (this.upgradeLoadingInstance) {
                this.upgradeLoadingInstance.close();
                this.upgradeLoadingInstance = null;
            }
            this.upgradeLoadingInstance = this.$loading({
                lock: true,
                text: '正在回滚...',
                spinner: 'el-icon-loading',
                background: 'rgba(0, 0, 0, 0.8)'
            });
        },
        // 轮询回滚状态
        pollRollbackStatus() {
            system.getUpgradeStatus().then(res => {
                if (res.code === 200) {
                    const data = res.data;
                    const state = data.state;
                    const progress = data.progress;
                    const message = data.message || '正在回滚...';
                    if (this.upgradeLoadingInstance) {
                        this.upgradeLoadingInstance.setText(`${message} (${progress}%)`);
                    }
                    const stateStr = typeof state === 'string' ? state.toUpperCase() : state;
                    if (state === 3 || stateStr === 'SUCCESS' || stateStr === 'ROLLBACK_SUCCESS') {
                        if (this.upgradeLoadingInstance) {
                            this.upgradeLoadingInstance.close();
                            this.upgradeLoadingInstance = null;
                        }
                        const backupInfo = data.backup_dir ? `\n备份目录：${data.backup_dir}` : '';
                        const endTimeInfo = data.end_time ? `\n结束时间：${data.end_time}` : '';
                        this.$alert(`回滚完成！请重新登录以确保功能正常。${backupInfo}${endTimeInfo}`, '回滚完成', {
                            confirmButtonText: '重新登录',
                            type: 'success',
                            showClose: false,
                            callback: () => {
                                this.$router.push({ path: "/login" });
                            }
                        });
                    } else if ([4, 5, 6].includes(state) || (typeof state === 'string' && (state.toUpperCase().includes('FAIL') || state.toUpperCase().includes('ERROR')))) {
                        if (this.upgradeLoadingInstance) {
                            this.upgradeLoadingInstance.close();
                            this.upgradeLoadingInstance = null;
                        }
                        let errorMsg = data.error || '回滚失败';
                        if (state === 5) errorMsg += '，正在回滚...';
                        if (state === 6) errorMsg += '，回滚成功';
                        this.$alert(errorMsg, '回滚中止', {
                            confirmButtonText: '确定',
                            type: 'error'
                        });
                    } else {
                        setTimeout(() => {
                            this.pollRollbackStatus();
                        }, 1500);
                    }
                } else {
                    setTimeout(() => {
                        this.pollRollbackStatus();
                    }, 2000);
                }
            }).catch(err => {
                console.error('获取回滚状态失败', err);
                setTimeout(() => {
                    this.pollRollbackStatus();
                }, 2000);
            });
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
                this.productinfo.software_version = dt.data.productID;
                this.productinfo.software_display_version = dt.data.softwareVersion;
                this.productinfo.time = dt.data.use_time;
                this.productinfo.authTime = dt.data.authTime;
                this.productinfo.authDays = dt.data.authDays;
                this.productinfo.leftDays = dt.data.leftDays;
                const authTime = (dt.data.authTime || '').trim();
                const authDays = String(dt.data.authDays || '').trim();
                const leftDays = String(dt.data.leftDays || '').trim();
                const inferredAuthorized = Boolean(dt.data.status) || (
                    authTime &&
                    authTime !== '未授权' &&
                    !leftDays.includes('过期') &&
                    authDays !== '--'
                );
                this.$store.commit('setSystemAuthorized', inferredAuthorized);
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
        
        benSavesecurity(){ //保存安全设置
            this.$refs.securityform.validate( async valid =>{ 
                if(!valid) return ;

            if(this.securityform.password_rank == '低'){
                    this.securityform.password_rank = 1;
                }
                else if(this.securityform.password_rank == '中'){
                    this.securityform.password_rank = 2;
                }
                else if(this.securityform.password_rank == '高'){
                    this.securityform.password_rank = 3;
            }




                const dt = await system.saveSecurity({
                    systemSecurityPasswordRank:this.securityform.password_rank,
                    systemSecurityLoginTimeout:this.securityform.login_timeout,
                    systemSecurityUserLimit:this.securityform.user_limit,
                    systemSecurityAdminLimit :this.securityform.admin_limit,
                    systemSecurityBanTime:this.securityform.ban_time,
                    systemSecurityPasswordValid :this.securityform.password_valid ,
                    systemSecurityExpireUnused :this.securityform.expire_unused  ,

                });
                if(dt.code == 200){  
                    this.$message({
                        message:dt.msg||'保存成功',
                        type: 'success'
                    });
                }else{
                    this.$message({
                        message:dt.msg||'保存失败',
                        type: 'error'
                    });
                }  
            })
        },
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
     

            const res = await system.getSysteminfo();
            if(res.code==  200){ 
                this.securityform = res.data.security;
                if(this.securityform.password_rank == 1){
                    this.securityform.password_rank = '低';
                }
                else if(this.securityform.password_rank == 2){
                    this.securityform.password_rank = '中';
             }else if(this.securityform.password_rank == 3){
                    this.securityform.password_rank = '高';
                }
            }
            else{
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
            this.$ajax.get('/smart/system/systemversion',{
                params: {}
            }).then(dt=>{
                let res = dt.data;
                if(res.code === 200){ 
                    this.systemupdate = res.data; 
                    // this.systemupdate.vuln_version = res.vuln_version;
                    // this.systemupdate.starttime = res.use_time;
                }else{
                    this.$message({
                        message:res.msg,
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
            // this.getsystem_uid()
        },
        // 确定上传离线更新包
        sureUpload () {
            document.querySelector(".btnUploadID2").click();
        },
        getoffuploadStatus(){
            system.getSystemUpgradeStatus().then(res=>{
                let dt = res.data;
                this.updateStatus = dt.status
                // console.log(dt.success,1827);
                if(dt.code === 200){ 
                    switch (dt.status) {
                        case 'running':
                            // this.getoffuploadStatus() // runing时，持续请求该路由，直到finish结束
                             this.timer2 = setTimeout(()=>{
                                this.getoffuploadStatus();
                              },2000)
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
                this.$message.error('获取更新状态失败')
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

            // 计算文件大小
            let sizeStr = '';
            if (systemfile.size < 1024 * 1024) {
                sizeStr = (systemfile.size / 1024).toFixed(2) + ' KB';
            } else {
                sizeStr = (systemfile.size / (1024 * 1024)).toFixed(2) + ' MB';
            }

            this.uploadFileName = systemfileName
            let formData = new FormData(); 
            formData.append('file',systemfile);
            // formData.append('system_uid',this.system_uid);
            this.dialogSysVisible = false;
            
            const loading = this.$loading({
                lock: true,
                text: '正在上传并校验升级包...',
                spinner: 'el-icon-loading',
                background: 'rgba(0, 0, 0, 0.7)'
            });

            try {
                const { data } = await updateFile(formData)
                loading.close();
                $('#input-file-ID2').val('')
            
                if (data.code === 200) {
                     // Store metadata and show confirm dialog
                     this.upgradeInfo = data.data || {}; 
                     if(!this.upgradeInfo.version) this.upgradeInfo.version = "未知版本";
                     
                     // 填充 size
                     if (!this.upgradeInfo.size) {
                         this.upgradeInfo.size = sizeStr;
                     }

                     // 填充 releaseTime
                     if (this.upgradeInfo.buildTime && !this.upgradeInfo.releaseTime) {
                         let timestamp = this.upgradeInfo.buildTime;
                         if (timestamp.toString().length === 10) {
                             timestamp *= 1000;
                         }
                         
                         const date = new Date(timestamp);
                         const year = date.getFullYear();
                         const month = (date.getMonth() + 1).toString().padStart(2, '0');
                         const day = date.getDate().toString().padStart(2, '0');
                         const hour = date.getHours().toString().padStart(2, '0');
                         const minute = date.getMinutes().toString().padStart(2, '0');
                         const second = date.getSeconds().toString().padStart(2, '0');
                         this.upgradeInfo.releaseTime = `${year}-${month}-${day} ${hour}:${minute}:${second}`;
                     }

                     this.upgradeConfirmDialogVisible = true;
                }else{
                    this.$message.error(data.msg || '上传校验失败')
                }
            } catch (error) {
                loading.close();
                this.$message.error('上传失败')
                console.error(error)
            }
        },
        
        // 确认执行升级
        executeUpgrade() {
            this.offupdateloading = true;
            // 从上传响应中获取 filename
            const filename = this.upgradeInfo.filename;
            if (!filename) {
                this.offupdateloading = false;
                this.$message.error('无法获取升级包文件名，请重新上传');
                return;
            }
            system.getSystemoffupgrade({ filename: filename }).then(res => {
                this.offupdateloading = false;
                if (res.code === 200) {
                    this.upgradeConfirmDialogVisible = false;
                    this.$message.success('升级任务已启动');
                    this.updateStatus = 'running';
                    
                    // 开始轮询状态
                    this.pollUpgradeStatus();
                } else {
                    this.$message.error(res.msg || '执行升级失败');
                }
            }).catch(err => {
                this.offupdateloading = false;
                this.$message.error('请求升级接口失败');
                console.error(err);
            });
        },
        
        // 轮询升级状态
        pollUpgradeStatus() {
            // 如果已有 loading 实例，复用；否则创建
            if (!this.upgradeLoadingInstance) {
                this.upgradeLoadingInstance = this.$loading({
                    lock: true,
                    text: '正在初始化升级...',
                    spinner: 'el-icon-loading',
                    background: 'rgba(0, 0, 0, 0.8)'
                });
            }

            system.getUpgradeStatus().then(res => {
                if (res.code === 200) {
                    const data = res.data;
                    const state = data.state; // 兼容数字状态(0-6)和字符串状态("SUCCESS"/"FAIL"等)
                    const progress = data.progress;
                    const message = data.message || '正在处理中...';

                    // 更新 loading 文字
                    if (this.upgradeLoadingInstance) {
                        this.upgradeLoadingInstance.setText(`${message} (${progress}%)`);
                    }

                    // 成功判断: 状态码3 或 字符串"SUCCESS"
                    if (state === 3 || state === 'SUCCESS') {
                        // 成功
                        if (this.upgradeLoadingInstance) {
                            this.upgradeLoadingInstance.close();
                            this.upgradeLoadingInstance = null;
                        }
                        
                        // 弹出模态框告知成功，用户确认后跳转登录
                        this.$alert('系统升级成功！为确保功能正常使用，请重新登录。', '升级完成', {
                            confirmButtonText: '重新登录',
                            type: 'success',
                            showClose: false,
                            callback: action => {
                                // 清除可能存在的旧token等
                                // localStorage.removeItem('token'); // 可选，视具体需求而定
                                this.$router.push({ path: "/login" });
                            }
                        });
                    } 
                    // 失败判断: 状态码4/5/6 或 包含"FAIL"/"ERROR"
                    else if ([4, 5, 6].includes(state) || (typeof state === 'string' && (state.includes('FAIL') || state.includes('ERROR')))) {
                        // 失败/回滚
                        if (this.upgradeLoadingInstance) {
                            this.upgradeLoadingInstance.close();
                            this.upgradeLoadingInstance = null;
                        }
                        let errorMsg = data.error || '升级失败';
                        if (state === 5) errorMsg += '，正在回滚...';
                        if (state === 6) errorMsg += '，回滚成功';
                        
                        this.$alert(errorMsg, '升级中止', {
                            confirmButtonText: '确定',
                            type: 'error'
                        });
                    } else {
                        // 继续轮询
                        setTimeout(() => {
                            this.pollUpgradeStatus();
                        }, 1500);
                    }
                } else {
                    // 接口异常，继续重试
                     setTimeout(() => {
                        this.pollUpgradeStatus();
                    }, 2000);
                }
            }).catch(err => {
                console.error('获取升级状态失败', err);
                // 网络错误也重试
                setTimeout(() => {
                    this.pollUpgradeStatus();
                }, 2000);
            });
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
        //  系统授权
        async handleFileUpload(){
            const res = await system.getSystemAuthSave({
                authCode:this.systemSTR
            });
            if(res.code == 200){
                this.$message.success('授权成功');
                this.fileDialog = false;
                this.$store.commit('setSystemAuthorized', true);
                await this.get_product_info();
                this.systemSTR = ''
            }else{
                this.$message.error(res.msg);
                this.fileDialog = false;
                this.systemSTR = ''
            }
        },

        // API秘钥
        async fnShowTable(){//展示列表 
            let params ={
                page: this.currentpage,
                size: this.pageSize
            }
            const res = await system.getTokens(params);
            if(res.code === 200){  
                this.keyTableData = res.data.list;
                this.totalpage = res.data.total;
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        },
        handleSizeChange(t){
			this.page_num = 1;
            this.pageSize = t;
            this.fnShowTable();
		},
        async generateKey(){//获取用户名
            this.dialogKeyVisible = true;

            const res = await system.getusernames({page:1,size:10000});
            if(res.code === 200){  
                this.userNames = res.data.list;
                if(this.userNames){
                    this.userKey = this.userNames[0].username
                }else{
                    this.userKey = '';
                }
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
            }
        },
        async submitKey(){//生成秘钥
            let _userName = this.userKey; 
            const res = await system.submitgenerateToken({ username:_userName });
            if(res.code === 200){  
                this.dialogKeyVisible = false;
                this.$message({
                    message:'生成成功',
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

        },
        async btnDel(scope){ //
            const res = await system.tokenDel({
                username:scope.row.username,
                token:scope.row.token
            })
            if(res.code == 200){
                this.fnShowTable();
                scope._self.$refs[`popover_id-${scope.row.token}`].doClose();
            }else{
                this.$message({
                    message:res.error,
                    type: 'error'
                });
            }
        },

        

        // 应用场景配置相关方法
        async loadScenarioConfig() {
            try {
                // 暂时使用固定数据，等后端接口完成后再切换到API调用
                // const res = await system.getScenarioConfig();
                
                // 模拟API响应数据
                const mockData = {
                    code: 200,
                    data: {
                        captchaModel: '', // 默认为空，用户可以选择
                        webStructureModel: '' // 默认为空，用户可以选择
                    }
                };
                
                if(mockData.code === 200) {
                    this.sceneConfig = {
                        captchaModel: mockData.data.captchaModel || '',
                        webStructureModel: mockData.data.webStructureModel || ''
                    };
                } else {
                    console.log('获取场景配置失败:', mockData.msg);
                }
            } catch (error) {
                console.error('获取场景配置失败:', error);
            }
        },

        async saveSceneConfig() {
            this.savingSceneConfig = true;
            try {
                // 暂时使用固定数据，等后端接口完成后再切换到API调用
                // const res = await system.saveScenarioConfig(this.sceneConfig);
                
                // 模拟API保存成功响应
                const mockResponse = {
                    code: 200,
                    msg: '保存成功'
                };
                
                // 模拟网络延迟
                await new Promise(resolve => setTimeout(resolve, 500));
                
                if(mockResponse.code === 200) {
                    this.$message.success('场景配置保存成功');
                    console.log('保存的配置:', this.sceneConfig);
                } else {
                    this.$message.error(mockResponse.msg || '场景配置保存失败');
                }
            } catch (error) {
                console.error('保存场景配置失败:', error);
                this.$message.error('场景配置保存失败');
            } finally {
                this.savingSceneConfig = false;
            }
        },

        resetSceneConfig() {
            this.sceneConfig = {
                captchaModel: '',
                webStructureModel: ''
            };
            this.$message.success('配置已重置');
        },

        // 获取当前选中的模型信息
        getCurrentModel(type) {
            const modelId = type === 'captcha' ? this.sceneConfig.captchaModel : this.sceneConfig.webStructureModel;
            return this.aiModelTableData.find(model => model.id == modelId);
        },
    }
}
 
</script>

<style scoped>
/* AI模型配置对话框样式 */
.ai-model-dialog /deep/ .el-dialog {
    min-height: auto !important;
    max-height: 80vh !important;
    height: auto !important;
}

.ai-model-dialog /deep/ .el-dialog__body {
    padding: 0;
    max-height: none;
    overflow-y: visible;
}

.ai-model-dialog /deep/ .el-dialog__header {
    padding: 12px 30px !important;
    background: linear-gradient(135deg, #4c7ae3 0%, #6b8ce8 100%);
    color: white;
    border-bottom: 1px solid #e6f0ff;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.ai-model-dialog .dialog-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
}

.ai-model-dialog .dialog-title {
    color: white;
    font-size: 16px;
    font-weight: 500;
    border-left: 2px solid #fff;
    padding-left: 10px;
    height: 20px;
    line-height: 20px;
}

.ai-model-dialog .dialog-header-buttons {
    display: flex;
    gap: 8px;
}

.ai-model-dialog .header-btn {
    padding: 6px 16px !important;
    border-radius: 3px !important;
    font-size: 13px !important;
    font-weight: 400 !important;
    min-width: 60px !important;
    height: 28px !important;
    line-height: 16px !important;
}

.ai-model-dialog .cancel-btn {
    background: rgba(255, 255, 255, 0.2) !important;
    border: 1px solid rgba(255, 255, 255, 0.4) !important;
    color: white !important;
}

.ai-model-dialog .cancel-btn:hover {
    background: rgba(255, 255, 255, 0.3) !important;
    border-color: rgba(255, 255, 255, 0.6) !important;
}

.ai-model-dialog .confirm-btn {
    background: rgba(255, 255, 255, 0.9) !important;
    border: 1px solid rgba(255, 255, 255, 0.9) !important;
    color: #4c7ae3 !important;
}

.ai-model-dialog .confirm-btn:hover {
    background: white !important;
    border-color: white !important;
    color: #4c7ae3 !important;
}

.ai-model-dialog /deep/ .el-dialog__title {
    color: white;
    font-size: 16px;
    font-weight: 500;
}

.ai-model-dialog /deep/ .el-dialog__headerbtn .el-dialog__close {
    color: white;
    font-size: 18px;
}

.ai-model-dialog /deep/ .el-dialog__headerbtn .el-dialog__close:hover {
    color: #f0f0f0;
}

.ai-model-dialog .el-form {
    background: #ffffff;
    margin: 0;
    padding: 20px 30px 20px 30px;
    min-height: auto;
}

.ai-model-dialog .el-form-item {
    margin-bottom: 18px;
}

.ai-model-dialog .el-form-item:last-child {
    margin-bottom: 0;
}

.ai-model-dialog .el-form-item__label {
    font-weight: 500;
    color: #606266;
    font-size: 14px;
}

.ai-model-dialog .el-input__inner,
.ai-model-dialog .el-select .el-input__inner {
    border-radius: 4px;
    border: 1px solid #dcdfe6;
    transition: all 0.3s;
    font-size: 14px;
    height: 40px;
    line-height: 40px;
}

.ai-model-dialog .el-input__inner:focus,
.ai-model-dialog .el-select .el-input__inner:focus {
    border-color: #4c7ae3;
    box-shadow: 0 0 0 2px rgba(76, 122, 227, 0.1);
}

.ai-model-dialog .el-select .el-input__inner {
    cursor: pointer;
}

/* API地址文本截断样式 */
.api-url-text {
    display: inline-block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: top;
}

/* 平台类型文本样式 */
.platform-text {
    display: inline-block;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    vertical-align: top;
}

.ai-model-dialog .el-switch.is-checked .el-switch__core {
    background-color: #4c7ae3;
    border-color: #4c7ae3;
}

/* 应用场景配置样式 */
.scene-config-container {
    background: #f8f9fa;
    border-radius: 8px;
    padding: 20px;
    border: 1px solid #e9ecef;
}

.scene-config-card {
    background: white;
    border-radius: 6px;
    padding: 16px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    border: 1px solid #e9ecef;
    transition: all 0.3s ease;
}

.scene-config-card:hover {
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
    border-color: #d0d7de;
}

.scene-header {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid #f0f0f0;
}

.scene-title {
    font-weight: 600;
    font-size: 14px;
    color: #303133;
}

.scene-content {
    margin-top: 12px;
}

.current-model {
    margin-top: 8px;
    padding: 8px 12px;
    background: #f0f9ff;
    border-radius: 4px;
    border-left: 3px solid #409EFF;
    font-size: 12px;
}

.current-label {
    color: #606266;
    margin-right: 8px;
}

.model-id {
    color: #909399;
    margin-left: 8px;
    font-size: 11px;
}

.scene-actions {
    text-align: center;
    padding-top: 16px;
    border-top: 1px solid #f0f0f0;
}

.auth-divider {
    height: 1px;
    background: #2a2a4a;
    margin: 24px 0;
}

.section-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;
}

.section-header i {
    color: #4c7ae3;
    font-size: 16px;
    margin-right: 8px;
}

.section-title {
    font-size: 15px;
    font-weight: 600;
    color: #e0e0e0;
}

.info-row {
    display: flex;
    align-items: center;
    padding: 10px 16px;
    border-bottom: 1px solid #2a2a4a;
}

.info-row:last-child {
    border-bottom: none;
}

.info-label {
    width: 120px;
    font-size: 14px;
    color: #909399;
    font-weight: 500;
}

.info-value {
    flex: 1;
    font-size: 14px;
    color: #e0e0e0;
    font-weight: 400;
}

.danger-text {
    color: #f56c6c !important;
    font-weight: 600 !important;
}

.auth-code-input /deep/ .el-textarea__inner {
    background: #1a1a2e;
    border: 1px solid #2a2a4a;
    color: #e0e0e0;
    font-family: 'Courier New', monospace;
    font-size: 13px;
}

.auth-code-input /deep/ .el-textarea__inner:focus {
    border-color: #4c7ae3;
}

</style>
