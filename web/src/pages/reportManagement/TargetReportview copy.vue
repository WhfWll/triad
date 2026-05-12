<template>
    <div> 
        <div class="exportbutton" id="exportbutton">  
            <div  v-if="type=='pdf'"  >
                <el-button size="small" style="margin-right:8px" @click="exportPdf">导出PDF</el-button> 
            </div> 
        </div> 
        <div id="file">
            <div class="reportbox" style="padding-bottom: 50px;">
                <div class="boxtitlebg" >
                    <div class="largetitle">{{report_cover_title}}</div> 
                    <div class="smalltitle">{{report_cover_create_time}}</div>
                </div>
                <div class="firstpart">
                    <div class="reportcontent">
                        <div>
                            <div class="report-item">
                                <div>
                                    <div class="item_title" id="reportCatalog" >
                                        <span>报告目录</span>  
                                    </div>
                                </div>
                                <div class="catalog">
                                    <ul id="catalogUl" style="margin: 0 25px 10px 25px;"> 
                                        <div  v-for="(catalog,i) in cataloglist " :key="i">
                                            <li class="level1" >
                                                <span class="spanText"><a :href="'#'+catalog.id">{{ catalog.name }}</a></span>
                                                <span class="spanShengluehao"></span>
                                            </li>   
                                        </div> 
                                    </ul>
                                </div> 
                            </div> 
                        </div>
                    </div>
                </div>

                <div style="background: #fff;margin-top: 25px;">
                    <div class="commontitle 2_1part" style="margin-top:0;margin-bottom: 25px;width:145px;" 
                        id="targetOverview" @click="foldingBtn('targetOverview','targetOverview')">
                        <span id="span-targetOverview">报告摘要</span>
                        <span class="iconfont  iconxialashixintop"> </span>
                    </div>
                    <div class="targetOverview" style="padding: 0 25px;">
                        <table class="firsttable">
                            <tbody>
                            <tr>
                                <td class="nodarkback renwugyArea-tdLeft">测试目标：</td>
                                <td class="lightback renwugyArea-tdRight" id="task-name">{{ targetOverview.target }}</td>
                            </tr>
                            <tr>
                                <td class="darkback renwugyArea-tdLeft">风险等级：</td>
                                <td class="darkback renwugyArea-tdRight" id="task-risk">
                                    <span class="levelcolor " style="width: 8px;"
                                            :class="[
                                                {'highcolor':targetOverview.risk == '高危'},
                                                {'middlecolor':targetOverview.risk == '中危'},
                                                {'lowcolor':targetOverview.risk == '低危'},
                                                {'infocolor':targetOverview.risk == '安全' || targetOverview.risk == '未发现'}
                                            ]"></span>{{ targetOverview.risk }}</td>
                            </tr>
                            <tr>
                                <td class="nodarkback renwugyArea-tdLeft">漏洞分布：</td>
                                <td class="lightback renwugyArea-tdRight" id="loudongfenbu">
                                    <div><span>漏洞总数：</span><span id="bug-count">{{ targetOverview.vulnStat.total }}</span></div>
                                    <div>
                                        <span>致命漏洞：</span>
                                        <span id="high-bugCount" class="spanR">{{ targetOverview.vulnStat.deadlyNumber }}</span>
                                        <span>高危漏洞：</span>
                                        <span id="middle-bugCount"  class="spanR">{{ targetOverview.vulnStat.highNumber }}</span>
                                        <span>中危漏洞：</span>
                                        <span id="low-bugCount"  class="spanR">{{ targetOverview.vulnStat.middleNumber }}</span>
                                        <span>低危漏洞：</span>
                                        <span id="info-count">{{ targetOverview.vulnStat.lowNumber }}</span></div>
                                </td>
                            </tr>
                            <tr>
                                <td class="darkback renwugyArea-tdLeft">漏洞状态：</td>
                                <td class="darkback renwugyArea-tdRight">
                                    <div><span>验证成功：</span>
                                        <span id="validateSuccess-count"  class="spanR">{{targetOverview.vulnVerify.verifySuccess }}</span>
                                        <span>利用成功：</span>
                                        <span id="useSuccess-count" class="spanR">{{targetOverview.vulnVerify.useSuccess }}</span>
                                        <span>未验证：</span>
                                        <span id="unvalidate-count" class="spanR">{{targetOverview.vulnVerify.repairSuccess }}</span>
                                    </div>
                                </td>
                            </tr>
                            <tr>
                                <td class="nodarkback renwugyArea-tdLeft">测试时间：</td>
                                <td class="lightback renwugyArea-tdRight" id="test-time">{{ targetOverview.createDate }}</td>
                            </tr> 
                            </tbody>
                        </table>
                    </div>
                    <div class="commontitle 4part zichanxx" style="margin-bottom: 25px;" id="assetInfo" @click="foldingBtn('assetInfo','assetInfo')">
                        <span id="span-assetInfo">资产信息</span>
                        <span class="iconfont iconxialashixintop"> </span>
                    </div>
                    <div class="assetInfo "  style="padding: 0 25px;"> 
                        <table class="firsttable " style="width: 100%;">
                            <tbody> 
                                <tr>
                                    <td class="darkback renwugyArea-tdLeft" >组件/指纹</td>
                                    <td class="darkback renwugyArea-tdRight"  >{{ assetInfo.component }}</td>
                                </tr>
                                <tr>
                                    <td class="nodarkback renwugyArea-tdLeft"  >服务</td>
                                    <td class="lightback renwugyArea-tdRight"  >{{ assetInfo.service }}</td>
                                </tr>
                                <tr>
                                    <td class="darkback renwugyArea-tdLeft"  >IP/域名</td>
                                    <td class="darkback renwugyArea-tdRight"  >{{ assetInfo.ipOrUrl }}</td>
                                </tr>
                                <tr>
                                    <td class="nodarkback renwugyArea-tdLeft"  >操作系统</td>
                                    <td class="lightback renwugyArea-tdRight"   >{{ assetInfo.system }}</td>
                                </tr> 
                            </tbody>
                        </table>
                    </div>
                    <div class="commontitle 4part"  style="width:145px; margin-bottom:25px" id="vulInfo" @click="foldingBtn('vulInfo','vulInfo')">
                        <span id="span-vulInfo">漏洞信息</span>
                        <span class="iconfont iconfont iconxialashixintop"> </span>
                    </div>
                    <div class="spotdistance vulInfo" style="padding: 0 25px; padding-bottom: 25px;">  
                        <div class="4_1box loudongxiangqingArea jcczloudong"> 
                            <el-table 
                                ref="multipleTable"
                                :data="data.vulInfo" 
                                tooltip-effect="dark" 
                                style=" margin-right: 25px;" 
                                    >
                                <el-table-column type="expand">
                                    <template slot-scope="scope">
                                        <div class="displaybox">
                                            <p><span class="loopname">漏洞类型：</span><span>{{ scope.row.type }}</span></p>
                                            <p><span class="loopname">漏洞编号：</span><span>{{ scope.row.cve }}</span></p>
                                            <p><span class="loopname">披露日期：</span><span>{{ scope.row.publishDate}}</span></p> 
                                            <p><span class="loopname">漏洞描述：</span><span>{{ scope.row.describe }}</span></p>
                                            <p><span class="loopname">漏洞结果：</span><span class="res_style">{{ scope.row.res }}</span></p>
                                            <p ><span class="loopname">修复建议：</span><span>{{ scope.row.fix }}</span></p>
                                            <p><span class="loopname">影响范围：</span><span>{{ scope.row.affectRange }}</span></p>
                                            <p><span class="loopname">漏洞位置：</span><span>{{ scope.row.location }}</span></p>
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
.res_style{
    white-space: normal; 
    word-break: break-word;
}
.catalog ul li a {
    text-decoration:none
}
</style>

