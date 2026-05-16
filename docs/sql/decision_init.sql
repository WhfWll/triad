-- Database: decision
-- All tables use CREATE TABLE IF NOT EXISTS so the script can be re-run without error 1050.
-- Note: existing tables are not altered; use migrations to change column definitions.
CREATE DATABASE IF NOT EXISTS `decision` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE `decision`;

CREATE TABLE IF NOT EXISTS `cvss` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
  `cve_id` varchar(255) NOT NULL DEFAULT '' COMMENT 'cveid',
  `cvss` int(11) unsigned NOT NULL DEFAULT '0' COMMENT 'cvss评分',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=391770 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='cvss评分表';

CREATE TABLE IF NOT EXISTS `dictionary` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `sources` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '来源，1系统，2手动',
  `service` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '适用范围',
  `name` varchar(50) NOT NULL DEFAULT '' COMMENT '字典名称',
  `types` tinyint(4) unsigned NOT NULL DEFAULT '1' COMMENT '字典类型：1-用户字典，2-密码字典，3-wifi，4-路径字典，5-子域名字典',
  `is_default` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '是否是默认，1-是，2-否',
  `content` longtext NOT NULL COMMENT '字典内容,换行符隔开',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=90 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='字典表';

CREATE TABLE IF NOT EXISTS `finger` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `app_version` varchar(100) NOT NULL DEFAULT '' COMMENT '应用版本',
  `cn_name` varchar(200) NOT NULL DEFAULT '' COMMENT '厂商名称-指纹名称',
  `app_name` varchar(200) NOT NULL DEFAULT '' COMMENT '应用名称-CPE',
  `flag` varchar(2000) NOT NULL DEFAULT '' COMMENT '匹配内容 指纹规则',
  `is_dev` tinyint(6) unsigned NOT NULL DEFAULT '0' COMMENT '是否设备指纹，1是，2否 [0软件 1 硬件]',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `app_class` int(11) NOT NULL DEFAULT '0' COMMENT '指纹分类',
  `desc` varchar(500) DEFAULT NULL COMMENT '指纹描述',
  `level` varchar(10) NOT NULL DEFAULT '0' COMMENT '服务分层	硬件层、系统层、服务层、支撑层、应用层',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_app_name` (`app_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=29108 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='指纹表';

