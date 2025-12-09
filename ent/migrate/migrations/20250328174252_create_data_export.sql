-- Modify "orders" table
ALTER TABLE `orders` ADD INDEX `order_store_id_deleted_at` (`store_id`, `deleted_at`), ADD INDEX `order_store_id_status_deleted_at` (`store_id`, `status`, `deleted_at`);
-- Modify "payments" table
ALTER TABLE `payments` DROP INDEX `payment_pay_biz_type_biz_id`, ADD INDEX `payment_pay_biz_type_biz_id_finished_at_deleted_at` (`pay_biz_type`, `biz_id`, `finished_at`, `deleted_at`);
-- Create "data_exports" table
CREATE TABLE `data_exports` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `store_id` bigint NOT NULL DEFAULT 0,
  `type` enum('order_list') NOT NULL,
  `status` enum('pending','success','failed') NOT NULL,
  `params` json NOT NULL,
  `failed_reason` varchar(255) NULL,
  `operator_type` enum('frontend','backend','admin','system') NOT NULL,
  `operator_id` bigint NOT NULL DEFAULT 0,
  `operator_name` varchar(255) NOT NULL,
  `file_name` varchar(255) NOT NULL,
  `url` varchar(255) NULL,
  PRIMARY KEY (`id`),
  INDEX `dataexport_deleted_at` (`deleted_at`),
  INDEX `dataexport_store_id_deleted_at` (`store_id`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
