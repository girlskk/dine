-- Modify "backend_users" table
ALTER TABLE `backend_users` DROP FOREIGN KEY `backend_users_stores_backend_users`, ADD CONSTRAINT `backend_users_stores_backend_user` FOREIGN KEY (`store_id`) REFERENCES `stores` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION;
