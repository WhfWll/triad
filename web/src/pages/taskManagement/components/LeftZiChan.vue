<template>
  <div style="height: calc(100vh - 200px);" class="dashboard-component">
    <!-- 资产统计 -->
    <div class="asset-statistics">
      <!-- 第一行 - 3列 -->
      <el-row :gutter="2" style="padding: 0 1px">
        <el-col :span="8">
          <div class="vuln_evidence" style="background: #f1f3ff">
            <div class="risk_icon" style="background: #e2e6ff">
              <img src="../../../assets/images/folder@2x.png" />
            </div>
            <div class="risk_content" style="padding-left:5px ">
              <p class="risk_name">{{ statisticsRow1[0].label }}</p>
              <p class="risk_number">
                <label>{{ statisticsRow1[0].value }}</label>
              </p>
            </div>
          </div>
        </el-col>
        <el-col :span="7">
          <div class="vuln_evidence" style="background: #f2faf4">
            <div class="risk_icon" style="background: #e1f4e6">
              <img src="../../../assets/images/server@2x.png" />
            </div>
            <div class="risk_content" style="padding-left:5px ">
              <p class="risk_name">{{ statisticsRow1[1].label }}</p>
              <p class="risk_number">
                <label>{{ statisticsRow1[1].value }}</label>
              </p>
            </div>
          </div>
        </el-col>
        <el-col :span="9">
          <div class="vuln_evidence" style="background: #faf5fd">
            <div class="risk_icon" style="background: #f5e8ff">
              <img src="../../../assets/images/file-text@2x.png" />
            </div>
            <div class="risk_content" style="padding-left:5px ">
              <p class="risk_name">{{ statisticsRow1[2].label }}</p>
              <p class="risk_number">
                <label>{{ statisticsRow1[2].value }}</label>
              </p>
            </div>
          </div>
        </el-col>
      </el-row>
      <!-- <div class="row">
        <div
          class="column"
          v-for="(item, index) in statisticsRow1"
          :key="`row1-col${index}`"
        >
          <i :class="`el-icon-${item.icon}`"></i>
          <span  class="label">{{ item.label }}</span>
          <span class="value">{{ item.value }}</span>
        </div>
      </div> -->
      <!-- 第二行和第三行 - 每行4列 -->
      <div style=" background: linear-gradient(135deg, #f7f9fc 0%, #f7f9fc 100%);
    padding:0 15px;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);margin-top:20px" class="row" v-for="row in [statisticsRow2]" :key="row[0].label">
        <div
        
          class="column"
          v-for="(item, index) in row"
          :key="`row-col${index}`"
        >
          <i :class="`icon ${item.icon}`"></i>
        
          <div class="vuln_evidence flex-center-vertical-horizontal" >
            <div class="risk_icon">
              <svg
                v-if="index == 0"
                t="1709180679312"
                class="icon"
                viewBox="0 0 1024 1024"
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
                p-id="3237"
                width="25"
                height="25"
              >
                <path
                  d="M512.1 512.4m-448 0a448 448 0 1 0 896 0 448 448 0 1 0-896 0Z"
                  fill="#ffe7e6"
                  p-id="3238"
                ></path>
                <path
                  d="M739.3 622.8H284.8c-15.3 0-27.8-12.5-27.8-27.8V329.6c0-15.3 12.5-27.8 27.8-27.8h454.5c15.3 0 27.8 12.5 27.8 27.8V595c0 15.3-12.5 27.8-27.8 27.8zM840.6 712.3h-657c-4.1 0-6.2-5-3.3-7.9l70.9-72.2c0.9-0.9 2.1-1.4 3.3-1.4h506.3c1.2 0 2.3 0.4 3.1 1.2l79.7 72.2c3.2 2.9 1.2 8.1-3 8.1z"
                  fill="#f5212d"
                  p-id="3239"
                ></path>
                <path
                  d="M840.6 712.3h-657c-4.1 0-6.2 0.6-3.3 1l70.9 9.4c0.9 0.1 2.1 0.2 3.3 0.2h506.3c1.2 0 2.3-0.1 3.1-0.2l79.7-9.4c3.2-0.3 1.2-1-3-1z"
                  fill="#ffe7e6"
                  p-id="3240"
                ></path>
                <path
                  d="M638.9 634.8h-46.7c-4.4 0-8-3.6-8-8s3.6-8 8-8h46.7c4.4 0 8 3.6 8 8s-3.6 8-8 8zM431.9 634.8h-46.7c-4.4 0-8-3.6-8-8s3.6-8 8-8h46.7c4.4 0 8 3.6 8 8s-3.6 8-8 8z"
                  fill="#FFFFFF"
                  p-id="3241"
                ></path>
                <path
                  d="M763.1 673.1h-502l40-28.2h422z"
                  fill="#ff3040"
                  p-id="3242"
                ></path>
                <path
                  d="M597.1 703h-170c-6 0-9.4-7-5.6-11.7l8.4-10.4c1.4-1.7 3.4-2.6 5.6-2.6h153.2c2.2 0 4.2 1 5.6 2.6l8.4 10.4c3.7 4.7 0.4 11.7-5.6 11.7z"
                  fill="#ffe7e6"
                  p-id="3243"
                ></path>
                <path
                  d="M731.1 604h-438c-8.2 0-14.8-6.7-14.8-14.8V335.5c0-8.2 6.7-14.8 14.8-14.8h438c8.2 0 14.8 6.7 14.8 14.8v253.6c0 8.2-6.6 14.9-14.8 14.9z"
                  fill="#ffe7e6"
                  p-id="3244"
                ></path>
              </svg>
              <svg
                v-if="index == 1"
                t="1709180679312"
                class="icon"
                viewBox="0 0 1024 1024"
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
                p-id="3237"
                width="25"
                height="25"
              >
                <path
                  d="M512.1 512.4m-448 0a448 448 0 1 0 896 0 448 448 0 1 0-896 0Z"
                  fill="#fd8636"
                  p-id="3238"
                ></path>
                <path
                  d="M739.3 622.8H284.8c-15.3 0-27.8-12.5-27.8-27.8V329.6c0-15.3 12.5-27.8 27.8-27.8h454.5c15.3 0 27.8 12.5 27.8 27.8V595c0 15.3-12.5 27.8-27.8 27.8zM840.6 712.3h-657c-4.1 0-6.2-5-3.3-7.9l70.9-72.2c0.9-0.9 2.1-1.4 3.3-1.4h506.3c1.2 0 2.3 0.4 3.1 1.2l79.7 72.2c3.2 2.9 1.2 8.1-3 8.1z"
                  fill="#fff"
                  p-id="3239"
                ></path>
                <path
                  d="M840.6 712.3h-657c-4.1 0-6.2 0.6-3.3 1l70.9 9.4c0.9 0.1 2.1 0.2 3.3 0.2h506.3c1.2 0 2.3-0.1 3.1-0.2l79.7-9.4c3.2-0.3 1.2-1-3-1z"
                  fill="#ffff"
                  p-id="3240"
                ></path>
                <path
                  d="M638.9 634.8h-46.7c-4.4 0-8-3.6-8-8s3.6-8 8-8h46.7c4.4 0 8 3.6 8 8s-3.6 8-8 8zM431.9 634.8h-46.7c-4.4 0-8-3.6-8-8s3.6-8 8-8h46.7c4.4 0 8 3.6 8 8s-3.6 8-8 8z"
                  fill="#FFFFFF"
                  p-id="3241"
                ></path>
                <path
                  d="M763.1 673.1h-502l40-28.2h422z"
                  fill="#ff3040"
                  p-id="3242"
                ></path>
                <path
                  d="M597.1 703h-170c-6 0-9.4-7-5.6-11.7l8.4-10.4c1.4-1.7 3.4-2.6 5.6-2.6h153.2c2.2 0 4.2 1 5.6 2.6l8.4 10.4c3.7 4.7 0.4 11.7-5.6 11.7z"
                  fill="#ffe7e6"
                  p-id="3243"
                ></path>
                <path
                  d="M731.1 604h-438c-8.2 0-14.8-6.7-14.8-14.8V335.5c0-8.2 6.7-14.8 14.8-14.8h438c8.2 0 14.8 6.7 14.8 14.8v253.6c0 8.2-6.6 14.9-14.8 14.9z"
                  fill="#fd8636"
                  p-id="3244"
                ></path>
              </svg>
              <svg
                v-if="index == 2"
                t="1709180679312"
                class="icon"
                viewBox="0 0 1024 1024"
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
                p-id="3237"
                width="25"
                height="25"
              >
                <path
                  d="M512.1 512.4m-448 0a448 448 0 1 0 896 0 448 448 0 1 0-896 0Z"
                  fill="#326fee"
                  p-id="3238"
                ></path>
                <path
                  d="M739.3 622.8H284.8c-15.3 0-27.8-12.5-27.8-27.8V329.6c0-15.3 12.5-27.8 27.8-27.8h454.5c15.3 0 27.8 12.5 27.8 27.8V595c0 15.3-12.5 27.8-27.8 27.8zM840.6 712.3h-657c-4.1 0-6.2-5-3.3-7.9l70.9-72.2c0.9-0.9 2.1-1.4 3.3-1.4h506.3c1.2 0 2.3 0.4 3.1 1.2l79.7 72.2c3.2 2.9 1.2 8.1-3 8.1z"
                  fill="#fff"
                  p-id="3239"
                ></path>
                <path
                  d="M840.6 712.3h-657c-4.1 0-6.2 0.6-3.3 1l70.9 9.4c0.9 0.1 2.1 0.2 3.3 0.2h506.3c1.2 0 2.3-0.1 3.1-0.2l79.7-9.4c3.2-0.3 1.2-1-3-1z"
                  fill="#ffff"
                  p-id="3240"
                ></path>
                <path
                  d="M638.9 634.8h-46.7c-4.4 0-8-3.6-8-8s3.6-8 8-8h46.7c4.4 0 8 3.6 8 8s-3.6 8-8 8zM431.9 634.8h-46.7c-4.4 0-8-3.6-8-8s3.6-8 8-8h46.7c4.4 0 8 3.6 8 8s-3.6 8-8 8z"
                  fill="#FFFFFF"
                  p-id="3241"
                ></path>
                <path
                  d="M763.1 673.1h-502l40-28.2h422z"
                  fill="#326fee"
                  p-id="3242"
                ></path>
                <path
                  d="M597.1 703h-170c-6 0-9.4-7-5.6-11.7l8.4-10.4c1.4-1.7 3.4-2.6 5.6-2.6h153.2c2.2 0 4.2 1 5.6 2.6l8.4 10.4c3.7 4.7 0.4 11.7-5.6 11.7z"
                  fill="#326fee"
                  p-id="3243"
                ></path>
                <path
                  d="M731.1 604h-438c-8.2 0-14.8-6.7-14.8-14.8V335.5c0-8.2 6.7-14.8 14.8-14.8h438c8.2 0 14.8 6.7 14.8 14.8v253.6c0 8.2-6.6 14.9-14.8 14.9z"
                  fill="#326fee"
                  p-id="3244"
                ></path>
              </svg>
              <svg
                v-if="index == 3"
                t="1709180679312"
                class="icon"
                viewBox="0 0 1024 1024"
                version="1.1"
                xmlns="http://www.w3.org/2000/svg"
                p-id="3237"
                width="25"
                height="25"
              >
                <path
                  d="M512.1 512.4m-448 0a448 448 0 1 0 896 0 448 448 0 1 0-896 0Z"
                  fill="#babc82"
                  p-id="3238"
                ></path>
                <path
                  d="M739.3 622.8H284.8c-15.3 0-27.8-12.5-27.8-27.8V329.6c0-15.3 12.5-27.8 27.8-27.8h454.5c15.3 0 27.8 12.5 27.8 27.8V595c0 15.3-12.5 27.8-27.8 27.8zM840.6 712.3h-657c-4.1 0-6.2-5-3.3-7.9l70.9-72.2c0.9-0.9 2.1-1.4 3.3-1.4h506.3c1.2 0 2.3 0.4 3.1 1.2l79.7 72.2c3.2 2.9 1.2 8.1-3 8.1z"
                  fill="#fff"
                  p-id="3239"
                ></path>
                <path
                  d="M840.6 712.3h-657c-4.1 0-6.2 0.6-3.3 1l70.9 9.4c0.9 0.1 2.1 0.2 3.3 0.2h506.3c1.2 0 2.3-0.1 3.1-0.2l79.7-9.4c3.2-0.3 1.2-1-3-1z"
                  fill="#ffff"
                  p-id="3240"
                ></path>
                <path
                  d="M638.9 634.8h-46.7c-4.4 0-8-3.6-8-8s3.6-8 8-8h46.7c4.4 0 8 3.6 8 8s-3.6 8-8 8zM431.9 634.8h-46.7c-4.4 0-8-3.6-8-8s3.6-8 8-8h46.7c4.4 0 8 3.6 8 8s-3.6 8-8 8z"
                  fill="#FFFFFF"
                  p-id="3241"
                ></path>
                <path
                  d="M763.1 673.1h-502l40-28.2h422z"
                  fill="#babc82"
                  p-id="3242"
                ></path>
                <path
                  d="M597.1 703h-170c-6 0-9.4-7-5.6-11.7l8.4-10.4c1.4-1.7 3.4-2.6 5.6-2.6h153.2c2.2 0 4.2 1 5.6 2.6l8.4 10.4c3.7 4.7 0.4 11.7-5.6 11.7z"
                  fill="#babc82"
                  p-id="3243"
                ></path>
                <path
                  d="M731.1 604h-438c-8.2 0-14.8-6.7-14.8-14.8V335.5c0-8.2 6.7-14.8 14.8-14.8h438c8.2 0 14.8 6.7 14.8 14.8v253.6c0 8.2-6.6 14.9-14.8 14.9z"
                  fill="#babc82"
                  p-id="3244"
                ></path>
              </svg>
              <!-- <img :src="require(`@/assets/icon/assetManage/${index}.png`)" /> -->
            </div>
            <div class="risk_content">
              <p class="risk_name">{{ item.label }}</p>
              <p class="risk_number">
                <label style="font-size: 14px">{{ item.value }}</label>
              </p>
            </div>
          </div>
        </div>
      </div>
      <div style=" background: linear-gradient(135deg, #f7f9fc 0%, #f7f9fc 100%);
  
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.1);" class="row" v-for="row in [statisticsRow3]" :key="row[0].label">
        <div
          class="column"
          v-for="(item, index) in row"
          :key="`row-col${index}`"
        >
          <i :class="`icon ${item.icon}`"></i>
          <!-- <span class="label">{{ item.label }}</span>
          <span class="value">{{ item.value }}</span> -->
          <div class="vuln_evidence" style="padding-top:0">
            <div class="risk_icon">
              <img :src="require(`@/assets/icon/assetManage/${index}.png`)" />
            </div>
            <div class="risk_content">
              <p class="risk_name">{{ item.label }}</p>
              <p class="risk_number">
                <label style="font-size: 14px">{{ item.value }}</label>
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 最新资产风险变化 -->

    <div v-if="riskChanges&&riskChanges.length" class="asset-risk-changes">
      <div class="list">
        <label style="font-size: 14px;color:#484866A3"> 最新资产风险变化</label>
        <div
          style="font-size: 12px"
          class="list-item"
          v-for="(item, index) in riskChanges"
          :key="index"
        >
          <span class="ip" @click="ipClick(item)">{{ item.ip }}</span>
          <span class="status">{{ item.riskLevelName }}</span>
          <!-- 时间元素 -->
          <span class="date">{{ item.updateTime }}</span>
          <!-- 忽略按钮，初始时不可见 -->
          <span class="ignore-btn" @click="ignoreRisk(item)">忽略</span>
        </div>
      </div>
      <el-button style="font-size: 12px" type="text" @click="loadMore"
        >加载更多</el-button
      >
    </div>
  </div>
