<template>
  <!-- 新建场景 -->
  <div>
    <div class="main-title">
      <router-link
        :underline="false"
        class="classA"
        :to="{ path: '/taskscenario' }"
        >任务场景
      </router-link>
      <label class="currentpagetitle">
        <el-tooltip
          class="item"
          effect="dark"
          :content="title"
          placement="bottom"
        >
          <span>{{ title }}</span>
        </el-tooltip>
      </label>
    </div>
    <BannerBox
      tips="任务场景用于设置新建渗透任务的默认参数。"
      style="margin-bottom: 8px"
    >
      <el-button
        type="primary"
        size="small"
        @click="btnUpdate"
        v-if="isAdd == 0 && !isUpdate"
        >编辑场景</el-button
      >
      <el-button
        type="primary"
        size="small"
        @click="btnSave"
        v-if="isUpdate && isAdd != 0"
        >保存场景</el-button
      >
      <el-button
        type="primary"
        size="small"
        @click="btnSave"
        v-if="isAdd == 0 && isUpdate"
        >保存场景</el-button
      >
    </BannerBox>
    <div class="tastBtnCont">
      <div class="step_box">
        <div :class="class1" @click="btnStep(1)">
          <div>
            <i></i>
            <label for="" class="s_label">信息收集</label>
          </div>
        </div>
        <div :class="class2" @click="btnStep(2)">
          <div>
            <i></i>
            <label for="" class="s_label">工具配置</label>
          </div>
        </div>
        <div :class="class3" @click="btnStep(3)">
          <div>
            <i></i>
            <label for="" class="s_label">漏洞测试</label>
          </div>
        </div>
        <div :class="class4" @click="btnStep(4)">
          <div>
            <i></i>
            <label for="" class="s_label">通用参数</label>
          </div>
        </div>
      </div>
    </div>
    <div class="step_content_box">
      <el-form
        :model="sceneform"
        ref="form"
        :rules="rules"
        style="height: 100%"
      >
        <div v-show="class1 == 'active'" class="step">
          <el-form-item>
            <label class="dialog_item_label" style="display: block"
              >端口扫描</label
            >
            <portModule ref="port"></portModule>
          </el-form-item>
          <el-form-item>
            <label for="" class="dialog_item_label">动态爬虫</label>
            <el-switch
              v-model="isReptile"
              class="elSwitch"
              :disabled="!isUpdate"
            >
            </el-switch>
            <div v-show="isReptile">
              <reptileModule
                ref="web_crawler"
                :crawler="crawler"
              ></reptileModule>
            </div>
          </el-form-item>
          <!-- <el-form-item>
            <label for="" class="dialog_item_label">Web路径爆破</label>
            <el-switch
              v-model="isPath"
              class="elSwitch"
              :disabled="!isUpdate"
            ></el-switch>
            <div v-show="isPath">
              <pathModule
                ref="web_path"
                :webPathScan="webPathScan"
              ></pathModule>
            </div>
          </el-form-item> -->
          <el-form-item>
            <label for="" class="dialog_item_label">子域名</label>
            <el-switch
              v-model="isSubdomain"
              class="elSwitch"
              :disabled="!isUpdate"
            >
            </el-switch>
            <div v-show="isSubdomain">
              <label class="dialog_item_label_m">子域名字典</label>
              <el-select
                :disabled="!isUpdate"
                v-model="subdomain.subdomainDict"
                size="small"
                placeholder="请选择"
                style="width: 720px"
              >
                <el-option
                  v-for="(item, index) in subdomainlist"
                  :key="index"
                  :label="item.label"
                  :value="item.value"
                ></el-option>
              </el-select>
            </div>
          </el-form-item>
        </div>
        <div v-show="class2 == 'active'" class="step">
          <el-form-item>
            <label for="" class="dialog_item_label">漏洞利用</label>
            <el-switch v-model="isUse" class="elSwitch" :disabled="!isUpdate">
            </el-switch>
          </el-form-item>
          <el-form-item>
            <label for="" class="dialog_item_label">安全测试</label>
            <el-switch
              v-model="isSafetyTesting"
              class="elSwitch"
              :disabled="!isUpdate"
            >
            </el-switch>
          </el-form-item>
          <el-form-item>
            <label for="" class="dialog_item_label">口令爆破</label>
            <el-switch
              v-model="isGuess"
              class="elSwitch"
              :disabled="!isUpdate"
            ></el-switch>
            <div v-show="isGuess">
              <guessModule ref="web_guess"></guessModule>
            </div>
          </el-form-item>
          <el-form-item>
            <label for="" class="dialog_item_label">横向移动</label>
            <el-switch
              v-model="isPostPenetration"
              class="elSwitch"
              :disabled="!isUpdate"
            >
            </el-switch>
            <div v-show="isPostPenetration">
              <postPenetrationModule
                ref="postPenetration"
              ></postPenetrationModule>
            </div>
          </el-form-item>
        </div>
        <div v-show="class3 == 'active'" class="step step2">
          <el-form-item>
            <div class="style_sw">
              <el-select
                size="small"
                v-model="form_search.type"
                placeholder="请选择"
                :disabled="!isUpdate"
                @change="showSelected"
              >
                <el-option :disabled='form_search.type == 3' value="1" label="全选">全选</el-option>
                <el-option value="2" label="取消全选">取消全选</el-option>
                <el-option value="3" label="显示已选">显示已选</el-option>
                <el-option value="4" label="显示全部">显示全部</el-option>
              </el-select>
            </div>
            <div class="serach-condition" v-show="!isShowHighSerch">
              <div class="search-text">
                <el-input
                  placeholder="请输入查找内容"
                  @keydown.enter.native="handlesearch"
                  v-model="form_search.search"
                  class="input-with-select"
                  size="small"
                  clearable
                >
                </el-input>
                <el-button type="primary" size="small" @click="handlesearch"
                  >搜索</el-button
                >
              </div>
              <div>
                <el-button
                  type="primary"
                  size="small"
                  @click="isShowHighSerch = true"
                  >高级搜索</el-button
                >
              </div>
              <div>
                <el-button type="primary" size="small" @click="handleReset"
                  >重置</el-button
                >
              </div>
            </div>
            <div
              v-show="isShowHighSerch"
              style="height: 300px;"
            >
              <div class="serach-condition">
                <div class="search-text">
                  <el-button type="primary" size="small" @click="handlesearch"
                    >搜索</el-button
                  >
                </div>
                <div>
                  <el-button
                    type="primary"
                    size="small"
                    @click="isShowHighSerch = false"
                    >隐藏高级搜索</el-button
                  >
                </div>
                <div>
                  <el-button type="primary" size="small" @click="handleResetByW"
                    >重置</el-button
                  >
                </div>
              </div>
              <div style="margin-top: 40px">
                <!-- 第一行全选 -->
                <div style="display:flex">
                    <div style="width:100px">漏洞对象 </div>
                  <el-checkbox
                  style="margin:0 10px"
                    :indeterminate="highSelectOne"
                    v-model="highCheckAll1"
                    @change="handleCheckAllChange"
                    >全选</el-checkbox
                  >
                  <el-checkbox-group
                    v-model="checks1"
                    @change="handleCheckedCitiesChange"
                  >
                    <el-checkbox
                      v-for="city in vulnenum.class"
                      v-show="city.label!='无'"
                      :label="city.value"
                      :key="city.value"
                      >{{ city.label }}</el-checkbox
                    >
                  </el-checkbox-group>
                </div>
                <!-- 第2行全选 -->
                <div style="display:flex">
                    <div style="width:100px">操作系统 </div>
                  <el-checkbox
                  style="margin:0 10px"
                     :indeterminate="highSelectOne2"
                    v-model="highCheckAll2"
                    @change="handleCheckAllChange2"
                    >全选</el-checkbox
                  >
                  <el-checkbox-group
                    v-model="checks2"
                    @change="handleCheckedCitiesChange2"
                  >
                    <el-checkbox
                    v-show="city.label!='无'"
                      v-for="city in vulnenum.operateSystem"
                      :label="city.value"
                      :key="city.value"
                      >{{ city.label }}</el-checkbox
                    >
                  </el-checkbox-group>
                </div>
                <!-- 第3行全选 -->
                <div style="display:flex">
                    <div style="width:100px">漏洞风险 </div>
                  <el-checkbox
                  style="margin:0 10px"
                    :indeterminate="highSelectOne3"
                    v-model="highCheckAll3"
                    @change="handleCheckAllChange3"
                    >全选</el-checkbox
                  >
                  <el-checkbox-group
                    v-model="checks3"
                    @change="handleCheckedCitiesChange3"
                  >
                    <el-checkbox
                    v-show="city.label!='无'"
                      v-for="city in vulnenum.risk"
                      :label="city.value"
                      :key="city.value"
                      >{{ city.label }}</el-checkbox
                    >
                  </el-checkbox-group>
                </div>
                <!-- 第4行全选 -->
                <div style="display:flex">
                    <div style="width:100px">利用影响 </div>
                  <el-checkbox
                  style="margin:0 10px"
                    :indeterminate="highSelectOne4"
                    v-model="highCheckAll4"
                    @change="handleCheckAllChange4"
                    >全选</el-checkbox
                  >
                  <el-checkbox-group
                    v-model="checks4"
                    @change="handleCheckedCitiesChange4"
                  >
                    <el-checkbox
                      v-for="city in vulnenum.exploitImpact"
                      :label="city.value"
                      :key="city.value"
                      >{{ city.label }}</el-checkbox
                    >
                  </el-checkbox-group>
                </div>
                <!-- 第5行全选 -->
                <div style="display:flex">
                    <div style="width:100px">检测方式 </div>
                  <el-checkbox
                  style="margin:0 10px"
                    :indeterminate="highSelectOne5"
                    v-model="highCheckAll5"
                    @change="handleCheckAllChange5"
                    >全选</el-checkbox
                  >
                  <el-checkbox-group
                    v-model="checks5"
                    @change="handleCheckedCitiesChange5"
                  >
                    <el-checkbox
                      v-for="city in vulnenum.vulScriptVerifyType"
                      :label="city.value"
                      :key="city.value"
                      >{{ city.label }}</el-checkbox
                    >
                  </el-checkbox-group>
                </div>
              </div>
            </div>
          </el-form-item>
          <div>
            <el-table
              ref="multipleTable"
              :data="tableData"
              tooltip-effect="dark"
              style="width: 100%"
              @select="handleSelect"
              @select-all="handleSelectAll"
            >
              <el-table-column
                type="selection"
                width="55"
                :selectable="checkboxT"
              >
              </el-table-column>
              <el-table-column prop="name" label="漏洞名称"> </el-table-column>
              <el-table-column prop="risk" label="漏洞风险" width="120">
                <template slot-scope="scope">
                  <span
                    :class="[
                      { 'riskstyle risk_info2': scope.row.risk == 0 },
                      { 'riskstyle risk_hight': scope.row.risk == 1 },
                      { 'riskstyle risk_middle': scope.row.risk == 2 },
                      { 'riskstyle risk_low': scope.row.risk == 3 },
                      { 'riskstyle risk_nofind': scope.row.risk == 4 },
                      { 'riskstyle risk_info': scope.row.risk == 5 },
                    ]"
                    ><i></i>{{ scope.row.riskName }}</span
                  >
                </template>
              </el-table-column>
              <el-table-column prop="scriptType" label="漏洞类型">
                <template slot-scope="scope">
                  <span>{{ scope.row.scriptType === 'mitm' || scope.row.scriptType === 'universal' ? '通用漏洞' : '专用漏洞' }}</span>
                </template>
              </el-table-column>
              <el-table-column
                prop="description"
                label="漏洞描述"
                show-overflow-tooltip
              >
              </el-table-column>
            </el-table>
            <el-pagination
              background
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
              :current-page="currentPage"
              :page-size="pageSize"
              layout=" total,  prev, pager, next, sizes,jumper"
              :total="total"
            >
            </el-pagination>
          </div>
        </div>
        <div v-show="class4 == 'active'" class="step currency">
          <el-form-item label="" prop="template">
            <label class="dialog_item_label">场景名称 </label>
            <el-input
              v-model="sceneform.name"
              style="width: 320px"
              :disabled="!isUpdate"
              placeholder="请输入场景名称"
            ></el-input>
          </el-form-item>
          <el-form-item>
            <label class="dialog_item_label">场景描述</label>
            <el-input
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 5 }"
              placeholder="请输入场景描述"
              v-model="sceneform.describe"
              :disabled="!isUpdate"
              style="
                width: 720px;
                margin-bottom: 10px;
                vertical-align: text-top;
              "
              resize="none"
            >
            </el-input>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>
