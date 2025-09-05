DELIMITER $$
DROP PROCEDURE IF EXISTS generate_test_data $$
CREATE PROCEDURE generate_test_data (
  IN p_user_id VARCHAR(255), -- 用户ID作为参数传入
  IN total_count INT, -- 总生成条数
  IN batch_size INT -- 每批插入多少条
) BEGIN DECLARE i INT DEFAULT 0;

DECLARE current_batch INT;

DECLARE total_amount BIGINT DEFAULT 0;

DECLARE user_exists INT DEFAULT 0;

DECLARE user_phone VARCHAR(255) DEFAULT '';

-- ✳️ Step 1：清空该用户明细
DELETE FROM llm_user_coins_detail
WHERE
  user_id = p_user_id;

DELETE FROM llm_user_consume_record
WHERE
  user_id = p_user_id;

DELETE FROM llm_user_consume_record_detail
WHERE
  user_id = p_user_id;

select
  phone INTO user_phone
from
  llm_user
where
  user_id = p_user_id;

-- ✳️ Step 2：生成明细数据
WHILE i < total_count DO
SET
  current_batch = 0;

SET
  @sql := 'INSERT INTO llm_user_coins_detail (nick_name, user_id, phone_num, source_coin_type, amount, remaining_amount, expiration_time, invite_code) VALUES ';

WHILE current_batch < batch_size
AND i < total_count DO
SET
  @nick := '';

SET
  @phone := user_phone;

SET
  @source_coin_type := FLOOR(1 + RAND() * 3);

SET
  @amount := FLOOR(10000 + RAND() * 10000) * 1000;

SET
  @expire := DATE_ADD(NOW(), INTERVAL FLOOR(1 + RAND() * 365) DAY);

SET
  @invite := UUID();

SET
  total_amount = total_amount + @amount;

SET
  @sql := CONCAT(
    @sql,
    '(',
    QUOTE(@nick),
    ',',
    QUOTE(p_user_id),
    ',',
    QUOTE(@phone),
    ',',
    @source_coin_type,
    ',',
    @amount,
    ',',
    @amount,
    ',',
    QUOTE(@expire),
    ',',
    QUOTE(@invite),
    ')',
    IF(
      current_batch < batch_size - 1
      AND i < total_count - 1,
      ',',
      ''
    )
  );

SET
  current_batch = current_batch + 1;

SET
  i = i + 1;

END
WHILE;

PREPARE stmt
FROM
  @sql;

EXECUTE stmt;

DEALLOCATE PREPARE stmt;

END
WHILE;

-- ✳️ Step 3：判断 user_coins 中是否已存在该用户
SELECT
  COUNT(*) INTO user_exists
FROM
  llm_user_coins
WHERE
  user_id = p_user_id;

SET
  @recharge := FLOOR(10000 + RAND() * 90000) * 1000;

IF user_exists > 0 THEN
-- 用户存在，执行 UPDATE
SET
  @update_sql := CONCAT(
    'UPDATE llm_user_coins SET ',
    'reward_coins_balance = ',
    total_amount,
    ', ',
    'recharge_coin_balance = ',
    @recharge,
    ' ',
    'WHERE user_id = ',
    QUOTE(p_user_id)
  );

ELSE
-- 用户不存在，执行 INSERT
SET
  @update_sql := CONCAT(
    'INSERT INTO llm_user_coins (user_id, nick_name, phone_num, reward_coins_balance, recharge_coin_balance) VALUES (',
    QUOTE(p_user_id),
    ', ',
    QUOTE('测试用户'),
    ', ',
    QUOTE('13900000000'),
    ', ',
    total_amount,
    ', ',
    @recharge,
    ')'
  );

END IF;

PREPARE update_stmt
FROM
  @update_sql;

EXECUTE update_stmt;

DEALLOCATE PREPARE update_stmt;

END $$ DELIMITER;

-- 为用户 mg1750851115730592 生成 500 条测试数据，每批 100 条 2180b152775c4f6cb8a9c78ac6ea2cab
CALL generate_test_data ('3483956f7ba843c4a2285619c57f6b40', 1, 1);
