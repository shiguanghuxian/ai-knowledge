CREATE TABLE `knowledge` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `question` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '问题',
  `answer` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci COMMENT '答案',
  `text` text COLLATE utf8mb4_general_ci COMMENT '知识内容',
  `vector_id` bigint DEFAULT NULL COMMENT '向量id',
  `type` tinyint NOT NULL DEFAULT '0' COMMENT '类型 0问答 1纯知识',
  `created_at` bigint DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint DEFAULT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_vector_id` (`vector_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='知识库';
