-- Modify "categories" table
ALTER TABLE `categories`
ADD COLUMN `stall_categories` char(36) NULL,
ADD COLUMN `tax_fee_categories` char(36) NULL,
ADD INDEX `categories_stalls_categories` (`stall_categories`),
ADD INDEX `categories_stalls_stall` (`stall_id`),
ADD INDEX `categories_tax_fees_categories` (`tax_fee_categories`),
ADD INDEX `categories_tax_fees_tax_rate` (`tax_rate_id`);

-- Modify "menu_items" table
ALTER TABLE `menu_items`
DROP COLUMN `sale_rule`;

-- Modify "menus" table
ALTER TABLE `menus`
DROP COLUMN `distribution_rule`,
ADD COLUMN `store_id` char(36) NOT NULL,
ADD UNIQUE INDEX `menu_merchant_id_store_id_name_deleted_at` (`merchant_id`, `store_id`, `name`, `deleted_at`),
ADD INDEX `menu_store_id` (`store_id`);

-- Modify "product_spec_relations" table
ALTER TABLE `product_spec_relations` MODIFY COLUMN `packing_fee_id` char(36) NULL,
ADD INDEX `product_spec_relations_additional_fees_product_specs` (`packing_fee_id`),
ADD UNIQUE INDEX `productspecrelation_product_id_spec_id_deleted_at` (`product_id`, `spec_id`, `deleted_at`);

-- Modify "products" table
ALTER TABLE `products` ADD INDEX `products_stalls_products` (`stall_id`),
ADD INDEX `products_tax_fees_products` (`tax_rate_id`);

-- Create "profit_distribution_bills" table
CREATE TABLE `profit_distribution_bills` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `no` varchar(64) NOT NULL,
  `receivable_amount` decimal(19, 4) NOT NULL,
  `payment_amount` decimal(19, 4) NOT NULL,
  `status` enum ('unpaid', 'paid') NOT NULL DEFAULT "unpaid",
  `bill_date` date NOT NULL,
  `start_date` date NOT NULL,
  `end_date` date NOT NULL,
  `rule_snapshot` json NOT NULL,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `no` (`no`),
  INDEX `profitdistributionbill_deleted_at` (`deleted_at`),
  INDEX `profitdistributionbill_merchant_id` (`merchant_id`),
  INDEX `profitdistributionbill_store_id` (`store_id`),
  UNIQUE INDEX `profitdistributionbill_store_id_bill_date_deleted_at` (`store_id`, `bill_date`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "profit_distribution_rules" table
CREATE TABLE `profit_distribution_rules` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `merchant_id` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `split_ratio` decimal(19, 4) NOT NULL,
  `billing_cycle` enum ('daily', 'monthly') NOT NULL DEFAULT "daily",
  `effective_date` timestamp NOT NULL,
  `expiry_date` timestamp NOT NULL,
  `bill_generation_day` bigint NOT NULL DEFAULT 1,
  `status` enum ('enabled', 'disabled') NOT NULL DEFAULT "disabled",
  `store_count` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `profitdistributionrule_deleted_at` (`deleted_at`),
  INDEX `profitdistributionrule_merchant_id` (`merchant_id`),
  UNIQUE INDEX `profitdistributionrule_merchant_id_name_deleted_at` (`merchant_id`, `name`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "profit_distribution_rule_store_relations" table
CREATE TABLE `profit_distribution_rule_store_relations` (
  `profit_distribution_rule_id` char(36) NOT NULL,
  `store_id` char(36) NOT NULL,
  PRIMARY KEY (`profit_distribution_rule_id`, `store_id`),
  INDEX `profit_distribution_rule_store_relations_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
