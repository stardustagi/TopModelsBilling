-- 模型阶梯价格配置表
CREATE TABLE models_tiered_pricing (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  model_id BIGINT NOT NULL COMMENT '模型ID，关联models_info.id',
  tier_start BIGINT NOT NULL COMMENT '阶梯起始token数（含）',
  tier_end BIGINT NOT NULL COMMENT '阶梯结束token数（含），-1表示无上限',
  input_price INT NOT NULL COMMENT '输入token价格',
  output_price INT NOT NULL COMMENT '输出token价格',
  cache_price INT NOT NULL COMMENT '缓存token价格',
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL,
  INDEX idx_model_id (model_id),
  INDEX idx_tier (model_id, tier_start)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模型阶梯价格配置';

-- 用户折扣率配置表
CREATE TABLE user_discount (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL COMMENT '用户ID',
  discount_rate INT NOT NULL DEFAULT 100 COMMENT '折扣率，100表示无折扣，90表示9折',
  created_at BIGINT NOT NULL,
  updated_at BIGINT NOT NULL,
  UNIQUE INDEX idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户折扣率配置';
