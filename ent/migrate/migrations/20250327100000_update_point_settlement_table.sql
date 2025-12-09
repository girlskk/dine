-- Modify "point_settlements" table
ALTER TABLE `point_settlements` ADD COLUMN `approved_at` timestamp NULL, ADD COLUMN `approver_id` bigint NULL;
