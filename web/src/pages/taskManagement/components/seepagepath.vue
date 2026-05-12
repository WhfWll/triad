<template>
    <div class="target_box" id="con_div">
            <div class="path_info">
              <div style="position: relative; z-index: 2">
                <div class="path_row clearfix">
                  <div class="path_msg_item">
                    <i class="micon iconfont iconzhuangtai"></i>
                    <label for="" class="lbname">任务状态</label>
                    <span class="spvalue">{{ basicinfo.progress }}</span>
                  </div>
                  <div class="path_msg_item">
                    <i class="micon iconfont iconwangluo"></i>

                    <el-popover
                      placement="bottom-start"
                      width="320"
                      :visible-arrow="false"
                      trigger="click"
                    >
                      <div class="networkpopover">
                        <ul>
                          <label class="part_title">目标</label>
                          <li>
                            网络质量：<span
                              :class="[
                                {
                                  networkcolor0: network.target_quality == '优',
                                },
                                {
                                  networkcolor1: network.target_quality == '良',
                                },
                                {
                                  networkcolor2: network.target_quality == '中',
                                },
                                {
                                  networkcolor3: network.target_quality == '差',
                                },
                                {
                                  networkcolor0:
                                    network.target_quality == '断网',
                                },
                              ]"
                              >{{ network.target_quality }}</span
                            >
                          </li>
                          <li>
                            网络延时：<span
                              :class="[
                                {
                                  networkcolor0: network.target_quality == '优',
                                },
                                {
                                  networkcolor1: network.target_quality == '良',
                                },
                                {
                                  networkcolor2: network.target_quality == '中',
                                },
                                {
                                  networkcolor3: network.target_quality == '差',
                                },
                                {
                                  networkcolor0:
                                    network.target_quality == '断网',
                                },
                              ]"
                              >{{ network.target_delay }}</span
                            >
                          </li>
                        </ul>
                        <ul>
                          <label class="part_title">本地</label>
                          <li>
                            网络质量：<span
                              :class="[
                                {
                                  networkcolor0: network.local_quality == '优',
                                },
                                {
                                  networkcolor1: network.local_quality == '良',
                                },
                                {
                                  networkcolor2: network.local_quality == '中',
                                },
                                {
                                  networkcolor3: network.local_quality == '差',
                                },
                                {
                                  networkcolor0:
                                    network.local_quality == '断网',
                                },
                              ]"
                              >{{ network.local_quality }}</span
                            >
                          </li>
                          <li>
                            网络延时：<span
                              :class="[
                                {
                                  networkcolor0: network.local_quality == '优',
                                },
                                {
                                  networkcolor1: network.local_quality == '良',
                                },
                                {
                                  networkcolor2: network.local_quality == '中',
                                },
                                {
                                  networkcolor3: network.local_quality == '差',
                                },
                                {
                                  networkcolor0:
                                    network.local_quality == '断网',
                                },
                              ]"
                              >{{ network.local_delay }}</span
                            >
                          </li>
                        </ul>
                      </div>
                      <div
                        slot="reference"
                        style="cursor: pointer; display: inline-block"
                      >
                        <label class="lbname" style="cursor: pointer"
                          >网络</label
                        >
                        <span class="network" style="position: relative">
                          <i
                            v-for="count in 5"
                            :class="[count > network.network ? '' : 'wifig']"
                            :key="count"
                          >
                          </i>
                        </span>
                      </div>
                    </el-popover>
                  </div>
                  <div class="path_msg_item">
                    <i class="micon iconfont iconyuankong"></i>
                    <label for="" class="lbname">远控</label>
                    <span class="spvalue">{{
                      basicinfo.remote_control_number
                    }}</span>
                  </div>
                  <div style="float: right; width: 210px">
                    <div class="unfold" style="margin-right: 48px">
                      <label for="">简化路径</label>
                      <el-switch
                        v-model="isSimplifyPath"
                        @change="pathChange(isSimplifyPath)"
                      >
                      </el-switch>
                    </div>
                    <div class="unfold" @click="screen" >
                      <label for="">{{ showTxt }}</label>
                      <i class="iconfont iconquanping" v-if="!fullscreen"></i>
                      <i class="iconfont icontuichuquanping" v-else></i>
                    </div>
                  </div>
                </div>
              </div>
              <span class="iconcircular circular1"></span>
              <span class="iconcircular circular2"></span>
              <span class="iconcircular circular3"></span>
            </div>
            <div style="">
              <el-row :gutter="24">
                <el-col :span="12">
                  <div style="padding-left: 48px; height: 24px">
                    <el-row v-show="fullscreen">
                      <el-col :span="12">
                        <label class="part_title" style="margin-right: 24px"
                          >目标</label
                        >
                        <div class="networkpopover_1">
                          <label
                            >网络质量：<span
                              :class="[
                                {
                                  networkcolor0: network.target_quality == '优',
                                },
                                {
                                  networkcolor1: network.target_quality == '良',
                                },
                                {
                                  networkcolor2: network.target_quality == '中',
                                },
                                {
                                  networkcolor3: network.target_quality == '差',
                                },
                                {
                                  networkcolor0:
                                    network.target_quality == '断网',
                                },
                              ]"
                              >{{ network.target_quality }}</span
                            ></label
                          >
                          <label
                            >网络延时：<span
                              :class="[
                                {
                                  networkcolor0: network.target_quality == '优',
                                },
                                {
                                  networkcolor1: network.target_quality == '良',
                                },
                                {
                                  networkcolor2: network.target_quality == '中',
                                },
                                {
                                  networkcolor3: network.target_quality == '差',
                                },
                                {
                                  networkcolor0:
                                    network.target_quality == '断网',
                                },
                              ]"
                              >{{ network.target_delay }}</span
                            ></label
                          >
                        </div>
                      </el-col>
                      <el-col :span="12">
                        <label class="part_title" style="margin-right: 24px"
                          >本地</label
                        >
                        <div class="networkpopover_1">
                          <label
                            >网络质量：<span
                              :class="[
                                {
                                  networkcolor0: network.local_quality == '优',
                                },
                                {
                                  networkcolor1: network.local_quality == '良',
                                },
                                {
                                  networkcolor2: network.local_quality == '中',
                                },
                                {
                                  networkcolor3: network.local_quality == '差',
                                },
                                {
                                  networkcolor0:
                                    network.local_quality == '断网',
                                },
                              ]"
                              >{{ network.local_quality }}</span
                            ></label
                          >
                          <label
                            >网络延时：<span
                              :class="[
                                {
                                  networkcolor0: network.local_quality == '优',
                                },
                                {
                                  networkcolor1: network.local_quality == '良',
                                },
                                {
                                  networkcolor2: network.local_quality == '中',
                                },
                                {
                                  networkcolor3: network.local_quality == '差',
                                },
                                {
                                  networkcolor0:
                                    network.local_quality == '断网',
                                },
                              ]"
                              >{{ network.local_delay }}</span
                            ></label
                          >
                        </div>
                      </el-col>
                    </el-row>
                  </div>
                </el-col>
                <el-col :span="12">
                  <div style="text-align: right">
                    <div class="vulnstatistics">
                      <i class="port port1"></i>
                      <label for="">致命漏洞</label>
                      <span>{{ basicinfo.high_number }}</span>
                    </div>
                    <div class="vulnstatistics">
                      <i class="port port2"></i>
                      <label for="">高危漏洞</label>
                      <span>{{ basicinfo.middle_number }}</span>
                    </div>
                    <div class="vulnstatistics">
                      <i class="port port3"></i>
                      <label for="">中危漏洞</label>
                      <span>{{ basicinfo.low_number }}</span>
                    </div>
                    <div class="vulnstatistics">
                      <i class="port port4"></i>
                      <label for="">低危漏洞</label>
                      <span>{{ basicinfo.info_number }}</span>
                    </div>
                  </div>
                </el-col>
              </el-row>
            </div>

            <div
              style="width: 100%"
              id="svgdiv"
              :class="{ finished: basicinfo.progress === '已完成' }"
              v-loading="treeloading"
            >
              <svg
                id="svgCanvas"
                
                @contextmenu.prevent="contextmenuclick($event)"
                @mousedown.prevent="menuflag = false"
                @click.prevent="svgclick($event)"
              ></svg>
              <div
                id="myMenu"
                class="myMenu"
                v-if="menuflag"
                :style="menustyle"
                ref="myMenu"
                @scroll.stop="handleScroll($event)"
              >
                <div style="overflow-y: hidden;">
                  <!-- <div>{{nodeDetail}}</div> -->
                  <!-- 头部 -->
                  <div style="border-bottom:1px solid #4d79e4;display:flex;justify-content: space-between; align-items: center;">
                    <div>{{detailTYPE == 'webdirpathscan'?'web路径爆破':nodeDetail.title}}</div>
                    <div style="margin-bottom:5px;"><el-button size="small" @click="menuflag = false">关闭</el-button></div>
                  </div>
                  <!-- 头部 -->
                  <!-- 这个是图标 要加的 -->
                  <!-- <label class="name"
                    >{{nodeDetail.templateName}}
                    <span class="spcor1" 
                      >信息</span
                    >
                    <span class="spcor2" 
                      >漏洞验证</span
                    >
                    <span class="spcor3" 
                      >漏洞利用</span
                    >
                  </label> -->
             
                  <!-- type  root  -->
                  <div v-if="detailTYPE == 'root'">
                    <label for="" class="lbtxt">任务场景</label>
                    <div style="width: 100%">
                      <span>{{ nodeDetail.templateName }}</span>
                    </div>
                    <label v-if="detailTYPE == 'root'" for="" class="lbtxt">扫描端口：</label>
                    <div
                      v-if="detailTYPE == 'root'"
                      style="width: 100%"
                    >
                      <div style="color:#87879a;font-size:12px">{{ nodeDetail.portScan }}</div>
                    </div>
                  </div>
                  <!-- type  root  -->

                  <!-- type  portScan  -->
                  <div v-if="detailTYPE == 'portScan'">
                    
                    <label v-if="detailTYPE == 'portScan'" for="" class="lbtxt">扫描端口：</label>
                    <div
                      v-show="detailTYPE == 'portScan'"
                      style="width: 100%;display:flex;color:#87879a;font-size:12px"
                      v-for="(item, key) in nodeDetail.ports" 
                      :key="key"
                    >
                    
                    <div style="flex:1;">{{item.port}}</div>
                    <div style="flex:1;">{{item.service}}</div>
                    </div>
                  </div>
                  <div v-if="detailTYPE == 'whois'">
                    
                    <label v-if="detailTYPE == 'whois'" for="" class="lbtxt">whois：</label>
                    <div
                      v-show="detailTYPE == 'whois'"
                      style="width: 100%;display:flex;flex-direction:column;margin-bottom:20px;color:#87879a;font-size:12px"
                      v-for="(item, key) in nodeDetail.content" 
                      :key="key"
                    >
                    <div style="display:flex;" v-for="(item2, key2) in item.list" :key="key2">
                      <div style="flex:1;">{{item2.key}}</div>
                      <div style="flex:2;">{{item2.value}}</div>
                    </div>
                    
                    </div>
                  </div>
                  <!-- type  portScan  -->

                  <!-- type  port  -->
                  <div v-if="detailTYPE == 'port'">
                    <label for="" class="lbtxt">服务</label>
                    <div style="width: 100%">
                      <span>{{ nodeDetail.service }}</span>
                    </div>
                    <label v-if="detailTYPE == 'port'" for="" class="lbtxt">扫描端口：</label>
                    <div
                      v-show="detailTYPE == 'port'"
                      style="width: 100%;display:flex;"
                      v-for="(item, key) in nodeDetail.result" 
                      :key="key"
                    >
                    <div style="flex:1;">
                      <label>{{key}}：</label> 
                    </div>
                     <div style="flex:1;"> 
                       <span>{{item}}</span>
                      </div>
                    </div>
                  </div>
                  <!-- type  port  -->

                  <!-- type  vul  -->
                  <div v-if="detailTYPE == 'vul'">
                    <label for="" class="lbtxt">类型</label>
                    <div style="width: 100%">
                      <span>{{ nodeDetail.typeEnum }}</span>
                    </div>
                    <label for="" class="lbtxt">状态</label>
                    <div style="width: 100%">
                      <span>{{ nodeDetail.statusEnum }}</span>
                    </div>
                    <label for="" class="lbtxt">风险</label>
                    <div style="width: 100%">
                      <span>{{ nodeDetail.riskEnum }}</span>
                    </div>
                    <label for="" class="lbtxt">描述</label>
                    <div style="width: 100%">
                      <span>{{ nodeDetail.Description }}</span>
                    </div>
                    <label for="" class="lbtxt">建议</label>
                    <div style="width: 100%">
                      <span>{{ nodeDetail.fixSuggest }}</span>
                    </div>
                  </div>
                  <!-- type  vul  -->
                  <!-- type  cdn  -->
                  <div v-if="detailTYPE == 'cdn'">
                    <label for="" class="lbtxt">CDN</label>
                    <div style="width: 100%">
                      <span>{{ nodeDetail.resovleIP   }}</span>
                    </div>
                    
                  </div>
                  <!-- type  cdn  -->
                  <!-- type  webdirpathscan  -->
                  <div v-if="detailTYPE == 'webdirpathscan'">
                   <div v-for="(item,index) in nodeDetail" :key="index">
                     <label for="" class="lbtxt">{{item.title}}</label>
                      <div style="display:flex;justify-content: space-between; align-items: center;">
                        <span>{{ item.url   }}</span>
                        <span>{{ item.status   }}</span>
                      </div>
                   </div>
                    
                  </div>
                  <!-- type  webdirpathscan  -->


                </div>
              </div>
            </div>
          </div>
