<template>
	<div class="content">
        <div class="exportbutton" id="exportbutton"> 
            <div v-if="type=='pdf'" >
                <el-button size="small" style="margin-right:8px" @click="exportPdf">导出PDF</el-button> 
            </div> 
        </div> 
		<div id="file">
            <div class="reportbox">
                <div class="boxtitlebg" :style="backgroundStyle">
                    <div class="largetitle">{{ report_cover_title }}</div> 
                    <div class="smalltitle">{{ report_cover_create_time }}</div>
                </div>
                <div class="reportcontent">
                    <div>
                        <div class="report-item">
                            <div>
                                <div class="item_title" id="reportCatalog" @click="foldingBtn('catalog', 'reportCatalog')">
                                    <span>报告目录</span>
                                    <!-- <i class="icon iconfont"></i> -->
                                    <span class="iconfont iconxialashixintop"> </span>
                                </div>
                            </div>
                            <div class="catalog">
                                <ul id="catalogUl" style="margin: 0 25px 10px 25px;"> 
                                    <div  v-for="(catalog,i) in cataloglist " :key="i">
                                        <li class="level1" >
                                            <span class="spanText"><a :href="'#'+catalog.id">{{ catalog.name }}</a></span>
                                            <span class="spanShengluehao"></span>
                                        </li>  
                                        <li class="level2"   v-for="(catalog_child,child_i) in catalog.catalog" :key="child_i">
                                            <span class="spanText">
                                                <a :href="'#'+catalog_child.id" > {{ catalog_child.name }} </a>
                                            </span>
                                            <span class="spanShengluehao"></span>
                                        </li>
                                    </div> 
                                </ul>
                            </div> 
                        </div> 
                    </div>
                </div>
                <div class="thirdpart">
                    <div >
                        <div class="commontitle 2_1part" style="margin-bottom: 25px;" id="taskOverview" @click="foldingBtn('taskOverview', 'taskOverview')">
                            <span id="span-taskOverview" >任务概述</span>
                            <span class="iconfont iconxialashixintop"> </span>
                        </div>
                        <div class="taskOverview" style="padding:0 25px ;">
                            <table class=" firsttable" style="width:100%">
                                <tbody>
                                <tr>
                                    <td class="nodarkback baogaozhaiyaoArea-tdLeft">任务名称：</td>
                                    <td class="lightback baogaozhaiyaoArea-tdRight" id="task-name">{{taskOverview.taskName}}</td> 
                                </tr>
                                <tr>
                                    <td class="darkback baogaozhaiyaoArea-tdLeft ">风险等级：</td>
                                    <td class="darkback baogaozhaiyaoArea-tdRight redpot" id="task-risk">
                                        <span class="levelcolor " style="width: 8px;"
                                            :class="[
                                                {'highcolor':taskOverview.taskRiskStr == '高危'},
                                                {'middlecolor':taskOverview.taskRiskStr == '中危'},
                                                {'lowcolor':taskOverview.taskRiskStr == '低危'},
                                                {'infocolor':taskOverview.taskRiskStr == '安全' || taskOverview.taskRiskStr == '未发现'}
                                            ]"></span>
                                            {{taskOverview.taskRiskStr}}</td>
                                </tr>
                                <tr>
                                    <td class="nodarkback baogaozhaiyaoArea-tdLeft">目标分布：</td>
                                    <td class="lightback baogaozhaiyaoArea-tdRight" id="task-plan">
                                        <div>
                                            <span>总目标：</span>
                                            <span id="target-count" class="spanR">{{taskOverview.targetStat.total}}</span>
                                            <span>存活目标：</span>
                                            <span id="alive-count" class="spanR">{{taskOverview.targetStat.liveTarget}}</span> 
                                        </div>
                                        <div>
                                            <span>高危目标：</span>
                                            <span id="high-targetCount" class="spanR">{{taskOverview.targetStat.HighTarget}}</span>
                                            <span>中危目标：</span>
                                            <span id="middle-targetCount" class="spanR">{{taskOverview.targetStat.middleTarget}}</span>
                                            <span>低危目标：</span>
                                            <span id="low-targetCount" class="spanR">{{taskOverview.targetStat.lowTarget}}</span>
                                            <span>安全目标：</span>
                                            <span id="safe-targetCount">{{taskOverview.targetStat.safeTarget}}</span>
                                        </div>
                                    </td>
                                </tr>
                                <tr>
                                    <td class="darkback baogaozhaiyaoArea-tdLeft">漏洞分布：</td>
                                    <td class="darkback baogaozhaiyaoArea-tdRight" id="task-count">
                                        <div><span>漏洞总数：</span><span id="bug-count">{{taskOverview.vulnStat.total}}</span></div>
                                        <div>
                                            <span>致命漏洞：</span>
                                            <span id="high-bugCount" class="spanR">{{taskOverview.vulnStat.deadlyNumber}}</span>
                                            <span>高危漏洞：</span>
                                            <span id="middle-bugCount" class="spanR">{{taskOverview.vulnStat.highNumber}}</span>
                                            <span>中危漏洞：</span>
                                            <span id="low-bugCount" class="spanR">{{taskOverview.vulnStat.middleNumber}}</span>
                                            <span>低危漏洞：</span>
                                            <span id="info-count">{{taskOverview.vulnStat.lowNumber}}</span>
                                        </div>
                                    </td>
                                </tr>
                                <tr>
                                    <td class="nodarkback baogaozhaiyaoArea-tdLeft">漏洞验证：</td>
                                    <td class="lightback baogaozhaiyaoArea-tdRight">
                                        <div>
                                            <span>验证成功：</span>
                                            <span id="validateSuccess-count" class="spanR">{{taskOverview.vulnVerify.verifySuccess}}</span>
                                            <span>利用成功：</span><span id="useSuccess-count" class="spanR">{{taskOverview.vulnVerify.useSuccess}}</span>
                                            <span>未验证：</span><span id="unvalidate-count" class="spanR">{{ taskOverview.vulnVerify.repairSuccess }}</span> 
                                        </div>
                                    </td>
                                </tr>
                                
                                <tr>
                                    <td class="nodarkback baogaozhaiyaoArea-tdLeft">任务场景：</td>
                                    <td class="lightback baogaozhaiyaoArea-tdRight" id="task-fangan">{{ taskOverview.templateName }}</td>
                                </tr>
                                <tr>
                                    <td class="darkback baogaozhaiyaoArea-tdLeft">测试时间：</td>
                                    <td class="darkback baogaozhaiyaoArea-tdRight" id="task-shijian">{{taskOverview.date}}</td>
                                </tr> 
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <div >
                        <div class="commontitle" id="taskStat" @click="foldingBtn('taskStat', 'taskStat')">
                            <span id="span-taskStat">信息统计</span>
                            <span class="iconfont iconxialashixintop"> </span>
                        </div>
                        <div class="taskStat">
                            <div class="commontitle width150 xinxitjArea spacialtitle" id="targetRisk">
                                <span id="span-targetRisk">目标风险统计</span>
                            </div> 
                            <div class="targetRisk">
                                <div class="loopStatisticschart  " style="position: relative;">
                                    <div class="totaltarget">
                                        <span>合计 : </span>
                                        <span class="totalnumber targetStatisticCount">{{ risk_total }}</span>
                                        <span>个目标</span>
                                    </div>
                                    <div class="el-row">
                                        <div id="targetschartpie" class="Statisticschart">
                                                                                
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        
                        
                        <div class="commontitle width150 2_2part xinxitjArea" style="margin-bottom: 25px;" id="vulRisk">
                            <span id="span-vulRisk">漏洞风险统计</span>
                        </div>
                        <div class="vulRisk loudongfxtj" style="padding: 0 25px;"> 
                            <el-table  
                                :data="data.vulRisk" 
                                tooltip-effect="dark" 
                                style=" margin-right: 25px;">
                                <el-table-column   prop="riskType" label="风险类型"  > 
                                </el-table-column>
                                <el-table-column   prop="verifySuccess" label="验证存在"  > 
                                </el-table-column>
                                <el-table-column   prop="repairSuccess" label="未验证"  > 
                                </el-table-column>
                                <el-table-column   prop="useSuccess" label="利用成功"  > 
                                </el-table-column>
                                <el-table-column   prop="total" label="漏洞总数"  > 
                                </el-table-column>
                                <el-table-column   prop="percent" label="漏洞类型占比"  > 
                                </el-table-column>
                            </el-table>
                        </div> 
                        <div class="commontitle 2_5part xinxitjArea" style="margin-bottom: 25px; width:164px;" id="vulType">
                            <span id="span-vulType">漏洞类型统计</span>
                        </div>
                        <div class="" style="padding: 0 25px;">
                            <el-table  
                                :data="data.vulType" 
                                tooltip-effect="dark" 
                                style=" margin-right: 25px;">
                                <el-table-column   prop="vulnType" label="漏洞类型"  > 
                                </el-table-column>
                                <el-table-column   prop="total" label="数量"  > 
                                </el-table-column>
                                <el-table-column   prop="percent" label="占比"  > 
                                </el-table-column>
                                <el-table-column   prop="targetNumber" label="影响目标数量"  > 
                                </el-table-column>
                            </el-table>
                        </div> 
                        <div class="commontitle 2_6part xinxitjArea" style="margin-bottom: 25px; width:164px;" id="topVulRisk">
                            <span id="span-topVulRisk">Top 危险漏洞 </span>
                        </div>
                        <div class="xinxitjArea topVulRisk" style="padding: 0 25px;"> 
                            <el-table  
                                :data="data.topVulRisk" 
                                tooltip-effect="dark" 
                                style=" margin-right: 25px;">
                                <el-table-column   prop="vulName" label="漏洞名称"  > 
                                </el-table-column>
                                <el-table-column   prop="risk" label="漏洞风险"  > 
                                </el-table-column>
                                <el-table-column   prop="number" label="出现次数"  > 
                                </el-table-column>
                                <el-table-column   prop="affectTargets" label="影响目标"  > 
                                </el-table-column>
                            </el-table>
                        </div>
                    </div>
                    <div>
                        <div class="commontitle 3part" style="" id="targetDetails" @click="foldingBtn('targetDetails', 'targetDetails')">
                                <span id="span-targetDetails">目标风险</span>
                                <span class="iconfont iconxialashixintop"> </span>
                            </div>
                            <div class="targetDetails" style="margin-top:25px;padding: 0 25px;"> 
                                <el-table  
                                    :data="data.targetDetails" 
                                    tooltip-effect="dark" 
                                    style=" margin-right: 25px;">
                                    <el-table-column   prop="target" label="测试目标"  > 
                                    </el-table-column>
                                    <el-table-column   prop="risk" label="目标风险"  > 
                                    </el-table-column>
                                    <el-table-column   prop="deadlyNumber" label="致命漏洞"  > 
                                    </el-table-column>
                                    <el-table-column   prop="highNumber" label="高危漏洞"  > 
                                    </el-table-column>
                                    <el-table-column   prop="middleNumber" label="中危漏洞"  > 
                                    </el-table-column>
                                    <el-table-column   prop="lowNumber" label="低危漏洞"  > 
                                    </el-table-column>
                                    <el-table-column   prop="vulStatus" label="漏洞状态"  > 
                                    </el-table-column>
                                </el-table>
                            </div>
                    </div>
                    <div  >
                        <div class="commontitle 4part" style="display: inline-block;margin-bottom: 25px;" id="vulDetails" @click="foldingBtn('vulDetails', 'vulDetails')">
                                <span id="span-vulDetails">漏洞详情</span>
                                <span class="iconfont iconxialashixintop"> </span>
                        </div> 
                        <div class="4_1box vulDetails" style="padding: 0 25px;" > 
                            <el-table
                                    id="file"
                                    ref="multipleTable"
                                    :data="data.vulDetails" 
                                    tooltip-effect="dark" 
                                    style=" margin-right: 25px;"
                                     >
                                    <el-table-column type="expand">
                                        <template slot-scope="scope">
                                            <div class="displaybox">
                                                <p><span class="loopname">漏洞类型：</span><span>{{ scope.row.type }}</span></p>
                                                <p><span class="loopname">漏洞编号：</span><span>{{ scope.row.cve }}</span></p>
                                                <p><span class="loopname">公开日期：</span><span>{{ scope.row.publishDate}}</span></p> 
                                                <p><span class="loopname">漏洞描述：</span><span>{{ scope.row.describe }}</span></p>
                                                <p ><span class="loopname">修复建议：</span><span>{{ scope.row.fix }}</span></p>
                                                <p><span class="loopname">影响范围：</span><span>{{ scope.row.affectRange }}</span></p>
                                                <p><span class="loopname">影响目标：</span><span>{{ scope.row.AffectTargets }}</span></p>
                                                <p><span class="loopname">参考链接：</span><span>{{ scope.row.link }}</span></p>
                                            </div> 
                                        </template> 
                                    </el-table-column>
                                    <el-table-column
                                        prop="vulName"
                                        label="漏洞名称"   width ="700">  
                                        <template slot-scope="scope">
                                            <div :class="[
                                                { 'highfont': scope.row.risk == '致命' },
                                                { 'middlefont': scope.row.risk == '高危' },
                                                { 'lowfont': scope.row.risk == '中危' },
                                                { 'infofont': scope.row.risk == '低危' }
                                                ]">{{ scope.row.vulName }}</div>
                                        </template>
                                    </el-table-column>
                                    <el-table-column
                                        prop="risk"
                                        label="漏洞风险" >  
                                        <template slot-scope="scope">
                                            <div >
                                                <span :class="[
                                                { 'highcolor': scope.row.risk == '致命' },
                                                { 'middlecolor': scope.row.risk == '高危' },
                                                { 'lowcolor': scope.row.risk == '中危' },
                                                { 'infocolor': scope.row.risk == '低危' }
                                                ]" class="spotclor2 "></span>
                                                {{ scope.row.risk }}</div>
                                        </template>
                                    </el-table-column>
                                    <el-table-column
                                        prop="vulStatus"
                                        label="漏洞状态" >  
                                        <template slot-scope="scope">
                                            <div  :class="[
                                                { 'highcolor': scope.row.vulStatus == '未能验证' },
                                                { 'middlecolor': scope.row.vulStatus == '验证失败' },
                                                { 'lowcolor': scope.row.vulStatus == '验证成功' },
                                                { 'infocolor': scope.row.vulStatus == '利用成功' },
                                                { 'lightbluecolor': scope.row.vulStatus == '待验证' }
                                                ]" class="spotclor3 ">
                                                {{ scope.row.vulStatus }}
                                            </div>
                                        </template>
                                    </el-table-column> 
                                    <el-table-column
                                        prop="number"
                                        label="出现次数" >  
                                    </el-table-column> 
                                </el-table> 
                        </div>  
                    </div>

                </div> 
            </div> 
		</div>
	</div>
