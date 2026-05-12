<template>
  <div class="gailan-container">
    <div class="main-title  ">  
      <i class="nav_icon"></i>
      <i style="margin: 0 8px; vertical-align: text-top; color: #BCC4D3;">|</i> 
      <label class="taskSceneBtn" >资产概览</label>
        
    </div> 
    <!-- 顶部区域 -->
    <div class="top-section">
      <!-- 左侧卡片区域 -->
      <div class="left-cards">
        <div class="card-item">
          <div class="card-icon">
            <el-icon :size="38" color="#E6A23C">
              <TrendCharts />
            </el-icon>
          </div>
          <div class="card-content">
            <div class="card-title">资产总数</div>
            <div class="card-value">{{ (gailan.assetRiskStat?.reduce((sum, item) => sum + item.count, 0)) || 0 }}</div>
          </div>
        </div>
                <div class="card-item">
          <div class="card-icon">
            <el-icon :size="38" color="#409EFF">
              <DataBoard />
            </el-icon>
          </div>
          <div class="card-content">
            <div class="card-title">风险资产数</div>
            <div class="card-value">{{ gailan.safeLoopholeTotal|| 0 }}</div>

          </div>
        </div>
<!-- 
        <div class="card-item">
          <div class="card-icon">
            <el-icon :size="38" color="#67C23A">
              <TrendCharts />
            </el-icon>
          </div>
          <div class="card-content">
            <div class="card-title">发现配置问题</div>
            <div class="card-value">{{ gailan.configErrTotal || 0 }}</div>

          </div>
        </div> -->



      </div>

      <!-- 右侧图表区域 -->
      <div class="right-charts">
        <div class="chart-container">
          <div class="chart-title">资产风险统计</div>
          <div ref="riskStatChart" class="chart-box" id="riskStatChart"></div>
        </div>

        <div class="chart-container">
          <div class="chart-title">资产风险趋势</div>
         <div ref="riskTrendChart" class="chart-box" id="riskTrendChart"></div>
        </div>
      </div>
    </div>

    <!-- 中间区域 -->
    <!-- <div class="middle-section">
      <div class="chart-container full-width">
        <div class="chart-title">资产风险分布趋势</div>
        <div ref="riskTrendChart" class="chart-box large" id="riskTrendChart"></div>
      </div>

      <div class="chart-container full-width">
        <div class="chart-title">资产风险趋势</div>
        <div ref="assetTrendChart" class="chart-box large" id="assetTrendChart"></div>
      </div>
    </div> -->

    <!-- 底部表格区域 -->
    <div class="bottom-section">
      <div class="table-title">近期危险资产</div>
            <el-table
              :data="tableData"
              style="width: 100%"
              height="calc(100% - 26px)"
              :header-cell-style="{ color: '#333333' }"
              :row-style="{ backgroundColor: '#FFFFFF', color: '#333333' }"
            >
              <el-table-column prop="ip" label="IP" min-width="140" />
              <el-table-column prop="assetType" label="资产类型" min-width="120" />
              <el-table-column prop="system" label="操作系统" min-width="120" />
              <el-table-column prop="ports" label="开放端口" min-width="180" />
              <el-table-column prop="riskLevel" label="风险等级" min-width="100">
                <template #default="scope">
                  <span :class="getRiskClass(scope.row.riskLevel)">{{ scope.row.riskLevel }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="tags" label="漏洞" min-width="80" >
                <template #default="scope">
                <div style="display: flex;justify-content: flex-start;align-items: center;gap: 10px;">
                  <span style="color: red;" >{{ scope.row.tags.deadlyVul || 0 }}</span>
                  <span style="color: orange;" >{{ scope.row.tags.highRiskVul || 0 }}</span>
                  <span style="color: blue;" >{{ scope.row.tags.mediumRiskVul || 0 }}</span>
                  <span style="color: #16c60c;" >{{ scope.row.tags.lowRiskVul || 0 }}</span>
                </div>
                </template>
              </el-table-column>
              <el-table-column prop="updateTime" label="更新时间" min-width="160" />
            </el-table>
    </div>
  </div>
</template>
 
<script> 
import asset from '@/api/asset.js'; 
var echarts = require('echarts');

export default({
    name:'assetview',  
    data(){  
    	return{  
            gailan:{
                assetRiskStat:[],
                assetRiskLevelTrendRes:[],
                assetReduceTrendChangeRes:[],
                assetRiskTypeTrendRes:[], 
                assetTrendChangeRe:[],
            }, 
            tableData:[],

        }
    },
     created:function(){ 
        this.$store.state.activefirstMenu="/assetview";  
    },
    mounted:function(){
        this.getGailan();
    },
    methods:{
        async getGailan(){
          try {
            const res = await asset.statisticsVal();
            this.gailan = res.data || {};

            // 处理表格数据，添加兜底处理
            if (this.gailan.recentDangerAsset && Array.isArray(this.gailan.recentDangerAsset)) {
              this.tableData = this.gailan.recentDangerAsset.map(item => ({
                ip: item.ip || '未知',
                assetType: item.assetTypeName || '未知',
                system: item.os || '未知',
                ports: item.openPort || '未知',
                riskLevel: item.riskLevel || '未知',
                tags: item.vulStatics || { deadlyVul: 0, highRiskVul: 0, mediumRiskVul: 0, lowRiskVul: 0 },
                updateTime: item.time || '暂无时间'
              }));
            } else {
              // 表格数据兜底
              this.tableData = [];
            }

            // 为图表数据添加兜底处理
            if (!this.gailan.assetRiskStat || !Array.isArray(this.gailan.assetRiskStat)) {
              this.gailan.assetRiskStat = [{ assetType: '暂无数据', count: 0 }];
            }

            if (!this.gailan.riskStatics || !Array.isArray(this.gailan.riskStatics)) {
              this.gailan.value.riskStatics = [{ assetType: '暂无数据', count: 0 }];
            }

            if (!this.gailan.assetRiskLevelTrendRes || !Array.isArray(this.gailan.assetRiskLevelTrendRes)) {
              this.gailan.assetRiskLevelTrendRes = [{
                date: '暂无数据',
                high: 0,
                medium: 0,
                low: 0,
                safe: 0
              }];
            }

            if (!this.gailan.assetRiskTypeTrendRes || !Array.isArray(this.gailan.assetRiskTypeTrendRes)) {
              this.gailan.assetRiskTypeTrendRes = [{
                date: '暂无数据',
                unSafe: 0,
                safe: 0
              }];
            }

            // 数据加载完成后初始化图表
            // await nextTick();
            this.initCharts();
          } catch (error) {
            console.log(error);
            // 请求失败时设置兜底数据
            this.setDefaultGailanData();
          }
        },
        initCharts(){
            this.initRiskStatChart();
            // this.initRiskTypeChart();
            this.initRiskTrendChart();
            // this.initAssetTrendChart();
        },
        async setDefaultGailanData(){
          this.gailan = {
            assetRiskStat: [{ assetType: '暂无数据', count: 0 }],
            riskStatics: [{ assetType: '暂无数据', count: 0 }],
            assetRiskLevelTrendRes: [{ date: '暂无数据', high: 0, medium: 0, low: 0, safe: 0 }],
            assetRiskTypeTrendRes: [{ date: '暂无数据', unSafe: 0, safe: 0 }],
            safeLoopholeTotal: 0,
            configErrTotal: 0,
            recentDangerAsset: []
          };

          this.tableData = [];

          // 初始化图表
        //   await nextTick();
          this.initCharts();
        },
        // 资产风险统计柱状图
        initRiskStatChart(){
           if (!this.gailan.assetRiskStat || 
                !Array.isArray(this.gailan.assetRiskStat) || 
              this.gailan.assetRiskStat.length === 0) {
                // 如果没有数据，设置默认数据
                gailan.value.assetRiskStat = [{ assetType: '暂无数据', count: 0 }];
            }
            const chart = echarts.init(document.getElementById('riskStatChart'));

            // let data = [
            //     {
            //         value: dt.data.free,
            //         percent: dt.data.freePercent,
            //         name: '可用空间'
            //     },
            //     {
            //         value: dt.data.used,
            //         percent: dt.data.usedPercent,
            //         name: '已用空间'
            //     },
            // ]
            

            // 处理数据
            const chartData = this.gailan.assetRiskStat.map(item => ({
                name: item.assetType || '未知',
                value: item.count || 0
            }));

            // 计算背景条的最大值，确保至少为1
            const maxValue = Math.max(Math.max(...chartData.map(item => item.value)) * 1.2, 1);
            const backgroundData = new Array(chartData.length).fill(maxValue);
            // 定义多种颜色
            const colors = ['#fa5949','#FFD93D',  '#409EFF', '#67C23A',];
            const option = {
                color: colors,
                grid: { top: 10, right: 10, bottom: 30, left: 30 },
                tooltip: { trigger: 'axis' },
                  legend: {
                    orient: 'vertical',
                    // left: 'right',
                    x:'60%',//水平安放位置，默认为'left'，可选为：'center' | 'left' | 'right' | {number}（x坐标，单位px）
                    y: 'center',//垂直安放位置，默认为top，可选为：'top' | 'bottom' | 'center' | {number}（y坐标，单位px）
                    icon: "circle",   //  这个字段控制形状  类型包括 circle 圆形，triangle 三角形，diamond 四边形，arrow 变异三角形，none 无
                    itemWidth: 10,  // 设置宽度
                    itemHeight: 10, // 设置高度
                    itemGap: 30 ,// 设置间距，
                    formatter: function (name) {
                        // console.log(data, 'data')
                        let total = 0
                        let tarValue
                        for (let i = 0; i < chartData.length; i++) {
                            total += chartData[i].value
                            if (chartData[i].name == name) {
                            tarValue = chartData[i].value
                            }
                        }
                        let v = tarValue 
                        //计算出百分比
                        let p = Math.round((tarValue / total) * 100) + '%'  
                        return `${name}  ${p}   ${v}`
                        //name是名称，v是数值
                    } 
                }, 
                series: [
                  {
                    name: '资产风险',
                    type: 'pie',
                    radius: ['40%', '70%'],
                    center:['25%','50%'], 
                    avoidLabelOverlap: false,
                    label: {//去除饼图的指示折线label
                        normal: {
                          show: false,
                          position: 'inside',
                          formatter:"{b}:{d}%"
                        },
                    }, 
                    labelLine: {
                      show: false
                    },
                    data: chartData
                  }
                ]
            };
            chart.setOption(option);
            window.addEventListener('resize', function () {
                chart.resize()
            })
        },
        // 风险类型柱状图
        initRiskTypeChart(){
            if (!this.gailan.riskStatics || !Array.isArray(this.gailan.riskStatics) || this.gailan.riskStatics.length === 0) {
                // 如果没有数据，设置默认数据
                this.gailan.riskStatics = [{ assetType: '暂无数据', count: 0 }];
            }
                        
            const chart = echarts.init(document.getElementById('riskTypeChart'));
            const data = this.gailan.riskStatics.map(item => ({
                name: item.assetType || '未知',
                value: item.count || 0
            }));

            // 计算背景条的最大值，确保至少为1
            const maxValue = Math.max(Math.max(...data.map(item => item.value)) * 1.2, 1);
            const backgroundData = new Array(data.length).fill(maxValue);

            const option = {
                color: ['#00b485'],
                grid: { top: 10, right: 10, bottom: 80, left: 50 },
                tooltip: { trigger: 'axis' },
                xAxis: {
                type: 'category',
                data: data.map(item => item.name),
                boundaryGap: true,
                axisTick: { show: false },
                axisLabel: {
                    rotate: 45,
                    interval: 0,
                    color: '#ADB2B9',
                },
                axisLine: {
                    lineStyle: {
                    color: '#E8EAED',
                    }
                }
                },
                yAxis: {
                type: 'value',
                splitLine: { show: true, lineStyle: { type: 'solid', color: '#E8EAED' } },
                axisLabel: {
                    color: '#ADB2B9',
                },
                },
                series: [
                // 背景轨道
                {
                    name: 'background',
                    type: 'bar',
                    barWidth: 13,
                    barGap: '-100%',
                    itemStyle: {
                    color: '#F8F9FA',
                    borderRadius: [6, 6, 6, 6],
                    },
                    data: backgroundData,
                    silent: true,
                    z: 1,
                    tooltip: {
                    show: false
                    }
                },
                {
                    name: '风险数量',
                    type: 'bar',
                    barWidth: 13,
                    itemStyle: {
                    color: '#00b485',
                    borderRadius: [6, 6, 0, 0],
                    },
                    data: data.map(item => item.value)
                }
                ]
            };
            chart.setOption(option);
            window.addEventListener('resize', function () {
                chart.resize()
            })
        },
        // 资产风险分布趋势折线图
        initRiskTrendChart(){
            if (!this.gailan.assetRiskLevelTrendRes || !Array.isArray(this.gailan.assetRiskLevelTrendRes) || this.gailan.assetRiskLevelTrendRes.length === 0) {
                // 如果没有数据，设置默认数据
                this.gailan.assetRiskLevelTrendRes = [{ date: '暂无数据', high: 0, medium: 0, low: 0, safe: 0 }];
            }
            
            const chart = echarts.init(document.getElementById('riskTrendChart'));

            // 提取日期和数据，添加兜底处理
            const dates = this.gailan.assetRiskLevelTrendRes.map(item => item.date || '暂无数据');
            const highData = this.gailan.assetRiskLevelTrendRes.map(item => item.highAsset || 0);
            const mediumData = this.gailan.assetRiskLevelTrendRes.map(item => item.mediumAsset || 0);
            const lowData = this.gailan.assetRiskLevelTrendRes.map(item => item.lowAsset || 0);
            const safeData = this.gailan.assetRiskLevelTrendRes.map(item => item.safeAsset || 0);

            const option = {
                tooltip: {
                trigger: 'axis',
                formatter: function(params) {
                    let result = params[0].name + '<br/>';
                    params.forEach(param => {
                    result += param.marker + param.seriesName + ': ' + param.value + '<br/>';
                    });
                    return result;
                }
                },
                grid: {
                left: '10',
                top: '40',
                bottom: '10',
                right: '5%',
                containLabel: true
                },
                legend: {
                data: ['高危', '中危', '低危', '安全'],
                textStyle: { color: '#ADB2B9' },
                top: 0,
                right: 0
                },
                xAxis: {
                type: 'category',
                boundaryGap: false,
                data: dates,
                axisLabel: {
                    color: '#ADB2B9',
                },
                axisLine: {
                    lineStyle: {
                    color: '#E8EAED',
                    }
                }
                },
                yAxis: {
                type: 'value',
                min: 0,
                splitLine: { show: true, lineStyle: { type: 'solid', color: '#E8EAED' } },
                axisLabel: {
                    color: '#ADB2B9',
                },
                },
                series: [
                {
                    name: '高危',
                    data: highData,
                    type: 'line',
                    smooth: true,
                    lineStyle: { color: '#fa5949' },
                    itemStyle: { color: '#fa5949' },
                    areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(250,89,73,0.3)' },
                        { offset: 1, color: 'rgba(250,89,73,0)' },
                    ]),
                    }
                },
                {
                    name: '中危',
                    data: mediumData,
                    type: 'line',
                    smooth: true,
                    lineStyle: { color: '#FFD93D' },//黄色
                    itemStyle: { color: '#FFD93D' },
                    areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(255,217,61,0.3)' },
                        { offset: 1, color: 'rgba(255,217,61,0)' },
                    ]),
                    }
                },
                {
                    name: '低危',
                    data: lowData,
                    type: 'line',
                    smooth: true,

                },
                {
                    name: '安全',
                    data: safeData,
                    type: 'line',
                    smooth: true,
                    lineStyle: { color: '#52E39D' }, //绿色
                    itemStyle: { color: '#52E39D' },
                    areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(82,227,157,0.3)' },
                        { offset: 1, color: 'rgba(82,227,157,0)' },
                    ]),
                    }
                }
                ]
            };
            chart.setOption(option);
            window.addEventListener('resize', function () {
                chart.resize()
            })
        },
        // 资产风险趋势折线图
        initAssetTrendChart(){
            if (!this.gailan.assetRiskTypeTrendRes || !Array.isArray(this.gailan.assetRiskTypeTrendRes) || this.gailan.assetRiskTypeTrendRes.length === 0) {
                // 如果没有数据，设置默认数据
                this.gailan.assetRiskTypeTrendRes = [{ date: '暂无数据', unSafe: 0, safe: 0 }];
            } 
            const chart = echarts.init(document.getElementById('assetTrendChart'));

            // 提取日期和数据，添加兜底处理
            const dates = this.gailan.assetRiskTypeTrendRes.map(item => item.date || '暂无数据');
            const unSafeData = this.gailan.assetRiskTypeTrendRes.map(item => item.unSafe || 0);
            const safeData = this.gailan.assetRiskTypeTrendRes.map(item => item.safe || 0);

            const option = {
                tooltip: {
                trigger: 'axis',
                formatter: function(params) {
                    let result = params[0].name + '<br/>';
                    params.forEach(param => {
                    result += param.marker + param.seriesName + ': ' + param.value + '台<br/>';
                    });
                    return result;
                }
                },
                grid: {
                left: '10',
                top: '40',
                bottom: '10',
                right: '5%',
                containLabel: true
                },
                legend: {
                data: ['风险主机', '安全主机'],
                textStyle: { color: '#ADB2B9' },
                    top: 0,
                    right: 0
                },
                xAxis: {
                type: 'category',
                boundaryGap: false,
                data: dates,
                axisLabel: {
                    color: '#ADB2B9',
                },
                axisLine: {
                    lineStyle: {
                    color: '#E8EAED',
                    }
                }
                },
                yAxis: {
                type: 'value',
                min: 0,
                splitLine: { show: true, lineStyle: { type: 'solid', color: '#E8EAED' } },
                axisLabel: {
                    color: '#ADB2B9',
                },
                },
                series: [
                {
                    name: '风险主机',
                    data: unSafeData,
                    type: 'line',
                    smooth: true,
                    lineStyle: { color: '#E35282' },
                    itemStyle: { color: '#E35282' },
                    areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(227,82,130,0.3)' },
                        { offset: 1, color: 'rgba(227,82,130,0)' },
                    ]),
                    }
                },
                {
                    name: '安全主机',
                    data: safeData,
                    type: 'line',
                    smooth: true,
                    lineStyle: { color: '#52E39D' },
                    itemStyle: { color: '#52E39D' },
                    areaStyle: {
                    color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
                        { offset: 0, color: 'rgba(82,227,157,0.3)' },
                        { offset: 1, color: 'rgba(82,227,157,0)' },
                    ]),
                    }

                }
                ]
            };
            chart.setOption(option);
            window.addEventListener('resize', function () {
                chart.resize()
            })
        },  
        // 获取风险等级样式
        getRiskClass(level){
            const classes = {
                '高危': 'risk-high',
                '高危资产': 'risk-high',
                '中危': 'risk-medium',
                '中危资产': 'risk-medium',
                '低危': 'risk-low',
                '低危资产': 'risk-low',
                '安全': 'risk-safe',
                '安全资产': 'risk-safe'
            };
            return classes[level] || '';
        },

    }
}) 
 