</template>

<script>

import * as d3 from 'd3'
import dagreD3 from 'dagre-d3'
import { task } from '@/api/task.js'
export default ({
  name: 'seepagepath',
  data () {
    return {
      detailTYPE: '',
      showTxt: '全屏',
      nodeDetail: { //右键节点详情
        name: '',
        description: '',
        loopdetail_port_scan: '',
        type: '',
        content: '',
        icon: '',
      },
      basicinfo: {
        break_system: "",
        check_target: "",
        high_number: 0,
        low_number: 0,
        middle_number: 0,
        info_number: 0,
        risk_level: "",
        use_time: "",
        progress: 0,
        create_time: '',
        remote_control_number: '',
        strenth: '',
        target_type: '',
        operating_system: '',
        geographic_location: ''
      },
      network: {
        local_delay: "",
        local_quality: "",
        target_delay: "",
        target_quality: "",
        network: '',
      },
      isSimplifyPath: false,
      fullscreen: false,
      yzloading: false,
      target_result: '',
      verify_result: [],
      pocname: '',
      target_result_id: '',
      task_Id: this.$route.query.taskId,
      task_name: this.$route.query.taskname,
      target_id: this.$route.query.id,
      target_name: this.$route.query.name,
      task_type: this.$route.query.tasktype,
      risk_level: this.$route.query.risk_level,
      tabs: this.$route.query.tabs,
      treeloading: true,
      show: false,
      activeName: 'tabs1',
      target_quality: '',
      target_delay: '',
      local_quality: '',
      local_delay: '',
      user_id: 0,
      menuflag: false,
      menucontent: {},
      menustyle: '',
      websocket: null,
  
      state: [],
      edg: [ ],
      diagGraph: {},
      flag: 0,
      currentPage: 1,
      total: 0,
      pageSize: 10,

      logtableData: [],
      targetpage_size: 50,
      timer1: null,
      timer2: null,
      statustimer: null,
      logpage: 1,
      drawer_attack: false,
      drawer_style: '',
      gridData: [],
      attackform: {
        task_id: '',
        attack_type: 1,
        remote_control: {
          is_open: false,
          listen_ip: '',
          listen_port: '',
          post_plan: '',
        },
        target: '',
        port: '',
        auth_type: 1,
        scheme: 1,
        service: '',
      },
      rules: {
        attack_type: [
          { required: true, message: '请选择攻击面类型', trigger: 'change' },
        ],
        port: [
          { required: true, message: '请输入端口', trigger: 'blur' },
        ],
        service: [
          { required: true, message: '请选择服务类型', trigger: 'change' },
        ],
        auth_type: [
          { required: true, message: '请选择认证方式', trigger: 'change' },
        ],
        path_type: [
          { required: true, message: '请选择敏感路径', trigger: 'change' },
        ]
      },
      attack_type_list: [{ id: 1, name: '服务' }, { id: 2, name: '敏感路径' }, { id: 3, name: '登录凭证' }],
      path_type_list: [{ id: 1, name: '登录入口' }, { id: 2, name: '文件上传' }],
      auth_type_list: [{ id: 1, name: 'cookie' }, { id: 2, name: 'headers' }],
      scheme_list: [{ value: 1, label: "HTTP" }, { value: 2, label: "HTTPS" }],
      services_list: [],
      looptableData: [],
      looptotalpage: 0,
      looppageSize: 10,
      loopcurrentpage: 1,
      loopformData: {
        page: 1,
        loopsearch: '',
      },
      loopexecute_dialog_title: '',
      loopexecute_dialog_desc: '',
      testform: {
        task_id: '',
        attack_type: 4,
        target: '',
        port: '',
        path: '',
        vul_id: '',
        cookie: '',
      },
      looprules: {
        target: [
          { required: true, message: '请输入目标IP', trigger: 'blur' },
        ],
        port: [
          { required: true, message: '请输入端口', trigger: 'blur' },
        ],
      },
      requestpack: '',
      responsepack: '',
      bugmessage: [],
      expands: [],
    }
  },
  created: function () {
    // this.$store.state.activefirstMenu = "/task"
    this.user_id = this.commonjs.decryptCBC(localStorage.getItem('user_id'), this.$commonjs.myKey)

    this.task_Id = this.$route.query.taskId
    this.task_name = this.$route.query.taskname
    this.target_id = this.$route.query.id
    this.target_name = this.$route.query.name
    this.tabs = this.$route.query.tabs
    this.risk_level = this.$route.query.risk_level


  },

  watch: {
    '$route': function (to, from) {
      this.flag = 0
      this.task_Id = this.$route.query.taskId
      this.task_name = this.$route.query.taskname
      this.target_id = this.$route.query.id
      this.target_name = this.$route.query.name
      this.tabs = this.$route.query.tabs
      this.risk_level = this.$route.query.risk_level

      if (this.$route.query.targetTabs) {
        this.activeName = this.$route.query.targetTabs
      }
      if (this.activeName == 'tabs1') {

      }

    },

  },
  mounted: function () {
    this.initdagreD3()
    this.getData();

    // 添加定时器,每5秒刷新一次
    this.timer = setInterval(() => {
      // 如果任务已完成,则清除定时器
      if(this.basicinfo.progress === '已结束') {
        if(this.timer) {
          clearInterval(this.timer)
        }
        return
      }
      this.getData()
    }, 5000)
  },
  // 组件销毁时清除定时器
  beforeDestroy() {
    if(this.timer) {
      clearInterval(this.timer)
    }
  },
  methods: {
    //-------------------------渗透路径 -----------------
    handleShowHide () {
      if (this.show) {
        this.showTxt = '全屏'
      } else {
        this.showTxt = '退出'
      }
      this.show = !this.show
    },
    initWebSocket () { //初始化weosocket     
      var wsurl = ''
      //    if (process.env.NODE_ENV === 'development') {
      //     wsurl = "ws://192.168.0.203:8011/remote/control/graph/"
      //     } else {
      //     wsurl = "ws://" + window.location.host + ':8011/remote/control/graph/'
      //     }

      if (process.env.NODE_ENV === 'development') {
        wsurl = "ws://" + process.env.VUE_APP_API_URL
      } else {
        wsurl = "ws://" + window.location.host
      }
      this.websocket = new WebSocket(wsurl)
      var that = this.websocket
      that.onopen = this.websocketonopen
      that.onerror = this.websocketonerror
      that.onmessage = this.websocketonmessage
      that.onclose = this.websocketclose

      // console.log(this.$router);
      this.$router.afterEach((to, from, next) => {
        // console.log(111)
        this.websocket.close()
      })
      // this.over = () =>{
      //     that.close()
      // }
    },
    initWebSocket1 () { //初始化weosocket  测试数据   
      let that = this
      this.treeloading = false
      that.state = [
        { id: "0", label: "192.168.0.68", class: "info", hide: false, status: "end", success: true },
        { id: "87c1c5ba-dfef-11ec-9981-0242ac130070", pocname: "device_os_info", label: "操作系统识别", class: "info", hide: false, status: "end", success: true },
        { id: "5662bcfe-dfef-11ec-9981-0242ac130070", pocname: "bf_redis", label: "redis弱口令漏洞", class: "info", hide: false, status: "end", success: true },
        { id: "1", label: "端口扫描", class: "info", hide: false, pocname: "port_scan", status: "end", success: true },
        { id: "50523830-dfef-11ec-9981-0242ac130070", label: "6379", class: "info", pocname: "port_scan", hide: false, status: "end", success: true },
      ]
      that.edg = [
        { start: "0", end: "87c1c5ba-dfef-11ec-9981-0242ac130070", hide: false },
        { start: "50523830-dfef-11ec-9981-0242ac130070", end: "5662bcfe-dfef-11ec-9981-0242ac130070", hide: false },
        { start: "0", end: "1", hide: false },
        { start: "1", end: "50523830-dfef-11ec-9981-0242ac130070", hide: false }
      ]
      that.treeTodagreD3()
    },
    websocketonopen () {
      this.websocket.send('ws_result.' + this.user_id + '.' + this.target_id)
    },
    websocketonerror (e) { //错误
      console.log("WebSocket连接发生错误")
    },
    websocketonmessage (event) { //数据接收  
      let that = this
      let node_message = JSON.parse(event.data)
      if (!node_message) return
      else if (node_message.sign == "tree") {
        this.treeloading = false
        // that.flag ++; 
        // if(that.flag == 1){
        that.state = node_message.state
        that.edg = node_message.edg
        that.treeTodagreD3()
        // }
        // else{
        //      that.diagGraph.changerender(node_message.state,node_message.edg);
        // }   
      }
      else if (node_message.sign == "info") {
        // this.basicinfo = node_message.result;
        // this.basicinfo.progress = parseInt(node_message.result.progress)
        if (node_message.result.hasOwnProperty('break_system')) {
          this.basicinfo.break_system = node_message.result.break_system
        }
        if (node_message.result.hasOwnProperty('check_target')) {
          this.basicinfo.check_target = node_message.result.check_target
        }
        if (node_message.result.hasOwnProperty('high_number')) {
          this.basicinfo.high_number = node_message.result.high_number
        }
        if (node_message.result.hasOwnProperty('low_number')) {
          this.basicinfo.low_number = node_message.result.low_number
        }
        if (node_message.result.hasOwnProperty('middle_number')) {
          this.basicinfo.middle_number = node_message.result.middle_number
        }
        if (node_message.result.hasOwnProperty('info_number')) {
          this.basicinfo.info_number = node_message.result.info_number
        }
        if (node_message.result.hasOwnProperty('progress')) {
          this.basicinfo.progress = node_message.result.progress
        }
        if (node_message.result.hasOwnProperty('risk_level')) {
          this.basicinfo.risk_level = node_message.result.risk_level
        }
        if (node_message.result.hasOwnProperty('use_time')) {
          this.basicinfo.use_time = node_message.result.use_time
        }
        if (node_message.result.hasOwnProperty('create_time')) {
          this.basicinfo.create_time = node_message.result.create_time
        }
        if (node_message.result.hasOwnProperty('remote_control_number')) {
          this.basicinfo.remote_control_number = node_message.result.remote_control_number
        }
        if (node_message.result.hasOwnProperty('strenth')) {
          this.basicinfo.strenth = node_message.result.strenth
        }
        if (node_message.result.hasOwnProperty('target_type')) {
          this.basicinfo.target_type = node_message.result.target_type
        }
        if (node_message.result.hasOwnProperty('operating_system')) {
          this.basicinfo.operating_system = node_message.result.operating_system
        }
        if (node_message.result.hasOwnProperty('geographic_location')) {
          this.basicinfo.geographic_location = node_message.result.geographic_location
        }
      }
      else if (node_message.sign == "network") {
        this.network = node_message
      }
    },
    //TODO:
    websocketclose (e) { //关闭 
      this.websocket.close()
      console.log("connection closed (" + e.code + ")")
    },
    // TODO:
    initdagreD3 () {
      let _that = this
      this.diagGraph = { //diag图数据操作 
        state: [],
        edg: [],
        statePoint: '',
        g: '',
        init: function (statePoint, state, edg) {
          this.statePoint = statePoint
          this.state = state
          this.edg = edg
          this.createG()
          this.renderG()
        },
        createG: function () {
          this.g = new dagreD3.graphlib.Graph()
            .setGraph({
              rankdir: 'LR', //设置方向
              // align: '', // 节点的对齐方式。有4个值: UL,UR,DL,DR。其中U是上(UP)，D是下(down)，L是左(left)，R是右(Right)
              ranksep: 100, // 每个层级间的距离
              edgesep: 10, // 在水平方向上, 线段间的距离
              nodesep: 20,// 水平方向上, 分隔节点的距离(节点之间的间距),
            })
            .setDefaultEdgeLabel(function () { return {} })

        },
        drawNode: function () {
          for (let i in this.state) { //画点
            let el = this.state[i]
            // if(!el.hide  || !_that.isSimplifyPath) {  
            if (el.hide && _that.isSimplifyPath) {

            }
            else {
              if (el.status == 'start') {
                this.g.setNode(el.id, {
                  id: el.id,
                  label: el.label || 'label',
                  // labelType:"html",
                  // label: el.label +'<i class=\"iconfont iconshouye\"></i>',
                  class: 'progress',
                  node_content: 'el.node_content',
                })
              }
              else {
                // console.log(444, el)
                this.g.setNode(el.id, {
                  id: el.id,
                  label: el.label || 'label',
                  // labelType:"html",
                  // label: el.label+'<i class=\"iconfont iconshouye\"></i> ',
                  class: el.class,
                  // style: style,
                  node_content: 'el.node_content',
                })

              }

            }



          }
          this.g.nodes().forEach((v) => { //画圆角
            var node = this.g.node(v)
            node.rx = node.ry = 20 // 所有节点的圆角
            // if (node.class === 'info') { // 灰色节点
            //   console.log(33, node)
            // }
          })
        },
        drawEdg: function () {
          for (let i in this.edg) { // 画连线
            let el = this.edg[i]
            // if(!el.hide  || !_that.isSimplifyPath) {  
            if (el.hide && _that.isSimplifyPath) {

            }
            else {

              this.g.setEdge(el.start, el.end, {
                style: "stroke: #c4c4ce; fill: none;",
                arrowheadStyle: 'fill: none;stroke: none;',
                arrowhead: 'undirected',
                curve: d3.curveBasis,

              })

            }

          }
        },
        addNode: function () {
          for (let i in this.state) {
            let el = this.state[i]
            if (el.hide && !_that.isSimplifyPath) {
              this.g.setNode(el.id, {
                id: el.id,
                label: el.label,
                class: el.class,
              })
            }
          }
          this.g.nodes().forEach((v) => { //画圆角
            var node = this.g.node(v)
            node.rx = node.ry = 20
          })
        },
        addEdge: function () {
          for (let i in this.edg) { //画点
            let el = this.edg[i]
            if (el.hide && !_that.isSimplifyPath) { //hide字段为true为可隐藏，为false则不可以
              this.g.setEdge(el.start, el.end, {
                style: "stroke: #c4c4ce; fill: none;",
                arrowheadStyle: 'fill: none;stroke: none;',
                arrowhead: 'undirected',
                curve: d3.curveBasis,
              })
            }

          }
          this.g.nodes().forEach((v) => { //画圆角
            var node = this.g.node(v)
            node.rx = node.ry = 20
          })
        },
        activeNode () {
          for (var n = 0; n < this.spot.length; n++) {
            for (let i in this.state) {
              let el = this.state[i]
              if (el.id === this.statePoint || this.spot[n] === el.id) {
                this.lineStyleClass = el.class
              }
              if (el.id == '0') {
                this.g.setNode(el.id, {
                  id: el.id,
                  label: el.label,
                  class: 'activeFnode',
                  style: 'fill: #4C7AE3;stroke: #4C7AE3;color:#fff',
                  node_content: 'el.node_content',
                })
              }
            }
          }


          this.g.nodes().forEach((v) => { //画圆角
            var node = this.g.node(v)
            node.rx = node.ry = 20
          })
        },
        activeEdg: function () { //高亮显示
          for (var n = 0; n < this.spot.length; n++) {
            for (let i in this.edg) { // 画连线
              let el = this.edg[i]
              if (el.end === this.statePoint || this.spot[n] === el.end) {
                let style = ''
                let styleColor = ''
                if (this.lineStyleClass == 'info') {
                  styleColor = '#D8D8D8'
                }
                else if (this.lineStyleClass == 'deadly') {
                  styleColor = '#B82B32'
                }
                else if (this.lineStyleClass == 'high') {
                  styleColor = '#e6a23c'
                }
                else if (this.lineStyleClass == 'middle') {
                  styleColor = '#4C7AE3'
                }
                if (this.lineStyleClass == 'low') {
                  styleColor = '#3DAF4E'
                }
                style = 'stroke: ' + styleColor + '; fill: none; stroke-width: 2px;'
                this.g.setEdge(el.start, el.end, {
                  //0fb2cc
                  style: style,
                  arrowheadStyle: 'fill: none;stroke: none;',
                  curve: d3.curveBasis,
                  arrowhead: 'undirected'
                })
              }
            }
          }
        },
        activeEdgNode (p) {
          for (let i in this.edg) {
            let el = this.edg[i]
            if (el.end === p) {
              this.spot.push(el.start)
              this.activeEdgNode(el.start)
            }
          }
        },
        removeNode () {
          for (let i in this.state) {
            let el = this.state[i]
            if (el.hide) {
              this.g.removeNode(el.id)
            }
          }
        },
        removeEdg () {
          for (var n = 0; n < this.spot.length; n++) {
            for (let i in this.edg) { // 画连线
              let el = this.edg[i]
              if (el.hide) {
                this.g.removeEdge(el.id)
              }

            }
          }
        },
        renderG: function () {
          var render = new dagreD3.render()
          var svg = d3.select("#svgCanvas")
          svg.select("g").remove() //删除以前的节点
          var svgGroup = svg.append("g")
          this.drawNode()
          this.drawEdg()
          try {
            //    this.g.nodes().forEach((v) => { //画圆角
            //         var node = this.g.node(v); 
            //         // node.rx = node.ry =20; // 所有节点的圆角
            //         // node.x = 0
            //         // node.width = 200
            //         // console.log(33, node)
            //     });
            render(d3.select("svg g"), this.g) //渲染节点

            this.setRender(svg)
          } catch (err) {
            console.log(err)
          }


        },
        setRender (svg) {
          var inner = svg.select("g")
          var zoom = d3.zoom().on("zoom", function () { //添加鼠标滚轮放大缩小事件
            // console.log(d3)
            // console.log(d3.event)
            // console.log('32', d3.event.transform)
            inner.attr("transform", d3.event.transform)
            _that.menuflag = false
          })
          svg.call(zoom)
          var max = svg._groups[0][0].clientWidth > svg._groups[0][0].clientHeight ? svg._groups[0][0].clientWidth : svg._groups[0][0].clientHeight
          // var initialScale =0.8;
          var initialScale = max / 2000
          var tWidth = (svg._groups[0][0].clientWidth - this.g.graph().width * initialScale) / 2
          var tHeight = (svg._groups[0][0].clientHeight - this.g.graph().height * initialScale) / 2

          var trans = d3.zoomIdentity.translate(50, tHeight).scale(initialScale)
          svg.call(zoom.transform, trans) //元素居中 
        },
        changePoint: function (point) {
          this.spot = []
          this.statePoint = point
          this.activeEdgNode(this.statePoint) //获得节点 

          this.drawNode()//画点
          this.drawEdg()// 画连线 
          this.activeNode() //高亮节点class 
          this.activeEdg() //高亮线  

          var render = new dagreD3.render()
          render(d3.select("svg g"), this.g) //渲染节点

        },
        changerender (state, edg) {
          this.state = state
          this.edg = edg
          var render = new dagreD3.render()
          var svg = d3.select("#svgCanvas")
          // svg.select("g").remove(); //删除以前的节点
          var inner = svg.select("g")
          var zoom = d3.zoom().on("zoom", function () { //放大
            inner.attr("transform", d3.event.transform)
          })
          svg.call(zoom)
          this.drawNode()
          this.drawEdg()
          render(d3.select("svg g"), this.g) //渲染节点

        },
        showhidechangerender () {
          var render = new dagreD3.render()
          var svg = d3.select("#svgCanvas")
          //开启简化路径是 true，关闭是false
          if (_that.isSimplifyPath) {
            this.removeNode()
            this.drawEdg()// 画连线 
          } else {
            this.addNode()
            this.addEdge()
          }
          render(d3.select("svg g"), this.g) //渲染节点

          this.setRender(svg)

        },
      }

    },
    svgclick (e) {
      this.menuflag = false
      e.preventDefault()
      if (e.target.tagName === 'rect') {
        this.diagGraph.changePoint(e.target.parentNode.id)
      } else {
        this.diagGraph.changePoint(0)
      }
    },
    treeTodagreD3 () {
      var statePoint = 0 // 当前选中的点  
      try {
        this.diagGraph.init(statePoint, this.state, this.edg) //创建关系图 
      } catch (err) {
        console.log(err)
      }

    },
    pathChange () { //简化路径切换 
      this.diagGraph.showhidechangerender()
    },
    contextmenuclick (e,o) { //右键弹出层 
      let that = this
      let detailID = ''
      console.log(e,'e.target',that.state);
      this.state.forEach((item)=>{
        if(item.id == e.target.__data__){
          detailID = item.detailId
        }
        })
      if (e.target.tagName === 'rect' || e.target.tagName === 'text') {
        var value = e.target.__data__
        this.$ajax.get('/smart/task/tasktargetmapnodedetail', {
          params: {
            detailId:detailID,
          }
        })
          .then((data) => {
            var dt = data.data
            if (dt.code == 200) {
              this.detailTYPE = dt.data.type
              // console.log(dt,'dt');
              this.menustyle = 'top:0px;left:0px'
               if(dt.data.type == "portScan"){
                // console.log("portScanportScanportScan");
                this.nodeDetail = dt.data.data
              }else if(dt.data.type == "root"){
                // console.log("rootrootrootroot");
                this.nodeDetail = dt.data.data
              // this.menucontent = dt.content
              }else if(dt.data.type == "port"){
                // console.log("portportportportport");
                this.nodeDetail = dt.data.data
                if(this.nodeDetail.result){
                  this.nodeDetail.result = JSON.parse(this.nodeDetail.result)
                }
              }else if(dt.data.type == "vul"){
                // console.log("vulvulvulvulvul");
                this.nodeDetail = dt.data.data
              }else if(dt.data.type == "whois"){
                // console.log("vulvulvulvulvul");
                this.nodeDetail = {
                  title: dt.data.data.title,
                  content: []
                }
                dt.data.data.dataList.forEach((jsonItem, index) => {
                  let list = []
                  for (let key in jsonItem) {
                    list.push({
                      key: key,
                      value: jsonItem[key]
                    })
                  }
                  this.nodeDetail.content.push({
                    list: list
                  })
                })
              }else{
                this.nodeDetail = dt.data.data
              }
              this.menuflag = true
          
             
              // this.loopdetailname = dt.name;
              // this.nodeDetail.desc = dt.description;
              // this.loopdetail_port_scan = dt.port_scan;
              let clientY = e.clientY
              let clientX = e.clientX
              let svgCanvas = document.getElementById('svgdiv')
              let svgh = svgCanvas.offsetHeight //canvas高
              let svgw = svgCanvas.offsetWidth  //canvas宽


              this.$nextTick(() => {
                var menu = document.getElementById('myMenu')
                let menuw = menu.offsetWidth //弹层宽
                let menuh = menu.offsetHeight  //弹出高
                let style = ''
                if ((clientY - 280 + menuh) > svgh && (clientY - 280 - menuh) > 0) {
                  style = 'top:' + (clientY - menuh) + 'px;'
                } else {
                  style = 'top:' + (clientY) + 'px;'
                }
                if ((clientX - 200 + menuw) > svgw && (clientX - 200 - menuw)) {
                  style += 'left:' + (clientX - menuw) + 'px;'
                } else {
                  style += 'left:' + (clientX) + 'px;'
                }

                that.menustyle = style
              })

            } else {
              this.$message({
                message: dt.msg,
                type: 'error'
              })
            }
          })
          .catch(function (error) {
            console.log(error)
          })

      } else {
        this.menuflag = false
      }

    },
    contextmenuclick1 (e) {
      // console.log(e);
      if (e.target.tagName === 'rect' || e.target.tagName === 'text') {
        var value = e.target.__data__
   
        // console.log(value);
        for (let index = 0; index < this.state.length; index++) {
          const element = this.state[index]
          if (value == element.id) {
            if (element.node_value && element.node_value.web_info) {
              this.menustyle = 'top:' + e.clientY + 'px;left:' + e.clientX + 'px'
              this.menuflag = true
              this.menucontent = element.node_value.web_info
              break
            }
          }

        }
      }
    },
    handleScroll (e) { //滚动，清空右键东西
      // this.menuflag=false; 
      if (e.target.id && e.target.id == 'myMenu') {

      } else {
        this.menuflag = false
      }
    },
    screen () { //全屏
      // let element = document.documentElement;//设置后就是我们平时的整个页面全屏效果
      let element = document.getElementById('con_div')//设置后就是 id==con_lf_top_div 的容器全屏

      if (this.fullscreen) {
        // 如果已经全屏了就退出全屏
        if (document.exitFullscreen) {
          document.exitFullscreen()
        }
        else if (document.webkitCancelFullScreen) {
          document.webkitCancelFullScreen()
        } else if (document.mozCancelFullScreen) {
          document.mozCancelFullScreen()
        } else if (document.msExitFullscreen) {
          document.msExitFullscreen()
        }
        this.fullscreen = false
        this.showTxt = '全屏'
      } else {
        // 如果不是全屏就变成全屏
        if (element.requestFullscreen) {
          element.requestFullscreen()
        } else if (element.webkitRequestFullScreen) {
          element.webkitRequestFullScreen()
        } else if (element.mozRequestFullScreen) {
          element.mozRequestFullScreen()
        } else if (element.msRequestFullscreen) {
          // IE11
          element.msRequestFullscreen()
        }
        this.fullscreen = true
        this.showTxt = '退出'
      }
    },
    //-------------------------渗透路径-end---------------- 
    async getSeriveList () {
      // const res = await attack.attackService()
      // if (res.success) {
      //   this.services_list = res.result
      // } else {
      //   this.$message({
      //     message: res.error,
      //     type: 'error'
      //   })
      // }
    },

    // TODO:
    async getData (tabname) {
      let that = this
      const res = await task.tasktargetmap({
        targetId: this.target_id
      })
      if (res.code == 200) {
        this.treeloading = false
        this.basicinfo.remote_control_number = res.data.remoteControlNum
        this.basicinfo.progress = res.data.progress

        this.basicinfo.high_number = res.data.risk.deadly
        this.basicinfo.middle_number = res.data.risk.high
        this.basicinfo.low_number = res.data.risk.middle
        this.basicinfo.info_number = res.data.risk.low

        this.network = res.data.network;
        let state_Arr = res.data.state||[];
        //遍历state,如果id为空则删除这一项
        for (let i = 0; i < state_Arr.length; i++) {
            if (state_Arr[i].id == '') {
              state_Arr.splice(i, 1)
              i--
            }
          }
          that.state = res.data.state || []
          that.edg = res.data.edg ||[]
          if(that.state.length>0 && that.edg.length>0){
            that.treeTodagreD3()
          }

        // 如果任务已完成,清除定时器
        if(this.basicinfo.progress === '已结束') {
          if(this.timer) {
            clearInterval(this.timer)
          }
        }
          
      } else {
        this.$message({
          message: res.error,
          type: 'error'
        })
      }
      console.log(res.data, 'res----------')


    },

   
  }
})
</script>
<style>
.examplecode {
  color: #4c7ae3;
  padding-left: 10px;
  margin: 10px 0;
  font-style: italic;
}
.link_default {
  padding-left: 10px;
}
.link_default span {
  display: inline-block;
  min-width: 26px;
  text-align: center;
  color: rgba(72, 72, 102, 0.32);
}
.packtxt textarea {
  resize: none !important;
}
</style>
<style lang="less" scoped>
.loglist {
  list-style: none;
  margin-top: 24px;

  li {
    color: rgba(72, 72, 102, 0.64);
    font-size: 13px;
  }
}
.network i {
  display: inline-block;
  // width: 6px;
  // height: 6px;
  // background-color: #ebeef5;
  // background: #15C53D;
  margin: 0 2px;
}
.network i.wifig {
  background: #fff;
}
.network {
  i {
    position: absolute;
    width: 4px;
    background: rgba(255, 255, 255, 0.7);
    &:nth-child(1) {
      top: 10px;
      left: 0;
      height: 4px;
    }
    &:nth-child(2) {
      top: 8px;
      left: 6px;
      height: 6px;
    }
    &:nth-child(3) {
      top: 6px;
      left: 12px;
      height: 8px;
    }
    &:nth-child(4) {
      top: 4px;
      left: 18px;
      height: 10px;
    }
    &:nth-child(5) {
      top: 2px;
      left: 24px;
      height: 12px;
    }
  }
}
/deep/ .el-switch.is-checked .el-switch__core {
  border-color: #fff;
  background-color: #fff;
}
/deep/ .el-switch__core:after {
  background-color: #4c7ae3;
}
#svgdiv {
  height: calc(100% - 180px);
}
.vulnstatistics {
  font-size: 13px;
  display: inline-block;
  margin-right: 54px;
  i {
    margin-right: 5px;
    vertical-align: baseline;
  }
  label {
    color: rgba(72, 72, 102, 0.64);
    margin-right: 16px;
  }
  span {
    color: rgba(72, 72, 102, 0.87);
  }
}
.port {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 4px;
}
.port1 {
  background: #f87d7d;
}
.port2 {
  background: #f9b640;
}
.port3 {
  background: #4c7ae3;
}
.port4 {
  background: #15c53d;
}
.port1,
.port2,
.port3,
.port4 {
  vertical-align: super;
}
// .iconduocijiance{
//     vertical-align: bottom;
// }
.tag_status {
  width: auto;
  padding: 0 8px;
}
// .selstatus /deep/ .el-input__inner{
//     height:40px !important;
// }
.updatestatus {
  /deep/ .el-dialog__body {
    height: 192px !important;
  }
  /deep/ .el-dialog {
    height: auto !important;
  }
  /deep/ .el-dialog__body {
    padding: 72px 152px !important;
  }
}

