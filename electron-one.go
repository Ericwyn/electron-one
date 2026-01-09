package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	Version = "0.1.2"
	cjsTemplate = `const { app, BrowserWindow } = require('electron')
const path = require('path')

function createWindow () {
    // 1. 在这里配置窗口的大小和图标
    const win = new BrowserWindow({
        width: 1400,                  // 设置宽度
        height: 900,                  // 设置高度
        // resizable: false,          // 如果不想让用户调整大小，取消注释这行
        {{if .IconPath}}icon: "{{.IconPath}}", // 设置图标 (确保目录下有 icon.png){{end}}
        webPreferences: {
            nodeIntegration: true,      // 根据你的需求开启/关闭
            contextIsolation: false
        }
    })

    // 2. 加载你的 index.html
    win.loadFile('{{.IndexPath}}')

    {{if .DebugMode}}// 调试模式：自动打开开发者工具
    win.webContents.openDevTools()
	{{end}}

    {{if .DisableMenu}}// 可选：启动时移除默认菜单栏
    win.setMenu(null)
	{{end}}

    {{if .ZoomFactor}}// 在加载文件后设置缩放
    // 注意：有时为了确保生效，建议放在 'dom-ready' 事件里
    win.webContents.on('dom-ready', () => {
        win.webContents.setZoomFactor({{.ZoomFactor}}) // 设置为 {{.Percentage}}%
    })
	{{end}}
}

app.whenReady().then(() => {
    createWindow()

    app.on('activate', () => {
        if (BrowserWindow.getAllWindows().length === 0) {
            createWindow()
        }
    })
})

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit()
    }
})
`
)

func main() {
	showVersion := flag.Bool("v", false, "显示版本信息")
	indexPath := flag.String("index", "", "启动的 html 文件 (必需)")
	iconPath := flag.String("icon", "", "图标地址 (可选)")
	disableMenu := flag.Bool("disable-menu", false, "禁用基础菜单")
	zoomFactor := flag.Float64("zoom", 0, "缩放比例 (可选)")
	debugMode := flag.Bool("debug", false, "启用调试模式 (自动打开开发者工具)")
	electronBin := flag.String("electron", "", "electron 二进制文件地址 (可选，也可通过 ELECTRON_BIN_PATH 环境变量设置)")

	flag.Parse()

	if *showVersion {
		electronBinPath := *electronBin
		if electronBinPath == "" {
			electronBinPath = os.Getenv("ELECTRON_BIN_PATH")
		}
		var electronVersion string
		if electronBinPath != "" {
			cmd := exec.Command(electronBinPath, "--version")
			output, err := cmd.Output()
			if err == nil {
				electronVersion = strings.TrimSpace(string(output))
			}
		}

		if electronVersion != "" {
			fmt.Printf("electron-one version %s (Electron %s)\n", Version, electronVersion)
		} else {
			fmt.Printf("electron-one version %s\n", Version)
		}
		os.Exit(0)
	}

	if *indexPath == "" {
		fmt.Fprintln(os.Stderr, "错误: -index 参数不能为空")
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*indexPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "错误: index 文件不存在: %s\n", *indexPath)
		os.Exit(1)
	}

	if *iconPath != "" {
		if _, err := os.Stat(*iconPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "警告: 图标文件不存在: %s\n", *iconPath)
		}
	}

	electronBinPath := *electronBin
	if electronBinPath == "" {
		electronBinPath = os.Getenv("ELECTRON_BIN_PATH")
	}
	if electronBinPath == "" {
		fmt.Fprintln(os.Stderr, "错误: 需要通过 -electron 参数设置或设置 ELECTRON_BIN_PATH 环境变量")
		os.Exit(1)
	}

	if _, err := os.Stat(electronBinPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "错误: electron 二进制文件不存在: %s\n", electronBinPath)
		os.Exit(1)
	}

	indexDir := filepath.Dir(*indexPath)
	cjsFilePath := filepath.Join(indexDir, "electron-one.cjs")

	var percentage string
	if *zoomFactor > 0 {
		percentage = fmt.Sprintf("%.0f", *zoomFactor*100)
	}

	templateData := struct {
		IndexPath     string
		IconPath      string
		DisableMenu   bool
		ZoomFactor    float64
		Percentage    string
		DebugMode     bool
	}{
		IndexPath:     *indexPath,
		IconPath:      *iconPath,
		DisableMenu:   *disableMenu,
		ZoomFactor:    *zoomFactor,
		Percentage:    percentage,
		DebugMode:     *debugMode,
	}

	tmpl, err := template.New("electron").Parse(cjsTemplate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 解析模板失败: %v\n", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 生成脚本失败: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(cjsFilePath, buf.Bytes(), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 写入脚本文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("已生成 electron-one.cjs 文件: %s\n", cjsFilePath)

	cmd := exec.Command(electronBinPath, cjsFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("正在启动 electron: %s\n", electronBinPath)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 启动 electron 失败: %v\n", err)
		os.Exit(1)
	}
}
