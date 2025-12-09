-- Modify "dine_tables" table
ALTER TABLE `dine_tables` DROP INDEX `dine_tables_orders_dinetables`, ADD UNIQUE INDEX `order_id` (`order_id`), DROP FOREIGN KEY `dine_tables_orders_dinetables`, ADD CONSTRAINT `dine_tables_orders_dinetable` FOREIGN KEY (`order_id`) REFERENCES `orders` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
