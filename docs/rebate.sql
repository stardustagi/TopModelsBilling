-- 用户消费返点阶梯配置表
CREATE TABLE IF NOT EXISTS `user_rebate_config` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `tier_start` bigint NOT NULL COMMENT '阶梯起始金额',
  `tier_end` bigint NOT NULL COMMENT '阶梯结束金额，-1表示无上限',
  `rebate_rate` int DEFAULT 0 COMMENT '返点比例，如10表示10%',
  `status` int DEFAULT 1 COMMENT '状态1启用0禁用',
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户消费返点阶梯配置表';

-- 用户月度返点记录表
CREATE TABLE IF NOT EXISTS `user_rebate_monthly` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `month` varchar(7) NOT NULL COMMENT '月份格式2024-01',
  `total_consumed` bigint DEFAULT 0 COMMENT '当月消费总额',
  `rebate_amount` bigint DEFAULT 0 COMMENT '已返点金额',
  `rebate_rate` int DEFAULT 0 COMMENT '返点比例快照',
  `status` int DEFAULT 0 COMMENT '0未返点1已返点',
  `created_at` bigint DEFAULT NULL,
  `updated_at` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_month` (`user_id`, `month`),
  KEY `idx_month` (`month`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户月度返点记录表';

-- user_wallet表添加返点余额字段
ALTER TABLE user_wallet
ADD COLUMN rebate_balance BIGINT(12) DEFAULT 0 COMMENT '返点余额' AFTER balance;

ALTER TABLE user_consume_record
ADD COLUMN rebate_deducted BIGINT DEFAULT 0 COMMENT '返点扣除金额' AFTER total_consumed;

ALTER TABLE user_consume_record
ADD COLUMN balance_deducted BIGINT DEFAULT 0 COMMENT '余额扣除金额' AFTER rebate_deducted;
