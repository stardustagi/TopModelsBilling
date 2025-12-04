-- 给 user_consume_detail_image 表添加 consume_id 字段（如果表已存在）
ALTER TABLE user_consume_detail_image ADD COLUMN consume_id BIGINT COMMENT '消费记录id' AFTER id;

-- 完整建表语句
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
