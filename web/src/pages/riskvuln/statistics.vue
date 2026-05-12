<template>
    <div>
        <div class="main-title">  
            漏洞统计
        </div> 
        <div v-loading="loading">
            <!-- 顶部区域 -->
            <div class="top-section" >
                <!-- 左侧卡片区域 -->
                <div class="left-cards">
                    <div class="card-item">
                        <div class="card-icon">
                            <el-icon :size="38" color="#E6A23C">
                            <TrendCharts />
                            </el-icon>
                        </div>
                        <div class="card-content">
                            <div class="card-title">总漏洞数</div>
                            <div class="card-value">{{ totalVulnerabilities }}</div>
                        </div>
                    </div>
                    <div class="card-item">
                        <div class="card-icon">
                            <el-icon :size="38" color="#409EFF">
                            <DataBoard />
                            </el-icon>
                        </div>
                        <div class="card-content">
                            <div class="card-title">漏洞平均修复时间</div>
                            <div class="card-value">{{ averageFixTime }}</div>

                        </div>
                    </div>

                    <div class="card-item">
                        <div class="card-icon">
                            <el-icon :size="38" color="#67C23A">
                            <TrendCharts />
                            </el-icon>
                        </div>
                        <div class="card-content">
                            <div class="card-title">高危漏洞数</div>
                            <div class="card-value">{{ highRiskCount }} </div>

                        </div>
                    </div>  
                </div>

                <!-- 右侧图表区域 -->
                <div class="right-charts">
                    <div class="chart-container">
                    <div class="chart-title">漏洞状态统计</div>
                    <div ref="riskStatChart" class="chart-box" id="riskStatChart"></div>
                    </div>

                    <div class="chart-container">
                    <div class="chart-title">漏洞等级统计</div>
                    <div ref="riskLevelChart" class="chart-box" id="riskLevelChart"></div>
                    </div>
                </div>
            </div> 
            <!-- 中间区域 -->
            <div class="middle-section">
                <div class="chart-container full-width">
                    <div class="chart-title">漏洞类型统计</div>
                    <div ref="riskTypeChart" class="chart-box large" id="riskTypeChart"></div>
                </div>

                <div class="chart-container full-width">
                    <div class="chart-title">TOP10漏洞统计</div> 
                    <el-table :data="top10Vulnerabilities" style="width: 100%" size="mini">
                        <el-table-column prop="name" label="名称" >
                        </el-table-column>
                        <el-table-column prop="num" label="数量" width="100" >
                        </el-table-column>
                        <el-table-column prop="riskLevel" label="等级" width="80">
                            <template #default="scope"> 
                                <div class="risk_style">
                                    <el-tag type="danger"  class="risk risk4" size="mini" v-if="scope.row.riskLevel==0"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                                    <el-tag type="danger" class="risk risk0" size="mini" v-if="scope.row.riskLevel==1"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                                    <el-tag type="warning" class="risk risk1" size="mini" v-if="scope.row.riskLevel==2"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                                    <el-tag type="primary" class="risk risk2"  size="mini"  v-if="scope.row.riskLevel==3"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                                    <el-tag type="success" class="risk risk3" size="mini"  v-if="scope.row.riskLevel==4"><i></i>{{scope.row.riskLevelStr}}</el-tag> 
                                    <el-tag type="info"  class="risk risk4" size="mini"  v-if="scope.row.riskLevel==5"><i></i>{{scope.row.riskLevelStr}}</el-tag>
                                </div>  
                            </template>
                        </el-table-column>
                    </el-table>
    
                </div>
            </div> 
        </div> 
    </div>
