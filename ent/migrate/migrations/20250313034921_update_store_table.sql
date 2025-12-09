-- Modify "stores" table
ALTER TABLE `stores` MODIFY COLUMN `need_audit` bool NOT NULL, MODIFY COLUMN `enabled` bool NOT NULL, ADD COLUMN `point_settlement_rate` decimal(5,4) NOT NULL, ADD COLUMN `point_withdrawal_rate` decimal(5,4) NOT NULL;
-- Modify "backend_users" table
ALTER TABLE `backend_users` RENAME COLUMN `store_backend_user` TO `store_id`, DROP INDEX `store_backend_user`, ADD UNIQUE INDEX `store_id` (`store_id`), DROP FOREIGN KEY `backend_users_stores_backend_user`, ADD CONSTRAINT `backend_users_stores_backend_users` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create "store_finances" table
CREATE TABLE `store_finances` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `bank_account` varchar(255) NULL,
  `bank_card_name` varchar(255) NULL,
  `bank_name` varchar(255) NULL,
  `branch_name` varchar(255) NULL,
  `public_account` varchar(255) NULL,
  `company_name` varchar(255) NULL,
  `public_bank_name` varchar(255) NULL,
  `public_branch_name` varchar(255) NULL,
  `credit_code` varchar(255) NULL,
  `store_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `store_id` (`store_id`),
  CONSTRAINT `store_finances_stores_store_finance` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "store_infos" table
CREATE TABLE `store_infos` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `city` varchar(255) NULL,
  `address` varchar(255) NULL,
  `contact_name` varchar(255) NULL,
  `contact_phone` varchar(255) NULL,
  `images` json NULL,
  `store_id` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `store_id` (`store_id`),
  CONSTRAINT `store_infos_stores_store_info` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
