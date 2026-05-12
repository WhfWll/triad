/* 整体页面 */
<template>
  <div class="neo4j">
    <!-- 面包屑 -->
    <div class="main-title">渗透图谱</div>
    <fullscreen class="neo4j2" :fullscreen.sync="fullscreen">
      <!-- 蓝色框title  Start-->
      <div class="blueChannel">
        <!-- 左侧 -->
        <div class="blueChannelLeft">
          <div class="showShan boxShaw">
            <div>
              <i class="el-icon-postcard"></i
              ><span class="spanNum">信息收集</span
              ><span>{{ chanelCount.info_collect_count }}</span>
            </div>
            <div>
              <i class="el-icon-stopwatch"></i
              ><span class="spanNum">漏洞测试</span
              ><span>{{ chanelCount.poc_count }}</span>
            </div>
            <div>
              <i class="el-icon-bangzhu"></i
              ><span class="spanNum">漏洞利用</span
              ><span>{{ chanelCount.exp_count }}</span>
            </div>
            <div>
              <i class="el-icon-full-screen"></i
              ><span class="spanNum">扫描探测</span
              ><span>{{ chanelCount.detect_count }}</span>
            </div>
            <div>
              <i class="el-icon-takeaway-box"></i
              ><span class="spanNum">规则总数</span
              ><span>{{ chanelCount.rule_total }}</span>
            </div>
            <!--  -->
          </div>
        </div>
        <!-- 中间下拉 -->
        <div class="MiddleXia" style="cursor: pointer">
          <div
            style="cursor: pointer; margin: 8px 50px 0 0"
            @click="isThree = !isThree"
          >
            <button type="button" class="btn cube cube-hover">
              <div class="bg-top">
                <div class="bg-inner"></div>
              </div>
              <div class="bg-right">
                <div class="bg-inner"></div>
              </div>
              <div class="bg">
                <div class="bg-inner"></div>
              </div>
              <div class="text">{{ isThree ? "返回" : "3D展示" }}</div>
            </button>
          </div>
          <el-dropdown v-model="MiddleSlectValue" :hide-on-click="false">
            <span class="el-dropdown-link">
              <!-- <button class="MiddleXiaButton">
                {{ MiddleSlectValue
                }}<i
                  class="el-icon-arrow-down el-icon--right"
                  style="font-size: 16px"
                ></i>
              </button> -->
              <button class="button1">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  viewBox="0 0 24 24"
                  width="24"
                  height="24"
                >
                  <path fill="none" d="M0 0h24v24H0z"></path>
                  <path
                    fill="currentColor"
                    d="M5 13c0-5.088 2.903-9.436 7-11.182C16.097 3.564 19 7.912 19 13c0 .823-.076 1.626-.22 2.403l1.94 1.832a.5.5 0 0 1 .095.603l-2.495 4.575a.5.5 0 0 1-.793.114l-2.234-2.234a1 1 0 0 0-.707-.293H9.414a1 1 0 0 0-.707.293l-2.234 2.234a.5.5 0 0 1-.793-.114l-2.495-4.575a.5.5 0 0 1 .095-.603l1.94-1.832C5.077 14.626 5 13.823 5 13zm1.476 6.696l.817-.817A3 3 0 0 1 9.414 18h5.172a3 3 0 0 1 2.121.879l.817.817.982-1.8-1.1-1.04a2 2 0 0 1-.593-1.82c.124-.664.187-1.345.187-2.036 0-3.87-1.995-7.3-5-8.96C8.995 5.7 7 9.13 7 13c0 .691.063 1.372.187 2.037a2 2 0 0 1-.593 1.82l-1.1 1.039.982 1.8zM12 13a2 2 0 1 1 0-4 2 2 0 0 1 0 4z"
                  ></path>
                </svg>
                <span>{{ MiddleSlectValue }}</span>
              </button>
            </span>
            <el-dropdown-menu slot="dropdown">
              <el-dropdown-item
                @click.native="clickDropdown(item)"
                v-for="item in ['查看100条', '查看全部']"
                :key="item"
                >{{ item }}</el-dropdown-item
              >
            </el-dropdown-menu>
          </el-dropdown>
        </div>
        <!-- 右边 -->
        <div class="blueChannel-right" style="padding: 0 20px">
          <div @click="fullscreen = !fullscreen">
            <el-tooltip
              popper-class="tooltipBlue"
              class="item"
              effect="light"
              content="全屏"
              placement="top"
            >
              <i
                style="font-size: 20px"
                class="iconfont iconquanping boxShaw"
                v-if="!fullscreen"
              ></i>
              <i
                style="font-size: 20px"
                class="iconfont icontuichuquanping boxShaw"
                v-else
              ></i>
            </el-tooltip>
          </div>

          <div class="rightIcon boxShaw">
            <div @click="finger">
              <el-tooltip
                popper-class="tooltipBlue"
                class="item"
                effect="light"
                content="新建指纹"
                placement="top"
              >
                <img style="width: 25px" src="@/assets/icon/finger.svg" />
              </el-tooltip>
            </div>
            <div @click="addInfo" v-if="false">
              <el-tooltip
                popper-class="tooltipBlue"
                class="item"
                effect="light"
                content="新建信息"
                placement="top"
              >
                <img style="width: 25px" src="@/assets/icon/cube.svg" />
              </el-tooltip>
            </div>
            <div @click="addShouJiTI">
              <el-tooltip
                popper-class="tooltipBlue"
                class="item"
                effect="light"
                content="新建信息收集实体"
                placement="top"
              >
                <img style="width: 25px" src="@/assets/icon/coll.svg" />
              </el-tooltip>
            </div>
            <div @click="addLouDOng">
              <el-tooltip
                popper-class="tooltipBlue"
                class="item"
                effect="light"
                content="新建漏洞实体"
                placement="top"
              >
                <img style="width: 25px" src="@/assets/icon/bug.svg" />
              </el-tooltip>
            </div>
          </div>
        </div>
      </div>
      <!-- 蓝色框title  End-->

      <!-- 页面主体部分 -->
      <div class="userbox context_box_bg">
        <div class="search-box">
          <div style="display: flex; height: 100%; align-items: center">
            <!-- 关系图     :Neo4jD3JsonList="trafficView"
              :Neo4jD3JsonList2="sessionView"   :tabName="radio3" -->
            <div style="position: relative; display: flex">
              <!-- <Neo4j @mouseNode="mouseNode"></Neo4j> -->
              <!-- <div class="cardX" v-show="nodeData.labels">
              <div
                style="
                  position: absolute;
                  top: 10%;
                  left: 10%;
                  width: 1190px;
                  height: 180px;
                  color: #409eff;
                  z-index: 99999;
                "
              >
                <div
                  style="
                    border-bottom: 1px dashed #01fef4;
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                  "
                >
                  <span style="color: #000">ID</span><span>{{ nodeData }}</span>
                </div>
                <div
                  style="
                    border-bottom: 1px dashed #01fef4;
                    margin-top: 20px;
                    display: flex;
                    justify-content: space-between;
                    align-items: center;
                  "
                >
                  <span style="color: #000">label</span
                  ><span v-if="nodeData.labels" style="margin-left: 20px">{{
                    nodeData.labels[0]
                  }}</span>
                </div>
              </div>
            </div> -->
            </div>
            <!-- three.js组件 -->
            <HaoYing
              @initImIt="initImIt"
              v-if="isThree"
              :threeData="threeData"
            ></HaoYing>
            <!-- combo 图 -->
            <G6
              v-if="!isThree"
              :MiddleSlectValue="MiddleSlectValue"
              :G6Data="G6Data"
              @initG6="initG6"
            ></G6>
            <!-- 新建指纹 -->
            <Drawer ref="addTa">
              <div class="DrawerTitleLeft boxShaw" slot="header-left">
                <span class="textSL">新建指纹</span>
              </div>
              <div class="DrawerTitleRight" slot="header-right">
                <button
                  class="btnFwa"
                  style="margin-right: 10px"
                  @click="saveDraw"
                >
                  <span class="shadow"></span>
                  <span class="edge"></span>
                  <span class="front text boxShaw" style="background: #67c23a">
                    保存
                  </span>
                </button>
                <!-- @click="closeDraw" -->
                <button class="btnFwa" @click="closeDraw">
                  <span class="shadow"></span>
                  <span class="edge"></span>
                  <span class="front text boxShaw"> 关闭 </span>
                </button>
              </div>
              <div slot="content" style="padding: 0 20px">
                <el-form
                  ref="createFingerForm"
                  :model="createFingerForm"
                  label-width="80px"
                >
                  <el-form-item label="产品名称">
                    <el-input v-model="createFingerForm.app_name"></el-input>
                  </el-form-item>
                  <el-form-item label="厂商  ">
                    <el-input v-model="createFingerForm.cn_name"></el-input>
                  </el-form-item>
                  <el-form-item label="类型  ">
                    <el-input v-model="createFingerForm.app_class"></el-input>
                  </el-form-item>
                  <el-form-item label="应用版本">
                    <el-input v-model="createFingerForm.app_version"></el-input>
                  </el-form-item>
                  <el-form-item label="指纹内容">
                    <el-input
                      v-model="createFingerForm.flag"
                      type="textarea"
                    ></el-input>
                  </el-form-item>
                </el-form>
              </div>
            </Drawer>
            <!-- 新建信息 -->
            <Drawer ref="addInfo">
              <div class="DrawerTitleLeft boxShaw" slot="header-left">
                <span class="textSL">新建信息</span>
              </div>
              <div class="DrawerTitleRight" slot="header-right">
                <!-- @click="saveBtn" -->
                <button
                  class="btnFwa"
                  style="margin-right: 10px"
                  @click="saveInfo"
                >
                  <span class="shadow"></span>
                  <span class="edge"></span>
                  <span class="front text boxShaw" style="background: #67c23a">
                    保存
                  </span>
                </button>
                <!-- @click="closeDraw" -->
                <button class="btnFwa" @click="closeInfo">
                  <span class="shadow"></span>
                  <span class="edge"></span>
                  <span class="front text boxShaw"> 关闭 </span>
                </button>
              </div>
              <div slot="content" style="padding: 0 20px">
                <el-form
                  ref="createFingerForm"
                  :model="createFingerForm"
                  label-width="80px"
                >
                  <el-form-item label="信息名称">
                    <el-input v-model="createFingerForm.app_name"></el-input>
                  </el-form-item>
                  <!-- <el-form-item label="厂商">
                    <el-input v-model="createFingerForm.cn_name"></el-input>
                  </el-form-item> -->
                  <el-form-item label="信息类型">
                    <el-input v-model="createFingerForm.app_class"></el-input>
                  </el-form-item>
                  <el-form-item label="输出脚本">
                    <el-input v-model="createFingerForm.app_version"></el-input>
                  </el-form-item>
                  <el-form-item label="信息描述">
                    <el-input
                      v-model="createFingerForm.flag"
                      type="textarea"
                    ></el-input>
                  </el-form-item>
                </el-form>
              </div>
            </Drawer>
            <!-- 新建漏洞实体 -->
            <Drawer :size="widthSize" ref="addLouDong">
              <div class="DrawerTitleLeft boxShaw" slot="header-left">
                <span class="textSL">新建漏洞实体</span>
              </div>
              <div class="DrawerTitleRight" slot="header-right">
                <!-- @click="saveBtn" -->
                <button
                  class="btnFwa"
                  style="margin-right: 10px"
                  @click="saveLouDongShiTi"
                >
                  <span class="shadow"></span>
                  <span class="edge"></span>
                  <span class="front text boxShaw" style="background: #67c23a">
                    保存
                  </span>
                </button>
                <!-- @click="closeDraw" -->
                <button class="btnFwa" @click="closeLouDongShiTi">
                  <span class="shadow"></span>
                  <span class="edge"></span>
                  <span class="front text boxShaw"> 关闭 </span>
                </button>
              </div>
              <div slot="content" style="padding: 0">
                <AddDetail ref="AddDetail" />
              </div>
            </Drawer>
            <!-- 新建信息收集实体 -->
            <Drawer :size="widthSize" ref="addShouJi">
              <div class="DrawerTitleLeft boxShaw" slot="header-left">
                <span class="textSL">新建信息收集实体</span>
              </div>
              <div class="DrawerTitleRight" slot="header-right">
                <!-- @click="saveBtn" -->
                <button
                  class="btnFwa"
                  style="margin-right: 10px"
                  @click="saveXinXiShouJi"
                >
                  <span class="shadow"></span>
                  <span class="edge"></span>
                  <span class="front text boxShaw" style="background: #67c23a">
                    保存
                  </span>
                </button>
                <!-- @click="closeDraw" -->
                <button class="btnFwa" @click="closeXinXiShouJi">
                  <span class="shadow"></span>
                  <span class="edge"></span>
                  <span class="front text boxShaw"> 关闭 </span>
                </button>
              </div>
              <div slot="content" style="padding: 0">
                <AddInfoDetail ref="ShouJiDetail" />
              </div>
            </Drawer>

            <!-- Echart 展示1万条以上的数据 -->
            <!-- <Echart
              ref="Ech"
              v-else-if="MiddleSlectValue == '大数据(万)'"
            ></Echart> -->
            <!-- loading 效果 -->
            <!-- <div style="flex: 1" v-if="false">
              <Loading></Loading>
            </div> -->
          </div>
        </div>
      </div>
    </fullscreen>
  </div>
