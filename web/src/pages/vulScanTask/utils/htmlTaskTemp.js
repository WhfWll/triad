export const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0"> 
    <link href = "./css/bootstrap.css" rel="stylesheet" >
    <link href = "./css/report.css" rel="stylesheet" type="text/css" />
    <link href = "./css/iconfont.css" rel="stylesheet" >
    
    <style type="text/css">
            .colortable2 th {
                width: 10%;
            }
            .spotclor3 {
                display: inline-block;
                width: 60px;
                height: 20px;
                line-height: 20px;
                text-align: center;
                border-radius: 12px;
            }
            .lightbluecolor{
                color: #fff;
                background-color: #09C1F7;
            }
  /*          .el-table .spanR{
                margin-right: 60px;
                }*/
            .el-table .lightback span {
                float: initial;
                margin-left: 0;                
            }
            .el-table .darkback span {
                float: initial;
                margin-left: 0;
            }
            .firsttable .lightback span {
                display: inline-block;
                width: 85px;
            }
            .firsttable .darkback span {
                display: inline-block;
                width: 85px;
            }
            #task-risk span{
                display: inline-block;
                width: 8px!important;
            }
            .baogaozhaiyaoArea-tdLeft{
                width: 30%
            }
            .baogaozhaiyaoArea-tdRight{
                width: 70%!important;
            }
        </style>
