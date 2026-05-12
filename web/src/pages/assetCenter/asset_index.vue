<template>
  <div class="home-container" style="overflow-x: hidden;">
    <div class="mt20" style="">
      <el-row :gutter="14">
        <el-col :span="4" style="height: 100%;">
          <div class="item_content_box">
            <div class="box_2">
              <div class="title">
              </div>
              <div>
                <div style="display: flex">
                  <div style="position: relative;flex: 1;">
                    <img style="position: absolute;top:20%;left: 0;" src="/@/assets/_icon/icon_38X38_IP@2x.png" />
                  </div>
                  <div style="flex: 3;">
                    <div class="mr40">
                      <p for="">资产总数</p>
                      <p for="">{{ chartsList.assetTotal }}</p>
                    </div>
                    <div class="mt4">
                      <p for="">同比上周<span style="margin-left: 10px;color: greenyellow;">+{{ chartsList.yoyLastWeekAssetTotal }}</span></p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div style="margin-top: 20px;" class="item_content_box">
            <div class="box_2">
              <div class="title">
              </div>
              <div>
                <div style="display: flex">
                  <div style="position: relative;flex: 1;">
                    <img style="position: absolute;top:20%;left: 0;" src="/@/assets/_icon/icon_38X38_IP@2x.png" />
                  </div>
                  <div style="flex: 3;">
                    <div class="mr40">
                      <p for="">新增IP</p>
                      <p for="">{{ chartsList.newAddIpNum }}</p>
                    </div>
                    <div class="mt4">
                      <p for="">本周新增IP<span style="color: #989898;margin-left: 10px;font-weight: 700;font-size:16px;">{{ chartsList.tkNewAddIpNum }}</span></p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div style="margin-top: 20px;" class="item_content_box">
            <div class="box_2">
              <div class="title">
              </div>
              <div>
                <div style="display: flex">
                  <div style="position: relative;flex: 1;">
                    <img style="position: absolute;top:20%;left: 0;" src="/@/assets/_icon/icon_38X38_IP@2x.png" />
                  </div>
                  <div style="flex: 3;">
                    <div class="mr40">
                      <p for="">减少IP</p>
                      <p for="">{{ chartsList.newReduceIpNum }}</p>
                    </div>
                    <div class="mt4">
                      <p for="">本周减少IP<span style="margin-left: 10px;color: #989898;font-weight: 700;font-size:16px;">  {{ chartsList.tkNewReduceIpNum }}</span></p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-col>
        <el-col :span="20" style="">
          <div class="chart_bg">
            <div class="chart_title mb12 flex">
              <label for="">资产变化趋势</label>
              <div class="btnlist">
                <!-- <el-button size="default" class="active " type="primary">攻击成功</el-button> -->
                <!-- <el-button size="default" type="primary" class="default" >攻击失败</el-button> -->
                <!-- <el-button size="default" type="primary">15天</el-button> -->
                <!-- <el-button size="default" type="primary" class="default">本周</el-button>
                  <el-button size="default" type="primary" class="default">本月</el-button>
                  <el-button size="default" type="primary" class="default">今年</el-button>  -->
                <el-date-picker v-if="false" v-model="timevalue" type="datetimerange" start-placeholder="开始时间" end-placeholder="结束时间"
                  format="YYYY-MM-DD HH:mm:ss" date-format="YYYY/MM/DD ddd" time-format="A hh:mm:ss" size="default"
                  class="ml40" style="width: 340px;flex-grow:0" />
              </div>
            </div>

            <div ref="homeLineRef" class="homeLineRef">

            </div>
          </div>

        </el-col>

      </el-row>



    </div>

    <div class="chart_bg mt20">
      <div class="chart_title mb12 flex">
        <label for="">资产类型统计</label>
        <div class="btnlist">
          <!-- <el-button size="default" type="primary">24h</el-button>
          <el-button size="default" type="primary" class="default">本周</el-button>
          <el-button size="default" type="primary" class="default">本月</el-button>
          <el-button size="default" type="primary" class="default">今年</el-button> -->
        </div>
      </div>
      <div ref='homeBarRef' class="homeBarRef">
      </div>
    </div>


  </div>