<script>
import $ from "jquery" 
import {report} from "@/api/report.js"; 
export default {
    name:'targetreportview',
    data () { 
        return {  
            report_id:this.$route.query.id, 
            type:this.$route.query.type,
            report_cover_title:'',
            report_cover_create_time:'',
            cataloglist:[],
            targetOverview:{
                vulnStat:{},
                vulnVerify:{}
            },
            assetInfo:{},
            data:{
                vulInfo:[]
            },
            allcatalogul:['targetOverview','assetInfo','vulInfo'], 
        }
    },
    mounted(){
        this.getData();   
    },
    methods:{
        async getData(){  
            const res = await report.downLoadfile({
                reportId:this.report_id, 
            })
            if (res.code === 200) {
                this.data = JSON.parse(res.data.content); 
                if (this.data && this.data.length > 0) {
                    //只保留this.data得第一项且是数组得话

                    this.data = this.data[0];
                }
                    this.report_cover_title = res.data.name||''; 
                 this.report_cover_create_time = res.data.createTime||''; 
            }else{
                this.$message({
                    message:res.msg,
                    type: 'error'
                });
                return;
            } 
            console.log(this.data)
            // 封皮
            
           
            // 目录 
            this.cataloglist = this.data?.catalogParent || []; 
            // 任务概述
            this.targetOverview = this.data?.targetOverview || {}; 
            this.assetInfo = this.data?.assetInfo || {};
            this.showTitlebyCatalogUl();
             
        },
        showTitlebyCatalogUl(){ //根据目录显示隐藏模块
            let _data = this.data.catalogParent || [];  
            let catalogarr = [];
            if (_data && _data.length > 0) {
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