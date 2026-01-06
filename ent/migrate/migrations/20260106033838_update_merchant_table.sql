-- Modify "admin_users" table
ALTER TABLE `admin_users`
DROP COLUMN `account_type`;

-- Modify "backend_users" table
ALTER TABLE `backend_users`
ADD COLUMN `merchant_backend_users` char(36) NULL,
ADD INDEX `backend_users_merchants_backend_users` (`merchant_backend_users`);

-- Modify "merchants" table
ALTER TABLE `merchants` MODIFY COLUMN `merchant_code` varchar(255) NULL DEFAULT "",
MODIFY COLUMN `merchant_short_name` varchar(50) NULL DEFAULT "",
MODIFY COLUMN `brand_name` varchar(255) NULL DEFAULT "",
MODIFY COLUMN `description` varchar(255) NULL DEFAULT "",
MODIFY COLUMN `address` varchar(255) NULL DEFAULT "",
MODIFY COLUMN `lng` varchar(255) NULL DEFAULT "",
MODIFY COLUMN `lat` varchar(255) NULL DEFAULT "",
DROP COLUMN `admin_user_id`,
ADD COLUMN `super_account` varchar(255) NOT NULL;

-- Modify "remarks" table
ALTER TABLE `remarks` MODIFY COLUMN `remark_type` enum ('system', 'brand', 'store') NOT NULL;

-- Modify "stores" table
ALTER TABLE `stores` MODIFY COLUMN `store_short_name` varchar(30) NULL DEFAULT "",
MODIFY COLUMN `store_code` varchar(255) NULL DEFAULT "",
MODIFY COLUMN `contact_name` varchar(20) NULL DEFAULT "",
MODIFY COLUMN `contact_phone` varchar(20) NULL DEFAULT "",
MODIFY COLUMN `unified_social_credit_code` varchar(50) NULL DEFAULT "",
MODIFY COLUMN `store_logo` varchar(500) NULL DEFAULT "",
MODIFY COLUMN `business_license_url` varchar(500) NULL DEFAULT "",
MODIFY COLUMN `storefront_url` varchar(500) NULL DEFAULT "",
MODIFY COLUMN `cashier_desk_url` varchar(500) NULL DEFAULT "",
MODIFY COLUMN `dining_environment_url` varchar(500) NULL DEFAULT "",
MODIFY COLUMN `food_operation_license_url` varchar(500) NULL DEFAULT "",
MODIFY COLUMN `business_hours` json NOT NULL,
MODIFY COLUMN `dining_periods` json NOT NULL,
MODIFY COLUMN `shift_times` json NOT NULL,
MODIFY COLUMN `lng` varchar(50) NULL DEFAULT "",
MODIFY COLUMN `lat` varchar(50) NULL DEFAULT "",
DROP COLUMN `admin_user_id`,
ADD COLUMN `super_account` varchar(255) NOT NULL;