</template>
<script>
import risk from '@/api/risk.js'
var echarts = require('echarts');
export default {
    data(){
        return{
            loading:false,
            averageFixTime:'', //漏洞平均修复时间
            highRiskCount:0,
            totalVulnerabilities:'',

            statusStatistics:[], //漏洞状态统计
            riskLevelStatistics:[],//漏洞等级统计
            typeStatistics:[],//漏洞类型统计
            top10Vulnerabilities:[],//top10漏洞
        }
    },
    created(){
        this.$store.state.activefirstMenu = "/riskvulnStatistics";
    },
    mounted(){
        this.getData();
    },
    methods:{
        async getData(){
            this.loading = true;
            var res = await risk.riskvulstatistics();
            if(res.code == 200){
                this.averageFixTime =res.data.averageFixTime;
                this.highRiskCount =res.data.highRiskCount;
                this.totalVulnerabilities = res.data.totalVulnerabilities;
                this.top10Vulnerabilities = res.data.top10Vulnerabilities;

                this.statusStatistics = res.data.statusStatistics;
                this.riskLevelStatistics = res.data.riskLevelStatistics;
                this.typeStatistics = res.data.typeStatistics;

                this.initCharts();
                 this.loading = false;
            }
        },
        initCharts(){
            this.initRiskStatChart(); //漏洞状态
            this.initRisklevelChart(); //漏洞等级
            this.initRiskTypeChart(); //漏洞类型
        },
        //漏洞状态
        initRiskStatChart(){
            const chart = echarts.init(document.getElementById('riskStatChart'));

            // 处理数据
            const chartData = this.statusStatistics.map(item => ({
                name: item.statusName || '未知',
                value: item.count || 0
            }));

            // 计算背景条的最大值，确保至少为1
            const maxValue = Math.max(Math.max(...chartData.map(item => item.value)) * 1.2, 1);
            const backgroundData = new Array(chartData.length).fill(maxValue);
            // 定义多种颜色
            const colors = ['#fa5949','#67C23A',  ];
            const option = {
                color: ['#fa5949'],
                grid: { top: 10, right: 10, bottom: 50, left: 50 },
                tooltip: { trigger: 'item' },
                legend: {
                    orient: 'vertical',
                    left: 'left'
                }, 
                series: [ 
                    {
                        name: '漏洞状态',
                        type: 'pie',
                        radius: '70%',
                        
                        itemStyle: {
                            color: function(params) {
                                // 为每个柱子分配不同的颜色
                                return colors[params.dataIndex % colors.length];
                            },
                            borderRadius: [6, 6, 0, 0],
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
        //漏洞等级
        initRisklevelChart(){
            const chart = echarts.init(document.getElementById('riskLevelChart'));

            // 处理数据
            const chartData = this.riskLevelStatistics.map(item => ({
                name: item.riskName || '未知',
                value: item.count || 0
            }));

            // 计算背景条的最大值，确保至少为1
            const maxValue = Math.max(Math.max(...chartData.map(item => item.value)) * 1.2, 1);
            const backgroundData = new Array(chartData.length).fill(maxValue);
            // 定义多种颜色
            const colors = ['#fa5949','#FFD93D','#409EFF', '#67C23A','#484866a3','#dcdfe6',];
            const option = {
                color: ['#fa5949'],
                grid: { top: 10, right: 10, bottom: 50, left: 50 },
                tooltip: { trigger: 'axis' },
                xAxis: {
                type: 'category',
                data: chartData.map(item => item.name),
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
                
                {
                    name: '漏洞状态',
                    type: 'bar',
                    barWidth: 13,
                    itemStyle: {
                    color: function(params) {
                        // 为每个柱子分配不同的颜色
                        return colors[params.dataIndex % colors.length];
                    },
                    borderRadius: [6, 6, 0, 0],
                    },
                    data: chartData.map(item => item.value)
                }
                ]
            };
            chart.setOption(option);
            window.addEventListener('resize', function () {
                chart.resize()
            })
        },
        //漏洞类型
        initRiskTypeChart(){
            const chart = echarts.init(document.getElementById('riskTypeChart'));
            const data = this.typeStatistics.map(item => ({
                name: item.typeName || '未知',
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
        }
    }
}
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
    gap: 20px;
    margin-bottom: 20px;
  }

  .left-cards {
    display: flex;
    flex-direction: column;
    gap: 16px;
    width: 300px;
    height: 400px; // 与右侧图表高度一致
    justify-content: space-between;
  }

  .card-item {
    display: flex;
    align-items: center;
    background: #fff;
    border: 1px solid #E9ECEF;
    border-radius: 12px;
    padding: 20px;
    height: 110px;
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
        margin-bottom: 8px;
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
    gap: 20px;

    .chart-container { 
      flex: 1;
      height: 360px; // 增加高度以避免图表压扁
    }
  }

  .middle-section {
    display: flex;
    gap: 20px;
    margin-bottom: 20px;
  }

  .chart-container {
    background: #fff;
    border: 1px solid #E9ECEF;
    border-radius: 12px;
    padding: 20px;

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
      height: 340px;

      &.large {
        height: 420px;
      }
    }
  }

  .bottom-section {
    background: #fff;
    border: 1px solid #E9ECEF;
    border-radius: 12px;
    padding: 20px;

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
