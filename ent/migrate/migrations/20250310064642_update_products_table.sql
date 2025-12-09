-- Modify "products" table
ALTER TABLE `products` DROP FOREIGN KEY `products_units_products`;
-- Modify "products" table
ALTER TABLE `products` MODIFY COLUMN `unit_id` bigint NULL, ADD CONSTRAINT `products_units_products` FOREIGN KEY (`unit_id`) REFERENCES `units` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL;
