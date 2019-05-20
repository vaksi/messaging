CREATE TABLE `messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `to_user_id` int(11) NOT NULL,
  `message_text` text NOT NULL,
  `subject` varchar(50) NOT NULL DEFAULT '',
  `status` tinyint(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp(),
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `messages` (`id`, `user_id`, `to_user_id`, `message_text`, `subject`, `status`, `created_at`, `updated_at`)
VALUES
	(82,1,2,'blablabla','okem',4,'2019-05-20 06:21:30','2019-05-20 06:21:39'),
	(83,1,2,'blablabla','okem',4,'2019-05-20 06:21:31','2019-05-20 06:21:42'),
	(84,1,2,'blablabla','okem',4,'2019-05-20 06:21:33','2019-05-20 06:21:45'),
	(85,1,2,'blablabla','okem',3,'2019-05-20 12:46:19','2019-05-20 12:48:20');