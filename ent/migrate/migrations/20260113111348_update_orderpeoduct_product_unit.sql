-- Modify "order_products" table
ALTER TABLE `order_products`
DROP COLUMN `unit_id`,
ADD COLUMN `product_unit` json NULL;

-- Modify "refund_order_products" table
ALTER TABLE `refund_order_products`
DROP COLUMN `unit_id`,
ADD COLUMN `product_unit` json NULL;
