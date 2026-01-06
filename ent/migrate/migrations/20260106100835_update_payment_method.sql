-- Modify "payment_methods" table
ALTER TABLE `payment_methods` MODIFY COLUMN `store_id` char(36) NULL,
MODIFY COLUMN `display_channels` json NULL,
MODIFY COLUMN `source` enum ('brand', 'store', 'system') NULL,
DROP COLUMN `store_ids`,
ADD COLUMN `source_payment_method_id` char(36) NULL;