</head> 
<body>
    <div class="bigreport">
        <div class="reportbox">
            <div class="boxtitlebg" id="boxtitlebg">
                <div class="largetitle"></div>
                <div class="middletitle"></div>
                <div class="smalltitle"></div>
            </div> 
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
                            </ul>
                        </div> 
                    </div>
                    <div class="systemdetection">
            
                    </div>
                </div>
            </div>
            <div class="thirdpart">
                <div class="task_title">
                    <div class="commontitle 2_1part" style="margin-bottom: 25px;"  id="taskOverview" onclick="foldingBtn('taskOverview', 'taskOverview')">
                        <span id="span-taskOverview">任务概述</span>
                        <span class="iconfont iconxialashixintop"> </span>
                    </div>
                    <div class="taskOverview">
                        <table class="el-table firsttable">
                            <tbody>
                                <tr>
                                    <td class="nodarkback baogaozhaiyaoArea-tdLeft">任务名称：</td>
                                    <td class="lightback baogaozhaiyaoArea-tdRight" id="task-name"></td>
                                    
                                </tr>
                                <tr>
                                    <td class="darkback baogaozhaiyaoArea-tdLeft ">任务风险：</td>
                                    <td class="darkback baogaozhaiyaoArea-tdRight redpot"  id="task-risk"></td>
                                </tr>
                                <tr>
                                    <td class="nodarkback baogaozhaiyaoArea-tdLeft">目标分布：</td>
                                    <td class="lightback baogaozhaiyaoArea-tdRight" id="task-plan">
                                        <div>
                                            <span>总目标：</span><span id="target-count" class="spanR"></span>
                                            <span>存活目标：</span><span id="alive-count" class="spanR"></span> 
                                        <div>
                                            <span>高危目标：</span><span id="high-targetCount" class="spanR"></span>
                                            <span>中危目标：</span><span id="middle-targetCount"  class="spanR"></span>
                                            <span>低危目标：</span><span id="low-targetCount"  class="spanR"></span>
                                            <span>安全目标：</span><span id="safe-targetCount"></span>
                                        </div>
                                    </td>
                                </tr>
                                <tr>
                                    <td class="darkback baogaozhaiyaoArea-tdLeft">漏洞分布：</td>
                                    <td class="darkback baogaozhaiyaoArea-tdRight" id="task-count">
                                        <div><span>漏洞总数：</span><span id="bug-count"></span></div>
                                        <div>
                                            <span>致命漏洞：</span><span id="high-bugCount" class="spanR"></span>
                                            <span>高危漏洞：</span><span id="middle-bugCount"  class="spanR"></span>
                                            <span>中危漏洞：</span><span id="low-bugCount"  class="spanR"></span>
                                            <span>低危漏洞：</span><span id="info-count"></span>
                                        </div>
                                    </td>
                                </tr>
                                <tr>
                                    <td class="nodarkback baogaozhaiyaoArea-tdLeft">漏洞验证：</td>
                                    <td class="lightback baogaozhaiyaoArea-tdRight">
                                        <div>
                                            <span>验证存在：</span><span id="validateSuccess-count"  class="spanR"></span>
                                            <span>利用成功：</span><span id="useSuccess-count" class="spanR"></span>
                                            <span>未验证：</span><span id="unvalidate-count" class="spanR"></span> 
                                        </div>
                                    </td>
                                </tr> 
                                <tr>
                                    <td class="nodarkback baogaozhaiyaoArea-tdLeft">任务场景：</td>
                                    <td class="lightback baogaozhaiyaoArea-tdRight" id="task-fangan"></td>
                                </tr>
                                <tr>
                                    <td class="darkback baogaozhaiyaoArea-tdLeft">测试时间：</td>
                                    <td class="darkback baogaozhaiyaoArea-tdRight" id="task-shijian"></td>
                                </tr> 
                            </tbody>
                        </table>
                    </div>


                    <div class="commontitle" id="taskStat"  onclick="foldingBtn('taskStat', 'taskStat')">
                        <span id="span-taskStat">信息统计</span>
                        <span class="iconfont iconxialashixintop"> </span>
                    </div>
                    <div class="taskStat">
                        <div class="commontitle width152 xinxitjArea spacialtitle"  id="targetRisk" style=" width:164px;">
                            <span id="span-targetRisk">目标风险统计</span>
                        </div>
                        <div class="targetRisk"> 
                            <div class="loopStatisticschart  " style="position: relative;">
                                <div class="totaltarget"><span>合计 : </span><span class="totalnumber targetStatisticCount" 
                                    id="target_Statistic_total" ></span><span>个目标</span></div>
                                <div class="el-row">
                                    <div id="targetschartpie" class="Statisticschart">
                                                                            
                                    </div>
                                </div>
                            </div>
                        </div>
                        <!-- ........... -->
                        <div class="commontitle width150 2_2part xinxitjArea" style="margin-bottom: 25px;width:164px;" id="vulRisk">
                            <span id="span-vulRisk">漏洞风险统计</span>
                        </div>
                        <div class="vulRisk loudongfxtj">
                            <table class="el-table nobordertable">
                                <tbody id="tbody-loudongfxtj">
                                <tr>
                                    <th class="darkback">风险类型</th>
                                    <th class="darkback">验证存在</th> 
                                    <th class="darkback">未验证</th>
                                    <th class="darkback">利用成功</th>
                                    <th class="darkback">漏洞总数</th>
                                    <th class="darkback">漏洞占比</th>
                                </tr>
                                </tbody>
                            </table>
                        </div> 
                        <!-- ........... -->
                        <div class="commontitle 2_5part xinxitjArea" style="margin-bottom: 25px; width:164px;" id="vulType">
                            <span id="span-vulType">漏洞类型统计</span>
                        </div>
                        <table class="el-table nobordertable xinxitjArea hasA vulType">
                            <tbody id="hignDangerBug">
                            <tr>
                                <th class="darkback">漏洞类型</th>
                                <th class="darkback">数量</th>
                                <th class="darkback">占比</th> 
                                <th class="darkback">影响目标数</th>
                            </tr>
                            </tbody>
                        </table>
                        <div class="commontitle 2_6part xinxitjArea" style="margin-bottom: 25px; width:164px;" id="topVulRisk">
                            <span id="span-topVulRisk">TOP危险漏洞</span>
                        </div>
                        <div class="topVulRisk ">
                            <table class="el-table nobordertable hasA">
                                <tbody id="mostBug">
                                <tr>
                                    <th class="darkback loopname">漏洞名称</th>
                                    <th class="darkback looprisk">漏洞风险</th> 
                                    <th class="darkback loopcount">出现次数</th> 
                                    <th class="darkback looptarget">影响目标</th>
                                </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <!-- ......................................................................................................... -->
                    <div class="commontitle 3part" style="" id="targetDetails" onclick="foldingBtn('targetDetails', 'targetDetails')">
                        <span id="span-targetDetails">目标详情</span>
                        <span class="iconfont iconxialashixintop"> </span>
                    </div>
                    <div class="targetDetails" style="margin-top:25px;">
                        <table class="el-table nobordertable colortable2">
                            <tbody id="goalDetail">
                            <tr>
                                <th class="thbg greybg cetarget">测试目标</th> 
                                <th class="thbg greybg goalrisk">目标风险</th>
                                <th class="thbg highcolor">致命漏洞</th>
                                <th class="thbg middlecolor">高危漏洞</th>
                                <th class="thbg lowcolor">中危漏洞</th>
                                <th class="thbg infocolor">低危漏洞</th>
                                <th class="thbg greybg servercount">漏洞状态</th>
                            </tr>
                            </tbody>
                        </table>
                    </div>
                     
                    <!-- ......................................................................................................... -->
                    <div class="commontitle 4part" style="display: inline-block;margin-bottom: 25px;" id="vulDetails"  onclick="foldingBtn('vulDetails', 'vulDetails')">
                            <span id="span-vulDetails">漏洞详情</span>
                            <span class="iconfont iconxialashixintop"> </span>
                    </div> 
                    <div class="4_1box vulDetails" style="padding: 0 25px;">
                        <table class="el-table nobordertable zhedie specialTable">
                            <tbody id="hadFoundBug">
                            <tr>                              
                                <th class="darkback th70">漏洞名称</th>
                                <th class="darkback th15">漏洞风险</th>
                                <th class="darkback th15">漏洞状态</th>
                                <th class="darkback th15">出现次数</th>
                            </tr>
                            </tbody>
                        </table>
                    </div>  

                </div>

            </div>
        </div>
    </div> 
