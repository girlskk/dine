-- Modify "order_finance_logs" table
ALTER TABLE `order_finance_logs` MODIFY COLUMN `channel` enum('cash','wechat','alipay','point','point_wallet') NOT NULL;
-- Modify "orders" table
ALTER TABLE `orders` ADD COLUMN `points_wallet_paid` decimal(10,2) NOT NULL AFTER `points_refunded`, ADD COLUMN `points_wallet_refunded` decimal(10,2) NOT NULL AFTER `points_wallet_paid`;
-- Modify "payment_callbacks" table
ALTER TABLE `payment_callbacks` MODIFY COLUMN `provider` enum('zxh','zxh_wallet','huifu') NOT NULL;
-- Modify "payments" table
ALTER TABLE `payments` MODIFY COLUMN `provider` enum('zxh','zxh_wallet','huifu') NOT NULL, MODIFY COLUMN `channel` enum('wxpay','alipay','point','point_wallet') NOT NULL;
-- Modify "reconciliation_records" table
ALTER TABLE `reconciliation_records` MODIFY COLUMN `channel` enum('cash','wechat','alipay','point','point_wallet') NOT NULL;
