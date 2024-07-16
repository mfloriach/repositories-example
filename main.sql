CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4;

INSERT INTO users (id,name,created_at,updated_at) VALUES
	 (1,'first','2024-04-10 23:00:02','2024-04-10 23:00:02'),
	 (2,'second','2024-04-11 23:00:02','2024-04-11 23:00:02'),
	 (3,'third','2024-04-12 23:00:20','2024-04-12 23:00:22'),
	 (4,'forth','2024-04-13 23:00:20','2024-04-13 23:00:22'),
	 (5,'five','2024-04-14 23:00:20','2024-04-14 23:00:20'),
	 (6,'six','2024-04-15 23:00:20','2024-04-15 23:00:20');
