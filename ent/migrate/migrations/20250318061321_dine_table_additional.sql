-- Modify "dine_tables" table
ALTER TABLE `dine_tables` DROP FOREIGN KEY `dine_tables_orders_dinetable`;
-- Modify "orders" table
ALTER TABLE `orders` ADD INDEX `orders_dine_tables_orders` (`table_id`);
-- Modify "dine_tables" table
ALTER TABLE `dine_tables` ADD CONSTRAINT `dine_tables_orders_current_dinetable` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
-- Modify "orders" table
ALTER TABLE `orders` ADD CONSTRAINT `orders_dine_tables_orders` FOREIGN KEY (`table_id`) REFERENCES `dine_tables` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