</template>

<script>
import { traffic } from '@/api/assetManagement.js'

export default {
  name: 'DashboardComponent',
  props: {
    zichanObj: {
      type: Object,
    },
  },
  data () {
    return {
      // 资产统计数据
      statisticsRow1: [{
        icon: `edit`,
        label: `总资产`,
        value: 0,
      }, {
        icon: `share`,
        label: `总漏洞`,
        value: 0,
      }, {
        icon: `share`,
        label: `月新增漏洞`,
        value: 0,
      }],
      statisticsRow2: [{
        icon: `edit`,
        label: `高危资产`,
        value: 0,
      }, {
        icon: `share`,
        label: `中危资产`,
        value: 0,
      },
      {
        icon: `share`,
        label: `低危资产`,
        value: 0,
      },
      {
        icon: `share`,
        label: `安全资产`,
        value: 0,
      }
      ],
      statisticsRow3: [{
        icon: `edit`,
        label: `致命漏洞`,
        value: 0,
      }, {
        icon: `share`,
        label: `高危漏洞`,
        value: 0,
      }, {
        icon: `share`,
        label: `中危漏洞`,
        value: 0,
      }, {
        icon: `share`,
        label: `低危漏洞`,
        value: 0,
      }
      ],
      // 最新资产风险变化列表
      riskChanges: [],
      // 模拟数据页面
      currentPage: 1,
      size: 100
    }
  },
  created () {
    this.getZCList()
  },
  methods: {

    ipClick (item) {
      console.log(item)
      // 跳转到概览
      this.$emit('clickSOnTOTab','1')

      this.$store.commit('setgroupID', [item.id, 2]);
    },
    ignoreRisk (item) {
      // 忽略风险的逻辑处理
      console.log("忽略风险：", item)
      this.uptodateignore(item.id)
    },
    async getZCList () {
      const res = await traffic.addTraffic({ page: 1, size: this.size })
      if (res.code == 200) {
        this.riskChanges = res.data.list
        console.log(this.riskChanges, 'this.riskChangesthis.riskChanges')
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    // 4、最新资产风险变化忽略（完成）
    async uptodateignore (id) {
      const res = await traffic.delTraffic({ assetId: id })
      if (res.code == 200) {
        this.$message({
          message: res.msg || "操作成功",
          type: 'success'
        })
      } else {
        this.$message({
          message: res.msg,
          type: 'error'
        })
      }
    },
    generateStatisticsData (columns) {
      // 生成统计数据
      return Array.from({ length: columns }, (_, index) => ({
        icon: `edit`,
        label: `项目${index + 1}`,
        value: (index + 1) * 1000,
      }))
    },
    generateRiskChanges (rows) {
      // 生成风险变化数据
      return Array.from({ length: rows }, (_, index) => ({
        ip: `192.168.0.${index + 1}`,
        status: ['正常', '警告', '异常'][index % 3],
        date: new Date().toISOString().split('T')[0],
      }))
    },
    loadMore () {
      // 加载更多功能的实现
      this.size += 20
      this.getZCList()
    },
  },
  watch: {
    zichanObj (newValue) {
      this.statisticsRow1[0].value = newValue.assetTotal
      this.statisticsRow1[1].value = newValue.assetVulTotal
      this.statisticsRow1[2].value = newValue.assetVulMonth

      this.statisticsRow2[0].value = newValue.assetRisk[0]
      this.statisticsRow2[1].value = newValue.assetRisk[1]
      this.statisticsRow2[2].value = newValue.assetRisk[2]
      this.statisticsRow2[3].value = newValue.assetRisk[3]

      this.statisticsRow3[0].value = newValue.assetVulRisk[0]
      this.statisticsRow3[1].value = newValue.assetVulRisk[1]
      this.statisticsRow3[2].value = newValue.assetVulRisk[2]
      this.statisticsRow3[3].value = newValue.assetVulRisk[3]

    }
  },
};
</script>
<style scoped>
.asset-risk-changes .list-item {
  /* ...其他样式 */
  position: relative; /* 设置相对定位 */
  padding: 10px; /* 或者你需要的任何内边距 */
}

.asset-risk-changes .list-item .date {
  /* ...其他样式 */
  transition: opacity 0.3s ease; /* 添加过渡效果 */
  opacity: 1; /* 默认完全不透明 */
}

.asset-risk-changes .list-item .ignore-btn {
  /* 初始样式，使按钮不可见 */
  opacity: 0; /* 完全透明 */
  transition: opacity 0.3s ease; /* 添加过渡效果 */
  cursor: pointer; /* 鼠标悬停时显示手形光标 */
  position: absolute; /* 绝对定位 */
  right: 10px; /* 或根据需要调整 */
  top: 50%; /* 垂直居中 */
  transform: translateY(-50%); /* 向上偏移自身高度的一半 */
  color: #004494;
}

.asset-risk-changes .list-item:hover .date {
  /* 鼠标悬停时变透明 */
  opacity: 0;
}

.asset-risk-changes .list-item:hover .ignore-btn {
  /* 鼠标悬停时变为不透明 */
  opacity: 1;
}
</style>
<style scoped>
.flex-center-vertical-horizontal {
  display: flex;
  flex-direction: column; 
  align-items: center;
  justify-content: center; 
  height: 100%; 
  
}
.dashboard-component {
  font-family: "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
  color: #333;
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  background: #fff;
  border-radius: 8px;
}

/* 资产统计样式 */
.asset-statistics .row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.asset-statistics .column {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  font-weight: 600;

  width: 100%;
  margin: 5px;
}

.asset-statistics .icon {
  font-size: 2em;
  margin-bottom: 8px;
}

.asset-statistics .label {
  font-size: 0.9em;
  color: #888;
  margin-bottom: 4px;
}

.asset-statistics .value {
  font-size: 1.4em;
  font-weight: 600;
  color: #333;
}

/* 风险变化列表样式 */
.asset-risk-changes {
  margin-top: 20px;
}

.asset-risk-changes .list {
  max-height: 400px;
  overflow-y: auto;
}

.asset-risk-changes .list-item {
  display: flex;
  justify-content: space-between;
  padding: 10px;
  border-bottom: 1px solid #eee;
  align-items: center;
}

.asset-risk-changes .list-item:not(:last-child) {
  border-bottom: 1px solid #eaeaea;
}

.asset-risk-changes .list-item span {
  display: block;
  text-align: center;
  margin: 0 5px;
}

.asset-risk-changes .ip {
  flex: 1;
}

.asset-risk-changes .status {
  flex: 1;
  text-transform: capitalize;
}

.asset-risk-changes .date {
  flex: 1;
}

.asset-risk-changes button {
  /* background-color: #0056b3;
  color: white;
  border: none;
  padding: 10px 20px;
  margin-top: 10px;
  cursor: pointer;
  border-radius: 4px;
  font-size: 0.9em; */
}

.asset-risk-changes button:hover {
  /* background-color: #004494; */
}

/* 响应式设计 */
@media (max-width: 768px) {
  .asset-statistics .row {
    flex-direction: column;
  }

  .asset-risk-changes .list-item {
    flex-direction: column;
    text-align: left;
  }

  .asset-risk-changes .list-item span {
    margin: 5px 0;
  }
}
</style>





<style  lang="less" scoped>
.box {
  background: #fff;
  box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
  border-radius: 4px;
}
.bog_bg {
  // background: linear-gradient(0deg, #FFFFFF 1%, #EBEDFA 100%);
}
.box_background {
  background-image: url(../../../assets/images/index_box_bg.png) !important ;
  background-repeat: no-repeat;
  background-size: 100%;
}
.task_count_color {
  color: #0256ff;
}
.task_week_color {
  color: #1a1a1a;
}
.legend {
  margin-top: 26px;
  padding: 0 16px;
  list-style: none;
  li {
    padding-left: 12px;
    margin-bottom: 10px;
    i {
      display: inline-block;
      width: 8px;
      height: 8px;
      border-radius: 50%;
    }
    i.hight {
      background: #ff718b;
    }
    i.middle {
      background: #ff9f43;
    }
    i.low {
      background: #ffd258;
    }
    i.safe {
      background: #66c681;
    }
    span {
      display: inline-block;
      font-size: 13px;
      color: rgba(72, 72, 102, 0.64);
      margin-left: 4px;
    }
    label {
      margin-left: 26px;
      display: inline-block;
      font-size: 14px;
      color: #1a1a1a !important;
      font-weight: 500;
    }
  }
}
.time_group {
  float: right;
  background: #fff;
  border-radius: 8px;
  padding: 4px 0;
  width: 110px;
  span {
    display: inline-block;
    width: 28px;
    height: 20px;
    text-align: center;
    line-height: 20px;
    border-radius: 6px;
    color: rgba(72, 72, 102, 0.64);
    margin: 0 4px;
    font-size: 12px;
    cursor: pointer;
  }
  span.active {
    background: #0256ff;
    color: #fff;
  }
}
.box_title {
  // height: 80px;
  position: relative;
  .box_title_icon {
    display: inline-block;
    width: 70px;
    height: 72px;
    vertical-align: middle;
    margin-top: -10px;
  }
  .icon_task {
    background: url(../../../assets/images/tick@2x.png) no-repeat;
    background-size: 100% 100%;
  }
  .icontarget {
    background: url(../../../assets/images/Bookmark@2x.png) no-repeat;
    background-size: 100% 100%;
  }
  .iconvuln {
    background: url(../../../assets/images/Danger@2x.png) no-repeat;
    background-size: 100% 100%;
  }
  .box_title_name {
    display: inline-block;
    color: rgba(72, 72, 102, 0.87);
    font-size: 14px;
    margin-left: -6px;
    font-weight: 500;
  }
  .box_search_time {
    position: absolute;
    top: 0;
    right: 0;
    padding: 18px 16px;
    width: 330px;
    height: 100px;
    // background: linear-gradient(-270deg, #F4D9FF 0%, #CDD4FF 100%);
    // filter: blur(78px);
  }
}
.box_no_icon_title {
  padding: 16px;
  box-sizing: border-box;
  display: inline-block;
  .box_title_name {
    display: inline-block;
    color: rgba(72, 72, 102, 0.87);
    font-size: 14px;
    font-weight: 500;
  }
}
.task_statistics {
  height: 50px;
  padding: 0 16px;
  > span {
    display: inline-block;
    > label {
      display: block;
      font-size: 13px;
      color: rgba(72, 72, 102, 0.64);
      margin-bottom: 12px;
    }
    > span {
      display: block;
      font-size: 24px;
    }
  }
}
.vuln_tj {
  .name {
    font-size: 14px;
    color: rgba(72, 72, 102, 0.64);
  }
  .number {
    font-size: 14px;
    margin: 6px 0;
  }
  .vuln_D {
    color: #ff718b;
  }
  .vuln_h {
    color: #fcb5c3;
  }
  .vuln_m {
    color: #ff9f43;
  }
  .vuln_l {
    color: #ffd258;
  }
  .percentage {
    font-size: 12px;
    color: #1a1a1a;
  }
}
.vuln_evidence {
  height: 66px;
  border-radius: 8px;
  padding-top: 12px;
  padding-left: 12px;
  box-sizing: border-box;
  .risk_icon {
    display: inline-block;
    width: 32px;
    height: 32px;
    line-height: 32px;
    text-align: center;
    // margin-right: 8px;
    border-radius: 6px;
    img {
      width: 20px;
      height: 20px;
      vertical-align: middle;
    }
  }
  .risk_content {
    
    display: inline-block;
    vertical-align: middle;
    .risk_number {
      font-size: 14px;
      margin-bottom: 8px;
      line-height: 14px;
      color: #1a1a1a;
      font-weight: 500;
      label {
        vertical-align: middle;
      }
      i {
        font-size: 12px;
      }
    }
    .risk_name {
      font-size: 12px;
      color: rgba(72, 72, 102, 0.64) !important;
      margin-bottom: 4px;
    }
  }
}
.permeate {
  list-style: none;
  // padding: 0 16px;
  li {
    margin-bottom: 16px;
    line-height: 33px;
    > div {
      display: inline-block;
      width: 32px;
      height: 32px;
      line-height: 32px;
      text-align: center;
      border-radius: 50%;
      border: 1px solid #eeeeee;
      img {
        width: 20px;
        height: 20px;
        vertical-align: middle;
        margin-top: -2px;
      }
    }
    label {
      display: inline-block;
      font-size: 13px;
      color: rgba(72, 72, 102, 0.64);
      margin-left: 8px;
    }
    span {
      float: right;
      display: inline-block;
      color: #1a1a1a;
      font-size: 14px;
      font-weight: 500;
    }
  }
}
.msgstat {
  padding: 0 16px;
  list-style: none;
  overflow-y: scroll;
  height: 170px;
  li {
    font-size: 12px;
    margin-bottom: 12px;
    overflow: hidden;
    > span:nth-child(1) {
      width: calc(100% - 177px);
      display: inline-block;
      vertical-align: middle;
      margin-right: 12px;
    }
    > span .taskname {
      width: 100%;
      color: #1a1a1a;
    }
    > span.time {
      //    float: right;
      //    margin-right: 24px;
      color: rgba(72, 72, 102, 0.64);
      height: 16px;
      line-height: 16px;
      display: inline-block;
      width: 116px;
    }
    > span.content {
      display: inline-block;
    }
    > label {
      display: inline-block;
      color: #4c7ae3;
      width: 55px;
      text-align: right;
    }
  }
}
::-webkit-scrollbar-track {
  background-color: #fff !important;
}
</style>