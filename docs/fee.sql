CREATE TABLE llm_user_consume_record (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键，自增',
  user_id VARCHAR(255) NOT NULL COMMENT '用户ID',
  used_reward_coins BIGINT DEFAULT 0 COMMENT '本次使用的奖励币数量',
  used_recharge_coins BIGINT DEFAULT 0 COMMENT '本次使用的充值币数量',
  caller_aid VARCHAR(64) DEFAULT NULL COMMENT '调用方AID',
  model VARCHAR(64) DEFAULT NULL COMMENT '模型',
  provider VARCHAR(64) DEFAULT NULL COMMENT '服务商',
  provider_aid VARCHAR(64) DEFAULT NULL COMMENT '服务商AID',
  input_tokens BIGINT DEFAULT 0 COMMENT '输入token数',
  output_tokens BIGINT DEFAULT 0 COMMENT '输出token数',
  cache_tokens BIGINT DEFAULT 0 COMMENT '缓存token数',
  input_price INT DEFAULT 0 COMMENT '输入token价格',
  output_price INT DEFAULT 0 COMMENT '输出token价格',
  cache_price INT DEFAULT 0 COMMENT '缓存token价格',
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  INDEX idx_user_id (user_id),
  INDEX idx_caller_aid (caller_aid)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户消费记录';

CREATE TABLE llm_user_consume_record_detail (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键，自增',
  user_id VARCHAR(255) NOT NULL COMMENT '用户ID',
  record_id BIGINT NOT NULL COMMENT '关联 user_consume_record.id',
  source_id BIGINT NOT NULL COMMENT '扣费来源ID（source_coin_records 表主键ID）',
  source_coin_type BIGINT DEFAULT NULL COMMENT '扣费来源类型（RewardType 主键ID）',
  consumed BIGINT DEFAULT 0 COMMENT '消费金额',
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX idx_record (record_id),
  INDEX idx_source (source_coin_type, source_id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户消费记录明细';

ALTER TABLE llm_user_consume_record
ADD COLUMN recharge_coins_after BIGINT DEFAULT 0 COMMENT '扣费后充值代币余额' after cache_price;

ALTER TABLE llm_user_consume_record
ADD COLUMN reward_coins_after BIGINT DEFAULT 0 COMMENT '扣费后奖励代币余额' after recharge_coins_after;

CREATE TABLE llm_user_token_usage (
  id VARCHAR(128) PRIMARY KEY COMMENT '请求ID',
  user_id VARCHAR(255) NOT NULL COMMENT '用户ID',
  model VARCHAR(64) COMMENT '模型名称',
  key_id VARCHAR(255) COMMENT 'API密钥ID',
  provider VARCHAR(64) COMMENT '服务提供商',
  caller_aid VARCHAR(64) COMMENT '调用方AID',
  provider_aid VARCHAR(64) COMMENT '提供商AID',
  token_usage JSON COMMENT 'token使用详情',
  status INT COMMENT '状态（1成功，0失败）',
  error VARCHAR(1024) COMMENT '错误信息',
  created TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COMMENT = '用户用量明细';

-- 奖励余额对账
SELECT
  user_id,
  phone_num,
  detail_total,
  reward_coins_balance
FROM
  (
    SELECT
      SUM(remaining_amount) AS detail_total,
      u.user_id,
      u.phone_num,
      u.reward_coins_balance
    FROM
      llm_user_coins u
      LEFT JOIN llm_user_coins_detail d ON u.user_id = d.user_id
    WHERE
      d.user_id IS NOT NULL
      AND d.source_coin_type != 0
    GROUP BY
      u.user_id,
      u.phone_num,
      u.reward_coins_balance
    HAVING
      SUM(remaining_amount) > 0
    ORDER BY
      u.user_id
  ) AS t
WHERE
  reward_coins_balance - detail_total != 0;

-- 充值余额对账
SELECT
  user_id,
  phone_num,
  detail_total,
  reward_coins_balance
FROM
  (
    SELECT
      SUM(remaining_amount) AS detail_total,
      u.user_id,
      u.phone_num,
      u.reward_coins_balance
    FROM
      llm_user_coins u
      LEFT JOIN llm_user_coins_detail d ON u.user_id = d.user_id
    WHERE
      d.user_id IS NOT NULL
      AND d.source_coin_type != 0
    GROUP BY
      u.user_id,
      u.phone_num,
      u.reward_coins_balance
    HAVING
      SUM(remaining_amount) > 0
    ORDER BY
      u.user_id
  ) AS t
WHERE
  reward_coins_balance - detail_total != 0;
