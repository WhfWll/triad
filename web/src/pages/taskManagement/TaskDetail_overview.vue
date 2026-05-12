<template>
    <!-- 概览 -->
    <div> 
        <div class="task_detail">
            <BannerBox tips=""  >
                <div class="banner">
                    <i class="banner_icon"></i>
                    <span>
                        <div> {{task_name}}</div>
                        <p>任务名称</p>
                    </span> 
                </div>
            </BannerBox>
            <div class="basic_info">
                <ul>
                    <li>
                        <div>
                            <label for="">任务场景</label>
                            <span>{{basic_info.taskTemplateName}}</span>
                        </div> 
                        <div>
                            <label for="">执行方式</label> 
                            <!-- 立即执行 -->
                            <span v-if="basic_info.executeType == 1">{{basic_info.executeTypeName}}</span> 

                            <!-- 定时任务 -->
                            <span v-if="basic_info.executeType == 2">
                                <el-tooltip class="item" effect="dark"  :content="'计划时间：'+basic_info.cyclePlanningName "  placement="bottom" v-if="ishowtip">
                                    <span> {{basic_info.executeTypeName}} </span>
                                </el-tooltip>  
                                <!-- <span style="margin-left: 20px;">计划时间：{{basic_info.cyclePlanningName}}</span> -->
                            </span> 
                            <!-- 周期任务 -->
                            <span v-if="basic_info.executeType == 3">

                                <el-tooltip class="item" effect="dark"  :content="'计划时间：'+basic_info.cyclePlanningName+'; 终止时间：'+ basic_info.endTimeName"  placement="bottom" v-if="ishowtip">
                                    <span>{{basic_info.executeTypeName}} </span>
                                </el-tooltip>  

                                <!-- <span  style="margin-left: 22px;">计划时间：{{basic_info.cyclePlanningName}} </span>
                                <span  style="margin-left: 20px;">终止时间：{{basic_info.endTimeName}}</span> -->
                            </span>
                        </div>
                    </li>
                    <li>
                        <div>
                            <label for="">创建时间</label>
                            <span>{{basic_info.createTime}}</span>
                        </div> 
                        <div>
                            <label for="">结束时间</label>
                            <span>{{basic_info.updateTime}}</span>
                        </div>
                    </li>
                    <li>
                        <div>
                            <label for="">执行时间段</label>
                            <span  > 
                                <el-tooltip class="item" effect="dark"  :content="runtimePeriod"  placement="bottom" v-if="ishowtip">
                                    <span>{{ basic_info.runtimePeriod }}</span>
                                </el-tooltip>
                                <span v-else>{{ basic_info.runtimePeriod }}</span>
                            </span>
                        </div> 
                        <div>
                            <label for="">风险等级</label>
                            <span :class="[ 
                                    { 'riskstyle risk_hight': basic_info.riskLevel == 1 } ,
                                    { 'riskstyle risk_middle': basic_info.riskLevel == 2 },
                                    { 'riskstyle risk_low': basic_info.riskLevel == 3 },
                                    { 'riskstyle risk_nofind': basic_info.riskLevel == 4 }]"><i></i>

                            {{ basic_info.riskLevelName }}</span>
                        </div>
                    </li>
                </ul>
                
            </div>
        </div>
        <div class="target_risk box_bg_shadow">
            <div class="box_title">
                <label>目标风险</label>
                <i></i>
            </div>
            <div class="box_content">
                <el-row>
                    <el-col :span="6"> 
                        <img class="risk_icon" src="../../assets/images/risk_height@2x.png" />
                        <div class="risk_content risk_content_height">
                            <p class="risk_number">
                                <label style="cursor: pointer;" @click="$emit('act',{tab:'tabs2',level:1})">{{  target_risk[0] }}</label>
                                <!-- <i class="el-icon-caret-right"></i> --></p>
                            <p class="risk_name" >高危目标</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_middle@2x.png" />
                        <div class="risk_content risk_content_middle">
                            <p class="risk_number">
                                <label style="cursor: pointer;" @click="$emit('act',{tab:'tabs2',level:2})">{{  target_risk[1] }}</label> 
                                <!-- <i class="el-icon-caret-right"></i> -->
                            </p>
                            <p class="risk_name" >中危目标</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_low@2x.png" />
                        <div class="risk_content risk_content_low">
                            <p class="risk_number">
                                <label style="cursor: pointer;" @click="$emit('act',{tab:'tabs2',level:3})">{{  target_risk[2] }}</label>
                                 <!-- <i class="el-icon-caret-right"></i> -->
                                </p>
                            <p class="risk_name" >低危目标</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_safe@2x.png" />
                        <div class="risk_content risk_content_safe">
                            <p class="risk_number"><label style="cursor: pointer;"  @click="$emit('act',{tab:'tabs2',level:4})">{{  target_risk[3] }}</label> 
                                <!-- <i class="el-icon-caret-right"></i> -->
                            </p>
                            <p class="risk_name" >安全目标</p>
                        </div>
                    </el-col>
                </el-row>
            </div>
        </div>
        <div class="target_survival box_bg_shadow">
            <div class="box_title">
                <label>目标存活</label>
                <i></i>
            </div>
            <div class="box_content">
                <el-row>
                    <el-col  :span="6">
                        <img class="risk_icon" src="../../assets/images/cunhuo@2x.png" />
                        <div class="risk_content ">
                            <p class="risk_number"><label  >{{target_survival[0]}}</label></p>
                            <p class="risk_name" >存活目标 
                                <!-- <i class="el-icon-arrow-down"></i> -->
                            </p>
                        </div>
                    </el-col>
                    <el-col  :span="6">
                        <img class="risk_icon" src="../../assets/images/bucunhuo@2x.png" />
                        <div class="risk_content ">
                            <p class="risk_number"><label  >{{target_survival[1]}}</label></p>
                            <p class="risk_name" >不存活目标 
                                <!-- <i class="el-icon-arrow-down"></i> -->
                            </p>
                        </div>
                    </el-col>
                    <el-col  :span="6">
                        <img class="risk_icon" src="../../assets/images/daijiance@2x.png" />
                        <div class="risk_content ">
                            <p class="risk_number"><label  >{{target_survival[2]}}</label></p>
                            <p class="risk_name" >待检测目标 
                                <!-- <i class="el-icon-arrow-down"></i> -->
                            </p>
                        </div>
                    </el-col>
                    <el-col  :span="6">
                        <img class="risk_icon" src="../../assets/images/targetTotal@2x.png" />
                        <div class="risk_content ">
                            <p class="risk_number"><label  >{{target_survival[3]}}</label></p>
                            <p class="risk_name" >目标总数
                                 <!-- <i class="el-icon-arrow-down"></i> -->
                                </p>
                        </div>
                    </el-col>
                </el-row>
            </div>
            <!-- <div class="target_detail">
                195.1621

            </div> -->
        </div>
        <div >
            <el-row :gutter="16">
                <el-col :span="8">
                    <div class="box_bg_shadow">
                        <div class="box_title">
                            <label>漏洞风险</label> 
                        </div>
                        <div class="vuln_risk" >
                            <div id="vuln_risk">

                            </div>
                            <div class="legend">
                                <el-row >
                                    <el-col :span="6" class="vlun_risk_deadly">
                                        <div class="name"><i></i>致命漏洞</div>
                                        <div class="number" style="cursor: pointer;" @click="$emit('actLD',{tab:'tabs4',level:1})">{{risk_num[0]}}</div>
                                        <div class="percentage">{{risk_percentage[0]}}%</div>
                                    </el-col>
                                    <el-col :span="6" class="vlun_risk_height">
                                        <div class="name"><i></i>高危漏洞</div>
                                        <div class="number" style="cursor: pointer;" @click="$emit('actLD',{tab:'tabs4',level:2})"  >{{risk_num[1]}}</div>
                                        <div class="percentage">{{risk_percentage[1]}}%</div>
                                    </el-col>
                                    <el-col :span="6" class="vlun_risk_middle">
                                        <div class="name"><i></i>中危漏洞</div>
                                        <div class="number"  style="cursor: pointer;" @click="$emit('actLD',{tab:'tabs4',level:3})">{{risk_num[2]}}</div>
                                        <div class="percentage">{{risk_percentage[2]}}%</div>
                                    </el-col>
                                    <el-col :span="6" class="vlun_risk_low">
                                        <div class="name"><i></i>低危漏洞</div>
                                        <div class="number" style="cursor: pointer;" @click="$emit('actLD',{tab:'tabs4',level:4})" >{{risk_num[3]}}</div>
                                        <div class="percentage">{{risk_percentage[3]}}%</div>
                                    </el-col>
                                </el-row>
                            </div>
                        </div>
                    </div>
                </el-col>
                <el-col :span="8">
                    <div class="box_bg_shadow">
                        <div class="box_title">
                            <label>漏洞状态</label> 
                        </div>
                        <div class="vuln_status" >
                            <div id="vuln_status">

                            </div>
                            <div class="legend" style="padding: 0 40px;">
                                <el-row > 
                                    <el-col :span="8" class="use_achievement" v-for="(item,i ) in vulExploitImpact" :key="i"> 
                                        <div class="name"><i ></i>{{ item.label }}</div>
                                        <div class="number"  >{{ item.value }}</div>
                                        <div class="percentage">{{item.percentage}}%</div>
                                    </el-col>
                                </el-row>
                            </div>
                        </div>
                    </div>
                </el-col>
                <el-col :span="8">
                    <div class="box_bg_shadow">
                        <div class="box_title">
                            <label>漏洞取证</label> 
                        </div>
                        <div class="vuln_obtain_evidence" >
                            <div id="vuln_obtain_evidence"> 
                            </div>
                            <div class="legend" style="padding: 0 40px;">
                                <el-row >
                                    <!-- <el-col :span="8" class="Remote_control">
                                        <div class="name"><i ></i>远程控制</div>
                                        <div class="number"  >{{obtain_evidence_num[0]}}</div> 
                                    </el-col>
                                    <el-col :span="8" class="data_breach">
                                        <div class="name"><i ></i>数据泄露</div>
                                        <div class="number"  >{{obtain_evidence_num[1]}}</div> 
                                    </el-col>
                                    <el-col :span="8" class="Logo_Credentials">
                                        <div class="name"><i ></i>登录凭证</div>
                                        <div class="number"  >{{obtain_evidence_num[2]}}</div> 
                                    </el-col>  -->
                                    <el-col :span="8" class="Logo_Credentials" v-for="(item,i) in obtain_evidence" :key="i">
                                        <div class="name"><i ></i>{{item.label}}</div>
                                        <div class="number"  >{{ item.value}}</div> 
                                    </el-col>
                                </el-row>
                            </div>
                        </div>
                    </div>
                </el-col>
            </el-row>
        </div>
        <div class="box_bg_shadow" style="max-height: 410px;margin-bottom: 16px;">
            <el-row>
                <el-col :span="8">
                    <div class="box_title"> 
                        <img class="title_icon" src="../../assets/images/port@2x.png" />
                        <span class="spans"> 端口</span> 
                    </div>
                   
                    <!-- <ul class="ul_box">
                        <li v-for="(item,i) in port_list">
                            <label>{{item.label}}</label>
                            <span>{{item.value}}</span>
                        </li>
                    </ul> -->
                 
                  <!-- <div>  <dv-scroll-ranking-board  class="rankingWa" v-if="typeof configChart.data == 'object'&&configChart.data.length>0" :config="configChart" style="width:80%;height:300px" /></div> -->
                        <div style="height:400px;overflow:auto;padding:0 50px">  
                        <div v-for="(item,index) in configChart.data" :key="index" style="display: flex;justify-content: space-between;align-items: center;">
                            <label for="" class="lbtxt">{{ item.name }}：
                            </label>
                            <span class="spans">{{ item.value }}</span>
                        </div>
                    </div>
                
                    </el-col>
                <el-col :span="8">
                    <div class="box_title"> 
                        <img class="title_icon" src="../../assets/images/service@2x.png" />
                        <span class="spans">服务</span> 
                    </div>
                    <!-- <ul class="ul_box">
                        <li v-for="(item,i) in service_list">
                            <label>{{item.label}}</label>
                            <span>{{item.value}}</span>
                        </li>
                    </ul> -->
                  <!-- <div>  <dv-scroll-ranking-board  class="rankingWa2" v-if="typeof configChart2.data == 'object'&&configChart2.data.length>0" :config="configChart2" style="width:80%;height:300px" /></div> -->
                    <div style="height:400px;overflow:auto;padding:0 50px">  
                        <div v-for="(item,index) in configChart2.data" :key="index" style="display: flex;justify-content: space-between;align-items: center;">
                            <label for="" class="lbtxt">{{ item.name }}：
                            </label>
                            <span class="spans">{{ item.value }}</span>
                        </div>
                    </div>
                </el-col>
                <el-col :span="8">
                    <div class="box_title"> 
                        <img class="title_icon" src="../../assets/images/assembly@2x.png" />
                        <span class="spans">组件</span> 
                    </div>
                    <!-- <ul class="ul_box">
                        <li v-for="(item,i) in component_list">
                            <label>{{item.label}}</label>
                            <span>{{item.value}}</span>
                        </li>
                    </ul> -->
                    <div style="height:400px;overflow:auto;padding:0 50px">  
                        <div v-for="(item,index) in configChart3.data" :key="index" style="display: flex;justify-content: space-between;align-items: center;">
                            <label for="" class="lbtxt">{{ item.name }}：
                            </label>
                            <span class="spans">{{ item.value }}</span>
                        </div>
                    </div>
                    <!-- <div>  <dv-scroll-ranking-board  class="rankingWa3" v-if="typeof configChart3.data == 'object'&&configChart3.data.length>0" :config="configChart3" style="width:80%;height:300px" /></div> -->
                </el-col>
            </el-row>
            
        </div>
        <div class="box_bg_shadow" style="max-height: 410px;">
            <el-row>
                <el-col :span="8">
                    <div class="box_title"> 
                        <img class="title_icon" src="../../assets/images/system@2x.png" />
                        <span class="spans">系统</span> 
                    </div>
                    <!-- <ul class="ul_box">
                        <li v-for="(item,i) in opSys_list">
                            <label>{{item.label}}</label>
                            <span>{{item.value}}</span>
                        </li>
                    </ul> -->
                    <!-- <div>  <dv-scroll-ranking-board  class="rankingWa4" v-if="typeof configChart4.data == 'object'&&configChart4.data.length>0" :config="configChart4" style="width:80%;height:300px" /></div> -->
                     <div style="height:400px;overflow:auto;padding:0 50px">  
                        <div v-for="(item,index) in configChart4.data" :key="index" style="display: flex;justify-content: space-between;align-items: center;">
                            <label for="" class="lbtxt">{{ item.name }}：
                            </label>
                            <span class="spans">{{ item.value }}</span>
                        </div>
                    </div>
                
                </el-col>
                <el-col :span="8">
                    <div class="box_title"> 
                        <img class="title_icon" src="../../assets/images/subdomain@2x.png" />
                        <span class="spans">子域名</span> 
                    </div>
                    <!-- <ul class="ul_box"> 
                        <li v-for="(item,i) in subDomain_list">
                            <label>{{item.label}}</label>
                            <span>{{item.value}}</span>
                        </li>
                    </ul> -->
                    <!-- <div>  <dv-scroll-ranking-board  class="rankingWa5" v-if="typeof configChart5.data == 'object'&&configChart5.data.length>0" :config="configChart5" style="width:80%;height:300px" /></div> -->
                           <div style="height:400px;overflow:auto;padding:0 50px">  
                        <div v-for="(item,index) in configChart5.data" :key="index" style="display: flex;justify-content: space-between;align-items: center;">
                            <label for="" class="lbtxt">{{ item.name }}：
                            </label>
                            <span class="spans">{{ item.value }}</span>
                        </div>
                    </div>
                </el-col>
                <el-col :span="8">
                    <div class="box_title"> 
                        <img class="title_icon" src="../../assets/images/url@2x.png" />
                        <span class="spans">URL标签</span> 
                    </div>
                    <!-- <ul class="ul_box">
                        <li v-for="(item,i) in urlTags_list">
                            <label>{{item.label}}</label>
                            <span>{{item.value}}</span>
                        </li>
                    </ul> -->
                    <!-- <div>  <dv-scroll-ranking-board  class="rankingWa6" v-if="typeof configChart6.data == 'object'&&configChart6.data.length>0" :config="configChart6" style="width:80%;height:300px" /></div> -->
                           <div style="height:400px;overflow:auto;padding:0 50px">  
                        <div v-for="(item,index) in configChart6.data" :key="index" style="display: flex;justify-content: space-between;align-items: center;">
                            <label for="" class="lbtxt">{{ item.name }}：
                            </label>
                            <span class="spans">{{ item.value }}</span>
                        </div>
                    </div>
                
                </el-col>
            </el-row> 
        </div>
        <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
            <label style="margin-left:10px;color:#484866A3">任务配置</label>
        </div>
        <!-- 端口扫描 -->
        <div v-if="configData.portScanConfig.isOpen" class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:20px" >
                <label style="margin-left:10px;color:#484866A3">端口扫描</label>
            </div>
            <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                <div>
                    <span >扫描方式</span>
                    <span style="margin-left:10px">{{configData.portScanConfig.tcpScanTypeZh || ''}}</span>
                </div>
                <div>
                    <span >端口扫描范围</span>
                    <span style="margin-left:10px">{{configData.portScanConfig.portScanTypeZh|| ''}}</span>
                </div>
                
            </div>
            <div style="font-size:12px;color:#484866A3;margin-bottom:10px">
                扫描端口
            </div>
            <el-input
            type="textarea"
            :rows="3"
            placeholder="请输入内容"
            v-model="configData.portScanConfig.scanPort">
            </el-input>
        </div>
        <!-- 动态爬虫 -->
        <div v-if="configData.webCrawlerConfig.isOpen" class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">动态爬虫</label>
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                <div>
                    <span >爬取范围</span>
                    <span style="margin-left:10px">{{configData.webCrawlerConfig.scanRangeZh}}</span>
                </div>
                <div>
                    <span >爬取深度</span>
                    <span style="margin-left:10px">{{configData.webCrawlerConfig.maxDepthZh}}</span>
                </div>
                <div>
                    <span >最大链接数</span>
                    <span style="margin-left:10px">{{configData.webCrawlerConfig.maxUrlZh}}</span>
                </div>
                
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                <div>
                    <span >URL去重</span>
                    <span style="margin-left:10px">{{configData.webCrawlerConfig.scanRepeatZh}}</span>
                </div>
                <div>
                    <span >单链接超时</span>
                    <span style="margin-left:10px">{{configData.webCrawlerConfig.timeoutZh}}</span>
                </div>
                <div>
                    <span >最大爬取时长</span>
                    <span style="margin-left:10px">{{configData.webCrawlerConfig.fullTimeoutZh}}</span>
                </div>
                
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>后缀过滤</span>
                <span style="margin-left:10px">{{configData.webCrawlerConfig.suffixFilter}}</span>
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>关键字白名单</span>
                <span style="margin-left:10px">{{configData.webCrawlerConfig.blackList}}</span>
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>关键字黑名单</span>
                <span style="margin-left:10px">{{configData.webCrawlerConfig.whiteList}}</span>
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>http请求头</span>
                <span style="margin-left:10px">{{configData.webCrawlerConfig.headers}}</span>
            </div>
        </div>
        <!-- web路径爆破 -->
        <div v-if="configData.webPathScanConfig.isOpen" class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">web路径爆破</label>
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                <div>
                    <span >猜测速率</span>
                    <span style="margin-left:10px">{{configData.webPathScanConfig.guessRateZh}}</span>
                </div>
                <div>
                    <span >猜测时长</span>
                    <span style="margin-left:10px">{{configData.webPathScanConfig.guessTimeoutZh}}</span>
                </div>
               
                
            </div>
          
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>路径字典</span>
                <span  style="margin-left:10px" v-for="(it,id) in configData.webPathScanConfig.dickNames" :key="id">{{it}}</span>
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>排除标题黑名单</span>
                <span style="margin-left:10px">{{configData.webPathScanConfig.titleBlack}}</span>
            </div>
      
        </div>
        <!-- 子域名收集 -->
        <div v-if="configData.subdomainCollectConfig.isOpen" class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">子域名收集</label>
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                <div>
                    <span >子域名字典</span>
                    <span style="margin-left:10px">{{configData.subdomainCollectConfig.subdomainDictZh}}</span>
                </div>
            
            </div>
          
        </div>
         <!-- 口令爆破 -->
        <div v-if="configData.weakPassConfig.isOpen" class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">口令爆破</label>
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                <div>
                    <span >猜测次数</span>
                    <span style="margin-left:10px">{{configData.weakPassConfig.guessNumZh}}</span>
                </div>
                <div>
                    <span >猜测时间</span>
                    <span style="margin-left:10px">{{configData.weakPassConfig.guessTimeoutZh}}</span>
                </div>
                <div>
                    <span >猜测速率</span>
                    <span style="margin-left:10px">{{configData.weakPassConfig.guessRateZh}}</span>
                </div>
               
                
            </div>
          
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>爆破服务</span>
                <span style="margin-left:10px" v-for="(item,index) in configData.weakPassConfig.servicesZh" :key="index">{{item}}</span>
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>字典类型</span>
                <span style="margin-left:10px">{{configData.weakPassConfig.dictTypeZh}}</span>
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>账号字典</span>
                <span style="margin-left:10px">{{configData.weakPassConfig.commonUserDictZh}}</span>
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>密码字典</span>
                <span style="margin-left:10px">{{configData.weakPassConfig.commonPassDictZh}}</span>
            </div>
      
        </div>
        
        <!-- 任务优先级 -->
        <div class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">任务优先级</label>
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                
            </div>
          
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span style="visibility: hidden;">横向</span>
                <span style="margin-left:10px">{{configData2}}</span>
            </div>
            
      
        </div>
        <!-- 漏洞利用 -->
        <div class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">漏洞利用</label>
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                
            </div>
          
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span style="visibility: hidden;">横向</span>
                <span style="margin-left:10px">{{configData.vulExploit?'功能开启':'功能关闭'}}</span>
            </div>
            
      
        </div>
        <!-- 安全测试 -->
        <div class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">安全测试</label>
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                
            </div>
          
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span style="visibility: hidden;">横向</span>
                <span style="margin-left:10px">{{configData.safeTest?'功能开启':'功能关闭'}}</span>
            </div>
            
      
        </div>
        <!-- 横向移动（状态） -->
        <div class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">横向移动</label>
            </div>
            <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span style="visibility: hidden;">横向</span>
                <span style="margin-left:10px">{{(configData.lateralMove && configData.lateralMove.isOpen) ? '功能开启' : '功能关闭'}}</span>
            </div>
        </div>
        <!-- 网站登录凭证 -->
        <div v-if="configData.websiteLoginConfig.isOpen" class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">网站登录凭证</label>
            </div>
        <div v-if="configData.websiteLoginConfig&&configData.websiteLoginConfig.list">
                <div  v-for="(item,index) in configData.websiteLoginConfig.list" :key="index">
                <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
                <div>
                    <span >检测目标</span>
                    <span style="margin-left:10px">{{item.target}}</span>
                </div>
                <div>
                    <span >协议</span>
                    <span style="margin-left:10px">{{item.scheme}}</span>
                </div>
                <div>
                    <span >认证方式</span>
                    <span style="margin-left:10px">{{item.verifyStatusZh}}</span>
                </div>
               
                
                </div>
            
                <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                    <span>凭证</span>
                    <span style="margin-left:10px">{{item.verifyValue}}</span>
                </div>
            </div>
        </div>
             
            
      
        </div>
         <!-- local Storage -->
        <div v-if="configData.webCrawlerConfig.isOpen" class="box_bg_shadow" style="max-height: 410px;">
            <div style="border-left:3px solid #4c7ae3;margin-bottom:10px" >
                <label style="margin-left:10px;color:#484866A3">local Storage</label>
            </div>
             <div style="margin-bottom:10px;display:flex;justify-content: space-between;font-size:12px ;color:#484866A3">
            </div>
          
            <div v-if="configData.webCrawlerConfig&&configData.webCrawlerConfig.localStorage.list">
                  <div v-for="(item,index) in configData.webCrawlerConfig.localStorage.list" :key="index">
                 <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>key</span>
                <span style="margin-left:10px">{{item.key}}</span>
            </div>
            <div style="margin-bottom:10px;font-size:12px ;color:#484866A3">
                <span>value</span>
                <span style="margin-left:10px">{{item.value}}</span>
            </div>
           </div>
            </div>
      
        </div>
    </div>
