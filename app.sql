CREATE DATABASE IF NOT EXISTS `app`;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
    `id` varchar(36) DEFAULT '',
    `firstname` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
    `lastname` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
    `username` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
    `email` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
    `password` varchar(128) COLLATE utf8_unicode_ci NOT NULL,
    `updated_at` datetime DEFAULT NULL,
    `created_at` datetime DEFAULT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# LOCK TABLES `user` WRITE;
# INSERT INTO `user` VALUES
#     (NULL, 'Luffy', 'Monkey', 'luffy.monkey', 'luffy.monkey@onepiece.com'),
#     (NULL, 'Sanji', 'Vinsmoke', 'sanji.vinsmoke', 'sanji.vinsmoke@onepiece.com'),
#     (NULL, 'Roronoa', 'Zoro', 'roronoa.zoro', 'roronoa.zoro@onepiece.com');
# UNLOCK TABLES;

DROP TABLE IF EXISTS `article`;
CREATE TABLE `article` (
  `id` varchar(36) DEFAULT '',
  `title` varchar(45) COLLATE utf8_unicode_ci NOT NULL,
  `content` longtext COLLATE utf8_unicode_ci NOT NULL,
  `author_id` varchar(36) DEFAULT '',
  `updated_at` datetime DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

LOCK TABLES  `article` WRITE;
INSERT INTO  `article` VALUES
  ('uniqueid1', 'Pirate King', 'Luffy will be the next Pirate King after Roger', 1, '2019-08-18 14:47:51', '2019-08-18 14:47:51'),
  ('uniqueid2', 'The Great Swordsman', 'Roronoa Zoro will be the next The Great Swordsman', 2, '2019-08-18 14:50:51', '2019-08-18 14:50:51'),
  ('uniqueid3', 'The Great Cook', 'Sanji Vinsmoke will be the next The Great Cook', 3, '2019-08-18 14:51:51', '2019-08-18 14:51:51');
UNLOCK TABLES;

DROP TABLE IF EXISTS `author`;
CREATE TABLE `author`(
  `id` varchar(36) DEFAULT '',
  `name` varchar(45) COLLATE utf8_unicode_ci DEFAULT '',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

# LOCK TABLES `author` WRITE;
# INSERT INTO `author` VALUES
#     (NULL, 'Luffy Monkey', '2019-08-18 14:47:51', '2019-08-18 14:47:51'),
#     (NULL, 'Roronoa Zoro', '2019-08-18 14:48:51', '2019-08-18 14:48:51'),
#     (NULL, 'Sanji Vinsmoke', '2019-08-18 14:49:51', '2019-08-18 14:49:51');
# UNLOCK TABLES;
