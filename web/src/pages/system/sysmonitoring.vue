//系统设置---系统工具---倒数第二个Tab  系统监控 页面代码
<template>
    <div style="height: 100%;overflow-x: hidden; overflow-y: auto;">
        <el-row :gutter="20" >
            <el-col :span="12">
                <div class="div_block">
                    <strong class="block_box_title">最近一小时CPU使用率</strong>
                    <div id="line1"></div>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="div_block">
                    <strong class="block_box_title">最近一小时内存使用率</strong>
                    <div id="line2"></div>
                </div>
            </el-col>
            <el-col :span="12">
                <div class="div_block">
                    <strong class="block_box_title">硬盘使用情况</strong>
                    <div id="pie1"></div>
                </div>
            </el-col>
        </el-row>
    </div>
</template>
<style scoped lang="less">
@import '../security/css/appsec-tokens.less';

     #line1, #line2, #pie1{
        height: 100%;
    }
    .block_box_title{
        position: relative;
        font-size: @appsec-font-size-base;
        padding-left: 10px;
        font-weight: @appsec-font-weight-semibold;
        color: @appsec-text-primary;
        border-left: 3px solid @appsec-accent;
    }
    .div_block{
        width: 100%;
        height: 298px;
        background: @appsec-bg-surface;
        box-shadow: @appsec-shadow-card;
        border: 1px solid @appsec-border-default;
        border-radius: @appsec-radius-md;
        margin-bottom: 24px;
        padding: 24px;
        box-sizing: border-box;
        /deep/ .el-form-item__label{
            height:32px;
            line-height: 32px;
        }
        /deep/ .el-form-item__content{
            height:32px;
            line-height: 32px;
        }
    }
