-- Modify "orders" table
ALTER TABLE `orders` ADD COLUMN `paid` decimal(10,2) NOT NULL, ADD COLUMN `refunded` decimal(10,2) NOT NULL, ADD COLUMN `cash_refunded` decimal(10,2) NOT NULL, ADD COLUMN `wechat_refunded` decimal(10,2) NOT NULL, ADD COLUMN `alipay_refunded` decimal(10,2) NOT NULL, ADD COLUMN `points_refunded` decimal(10,2) NOT NULL;
-- Create "order_finance_logs" table
CREATE TABLE `order_finance_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `order_id` bigint NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `type` enum('paid','refund') NOT NULL,
  `channel` enum('cash','wechat','alipay','point') NOT NULL,
  `seq_no` varchar(255) NOT NULL,
  `creator_type` enum('frontend','backend','system') NOT NULL,
  `creator_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `orderfinancelog_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "payments" table
CREATE TABLE `payments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `seq_no` varchar(255) NOT NULL,
  `provider` enum('zxh') NOT NULL,
  `channel` enum('wxpay','alipay','point') NOT NULL,
  `state` enum('U','P','S','F','W') NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `goods_desc` varchar(255) NOT NULL,
  `mch_id` varchar(255) NOT NULL,
  `ip_addr` varchar(255) NOT NULL,
  `req` json NOT NULL,
  `resp` json NOT NULL,
  `callback` json NOT NULL,
  `finished_at` timestamp NULL,
  `refunded` decimal(10,2) NOT NULL,
  `fail_reason` varchar(255) NULL,
  `pay_biz_type` enum('order') NOT NULL,
  `biz_id` bigint NOT NULL,
  `creator` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `payment_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