.rotateZ {
  transform: rotateZ(180deg);
}
/deep/ .el-dropdown-menu__item:hover {
  background: #f7f7fb !important;
}
.notarget {
  color: rgba(72, 72, 102, 0.32);
  cursor: default;
  min-width: 160px;
  box-sizing: border-box;
}
.notarget:hover {
  color: rgba(72, 72, 102, 0.32);
  background: #fff !important;
}
.targetitem {
  min-width: 200px;
  span {
    i {
      display: inline-block;
      width: 8px;
      height: 8px;
      border-radius: 50%;
    }
    margin-right: 16px;
    width: 55px;
    display: inline-block;
    height: 17px;
    line-height: 17px;
  }
  .target_url {
    text-decoration: none;
    color: rgba(72, 72, 102, 0.64);
    display: inline-block;
    // min-width:100%;
  }
  a:hover {
    background: #f7f7fb !important;
  }
  label {
    display: inline-block;
    overflow: hidden;
    white-space: nowrap;
    text-overflow: ellipsis;
    width: 110px;
    height: 17px;
    line-height: 17px;
    vertical-align: sub;
    cursor: pointer;
  }
}
.targetitem:hover {
  a {
    color: #4c7ae3 !important;
  }
  background: #f7f7fb !important;
}
.ykbg {
  display: inline-block;
  width: 14px;
  height: 14px;
  background-size: cover;
  vertical-align: middle;
  margin-left: 16px;
}
.yes {
}
.textarea /deep/ textarea {
  resize: none !important;
}
.targetaddress {
  font-size: 13px;
  font-weight: 500;
  color: rgba(72, 72, 102, 0.64);
  line-height: 14px;
}
.time {
  font-size: 13px;
  font-weight: 500;
  color: rgba(72, 72, 102, 0.64);
  line-height: 14px;
  margin-top: 8px;
  margin-bottom: 8px;
}
.usebox {
  margin-top: 16px;
  color: rgba(72, 72, 102, 0.64);
  .cmdresult {
    padding: 16px 0;
    word-wrap: break-word;
    word-break: normal;
  }
}
.controlinput {
  width: 90% !important;
  box-sizing: border-box;

  /deep/ .el-input__inner {
    border: none !important;
    padding-left: 0;
  }
}
.controlbox {
  margin-top: 16px;
  color: rgba(72, 72, 102, 0.64);
  .cmdresult {
    padding: 16px 0;
    word-wrap: break-word;
    word-break: normal;
  }
}
.useinput {
  width: 90% !important;
  box-sizing: border-box;

  /deep/ .el-input__inner {
    border: none !important;
    padding-left: 0;
  }
}
/deep/ .el-progress-bar {
  width: 110px !important;
}
/deep/ .el-table td:not(.el-table-column--selection):first-child .cell,
/deep/ .el-table th:not(.el-table-column--selection):first-child .cell {
  padding-left: 32px !important;
}
.bugbasicinfo /deep/ .el-table--enable-row-transition .el-table__body td {
  transition: none;
}
.bugbasicinfo /deep/ .el-table--enable-row-hover .el-table__body tr:hover > td {
  background: none !important;
}
.iconcircular {
  position: absolute;
  border-radius: 50%;
}
.circular1 {
  width: 30px;
  height: 30px;
  background: #3768d6;
  box-shadow: 0px -2px 2px 0px #3157af;
  bottom: -19px;
  left: 131px;
}
.circular2 {
  width: 86px;
  height: 86px;
  box-shadow: 0px 2px 2px 0px #2e5dc8;
  border: 20px solid #3768d6;
  top: -82px;
  left: 50%;
}
.circular3 {
  width: 67px;
  height: 67px;
  background: #3768d6;
  box-shadow: 0px -2px 2px 0px #3157af;
  bottom: -44px;
  right: 100px;
}