</template>
<style  lang="less" scoped>
.item_content_title {
  display: block;
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
  // margin-top: 24px;
  margin-bottom: 16px;
}

.item_content_box {
  // border: 1px solid red;
  background: #fff;
  height: 130px;
  border-radius: 10px;

  >div {
    border-radius: 8px;

    padding: 0 20px;

    img {
      width: 38px;
      height: 38px;
      margin-right: 14px;
    }

    >div {
      p:nth-child(1) {
        font-size: 12px;
        color: rgba(255, 255, 255, 0.65);
        margin-bottom: 8px;
      }

      p:nth-child(2) {
        font-size: 18px;
        color: rgba(255, 255, 255, 0.85);
      }
    }

    &.box_2 {
      width: 100%;
      height: 100%;
      margin-right: 16px;
      display: inline-block;
      margin-bottom: 16px;

      .title {
        font-size: 14px;
        color: #FFFFFF;
        margin-top: 20px;
        margin-bottom: 24px;
      }

      >div {
        // display: flex;
        // // 竖着排列
        // flex-direction: column;
        // // 水平居中
        // justify-content: flex-start;
        // // 垂直居中
        // align-items: center;

      }
    }

    &.box_3 {
      width: 264px;
      height: 124px;
      margin-right: 16px;
      display: inline-block;
      margin-bottom: 16px;
      padding: 18px 16px;
      position: relative;

      label {
        color: rgba(255, 255, 255, 0.85);
        font-size: 14px;

      }

      >div p {
        color: rgba(255, 255, 255, 0.65);
        font-size: 12px;
      }

      .label_s {
        position: absolute;

        width: 48px;
        height: 20px;
        font-size: 12px;
        text-align: center;
        line-height: 20px;
        border-radius: 4px 0 4px 4px;
        top: 20px;
        right: -4px;

        i {
          position: absolute;
          top: -4px;
          right: 0;
          width: 4px;
          height: 4px;
          border-radius: 0px 8px 0px 0px;

        }

        &.yes {
          color: #fff;
          background: #56BD63;

          i {
            background: #56BD63;
          }
        }

        &.no {
          color: #60646C;
          ;
          background: #ABAFB5;

          i {
            background: #ABAFB5;
          }
        }
      }

    }

    &.disabled {
      opacity: 0.465;
    }

    &.active {
      border: 1px solid rgba(2, 109, 255, 0.6);
      background: rgba(2, 109, 255, 0.2);
    }
  }
}

.tools-list {
  background-color: var(--el-bg-color);
  padding: 14px;
  border-radius: 12px;
}

.chart_bg {
  padding: 22px 24px;
  box-sizing: border-box;
  background: #fff;
  border-radius: 12px;

  .chart_title {
    align-items: center;

    label {
      font-size: 16px;
      color: #FFFFFF;
      display: inline-flex;
    }

    .btnlist {
      display: flex;
      align-items: center;
      justify-content: flex-end;
      flex: 1 1 0%;

      >button {
        font-size: 12px;

        &.default {
          background: rgba(255, 255, 255, 0.15);
          color: rgba(255, 255, 255, 0.65);
          border: 0
        }
      }

      .el-button--primary.default {}
    }
  }
}

.homeLineRef,
.homeBarRef {
  height: 350px;
  width: 100%;
}
</style>
<script>
import { traffic } from '@/api/asset.js'; 
var echarts = require('echarts');