</template>

<script>
import { getCommonOptions, creFinger } from "@/api/vul/scripts"

// import Neo4j from "@/components/Neo4j/neo4jD3/index.vue"
import HaoYing from "@/components/Neo4j/Three/index.vue"
import G6 from "@/components/Neo4j/G6antv/index.vue"
import Echart from "@/components/Neo4j/Echart/index.vue"
import Loading from '@/components/Neo4j/loading/index.vue'
import CorrelationAtlas from '@/api/CorrelationAtlas'
import Drawer from '@/components/Neo4j/G6antv/Drawer.vue'
import AddDetail from '@/components/Neo4j/G6antv/addDetail.vue'
import AddInfoDetail from '@/components/Neo4j/G6antv/AddInfoDetail.vue'
export default ({
  inject: ['reload'],
  components: {
    // Neo4j,
    HaoYing,
    G6,
    Echart,
    Loading,
    Drawer,
    AddDetail,
    AddInfoDetail
  },
  data () {
    return {
      widthSize: '40%',
      isThree: false,
      chanelCount: {},
      fullscreen: false,
      nodeData: {},
      MiddleSlectValue: '查看全部',
      threeData: {},
      G6Data: {},
      createFingerForm: {
        "app_name": "",
        "app_class": "",
        "app_version": "",
        "cn_name": "",
        "flag": "",
      },
    }
  },
  methods: {
    // 点击抽屉保存-新建指纹
    saveDraw () {
      this.addFinger(this.createFingerForm)
      this.createFingerForm = {}
      this.$refs.addTa.table = false
    },
    // 点击抽屉关闭-新建指纹
    closeDraw () {
      this.createFingerForm = {}
      this.$refs.addTa.table = false
    },
    // 点击信息收集抽屉保存
    saveInfo () {
      // this.addFinger(this.createFingerForm)
      // this.createFingerForm = {}
      this.$refs.addInfo.table = false
    },
    // 点击信息收集抽屉关闭
    closeInfo () {
      // this.createFingerForm = {}
      this.$refs.addInfo.table = false
    },
    // 点击新建漏洞实体抽屉保存
    saveLouDongShiTi () {
      // this.addFinger(this.createFingerForm)
      // this.createFingerForm = {}
      this.$refs.addLouDong.table = false
      this.$refs.AddDetail.save()
      // this.reload()
    },
    // 点击新建漏洞实体抽屉关闭
    closeLouDongShiTi () {
      // this.createFingerForm = {}
      this.$refs.addLouDong.table = false
    },
    // 点击信息收集实体抽屉保存
    saveXinXiShouJi () {
      this.$refs.addShouJi.table = false
      this.$refs.ShouJiDetail.save()
      // this.reload()
    },
    // 点击信息收集实体抽屉关闭
    closeXinXiShouJi () {
      // this.createFingerForm = {}
      this.$refs.addShouJi.table = false
    },
    // 点击指纹
    finger () {
      this.$refs.addTa.table = true
    },
    // 点击信息收集实体
    addShouJiTI () {
      this.$refs.addShouJi.table = true
    },
    // 点击漏洞实体
    addLouDOng () {
      this.$refs.addLouDong.table = true
    },
    // 点击信息收集
    addInfo () {
      this.$refs.addInfo.table = true
    },
    // neo4jd3初始化
    iniGlobalData () {
      getCommonOptions({}).then(res => {
        localStorage.setItem("global_options", JSON.stringify(res))
      })
    },
    // neo4jd3-鼠标经过节点拿到的数据
    mouseNode (data) {
      this.nodeData = data
      console.log(data, '--------------')
    },
    // neo4jd3-鼠标离开节点事件--隐藏显示框
    mouseLeave () {
      this.nodeData = {}
    },
    // three组件通知父组件发请求
    initImIt () {
      if (this.MiddleSlectValue == '查看100条') {
        this.getTuPu(100)
      } else if (this.MiddleSlectValue == '查看全部') {
        this.getTuPu(0)
      }
    },
    // G6组件通知父组件发请求
    initG6 () {
      if (this.MiddleSlectValue == '查看100条') {
        this.getTuPu(100)
      } else if (this.MiddleSlectValue == '查看全部') {
        this.getTuPu(0)
      }
    },
    /* 点击下拉
       应该发请求
    */
    clickDropdown (data) {
      // console.log(data, '下拉下拉下拉下拉下拉下拉')
      this.MiddleSlectValue = data
      if (data == '查看100条') {
        this.getTuPu(100, data)
      } else if (data == '查看全部') {
        this.getTuPu(0, data)
      }
    },
    /* 关联图谱请求接口
      0代表全部数据
      其他数字代表具体请求的数量
    */
    async getTuPu (num = 0, name) {
      try {

        const res = await CorrelationAtlas.getAtlas({ num })
        if (num == 0) {
          this.threeData = res.data
          this.G6Data = res.data
          this.chanelCount = res.data.statistic_count
        } else {
          this.threeData = res.data
          this.G6Data = res.data
          this.chanelCount = res.data.statistic_count
        }
      } catch (error) {
        console.log(error)
      }
    },
    /* 新建指纹
    */
    async addFinger (obj) {
      try {
        const res = await creFinger(obj)
        if (res.msg == '成功') {
          this.$message.success('新建指纹-成功!')
        }
      } catch (error) {
        console.log(error)
      }
    },

  },
  mounted () {
    try {
      this.iniGlobalData()
      /* 默认情况下是小数据 */
      this.getTuPu(0, '查看全部')
    } catch (error) {
      console.log(error, '595------------')
    }
  },
})

