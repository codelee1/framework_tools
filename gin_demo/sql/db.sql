
-- 导出 db 的数据库结构
CREATE DATABASE IF NOT EXISTS `demo` /*!40100 DEFAULT CHARACTER SET utf8mb4 */;
USE `demo`;

-- 导出  表 db.control_notice 结构
CREATE TABLE IF NOT EXISTS `control_notice` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(50) DEFAULT NULL COMMENT '创建时间',
  `content` text,
  `url` varchar(150) DEFAULT NULL COMMENT '超链接',
  `create_time` bigint(20) DEFAULT NULL COMMENT '创建时间',
  `update_time` bigint(20) DEFAULT NULL COMMENT '更新时间',
  `del` bit(1) NOT NULL DEFAULT b'0' COMMENT '0为未删除，1已删除',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
