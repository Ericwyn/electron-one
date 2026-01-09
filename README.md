# electron-one

使用 Golang 编写的 Electron 快速启动工具，可以快速将本地 HTML 项目以 Electron 形式运行。

## 功能特性

- 快速启动本地 HTML 项目为 Electron 应用
- 支持自定义图标
- 支持禁用菜单栏
- 支持页面缩放
- 自动生成 Electron 启动脚本

## 安装

### 从源码构建

```bash
git clone https://github.com/Ericwyn/electron-one.git
cd electron-one
go build -o electron-one .
```

### 前置要求

在使用 electron-one 之前，需要先安装 Electron。可以通过以下两种方式配置 Electron 二进制文件路径：

**方式一：使用环境变量**

```bash
# 设置 Electron 二进制文件路径
export ELECTRON_BIN_PATH=/path/to/electron
```

建议将此环境变量添加到你的 shell 配置文件中（如 `~/.bashrc` 或 `~/.zshrc`）。

**方式二：使用命令行参数**

```bash
./electron-one -index index.html -electron /path/to/electron
```

## 使用方法

### 基本用法

```bash
# 查看版本号
./electron-one -v

# 启动应用
./electron-one -index /path/to/your/index.html
```

### 完整参数示例

```bash
./electron-one \
    -index /home/ericwyn/dev/nodejs/mcp-partner-desktop/mcp-partner/dist/index.html \
    -electron /usr/local/bin/electron \
    -icon /home/ericwyn/dev/nodejs/mcp-partner-desktop/mcp-partner/public/icon_512px.png \
    -disable-menu \
    -zoom 1.25
```

### 参数说明

| 参数 | 必需 | 说明 |
|------|------|------|
| `-v` | 否 | 显示版本信息 |
| `-index` | 是 | 启动的 HTML 文件路径 |
| `-electron` | 否 | Electron 二进制文件地址（也可通过 ELECTRON_BIN_PATH 环境变量设置） |
| `-icon` | 否 | 应用图标地址 |
| `-disable-menu` | 否 | 禁用基础菜单栏（启用此选项） |
| `-zoom` | 否 | 页面缩放比例（例如：1.25 表示 125%） |

## 工作原理

1. electron-one 根据提供的参数生成一个 `electron-one.cjs` 脚本文件（位于 index.html 所在目录）
2. 调用 Electron 二进制文件（优先使用 `-electron` 参数，如果没有设置则使用 `ELECTRON_BIN_PATH` 环境变量），传入生成的脚本
3. Electron 应用启动，加载指定的 HTML 文件

生成的 `electron-one.cjs` 脚本包含了窗口配置、图标、菜单和缩放设置。

## 配置示例

生成的 Electron 脚本默认配置：
- 窗口大小：1400x900
- Node Integration：启用
- Context Isolation：禁用

如需修改默认配置，可以编辑 `electron-one.go` 中的 `cjsTemplate` 常量。

## 故障排除

### 错误：需要通过 -electron 参数设置或设置 ELECTRON_BIN_PATH 环境变量

确保已通过以下任一方式配置 Electron 二进制文件路径：

**方式一：使用命令行参数**

```bash
./electron-one -index index.html -electron /path/to/electron
```

**方式二：使用环境变量**

```bash
echo $ELECTRON_BIN_PATH
```

### 错误：index 文件不存在

检查 `-index` 参数指定的文件路径是否正确。

### 错误：electron 二进制文件不存在

检查 `ELECTRON_BIN_PATH` 环境变量指向的文件是否存在。

## 开发

```bash
# 安装依赖
go mod tidy

# 构建
go build -o electron-one .

# 运行
./electron-one -index /path/to/index.html
```

## 许可证

MIT

## 作者

Ericwyn