/deep/ .el-progress__text {
  color: #fff;
  font-size: 13px !important;
}
/deep/ .el-progress-bar__outer {
  background: rgba(255, 255, 255, 0.3);
}
.part_title {
  font-size: 14px;
  margin-bottom: 16px;
  font-weight: 800;
  border-left: 3px solid #4c7ae3;
  padding-left: 10px;
  height: 14px;
  line-height: 14px;
  color: rgba(72, 72, 102, 0.87);
}
@media (max-width: 1440px) {
  /deep/ .el-dialog {
    height: calc(100% - 96px);
  }
}
@media (min-width: 1440px) {
  /deep/ .el-dialog {
    height: calc(100% - 176px);
  }
}
/deep/ .el-tabs__item {
  height: 48px;
  line-height: 48px;
  padding: 0 24px;
}
/deep/ .el-tabs__item.is-active {
  color: #4c7ae3;
  font-weight: 500;
}
/deep/ .el-tabs__nav-wrap {
  padding: 0 24px;
}
/deep/ .el-tabs__nav-wrap::after {
  background: #e8e8f5;
  height: 1px;
}
/deep/ .el-tabs__header {
  margin: 0 0 24px;
}
.targetBluebg {
  background: #4c7ae3;
  font-size: 13px;
  font-weight: 500;
  border-radius: 4px;
  padding: 24px;
  box-sizing: border-box;
  button {
    margin-right: 48px;
  }
}
.path_info {
  background: #4c7ae3;
  margin-bottom: 24px;
  border-radius: 4px;
  padding: 24px 24px 0;
  box-sizing: border-box;
  font-size: 13px;
  font-weight: 500;
  position: relative;
  overflow: hidden;
  .path_row {
    width: 100%;
    .path_msg_item {
      float: left;
      width: 200px;
      margin-bottom: 24px;
      .micon {
        color: rgba(255, 255, 255, 0.7);
        vertical-align: top;
        height: 20px;
        line-height: 19px;
        display: inline-block;
      }
      .lbname {
        display: inline-block;
        // width: 54px;
        margin-right: 18px;
        margin-left: 8px;
        color: rgba(255, 255, 255, 0.7);
        height: 20px;
        line-height: 20px;
        vertical-align: middle;
      }
      .spvalue {
        width: calc(100% - 100px) !important;
        display: inline-block;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        color: #fff;
        vertical-align: text-top;
        height: 20px;
        line-height: 20px;
        position: relative;
        vertical-align: middle;
      }
      .spvaluerisklevel {
        position: absolute;
        display: inline-block;
        // width: 40px;
        padding: 0 8px;
        height: 20px;
        line-height: 20px;
        text-align: center;
        background: #fff;
        border-radius: 20px;
        z-index: 10;
        font-size: 12px;
        box-sizing: border-box;
      }
      .tag_high {
        color: #f87d7d !important;
      }
      .tag_primary {
        color: #4c7ae3 !important;
      }
      .tag_success {
        color: #15c53d !important;
      }
      .tag_warning {
        color: #f9b640 !important;
      }
    }
    .path_bug_item {
      float: left;
      width: 23%;
      margin-bottom: 24px;
      i {
        display: inline-block;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        margin-left: 4px;
        margin-right: 4px;
      }
      label {
        display: inline-block;
        font-size: 13px;
        color: rgba(255, 255, 255, 0.7);
        font-weight: 500;
        padding: 0 16px 0 8px;
      }
      span {
        display: inline-block;
        font-size: 13px;
        font-weight: 500;
        color: #fff;
      }
    }
    .unfold {
      float: left;
      // width: 8%;
      color: #ffffff;
      text-align: right;
      label {
        margin-right: 10px;
      }
      cursor: pointer;
      * {
        cursor: pointer;
      }
    }
  }
}
.networkpopover > ul {
  list-style: none;
  display: inline-block;
  width: 49%;
}
.networkpopover > ul li {
  margin-bottom: 8px;
  padding-left: 13px;
  font-size: 13px;
}
/*优*/
.networkcolor0 {
  color: #67c23a;
}
/*良*/
.networkcolor1 {
  color: #5cb85c;
}
/*中*/
.networkcolor2 {
  color: #e6a23c;
}
/*差*/
.networkcolor3 {
  color: #f56c6c;
}
/*断网*/
.networkcolor4 {
  color: #909399;
}
.part_title {
  display: inline-block;
  font-size: 14px;
  margin-bottom: 16px;
  font-weight: 800;
  border-left: 3px solid #4c7ae3;
  padding-left: 10px;
  height: 14px;
  line-height: 14px;
  color: rgba(72, 72, 102, 0.87);
}
.myMenu {
  position: fixed;
  padding: 14px 14px;
  border-radius: 4px;
  background: #fff;
  box-shadow: 0px 4px 8px 0px rgba(72, 72, 102, 0.32);
  max-height: 300px;
  max-width: 500px;
  min-width: 300px;
  word-wrap: break-word;
  overflow-y: auto;
  overflow-x: hidden;
  z-index: 99999;
  > div {
    font-size: 16px;
    margin-bottom: 10px;
    .name {
      display: block;
      font-size: 14px;
      width: 100%;
      margin-bottom: 10px;
      border-left: 3px solid #4c7ae3;
      padding-left: 10px;
      padding-right: 10px;
      box-sizing: border-box;
      span {
        display: block;
        float: right;
        font-size: 12px;
        border-radius: 10px;
        color: #fff;
        padding: 2px 8px;
      }
      .spcor1 {
        background: #09c1f7;
      }
      .spcor2 {
        background: #15c53d;
      }
      .spcor3 {
        background: #f35f28;
      }
    }
    .lbtxt {
      display: block;
      font-size: 14px;
      margin: 10px 0;
      border-left: 3px solid #4c7ae3;
      padding-left: 10px;
    }
    p {
      font-size: 12px;
    }
    div {
      label {
        display: inline-block;
        min-width: 60px;
        font-size: 12px;
        padding-left: 10px;
        height: 14px;
        line-height: 14px;
        color: rgba(72, 72, 102, 0.64);
        // vertical-align: sub;
      }
      span {
        font-size: 12px;
        display: inline-block;
        // width: calc(100% - 130px);
        color: rgba(72, 72, 102, 0.64);
        vertical-align: text-top;
      }
    }
  }
  > div:last-child {
    margin-bottom: 0;
  }
}

