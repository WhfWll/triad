export const htmlTargetTemplate = `<!DOCTYPE html>
<html  lang="en"> 
<head>
    <title>目标报告</title>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0"> 
    <link href = "./css/bootstrap.css" rel="stylesheet" >
    <link href = "./css/report.css" rel="stylesheet" type="text/css" />
    <link href = "./css/iconfont.css" rel="stylesheet" >
    
<style type="text/css">
.hyZheDie{
    margin-top:20px !important;
    width: 1000px !important;
    margin-bottom: 0px;
}
    .el-table .spanR{
            margin-right: 60px;
            }
        .el-table .lightback span {
            float: initial;
            margin-left: 0;
        }
        .el-table .darkback span {
            float: initial;
            margin-left: 0;
        }
        .renwugyArea-tdLeft{
                width: 30%
        }
        .renwugyArea-tdRight{
            width: 70%!important;
        }
        .unValidateBg {
            color: #fff;
            background: #F35F28;
        }
            .validateSuccessBg {
            color: #fff;
            background: #4C7AE3;
            }
            .validateFailBg {
            color: #fff;
            background: #FFB700;
            }
            .useSuccessBg {
            color: #fff;
            background: #15C53D;
            }
            .waitValidate {
            color: #fff;
            background: #09C1F7;
            }
            .statusTd{
            padding: 2px 6px;
            border-radius: 12px;
            font-size: 13px;
            }
            .displaybox{
                width:1100px;
            }
            .displaybox >.div1 {
                margin: 0px;
                width: 100%;
                overflow: hidden;
            }
            .displaybox .div1.gongjiDiv{
                display: flex;
                align-items: center;
            padding-bottom: 16px;
            }
            .displaybox .div1 .gongjilian{
            width: 54%;
            float:left;
            }
            .displaybox .div1 .gongjilian div{
            width: 400px;
            height: 72px;
            margin: 0 auto;
            margin-top: 16px;
            text-align: center;
            }
            .displaybox .div1 .gongjilian div.whiteBg{
            background-color: #FFF;
            padding-top: 15px;
            box-sizing: border-box;
            border-radius: 8px;
            border: 1px solid #E8E8F5;
            }
            .displaybox .div1 .gongjilian div.arrowDiv{
            height: 70px;
            margin-top: 0;
            }
            .displaybox .div1 .gongjilian div.arrowDiv img{
            width:12px;
            margin:  22px 0;
            }
            .displaybox .div1 .gongjilian div span:nth-child(1) {
            height: 18px;
            line-height: 18px;
            font-size: 13px;
            font-weight: 500;
            width: 100%;
            display: block;
            }
            .displaybox .div1 .gongjilian div span:nth-child(2) {
            height: 18px;
            line-height: 18px;
            margin-top: 8px;
            width: 100%;
            font-size: 13px;
            display: block;
            color: #A1A5B4;
            }
            .displaybox .div1:nth-child(odd) {
                background-color: #fff;
            }
            
            
</style>
</head>

<body>
    <div class="bigreport" id="top">
        <div class="reportbox" id="articleTop">
            <div class="boxtitlebg" id="boxtitlebg">
                <div class="largetitle"></div>
                <div class="middletitle"></div>
                <div class="smalltitle"></div>
            </div>
         <div class="commontitle 2_1part hyZheDie" style="margin-top:0;margin-bottom: 25px;width:145px;" 
            id="part1" onclick="changeIcon('part1')">
                <span  id="hyTitle">…</span>
                <span class="iconfont  iconxialashixintop"> </span>
        </div>
         <div class="part1">
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
                </div>
            </div>
            <div style="background: #fff;margin-top: 25px;">
                <div class="commontitle 2_1part" style="margin-top:0;margin-bottom: 25px;width:145px;" 
                id="targetOverview" onclick="changeIcon('targetOverview')">
                    <span id="span-targetOverview">报告摘要</span>
                    <span class="iconfont  iconxialashixintop"> </span>
                </div>
                <div class="targetOverview renwugyArea">
                    <table class="el-table">
                        <tbody>
                        <tr>
                            <td class="nodarkback renwugyArea-tdLeft">测试目标：</td>
                            <td class="lightback renwugyArea-tdRight" id="task-name"></td>
                        </tr>
                        <tr>
                            <td class="darkback renwugyArea-tdLeft">风险等级：</td>
                            <td class="darkback renwugyArea-tdRight" id="task-risk"></td>
                        </tr>
                        <tr>
                            <td class="nodarkback renwugyArea-tdLeft">漏洞分布：</td>
                            <td class="lightback renwugyArea-tdRight" id="loudongfenbu">
                                <div><span>漏洞总数：</span><span id="bug-count"></span></div>
                                <div>
                                    <span>致命漏洞：</span>
                                    <span id="high-bugCount" class="spanR"></span>
                                    <span>高危漏洞：</span>
                                    <span id="middle-bugCount"  class="spanR"></span>
                                    <span>中危漏洞：</span>
                                    <span id="low-bugCount"  class="spanR"></span>
                                    <span>低危漏洞：</span>
                                    <span id="info-count"></span></div>
                            </td>
                        </tr>
                        <tr>
                            <td class="darkback renwugyArea-tdLeft">漏洞状态：</td>
                            <td class="darkback renwugyArea-tdRight">
                                <div><span>验证成功：</span>
                                    <span id="validateSuccess-count"  class="spanR"></span>
                                    <span>利用成功：</span>
                                    <span id="useSuccess-count" class="spanR"></span>
                                    <span>未验证：</span>
                                    <span id="unvalidate-count" class="spanR"></span>
                                    
                            </td>
                        </tr>
                        <tr>
                            <td class="nodarkback renwugyArea-tdLeft">测试时间：</td>
                            <td class="lightback renwugyArea-tdRight" id="test-time"></td>
                        </tr> 
                        </tbody>
                    </table>
                </div>
                <div class="commontitle 4part zichanxx" style="margin-bottom: 25px;" id="assetInfo" onclick="changeIcon('assetInfo')">
                    <span id="span-assetInfo">资产信息</span>
                    <span class="iconfont iconxialashixintop"> </span>
                </div>
                <div class="assetInfo zichanxxArea">
                    <table class="el-table tableLayoutFixed">
                        <tbody> 
                        <tr>
                            <td class="darkback" style="width: 30%">组件/指纹</td>
                            <td class="darkback" id="zichan-zujian" style="width: 70%"></td>
                        </tr>
                        <tr>
                            <td class="nodarkback" style="width: 30%">服务</td>
                            <td class="lightback" id="zichan-fuwu" style="width: 70%"></td>
                        </tr>
                        <tr>
                            <td class="darkback " style="width: 30%">IP/域名</td>
                            <td class="darkback" id="zichan-ip" style="width: 70%"></td>
                        </tr>
                        <tr>
                            <td class="nodarkback " style="width: 30%">操作系统</td>
                            <td class="lightback" id="zichan-caozuo" style="width: 70%"></td>
                        </tr> 
                        </tbody>
                    </table>
                </div>
                <div class="commontitle 4part"  style="width:145px; margin-bottom:25px" id="vulInfo" onclick="changeIcon('vulInfo')">
                    <span id="span-vulInfo">漏洞信息</span>
                    <span class="iconfont iconfont iconxialashixintop"> </span>
                </div>
                <div class="spotdistance vulInfo" style="padding-bottom:25px">  
                    <div class="4_1box loudongxiangqingArea jcczloudong">
                        <table class="el-table nobordertable zhedie  ">
                            <tbody id="vuln_existTbody">
                                <tr>
                                    <th class="darkback th70">漏洞名称</th>
                                    <th class="darkback ">漏洞风险</th>
                                    <th class="darkback">漏洞状态</th>
                                </tr>
                            </tbody>
                        </table>
                    </div> 
                </div>  
            </div>
         </div>
        </div>
        <div class="bottombg2">
             <div class="bottombg"></div>
        </div>
    </div>


</body> 
<script src="./js/echarts_min.js"></script>
<script> 
    window.onload = function () {
        var allcatalogul = ['targetOverview','assetInfo','vulInfo']; 
        var data = %data%;  
        var real_catalog_arr = [];
        // 判断data的类型
        // console.log('data---SSS444444444444444444SSS',Object.prototype.toString.call(data));
        if(Object.prototype.toString.call(data) === '[object Object]'){
            let arr = [];
            arr.push(data);
            data = arr;
            console.log('data---SSSSSS',data);
        }
        // console.log('data---目标！！！！！！！！！！',data)
        // 封面表皮
        document.getElementsByClassName('largetitle')[0].innerHTML = data[0].reportCover.title; 
        document.getElementsByClassName('smalltitle')[0].innerHTML = data[0].reportCover.createTime;
       
            // setTimeout(()=>{
            //     data.forEach((it,id)=>{
            //         // 设置准备的title

            //         // if(id== 2){
            //         //     var hyTitle = document.querySelector('#hyTitle');
            //         //     hyTitle.innerHTML = it.targetOverview.target;
            //         // }else if(id== 0){
            //         //     var hyTitle = document.querySelector('#part113 #hyTitle');
            //         //     hyTitle.innerHTML = it.targetOverview.target;
            //         //     console.log('hyTitle',hyTitle); 
            //         // }
            //         // else if(id== 1){
            //         //     // var hyTitle = document.querySelector('#part123 #hyTitle');
            //         //     // hyTitle.innerHTML = it.targetOverview.target;
            //         //     // console.log('hyTitle',hyTitle); 
            //         // }

                    
            //     })
            //    },1000)
        
        data.forEach((item,index)=>{
           if( index == 0){
            //设置第一个title
            var hyTitle = document.querySelector('#part1 #hyTitle');
            var taskNameInput = document.querySelector('.part1 #task-name'); 
           setTimeout(()=>{
            console.log('taskNameInput---真他吗奇怪',taskNameInput); 
            var taskNameValue = taskNameInput.innerHTML;
            console.log('taskNameValue',taskNameValue);
            hyTitle.innerHTML = taskNameValue || '目标报告'+index;
           },1000)
           }
            if(index > 0 ){
           
             var originalDiv = document.querySelector('.part1');
             var copiedDiv = originalDiv.cloneNode(true);
             copiedDiv.className = 'part1'+index +3;
             originalDiv.insertAdjacentElement('afterend', copiedDiv);
             console.log(item,'item++++++++++++');

           
             // 加一个折叠
             var hyZheDie = document.querySelector('.hyZheDie');
             var hyZheDieCopy = hyZheDie.cloneNode(true);
             hyZheDieCopy.id = 'part1'+index +3;
             hyZheDie.insertAdjacentElement('afterend', hyZheDieCopy);
             // 设置折叠的title
                
                var hyTitle = document.querySelector('#part1'+index +3+' #hyTitle');
                var taskNameInput = document.querySelector('.part1'+index +3+' #task-name');  
                var taskNameValue = taskNameInput.innerHTML;
                console.log('taskNameValue',taskNameValue);
                hyTitle.innerHTML = taskNameValue || '目标报告'+index;
                // hyTitle.innerHTML = item.targetOverview.target || '目标报告'+index;
                
             //先隐藏所有
             document.querySelector('.'+copiedDiv.className).style.display = 'none';
             document.querySelector('.'+originalDiv.className).style.display = 'none';


             hyZheDieCopy.onclick = () => {
                 console.log('点击了折叠-hyZheDieCopy',copiedDiv.className);






                 if( document.querySelector('.'+copiedDiv.className).style.display == 'none'){
                   
                 // 获取父元素
                 var parent = document.querySelector('.reportbox');
 
                 // 获取父元素下所有类名中包含"part1"的子元素
                 var children = parent.querySelectorAll('[class*=part1]');
 
                 // 将子元素设置为display: none
                 for (var i = 0; i < children.length; i++) {
                 var child = children[i];
                 child.style.display = 'none';
                 }



                 document.querySelector('.'+copiedDiv.className).style.display = 'block';
                 }else{
                    
                 // 获取父元素
                 var parent = document.querySelector('.reportbox');
 
                 // 获取父元素下所有类名中包含"part1"的子元素
                 var children = parent.querySelectorAll('[class*=part1]');
 
                 // 将子元素设置为display: none
                 for (var i = 0; i < children.length; i++) {
                 var child = children[i];
                 child.style.display = 'none';
                 }



                 }

                //  changeIcon( copiedDiv.className)
             }
             hyZheDie.onclick = () => {
                 console.log('点击了折叠-hyZheDie',originalDiv.className);



                

                    if( document.querySelector('.'+originalDiv.className).style.display == 'none'){
                        // 获取父元素
                        var parent = document.querySelector('.reportbox');

                        // 获取父元素下所有类名中包含"part1"的子元素
                        var children = parent.querySelectorAll('[class*=part1]');

                        // 将子元素设置为display: none
                        for (var i = 0; i < children.length; i++) {
                        var child = children[i];
                        child.style.display = 'none';
                        }

                        document.querySelector('.'+originalDiv.className).style.display = 'block';
                    }else{
                        // 获取父元素
                        var parent = document.querySelector('.reportbox');
        
                        // 获取父元素下所有类名中包含"part1"的子元素
                        var children = parent.querySelectorAll('[class*=part1]');
        
                        // 将子元素设置为display: none
                        for (var i = 0; i < children.length; i++) {
                        var child = children[i];
                        child.style.display = 'none';
                        }
        
                    }
             }
             
            }

            getCatalogUl(item.catalogParent); //创建目录
            getTargetOverview(item.targetOverview) //报告摘要 
            getZichanxxData(item.assetInfo); //资产
            jcczloudong(item.vulInfo) //存在漏洞
         })


       
  

           
         
        // 创建目录
        function getCatalogUl(catalogParent){ 
            let ul = ''
            let startNum = 0;
            let catalogarr = [];
            if (catalogParent) {
                catalogParent.forEach((item, index) => {
                    if (item.isShow) {
                        startNum++
                        catalogarr.push(item.id);
                        //  目录列表
                        ul += \`<li class="level1">
                                    <span class="spanText"><a href="#\`+item.id+\`">\`+item.name+\`</a></span>
                                    <span class="spanShengluehao"></span>
                                </li>\`;  
                        let obj = document.getElementById('span-'+item.id);
                        if(obj){
                            let text = obj.innerHTML;
                            obj.innerHTML = item.name;
                        } 
                    } else { 
                        document.getElementById( item.id ).style.display = 'none';
                        document.querySelector(item.id).style.display = 'none';
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
       
        //报告摘要
        function getTargetOverview(target_summary) {
            if (target_summary) {
                document.getElementById('task-name').innerHTML = target_summary.target;
                let statusClass = ''
                switch (target_summary.risk) {
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
                document.getElementById('task-risk').innerHTML  = '<span class="levelcolor '+statusClass+'"></span>'+ target_summary.risk;

                document.getElementById('bug-count').innerHTML = target_summary.vulnStat.total;
                document.getElementById('high-bugCount').innerHTML = target_summary.vulnStat.deadlyNumber;
                document.getElementById('middle-bugCount').innerHTML = target_summary.vulnStat.highNumber;
                document.getElementById('low-bugCount').innerHTML = target_summary.vulnStat.middleNumber;
                document.getElementById('info-count').innerHTML = target_summary.vulnStat.lowNumber;
 
               
                document.getElementById('validateSuccess-count').innerHTML = target_summary.vulnVerify.verifySuccess;
                document.getElementById('useSuccess-count').innerHTML = target_summary.vulnVerify.useSuccess;
                document.getElementById('unvalidate-count').innerHTML = target_summary.vulnVerify.repairSuccess;

                document.getElementById('test-time').innerHTML = target_summary.createDate; 
                
            }
            
        }
        

        // 资产
        function getZichanxxData(ziChan) {
            if (ziChan) {
                // document.getElementById('zichan-yingyong').innerHTML = ziChan.application;
                document.getElementById('zichan-zujian').innerHTML = ziChan.component;
                document.getElementById('zichan-fuwu').innerHTML = ziChan.service;
                document.getElementById('zichan-ip').innerHTML = ziChan.ipOrUrl;
                document.getElementById('zichan-caozuo').innerHTML = ziChan.system; 
            }
        }
        //  漏洞
        function jcczloudong(vulInfo) {
            console.log('vulInfo!!!!!!!!',vulInfo);
            let trs = ''
            if (vulInfo) {
                vulInfo.forEach((item, index) => {
                    let fontClass = ''
                    let riskClass = ''
                    switch (item.risk) {
                        case '致命':
                            fontClass="highfont";
                            riskClass="highcolor";
                            break;
                        case '高危':
                            fontClass="middlefont";
                            riskClass="middlecolor";
                            break;
                        case '中危':
                            fontClass="lowfont";
                            riskClass="lowcolor";
                            break;
                        case '低危':
                            fontClass="infofont";
                            riskClass="infocolor";
                            break;
                    }
                    let bgClass = ''
                    switch (item.vulStatus) {
                        case '未能验证':
                            bgClass="unValidateBg";
                            break;
                        case '验证成功':
                            bgClass="validateSuccessBg";
                            break;
                        case '验证失败':
                            bgClass="validateFailBg";
                            break;
                        case '利用成功':
                            bgClass="useSuccessBg";
                            break;
                        case '待验证':
                            bgClass="waitValidate";
                            break;
                    }
                    trs += \` <tr data-attr_index="\`+index+\`" >
                                <td class="\`+fontClass+\`" style="width:70%;cursor: pointer;">
                                    <span class="iconfont iconunfoldcross foundBugs_\`+index+\`"  style="cursor: pointer;"> </span>\`+item.vulName+\`
                                </td>
                                <td><span class="spotclor2 \`+riskClass+\`"></span>\`+item.risk+\`</td>
                                <td><span class="\`+bgClass+\` statusTd">\`+item.vulStatus+\`</span></td>
                            </tr>\`;
                })
            }


            var div = document.getElementById('vuln_existTbody');
            div.innerHTML+=trs;    
            console.log('trs---------',trs); 
            div.addEventListener('click',function(e){ 
                if(e.target && e.target.parentNode.nodeName =='TR'){
                    
                    let trobj =  e.target.parentNode;
                    let index = trobj.getAttribute('data-attr_index'); 
                    shousuo(index,trobj,vulInfo);
                }
                if(e.target && e.target.parentNode.nodeName =='TD'){
                    
                    let trobj =  e.target.parentNode.parentNode;
                    let index = trobj.getAttribute('data-attr_index'); 
                    shousuo(index,trobj,vulInfo);
                }
            })
        }
        function shousuo(index, obj,vulInfo){ 
            var _class = 'foundBugs_' + index;
            var div = document.getElementsByClassName(_class)[0];
            if(div.classList.contains('iconunfoldcross')){ 
                var nextDiv = obj.nextSibling;
                console.log('index!!!!!!!!!!!!!!!!',index);
                let detail =  {}
                if(vulInfo&&vulInfo.length>0&&vulInfo[index]){
                    detail = vulInfo[index];
                }else{
                    detail = {};
                }
                
                var res = highlightedRequest(detail.verMsg.request, detail.verMsg.payload) 
                var rep = highlightedRequest(detail.verMsg.response, detail.verMsg.payload_success_flag) 
 
                
                var  $el = document.createElement('tr');
                $el.classList.add('trContent');
                $el.classList.add('trContent_'+index+'vuln');

                let gongjilian = '';
                detail.attackChain && detail.attackChain.forEach((item,i)=>{
                    if(i!=0){
                        gongjilian+='<div class="arrowDiv"  >'+
                                    '<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABQAAABECAYAAACBBxX6AAAAAXNSR0IArs4c6QAAAERlWElmTU0AKgAAAAgAAYdpAAQAAAABAAAAGgAAAAAAA6ABAAMAAAABAAEAAKACAAQAAAABAAAAFKADAAQAAAABAAAARAAAAACKzj82AAABHElEQVRYCe2X2w3CMAxFS2ENZoEBYKF+8cVCLACzdBDoqZQqbR3XKZF4yJYqUfte58gxH62qT8SpaZ88lrNriyhH4w1zpvXXWt9D9Xq//5+i4v9O0fdQvSvfQ3U85Yq+h+osfQ/V8ZQr+h6qs/Q9VMdTruh7qM7S91Adj6m4QXVu2nv3MXwwORKirtHjdt0f+0vZ1tUloTOnQ4+eENc7lIGOPsPahBNI5kbsHQhpsoYypqPHQMhLfBLvlph6RoQ0yKGc0uEfEZKYnkguFZJ2RojZQinR4Z0RkpROJh9HSiMSYtQoU3T4REIKKYKlWpIQo0Sp0eFJElKUKKUc2hAqIaKYcokOvUqIICaKf1NbHVDyWBrsLKJiZJbDppoXhFhd4t1MDAAAAAAASUVORK5CYII=" />'+
                                    '</div>'
                    }
                    
                    gongjilian +=  '<div class="whiteBg"><span>'+item.name+'</span><span>'+item.value +'</span></div>' ; 
                }) 

                let _html = '<td colspan="4">'+
                    '<div class="displaybox" >'+
                    '<p><span class="loopname">漏洞类型：</span><span>'+detail.type+'</span></p>'+
                    '<p><span class="loopname">漏洞编号：</span><span>'+detail.cve+'</span></p>'+
                    '<p><span class="loopname">披露日期：</span><span>'+detail.publishDate+'</span></p>'+
                    '<p><span class="loopname">漏洞描述：</span><span>'+detail.describe+'</span></p>'+
                    '<p><span class="loopname">漏洞结果：</span><span>'+detail.res+'</span></p>';
                    
                _html+='<div><span class="loopname">修复建议：</span><span>'+detail.fix+'</span></div>';
              
                _html+=  '<p><span class="loopname">影响范围：</span><span>'+detail.affectRange+'</span></p>'+
                    '<p><span class="loopname">漏洞位置：</span><span class="bluetarget">'+detail.location+'</span></p>'+
                    '<p><span class="loopname">参考链接：</span><span>'+detail.link+'</span></p>'+ 
                    '<p><span class="loopname">payload：</span><span>'+detail.verMsg.payload+'</span></p>'+ 
                    '<p><span class="loopname">请求报文：</span><span>'+res+'</span></p>'+ 
                    '<p><span class="loopname">响应报文：</span><span>'+rep+'</span></p>'+ 
                    '</div></td>';

                $el.innerHTML = _html ;
                    insertAfter( $el,obj); 
               
                let _obj = document.getElementsByClassName(_class)[0];
                _obj.classList.add('iconfoldcross');
                _obj.classList.remove('iconunfoldcross'); 
            }else{
                let _obj = document.getElementsByClassName(_class)[0];
                _obj.classList.add('iconunfoldcross');
                _obj.classList.remove('iconfoldcross'); 
            }
            let dev = document.querySelector('.trContent_'+index+'vuln');
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

    }  

    function highlightedRequest(request,payload) { 
        let safeResult = request;
        if(payload == ''){
            return safeResult
        }

        let regex  = new RegExp('('+payload+')', 'g');
        let replacestr=  safeResult.replace(regex, '<span style="color: red;font-weight:700">$1</span>');
        return replacestr
    }
    function changeIcon(dom){ 
        let child = document.getElementById(dom).childNodes;
        let childrenarr = [];
        child.forEach(item =>{
            if(item.nodeName == 'SPAN'){
                childrenarr.push(item)
            }
        })
        let children = childrenarr[1];
        // const children  = child[0];
        if (children.classList.contains('iconxialashixintop')) {
            children.classList.remove('iconxialashixintop')
            children.classList.add('iconxialashixinright')
        } else {
            children.classList.remove('iconxialashixinright')
            children.classList.add('iconxialashixintop')
        } 
 
        let dev = document.getElementsByClassName(dom);
        for(var i=0;i<dev.length;i++){
            let _obj = dev[i];

            if( _obj.style.display =='' ||_obj.style.display == 'block' ){
                _obj.style.display  = 'none'
            }else{
                _obj.style.display  = 'block'
            } 
            
        } 
    }
</script>
</html>
`;