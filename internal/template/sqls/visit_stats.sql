DROP TABLE IF EXISTS `visit_stats`;
CREATE TABLE `visit_stats` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `short_link_id` bigint unsigned NOT NULL COMMENT '短链ID',
  `visit_count` bigint unsigned NOT NULL DEFAULT 0 COMMENT '访问次数',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_short_link_id` (`short_link_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='短链访问统计表';
