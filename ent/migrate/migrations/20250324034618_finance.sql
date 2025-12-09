-- Modify "reconciliation_records" table
ALTER TABLE `reconciliation_records` ADD COLUMN `no` varchar(255) NOT NULL, ADD UNIQUE INDEX `no` (`no`);
-- Create "point_settlements" table
CREATE TABLE `point_settlements` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `no` varchar(255) NOT NULL,
  `store_id` bigint NOT NULL,
  `store_name` varchar(255) NOT NULL,
  `order_count` bigint NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `total_points` decimal(10,2) NOT NULL,
  `date` date NOT NULL,
  `status` bigint NOT NULL DEFAULT 1,
  `point_settlement_rate` decimal(5,4) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `no` (`no`),
  UNIQUE INDEX `pointsettlement_date_store_id` (`date`, `store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "store_account_transactions" table
CREATE TABLE `store_account_transactions` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `no` varchar(255) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `after` decimal(10,2) NOT NULL,
  `type` bigint NOT NULL,
  `store_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `no` (`no`),
  INDEX `store_account_transactions_stores_store_account_transactions` (`store_id`),
  INDEX `storeaccounttransaction_deleted_at` (`deleted_at`),
  CONSTRAINT `store_account_transactions_stores_store_account_transactions` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "store_accounts" table
CREATE TABLE `store_accounts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `balance` decimal(10,2) NOT NULL,
  `pending_withdraw` decimal(10,2) NOT NULL,
  `withdrawn` decimal(10,2) NOT NULL,
  `total_amount` decimal(10,2) NOT NULL,
  `store_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `store_id` (`store_id`),
  INDEX `storeaccount_deleted_at` (`deleted_at`),
  CONSTRAINT `store_accounts_stores_store_account` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "store_withdraws" table
CREATE TABLE `store_withdraws` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `store_name` varchar(255) NOT NULL,
  `no` varchar(255) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `point_withdrawal_rate` decimal(5,4) NOT NULL,
  `actual_amount` decimal(10,2) NOT NULL,
  `account_type` enum('publish','private') NOT NULL,
  `bank_account` varchar(255) NOT NULL,
  `bank_card_name` varchar(255) NOT NULL,
  `bank_name` varchar(255) NOT NULL,
  `bank_branch` varchar(255) NOT NULL,
  `invoice_amount` decimal(10,2) NOT NULL,
  `status` bigint NOT NULL,
  `store_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `no` (`no`),
  INDEX `store_withdraws_stores_store_withdraws` (`store_id`),
  INDEX `storewithdraw_deleted_at` (`deleted_at`),
  CONSTRAINT `store_withdraws_stores_store_withdraws` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