</script>


<style scoped lang="less">
@import "../../components/Neo4j/css/btn.less";
/deep/.el-form-item__label {
  text-align: left;
}
/deep/.blueChannel .el-tabs__item {
  padding: 0 20px 0 0;
}
/deep/ .el-input {
  width: 217px;
}
/deep/.el-tabs__content {
  font-size: 12px;
}
.search-box,
.userbox,
.neo4j,
.neo4j2 {
  height: calc(100% - 35px);
}
.u-txt {
  display: inline;
}

.u-btn {
  display: none;
}

.el-table__body tr:hover {
  .u-btn {
    display: inline;
  }

  .u-txt {
    display: none;
  }
}

.learnMore {
  min-width: 100px !important;
  padding: 16px 0 !important;
}

.userbox {
  padding: 24px;
  background: #fff;
  min-height: calc(100% - 39px);
  box-sizing: border-box;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
}

.txtareacontent /deep/ textarea {
  resize: none !important;
}

.newUserDialog /deep/ .el-form-item__label {
  position: absolute;
  left: 127px;
}

.newUserDialog /deep/ .el-form-item__error {
  left: 112px;
}

.popover_block {
  list-style: none;

  li {
    text-align: center;
    margin-top: 8px;
  }
}

.userblock {
  display: inline-block;
}

