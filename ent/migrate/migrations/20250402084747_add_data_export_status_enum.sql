-- Modify "data_exports" table
ALTER TABLE `data_exports` MODIFY COLUMN `type` enum('order_list','reconciliation_record_list','reconciliation_record_details','point_settlement_list','point_settlement_details') NOT NULL;
