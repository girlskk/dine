-- Modify "stores" table
ALTER TABLE `stores` ADD COLUMN `type` enum('restaurant','cafeteria') NOT NULL AFTER `name`;
-- Modify "orders" table
ALTER TABLE `orders` DROP FOREIGN KEY `orders_dine_tables_orders`;
-- Modify "orders" table
ALTER TABLE `orders` MODIFY COLUMN `table_id` bigint NULL, ADD CONSTRAINT `orders_dine_tables_orders` FOREIGN KEY (`table_id`) REFERENCES `dine_tables` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