export default({
    name:'assetview',  
    data(){  
    	return{  
            global:{
                homeChartOne: null,
                homeChartTwo: null,
                homeCharThree: null,
                dispose: [null, '', undefined],
            },
            timevalue:'',
            state:{
                homeThree: [
                    {
                    icon: 'iconfont icon-yangan',
                    label: '浅粉红',
                    value: '2.1%OBS/M',
                    iconColor: '#F72B3F',
                    },
                    {
                    icon: 'iconfont icon-wendu',
                    label: '深红(猩红)',
                    value: '30℃',
                    iconColor: '#91BFF8',
                    },
                    {
                    icon: 'iconfont icon-shidu',
                    label: '淡紫红',
                    value: '57%RH',
                    iconColor: '#88D565',
                    },
                    {
                    icon: 'iconfont icon-shidu',
                    label: '弱紫罗兰红',
                    value: '107w',
                    iconColor: '#88D565',
                    },
                    {
                    icon: 'iconfont icon-zaosheng',
                    label: '中紫罗兰红',
                    value: '57DB',
                    iconColor: '#FBD4A0',
                    },
                    {
                    icon: 'iconfont icon-zaosheng',
                    label: '紫罗兰',
                    value: '57PV',
                    iconColor: '#FBD4A0',
                    },
                    {
                    icon: 'iconfont icon-zaosheng',
                    label: '暗紫罗兰',
                    value: '517Cpd',
                    iconColor: '#FBD4A0',
                    },
                    {
                    icon: 'iconfont icon-zaosheng',
                    label: '幽灵白',
                    value: '12kg',
                    iconColor: '#FBD4A0',
                    },
                    {
                    icon: 'iconfont icon-zaosheng',
                    label: '海军蓝',
                    value: '64fm',
                    iconColor: '#FBD4A0',
                    },
                ],
                myCharts: [],
                charts: {
                    theme: '',
                    bgColor: '',
                    color: '#303133',
                },
                },
            chartsList:[],
            zhuKey:[],
            zhuvalue:[],
        }
    },
    mounted:function(){

    },
    methods:{
         // 折线图
        initLineChart(){
            if (!global.dispose.some((b) => b === global.homeChartOne)) global.homeChartOne.dispose();
                global.homeChartOne = echarts.init(homeLineRef.value, state.charts.theme);
                const option = {
                    backgroundColor: state.charts.bgColor,
                    title: {
                    text: '',
                    x: 'left',
                    textStyle: { fontSize: '15', color: state.charts.color },
                    },
                    grid: { top: 30, right: 20, bottom: 30, left: 30 },
                    tooltip: { trigger: 'axis' },
                    legend: {
                    data: ['新增IP', '减少IP'], left: 0
                    },
                    xAxis: {
                    data: chartsList.value?.assetReduceTrendChangeRes.name||[],
                    },
                    yAxis: [
                    {
                        type: 'value',
                        name: '',
                        splitLine: { show: true, lineStyle: { type: 'dashed', color: 'rgba(255,255,255,0.15)' } },
                    },
                    ],
                    series: [
                    {
                        name: '减少IP',
                        type: 'line',
                        symbolSize: 6,
                        symbol: 'circle',
                        smooth: true,
                        data: chartsList.value.assetReduceTrendChangeRes.value||[],
                        lineStyle: { color: '#F9242F' },
                        itemStyle: { color: '#F9242F', borderColor: '#fe9a8b' },
                        areaStyle: {
                        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                            { offset: 0, color: 'rgba(249,36,47,0.45)' },
                            { offset: 1, color: 'rgba(249,36,47,0)' },
                        ]),
                        },
                    },

                    {
                        name: '新增IP',
                        type: 'line',
                        symbolSize: 6,
                        symbol: 'circle',
                        smooth: true,
                        data:chartsList.value.assetTrendChangeRes.value||[],
                        lineStyle: { color: '#026DFF' },
                        itemStyle: { color: '#026DFF', borderColor: '#9E87FF' },
                        areaStyle: {
                        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                            { offset: 0, color: ' rgba(2,109,255,0.55)' },
                            { offset: 1, color: 'rgba(2,109,255,0)' },
                        ]),
                        },
                    }
                    ],
                };
                global.homeChartOne.setOption(option);
                state.myCharts.push(global.homeChartOne);
        },
         // 柱状图
        initBarChart(){
             if (!global.dispose.some((b) => b === global.homeCharThree)) global.homeCharThree.dispose();
                global.homeCharThree = echarts.init(homeBarRef.value, state.charts.theme);
                const option = {
                    backgroundColor: state.charts.bgColor,

                    tooltip: { trigger: 'axis' },
                    legend: { data: ['其他', '办公设备'], left: 0 },
                    grid: { top: 40, right: 20, bottom: 20, left: 40 },
                    xAxis: [
                    {
                        type: 'category',
                        data: zhuKey.value||[],
                        // data: ['防火墙', '账务系统waf', '办公防火墙', '数据库审计', '5km', '6km'],
                        boundaryGap: true,
                        axisTick: { show: false },
                    },
                    ],
                    yAxis: [
                    {
                        type: 'value',
                        splitLine: { show: true, lineStyle: { type: 'dashed', color: 'rgba(255,255,255,0.15)' } },
                    },
                    ],
                    series: [
                    {
                        name: '类型数量',
                        type: 'bar',
                        barWidth: 10,
                        stack: 'total',
                        itemStyle: {
                        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                            { offset: 0, color: '#1369DE' },
                            { offset: 1, color: '#2A84FF' },
                        ]),
                        //柱状图圆角
                        borderRadius: [30, 30, 0, 0],
                        },
                        // data: [
                        //   { value: 11, stationName: 's1' },
                        //   { value: 34, stationName: 's2' },
                        //   { value: 54, stationName: 's3' },
                        // ],
                        data:zhuvalue.value,
                    },
                    // {
                    //   name: '错误告警',
                    //   type: 'bar',
                    //   stack: 'total',
                    //   barWidth: 10,
                    //   itemStyle: {
                    //     color: '#363E4E',
                    //     //柱状图圆角
                    //     borderRadius: [30, 30, 0, 0],
                    //   },
                    //   data: [
                    //     { value: 111, stationName: 's1' },
                    //     { value: 34, stationName: 's2' },
                    //     { value: 54, stationName: 's3' },
                    //   ],
                    // },

                    ],
                };
                global.homeCharThree.setOption(option);
                state.myCharts.push(global.homeCharThree);
        },
        // 请求数据
        async getData(){
            try {
                const res = await traffic.statisticsVal()
                chartsList.value = res.data
                zhuKey.value = chartsList.value.assetTypeStatics.map(item => item.AssetType)
                zhuvalue.value = chartsList.value.assetTypeStatics.map(item => {
                return {value:item.Count}
                })
            } catch (error) {
                console.log(error);

            }
        },
    }
})

 

 
//   setup() {  
//     const homeLineRef = ref();
//     const homeBarRef = ref();
//     const chartsList = ref([]);
//     const zhuKey = ref([]);
//     const zhuvalue = ref([]);
 
   
 
// getData()
//     // 监听 vuex 中是否开启深色主题
//     watch(
//       () => true,
//       (isIsDark) => {
//         nextTick(() => {
//           state.charts.theme = isIsDark ? 'dark' : '';
//           state.charts.bgColor = isIsDark ? 'transparent' : '';
//           state.charts.color = isIsDark ? '#dadada' : '#303133';
//           setTimeout(() => {
//             initLineChart();
//           }, 500);
//           // setTimeout(() => {
//           // 	initPieChart();
//           // }, 700);
//           setTimeout(() => {
//             initBarChart();
//           }, 1000);
//         });
//       },
//       {
//         deep: true,
//         immediate: true,
//       }
//     );

//     return {
//       timevalue,
//       homeLineRef,
//       homeBarRef,
//       initLineChart,
//       initBarChart,
//       chartsList
//     }
//   }
// };
</script>
