-- Create "orders" table
CREATE TABLE `orders` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `merchant_id` char(36) NOT NULL,
  `store_id` char(36) NOT NULL,
  `business_date` varchar(255) NOT NULL,
  `shift_no` varchar(255) NULL,
  `order_no` varchar(255) NOT NULL,
  `order_type` enum('SALE','REFUND','PARTIAL_REFUND') NOT NULL DEFAULT "SALE",
  `origin_order_id` varchar(255) NULL,
  `refund` json NULL,
  `opened_at` timestamp NULL,
  `placed_at` timestamp NULL,
  `paid_at` timestamp NULL,
  `completed_at` timestamp NULL,
  `opened_by` varchar(255) NULL,
  `placed_by` varchar(255) NULL,
  `paid_by` varchar(255) NULL,
  `dining_mode` enum('DINE_IN','TAKEAWAY') NOT NULL,
  `order_status` enum('DRAFT','PLACED','IN_PROGRESS','READY','COMPLETED','CANCELLED','VOIDED','MERGED') NOT NULL DEFAULT "DRAFT",
  `payment_status` enum('UNPAID','PAYING','PARTIALLY_PAID','PAID','PARTIALLY_REFUNDED','REFUNDED') NOT NULL DEFAULT "UNPAID",
  `fulfillment_status` enum('NONE','IN_RESTAURANT','SERVED','PICKUP_PENDING','PICKED_UP','DELIVERING','DELIVERED') NULL,
  `table_status` enum('OPENED','TRANSFERRED','RELEASED') NULL,
  `table_id` varchar(255) NULL,
  `table_name` varchar(255) NULL,
  `table_capacity` bigint NULL,
  `guest_count` bigint NULL,
  `merged_to_order_id` varchar(255) NULL,
  `merged_at` timestamp NULL,
  `store` json NOT NULL,
  `channel` json NOT NULL,
  `pos` json NOT NULL,
  `cashier` json NOT NULL,
  `member` json NULL,
  `takeaway` json NULL,
  `cart` json NOT NULL,
  `products` json NOT NULL,
  `promotions` json NULL,
  `coupons` json NULL,
  `tax_rates` json NULL,
  `fees` json NULL,
  `payments` json NULL,
  `refunds_products` json NULL,
  `amount` json NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `order_deleted_at` (`deleted_at`),
  INDEX `order_merchant_id` (`merchant_id`),
  INDEX `order_merged_to_order_id` (`merged_to_order_id`),
  INDEX `order_origin_order_id` (`origin_order_id`),
  INDEX `order_store_id` (`store_id`),
  INDEX `order_store_id_business_date` (`store_id`, `business_date`),
  UNIQUE INDEX `order_store_id_order_no_deleted_at` (`store_id`, `order_no`, `deleted_at`),
  INDEX `order_store_id_order_status` (`store_id`, `order_status`),
  INDEX `order_store_id_payment_status` (`store_id`, `payment_status`),
  INDEX `order_table_id` (`table_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Modify "provinces" table
ALTER TABLE `provinces` ADD CONSTRAINT `provinces_countries_provinces` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "cities" table
ALTER TABLE `cities` ADD CONSTRAINT `cities_countries_cities` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `cities_provinces_cities` FOREIGN KEY (`province_id`) REFERENCES `provinces` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "districts" table
ALTER TABLE `districts` ADD CONSTRAINT `districts_cities_districts` FOREIGN KEY (`city_id`) REFERENCES `cities` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `districts_countries_districts` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `districts_provinces_districts` FOREIGN KEY (`province_id`) REFERENCES `provinces` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "merchants" table
ALTER TABLE `merchants` ADD CONSTRAINT `merchants_admin_users_merchant` FOREIGN KEY (`admin_user_id`) REFERENCES `admin_users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `merchants_cities_merchants` FOREIGN KEY (`city_id`) REFERENCES `cities` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `merchants_countries_merchants` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `merchants_districts_merchants` FOREIGN KEY (`district_id`) REFERENCES `districts` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `merchants_merchant_business_types_merchants` FOREIGN KEY (`business_type_id`) REFERENCES `merchant_business_types` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `merchants_provinces_merchants` FOREIGN KEY (`province_id`) REFERENCES `provinces` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
-- Modify "merchant_renewals" table
ALTER TABLE `merchant_renewals` ADD CONSTRAINT `merchant_renewals_merchants_merchant_renewals` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "product_attr_items" table
ALTER TABLE `product_attr_items` ADD CONSTRAINT `product_attr_items_product_attrs_items` FOREIGN KEY (`attr_id`) REFERENCES `product_attrs` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "categories" table
ALTER TABLE `categories` ADD CONSTRAINT `categories_categories_children` FOREIGN KEY (`parent_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
-- Modify "products" table
ALTER TABLE `products` ADD CONSTRAINT `products_categories_products` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `products_product_units_products` FOREIGN KEY (`unit_id`) REFERENCES `product_units` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "product_attr_relations" table
ALTER TABLE `product_attr_relations` ADD CONSTRAINT `product_attr_relations_product_attr_items_product_attrs` FOREIGN KEY (`attr_item_id`) REFERENCES `product_attr_items` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `product_attr_relations_product_attrs_product_attrs` FOREIGN KEY (`attr_id`) REFERENCES `product_attrs` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `product_attr_relations_products_product_attrs` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "product_spec_relations" table
ALTER TABLE `product_spec_relations` ADD CONSTRAINT `product_spec_relations_product_specs_product_specs` FOREIGN KEY (`spec_id`) REFERENCES `product_specs` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `product_spec_relations_products_product_specs` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "product_tag_relations" table
ALTER TABLE `product_tag_relations` ADD CONSTRAINT `product_tag_relations_product_id` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE, ADD CONSTRAINT `product_tag_relations_tag_id` FOREIGN KEY (`tag_id`) REFERENCES `product_tags` (`id`) ON UPDATE NO ACTION ON DELETE CASCADE;
-- Modify "remark_categories" table
ALTER TABLE `remark_categories` ADD CONSTRAINT `remark_categories_merchants_remark_categories` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
-- Modify "stores" table
ALTER TABLE `stores` ADD CONSTRAINT `stores_admin_users_store` FOREIGN KEY (`admin_user_id`) REFERENCES `admin_users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `stores_cities_stores` FOREIGN KEY (`city_id`) REFERENCES `cities` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `stores_countries_stores` FOREIGN KEY (`country_id`) REFERENCES `countries` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `stores_districts_stores` FOREIGN KEY (`district_id`) REFERENCES `districts` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `stores_merchant_business_types_stores` FOREIGN KEY (`business_type_id`) REFERENCES `merchant_business_types` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `stores_merchants_stores` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `stores_provinces_stores` FOREIGN KEY (`province_id`) REFERENCES `provinces` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
-- Modify "remarks" table
ALTER TABLE `remarks` ADD CONSTRAINT `remarks_merchants_remarks` FOREIGN KEY (`merchant_id`) REFERENCES `merchants` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT `remarks_remark_categories_remarks` FOREIGN KEY (`category_id`) REFERENCES `remark_categories` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `remarks_stores_remarks` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
-- Modify "set_meal_groups" table
ALTER TABLE `set_meal_groups` ADD CONSTRAINT `set_meal_groups_products_set_meal_groups` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Modify "set_meal_details" table
ALTER TABLE `set_meal_details` ADD CONSTRAINT `set_meal_details_products_set_meal_details` FOREIGN KEY (`product_id`) REFERENCES `products` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT `set_meal_details_set_meal_groups_details` FOREIGN KEY (`group_id`) REFERENCES `set_meal_groups` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
