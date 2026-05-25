-- 安全检查报告中心表
-- 用于存储主机安全、应用安全、数据安全生成的 HTML 报告

CREATE TABLE IF NOT EXISTS `security_report` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '报告标题',
  `module` varchar(32) NOT NULL DEFAULT '' COMMENT '所属模块：host/app/data',
  `task_id` int(11) NOT NULL DEFAULT 0 COMMENT '关联任务 ID',
  `task_name` varchar(255) NOT NULL DEFAULT '' COMMENT '关联任务名称',
  `content` longtext NOT NULL COMMENT 'HTML 报告内容',
  `create_by` int(11) NOT NULL DEFAULT 0 COMMENT '创建人用户 ID',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_module` (`module`),
  KEY `idx_task_id` (`task_id`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='安全检查报告';
