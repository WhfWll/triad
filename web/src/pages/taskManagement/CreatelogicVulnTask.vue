<template>
    <!-- 创建逻辑漏洞 -->
    <div>
        <div class="main-title" > 
            <router-link :underline="false" class="classA" :to="{ path: '/logicvuln' }"  >逻辑漏洞
            </router-link> 
            <label class="currentpagetitle" >   
                <span>创建任务</span>
            </label> 
        </div>
        <BannerBox ref="BannerBox" tips="任务参数通过任务场景进行默认配置，用户可以在创建任务时对默认参数进行调整。" style="margin-bottom: 16px;">
            <el-button type="primary" size="small" @click="submithandle" > 执行任务 </el-button> 
        </BannerBox>
        <div class="createtask_box">
            <el-form :model="taskform" ref="form" :rules="rules" style="height: 100%">
                <div class="basic_config">
                    <el-form-item label="" prop="name"  class="taskNameClass">
                        <label class="dialog_item_label">任务名称<i class="is-required" style="float: right">*</i></label>
                        <el-input   v-model="taskform.name"   size="small"
                        class="form_item_width"   placeholder="请输入任务名称" maxlength="50"></el-input>
                    </el-form-item> 
                    <el-form-item label="" prop="template"  >
                        <label class="dialog_item_label">检测类型</label>
                        <el-select  
                            v-model="taskform.type" 
                            size="small" placeholder="请选择" class="form_item_width"  
                            >
                            <el-option v-for="(item, index) in typelist" :key="index" :label="item.label"
                                :value="item.value"></el-option>
                        </el-select>  
                    </el-form-item>
                    <div  style="position: relative">
                        <label class="dialog_item_label" >测试范围 <i class="is-required" style="float: right">*</i></label>
                        <el-form-item label="" prop="target" label-width="0" class="target_box" 
                            style="display: inline-block; margin-right: 0">
                            <el-input    v-model="taskform.targetUrl" size="small"
                             placeholder="测试范围不能为空" class="form_item_width"  
                                ></el-input>
                            
                        </el-form-item> 
                        <el-tooltip placement="right-start">
                            <div slot="content">
                                测试范围支持接口和站点地址
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                    </div>
                  
                    <el-form-item label="" prop="template"  >
                        <label class="dialog_item_label">登录凭证</label>
                        <el-select  
                            v-model="taskform.loginCred_pattern" 
                            size="small" placeholder="请选择" class="form_item_width"
                            >
                            <el-option v-for="(item, index) in credPatternlist" :key="index" :label="item.label"
                                :value="item.value"></el-option>
                        </el-select> 
                        <div style="position: relative; margin-left: 112px;">
                            <el-input  type="textarea" rows="3" resize="none"   v-model="taskform.loginCred_value"   size="small" class="form_item_width"
                                ></el-input>
                            <el-tooltip placement="right-start">
                                <div slot="content">
                                    cookie/header名和cookie/header值用英文冒号隔开，不同组数据用英文分号隔开 <br /> 
                                    如 Cookie: lang=zh-CN; PHPSESSID=jvfbua8qs3vb81rapv45d2134q
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </div> 
                    </el-form-item>
                 
                    <el-form-item label="" prop="template"  v-if="taskform.type == 1">
                        <label class="dialog_item_label">待测凭证</label>
                        <el-select  
                            v-model="taskform.waitCred_pattern" 
                            size="small" placeholder="请选择" class="form_item_width"   
                           >
                            <el-option v-for="(item, index) in credPatternlist" :key="index" :label="item.label"
                                :value="item.value"></el-option>
                        </el-select>  
                        <div style="position: relative; margin-left: 112px;">
                            <el-input type="textarea" rows="3" resize="none" v-model="taskform.waitCred_value" 
                            size="small" class="form_item_width"
                                ></el-input>
                            <el-tooltip placement="right-start">
                                <div slot="content">
                                    cookie/header名和cookie/header值用英文冒号隔开，不同组数据用英文分号隔开 <br /> 
                                    如 Cookie: lang=zh-CN; PHPSESSID=jvfbua8qs3vb81rapv45d2134q
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute;  "></i>
                            </el-tooltip>
                        </div>
                    </el-form-item>
                </div>
                <div class="more_config">
                    <label  @click="showMoreconfig" style="cursor: pointer;">更多配置</label>
                    <i></i>
                </div>
               <div v-if="taskform.type == 1 && isShowMore">
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label">路径白名单</label>
                        <!-- <el-select size="small" 
                            v-model="taskform.whitePath"
                            multiple
                            filterable
                            allow-create
                            default-first-option 
                            placeholder=""
                            class="form_item_width"  >
                            <el-option
                                v-for="(item,i) in whitePathlist"
                                :key="i"
                                :label="item"
                                :value="item">
                            </el-option>
                        </el-select> -->
                        <el-input   v-model="taskform.whitePath" size="small" 
                            class="form_item_width"  rows="4"  
                            autocomplete="off" type="textarea" resize="none" placeholder="请输入路径白名单"></el-input>
                        <el-tooltip placement="right-start">
                            <div slot="content">
                                进行越权检测的url路径关键字
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                        <div style="padding-left: 112px;margin-top: 4px;">
                            <el-button type="primary" size="mini" style=" margin-right: 12px"
                            @click="clickuploadPathWhite()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px;  ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadIDPathWhite" ref="upload" @change="changeuploaIDPathWhite($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                    </el-form-item>
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label">路径黑名单</label>
                        <!-- <el-select size="small" 
                            v-model="taskform.blackPath"
                            multiple
                            filterable
                            allow-create
                            default-first-option 
                            placeholder=""
                            class="form_item_width"  >
                            <el-option
                                v-for="(item,i) in blackPathlist"
                                :key="i"
                                :label="item"
                                :value="item">
                            </el-option>
                        </el-select> -->
                        <el-input   v-model="taskform.blackPath" size="small" 
                            class="form_item_width"  rows="4"  
                            autocomplete="off" type="textarea" resize="none" placeholder="请输入路径黑名单"></el-input>
                        <el-tooltip placement="right-start">
                            <div slot="content">
                                用于二次筛选检测路径，<br/>例如：检测路径白名单配置id，检测路径黑名单配置wide，表示对含有id但不含wide的路径进行检测
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                        <div style="padding-left: 112px;margin-top: 4px;">
                            <el-button type="primary" size="mini" style=" margin-right: 12px"
                            @click="clickuploadPathBlack()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px;  ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadIDPathBlack" ref="upload" @change="changeuploaIDPathBlack($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                    </el-form-item>
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label" style="vertical-align: top;  ">关键字字典</label>
                        <el-input   v-model="taskform.keywords" size="small" 
                            class="form_item_width"  rows="4"  
                            autocomplete="off" type="textarea" resize="none" placeholder="请输入关键字字典"></el-input>
                        <el-tooltip placement="right-start">
                            <div slot="content">
                                作为判定漏洞是否存在的检测依据，用于提高准确率
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                        <div style="padding-left: 112px;margin-top: 4px;">
                            <el-button type="primary" size="mini" style=" margin-right: 12px"
                            @click="clickupload()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px;  ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadID" ref="upload" @change="changeuploaID($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                        
                    </el-form-item>
                </div>
                <div v-if="taskform.type == 2 && isShowMore">
                    <div>
                        <label for="" class="dialog_item_label">fuzz参数</label> 
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline"  >字符型</label>
                            <el-select size="small" 
                                v-model="taskform.fuzzParam.character"
                                multiple
                                filterable
                                allow-create
                                default-first-option 
                                placeholder=""
                                class="form_item_width"  >
                                <el-option
                                    v-for="(item,i) in characterlist"
                                    :key="i"
                                    :label="item"
                                    :value="item">
                                </el-option>
                            </el-select>
                            <el-tooltip placement="right-start">
                            <div slot="content">
                                对字符串类型的参数名进行fuzz，<br/>如 /vul/sqli/sqli_str.php?name=zhangsan 中的name
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                        </el-form-item>
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline"  >数字型</label>
                            <el-select size="small" 
                                v-model="taskform.fuzzParam.number"
                                multiple
                                filterable
                                allow-create
                                default-first-option 
                                placeholder=""
                                class="form_item_width"  >
                                <el-option
                                    v-for="(item,i) in numberlist"
                                    :key="i"
                                    :label="item"
                                    :value="item">
                                </el-option>
                            </el-select>
                            <el-tooltip placement="right-start">
                            <div slot="content">
                                对数字类型的参数名进行fuzz，<br/>如 /vul/sqli/sqli_id.php?id=1 中的id
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                        </el-form-item>
                    </div> 
                    <div>
                        <label for="" class="dialog_item_label">fuzz字典</label> 
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline" style="vertical-align: top;  " >字符型</label>
                            <el-input   v-model="taskform.fuzzDict.character" size="small" 
                                class="form_item_width"  rows="4"  
                                autocomplete="off" type="textarea" resize="none" placeholder=""></el-input>
                            <el-tooltip placement="right-start">
                                <div slot="content">
                                    使用该字典中的字符串进行fuzz，如 /vul/sqli/sqli_str.php?name=zhangsan 中的zhangsan
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <div style="padding-left: 114px;margin-top: 4px;margin-bottom: 4px;">
                            <el-button type="primary" size="mini" style="  margin-right: 12px"
                            @click="clickupload1()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px; ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadID1" ref="upload" @change="changeuploaID1($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline" style="vertical-align: top;  " >数字型</label>
                            <el-input   v-model="taskform.fuzzDict.number" size="small" 
                                class="form_item_width"  rows="4"  
                                autocomplete="off" type="textarea" resize="none" placeholder=""></el-input>
                            <el-tooltip placement="right-start">
                                <div slot="content">
                                    使用该字典中的字符串进行fuzz，如 /vul/sqli/sqli_id.php?id=1 中的 1
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <div style="padding-left: 114px;margin-top: 4px;">
                            <el-button type="primary" size="mini" style="  margin-right:12px"
                            @click="clickupload2()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px; ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadID2" ref="upload" @change="changeuploaID2($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                    </div>
                    <div>
                        <label for="" class="dialog_item_label" style="margin-bottom: 12px;">响应字典</label> 
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline" style="vertical-align: top;  " >json关键字</label>
                            <el-input   v-model="taskform.response.jsonKeyword" size="small" 
                                class="form_item_width"  rows="4"  
                                autocomplete="off" type="textarea" resize="none" placeholder="请输入json关键字"></el-input>
                            <el-tooltip placement="right-start">
                                <div slot="content">
                                    作为json响应报文的判断依据，用于提高准确率
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </el-form-item>
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline"  >非json关键字</label>
                            <el-switch
                                v-model="taskform.response.noJsonSwitch"
                                class="elSwitch"  >
                            </el-switch> 
                        </el-form-item>
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px" v-show="taskform.response.noJsonSwitch"> 
                            <el-input   v-model="taskform.response.noJsonKeyword" size="small"  style="margin-left: 116px;"
                                class="form_item_width"  rows="4"  
                                autocomplete="off" type="textarea" resize="none" ></el-input>
                                <el-tooltip placement="right-start">
                                <div slot="content">
                                    作为非json响应报文的判断依据，用于提高准确率
                                </div>
                                <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                            </el-tooltip>
                        </el-form-item>
                    </div> 
                </div>
                <div v-if="taskform.type == 3 && isShowMore">
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label">路径白名单</label>
                        <el-select size="small" 
                            v-model="taskform.whitePath"
                            multiple
                            filterable
                            allow-create
                            default-first-option 
                            placeholder=""
                            class="form_item_width"  >
                            <el-option
                                v-for="(item,i) in whitePathlist"
                                :key="i"
                                :label="item"
                                :value="item">
                            </el-option>
                        </el-select>
                        <el-tooltip placement="right-start">
                            <div slot="content">
                                进行越权检测的url路径关键字
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                    </el-form-item>
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label">路径黑名单</label>
                        <el-select size="small" 
                            v-model="taskform.blackPath"
                            multiple
                            filterable
                            allow-create
                            default-first-option 
                            placeholder=""
                            class="form_item_width"  >
                            <el-option
                                v-for="(item,i) in blackPathlist"
                                :key="i"
                                :label="item"
                                :value="item">
                            </el-option>
                        </el-select>
                        <el-tooltip placement="right-start">
                            <div slot="content">
                                用于二次筛选检测路径，<br/>例如：检测路径白名单配置id，检测路径黑名单配置wide，表示对含有id但不含wide的路径进行检测
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                    </el-form-item>
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label" style="vertical-align: top;  ">关键字字典</label>
                        <el-input   v-model="taskform.keywords" size="small" 
                            class="form_item_width"  rows="4"  
                            autocomplete="off" type="textarea" resize="none" placeholder="请输入关键字字典"></el-input>
                        <el-tooltip placement="right-start">
                            <div slot="content">
                                作为判定漏洞是否存在的检测依据，用于提高准确率
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip>
                        <div style="padding-left: 112px;margin-top: 4px;">
                            <el-button type="primary" size="mini" style=" margin-right: 12px"
                            @click="clickupload()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px;  ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadID" ref="upload" @change="changeuploaID($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                        
                    </el-form-item>
                    <el-form-item label="" prop="port_scan_type">
                        <label class="dialog_item_label" style="vertical-align: top;  ">凭证身份标识字典</label>
                        <el-input   v-model="taskform.credIdentifyList" size="small" 
                            class="form_item_width"  rows="4"  
                            autocomplete="off" type="textarea" resize="none" placeholder="请输入凭证身份标识字典"></el-input>
                        <!-- <el-tooltip placement="right-start">
                            <div slot="content">
                                作为判定漏洞是否存在的检测依据，用于提高准确率
                            </div>
                            <i class="iconfont icontishi icontsstyle" style="position: absolute; "></i>
                        </el-tooltip> -->
                        <div style="padding-left: 112px;margin-top: 4px;">
                            <el-button type="primary" size="mini" style=" margin-right: 12px"
                            @click="clickuploadcredIdentify()">导入</el-button>
                            <span style="color: rgba(72, 72, 102, 0.32);font-size: 12px;  ">只能上传.txt格式文件</span>
                            <input type="file" class="btnUploadIDcredIdentify" ref="upload" @change="changeuploaIDcredIdentify($event)"
                                style="display: none" id="input-file-ID" />
                        </div>
                        
                    </el-form-item>
                </div>
                <div  v-if=" isShowMore " >
                    <label for="" class="dialog_item_label">爬虫</label> 
                    <div  >
                        <el-form-item prop = 'scope' label=" " style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline"  >爬取范围</label>
                            <el-select v-model="taskform.crawler.range"  class="frame_width"  
                                clearable placeholder="爬取范围"  size="small" ref="vulSelect">  
                                <el-option
                                    v-for="(item, i) in scanRangelist"
                                    :key="i"
                                    :label="item.label"
                                    :value="item.value"> 
                                </el-option>
                            </el-select>  
                        </el-form-item>
                        <el-form-item prop = 'depth' label=" "   style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline">爬取深度</label>
                                <el-select v-model="taskform.crawler.depth"   class="frame_width"   clearable placeholder="爬取深度"  size="small" ref="vulSelect"  >  
                                    <el-option
                                        v-for="(item, i) in maxDeeplist"
                                        :key="i"
                                        :label="item.label"
                                        :value="item.value"> 
                                    </el-option>
                                </el-select> 
                        </el-form-item> 
                        <el-form-item prop = 'speed' label=" "   style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline">最大链接数</label>
                            <el-select v-model="taskform.crawler.maxlink"   class="frame_width"  clearable  size="small"   >  
                                <el-option
                                    v-for="(item, i) in maxlinklist"
                                    :key="i"
                                    :label="item.label"
                                    :value="item.value"> 
                                </el-option>
                            </el-select>  
                        </el-form-item> 
                        <el-form-item prop = 'timeout' label=" "   style="margin-bottom:10px"> 
                            <label class="dialog_item_label_m topline">单链接超时</label>
                                <el-select v-model="taskform.crawler.singleLink"   class="frame_width"   clearable placeholder="单链接超时"  size="small" ref="vulSelect"  >   
                                    <el-option
                                        v-for="(item, i) in timeoutlist"
                                        :key="i"
                                        :label="item.label"
                                        :value="item.value"> 
                                    </el-option>
                                </el-select> 
                        </el-form-item>
                        <el-form-item label=""    style="margin-bottom:10px">
                            <label class="dialog_item_label_m">敏感词</label> 
                            <el-select size="small" 
                                v-model="taskform.crawler.sensitive"
                                multiple
                                filterable
                                allow-create
                                default-first-option 
                                placeholder=""
                                class="form_item_width"  >
                                <el-option
                                    v-for="(item,i) in sensitivelist"
                                    :key="i"
                                    :label="item"
                                    :value="item">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item label=""    style="margin-bottom:10px">
                            <label class="dialog_item_label_m">关键字白名单</label>
                        
                            <el-select size="small" 
                                v-model="taskform.crawler.whiteWord"
                                multiple
                                filterable
                                allow-create
                                default-first-option 
                                placeholder=""
                                class="form_item_width"  >
                                <el-option
                                    v-for="(item,i) in whiteWordlist"
                                    :key="i"
                                    :label="item"
                                    :value="item">
                                </el-option>
                            </el-select>
                        </el-form-item>
                        <el-form-item label=""   style="margin-bottom:10px">
                            <label class="dialog_item_label_m">关键字黑名单</label>
                         
                            <el-select size="small" 
                                v-model="taskform.crawler.blackWord"
                                multiple
                                filterable
                                allow-create
                                default-first-option 
                                placeholder=""
                                class="form_item_width"  >
                                <el-option
                                    v-for="(item,i) in blackWordlist"
                                    :key="i"
                                    :label="item"
                                    :value="item">
                                </el-option>
                            </el-select>
                        </el-form-item>
                    </div>
                </div> 
               
            </el-form>
        </div>
    </div>
