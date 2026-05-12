<template>
  <div class="overview-module">
    <!-- 信息卡片区 -->
      <h3 style="margin:10px 0;color:#484866A3;font-size:16px" >资产组信息</h3>

    <div class="card-section">
      <!-- 这里可以根据实际需要循环渲染卡片 -->
      <div class="info-card" >
      
        <div class="card-content">
          <h3>资产组</h3>
          <p>{{allDAta.assetGroupInfo?.name||'无'}}</p>
        </div>
        <div class="card-content" style="margin-left:20px">
          <h3>资产组风险</h3>
          <p>{{allDAta.assetGroupInfo.riskLevelName||'无'}}</p>
        </div>
        <div class="card-content" style="margin-left:20px">
          <h3>子资产组数量</h3>
         <p>{{allDAta.assetGroupInfo.subAssetCount||'无'}}</p>
        </div>
        <div class="card-content" style="margin-left:20px">
          <h3>资产总数</h3>
          <p>{{allDAta.assetGroupInfo.assetCount||'无'}}</p>
        </div>
        <div class="card-content" style="margin-left:20px">
          <h3>资产组路径</h3>
          <p>{{allDAta.assetGroupInfo.assetGroupPath||'无'}}</p>
        </div>
        <div class="card-content" style="margin-left:20px">
          <h3>资产组说明</h3>
         <p>{{allDAta.assetGroupInfo.assetGroupRemark||'无'}}</p>
        </div>
       
      </div>
    </div>
      <h3 style="margin:20px 0 ;color:#484866A3;font-size:16px">资产统计</h3>

    <div class="card-section">
      <!-- 这里可以根据实际需要循环渲染卡片 -->
      <div class="info-card" >
        <div class="card-content">
          <img style="width:30px" src="../../../../assets/icon/assetManage/gaowei.png" alt="" ></img>
          <h4>高危资产</h4>
          <h3 style="color:#ff687a">{{allDAta.assetStatics.highRiskAsset}}</h3>
        </div>
        <div class="card-content" style="margin-left:20px">
          <img style="width:30px" src="../../../../assets/icon/assetManage/zhongwei.png" alt="" ></img>

          <h3>中危资产</h3>
          <h3 style="color:#ffbb44">{{allDAta.assetStatics.middleRiskAsset}}</h3>
        </div>
        <div class="card-content" style="margin-left:20px">
          <img style="width:30px" src="../../../../assets/icon/assetManage/diwei.png" alt="" ></img>

          <h4>低危资产</h4>
          <h3 style="color:#33a2ef">{{allDAta.assetStatics.lowRiskAsset}}</h3>
        </div>
        <div class="card-content" style="margin-left:20px">
          <img style="width:30px" src="../../../../assets/icon/assetManage/anquan.png" alt="" ></img>

          <h3>安全资产</h3>
         <h3 style="color:#16c4af">{{allDAta.assetStatics.safeAsset}}</h3>
        </div>
   
       
      </div>
      <div class="info-card" >
      
       
        <div class="card-content" style="margin-left:20px">
          <img style="width:30px" src="../../../../assets/icon/assetManage/0.png" alt="" ></img>

          <h3>致命漏洞</h3>
          <h3 style="color:#d4237a">{{allDAta.vulStatics.deadlyVul}}</h3>
        </div>
        <div class="card-content" style="margin-left:20px">
          <img style="width:30px" src="../../../../assets/icon/assetManage/1.png" alt="" ></img>

          <h3>高危漏洞</h3>
          <h3 style="color:#ff7b20">{{allDAta.vulStatics.highRiskVul}}</h3>
        </div>
        <div class="card-content" style="margin-left:20px">
          <img style="width:30px" src="../../../../assets/icon/assetManage/2.png" alt="" ></img>

          <h3>中危漏洞</h3>
          <h3 style="color:#1b69fd">{{allDAta.vulStatics.mediumRiskVul}}</h3>
        </div>
        <div class="card-content" style="margin-left:20px">
          <img style="width:30px" src="../../../../assets/icon/assetManage/3.png" alt="" ></img>

          <h3>低危漏洞</h3>
          <h3 style="color:#5bb878">{{allDAta.vulStatics.lowRiskVul}}</h3>
        </div>
       
      </div>
    </div>
    
    <!-- 图表区 -->
   <div  style="height: 270px;margin-top:20px" class="box">
      <div v-if="allDAta.vulTypeStatics.vulType" class="box_title" style="position: relative;">
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
  </div>