CREATE TABLE IF NOT EXISTS `finger_result` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '插入时间',
  `cert` varchar(2000) NOT NULL DEFAULT '' COMMENT '证书内容',
  `favicon_hash` varchar(1000) NOT NULL DEFAULT '' COMMENT '图标hash',
  `img_ocr` varchar(100) NOT NULL DEFAULT '' COMMENT '图片ocr',
  `js_name` varchar(100) NOT NULL DEFAULT '' COMMENT 'js文件名',
  `tls_jarm_hash` varchar(100) NOT NULL DEFAULT '' COMMENT 'tls jarm 哈希',
  `url` varchar(100) NOT NULL DEFAULT '' COMMENT 'url',
  `web_body` text NOT NULL COMMENT 'web响应报文内容',
  `web_header` text NOT NULL COMMENT 'web响应头内容',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `history_tasks` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` varchar(50) NOT NULL COMMENT '任务id',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `task_data` longtext NOT NULL COMMENT '任务数据',
  `status` varchar(20) NOT NULL COMMENT '任务结束状态',
  `check_target` varchar(200) NOT NULL COMMENT '测试目标',
  `type` varchar(30) NOT NULL COMMENT '测试类型',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=39383 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='历史任务检测结果表';

CREATE TABLE IF NOT EXISTS `ip_host` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `ip` varchar(36) NOT NULL DEFAULT '' COMMENT 'Ip',
  `hosts` varchar(2056) NOT NULL DEFAULT '' COMMENT 'url',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='IP 域名绑定表';

CREATE TABLE IF NOT EXISTS `map_set` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `estate` varchar(64) NOT NULL DEFAULT 'valid' COMMENT '数据状态valid/deleted',
  `obj_key` varchar(255) NOT NULL DEFAULT '' COMMENT '字典键key',
  `obj_value` longtext NOT NULL COMMENT '字典值value',
  `content` varchar(255) NOT NULL DEFAULT '' COMMENT '备注说明',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_obj_key` (`obj_key`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8 COMMENT='枚举map表';

CREATE TABLE IF NOT EXISTS `node` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `node_id` varchar(20) NOT NULL DEFAULT '' COMMENT '节点id',
  `node_type` varchar(20) NOT NULL DEFAULT '' COMMENT '节点类型',
  `token` varchar(255) NOT NULL DEFAULT '' COMMENT '节点token',
  `main_ip_addr` varchar(50) NOT NULL DEFAULT '' COMMENT '主ip地址列表',
  `flow_task_id` varchar(255) NOT NULL DEFAULT '' COMMENT '流量分析任务id',
  `flow_server_id` varchar(50) NOT NULL DEFAULT '' COMMENT '流量服务id',
  `node_status` varchar(20) NOT NULL DEFAULT '' COMMENT '节点状态，开启代理 flow openvpn passive',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_node_id` (`node_id`),
  KEY `idx_flow_task_id` (`flow_task_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='节点表';

CREATE TABLE IF NOT EXISTS `sbtest1` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `k` int(11) NOT NULL DEFAULT '0',
  `c` char(120) NOT NULL DEFAULT '',
  `pad` char(60) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `k_1` (`k`)
) ENGINE=InnoDB AUTO_INCREMENT=10001 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `script_rule` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `script_name` varchar(255) DEFAULT NULL,
  `precondition` text,
  `consequence` text,
  `app_name` varchar(255) DEFAULT NULL,
  `remark` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=11320 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;