</template>
<style lang="less" scoped>
.createtask_box{
    padding: 24px;
    background: #fff;
    min-height: calc(100% - 39px);
    box-sizing: border-box;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}
.form_item_width{
    width: 320px;
}
.frame_width{
    width: 318px;
}
.basic_config{
    position:relative;
    padding-bottom: 10px;
}
.basic_config::before{
    left: 0;
    bottom: 0;
    width: 100%;
    height: 1px;
    content: "";
    position: absolute; 
    background: #E8E8F5;
    z-index: 1;
}
.more_config{
    margin:24px 0;
    label{ 
        color: #4C7AE3;
        font-size:13px;
    }
    i{
        display:inline-block;
        width:16px;
        height:16px;
        background:url(../../assets/images/show@2x.png);
        background-size: cover;
        vertical-align: middle;
    }
}
.dialog_item_label_m {
    display: inline-block;
    font-size: 14px;
    font-weight: 500;
    color: rgba(72, 72, 102, 0.87);
    padding-left: 11px;
    width: 104px; 
    line-height: 16px;
}
i.is-required {
    margin-right: 4px;
    color: #f56c6c;
    font-size: 12px; 
}
</style>
<script>
import BannerBox from "@/components/BannerBox.vue"; 
import { logic } from '@/api/task.js'
export default {
    name:'createlogicvuln',
    components:{
        BannerBox
    },
    data(){
        return{
            flag: this.$route.query.flag, //1新建，2复制
            id:this.$route.query.id, //1新建，2复制
            taskform:{
                name:'',
                type:1,
                targetUrl:'',
                loginCred_pattern:1,
                waitCred_pattern:1,
                loginCred_value:'',
                waitCred_value:'',
                whitePath:'',
                blackPath:[],
                crawler:{
                    "range": "",
                    "depth": "",
                    "maxlink": '',
                    "singleLink": '',
                    "sensitive": [],
                    "blackWord": [],
                    "whiteWord": []
                },
                character:[],
                number:[],
                response:{
                    jsonKeyword:'',
                    noJsonSwitch:true,
                    noJsonKeyword:'',
                },
                fuzzParam:{
                    character:[],
                    number:[],
                },
                fuzzDict:{
                    character:'',
                    number:'',
                },
                credIdentifyList:'',
                keywords:'',

            },
            typelist:[],
            credPatternlist:[],
            rules:{
                
            },
            isShowMore:false,
            whitePathlist:[],
            blackPathlist:[],
            scanRangelist:[],
            maxDeeplist:[],
            maxlinklist:[],
            timeoutlist:[],
            characterlist:[],
            numberlist:[],
            sensitivelist:[],
            whiteWordlist:[],
            blackWordlist:[],
        }
    },
    created(){
        this.$store.state.activefirstMenu = '/logicvuln';
    },
    mounted(){
        if(this.flag == 2){
            this.getInfo();
        }
        this.getEnum();
    },
    methods:{
        async getEnum(){
            const res = await logic.logicEnum();
            if(res.code == 200){
                this.typelist = res.data.scanType;
                this.credPatternlist = res.data.credPattern;
                //this.whitePathlist = res.data.whitePath == '' ? []:res.data.whitePath.split(',');
                // this.blackPathlist = res.data.blackPath == '' ? []:res.data.blackPath.split(',');
                this.scanRangelist = res.data.crawler.scanRange;
                this.maxDeeplist = res.data.crawler.maxDeep;
                this.maxlinklist = res.data.crawler.maxUrl;
                this.timeoutlist = res.data.crawler.timeout;
                this.characterlist = res.data.fuzzParam.character == ''?[]:res.data.fuzzParam.character.split(',');
                this.numberlist =  res.data.fuzzParam.number == ''?[]:res.data.fuzzParam.number.split(',');
                this.sensitivelist = res.data.crawler.sensitive == ''?[]:res.data.crawler.sensitive.split(',');
                this.whiteWordlist = res.data.crawler.whiteList == ''?[]:res.data.crawler.whiteList.split(',');
                this.blackWordlist = res.data.crawler.blackList ==''?[]:res.data.crawler.blackList.split(',');

                // 设置默认值
                console.log(res.data.whitePath)
                this.taskform.blackPath =  res.data.blackPath;
                this.taskform.whitePath =  res.data.whitePath
                this.taskform.fuzzParam.character = this.characterlist;
                this.taskform.fuzzParam.number = this.numberlist;
                this.taskform.response.jsonKeyword = res.data.response.jsonKeyword;
                this.taskform.response.noJsonKeyword = res.data.response.noJsonKeyword;

                this.scanRangelist.forEach(item=>{
                    if( item.isDefault==true){
                        this.taskform.crawler.range = item.value
                    }
                })
                this.maxDeeplist.forEach(item=>{
                    if( item.isDefault==true){
                        this.taskform.crawler.depth = item.value
                    }
                })
                this.maxlinklist.forEach(item=>{
                    if( item.isDefault==true){
                        this.taskform.crawler.maxlink = item.value
                    }
                })
                this.timeoutlist.forEach(item=>{
                    if( item.isDefault==true){
                        this.taskform.crawler.singleLink = item.value
                    }
                })
                this.taskform.crawler.sensitive = this.sensitivelist ;
                this.taskform.crawler.whiteWord = this.whiteWordlist ;
                this.taskform.crawler.blackWord = this.blackWordlist 

            }
        },
        async getInfo(){
            const dt = await logic.logicTaskcopy({
                id: this.id, 
              });
              if (dt.code == 200) { 
                console.log(dt.data)
                this.getShowinfo(dt.data);
              } else {
                this.$message({
                  message: dt.msg,
                  type: "error"
                });
              }
        },
        getShowinfo(data){ //渲染信息
            this.taskform.name = data.name;
            this.taskform.type = data.type;
            this.taskform.targetUrl = data.targetUrl;

            let config = JSON.parse(data.scanConfig); 

            this.taskform.crawler.range = config.crawler.range;
            this.taskform.crawler.depth = config.crawler.depth;
            this.taskform.crawler.maxlink = config.crawler.maxlink;
            this.taskform.crawler.singleLink = config.crawler.singleLink;
            this.taskform.crawler.sensitive = config.crawler.sensitive;
            this.taskform.crawler.blackWord = config.crawler.blackWord; 
            this.taskform.crawler.whiteWord = config.crawler.whiteWord; 

            this.taskform.loginCred_pattern = config.loginCred.pattern;
            this.taskform.loginCred_value =  config.loginCred.value;

            this.taskform.waitCred_pattern = config.waitCred.pattern;
            this.taskform.waitCred_value = config.waitCred.value;

            this.taskform.fuzzParam.character = config.fuzzParam.character.split(',')
            this.taskform.fuzzParam.number = config.fuzzParam.number.split(',')

            this.taskform.fuzzDict.character = config.fuzzDict.character
            this.taskform.fuzzDict.number = config.fuzzDict.number

            this.taskform.response.jsonKeyword = config.response.jsonKeyword
            this.taskform.response.noJsonSwitch = config.response.noJsonSwitch
            this.taskform.response.noJsonKeyword = config.response.noJsonKeyword

            this.taskform.whitePath =  config.whitePath.split(',')
            this.taskform.blackPath =  config.blackPath.split(',')
            this.taskform.keywords =  config.keywords;
            this.taskform.credIdentifyList = config.credIdentifyList;

        },
        showMoreconfig(){
            this.isShowMore  = !this.isShowMore;
        },
        async submithandle(){
            // 爬虫
            let _crawler = {
                range: this.taskform.crawler.range,
                    depth: this.taskform.crawler.depth,
                    maxlink: this.taskform.crawler.maxlink,
                    singleLink: this.taskform.crawler.singleLink,
                    sensitive: this.taskform.crawler.sensitive,
                    blackWord: this.taskform.crawler.blackWord,
                    whiteWord: this.taskform.crawler.whiteWord
            };
            // config配置
            let config = {
                loginCred:{
                    pattern:this.taskform.loginCred_pattern,
                    value:this.taskform.loginCred_value
                },
                waitCred:{
                    pattern:this.taskform.waitCred_pattern,
                    value:this.taskform.waitCred_value
                },
                fuzzParam:{
                    character:this.taskform.fuzzParam.character.join(','),
                    number:this.taskform.fuzzParam.number.join(','),
                },
                fuzzDict:{
                    character:this.taskform.fuzzDict.character,
                    number:this.taskform.fuzzDict.number,
                },
                response:{
                    jsonKeyword:this.taskform.response.jsonKeyword,
                    noJsonSwitch:this.taskform.response.noJsonSwitch,
                    noJsonKeyword:this.taskform.response.noJsonKeyword
                },
                whitePath:this.taskform.whitePath.replaceAll('\n', ",").replaceAll('\r', ","),
                blackPath:this.taskform.blackPath.replaceAll('\n', ",").replaceAll('\r', ","),
                keywords:this.taskform.keywords,
                credIdentifyList:this.taskform.credIdentifyList,
                crawler:_crawler
            }
            // 提交参数
            let param = {
                name:this.taskform.name,
                type:this.taskform.type,
                targetUrl:this.taskform.targetUrl,
                scanConfig:JSON.stringify(config)
            };   
 
            // return ; 
            const res = await logic.createtask(param);
            if(res.code == 200){
                this.$router.push({
                    path: `/logicvuln`,
                    query: {  
                    }
                });
            }else{
                this.$message({
                    message: res.msg,
                    type: "error",
                });
            }

        },
        clickuploadPathWhite(){
            document.querySelector(".btnUploadIDPathWhite").click();
        },
        clickuploadPathBlack(){
            document.querySelector(".btnUploadIDPathBlack").click();
        },
        clickupload(){
            document.querySelector(".btnUploadID").click();
        },
        changeuploaID(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1); 
            if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) { 
                        that.taskform.keywords = reader.result; 
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        changeuploaIDPathWhite(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1); 
            if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) { 
                        that.taskform.whitePath = reader.result; 
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        changeuploaIDPathBlack(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1); 
            if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) { 
                        that.taskform.blackPath = reader.result; 
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        clickupload1(){
            document.querySelector(".btnUploadID1").click();
        },
        changeuploaID1(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1); 
            if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) { 
                        that.taskform.fuzzDict.character = reader.result; 
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        clickupload2(){
            document.querySelector(".btnUploadID2").click();
        },
        changeuploaID2(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1); 
            if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) { 
                        that.taskform.fuzzDict.number = reader.result; 
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
        clickuploadcredIdentify(){
            document.querySelector(".btnUploadIDcredIdentify").click();
        },
        changeuploaIDcredIdentify(e){
            var that = this;
            var f = e.target.files[0];
            if (!f) return; 
            let fileSuffix = f.name.substr(f.name.lastIndexOf(".") + 1); 
            if (fileSuffix.indexOf("txt") != -1) {
                var input = e.target;
                var reader = new FileReader();
                reader.onload = function () {
                    if (reader.result) { 
                        that.taskform.credIdentifyList = reader.result; 
                    }
                };
                reader.readAsText(input.files[0]);
            }
        },
    }
}
</script>