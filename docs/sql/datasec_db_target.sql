-- 数据安全 - 数据库扫描目标库（连接信息加密存储）
USE triad;

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
