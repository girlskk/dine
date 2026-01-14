-- Modify "business_configs" table
ALTER TABLE `business_configs` MODIFY COLUMN `store_id` char(36) NOT NULL DEFAULT "",
MODIFY COLUMN `group` enum ('order', 'payment', 'kitchen', 'refund', 'print') NULL,
MODIFY COLUMN `config_type` enum (
  'string',
  'int',
  'uint',
  'datetime',
  'date',
  'bool'
) NULL,
MODIFY COLUMN `is_default` bool NOT NULL DEFAULT 0,
ADD COLUMN `modify_status` bool NOT NULL DEFAULT 1,
DROP INDEX `businessconfig_group`,
DROP INDEX `businessconfig_key`,
ADD UNIQUE INDEX `businessconfig_merchant_id_store_id_group_key_deleted_at` (
  `merchant_id`,
  `store_id`,
  `group`,
  `key`,
  `deleted_at`
);

-- Modify "tax_fees" table
ALTER TABLE `tax_fees` MODIFY COLUMN `tax_fee_type` enum ('system', 'merchant', 'store') NOT NULL;
