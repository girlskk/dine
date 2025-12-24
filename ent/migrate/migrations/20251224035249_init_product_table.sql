-- Create "products" table
CREATE TABLE `products` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `type` enum ('normal', 'set_meal') NOT NULL DEFAULT "normal",
  `name` varchar(255) NOT NULL,
  `menu_id` char(36) NULL,
  `mnemonic` varchar(255) NOT NULL DEFAULT "",
  `shelf_life` bigint NOT NULL DEFAULT 0,
  `support_types` json NOT NULL,
  `sale_status` enum ('on_sale', 'off_sale') NOT NULL DEFAULT "on_sale",
  `sale_channels` json NOT NULL,
  `effective_date_type` enum ('daily', 'custom') NULL,
  `effective_start_time` timestamp NULL,
  `effective_end_time` timestamp NULL,
  `min_sale_quantity` bigint NULL,
  `add_sale_quantity` bigint NULL,
  `inherit_tax_rate` bool NOT NULL DEFAULT 1,
  `tax_rate_id` char(36) NULL,
  `inherit_stall` bool NOT NULL DEFAULT 1,
  `stall_id` char(36) NULL,
  `main_image` varchar(512) NOT NULL DEFAULT "",
  `detail_images` json NULL,
  `description` varchar(2000) NOT NULL DEFAULT "",
  `estimated_cost_price` decimal(10, 2) NULL,
  `delivery_cost_price` decimal(10, 2) NULL,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NULL,
  `category_id` char(36) NOT NULL,
  `unit_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `product_category_id` (`category_id`),
  INDEX `product_deleted_at` (`deleted_at`),
  INDEX `product_merchant_id` (`merchant_id`),
  INDEX `product_store_id` (`store_id`),
  INDEX `products_product_units_products` (`unit_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "product_attrs" table
CREATE TABLE `product_attrs` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `channels` json NOT NULL,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NULL,
  `product_count` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `productattr_deleted_at` (`deleted_at`),
  INDEX `productattr_merchant_id` (`merchant_id`),
  INDEX `productattr_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "product_attr_items" table
CREATE TABLE `product_attr_items` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `image` varchar(512) NOT NULL DEFAULT "",
  `base_price` decimal(10, 2) NOT NULL,
  `product_count` bigint NOT NULL DEFAULT 0,
  `attr_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `productattritem_attr_id` (`attr_id`),
  INDEX `productattritem_deleted_at` (`deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "product_attr_relations" table
CREATE TABLE `product_attr_relations` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `is_default` bool NOT NULL DEFAULT 0,
  `product_id` char(36) NOT NULL,
  `attr_id` char(36) NOT NULL,
  `attr_item_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `productattrrelation_attr_id` (`attr_id`),
  INDEX `productattrrelation_attr_item_id` (`attr_item_id`),
  INDEX `productattrrelation_deleted_at` (`deleted_at`),
  INDEX `productattrrelation_product_id` (`product_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "product_specs" table
CREATE TABLE `product_specs` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NULL,
  `product_count` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `productspec_deleted_at` (`deleted_at`),
  INDEX `productspec_merchant_id` (`merchant_id`),
  INDEX `productspec_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "product_spec_relations" table
CREATE TABLE `product_spec_relations` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `base_price` decimal(10, 2) NOT NULL,
  `member_price` decimal(10, 2) NULL,
  `packing_fee_id` char(36) NOT NULL,
  `estimated_cost_price` decimal(10, 2) NULL,
  `other_price1` decimal(10, 2) NULL,
  `other_price2` decimal(10, 2) NULL,
  `other_price3` decimal(10, 2) NULL,
  `barcode` varchar(255) NOT NULL DEFAULT "",
  `is_default` bool NOT NULL DEFAULT 0,
  `product_id` char(36) NOT NULL,
  `spec_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `productspecrelation_deleted_at` (`deleted_at`),
  INDEX `productspecrelation_product_id` (`product_id`),
  INDEX `productspecrelation_spec_id` (`spec_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "product_tags" table
CREATE TABLE `product_tags` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NULL,
  `product_count` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `producttag_deleted_at` (`deleted_at`),
  INDEX `producttag_merchant_id` (`merchant_id`),
  INDEX `producttag_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "product_units" table
CREATE TABLE `product_units` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `type` enum ('quantity', 'weight') NOT NULL,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NULL,
  `product_count` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `productunit_deleted_at` (`deleted_at`),
  INDEX `productunit_merchant_id` (`merchant_id`),
  INDEX `productunit_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "set_meal_details" table
CREATE TABLE `set_meal_details` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `quantity` bigint NOT NULL DEFAULT 1,
  `is_default` bool NOT NULL DEFAULT 0,
  `optional_product_ids` json NULL,
  `product_id` char(36) NOT NULL,
  `group_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `setmealdetail_deleted_at` (`deleted_at`),
  INDEX `setmealdetail_group_id` (`group_id`),
  INDEX `setmealdetail_product_id` (`product_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "set_meal_groups" table
CREATE TABLE `set_meal_groups` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL,
  `selection_type` enum ('fixed', 'optional') NOT NULL DEFAULT "fixed",
  `product_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `setmealgroup_deleted_at` (`deleted_at`),
  INDEX `setmealgroup_product_id` (`product_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "product_tag_relations" table
CREATE TABLE `product_tag_relations` (
  `product_id` char(36) NOT NULL,
  `tag_id` char(36) NOT NULL,
  PRIMARY KEY (`product_id`, `tag_id`),
  INDEX `product_tag_relations_tag_id` (`tag_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
