-- License: MIT
-- Authors:
-- 		- Josep Bigorra (averageflow)
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `logs`;
CREATE TABLE `logs` (
  `uid` varchar(36) NOT NULL,
  `application` varchar(255) NOT NULL,
  `error` longtext,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `requests`;
CREATE TABLE `requests` (
  `uid` varchar(36) NOT NULL,
  `application` varchar(255) NOT NULL,
  `client_ip` varchar(255) DEFAULT NULL,
  `method` varchar(255) DEFAULT NULL,
  `path` varchar(255) DEFAULT NULL,
  `url` longtext DEFAULT NULL,
  `host` varchar(255) DEFAULT NULL,
  `headers` text,
  `body` longtext,
  `referrer` varchar(255) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `responses`;
CREATE TABLE `responses` (
  `uid` varchar(36) NOT NULL,
  `request_uid` varchar(36) NOT NULL,
  `application` varchar(255) NOT NULL,
  `client_ip` varchar(255) DEFAULT NULL,
  `status` int(255) DEFAULT NULL,
  `body` longtext,
  `path` varchar(255) DEFAULT NULL,
  `headers` text,
  `size` bigint(20) DEFAULT NULL,
  `time` bigint(20) NOT NULL,
  PRIMARY KEY (`uid`) USING BTREE,
  KEY `request_uid_foreign` (`request_uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;