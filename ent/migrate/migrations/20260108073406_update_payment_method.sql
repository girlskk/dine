-- Modify "orders" table
ALTER TABLE `orders`
ADD COLUMN `remark` varchar(255) NULL;

-- Modify "payment_methods" table
ALTER TABLE `payment_methods` MODIFY COLUMN `store_id` char(36) NULL,
MODIFY COLUMN `payment_type` enum (
  'cash',
  'online_payment',
  'member_card',
  'custom_coupon',
  'partner_coupon',
  'bank_card'
) NOT NULL DEFAULT "cash",
MODIFY COLUMN `display_channels` json NULL,
MODIFY COLUMN `source` enum ('brand', 'store', 'system') NULL,
DROP COLUMN `store_ids`,
ADD COLUMN `source_payment_method_id` char(36) NULL;
