DROP TABLE IF EXISTS `user_operation_logs`;
CREATE TABLE `user_operation_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `action` varchar(32) NOT NULL COMMENT '动作类型: login/create_link',
  `result` varchar(16) NOT NULL COMMENT '结果: success',
  `target_id` bigint unsigned DEFAULT NULL COMMENT '目标对象ID, 如短链ID',
  `target_code` varchar(32) DEFAULT NULL COMMENT '目标对象编码, 如短码',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id_id` (`user_id`, `id`),
  KEY `idx_user_id_action_id` (`user_id`, `action`, `id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户操作日志表';
