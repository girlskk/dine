-- Modify "store_withdraws" table
ALTER TABLE `store_withdraws` MODIFY COLUMN `account_type` enum('public','private') NOT NULL;
