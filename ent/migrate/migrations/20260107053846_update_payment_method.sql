-- Modify "payment_methods" table
ALTER TABLE `payment_methods` MODIFY COLUMN `payment_type` enum (
  'cash',
  'online_payment',
  'member_card',
  'custom_coupon',
  'partner_coupon',
  'bank_card'
) NOT NULL DEFAULT "cash";
