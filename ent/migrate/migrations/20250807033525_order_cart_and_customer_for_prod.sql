-- Modify "data_exports" table
ALTER TABLE
  `data_exports`
MODIFY
  COLUMN `operator_type` enum('frontend', 'backend', 'admin', 'system', 'customer') NOT NULL;

-- Modify "order_finance_logs" table
ALTER TABLE
  `order_finance_logs`
MODIFY
  COLUMN `creator_type` enum('frontend', 'backend', 'admin', 'system', 'customer') NOT NULL;

-- Modify "order_logs" table
ALTER TABLE
  `order_logs`
MODIFY
  COLUMN `operator_type` enum('frontend', 'backend', 'admin', 'system', 'customer') NOT NULL;

-- Modify "orders" table
ALTER TABLE
  `orders`
MODIFY
  COLUMN `source` enum('offline', 'mini_program') NOT NULL,
ADD
  COLUMN `creator_type` enum('frontend', 'backend', 'admin', 'system', 'customer') NOT NULL DEFAULT "frontend"
AFTER
  `creator_name`;

-- Modify "payments" table
ALTER TABLE
  `payments`
MODIFY
  COLUMN `creator_type` enum('frontend', 'backend', 'admin', 'system', 'customer') NOT NULL;

-- Create "customers" table
CREATE TABLE `customers` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `nickname` varchar(100) NOT NULL,
  `phone` varchar(20) NOT NULL,
  `avatar` varchar(255) NULL,
  `gender` enum('male', 'female', 'unknown') NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `customer_deleted_at` (`deleted_at`),
  UNIQUE INDEX `customer_phone_deleted_at` (`phone`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "order_carts" table
CREATE TABLE `order_carts` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `table_id` bigint NOT NULL,
  `quantity` decimal(10, 2) NOT NULL,
  `attr_id` bigint NULL,
  `product_id` bigint NOT NULL,
  `product_spec_id` bigint NULL,
  `recipe_id` bigint NULL,
  PRIMARY KEY (`id`),
  INDEX `order_carts_attrs_order_cart_items` (`attr_id`),
  INDEX `order_carts_product_specs_order_cart_items` (`product_spec_id`),
  INDEX `order_carts_products_order_cart_items` (`product_id`),
  INDEX `order_carts_recipes_order_cart_items` (`recipe_id`),
  INDEX `ordercart_deleted_at` (`deleted_at`),
  INDEX `ordercart_table_id` (`table_id`),
  CONSTRAINT `order_carts_attrs_order_cart_items` FOREIGN KEY (`attr_id`) REFERENCES `attrs` (`id`) ON UPDATE NO ACTION ON DELETE
  SET
    NULL,
    CONSTRAINT `order_carts_product_specs_order_cart_items` FOREIGN KEY (`product_spec_id`) REFERENCES `product_specs` (`id`) ON UPDATE NO ACTION ON DELETE
  SET
    NULL,
    CONSTRAINT `order_carts_products_order_cart_items` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
    CONSTRAINT `order_carts_recipes_order_cart_items` FOREIGN KEY (`recipe_id`) REFERENCES `recipes` (`id`) ON UPDATE NO ACTION ON DELETE
  SET
    NULL
) CHARSET utf8mb4 COLLATE utf8mb4_bin;