.svgCanvas /deep/ g.node rect {
  padding: 5px;
  fill: #dbe4f9;
  stroke: #4c7ae3;
  stroke-width: 1px;
}
.svgCanvas /deep/ g.node text {
  color: #4c7ae3;
  font: 14px sans-serif;
  /*font-weight:700;*/
  font-size: 12px;
  cursor: pointer;
}
#svgCanvas {
  height: 100%;
  min-height: 590px;
  width: 100%;
}
.svgCanvas .node text {
  font-weight: bold;
  font-family: "Microsoft YaHei";
  font-size: 12px;
  pointer-events: none;
  text-anchor: middle;
  fill: white;
}
.svgCanvas .label g {
  transform: translate(0, -13px);
}

.svgCanvas /deep/ .node rect {
  fill: white;
  stroke-width: 0px;
  color: white;
}
.svgCanvas g > rect {
  stroke: #dee6e8;
  stroke-width: 1px;
}
.svgCanvas g text {
  fill: #687386 !important;
  font-size: 14px !important;
}
.svgCanvas /deep/ .edgePath path {
  stroke: #ccc;
  stroke-width: 1px;
}
/deep/ .el-tabs__header {
  margin: 0 0 15px;
  background: #fff;
  border-radius: 4px;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
  border: none;
}
.targetlist {
  min-height: calc(100% - 39px);
  box-sizing: border-box;
  // background: #fff;
  // box-shadow:0px 2px 4px 0px rgba(76,122,227,0.12);
  position: relative;
}
.context_box_bg {
  background: none;
}
.target_box {
  padding: 24px 24px;
  box-sizing: border-box;
  position: relative;
  height: 100%;
  background: #fff;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
  border-radius: 4px;
  margin-bottom: 10px;
  .el-table__body-wrapper {
    height: calc(100% - 54px);
  }
}
.networkpopover_1 {
  display: inline-block;
}
.networkpopover_1 > label {
  margin-bottom: 8px;
  padding-left: 13px;
  font-size: 13px;
  color: rgba(72, 72, 102, 0.64);
  margin-right: 24px;
}
.targettypelist {
  position: absolute;
  right: 24px;
  top: 8px;
  z-index: 9;
  .el-tabs {
    height: 100%;
    position: absolute;
    .el-tabs__content {
      height: calc(100% - 70px);
      .el-tab-pane {
        height: 100%;
      }
    }
  }
}
.el-tabs__content {
  overflow: auto !important;
}
.el-dropdown-link {
  background: rgba(76, 122, 227, 0.12);
  display: inline-block;
  height: 32px;
  width: 160px;
  line-height: 32px;
  padding-left: 12px;
  box-sizing: border-box;
  font-weight: 500;
  color: rgba(76, 122, 227, 1);
  cursor: pointer;
  outline: none;
  i {
    float: right;
    margin-right: 12px;
    // margin-top: 10px;
  }
}
.attack_tabs > div {
  float: left;
  height: 32px;
  line-height: 32px;
  border-radius: 2px;
  border: 1px solid rgba(232, 232, 245, 1);
  padding: 0 24px;
  font-weight: 500;
  color: rgba(72, 72, 102, 0.64);
  font-size: 14px;
  margin-left: -1px;
  -moz-user-select: none;
  -khtml-user-select: none;
  user-select: none;
}
.attack_tabs > div.active {
  background: #4c7ae3;
  color: #fff;
  border: 1px solid #4c7ae3;
}
.attack_tabs_content {
  margin-top: 24px;
}
.target_box /deep/ .el-tabs__nav-wrap {
  padding: 0;
}
.target_box /deep/ .el-tabs--card > .el-tabs__header .el-tabs__item {
  height: 32px;
  line-height: 32px;
  padding: 0 24px;
  border: 1px solid #e8e8f5;
  margin-left: -1px;
  transition: none;
}
.target_box /deep/ .el-tabs--card > .el-tabs__header {
  border: 0;
}
.target_box /deep/ .el-tabs--card > .el-tabs__header .el-tabs__item.is-active {
  background: #4c7ae3;
  color: #fff;
  border: 1px solid #4c7ae3;
}
.target_box /deep/ .el-tabs--card > .el-tabs__header .el-tabs__nav {
  border-left: 1px solid #e8e8f5;
  border-top: none;
  border-right: none;
}
.linkafter {
  display: inline-block;
  border-right: 1px solid #e8e8f5;
  padding-right: 10px !important;
  height: 14px;
  line-height: 16px;
  padding-left: 10px;
}
.linkbefore {
  display: inline-block;
  border-left: 1px solid #e8e8f5;
  padding-left: 10px !important;
  height: 14px;
  line-height: 16px;
  padding-left: 10px;
}
.dialog_b_btn {
  position: absolute;
  top: 15px;
  right: 24px;
  font-size: 14px;
  button {
    color: #4c7ae3;
  }
}
.buginfo_box {
  padding: 24px;
}
.bugbasicinfo {
  padding: 24px;
  background: #fff;
  border: 1px solid rgba(232, 232, 245, 1);
}
.bugotherinfo {
  margin-top: 32px;
  .part_title {
    margin-bottom: 8px;
  }
  .content {
    background: rgba(255, 255, 255, 1);
    border-radius: 2px;
    border: 1px solid rgba(232, 232, 245, 1);
    padding: 12px 16px;
    color: rgba(72, 72, 102, 0.64);
    font-size: 13px;
  }
}
.title_bg {
  width: 84px;
  height: 32px;
  font-size: 13px;
  font-weight: 500;
}
.title_bg1 {
  background-color: rgba(243, 95, 40, 0.12) !important;
  border: 1px solid rgba(24, 144, 255, 0.08);
  color: #f35f28 !important;
  border-left: 3px solid #f35f28;
}
.title_bg2 {
  background-color: rgba(76, 122, 227, 0.12) !important;
  border: 1px solid rgba(24, 144, 255, 0.08);
  color: #4c7ae3 !important;
  border-left: 3px solid #4c7ae3;
}
.message > div {
  margin-bottom: 24px;
  background: #f7f7fb;
  border-radius: 4px;
  border: 1px solid #e8e8f5;
  padding: 16px;
  box-sizing: border-box;
}
.requestpack > div {
  background: #fff !important;
  padding: 0 !important;
}
.packbtn {
  border-top: 1px solid #e8e8f5;
  height: 65px;
  box-sizing: border-box;
  padding: 16px;
}

