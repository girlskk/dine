-- Modify "order_products" table
ALTER TABLE `order_products`
DROP COLUMN `menu_id`,
DROP COLUMN `support_types`,
DROP COLUMN `sale_status`,
DROP COLUMN `sale_channels`,
MODIFY COLUMN `refunded_by` char(36) NULL,
DROP COLUMN `estimated_cost_price`,
DROP COLUMN `delivery_cost_price`,
DROP COLUMN `set_meal_groups`,
ADD COLUMN `is_gift` bool NOT NULL DEFAULT 0,
ADD COLUMN `gift_qty` bigint NOT NULL DEFAULT 0,
ADD COLUMN `price` decimal(10, 4) NULL,
ADD COLUMN `groups` json NULL;

-- Modify "orders" table
ALTER TABLE `orders`
DROP COLUMN `refund`,
MODIFY COLUMN `placed_by` char(36) NULL,
MODIFY COLUMN `table_id` char(36) NULL,
ADD COLUMN `placed_by_name` varchar(255) NULL;
