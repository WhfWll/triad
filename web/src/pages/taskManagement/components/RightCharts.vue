<template>
  <div class="indexbox">
   
   <div style="display:flex;justify-content: space-between;align-items: center;">
     <div style="height: 230px;flex:1" class="box">
      <div class="box_title">
        <div
          style="padding: 16px; box-sizing: border-box; display: inline-block"
        >
          <label for="" class="">资产渗透测试趋势</label>
        </div>

        <div
          class="time_group"
          style="
            border: 1px solid #e8e8f5;
            margin-right: 16px;
            margin-top: 16px;
          "
        >
          <span :class="{ active: zichanceshiNum === 1 }" @click="clickZYN(1)">周</span>
          <span :class="{ active: zichanceshiNum === 2 }" @click="clickZYN(2)">月</span>
          <span :class="{ active: zichanceshiNum === 3 }" @click="clickZYN(3)">年</span>
        </div>
      </div>
      <div id="vuln_line" style="height: 180px"></div>
     </div>
     <div style="height: 230px;flex:1" class="box">
      <div class="box_title">
        <div
          style="padding: 16px; box-sizing: border-box; display: inline-block"
        >
          <label for="" class="">资产漏洞变化趋势</label>
        </div>

        <div
          class="time_group"
          style="
            border: 1px solid #e8e8f5;
            margin-right: 16px;
            margin-top: 16px;
          "
        >
          <span :class="{ active: zichanceshiNum2 === 1 }" @click="clickZYN2(1)">周</span>
          <span :class="{ active: zichanceshiNum2 === 2 }" @click="clickZYN2(2)">月</span>
          <span :class="{ active: zichanceshiNum2 === 3 }" @click="clickZYN2(3)">年</span>
        </div>
      </div>
      <div id="vuln_line2" style="height: 180px"></div>
     </div>
   </div>
    <div style="height: 270px" class="box">
      <div class="box_title" style="position: relative;">
        <div
          style="padding: 16px; box-sizing: border-box; display: inline-block;"
        >
          <label for="" class="">漏洞类型统计</label>
          <div style="position: absolute;bottom:0%;left:40%;display:flex;  align-items: center;">
            <div style="display:flex;">
              <div style="background:red;width:10px;height:10px;border-radius:50%;margin-right:5px"></div>
              <span style="font-size:12px">致命</span>
            </div>
            <div style="display:flex;margin-left:10px">
              <div style="background:#f8b761;width:10px;height:10px;border-radius:50%;margin-right:5px"></div>
              <span style="font-size:12px">高危</span>
            </div>
            <div style="display:flex;margin-left:10px">
              <div style="background:#4c7ae3;width:10px;height:10px;border-radius:50%;margin-right:5px"></div>
              <span style="font-size:12px">中危</span>
            </div>
            <div style="display:flex;margin-left:10px">
              <div style="background:#70b603;width:10px;height:10px;border-radius:50%;margin-right:5px"></div>
              <span style="font-size:12px">低危</span>
            </div>
           
          </div>
        </div>
        <div
          class="time_group"
          style="
            border: 1px solid #e8e8f5;
            margin-right: 16px;
            margin-top: 16px;
          "
        >
         
        </div>
      </div>
      <div id="vuln_type_bar" style="height: 210px"></div>
    </div>
    <div style="height: 270px" class="box">
      <div class="box_title">
        <div
          style="padding: 16px; box-sizing: border-box; display: inline-block"
        >
          <label for="" class="">服务类型统计</label>
        </div>
        <div
          class="time_group"
          style="
            border: 1px solid #e8e8f5;
            margin-right: 16px;
            margin-top: 16px;
          "
        >
         
        </div>
      </div>
      <div id="vuln_type_bar2" style="height: 210px"></div>
    </div>

  
  </div>
</template>  
<style  lang="less" scoped>
.active {
  background-color: #0256ff;
  color: white;
}
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
    margin-right: 8px;
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
      font-size: 13px;
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
<script>   
var echarts = require('echarts')
import index from '@/api/index.js'
import { traffic } from '@/api/assetManagement.js'