</style>
<script> 
var echarts = require('echarts');
import { system } from '@/api/system.js'
import $ from 'jquery'
export default({
    name:'sysmonitoring',
    components: {
  	},
    data(){ 
        //系统监控

    	return{
          color: ['#34d399', '#00d4aa'],
          chartColors: {
            axisLabel: '#64748b',
            axisLine: 'rgba(255, 255, 255, 0.12)',
            splitLine: 'rgba(255, 255, 255, 0.08)',
            line: '#00d4aa',
            legend: '#94a3b8',
          },
          timer1: null,
          Loading:false,
          timermillisec:0,
          total: 0
    	}
    }, 
    created:function(){
        this.timermillisec = this.commonjs.timermillisec;
    },
    mounted:function(){
           this.getvul_database_trend();
            this.getvul_database_trend1();
            this.getpiedata3();
        this.timer1 = setInterval(() => {
			this.getvul_database_trend(true);
            this.getvul_database_trend1(true);
            this.getpiedata3(true);
    	}, 120000) 
        
    },
    beforeDestroy(){
            //页面销毁清除定时器
    this.timer1 ? clearInterval(this.timer1) : "";
    clearTimeout(this.commonjs.timeer);
        this.commonjs.timeer = null;
    },
    methods:{
        chartTooltip(extra) {
            return Object.assign({
                backgroundColor: 'rgba(26, 26, 46, 0.95)',
                borderColor: 'rgba(0, 212, 170, 0.2)',
                textStyle: { color: '#e2e8f0' },
            }, extra || {})
        },
        chartXAxis(xdt) {
            const c = this.chartColors
            return {
                type: 'category',
                boundaryGap: false,
                data: xdt,
                axisLine: { lineStyle: { color: c.axisLine } },
                axisLabel: { color: c.axisLabel },
                axisTick: { show: false },
            }
        },
        chartYAxis() {
            const c = this.chartColors
            return {
                type: 'value',
                max: 100,
                min: 0,
                interval: 20,
                axisLine: { show: false, lineStyle: { color: c.axisLine } },
                axisLabel: {
                    color: c.axisLabel,
                    formatter: function (val) {
                        return val + '%'
                    },
                },
                axisTick: { show: false },
                splitLine: { lineStyle: { type: 'dashed', color: c.splitLine } },
            }
        },
        //系统监控 折线图1  获取接口
          async getvul_database_trend1(){ 
            const dt = await system.getDatabasetrend();
            if(dt.code === 200){
                let xData = []
                let yData = []
                dt.data.list.forEach(item => {
                    xData.push(item.xData)
                    yData.push(item.yData)
                })
                this.echartline('line1', xData,yData);
            }else{
                this.$message({
                    message:dt.msg,
                    type: 'error'
                });
            }  
        },
        //系统监控 折线图2 获取接口
        async getvul_database_trend(){ 
            const dt = await system.getDataneicun();
            if(dt.code ===200){
                let xData = []
                let yData = []
                dt.data.list.forEach(item => {
                    xData.push(item.xData)
                    yData.push(item.yData)
                })
                this.echartline2('line2', xData,yData);
            }else{
                this.$message({
                    message:dt.msg,
                    type: 'error'
                });
            }  
        },
        //系统监控 饼图 获取接口
        async getpiedata3(){ 
            const dt = await system.getpiedata();
            if(dt.code === 200){
                this.total = dt.data.total
                let data = [
                    {
                        value: dt.data.free,
                        percent: dt.data.freePercent,
                        name: '可用空间'
                    },
                    {
                        value: dt.data.used,
                        percent: dt.data.usedPercent,
                        name: '已用空间'
                    },
                ]
                
                this.getpie_database('pie1', data);
            }else{
                this.$message({
                    message:dt.msg,
                    type: 'error'
                });
            }  
        },
        echartline:function(id, xdt,ydt){
            var myChart = echarts.init(document.getElementById(id));
            var lineColor = this.chartColors.line
            myChart.setOption({
                tooltip: this.chartTooltip({
                    trigger: 'axis',
                    axisPointer: { type: 'line' },
                    formatter: 'CPU使用率 : {c}%  <br/>{b}<br/>',
                }),
                grid:{
                    left:'20',
                    top:'10%',
                    bottom:'24',
                    right:'5%',
                    containLabel: true
                },
                xAxis: this.chartXAxis(xdt),
                yAxis: this.chartYAxis(),
                series: [{
                    data: ydt,
                    type: 'line',
                    smooth: true,
                    lineStyle:{ color: lineColor },
                    itemStyle: {
                        borderWidth:1,
                        color: lineColor
                    }
                }]
            });
        },
        echartline2:function(id, xdt,ydt){
            var myChart = echarts.init(document.getElementById(id));
            var lineColor = this.chartColors.line
            myChart.setOption({
                tooltip: this.chartTooltip({
                    trigger: 'axis',
                    axisPointer: { type: 'line' },
                    formatter: '内存使用率 : {c}%  <br/>{b}<br/>',
                }),
                grid:{
                    left:'20',
                    top:'10%',
                    bottom:'24',
                    right:'5%',
                    containLabel: true
                },
                xAxis: this.chartXAxis(xdt),
                yAxis: this.chartYAxis(),
                series: [{
                    data: ydt,
                    type: 'line',
                    smooth: true,
                    lineStyle:{ color: lineColor },
                    itemStyle: {
                        borderWidth:1,
                        color: lineColor
                    }
                }]
            });
        },
        //饼图 获取接口
         getpie_database(id, data){ 
            var myChart = echarts.init(document.getElementById(id));
            myChart.setOption({ 
                tooltip: this.chartTooltip({
                    trigger: 'item',
                    axisPointer: { type: 'shadow' },
                    formatter:function(data2){
                        return  data2.name+ " : " + " "+data2.percent.toFixed()+"%";
                        }
                }),
                legend: {
                    orient: 'vertical',
                    x:'50%',
                    y: 'center',
                    icon: "circle",
                    itemWidth: 10,
                    itemHeight: 10,
                    itemGap: 35,
                    textStyle: {
                        color: this.chartColors.legend,
                        fontSize: 12,
                    },
                    formatter: function (name) {
                        // console.log(data, 'data')
                        let total = 0
                        let tarValue
                        for (let i = 0; i < data.length; i++) {
                            total += data[i].value
                            if (data[i].name == name) {
                            tarValue = data[i].value
                            }
                        }
                        let v = tarValue + 'G'
                        //计算出百分比
                        let p = Math.round((tarValue / total) * 100) + '%'  
                        return `${name}  ${p}   ${v}`
                        //name是名称，v是数值
                    }
                    // textStyle: {  // 图例文字的样式
                    //     color: '#fff',
                    //     fontSize: 16,
                    // }
                },
                color:this.color,
                series: [
                    {
                    // name: 'Access From',
                    type: 'pie',
                    radius: '70%',
                    center:['20%','50%'],
                    label: {//去除饼图的指示折线label
                        normal: {
                        show: false,
                        position: 'inside',
                        formatter:"{b}:{d}%"
                        },
                    },
                    data: data,
                    emphasis: {
                        itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                        }
                    }
                    }
                ]
            });
        },
    }
})
 
</script>
