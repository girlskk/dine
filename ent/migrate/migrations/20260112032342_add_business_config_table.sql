-- Create "business_configs" table
CREATE TABLE `business_configs` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `source_config_id` char(36) NULL,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  `group` enum ('print') NULL,
  `name` varchar(100) NULL DEFAULT "",
  `config_type` enum ('string', 'int', 'uint', 'datetime', 'date') NULL,
  `key` varchar(100) NULL DEFAULT "",
  `value` varchar(500) NULL DEFAULT "",
  `sort` int NOT NULL DEFAULT 0,
  `tip` varchar(500) NULL DEFAULT "",
  `is_default` bool NOT NULL DEFAULT 1,
  `status` bool NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `businessconfig_deleted_at` (`deleted_at`),
  INDEX `businessconfig_group` (`group`),
  INDEX `businessconfig_key` (`key`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
