<template>
    <div>
        <div class="  box_bg_shadow">
            <div class="box_title">
                <label>任务组信息</label>
                <i></i>
            </div>
            <div class="basic_info">
                <ul>
                    <li>
                        <div>
                            <label for="">任务组</label>
                            <span>{{basic_info.Name}}</span>
                        </div>  
                        <div>
                            <label for="">任务组风险</label>
                            <span :class="[ 
                                    { 'riskstyle risk_hight': basic_info.riskLevel == 1 } ,
                                    { 'riskstyle risk_middle': basic_info.riskLevel == 2 },
                                    { 'riskstyle risk_low': basic_info.riskLevel == 3 },
                                    { 'riskstyle risk_nofind': basic_info.riskLevel == 4 }]"><i></i>

                            {{ basic_info.Risk }}</span>
                        </div> 
                        <div>
                            <label for="">创建时间</label>
                            <span>{{basic_info.CreateTime}}</span>
                        </div> 
                        <div>
                            <label for="">任务组描述</label>
                            <span>{{basic_info.Describe}}</span>
                        </div>
                    </li> 
                </ul>
                
            </div>
        </div>
        <div class=" box_bg_shadow">
            <div class="box_title">
                <label>任务统计</label>
                <i></i>
            </div>
            <div class="box_content">
                <el-row>
                    <el-col :span="6"> 
                        <img class="risk_icon" src="../../assets/images/risk_height@2x.png" />
                        <div class="risk_content risk_content_height">
                            <p class="risk_number">
                                <label  >{{  TaskStat[0] }}</label> 
                            </p>
                            <p class="risk_name" >高危任务</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_middle@2x.png" />
                        <div class="risk_content risk_content_middle">
                            <p class="risk_number">
                                <label   >{{  TaskStat[1] }}</label>  
                            </p>
                            <p class="risk_name" >中危任务</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_low@2x.png" />
                        <div class="risk_content risk_content_low">
                            <p class="risk_number">
                                <label   >{{  TaskStat[2] }}</label> 
                                </p>
                            <p class="risk_name" >低危任务</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_safe@2x.png" />
                        <div class="risk_content risk_content_safe">
                            <p class="risk_number"><label    >{{  TaskStat[3] }}</label>  
                            </p>
                            <p class="risk_name" >安全任务</p>
                        </div>
                    </el-col>
                </el-row>
            </div>
        </div>
        <div class=" box_bg_shadow">
            <div class="box_title">
                <label>目标统计</label>
                <i></i>
            </div>
            <div class="box_content">
                <el-row>
                    <el-col :span="6"> 
                        <img class="risk_icon" src="../../assets/images/risk_height@2x.png" />
                        <div class="risk_content risk_content_height">
                            <p class="risk_number">
                                <label  >{{  TargetStat[0] }}</label> 
                            </p>
                            <p class="risk_name" >高危目标</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_middle@2x.png" />
                        <div class="risk_content risk_content_middle">
                            <p class="risk_number">
                                <label   >{{  TargetStat[1] }}</label>  
                            </p>
                            <p class="risk_name" >中危目标</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_low@2x.png" />
                        <div class="risk_content risk_content_low">
                            <p class="risk_number">
                                <label   >{{  TargetStat[2] }}</label> 
                                </p>
                            <p class="risk_name" >低危目标</p>
                        </div>
                    </el-col>
                    <el-col :span="6">
                        <img class="risk_icon" src="../../assets/images/risk_safe@2x.png" />
                        <div class="risk_content risk_content_safe">
                            <p class="risk_number"><label    >{{  TargetStat[3] }}</label>  
                            </p>
                            <p class="risk_name" >安全目标</p>
                        </div>
                    </el-col>
                </el-row>
            </div>
        </div>
        <div >
            <el-row :gutter="16">
                <el-col :span="8">
                    <div class="box_bg_shadow">
                        <div class="box_title">
                            <label>漏洞风险</label> 
                        </div>
                        <div class="vuln_risk" >
                            <div id="vuln_risk">

                            </div>
                            <div class="legend">
                                <el-row >
                                    <el-col :span="6" class="vlun_risk_deadly">
                                        <div class="name"><i></i>致命漏洞</div>
                                        <div class="number" style="cursor: pointer;" @click="$emit('actLD',{tab:'tabs4',level:1})">{{risk_num[0]}}</div>
                                        <div class="percentage">{{risk_percentage[0]}}%</div>
                                    </el-col>
                                    <el-col :span="6" class="vlun_risk_height">
                                        <div class="name"><i></i>高危漏洞</div>
                                        <div class="number" style="cursor: pointer;" @click="$emit('actLD',{tab:'tabs4',level:2})"  >{{risk_num[1]}}</div>
                                        <div class="percentage">{{risk_percentage[1]}}%</div>
                                    </el-col>
                                    <el-col :span="6" class="vlun_risk_middle">
                                        <div class="name"><i></i>中危漏洞</div>
                                        <div class="number"  style="cursor: pointer;" @click="$emit('actLD',{tab:'tabs4',level:3})">{{risk_num[2]}}</div>
                                        <div class="percentage">{{risk_percentage[2]}}%</div>
                                    </el-col>
                                    <el-col :span="6" class="vlun_risk_low">
                                        <div class="name"><i></i>低危漏洞</div>
                                        <div class="number" style="cursor: pointer;" @click="$emit('actLD',{tab:'tabs4',level:4})" >{{risk_num[3]}}</div>
                                        <div class="percentage">{{risk_percentage[3]}}%</div>
                                    </el-col>
                                </el-row>
                            </div>
                        </div>
                    </div>
                </el-col>
                <el-col :span="8">
                    <div class="box_bg_shadow">
                        <div class="box_title">
                            <label>漏洞状态</label> 
                        </div>
                        <div class="vuln_status" >
                            <div id="vuln_status">
                            </div>
                            <div class="legend" style="padding: 0 40px;">
                                <el-row > 
                                    <el-col :span="8" class="use_achievement" v-for="(item,i ) in vulExploitImpact" :key="i"> 
                                        <div class="name"><i ></i>{{ item.label }}</div>
                                        <div class="number"  >{{ item.value }}</div>
                                        <div class="percentage">{{item.percentage}}%</div>
                                    </el-col>
                                </el-row>
                            </div>
                        </div>
                    </div>
                </el-col>
                <el-col :span="8">
                    <div class="box_bg_shadow">
                        <div class="box_title">
                            <label>漏洞取证</label> 
                        </div>
                        <div class="vuln_obtain_evidence" >
                            <div id="vuln_obtain_evidence"> 
                            </div>
                            <div class="legend" style="padding: 0 40px;">
                                <el-row > 
                                    <el-col :span="8" class="Logo_Credentials" v-for="(item,i) in obtain_evidence" :key="i">
                                        <div class="name"><i ></i>{{item.label}}</div>
                                        <div class="number"  >{{ item.value}}</div> 
                                    </el-col>
                                </el-row>
                            </div>
                        </div>
                    </div>
                </el-col>
            </el-row>
            <el-row :gutter="16">
                <el-col :span="12">
                    <div style="height: 230px;" class="box_bg_shadow">
                        <div class="box_title">
                            <label for="" class="">漏洞类型统计</label> 
                        </div>
                        <div id="vuln_type_bar" style="height: 210px;"> </div>
                    </div>
                </el-col>
                <el-col :span="12">
                    <div style="height: 230px;" class="box_bg_shadow">
                        <div class="box_title">
                            <label for="" class="">任务趋势</label>
                        </div>
                        <div id="vuln_line" style="height: 210px;"> </div>
                    </div>
                </el-col>
            </el-row>
        </div>
    </div>
