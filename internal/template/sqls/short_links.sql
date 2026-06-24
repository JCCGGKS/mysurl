DROP TABLE IF EXISTS `short_links`;
CREATE TABLE `short_links` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '创建者用户ID, V4新增',
  `short_code` varchar(11) DEFAULT NULL COMMENT '短码',
  `original_url` varchar(250) NOT NULL COMMENT '规范化后的长链',
  `url_hash` char(64) NOT NULL COMMENT '规范化长链的哈希值, 用于辅助查重',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间, 预留后续扩展',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态: 不使用0; 1=active, 2=disabled',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间, 预留字段',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_short_code` (`short_code`),
  UNIQUE KEY `uk_user_original_url` (`user_id`, `original_url`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_url_hash` (`url_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='短链主表';
