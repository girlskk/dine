-- Modify "orders" table
ALTER TABLE `orders` MODIFY COLUMN `channel` enum ('POS', 'H5', 'APP') NOT NULL DEFAULT "POS",
ADD COLUMN `operation_logs` json NULL;

-- Modify "refund_orders" table
ALTER TABLE `refund_orders` MODIFY COLUMN `channel` enum ('POS', 'H5', 'APP') NOT NULL DEFAULT "POS";
