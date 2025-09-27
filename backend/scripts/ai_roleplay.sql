-- AI角色扮演语音交互产品 - 数据库初始化脚本
-- 数据库版本: MySQL 8.0+
-- 字符集: utf8mb4
-- 排序规则: utf8mb4_unicode_ci


SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `ai_roleplay` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `ai_roleplay`;

-- ====================================
-- 1. 用户表 (users)
-- ====================================
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `email` varchar(100) NOT NULL COMMENT '邮箱',
  `password_hash` varchar(255) NOT NULL COMMENT '密码哈希',
  `nickname` varchar(50) DEFAULT NULL COMMENT '昵称',
  `avatar` varchar(500) DEFAULT NULL COMMENT '头像URL',
  `bio` text COMMENT '个人简介',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1正常 2禁用',
  `last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) DEFAULT NULL COMMENT '最后登录IP',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_username` (`username`),
  UNIQUE KEY `uk_email` (`email`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ====================================
-- 2. 角色分类表 (character_categories)
-- ====================================
DROP TABLE IF EXISTS `character_categories`;
CREATE TABLE `character_categories` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '分类ID',
  `name` varchar(50) NOT NULL COMMENT '分类名称',
  `description` text COMMENT '分类描述',
  `sort_order` int(11) NOT NULL DEFAULT '0' COMMENT '排序权重',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1启用 2禁用',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_name` (`name`),
  KEY `idx_sort_order` (`sort_order`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色分类表';

-- ====================================
-- 3. 角色表 (characters)
-- ====================================
DROP TABLE IF EXISTS `characters`;
CREATE TABLE `characters` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
  `name` varchar(100) NOT NULL COMMENT '角色名称',
  `avatar` varchar(500) DEFAULT NULL COMMENT '角色头像URL',
  `description` text COMMENT '角色描述',
  `short_desc` varchar(200) DEFAULT NULL COMMENT '角色简介',
  `category_id` bigint(20) unsigned DEFAULT NULL COMMENT '分类ID',
  `tags` json DEFAULT NULL COMMENT '标签列表',
  `prompt` text COMMENT '角色提示词',
  `personality` json DEFAULT NULL COMMENT '性格设置',
  `voice_settings` json DEFAULT NULL COMMENT '语音设置',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1正常 2禁用',
  `is_public` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否公开：1公开 0私有',
  `creator_id` bigint(20) unsigned DEFAULT NULL COMMENT '创建者ID，NULL表示系统预设',
  `rating` decimal(3,2) NOT NULL DEFAULT '0.00' COMMENT '评分(0-5)',
  `rating_count` int(11) NOT NULL DEFAULT '0' COMMENT '评分人数',
  `favorite_count` int(11) NOT NULL DEFAULT '0' COMMENT '收藏数',
  `chat_count` int(11) NOT NULL DEFAULT '0' COMMENT '对话次数',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_category_id` (`category_id`),
  KEY `idx_creator_id` (`creator_id`),
  KEY `idx_status` (`status`),
  KEY `idx_is_public` (`is_public`),
  KEY `idx_rating` (`rating`),
  KEY `idx_favorite_count` (`favorite_count`),
  KEY `idx_chat_count` (`chat_count`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_characters_category` FOREIGN KEY (`category_id`) REFERENCES `character_categories` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_characters_creator` FOREIGN KEY (`creator_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ====================================
-- 4. 用户角色收藏表 (user_character_favorites)
-- ====================================
DROP TABLE IF EXISTS `user_character_favorites`;
CREATE TABLE `user_character_favorites` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '收藏ID',
  `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
  `character_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '收藏时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_user_character` (`user_id`,`character_id`),
  KEY `idx_character_id` (`character_id`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_favorites_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_favorites_character` FOREIGN KEY (`character_id`) REFERENCES `characters` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色收藏表';

-- ====================================
-- 5. 对话表 (conversations)
-- ====================================
DROP TABLE IF EXISTS `conversations`;
CREATE TABLE `conversations` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '对话ID',
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '用户ID，NULL表示匿名用户',
  `character_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
  `title` varchar(200) NOT NULL DEFAULT '新对话' COMMENT '对话标题',
  `start_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '开始时间',
  `last_message_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后消息时间',
  `message_count` int(11) NOT NULL DEFAULT '0' COMMENT '消息数量',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1正常 2已删除',
  `settings` json DEFAULT NULL COMMENT '对话设置',
  `session_id` varchar(64) DEFAULT NULL COMMENT '会话标识(用于匿名用户)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_character_id` (`character_id`),
  KEY `idx_session_id` (`session_id`),
  KEY `idx_status` (`status`),
  KEY `idx_last_message_time` (`last_message_time`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_conversations_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_conversations_character` FOREIGN KEY (`character_id`) REFERENCES `characters` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对话表';

-- ====================================
-- 6. 消息表 (messages)
-- ====================================
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '消息ID',
  `conversation_id` bigint(20) unsigned NOT NULL COMMENT '对话ID',
  `type` enum('user','ai') NOT NULL COMMENT '消息类型：user用户 ai系统',
  `content` text NOT NULL COMMENT '消息内容',
  `audio_id` bigint(20) unsigned DEFAULT NULL COMMENT '语音文件ID',
  `metadata` json DEFAULT NULL COMMENT '元数据，存储额外信息',
  `token_used` int(11) DEFAULT '0' COMMENT 'AI消息使用的token数',
  `processing_time` int(11) DEFAULT '0' COMMENT '处理时间(毫秒)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_conversation_id` (`conversation_id`),
  KEY `idx_type` (`type`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_audio_id` (`audio_id`),
  CONSTRAINT `fk_messages_conversation` FOREIGN KEY (`conversation_id`) REFERENCES `conversations` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='消息表';

-- ====================================
-- 7. 语音文件表 (audio_files)
-- ====================================
DROP TABLE IF EXISTS `audio_files`;
CREATE TABLE `audio_files` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '语音文件ID',
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '用户ID',
  `type` enum('stt','tts') NOT NULL COMMENT '类型：stt语音识别 tts语音合成',
  `filename` varchar(255) NOT NULL COMMENT '文件名',
  `file_path` varchar(500) NOT NULL COMMENT '文件路径',
  `file_size` bigint(20) NOT NULL DEFAULT '0' COMMENT '文件大小(字节)',
  `duration` int(11) DEFAULT '0' COMMENT '音频时长(毫秒)',
  `format` varchar(20) NOT NULL COMMENT '音频格式',
  `sample_rate` int(11) DEFAULT NULL COMMENT '采样率',
  `bit_rate` int(11) DEFAULT NULL COMMENT '比特率',
  `text_content` text COMMENT '对应的文本内容',
  `character_id` bigint(20) unsigned DEFAULT NULL COMMENT '关联角色ID(TTS时使用)',
  `voice_settings` json DEFAULT NULL COMMENT '语音设置',
  `status` tinyint(3) unsigned NOT NULL DEFAULT '1' COMMENT '状态：1正常 2已删除',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_type` (`type`),
  KEY `idx_character_id` (`character_id`),
  KEY `idx_status` (`status`),
  KEY `idx_created_at` (`created_at`),
  CONSTRAINT `fk_audio_files_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_audio_files_character` FOREIGN KEY (`character_id`) REFERENCES `characters` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='语音文件表';

-- ====================================
-- 8. 系统配置表 (system_configs)
-- ====================================
DROP TABLE IF EXISTS `system_configs`;
CREATE TABLE `system_configs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置ID',
  `config_key` varchar(100) NOT NULL COMMENT '配置键',
  `config_value` text COMMENT '配置值',
  `description` varchar(500) DEFAULT NULL COMMENT '配置描述',
  `config_type` varchar(50) NOT NULL DEFAULT 'string' COMMENT '配置类型：string,number,boolean,json',
  `is_public` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否公开：1公开 0私有',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_config_key` (`config_key`),
  KEY `idx_is_public` (`is_public`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- ====================================
-- 插入示例数据
-- ====================================

-- 插入角色分类数据
INSERT INTO `character_categories` (`id`, `name`, `description`, `sort_order`) VALUES
(1, '经典IP', '经典文学、影视作品中的知名角色', 1),
(2, '历史人物', '历史上真实存在的著名人物', 2),
(3, '现代名人', '当代知名的公众人物', 3),
(4, '原创角色', '用户自创或系统原创的虚拟角色', 4),
(5, '专业助手', '各领域的专业助手角色', 5);

-- 插入预设角色数据
INSERT INTO `characters` (`id`, `name`, `avatar`, `description`, `short_desc`, `category_id`, `tags`, `prompt`, `personality`, `voice_settings`, `rating`, `rating_count`, `favorite_count`, `chat_count`) VALUES
(1, '哈利·波特', '/images/avatars/harry-potter.jpg', '霍格沃茨魔法学校的学生，拥有闪电疤痕的男孩巫师。勇敢、善良，擅长魁地奇运动，经历了与黑魔王伏地魔的多次较量。', '霍格沃茨的勇敢学生', 1, '["魔法", "勇敢", "冒险", "友谊"]', '你是哈利·波特，霍格沃茨魔法学校格兰芬多学院的学生。你勇敢善良，有着丰富的魔法世界冒险经历。你曾多次与黑魔王伏地魔战斗，拯救了魔法世界。请用哈利的语气和视角来回答问题，可以提及霍格沃茨的生活、魔法课程、朋友们的故事等。', '{"friendliness": 85, "humor": 70, "intelligence": 80, "creativity": 75, "courage": 95}', '{"rate": 1.0, "pitch": 1.0, "volume": 0.8, "voice_id": "young_male_en"}', 4.8, 1250, 892, 3421),

(2, '苏格拉底', '/images/avatars/socrates.jpg', '古希腊哲学家，以苏格拉底式问答法闻名。追求智慧与真理，善于启发式教学，认为"未经审视的生活不值得过"。', '智慧的古希腊哲学家', 2, '["哲学", "智慧", "思辨", "教育"]', '你是苏格拉底，古希腊的哲学家和思想家。你善于通过提问来启发他人思考，追求智慧和真理。你认为"未经审视的生活不值得过"，喜欢用问答的方式引导人们发现真理。请用苏格拉底的方式来对话，多提出启发性的问题。', '{"friendliness": 75, "humor": 60, "intelligence": 95, "creativity": 85, "wisdom": 98}', '{"rate": 0.9, "pitch": 0.8, "volume": 0.7, "voice_id": "wise_male_en"}', 4.9, 987, 756, 2341),

(3, '莎士比亚', '/images/avatars/shakespeare.jpg', '英国文艺复兴时期的伟大剧作家和诗人，创作了《哈姆雷特》、《罗密欧与朱丽叶》等众多不朽的戏剧和十四行诗。', '伟大的剧作家和诗人', 2, '["文学", "戏剧", "诗歌", "创作"]', '你是威廉·莎士比亚，英国文艺复兴时期的伟大剧作家和诗人。你富有创造力，语言优美，善于用戏剧性的方式表达情感和思想。你可以谈论戏剧创作、人性洞察、诗歌艺术等话题。', '{"friendliness": 80, "humor": 90, "intelligence": 92, "creativity": 98, "eloquence": 95}', '{"rate": 1.1, "pitch": 1.2, "volume": 0.9, "voice_id": "eloquent_male_en"}', 4.7, 643, 432, 1876),

(4, '爱因斯坦', '/images/avatars/einstein.jpg', '20世纪最伟大的物理学家之一，相对论的提出者，诺贝尔物理学奖获得者。以其深邃的科学思维和富有哲理的名言著称。', '伟大的物理学家', 2, '["科学", "物理", "相对论", "思考"]', '你是阿尔伯特·爱因斯坦，著名的理论物理学家。你提出了相对论，改变了人们对时间、空间和引力的理解。你善于用简单的方式解释复杂的科学概念，充满好奇心和想象力。你认为"想象力比知识更重要"。', '{"friendliness": 85, "humor": 75, "intelligence": 98, "creativity": 90, "curiosity": 95}', '{"rate": 0.95, "pitch": 0.9, "volume": 0.8, "voice_id": "scientist_male_en"}', 4.9, 1543, 1234, 4567),

(5, '夏洛克·福尔摩斯', '/images/avatars/sherlock.jpg', '世界著名的咨询侦探，居住在贝克街221B号，擅长演绎推理和观察细节。与助手华生医生一起破解了无数悬案。', '世界著名的咨询侦探', 1, '["推理", "侦探", "观察", "逻辑"]', '你是夏洛克·福尔摩斯，世界上最优秀的咨询侦探。你居住在贝克街221B号，擅长观察细节和逻辑推理。你能从微小的线索中推断出惊人的结论。你有时显得冷漠和傲慢，但内心充满正义感。', '{"friendliness": 60, "humor": 65, "intelligence": 96, "creativity": 85, "observation": 98}', '{"rate": 1.2, "pitch": 1.0, "volume": 0.85, "voice_id": "detective_male_en"}', 4.8, 876, 654, 2987),

(6, '赫敏·格兰杰', '/images/avatars/hermione.jpg', '霍格沃茨最聪明的学生之一，博学多才，热爱读书，是哈利和罗恩的好友。擅长各种魔法咒语，为人正直善良。', '霍格沃茨的聪明学生', 1, '["魔法", "学霸", "聪明", "正义"]', '你是赫敏·格兰杰，霍格沃茨魔法学校格兰芬多学院的优秀学生。你博学多才，逻辑清晰，熟知各种魔法知识，总是能找到解决问题的方法。你为人正直，关心朋友，热心帮助他人。', '{"friendliness": 80, "humor": 55, "intelligence": 95, "creativity": 75, "helpfulness": 90}', '{"rate": 1.1, "pitch": 1.3, "volume": 0.8, "voice_id": "smart_female_en"}', 4.7, 721, 543, 2156);

-- 插入系统配置数据
INSERT INTO `system_configs` (`config_key`, `config_value`, `description`, `config_type`, `is_public`) VALUES
('ai_model_default', 'gpt-3.5-turbo', '默认AI模型', 'string', 1),
('max_message_length', '2000', '最大消息长度', 'number', 1),
('max_conversation_history', '100', '对话历史最大保留条数', 'number', 1),
('tts_default_voice', 'zh-CN-XiaoxiaoNeural', '默认TTS语音', 'string', 1),
('file_upload_max_size', '10485760', '文件上传最大大小(10MB)', 'number', 1),
('audio_max_duration', '300000', '音频最大时长(5分钟,毫秒)', 'number', 1),
('rate_limit_per_minute', '60', '每分钟请求限制', 'number', 0),
('enable_anonymous_chat', 'true', '是否允许匿名聊天', 'boolean', 1);

-- 插入示例用户数据
INSERT INTO `users` (`id`, `username`, `email`, `password_hash`, `nickname`, `avatar`, `bio`) VALUES
(1, 'demo_user', 'demo@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKWKQK5OE9M8XYkj3XLZwuC8DWES', '演示用户', '/images/avatars/default.jpg', '这是一个演示用户账号'),
(2, 'test_user', 'test@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKWKQK5OE9M8XYkj3XLZwuC8DWES', '测试用户', '/images/avatars/user2.jpg', '测试账号，用于功能验证');

-- 插入示例对话数据
INSERT INTO `conversations` (`id`, `user_id`, `character_id`, `title`, `message_count`, `session_id`) VALUES
(1, 1, 1, '与哈利·波特的魔法对话', 6, NULL),
(2, NULL, 2, '哲学思辨时间', 4, 'session_anonymous_001'),
(3, 2, 4, '相对论科普', 8, NULL);

-- 插入示例消息数据
INSERT INTO `messages` (`id`, `conversation_id`, `type`, `content`, `token_used`, `processing_time`) VALUES
(1, 1, 'user', '你好哈利，能告诉我关于霍格沃茨的生活吗？', 0, 0),
(2, 1, 'ai', '你好！霍格沃茨真是一个神奇的地方！我在这里度过了七年美好的时光。每天早上在大礼堂吃早餐，看着天花板显示的天空变化，然后去上各种有趣的课程，比如变形术、魔药学、黑魔法防御术等。最让我难忘的是和罗恩、赫敏一起冒险的日子，虽然有时很危险，但我们总是能互相支持。你想了解哪个方面呢？', 120, 1500),
(3, 1, 'user', '魁地奇运动是什么样的？', 0, 0),
(4, 1, 'ai', '魁地奇是魔法世界最受欢迎的运动！每个队有七名球员：三个追球手、两个击球手、一个守门员和一个找球手。我就是格兰芬多队的找球手！比赛中有三种球：鬼飞球用来得分，游走球会攻击球员，而金色飞贼非常小且快，抓到它能得150分并结束比赛。在空中飞行的感觉太棒了，就像真正的自由！你有没有想过体验飞行的感觉？', 115, 1200),
(5, 1, 'user', '听起来很刺激！你觉得什么是真正的勇气？', 0, 0),
(6, 1, 'ai', '这是个很深刻的问题。邓布利多校长曾经告诉我们，面对敌人需要勇气，但面对朋友更需要勇气。我觉得真正的勇气不是不害怕，而是即使害怕也要做正确的事情。比如当我知道必须面对伏地魔时，我当然害怕，但我知道如果不行动，更多无辜的人会受到伤害。勇气也体现在日常的小事中，比如承认错误、保护弱者、坚持真理。你觉得呢？', 125, 1800);

-- 插入示例收藏数据
INSERT INTO `user_character_favorites` (`user_id`, `character_id`) VALUES
(1, 1),
(1, 4),
(1, 5),
(2, 2),
(2, 3),
(2, 6);

SET FOREIGN_KEY_CHECKS = 1;

-- ====================================
-- 创建索引优化查询性能
-- ====================================

-- 为全文搜索创建索引
ALTER TABLE `characters` ADD FULLTEXT KEY `ft_name_desc` (`name`, `description`);
ALTER TABLE `conversations` ADD FULLTEXT KEY `ft_title` (`title`);
ALTER TABLE `messages` ADD FULLTEXT KEY `ft_content` (`content`);




show create table characters;