.packinput {
  padding: 16px;
  height: 185px;
  box-sizing: border-box;
}
.message .title_bg {
  margin-bottom: 8px;
}
.message > label {
  display: inline-block;
  width: 80px;
  text-align: center;
  height: 26px;
  line-height: 26px;
  color: #fff;
  background-color: #4c7ae3;
  font-weight: bold;
  font-size: 12px;
}
.message > div {
  height: 253px;
  overflow-y: auto;
}
.delButton_popper {
  padding: 16px !important;
  .el-button--mini {
    padding: 5px 10px;
    border-radius: 2px;
  }
}
.delText {
  margin-bottom: 16px;
  color: rgba(72, 72, 102, 0.64);
  i {
    color: #f9b640;
    margin-right: 10px;
  }
}
.filetree {
  padding: 16px;
  float: left;
  width: 310px;
  border: 1px solid #e8e8f5;
  box-sizing: border-box;
  border-radius: 4px;
  // display: flex;
  // align-items: stretch;
  /deep/ .el-button--small {
    padding: 9px;
  }
  /deep/ .el-button--small.is-round {
    padding: 9px;
  }
  .el-button + .el-button {
    margin-left: 8px;
  }
}
.filelist {
  float: left;
  width: calc(100% - 334px);
  margin-left: 24px;
}
.equipmentinfo {
  margin-bottom: 24px;
}
.targetNofind {
  display: inline-block;
  background: rgba(255, 255, 255, 0.12);
  border-radius: 2px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
  height: 32px;
  line-height: 32px;
  padding: 0 17px;
  i {
    margin-right: 16px;
  }
}
/*优*/
.networkcolor0 {
  color: #67c23a;
}
/*良*/
.networkcolor1 {
  color: #5cb85c;
}
/*中*/
.networkcolor2 {
  color: #e6a23c;
}
/*差*/
.networkcolor3 {
  color: #f56c6c;
}
/*断网*/
.networkcolor4 {
  color: #909399;
}
</style>
<style >
.finished g.node rect {
  fill: #f6f7f7;
  stroke: #d8d8d8;
}
.finished .node text tspan {
  fill: rgba(72, 72, 102, 0.6399999856948853);
}
g.node rect {
  /* padding: 5px; */
  fill: #dbe4f9;
  stroke: none;
  stroke-width: 0.5px;
}
g.node text {
  color: #4c7ae3;
  font: 14px;
  /*font-weight:700;*/
  font-size: 12px;
  cursor: pointer;
}
.node text {
  font-weight: bold;
  font-size: 12px;
  pointer-events: none;
  text-anchor: middle;
  fill: white;
}

