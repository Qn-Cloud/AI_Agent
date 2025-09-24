-- AI 角色扮演语音交互产品数据库初始化脚本
SET NAMES utf8mb4;

-- 用户�?
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL UNIQUE,
  `email` varchar(100) NOT NULL UNIQUE,
  `password_hash` varchar(255) NOT NULL,
  `avatar` varchar(255) DEFAULT '',
  `nickname` varchar(50) DEFAULT '',
  `bio` text,
  `status` tinyint DEFAULT '1' COMMENT '1:正常 0:禁用',
  `last_login_at` timestamp NULL,
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户�?;

-- 角色�?
CREATE TABLE IF NOT EXISTS `characters` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `avatar` varchar(255) DEFAULT '',
  `description` text,
  `prompt` text NOT NULL,
  `tags` json,
  `category` varchar(50) DEFAULT '',
  `rating` decimal(3,2) DEFAULT '0.00',
  `rating_count` int DEFAULT '0',
  `favorite_count` int DEFAULT '0',
  `chat_count` int DEFAULT '0',
  `status` tinyint DEFAULT '1' COMMENT '1:启用 0:禁用',
  `is_public` tinyint DEFAULT '1' COMMENT '1:公开 0:私有',
  `creator_id` bigint DEFAULT '0',
  `created_at` timestamp DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_name` (`name`),
  KEY `idx_category` (`category`),
  KEY `idx_creator` (`creator_id`),
  KEY `idx_status` (`status`),
  KEY `idx_public` (`is_public`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色�?;

-- 插入初始数据
INSERT IGNORE INTO `characters` (`id`, `name`, `avatar`, `description`, `prompt`, `tags`, `category`, `creator_id`) VALUES
(1, '哈利·波特', '/images/avatars/harry-potter.jpg', '霍格沃茨魔法学校的学生，拥有闪电疤痕的男孩巫师。勇敢、善良，擅长魁地奇运动�?, '你是哈利·波特，霍格沃茨的学生。你勇敢善良，有着丰富的魔法世界冒险经历。请用哈利的语气和视角来回答问题�?, '["魔法", "勇敢", "冒险", "友谊"]', '经典IP', 0),
(2, '苏格拉底', '/images/avatars/socrates.jpg', '古希腊哲学家，以苏格拉底式问答法闻名。追求智慧与真理，善于启发式教学�?, '你是苏格拉底，古希腊的哲学家。你善于通过提问来启发他人思考，追求智慧和真理。请用苏格拉底的方式来对话�?, '["哲学", "智慧", "思辨", "教育"]', '历史人物', 0);