/deep/ .userblock .el-date-editor.el-input,
.el-date-editor.el-input__inner {
  width: 320px !important;
}

/deep/ .el-popconfirm {
  .el-popconfirm__main {
    margin-bottom: 16px;
    color: rgba(72, 72, 102, 0.64);
  }
}
</style>

/* 蓝色面板样式 */
<style scoped lang='less'>
.blueChannel {
  font-size: 14px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 68px;
  background: #4c7ae3;
  border-radius: 5px;
  .blueChannelLeft {
    padding: 0 20px;
    .showShan {
      width: 700px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      color: #fff;

      .spanNum {
        margin: 0 10px 0 5px;
      }
    }
  }
  .blueChannel-right {
    color: #fff;
    width: 600px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    color: #fff;
    .rightIcon {
      font-size: 20px;
      width: 260px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-top: 10px;
    }
  }
  .MiddleXia {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .MiddleXiaButton {
      --color: #0077ff;
      display: inline-block;
      width: 100px;
      height: 38px;
      line-height: 35px;
      overflow: hidden;
      // margin: 20px;
      font-size: 14px;
      z-index: 1;
      color: var(--color);
      border: 2px solid #4c7ae3;
      border-radius: 3px;
      position: relative;
    }

    .MiddleXiaButton::before {
      position: absolute;
      content: "";
      background: var(--color);
      width: 150px;
      height: 200px;
      z-index: -1;
      border-radius: 50%;
    }

    .MiddleXiaButton:hover {
      color: white;
    }

    .MiddleXiaButton:before {
      top: 100%;
      left: 100%;
      transition: 0.3s all;
    }

    .MiddleXiaButton:hover::before {
      top: -30px;
      left: -30px;
    }
  }
  .boxShaw {
    text-shadow: 0 0 1px #ddd, 0 0 1px #ddd, 0 0 2px #ddd, 0 0 2px #ddd !important;
  }
}

