-- Create "payment_methods" table
CREATE TABLE `payment_methods` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `accounting_rule` enum ('income', 'discount') NOT NULL DEFAULT "income",
  `payment_type` enum (
    'other',
    'cash',
    'offline_card',
    'custom_coupon',
    'partner_coupon'
  ) NOT NULL DEFAULT "other",
  `fee_rate` decimal(10, 2) NULL,
  `invoice_rule` enum ('no_invoice', 'actual_amount') NULL,
  `cash_drawer_status` bool NOT NULL DEFAULT 0,
  `display_channels` json NOT NULL,
  `source` enum ('brand', 'store', 'system') NOT NULL DEFAULT "brand",
  `store_ids` json NULL,
  `status` bool NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `paymentmethod_deleted_at` (`deleted_at`),
  INDEX `paymentmethod_merchant_id_store_id` (`merchant_id`, `store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
