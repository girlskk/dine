-- Modify "devices" table
ALTER TABLE `devices` MODIFY COLUMN `connect_type` enum ('inside', 'outside', 'network') NULL;

-- Modify "tax_fees" table
ALTER TABLE `tax_fees` MODIFY COLUMN `tax_fee_type` enum ('system', 'merchant', 'store') NOT NULL;