/* 下拉 */
.el-dropdown-link {
  cursor: pointer;
  color: #fff;
}
.el-icon-arrow-down {
  font-size: 12px;
}
.button1 {
  display: flex;
  width: 120px;
  height: 40px;
  align-items: center;
  font-family: inherit;
  font-weight: 500;
  font-size: 15px;
  padding: 9px;
  color: white;
  background: #ad5389;
  background: linear-gradient(to right, #0f0c29, #302b63, #24243e);
  border: none;
  letter-spacing: 0.05em;
  border-radius: 16px;
}

.button1 svg {
  margin-right: 3px;
  transform: rotate(30deg);
  transition: transform 0.5s cubic-bezier(0.76, 0, 0.24, 1);
}

.button1 span {
  transition: transform 0.5s cubic-bezier(0.76, 0, 0.24, 1);
}

.button1:hover svg {
  transform: translateX(5px) rotate(90deg);
}

.button1:hover span {
  transform: translateX(7px);
}
</style>
/* 
  card组件样式
*/
 <style  scoped lang="less">
.cardX {
  width: 200px;
  height: 200px;
  position: absolute;
  top: 0;
  right: 0%;
  background: #bbb;
  display: flex;
  padding: 10px;
  place-content: center;
  place-items: center;
  overflow: hidden;
  border-radius: 25px;
  z-index: 999;
}

.cardX h2 {
  z-index: 10;
  color: white;
  font-size: 2em;
}

.cardX::before {
  content: "";
  position: absolute;
  width: 100px;
  background-image: linear-gradient(
    180deg,
    rgb(0, 183, 255),
    rgb(255, 48, 255)
  );
  height: 130%;
  animation: rotBGimg 3s linear infinite;
  transition: all 0.2s linear;
}

@keyframes rotBGimg {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

.cardX::after {
  content: "";
  position: absolute;
  background: #bbb;
  inset: 5px;
  border-radius: 15px;
}
.el-card {
  color: #687075 !important;
}
/deep/ .el-checkbox__input.is-checked .el-checkbox__inner,
.el-checkbox__input.is-indeterminate .el-checkbox__inner {
  // background: red;
  // border-color: red;
}
</style>

<style lang='less' scoped>
.btn {
  display: block;
  width: 90px;
  height: 32px;
  background: transparent;
  outline: none;
  border: 0;
  color: #fff;
  letter-spacing: 0.1em;
  font-family: monospace;
  font-size: 14px;
  font-weight: bold;
  z-index: 1;
}

.cube {
  position: relative;
  transition: all 0.5s;
}

.cube .bg-top {
  position: absolute;
  height: 10px;
  background: #d4af37;
  bottom: 100%;
  left: 5px;
  right: -5px;
  transform: skew(-45deg, 0);
  margin: 0;
  transition: all 0.4s;
}

.cube .bg-top .bg-inner {
  bottom: 0;
}

.cube .bg {
  position: absolute;
  left: 0;
  bottom: 0;
  top: 0;
  right: 0;
  background: #d4af37;
  transition: all 0.4s;
}

.cube .bg-right {
  position: absolute;
  background: #d4af37;
  top: -5px;
  z-index: 0;
  bottom: 5px;
  width: 10px;
  left: 100%;
  transform: skew(0, -45deg);
  transition: all 0.4s;
}

.cube .bg-right .bg-inner {
  left: 0;
}

.cube .bg-inner {
  background: #28282d;
  position: absolute;
  left: 2px;
  right: 2px;
  top: 2px;
  bottom: 2px;
}

.cube .text {
  position: relative;
  transition: all 0.4s;
}

.cube:hover .bg-inner {
  background: #d4af37;
  transition: all 0.4s;
}

.cube:hover .text {
  color: #28282d;
  transition: all 0.4s;
}

.cube:hover .bg-right,
.cube:hover .bg,
.cube:hover .bg-top {
  background: #28282d;
}

.cube:active {
  z-index: 9999;
  animation: bounce 0.1s linear;
}

@keyframes bounce {
  50% {
    transform: scale(0.9);
  }
}
</style>

/* 按钮 */
<style scoped>
.btnFwa {
  position: relative;
  border: none;
  background: transparent;
  padding: 0;
  cursor: pointer;
  outline-offset: 4px;
  transition: filter 250ms;
  user-select: none;
  touch-action: manipulation;
}

.shadow {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 10px;
  background: hsl(0deg 0% 0% / 0.25);
  will-change: transform;
  transform: translateY(2px);
  transition: transform 600ms cubic-bezier(0.3, 0.7, 0.4, 1);
}

.edge {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 12px;
  background: linear-gradient(
    to left,
    hsl(340deg 100% 16%) 0%,
    hsl(340deg 100% 32%) 8%,
    hsl(340deg 100% 32%) 92%,
    hsl(340deg 100% 16%) 100%
  );
}

.front {
  display: block;
  position: relative;
  padding: 6px 14px;
  border-radius: 12px;
  font-size: 14px;
  color: white;
  background: hsl(345deg 100% 47%);
  will-change: transform;
  transform: translateY(-4px);
  transition: transform 600ms cubic-bezier(0.3, 0.7, 0.4, 1);
}

button:hover {
  filter: brightness(110%);
}

button:hover .front {
  transform: translateY(-6px);
  transition: transform 250ms cubic-bezier(0.3, 0.7, 0.4, 1.5);
}

button:active .front {
  transform: translateY(-2px);
  transition: transform 34ms;
}

button:hover .shadow {
  transform: translateY(4px);
  transition: transform 250ms cubic-bezier(0.3, 0.7, 0.4, 1.5);
}

button:active .shadow {
  transform: translateY(1px);
  transition: transform 34ms;
}

button:focus:not(:focus-visible) {
  outline: none;
}
</style>
/* 抽屉样式 */
<style scoped lang='less'>
.DrawerTitleLeft {
  color: #fff;
  font-size: 15px;
  .textSL {
    color: #fff;
    font-size: 15px;
  }
}
.DrawerTitleRight {
}
.boxShaw {
  text-shadow: 0 0 1px #ddd, 0 0 1px #ddd, 0 0 2px #ddd, 0 0 2px #ddd !important;
}
/* 给抽屉里tabs组件写 */
.tabs {
  width: 70%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  .tabs-name {
    width: 100px;
    height: 31px;
    line-height: 31px;
  }
}
</style>