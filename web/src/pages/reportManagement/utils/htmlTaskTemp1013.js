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
            .hyZheDie{
                width: 98%!important;
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
        
           <div   >
            <div class="commontitle 2_1part hyZheDie" style="margin-bottom: 25px;"  id="bluebar2" onclick="foldingBtn('bluebar2', 'bluebar2')">
                <span id="hyTitle">原因是这个被js赋值了-应该是url</span>
                <span class="iconfont iconxialashixintop"> </span>
             </div>
        <div class="bluebar2 ">
            <div class="reportcontent " style="color:#FFF">
                <div>
                <div class="wozijiasdass ">
                    
                </div>
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
        //  var data =[
        //     {"reportId":145,"targetOverview":{"target":"http://www.chinaqkin_20230927170036"},"reportCover":{"title":"目标测试报告数组","createTime":"2023-09-25","backgroundImg":""},"catalogParent":[{"name":"1. 任务概述111","id":"taskOverview","isShow":true,"catalog":null},{"name":"2. 信息统计111","id":"taskStat","isShow":true,"catalog":[{"name":"2.1 目标风险统计","id":"targetRisk","isShow":true,"catalog":null},{"name":"2.2 漏洞风险统计","id":"vulRisk","isShow":true,"catalog":null},{"name":"2.3 漏洞类型统计","id":"vulType","isShow":true,"catalog":null},{"name":"2.4 Top 危险漏洞","id":"topVulRisk","isShow":true,"catalog":null}]},{"name":"3. 目标详情","id":"targetDetails","isShow":true,"catalog":null},{"name":"4. 漏洞详情","id":"vulDetails","isShow":true,"catalog":null}],"taskOverview":{"taskName":"127.0.0.1192.168.0._20230925150111","taskRiskStr":"高危","targetStat":{"total":3,"liveTarget":3,"HighTarget":2,"middleTarget":0,"lowTarget":0,"safeTarget":1},"vulnStat":{"total":0,"deadlyNumber":2,"highNumber":0,"middleNumber":0,"lowNumber":2},"vulnVerify":{"verifySuccess":4,"useSuccess":0,"repairSuccess":0},"templateName":"全场景测试","date":"2023-09-25 15:01至2023-09-25 15:03"},"targetRisk":{"total":3,"highNumber":2,"highNumberRate":"66.66%","MiddleNumber":0,"MiddleNumberRate":"0%","lowNumber":0,"lowNumberRate":"0%","safeNumber":1,"safeNumberRate":"33.33%"},"vulRisk":[{"riskType":"致命漏洞","verifySuccess":2,"repairSuccess":0,"useSuccess":0,"total":2,"percent":"50%"},{"riskType":"高危漏洞","verifySuccess":0,"repairSuccess":0,"useSuccess":0,"total":0,"percent":"0%"},{"riskType":"中危漏洞","verifySuccess":0,"repairSuccess":0,"useSuccess":0,"total":0,"percent":"0%"},{"riskType":"低危漏洞","verifySuccess":2,"repairSuccess":0,"useSuccess":0,"total":2,"percent":"50%"},{"riskType":"统计","verifySuccess":4,"repairSuccess":0,"useSuccess":0,"total":4,"percent":"100%"}],"vulType":[{"vulnType":"弱密码","total":2,"percent":"50%","targetNumber":2},{"vulnType":"其他","total":2,"percent":"50%","targetNumber":1}],"topVulRisk":[{"vulName":"ssh弱口令","risk":"致命","number":2,"affectTargets":"127.0.0.1,192.168.0.61"},{"vulName":"\\"无HTTP重定向\\"","risk":"低危","number":1,"affectTargets":"127.0.0.1"},{"vulName":"Unencrypted connection的中文翻译为“非加密连接”。","risk":"低危","number":1,"affectTargets":"127.0.0.1"}],"targetDetails":[{"target":"127.0.0.1","risk":"高危","deadlyNumber":1,"highNumber":0,"middleNumber":0,"lowNumber":2,"vulStatus":"验证存在"},{"target":"192.168.0.123","risk":"安全","deadlyNumber":0,"highNumber":0,"middleNumber":0,"lowNumber":0,"vulStatus":""},{"target":"192.168.0.61","risk":"高危","deadlyNumber":1,"highNumber":0,"middleNumber":0,"lowNumber":0,"vulStatus":"验证存在"}],"vulDetails":[{"vulName":"ssh弱口令","risk":"致命","vulStatus":"验证存在","number":2,"type":"弱密码","cve":"","publishDate":"2022-05-16 00:00:00","describe":"弱口令没有严格和准确的定义,通常认为容易被别人(它们有可能对你很了解)猜测或被破解工具破解的口令均为弱口令。","res":"{\\"ip\\":\\"127.0.0.1\\",\\"password\\":\\"admin123321\\",\\"service\\":\\"ssh\\",\\"target\\":\\"127.0.0.1:22\\",\\"username\\":\\"dogs\\"}","fix":"1. SSH修改默认端口 2. SSH防御暴力破解用户账号 3. SSH设置PGP登录 4. Iptables设置阈值防止暴力破解","affectRange":"","AffectTargets":"127.0.0.1,192.168.0.61","location":"","link":""},{"vulName":"\\"无HTTP重定向\\"","risk":"低危","vulStatus":"验证存在","number":1,"type":"其他","cve":"","publishDate":"","describe":"检测到您的Web应用程序使用HTTP协议，但不会自动重定向用户到HTTPS。","res":"{\\"pocname\\":\\"3f6a8a0e-07f2-af81-54ff-61020299caeb\\"}","fix":"该漏洞描述了一个常见的web应用程序安全问题，即应用程序使用HTTP协议而不是HTTPS协议。HTTP协议是一种不安全的传输协议，未加密的数据包会通过网络传递。这使得攻击者可以轻易地读取、修改和窃取这些数据，并对用户隐私造成威胁。因此，为了确保应用程序使用最高级别的安全措施，就需要使用HTTPS协议代替HTTP协议。HTTPS协议使用传输层安全（TLS）或安全套接字层（SSL）协议来加密数据传输。这样，只有受信任的接收方可以读取和使用传输数据。\\n\\n为了修复这个漏洞，应用程序开发人员需要将HTTP协议更改为HTTPS协议。他们可以在应用程序配置文件中为使用HTTPS协议的路径添加重定向规则，这将确保用户始终通过HTTPS协议访问应用程序。此外，开发人员应可以使用SSL证书或其他安全措施来保护传输过程中的数据。这可以通过将SSL证书绑定到应用程序域名上来实现，或使用基于云的安全服务来实现。这些安全措施将确保用户隐私的安全，并使攻击者无法查看或窃取数据。在应用程序修复这个漏洞后，网站管理员应该定期进行安全测试和漏洞扫描，以确保应用程序继续保持安全并且满足最高级别的安全标准。","affectRange":"","AffectTargets":"127.0.0.1","location":"","link":""},{"vulName":"Unencrypted connection的中文翻译为“非加密连接”。","risk":"低危","vulStatus":"验证存在","number":1,"type":"其他","cve":"","publishDate":"","describe":"这个扫描目标连接是未加密的。一个潜在的攻击者可以拦截和修改从这个站点发送和接收的数据。","res":"{\\"pocname\\":\\"bca221d1-8581-3375-4097-66b0048ed088\\"}","fix":"该漏洞是因为扫描目标连接未加密，导致攻击者可以拦截和修改该站点发送和接收的数据，因此需要采取一定的措施进行修复。通常，可以通过SSL（安全套接字层）协议来解决这个问题，SSL协议是一种加密协议，可以保护网络传输过程中的数据安全。\\n\\n对于站点管理者来说，他们需要考虑启用SSL证书以支持HTTPS协议，HTTPS协议相比HTTP协议是基于SSL加密的，在传输过程中可以将数据进行加密传输，确保数据的安全性。并且站点管理者可以通过定期的安全检查以及更新站点的安全软件，及时发现和修复可能存在的漏洞，保证站点的安全性。\\n\\n对于用户来说，他们需要注意防范网络钓鱼以及恶意软件的攻击，尽可能选择使用SSL加密的站点，同时多使用杀毒软件和防火墙等安全软件，增强自身网络安全意识。在使用公共网络时，为了避免个人账号和密码等敏感信息的泄露，建议不要随意使用未知的公共网络，并且尽可能不在公共网络中登录个人账号。\\n\\n综上所述，修复该漏洞的关键在于采取措施来加强站点的安全保障体系，包括启用SSL证书、定期更新安全软件、及时发现和修复可能存在的漏洞等。同时，用户也需要更加注重自身的网络安全，选择SSL加密的站点、使用安全软件，避免在公共网络中登录敏感信息，从而保障网络的安全。","affectRange":"","AffectTargets":"127.0.0.1","location":"","link":""}]},
        //     {"reportId":146,"targetOverview":{"target":"目标测试报告数组222","createTime":"2023-09-25","backgroundImg":""},"catalogParent":[{"name":"1. 任务概述2222","id":"taskOverview","isShow":true,"catalog":null},{"name":"2. 信息统计222","id":"taskStat","isShow":true,"catalog":[{"name":"2.1 目标风险统计","id":"targetRisk","isShow":true,"catalog":null},{"name":"2.2 漏洞风险统计","id":"vulRisk","isShow":true,"catalog":null},{"name":"2.3 漏洞类型统计","id":"vulType","isShow":true,"catalog":null},{"name":"2.4 Top 危险漏洞","id":"topVulRisk","isShow":true,"catalog":null}]},{"name":"3. 目标详情","id":"targetDetails","isShow":true,"catalog":null},{"name":"4. 漏洞详情","id":"vulDetails","isShow":true,"catalog":null}],"taskOverview":{"taskName":"127.0.0.1192.168.0._20230925150111","taskRiskStr":"高危","targetStat":{"total":3,"liveTarget":3,"HighTarget":2,"middleTarget":0,"lowTarget":0,"safeTarget":1},"vulnStat":{"total":0,"deadlyNumber":2,"highNumber":0,"middleNumber":0,"lowNumber":2},"vulnVerify":{"verifySuccess":4,"useSuccess":0,"repairSuccess":0},"templateName":"全场景测试","date":"2023-09-25 15:01至2023-09-25 15:03"},"targetRisk":{"total":3,"highNumber":2,"highNumberRate":"66.66%","MiddleNumber":0,"MiddleNumberRate":"0%","lowNumber":0,"lowNumberRate":"0%","safeNumber":1,"safeNumberRate":"33.33%"},"vulRisk":[{"riskType":"致命漏洞","verifySuccess":2,"repairSuccess":0,"useSuccess":0,"total":2,"percent":"50%"},{"riskType":"高危漏洞","verifySuccess":0,"repairSuccess":0,"useSuccess":0,"total":0,"percent":"0%"},{"riskType":"中危漏洞","verifySuccess":0,"repairSuccess":0,"useSuccess":0,"total":0,"percent":"0%"},{"riskType":"低危漏洞","verifySuccess":2,"repairSuccess":0,"useSuccess":0,"total":2,"percent":"50%"},{"riskType":"统计","verifySuccess":4,"repairSuccess":0,"useSuccess":0,"total":4,"percent":"100%"}],"vulType":[{"vulnType":"弱密码","total":2,"percent":"50%","targetNumber":2},{"vulnType":"其他","total":2,"percent":"50%","targetNumber":1}],"topVulRisk":[{"vulName":"ssh弱口令","risk":"致命","number":2,"affectTargets":"127.0.0.1,192.168.0.61"},{"vulName":"\\"无HTTP重定向\\"","risk":"低危","number":1,"affectTargets":"127.0.0.1"},{"vulName":"Unencrypted connection的中文翻译为“非加密连接”。","risk":"低危","number":1,"affectTargets":"127.0.0.1"}],"targetDetails":[{"target":"127.0.0.1","risk":"高危","deadlyNumber":1,"highNumber":0,"middleNumber":0,"lowNumber":2,"vulStatus":"验证存在"},{"target":"192.168.0.123","risk":"安全","deadlyNumber":0,"highNumber":0,"middleNumber":0,"lowNumber":0,"vulStatus":""},{"target":"192.168.0.61","risk":"高危","deadlyNumber":1,"highNumber":0,"middleNumber":0,"lowNumber":0,"vulStatus":"验证存在"}],"vulDetails":[{"vulName":"ssh弱口令","risk":"致命","vulStatus":"验证存在","number":2,"type":"弱密码","cve":"","publishDate":"2022-05-16 00:00:00","describe":"弱口令没有严格和准确的定义,通常认为容易被别人(它们有可能对你很了解)猜测或被破解工具破解的口令均为弱口令。","res":"{\\"ip\\":\\"127.0.0.1\\",\\"password\\":\\"admin123321\\",\\"service\\":\\"ssh\\",\\"target\\":\\"127.0.0.1:22\\",\\"username\\":\\"dogs\\"}","fix":"1. SSH修改默认端口 2. SSH防御暴力破解用户账号 3. SSH设置PGP登录 4. Iptables设置阈值防止暴力破解","affectRange":"","AffectTargets":"127.0.0.1,192.168.0.61","location":"","link":""},{"vulName":"\\"无HTTP重定向\\"","risk":"低危","vulStatus":"验证存在","number":1,"type":"其他","cve":"","publishDate":"","describe":"检测到您的Web应用程序使用HTTP协议，但不会自动重定向用户到HTTPS。","res":"{\\"pocname\\":\\"3f6a8a0e-07f2-af81-54ff-61020299caeb\\"}","fix":"该漏洞描述了一个常见的web应用程序安全问题，即应用程序使用HTTP协议而不是HTTPS协议。HTTP协议是一种不安全的传输协议，未加密的数据包会通过网络传递。这使得攻击者可以轻易地读取、修改和窃取这些数据，并对用户隐私造成威胁。因此，为了确保应用程序使用最高级别的安全措施，就需要使用HTTPS协议代替HTTP协议。HTTPS协议使用传输层安全（TLS）或安全套接字层（SSL）协议来加密数据传输。这样，只有受信任的接收方可以读取和使用传输数据。\\n\\n为了修复这个漏洞，应用程序开发人员需要将HTTP协议更改为HTTPS协议。他们可以在应用程序配置文件中为使用HTTPS协议的路径添加重定向规则，这将确保用户始终通过HTTPS协议访问应用程序。此外，开发人员应可以使用SSL证书或其他安全措施来保护传输过程中的数据。这可以通过将SSL证书绑定到应用程序域名上来实现，或使用基于云的安全服务来实现。这些安全措施将确保用户隐私的安全，并使攻击者无法查看或窃取数据。在应用程序修复这个漏洞后，网站管理员应该定期进行安全测试和漏洞扫描，以确保应用程序继续保持安全并且满足最高级别的安全标准。","affectRange":"","AffectTargets":"127.0.0.1","location":"","link":""},{"vulName":"Unencrypted connection的中文翻译为“非加密连接”。","risk":"低危","vulStatus":"验证存在","number":1,"type":"其他","cve":"","publishDate":"","describe":"这个扫描目标连接是未加密的。一个潜在的攻击者可以拦截和修改从这个站点发送和接收的数据。","res":"{\\"pocname\\":\\"bca221d1-8581-3375-4097-66b0048ed088\\"}","fix":"该漏洞是因为扫描目标连接未加密，导致攻击者可以拦截和修改该站点发送和接收的数据，因此需要采取一定的措施进行修复。通常，可以通过SSL（安全套接字层）协议来解决这个问题，SSL协议是一种加密协议，可以保护网络传输过程中的数据安全。\\n\\n对于站点管理者来说，他们需要考虑启用SSL证书以支持HTTPS协议，HTTPS协议相比HTTP协议是基于SSL加密的，在传输过程中可以将数据进行加密传输，确保数据的安全性。并且站点管理者可以通过定期的安全检查以及更新站点的安全软件，及时发现和修复可能存在的漏洞，保证站点的安全性。\\n\\n对于用户来说，他们需要注意防范网络钓鱼以及恶意软件的攻击，尽可能选择使用SSL加密的站点，同时多使用杀毒软件和防火墙等安全软件，增强自身网络安全意识。在使用公共网络时，为了避免个人账号和密码等敏感信息的泄露，建议不要随意使用未知的公共网络，并且尽可能不在公共网络中登录个人账号。\\n\\n综上所述，修复该漏洞的关键在于采取措施来加强站点的安全保障体系，包括启用SSL证书、定期更新安全软件、及时发现和修复可能存在的漏洞等。同时，用户也需要更加注重自身的网络安全，选择SSL加密的站点、使用安全软件，避免在公共网络中登录敏感信息，从而保障网络的安全。","affectRange":"","AffectTargets":"127.0.0.1","location":"","link":""}]},
        //     {"reportId":147,"targetOverview":{"target":"目标测试报告数组333","createTime":"2023-09-25","backgroundImg":""},"catalogParent":[{"name":"1. 任务概述333","id":"taskOverview","isShow":true,"catalog":null},{"name":"2. 信息统计333","id":"taskStat","isShow":true,"catalog":[{"name":"2.1 目标风险统计","id":"targetRisk","isShow":true,"catalog":null},{"name":"2.2 漏洞风险统计","id":"vulRisk","isShow":true,"catalog":null},{"name":"2.3 漏洞类型统计","id":"vulType","isShow":true,"catalog":null},{"name":"2.4 Top 危险漏洞","id":"topVulRisk","isShow":true,"catalog":null}]},{"name":"3. 目标详情","id":"targetDetails","isShow":true,"catalog":null},{"name":"4. 漏洞详情","id":"vulDetails","isShow":true,"catalog":null}],"taskOverview":{"taskName":"127.0.0.1192.168.0._20230925150111","taskRiskStr":"高危","targetStat":{"total":3,"liveTarget":3,"HighTarget":2,"middleTarget":0,"lowTarget":0,"safeTarget":1},"vulnStat":{"total":0,"deadlyNumber":2,"highNumber":0,"middleNumber":0,"lowNumber":2},"vulnVerify":{"verifySuccess":4,"useSuccess":0,"repairSuccess":0},"templateName":"全场景测试","date":"2023-09-25 15:01至2023-09-25 15:03"},"targetRisk":{"total":3,"highNumber":2,"highNumberRate":"66.66%","MiddleNumber":0,"MiddleNumberRate":"0%","lowNumber":0,"lowNumberRate":"0%","safeNumber":1,"safeNumberRate":"33.33%"},"vulRisk":[{"riskType":"致命漏洞","verifySuccess":2,"repairSuccess":0,"useSuccess":0,"total":2,"percent":"50%"},{"riskType":"高危漏洞","verifySuccess":0,"repairSuccess":0,"useSuccess":0,"total":0,"percent":"0%"},{"riskType":"中危漏洞","verifySuccess":0,"repairSuccess":0,"useSuccess":0,"total":0,"percent":"0%"},{"riskType":"低危漏洞","verifySuccess":2,"repairSuccess":0,"useSuccess":0,"total":2,"percent":"50%"},{"riskType":"统计","verifySuccess":4,"repairSuccess":0,"useSuccess":0,"total":4,"percent":"100%"}],"vulType":[{"vulnType":"弱密码","total":2,"percent":"50%","targetNumber":2},{"vulnType":"其他","total":2,"percent":"50%","targetNumber":1}],"topVulRisk":[{"vulName":"ssh弱口令","risk":"致命","number":2,"affectTargets":"127.0.0.1,192.168.0.61"},{"vulName":"\\"无HTTP重定向\\"","risk":"低危","number":1,"affectTargets":"127.0.0.1"},{"vulName":"Unencrypted connection的中文翻译为“非加密连接”。","risk":"低危","number":1,"affectTargets":"127.0.0.1"}],"targetDetails":[{"target":"127.0.0.1","risk":"高危","deadlyNumber":1,"highNumber":0,"middleNumber":0,"lowNumber":2,"vulStatus":"验证存在"},{"target":"192.168.0.123","risk":"安全","deadlyNumber":0,"highNumber":0,"middleNumber":0,"lowNumber":0,"vulStatus":""},{"target":"192.168.0.61","risk":"高危","deadlyNumber":1,"highNumber":0,"middleNumber":0,"lowNumber":0,"vulStatus":"验证存在"}],"vulDetails":[{"vulName":"ssh弱口令","risk":"致命","vulStatus":"验证存在","number":2,"type":"弱密码","cve":"","publishDate":"2022-05-16 00:00:00","describe":"弱口令没有严格和准确的定义,通常认为容易被别人(它们有可能对你很了解)猜测或被破解工具破解的口令均为弱口令。","res":"{\\"ip\\":\\"127.0.0.1\\",\\"password\\":\\"admin123321\\",\\"service\\":\\"ssh\\",\\"target\\":\\"127.0.0.1:22\\",\\"username\\":\\"dogs\\"}","fix":"1. SSH修改默认端口 2. SSH防御暴力破解用户账号 3. SSH设置PGP登录 4. Iptables设置阈值防止暴力破解","affectRange":"","AffectTargets":"127.0.0.1,192.168.0.61","location":"","link":""},{"vulName":"\\"无HTTP重定向\\"","risk":"低危","vulStatus":"验证存在","number":1,"type":"其他","cve":"","publishDate":"","describe":"检测到您的Web应用程序使用HTTP协议，但不会自动重定向用户到HTTPS。","res":"{\\"pocname\\":\\"3f6a8a0e-07f2-af81-54ff-61020299caeb\\"}","fix":"该漏洞描述了一个常见的web应用程序安全问题，即应用程序使用HTTP协议而不是HTTPS协议。HTTP协议是一种不安全的传输协议，未加密的数据包会通过网络传递。这使得攻击者可以轻易地读取、修改和窃取这些数据，并对用户隐私造成威胁。因此，为了确保应用程序使用最高级别的安全措施，就需要使用HTTPS协议代替HTTP协议。HTTPS协议使用传输层安全（TLS）或安全套接字层（SSL）协议来加密数据传输。这样，只有受信任的接收方可以读取和使用传输数据。\\n\\n为了修复这个漏洞，应用程序开发人员需要将HTTP协议更改为HTTPS协议。他们可以在应用程序配置文件中为使用HTTPS协议的路径添加重定向规则，这将确保用户始终通过HTTPS协议访问应用程序。此外，开发人员应可以使用SSL证书或其他安全措施来保护传输过程中的数据。这可以通过将SSL证书绑定到应用程序域名上来实现，或使用基于云的安全服务来实现。这些安全措施将确保用户隐私的安全，并使攻击者无法查看或窃取数据。在应用程序修复这个漏洞后，网站管理员应该定期进行安全测试和漏洞扫描，以确保应用程序继续保持安全并且满足最高级别的安全标准。","affectRange":"","AffectTargets":"127.0.0.1","location":"","link":""},{"vulName":"Unencrypted connection的中文翻译为“非加密连接”。","risk":"低危","vulStatus":"验证存在","number":1,"type":"其他","cve":"","publishDate":"","describe":"这个扫描目标连接是未加密的。一个潜在的攻击者可以拦截和修改从这个站点发送和接收的数据。","res":"{\\"pocname\\":\\"bca221d1-8581-3375-4097-66b0048ed088\\"}","fix":"该漏洞是因为扫描目标连接未加密，导致攻击者可以拦截和修改该站点发送和接收的数据，因此需要采取一定的措施进行修复。通常，可以通过SSL（安全套接字层）协议来解决这个问题，SSL协议是一种加密协议，可以保护网络传输过程中的数据安全。\\n\\n对于站点管理者来说，他们需要考虑启用SSL证书以支持HTTPS协议，HTTPS协议相比HTTP协议是基于SSL加密的，在传输过程中可以将数据进行加密传输，确保数据的安全性。并且站点管理者可以通过定期的安全检查以及更新站点的安全软件，及时发现和修复可能存在的漏洞，保证站点的安全性。\\n\\n对于用户来说，他们需要注意防范网络钓鱼以及恶意软件的攻击，尽可能选择使用SSL加密的站点，同时多使用杀毒软件和防火墙等安全软件，增强自身网络安全意识。在使用公共网络时，为了避免个人账号和密码等敏感信息的泄露，建议不要随意使用未知的公共网络，并且尽可能不在公共网络中登录个人账号。\\n\\n综上所述，修复该漏洞的关键在于采取措施来加强站点的安全保障体系，包括启用SSL证书、定期更新安全软件、及时发现和修复可能存在的漏洞等。同时，用户也需要更加注重自身的网络安全，选择SSL加密的站点、使用安全软件，避免在公共网络中登录敏感信息，从而保障网络的安全。","affectRange":"","AffectTargets":"127.0.0.1","location":"","link":""}]}
        // ]; 
       console.log(data,'data!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!');
        // 封面表皮
        document.getElementsByClassName('largetitle')[0].innerHTML = data[0].reportCover.title; 
        document.getElementsByClassName('smalltitle')[0].innerHTML = data[0].reportCover.createTime;

        setTimeout(()=>{
            data.forEach((it,id)=>{
                if(id== 2){
                    var hyTitle = document.querySelector('#hyTitle');
                    hyTitle.innerHTML = it.targetOverview.target || '报告3';
                }else if(id== 0){
                    var hyTitle = document.querySelector('#bluebar23 #hyTitle');
                    hyTitle.innerHTML = it.targetOverview.target || '报告1';
                    console.log('hyTitle',hyTitle); 
                }
                else if(id== 1){
                    var hyTitle = document.querySelector('#bluebar13 #hyTitle');
                    hyTitle.innerHTML = it.targetOverview.target || '报告2';
                    console.log('hyTitle',hyTitle); 
                }

                
            })
           },1000)


        data.forEach((item,index)=>{
           if(index > 0 ){
            var originalDiv = document.querySelector('.bluebar2');
            var copiedDiv = originalDiv.cloneNode(true);
            copiedDiv.className = 'bluebar'+index +3;
            originalDiv.insertAdjacentElement('afterend', copiedDiv);
            console.log(item,'item++++++++++++');
            // 加一个折叠
            var hyZheDie = document.querySelector('.hyZheDie');
            var hyZheDieCopy = hyZheDie.cloneNode(true);
            hyZheDieCopy.id = 'bluebar'+index +3;
            hyZheDie.insertAdjacentElement('afterend', hyZheDieCopy);
            hyZheDieCopy.onclick = () => {
                console.log('点击了折叠',copiedDiv.className);
                foldingBtn( copiedDiv.className, copiedDiv.className)
            }
            
           }
            //目录
        getCatalogUl(item.catalogParent);  
             //目标统计
        target_statistics(%chart_tmp%,item.targetRisk);
             //任务概述 
        getTaskOverview(item.taskOverview);
            //漏洞风险统计
        getloudongfxtj(item.vulRisk);
            // 漏洞类型统计
        getvuln_type_tatistics(item.vulType);
            //top危险
        gettopVulRisk(item.topVulRisk); 
        //目标详情
        gettarget_detail(item.targetDetails);   
        // 漏洞详情
        getvulDetails(item.vulDetails);  
        })
 
       
       
       
       
       
       
       
       

        function getCatalogUl(catalogParent){ //目录
            let _data = catalogParent;
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
       
        function target_statistics(charsOptionStr,targetRisk){ 
            let target_stat = targetRisk;
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

            let txt = targetRisk.total || 0;
            document.getElementById("target_Statistic_total").innerHTML = txt;
        }
        function  getTaskOverview(taskOverview){//任务概述
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

        function getloudongfxtj(vulRisk){
            let trs = '';
            let _dt = vulRisk;
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
        function getvuln_type_tatistics(vulType){
            let trs = ''
            let _dt = vulType;
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
        function gettopVulRisk(topVulRisk){
            let trs = ''
            let _dt = topVulRisk;
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
        function gettarget_detail(targetDetails){
            let trs = ''
            let _dt = targetDetails;
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
        function getvulDetails (vulDetails) {
                let trs = ''
                if (vulDetails) {
                    vulDetails.forEach((item, index) => {
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