/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
/*Table structure for table `access` */

DROP TABLE IF EXISTS `access`;

CREATE TABLE `access` (
  `controllers_id` bigint(20) NOT NULL,
  `rules_id` bigint(20) NOT NULL,
  `method` varchar(50) NOT NULL,
  KEY `controllers_id` (`controllers_id`),
  KEY `rules_id` (`rules_id`),
  CONSTRAINT `access_ibfk_1` FOREIGN KEY (`controllers_id`) REFERENCES `controllers` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `access_ibfk_2` FOREIGN KEY (`rules_id`) REFERENCES `rules` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 CHECKSUM=1 DELAY_KEY_WRITE=1 ROW_FORMAT=FIXED COMMENT='права';

/*Data for the table `access` */

/*Table structure for table `controllers` */

DROP TABLE IF EXISTS `controllers`;

CREATE TABLE `controllers` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'Идентификатор',
  `nam` varchar(50) NOT NULL COMMENT 'Название контроллера',
  `controller` varchar(50) NOT NULL COMMENT 'Контроллер',
  `layout` varchar(50) DEFAULT NULL COMMENT 'Макет',
  `url` varchar(100) DEFAULT NULL COMMENT 'Урл контроллера',
  `url_redirect` varchar(100) DEFAULT NULL COMMENT 'Редирект',
  `minute` varchar(50) DEFAULT NULL COMMENT 'Минуты',
  `hour` varchar(50) DEFAULT NULL COMMENT 'Часы',
  `day` varchar(50) DEFAULT NULL COMMENT 'Дни',
  `month` varchar(50) DEFAULT NULL COMMENT 'Месяцы',
  `week` varchar(50) DEFAULT NULL COMMENT 'День недели',
  `is_execute` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Используется',
  `is_authorized` tinyint(1) NOT NULL DEFAULT '0' COMMENT 'Авторизованный',
  `typ` enum('Web','Api','Console') NOT NULL DEFAULT 'Web' COMMENT 'Тип контроллера',
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `url` (`url`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=FIXED COMMENT='контроллеры';

/*Data for the table `controllers` */

/*Table structure for table `rules` */

DROP TABLE IF EXISTS `rules`;

CREATE TABLE `rules` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `nam` varchar(50) NOT NULL,
  `is_access` tinyint(1) NOT NULL DEFAULT '0',
  `description` varchar(250) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 ROW_FORMAT=FIXED COMMENT='роли';

/*Data for the table `rules` */

insert  into `rules`(`id`,`nam`,`is_access`,`description`,`created_at`,`updated_at`,`deleted_at`) values (1,'dev',1,'Full access','2019-01-01 17:23:23','2019-01-01 17:23:23','2019-01-01 17:23:23');

/*Table structure for table `users` */

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `fio` varchar(50) NOT NULL,
  `email` varchar(50) NOT NULL,
  `password` varchar(50) NOT NULL,
  `is_access` tinyint(1) NOT NULL DEFAULT '0',
  `is_condition` tinyint(1) NOT NULL DEFAULT '1',
  `is_online` tinyint(1) NOT NULL DEFAULT '0',
  `phone` varchar(50) DEFAULT NULL,
  `img_avatar` varchar(150) DEFAULT NULL,
  `address` varchar(250) DEFAULT NULL,
  `date_online` datetime DEFAULT NULL,
  `date_reg` datetime DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 ROW_FORMAT=FIXED COMMENT='пользователи';

/*Data for the table `users` */

insert  into `users`(`id`,`fio`,`email`,`password`,`is_access`,`is_condition`,`is_online`,`phone`,`img_avatar`,`address`,`date_online`,`date_reg`,`created_at`,`updated_at`,`deleted_at`) values (1,'dev','dev@dev.dev','dev',1,0,0,NULL,NULL,NULL,NULL,'2019-01-01 17:23:23','2019-01-01 17:23:23','2019-01-01 17:23:23',NULL);

/*Table structure for table `usrules` */

DROP TABLE IF EXISTS `usrules`;

CREATE TABLE `usrules` (
  `users_id` bigint(20) NOT NULL,
  `rules_id` bigint(20) NOT NULL,
  KEY `users_id` (`users_id`),
  KEY `rules_id` (`rules_id`),
  CONSTRAINT `usrules_ibfk_1` FOREIGN KEY (`users_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `usrules_ibfk_2` FOREIGN KEY (`rules_id`) REFERENCES `rules` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

/*Data for the table `usrules` */

insert  into `usrules`(`users_id`,`rules_id`) values (1,1);

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
