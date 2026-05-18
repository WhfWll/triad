import Vue from 'vue'
// import { createRouter, createWebHistory } from 'vue-router'
import VueRouter from 'vue-router'

// var Login = () => import('@/views/Login') 
// var Login2 = () => import('@/views/Login2.vue') 

let Login = '';
if(process.env.VUE_APP_FREELOGIN>=5){
  Login = () => import('@/views/Login2.vue');
}else{
  Login = () => import('@/views/Login');
}


var Home = () => import('@/views/Home')
var Index = () => import('@/views/Index')
//用户管理
var usermanagement = () => import('@/pages/user/Usermanagementnew')
var usergroup = () => import('@/pages/user/usergroup')
var groupmanagement = () => import('@/pages/user/groupmanagement')
var xray = () => import('@/pages/othner/x-ray')
var burpsuite = () => import('@/pages/othner/burpsuite')
var detectionTask = () => import('@/pages/othner/detectionTask')
//系统管理
var systemsetting = () => import('@/pages/system/Systemsettingnew')
var systemtip = () => import('@/pages/system/systemtip')
var node = () => import('@/pages/system/Node.vue')

 //工具管理
var vulnerability = () => import('@/pages/tool/VulnerabilityLibrary') //漏洞库
var auxiliarytool = () => import('@/pages/tool/Auxiliarytool') //辅助工具
var fingerprint = () => import('@/pages/tool/FingerprintLib') //指纹库
var systemsetting = () => import('@/pages/system/Systemsetting')
var log = () => import('@/pages/log/Log')
var logconfig = () => import('@/pages/log/Logconfig')
Vue.use(VueRouter);
// 解决重复路由问题
const originalPush = VueRouter.prototype.push
 
