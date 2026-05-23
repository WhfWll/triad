package routers

import (
	"smart/api/rest"
	"smart/middleware"
	"smart/tools/enums"

	"github.com/gin-gonic/gin"
)

func RegisterRoute() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.LoginAuth, middleware.AccessWhite, middleware.RecordLogin, middleware.RecordUserLog)
	smartRouterGroup := router.Group("/smart")

	// 加载go-phish管理模块
	GoPhishRouter(smartRouterGroup)

	// 全局接口
	smartRouterGroup.GET("/global/options", rest.GlobalOptions) // 全局接口

	// 用户管理
	smartRouterGroup.GET("/user/logincaptcha", rest.UserLoginCaptcha)                                                                             // 用户 - 登陆验证码
	smartRouterGroup.POST("/user/login", rest.UserLogin)                                                                                          // 用户 - 登陆
	smartRouterGroup.GET("/user/enum", middleware.RoleAuth(enums.UserRoleAdministrator, enums.UserRoleSuperAdmin), rest.UserEnumList)             // 用户 - 枚举列表
	smartRouterGroup.GET("/user/list", middleware.RoleAuth(enums.UserRoleAdministrator, enums.UserRoleSuperAdmin), rest.UserList)                 // 用户 - 用户列表及筛选
	smartRouterGroup.GET("/user/info", rest.UserDetail)                                                                                           // 用户 - 用户详情
	smartRouterGroup.POST("/user/save", rest.UserManageOp)                                                                                        // 用户 - 新增/编辑用户
	smartRouterGroup.GET("/user/del", middleware.RoleAuth(enums.UserRoleAdministrator, enums.UserRoleSuperAdmin), rest.DelUser)                   // 用户 - 删除
	smartRouterGroup.POST("/user/updatepw", rest.UpdatePassWord)                                                                                  // 用户 - 修改密码
	smartRouterGroup.POST("/user/resetpw", middleware.RoleAuth(enums.UserRoleAdministrator, enums.UserRoleSuperAdmin), rest.ResetPassWord)        // 用户 - 重置密码
	smartRouterGroup.POST("/user/updateuserexp", middleware.RoleAuth(enums.UserRoleAdministrator, enums.UserRoleSuperAdmin), rest.UpdateUserExp)  // 用户 - 修改有效期
	smartRouterGroup.GET("/user/changestatus", middleware.RoleAuth(enums.UserRoleAdministrator, enums.UserRoleSuperAdmin), rest.ChangeUserStatus) // 用户 - 切换用户状态
	smartRouterGroup.GET("/user/logout", rest.LogOut)                                                                                             // 用户 - 退出登录
	smartRouterGroup.POST("/user/loginfreepassa", rest.LoginBy15Suo)                                                                              // 用户 - 15所项目免密登录
	smartRouterGroup.GET("/user/passwordcheck", rest.PasswordCheck)                                                                               // 用户 - 密码检查接口
	smartRouterGroup.POST("/user/loginfreepassb", rest.LoginByLianTong)                                                                           // 用户 - 靶场导调平台免密登录接口
	smartRouterGroup.POST("/user/loginfreepassc", rest.LoginBySiYuan)                                                                             // 用户 - 四院导调项目免密登录接口
	smartRouterGroup.POST("/user/loginfreepassd", rest.LoginByApiToken)                                                                           // 用户 - 中测yakit项目登录
	smartRouterGroup.POST("/user/loginfreepasse", rest.LoginByHTYZ)                                                                               // 用户 - 航天运载靶场免密登录
	smartRouterGroup.POST("/user/loginfreepassf", rest.LoginByGAYS)                                                                               // 用户 - 公安一所免密登录

	// 系统管理
	smartRouterGroup.GET("/system/authinfo", rest.AuthInfo)                                    // 产品授权 - 查询
	smartRouterGroup.POST("/system/authsave", rest.AuthSave)                                   // 产品授权 - 授权
	smartRouterGroup.GET("/system/generateProductID", rest.GenerateProductID)                  // 产品授权 - 生成产品id
	smartRouterGroup.POST("/system/securities", rest.CreateSystemSecurity)                     // 系统管理 - 创建安全配置
	smartRouterGroup.GET("/system/initialinfo", rest.SystemSecurityInfo)                       // 系统管理 - 安全配置初始化信息
	smartRouterGroup.POST("/system/configbackupconfigsave", rest.SystemConfigBackupConfigSave) // 配置备份 - 保存配置
	smartRouterGroup.GET("/system/configbackupconfiginfo", rest.SystemConfigBackupConfigInfo)  // 配置备份 - 配置信息
	smartRouterGroup.POST("/system/configbackupnow", rest.SystemConfigBackupNow)               // 配置备份 - 立即备份
	smartRouterGroup.GET("/system/configbackuplist", rest.SystemConfigBackupList)              // 配置备份 - 列表
	smartRouterGroup.GET("/system/configbackupdownload", rest.SystemConfigBackupDownload)      // 配置备份 - 下载
	smartRouterGroup.POST("/system/configbackupdelete", rest.SystemConfigBackupDelete)         // 配置备份 - 删除
	smartRouterGroup.POST("/system/configbackuprestore", rest.SystemConfigBackupRestore)       // 配置备份 - 恢复
	smartRouterGroup.POST("/system/ipwhitesave", rest.SystemSettingIpWhiteSave)                // 系统设置 - ip白名单 - 保存
	smartRouterGroup.GET("/system/ipwhiteinfo", rest.SystemSettingIpWhiteInfo)                 // 系统设置 - ip白名单 - 信息
	smartRouterGroup.POST("/system/syslogsave", rest.SystemSettingSyslogSave)                  // 系统设置 - syslog服务 - 保存
	smartRouterGroup.GET("/system/sysloginfo", rest.SystemSettingSyslogInfo)                   // 系统设置 - syslog服务 - 信息
	smartRouterGroup.POST("/system/mailsave", rest.SystemSettingMailSave)                      // 系统设置 - 邮箱配置 - 保存
	smartRouterGroup.GET("/system/mailinfo", rest.SystemSettingMailInfo)                       // 系统设置 - 邮箱配置 - 信息
	smartRouterGroup.POST("/system/mailverify", rest.SystemSettingMailVerify)                  // 系统设置 - 邮箱配置 - 验证
	smartRouterGroup.POST("/system/networkconfigsave", rest.SystemSettingNetworkConfigSave)    // 系统设置 - 网络配置 - 保存
	smartRouterGroup.GET("/system/networkconfiginfo", rest.SystemSettingNetworkConfigInfo)     // 系统设置 - 网络配置 - 信息
	smartRouterGroup.GET("/system/routelist", rest.RouteList)                                  // 系统设置 - 配置路由 - 列表
	smartRouterGroup.POST("/system/routeadd", rest.RouteAdd)                                   // 系统设置 - 配置路由 - 增加
	smartRouterGroup.POST("/system/routedelete", rest.RouteDelete)                             // 系统设置 - 配置路由 - 删除
	smartRouterGroup.GET("/system/monitorwarninfo", rest.SystemSettingMonitorWarnInfo)         // 系统设置 - 系统监控告警 - 信息
	smartRouterGroup.POST("/system/monitorwarnsave", rest.SystemSettingMonitorWarnSave)        // 系统设置 - 系统监控告警 - 保存
	smartRouterGroup.POST("/system/targetipsave", rest.TargetIpSave)                           // 业务设置 - 测试目标黑白名单更新
	smartRouterGroup.GET("/system/targetiplist", rest.TargetIpList)                            // 业务设置 - 测试目标黑白名单查询
	smartRouterGroup.GET("/system/getreverseiphost", rest.GetReverseIpHost)                    // 业务设置 - 远程监听查询
	smartRouterGroup.POST("/system/reverseiphostsave", rest.ReverseIpHostSave)                 // 业务设置 - 远程监听修改
	smartRouterGroup.POST("/system/tcpblindtestsave", rest.BusinessSettingTcpBlindTestSave)    // 业务设置 - tcp盲测平台 - 保存
	smartRouterGroup.GET("/system/tcpblindtestinfo", rest.BusinessSettingTcpBlindTestInfo)     // 业务设置 - tcp盲测平台 - 信息
	smartRouterGroup.POST("/system/httpblindtestsave", rest.BusinessSettingHttpBlindTestSave)  // 业务设置 - http盲测平台 - 保存
	smartRouterGroup.GET("/system/httpblindtestinfo", rest.BusinessSettingHttpBlindTestInfo)   // 业务设置 - http盲测平台 - 信息
	smartRouterGroup.POST("/system/dnsblindtestsave", rest.BusinessSettingDnsBlindTestSave)    // 业务设置 - dns盲测平台 - 保存
	smartRouterGroup.GET("/system/dnsblindtestinfo", rest.BusinessSettingDnsBlindTestInfo)     // 业务设置 - dns盲测平台 - 信息
	smartRouterGroup.POST("/system/icmpblindtestsave", rest.BusinessSettingIcmpBlindTestSave)  // 业务设置 - icmp盲测平台 - 保存
	smartRouterGroup.GET("/system/icmpblindtestinfo", rest.BusinessSettingIcmpBlindTestInfo)   // 业务设置 - icmp盲测平台 - 信息
	smartRouterGroup.GET("/system/curtasksinfo", rest.CurTasksInfo)                            // 业务设置 - 任务并发配置 - 信息
	smartRouterGroup.GET("/system/curtaskssave", rest.CurTasksSave)                            // 业务设置 - 任务并发配置 - 保存
	smartRouterGroup.GET("/system/securityscanconcurrencyinfo", rest.SecurityScanConcurrencyInfo)   // 安全检查并发配置 - 信息
	smartRouterGroup.POST("/system/securityscanconcurrencysave", rest.SecurityScanConcurrencySave) // 安全检查并发配置 - 保存
	smartRouterGroup.GET("/system/usescoreinfo", rest.UseScoreInfo)                            // 业务设置 - 可以利用评分 - 信息
	smartRouterGroup.POST("/system/usescoresave", rest.UseScoreSave)                           // 业务设置 - 可以利用评分 - 保存
	smartRouterGroup.GET("/system/testscopeinfo", rest.TestScopeInfo)                          // 业务设置 - 测试范围校验开关 - 信息
	smartRouterGroup.POST("/system/testscopesave", rest.TestScopeSave)                         // 业务设置 - 测试范围校验开关 - 保存

	smartRouterGroup.GET("/system/cpuinfo", rest.CpuInfo)                            // 系统监控 - cpu
	smartRouterGroup.GET("/system/memoryinfo", rest.MemoryInfo)                      // 系统监控 - 内存
	smartRouterGroup.GET("/system/diskinfo", rest.DiskInfo)                          // 系统监控 - 磁盘
	smartRouterGroup.POST("/system/uploadupgradefile", rest.UploadUpgradeFile)       // 升级还原 - 上传文件
	smartRouterGroup.POST("/system/confirmupgrade", rest.ConfirmUpgrade)             // 升级还原 - 确认升级
	smartRouterGroup.GET("/system/upgradestatus", rest.GetUpgradeStatus)             // 升级还原 - 获取升级进度
	smartRouterGroup.POST("/system/manualrollback", rest.SystemManualRollback)       // 升级还原 - 手动回滚
	smartRouterGroup.GET("/system/systemversion", rest.SystemVersion)                // 升级还原 - 系统版本信息
	smartRouterGroup.POST("/system/generatetoken", rest.GenerateToken)               // 系统管理 - 秘钥生成
	smartRouterGroup.GET("/system/tokenlist", rest.TokenList)                        // 系统管理 - 秘钥列表
	smartRouterGroup.POST("/system/tokendel", rest.TokenDelete)                      // 系统管理 - 秘钥删除
	smartRouterGroup.GET("/system/nodedownload", rest.SystemNodeDownload)            // 节点管理 - 下载节点安装包
	smartRouterGroup.POST("/system/nodesetdistribute", rest.SystemNodeSetDistribute) // 节点管理 - 设置是否开启分布式
	smartRouterGroup.GET("/system/nodegetdistribute", rest.SystemNodeGetDistribute)  // 节点管理 - 获取是否开启分布式
	smartRouterGroup.POST("/system/nodeadd", rest.SystemNodeAdd)                     // 节点管理 - 节点添加
	smartRouterGroup.POST("/system/nodeedit", rest.SystemNodeEdit)                   // 节点管理 - 节点编辑
	smartRouterGroup.GET("/system/nodeinfo", rest.SystemNodeInfo)                    // 节点管理 - 节点详情
	smartRouterGroup.GET("/system/nodelist", rest.SystemNodeList)                    // 节点管理 - 节点列表
	smartRouterGroup.POST("/system/nodedel", rest.SystemNodeDel)                     // 节点管理 - 节点删除
	smartRouterGroup.POST("/system/nodedisorenable", rest.SystemNodeDisOrEnable)     // 节点管理 - 节点是否启用
	smartRouterGroup.GET("/system/nodeallenable", rest.SystemNodeAllEnable)          // 节点管理 - 所有可用节点

	// 大模型管理
	smartRouterGroup.GET("/llmmodel/list", rest.LlmModelList)              // 大模型 - 列表
	smartRouterGroup.GET("/llmmodel/detail", rest.LlmModelDetail)          // 大模型 - 详情
	smartRouterGroup.POST("/llmmodel/save", rest.LlmModelSave)             // 大模型 - 保存（添加/编辑）
	smartRouterGroup.POST("/llmmodel/delete", rest.LlmModelDelete)         // 大模型 - 删除
	smartRouterGroup.POST("/llmmodel/setdefault", rest.LlmModelSetDefault) // 大模型 - 设置默认
	smartRouterGroup.POST("/llmmodel/check", rest.LlmModelEnabledTest)     // 大模型 - 测试是否可用
	smartRouterGroup.GET("/llmmodel/enums", rest.LlmModelEnums)            // 大模型 - 枚举

	smartRouterGroup.GET("/llmmodel/enhance/detail", rest.LlmModelEnhancementDetail) // 大模型 - 模型增强详情
	smartRouterGroup.POST("/llmmodel/enhance/edit", rest.LlmModelEnhancementEdit)    // 大模型 - 模型增强切换

	// AI应用场景管理
	smartRouterGroup.GET("/aiscenario/list", rest.AiScenarioList)      // AI应用场景 - 列表
	smartRouterGroup.POST("/aiscenario/config", rest.AiScenarioConfig) // AI应用场景 - 配置

	// 反连服务器
	smartRouterGroup.GET("/reverse/status", rest.ReverseServerStatus)       // 反连服务器 - 当前状态
	smartRouterGroup.POST("/reverse/start", rest.ReverseServerStart)        // 反连服务器 - 开启反连服务器
	smartRouterGroup.POST("/reverse/stop", rest.ReverseServerStop)          // 反连服务器 - 关闭反连服务器
	smartRouterGroup.GET("/reverse/message", rest.ReverseServerMessages)    // 反连服务器 - 当前信息
	smartRouterGroup.POST("/reverse/clear", rest.ReverseServerClearMessage) // 反连服务器 - 清空消息

	// 工具管理
	smartRouterGroup.GET("/tools/vulenum", rest.ToolsVulEnum)                                // 漏洞库 - 枚举
	smartRouterGroup.GET("/tools/vullist", rest.ToolsVulList)                                // 漏洞库 - 列表
	smartRouterGroup.GET("/tools/vuldetail", rest.ToolsVulDetail)                            // 漏洞库 - 详情
	smartRouterGroup.POST("/tools/vuledit", rest.ToolsVulEdit)                               // 漏洞库 - 编辑
	smartRouterGroup.POST("/tools/vuleditstatus", rest.ToolsVulEditStatus)                   // 漏洞库 - 修改状态
	smartRouterGroup.POST("/tools/importvulnvulkit", rest.ImportVulnFromVulKit)              // 漏洞库 - 导入VulKit漏洞脚本
	smartRouterGroup.GET("/tools/enum", rest.ToolsDictionaryEnum)                            // 字典库 - 枚举
	smartRouterGroup.GET("/tools/dictlist", rest.ToolsDictionaryList)                        // 字典库 - 列表
	smartRouterGroup.GET("/tools/dictinfo", rest.ToolsDictionaryDetail)                      // 字典库 - 详情
	smartRouterGroup.GET("/tools/dictdel", rest.ToolsDictionaryDelete)                       // 字典库 - 删除
	smartRouterGroup.GET("/tools/dictsetdefault", rest.ToolsDictionarySetDefault)            // 字典库 - 设置默认字典
	smartRouterGroup.POST("/tools/dictaddoredit", rest.ToolsDictionaryAddOrEdit)             // 字典库 - 新增或编辑
	smartRouterGroup.GET("/tools/fingerenum", rest.ToolsFingerEnum)                          // 指纹库 - 枚举
	smartRouterGroup.GET("/tools/fingerlist", rest.ToolsFingerList)                          // 指纹库 - 列表
	smartRouterGroup.POST("/tools/addfinger", rest.AddToolsFinger)                           // 指纹库 - 添加
	smartRouterGroup.POST("/tools/editfinger", rest.EditToolsFinger)                         // 指纹库 - 修改
	smartRouterGroup.GET("/tools/fingerinfo", rest.ToolsFingerDetail)                        // 指纹库 - 详情
	smartRouterGroup.GET("/tools/delfinger", rest.DelToolsFinger)                            // 指纹库 - 删除
	smartRouterGroup.POST("/tools/testfinger", rest.TestToolsFinger)                         // 指纹库 - 测试
	smartRouterGroup.GET("/tools/testfingerresult", rest.ToolsFingerTestResult)              // 指纹库 - 测试结果
	smartRouterGroup.POST("/tools/supportiphostcreate", rest.ToolsIpHostBindCreate)          // 辅助工具 - IP域名绑定 - 创建/编辑
	smartRouterGroup.GET("/tools/supportiphostlist", rest.ToolsIpHostBindList)               // 辅助工具 - IP域名绑定 - 列表
	smartRouterGroup.POST("/tools/supportiphostdel", rest.ToolsIpHostBindDel)                // 辅助工具 - IP域名绑定 - 删除
	smartRouterGroup.POST("/tools/supportpingcreate", rest.ToolsPingCreate)                  // 辅助工具 - Ping - 创建
	smartRouterGroup.POST("/tools/supportpingstop", rest.ToolsPingStop)                      // 辅助工具 - Ping - 停止
	smartRouterGroup.GET("/tools/supportpingresult", rest.ToolsPingResultUnread)             // 辅助工具 - Ping - 结果
	smartRouterGroup.POST("/tools/supporttraceroutecreate", rest.ToolsTracerouteCreate)      // 辅助工具 - traceroute - 创建
	smartRouterGroup.POST("/tools/supporttraceroutestop", rest.ToolsTracerouteStop)          // 辅助工具 - traceroute - 停止
	smartRouterGroup.GET("/tools/supporttracerouteresult", rest.ToolsTracerouteResultUnread) // 辅助工具 - traceroute - 结果
	smartRouterGroup.GET("/tools/toolfilelist", rest.ToolFileList)                           // 辅助工具 - 工具库 - 列表
	smartRouterGroup.GET("/tools/toolfiledownload", rest.ToolFileDownload)                   // 辅助工具 - 工具库 - 下载

	// 用户组管理
	smartRouterGroup.GET("/usergroup/page", rest.UserGroupList)                  // 用户组列表
	smartRouterGroup.POST("/usergroup/create", rest.UserGroupCreate)             // 新建用户组
	smartRouterGroup.POST("/usergroup/update", rest.UserGroupUpdate)             // 编辑用户组
	smartRouterGroup.GET("/usergroup/select", rest.GroupSelect)                  // 用户组选择
	smartRouterGroup.POST("/usergroup/updatestatus", rest.UserGroupUpdateStatus) // 更新用户组状态-删除等
	// 用户组组员管理 用户组 与 用户 关联
	smartRouterGroup.POST("/usergroup/preselectionpage", rest.GroupUserPreselectionList) // 组员列表
	smartRouterGroup.POST("/usergroup/alreadypage", rest.GroupUserAlreadyList)           // 组内成员
	smartRouterGroup.POST("/usergroup/relation", rest.GroupUserRelation)                 // 保存至用户组

	// 场景管理
	smartRouterGroup.GET("/scene/enum", rest.SceneEnums)                                     // 任务场景 - 枚举值
	smartRouterGroup.POST("/scene/save", rest.SceneTaskTemplateSave)                         // 任务场景 - 创建/更新
	smartRouterGroup.GET("/scene/list", rest.SceneTaskTemplateList)                          // 任务场景 - 列表
	smartRouterGroup.POST("/scene/vullist", rest.ToolsVulListToTaskTemplate)                 // 任务场景 - 漏洞列表
	smartRouterGroup.GET("/scene/info", rest.SceneTaskTemplateDetail)                        // 任务场景 - 详情
	smartRouterGroup.POST("/scene/copy", rest.SceneTaskTemplateCopy)                         // 任务场景 - 拷贝
	smartRouterGroup.POST("/scene/setdefault", rest.SceneTaskTemplateSetDefault)             // 任务场景 - 设置默认
	smartRouterGroup.POST("/scene/del", rest.SceneTaskTemplateDel)                           // 任务场景 - 删除
	smartRouterGroup.GET("/scene/sceneoptions", rest.SceneTaskTemplateToTaskTemplateOptions) // 渗透任务 - 场景选项
	smartRouterGroup.GET("/scene/graph", rest.Graph)                                         // 渗透任务 - 场景选项
	//smartRouterGroup.POST("/scene/tasktemplateorphan", rest.ToolsVulnerabilityToTaskTemplateCheckVulIds) // 任务场景查询哪些漏洞为孤独节点??

	// 任务中心
	smartRouterGroup.GET("/task/enum", rest.TaskTaskEnum)                          // 渗透任务 - 枚举
	smartRouterGroup.POST("/task/websitelogincheck", rest.TaskTaskWebLoginCheck)   // 渗透任务 - 验证站点是否能登陆
	smartRouterGroup.POST("/task/save", rest.TaskSave)                             // 渗透任务 - 创建
	smartRouterGroup.GET("/task/list", rest.TaskTaskList)                          // 渗透任务 - 列表
	smartRouterGroup.POST("/task/del", rest.TaskTaskDel)                           // 渗透任务 - 删除
	smartRouterGroup.GET("/task/copy", rest.TaskTaskCopy)                          // 渗透任务 - 拷贝
	smartRouterGroup.GET("/task/changestate", rest.TaskTaskChangeState)            // 渗透任务 - 修改任务运行状态
	smartRouterGroup.GET("/task/getstate", rest.GetState)                          // 渗透任务 - 获取任务状态
	smartRouterGroup.GET("/task/overview", rest.OverView)                          // 渗透任务 - 任务概览
	smartRouterGroup.GET("/task/configinfo", rest.TaskConfigInfo)                  // 渗透任务 - 任务配置信息
	smartRouterGroup.GET("/task/targetchangestate", rest.TargetChangeState)        // 渗透任务 - 修改测试目标状态
	smartRouterGroup.GET("/task/targetlist", rest.TargetList)                      // 渗透任务 - 测试目标列表及筛选
	smartRouterGroup.POST("/task/updatetargetusescore", rest.UpdateTargetUseScore) // 渗透任务 - 批量修改目标的利用评分和状态
	smartRouterGroup.GET("/task/targetdel", rest.TargetDel)                        // 渗透任务 - 测试目标删除
	smartRouterGroup.GET("/task/taskresultlist", rest.TaskResultList)              // 渗透任务 - 信息收集列表及筛选
	smartRouterGroup.GET("/task/taskresulturltree", rest.TaskResultUrlTree)        // 渗透任务 - 信息收集列表及筛选
	smartRouterGroup.GET("/task/taskresultdel", rest.TaskResultDel)                // 渗透任务 - 信息收集删除
	smartRouterGroup.GET("/task/taskresultdetail", rest.TaskResultDetail)          // 渗透任务 - 信息收集详情
	smartRouterGroup.GET("/task/vullist", rest.VulList)                            // 渗透任务 - 漏洞测试列表及筛选
	smartRouterGroup.GET("/task/vulinfo", rest.VulInfo)                            // 渗透任务 - 漏洞测试详情
	smartRouterGroup.GET("/task/checkvul", rest.CheckVul)                          // 渗透任务 - 漏洞测试qianji
	smartRouterGroup.GET("/task/getvulsnapshot", rest.GetVulSnapshot)              // 渗透任务 - 查看漏洞截图
	smartRouterGroup.GET("/task/vuldel", rest.VulDel)                              // 渗透任务 - 漏洞测试删除
	smartRouterGroup.POST("/task/vultest", rest.VulTest)                           // 渗透任务 - 漏洞测试测试
	smartRouterGroup.POST("/task/vulverify", rest.VulVerify)                       // 渗透任务 - 漏洞测试验证
	smartRouterGroup.POST("/task/asyncvulverify", rest.AsyncVulVerify)             // 渗透任务 - 漏洞测试验证
	smartRouterGroup.GET("/task/testvultest", rest.TestVulTest)                    // 渗透任务 - 待测漏洞测试	smartRouterGroup.GET("/task/attacklink", rest.AttackLink)                               // 攻击链路图
	smartRouterGroup.GET("/task/attacklink", rest.AttackLink)                      // 渗透任务 - 攻击链路图

	smartRouterGroup.GET("/task/loglist", rest.LogList)                                     // 渗透任务 - 测试日志列表及筛选
	smartRouterGroup.GET("/task/loginfo", rest.LogInfo)                                     // 渗透任务 - 测试日志详情
	smartRouterGroup.POST("/task/apisave", rest.TaskTaskApiSave)                            // 渗透任务 - 第三方创建接口
	smartRouterGroup.GET("/task/progress", rest.TaskProgress)                               // 渗透任务 - 获取进度
	smartRouterGroup.GET("/task/apivullist", rest.ApiVulList)                               // 渗透任务 - 漏洞测试列表及筛选
	smartRouterGroup.GET("/task/addtarget", rest.AddTarget)                                 // 渗透任务 - 动态添加目标
	smartRouterGroup.POST("/task/addattackface", rest.AddAttackFace)                        // 渗透任务 - 动态添加攻击面
	smartRouterGroup.POST("/task/addvul", rest.AddVul)                                      // 渗透任务 - 动态添加漏洞
	smartRouterGroup.GET("/task/tasktargetmap", rest.TaskTargetMap)                         // 渗透任务 - 路径图
	smartRouterGroup.GET("/task/tasktargetmapnodedetail", rest.TaskTargetMapNodeDetail)     // 渗透任务 - 路径图节点详情
	smartRouterGroup.GET("/task/taskthreeexport", rest.TaskThreeExport)                     // 渗透任务 - 三方数据导出
	smartRouterGroup.GET("/task/tasktopologymap", rest.TaskTopologyMap)                     // 渗透任务 - 攻击拓扑图
	smartRouterGroup.GET("/task/tasktopologymapnodedetail", rest.TaskTopologyMapNodeDetail) // 渗透任务 - 攻击拓扑图 - 节点详情
	smartRouterGroup.GET("/task/taskalltaskvulbypage", rest.TaskAllTaskVulByPage)           // 渗透任务 - 发现的所有漏洞信息
	smartRouterGroup.GET("/task/vulevidencelist", rest.VulEvidenceList)                     // 漏洞取证 - 列表及筛选
	smartRouterGroup.GET("/task/risktypeenum", rest.RiskTypeInfoEnum)                       // 漏洞取证 - 风险类型枚举
	smartRouterGroup.GET("/task/delvulevidence", rest.VulEvidenceDel)                       // 漏洞取证 - 删除漏洞取证列表信息
	smartRouterGroup.GET("/task/vulevidenceinfo", rest.VulEvidenceInfo)                     // 漏洞取证 - 漏洞取证详情信息
	smartRouterGroup.GET("/task/remotesessionlist", rest.RemoteSessionList)                 // 远程会话 - 列表及筛选
	smartRouterGroup.GET("/task/delremotesession", rest.RemoteSessionDel)                   // 远程会话 - 删除列表
	smartRouterGroup.GET("/task/remotesessioninfo", rest.RemoteSessionInfo)                 // 远程会话 - 详情
	smartRouterGroup.GET("/task/remotesessiondir", rest.RemoteSessionDir)                   // 远程会话 - 列出目录
	smartRouterGroup.GET("/task/downloadfile", rest.FileDownload)                           // 远程会话 - 下载已下载文件
	smartRouterGroup.GET("/task/delfile", rest.DelFile)                                     // 远程会话 - 删除已下载文件
	smartRouterGroup.GET("/task/captureinfoenum", rest.CaptureInfoEnum)                     // 远程会话 - 抓取信息枚举
	smartRouterGroup.GET("/task/break", rest.BreakShell)                                    // 远程会话 - 断开
	smartRouterGroup.POST("/task/exceshell", rest.VulEvidenceUse)                           // 远程会话 - 利用shell
	smartRouterGroup.GET("/task/captureinfo", rest.ToCaptureInfo)                           // 远程会话 - 抓取信息
	smartRouterGroup.POST("/task/exceshellmany", rest.ExceShellMany)                        // 远程会话 - 批量收集信息
	smartRouterGroup.POST("/task/filemanagement", rest.FileManagement)                      // 远程会话 -文件管理
	smartRouterGroup.GET("/task/shellfiledownload", rest.ShellFileDownload)                 // 远程会话 -文件管理
	smartRouterGroup.GET("/task/flowtaskenum", rest.FlowTaskEnum)                           // 流量分析 - 枚举
	smartRouterGroup.GET("/task/flowtasklist", rest.FlowTaskList)                           // 流量分析 - 任务列表
	smartRouterGroup.GET("/task/flowtaskdel", rest.FlowTaskDel)                             // 流量分析 - 任务删除
	smartRouterGroup.POST("/task/flowtaskadd", rest.FlowTaskAdd)                            // 流量分析 - 创建任务
	smartRouterGroup.GET("/task/changeflowtaskstatus", rest.ChangeFlowTaskStatus)           // 流量分析 - 任务操作
	smartRouterGroup.GET("/task/flowtaskstatus", rest.FlowTaskStatus)                       // 流量分析 - 任务状态
	smartRouterGroup.GET("/task/flowtaskinfo", rest.FlowTaskInfo)                           // 流量分析 - 任务详情
	smartRouterGroup.GET("/task/httpscert", rest.HttpsCert)                                 // 流量分析 - https证书下载
	smartRouterGroup.GET("/task/flowrisklist", rest.FlowRiskList)                           // 流量分析 - 漏洞列表
	smartRouterGroup.GET("/task/flowriskinfo", rest.FlowRiskInfo)                           // 流量分析 - 漏洞详情
	smartRouterGroup.GET("/task/flowriskdel", rest.FlowRiskDel)                             // 流量分析 - 漏洞删除
	smartRouterGroup.GET("/task/flowbaselist", rest.FlowBaseList)                           // 流量分析 - 被动流量列表
	smartRouterGroup.GET("/task/flowbaseinfo", rest.FlowBaseInfo)                           // 流量分析 - 被动流量详情
	smartRouterGroup.GET("/task/flowbasedel", rest.FlowBaseDel)                             // 流量分析 - 被动流量删除
	smartRouterGroup.GET("/task/flowloginfo", rest.FlowLogInfo)                             // 流量分析 - 被动流量日志查询
	smartRouterGroup.GET("/task/flowlogdel", rest.FlowLogDel)                               // 流量分析 - 被动流量日志清除
	smartRouterGroup.POST("/task/flowtaskedit", rest.FlowTaskEdit)                          // 流量分析 - 被动流量日志清除
	smartRouterGroup.GET("/task/flowtaskexport", rest.FlowTaskExport)                       // 流量分析 - 流量导出

	//报告中心
	smartRouterGroup.GET("/report/enum", rest.ReportEnum)         // 报告枚举
	smartRouterGroup.GET("/report/list", rest.ReportList)         // 报告清单-列表及筛选
	smartRouterGroup.GET("/report/download", rest.ReportDownload) // 报告清单-下载
	smartRouterGroup.GET("/report/del", rest.ReportDel)           // 报告清单-删除
	smartRouterGroup.POST("/report/save", rest.ReportSave)        // 生成报告-保存

	//日志管理
	smartRouterGroup.GET("/logs/enum", rest.LogsEnum)                           //日志管理-日志枚举
	smartRouterGroup.GET("/logs/logauditlist", rest.LogAuditList)               //日志管理-审计日志-列表
	smartRouterGroup.POST("/logs/logauditempty", rest.LogAuditEmpty)            //日志管理-审计日志-清空
	smartRouterGroup.POST("/logs/logbackupconfig", rest.LogBackupConfig)        //日志管理-日志备份-配置
	smartRouterGroup.GET("/logs/logbackupconfiginfo", rest.LogBackupConfigInfo) //日志管理-日志备份-配置信息
	smartRouterGroup.POST("/logs/setexptime", rest.SetLogExpirationTime)        //日志管理-日志备份-设置日志过期时间
	smartRouterGroup.GET("/logs/getexptime", rest.GetLogExpirationTime)         //日志管理-日志备份-查询日志过期时间
	smartRouterGroup.GET("/logs/logbackuplist", rest.LogBackupList)             //日志管理-日志备份-列表
	smartRouterGroup.GET("/logs/logbackupdownload", rest.LogBackupDownload)     //日志管理-日志备份-下载
	smartRouterGroup.POST("/logs/logbackupdelete", rest.LogBackupDelete)        //日志管理-日志备份-删除
	smartRouterGroup.POST("/logs/logbackupnow", rest.LogBackupNow)              //日志管理-日志备份-立即备份

	// 三方工具
	smartRouterGroup.POST("/tripartite/xraysave", rest.TripartiteToolsXrayCreate)                    // xray - 任务创建
	smartRouterGroup.POST("/tripartite/xrayupload", rest.TripartiteToolsXrayUpload)                  // xray - 任务导入
	smartRouterGroup.POST("/tripartite/xraydel", rest.TripartiteToolsXrayDel)                        // xray - 任务删除
	smartRouterGroup.GET("/tripartite/xraylist", rest.TripartiteToolsXRayPage)                       // xray - 任务列表
	smartRouterGroup.GET("/tripartite/xraydetaillist", rest.TripartiteToolsXRayDetailPage)           // xray - 任务详情详情
	smartRouterGroup.POST("/tripartite/burpsuitesave", rest.TripartiteToolsBurpsuiteCreate)          // burpsuite - 任务创建
	smartRouterGroup.POST("/tripartite/burpsuiteupload", rest.TripartiteToolsBurpsuiteUpload)        // burpsuite - 任务导入
	smartRouterGroup.POST("/tripartite/burpsuitedel", rest.TripartiteToolsBurpsuiteDel)              // burpsuite - 任务删除
	smartRouterGroup.GET("/tripartite/burpsuitelist", rest.TripartiteToolsBurpsuitePage)             // burpsuite - 任务列表
	smartRouterGroup.GET("/tripartite/burpsuitedetaillist", rest.TripartiteToolsBurpsuiteDetailPage) // burpsuite - 任务详情列表
	smartRouterGroup.GET("/tripartite/wifiaplist", rest.TripartiteToolsWifiApList)                   // wifi - 可用wifi列表
	smartRouterGroup.POST("/tripartite/wificreate", rest.TripartiteToolsWifiCreate)                  // wifi - 任务创建
	smartRouterGroup.GET("/tripartite/wifilist", rest.TripartiteToolsWifiPage)                       // wifi - 任务列表
	smartRouterGroup.POST("/tripartite/wifidel", rest.TripartiteToolsWifiDel)                        // wifi - 任务删除

	//bas
	smartRouterGroup.GET("/bas/enum", rest.BasEnum)                                 // bas - 枚举接口
	smartRouterGroup.POST("/bas/taskcreate", rest.BasCreateTask)                    // bas - 任务创建
	smartRouterGroup.POST("/bas/receivresult", rest.BasReceivResult)                // bas - bas心跳及检测结果接收
	smartRouterGroup.GET("/bas/basvulstat", rest.BasVulStat)                        // bas - 漏洞测试统计
	smartRouterGroup.GET("/bas/basvullist", rest.BasVulList)                        // bas - 漏洞测试列表
	smartRouterGroup.GET("/bas/basvuldel", rest.BasVulDel)                          // bas - 漏洞测试删除
	smartRouterGroup.POST("/bas/bastruleimpor", rest.BasRuleImport)                 // bas - 规则导入
	smartRouterGroup.GET("/bas/basruleenum", rest.BasRuleEnum)                      // bas - 规则枚举
	smartRouterGroup.GET("/bas/basruleget", rest.BasRuleGet)                        // bas - 规则列表
	smartRouterGroup.GET("/bas/basruleinfo", rest.BasRuleInfo)                      // bas - 规则详情
	smartRouterGroup.POST("/bas/basruleedit", rest.BasRuleEdit)                     // bas - 规则编辑
	smartRouterGroup.POST("/bas/bastemplatecreate", rest.BasTemplateCreate)         // bas - 剧本集创建
	smartRouterGroup.GET("/bas/bastemplateget", rest.BasTemplateGet)                // bas - 剧本集列表
	smartRouterGroup.GET("/bas/bastemplatebyid", rest.BasGetTemplateById)           // bas - 剧本集详情
	smartRouterGroup.POST("/bas/bastemplatedel", rest.BasDelTemplateById)           // bas - 剧本集删除
	smartRouterGroup.POST("/bas/bastemplatesetdefault", rest.BasTemplateSetDefault) // bas - 剧本集设置默认
	smartRouterGroup.POST("/bas/basagentisonline", rest.BasAgentIsOnline)           // bas - agent是否在线
	smartRouterGroup.GET("/bas/bastaskget", rest.BasGetTask)                        // bas - 任务列表
	smartRouterGroup.POST("/bas/bastaskend", rest.BasEndTaskById)                   // bas - 任务结束
	smartRouterGroup.POST("/bas/bastaskdel", rest.BasDelTask)                       // bas - 任务删除
	smartRouterGroup.GET("/bas/bastasktargetpage", rest.BasGetTaskTargetPage)       // bas - 任务详情
	smartRouterGroup.GET("/bas/bastasktargetlog", rest.BasGetTargetLogs)            // bas - 任务目标日志
	smartRouterGroup.POST("/bas/bastasktargetdel", rest.BasGetTargetDel)            // bas - 任务目标删除
	smartRouterGroup.GET("/bas/basagentdownload", rest.BasAgentDownload)            // bas - agent- 下载
	smartRouterGroup.GET("/bas/basagentlist", rest.BasAgentList)                    // bas - agent- 列表
	smartRouterGroup.GET("/bas/basagentlive", rest.BasAgentLive)                    // bas - agent- 可用节点列表
	smartRouterGroup.POST("/bas/basagentstatusedit", rest.BasAgentStatusEdit)       // bas - agent- 状态修改

	//首页
	smartRouterGroup.GET("/homepage/taskinfostat", rest.TaskInfoStat)         //首页 - 任务统计
	smartRouterGroup.GET("/homepage/vulevidencestat", rest.VulEvidenceStat)   //首页 - 漏洞取证
	smartRouterGroup.GET("/homepage/targetriskstat", rest.TargetRiskStat)     //首页 - 目标统计
	smartRouterGroup.GET("/homepage/taskvulriskstat", rest.TaskVulRiskStat)   //首页 - 任务漏洞统计
	smartRouterGroup.GET("/homepage/toolinfostat", rest.ToolInfoStat)         //首页 - 渗透支撑
	smartRouterGroup.GET("/homepage/vultypestat", rest.VulTypeStat)           //首页 - 漏洞类型统计
	smartRouterGroup.GET("/homepage/vulfindtrendstat", rest.VulFindTrendStat) //首页 - 漏洞发现趋势
	smartRouterGroup.GET("/homepage/messagestat", rest.MessageStat)           //首页 - 最新消息统计模块

	//报告验证接口
	smartRouterGroup.POST("/reportverify/upload", rest.ReportVerifyUpload)            // 报告验证 - 报告上传
	smartRouterGroup.GET("/reportverify/tasklist", rest.ReportVerifyTaskList)         // 报告验证 - 任务列表
	smartRouterGroup.GET("/reportverify/taskdetail", rest.ReportVerifyTaskDetail)     // 报告验证 - 任务详情
	smartRouterGroup.GET("/reportverify/targetlist", rest.ReportVerifyTargetList)     // 报告验证 - 目标列表
	smartRouterGroup.GET("/reportverify/portlist", rest.ReportVerifyPortList)         // 报告验证 - 端口列表
	smartRouterGroup.GET("/reportverify/vullist", rest.ReportVerifyVulList)           // 报告验证 - 漏洞列表
	smartRouterGroup.GET("/reportverify/enum", rest.ReportVerifyEnum)                 // 报告验证 - 枚举接口
	smartRouterGroup.GET("/reportverify/stats", rest.ReportVerifyStatsInfo)           // 报告验证 - 统计信息
	smartRouterGroup.GET("/reportverify/taskstop", rest.ReportVerifyTaskStop)         // 报告验证 - 结束任务
	smartRouterGroup.GET("/reportverify/taskdelete", rest.ReportVerifyTaskDelete)     // 报告验证 - 删除任务
	smartRouterGroup.GET("/reportverify/targetdelete", rest.ReportVerifyTargetDelete) // 报告验证 - 删除目标
	smartRouterGroup.GET("/reportverify/vuldelete", rest.ReportVerifyVulDelete)       // 报告验证 - 删除漏洞
	smartRouterGroup.GET("/reportverify/vuldetail", rest.ReportVerifyVulDetail)       // 报告验证 - 漏洞详情

	//任务组
	smartRouterGroup.POST("/taskgroup/create", rest.TaskGroupCreate)       // 任务组 - 任务组新建
	smartRouterGroup.GET("/taskgroup/list", rest.TaskGroupList)            // 任务组 - 任务组列表
	smartRouterGroup.GET("/taskgroup/delete", rest.TaskGroupDelete)        // 任务组 - 任务组删除
	smartRouterGroup.POST("/taskgroup/groupbind", rest.TaskGroupGroupBind) // 任务组 - 任务组内任务新建
	smartRouterGroup.GET("/taskgroup/tasklist", rest.TaskGroupTaskList)    // 任务组 - 任务组内任务列表
	smartRouterGroup.GET("/taskgroup/overview", rest.TaskGroupOverView)    // 任务组 - 任务组内任务统计
	smartRouterGroup.GET("/taskgroup/status", rest.TaskGroupStatus)        // 任务组 - 任务组内任务状态
	smartRouterGroup.POST("/taskgroup/groupedit", rest.TaskGroupEdit)      // 任务组 - 任务组编辑功能

	// 资产中心 - 资产树【资产组】
	smartRouterGroup.GET("/asset/statistics", rest.Statistics)         // 资产中心 - 资产概览统计
	smartRouterGroup.GET("/assetgroup/enums", rest.GetAssetGroupEnums) // 资产中心 - 资产组枚举
	smartRouterGroup.GET("/assettree/list", rest.AssetTree)            // 资产中心 - 资产树列表
	smartRouterGroup.GET("/assetgroup/list", rest.AssetGroupList)      // 资产中心 - 资产组列表
	smartRouterGroup.POST("/assetgroup/add", rest.AssetGroupAdd)       // 资产中心 - 资产组新增
	smartRouterGroup.GET("/assetgroup/info", rest.AssetGroupInfo)      // 资产中心 - 资产组详情
	smartRouterGroup.POST("/assetgroup/edit", rest.AssetGroupEdit)     // 资产中心 - 资产组编辑
	smartRouterGroup.GET("/assetgroup/del", rest.AssetGroupDel)        // 资产中心 - 资产组删除
	// 资产中心 - 资产信息【列表】
	smartRouterGroup.GET("/asset/list", rest.AssetList)           // 资产中心 - 资产信息列表
	smartRouterGroup.GET("/asset/enums", rest.GetAssetEnums)      // 资产中心 - 资产信息枚举
	smartRouterGroup.POST("/asset/add", rest.AssetAdd)            // 资产中心 - 资产信息新增
	smartRouterGroup.POST("/asset/edit", rest.AssetEdit)          // 资产中心 - 资产信息编辑
	smartRouterGroup.GET("/asset/del", rest.AssetDel)             // 资产中心 - 资产信息删除
	smartRouterGroup.GET("/asset/detail", rest.AssetDetail)       // 资产中心 - 资产信息资产详情
	smartRouterGroup.POST("/asset/import", rest.Import)           // 资产中心 - 资产信息导入
	smartRouterGroup.POST("/asset/export", rest.Export)           // 资产中心 - 资产信息导出
	smartRouterGroup.GET("/asset/selectall", rest.SelectAllAsset) // 资产中心 - 资产信息列表全选功能
	smartRouterGroup.GET("/assetconn/list", rest.AssetConnList)   // 资产中心 - 资产连接方式

	//逻辑漏洞检测模块
	smartRouterGroup.POST("/logic/taskcreate", rest.LogicTaskCreate)        // 逻辑漏洞 - 任务新建
	smartRouterGroup.GET("/logic/taskstop", rest.LogicTaskStop)             // 逻辑漏洞 - 任务结束
	smartRouterGroup.GET("/logic/taskdel", rest.LogicTaskDel)               // 逻辑漏洞 - 任务删除
	smartRouterGroup.GET("/logic/tasklist", rest.LogicTaskList)             // 逻辑漏洞 - 任务列表
	smartRouterGroup.GET("/logic/targetlist", rest.LogicTargetList)         // 逻辑漏洞 - 目标列表
	smartRouterGroup.GET("/logic/vullist", rest.LogicVulList)               // 逻辑漏洞 - 任务列表
	smartRouterGroup.GET("/logic/loglist", rest.LogicLogList)               // 逻辑漏洞 - 日志列表
	smartRouterGroup.GET("/logic/vulinfo", rest.LogicVulInfo)               // 逻辑漏洞 - 漏洞详情
	smartRouterGroup.GET("/logic/loginfo", rest.LogicLogInfo)               // 逻辑漏洞 - 日志详情
	smartRouterGroup.GET("/logic/enum", rest.LogicEnum)                     // 逻辑漏洞 - 参数枚举
	smartRouterGroup.GET("/logic/vuldel", rest.LogicVulDelete)              // 逻辑漏洞 - 漏洞删除
	smartRouterGroup.POST("/logic/vultest", rest.LogicVulTest)              // 逻辑漏洞 - 漏洞删除
	smartRouterGroup.GET("/logic/taskcopy", rest.LogicTaskCopy)             // 逻辑漏洞 - 任务复制
	smartRouterGroup.POST("/logic/reportsave", rest.LogicReportSave)        // 逻辑漏洞 - 报告生成
	smartRouterGroup.GET("/logic/flowbaselist", rest.LogicFlowBaseList)     // 逻辑漏洞 - 流量接口
	smartRouterGroup.GET("/logic/flowbaseinfo", rest.LogicFlowBaseInfo)     // 逻辑漏洞 - 流量详情
	smartRouterGroup.GET("/logic/flowbaseexport", rest.LogicFlowBaseExport) // 逻辑漏洞 - 流量导出

	// 漏洞扫描匹配模块
	smartRouterGroup.GET("/vulscan/tasklist", rest.VulScanTaskList)         // 漏洞扫描 - 任务列表
	smartRouterGroup.POST("/vulscan/tasksave", rest.VulScanTaskSave)        // 漏洞扫描 - 任务创建
	smartRouterGroup.GET("/vulscan/taskstop", rest.VulScanTaskStop)         // 漏洞扫描 - 任务结束
	smartRouterGroup.GET("/vulscan/taskdelete", rest.VulScanTaskDelete)     // 漏洞扫描 - 任务删除
	smartRouterGroup.GET("/vulscan/targetlist", rest.VulScanTargetList)     // 漏洞扫描 - 目标列表
	smartRouterGroup.GET("/vulscan/vullist", rest.VulScanVulList)           // 漏洞扫描 - 漏洞列表
	smartRouterGroup.GET("/vulscan/vuldetail", rest.VulScanVulDetail)       // 漏洞扫描 - 漏洞详情
	smartRouterGroup.GET("/vulscan/cvelist", rest.VulScanCveList)           // 漏洞扫描 - cve列表
	smartRouterGroup.GET("/vulscan/cvedetail", rest.VulScanCveDetail)       // 漏洞扫描 - cve详情
	smartRouterGroup.GET("/vulscan/taskoverview", rest.VulScanTaskOverview) // 漏洞扫描 - 任务详情
	smartRouterGroup.GET("/vulscan/taskstate", rest.VulScanTaskState)       // 漏洞扫描 - 任务状态

	// 风险管理 - 漏洞风险
	smartRouterGroup.GET("/riskvul/statistics", rest.VulRiskStatistics) // 漏洞风险 - 漏洞风险统计
	smartRouterGroup.GET("/riskvul/list", rest.VulRiskList)             // 漏洞风险 - 漏洞风险列表
	smartRouterGroup.GET("/riskvul/updatestatus", rest.UpdateVulRisk)   // 漏洞风险 - 漏洞风险状态 包括已修复等
	smartRouterGroup.GET("/riskvul/delete", rest.DelVulRisk)            // 漏洞风险 - 删除漏洞风险
	smartRouterGroup.GET("/riskvul/detail", rest.VulRiskDetail)         // 漏洞风险 - 漏洞风险详情
	smartRouterGroup.GET("/riskvul/test", rest.VulRiskTest)             // 漏洞风险 - 漏洞风险测试
	smartRouterGroup.GET("/riskvul/verify", rest.VulRiskVerify)         // 漏洞风险 - 漏洞风险验证
	smartRouterGroup.GET("/riskvul/enums", rest.RiskManageEnum)         // 漏洞风险 - 漏洞风险枚举

	// 安全检查 - 安全配置核查
	smartRouterGroup.POST("/baseline/check", rest.BaselineCheckRun)                   // 安全配置核查 - 执行主机基线检查（单目标）
	smartRouterGroup.POST("/baseline/check/batch", rest.BaselineBatchCheckRun)        // 安全配置核查 - 批量多目标核查（异步）
	smartRouterGroup.GET("/baseline/check/progress", rest.BaselineBatchCheckProgress) // 安全配置核查 - 批量任务进度
	smartRouterGroup.GET("/baseline/result", rest.BaselineCheckResultList)            // 安全配置核查 - 检查结果列表
	smartRouterGroup.GET("/baseline/stat", rest.BaselineCheckStat)                    // 安全配置核查 - 检查统计
	smartRouterGroup.GET("/baseline/tasks", rest.BaselineTaskList)                    // 安全配置核查 - 核查批次列表
	smartRouterGroup.GET("/baseline/task/targets", rest.BaselineTaskTargets)          // 安全配置核查 - 任务目标列表
	smartRouterGroup.GET("/baseline/rules", rest.BaselineRulesList)                   // 安全配置核查 - 规则列表
	smartRouterGroup.POST("/baseline/rules/reload", rest.BaselineRulesReload)         // 安全配置核查 - 从库重载规则（导入 host_baseline_rule 后调用）
	smartRouterGroup.POST("/baseline/rules/import", rest.BaselineRulesImport)         // 安全配置核查 - 导入规则（JSON格式）
	smartRouterGroup.GET("/baseline/rules/db", rest.BaselineRulesListFromDB)          // 安全配置核查 - 从数据库获取规则列表
	smartRouterGroup.GET("/baseline/rule/detail", rest.BaselineRuleDetail)            // 安全配置核查 - 规则详情
	smartRouterGroup.POST("/baseline/rule/create", rest.BaselineRuleCreate)           // 安全配置核查 - 新增规则
	smartRouterGroup.POST("/baseline/rule/update", rest.BaselineRuleUpdate)           // 安全配置核查 - 编辑规则
	smartRouterGroup.GET("/baseline/rule/delete", rest.BaselineRuleDelete)            // 安全配置核查 - 删除规则
	// 安全检查 - CVE漏洞扫描（基于软件版本匹配）
	smartRouterGroup.POST("/vulnscan/cve/run", rest.VulnScanCveRun)          // CVE漏洞扫描 - 单目标执行
	smartRouterGroup.POST("/vulnscan/cve/batch", rest.VulnScanCveBatchRun)   // CVE漏洞扫描 - 批量多目标（异步）
	smartRouterGroup.GET("/vulnscan/cve/progress", rest.VulnScanCveProgress) // CVE漏洞扫描 - 批量任务进度
	// 安全检查 - 恶意代码检测
	smartRouterGroup.POST("/malware/scan", rest.MalwareScanRun)      // 恶意代码检测 - 执行扫描
	smartRouterGroup.GET("/malware/result", rest.MalwareResultList)  // 恶意代码检测 - 扫描结果
	smartRouterGroup.GET("/malware/tasks", rest.MalwareScanTaskList) // 恶意代码检测 - 按批次聚合列表
	// 安全检查 - YARA恶意代码检测（基于YARA规则引擎）
	smartRouterGroup.POST("/malware/yara/run", rest.MalwareYaraRun)          // YARA恶意代码检测 - 单目标执行
	smartRouterGroup.POST("/malware/yara/batch", rest.MalwareYaraBatchRun)   // YARA恶意代码检测 - 批量多目标（异步）
	smartRouterGroup.GET("/malware/yara/progress", rest.MalwareYaraProgress) // YARA恶意代码检测 - 批量任务进度
	smartRouterGroup.GET("/malware/yara/result", rest.MalwareYaraResultList) // YARA恶意代码检测 - 扫描结果
	smartRouterGroup.GET("/malware/yara/tasks", rest.MalwareYaraTaskList)    // YARA恶意代码检测 - 按批次聚合列表
	// 安全检查 - 病毒库规则管理
	smartRouterGroup.POST("/malware/rule/create", rest.MalwareRuleCreate) // 病毒库规则 - 新增
	smartRouterGroup.POST("/malware/rule/update", rest.MalwareRuleUpdate) // 病毒库规则 - 编辑
	smartRouterGroup.GET("/malware/rule/delete", rest.MalwareRuleDelete)  // 病毒库规则 - 删除
	smartRouterGroup.GET("/malware/rule/detail", rest.MalwareRuleDetail)  // 病毒库规则 - 详情
	smartRouterGroup.GET("/malware/rules", rest.MalwareRuleList)          // 病毒库规则 - 列表
	smartRouterGroup.POST("/malware/rule/import", rest.MalwareRuleImport) // 病毒库规则 - 导入（.yar / .json）
	// 安全检查 - CVE漏洞库查询
	smartRouterGroup.GET("/cvedb/info", rest.CveDBInfo)   // CVE库 - 统计信息
	smartRouterGroup.GET("/cvedb/query", rest.CveDBQuery) // CVE库 - 搜索
	// 安全检查 - 数据库安全检查
	smartRouterGroup.POST("/db/check", rest.DBCheckRun)        // 数据库检查 - 执行数据库安全基线检查
	smartRouterGroup.GET("/db/result", rest.DBCheckResultList) // 数据库检查 - 检查结果列表
	// 安全检查 - 敏感数据发现
	smartRouterGroup.POST("/sensitive/scan", rest.SensitiveDataScanRun)     // 敏感数据发现 - 执行扫描
	smartRouterGroup.GET("/sensitive/result", rest.SensitiveDataResultList) // 敏感数据发现 - 扫描结果
	smartRouterGroup.GET("/sensitive/stat", rest.SensitiveDataStat)         // 敏感数据发现 - 统计
	// 安全检查通用枚举
	smartRouterGroup.GET("/baseline/enums", rest.BaselineEnums) // 安全检查通用枚举

	// 应用安全扫描（Web 动态扫描 / 专项应用检测，调用 scanner.exe）
	smartRouterGroup.POST("/appsec/dynamic/run", rest.AppSecDynamicScanRun)            // 应用安全 - 动态扫描创建
	smartRouterGroup.GET("/appsec/dynamic/list", rest.AppSecDynamicScanList)           // 应用安全 - 动态扫描列表
	smartRouterGroup.GET("/appsec/dynamic/detail", rest.AppSecDynamicScanDetail)       // 应用安全 - 动态扫描详情
	smartRouterGroup.POST("/appsec/appspecific/run", rest.AppSecAppSpecificScanRun)    // 应用安全 - 专项应用检测创建
	smartRouterGroup.GET("/appsec/appspecific/list", rest.AppSecAppSpecificScanList)   // 应用安全 - 专项应用检测列表
	smartRouterGroup.GET("/appsec/appspecific/detail", rest.AppSecAppSpecificScanDetail) // 应用安全 - 专项应用检测详情

	// 数据安全（任务化：数据库基线检查 / 敏感数据发现）
	smartRouterGroup.POST("/datasec/db/test-conn", rest.DataSecDBTestConn)           // 数据安全 - 数据库连接测试
	smartRouterGroup.POST("/datasec/db/run", rest.DataSecDBCheckRun)                 // 数据安全 - 数据库检查创建
	smartRouterGroup.GET("/datasec/db/list", rest.DataSecDBCheckList)               // 数据安全 - 数据库检查任务列表
	smartRouterGroup.GET("/datasec/db/detail", rest.DataSecDBCheckDetail)           // 数据安全 - 数据库检查任务详情
	smartRouterGroup.POST("/datasec/sensitive/run", rest.DataSecSensitiveScanRun)   // 数据安全 - 敏感数据扫描创建
	smartRouterGroup.GET("/datasec/sensitive/list", rest.DataSecSensitiveScanList)   // 数据安全 - 敏感数据任务列表
	smartRouterGroup.GET("/datasec/sensitive/detail", rest.DataSecSensitiveScanDetail) // 数据安全 - 敏感数据任务详情
	smartRouterGroup.GET("/datasec/task/clone-targets", rest.DataSecTaskCloneTargets)   // 数据安全 - 复制历史任务目标
	smartRouterGroup.POST("/datasec/task/rerun", rest.DataSecTaskRerun)                 // 数据安全 - 再次检测
	smartRouterGroup.GET("/datasec/task/delete", rest.DataSecTaskDelete)                 // 数据安全 - 删除任务
	smartRouterGroup.GET("/datasec/target/list", rest.DataSecDBTargetList)              // 数据安全 - 目标库列表
	smartRouterGroup.POST("/datasec/target/save", rest.DataSecDBTargetSave)             // 数据安全 - 目标库保存
	smartRouterGroup.GET("/datasec/target/delete", rest.DataSecDBTargetDelete)          // 数据安全 - 目标库删除
	smartRouterGroup.POST("/datasec/target/import", rest.DataSecDBTargetImport)         // 数据安全 - 目标库导入
	smartRouterGroup.GET("/datasec/target/export", rest.DataSecDBTargetExport)          // 数据安全 - 目标库导出
	smartRouterGroup.POST("/datasec/target/save-from-task", rest.DataSecSaveTargetsToLibrary) // 保存任务目标到库
	smartRouterGroup.POST("/datasec/target/test-conn", rest.DataSecDBTargetTestConn)         // 目标库单条连接测试
	smartRouterGroup.POST("/datasec/target/batch-test-conn", rest.DataSecDBTargetBatchTestConn) // 目标库批量连接测试
	// 数据安全 - 检测规则库
	smartRouterGroup.GET("/datasec/rules", rest.DataSecRulesList)                        // 规则列表与统计
	smartRouterGroup.POST("/datasec/rules/reload", rest.DataSecRulesReload)            // 从库重载规则
	smartRouterGroup.POST("/datasec/rules/import", rest.DataSecRulesImport)             // JSON 导入规则
	smartRouterGroup.POST("/datasec/rules/import-builtin", rest.DataSecRulesImportBuiltin) // 导入内置规则到库
	smartRouterGroup.GET("/datasec/rules/cve-preview", rest.DataSecCveImportPreview)       // CVE 库可导入条数预览
	smartRouterGroup.POST("/datasec/rules/import-cve", rest.DataSecRulesImportFromCve)       // 从 default-cve.db 导入
	smartRouterGroup.GET("/datasec/rule/detail", rest.DataSecRuleDetail)               // 规则详情
	smartRouterGroup.POST("/datasec/rule/create", rest.DataSecRuleCreate)              // 新增规则
	smartRouterGroup.POST("/datasec/rule/update", rest.DataSecRuleUpdate)              // 编辑规则
	smartRouterGroup.GET("/datasec/rule/delete", rest.DataSecRuleDelete)               // 删除规则

	return router
}
