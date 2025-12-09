-- Modify "order_finance_logs" table
ALTER TABLE `order_finance_logs` MODIFY COLUMN `creator_type` enum('frontend','backend','admin','system') NOT NULL, MODIFY COLUMN `creator_id` bigint NOT NULL DEFAULT 0, ADD COLUMN `creator_name` varchar(255) NOT NULL;
-- Modify "order_logs" table
ALTER TABLE `order_logs` MODIFY COLUMN `operator_type` enum('frontend','backend','admin','system') NOT NULL;
-- Modify "payments" table
ALTER TABLE `payments` DROP COLUMN `creator`, ADD COLUMN `creator_type` enum('frontend','backend','admin','system') NOT NULL, ADD COLUMN `creator_id` bigint NOT NULL DEFAULT 0, ADD COLUMN `creator_name` varchar(255) NOT NULL, ADD UNIQUE INDEX `seq_no` (`seq_no`);
-- Create "payment_callbacks" table
CREATE TABLE `payment_callbacks` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `seq_no` varchar(255) NOT NULL,
  `type` enum('pay','refund') NOT NULL,
  `raw` json NOT NULL,
  `provider` enum('zxh') NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `paymentcallback_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
