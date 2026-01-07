-- Modify "admin_users" table
ALTER TABLE `admin_users`
ADD COLUMN `code` varchar(255) NOT NULL,
ADD COLUMN `department_id` char(36) NOT NULL,
ADD INDEX `admin_users_departments_admin_users` (`department_id`);

-- Modify "backend_users" table
ALTER TABLE `backend_users`
ADD COLUMN `code` varchar(255) NOT NULL,
ADD COLUMN `department_id` char(36) NOT NULL,
ADD INDEX `backend_users_departments_backend_users` (`department_id`);

-- Modify "devices" table
ALTER TABLE `devices`
ADD COLUMN `connect_type` enum ('inside', 'outside') NULL;

-- Modify "merchant_business_types" table
ALTER TABLE `merchant_business_types`
ADD COLUMN `merchant_id` char(36) NOT NULL,
ADD UNIQUE INDEX `merchantbusinesstype_type_code_deleted_at` (`type_code`, `deleted_at`);

-- Modify "merchants" table
ALTER TABLE `merchants`
DROP COLUMN `business_type_id`,
ADD COLUMN `business_type_code` varchar(255) NOT NULL;

-- Modify "roles" table
ALTER TABLE `roles`
ADD COLUMN `login_channels` json NULL;

-- Modify "router_menus" table
ALTER TABLE `router_menus` MODIFY COLUMN `path` varchar(255) NOT NULL,
ADD COLUMN `layer` bigint NOT NULL DEFAULT 1;

-- Modify "store_users" table
ALTER TABLE `store_users`
ADD COLUMN `code` varchar(255) NOT NULL,
ADD COLUMN `department_id` char(36) NOT NULL,
ADD INDEX `store_users_departments_store_users` (`department_id`);

-- Modify "stores" table
ALTER TABLE `stores`
DROP COLUMN `business_type_id`,
ADD COLUMN `business_type_code` varchar(255) NOT NULL;
