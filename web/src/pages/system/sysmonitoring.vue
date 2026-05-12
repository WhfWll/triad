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
     #line1, #line2, #pie1{
        height: 100%;
    }
    .block_box_title{
        position: relative;
        font-size: 14px;
        padding-left: 10px;
        font-weight: 800;
        color: rgba(72,72,102,.87) ;
        border-left:3px solid #4c7ae3;
    }
    .div_block{
        width: 100%;
        height: 298px;
        background-color: #fff;
        box-shadow: 0px 2px 4px 1px rgba(76, 122, 227, 0.11);
        border-radius: 4px;
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
          color:['#65C680 ','#4C7AE3'],  
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
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById(id)); 
            // 绘制图表
            myChart.setOption({ 
                tooltip: {
                    // trigger: 'axis',
					trigger: 'axis', //item数据项图形触发，主要在散点图，饼图等无类目轴的图表中使用。
					axisPointer: {
						// 坐标轴指示器，坐标轴触发有效
						type: 'line' // 默认为直线，可选为：'line' | 'shadow'
					},
					formatter: 'CPU使用率 : {c}%  <br/>{b}<br/>' //{a}（系列名称），{b}（数据项名称），{c}（数值）, {d}（百分比）
                },
                grid:{
                    left:'20',
                    top:'10%',
                    bottom:'24',
                    right:'5%',
                    containLabel: true
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    axisLine:{
                        lineStyle:{
                            color:'#d9e1e4'
                        },
                    },
                    axisLabel:{
                        color:'#4e5b5f'
                    },
                    axisTick:{
                        show:false,
                    },
                    data:xdt
                },
                yAxis: {
                    type: 'value',
                    max: 100,
                    min: 0,
                    interval:20,
                    axisLine:{
                        show: false,
                        lineStyle:{
                            color:'#d9e1e4'
                        },
                    },
                    axisLabel:{
                        color:'#4e5b5f',
                        formatter: function (val) {//百分比显示
                            return val + '%';
                        }
                    },
                    axisTick:{
                        show:false,
                    },
                    splitLine:{
                        lineStyle:{
                            type:'dashed'
                        }
                    }

                },
                series: [{
                    data: ydt,
                    type: 'line',
                    smooth: true,
                    // areaStyle: {
                    //     color: 'rgba(103,194,58,.3)',
                    // },
                    lineStyle:{
                        color:'#4C7AE3'
                    },
                    itemStyle: {
                        borderWidth:1,
                        color:'#4C7AE3'
                    }
                }]
            });
        },
        echartline2:function(id, xdt,ydt){
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById(id)); 
            // 绘制图表
            myChart.setOption({ 
                tooltip: {
                    // trigger: 'axis',
					trigger: 'axis', //item数据项图形触发，主要在散点图，饼图等无类目轴的图表中使用。
					axisPointer: {
						// 坐标轴指示器，坐标轴触发有效
						type: 'line' // 默认为直线，可选为：'line' | 'shadow'
					},
					formatter: '内存使用率 : {c}%  <br/>{b}<br/>' //{a}（系列名称），{b}（数据项名称），{c}（数值）, {d}（百分比）
                },
                grid:{
                    left:'20',
                    top:'10%',
                    bottom:'24',
                    right:'5%',
                    containLabel: true
                },
                xAxis: {
                    type: 'category',
                    boundaryGap: false,
                    axisLine:{
                        lineStyle:{
                            color:'#d9e1e4'
                        },
                    },
                    axisLabel:{
                        color:'#4e5b5f'
                    },
                    axisTick:{
                        show:false,
                    },
                    data:xdt
                },
                yAxis: {
                    type: 'value',
                    max: 100,
                    min: 0,
                    interval:20,
                    axisLine:{
                        show: false,
                        lineStyle:{
                            color:'#d9e1e4'
                        },
                    },
                    axisLabel:{
                        color:'#4e5b5f',
                        formatter: function (val) {//百分比显示
                            return val + '%';
                        }
                    },
                    axisTick:{
                        show:false,
                    },
                    splitLine:{
                        lineStyle:{
                            type:'dashed'
                        }
                    }

                },
                series: [{
                    data: ydt,
                    type: 'line',
                    smooth: true,
                    // areaStyle: {
                    //     color: 'rgba(103,194,58,.3)',
                    // },
                    lineStyle:{
                        color:'#4C7AE3'
                    },
                    itemStyle: {
                        borderWidth:1,
                        color:'#4C7AE3'
                    }
                }]
            });
        },
        //饼图 获取接口
         getpie_database(id, data){ 
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById(id)); 
            // 绘制图表
            myChart.setOption({ 
                tooltip: {
                    trigger: 'item', //item数据项图形触发，主要在散点图，饼图等无类目轴的图表中使用。
					axisPointer: {
						// 坐标轴指示器，坐标轴触发有效
						type: 'shadow' // 默认为直线，可选为：'line' | 'shadow'
					},
                    //formatter: "{a} <br/>{b} : {c} ({d}%)",
                    formatter:function(data2){
                        // console.log(data)
                        return  data2.name+ " : " + " "+data2.percent.toFixed()+"%";
                        }
                },
                legend: {
                    orient: 'vertical',
                    // left: 'right',
                    x:'50%',//水平安放位置，默认为'left'，可选为：'center' | 'left' | 'right' | {number}（x坐标，单位px）
                    y: 'center',//垂直安放位置，默认为top，可选为：'top' | 'bottom' | 'center' | {number}（y坐标，单位px）
                    icon: "circle",   //  这个字段控制形状  类型包括 circle 圆形，triangle 三角形，diamond 四边形，arrow 变异三角形，none 无
                    itemWidth: 10,  // 设置宽度
                    itemHeight: 10, // 设置高度
                    itemGap: 35 ,// 设置间距，
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
