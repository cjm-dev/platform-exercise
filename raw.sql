CREATE TABLE `users` (`id` int(11) NOT NULL AUTO_INCREMENT, `name` varchar(255) DEFAULT NULL,
                      `email` varchar(255) NOT NULL, `password` varchar(255) NOT NULL,
                      `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
                      `updated_at` datetime DEFAULT current_timestamp() ON UPDATE current_timestamp(),
                      PRIMARY KEY (`id`), UNIQUE KEY `email` (`email`));

CREATE TABLE `authentication` (`id` int(11) NOT NULL AUTO_INCREMENT, `user_id` int(11) NOT NULL,
                               `auth_uuid` varchar(255) NOT NULL,
                                PRIMARY KEY (`Id`));