export default ({
  name: 'RightCharts',
  data () {
    return {
      zichanceshiNum:2,
      zichanceshiNum2:2,
      class1: '',
      class2: 'active',
      class3: '',
      vulntype1: '',
      vulntype2: 'active',
      vulntype3: '',
      findtrend1: '',
      findtrend2: 'active',
      findtrend3: '',
      vulExploitImpact: [
        {
          label: '高危',
          value: 12,
          percentage: 50
        },
        {
          label: '中危',
          value: 12,
          percentage: 50
        },
        {
          label: '低危',
          value: 12,
          percentage: 50
        }
        , {
          label: '安全',
          value: 12,
          percentage: 50
        }
      ],
      taskCount: 0,
      taskmodel: 2,
      vultypemode: 2, //漏洞类型 统计模式 1-周 2-月 3-年，默认为2
      vulfindtrend: 2,
      tool: {
        "vulCount": 0,
        "fingerCount": 0,
        "taskSceneCount": 0
      },
      taskvulriskstat: {
        "fatalCount": 0,
        "highCount": 0,
        "mediumCount": 0,
        "lowCount": 0,
        "taskVulRiskCount": 0
      },
      targetriskstat: {
        "highCount": 0,
        "mediumCount": 0,
        "lowCount": 0,
        "safeCount": 0,
        "targetRiskCount": 0
      },
      vulevidencestat: {
        dbCount: 0,
        fileLeakCount: 0,
        infoLeakCount: 0,
        loginCredentialsCount: 0,
        remoteControlCount: 0,
      },
      latestWeekTaskCount: 0,
      msgstat: [],
      timer: null,
      chart1: '',
      ebartarget: '',
      evulfindtrend: '',
      eVulnType: '',
      epie: '',
    }
  },
  created: function () {
    // this.$store.state.activefirstMenu="/index"; 

  },
  beforeDestroy () {
    //页面销毁清除定时器
    this.timer ? clearInterval(this.timer) : ""

  },
  mounted: function () {
    this.init()
    // 5分钟请求一次
    this.timer = setInterval(() => {
      this.init()
    }, 300000)

    window.onresize = () => {
      return (() => {
        this.$nextTick(() => {
          // this.chart1.resize()
          // this.ebartarget.resize()
          this.evulfindtrend.resize()
          this.eVulnType.resize()
          // this.epie.resize()
        })
      })()
    }

  },
  methods: {
    // 点击资产渗透测试趋势的--周月年
    clickZYN(num){
      this.zichanceshiNum = num
      this.getvulfindtrendstat()
    },
    clickZYN2(num){
      this.zichanceshiNum2 = num
      this.getvulfindtrendstat2()
    },
    init: function () {
    
      this.getvultypestat()
      this.getvultypestat2()
      this.getvulfindtrendstat() //资产渗透测试趋势
      this.getvulfindtrendstat2()



    },
    async gettaskinfostat () { //任务统计nt,必须,1-周 2-月 3-年
      const res = await traffic.trafficStatus({
        mode: this.zichanceshiNum
      })
      if (res.code == 200) {
        let xdt = res.data.date
        let xdy = res.data.count
        this.echartline1(xdt, xdy)
        this.taskCount = res.data.taskCount
        this.latestWeekTaskCount = res.data.latestWeekTaskCount
      }
    },
    async gettargetriskstat () {
      const res = await index.gettargetriskstat()
      if (res.code == 200) {
        this.targetriskstat = res.data
        let total = this.targetriskstat.targetRiskCount
        let h_p = ((this.targetriskstat.highCount / total) * 100).toFixed(0)
        let m_p = ((this.targetriskstat.mediumCount / total) * 100).toFixed(0)
        let l_p = ((this.targetriskstat.lowCount / total) * 100).toFixed(0)
        let s_p = 100 - h_p - m_p - l_p
        let _data = [
          { label: '高危', value: this.targetriskstat.highCount, p: h_p },
          { label: '中危', value: this.targetriskstat.mediumCount, p: m_p },
          { label: '低危', value: this.targetriskstat.lowCount, p: l_p },
          { label: '安全', value: this.targetriskstat.safeCount, p: s_p }
        ]
        this.echartbartarget(total, _data)
      }
    },
    async getvultypestat () { //漏洞类型统计 
      const res = await traffic.downloadhttps({
        mode: this.vultypemode
      })
      if (res.code == 200) {
        let xdt = res.data.vulType
        let xdy = res.data.count
        this.echartVulnType(xdt, xdy)
      }
    },
    async getvultypestat2 () { //fuwu类型统计 
      const res = await traffic.trafficlist({
        mode: this.vultypemode
      })
      if (res.code == 200) {
        let xdt = res.data.vulType
        let xdy = res.data.count
        this.echartVulnType2(xdt, xdy)
      }
    },
    async getvulfindtrendstat () { //资产渗透测试趋势
     const res = await traffic.trafficStatus({
        mode: this.zichanceshiNum
      })
      if (res.code == 200) {
        let xdt = res.data.date
        let xdy = res.data.count
        this.echartvulfindtrend(xdt, xdy)
      }
    },
    async getvulfindtrendstat2 () { //6、漏洞发现趋势
     const res = await traffic.trafficbase({
        mode: this.zichanceshiNum2
      })
      if (res.code == 200) {
        let xdt = res.data.date
        let backendData = res.data
        this.echartvulfindtrend2(xdt, backendData)
      }
    },
    async getmessagestat () { //最新消息
      const res = await index.getmessagestat()
      if (res.code == 200) {
        this.msgstat = res.data.list
      }
    },
    async gettaskvulriskstat () { //漏洞统计
      const res = await index.gettaskvulriskstat()
      if (res.code == 200) {
        this.taskvulriskstat = res.data
        this.vulExploitImpact = []
        this.vulExploitImpact.push({
          color: 'vuln_D',
          label: '致命',
          value: this.taskvulriskstat.fatalCount,
          percentage: ((this.taskvulriskstat.fatalCount / this.taskvulriskstat.taskVulRiskCount) * 100).toFixed(0)
        })
        this.vulExploitImpact.push({
          color: 'vuln_h',
          label: '高危',
          value: this.taskvulriskstat.highCount,
          percentage: ((this.taskvulriskstat.highCount / this.taskvulriskstat.taskVulRiskCount) * 100).toFixed(0)
        })
        this.vulExploitImpact.push({
          color: 'vuln_m',
          label: '中危',
          value: this.taskvulriskstat.mediumCount,
          percentage: ((this.taskvulriskstat.mediumCount / this.taskvulriskstat.taskVulRiskCount) * 100).toFixed(0)
        })
        this.vulExploitImpact.push({
          color: 'vuln_l',
          label: '低危',
          value: this.taskvulriskstat.lowCount,
          percentage: ((this.taskvulriskstat.lowCount / this.taskvulriskstat.taskVulRiskCount) * 100).toFixed(0)
        })
        this.echartpie()
      }
    },
    async getvulevidencestat () { //漏洞取证
      const res = await index.getvulevidencestat()
      if (res.code == 200) {
        this.vulevidencestat = res.data
      }
    },
    async gettoolinfostat () {
      const res = await index.gettoolinfostat()
      if (res.code == 200) {
        this.tool = res.data
      }
    },

    echartline1: function (xdt, ydt) {
      // 基于准备好的dom，初始化echarts实例
      var myChart = echarts.init(document.getElementById('taskTj'))
      this.chart1 = myChart
      // 绘制图表
      myChart.setOption({
        grid: {
          left: '20',
          top: '15%',
          bottom: '14',
          right: '5%',
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          // formatter :'{b} <br />任务数量：{c} '
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          axisLine: {
            lineStyle: {
              color: '#d9e1e4'
            },
          },
          axisLabel: {
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          axisTick: {
            show: false,
          },
          data: xdt
        },
        yAxis: {
          axisLine: {
            lineStyle: {
              color: '#d9e1e4'
            },
          },
          axisLabel: {
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,

          },
          axisTick: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              type: 'dashed'
            }
          }
        },
        series: [{
          data: ydt,
          type: 'line',
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              {
                offset: 0,
                color: '#E0EAFF'
              },
              {
                offset: 1,
                color: '#fff'
              }
            ])
          },
          lineStyle: {
            color: '#4C7AE3'
          },
          itemStyle: {

            color: '#4C7AE3'
          },
          showSymbol: false,

        }]
      })
    },

    echartbartarget (total, _data) {
      var myChart = echarts.init(document.getElementById('bar'))
      this.ebartarget = myChart
      // 绘制图表
      myChart.setOption({
        color: ['#FF718B', '#FF9F43', '#FFD258', '#66C681',],
        grid: {
          left: '10',
          top: '18',
          bottom: '8',
          right: '5',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: ['高危', '中危', '低危', '安全'],
          axisLabel: {
            show: true,
          },
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          splitLine: {
            show: false
          }
        },
        yAxis: {
          type: 'value',
          axisLabel: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          axisLine: {
            show: false,
          },
          splitLine: {
            show: false
          }
        },
        series: [
          {
            colorBy: 'data',
            data: _data,
            type: 'bar',
            // barGap:'30%',
            barWidth: '30',
            label: {
              show: true,
              position: 'top',
              fontWeight: '500',
              fontSize: 12,
              color: '#1A1A1A',
              formatter: function (name) {
                return name.data.p + '%'
              },

            },

          }
        ]
      })
    },

    echartvulfindtrend: function (xdt, ydt) { //漏洞趋势
      // 基于准备好的dom，初始化echarts实例
      var myChart = echarts.init(document.getElementById('vuln_line'))
      this.evulfindtrend = myChart
      // 绘制图表
      myChart.setOption({
        grid: {
          left: '20',
          top: '5%',
          bottom: '14',
          right: '5%',
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',

        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          axisLine: {
            lineStyle: {
              color: '#d9e1e4'
            },
          },
          axisLabel: {
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          axisTick: {
            show: false,
          },
          data: xdt
        },
        yAxis: {
          type: 'value',
          axisLine: {
            lineStyle: {
              color: '#d9e1e4'
            },
          },
          axisLabel: {
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          axisTick: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              type: 'dashed'
            }
          }
        },
        series: [{
          data: ydt,
          type: 'line',
          showSymbol: false,
          lineStyle: {
            color: '#4C7AE3'
          },
          itemStyle: {
            borderWidth: 1,
            color: '#4C7AE3'
          },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              {
                offset: 0,
                color: '#E0EAFF'
              },
              {
                offset: 1,
                color: '#fff'
              }
            ])
          },
        }]
      })
    },
    echartvulfindtrend2: function (xdt, backendData) { //漏洞趋势
      // 基于准备好的dom，初始化echarts实例
      var myChart = echarts.init(document.getElementById('vuln_line2'))
      this.evulfindtrend = myChart
      // 绘制图表
      myChart.setOption({
        legend: {
    data: ['致命', '高危', '中危', '低危'],
    // 可以根据需要调整图例的位置和样式
    top: 10
  },
        grid: {
          left: '20',
          top: '5%',
          bottom: '14',
          right: '5%',
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',

        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          axisLine: {
            lineStyle: {
              color: '#d9e1e4'
            },
          },
          axisLabel: {
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          axisTick: {
            show: false,
          },
          data: xdt
        },
        yAxis: {
          type: 'value',
          axisLine: {
            lineStyle: {
              color: '#d9e1e4'
            },
          },
          axisLabel: {
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          axisTick: {
            show: false,
          },
          splitLine: {
            lineStyle: {
              type: 'dashed'
            }
          }
        },
        series: [
    {
      name: '致命',
      data: backendData.fatalCount,
      type: 'line',
      // 其他样式设置...
    },
    {
      name: '高危',
      data: backendData.highCount,
      type: 'line',
      // 其他样式设置...
    },
    {
      name: '中危',
      data: backendData.middleCount,
      type: 'line',
      // 其他样式设置...
    },
    {
      name: '低危',
      data: backendData.lowCount,
      type: 'line',
      // 其他样式设置...
    }
  ]
        // series: [{
        //   data: ydt,
        //   type: 'line',
        //   showSymbol: false,
        //   lineStyle: {
        //     color: '#4C7AE3'
        //   },
        //   itemStyle: {
        //     borderWidth: 1,
        //     color: '#4C7AE3'
        //   },
        //   areaStyle: {
        //     color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        //       {
        //         offset: 0,
        //         color: '#E0EAFF'
        //       },
        //       {
        //         offset: 1,
        //         color: '#fff'
        //       }
        //     ])
        //   },
        // }]
      })
    },
    echartVulnType (xdt, xdy) {
      var myChart = echarts.init(document.getElementById('vuln_type_bar'))
      this.eVulnType = myChart
      let count = 66
      let _data = [
        { label: '高危', value: 1 },
        { label: '中危', value: 4 },
        { label: '低危', value: 10 },
        { label: '安全', value: 50 },
        { label: '中危1', value: 4 },
        { label: '低危2', value: 10 },
        { label: '安全3', value: 50 }
      ]
      let dtx = []
      for (var i = 0; i < _data.length; i++) {
        dtx.push(_data[i].label)
      }
       // 添加数据检查
      if (!xdt || !Array.isArray(xdt)) {
        console.warn('xdt 数据无效');
        xdt = [];
      }
      // 绘制图表
      myChart.setOption({
        legend: {
        data: xdy,
        // 可以根据需要调整图例的位置和样式
        top: 30
      },
        grid: {
          left: '10',
          top: '5%',
          bottom: '8',
          right: '26',
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          },
        },
        xAxis: {
          data: xdt?.map((item)=> item.name) || [],
          axisLabel: {
            show: true,
            interval: 0,
            rotate: 30,
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
        },
        yAxis: {
          type: 'value',
          axisTick: {
            show: false,
          },
          axisLine: {
            show: false,
          },
          axisLabel: {
            show: true,
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          splitLine: {
            lineStyle: {
              type: 'dashed'
            }
          }
        },
        // dataZoom: [{
        //     type: 'inside',
        //     start: 0,
        //     end: 10
        // }, {
        //     start: 0,
        //     end: 10,
        //     // handleIcon: 'M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z',
        //     handleSize: '80%',
        //     handleStyle: {
        //         color: '#fff',
        //         shadowBlur: 3,
        //         shadowColor: 'rgba(0, 0, 0, 0.6)',
        //         shadowOffsetX: 2,
        //         shadowOffsetY: 2
        //     }
        // }],
        // dataZoom:[
        //     {
        //         type: "inside", // 内嵌入图标中
        //         startValue: 0, // 用绝对值的方式定义了窗口数据默认展示的数据开始范围
        //         endValue: 10, // 用绝对值的方式定义了窗口数据默认展示的数据终止范围
        //         filterMode: 'filter',
        //         /**
        //         不指定时，当 dataZoom-slider.orient 为 'horizontal'时，
        //         默认控制和 dataZoom 平行的第一个 xAxis。但是不建议使用默认值，建议显式指定。
        //         如果是 number 表示控制一个轴，如果是 Array 表示控制多个轴。
        //         */
        //         xAxisIndex: [0]
        //     }
        // ],
        dataZoom: [
          {
            id: 'dataZoomZ',
            xAxisIndex: [0],
            show: true, //是否显示滑动条，不影响使用
            type: 'slider', // 这个 dataZoom 组件是 slider 型 dataZoom 组件
            startValue: 0, // 从头开始。
            endValue: 9, // 一次性展示5个
            height: 6,
            borderColor: 'transparent',
            fillerColor: 'rgba(205,205,205,1)',
            zoomLock: true,
            showDataShadow: false, //是否显示数据阴影 默认auto
            backgroundColor: 'rgba(0,0,0,0)',
            showDetail: false, //即拖拽时候是否显示详细数值信息 默认true
            realtime: true, //是否实时更新
            filterMode: 'filter',
            handleIcon: 'circle',
            handleStyle: {
              color: 'rgba(205,205,205,1)',
              borderColor: 'rgba(205,205,205,1)',
            },
            handleSize: '80%',
            moveHandleSize: 0,
            maxValueSpan: 9,
            minValueSpan: 9,
            brushSelect: false, //刷选功能，设为false可以防止拖动条长度改变 ************（这是一个坑）
            bottom: 0,
          }, {
            type: 'inside',
            xAxisIndex: 0,
            zoomOnMouseWheel: false,  //滚轮是否触发缩放
            moveOnMouseMove: true,  //鼠标滚轮触发滚动
            moveOnMouseWheel: true
          }],
        series: [
          {
            data: xdy,
            type: 'bar',
            barGap: '10%',
            barCategoryGap: '20%',
            barWidth: '30',
            itemStyle: {
      // 为每个柱子设置不同的颜色
        color: function (params) {
          if (xdt[params.dataIndex].risk == 1) {
            return "#e4485b";
          }else if(xdt[params.dataIndex].risk==2 ){
            
             return '#f8b761'; 
          }
          else if(xdt[params.dataIndex].risk==3){
              return '#4c7ae3';   
          }
          else if(xdt[params.dataIndex].risk>3){
           
             return '#70b603';  
          }
          // params 是一个对象，包含了当前柱子的数据，索引等信息
          
          
        }
    }
          }
        ]
      })
    },
    echartVulnType2 (xdt, xdy) {
      var myChart = echarts.init(document.getElementById('vuln_type_bar2'))
      this.eVulnType = myChart
      let count = 66
      let _data = [
        { label: '高危', value: 1 },
        { label: '中危', value: 4 },
        { label: '低危', value: 10 },
        { label: '安全', value: 50 },
        { label: '中危1', value: 4 },
        { label: '低危2', value: 10 },
        { label: '安全3', value: 50 }
      ]
      let dtx = []
      for (var i = 0; i < _data.length; i++) {
        dtx.push(_data[i].label)
      }
      // 绘制图表
      myChart.setOption({
        legend: {
        data: ['致命', '高危', '中危', '低危'],
        // 可以根据需要调整图例的位置和样式
        top: 10
      },
        grid: {
          left: '10',
          top: '5%',
          bottom: '8',
          right: '26',
          containLabel: true
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          },
        },
        xAxis: {
          data: xdt,
          axisLabel: {
            show: true,
            interval: 0,
            rotate: 30,
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
        },
        yAxis: {
          type: 'value',
          axisTick: {
            show: false,
          },
          axisLine: {
            show: false,
          },
          axisLabel: {
            show: true,
            color: 'rgba(72,72,102,0.64)',
            fontSize: 12,
          },
          splitLine: {
            lineStyle: {
              type: 'dashed'
            }
          }
        },
        // dataZoom: [{
        //     type: 'inside',
        //     start: 0,
        //     end: 10
        // }, {
        //     start: 0,
        //     end: 10,
        //     // handleIcon: 'M10.7,11.9v-1.3H9.3v1.3c-4.9,0.3-8.8,4.4-8.8,9.4c0,5,3.9,9.1,8.8,9.4v1.3h1.3v-1.3c4.9-0.3,8.8-4.4,8.8-9.4C19.5,16.3,15.6,12.2,10.7,11.9z M13.3,24.4H6.7V23h6.6V24.4z M13.3,19.6H6.7v-1.4h6.6V19.6z',
        //     handleSize: '80%',
        //     handleStyle: {
        //         color: '#fff',
        //         shadowBlur: 3,
        //         shadowColor: 'rgba(0, 0, 0, 0.6)',
        //         shadowOffsetX: 2,
        //         shadowOffsetY: 2
        //     }
        // }],
        // dataZoom:[
        //     {
        //         type: "inside", // 内嵌入图标中
        //         startValue: 0, // 用绝对值的方式定义了窗口数据默认展示的数据开始范围
        //         endValue: 10, // 用绝对值的方式定义了窗口数据默认展示的数据终止范围
        //         filterMode: 'filter',
        //         /**
        //         不指定时，当 dataZoom-slider.orient 为 'horizontal'时，
        //         默认控制和 dataZoom 平行的第一个 xAxis。但是不建议使用默认值，建议显式指定。
        //         如果是 number 表示控制一个轴，如果是 Array 表示控制多个轴。
        //         */
        //         xAxisIndex: [0]
        //     }
        // ],
        dataZoom: [
          {
            id: 'dataZoomZ',
            xAxisIndex: [0],
            show: true, //是否显示滑动条，不影响使用
            type: 'slider', // 这个 dataZoom 组件是 slider 型 dataZoom 组件
            startValue: 0, // 从头开始。
            endValue: 9, // 一次性展示5个
            height: 6,
            borderColor: 'transparent',
            fillerColor: 'rgba(205,205,205,1)',
            zoomLock: true,
            showDataShadow: false, //是否显示数据阴影 默认auto
            backgroundColor: 'rgba(0,0,0,0)',
            showDetail: false, //即拖拽时候是否显示详细数值信息 默认true
            realtime: true, //是否实时更新
            filterMode: 'filter',
            handleIcon: 'circle',
            handleStyle: {
              color: 'rgba(205,205,205,1)',
              borderColor: 'rgba(205,205,205,1)',
            },
            handleSize: '80%',
            moveHandleSize: 0,
            maxValueSpan: 9,
            minValueSpan: 9,
            brushSelect: false, //刷选功能，设为false可以防止拖动条长度改变 ************（这是一个坑）
            bottom: 0,
          }, {
            type: 'inside',
            xAxisIndex: 0,
            zoomOnMouseWheel: false,  //滚轮是否触发缩放
            moveOnMouseMove: true,  //鼠标滚轮触发滚动
            moveOnMouseWheel: true
          }],
        series: [
          {
            data: xdy,
            type: 'bar',
            barGap: '10%',
            barCategoryGap: '20%',
            barWidth: '30',
          }
        ]
      })
    },
    echartpie () {
      let totalvuln = this.taskvulriskstat.taskVulRiskCount
      let fatalCount = this.taskvulriskstat.fatalCount
      let highCount = this.taskvulriskstat.highCount
      let mediumCount = this.taskvulriskstat.mediumCount
      let lowCount = this.taskvulriskstat.lowCount

      let totalvulnPie = totalvuln
      if (totalvuln == 0) {
        fatalCount = 1
        highCount = 1
        mediumCount = 1
        lowCount = 1
        totalvulnPie = 4
      }

      var myChart = echarts.init(document.getElementById('vuln_pie'))
      this.epie = myChart
      myChart.setOption({
        color: ['#FF718B', '#FCB5C3', '#FF9F43', '#FFD258'],
        grid: {
          left: '0',
          top: '0',
          bottom: '0',
          right: '0',
          containLabel: true
        },
        title: {
          text: totalvuln,
          textStyle: {
            fontSize: 40,
            fontWeight: 400,
            color: '#1A1A1A',
          },
          subtext: '漏洞总数',
          subtextStyle: {
            fontSize: 13,
            color: 'rgba(72,72,102,0.64)',
          },
          textAlign: "center",
          x: "49%",
          y: "26%"

        },
        series: [
          {
            name: '背景',
            type: 'pie',
            radius: ['72%', '85%'],
            center: ['50%', '43%'],
            startAngle: 180,
            itemStyle: {
              borderRadius: 4,
              // borderColor: '#E9EDF9',
              // borderWidth: 2
            },
            label: {
              show: false,
            },
            data: [
              { value: 1, itemStyle: { color: '#E9EDF9' } },

              {
                // make an record to fill the bottom 50%
                value: 1,
                itemStyle: {
                  // stop the chart from rendering this piece
                  color: 'none',
                  decal: {
                    symbol: 'none'
                  }
                },
                label: {
                  show: false
                }
              }
            ],
            silent: true,
          },
          {
            name: '漏洞统计',
            type: 'pie',
            radius: ['75%', '82%'],
            center: ['50%', '43%'],
            startAngle: 180,
            itemStyle: {
              borderRadius: 4,
              borderColor: 'rgba(255,255,255,0)',
              borderWidth: 2
            },
            label: {
              show: false,
              formatter (param) {
                // correct the percentage
                return param.name + ' (' + param.percent * 2 + '%)'
              }
            },
            data: [
              { value: fatalCount, },
              { value: highCount, },
              { value: mediumCount, },
              { value: lowCount, },
              {
                value: totalvulnPie,
                itemStyle: {
                  // stop the chart from rendering this piece
                  color: 'none',
                  decal: {
                    symbol: 'none'
                  }
                },
                label: {
                  show: false
                }
              }
            ],
            silent: true,
          },

        ]
      })
    },

  }
})

</script>