</template>
<style lang="less" scoped>
@import "./css/createscene.less";
</style>

<style lang="less">
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
  fill: #f87d7d;
  stroke: #f87d7d;
  stroke-width: 0.5px;
}

.deadly text tspan {
  fill: #fff;
  font-size: 14px;
}

/* //高危 */
g.high rect {
  fill: #f9b640;
  stroke: #f9b640;
  stroke-width: 0.5px;
}

.high text tspan {
  fill: #fff;
  font-size: 14px;
}

/* //中危 */
g.middle > rect {
  fill: #4c7ae3;
  stroke: #4c7ae3;
  stroke-width: 0.5px;
}

.middle text tspan {
  fill: #fff;
  font-size: 14px;
}

/* //低危 */
g.low > rect {
  fill: #65c680;
  stroke: #65c680;
  stroke-width: 0.5px;
}

.low text tspan {
  fill: #fff;
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
  fill: #f6f7f7;
  stroke: #4c7ae3;
  stroke-width: 1px;
}

.progress text tspan {
  fill: rgba(72, 72, 102, 0.6399999856948853);
  font-size: 14px;
}

.progress i {
  display: inline-block;
  width: 20px;
  height: 20px;
  color: #4c7ae3;
}

.activeFnode text tspan {
  fill: #fff;
  font-size: 14px;
}
</style> 
<script>   
import portModule from "./components/port_module.vue"
import reptileModule from "./components/reptile_module.vue"
import pathModule from "./components/path_module.vue"
import guessModule from "./components/guess_module.vue"
import postPenetrationModule from "./components/post_penetration_module.vue"
import BannerBox from "@/components/BannerBox.vue"
import scene from '@/api/scene.js'
import { set } from 'vue'
export default {
  name: 'createscene',
  components: {
    portModule,
    reptileModule,
    pathModule,
    guessModule,
    postPenetrationModule,
    BannerBox
  },
  data () {
    return {
        highCheckAll1: false,
        checks1: [],
        cities1: [],
        highSelectOne: false,

        highCheckAll2: false,
        checks2: [],
        cities2: [],
        highSelectOne2: false,

        highCheckAll3: false,
        checks3: [],
        cities3: ['上海', '北京', '广州', '深圳'],
        highSelectOne3: false,

        highCheckAll4: false,
        checks4: [],
        cities4: ['上海', '北京', '广州', '深圳'],
        highSelectOne4: false,

        highCheckAll5: false,
        checks5: [],
        cities5: ['上海', '北京', '广州', '深圳'],
        highSelectOne5: false,
      isShowHighSerch: false,
      title: '新增场景',
      isAdd: this.$route.query.isAdd,//1,创建，0,编辑,2:复制
      sceneid: this.$route.query.sceneid,
      scene_name: this.$route.query.scene_name,
      class1: 'active',
      class2: '',
      class3: '',
      class4: '',
      isShowchart: false,
      state: [
      ],
      edg: [
      ],
      diagGraph: {},
      menuflag: false,
      tabLabellist: ['信息收集', '漏洞测试', '工具关联', '渗透参数', '通用参数'],
      sceneform: {
        // scanportdefaulttype: ['type1', 'type2'],
        scanporttype: ['type1', 'type2'],
        other_tools: [],
        vul_use: false,
        // remote_control:false,
        // post_penetration:false,
        info: [],
        scan: [],
        test: [],
        use: [],
        name: '',
        describe: '',
        hx: false,
        kl: false,

      },
      rules: {
        name: { required: true, message: '场景名称不能为空', trigger: 'blur' },
      },
      form_search: {
        type: '',
        searchinfo: '',
        search: '',
        object: [],
        system: [],
        risk: [],
        stage: [],
        useinfluence: [],
        mode: [],
        isSelected: false, //是否已选择漏洞
        page: 1,
      },
      searchinfoShow: false,
      searchtestShow: false,
      searchSacnShow: false,
      searchUseShow: false,
      advanced_aearch_txt: '高级搜索',
      isShowAdvanced: false,
      isIndeterminateobject: false,
      checkAllobject: false,
      objectlist: [
        {
          label: '主机服务',
          value: 1,
        },
        {
          label: 'WEB应用',
          value: 2,
        },
        {
          label: '数据库',
          value: 3,
        },
        {
          label: '虚拟化',
          value: 4,
        },
        {
          label: '网络设备',
          value: 5,
        },
        {
          label: '安全设备',
          value: 6,
        },
        {
          label: '办公设备',
          value: 7,
        },
        {
          label: '视频监控',
          value: 8,
        },
        {
          label: '大数据',
          value: 9,
        },
        {
          label: 'IOT',
          value: 10,
        },
        {
          label: '其他',
          value: 11,
        },
      ],
      checkobjectlist: [],
      isIndeterminatesystem: false,
      checkAllsystem: false,
      systemlist: [
        {
          label: 'Windows',
          value: 1
        },
        {
          label: 'Linux',
          value: 2
        },
        {
          label: 'Unix',
          value: 3
        },
        {
          label: '其他',
          value: 4
        },
      ],
      checksystemlist: [],
      isIndeterminaterisk: false,
      checkAllrisk: false,
      risklist: [
        {
          label: '致命',
          value: 1,
        },
        {
          label: '高危',
          value: 2,
        },
        {
          label: '中危',
          value: 3,
        },
        {
          label: '低危',
          value: 4,
        },
      ],
      checkrisklist: [],
      isIndeterminatestage: false,
      checkAllstage: false,
      stagelist: [
        {
          label: '信息收集',
          value: 1,
        },
        {
          label: '扫描探测',
          value: 2,
        },
        {
          label: '漏洞检测',
          value: 3,
        },
        {
          label: '漏洞利用',
          value: 4,
        },
      ],
      checkstagelist: [],
      isIndeterminateuseinfluence: false,
      checkAlluseinfluence: false,
      useinfluencelist: [
        {
          label: '系统宕机',
          value: 1,
        },
        {
          label: '服务崩溃',
          value: 2,
        },
        {
          label: '拒绝服务',
          value: 3,
        },
        {
          label: '运行缓慢',
          value: 4,
        },
        {
          label: '无影响',
          value: 5,
        },
      ],
      checkuseinfluencelist: [],
      isIndeterminatemode: false,
      checkAllmode: false,
      modelist: [
        {
          label: 'POC',
          value: 1,
        },
        {
          label: 'EXP',
          value: 2,
        },
        {
          label: 'POC/EXP',
          value: 3,
        },
        {
          label: 'Version',
          value: 4,
        },
      ],
      checkmodelist: [],
      tableData: [],
      currentPage: 1,
      total: 0,
      infolist: [],
      scanlist: [],
      testlist: [],
      uselist: [],
      testlistback: [],
      uselistback: [],
      formuse: {
        searchtest: '',
        searchuse: ''
      },
      dialogaddFormVisible: false,
      dialogType: 0,
      isUpdate: false,
      isPortScan: false, // 是否选中端口扫描。true选中，false未选中
      isReptile: false, //是否爬虫
      isPath: false, //路径爆破
      isUse: false,//渗透利用
      isGuess: false,//口令 
      isPostPenetration: false,//横向移动
      isSafetyTesting: false,//安全测试
      isSubdomain: false, //子域名 
      subdomainlist: [],
      portRang: [], //端口范围 
      portRangeValue: [],
      tcp_scan_type: [],
      firstIN: 1,
      tableIds: [],
      vulnlist: [],
      vulnlistold: [],
      pocnum: 0,
      expnum: 0,
      web_crawler: {},
      checkvulnlist: [],
      user_id: '',
      infocheckedvuln: [],//场景详情，已选中的漏洞
      subdomain: { //子域名
        isOpen: true,
        subdomainDict: '',
      },
      STEP: 4,
      config: {},
      pageSize: 10,
      crawler: {},
      webPathScan: {},
      searchIDs: [], //所有符合条件的id
      vulnenum: {},
      web_guess: {},
      postPenetration: {},
    }
  },
  created () {
    this.$store.state.activefirstMenu = '/taskscenario'
    this.pageSize = this.$commonjs.pageSize
    this.user_id = this.$commonjs.decryptCBC(localStorage.getItem('user_id'), this.$commonjs.myKey)
  },
  mounted () {
    this.vulnEnum()

    this.getSceneEnum()
    let index = 1
    if (this.isAdd == 1) { //新增

      this.isUpdate = true
      this.title = '新增场景'
      this.$nextTick(() => {
        this.$refs.port.getIsUpdate(true)
        this.$refs.web_crawler.getIsUpdate(this.isUpdate)
        // this.$refs.web_path.getIsUpdate(this.isUpdate)
        this.$refs.web_guess.getIsUpdate(this.isUpdate)
        this.$refs.postPenetration.getIsUpdate(this.isUpdate)

      })

    } else {
      this.title = '编辑场景：' + this.scene_name
      this.isUpdate = false;
      this.form_search.type = '3';
    }
    for (var i = 0; i < this.STEP; i++) {
      let _i = i + index
      if (_i == index) {
        this['class' + index] = 'active'
      } else {
        this['class' + _i] = ''
      }
    }


  },
  methods: {
    handleCheckAllChange(val) {
        let arr = []
        for (const iterator of this.vulnenum.class) {
           arr.push(iterator.value) 
        } 
    
        this.checks1 = val ? arr : [];
        this.highSelectOne = false;
      },
      handleCheckedCitiesChange(value) {
        let checkedCount = value.length;
        this.highCheckAll1 = checkedCount === this.vulnenum.class.length;
        this.highSelectOne = checkedCount > 0 && checkedCount < this.vulnenum.class.length;
      },
    handleCheckAllChange2(val) {
        let arr = []
        for (const iterator of this.vulnenum.operateSystem) {
           arr.push(iterator.value) 
        } 
    
        this.checks2 = val ? arr : [];
        this.highSelectOne2 = false;
      },
      handleCheckedCitiesChange2(value) {
        console.log(value, 'value', this.cities2, 'this.cities2');
        let checkedCount = value.length;
        this.highCheckAll2 = checkedCount === this.vulnenum.operateSystem.length;
        this.highSelectOne2 = checkedCount > 0 && checkedCount < this.vulnenum.operateSystem.length;
      },
    handleCheckAllChange3(val) {
        let arr = []
        for (const iterator of this.vulnenum.risk) {
           arr.push(iterator.value) 
        } 
    
        this.checks3 = val ? arr : [];
        this.highSelectOne3 = false;
      },
      handleCheckedCitiesChange3(value) {
        let checkedCount = value.length;
        this.highCheckAll3 = checkedCount === this.vulnenum.risk.length;
        this.highSelectOne3 = checkedCount > 0 && checkedCount < this.vulnenum.risk.length;
      },
    handleCheckAllChange4(val) {
        let arr = []
        for (const iterator of this.vulnenum.exploitImpact) {
           arr.push(iterator.value) 
        } 
    
        this.checks4 = val ? arr : [];
        this.highSelectOne4 = false;
      },
      handleCheckedCitiesChange4(value) {
        let checkedCount = value.length;
        this.highCheckAll4 = checkedCount === this.vulnenum.exploitImpact.length;
        this.highSelectOne4 = checkedCount > 0 && checkedCount < this.vulnenum.exploitImpact.length;
      },
    handleCheckAllChange5(val) {
        let arr = []
        for (const iterator of this.vulnenum.vulScriptVerifyType) {
           arr.push(iterator.value) 
        } 
    
        this.checks5 = val ? arr : [];
        this.highSelectOne5 = false;
      },
      handleCheckedCitiesChange5(value) {
        let checkedCount = value.length;
        this.highCheckAll5 = checkedCount === this.vulnenum.vulScriptVerifyType.length;
        this.highSelectOne5 = checkedCount > 0 && checkedCount < this.vulnenum.vulScriptVerifyType.length;
      },








    async vulnEnum () {
      scene.vulnEnum().then(res => {
        if (res.code == 200) {
          this.vulnenum = res.data
          this.vulnenum.class.map((item,index)=>{
              //返回的class数组中的某项value为0，则删除该项
              if(item.value==0){
                this.vulnenum.class.splice(index,1)
            }
          })
          this.vulnenum.operateSystem.map((item,index)=>{
              //返回的class数组中的某项value为0，则删除该项
              if(item.value==0){
                this.vulnenum.operateSystem.splice(index,1)
            }
          })
          this.vulnenum.risk.map((item,index)=>{
              //返回的class数组中的某项value为0，则删除该项
              if(item.value==0){
                this.vulnenum.risk.splice(index,1)
            }
          })
          this.vulnenum.type.map((item,index)=>{
              //返回的class数组中的某项value为0，则删除该项
              if(item.value==0){
                this.vulnenum.type.splice(index,1)
            }
          })
     
          
        }
      })
    },
    async getSceneEnum () {
      let res = await scene.getSceneEnum()
      if (res.code == 200) {
        this.subdomainlist = res.data.subdomainDictCollect //子域名列表  
        if (this.isAdd == 1) {
          this.subdomainlist.forEach(item => {
            if (item.isDefault == true) {
              this.subdomain.subdomainDict = item.value
            }

          })
        }
        this.crawler = res.data.crawler
        this.webPathScan = res.data.webPathScan
        this.$refs.port.getEnum(res.data.portScan, this.isAdd)
        this.$refs.web_crawler.getEnum(res.data.crawler, this.isAdd)
        // this.$refs.web_path.getEnum(this.webPathScan, this.isAdd) 
        this.$refs.web_guess.getEnum(res.data.weakPass, this.isAdd)
        this.$refs.postPenetration.getEnum(res.data.lateralMove, this.isAdd)

        if (this.isAdd != 1) { //编辑
          this.getInfo()
        }
      }
    },
    async getInfo () { //详情
      const res = await scene.getSceneinfo({
        taskTemplateId: this.sceneid
      })
      if (res.code == 200) {
        this.sceneform.name = res.data.templateName
        this.sceneform.describe = res.data.describe

        let _config = res.data.config

        this.config = _config

        if (_config.webPathScanConfig.isOpen) { //路径爆破
          this.isPath = true
        }
        if (_config.subdomainCollectConfig.isOpen) {
          this.isSubdomain = true
          this.subdomain.subdomainDict = _config.subdomainCollectConfig.subdomainDict
        }
        //端口扫描  
        this.$refs.port.getConifg(_config.portScanConfig)

        //爬虫 
        this.isReptile = _config.webCrawlerConfig.isOpen
        this.$refs.web_crawler.getConifg(_config.webCrawlerConfig)

        //路径猜测 
        // this.$refs.web_path.getConifg(_config.webPathScanConfig)

        //口令
        this.$refs.web_guess.getConifg(_config.weakPassConfig)
        this.isGuess = _config.weakPassConfig.isOpen

        //横向移动
        this.$refs.postPenetration.getConifg(_config.lateralMove)
        this.isPostPenetration = _config.lateralMove.isOpen
        //安全测试
        this.isSafetyTesting = _config.safeTest
        //漏洞利用
        this.isUse = _config.vulExploit

        this.tableIds = _config.vulIdsConfig
      }
    },
    btnUpdate () { //编辑
      this.isUpdate = true

      this.$nextTick(() => {
        this.$refs.port.getIsUpdate(this.isUpdate)
        this.$refs.web_crawler.getIsUpdate(this.isUpdate)
        // this.$refs.web_path.getIsUpdate(this.isUpdate)
        this.$refs.web_guess.getIsUpdate(this.isUpdate)
        this.$refs.postPenetration.getIsUpdate(this.isUpdate)
      })
    },
    checkboxT (row, index) {
      if (!this.isUpdate) {
        return 0
      } else {
        return 1
      }
    },

    btnStep (index) { //编辑点击第几步 
      for (var i = 0; i < this.STEP; i++) {
        let _i = i + 1
        if (_i == index) {
          this['class' + index] = 'active'
        } else {
          this['class' + _i] = 'finish'
        }
      }
      if (this.class1 == 'active') {
        this.$nextTick(() => {
          this.$refs.port.getIsUpdate(this.isUpdate)
          this.$refs.web_crawler.getIsUpdate(this.isUpdate)
          // this.$refs.web_path.getIsUpdate(this.isUpdate)
        })

      }
      if (this.class2 == 'active') {
        this.$nextTick(() => {
          this.$refs.postPenetration.getIsUpdate(this.isUpdate)

        })

      }
      if (this.class3 == 'active') {
        this.getVulnlist()
      }
    },
    toggleSelection (rows) {   // 默认选中项
      if (rows) {
        rows.forEach(row => {
          this.$refs.multipleTable.toggleRowSelection(row, true)
        })
      } else {
        this.$refs.multipleTable.clearSelection()
      }
    },
    async btnSave () { //保存 创建/编辑  
      if (this.sceneform.name == '') {
        this.$message({
          message: '场景名称不能为空',
          type: "error",
        })
        return
      }

      let param = this.fnsubmitparm()
      if (this.isAdd == 0) {
        param.taskTemplateId = this.sceneid
      }
      const res = await scene.saveScene(param)

      if (res.code == 200) {
        this.$message({
          message: "保存场景成功",
          type: "success",
        })
        this.$router.push({
          path: `/taskscenario`,
        })
      } else {
        this.$message({
          message: res.msg,
          type: "error",
        })
      }
    },
    fnsubmitparm () { //保存场景参数

      let _json = {}
      _json.userId = Number(this.user_id)
      _json.templateName = this.sceneform.name
      _json.describe = this.sceneform.describe
      let config = {
        portScanConfig: {},
        webCrawlerConfig: {},
        webPathScanConfig: {},
        subdomainCollectConfig: {},
        websiteLoginConfig: {}
      }
      config.websiteLoginConfig.isOpen = false
      config.websiteLoginConfig.list = []
      //端口扫描
      this.portconfig = this.$refs.port.getAllData()
      this.portconfig.isOpen = true
      config.portScanConfig = this.portconfig

      //动态爬虫 web_crawler_config  
      this.web_crawler = this.$refs.web_crawler.getAllData()
      this.web_crawler.isOpen = this.isReptile
      config.webCrawlerConfig = this.web_crawler

      //路径爆破 web_path_scan_config  
      // this.web_path = this.$refs.web_path.getAllData()
      // this.web_path.isOpen = this.isPath
      // config.webPathScanConfig = this.web_path
      config.webPathScanConfig = {
        isOpen: this.isPath,
        pathDict: 0,
        threadNum: 10,
        timeout: 5
      }

      //子域名收集 subdomain_collect_config

      config.subdomainCollectConfig.isOpen = this.isSubdomain
      config.subdomainCollectConfig.subdomainDict = !this.subdomain.subdomainDict ? 0 : this.subdomain.subdomainDict

      //口令爆破
      this.web_guess = this.$refs.web_guess.getAllData()
      this.web_guess.isOpen = this.isGuess
      config.weakPassConfig = this.web_guess

      // 横向移动
      this.postPenetration = this.$refs.postPenetration.getAllData()
      this.postPenetration.isOpen = this.isPostPenetration
      config.lateralMove = this.postPenetration

      //安全测试
      config.safeTest = this.isSafetyTesting
      //漏洞利用
      config.vulExploit = this.isUse

      config.vulIdsConfig = this.tableIds

      _json.config = JSON.stringify(config)

      return _json
    },
    //第三步漏洞搜索
    handlesearch () {
      this.form_search.page = 1
      this.getVulnlist()
    },
    handleReset () { //重置
      this.form_search.search = ''
      this.form_search.type = ''
      this.currentPage = 1
      this.form_search.page = 1
      this.pageSize = 10
       this.tableIds = []

      this.getVulnlist()

    },
    handleResetByW () { //重置
      this.form_search.search = ''
      this.form_search.type = ''
      this.currentPage = 1
      this.form_search.page = 1
      this.tableIds = []
      this.pageSize = 10
      this.checks1 = []
      this.checks2 = []
      this.checks3 = []
      this.checks4 = []
      this.checks5 = []

      //全选设置为false
      this.highCheckAll1 = false
      this.highCheckAll2 = false
      this.highCheckAll3 = false
      this.highCheckAll4 = false
      this.highCheckAll5 = false

      this.getVulnlist()

    },

    handleSelectAll (val) {
      if (this.firstIN === 1) { // 意思第一次点击不会执行hangleSelectAll里面的方法
        var v = this
        // remove
        if (val.length > 0) {
          for (const n in val) {
            if (this.tableIds.indexOf(val[n].id) == -1) { //存在
              this.tableIds.push(val[n].id)
            }

          }
        }
        if (val.length === 0) { //取消 
          for (var i = 0; i < v.tableData.length; i++) {
            for (var j in this.tableIds) {
              if (v.tableData[i].id === this.tableIds[j]) {
                v.vulnlist.splice(j, 1)
                this.tableIds.splice(j, 1)
                break
              }
            }
          }

        }
        if (v.vulnlist.length === 0) {
          for (const i in val) {
            v.tableIds.push(val[i].id)
          }
        } else {
          for (const i in val) {
            let flag = false
            for (const j in v.tableIds) {
              if (v.tableIds[j] === val[i].id) {
                flag = true
                break
              }
            }
            if (!flag) {
              v.tableIds.push(val[i].id)
            }
          }
        }
      }
    },
    handleSelect (val, row) {

      if (this.firstIN === 1) { // 设置第一次进来才回触发事件
        /* 1 => add ; 0 => remove*/
        let flag = 0
        for (const i in val) {
          if (row.id === val[i].id) {
            flag = 1
            break
          }
        }
        if (flag === 1) {
          // 如果判断当前为添加则将当前勾选数据push到指定数组中 
          this.tableIds.push(row.id)
        } else {
          // 否则从数组中删除当前行数据 
          for (const i in this.tableIds) {
            if (this.tableIds[i] === row.id) {
              this.tableIds.splice(i, 1)
            }
          }
        }
      }
    },
    handleSizeChange (t) {
      this.currentPage = 1
      this.form_search.page = 1
      this.pageSize = t
      this.getVulnlist()
    },
    handleCurrentChange (t) {
      this.currentPage = t
      this.form_search.page = t
      // this.pageSize = 10;
      this.getVulnlist()
    },
    async getVulnlist () { //获得漏洞信息
      let parm = {
        page: this.form_search.page,
        size: this.pageSize,
        search: this.form_search.search,
        libClasses: this.checks1.join(','), //漏洞分类
        operatingSystem: this.checks2.join(','), //漏洞分类
        libRisks: this.checks3.join(','), //漏洞分类
        exploitImpact: this.checks4.join(','), //漏洞分类
        scriptVerifyTypes: this.checks5.join(','), //漏洞分类
      }
      if (this.form_search.type == 3) { //显示已选



        parm.libIds = this.tableIds.join(',')
        console.log(parm.libIds,'parm.libIds',this.tableIds);
        if(parm.libIds) localStorage.setItem('libIDS',parm.libIds)
        
      }
      const res = await scene.getVulnlist(parm)
      if (res.code == 200) {
        this.tableData = res.data.list
        this.total = res.data.total
        this.searchIDs = res.data.libIds
        let defaultSelected = []

        this.tableData && this.tableData.forEach((item, i) => {
          if (this.tableIds.indexOf(item.id) != -1) { //存在
            defaultSelected.push(item)
          }
        })
        this.$nextTick(() => {
          this.toggleSelection(defaultSelected)
        })
      } else {
        this.$message({
          message: res.msg,
          type: "error",
        })
      }
    },
    showSelectedVunl () { //显示已选中漏洞 
      this.currentPage = 1
      this.getVulnlist()
    },
    showSelected () {
      let _type = this.form_search.type
      if (_type == 1) { //全选
        
        

           this.tableIds = this.searchIDs
          this.form_search.page = 1
          this.currentPage = 1
        this.getVulnlist()
      }
      if (_type == 2) { //取消全选
        this.tableIds = []
        this.form_search.page = 1
        this.currentPage = 1
        this.getVulnlist()
      }
      if (_type == 3) { //显示已选 
        // this.tableIds = this.config.vulIdsConfig; 
        console.log(this.tableIds,'this.tableIds----');
        if(this.tableIds.length<1&&localStorage.getItem('libIDS')){
            this.tableIds = localStorage.getItem('libIDS').split(',').map(it=>Number(it)); 
            console.log(this.tableIds,'this.tableIds');
        }
        this.form_search.page = 1
        this.currentPage = 1
        this.getVulnlist()
        if (this.tableIds.length == 0) {
          this.$message({
            message: '没有已选择的漏洞',
            type: "error",
          })
        }
      }
      if(_type == 4){ //显示全部
        this.form_search.page = 1
        this.currentPage = 1
        this.getVulnlist()
      }

    },
  },
}
</script>