-- Modify "admin_users" table
ALTER TABLE `admin_users`
ADD COLUMN `account_type` varchar(255) NOT NULL DEFAULT "normal";

-- Create "cities" table
CREATE TABLE `cities` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `sort` bigint NOT NULL DEFAULT 0,
  `country_id` char(36) NOT NULL,
  `province_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `cities_countries_cities` (`country_id`),
  INDEX `city_deleted_at` (`deleted_at`),
  INDEX `city_province_id_country_id` (`province_id`, `country_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "countries" table
CREATE TABLE `countries` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `sort` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `country_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "districts" table
CREATE TABLE `districts` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `sort` bigint NOT NULL DEFAULT 0,
  `city_id` char(36) NOT NULL,
  `country_id` char(36) NOT NULL,
  `province_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `district_city_id_province_id_country_id` (`city_id`, `province_id`, `country_id`),
  INDEX `district_deleted_at` (`deleted_at`),
  INDEX `districts_countries_districts` (`country_id`),
  INDEX `districts_provinces_districts` (`province_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "merchants" table
CREATE TABLE `merchants` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `merchant_code` varchar(255) NOT NULL DEFAULT "",
  `merchant_name` varchar(50) NOT NULL DEFAULT "",
  `merchant_short_name` varchar(50) NOT NULL DEFAULT "",
  `merchant_type` enum ('brand', 'store') NOT NULL,
  `brand_name` varchar(255) NOT NULL DEFAULT "",
  `admin_phone_number` varchar(20) NOT NULL DEFAULT "",
  `expire_utc` timestamp NULL,
  `merchant_logo` varchar(500) NOT NULL DEFAULT "",
  `description` varchar(255) NOT NULL DEFAULT "",
  `status` enum ('active', 'expired', 'disabled') NOT NULL,
  `address` varchar(255) NOT NULL DEFAULT "",
  `lng` varchar(255) NOT NULL DEFAULT "",
  `lat` varchar(255) NOT NULL DEFAULT "",
  `admin_user_id` char(36) NOT NULL,
  `city_id` char(36) NULL,
  `country_id` char(36) NULL,
  `district_id` char(36) NULL,
  `business_type_id` char(36) NOT NULL,
  `province_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `merchant_deleted_at` (`deleted_at`),
  INDEX `merchants_admin_users_merchant` (`admin_user_id`),
  INDEX `merchants_cities_merchants` (`city_id`),
  INDEX `merchants_countries_merchants` (`country_id`),
  INDEX `merchants_districts_merchants` (`district_id`),
  INDEX `merchants_merchant_business_types_merchants` (`business_type_id`),
  INDEX `merchants_provinces_merchants` (`province_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "merchant_business_types" table
CREATE TABLE `merchant_business_types` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `type_code` varchar(50) NOT NULL DEFAULT "",
  `type_name` varchar(50) NOT NULL DEFAULT "",
  PRIMARY KEY (`id`),
  INDEX `merchantbusinesstype_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "merchant_renewals" table
CREATE TABLE `merchant_renewals` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `purchase_duration` bigint NOT NULL DEFAULT 0,
  `purchase_duration_unit` enum ('day', 'month', 'year', 'week') NOT NULL,
  `operator_name` varchar(50) NOT NULL DEFAULT "",
  `operator_account` varchar(50) NOT NULL DEFAULT "",
  `merchant_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `merchantrenewal_deleted_at` (`deleted_at`),
  INDEX `merchantrenewal_merchant_id` (`merchant_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "provinces" table
CREATE TABLE `provinces` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `sort` bigint NOT NULL DEFAULT 0,
  `country_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `province_deleted_at` (`deleted_at`),
  INDEX `provinces_countries_provinces` (`country_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "remarks" table
CREATE TABLE `remarks` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(50) NOT NULL,
  `remark_type` enum ('system', 'brand') NOT NULL,
  `enabled` bool NOT NULL DEFAULT 1,
  `sort_order` bigint NOT NULL DEFAULT 1000,
  `merchant_id` char(36) NULL,
  `category_id` char(36) NOT NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `remark_category_id` (`category_id`),
  INDEX `remark_deleted_at` (`deleted_at`),
  INDEX `remark_merchant_id` (`merchant_id`),
  INDEX `remark_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "remark_categories" table
CREATE TABLE `remark_categories` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(50) NOT NULL,
  `remark_scene` enum (
    'whole_order',
    'item',
    'cancel_reason',
    'discount',
    'gift',
    'rebill',
    'refund_reject'
  ) NOT NULL,
  `description` varchar(255) NOT NULL DEFAULT "",
  `sort_order` bigint NOT NULL DEFAULT 1000,
  `merchant_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `remarkcategory_deleted_at` (`deleted_at`),
  INDEX `remarkcategory_merchant_id` (`merchant_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "stores" table
CREATE TABLE `stores` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `admin_phone_number` varchar(20) NOT NULL DEFAULT "",
  `store_name` varchar(30) NOT NULL DEFAULT "",
  `store_short_name` varchar(30) NOT NULL DEFAULT "",
  `store_code` varchar(255) NOT NULL DEFAULT "",
  `status` enum ('open', 'closed') NOT NULL,
  `business_model` enum ('direct', 'franchisee') NOT NULL,
  `location_number` varchar(255) NOT NULL,
  `contact_name` varchar(20) NOT NULL DEFAULT "",
  `contact_phone` varchar(20) NOT NULL DEFAULT "",
  `unified_social_credit_code` varchar(50) NOT NULL DEFAULT "",
  `store_logo` varchar(500) NOT NULL DEFAULT "",
  `business_license_url` varchar(500) NOT NULL DEFAULT "",
  `storefront_url` varchar(500) NOT NULL DEFAULT "",
  `cashier_desk_url` varchar(500) NOT NULL DEFAULT "",
  `dining_environment_url` varchar(500) NOT NULL DEFAULT "",
  `food_operation_license_url` varchar(500) NOT NULL DEFAULT "",
  `business_hours` varchar(255) NOT NULL DEFAULT "",
  `dining_periods` varchar(255) NOT NULL,
  `shift_times` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL DEFAULT "",
  `lng` varchar(50) NOT NULL DEFAULT "",
  `lat` varchar(50) NOT NULL DEFAULT "",
  `admin_user_id` char(36) NOT NULL,
  `city_id` char(36) NULL,
  `country_id` char(36) NULL,
  `district_id` char(36) NULL,
  `merchant_id` char(36) NOT NULL,
  `business_type_id` char(36) NOT NULL,
  `province_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `store_deleted_at` (`deleted_at`),
  INDEX `store_merchant_id` (`merchant_id`),
  INDEX `stores_admin_users_store` (`admin_user_id`),
  INDEX `stores_cities_stores` (`city_id`),
  INDEX `stores_countries_stores` (`country_id`),
  INDEX `stores_districts_stores` (`district_id`),
  INDEX `stores_merchant_business_types_stores` (`business_type_id`),
  INDEX `stores_provinces_stores` (`province_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
