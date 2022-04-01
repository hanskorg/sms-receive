/*
 Navicat Premium Data Transfer

 Target Server Type    : MySQL
 Target Server Version : 80017
 File Encoding         : 65001

 Date: 12/03/2021 00:11:13
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `id` int(11) NOT NULL,
  `num_id` int(11) NOT NULL COMMENT '手机号ID',
  `from` int(20) NOT NULL COMMENT '发信人',
  `is_del` tinyint(2) NOT NULL COMMENT '显示否',
  `message` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '短信',
  `created_at` datetime NOT NULL COMMENT '收信时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for numbers
-- ----------------------------
DROP TABLE IF EXISTS `numbers`;
CREATE TABLE `numbers` (
  `id` int(5) NOT NULL AUTO_INCREMENT,
  `number` int(20) unsigned NOT NULL COMMENT '手机号码',
  `zone` tinyint(4) NOT NULL COMMENT '国际区号',
  `valid` tinyint(1) NOT NULL COMMENT '是否有效',
  `free` tinyint(1) NOT NULL COMMENT '是否为免费号码',
  `carrier` tinyint(4) NOT NULL COMMENT '运行商ID',
  `create_at` datetime NOT NULL,
  `release_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `z_n` (`zone`,`number`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