VueRouter.prototype.push = function push(location) {
  return originalPush.call(this, location).catch(err => err)
}
const routes = [
  {
    path: '/',
    redirect: '/login',
  },
  {
    path: '/index',
    name: 'index',
    component: Home, 
    children: [ 
      {
        path: '/index',
        name: 'index',
        component: () => import('@/views/Index.vue') ,
      }, 
    ]
  },
  {
    path: '/task',
    name: 'task',
    component: Home,
    children: [ 
      {
        path: '/logicvuln', //逻辑漏洞
        name: 'logicvuln',
        component: () => import('@/pages/taskManagement/logicVuln.vue') ,
      },
      {
        path:'/createlogicvuln',
        name: 'createlogicvuln',
        component: () => import('@/pages/taskManagement/CreatelogicVulnTask.vue') ,
      },
      {
        path: '/logicvulnDetail',
        name: 'logicvulnDetail',
        component: () => import('@/pages/taskManagement/logicVulnDetail.vue') ,
      },
      {
        path: '/taskgroup',
        name: 'taskgroup',
        component: () => import('@/pages/taskManagement/TaskGroup.vue') ,
      }, 
      {
        path: '/taskGroupDetail',
        name: 'taskGroupDetail',
        component: () => import('@/pages/taskManagement/TaskGroupDetail.vue') ,
      },
      {
        path: '/taskscenario',
        name: 'taskscenario',
        component: () => import('@/pages/sceneManagement/Taskscenario.vue') ,
      }, 
     
      {
        path: '/createscene',
        name: 'createscene',
        component: () => import('@/pages/sceneManagement/Createscene.vue') ,
      },
      {
        path: '/task',
        name: 'task',
        component: () => import('@/pages/taskManagement/TaskList.vue') ,
      },
      {
        path:'/traffic',
        name: 'traffic',
        component: () => import('@/pages/taskManagement/TrafficList.vue') ,
      },
      {
        path:'/taskForce', //++++任务组
        name: 'taskForce',
        component: () => import('@/pages/taskManagement/taskForce.vue') ,
      },
      {
        path:'/assetManagement', //++++资产管理
        name: 'assetManagement',
        component: () => import('@/pages/taskManagement/taskForce2.vue') ,
      },
      {
        path:'/assetStatistics', //++++资产统计
        name: 'assetStatistics',
        component: () => import('@/pages/taskManagement/components/LeftZiChan.vue') ,
      },  
      {
        path:'/addtraffic',
        name: 'addtraffic',
        component: () => import('@/pages/taskManagement/Addtraffic.vue') ,
      },
      {
        path:'/trafficinfo',
        name: 'trafficinfo',
        component: () => import('@/pages/taskManagement/TrafficInfo.vue') ,
      },
      {
        path: '/knowledgegraphsm',
        name: 'knowledgegraph',
        component: () => import('@/pages/sceneManagement/Knowledgegraph.vue') ,
      },
      
      // {
      //   path: '/validationReport',
      //   name: 'validationReport',
      //   component: () => import('@/pages/taskManagement/validationReport.vue') ,
      // }, 
      {
        path: '/createtask',
        name: 'createtask',
        component: () => import('@/pages/taskManagement/CreateTask.vue') ,
      },
      {
        path: '/createXray',
        name: 'createXray',
        component: () => import('@/pages/othner/createXray.vue') ,
      },
      {
        path: '/detailXray',
        name: 'detailXray',
        component: () => import('@/pages/othner/detailXray.vue') ,
      },
      {
        path: '/createBurp',
        name: 'createBurp',
        component: () => import('@/pages/othner/createBurp.vue') ,
      },
      {
        path: '/detailBurp',
        name: 'detailBurp',
        component: () => import('@/pages/othner/detailBurp.vue') ,
      },
      {
        path: '/taskDetail',
        name: 'taskDetail',
        component: () => import('@/pages/taskManagement/TaskDetail.vue') ,
      },
      {
        path:'/vulScanTask',
        name: 'vulScanTask',
        component: () => import('@/pages/vulScanTask/VulScanTaskList.vue') ,
      },
      {
        path:'/vulScanTaskDetail',
        name: 'vulScanTaskDetail',
        component: () => import('@/pages/vulScanTask/vulScanTaskDetail.vue') ,
      },
      {
        path:'/vulScanDatabase',
        name: 'vulScanDatabase',
        component: () => import('@/pages/vulScanTask/VulScanDatabaseList.vue') ,
      },
      {
        path:'/vulScanTaskCreate',
        name: 'vulScanTaskCreate',
        component: () => import('@/pages/vulScanTask/VulScanTaskCreate.vue') ,
      },
      {
        path:'/reportlist',
        name: 'reportlist',
        component: () => import('@/pages/reportManagement/reportList.vue') ,
      },
      {
        path:'/createreport',
        name: 'createreport',
        component: () => import('@/pages/reportManagement/Createreport.vue') ,
      },
      {
        path:'/reporttemplate',
        name: 'reporttemplate',
        component: () => import('@/pages/reportManagement/Reporttemplate.vue') ,
      },
      {
        path: '/targetDetail',
        name: 'targetDetail',
        component: () => import('@/pages/taskManagement/TargetDetail.vue') ,
      },
      {
        path: '/verificationReport',
        name: 'verificationReport',
        component: () => import('@/pages/taskManagement/verificationReport/reportList.vue') ,
      },{
        path:'/uploadReport',
        name: 'uploadReport',
        component: () => import('@/pages/taskManagement/verificationReport/UploadReport.vue') ,
      },
      {
        path:'/reportTaskInfo',
        name: 'reportTaskInfo',
        component: () => import('@/pages/taskManagement/verificationReport/ReportTaskInfo.vue') ,
      }
      
    ],
  },
    {
        path:'/codereview',
        name: 'codereview',
        component: Home, 
        children: [ 
          {
            path: '/codereview',
            name: 'codereview',
            component: () => import('@/pages/codeRe/index.vue') ,
    
          }, 
        ]
      },
  // {
  //   path:'/taskreportview',
  //   name: 'taskreportview',
  //   component: () => import('@/pages/reportManagement/ExportFile.vue') ,
  // },
  {
    path: '/assetview', //资产中心
    component: Home,
    name: 'assetview',
    children: [
      {
        path: '/assetview',
        name: 'assetview',
        component: () => import('@/pages/assetCenter/assetview.vue') ,
      },
      {
        path: '/assettree',
        name: 'assettree',
        component: () => import('@/pages/assetCenter/assettree.vue') ,
      },
    ],
  },
  {
    path: '/riskvuln', //漏洞风险中心
    component: Home,
    name: 'riskvuln',
    children: [
      {
        path: '/riskvuln',
        name: 'riskvuln',
        component: () => import('@/pages/riskvuln/index.vue') ,
      },
      {
        path: '/riskvulnStatistics',
        name: 'riskvulnStatistics',
        component: () => import('@/pages/riskvuln/statistics.vue') ,
      },
    ],
  },
  {
    path: '/taskreportview',
    component: () => import('@/pages/reportManagement/TaskReportview.vue') ,
  },

  {
    path:'/targetreportview', 
    component: () => import('@/pages/reportManagement/TargetReportview.vue') ,
  },
  {
    path:'/bastargetlog', 
    component: () => import('@/pages/bas/log.vue'),
  },
  {
    path: '/x-ray',
    component: Home,
    name: 'x-ray',
    children: [
      {
        path: '/x-ray',
        name: 'x-ray',
        component: xray,
      },
      {
        path: '/burpsuite',
        name: 'burpsuite',
        component: burpsuite,
      },
      {
        path: '/detectionTask',
        name: 'detectionTask',
        component: detectionTask,
      },
    ],
  },
  {
    path: '/usermanagement',
    component: Home,
    name: 'usermanagement',
    children: [
      {
        path: '/usermanagement',
        name: 'usermanagement',
        component: usermanagement,
      },
      {
        path: '/usergroup',
        name: 'usergroup',
        component: usergroup,
      },
      {
        path: '/groupmanagement',
        name: 'groupmanagement',
        component: groupmanagement,
      },
    ],
  },
  {
    path: '/systemsetting',
    component: Home,
    name: 'systemsetting',
    children: [
      {
        path: '/systemsetting',
        name: 'systemsetting',
        component: systemsetting,
      }, 
      {
        path: '/node',
        component: node,
        name: 'node',
      },
      {
        path: '/routerSetting',
        name: 'routerSetting',
        component: () => import('@/pages/routerSet.vue'),
      }, 
    ],
  },
  {
    path: '/toolmanagement',
    component: Home,
    name: 'toolmanagement',
    children: [
      // {
      //   path: '/toolmanagement',
      //   name: 'toolmanagement',
      //   component: toolmanagement,
      // },
      // {
      //   path: '/toolinfo',
      //   name: 'toolinfo',
      //   component: toolinfo,
      // },
      // {
      //   path: '/toolinfo/:id',
      //   name: 'toolinfo',
      //   component: toolinfo,
      // },
      {
        path: '/auxiliarytool',
        name: 'auxiliarytool',
        component: auxiliarytool,
      },
      {
        path: "/dictionary",
        name: "Dictionary",
        component: () => import("../pages/dictionaryLibrary/dictionary.vue")
      },
        {
    path: '/editTable',
    component: () => import('@/pages/editTable/index.vue'),
    name: 'editTable',
 
  },
      {
        path:'/vulnerability',
        name:'VulnerabilityLibrary',
        component: vulnerability,
      },
      // {
      //   path: '/vulnerability',
      //   name: 'VulnerabilityLibrary',
      //   component: vulnerability,
      // },
      // {
      //   path: '/addnewtool',
      //   name: 'addnewtool',
      //   component: addnewtool,
      // },
      // {
      //   path: '/addnewloop',
      //   name: 'addnewloop',
      //   component: addnewloop,
      // },
      {
        path: '/fingerprint',
        name: 'FingerprintLib',
        component: fingerprint,
      },
      // {
      //   path: '/neo4j',
      //   name: 'neo4j',
      //   component: Neo4j,
      // },
      // {
      //   path: '/neo4jTest',
      //   name: 'neo4jTest',
      //   component: Neo4jTest,
      // },
    ],
  },
   {
     path: '/systemtip',
     component: systemtip,
     name: 'systemtip',
     meta: {
     title: 'systemtip',
   },
   },
    
  {
    path:'/log',
    component: Home,
    name: 'log',
    children: [
		{
			path: '/log',
			name: 'log',
			component: log,
		},
    {
			path: '/logconfig',
			name: 'logconfig',
			component: logconfig,
		},
    ],
  },
  {
    path: '/login',
    component: Login,
    name: 'Login', 
  }, 
 
  {
    path: '/bastask',
    redirect: '/hostsec/tasks',
  },
  {
    path: '/bastaskDetail',
    redirect: '/hostsec/tasks',
  },
  {
    path: '/createbas',
    redirect: '/hostsec/tasks',
  },
  {
    path: '/scriptLibrary',
    redirect: '/hostsec/tasks',
  },
  {
    path: '/evaluationScheme',
    redirect: '/hostsec/tasks',
  },
  {
    path: '/createscheme',
    redirect: '/hostsec/tasks',
  },
  {
    path: '/agent',
    redirect: '/hostsec/tasks',
  },
  {
    path: '/security',
    component: Home,
    name: 'security',
    children: [
      {
        path: '/hostsec/tasks',
        name: 'hostsecTasks',
        component: () => import('@/pages/security/HostSecurityHub.vue'),
      },
      {
        path: '/hostsec/task-detail',
        name: 'hostsecTaskDetail',
        component: () => import('@/pages/security/HostTaskDetail.vue'),
      },
      {
        path: '/hostsec/rules',
        name: 'hostsecRules',
        component: () => import('@/pages/security/HostDetectionRules.vue'),
      },
      {
        path: '/hostsec/vuln-rules',
        name: 'hostsecVulnRules',
        component: () => import('@/pages/security/VulnDetectionRules.vue'),
      },
      {
        path: '/hostsec/malware-rules',
        name: 'hostsecMalwareRules',
        component: () => import('@/pages/security/MalwareDetectionRules.vue'),
      },
      {
        path: '/appsec/tasks',
        name: 'appsecTasks',
        component: () => import('@/pages/security/AppSecTaskHub.vue'),
      },
      {
        path: '/appsec/task/new',
        name: 'appsecTaskNew',
        component: () => import('@/pages/security/AppScanTaskNew.vue'),
      },
      {
        path: '/appsec/rules',
        name: 'appsecRules',
        component: () => import('@/pages/security/AppDetectionRules.vue'),
      },
      {
        path: '/appsec/vul-db',
        name: 'appsecVulDB',
        component: () => import('@/pages/security/AppVulnerabilityDB.vue'),
      },
      {
        path: '/appsec/strategy',
        name: 'appsecStrategy',
        component: () => import('@/pages/security/AppScanStrategy.vue'),
      },
      {
        path: '/appsec/strategy/new',
        name: 'appsecStrategyNew',
        component: () => import('@/pages/security/AppScanStrategyEdit.vue'),
      },
      {
        path: '/appsec/strategy/edit',
        name: 'appsecStrategyEdit',
        component: () => import('@/pages/security/AppScanStrategyEdit.vue'),
      },
      {
        path: '/datasec/tasks',
        name: 'datasecTasks',
        component: () => import('@/pages/security/DataSecTaskHub.vue'),
      },
      {
        path: '/datasec/rules',
        name: 'datasecRules',
        component: () => import('@/pages/security/DataDetectionRules.vue'),
      },
      {
        path: '/hostsecurityhub',
        redirect: '/hostsec/tasks',
      },
      {
        path: '/hostbaseline',
        name: 'hostbaseline',
        component: () => import('@/pages/security/HostBaselineAudit.vue'),
      },
      {
        path: '/malware',
        name: 'malware',
        component: () => import('@/pages/security/MalwareScan.vue'),
      },
      {
        path: '/dbcheck',
        redirect: '/datasec/tasks?tab=db',
      },
      {
        path: '/sensitive',
        redirect: '/datasec/tasks?tab=sensitive',
      },
      {
        path: '/appspecific',
        redirect: '/appsec/tasks?tab=app',
      },
      {
        path: '/dynamicscan',
        redirect: '/appsec/tasks?tab=dyn',
      },
    ],
  },
  {
    path: '/campaign',
    component: Home,
    name: 'campaign',
    children: [
		{
			path: '/campaign',
			name: 'campaign',
			 component: () => import("../pages/gophish/campaign.vue")
		},
    {
			path: '/campaignDetail',
			name: 'campaignDetail',
			 component: () => import("../pages/gophish/campaignDetail.vue")
		},
    {
			path: '/group',
			name: 'group',
			 component: () => import("../pages/gophish/group.vue")
		},
    {
			path: '/template',
			name: 'template',
			 component: () => import("../pages/gophish/template.vue")
		},
    {
			path: '/gophishpage',
			name: 'gophishpage',
			 component: () => import("../pages/gophish/gophishpage.vue")
		},
    {
			path: '/sendprofile',
			name: 'sendprofile',
			 component: () => import("../pages/gophish/sendprofile.vue")
		},
    ],
  },
  {
    path: '*',
    component: () => import('@/views/404.vue'),
    meta: {
      title: '404未找到',
    },
  },
]
const router = new VueRouter({  
    routes  
}) 

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.islogin)) {
    if (localStorage.getItem('user')) {
      if (localStorage.getItem('role') == 'BcJQ9Ute7p0=') { //日志
        next({
          path: '/log'
        }) //表示已经登录
      } else {
        next({
          path: '/index'
        }) //表示已经登录
      }
    } else {
      next({
        path: '/login'
      })
    }
  } else {
    //表示不需要登录
    next() //继续往后走
  }
})

export default router