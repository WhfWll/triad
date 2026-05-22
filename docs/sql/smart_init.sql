-- Database: smart
-- All tables use CREATE TABLE IF NOT EXISTS so the script can be re-run without error 1050.
-- Note: existing tables are not altered; use migrations to change column definitions.
CREATE DATABASE IF NOT EXISTS `smart` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `smart`;

CREATE TABLE IF NOT EXISTS `asset` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `asset_group_id` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '所属资产组，1表示未分组',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT '资产ip或域名',
  `ip_num` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT 'ip值',
  `ip_segment` varchar(255) NOT NULL COMMENT 'Ip段',
  `asset_type` int(11) NOT NULL DEFAULT '0' COMMENT '资产类型',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '资产名称',
  `operate_system` varchar(50) NOT NULL DEFAULT '' COMMENT '操作系统',
  `location` varchar(50) NOT NULL DEFAULT '' COMMENT '地理位置',
  `risk_level` int(10) unsigned NOT NULL DEFAULT '4' COMMENT '风险等级，1&2-高危,3-中危,4-低危,0-安全',
  `business_system` varchar(50) NOT NULL DEFAULT '' COMMENT '业务系统',
  `responsible_department` varchar(50) NOT NULL DEFAULT '' COMMENT '责任部门',
  `filing_level` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '备案等级，1-无,2-等保一级,3-等保二级,4-等保三级,5-等保四级,6-等保五级',
  `device_weight` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '设备权重 1高 2中 3低 4极高 5极低',
  `trust_level` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '可信权重 1 可信 2 未登记',
  `is_cloud_host` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否是云主机 true是 false否',
  `tags` text NOT NULL COMMENT '标签，多个用英文逗号隔开',
  `info` text NOT NULL COMMENT '资产其他信息,JSON格式，硬件/主机名/资产类型/虚拟资产/归属地/责任人/责任人邮箱',
  `is_ignore` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '最新资产排序是否忽略，1-不忽略，2-忽略',
  `asset_changes_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '资产变化类型:1-未变化,2-已减少IP，3-新增加IP，4-端口变化IP，5-服务变化IP，6-组件变化IP',
  `islive` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '存活状态:1-存活,2-不存活IP',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `target_ids` varchar(255) NOT NULL DEFAULT '' COMMENT '渗透目标id,多个用英文逗号隔开',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态，1-待同步、2-同步中、3-已完成',
  `risk_num` varchar(50) NOT NULL DEFAULT '' COMMENT '记录资产最近一次任务的风险漏洞情况',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_ip` (`ip`) USING BTREE,
  KEY `idx_update_time` (`update_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=400 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='资产表';

CREATE TABLE IF NOT EXISTS `asset_connection` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `asset_id` int(10) unsigned DEFAULT NULL COMMENT '关联资产ID，可为空',
  `ip` varchar(45) NOT NULL COMMENT '资产IP，唯一索引',
  `port` int(11) NOT NULL COMMENT '端口',
  `protocol` int(11) NOT NULL COMMENT '协议/连接方式，如 1 ssh、rdp、mysql',
  `username` varchar(100) NOT NULL COMMENT '用户名',
  `password` text NOT NULL COMMENT '加密存储的密码',
  `remark` varchar(255) DEFAULT NULL COMMENT '备注',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_ip` (`ip`),
  KEY `idx_asset_id` (`asset_id`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COMMENT='资产连接信息表';

CREATE TABLE IF NOT EXISTS `asset_group` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父级ID',
  `level` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '级别，1一级资产，2二级资产...最多六级',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '资产组名称',
  `remark` text NOT NULL COMMENT '资产组说明',
  `user_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '提交用户ID',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_pid_name` (`pid`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='资产组表';

CREATE TABLE IF NOT EXISTS `asset_log` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '日志类型，1-创建资产渗透',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='资产日志表';

CREATE TABLE IF NOT EXISTS `asset_port` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT 'ip',
  `port` varchar(50) NOT NULL DEFAULT '' COMMENT '端口',
  `protocol` varchar(50) NOT NULL DEFAULT '' COMMENT '协议',
  `service` varchar(100) NOT NULL DEFAULT '' COMMENT '服务',
  `assembly` text NOT NULL COMMENT '组件',
  `remark` text NOT NULL COMMENT '备注',
  `islive` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '存活状态:1-存活,2-不存活',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '任务ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_ip_port` (`ip`,`port`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1741 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='资产ip端口表';

CREATE TABLE IF NOT EXISTS `asset_risk_trend` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `asset_id` int(10) unsigned NOT NULL DEFAULT '1' COMMENT '资产ID',
  `ip` varchar(50) NOT NULL DEFAULT '' COMMENT '资产ip或域名',
  `risk_level` int(10) unsigned NOT NULL DEFAULT '4' COMMENT '风险等级，1&2-高危,3-中危,4-低危,0-安全',
  `risk_num` varchar(50) NOT NULL DEFAULT '' COMMENT '记录资产最近一次任务的风险漏洞情况',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '渗透目标id',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务ID',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ip` (`ip`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=768 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='资产风险趋势表';

CREATE TABLE IF NOT EXISTS `asset_task_result` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_result_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '所属任务结果id',
  `asset_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '所属资产id',
  `sub_obj_type` varchar(50) NOT NULL DEFAULT '' COMMENT '数据子类型',
  `obj_id` varchar(50) NOT NULL DEFAULT '' COMMENT '数据对象id',
  `sub_obj_id` varchar(50) NOT NULL DEFAULT '' COMMENT '数据子对象id',
  `identify` varchar(30) NOT NULL DEFAULT '' COMMENT '数据标识符,用来修改时方便定位数据的',
  `field1` varchar(255) NOT NULL DEFAULT '' COMMENT '筛选字段1',
  `field2` varchar(255) NOT NULL DEFAULT '' COMMENT '筛选字段2',
  `field3` varchar(255) NOT NULL DEFAULT '' COMMENT '筛选字段3',
  `field4` varchar(255) NOT NULL DEFAULT '' COMMENT '筛选字段4',
  `json_result` text NOT NULL COMMENT '各数据类型的json格式结果',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_asset_id` (`asset_id`) USING BTREE,
  KEY `idx_field1` (`field1`) USING BTREE,
  KEY `idx_field2` (`field2`) USING BTREE,
  KEY `idx_field3` (`field3`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='资产结果表';

CREATE TABLE IF NOT EXISTS `asset_vul` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_vul_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '所属任务漏洞id',
  `asset_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '所属资产id',
  `task_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '所属目标id',
  `target_url` varchar(255) NOT NULL DEFAULT '' COMMENT '测试目标地址',
  `pocname` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞标识',
  `name` varchar(200) NOT NULL DEFAULT '' COMMENT '漏洞名称',
  `class` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '漏洞分类',
  `type` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '漏洞类型',
  `risk` int(10) unsigned NOT NULL DEFAULT '4' COMMENT '风险等级，1-致命/2-高危/3-中危/4-低危/5-信息',
  `location` varchar(1000) NOT NULL DEFAULT '' COMMENT '漏洞位置或地址',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '漏洞状态，1-未验证，2-验证成功，3-利用成功，4-已修复，5-流量验证,6-未知,7-存在,8-不存在',
  `test_status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '测试状态，1-未测试，2-测试中，3-已测试',
  `exploit_impact` varchar(10) NOT NULL DEFAULT '' COMMENT '利用影响',
  `vul_id` varchar(30) NOT NULL DEFAULT '' COMMENT '漏洞id',
  `description` text NOT NULL COMMENT '漏洞描述',
  `fix_suggest` varchar(512) NOT NULL DEFAULT '' COMMENT '修复建议',
  `published_time` varchar(32) NOT NULL DEFAULT '' COMMENT '披露时间',
  `affect_range` varchar(255) NOT NULL DEFAULT '' COMMENT '影响范围',
  `target_result_id` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '检测结果id',
  `vul_number` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞编号',
  `vul_address` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞地址',
  `ref_url` varchar(255) NOT NULL DEFAULT '' COMMENT '参考链接',
  `cvss` varchar(255) NOT NULL DEFAULT '' COMMENT 'cvss评分',
  `vul_result` text NOT NULL COMMENT '漏洞结果',
  `vul_param` text NOT NULL COMMENT '漏洞请求参数',
  `ver_msg` text NOT NULL COMMENT '验证报文',
  `is_replace` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '是否替换漏洞，1-否，2-是',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_asset_id` (`asset_id`) USING BTREE,
  KEY `idx_target_id` (`target_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=14198 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='资产漏洞表';

CREATE TABLE IF NOT EXISTS `bas_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `bas_task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'bas任务ID',
  `bas_target_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'bas任务目标agent ID',
  `content` varchar(1024) NOT NULL DEFAULT '' COMMENT '日志内容',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_bas_target_id` (`bas_target_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='bas任务目标agent日志';

CREATE TABLE IF NOT EXISTS `bas_node` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '节点名称',
  `ip` varchar(20) NOT NULL DEFAULT '' COMMENT 'IP',
  `online_status` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '在线状态 1-在线 2-离线',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '节点状态 1-启用 2-禁用',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='bas节点表';

CREATE TABLE IF NOT EXISTS `bas_rules` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `content` varchar(2560) NOT NULL DEFAULT '' COMMENT '规则内容',
  `name` varchar(512) NOT NULL DEFAULT '' COMMENT '规则名称',
  `protocol` varchar(512) NOT NULL DEFAULT '' COMMENT '协议名称',
  `keywords` varchar(2560) NOT NULL DEFAULT '' COMMENT '关键字',
  `keywords_zh` varchar(2560) NOT NULL DEFAULT '' COMMENT '关键字中文描述',
  `description` varchar(2560) NOT NULL DEFAULT '' COMMENT '描述',
  `description_zh` varchar(2560) NOT NULL DEFAULT '' COMMENT '描述中文',
  `cve` varchar(255) NOT NULL DEFAULT '' COMMENT 'cve编号',
  `raw_traffic_beyond_ip_packet_base64` text COMMENT '基于ip数据包的原始流量',
  `raw_traffic_beyond_http_base64` text COMMENT 'http 之外的原始流量',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `name_zh` varchar(512) NOT NULL DEFAULT '' COMMENT '规则名称',
  `class_type` varchar(128) NOT NULL DEFAULT '' COMMENT '规则分类',
  `hash` varchar(128) NOT NULL DEFAULT '' COMMENT 'hash',
  `attack_stage` tinyint(4) unsigned NOT NULL COMMENT '攻击阶段',
  `risk_level` tinyint(4) unsigned NOT NULL DEFAULT '4' COMMENT '任务风险等级，1-高危、2-中危、3-低危、4-安全',
  `effect_target` varchar(128) NOT NULL DEFAULT '' COMMENT '影响目标',
  `effect_score` varchar(12) NOT NULL DEFAULT '' COMMENT '影响评分',
  `relation_attack_method` varchar(1024) NOT NULL DEFAULT '' COMMENT '关联攻击方式',
  `ref_url` varchar(255) NOT NULL DEFAULT '' COMMENT '修复建议',
  `fix_suggest` text COMMENT '修复建议',
  `attack_mode` int(11) NOT NULL DEFAULT '0' COMMENT '攻击方式',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4508 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC COMMENT='bas规则库';

CREATE TABLE IF NOT EXISTS `bas_target` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `bas_task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'bas任务ID',
  `addr` varchar(36) NOT NULL DEFAULT '' COMMENT '地址',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1待开始 2已开始 3已完成',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `bas_template_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '方案id',
  `bas_template_json` text NOT NULL COMMENT '方案数据json',
  PRIMARY KEY (`id`),
  KEY `idx_bas_task_id` (`bas_task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='bas任务目标agent';

CREATE TABLE IF NOT EXISTS `bas_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(36) NOT NULL DEFAULT '' COMMENT '名称',
  `bas_template_id` int(11) NOT NULL DEFAULT '0' COMMENT '方案ID',
  `bas_template_json` text COMMENT '方案数据json',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态：1待开始 2已开始 3已完成',
  `risk_level` tinyint(4) unsigned NOT NULL DEFAULT '4' COMMENT '任务风险等级，1-高危、2-中危、3-低危、4-安全',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `user` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='bas任务';

CREATE TABLE IF NOT EXISTS `bas_template` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(36) NOT NULL DEFAULT '' COMMENT '名称',
  `desc` varchar(256) NOT NULL DEFAULT '' COMMENT '方案描述',
  `rule_ids` text COMMENT '规则ID',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `is_default` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否默认 0否 1是',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COMMENT='bas规则集合 剧本集';

CREATE TABLE IF NOT EXISTS `bas_vul` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `bas_task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'bas任务ID',
  `bas_target_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'bas目标ID',
  `addr` varchar(36) NOT NULL DEFAULT '' COMMENT '攻击目标',
  `rule_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'bas目标ID',
  `rule_name` varchar(255) NOT NULL DEFAULT '' COMMENT '剧本名称',
  `attack_mode` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '攻击方式',
  `attack_stage` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '攻击阶段',
  `risk_level` tinyint(4) unsigned NOT NULL DEFAULT '4' COMMENT '任务风险等级，1-高危、2-中危、3-低危、4-安全',
  `md5_code` text NOT NULL COMMENT 'md5值，英文逗号隔开',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '攻击结果：1失败 2成功',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_bas_target_id` (`bas_target_id`),
  KEY `idx_addr` (`addr`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8 COMMENT='bas漏洞表';

CREATE TABLE IF NOT EXISTS `burpsuite_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `origin_task_id` int(11) NOT NULL COMMENT 'burpsuite软件生成的任务ID',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `target` varchar(5120) NOT NULL DEFAULT '' COMMENT '目标',
  `risk` tinyint(4) NOT NULL DEFAULT '0' COMMENT '风险等级',
  `is_crawler` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '是否开启爬虫 0为开始 1默认开启（调用API后它是默认开启的，不能更改）',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '状态 1待运行 2运行中 3已完成',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='burpsuite任务表';

CREATE TABLE IF NOT EXISTS `burpsuite_task_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `burpsuite_task_id` int(11) NOT NULL DEFAULT '0' COMMENT 'burpsuite任务表ID',
  `origin_result_id` int(11) NOT NULL DEFAULT '0' COMMENT 'burpsuite软件生成的任务结果ID',
  `action` varchar(50) NOT NULL DEFAULT '' COMMENT '对应json结果的外层type',
  `issue_type` varchar(128) NOT NULL DEFAULT '' COMMENT '问题类型 取原数据的issue.name',
  `host` varchar(50) NOT NULL DEFAULT '' COMMENT '主机',
  `path` varchar(128) NOT NULL DEFAULT '' COMMENT '路径',
  `insertion_point` varchar(128) NOT NULL DEFAULT '' COMMENT '插入点',
  `severity` varchar(128) NOT NULL DEFAULT '' COMMENT '严重程度 也是风险等级 取原数据的即可',
  `confidence` varchar(128) NOT NULL DEFAULT '' COMMENT 'confidence 取原数据的即可',
  `describe` varchar(5120) NOT NULL DEFAULT '' COMMENT '漏洞描述 取原数据即可',
  `issue_background` varchar(5120) NOT NULL DEFAULT '' COMMENT '问题背景 取原数据即可',
  `remediation_background` varchar(5120) NOT NULL DEFAULT '' COMMENT '补救背景 取原数据即可',
  `request_response` mediumtext COMMENT '请求与响应信息 取原数据即可',
  `internal_data` mediumtext COMMENT '内部数据 取原数据即可',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='burpsuite任务结果表';

CREATE TABLE IF NOT EXISTS `dictionary` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `sources` tinyint(3) unsigned NOT NULL DEFAULT '2' COMMENT '来源，1系统，2手动',
  `service` tinyint(3) unsigned NOT NULL DEFAULT '2' COMMENT '适用范围',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '字典名称',
  `types` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '字典类型：1-用户字典，2-密码字典，3-wifi，4-路径字典，5-子域名字典',
  `is_default` tinyint(3) unsigned NOT NULL DEFAULT '2' COMMENT '是否是默认，1-是，2-否',
  `content` longtext NOT NULL COMMENT '字典内容,换行符隔开',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=93 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='字典表';

CREATE TABLE IF NOT EXISTS `finger` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_name` varchar(200) NOT NULL DEFAULT '' COMMENT '应用名称',
  `cn_name` varchar(200) NOT NULL DEFAULT '' COMMENT '中文名称',
  `app_version` varchar(100) NOT NULL DEFAULT '' COMMENT '应用版本',
  `source` tinyint(4) NOT NULL DEFAULT '1' COMMENT '数据来源: 1-系统自带, 2-用户添加',
  `level` varchar(10) NOT NULL DEFAULT '0' COMMENT '服务分层	硬件层、系统层、服务层、支撑层、应用层',
  `finger_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '指纹类型: 1-Web指纹, 2-服务指纹, 3-硬件指纹, 4-其他',
  `app_class` int(11) NOT NULL DEFAULT '0' COMMENT '指纹分类',
  `flag` varchar(2000) NOT NULL DEFAULT '' COMMENT '匹配内容 指纹规则',
  `desc` varchar(500) DEFAULT NULL COMMENT '指纹描述',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_app_name` (`app_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=12907 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='指纹表';

CREATE TABLE IF NOT EXISTS `flow_base` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `flow_task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属流量分析任务id',
  `flow_target_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属流量分析目标id',
  `hash` varchar(64) NOT NULL DEFAULT '' COMMENT '唯一标识，用于区分数据是否重复',
  `yak_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'yak数据ID，用于增量数据标识',
  `host` varchar(64) NOT NULL DEFAULT '' COMMENT '请求地址',
  `ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `method` varchar(12) NOT NULL DEFAULT '' COMMENT '请求方法',
  `protocol` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '网络协议，1-http，2-https',
  `resp_title` varchar(128) NOT NULL DEFAULT '' COMMENT '响应title',
  `tags` varchar(255) NOT NULL DEFAULT '' COMMENT '标签',
  `resp_code` varchar(12) NOT NULL DEFAULT '' COMMENT '响应码',
  `resp_content_type` varchar(128) NOT NULL DEFAULT '' COMMENT '响应类型',
  `req_header` text NOT NULL COMMENT '源请求头',
  `resp_header` text NOT NULL COMMENT '源响应头',
  `resp_content` longtext NOT NULL COMMENT '源响应数据',
  `like_field` varchar(128) NOT NULL DEFAULT '' COMMENT '冗余like字段',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `url` varchar(500) NOT NULL DEFAULT '' COMMENT '请求url',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_hash` (`hash`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=540 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='流量分析被动流量表';

CREATE TABLE IF NOT EXISTS `flow_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `flow_task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属流量分析任务id',
  `content` text NOT NULL COMMENT '日志内容',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_flow_task_id` (`flow_task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1645 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='流量分析日志表';

CREATE TABLE IF NOT EXISTS `flow_risk` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `flow_task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属流量分析任务id',
  `flow_target_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属流量分析目标id',
  `hash` varchar(64) NOT NULL DEFAULT '' COMMENT '唯一标识，用于区分数据是否重复',
  `yak_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'yak数据ID，用于增量数据标识',
  `host` varchar(64) NOT NULL DEFAULT '' COMMENT '请求地址',
  `ip` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP地址',
  `ip_integer` varchar(64) NOT NULL DEFAULT '' COMMENT 'IP？Yak的记录，不知道是啥',
  `port` varchar(36) NOT NULL DEFAULT '' COMMENT '端口',
  `title` varchar(512) NOT NULL DEFAULT '' COMMENT '漏洞标题',
  `risk_type` varchar(128) NOT NULL DEFAULT '' COMMENT '漏洞类型',
  `risk_type_verbose` varchar(128) NOT NULL DEFAULT '' COMMENT '漏洞类型冗长',
  `payload` text NOT NULL COMMENT '内容',
  `detail` text NOT NULL COMMENT '详情',
  `risk_level` tinyint(4) unsigned NOT NULL DEFAULT '5' COMMENT '风险等级,1-致命,2-高危,3-中危,4-低位,5-安全',
  `request` text NOT NULL COMMENT '请求',
  `response` text NOT NULL COMMENT '响应',
  `parameter` varchar(255) NOT NULL DEFAULT '' COMMENT '参数',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_flow_task_id` (`flow_task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='流量分析风险表';

CREATE TABLE IF NOT EXISTS `flow_target` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `flow_task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属流量分析任务id',
  `target_url` varchar(255) NOT NULL DEFAULT '' COMMENT '测试目标地址',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '状态，1-禁用，2-正常',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_flow_task_id` (`flow_task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='流量分析目标表';

CREATE TABLE IF NOT EXISTS `flow_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `node_id` varchar(64) NOT NULL DEFAULT '' COMMENT '节点ID',
  `netword_card` varchar(36) NOT NULL DEFAULT '' COMMENT '节点网卡',
  `port` varchar(36) NOT NULL DEFAULT '' COMMENT '节点端口',
  `expire_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '分析结束时间',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '任务运行状态，1-已禁用，2-待执行 3-已开始 4-已结束',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '操作人id',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `other_config` varchar(1000) DEFAULT NULL COMMENT '其他配置',
  `vul_config` varchar(1000) DEFAULT NULL COMMENT '漏洞配置',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8 COMMENT='流量分析任务表';

CREATE TABLE IF NOT EXISTS `llm_models` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `model_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型名称',
  `platform` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台类型：openai、baidu、ali等',
  `api_url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'API地址',
  `api_key` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'API密钥',
  `model_id` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模型ID',
  `is_default` tinyint(4) DEFAULT '2' COMMENT '是否默认：1-是，2-否',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态：1-启用，2-禁用',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_platform` (`platform`),
  KEY `idx_status` (`status`),
  KEY `idx_is_default` (`is_default`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='大语言模型表';

CREATE TABLE IF NOT EXISTS `llm_scenarios` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL COMMENT '场景名称',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `icon` varchar(100) DEFAULT NULL COMMENT '图标',
  `is_enabled` tinyint(4) NOT NULL DEFAULT '1' COMMENT '是否启用：1-启用，2-禁用',
  `llm_model_id` int(11) NOT NULL DEFAULT '0' COMMENT '关联的大模型ID',
  `prompt` text COMMENT '场景提示词',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COMMENT='AI应用场景表';

CREATE TABLE IF NOT EXISTS `log_audit` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `log_type` int(11) NOT NULL DEFAULT '0' COMMENT '日志类型，1-登录日志，2-操作日志，3-告警日志',
  `content` text NOT NULL COMMENT '日志内容',
  `username` varchar(32) NOT NULL COMMENT '操作用户',
  `ip` varchar(32) NOT NULL DEFAULT '0' COMMENT '登录ip',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=782522 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='日志审计';

CREATE TABLE IF NOT EXISTS `log_backup` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(100) NOT NULL COMMENT '备份日志名称',
  `path` varchar(100) NOT NULL COMMENT '备份日志路径',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='备份日志';

CREATE TABLE IF NOT EXISTS `logic_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) DEFAULT NULL COMMENT '任务id',
  `target_id` int(11) DEFAULT NULL COMMENT '目标id',
  `target_url` varchar(255) DEFAULT NULL COMMENT '目标地址',
  `status` int(11) DEFAULT NULL COMMENT '状态',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `start_time` datetime DEFAULT NULL COMMENT '开始时间',
  `end_time` datetime DEFAULT NULL COMMENT '结束时间',
  `is_alive` tinyint(4) DEFAULT NULL COMMENT '是否存活',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `logic_log_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) DEFAULT NULL COMMENT '任务id',
  `target_id` int(11) DEFAULT NULL COMMENT '目标id',
  `log_id` int(11) DEFAULT NULL COMMENT '日志id',
  `target_url` varchar(255) DEFAULT NULL COMMENT '测试目标',
  `pocname` varchar(255) DEFAULT NULL COMMENT '漏洞标识',
  `result` varchar(255) DEFAULT NULL COMMENT '日志结果',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=216 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `logic_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '目标id',
  `node_id` varchar(50) NOT NULL DEFAULT '' COMMENT '节点id',
  `father_id` varchar(1024) NOT NULL DEFAULT '' COMMENT '父节点id',
  `pocname` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞标识',
  `vul_name` varchar(300) NOT NULL COMMENT '漏洞名称',
  `risk` tinyint(4) unsigned NOT NULL COMMENT '漏洞风险等级',
  `result` longtext NOT NULL COMMENT '漏洞结果',
  `request` longtext NOT NULL COMMENT '请求报文',
  `response` longtext NOT NULL COMMENT '响应报文',
  `checked` tinyint(4) unsigned NOT NULL COMMENT '检测结果 0未检测出来 1检测出来',
  `create_time` datetime(6) NOT NULL DEFAULT '1970-01-01 08:00:01.000000' COMMENT '创建时间',
  `decision_lib_id` int(11) NOT NULL DEFAULT '0' COMMENT '决策引擎漏洞库Id',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_pocname` (`pocname`) USING BTREE,
  KEY `target_id` (`target_id`) USING BTREE,
  KEY `node_id` (`node_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='逻辑漏洞结果表';

CREATE TABLE IF NOT EXISTS `logic_target` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `target_url` varchar(255) DEFAULT NULL COMMENT '目标地址',
  `task_id` int(11) DEFAULT NULL COMMENT '任务id',
  `status` int(11) DEFAULT NULL COMMENT '目标状态',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `type` varchar(255) DEFAULT NULL COMMENT '逻辑漏洞测试类型',
  `configJson` varchar(255) DEFAULT NULL COMMENT '扫描配置信息',
  `risk` int(11) DEFAULT NULL COMMENT '风险等级',
  `is_alive` int(11) DEFAULT NULL COMMENT '是否存活',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=35 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `logic_task` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(100) DEFAULT NULL COMMENT '名称',
  `target_url` varchar(500) DEFAULT NULL COMMENT '目标',
  `status` int(11) DEFAULT NULL COMMENT '状态',
  `target_num` int(11) DEFAULT NULL COMMENT '目标数量',
  `user_id` int(11) DEFAULT NULL COMMENT '用户id',
  `scan_config` text COMMENT '扫描配置',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `type` int(11) DEFAULT NULL COMMENT '扫描类型',
  `risk` int(11) DEFAULT NULL COMMENT '风险等级',
  `high_num` int(11) DEFAULT NULL,
  `middle_num` int(11) DEFAULT NULL,
  `low_num` int(11) DEFAULT NULL,
  `safe_num` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `logic_vul` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) DEFAULT NULL COMMENT '任务id',
  `target_id` int(11) DEFAULT NULL COMMENT '目标id',
  `pocname` varchar(255) DEFAULT NULL COMMENT '漏洞标识',
  `name` varchar(255) DEFAULT NULL COMMENT '漏洞名称',
  `class` varchar(255) DEFAULT NULL COMMENT '漏洞分类',
  `type` int(11) DEFAULT NULL COMMENT '漏洞类型',
  `risk` int(11) DEFAULT NULL COMMENT '风险等级',
  `location` varchar(255) DEFAULT NULL COMMENT '漏洞位置',
  `description` text COMMENT '漏洞描述',
  `fix_suggest` text COMMENT '修复建议',
  `vul_param` varchar(255) DEFAULT NULL COMMENT '漏洞参数',
  `vul_result` varchar(2000) DEFAULT NULL COMMENT '漏洞结果',
  `ver_msg` text COMMENT '验证报文',
  `decision_vul_id` varchar(50) DEFAULT NULL COMMENT '漏洞id',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `target_url` varchar(200) DEFAULT NULL COMMENT '测试地址',
  `vul_id` varchar(200) NOT NULL DEFAULT '' COMMENT '漏洞id',
  `status` int(11) DEFAULT NULL COMMENT '漏洞状态',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `map_set` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `estate` varchar(64) NOT NULL DEFAULT 'valid' COMMENT '数据状态valid/deleted',
  `obj_key` varchar(255) NOT NULL DEFAULT '' COMMENT '字典键key',
  `obj_value` longtext NOT NULL COMMENT '字典值value',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '备注说明',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_obj_key` (`obj_key`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=74 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='枚举map表';

CREATE TABLE IF NOT EXISTS `remote_session` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `estate` varchar(64) NOT NULL DEFAULT 'valid' COMMENT '数据状态valid/deleted',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '目标ID',
  `session_id` varchar(255) NOT NULL DEFAULT '0' COMMENT '会话编号',
  `target_url` varchar(255) NOT NULL DEFAULT '' COMMENT '目标地址',
  `route` varchar(200) NOT NULL DEFAULT '' COMMENT '路由',
  `remote_control` varchar(200) NOT NULL DEFAULT '' COMMENT '远程控制',
  `detail` varchar(1000) NOT NULL DEFAULT '' COMMENT '详情',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态',
  `downloaded_files` text NOT NULL COMMENT '文件下载',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '操作人id',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `session_key` varchar(255) NOT NULL COMMENT '会话key',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_task_id` (`task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=200 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='远程会话表';

CREATE TABLE IF NOT EXISTS `report_record` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(200) NOT NULL DEFAULT '' COMMENT '报告标题',
  `type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '报告类型，1-统计报告，2-目标报告',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '报告状态，1-生成中，2-已生成',
  `config_json` longtext NOT NULL COMMENT '报告配置参数',
  `format` tinyint(4) NOT NULL DEFAULT '1' COMMENT '报告格式，1-html，2-word,3-pdf,4-excel,5-csv',
  `content` longtext NOT NULL COMMENT '报告内容',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '提交者id',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=238 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='报告记录表';

CREATE TABLE IF NOT EXISTS `report_verify_port` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `target` varchar(255) DEFAULT NULL COMMENT '目标',
  `port` varchar(20) DEFAULT NULL COMMENT '端口',
  `scheme` varchar(20) DEFAULT NULL COMMENT '协议',
  `service` varchar(50) DEFAULT NULL COMMENT '服务',
  `component` varchar(255) DEFAULT NULL COMMENT '组件',
  `task_id` int(11) DEFAULT NULL COMMENT '任务id',
  `target_id` int(11) DEFAULT NULL COMMENT '目标id',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `report_verify_target` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) DEFAULT NULL COMMENT '任务id',
  `target` varchar(255) DEFAULT NULL COMMENT '目标',
  `os` varchar(255) DEFAULT NULL COMMENT '操作系统',
  `risk` int(11) DEFAULT NULL COMMENT '风险等级',
  `exp` int(11) DEFAULT NULL COMMENT '利用成功数',
  `verify` int(11) DEFAULT NULL COMMENT '验证成功数',
  `failed` int(11) DEFAULT NULL COMMENT '验证失败数',
  `unverify` int(11) DEFAULT NULL COMMENT '未能验证数',
  `status` int(11) DEFAULT NULL COMMENT '运行状态',
  `analysis_data` text COMMENT '分析数据',
  `is_alive` int(11) DEFAULT NULL COMMENT '是否存活',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=116 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `report_verify_task` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL COMMENT '任务名称',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `execute_type` int(11) DEFAULT NULL COMMENT '执行方式',
  `producer` int(50) DEFAULT NULL COMMENT '厂商',
  `user` int(50) DEFAULT NULL COMMENT '提交者',
  `status` int(11) DEFAULT NULL COMMENT '状态',
  `risk` int(11) DEFAULT NULL COMMENT '风险等级',
  `overview` text COMMENT '统计信息',
  `is_stats` int(11) DEFAULT NULL COMMENT '是否统计',
  `fileinfo` varchar(500) DEFAULT NULL COMMENT '文件信息',
  `exp` int(255) DEFAULT NULL COMMENT '利用成功',
  `verify` int(255) DEFAULT NULL COMMENT '验证成功',
  `failed` int(255) DEFAULT NULL COMMENT '验证失败',
  `unverify` int(255) DEFAULT NULL COMMENT '未能验证',
  `execute_time` datetime DEFAULT NULL COMMENT '执行时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `report_verify_vul` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) DEFAULT NULL COMMENT '漏洞名称',
  `risk` varchar(255) DEFAULT NULL COMMENT '风险等级',
  `status` varchar(255) DEFAULT NULL COMMENT '漏洞状态',
  `location` varchar(255) DEFAULT NULL COMMENT '漏洞位置',
  `task_id` int(11) DEFAULT NULL COMMENT '任务id',
  `target_id` int(11) DEFAULT NULL COMMENT '目标id',
  `cve` varchar(50) DEFAULT NULL COMMENT 'cve编号',
  `cnvd` varchar(50) DEFAULT NULL COMMENT 'cnvd编号',
  `cnnvd` varchar(50) DEFAULT NULL COMMENT 'cnnvd编号',
  `desc` varchar(2000) DEFAULT NULL COMMENT '漏洞描述',
  `fix` varchar(2000) DEFAULT NULL COMMENT '修复建议',
  `cvss` varchar(50) DEFAULT NULL COMMENT 'cvss评分',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4044 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `scanner_node` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(36) NOT NULL DEFAULT '' COMMENT '节点名称',
  `ip` varchar(36) NOT NULL DEFAULT '' COMMENT '节点IP',
  `port` varchar(36) NOT NULL DEFAULT '' COMMENT '节点端口',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态 0不在线 1在线',
  `is_disable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '禁用状态 0启用 1禁用',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='分布式扫描 节点表';

CREATE TABLE IF NOT EXISTS `system_config_backup` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(100) NOT NULL COMMENT '系统配置备份名称',
  `path` varchar(100) NOT NULL COMMENT '系统配置备份路径',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='系统配置备份';

CREATE TABLE IF NOT EXISTS `system_message` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `content` varchar(300) NOT NULL COMMENT '系统配置备份名称',
  `type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '消息类型 1-通知 2-告警 3-异常',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '操作人id',
  `status` tinyint(4) NOT NULL DEFAULT '2' COMMENT '消息读取状态 1-已读 2-未读',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=337735 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='系统消息表';

CREATE TABLE IF NOT EXISTS `task_configuration` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `obj_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务场景id',
  `config_key` varchar(30) NOT NULL DEFAULT '' COMMENT '配置项的键值名',
  `config_json` longtext NOT NULL COMMENT '配置json,包含所有配置',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '提交者id',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3770 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务场景配置';

CREATE TABLE IF NOT EXISTS `task_evidence` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `estate` varchar(64) NOT NULL DEFAULT 'valid' COMMENT '数据状态valid/deleted',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '目标ID',
  `target_url` varchar(255) NOT NULL DEFAULT '' COMMENT '关联目标地址',
  `vul_name` varchar(200) NOT NULL DEFAULT '' COMMENT '关联漏洞名称',
  `risk_type` int(11) NOT NULL DEFAULT '0' COMMENT '风险类型 2-登录凭证 3-敏感数据 4-敏感文件 5-数据库',
  `risk_detail` text NOT NULL COMMENT '风险详情',
  `downloaded_files` varchar(512) NOT NULL DEFAULT '' COMMENT '已下载文件',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '操作人ID',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_task_id` (`task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7492 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务证据表';

CREATE TABLE IF NOT EXISTS `task_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属目标id',
  `target_url` varchar(255) NOT NULL DEFAULT '' COMMENT '测试目标地址',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `start_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '开始时间',
  `end_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '结束时间',
  `is_alive` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '是否存活',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_task_id` (`task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7080 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务日志表';

CREATE TABLE IF NOT EXISTS `task_log_info` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_log_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务日志id',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属目标id',
  `target_url` varchar(255) NOT NULL DEFAULT '' COMMENT '测试目标地址',
  `pocname` varchar(255) NOT NULL DEFAULT '' COMMENT 'pocname',
  `result` longtext NOT NULL COMMENT '日志内容',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_task_log_id` (`task_log_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=420497 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务日志详情表';

CREATE TABLE IF NOT EXISTS `task_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '目标id',
  `node_id` varchar(50) NOT NULL DEFAULT '' COMMENT '节点id',
  `father_id` varchar(1024) NOT NULL DEFAULT '' COMMENT '父节点id',
  `pocname` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞标识',
  `vul_name` varchar(300) NOT NULL COMMENT '漏洞名称',
  `risk` tinyint(4) unsigned NOT NULL COMMENT '漏洞风险等级',
  `result` longtext NOT NULL COMMENT '漏洞结果',
  `request` longtext NOT NULL COMMENT '请求报文',
  `response` longtext NOT NULL COMMENT '响应报文',
  `checked` tinyint(4) unsigned NOT NULL COMMENT '检测结果 0未检测出来 1检测出来',
  `create_time` datetime(6) NOT NULL DEFAULT '1970-01-01 08:00:01.000000' COMMENT '创建时间',
  `decision_lib_id` int(11) NOT NULL DEFAULT '0' COMMENT '决策引擎漏洞库Id',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_pocname` (`pocname`) USING BTREE,
  KEY `target_id` (`target_id`) USING BTREE,
  KEY `node_id` (`node_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=93498 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务检测结果表-场景(渗透路径图)';

CREATE TABLE IF NOT EXISTS `task_target` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_url` varchar(255) NOT NULL DEFAULT '' COMMENT '测试目标地址',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '目标状态',
  `op_sys` varchar(200) NOT NULL DEFAULT '' COMMENT '操作系统',
  `risk_level` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '目标风险等级',
  `is_alive` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否存活',
  `target_type` int(11) NOT NULL DEFAULT '0' COMMENT '目标类型',
  `task_template_id` int(11) DEFAULT '0' COMMENT '所选择的任务场景id',
  `task_template_json` longtext NOT NULL COMMENT '任务场景参数',
  `is_remote_session` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否有远程会话:1否，2是',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '提交者id',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `end_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `weight` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '优先权重，1-低，2-中，3-高',
  `use_score` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '利用评分*100',
  `is_score` tinyint(4) unsigned NOT NULL DEFAULT '4' COMMENT '利用评分状态，1-未开始,2-进行中,3-完成,4-无需评分',
  `extend_field` varchar(512) NOT NULL DEFAULT '' COMMENT '其他数据',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_task_id` (`task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=60593 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='目标表';

CREATE TABLE IF NOT EXISTS `task_target_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `obj_type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '数据类型',
  `sub_obj_type` varchar(50) NOT NULL DEFAULT '' COMMENT '数据子类型',
  `obj_id` varchar(50) NOT NULL DEFAULT '' COMMENT '数据对象id',
  `sub_obj_id` varchar(50) NOT NULL DEFAULT '' COMMENT '数据子对象id',
  `identify` varchar(255) NOT NULL DEFAULT '' COMMENT '数据标识符',
  `field1` varchar(30) NOT NULL DEFAULT '' COMMENT '筛选字段1',
  `field2` varchar(30) NOT NULL DEFAULT '' COMMENT '筛选字段2',
  `field3` varchar(30) NOT NULL DEFAULT '' COMMENT '筛选字段3',
  `json_result` text NOT NULL COMMENT '各数据类型的json格式结果',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_obj_id` (`obj_id`) USING BTREE,
  KEY `idx_field1` (`field1`) USING BTREE,
  KEY `idx_field2` (`field2`) USING BTREE,
  KEY `idx_field3` (`field3`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='通用记录数据表';

CREATE TABLE IF NOT EXISTS `task_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_name` text NOT NULL COMMENT '任务名称',
  `risk_level` int(11) NOT NULL DEFAULT '0' COMMENT '该任务风险等级',
  `status` int(11) DEFAULT '0' COMMENT '任务运行状态',
  `is_stats` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '是否需要重新统计,1否,2是',
  `task_type` int(11) NOT NULL DEFAULT '0' COMMENT '任务类型',
  `execute_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '执行方法:1-即时执行,2-定时执行,3-周期执行,4-监控执行',
  `task_template_id` int(11) DEFAULT '0' COMMENT '所选择的任务场景id',
  `target_num` int(11) NOT NULL DEFAULT '0' COMMENT '该任务下目标数量',
  `hige_num` int(11) DEFAULT '0' COMMENT '高危目标个数',
  `middle_num` int(11) DEFAULT '0' COMMENT '中危目标个数',
  `low_num` int(11) DEFAULT '0' COMMENT '低危目标个数',
  `safe_num` int(11) DEFAULT '0' COMMENT '安全目标个数',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '提交者id',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `weight` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '优先权重，1-低，2-中，3-高',
  `pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '任务父级任务ID，用于任务验证测试',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2832 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务基础表';

CREATE TABLE IF NOT EXISTS `task_task_info` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `task_type` int(11) NOT NULL DEFAULT '0' COMMENT '任务类型',
  `check_target` longtext,
  `execute_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '执行方法 1即时执行、2定时执行 、3周期执行、4监控执行',
  `execute_last_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '任务上次执行时间',
  `execute_next_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '任务下次执行时间',
  `execute_json` text NOT NULL COMMENT '执行方式具体参数',
  `task_template_id` int(11) NOT NULL DEFAULT '0' COMMENT '所选择的任务场景id',
  `task_template_json` longtext NOT NULL COMMENT '任务场景参数',
  `overview` longtext NOT NULL COMMENT '概览统计',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '提交者id',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `status` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '任务运行状态',
  `weight` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '优先权重，1-低，2-中，3-高',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_task_id` (`task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1562 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务详情表';

CREATE TABLE IF NOT EXISTS `task_task_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `obj_type` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '数据类型：1-信息收集',
  `sub_obj_type` varchar(50) NOT NULL DEFAULT '' COMMENT '数据子类型',
  `obj_id` varchar(50) NOT NULL DEFAULT '' COMMENT '数据对象id',
  `sub_obj_id` varchar(50) NOT NULL DEFAULT '' COMMENT '数据子对象id',
  `identify` varchar(50) NOT NULL DEFAULT '' COMMENT '数据标识符,用来修改时方便定位数据的',
  `field1` varchar(200) NOT NULL DEFAULT '' COMMENT '筛选字段1',
  `field2` varchar(200) NOT NULL DEFAULT '' COMMENT '筛选字段2',
  `field3` varchar(200) NOT NULL DEFAULT '' COMMENT '筛选字段3',
  `field4` varchar(200) NOT NULL DEFAULT '' COMMENT '筛选字段4',
  `json_result` text NOT NULL COMMENT '各数据类型的json格式结果',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_obj_id` (`obj_id`) USING BTREE,
  KEY `idx_sub_obj_id` (`sub_obj_id`) USING BTREE,
  KEY `idx_field1` (`field1`) USING BTREE,
  KEY `idx_field2` (`field2`) USING BTREE,
  KEY `idx_field3` (`field3`) USING BTREE,
  KEY `idx_obj_sub_obj` (`obj_type`,`sub_obj_type`,`sub_obj_id`)
) ENGINE=InnoDB AUTO_INCREMENT=219740 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='通用结果数据表';

CREATE TABLE IF NOT EXISTS `task_template` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `estate` varchar(64) NOT NULL DEFAULT 'valid' COMMENT '数据状态valid/deleted',
  `template_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务场景名字',
  `describe` varchar(255) NOT NULL DEFAULT '' COMMENT '任务描述',
  `is_default` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '1:默认，2：不是默认',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '提交者id',
  `source` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '模板来源，1用户建的场景，2软件自带场景',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=291 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='任务场景模版表';

CREATE TABLE IF NOT EXISTS `task_vul` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `data_type` tinyint(4) NOT NULL DEFAULT '1' COMMENT '漏洞类型，1-漏洞测试,2-待测试漏洞',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属目标id',
  `target_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '测试目标地址',
  `pocname` varchar(1024) NOT NULL DEFAULT '' COMMENT '漏洞标识',
  `name` varchar(200) NOT NULL DEFAULT '' COMMENT '漏洞名称',
  `class` int(11) NOT NULL DEFAULT '0' COMMENT '漏洞分类',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '漏洞类型',
  `risk` int(11) NOT NULL DEFAULT '4' COMMENT '风险等级',
  `location` varchar(1000) NOT NULL DEFAULT '' COMMENT '漏洞位置或地址',
  `status` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '漏洞状态,1-未验证，2-验证成功，3-利用成功，4-已修复，5-流量验证,6-未知,7-存在,8-不存在',
  `test_status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '测试状态，1-未测试，2-测试中，3-已测试',
  `exploit_impact` varchar(10) NOT NULL DEFAULT '' COMMENT '利用影响',
  `vul_id` varchar(200) NOT NULL DEFAULT '' COMMENT '漏洞id',
  `description` text NOT NULL COMMENT '漏洞描述',
  `fix_suggest` varchar(2156) NOT NULL DEFAULT '' COMMENT '修复建议',
  `published_time` varchar(32) NOT NULL DEFAULT '' COMMENT '披露时间',
  `affect_range` varchar(255) NOT NULL DEFAULT '' COMMENT '影响范围',
  `target_result_id` int(11) NOT NULL DEFAULT '0' COMMENT '检测结果id',
  `vul_number` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞编号',
  `vul_address` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞地址',
  `ref_url` varchar(255) NOT NULL DEFAULT '' COMMENT '参考链接',
  `cvss` varchar(255) NOT NULL DEFAULT '' COMMENT 'cvss评分',
  `vul_param` text NOT NULL COMMENT '漏洞请求参数',
  `vul_result` longtext NOT NULL COMMENT '漏洞结果',
  `ver_msg` longtext NOT NULL COMMENT '验证报文',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `decision_vul_id` varchar(50) NOT NULL DEFAULT '' COMMENT '决策引擎的漏洞唯一ID',
  `snapshot` longtext NOT NULL COMMENT '截图',
  `patch_url` varchar(1000) DEFAULT NULL COMMENT '补丁链接',
  `ai_analysis_result` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'AI研判结果',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_task_id` (`task_id`) USING BTREE,
  KEY `idx_target_id` (`target_id`) USING BTREE,
  KEY `idx_group_maxid_fixed` (`pocname`(100),`name`(100),`type`,`risk`,`id`),
  KEY `idx_target_data_risk` (`target_id`,`data_type`,`risk`)
) ENGINE=InnoDB AUTO_INCREMENT=50147 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='测试漏洞表';

CREATE TABLE IF NOT EXISTS `tool_file` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '文件名称',
  `file_type` varchar(50) NOT NULL DEFAULT '' COMMENT '文件类型',
  `file_path` varchar(255) NOT NULL DEFAULT '' COMMENT '文件路径',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_file_path` (`file_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='工具库文件下载';

CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `estate` varchar(64) NOT NULL DEFAULT 'valid' COMMENT '数据状态valid/deleted',
  `type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '账号类型,1-普通用户，2-管理员，3-审核员，4-超级管理员',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '账号状态，1-账号正常，2-登录失败次数过多禁用，3-账号有效期超期禁用，4-密码长期未修改禁用，5-手动禁用，6-账号长时间未登录禁用',
  `is_alive` tinyint(4) NOT NULL DEFAULT '1' COMMENT '在线状态，1-在线，2-离线',
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '账号名称',
  `real_name` varchar(50) NOT NULL DEFAULT '' COMMENT '真实姓名',
  `email` varchar(50) NOT NULL DEFAULT '' COMMENT '邮箱',
  `department` varchar(50) NOT NULL DEFAULT '' COMMENT '所属部门',
  `password` varchar(128) NOT NULL DEFAULT '' COMMENT '密码',
  `remark` varchar(128) NOT NULL DEFAULT '' COMMENT '描述',
  `operator_id` int(11) NOT NULL DEFAULT '0' COMMENT '操作人id',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `last_login_time` datetime NOT NULL COMMENT '上次登录时间',
  `last_operate_time` datetime NOT NULL COMMENT '用户最后操作时间',
  `account_expire_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '账号过期时间',
  `password_expire_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '密码过期时间',
  `token` varchar(64) NOT NULL DEFAULT '' COMMENT 'API密钥',
  `token_create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '密钥生成时间',
  `token_expire_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '密钥过期时间',
  `password_change_time` datetime DEFAULT NULL COMMENT '密码更改时间',
  `password_already_used` varchar(2000) DEFAULT '' COMMENT '已经使用过的密码',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_email` (`email`) USING BTREE,
  KEY `idx_username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=85 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `user_gophish_key` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `gophish_key` varchar(255) NOT NULL DEFAULT '' COMMENT 'gophish_api_key',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_id` (`user_id`) USING BTREE,
  UNIQUE KEY `gophish_key` (`gophish_key`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COMMENT='钓鱼管理的 apikey 管理';

CREATE TABLE IF NOT EXISTS `user_groups` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '组名称',
  `pid` int(11) NOT NULL DEFAULT '0' COMMENT '父级ID',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '组状态',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '创建者id',
  `is_range_open` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否开启测试范围,1否，2是',
  `range` text NOT NULL COMMENT '测试范围',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户组表';

CREATE TABLE IF NOT EXISTS `user_groups_users` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_groups_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户组ID',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '用户ID',
  `submit_user_id` int(11) NOT NULL DEFAULT '0' COMMENT '提交人用户ID',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_user_id` (`user_id`) USING BTREE,
  KEY `idx_user_groups_id` (`user_groups_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='用户组与用户关联表';

CREATE TABLE IF NOT EXISTS `vul_check` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `vul_id` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '0',
  `check_result` tinyint(4) NOT NULL DEFAULT '0',
  `check_detail` text,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `vul_libraries` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `data_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '数据类型，1:真实数据 2:假数据',
  `vul_id` varchar(50) DEFAULT '' COMMENT '四维自定义漏洞id',
  `name` varchar(1000) DEFAULT '' COMMENT '漏洞名称',
  `cve` varchar(50) DEFAULT '' COMMENT 'cve编号',
  `risk` int(11) DEFAULT '0' COMMENT '漏洞风险',
  `type` int(11) DEFAULT '0' COMMENT '漏洞类型',
  `class` int(11) DEFAULT '0' COMMENT '漏洞分类',
  `published_time` varchar(32) DEFAULT '' COMMENT '公开时间',
  `description` text COMMENT '漏洞描述',
  `affect_range` text COMMENT '影响范围',
  `exploit_impact` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '利用影响',
  `fix_suggest` text COMMENT '修复建议',
  `cnvd` varchar(255) DEFAULT '' COMMENT 'cnvd编号',
  `cnnvd` varchar(255) DEFAULT '' COMMENT 'cnnvd编号',
  `component` varchar(1000) NOT NULL DEFAULT '' COMMENT '受影响的组建名称',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '脚本状态:0删除 1=已禁用，2=正常，其他状态请查看Enum文件',
  `create_time` datetime DEFAULT '1970-01-01 08:00:01' COMMENT '添加时间',
  `update_time` datetime DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
  `status_msg` varchar(128) NOT NULL DEFAULT '' COMMENT '脚本状态描述',
  `priority` varchar(20) NOT NULL DEFAULT '' COMMENT '修复优先级',
  `operating_system` int(11) NOT NULL DEFAULT '0' COMMENT '1:windows, 2:linux, 3:unix, 4:其他',
  `pocname` varchar(255) DEFAULT NULL COMMENT '漏洞标识',
  `verify_type` int(11) DEFAULT NULL COMMENT '漏洞验证类型  1 原理验证 2 版本匹配',
  `patch_url` varchar(1000) DEFAULT NULL COMMENT '补丁链接',
  `bugtraq` varchar(255) DEFAULT '' COMMENT 'Bugtraq编号',
  `cncve` varchar(255) DEFAULT '' COMMENT 'CNCVE编号',
  `cvss_version` varchar(20) DEFAULT '' COMMENT 'CVSS版本',
  `cvss_score` varchar(20) DEFAULT '' COMMENT 'cvss评分',
  `poc_or_exp` varchar(20) DEFAULT NULL COMMENT 'poc或者exp',
  `script_type` varchar(50) DEFAULT NULL COMMENT '脚本类型',
  `ms` varchar(50) DEFAULT NULL COMMENT 'ms漏洞编号',
  `vul_source` tinyint(4) NOT NULL DEFAULT '2' COMMENT '漏洞来源（是否国产化）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `uni_pocname` (`pocname`) USING BTREE COMMENT 'pocname唯一索引',
  KEY `vul_id_index` (`vul_id`) USING BTREE COMMENT '漏洞id索引',
  KEY `idx_name` (`name`(50)) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=294680 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='漏洞表';

CREATE TABLE IF NOT EXISTS `wifi_ap_info` (
  `source_mac` varchar(32) NOT NULL,
  `attack_type` int(11) NOT NULL,
  `manuf` varchar(255) DEFAULT NULL,
  `first_time` bigint(20) DEFAULT NULL,
  `last_time` bigint(20) DEFAULT NULL,
  `duration` int(11) DEFAULT NULL,
  `channel` int(11) NOT NULL,
  `maxseenrate` int(11) DEFAULT NULL,
  `carrierset` int(11) DEFAULT NULL,
  `encodingset` varchar(256) DEFAULT NULL,
  `llc_packets` bigint(20) NOT NULL,
  `data_packets` bigint(20) DEFAULT NULL,
  `datasize` bigint(20) DEFAULT NULL,
  `last_signal_rssi` int(11) DEFAULT NULL,
  `min_signal_rssi` int(11) DEFAULT NULL,
  `max_signal_rssi` int(11) DEFAULT NULL,
  `ssid_cloaked` int(11) DEFAULT NULL COMMENT '0：广播，1：隐藏',
  `ssid` varchar(255) DEFAULT NULL,
  `ssid_beacon_info` varchar(256) DEFAULT NULL,
  `ssid_first_time` bigint(20) NOT NULL,
  `ssid_last_time` bigint(20) DEFAULT NULL,
  `ssid_maxrate` double DEFAULT NULL,
  `ssid_beaconrate` int(11) DEFAULT NULL,
  `ssid_country` varchar(64) DEFAULT NULL,
  `ssid_startchan` int(11) DEFAULT NULL,
  `ssid_endchan` int(11) DEFAULT NULL,
  `ssid_txpower` int(11) DEFAULT NULL,
  `ssid_cryptset` int(11) DEFAULT NULL,
  `location_flag` int(11) DEFAULT NULL,
  `lat` double DEFAULT NULL,
  `lon` double DEFAULT NULL,
  `alt` double DEFAULT NULL,
  `drone_num` int(11) DEFAULT NULL,
  `nearest_drone` varchar(64) DEFAULT NULL,
  `nearest_drone_rssi` int(11) DEFAULT NULL,
  `nearest_drone_distance` double DEFAULT NULL,
  `freq_counter` varchar(1024) DEFAULT NULL,
  `list` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT '0' COMMENT '0：无，1：强制下线',
  `connect_cli_count` int(11) DEFAULT NULL,
  `area_id` int(11) DEFAULT NULL,
  `label_id` int(11) DEFAULT NULL,
  `is_mydevice` int(11) DEFAULT NULL COMMENT '0:非自有设备，1：自有设备',
  `network_type` int(11) DEFAULT NULL COMMENT '0：内网，1：外网',
  `tag` int(11) DEFAULT '0' COMMENT '1:标记为非私搭乱建',
  `device_type` int(11) DEFAULT '9' COMMENT '0:AP，  1:终端',
  `operate_user` varchar(25) DEFAULT NULL,
  `add_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`source_mac`) USING BTREE,
  KEY `dd` (`area_id`) USING BTREE,
  KEY `label_id` (`label_id`) USING BTREE,
  KEY `last_time` (`last_time`,`source_mac`,`status`,`ssid`,`channel`,`manuf`,`lat`,`lon`,`area_id`,`first_time`,`connect_cli_count`,`last_signal_rssi`,`freq_counter`(255),`attack_type`,`label_id`,`is_mydevice`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `wifi_cli_info` (
  `source_mac` varchar(32) NOT NULL,
  `attack_type` int(11) NOT NULL,
  `manuf` varchar(255) DEFAULT NULL,
  `first_time` bigint(20) DEFAULT NULL,
  `last_time` bigint(20) DEFAULT NULL,
  `duration` int(11) DEFAULT NULL,
  `channel` int(11) DEFAULT NULL,
  `maxseenrate` int(11) DEFAULT NULL,
  `carrierset` int(11) DEFAULT NULL,
  `encodingset` int(11) DEFAULT NULL,
  `llc_packets` bigint(20) DEFAULT NULL,
  `data_packets` bigint(20) DEFAULT NULL,
  `datasize` bigint(20) DEFAULT NULL,
  `last_signal_rssi` int(11) DEFAULT NULL,
  `min_signal_rssi` int(11) DEFAULT NULL,
  `max_signal_rssi` int(11) DEFAULT NULL,
  `bssid_mac` varchar(32) DEFAULT NULL,
  `bssid_first_time` bigint(20) DEFAULT NULL,
  `bssid_last_time` bigint(20) DEFAULT NULL,
  `ssid_flag` int(11) DEFAULT NULL,
  `ssid_cloaked` int(11) DEFAULT NULL COMMENT '0：广播，1：隐藏',
  `ssid` varchar(255) DEFAULT NULL,
  `ssid_beacon_info` varchar(256) DEFAULT NULL,
  `ssid_first_time` bigint(20) DEFAULT NULL,
  `ssid_last_time` bigint(20) DEFAULT NULL,
  `ssid_maxrate` double DEFAULT NULL,
  `ssid_beaconrate` int(11) DEFAULT NULL,
  `ssid_country` varchar(64) DEFAULT NULL,
  `ssid_startchan` int(11) DEFAULT NULL,
  `ssid_endchan` int(11) DEFAULT NULL,
  `ssid_txpower` int(11) DEFAULT NULL,
  `ssid_cryptset` int(11) DEFAULT NULL,
  `location_flag` int(11) DEFAULT NULL,
  `lat` double DEFAULT NULL,
  `lon` double DEFAULT NULL,
  `alt` double DEFAULT NULL,
  `drone_num` int(11) DEFAULT NULL,
  `nearest_drone` varchar(64) DEFAULT NULL,
  `nearest_drone_rssi` int(11) DEFAULT NULL,
  `nearest_drone_distance` double DEFAULT NULL,
  `freq_counter` varchar(1024) DEFAULT NULL,
  `list` int(11) DEFAULT NULL,
  `status` int(11) DEFAULT '0',
  `connect_ap_count` int(11) DEFAULT NULL,
  `area_id` int(11) DEFAULT NULL,
  `label_id` int(11) DEFAULT NULL,
  `is_mydevice` int(11) DEFAULT NULL COMMENT '1：自有设备',
  `network_type` int(11) DEFAULT '9' COMMENT '0:内网，1：外网',
  `device_type` int(11) DEFAULT NULL COMMENT '0:AP，  1:终端',
  `operate_user` varchar(25) DEFAULT NULL,
  `add_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`source_mac`) USING BTREE,
  KEY `cli_area` (`area_id`) USING BTREE,
  KEY `label_id` (`label_id`) USING BTREE,
  KEY `last_time` (`area_id`,`last_time`,`source_mac`,`attack_type`,`manuf`,`channel`,`ssid`,`last_signal_rssi`,`bssid_mac`,`first_time`,`status`,`label_id`,`is_mydevice`,`ssid_cloaked`,`bssid_first_time`,`network_type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `wifi_task` (
  `task_id` int(11) NOT NULL AUTO_INCREMENT,
  `mac` varchar(32) NOT NULL,
  `status` int(11) NOT NULL COMMENT '0：初始状态，等待开始检测\r\n1：开始进行报文收集\r\n2：正在进行密码爆破\r\n3：密码爆破失败\r\n4：密码爆破成功',
  `passwd_source` int(11) NOT NULL,
  `passwd_dict` varchar(255) DEFAULT NULL,
  `channel` int(11) NOT NULL,
  `encrypt` int(11) NOT NULL,
  `passwd` varchar(255) DEFAULT NULL,
  `start_time` bigint(20) NOT NULL,
  `end_time` bigint(20) DEFAULT NULL,
  `ssid` varchar(128) NOT NULL,
  `is_simulate` tinyint(4) DEFAULT '0' COMMENT '是否模拟：0：否，1：是',
  `simulate_duration` int(11) DEFAULT NULL COMMENT '模拟时长，单位秒',
  `model` varchar(128) DEFAULT NULL,
  `is_crack` tinyint(4) unsigned NOT NULL DEFAULT '0' COMMENT '是否爆破，0：否，1：是',
  `id` int(11) DEFAULT NULL,
  `is_embed` tinyint(4) DEFAULT '0' COMMENT '是否植入，0：否，1：是',
  `reason_code` tinyint(4) DEFAULT '0',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  `task_name` varchar(50) NOT NULL DEFAULT '' COMMENT '任务名称',
  `carrier` tinyint(4) NOT NULL DEFAULT '0' COMMENT '协议类型',
  PRIMARY KEY (`task_id`,`mac`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `wifi_task_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ap_mac` varchar(32) DEFAULT NULL,
  `task_id` int(11) DEFAULT NULL,
  `content` varchar(255) DEFAULT NULL,
  `generate_time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

-- ---------------------------------------------------------------------------
-- Security module result tables (defined in models/mysqls, pending migration)
-- ---------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS `baseline_check_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '目标id',
  `target_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '目标IP',
  `os_type` int(11) NOT NULL DEFAULT '0' COMMENT '操作系统类型',
  `scan_scene` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1=安全配置核查 2=主机漏洞检测(同引擎/规则库)',
  `rule_id` int(11) NOT NULL DEFAULT '0' COMMENT '规则ID',
  `rule_name` varchar(255) NOT NULL DEFAULT '' COMMENT '规则名称',
  `rule_category` int(11) NOT NULL DEFAULT '0' COMMENT '规则分类',
  `rule_risk` int(11) NOT NULL DEFAULT '0' COMMENT '规则风险',
  `check_result` int(11) NOT NULL DEFAULT '0' COMMENT '检查结果',
  `expected_value` varchar(512) NOT NULL DEFAULT '' COMMENT '期望值',
  `actual_value` varchar(512) NOT NULL DEFAULT '' COMMENT '实际值',
  `check_command` text COMMENT '检查命令',
  `fix_suggestion` text COMMENT '修复建议',
  `risk_description` text COMMENT '风险描述',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_target_id` (`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机基线核查结果';

-- 主机核查规则（导入后服务启动或 POST /smart/baseline/rules/reload 会加载；至少一条 enabled=1 时覆盖内置规则）
CREATE TABLE IF NOT EXISTS `host_baseline_rule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `rule_code` int(11) NOT NULL COMMENT '规则编号，与 baseline_check_result.rule_id 对应，全局唯一',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '规则名称',
  `description` varchar(512) NOT NULL DEFAULT '' COMMENT '规则描述',
  `category` int(11) NOT NULL DEFAULT '0' COMMENT '规则分类，见 enums.BaselineCategory*',
  `risk` int(11) NOT NULL DEFAULT '0' COMMENT '风险等级',
  `os_type` int(11) NOT NULL DEFAULT '0' COMMENT '适用系统类型 1=Linux 2=Windows 3=国产 4=嵌入式',
  `commands_json` text COMMENT '检查命令 JSON 数组，如 ["cmd1","cmd2"]',
  `expected_value` varchar(512) NOT NULL DEFAULT '' COMMENT '期望值',
  `match_type` varchar(32) NOT NULL DEFAULT 'contains' COMMENT '匹配方式 contains|exact|regex 等',
  `fix_suggestion` text COMMENT '修复建议',
  `risk_description` varchar(512) NOT NULL DEFAULT '' COMMENT '风险说明',
  `enabled` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1=启用 0=停用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rule_code` (`rule_code`),
  KEY `idx_enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='主机核查规则表';

-- 示例（可删）：Linux SSH Protocol 2
-- INSERT INTO `host_baseline_rule` (`rule_code`,`name`,`description`,`category`,`risk`,`os_type`,`commands_json`,`expected_value`,`match_type`,`fix_suggestion`,`risk_description`,`enabled`) VALUES
-- (100001,'示例-SSH协议2','示例规则',9,1,1,'["grep ''^Protocol'' /etc/ssh/sshd_config 2>/dev/null || echo missing"]','Protocol 2','contains','在 sshd_config 中设置 Protocol 2','',1);

-- 已部署库若缺列 baseline_check_result.scan_scene 请执行：
-- ALTER TABLE `baseline_check_result` ADD COLUMN `scan_scene` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1=安全配置核查 2=主机漏洞检测' AFTER `os_type`;

-- 病毒库规则表（用于 YARA 规则管理）
CREATE TABLE IF NOT EXISTS `malware_rule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '规则名称',
  `description` varchar(512) NOT NULL DEFAULT '' COMMENT '规则描述',
  `risk_level` int(11) NOT NULL DEFAULT '0' COMMENT '风险等级：1-高危 2-中危 3-低危',
  `rule_content` text COMMENT 'YARA 规则内容',
  `os_type` int(11) NOT NULL DEFAULT '0' COMMENT '适用系统类型：1-Linux 2-Windows 3-国产 4-嵌入式 0-通用',
  `category` varchar(255) NOT NULL DEFAULT '' COMMENT '规则分类，如：挖矿木马、Webshell、APT',
  `status` int(11) NOT NULL DEFAULT '1' COMMENT '1-启用 0-停用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_os_type` (`os_type`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='病毒库规则表';

CREATE TABLE IF NOT EXISTS `db_check_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '目标id',
  `target_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '目标IP',
  `db_type` int(11) NOT NULL DEFAULT '0' COMMENT '数据库类型',
  `db_name` varchar(128) NOT NULL DEFAULT '' COMMENT '数据库名',
  `check_category` int(11) NOT NULL DEFAULT '0' COMMENT '检查分类',
  `rule_id` int(11) NOT NULL DEFAULT '0' COMMENT '规则ID',
  `rule_name` varchar(255) NOT NULL DEFAULT '' COMMENT '规则名称',
  `check_result` int(11) NOT NULL DEFAULT '0' COMMENT '检查结果',
  `expected_value` text COMMENT '期望值',
  `actual_value` text COMMENT '实际值',
  `risk_level` int(11) NOT NULL DEFAULT '0' COMMENT '风险等级',
  `fix_suggestion` text COMMENT '修复建议',
  `risk_description` text COMMENT '风险描述',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_target_id` (`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据库配置核查结果';

CREATE TABLE IF NOT EXISTS `malware_check_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '目标id',
  `target_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '目标IP',
  `check_type` int(11) NOT NULL DEFAULT '0' COMMENT '检查类型',
  `risk_level` int(11) NOT NULL DEFAULT '0' COMMENT '风险等级',
  `match_rule` varchar(512) NOT NULL DEFAULT '' COMMENT '匹配规则',
  `file_path` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件路径',
  `process_info` text COMMENT '进程信息',
  `description` text COMMENT '描述',
  `fix_suggestion` text COMMENT '修复建议',
  `raw_output` longtext COMMENT '原始输出',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_target_id` (`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='恶意代码检测结果';

CREATE TABLE IF NOT EXISTS `sensitive_data_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) NOT NULL DEFAULT '0' COMMENT '目标id',
  `target_ip` varchar(50) NOT NULL DEFAULT '' COMMENT '目标IP',
  `db_type` int(11) NOT NULL DEFAULT '0' COMMENT '数据库类型',
  `db_name` varchar(128) NOT NULL DEFAULT '' COMMENT '数据库名',
  `table_name` varchar(128) NOT NULL DEFAULT '' COMMENT '表名',
  `column_name` varchar(128) NOT NULL DEFAULT '' COMMENT '字段名',
  `data_type` int(11) NOT NULL DEFAULT '0' COMMENT '数据类型',
  `data_level` int(11) NOT NULL DEFAULT '0' COMMENT '数据等级',
  `sample_data` varchar(512) NOT NULL DEFAULT '' COMMENT '样本数据',
  `match_rule` varchar(512) NOT NULL DEFAULT '' COMMENT '匹配规则',
  `match_type` int(11) NOT NULL DEFAULT '0' COMMENT '匹配类型',
  `total_rows` bigint(20) NOT NULL DEFAULT '0' COMMENT '总行数',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_target_id` (`target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='敏感数据发现结果';

-- 数据安全检测规则（导入后任务扫描会优先使用库中启用规则；无库规则时回退内置约 24 条）
CREATE TABLE IF NOT EXISTS `datasec_rule` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `rule_code` int(11) NOT NULL COMMENT '规则编号，与 db_check_result.rule_id 对应',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '规则名称',
  `description` varchar(512) NOT NULL DEFAULT '' COMMENT '规则描述',
  `category` int(11) NOT NULL DEFAULT '0' COMMENT '分类：身份认证/权限/SQL注入/敏感数据等',
  `risk` int(11) NOT NULL DEFAULT '0' COMMENT '风险等级 0-4',
  `db_type` int(11) NOT NULL DEFAULT '0' COMMENT '适用数据库 0=全部 1=MySQL 2=PG 3=Mongo 4=Redis 5=CouchDB',
  `queries_json` text COMMENT '检查 SQL/命令 JSON 数组',
  `expected_value` varchar(512) NOT NULL DEFAULT '' COMMENT '期望值说明',
  `match_type` varchar(32) NOT NULL DEFAULT 'contains' COMMENT 'contains|exact|not_contains|empty|always',
  `fix_suggestion` text COMMENT '修复建议',
  `risk_description` varchar(512) NOT NULL DEFAULT '' COMMENT '风险说明',
  `enabled` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1=启用 0=停用',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_rule_code` (`rule_code`),
  KEY `idx_db_type_enabled` (`db_type`, `enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据安全检测规则表';

CREATE TABLE IF NOT EXISTS `datasec_db_target` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int(11) NOT NULL DEFAULT '0' COMMENT '所属用户',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT '目标名称',
  `group_name` varchar(64) NOT NULL DEFAULT '' COMMENT '分组',
  `db_type` int(11) NOT NULL DEFAULT '1' COMMENT '1=MySQL 2=PG 3=Mongo 4=Redis 5=CouchDB',
  `db_host` varchar(255) NOT NULL DEFAULT '' COMMENT '地址',
  `db_port` int(11) NOT NULL DEFAULT '0' COMMENT '端口',
  `db_name` varchar(128) NOT NULL DEFAULT '' COMMENT '库名',
  `db_user` varchar(128) NOT NULL DEFAULT '' COMMENT '用户名',
  `db_password` varchar(512) NOT NULL DEFAULT '' COMMENT '密码(AES加密base64)',
  `remark` varchar(512) NOT NULL DEFAULT '' COMMENT '备注',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_db_type` (`user_id`, `db_type`),
  KEY `idx_user_group` (`user_id`, `group_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据安全-数据库扫描目标库';

