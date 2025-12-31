-- 供应商消费汇总表
CREATE TABLE IF NOT EXISTS `provider_consume_summary` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键，自增',
  `actual_provider_id` int NOT NULL COMMENT '实际服务商ID',
  `total_consumed` bigint DEFAULT 0 COMMENT '总消费金额',
  `total_cost` bigint DEFAULT 0 COMMENT '总成本',
  `updated_at` bigint DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_actual_provider_id` (`actual_provider_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '供应商消费汇总表';

-- 供应商模型日消费汇总表
CREATE TABLE IF NOT EXISTS `provider_model_daily_summary` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键，自增',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `actual_provider_id` int NOT NULL COMMENT '实际服务商ID',
  `model_id` int NOT NULL COMMENT '模型ID',
  `consume_type` varchar(32) NOT NULL COMMENT '消费类型',
  `date` varchar(10) NOT NULL COMMENT '日期YYYY-MM-DD',
  `total_consumed` bigint DEFAULT 0 COMMENT '总消费金额',
  `total_cost` bigint DEFAULT 0 COMMENT '总成本',
  `updated_at` bigint DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_provider_model_type_date` (
    `user_id`,
    `actual_provider_id`,
    `model_id`,
    `consume_type`,
    `date`
  ),
  KEY `idx_date` (`date`),
  KEY `idx_user_id` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '供应商模型日消费汇总表';

ALTER TABLE provider_consume_summary
ADD UNIQUE INDEX idx_provider_month (actual_provider_id, month);
