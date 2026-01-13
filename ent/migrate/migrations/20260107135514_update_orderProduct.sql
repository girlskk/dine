-- Modify "order_products" table
ALTER TABLE `order_products`
DROP COLUMN `category_id`,
ADD COLUMN `category` json NULL,
ADD COLUMN `attr_amount` decimal(10, 4) NULL,
ADD COLUMN `gift_amount` decimal(10, 4) NULL;

-- Modify "orders" table
ALTER TABLE `orders`
ADD COLUMN `remark` varchar(255) NULL;
