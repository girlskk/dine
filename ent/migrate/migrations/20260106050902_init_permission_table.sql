-- Modify "admin_users" table
ALTER TABLE `admin_users`
ADD COLUMN `real_name` varchar(100) NOT NULL,
ADD COLUMN `gender` enum ('male', 'female', 'other', 'unknown') NOT NULL,
ADD COLUMN `email` varchar(100) NULL,
ADD COLUMN `phone_number` varchar(20) NULL,
ADD COLUMN `enabled` bool NOT NULL DEFAULT 0,
ADD COLUMN `is_superadmin` bool NOT NULL DEFAULT 0;

-- Modify "backend_users" table
ALTER TABLE `backend_users`
DROP INDEX `backend_users_merchants_backend_users`;

-- Modify "backend_users" table
ALTER TABLE `backend_users`
DROP COLUMN `merchant_backend_users`,
ADD COLUMN `real_name` varchar(100) NOT NULL,
ADD COLUMN `gender` enum ('male', 'female', 'other', 'unknown') NOT NULL,
ADD COLUMN `email` varchar(100) NULL,
ADD COLUMN `phone_number` varchar(20) NULL,
ADD COLUMN `enabled` bool NOT NULL DEFAULT 0,
ADD COLUMN `is_superadmin` bool NOT NULL DEFAULT 0,
ADD INDEX `backend_users_merchants_backend_users` (`merchant_id`);

-- Modify "merchants" table
ALTER TABLE `merchants` MODIFY COLUMN `merchant_code` varchar(100) NULL DEFAULT "",
MODIFY COLUMN `merchant_name` varchar(100) NOT NULL DEFAULT "",
MODIFY COLUMN `merchant_short_name` varchar(100) NULL DEFAULT "",
MODIFY COLUMN `brand_name` varchar(100) NULL DEFAULT "";

-- Modify "roles" table
ALTER TABLE `roles`
ADD COLUMN `data_scope` enum (
  'all',
  'merchant',
  'store',
  'department',
  'self',
  'custom'
) NULL DEFAULT "all";

-- Modify "store_users" table
ALTER TABLE `store_users`
ADD COLUMN `real_name` varchar(100) NOT NULL,
ADD COLUMN `gender` enum ('male', 'female', 'other', 'unknown') NOT NULL,
ADD COLUMN `email` varchar(100) NULL,
ADD COLUMN `phone_number` varchar(20) NULL,
ADD COLUMN `enabled` bool NOT NULL DEFAULT 0,
ADD COLUMN `is_superadmin` bool NOT NULL DEFAULT 0;

-- Modify "stores" table
ALTER TABLE `stores` MODIFY COLUMN `store_name` varchar(100) NOT NULL DEFAULT "",
MODIFY COLUMN `store_short_name` varchar(100) NULL DEFAULT "",
MODIFY COLUMN `contact_name` varchar(100) NULL DEFAULT "",
MODIFY COLUMN `address` varchar(255) NULL DEFAULT "";

-- Create "permissions" table
CREATE TABLE `permissions` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `perm_code` varchar(150) NOT NULL,
  `name` varchar(150) NOT NULL,
  `method` varchar(10) NOT NULL,
  `path` varchar(255) NOT NULL,
  `enabled` bool NOT NULL DEFAULT 1,
  `menu_id` char(36) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `permission_deleted_at` (`deleted_at`),
  INDEX `permission_menu_id` (`menu_id`),
  UNIQUE INDEX `permission_method_path_deleted_at` (`method`, `path`, `deleted_at`),
  UNIQUE INDEX `permission_perm_code_deleted_at` (`perm_code`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "role_menus" table
CREATE TABLE `role_menus` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `role_type` enum ('admin', 'backend', 'store') NOT NULL,
  `role_id` char(36) NOT NULL,
  `menu_id` char(36) NOT NULL,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `role_menu_unique_idx` (
    `role_id`,
    `merchant_id`,
    `store_id`,
    `menu_id`,
    `role_type`,
    `deleted_at`
  ),
  INDEX `rolemenu_deleted_at` (`deleted_at`),
  INDEX `rolemenu_menu_id` (`menu_id`),
  INDEX `rolemenu_merchant_id_store_id` (`merchant_id`, `store_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "role_permissions" table
CREATE TABLE `role_permissions` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `role_type` enum ('admin', 'backend', 'store') NOT NULL,
  `role_id` char(36) NOT NULL,
  `permission_id` char(36) NOT NULL,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  INDEX `rolepermission_deleted_at` (`deleted_at`),
  INDEX `rolepermission_merchant_id_store_id` (`merchant_id`, `store_id`),
  INDEX `rolepermission_permission_id` (`permission_id`),
  UNIQUE INDEX `rolepermission_role_id_merchant_54b13b56d14df913f5c757e7528f92ab` (
    `role_id`,
    `merchant_id`,
    `store_id`,
    `permission_id`,
    `role_type`,
    `deleted_at`
  )
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "router_menus" table
CREATE TABLE `router_menus` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `user_type` enum ('admin', 'backend', 'store') NOT NULL,
  `parent_id` char(36) NULL,
  `name` varchar(100) NOT NULL,
  `path` varchar(255) NULL,
  `component` varchar(255) NULL,
  `icon` varchar(500) NULL,
  `sort` bigint NOT NULL DEFAULT 0,
  `enabled` bool NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `routermenu_deleted_at` (`deleted_at`),
  UNIQUE INDEX `routermenu_parent_id_name_deleted_at` (`parent_id`, `name`, `deleted_at`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;

-- Create "user_roles" table
CREATE TABLE `user_roles` (
  `id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` bigint NOT NULL DEFAULT 0,
  `user_type` enum ('admin', 'backend', 'store') NOT NULL,
  `user_id` char(36) NOT NULL,
  `role_id` char(36) NOT NULL,
  `merchant_id` char(36) NULL,
  `store_id` char(36) NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `role_user_unique_idx` (
    `role_id`,
    `user_type`,
    `user_id`,
    `merchant_id`,
    `store_id`,
    `deleted_at`
  ),
  INDEX `userrole_deleted_at` (`deleted_at`),
  INDEX `userrole_merchant_id_store_id` (`merchant_id`, `store_id`),
  INDEX `userrole_user_type_user_id` (`user_type`, `user_id`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
