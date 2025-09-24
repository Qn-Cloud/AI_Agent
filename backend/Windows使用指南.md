# AI 角色扮演语音交互产品 - Windows 使用指南

## 🎯 专为 Windows 用户设计

您提到的问题很对！我的 shell 脚本确实是为 Linux/macOS 设计的。为了让 Windows 用户也能方便使用，我专门创建了 Windows 版本的脚本。

## 🚀 Windows 快速开始

### 📋 准备工作

1. **确认 Windows 版本**
   - Windows 10 或 Windows 11
   - 建议开启 Windows Subsystem for Linux (WSL) 以获得更好的开发体验

2. **安装必要软件**

   **a) 安装 Go 语言**
   ```cmd
   # 下载 Go 安装包
   https://golang.org/dl/
   
   # 安装后验证
   go version
   ```

   **b) 安装 Git**
   ```cmd
   # 下载 Git for Windows
   https://git-scm.com/download/win
   ```

   **c) 安装 go-zero 工具**
   ```cmd
   go install github.com/zeromicro/go-zero/tools/goctl@latest
   ```

   **d) 安装 Docker Desktop (可选)**
   ```cmd
   # 下载 Docker Desktop for Windows
   https://www.docker.com/products/docker-desktop
   ```

### 🔧 三种运行方式

#### 方式一：PowerShell 脚本 (推荐)

```powershell
# 在 PowerShell 中运行
cd backend
.\scripts\init-project.ps1
```

如果提示执行策略错误，运行：
```powershell
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

#### 方式二：批处理文件 (最兼容)

```cmd
# 在命令提示符中运行
cd backend
scripts\init-project.bat
```

#### 方式三：WSL (Linux 体验)

```bash
# 在 WSL 中运行原始 bash 脚本
cd backend
./scripts/init-project.sh
```

## 📁 Windows 特有文件说明

```
backend/
├── scripts/
│   ├── init-project.sh     # Linux/macOS 版本
│   ├── init-project.ps1    # PowerShell 版本 ⭐
│   └── init-project.bat    # 批处理版本 ⭐
├── build.ps1               # PowerShell 构建脚本 ⭐
├── build.bat               # 批处理构建脚本 ⭐
├── README_Windows.md       # Windows 专用说明 ⭐
└── Windows使用指南.md      # 本文件 ⭐
```

## 🛠️ Windows 开发工作流

### 1. 项目初始化

**选择任一方式运行初始化脚本：**

```powershell
# PowerShell 方式 (推荐)
.\scripts\init-project.ps1

# 或者批处理方式
scripts\init-project.bat
```

### 2. 环境配置

```cmd
# 复制配置文件
copy .env.example .env

# 使用记事本编辑配置
notepad .env
```

### 3. 启动基础设施

**使用 Docker Desktop：**
```cmd
build.bat infra
```

**或手动启动（如果没有 Docker）：**
- 安装并启动 MySQL 8.0
- 安装并启动 Redis
- 导入数据库：`mysql -u root -p ai_roleplay < docs\sql\init.sql`

### 4. 构建项目

```cmd
# 构建所有服务
build.bat build

# 查看生成的可执行文件
dir bin\
```

### 5. 启动服务

```cmd
# 启动 API 网关
bin\gateway.exe

# 在新的命令行窗口中启动其他服务
bin\user-api.exe
bin\character-api.exe
bin\chat-api.exe
bin\ai-api.exe
bin\speech-api.exe
bin\storage-api.exe
```

## 🔧 Windows 开发工具推荐

### IDE 和编辑器
- **Visual Studio Code** + Go 插件
- **GoLand** (JetBrains)
- **LiteIDE**

### 终端工具
- **Windows Terminal** (推荐)
- **PowerShell 7**
- **Git Bash**

### 数据库工具
- **Navicat for MySQL**
- **DBeaver** (免费)
- **MySQL Workbench**

### API 测试
- **Postman**
- **Insomnia**
- **Thunder Client** (VS Code 插件)

## 🐛 Windows 常见问题解决

### 问题 1: PowerShell 执行策略限制

**错误信息：**
```
无法加载文件，因为在此系统上禁止运行脚本
```

**解决方案：**
```powershell
# 临时允许脚本执行
Set-ExecutionPolicy -ExecutionPolicy Bypass -Scope Process

