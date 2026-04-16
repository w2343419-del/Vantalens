# Vantalens（Hugo Blog + TalentWriter）

[![Hugo](https://img.shields.io/badge/Hugo-Extended-blueviolet?style=flat-square)](https://gohugo.io/)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue?style=flat-square)](LICENSE)

Vantalens 是一个基于 Hugo 的双语博客项目，配套本地管理工具 TalentWriter（Go）。

- 托管：GitHub Pages
- 统计：卜算子（Busuanzi）
- 评论：GitHub Issues 审核制（先审后显）

## 核心能力

- 双语内容管理（中文/英文）
- 本地可视化编辑与发布流程
- 评论审核、批量处理与导出
- 访问统计与访客 IP 统计（仅管理员可见）

## 快速开始

### 1. 运行 TalentWriter（双后端）

进入后端模块目录：

```bash
cd TalentWriter
```

启动总控后端（默认 9090）：

```bash
go run ./cmd/control
```

启动写作后端（默认 9091）：

```bash
go run ./cmd/writer
```

访问地址：

- 总控界面：http://127.0.0.1:9090/platform/control
- 写作界面：http://127.0.0.1:9091/platform/backend

说明：`TalentWriter.go` 属于历史单体入口，不再作为当前默认启动方式。

### 1.1 独立软件模式（推荐）

从当前版本开始，TalentWriter 支持独立软件模式，默认不自动拉起 Hugo，启动更快。

可通过环境变量控制启动行为：

- `TALENTWRITER_APP_MODE=standalone|dev`
- `TALENTWRITER_AUTOSTART_HUGO=true|false`
- `TALENTWRITER_AUTO_OPEN_BROWSER=true|false`
- `TALENTWRITER_SERVE_PUBLIC_SITE=true|false`

说明：

- `standalone`：默认不自动启动 Hugo，预览走内置 `/site/`（来自 `public/` 目录）
- `dev`：默认自动启动 Hugo（可被 `TALENTWRITER_AUTOSTART_HUGO` 覆盖）
- `TALENTWRITER_AUTO_OPEN_BROWSER=false` 可关闭自动弹浏览器，进一步减小启动干扰

Windows 示例（独立模式）：

```powershell
$env:TALENTWRITER_APP_MODE="standalone"
$env:TALENTWRITER_AUTOSTART_HUGO="false"
./TalentWriter.exe
```

### 2. 本地预览 Hugo 站点

```bash
hugo server
```

打开 http://localhost:1313/VantalensWeb/ 预览。

### 3. 构建桌面工具

```bash
go build -o TalentWriter.exe TalentWriter.go
```

### 3.1 构建桌面软件版（原生窗口）

桌面软件版会使用系统 WebView 打开管理界面，不依赖外部浏览器窗口。

如果需要把图标嵌入到 `TalentWriter-desktop.exe`，先生成 Windows 资源文件：

```bash
go run github.com/akavel/rsrc@latest -ico build/icons/wswriter.ico -o resource_windows.syso
```

构建命令（Windows）：

```bash
go build -tags desktop -o TalentWriter-desktop.exe .
```

运行示例：

```powershell
$env:TALENTWRITER_APP_MODE="desktop"
$env:TALENTWRITER_DESKTOP="true"
$env:TALENTWRITER_AUTOSTART_HUGO="false"
./TalentWriter-desktop.exe
```

可选项：

- `TALENTWRITER_EXIT_ON_WINDOW_CLOSE=true`：关闭窗口后退出程序（默认 true）
- `TALENTWRITER_EXIT_ON_WINDOW_CLOSE=false`：关闭窗口后后台继续运行服务

## 登录与权限

TalentWriter 登录使用本地后端鉴权（JWT）：

- 管理员账号由 .env 中 ADMIN_USERNAME / ADMIN_PASSWORD 配置
- JWT 密钥由 JWT_SECRET 配置
- 敏感接口（评论管理、统计、设置）需要已登录

## 评论工作流（GitHub Issues）

默认流程：访客提交 -> 生成 Issue（comment + pending）-> 管理员审核 -> approved 后展示。

配置文件：

- [config/_default/params.toml](config/_default/params.toml)
- [config/comment_settings.json](config/comment_settings.json)

## 统计方案

- 前台站点统计：卜算子脚本
- 管理后台统计：TalentWriter 聚合数据（含访客 IP）

详细说明见 [BUSUANZI_SETUP.md](BUSUANZI_SETUP.md)。

## 项目结构

```text
content/               # 博客内容（中英双语）
assets/                # 前端资源（JS/SCSS）
config/                # Hugo 与评论配置
layouts/               # 模板覆盖
static/                # 静态文件
TalentWriter.go            # TalentWriter 源码
TalentWriter.exe           # Windows 可执行文件
```

## 部署

当前推荐部署方式：

1. 使用 Hugo 构建静态站点
2. 推送到 GitHub 仓库
3. 由 GitHub Pages 托管

## 许可证

MIT License，详见 [LICENSE](LICENSE)。