</template>
<style lang="less" scoped>
/deep/.dv-scroll-ranking-board .ranking-column .shine{
   background: rgba(255,255,255,0);; 
}
.rankingWa{
    font-weight: 700;
    color: #35b0eb;
   
}
/deep/.rankingWa .ranking-column .inside-column{
    background-color: #35b0eb !important;
}
/deep/ .rankingWa .ranking-column {
    border-bottom: 2px solid #35b0eb;
}
/deep/ .rankingWa .ranking-info .rank {
    color: #484866DE;
}
.rankingWa2{
    color: #6c63f0;
    font-weight: 700;
}
/deep/.rankingWa2 .ranking-column .inside-column{
    background-color: #6c63f0 !important;
}
/deep/ .rankingWa2 .ranking-column {
    border-bottom: 2px solid #6c63f0;
}
/deep/ .rankingWa2 .ranking-info .rank {
    color: #484866DE;
}
.rankingWa3{
    font-weight: 700;
    color: #65c680;
   
}
    .lbtxt {
      display: block;
      font-size: 14px;
      margin: 10px 0;
      border-left: 3px solid #4c7ae3;
      padding-left: 10px;
    }
/deep/.rankingWa3 .ranking-column .inside-column{
    background-color: #65c680 !important;
}
/deep/ .rankingWa3 .ranking-column {
    border-bottom: 2px solid #65c680;
}
/deep/ .rankingWa3 .ranking-info .rank {
    color: #484866DE;
}
.rankingWa4{
    font-weight: 700;
    color: #4c7ae3;
   
}
/deep/.rankingWa4 .ranking-column .inside-column{
    background-color: #4c7ae3 !important;
}
/deep/ .rankingWa4 .ranking-column {
    border-bottom: 2px solid #4c7ae3;
}
/deep/ .rankingWa4 .ranking-info .rank {
    color: #484866DE;
}
.rankingWa5{
    font-weight: 700;
    color: #f6a623;
   
}
/deep/.rankingWa5 .ranking-column .inside-column{
    background-color: #f6a623 !important;
}
/deep/ .rankingWa5 .ranking-column {
    border-bottom: 2px solid #f6a623;
}
/deep/ .rankingWa5 .ranking-info .rank {
    color: #484866DE;
}
.rankingWa6{
    font-weight: 700;
    color: #f87d7d;
   
}
/deep/.rankingWa6 .ranking-column .inside-column{
    background-color: #f87d7d !important;
}
/deep/ .rankingWa6 .ranking-column {
    border-bottom: 2px solid #f87d7d;
}
/deep/ .rankingWa6 .ranking-info .rank {
    color: #484866DE;
}

