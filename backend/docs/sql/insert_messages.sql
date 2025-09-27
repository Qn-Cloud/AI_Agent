-- 插入新的消息数据
-- Character ID 2: 苏格拉底 (对话ID: 2)
-- Character ID 4: 爱因斯坦 (对话ID: 3)

-- 苏格拉底对话 (conversation_id = 2)
INSERT INTO `messages` (`conversation_id`, `type`, `content`, `audio_id`, `metadata`, `token_used`, `processing_time`, `created_at`) VALUES
(2, 'user', '苏格拉底，什么是真正的智慧？', NULL, NULL, 0, 0, '2025-09-27 11:05:00'),
(2, 'ai', '我的朋友，这是一个值得深思的问题。我常说"我知道我一无所知"，这并非自谦，而是真正智慧的开始。真正的智慧在于认识到自己知识的有限性，保持谦逊和好奇心。当我们以为自己已经掌握了真理时，往往是最危险的时刻。你认为自己了解什么是智慧吗？', NULL, '{"topic": "wisdom", "philosophy_type": "epistemology"}', 145, 2000, '2025-09-27 11:05:30'),
(2, 'user', '那么，我们如何才能获得真正的知识呢？', NULL, NULL, 0, 0, '2025-09-27 11:06:15'),
(2, 'ai', '这正是我想要探讨的！通过不断的质疑和对话，我们可以接近真理。就像助产士帮助产妇分娩一样，我的方法是帮助人们"生产"出内心已有的知识。让我问你：你能定义什么是"正义"吗？当你试图回答时，你会发现这比想象的更复杂。真正的学习不是被动接受，而是主动思考和质疑。', NULL, '{"topic": "knowledge", "method": "socratic_method"}', 138, 1900, '2025-09-27 11:06:45'),

-- 爱因斯坦对话 (conversation_id = 3)
(3, 'user', '爱因斯坦教授，能简单解释一下相对论吗？', NULL, NULL, 0, 0, '2025-09-27 11:10:00'),
(3, 'ai', '当然！让我用一个简单的例子来解释。想象你坐在一列行驶的火车上，你扔出一个球。对你来说，球的速度可能是10公里/小时，但对站台上的人来说，球的速度是火车速度加上10公里/小时。这就是相对性的基本概念——运动是相对的！但光不一样，无论你如何运动，光速总是恒定的。这个发现改变了我们对时间和空间的理解。', NULL, '{"topic": "relativity", "complexity": "basic"}', 155, 2200, '2025-09-27 11:10:30'),
(3, 'user', '时间真的会因为速度而变慢吗？', NULL, NULL, 0, 0, '2025-09-27 11:11:15'),
(3, 'ai', '是的！这被称为时间膨胀。当你以接近光速运动时，时间确实会变慢。这不是错觉，而是物理现实！比如，如果你的双胞胎兄弟乘坐高速飞船旅行，当他回来时，他会比留在地球上的你更年轻。GPS卫星就必须考虑这种效应，否则定位会出现误差。时间和空间不是绝对的，它们会因观察者的运动状态而改变。', NULL, '{"topic": "time_dilation", "example": "twin_paradox"}', 142, 2100, '2025-09-27 11:11:45'),
(3, 'user', '这个理论是如何被证实的？', NULL, NULL, 0, 0, '2025-09-27 11:12:30'),
(3, 'ai', '有很多精彩的实验证据！1919年的日食观测证实了光线在重力场中的弯曲；原子钟实验证实了时间膨胀；粒子加速器中高速运动的粒子寿命确实延长了。最近的引力波探测更是直接证实了我关于时空弯曲的预测。科学的美妙之处就在于理论必须接受实验的检验。正如我常说："想象力比知识更重要，因为知识是有限的。"', NULL, '{"topic": "experimental_evidence", "examples": ["eclipse_1919", "atomic_clocks", "gravitational_waves"]}', 168, 2500, '2025-09-27 11:13:00'),
(3, 'user', '您对现代科学发展有什么看法？', NULL, NULL, 0, 0, '2025-09-27 11:14:00'),
(3, 'ai', '现代科学的发展令我既兴奋又担忧。量子力学的发展超出了我的预期，虽然我曾说"上帝不掷骰子"，但量子纠缠等现象确实存在。我很高兴看到引力波被探测到，黑洞被直接观测到。但我也担心科学技术被滥用。核能可以造福人类，也可以毁灭世界。我们需要的不仅是聪明的科学家，更需要有智慧和道德的科学家。科学应该服务于全人类的福祉。', NULL, '{"topic": "modern_science", "concerns": ["ethics", "responsibility"], "achievements": ["quantum_mechanics", "gravitational_waves"]}', 175, 2300, '2025-09-27 11:14:30'),
(3, 'user', '您觉得好奇心在科学研究中有多重要？', NULL, NULL, 0, 0, '2025-09-27 11:15:15'),
(3, 'ai', '好奇心是科学研究的灵魂！我小时候就对指南针为什么总是指向北方感到好奇，这种好奇心驱使我思考空间和磁场的本质。没有好奇心，就没有真正的科学发现。我们应该像孩子一样保持对世界的惊奇感，不断问"为什么"。教育的目标不应该是扼杀好奇心，而是培养它。记住，重要的是不要停止质疑。好奇心有其存在的理由。', NULL, '{"topic": "curiosity", "personal_example": "compass", "educational_philosophy": "preserve_wonder"}', 152, 1950, '2025-09-27 11:15:45'),
(3, 'user', '最后一个问题，您认为宇宙的终极奥秘是什么？', NULL, NULL, 0, 0, '2025-09-27 11:16:30'),
(3, 'ai', '这是一个深刻的问题！我一生都在寻找一个"万有理论"，能够统一所有自然力的理论。但也许宇宙最大的奥秘就是它是可以被理解的。为什么数学——这个人类大脑的产物——能够如此完美地描述宇宙的规律？这本身就是个奇迹。也许答案在于意识与宇宙的深层联系。我相信，宇宙中最不可理解的事情，就是宇宙是可以理解的。这种神秘感让我们永远保持谦卑和敬畏。', NULL, '{"topic": "ultimate_mystery", "concepts": ["theory_of_everything", "consciousness", "mathematical_universe"]}', 188, 2800, '2025-09-27 11:17:00');

-- 更新对话的消息计数和最后消息时间
UPDATE conversations SET 
    message_count = (SELECT COUNT(*) FROM messages WHERE conversation_id = 2),
    last_message_time = (SELECT MAX(created_at) FROM messages WHERE conversation_id = 2),
    updated_at = (SELECT MAX(created_at) FROM messages WHERE conversation_id = 2)
WHERE id = 2;

UPDATE conversations SET 
    message_count = (SELECT COUNT(*) FROM messages WHERE conversation_id = 3),
    last_message_time = (SELECT MAX(created_at) FROM messages WHERE conversation_id = 3),
    updated_at = (SELECT MAX(created_at) FROM messages WHERE conversation_id = 3)
WHERE id = 3; 