</script>

<style scoped lang="less">
  .gailan-container {
    // padding: 20px;
    box-sizing: border-box;
    background-color: transparent;
    // min-height: 100vh;
    color: #333333;
    height: calc(100% - 24px);
  }

  .top-section {
    display: flex;
    gap: 16px;
    margin-bottom: 16px;
  }

  .left-cards {
    display: flex;
    flex-direction: column;
    gap: 16px;
    width: 300px;
    height: 300px; // 与右侧图表高度一致
    justify-content: space-between;
  }

  .card-item {
    display: flex;
    align-items: center;
    background: #fff;
    // border: 1px solid #E9ECEF;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    border-radius: 4px; 
    padding: 20px;
    height: 130px;
    flex: 1;

    .card-icon {
      margin-right: 16px;
      display: flex;
      align-items: center;
      justify-content: center;

      .el-icon {
        font-size: 38px;
      }
    }

    .card-content {
      flex: 1;

          .card-title {
        font-size: 16px;
        color: #666666;
        margin-bottom: 20px;
      }

      .card-value {
        font-size: 28px;
        font-weight: bold;
        color: #333333;
        margin-bottom: 8px;
      }

      .card-trend {
        font-size: 12px;
        color: #666666;

        .trend-up {
          color: #52c41a;
          margin-left: 10px;
          font-weight: bold;
        }

        .trend-value {
          color: #666666;
          margin-left: 10px;
          font-weight: bold;
          font-size: 16px;
        }
      }
    }
  }

  .right-charts {
    flex: 1;
    display: flex;
    gap: 16px;

    .chart-container { 
      flex: 1;
      height: 300px; // 增加高度以避免图表压扁
    }
  }

  .middle-section {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;
  }

  .chart-container {
    background: #fff;
    // border: 1px solid #E9ECEF;
    box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    border-radius: 4px;  
    padding: 20px;
    box-sizing: border-box;
    &.full-width {
      flex: 1;
    }

    .chart-title {
      font-size: 16px;
      color: #333333;
      margin-bottom: 16px;
      font-weight: 500;
    }

    .chart-box {
      height: 240px;

      &.large {
        height: 350px;
      }
    }
  }

  .bottom-section {
    height: calc(100% - 325px);
    background: #fff;
    // border: 1px solid #E9ECEF;
      box-shadow: 0px 2px 4px 0px rgba(76, 122, 227, 0.12);
    border-radius: 4px;   
    padding: 20px;
    box-sizing: border-box;
    .table-title {
      font-size: 16px;
      color: #333333;
      margin-bottom: 16px;
      font-weight: 500;
    }
  }

  .risk-high {
    color: #FF6B6B;
  }

  .risk-medium {
    color: #FFD93D;
  }

  .risk-low {
    color: blue;
  }

  .risk-safe {
    color: #16c60c;
  }

  :deep(.el-table) {
    background-color: transparent;

    .el-table__body tr {
      background-color: #FFFFFF;

      &:hover > td {
        background-color: #F8F9FA !important;
      }
    }

    .el-table__body td {
      border-color: #E9ECEF;
      color: #333333;
    }

    .el-table__header th {
      // background-color: #F1F3F4;
      border-color: #E9ECEF;
      color: #333333;
    }
  }
 
</style>