.node rect {
  fill: white;
  stroke-width: 0px;
  color: white;
}
g > rect {
  stroke: #dee6e8;
  stroke-width: 1px;
}
.label g {
  transform: translate(0, -9px);
}
/* 当 label标签为html */
g div {
  width: 100px;
  text-align: center;
  font-size: 14px;
  vertical-align: super;
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
g text {
  fill: #687386 !important;
  font-size: 14px !important;
}
.edgePath path {
  stroke: #c4c4ce;
  stroke-width: 1px;
}
/* common   high   middle low  info    */
/* g.common>rect{
        fill:rgba(76, 122, 227, 0.12);
        stroke: none;
        stroke-width: 0.5px;
    }
    .common text tspan{
        fill:rgba(72, 72, 102, 0.64); 
    }  */
/* 致命 */
g.deadly > rect {
  fill: #f87d7d !important;
  stroke: #f87d7d !important;
  stroke-width: 0.5px;
}
.deadly text tspan {
  fill: #fff !important;
  font-size: 14px;
}
/* //高危 */
g.high rect {
  fill: #f9b640 !important;
  stroke: #f9b640 !important;
  stroke-width: 0.5px;
}
.high text tspan {
  fill: #fff !important;
  font-size: 14px;
}
/* //中危 */
g.middle > rect {
  fill: #4c7ae3 !important;
  stroke: #4c7ae3 !important;
  stroke-width: 0.5px;
}
.middle text tspan {
  fill: #fff !important;
  font-size: 14px;
}
/* //低危 */
g.low > rect {
  fill: #65c680 !important;
  stroke: #65c680 !important;
  stroke-width: 0.5px;
}
.low text tspan {
  fill: #fff !important;
  font-size: 14px;
}
/* //信息 */

g.info > rect {
  fill: #f6f7f7;
  stroke: #d8d8d8;
}
.info text tspan {
  fill: rgba(72, 72, 102, 0.6399999856948853);
  font-size: 14px;
}

/* 正在运行状态节点 */
.progress .label g {
  /* transform: translate(-50px,-9px); */
}
g.progress > rect {
  fill: #f6f7f7 !important;
  stroke: #4c7ae3 !important;
  stroke-width: 1px;
}
.progress text tspan {
  fill: rgba(72, 72, 102, 0.6399999856948853) !important;
  font-size: 14px;
}
.progress i {
  display: inline-block;
  width: 20px;
  height: 20px;
  color: #4c7ae3;
}

.activeFnode text tspan {
  fill: #fff !important;
  font-size: 14px;
}
</style>
