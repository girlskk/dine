-- Modify "business_configs" table
ALTER TABLE `business_configs` MODIFY COLUMN `config_type` enum (
  'string',
  'int',
  'uint',
  'datetime',
  'date',
  'bool'
) NULL,
MODIFY COLUMN `is_default` bool NOT NULL DEFAULT 0,
DROP INDEX `businessconfig_group`,
DROP INDEX `businessconfig_key`,
ADD UNIQUE INDEX `businessconfig_merchant_id_group_key_deleted_at` (`merchant_id`, `group`, `key`, `deleted_at`),
ADD UNIQUE INDEX `businessconfig_store_id_group_key_deleted_at` (`store_id`, `group`, `key`, `deleted_at`);