ul{
    list-style: none;
}
.ul_box{
    padding-left: 70px;
    padding-top: 16px;
    max-height: 380px;
    overflow-y: auto;
    box-sizing: border-box;
    li{
        display: flex;
        font-size: 13px;
        color: rgba(72,72,102,0.64);
        label{
            flex: 1;
        }
        span{
            flex: 1;
        }
    }
}
.box_bg_shadow{
    border-radius: 4px;
    padding: 24px;
    background-color: #fff;
    box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12); 
    margin-bottom: 16px;
    .box_title{
        label{
            border-left: 2px solid #4C7AE3;
            padding-left: 8px;
            color: rgba(72,72,102,0.87);
            font-size: 13px;
        }
        .title_icon{
            display: inline-block;
            width: 32px;
            height: 32px; 
            margin-right: 16px;
            vertical-align: middle;
        }
        span{
            display: inline-block;
            font-size: 13px;
            color: rgba(72,72,102,0.87);
        }
    }
    .box_content{
        padding: 40px 0 16px;
        .risk_icon{
            display: inline-block;
            width: 48px;
            height: 48px; 
            margin-right: 24px;
            vertical-align: middle;
        }
        .risk_content{
            display: inline-block;
            vertical-align: middle;
            .risk_number{
                font-size: 18px;
                margin-bottom:8px; 
                line-height: 12px;
                color: rgba(72,72,102,0.87);
                label{
                    vertical-align: middle;
                    
                }
                i{
                    font-size: 12px;
                }
            }
            .risk_name{
                font-size: 13px;
                color: rgba(72,72,102,0.64);;
            }
        }
        .risk_content_height .risk_number{
            color: #F87D7D;
        }
        .risk_content_middle .risk_number{
            color: #F9B640;
        }
        .risk_content_low .risk_number{
            color: #4C7AE3;
        }
        .risk_content_safe .risk_number{
            color: #65C680;
        }
        
    }
    
}
.banner{
    .banner_icon{
        display: inline-block;
        width: 56px;
        height: 56px;
        background: url(../../assets/images/banner_icon@2x.png) ;
        background-size: cover;
        vertical-align: middle;
    }
    span{
        display: inline-block;
        vertical-align: middle;
        margin-left: 24px;
        >div{ 
            color: rgba(72,72,102,0.87); 
            font-size: 18px; 
            margin-bottom: 8px;
        }
        >p{
            font-size: 13px;
            color: rgba(72,72,102,0.87);
        }
    }
    
    
}
.tastBtnCont .search-box{
    padding: 16px 0 !important;
}
.task_detail{
    height: 226px; 
    background: #F7F8FB;
    box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12); 
    margin-bottom: 16px;
}
.basic_info{   
    padding: 24px 0 10px 24px; 
    box-sizing: border-box;
    ul{
        list-style: none;
    }
    li{
        margin-bottom: 20px;
        >div{
            min-width: 360px;
            display: inline-block; 
            label{ 
                display: inline-block;
                padding-left: 8px;
                border-left: 2px solid #4C7AE3;
                color: rgba(72,72,102,0.87);
                font-size: 13px;
                width: 90px;
                height: 14px ;
                line-height: 14px;
            }
            span{ 
                font-size: 13px;
                color: rgba(72,72,102,0.64);
            }
        }
    }
    
     
}
.target_risk{
    padding: 24px;
}
.target_detail{
    height: 118px;  
    background: #F7F8FB;
    box-shadow: inset 0px 2px 4px 0px rgba(76,122,227,0.12);
    color: rgba(72,72,102,0.64);
    padding: 24px;
    box-sizing: border-box;
    font-size: 13px;
}
#vuln_risk,
#vuln_status,
#vuln_obtain_evidence{
    height: 180px;
}
.legend{ 
    height: 70px;
    .name{
        color: rgba(72,72,102,0.64);
        font-size: 13px;
        i{
            display: inline-block;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            margin-right: 8px;
            vertical-align: middle;
        }
    }
    .number{
        font-size: 18px;
        margin: 10px 0 6px; 
        color: rgba(72,72,102,0.87);
    }
    .percentage{
        font-size: 13px;
        color: rgba(72,72,102,0.64);
    }
}
.vlun_risk_deadly { 
    .name i{
        background-color: #F87D7D;
    }
    .number{
        color: #F87D7D;
    }
}
.vlun_risk_height { 
    .name i{
        background-color: #F9B640;
    }
    .number{
        color: #F9B640;
    }
}
.vlun_risk_middle{
    .name i{
        background-color: @theme-BaseColor;
    }
    .number{
        color: @theme-BaseColor;
    }
}
.vlun_risk_low{
    .name i{
        background-color: #65C680;
    }
    .number{
        color: #65C680;
    }
}
.use_achievement{
    .name i{
        background-color: #6C63F0;
    }
}
.verify_existence{
    .name i{
        background-color: #3AACFF;
    }
}
.Not_verified{
    .name i{
        background-color: #FB466C;
    }
}
.Remote_control{
    .name i{
        background-color: #9A6BFF;
    }
}
.data_breach{
    .name i{
        background-color: #8ED5F1;
    }
}
.Logo_Credentials{
    .name i{
        background-color: #FF8787;
    }
}
</style>
<script>
import { task } from '@/api/task.js'
import BannerBox from "@/components/BannerBox.vue";
var echarts = require('echarts');
export default {
    name:'', 
    components:{
        BannerBox,
    },
    props:{
        task_name:{},
        task_id:{},
    },
    data(){
        return{
            configData:{
                portScanConfig:{
                    isOpen:false,
                },
                webCrawlerConfig:{
                    isOpen:false,
                },
                websiteLoginConfig:{
                    isOpen:false,
                },
                lateralMove:{
                    isOpen:false,
                },
                webPathScanConfig:{
                    titleBlack:'',
                    isOpen:false,
                },
                weakPassConfig:{
                    isOpen:false,
                },
                testIntensity:{
                    isOpen:false,
                },
                safeTest:{
                    isOpen:false,
                }, 
                subdomainCollectConfig:{
                    isOpen:false,
                }
            },
            configData2:{
                portScanConfig:{
                    isOpen:false,
                },
                webCrawlerConfig:{
                    isOpen:false,
                },
                websiteLoginConfig:{
                    isOpen:false,
                },
                lateralMove:{
                    isOpen:false,
                },
                webPathScanConfig:{
                    titleBlack:'',
                    isOpen:false,
                },
                weakPassConfig:{
                    isOpen:false,
                },
                testIntensity:{
                    isOpen:false,
                },
                safeTest:{
                    isOpen:false,
                },
                subdomainCollectConfig:{
                    isOpen:false,
                }
            },

            basic_info:{
                "taskName": "",//string,任务名称
                "executeType": 1,
                "executeTypeName": "",//string，执行方式
                "runtimePeriod": "",//string,执行时间段
                "riskLevel":0 ,
                "riskLevelName": "",//string,任务风险
                "taskTemplateId": 1,
                "taskTemplateName": "",//string,任务场景
                "status": 0,
                "statusName": "",//string,任务状态
                "createTime": "",//string,创建时间
                "updateTime": "",//string,结束时间
                cyclePlanningName:'',
                endTimeName:'',

            },
            runtimePeriod:'',
            executeTypeName:'', //执行方式
            ishowtip:false,
            target_risk:[],
            target_survival:[],
            risk_num:[],
            risk_total_num:0,
            risk_percentage:[],
            status_num:[],
            status_total_num:0,
            status_percentage:[],
            vulExploitImpact:[],
            // obtain_evidence_num:[],
            obtain_evidence_total_num:0,
            obtain_evidence:[],
            port_list:[],
            service_list:[],
            component_list:[],
            opSys_list:[],
            subDomain_list:[],
            urlTags_list:[],
            configChart:{
                data:  [ ],
            },
            configChart2:{
                data:  [ ],
            },
            configChart3:{
                data:  [ ],
            },
            configChart4:{
                data:  [ ],
            },
            configChart5:{
                data:  [ ],
            },
            configChart6:{
                data:  [ ],
            },
        }
    },
    created(){
        this.getConfigInfo();
    },
    mounted(){
        // this.getData(); 
    },
    methods:{ 
        async getConfigInfo(){
            const res = await task.configinfo({
                taskId:this.task_id,
            })
          if(res.code == 200){
            this.configData = res.data.config;
            this.configData2 = res.data.priorityZh;
            
            console.log(this.configData,'this.configData');
          }else{
            this.$message.error(res.message);
          }
        },
        async getData(){
            const res = await task.getTaskOverview({
                taskId:this.task_id,
            })
            if(res.code == 200){
                // localStorage.setItem('status', res.data.status);
                // localStorage.setItem('statusName', res.data.statusName); 
                this.vulExploitImpact = [];
                this.status_total_num = 0;
                let _this = this;
                // 基础信息
                this.basic_info.executeType = res.data.executeType;
                this.basic_info.executeTypeName = res.data.executeTypeName; 
                this.basic_info.cyclePlanningName = res.data.cyclePlanningName;
                this.basic_info.endTimeName =  res.data.endTimeName;
                //执行时间段
                this.runtimePeriod = res.data.runtimePeriod;
                // this.basic_info.runtimePeriod = res.data.runtimePeriod;
                let runtimePeriodarr = res.data.runtimePeriod.split(',');
                if(runtimePeriodarr.length >= 3){
                    this.ishowtip = true;
                    this.basic_info.runtimePeriod  = runtimePeriodarr[0]+','+runtimePeriodarr[1];
                }else{
                    this.ishowtip = false;
                    this.basic_info.runtimePeriod  = res.data.runtimePeriod;
                }
               

                this.basic_info.riskLevel = res.data.riskLevel;
                this.basic_info.riskLevelName = res.data.riskLevelName;
                this.basic_info.taskTemplateId = res.data.taskTemplateId;
                this.basic_info.taskTemplateName = res.data.taskTemplateName;
                this.basic_info.status = res.data.status;
                this.basic_info.statusName = res.data.statusName;
                this.basic_info.createTime = res.data.createTime;
                this.basic_info.updateTime = res.data.updateTime;
                // 目标风险
                this.target_risk = res.data.targetRisk ; 
                // 目标存活
                this.target_survival = res.data.targetNum;
                // 漏洞风险 - 添加空值检查
                this.risk_num = res.data.vulRisk || [];
                if (this.risk_num && this.risk_num.length > 0) {
                    this.risk_num.forEach(item =>{
                        _this.risk_total_num+=item
                    })  
                    this.risk_num.forEach(item =>{  
                        let item_num =_this.risk_total_num == 0? 0: (item/_this.risk_total_num*100).toFixed(2); 
                        _this.risk_percentage.push(item_num);
                    })
                }
                this.echart_vuln_risk(res.data.vulRisk || []);
                // 漏洞状态 - 添加空值检查
                if (res.data.vulExploitImpact && res.data.vulExploitImpact.length > 0) {
                    res.data.vulExploitImpact.forEach(item =>{
                        _this.status_total_num+=item.value;
                        this.status_num.push(item.value)
                    })  

                    res.data.vulExploitImpact.forEach(item =>{ 
                        let item_num =_this.status_total_num==0 ? 0: (item.value/_this.status_total_num*100).toFixed(2);
                        this.vulExploitImpact.push({label:item.label,value:item.value,percentage:item_num})
                    })
                }
                console.log(this.vulExploitImpact)
                this.echart_vuln_status(res.data.vulExploitImpact || []);
                // 漏洞取证 - 添加空值检查
                if (res.data.evidenceStat && res.data.evidenceStat.length > 0) {
                    res.data.evidenceStat.forEach(item =>{
                        // _this.statuobtain_evidence_total_nums_total_num+=item.value
                        // this.obtain_evidence_num.push(item.value)
                    })  
                }
                this.obtain_evidence = res.data.evidenceStat || [];
                this.echart_vuln_obtain_evidence(res.data.evidenceStat || []);
 
                // 端口数据 - 添加空值检查
                this.port_list = res.data.port || [];
                if (res.data.port && res.data.port.length > 0) {
                    let arrData = res.data.port.map(item => {
                        return {
                            name: item.label,
                            value: item.value
                        }
                    });
                    this.$set(this.configChart, 'data', arrData);
                } else {
                    this.$set(this.configChart, 'data', []);
                }
                
                // 服务数据 - 添加空值检查
                this.service_list = res.data.service || [];
                if (res.data.service && res.data.service.length > 0) {
                    let arrData2 = res.data.service.map(item => {
                        return {
                            name: item.label,
                            value: item.value
                        }
                    });
                    this.$set(this.configChart2, 'data', arrData2);
                } else {
                    this.$set(this.configChart2, 'data', []);
                }
                
                // 组件数据 - 添加空值检查
                this.component_list = res.data.component || [];
                if (res.data.component && res.data.component.length > 0) {
                    let arrData3 = res.data.component.map(item => {
                        return {
                            name: item.label,
                            value: item.value
                        }
                    });
                    this.$set(this.configChart3, 'data', arrData3);
                } else {
                    this.$set(this.configChart3, 'data', []);
                }
                
                // 系统数据 - 添加空值检查
                this.opSys_list = res.data.opSys || [];
                if (res.data.opSys && res.data.opSys.length > 0) {
                    let arrData4 = res.data.opSys.map(item => {
                        return {
                            name: item.label,
                            value: item.value
                        }
                    });
                    this.$set(this.configChart4, 'data', arrData4);
                } else {
                    this.$set(this.configChart4, 'data', []);
                }
                
                // 子域名数据 - 添加空值检查
                this.subDomain_list = res.data.subDomain || [];
                if (res.data.subDomain && res.data.subDomain.length > 0) {
                    let arrData5 = res.data.subDomain.map(item => {
                        return {
                            name: item.label,
                            value: item.value
                        }
                    });
                    this.$set(this.configChart5, 'data', arrData5);
                } else {
                    this.$set(this.configChart5, 'data', []);
                }
                
                // URL标签数据 - 添加空值检查
                this.urlTags_list = res.data.urlTags || [];
                if (res.data.urlTags && res.data.urlTags.length > 0) {
                    let arrData6 = res.data.urlTags.map(item => {
                        return {
                            name: item.label,
                            value: item.value
                        }
                    });
                    this.$set(this.configChart6, 'data', arrData6);
                } else {
                    this.$set(this.configChart6, 'data', []);
                }

            }else{
                // 处理API错误情况
                this.$message.error(res.message || '获取任务概览数据失败');
            }
        },
        echart_vuln_risk(vulRisk){
            let data = [];
            data.push({name:'致命漏洞',value:vulRisk[0]});
            data.push({name:'高危漏洞',value:vulRisk[1]});
            data.push({name:'中危漏洞',value:vulRisk[2]});
            data.push({name:'低危漏洞',value:vulRisk[3]});
            let totalvuln = vulRisk[0]+vulRisk[1]+vulRisk[2]+vulRisk[3];
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('vuln_risk')); 
            let option = {
                color:['#F87D7D','#F9B640','#4C7AE3','#65C680'],
                title:{
                    text:totalvuln,
                    subtext: '漏洞总数',
                    textStyle: {
                        fontsize: 18, 
                        fontWeight:400,
                        color: 'rgba(72,72,102,0.87)', 
                    },
                    subtextStyle: {
                        fontsize: 13,
                        color: 'rgba(72,72,102,0.64)', 
                    },
                    textAlign:"center", 
                    x:"49%", 
                    y:"37%" 
                    
                },
                series: [
                    {
                        name: ' ',
                        top:'top',
                        type: 'pie',
                        radius: ['50%', '70%'],
                        avoidLabelOverlap: false,
                        label: {
                            show: false,
                            position: 'center'
                        }, 
                        labelLine: {
                            show: false
                        },
                        data: data
                    }
                ]
            };
            // 绘制图表
            myChart.setOption(option);
        },
        echart_vuln_status(statuslist){
            let data = [];
            let totalvuln = 0;
            statuslist.forEach(item=>{
                data.push({name:item.label,value:item.value});
                totalvuln+=item.value
            });
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('vuln_status')); 
            let option = {
                color:['#6C63F0','#3AACFF','#FB466C' ],
                title:{
                    text:totalvuln,
                    subtext: '漏洞总数',
                    textStyle: {
                        fontsize: 18, 
                        fontWeight:400,
                        color: 'rgba(72,72,102,0.87)', 
                    },
                    subtextStyle: {
                        fontsize: 13,
                        color: 'rgba(72,72,102,0.64)', 
                    },
                    textAlign:"center", 
                    x:"49%", 
                    y:"37%" 
                    
                },
                series: [
                    {
                        name: ' ',
                        top:'top',
                        type: 'pie',
                        radius: ['50%', '70%'],
                        avoidLabelOverlap: false,
                        label: {
                            show: false,
                            position: 'center'
                        }, 
                        labelLine: {
                            show: false
                        },
                        data: data
                    }
                ]
            };
            // 绘制图表
            myChart.setOption(option);
        },
        echart_vuln_obtain_evidence(statuslist ){
            let data = [];
            let totalvuln = 0;
            statuslist.forEach(item=>{
                data.push({name:item.label,value:item.value});
                totalvuln+=item.value;
            });
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('vuln_obtain_evidence')); 
            let option = {
                color:['#9A6BFF','#8ED5F1','#FF8787'],
                title:{
                    text:totalvuln,
                    subtext: '漏洞取证',
                    textStyle: {
                        fontsize: 18, 
                        fontWeight:400,
                        color: 'rgba(72,72,102,0.87)', 
                    },
                    subtextStyle: {
                        fontsize: 13,
                        color: 'rgba(72,72,102,0.64)', 
                    },
                    textAlign:"center", 
                    x:"49%", 
                    y:"37%" 
                    
                },
                series: [
                    {
                        name: ' ',
                        top:'top',
                        type: 'pie',
                        radius: ['50%', '70%'],
                        avoidLabelOverlap: false,
                        label: {
                            show: false,
                            position: 'center'
                        }, 
                        labelLine: {
                            show: false
                        },
                        data: data
                    }
                ]
            };
            // 绘制图表
            myChart.setOption(option);
        },
    }
}
</script>
