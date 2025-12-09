-- Create "reconciliation_records" table
CREATE TABLE `reconciliation_records` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `store_id` bigint NOT NULL,
  `store_name` varchar(255) NOT NULL,
  `order_count` bigint NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `channel` enum('cash','wechat','alipay','point') NOT NULL,
  `date` date NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `reconciliationrecord_date_store_id_channel` (`date`, `store_id`, `channel`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