CREATE TABLE IF NOT EXISTS `use_score` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属目标id',
  `target_url` varchar(255) NOT NULL DEFAULT '' COMMENT '目标url',
  `tcp_scan_type` tinyint(4) unsigned NOT NULL DEFAULT '2' COMMENT '端口扫描类型',
  `to_scan_port` text NOT NULL COMMENT '端口扫描列表',
  `use_score` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '利用评分',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态，1-端口扫描进行中,2-端口扫描完成,3-评分计算完成',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_target_id` (`target_id`),
  KEY `idx_task_id` (`task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3587 DEFAULT CHARSET=utf8 COMMENT='利用评分表';

CREATE TABLE IF NOT EXISTS `use_score_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `task_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属任务id',
  `target_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属目标id',
  `content` text NOT NULL COMMENT '内容',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_task_id` (`task_id`) USING BTREE,
  KEY `idx_target_id` (`target_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=836 DEFAULT CHARSET=utf8 COMMENT='利用评分日志表';

CREATE TABLE IF NOT EXISTS `user` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '名字',
  `mobile` char(11) NOT NULL DEFAULT '0' COMMENT '手机号',
  `password` varchar(36) NOT NULL DEFAULT '' COMMENT '密码',
  `salf` varchar(36) NOT NULL DEFAULT '' COMMENT '密码盐值',
  `email` varchar(52) NOT NULL DEFAULT '' COMMENT 'email',
  `roles` text COMMENT '角色 JSON',
  `is_super` enum('Y','N') NOT NULL DEFAULT 'N' COMMENT '是否超管',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '激活状态:0删除 1=已禁用，2=已激活',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `mobile` (`mobile`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';

CREATE TABLE IF NOT EXISTS `vul_libraries` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `data_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '数据类型，1:真实数据 2:假数据',
  `vul_id` varchar(50) DEFAULT '' COMMENT '四维自定义漏洞id',
  `name` varchar(300) DEFAULT '' COMMENT '漏洞名称',
  `cve` varchar(30) DEFAULT '' COMMENT 'cve编号',
  `risk` int(11) DEFAULT '0' COMMENT '漏洞风险',
  `type` int(11) DEFAULT '0' COMMENT '漏洞类型',
  `class` int(11) DEFAULT '0' COMMENT '漏洞分类',
  `published_time` varchar(32) DEFAULT '' COMMENT '公开时间',
  `description` text COMMENT '漏洞描述',
  `affect_range` text COMMENT '影响范围',
  `exploit_impact` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '利用影响',
  `fix_suggest` text COMMENT '修复建议',
  `cnvd` varchar(255) DEFAULT '' COMMENT 'cnvd编号',
  `cnnvd` varchar(255) DEFAULT '' COMMENT 'cnnvd编号',
  `component` varchar(255) NOT NULL DEFAULT '' COMMENT '受影响的组建名称',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '脚本状态:0删除 1=已禁用，2=正常，其他状态请查看Enum文件',
  `create_time` datetime DEFAULT '1970-01-01 08:00:01' COMMENT '添加时间',
  `update_time` datetime DEFAULT '1970-01-01 08:00:01' COMMENT '更新时间',
  `status_msg` varchar(128) NOT NULL DEFAULT '' COMMENT '脚本状态描述',
  `priority` varchar(20) NOT NULL DEFAULT '' COMMENT '修复优先级',
  `operating_system` int(11) NOT NULL DEFAULT '0' COMMENT '1:windows, 2:linux, 3:unix, 4:其他',
  `pocname` varchar(255) DEFAULT NULL COMMENT '漏洞标识',
  `ms` varchar(50) DEFAULT NULL COMMENT 'ms漏洞编号',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `vul_id_index` (`vul_id`) USING BTREE COMMENT '漏洞id索引'
) ENGINE=InnoDB AUTO_INCREMENT=34154 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='漏洞表';

CREATE TABLE IF NOT EXISTS `vul_scripts` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` int(11) DEFAULT '0' COMMENT '作者ID',
  `script_name` varchar(255) DEFAULT '' COMMENT 'pocname',
  `type` varchar(50) DEFAULT '' COMMENT '工具类型',
  `lib_name` varchar(200) DEFAULT '' COMMENT '漏洞名称 冗余vul_libraries表的name值',
  `content` longtext,
  `vul_id` varchar(20) DEFAULT '' COMMENT '该脚本对应的漏洞编号',
  `verify_type` varchar(20) DEFAULT '' COMMENT '校验类型',
  `params` varchar(1024) NOT NULL DEFAULT '' COMMENT '最近一次的测试参数',
  `target_list` varchar(256) NOT NULL DEFAULT '' COMMENT '靶站地址 json格式',
  `target_status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '脚本靶站状态',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '脚本状态:0删除 1=已禁用，2=正常，其他状态请查看Enum文件',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `evidence_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '取证类型',
  `status_msg` varchar(128) NOT NULL DEFAULT '' COMMENT '脚本状态描述',
  `data_type` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '数据类型，1:真实数据 2:假数据',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `vul_id_index` (`vul_id`) USING BTREE COMMENT '漏洞id索引'
) ENGINE=InnoDB AUTO_INCREMENT=10702 DEFAULT CHARSET=utf8 ROW_FORMAT=DYNAMIC COMMENT='脚本表';

CREATE TABLE IF NOT EXISTS `yak_node` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `name` varchar(36) NOT NULL DEFAULT '' COMMENT '节点名称',
  `ip` varchar(36) NOT NULL DEFAULT '' COMMENT '节点IP',
  `port` varchar(36) NOT NULL DEFAULT '' COMMENT '节点端口',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态 0不在线 1在线',
  `is_disable` tinyint(4) NOT NULL DEFAULT '0' COMMENT '禁用状态 0启用 1禁用',
  `create_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  `update_time` datetime NOT NULL DEFAULT '1970-01-01 08:00:01' COMMENT '发生时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8 COMMENT='YAK分布式 节点表';