</template>
<style scoped lang="less" >
#vuln_risk,
#vuln_status,
#vuln_obtain_evidence{
    height: 180px;
}
.legend{ 
    height: 70px;
    .name{
        color: rgba(72,72,102,0.64);
        font-size: 13px;
        i{
            display: inline-block;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            margin-right: 8px;
            vertical-align: middle;
        }
    }
    .number{
        font-size: 18px;
        margin: 10px 0 6px; 
        color: rgba(72,72,102,0.87);
    }
    .percentage{
        font-size: 13px;
        color: rgba(72,72,102,0.64);
    }
}
.vlun_risk_deadly { 
    .name i{
        background-color: #F87D7D;
    }
    .number{
        color: #F87D7D;
    }
}
.vlun_risk_height { 
    .name i{
        background-color: #F9B640;
    }
    .number{
        color: #F9B640;
    }
}
.vlun_risk_middle{
    .name i{
        background-color: @theme-BaseColor;
    }
    .number{
        color: @theme-BaseColor;
    }
}
.vlun_risk_low{
    .name i{
        background-color: #65C680;
    }
    .number{
        color: #65C680;
    }
}
.use_achievement{
    .name i{
        background-color: #6C63F0;
    }
}
.verify_existence{
    .name i{
        background-color: #3AACFF;
    }
}
.Not_verified{
    .name i{
        background-color: #FB466C;
    }
}
.Remote_control{
    .name i{
        background-color: #9A6BFF;
    }
}
.data_breach{
    .name i{
        background-color: #8ED5F1;
    }
}
.Logo_Credentials{
    .name i{
        background-color: #FF8787;
    }
}
.basic_info{   
    padding: 24px 0 10px 24px; 
    box-sizing: border-box;
    ul{
        list-style: none;
    }
    li{
        // margin-bottom: 20px;
        >div{
            min-width: 25%;
            display: inline-block; 
            label{ 
                display: inline-block;
                padding-left: 8px;
                // border-left: 2px solid #4C7AE3;
                color: rgba(72,72,102,0.87);
                font-size: 13px;
                width: 70px;
                height: 14px ;
                line-height: 14px;
            }
            span{ 
                font-size: 13px;
                color: rgba(72,72,102,0.64);
            }
        }
    }
    
     
}
.box_bg_shadow{
    border-radius: 4px;
    padding: 24px;
    background-color: #fff;
    box-shadow: 0px 2px 4px 0px rgba(76,122,227,0.12); 
    margin-bottom: 16px;
    .box_title{
        label{
            border-left: 2px solid #4C7AE3;
            padding-left: 8px;
            color: rgba(72,72,102,0.87);
            font-size: 13px;
        }
        .title_icon{
            display: inline-block;
            width: 32px;
            height: 32px; 
            margin-right: 16px;
            vertical-align: middle;
        }
        span{
            display: inline-block;
            font-size: 13px;
            color: rgba(72,72,102,0.87);
        }
    }
    .box_content{
        padding: 40px 0 16px;
        .risk_icon{
            display: inline-block;
            width: 48px;
            height: 48px; 
            margin-right: 24px;
            vertical-align: middle;
        }
        .risk_content{
            display: inline-block;
            vertical-align: middle;
            .risk_number{
                font-size: 18px;
                margin-bottom:8px; 
                line-height: 12px;
                color: rgba(72,72,102,0.87);
                label{
                    vertical-align: middle;
                    
                }
                i{
                    font-size: 12px;
                }
            }
            .risk_name{
                font-size: 13px;
                color: rgba(72,72,102,0.64);;
            }
        }
        .risk_content_height .risk_number{
            color: #F87D7D;
        }
        .risk_content_middle .risk_number{
            color: #F9B640;
        }
        .risk_content_low .risk_number{
            color: #4C7AE3;
        }
        .risk_content_safe .risk_number{
            color: #65C680;
        }
        
    }
    
}
</style>
<script>
import { task_group } from '@/api/task.js'
var echarts = require('echarts');
export default {
    name:'',
    data(){
        return{
            basic_info:{
                CreateTime:'',
                Describe:'',
                Name:'',
                Risk:'',
            },
            TaskStat:[],
            TargetStat:[],
            risk_percentage:[],
            risk_num:[],
            risk_total_num:0,
            vulExploitImpact:[],
            status_num:[],
            status_total_num:0,
            statuobtain_evidence_total_nums_total_num:0,
            obtain_evidence:[], 
            ExpStat:[],
        }
    },
    props:{
        group_name:{},
        group_id:{},
    },
    created(){

    },
    mounted(){
        this.getinfo();
    },
    methods:{
        async getinfo(){
            const res = await task_group.getTaskGroupOverview({
                group_id:this.group_id,
            });
            if(res.code == 200){
                let overview =  JSON.parse(res.data.overview);
                console.log(overview)
                this.basic_info.CreateTime = overview.CreateTime
                this.basic_info.Describe = overview.Describe
                this.basic_info.Name = overview.Name
                this.basic_info.Risk = overview.Risk

                this.TaskStat = overview.TaskStat;
                this.TargetStat = overview.TargetStat;
                 
                let _this = this;
                // 漏洞风险 
                this.risk_num = overview.VulStat;
                this.risk_num.forEach(item =>{
                    _this.risk_total_num+=item
                })   
                this.risk_num.forEach(item =>{  
                    let item_num =_this.risk_total_num == 0? 0: (item/_this.risk_total_num*100).toFixed(2); 
                    _this.risk_percentage.push(item_num);
                }) 
                this.echart_vuln_risk(this.risk_num);
                // 漏洞状态 
                this.ExpStat =  overview.ExpStat;
                this.ExpStat.forEach(item =>{
                    _this.status_total_num+=item;
                    this.status_num.push(item)
                })   
                overview.ExpStat.forEach((item,i) =>{ 
                    let item_num =_this.status_total_num==0 ? 0: (item/_this.status_total_num*100).toFixed(2); 
                    if(i==0){
                        this.vulExploitImpact.push({label:'利用成功',value:item})
                    }
                    if(i==1){
                        this.vulExploitImpact.push({label:'验证存在',value:item})
                    }
                    if(i==2){
                        this.vulExploitImpact.push({label:'未验证',value:item})
                    }
                    
                }) 
                this.echart_vuln_status(this.vulExploitImpact);
                // 漏洞取证 
                overview.EvidenceStat.forEach((item,i) =>{
                    _this.statuobtain_evidence_total_nums_total_num+=item 
                    if(i==0){
                        this.obtain_evidence.push({label:'远程控制',value:item})
                    }
                    if(i==1){
                        this.obtain_evidence.push({label:'数据泄露',value:item})
                    }
                    if(i==2){
                        this.obtain_evidence.push({label:'登录凭证',value:item})
                    }
                })   

                this.echart_vuln_obtain_evidence(this.obtain_evidence);

                //
                let VulTypeStat = overview.VulTypeStat;
                this.echartVulnType(VulTypeStat.x,VulTypeStat.y);
                let TaskTrend =  overview.TaskTrend;
                this.echartvulfindtrend(TaskTrend.x,TaskTrend.y);
            }
            
        },
        echart_vuln_risk(vulRisk){
            let data = [];
            data.push({name:'致命漏洞',value:vulRisk[0]});
            data.push({name:'高危漏洞',value:vulRisk[1]});
            data.push({name:'中危漏洞',value:vulRisk[2]});
            data.push({name:'低危漏洞',value:vulRisk[3]});
            let totalvuln = vulRisk[0]+vulRisk[1]+vulRisk[2]+vulRisk[3];
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('vuln_risk')); 
            let option = {
                color:['#F87D7D','#F9B640','#4C7AE3','#65C680'],
                title:{
                    text:totalvuln,
                    subtext: '漏洞总数',
                    textStyle: {
                        fontsize: 18, 
                        fontWeight:400,
                        color: 'rgba(72,72,102,0.87)', 
                    },
                    subtextStyle: {
                        fontsize: 13,
                        color: 'rgba(72,72,102,0.64)', 
                    },
                    textAlign:"center", 
                    x:"49%", 
                    y:"37%" 
                    
                },
                series: [
                    {
                        name: ' ',
                        top:'top',
                        type: 'pie',
                        radius: ['50%', '70%'],
                        avoidLabelOverlap: false,
                        label: {
                            show: false,
                            position: 'center'
                        }, 
                        labelLine: {
                            show: false
                        },
                        data: data
                    }
                ]
            };
            // 绘制图表
            myChart.setOption(option);
        },
        echart_vuln_status(statuslist){
            let data = [];
            let totalvuln = 0;
            statuslist.forEach(item=>{
                data.push({name:item.label,value:item.value});
                totalvuln+=item.value
            });
            console.log(data)
            // data.push({name:'利用成功',value:statuslist[0]});
            // data.push({name:'验证存在',value:statuslist[1]});
            // data.push({name:'未验证',value:statuslist[2]});
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('vuln_status')); 
            let option = {
                color:['#6C63F0','#3AACFF','#FB466C' ],
                title:{
                    text:totalvuln,
                    subtext: '漏洞总数',
                    textStyle: {
                        fontsize: 18, 
                        fontWeight:400,
                        color: 'rgba(72,72,102,0.87)', 
                    },
                    subtextStyle: {
                        fontsize: 13,
                        color: 'rgba(72,72,102,0.64)', 
                    },
                    textAlign:"center", 
                    x:"49%", 
                    y:"37%" 
                    
                },
                series: [
                    {
                        name: ' ',
                        top:'top',
                        type: 'pie',
                        radius: ['50%', '70%'],
                        avoidLabelOverlap: false,
                        label: {
                            show: false,
                            position: 'center'
                        }, 
                        labelLine: {
                            show: false
                        },
                        data: data
                    }
                ]
            };
            // 绘制图表
            myChart.setOption(option);
        },
        echart_vuln_obtain_evidence(statuslist ){
            let data = [];
            let totalvuln = 0;
            statuslist.forEach(item=>{
                data.push({name:item.label,value:item.value});
                totalvuln+=item.value;
            });
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('vuln_obtain_evidence')); 
            let option = {
                color:['#9A6BFF','#8ED5F1','#FF8787'],
                title:{
                    text:totalvuln,
                    subtext: '漏洞取证',
                    textStyle: {
                        fontsize: 18, 
                        fontWeight:400,
                        color: 'rgba(72,72,102,0.87)', 
                    },
                    subtextStyle: {
                        fontsize: 13,
                        color: 'rgba(72,72,102,0.64)', 
                    },
                    textAlign:"center", 
                    x:"49%", 
                    y:"37%" 
                    
                },
                series: [
                    {
                        name: ' ',
                        top:'top',
                        type: 'pie',
                        radius: ['50%', '70%'],
                        avoidLabelOverlap: false,
                        label: {
                            show: false,
                            position: 'center'
                        }, 
                        labelLine: {
                            show: false
                        },
                        data: data
                    }
                ]
            };
            // 绘制图表
            myChart.setOption(option);
        },
        echartVulnType(xdt,xdy){ 
            var myChart = echarts.init(document.getElementById('vuln_type_bar')); 
            this.eVulnType = myChart;
            let count = 66;
            let _data = [
                {label: '高危',value:1},
                {label: '中危',value:4},
                {label: '低危',value:10},
                {label: '安全',value:50},
                {label: '中危1',value:4},
                {label: '低危2',value:10},
                {label: '安全3',value:50}
            ];
            let dtx = [];
            for(var i =0;i<_data.length;i++){
                dtx.push(_data[i].label)
            }
            // 绘制图表
            myChart.setOption( { 
                grid:{
                    left:'10',
                    top:'5%',
                    bottom:'8',
                    right:'26',
                    containLabel: true
                },
                tooltip: {
                    trigger: 'axis',
                    axisPointer: {
                        type: 'shadow'
                    },
                    // formatter: function (params) {
                    //     var tar = params[1];
                    //     console.log(params)
                    //     // return tar.name + '<br/>' + tar.seriesName + ' : ' + tar.value;
                    //     return 1
                    // }
                }, 
                xAxis: {  
                    data: xdt,
                    axisLabel:{
                        show:true, 
                        interval: 0, 
                        rotate: 30,
                        color: 'rgba(72,72,102,0.64)',
                        fontSize:12,
                    },
                    axisLine:{
                        show:false,
                    },
                    axisTick:{
                        show:false,
                    },  
                },
                yAxis: {
                    type: 'value', 
                    axisTick:{
                        show:false,
                    },
                    axisLine:{
                        show:false,
                    }, 
                    axisLabel:{
                        show:true, 
                        color: 'rgba(72,72,102,0.64)',
                        fontSize:12,
                    },
                    splitLine:{
                        lineStyle:{
                            type:'dashed'
                        }
                    }
                }, 
               
                // dataZoom:[
                // {
                // id: 'dataZoomZ',
                // xAxisIndex: [0],
                // show: true, //是否显示滑动条，不影响使用
                // type: 'slider', // 这个 dataZoom 组件是 slider 型 dataZoom 组件
                // startValue: 0, // 从头开始。
                // endValue: 9, // 一次性展示5个
                // height: 6,
                // borderColor: 'transparent',
                // fillerColor: 'rgba(205,205,205,1)',
                // zoomLock: true,
                // showDataShadow: false, //是否显示数据阴影 默认auto
                // backgroundColor: 'rgba(0,0,0,0)',
                // showDetail: false, //即拖拽时候是否显示详细数值信息 默认true
                // realtime: true, //是否实时更新
                // filterMode: 'filter',
                // handleIcon: 'circle',
                // handleStyle: {
                //     color: 'rgba(205,205,205,1)',
                //     borderColor: 'rgba(205,205,205,1)',
                // },
                // handleSize: '80%',
                // moveHandleSize: 0,
                // maxValueSpan: 9,
                // minValueSpan: 9,
                // brushSelect: false, //刷选功能，设为false可以防止拖动条长度改变 ************（这是一个坑）
                // bottom:0,
                // }, {
                //     type:'inside',
                //     xAxisIndex:0,
                //     zoomOnMouseWheel:false,  //滚轮是否触发缩放
                //     moveOnMouseMove:true,  //鼠标滚轮触发滚动
                //     moveOnMouseWheel:true
                // }],
                series: [
                    {  
                        data: xdy,
                        type: 'bar',
                        barGap:'10%',
                        barCategoryGap:'20%',
                        barWidth:'30', 
                    }
                ]
            })
        },
        echartvulfindtrend:function(xdt,ydt){ //趋势
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(document.getElementById('vuln_line')); 
            this.evulfindtrend = myChart;
            // 绘制图表
            myChart.setOption({ 
                grid:{
                    left:'20',
                    top:'5%',
                    bottom:'14',
                    right:'5%',
                    containLabel: true
                },
                tooltip: {
                    trigger: 'axis',
                   
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
                        color: 'rgba(72,72,102,0.64)',
                        fontSize:12,
                    },
                    axisTick:{
                        show:false,
                    },
                    data:xdt
                },
                yAxis: { 
                    type: 'value',
                    axisLine:{
                        lineStyle:{
                            color:'#d9e1e4'
                        },
                    },
                    axisLabel:{
                        color: 'rgba(72,72,102,0.64)',
                        fontSize:12,
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
                    showSymbol: false,  
                    lineStyle:{
                        color:'#4C7AE3'
                    },
                    itemStyle: {
                        borderWidth:1,
                        color:'#4C7AE3'
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
            });
        },
    }
}
</script>