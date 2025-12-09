-- Modify "dine_tables" table
ALTER TABLE `dine_tables` ADD COLUMN `order_id` bigint NULL, ADD INDEX `dine_tables_orders_dinetables` (`order_id`), ADD CONSTRAINT `dine_tables_orders_dinetables` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
