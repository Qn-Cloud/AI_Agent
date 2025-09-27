# AI角色扮演语音交互产品 - 数据库设计文档

## 概述

本文档描述了AI角色扮演语音交互产品的数据库设计，采用MySQL 8.0+，支持用户管理、角色管理、对话管理、语音处理等核心功能。

## 数据库配置

- **数据库名称**: `ai_roleplay`
- **字符集**: `utf8mb4`
- **排序规则**: `utf8mb4_unicode_ci`
- **引擎**: `InnoDB`

## 表结构设计

### 1. 用户表 (users)

存储用户基本信息和登录相关数据。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint(20) unsigned | 用户ID | 主键，自增 |
| username | varchar(50) | 用户名 | 唯一，非空 |
| email | varchar(100) | 邮箱 | 唯一，非空 |
| password_hash | varchar(255) | 密码哈希 | 非空 |
| nickname | varchar(50) | 昵称 | 可空 |
| avatar | varchar(500) | 头像URL | 可空 |
| bio | text | 个人简介 | 可空 |
| status | tinyint(3) unsigned | 状态：1正常 2禁用 | 默认1 |
| last_login_at | timestamp | 最后登录时间 | 可空 |
| last_login_ip | varchar(45) | 最后登录IP | 可空 |
| created_at | timestamp | 创建时间 | 自动填充 |
| updated_at | timestamp | 更新时间 | 自动更新 |

### 2. 角色分类表 (character_categories)

角色的分类管理，如经典IP、历史人物等。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint(20) unsigned | 分类ID | 主键，自增 |
| name | varchar(50) | 分类名称 | 唯一，非空 |
| description | text | 分类描述 | 可空 |
| sort_order | int(11) | 排序权重 | 默认0 |
| status | tinyint(3) unsigned | 状态：1启用 2禁用 | 默认1 |
| created_at | timestamp | 创建时间 | 自动填充 |
| updated_at | timestamp | 更新时间 | 自动更新 |

### 3. 角色表 (characters)

存储AI角色的详细信息，包括提示词、性格设置等。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint(20) unsigned | 角色ID | 主键，自增 |
| name | varchar(100) | 角色名称 | 非空 |
| avatar | varchar(500) | 角色头像URL | 可空 |
| description | text | 角色描述 | 可空 |
| short_desc | varchar(200) | 角色简介 | 可空 |
| category_id | bigint(20) unsigned | 分类ID | 外键 |
| tags | json | 标签列表 | 可空 |
| prompt | text | 角色提示词 | 可空 |
| personality | json | 性格设置 | 可空 |
| voice_settings | json | 语音设置 | 可空 |
| status | tinyint(3) unsigned | 状态：1正常 2禁用 | 默认1 |
| is_public | tinyint(1) | 是否公开：1公开 0私有 | 默认1 |
| creator_id | bigint(20) unsigned | 创建者ID，NULL表示系统预设 | 外键，可空 |
| rating | decimal(3,2) | 评分(0-5) | 默认0.00 |
| rating_count | int(11) | 评分人数 | 默认0 |
| favorite_count | int(11) | 收藏数 | 默认0 |
| chat_count | int(11) | 对话次数 | 默认0 |
| created_at | timestamp | 创建时间 | 自动填充 |
| updated_at | timestamp | 更新时间 | 自动更新 |

### 4. 用户角色收藏表 (user_character_favorites)

用户收藏角色的关联表。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint(20) unsigned | 收藏ID | 主键，自增 |
| user_id | bigint(20) unsigned | 用户ID | 外键，非空 |
| character_id | bigint(20) unsigned | 角色ID | 外键，非空 |
| created_at | timestamp | 收藏时间 | 自动填充 |

### 5. 对话表 (conversations)

存储对话会话信息。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint(20) unsigned | 对话ID | 主键，自增 |
| user_id | bigint(20) unsigned | 用户ID，NULL表示匿名用户 | 外键，可空 |
| character_id | bigint(20) unsigned | 角色ID | 外键，非空 |
| title | varchar(200) | 对话标题 | 默认'新对话' |
| start_time | timestamp | 开始时间 | 自动填充 |
| last_message_time | timestamp | 最后消息时间 | 自动填充 |
| message_count | int(11) | 消息数量 | 默认0 |
| status | tinyint(3) unsigned | 状态：1正常 2已删除 | 默认1 |
| settings | json | 对话设置 | 可空 |
| session_id | varchar(64) | 会话标识(用于匿名用户) | 可空 |
| created_at | timestamp | 创建时间 | 自动填充 |
| updated_at | timestamp | 更新时间 | 自动更新 |

### 6. 消息表 (messages)