</template>
<style lang="less"  >
@import 'style/export_vue.css';
</style>
<style lang="less" scoped> 
/deep/ .el-table::before {
    height: 1px !important;
}
/deep/ .el-table td.el-table__cell, 
/deep/.el-table th.el-table__cell.is-leaf{
  border-bottom: 1px solid #ebeef5 !important;
}
/deep/ .el-table__fixed-right::before, 
/deep/ .el-table__fixed::before{
  height:1px   !important;
}
 
.el-table th{
    padding-left: 0 !important;
}
.catalog ul li a {
    text-decoration:none
}
</style>
<script>
import $ from "jquery" 
import {report} from "@/api/report.js"; 
var echarts = require('echarts');  
export default {
    data () { 
        return {  
            report_id:this.$route.query.id, 
            type:this.$route.query.type,
            report_cover_title:'',
            report_cover_create_time:'',
            cataloglist:[],
            taskOverview:{
                targetStat:{},
                vulnVerify:{},
                vulnStat:{},
            },
            risk_total:0,
            data:{
                vulRisk:[],
                topVulRisk:[]
            },
            dev_title:'',
            allcatalogul :['taskOverview','taskStat','targetRisk','vulRisk','vulType','topVulRisk','targetDetails','vulDetails'],
            backgroundStyle: {
                background: `url(${require('../../assets/images/templatebg.png')})`
            }
        };
	}, 
	created() { 
	},
    mounted() {  
        this.getData();     
    },
    methods: { 
        async getData(){  
            const res = await report.downLoadfile({
                reportId:this.report_id, 
            })
            if (res.code === 200) {
                this.data = JSON.parse(res.data.content); 
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
                return;
            }  
            // 封皮
            this.report_cover_title = this.data.reportCover.title||''; 
            this.report_cover_create_time = this.data.reportCover.createTime; 
            // 目录 
            this.cataloglist = this.data.catalogParent; 
            // 任务概述
            this.taskOverview = this.data.taskOverview; 
            
            this.showTitlebyCatalogUl();
            this.$nextTick(()=>{
                this.targetschartpie(); //目标统计 
                this.risk_total = this.data.targetRisk.total;
            })

            var configJson = JSON.parse(res.data.configJson);
            if (configJson.cover) {
                this.backgroundStyle.background = `url(${configJson.cover})`;
            }
        },
        showTitlebyCatalogUl(){ //根据目录显示隐藏模块
            let _data = this.data.catalogParent;  
            let catalogarr = [];
            if (_data) {
                _data.forEach((item, index) => {
                    if (item.isShow) {  
                        catalogarr.push(item.id); 
                        $('#span-' + item.id).text(item.name)
                        if (item.catalog){ 
                            item.catalog.forEach((child, childIndex) => {
                                if (child.isShow) { 
                                    catalogarr.push(child.id); 
                                    $('#span-' + child.id).text(child.name)
                                } else {
                                    $('#' + child.id).addClass('notExist')
                                    $('.' + child.id).addClass('notExist')
                                    $('#' + child.id).hide()
                                    $('.' + child.id).hide()
                                }
                                
                            })
                        }
                    } else {
                        $('#' + item.id).hide()
                        $('.' + item.id).hide()
                        
                    }
                })
            } 
            // 不同
            let different = this.getNewArr(catalogarr,this.allcatalogul);

            different.forEach(item =>{
                if(item !=''){
                    $('#' + item).hide()
                    $('.' + item).hide() 
                } 
            }) 
            
        },
        getNewArr(a,b){
            const arr = [...a,...b];
            const newArr = arr.filter(item => {
                return !(a.includes(item) && b.includes(item));
            });
            return newArr;
        },
        targetschartpie(){ //目标统计  
            let _this = this;
            let target_stat = this.data.targetRisk;
            let _data = [
                {name: ' 高危目标      ',value:target_stat.highNumber,rate:target_stat.highNumberRate},
                {name: ' 中危目标      ',value:target_stat.MiddleNumber,rate:target_stat.MiddleNumberRate},
                {name: ' 低危目标      ',value:target_stat.lowNumber,rate:target_stat.lowNumberRate},
                {name: ' 安全目标      ',value:target_stat.safeNumber,rate:target_stat.safeNumberRate}
            ];
            console.log(_data)
            let dom = document.getElementById('targetschartpie');
       
            // 基于准备好的dom，初始化echarts实例
            var myChart = echarts.init(dom);
            // 绘制图表
            let target_chartsOptions = { 
                color:['#F87D7D','#F9B640','#4C7AE3','#15C53D'],
                tooltip: { 
                    trigger: 'item',
                    formatter: "{b}: {c} ({d}%)"
                },
                legend: {
                    type: 'scroll',
                    orient: 'vertical',
                    textStyle: {
                        color: '#484866'          // 图例文字颜色
                    },
                    right: 100,
                    top: 120,
                    bottom: 50,
                    itemGap: 35,
                    itemWidth:6,
                    itemHeight:6,
                    icon: "circle",  // 使用回调函数
                    formatter: function(name) {  
                        // var t_data = _data;
                        // var total = 0;
                        // var tarValue;
                        // for (var i = 0, l = t_data.length; i < l; i++) {
                        // total += t_data[i].value;
                        // if (t_data[i].name == name) {
                        //     tarValue = t_data[i].value;
                        // }
                        // }
                        // var p = ((tarValue / total) * 100).toFixed(2);
                        // return name + '         ' + tarValue + '              ' + p + "%"; 
                        var t_data = _data; 
                        var tarValue;
                        var t_rate='';
                        for (var i = 0, l = t_data.length; i < l; i++) { 
                            if (t_data[i].name == name) {
                                tarValue = t_data[i].value;
                                t_rate =  t_data[i].rate;
                            }
                        } 
                        return name + '         ' + tarValue + '              ' + t_rate; 
                    },
                    
                },          
                series: [
                    {
                        name:'',
                        type:'pie', 
                        radius: ['50%', '70%'],
                        hoverAnimation:false,
                        data: _data,
                        label: {
                            normal: {
                                formatter: function(params, ticket, callback) {
                                    return params.name + '：\n ('+params.percent+'%)';
                                },
                                position: 'inner',
                                show : false
                            },
                        },

                        tooltip: {
                            trigger: 'item',
                            formatter: "{b}: {c} ({d}%)"
                        },
                    }
                ]
            } 

            myChart.setOption(target_chartsOptions) 
            
            // $('.targetStatisticCount').text(target_stat.total)
        } ,
        exportPdf(){
            let obj = document.querySelector('#exportbutton');
            obj.style.display = 'none';
            window.print();

            obj.style.display = 'block';
        },
        foldingBtn(className, id) {
            if ($('#' + id + ' .iconfont').attr('class').includes('iconxialashixintop')) {
                $('#' + id + ' .iconfont').removeClass('iconxialashixintop')
                $('#' + id + ' .iconfont').addClass('iconxialashixinright')
            } else {
                $('#' + id + ' .iconfont').removeClass('iconxialashixinright')
                $('#' + id + ' .iconfont').addClass('iconxialashixintop')
            }
            $('.' + className + ':not(".notExist")').toggle()
        } 
    }
}
</script>