# 或永久设置为当前用户
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### 问题 2: goctl 命令找不到

**错误信息：**
```
'goctl' 不是内部或外部命令
```

**解决方案：**
```cmd
# 查看 Go 环境
go env GOPATH

# 将 %GOPATH%\bin 添加到系统 PATH 环境变量
# 或重新安装 goctl
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

### 问题 3: 端口被占用

**检查端口占用：**
```cmd
# 查看端口占用
netstat -ano | findstr :8888

# 结束占用进程
taskkill /PID <进程ID> /F
```

### 问题 4: MySQL 连接失败

**解决步骤：**
1. 确认 MySQL 服务已启动
2. 检查防火墙设置
3. 验证连接配置：
   ```cmd
   mysql -h localhost -u root -p
   ```

### 问题 5: Docker 启动失败

**常见原因和解决方案：**
- **Hyper-V 未启用**：在 Windows 功能中启用 Hyper-V
- **WSL2 未安装**：安装 WSL2 更新包
- **内存不足**：调整 Docker Desktop 内存限制

## 🔄 Windows 与 Linux 差异对照

| 功能 | Linux/macOS | Windows |
|------|-------------|---------|
| 脚本类型 | bash (.sh) | PowerShell (.ps1) / 批处理 (.bat) |
| 路径分隔符 | `/` | `\` |
| 可执行文件 | 无后缀 | `.exe` |
| 包管理 | apt/brew | chocolatey/winget |
| 终端 | Terminal | PowerShell/CMD |

## 🎯 开发建议

### 1. 使用 Windows Terminal
- 支持多标签
- 更好的字体渲染
- 支持 PowerShell、CMD、WSL

### 2. 配置 Git
```cmd
git config --global user.name "Your Name"
git config --global user.email "your.email@example.com"
git config --global core.autocrlf true
```

### 3. 使用 WSL (可选)
- 获得类 Linux 开发体验
- 运行原始 bash 脚本
- 更好的容器支持

```cmd
# 安装 WSL
wsl --install

# 安装 Ubuntu
wsl --install -d Ubuntu
```

## 📚 学习资源

### Go 语言
- [Go 官方文档](https://golang.org/doc/)
- [Go by Example](https://gobyexample.com/)
- [Go 语言圣经](https://gopl-zh.github.io/)

### Go-Zero 框架
- [Go-Zero 官方文档](https://go-zero.dev/)
- [Go-Zero 视频教程](https://www.bilibili.com/video/BV1bv4y1X7nS)

### Windows 开发
- [Windows Terminal 文档](https://docs.microsoft.com/zh-cn/windows/terminal/)
- [PowerShell 学习指南](https://docs.microsoft.com/zh-cn/powershell/)

## 🆘 获取帮助

1. **查看脚本帮助**
   ```cmd
   build.bat help
   ```

2. **查看详细日志**
   ```cmd
   # 运行时查看详细输出
   .\scripts\init-project.ps1 -Verbose
   ```

3. **联系支持**
   - GitHub Issues: 提交问题报告
   - 技术交流群: 加入开发者社区
   - 邮件支持: backend-team@example.com

---

## 🎉 总结

现在您有了三个选择：

1. **PowerShell 脚本** (`init-project.ps1`) - 功能最全，推荐使用
2. **批处理脚本** (`init-project.bat`) - 兼容性最好，适合老版本 Windows
3. **WSL + Bash** - 如果您希望获得 Linux 开发体验

无论选择哪种方式，都能完整地初始化和运行您的 AI 角色扮演后端项目！

**建议：** 初学者推荐使用批处理脚本（`.bat`），有经验的开发者推荐 PowerShell 脚本（`.ps1`）。 