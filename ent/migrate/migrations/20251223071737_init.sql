-- Create "admin_users" table
CREATE TABLE `admin_users` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `username` varchar(100) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `nickname` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `adminuser_deleted_at` (`deleted_at`),
  UNIQUE INDEX `adminuser_username_deleted_at` (`username`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "backend_users" table
CREATE TABLE `backend_users` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `username` varchar(100) NOT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `nickname` varchar(255) NOT NULL,
  `merchant_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `backenduser_deleted_at` (`deleted_at`),
  UNIQUE INDEX `backenduser_username_deleted_at` (`username`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "categories" table
CREATE TABLE `categories` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NULL,
  `inherit_tax_rate` bool NOT NULL DEFAULT 0,
  `tax_rate_id` char(36) NULL,
  `inherit_stall` bool NOT NULL DEFAULT 0,
  `stall_id` char(36) NULL,
  `product_count` bigint NOT NULL DEFAULT 0,
  `sort_order` bigint NOT NULL DEFAULT 32767,
  `parent_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `categories_categories_children` (`parent_id`),
  INDEX `category_deleted_at` (`deleted_at`),
  INDEX `category_merchant_id` (`merchant_id`),
  INDEX `category_store_id` (`store_id`),
  CONSTRAINT `categories_categories_children` FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
