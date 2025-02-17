# Drop It! 文件自动整理工具 📂

<p align="center">
  <img src="frontend/src/assets/images/logo.png" alt="Drop It Logo" width="200">
</p>

<p align="center">
  <strong>让文件整理变得简单而优雅</strong>
</p>

## ✨ 特性

- 🚀 **一键整理**：只需选择文件夹，点击按钮即可自动整理
- 🎯 **智能分类**：自动识别各种文件类型，包括：
  - 📸 图片（PNG, JPG, GIF等）
  - 📚 电子书（PDF, EPUB, MOBI等）
  - 📝 文档（Word, Excel, PPT等）
  - 🎵 音频（MP3, FLAC, WAV等）
  - 🎬 视频（MP4, MKV, AVI等）
  - 💾 程序（EXE, DMG, APK等）
  - 📦 压缩包（ZIP, RAR, 7Z等）
- 🔄 **自动备份**：原始文件会被保存在备份文件夹中
- ⚡ **高效性能**：基于 Go 语言开发，运行快速且占用资源少

## 🚀 快速开始

### 下载安装

访问 [Releases](https://github.com/yourusername/dropit/releases) 页面下载最新版本：

- Windows: `DropIt-windows-amd64.exe`

### 使用方法

1. 启动 Drop It
2. 点击"选择文件夹"按钮或使用文件夹选择器选择需要整理的文件夹
3. 点击"开始整理"按钮
4. 等待进度条完成即可

## 🛠️ 开发者指南

### 环境要求

- Go 1.18+
- Node.js 16+
- Wails CLI

### 本地开发

```bash
# 克隆项目
git clone https://github.com/yourusername/dropit.git

# 安装依赖
cd dropit
go mod tidy
cd frontend && npm install

# 启动开发服务器
wails dev
```
