-- Modify "order_items" table
ALTER TABLE `order_items` ADD COLUMN `remark` varchar(500) NOT NULL;
-- Modify "orders" table
ALTER TABLE `orders` RENAME COLUMN `total` TO `total_price`, MODIFY COLUMN `member_id` bigint NOT NULL DEFAULT 0, ADD COLUMN `real_price` decimal(10,2) NOT NULL, ADD COLUMN `points_available` decimal(10,2) NOT NULL, ADD COLUMN `cash_paid` decimal(10,2) NOT NULL, ADD COLUMN `wechat_paid` decimal(10,2) NOT NULL, ADD COLUMN `alipay_paid` decimal(10,2) NOT NULL, ADD COLUMN `points_paid` decimal(10,2) NOT NULL, ADD COLUMN `store_name` varchar(255) NOT NULL, ADD COLUMN `table_name` varchar(255) NOT NULL, ADD COLUMN `people_number` bigint NOT NULL, ADD COLUMN `creator_name` varchar(255) NOT NULL, DROP INDEX `orders_frontend_users_orders`, DROP INDEX `orders_stores_orders`, DROP FOREIGN KEY `orders_frontend_users_orders`, DROP FOREIGN KEY `orders_stores_orders`;
-- Modify "stores" table
ALTER TABLE `stores` MODIFY COLUMN `need_audit` bool NOT NULL DEFAULT 1, MODIFY COLUMN `enabled` bool NOT NULL DEFAULT 1;
-- Create "order_logs" table
CREATE TABLE `order_logs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `event` enum('create','append_item','remove_item','turn_table','turn_item','paid','cancel','finish') NOT NULL,
  `operator_type` enum('frontend','backend','system') NOT NULL,
  `operator_id` bigint NOT NULL DEFAULT 0,
  `operator_name` varchar(255) NOT NULL,
  `order_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `order_logs_orders_logs` (`order_id`),
  INDEX `orderlog_deleted_at` (`deleted_at`),
  CONSTRAINT `order_logs_orders_logs` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
