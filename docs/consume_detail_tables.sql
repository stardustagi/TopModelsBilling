-- user_consume_record 主表
CREATE TABLE IF NOT EXISTS user_consume_record (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键，自增',
  user_id INT NOT NULL COMMENT '用户ID',
  node_id INT COMMENT 'node编号',
  discount_amount BIGINT DEFAULT 0 COMMENT '折扣数量',
  total_consumed BIGINT DEFAULT 0 COMMENT '本次使用的币数量',
  caller VARCHAR(64) COMMENT '调用方',
  model VARCHAR(64) COMMENT '模型',
  model_id INT COMMENT '模型id',
  actual_provider VARCHAR(64) COMMENT '服务商',
  actual_provider_id VARCHAR(64) COMMENT '服务商id',
  consume_type VARCHAR(255) DEFAULT '' COMMENT '消费类型',
  total_cost BIGINT DEFAULT 0 COMMENT '陈本',
  created_at BIGINT COMMENT '创建时间',
  updated_at BIGINT COMMENT '更新时间',
  INDEX idx_user_id (user_id),
  INDEX idx_caller (caller)
);

-- 明细表
CREATE TABLE IF NOT EXISTS user_consume_detail_text (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键，自增',
  consume_id BIGINT COMMENT '消费记录id',
  input_tokens BIGINT DEFAULT 0 COMMENT '输入token数',
  output_tokens BIGINT DEFAULT 0 COMMENT '输出token数',
  cache_tokens BIGINT DEFAULT 0 COMMENT '缓存token数',
  input_price INT DEFAULT 0 COMMENT '输入token价格',
  output_price INT DEFAULT 0 COMMENT '输出token价格',
  cache_price INT DEFAULT 0 COMMENT '缓存token价格',
  created_at BIGINT COMMENT '创建时间'
);

CREATE TABLE IF NOT EXISTS user_consume_detail_image (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键，自增',
  consume_id BIGINT COMMENT '消费记录id',
  quality VARCHAR(64) COMMENT 'Quality',
  size VARCHAR(64) COMMENT 'Size',
  created_at BIGINT COMMENT '创建时间'
);

CREATE TABLE IF NOT EXISTS user_consume_detail_video (
  id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '主键，自增',
  consume_id BIGINT COMMENT '消费记录id',
  seconds DOUBLE COMMENT 'Seconds',
  size VARCHAR(64) COMMENT 'Size',
  created_at BIGINT COMMENT '创建时间'
);
