-- Create "payment_accounts" table
CREATE TABLE `payment_accounts` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `merchant_id` char(36) NOT NULL,
  `channel` enum ('rm') NOT NULL,
  `merchant_number` varchar(255) NOT NULL,
  `merchant_name` varchar(255) NOT NULL,
  `is_default` bool NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `paymentaccount_deleted_at` (`deleted_at`),
  INDEX `paymentaccount_merchant_id` (`merchant_id`),
  UNIQUE INDEX `paymentaccount_merchant_id_channel_deleted_at` (`merchant_id`, `channel`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "store_payment_accounts" table
CREATE TABLE `store_payment_accounts` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `merchant_id` char(36) NOT NULL,
  `merchant_number` varchar(255) NOT NULL,
  `payment_account_id` char(36) NOT NULL,
  `store_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `storepaymentaccount_deleted_at` (`deleted_at`),
  INDEX `storepaymentaccount_payment_account_id` (`payment_account_id`),
  INDEX `storepaymentaccount_store_id` (`store_id`),
  UNIQUE INDEX `storepaymentaccount_store_id_payment_account_id_deleted_at` (`store_id`, `payment_account_id`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
