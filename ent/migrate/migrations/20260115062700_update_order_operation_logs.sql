-- Modify "order_products" table
ALTER TABLE `order_products`
DROP COLUMN `unit_id`,
ADD COLUMN `product_unit` json NULL;

-- Modify "orders" table
ALTER TABLE `orders` MODIFY COLUMN `channel` enum ('POS', 'H5', 'APP') NOT NULL DEFAULT "POS",
ADD COLUMN `operation_logs` json NULL;

-- Modify "refund_order_products" table
ALTER TABLE `refund_order_products`
DROP COLUMN `unit_id`,
ADD COLUMN `product_unit` json NULL;

-- Modify "refund_orders" table
ALTER TABLE `refund_orders` MODIFY COLUMN `channel` enum ('POS', 'H5', 'APP') NOT NULL DEFAULT "POS";
