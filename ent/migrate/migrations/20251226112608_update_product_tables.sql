-- Modify "categories" table
ALTER TABLE `categories` MODIFY COLUMN `store_id` char(36) NOT NULL,
ADD UNIQUE INDEX `category_merchant_id_store_id_parent_id_name_deleted_at` (
  `merchant_id`,
  `store_id`,
  `parent_id`,
  `name`,
  `deleted_at`
);

-- Modify "product_attr_items" table
ALTER TABLE `product_attr_items` ADD UNIQUE INDEX `productattritem_attr_id_name_deleted_at` (`attr_id`, `name`, `deleted_at`);

-- Modify "product_attrs" table
ALTER TABLE `product_attrs` MODIFY COLUMN `store_id` char(36) NOT NULL,
ADD UNIQUE INDEX `productattr_merchant_id_store_id_name_deleted_at` (`merchant_id`, `store_id`, `name`, `deleted_at`);

-- Modify "product_specs" table
ALTER TABLE `product_specs` MODIFY COLUMN `store_id` char(36) NOT NULL,
ADD UNIQUE INDEX `productspec_merchant_id_store_id_name_deleted_at` (`merchant_id`, `store_id`, `name`, `deleted_at`);

-- Modify "product_tags" table
ALTER TABLE `product_tags` MODIFY COLUMN `store_id` char(36) NOT NULL,
ADD UNIQUE INDEX `producttag_merchant_id_store_id_name_deleted_at` (`merchant_id`, `store_id`, `name`, `deleted_at`);

-- Modify "product_units" table
ALTER TABLE `product_units` MODIFY COLUMN `store_id` char(36) NOT NULL,
ADD UNIQUE INDEX `productunit_merchant_id_store_id_name_deleted_at` (`merchant_id`, `store_id`, `name`, `deleted_at`);

-- Modify "products" table
ALTER TABLE `products` MODIFY COLUMN `estimated_cost_price` decimal(19, 4) NULL,
MODIFY COLUMN `delivery_cost_price` decimal(19, 4) NULL,
MODIFY COLUMN `store_id` char(36) NOT NULL,
ADD UNIQUE INDEX `product_merchant_id_store_id_name_deleted_at` (`merchant_id`, `store_id`, `name`, `deleted_at`);

-- Create "menus" table
CREATE TABLE `menus` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `merchant_id` char(36) NOT NULL,
  `name` varchar(255) NOT NULL,
  `distribution_rule` enum ('override', 'keep') NOT NULL DEFAULT "override",
  `store_count` bigint NOT NULL DEFAULT 0,
  `item_count` bigint NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `menu_deleted_at` (`deleted_at`),
  INDEX `menu_merchant_id` (`merchant_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "menu_items" table
CREATE TABLE `menu_items` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `sale_rule` enum ('keep_brand_status', 'keep_store_status') NOT NULL DEFAULT "keep_brand_status",
  `base_price` decimal(19, 4) NULL,
  `member_price` decimal(19, 4) NULL,
  `menu_id` char(36) NOT NULL,
  `product_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `menuitem_deleted_at` (`deleted_at`),
  INDEX `menuitem_menu_id` (`menu_id`),
  INDEX `menuitem_product_id` (`product_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "menu_store_relations" table
CREATE TABLE `menu_store_relations` (
  `menu_id` char(36) NOT NULL,
  `store_id` char(36) NOT NULL,
  PRIMARY KEY (`menu_id`, `store_id`),
  INDEX `menu_store_relations_store_id` (`store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
