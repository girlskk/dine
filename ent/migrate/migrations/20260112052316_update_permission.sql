-- Modify "merchant_renewals" table
ALTER TABLE `merchant_renewals` MODIFY COLUMN `operator_name` varchar(50) NULL DEFAULT "",
MODIFY COLUMN `operator_account` varchar(50) NULL DEFAULT "";

-- Modify "remarks" table
ALTER TABLE `remarks`
DROP COLUMN `category_id`,
ADD COLUMN `remark_scene` enum (
  'whole_order',
  'item',
  'cancel_reason',
  'discount',
  'gift',
  'rebill',
  'refund_reject'
) NOT NULL,
ADD INDEX `remark_remark_scene` (`remark_scene`);

-- Modify "role_menus" table
ALTER TABLE `role_menus`
DROP COLUMN `menu_id`,
ADD COLUMN `path` varchar(255) NOT NULL,
DROP INDEX `role_menu_unique_idx`,
ADD UNIQUE INDEX `role_path_unique_idx` (
  `role_id`,
  `merchant_id`,
  `store_id`,
  `path`,
  `role_type`,
  `deleted_at`
),
ADD INDEX `rolemenu_path` (`path`);

-- Modify "router_menus" table
ALTER TABLE `router_menus` ADD UNIQUE INDEX `routermenu_path_deleted_at` (`path`, `deleted_at`);
