CREATE TABLE `short_links` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `short_code` varchar(16) NOT NULL COMMENT '短码',
  `original_url` varchar(2048) NOT NULL COMMENT '原始长链',
  `url_hash` char(64) NOT NULL COMMENT 'original_url去除尾部/后的字符串哈希值, 用于查重',
  `visit_count` bigint unsigned NOT NULL DEFAULT 0 COMMENT '访问次数',
  `status` tinyint unsigned NOT NULL DEFAULT 1 COMMENT '状态: 不使用0; 1=active, 2=disabled',
  `expires_at` datetime DEFAULT NULL COMMENT '过期时间, 预留字段',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除时间, 预留字段',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_short_code` (`short_code`),
  KEY `idx_url_hash` (`url_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='短链主表';
