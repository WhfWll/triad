-- 主机基线核查：扩展 baseline_check_result 字段长度（解决 actual_value 写入失败）
USE triad;

ALTER TABLE `baseline_check_result`
  MODIFY COLUMN `expected_value` text COMMENT '期望值',
  MODIFY COLUMN `actual_value` text COMMENT '实际值';