</body> 

<script src="./js/echarts_min.js" ></script>
<script> 
    function foldingBtn(className, id) { 
        var idobj = document.getElementById(id);
        var classobj = document.getElementsByClassName(className);
        let spanObject = idobj.getElementsByClassName('iconfont');

        let span_obj = idobj.getElementsByTagName('span')[1]; 

        let obj_classlist  = span_obj.className;
        if(obj_classlist.includes('iconxialashixintop')){
            span_obj.classList.remove('iconxialashixintop');
            span_obj.classList.add('iconxialashixinright'); 

            classobj[0].style.display = 'none'; 


        }else{
            span_obj.classList.remove('iconxialashixinright');
            span_obj.classList.add('iconxialashixintop'); 
            classobj[0].style.display = 'block'; 
        } 
    }
    window.onload=function(){ 
        var allcatalogul = ['taskOverview','taskStat','targetRisk','vulRisk','vulType','topVulRisk','targetDetails','vulDetails']; 
        var data = %data%; 
        // 封面表皮
        document.getElementsByClassName('largetitle')[0].innerHTML = data.reportCover.title; 
        document.getElementsByClassName('smalltitle')[0].innerHTML = data.reportCover.createTime;
  
 
        //目录
        getCatalogUl();  
        //目标统计
        target_statistics(%chart_tmp%);
        //任务概述
        getTaskOverview();
        //漏洞风险统计
        getloudongfxtj();
        // 漏洞类型统计
        getvuln_type_tatistics();
        //top危险
        gettopVulRisk();
        //目标详情
        gettarget_detail();
        // 漏洞详情
        getvulDetails();  

        function getCatalogUl(){ //目录
            let _data = data.catalogParent;
            let ul = ''
            let startNum = 0
            let catalogarr = [];
            if (_data) {
                _data.forEach((item, index) => {
                    if (item.isShow) {
                        startNum++
                        catalogarr.push(item.id);
                        ul += '<li class="level1">'+
                            '<span class="spanText"><a href="#'+item.id+'">'+item.name+'</a></span>'+
                            '<span class="spanShengluehao"></span>'+
                            '</li>';
                        var obj = document.getElementById('span-'+item.id);
                        if(obj){
                            var spantitleText = obj.innerHTML;
                            obj.innerHTML =item.name;
                        } 
 
                        if (item.catalog){
                            let childStart = 0
                            item.catalog.forEach((child, childIndex) => { 
                                if (child.isShow) {
                                    catalogarr.push(child.id);
                                    childStart++
                                    ul += '<li class="level2">'+
                                        '<span class="spanText"><a href="#'+child.id+'">'+child.name+'</a></span>'+
                                        '<span class="spanShengluehao"></span>'+
                                        '</li>';
                                    var childobj = document.getElementById('span-'+child.id);
                                    if(childobj){
                                        var childtitleText = childobj.innerHTML; 
                                        childobj.innerHTML = child.name; 
                                    } 
                                } else {
                                    var childobj =  document.getElementById( child.id)
                                    var childobj_ = document.getElementsByClassName( child.id)[0];

                                    childobj.classList.add('notExist');
                                    childobj_.classList.add('notExist');

                                    childobj.style.display = 'none';
                                    childobj_.style.display = 'none'; 
                                }
                                
                            })
                        }
                    } else {
                        var obj = document.getElementById(item.id);
                        var obj_ = document.getElementsByClassName( item.id)[0];

                        obj.style.display = 'none';
                        obj_.style.display = 'none'; 
                    }
                })
            }
            document.getElementById('catalogUl').innerHTML = ul;

            var  different =  getNewArr(catalogarr,allcatalogul);
            different.forEach(item =>{
                if(item !=''){ 
                    document.getElementById(item).style.display='none';
                    var divobj = document.getElementsByClassName( item);
                    if(divobj && divobj.length > 0){
                        divobj[0].style.display='none';
                    }
                   
                } 
            })

        }
        function getNewArr(a,b){
            const arr = [...a,...b];
            const newArr = arr.filter(item => {
                return !(a.includes(item) && b.includes(item));
            });
            return newArr;
        }
       
        function target_statistics(charsOptionStr){ 
            let target_stat = data.targetRisk;
            let _data = [
                {name: ' 高危目标      ',value:target_stat.highNumber,rate:target_stat.highNumberRate},
                {name: ' 中危目标      ',value:target_stat.MiddleNumber,rate:target_stat.MiddleNumberRate},
                {name: ' 低危目标      ',value:target_stat.lowNumber,rate:target_stat.lowNumberRate},
                {name: ' 安全目标      ',value:target_stat.safeNumber,rate:target_stat.safeNumberRate}
            ] 
            charsOptionStr.legend.formatter = function(name) {  
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
            };
            var myChart = echarts.init(document.getElementById("targetschartpie")); 
            myChart.setOption(charsOptionStr); 

            let txt = data.targetRisk.total;
            document.getElementById("target_Statistic_total").innerHTML = txt;
        }
        function  getTaskOverview(){//任务概述
            let taskOverview = data.taskOverview;
            document.getElementById('task-name').innerHTML = taskOverview.taskName; 
            document.getElementById('task-fangan').innerHTML = taskOverview.templateName;  
            document.getElementById('task-shijian').innerHTML = taskOverview.date 
            let statusClass = ''
            switch (taskOverview.taskRiskStr) {
                case '高危':
                    statusClass="highcolor";
                    break;
                case '中危':
                    statusClass="middlecolor";
                    break;
                case '低危':
                    statusClass="lowcolor";
                    break;
                case '未发现':
                    statusClass="infocolor";
                    break;
                case '安全':
                    statusClass="infocolor";
                    break;
            }
            document.getElementById('task-risk').innerHTML = '<span class="levelcolor '+statusClass+'"></span>'+taskOverview.taskRiskStr
             
            document.getElementById('bug-count').innerHTML = taskOverview.vulnStat.total;
            document.getElementById('high-bugCount').innerHTML = taskOverview.vulnStat.deadlyNumber;
            document.getElementById('middle-bugCount').innerHTML = taskOverview.vulnStat.highNumber
            document.getElementById('low-bugCount').innerHTML = taskOverview.vulnStat.middleNumber
            document.getElementById('info-count').innerHTML = taskOverview.vulnStat.lowNumber
           
            document.getElementById('target-count').innerHTML = taskOverview.targetStat.total
            document.getElementById('alive-count').innerHTML = taskOverview.targetStat.liveTarget
            document.getElementById('high-targetCount').innerHTML = taskOverview.targetStat.HighTarget
            document.getElementById('middle-targetCount').innerHTML = taskOverview.targetStat.middleTarget
            document.getElementById('low-targetCount').innerHTML = taskOverview.targetStat.lowTarget
            document.getElementById('safe-targetCount').innerHTML = taskOverview.targetStat.safeTarget
           
            document.getElementById('validateSuccess-count').innerHTML = taskOverview.vulnVerify.verifySuccess
            document.getElementById('useSuccess-count').innerHTML = taskOverview.vulnVerify.useSuccess
            document.getElementById('unvalidate-count').innerHTML = taskOverview.vulnVerify.repairSuccess
           
        }

        function getloudongfxtj(){
            let trs = '';
            let _dt = data.vulRisk;
            if (_dt) {
                _dt.forEach(item => {
                    trs+='<tr><td>'+item.riskType
                        +'</td><td>'+item.verifySuccess
                        +'</td><td>'+item.useSuccess
                        +'</td><td>'+item.repairSuccess
                        +'</td><td>'+item.total
                        +'</td><td>'+item.percent+'</td></tr>'; 
                })
            }
            document.getElementById('tbody-loudongfxtj').innerHTML+=trs; 
        }
        function getvuln_type_tatistics(){
            let trs = ''
            let _dt = data.vulType;
            if (_dt) {
                _dt.forEach(item => {
                    trs+='<tr><td>'+item.vulnType
                        +'</td><td>'+item.total
                        +'</td><td>'+item.percent
                        +'</td><td>'+item.targetNumber
                        +'</td></tr>'; 
                })
            }
            document.getElementById('hignDangerBug').innerHTML+=trs;   
        }
        function gettopVulRisk(){
            let trs = ''
            let _dt = data.topVulRisk;
            if (_dt) {
                _dt.forEach(item => {
                    trs+='<tr><td>'+item.vulName
                        +'</td><td>'+item.risk
                        +'</td><td>'+item.number
                        +'</td><td>'+item.affectTargets
                        +'</td></tr>'; 
                    
                })
            }
            document.getElementById('mostBug').innerHTML+=trs;    
        }
        function gettarget_detail(){
            let trs = ''
            let _dt = data.targetDetails;
            if (_dt) {
                _dt.forEach(item => {
                    let statusClass = ''
                    switch (item.risk) {
                        case '高危':
                            statusClass="highcolor";
                            break;
                        case '中危':
                            statusClass="middlecolor";
                            break;
                        case '低危':
                            statusClass="lowcolor";
                            break;
                        case '信息':
                            statusClass="infocolor";
                            break;
                    }
                    let loopClass = ''
                    switch (item.vulStatus) {
                        case '未能验证':
                            loopClass="highcolor";
                            break;
                        case '验证失败':
                            loopClass="middlecolor";
                            break;
                        case '验证成功':
                            loopClass="lowcolor";
                            break;
                        case '待验证':
                            loopClass="lightbluecolor";
                            break;   
                        case '利用成功':
                            loopClass="infocolor";
                            break;
                    }
                    trs+='<tr><td class=" jumpa"> '+item.target+'</td>'+ 
                        '<td><span class="spotclor2 '+statusClass+'"></span>'+item.risk+'</td>'+
                        '<td>'+item.deadlyNumber+'</td>'+
                        '<td>'+item.highNumber+'</td>'+
                        '<td>'+item.middleNumber+'</td>'+
                        '<td>'+item.lowNumber+'</td>'+
                        '<td><span class="spotclor3 '+loopClass+'">'+item.vulStatus+'</span></td></tr>'; 
                })
            }
            document.getElementById('goalDetail').innerHTML+=trs;     
        }  
        function getvulDetails () {
                let trs = ''
                if (data.vulDetails) {
                    data.vulDetails.forEach((item, index) => {
                        let fontClass = ''
                        let statusClass = ''
                        switch (item.risk) {
                            case '致命':
                                fontClass="highfont";
                                statusClass="highcolor";
                                break;
                            case '高危':
                                fontClass="middlefont";
                                statusClass="middlecolor";
                                break;
                            case '中危':
                                fontClass="lowfont";
                                statusClass="lowcolor";
                                break;
                            case '低危':
                                fontClass="infofont";
                                statusClass="infocolor";
                                break;
                        }
                        let loopClass = ''
                        switch (item.vulStatus) {
                            case '未能验证':
                                loopClass="highcolor";
                                break;
                            case '验证失败':
                                loopClass="middlecolor";
                                break;
                            case '验证成功':
                                loopClass="lowcolor";
                                break;
                            case '待验证':
                                loopClass="lightbluecolor";
                                break;   
                            case '利用成功':
                                loopClass="infocolor";
                                break;
                        }
                        let _index = index+1;
                        // onclick="shousuo('+_index+', this, "foundBug"
                        trs+= '<tr data-attr_index = '+_index+' data-attr-class = "foundBug" >'+
                            '<td class="'+fontClass+'" style="width:65%;cursor: pointer;"><span class="iconfont iconunfoldcross expandBtn'+_index+'" > </span>'+item.vulName+'</td>'+
                            '<td><span class="spotclor2 '+statusClass+'"></span>'+item.risk+'</td>'+
                            '<td><span class="spotclor3 '+loopClass+'">'+item.vulStatus+'</span></td>'+
                            '<td style="padding-left:40px">'+item.number+'</td></tr>';
                    })
                }
                var div = document.getElementById('hadFoundBug');
                div.innerHTML+=trs;     
                div.addEventListener('click',function(e){ 
                    if(e.target && e.target.parentNode.nodeName =='TR'){
                      
                        let trobj =  e.target.parentNode;
                        let index = trobj.getAttribute('data-attr_index');
                        let classname = trobj.getAttribute('data-attr-class');
                        shousuo(index,trobj,classname);
                    }
                    if(e.target && e.target.parentNode.nodeName =='TD'){
                      
                        let trobj =  e.target.parentNode.parentNode;
                        let index = trobj.getAttribute('data-attr_index');
                        let classname = trobj.getAttribute('data-attr-class');
                        shousuo(index,trobj,classname);
                    }
                })

            }


        function shousuo(i, obj, which){ 
            var _class = 'expandBtn' + i;
            let index= parseInt(i);
            var div = document.getElementsByClassName(_class)[0];
            if(div.classList.contains('iconunfoldcross')){ 
                var nextDiv = obj.nextSibling;
                // if(nextDiv || !nextDiv.classList.contains(trContent)){
                  
                    let detail = ''; 
                        detail = data.vulDetails[index-1]; 
                    var  $el = document.createElement('tr');
                    $el.classList.add('trContent');
                    $el.classList.add('trContent'+index);

                    var _html = '<td colspan="4">'+
                        '<div class="displaybox">'+
                        '<p><span class="loopname">漏洞类型：</span><span>'+detail.type+'</span></p>'+
                        '<p><span class="loopname">漏洞编号：</span><span>'+detail.cve+'</span></p>'+
                        '<p><span class="loopname">公开日期：</span><span>'+detail.publishDate+'</span></p>'+
                        '<p><span class="loopname">漏洞描述：</span><span>'+detail.describe+'</span></p>'; 
                        _html+= '<p ><span class="loopname">修复建议：</span><span>'+detail.fix+'</span></p>';
                   
                    _html+= '<p><span class="loopname">影响范围：</span><span>'+detail.affectRange+'</span></p>'+
                        '<p><span class="loopname">影响目标：</span><span class="bluetarget">'+detail.AffectTargets+'</span></p>'+
                        '<p><span class="loopname">参考链接：</span><span>'+detail.link+'</span></p>'+
                        '</div></td>';

                    $el.innerHTML = _html;

                        insertAfter( $el,obj); 
                // }  
                let _obj = document.getElementsByClassName('expandBtn' + index)[0];
                _obj.classList.add('iconfoldcross');
                _obj.classList.remove('iconunfoldcross'); 
            }else{
                let _obj = document.getElementsByClassName('expandBtn' + index)[0];
                _obj.classList.add('iconunfoldcross');
                _obj.classList.remove('iconfoldcross'); 
            }
            let dev = document.querySelector('.trContent' + index);
            if(dev.style.display !='table-row'){
                dev.style.display  = 'table-row'
            }else{
                dev.style.display  = 'none'
            }
        }
       

        /*
        函数说明：用来在某个元素之后插入元素
        参数说明:
        newElement:表示被插入的元素
        targetElement:表示目标元素
        */
        function insertAfter(newElement, targetElement){
            //得到父节点
            var parent = targetElement.parentNode;
            //如果目标节点已经是最后一个元素，那么直接添加即可
            if(targetElement === parent.lastChild){
                parent.appendChild(newElement);
            }else 
                {
                //否则，在当前节点的下一个节点之前添加
                parent.insertBefore(newElement,targetElement.nextSibling);
            }
        }
    };
   
 
</script>
</html>`;