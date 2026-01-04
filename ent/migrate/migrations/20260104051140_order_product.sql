-- Modify "orders" table
ALTER TABLE `orders` MODIFY COLUMN `dining_mode` enum('DINE_IN') NOT NULL DEFAULT "DINE_IN";
