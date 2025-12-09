-- Modify "payments" table
ALTER TABLE `payments` MODIFY COLUMN `resp` json NULL, ADD COLUMN `store_id` bigint NOT NULL, ADD INDEX `payment_pay_biz_type_biz_id` (`pay_biz_type`, `biz_id`);
-- Modify "stores" table
ALTER TABLE `stores` ADD COLUMN `huifu_id` varchar(255) NOT NULL, ADD COLUMN `zxh_id` varchar(255) NOT NULL, ADD COLUMN `zxh_secret` varchar(255) NOT NULL;
