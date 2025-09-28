# API 配置说明

## 🚀 微服务架构配置

本项目已更新为支持微服务架构，不同的服务使用不同的端口。

### 📋 服务端口配置

| 服务 | 端口 | 说明 |
|------|------|------|
| 聊天服务 | 7001 | 对话管理、消息处理 |
| 角色服务 | 7002 | 角色管理、角色信息 |
| 用户服务 | 7003 | 用户认证、用户信息 |
| AI服务 | 7004 | AI对话、智能回复 |
| 语音服务 | 7005 | 语音识别、语音合成 |
| 存储服务 | 7006 | 文件上传、媒体存储 |
| 网关服务 | 8888 | API网关（可选） |

### ⚙️ 环境变量配置

在项目根目录创建 `.env.development` 文件：

```env
# 开发环境配置
VITE_APP_TITLE=AI角色扮演语音交互产品
VITE_APP_VERSION=1.0.0

# 微服务API地址
VITE_CHAT_API_URL=http://192.168.23.188:7001
VITE_CHARACTER_API_URL=http://192.168.23.188:7002
VITE_USER_API_URL=http://192.168.23.188:7003
VITE_AI_API_URL=http://192.168.23.188:7004
VITE_SPEECH_API_URL=http://192.168.23.188:7005
VITE_STORAGE_API_URL=http://192.168.23.188:7006

# 网关地址（如果有）
VITE_GATEWAY_API_URL=http://192.168.23.188:8888

# 默认API地址（向后兼容）
VITE_API_BASE_URL=http://192.168.23.188:8888
```

### 🔧 代码修改说明

#### 1. 配置文件更新
- `src/config/index.js` - 添加了微服务配置
- `src/services/apiFactory.js` - 新增API工厂管理不同服务

#### 2. 服务文件更新
- `src/services/chat.js` - 使用聊天服务API实例
- `src/services/character.js` - 使用角色服务API实例
- `src/services/index.js` - 重新组织导出结构

#### 3. 向后兼容
为保持现有代码正常运行，保留了原有的导出方式：
```javascript
import { chatApi, characterApi } from '@/services'
```

### 🌐 部署配置

#### 开发环境
服务器地址：`192.168.23.188`
- 聊天服务：`:7001`
- 角色服务：`:7002`

#### 生产环境
创建 `.env.production` 文件并修改相应的生产环境地址。

### 🔍 故障排查

1. **CORS问题**：确保后端服务配置了正确的CORS设置
2. **网络连接**：检查服务器地址和端口是否可访问
3. **服务状态**：确认各个微服务是否正常运行

### 📝 使用示例

```javascript
// 使用聊天服务
import { chatApiService } from '@/services'
await chatApiService.getConversationHistory(params)

// 使用角色服务  
import { characterApiService } from '@/services'
await characterApiService.getCharacterList(params)

// 或者使用原有方式（向后兼容）
import { chatApi, characterApi } from '@/services'
await chatApi.getConversationList(params)
await characterApi.getCharacterList(params)
``` 