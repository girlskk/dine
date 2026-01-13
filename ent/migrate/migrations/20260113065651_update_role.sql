-- Modify "departments" table
ALTER TABLE `departments`
DROP COLUMN `enable`,
ADD COLUMN `enabled` bool NOT NULL DEFAULT 1;

-- Modify "roles" table
ALTER TABLE `roles`
DROP COLUMN `enable`,
ADD COLUMN `enabled` bool NOT NULL DEFAULT 1;