</template>

<script>
var echarts = require('echarts')

import { traffic } from '@/api/assetManagement.js'

export default {
 name: 'gailan',
  data() {
    return {
 allDAta:{},
 eVulnType:null,
      // 图表配置
      chartOptions: {
        // 这里将是你的ECharts图表配置
      }
    };
  },
  methods: {
      echartVulnType (xdt, xdy) {
        try {
          if(!xdt) return
          this.$nextTick(()=>{
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
            // 绘制图表
            myChart.setOption({
              legend: {
              data: xdy || [],
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
                data: xdt?.map((item)=> item?.name || "") || [],
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
                  data: xdy || [],
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
          
          })
           
        } catch (error) {
          console.log(error,'290000000');
        }
    },
    async getvultypestat2 () { //fuwu类型统计 
    if(this.$store.state.groupID[1] == 2) return  // 左侧的树分为两个类型，这个是组类型左边为1 如果是2就是IP的详情 
      const res = await traffic.trafficVuln({
        groupID: this.$store.state.groupID[0]||-1
      })
      if (res.code == 200) {
        this.allDAta = res.data
        console.log(this.allDAta,'this.allDAta');
        let xdt = res.data.vulTypeStatics.vulType
        let xdy = res.data.vulTypeStatics.count
        this.echartVulnType(xdt, xdy)
      }
    },
  },
  mounted () {
    this.getvultypestat2();
  },
   watch: {
    // 监听 Vuex store 中的 groupID
    '$store.state.groupID'(newVal, oldVal) {
      // 当 groupID 发生变化时，执行一些操作
      // newVal 是新的 groupID 值，oldVal 是旧的值
     if(newVal){
      this.$nextTick(()=>{
        this.getvultypestat2();
      })
     }
    }
  },
};
</script>


<style scoped>
/* 基础样式调整 */
.overview-module {
  padding: 20px;
  background: #FFFFFF;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.section-title {
  margin: 20px 0;
  color: #484866A3; /* 更新为指定的颜色 */
  font-size: 16px; /* 字体大小调整 */
  font-weight: bold;
}

.card-section {
  display: flex;
  flex-wrap: wrap;
  gap: 20px; /* 控制卡片间的间距 */
}

.info-card {
  flex: 1;
  /* background: linear-gradient(145deg, #e6e9f0, #ffffff); */
  background: linear-gradient(135deg, #f7f9fc 0%, #f7f9fc 100%);

  padding: 15px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-around; /* 优化内部元素的分布 */
}

.card-content {
  margin: 10px 0;
  text-align: center;
}

.card-content h3, .card-content h4 {
  margin-bottom: 5px;
  color: #484866A3; /* 标题颜色调整为指定颜色 */
  font-size: 14px; /* 标题字体大小调整 */
}

.card-content p {
  font-size: 12px; /* 内容字体大小调整 */
  color: #666;
}

.chart-section {
  margin-top: 30px;
}

.box_title {
  margin-bottom: 20px;
}

.chart {
  height: 270px;
  background: #f0f2f5;
  border-radius: 8px;
  box-shadow: inset 0 4px 12px rgba(0,0,0,0.05);
}

/* 漏洞类型统计图表的自定义样式 */
.box_title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f0f2f5; /* 调整为更柔和的背景色 */
  color: #484866A3; /* 标题颜色调整 */
  border-radius: 8px 8px 0 0;
}

.time_group {
  display: flex;
  /* 添加必要的样式 */
}
</style>
