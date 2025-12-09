-- Modify "payment_callbacks" table
ALTER TABLE `payment_callbacks` MODIFY COLUMN `provider` enum('zxh','huifu') NOT NULL;
-- Modify "payments" table
ALTER TABLE `payments` MODIFY COLUMN `provider` enum('zxh','huifu') NOT NULL;