存储对话中的具体消息内容。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint(20) unsigned | 消息ID | 主键，自增 |
| conversation_id | bigint(20) unsigned | 对话ID | 外键，非空 |
| type | enum('user','ai') | 消息类型：user用户 ai系统 | 非空 |
| content | text | 消息内容 | 非空 |
| audio_id | bigint(20) unsigned | 语音文件ID | 外键，可空 |
| metadata | json | 元数据，存储额外信息 | 可空 |
| token_used | int(11) | AI消息使用的token数 | 默认0 |
| processing_time | int(11) | 处理时间(毫秒) | 默认0 |
| created_at | timestamp | 创建时间 | 自动填充 |

### 7. 语音文件表 (audio_files)

存储语音文件的元数据和路径信息。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint(20) unsigned | 语音文件ID | 主键，自增 |
| user_id | bigint(20) unsigned | 用户ID | 外键，可空 |
| type | enum('stt','tts') | 类型：stt语音识别 tts语音合成 | 非空 |
| filename | varchar(255) | 文件名 | 非空 |
| file_path | varchar(500) | 文件路径 | 非空 |
| file_size | bigint(20) | 文件大小(字节) | 默认0 |
| duration | int(11) | 音频时长(毫秒) | 默认0 |
| format | varchar(20) | 音频格式 | 非空 |
| sample_rate | int(11) | 采样率 | 可空 |
| bit_rate | int(11) | 比特率 | 可空 |
| text_content | text | 对应的文本内容 | 可空 |
| character_id | bigint(20) unsigned | 关联角色ID(TTS时使用) | 外键，可空 |
| voice_settings | json | 语音设置 | 可空 |
| status | tinyint(3) unsigned | 状态：1正常 2已删除 | 默认1 |
| created_at | timestamp | 创建时间 | 自动填充 |
| updated_at | timestamp | 更新时间 | 自动更新 |

### 8. 系统配置表 (system_configs)

存储系统级别的配置参数。

| 字段名 | 类型 | 说明 | 约束 |
|--------|------|------|------|
| id | bigint(20) unsigned | 配置ID | 主键，自增 |
| config_key | varchar(100) | 配置键 | 唯一，非空 |
| config_value | text | 配置值 | 可空 |
| description | varchar(500) | 配置描述 | 可空 |
| config_type | varchar(50) | 配置类型：string,number,boolean,json | 默认'string' |
| is_public | tinyint(1) | 是否公开：1公开 0私有 | 默认0 |
| created_at | timestamp | 创建时间 | 自动填充 |
| updated_at | timestamp | 更新时间 | 自动更新 |

## 预设数据

### 角色分类

1. **经典IP** - 经典文学、影视作品中的知名角色
2. **历史人物** - 历史上真实存在的著名人物
3. **现代名人** - 当代知名的公众人物
4. **原创角色** - 用户自创或系统原创的虚拟角色
5. **专业助手** - 各领域的专业助手角色

### 预设角色

1. **哈利·波特** - 霍格沃茨的勇敢学生
2. **苏格拉底** - 智慧的古希腊哲学家
3. **莎士比亚** - 伟大的剧作家和诗人
4. **爱因斯坦** - 伟大的物理学家
5. **夏洛克·福尔摩斯** - 世界著名的咨询侦探
6. **赫敏·格兰杰** - 霍格沃茨的聪明学生

每个角色都包含详细的描述、提示词、性格设置和语音配置。

### 系统配置

- AI模型设置
- 消息长度限制
- 语音相关配置
- 文件上传限制
- 速率限制等

## 索引设计

### 主要索引

1. **唯一索引**: 用户名、邮箱、分类名称、配置键
2. **外键索引**: 所有外键字段都创建了索引
3. **状态索引**: 状态字段用于快速过滤
4. **时间索引**: 创建时间、更新时间用于排序
5. **全文索引**: 角色名称描述、对话标题、消息内容支持全文搜索

### 复合索引

- `user_character_favorites`: (user_id, character_id) 唯一索引
- 其他根据查询需求优化的复合索引

## 视图设计

### v_character_details
包含角色详细信息，关联分类名称和创建者信息的视图。

### v_conversation_list
包含对话列表信息，关联角色名称和用户信息的视图。

## 数据备份与恢复

建议定期备份数据库，特别是：
- 用户数据和对话记录
- 自定义角色数据
- 语音文件元数据

## 性能优化建议

1. **索引优化**: 根据实际查询模式调整索引
2. **分区策略**: 对于消息表可考虑按时间分区
3. **缓存策略**: 热门角色和配置数据使用Redis缓存
4. **读写分离**: 高并发场景下可考虑主从复制
5. **定期清理**: 定期清理已删除的数据和过期文件

## 扩展性考虑

1. **水平扩展**: 设计支持分库分表
2. **多语言支持**: 字符集已支持多语言
3. **版本管理**: 角色和配置支持版本控制
4. **审计日志**: 可扩展操作日志记录 