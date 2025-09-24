# AI角色扮演语音交互前端项目

一个基于 Vue 3 + Vite + Element Plus 的AI角色扮演语音交互应用。

## 📋 项目概述

这是一个现代化的前端应用，提供AI角色扮演和语音交互功能。用户可以与不同的AI角色进行对话，支持语音输入和输出。

## 🛠️ 技术栈

- **框架**: Vue 3 (Composition API)
- **构建工具**: Vite 5.x
- **UI组件库**: Element Plus 2.x
- **状态管理**: Pinia 2.x
- **路由**: Vue Router 4.x
- **HTTP客户端**: Axios
- **样式**: Sass
- **代码规范**: ESLint + Prettier

## 📁 项目结构

```
frontend/
├── public/                 # 静态资源目录
│   └── images/
│       └── avatars/       # 角色头像图片
├── src/
│   ├── assets/            # 项目资源文件
│   ├── components/        # 可复用组件
│   │   ├── CharacterCard.vue    # 角色卡片组件
│   │   ├── SearchBar.vue        # 搜索栏组件
│   │   └── VoiceController.vue  # 语音控制组件
│   ├── router/            # 路由配置
│   ├── services/          # API服务
│   ├── stores/            # Pinia状态管理
│   ├── utils/             # 工具函数
│   ├── views/             # 页面组件
│   │   ├── Chat.vue       # 聊天页面
│   │   ├── History.vue    # 历史记录页面
│   │   ├── Home.vue       # 首页
│   │   └── Settings.vue   # 设置页面
│   ├── App.vue            # 根组件
│   └── main.js            # 应用入口
├── index.html             # HTML入口文件
├── package.json           # 依赖配置
├── vite.config.js         # Vite配置
└── README.md              # 项目文档
```

## 🚀 快速开始

### 环境要求

- Node.js >= 16.0.0
- npm >= 8.0.0

### 安装依赖

```bash
# 进入前端项目目录
cd frontend

# 安装项目依赖
npm install
```

### 开发环境

```bash
# 启动开发服务器
npm run dev
```

开发服务器将在 `http://localhost:3000` 启动，支持热重载。

### 生产构建

```bash
# 构建生产版本
npm run build
```

构建文件将输出到 `dist/` 目录。

### 预览生产构建

```bash
# 预览生产构建
npm run preview
```

## 🔧 开发配置

### Vite 配置

项目使用 Vite 作为构建工具，主要配置包括：

- **开发服务器**: 端口3000，自动打开浏览器，支持所有网络接口访问
- **路径别名**: `@` 指向 `src` 目录
- **构建优化**: 禁用 sourcemap，chunk大小警告限制1500KB

### 代码规范

```bash
# 运行ESLint检查并自动修复
npm run lint
```

支持的文件类型：`.vue`, `.js`, `.jsx`, `.ts`, `.tsx` 等

## 📱 功能模块

### 主要页面

1. **首页 (Home.vue)**: 展示可用的AI角色，支持角色选择
2. **聊天页面 (Chat.vue)**: 与选中角色进行对话交互
3. **历史记录 (History.vue)**: 查看历史对话记录
4. **设置页面 (Settings.vue)**: 应用设置和配置

### 核心组件

1. **CharacterCard.vue**: 角色展示卡片
2. **SearchBar.vue**: 角色搜索功能
3. **VoiceController.vue**: 语音输入输出控制

## 🌐 API 集成

项目使用 Axios 进行 HTTP 请求，API服务文件位于 `src/services/` 目录。

## 📦 依赖说明

### 生产依赖

- `vue`: Vue 3 核心框架
- `vue-router`: 客户端路由
- `pinia`: 状态管理
- `element-plus`: UI组件库
- `axios`: HTTP客户端
- `@element-plus/icons-vue`: Element Plus 图标

### 开发依赖

- `@vitejs/plugin-vue`: Vue 3 的 Vite 插件
- `vite`: 现代化构建工具
- `eslint`: 代码质量检查
- `eslint-plugin-vue`: Vue 文件的 ESLint 插件
- `prettier`: 代码格式化
- `sass`: CSS 预处理器

## 🎯 开发建议

1. **组件开发**: 遵循 Vue 3 Composition API 最佳实践
2. **状态管理**: 使用 Pinia 进行全局状态管理
3. **样式编写**: 优先使用 Element Plus 组件，自定义样式使用 Sass
4. **API调用**: 统一在 services 目录下管理 API 请求
5. **代码规范**: 提交前运行 `npm run lint` 确保代码质量

## 🐛 常见问题

### Q: 开发服务器无法启动？
A: 确保 Node.js 版本 >= 16.0.0，删除 `node_modules` 重新安装依赖。

### Q: 构建失败？
A: 检查是否有语法错误，运行 `npm run lint` 检查代码规范。

### Q: 样式不生效？
A: 确保 Sass 正确安装，检查样式文件路径是否正确。

## 📄 许可证

本项目采用 MIT 许可证，详见 LICENSE 文件。

## 🤝 贡献指南

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交变更 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

---

如有问题或建议，请创建 Issue 或联系开发团队。