-- Create "additional_fees" table
CREATE TABLE `additional_fees` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(50) NOT NULL,
  `fee_type` enum ('merchant', 'store') NOT NULL,
  `fee_category` enum ('service_fee', 'additional_fee', 'packing_fee') NOT NULL,
  `charge_mode` enum ('percent', 'fixed') NOT NULL,
  `fee_value` decimal(19, 4) NOT NULL,
  `include_in_receivable` bool NOT NULL DEFAULT 0,
  `taxable` bool NOT NULL DEFAULT 0,
  `discount_scope` enum ('before_discount', 'after_discount') NOT NULL,
  `order_channels` json NOT NULL,
  `dining_ways` json NOT NULL,
  `enabled` bool NOT NULL DEFAULT 1,
  `sort_order` bigint NOT NULL DEFAULT 1000,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `additionalfee_deleted_at` (`deleted_at`),
  INDEX `additionalfee_merchant_id` (`merchant_id`),
  INDEX `additionalfee_store_id` (`store_id`),
  UNIQUE INDEX `idx_additional_fee_name_merchant_store_deleted` (`name`, `merchant_id`, `store_id`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "departments" table
CREATE TABLE `departments` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `code` varchar(255) NOT NULL,
  `department_type` enum ('admin', 'backend', 'store') NOT NULL,
  `enable` bool NOT NULL DEFAULT 1,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `department_code_deleted_at` (`code`, `deleted_at`),
  INDEX `department_deleted_at` (`deleted_at`),
  INDEX `department_merchant_id_store_id` (`merchant_id`, `store_id`),
  INDEX `department_name` (`name`),
  INDEX `departments_stores_departments` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "devices" table
CREATE TABLE `devices` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(50) NOT NULL,
  `device_type` enum ('cashier', 'printer') NOT NULL,
  `device_code` varchar(100) NOT NULL DEFAULT "",
  `device_brand` varchar(255) NULL,
  `device_model` varchar(255) NULL,
  `location` enum ('front_hall', 'back_kitchen') NOT NULL,
  `enabled` bool NOT NULL DEFAULT 1,
  `status` enum ('online', 'offline') NULL DEFAULT "offline",
  `ip` varchar(50) NULL DEFAULT "",
  `sort_order` bigint NULL DEFAULT 1000,
  `paper_size` enum ('58mm', '80mm') NULL,
  `order_channels` json NULL,
  `dining_ways` json NULL,
  `device_stall_print_type` enum ('all', 'combined', 'separate') NULL,
  `device_stall_receipt_type` enum ('all', 'exclude') NULL,
  `open_cash_drawer` bool NULL,
  `merchant_id` char(36) NOT NULL,
  `stall_id` char(36) NULL,
  `store_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `device_deleted_at` (`deleted_at`),
  INDEX `device_merchant_id` (`merchant_id`),
  INDEX `device_stall_id` (`stall_id`),
  INDEX `device_store_id` (`store_id`),
  UNIQUE INDEX `idx_device_name_scope_deleted` (`name`, `merchant_id`, `store_id`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "roles" table
CREATE TABLE `roles` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `code` varchar(255) NOT NULL,
  `role_type` enum ('admin', 'backend', 'store') NOT NULL,
  `enable` bool NOT NULL DEFAULT 1,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `role_code_deleted_at` (`code`, `deleted_at`),
  INDEX `role_deleted_at` (`deleted_at`),
  INDEX `role_merchant_id_store_id` (`merchant_id`, `store_id`),
  INDEX `role_name` (`name`),
  INDEX `roles_stores_roles` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "stalls" table
CREATE TABLE `stalls` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(20) NOT NULL,
  `stall_type` enum ('system', 'brand', 'store') NOT NULL,
  `print_type` enum ('receipt', 'label') NOT NULL,
  `enabled` bool NOT NULL DEFAULT 1,
  `sort_order` bigint NOT NULL DEFAULT 0,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_stall_name_merchant_store_deleted` (`name`, `merchant_id`, `store_id`, `deleted_at`),
  INDEX `stall_deleted_at` (`deleted_at`),
  INDEX `stall_merchant_id` (`merchant_id`),
  INDEX `stall_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "store_users" table
CREATE TABLE `store_users` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `username` varchar(100) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `nickname` varchar(255) NOT NULL,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `store_users_merchants_store_users` (`merchant_id`),
  INDEX `store_users_stores_store_users` (`store_id`),
  INDEX `storeuser_deleted_at` (`deleted_at`),
  UNIQUE INDEX `storeuser_username_deleted_at` (`username`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "tax_fees" table
CREATE TABLE `tax_fees` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(50) NOT NULL,
  `tax_fee_type` enum ('merchant', 'store') NOT NULL,
  `tax_code` varchar(50) NOT NULL,
  `tax_rate_type` enum ('unified', 'custom') NOT NULL DEFAULT "unified",
  `tax_rate` decimal(19, 4) NOT NULL,
  `default_tax` bool NOT NULL DEFAULT 0,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_tax_fee_name_merchant_store_deleted` (`name`, `merchant_id`, `store_id`, `deleted_at`),
  INDEX `taxfee_deleted_at` (`deleted_at`),
  INDEX `taxfee_merchant_id` (`merchant_id`),
  INDEX `taxfee_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
