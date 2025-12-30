-- Modify "menu_items" table
ALTER TABLE `menu_items` ADD CONSTRAINT `menu_items_menus_items` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `menu_items_products_menu_items` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "menu_store_relations" table
ALTER TABLE `menu_store_relations` ADD CONSTRAINT `menu_store_relations_menu_id` FOREIGN KEY (`menu_id`) REFERENCES `menus` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT `menu_store_relations_store_id` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "orders" table
ALTER TABLE `orders` DROP COLUMN `origin_order_id`, DROP COLUMN `opened_at`, DROP COLUMN `opened_by`, DROP COLUMN `paid_by`, MODIFY COLUMN `dining_mode` enum('DINE_IN') NOT NULL, MODIFY COLUMN `order_status` enum('PLACED','COMPLETED','CANCELLED') NOT NULL DEFAULT "PLACED", MODIFY COLUMN `payment_status` enum('UNPAID','PAYING','PAID','REFUNDED') NOT NULL DEFAULT "UNPAID", DROP COLUMN `fulfillment_status`, DROP COLUMN `table_status`, DROP COLUMN `table_capacity`, DROP COLUMN `merged_to_order_id`, DROP COLUMN `merged_at`, MODIFY COLUMN `channel` enum('POS') NOT NULL DEFAULT "POS", DROP COLUMN `member`, DROP COLUMN `takeaway`, DROP COLUMN `cart`, DROP COLUMN `products`, DROP COLUMN `promotions`, DROP COLUMN `coupons`, DROP COLUMN `refunds_products`, DROP INDEX `order_store_id_business_date`, DROP INDEX `order_store_id_order_status`, DROP INDEX `order_store_id_payment_status`, ADD INDEX `order_order_status` (`order_status`), ADD INDEX `order_payment_status` (`payment_status`);
-- Create "order_products" table
CREATE TABLE `order_products` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `order_item_id` varchar(255) NOT NULL,
  `index` bigint NOT NULL DEFAULT 0,
  `product_id` char(36) NOT NULL,
  `product_name` varchar(255) NOT NULL,
  `product_type` enum('normal','set_meal') NOT NULL DEFAULT "normal",
  `category_id` char(36) NULL,
  `menu_id` char(36) NULL,
  `unit_id` char(36) NULL,
  `support_types` json NULL,
  `sale_status` enum('on_sale','off_sale') NULL,
  `sale_channels` json NULL,
  `main_image` varchar(512) NOT NULL DEFAULT "",
  `description` varchar(2000) NOT NULL DEFAULT "",
  `qty` bigint NOT NULL DEFAULT 1,
  `subtotal` decimal(10,4) NULL,
  `discount_amount` decimal(10,4) NULL,
  `amount_before_tax` decimal(10,4) NULL,
  `tax_rate` decimal(10,4) NULL,
  `tax` decimal(10,4) NULL,
  `amount_after_tax` decimal(10,4) NULL,
  `total` decimal(10,4) NULL,
  `promotion_discount` decimal(10,4) NULL,
  `void_qty` bigint NOT NULL DEFAULT 0,
  `void_amount` decimal(10,4) NULL,
  `refund_reason` varchar(255) NULL,
  `refunded_by` varchar(255) NULL,
  `refunded_at` timestamp NULL,
  `note` varchar(255) NULL,
  `estimated_cost_price` decimal(10,4) NULL,
  `delivery_cost_price` decimal(10,4) NULL,
  `set_meal_groups` json NULL,
  `spec_relations` json NULL,
  `attr_relations` json NULL,
  `order_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `orderproduct_deleted_at` (`deleted_at`),
  INDEX `orderproduct_order_id` (`order_id`),
  UNIQUE INDEX `orderproduct_order_id_order_item_id` (`order_id`, `order_item_id`),
  INDEX `orderproduct_product_id` (`product_id`),
  CONSTRAINT `order_products_orders_order_products` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
