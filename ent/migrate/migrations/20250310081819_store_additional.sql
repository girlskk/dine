-- Modify "stores" table
ALTER TABLE `stores` ADD COLUMN `cooperation_type` enum('join') NOT NULL, ADD COLUMN `need_audit` bool NOT NULL, ADD COLUMN `enabled` bool NOT NULL;
