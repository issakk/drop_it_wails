package main

import (
	"context"
	"dropit/model"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"  // 添加这行
	"strings"
	"sync"     // 添加这行
	"sync/atomic"  // 添加这行

	"github.com/samber/lo"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"  // 重命名避免冲突
)

var fileMap = map[string][]string{
    "图片": {
        "png", "jpg", "jpeg", "gif", "bmp", "tiff", "svg", "ico", "webp", "psd", "raw", "heic", 
        "cr2", "nef", "arw", "dng", "raf", "ai", "eps", "xcf", "sketch", "fig", "jfif", "avif",
        "jp2", "jxr", "hdp", "wdp", "exr", "dds",
    },
    "电子书": {
        "pdf", "mobi", "azw3", "epub", "djvu", "cbz", "cbr", "txt", "chm", "azw", "kfx", 
        "fb2", "lit", "prc", "tpz", "opf", "htmlz", "lrf", "tcr", "pdb", "caj", "pdg",
    },
    "文档资料": {
        "txt", "md", "doc", "docx", "xlsx", "csv", "ppt", "pptx", "rtf", "odt", "ods", "odp", 
        "pages", "numbers", "key", "one", "xls", "xmind", "mindnode", "wps", "et", "dps", "vsd",
        "vsdx", "mpp", "tex", "docm", "xlsm", "pptm", "dot", "dotx", "odf", "ott", "ots", "otp",
        "pub", "wpd", "oxps", "xps",
    },
    "压缩包": {
        "zip", "rar", "7z", "tar", "gz", "bz2", "xz", "iso", "pkg", "tgz", "tbz", "tbz2", 
        "txz", "cab", "ace", "lzh", "lha", "arj", "uue", "bz", "rz", "war", "ear", "sitx",
    },
    "音频": {
        "mp3", "wav", "m4a", "flac", "wav", "ogg", "aac", "m4b", "ape", "aiff", "wma", "opus",
        "mid", "midi", "amr", "ac3", "dsf", "dff", "mka", "ra", "rm", "tta", "wv", "mpc",
        "spx", "voc", "3ga", "m4r", "alac",
    },
    "视频": {
        "mp4", "mkv", "mov", "flv", "wmv", "rmvb", "avi", "m4v", "webm", "3gp", "ts", "f4v",
        "vob", "ogv", "mts", "m2ts", "mpg", "mpeg", "m1v", "m2v", "asf", "swf", "divx", "xvid",
        "mxf", "roq", "nsv", "mpv", "yuv", "rm", "qt", "amv",
    },
    "exe程序": {
        "exe", "sh", "bat", "py", "jar", "app", "dmg", "msi", "apk", "ipa", "deb", "rpm", "run",
        "bin", "com", "vbs", "js", "cmd", "ps1", "psm1", "psd1", "msp", "msix", "appx", "appxbundle",
        "air", "pkg", "flatpak", "snap", "AppImage",
    },
    "编程相关": {
        "out", "log", "json", "yml", "yaml", "xml", "sql", "db", "sqlite", "ini", "conf", "config",
        "properties", "env", "toml", "go", "cpp", "c", "h", "hpp", "cs", "java", "kt", "rs",
        "swift", "m", "mm", "rb", "php", "pl", "lua", "r", "dart", "ts", "jsx", "tsx", "vue",
        "css", "scss", "less", "html", "htm", "xhtml", "jsp", "asp", "aspx", "erb", "haml",
        "slim", "phtml", "gradle", "pom", "sln", "csproj", "vcxproj", "pbxproj", "xcodeproj",
        "gitignore", "dockerignore", "makefile", "cmake", "pro", "pri",
    },
}
// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
} /**/

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
func (a *App) OpenFileDialog() []model.TreeNode {
	fmt.Println("openDialog")
	path, _ := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{})
	//dir, _ := os.ReadDir(path)
	//dirList := lo.Filter(dir, func(item os.DirEntry, index int) bool {
	//	info, _ := item.Info()
	//	return info.IsDir()
	//})
	//return lo.Map(dirList, func(item os.DirEntry, index int) model.TreeNode {
	//	info, _ := item.Info()
	//	filePath := filepath.Join(path, info.Name())
	//	return model.TreeNode{Value: filePath, Label: info.Name()}
	//})

	node := NewNode(path)

	return []model.TreeNode{node}
}

func NewNode(path string) model.TreeNode {
	if path == "" {
		return model.TreeNode{}
	}
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	node := model.TreeNode{
		Value: path,
		Label: filepath.Base(path),
	}

	if info.IsDir() {
		childrenInfo, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		for _, childFileInfo := range childrenInfo {
			childPath := filepath.Join(path, childFileInfo.Name())
			if childFileInfo.IsDir() {
				node.Children = append(node.Children, NewNode(childPath))
			}
		}
	}

	return node
}

func (a *App) ListFileInfo(path string) []model.FileInfo {
	dir, _ := os.ReadDir(path)
	fileList := lo.Filter(dir, func(item os.DirEntry, index int) bool {
		info, _ := item.Info()
		return !info.IsDir()
	})
	return lo.Map(fileList, func(item os.DirEntry, index int) model.FileInfo {
		info, _ := item.Info()
		return model.FileInfo{Name: info.Name(), Size: info.Size(), Date: info.ModTime().Format("2006-01-02 15:04:05")}
	})
}
func (a *App) Drop(path string) (int64, error) {
    files := a.ListFileInfo(path)  // 移除 err 变量，因为 ListFileInfo 只返回一个值
    
    var wg sync.WaitGroup
    var processedCount int64
    
    // 创建备份目录
    backupPath := filepath.Join(path, "备份")
    if err := os.MkdirAll(backupPath, 0755); err != nil {
        return 0, fmt.Errorf("创建备份目录失败: %w", err)
    }
    
    semaphore := make(chan struct{}, runtime.NumCPU())
    
    for _, file := range files {
        wg.Add(1)
        semaphore <- struct{}{} // 获取信号量
        
        go func(fileInfo model.FileInfo) {
            defer func() {
                <-semaphore // 释放信号量
                wg.Done()
            }()
            // 获取文件扩展名并创建目标目录
            // 获取文件扩展名
            ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileInfo.Name), "."))
            if ext == "" {
                return
            }
            
            // 查找文件类别
            var category string
            for cat, extensions := range fileMap {
                if lo.Contains(extensions, ext) {
                    category = cat
                    break
                }
            }
            
            // 如果没找到对应分类则跳过
            if category == "" {
                return
            }
            
            // 使用分类名创建目标目录
            targetDir := filepath.Join(path, category)
            if err := os.MkdirAll(targetDir, 0755); err != nil {
                return
            }
            // 先复制到分类目录
            srcPath := filepath.Join(path, fileInfo.Name)
            if err := copyFileToDir(srcPath, targetDir); err != nil {
                return
            }
            
            // 再移动到备份目录
            if err := os.Rename(srcPath, filepath.Join(backupPath, fileInfo.Name)); err != nil {
                return
            }

            atomic.AddInt64(&processedCount, 1)
            progress := float64(atomic.LoadInt64(&processedCount)) / float64(len(files)) * 100
            wailsRuntime.EventsEmit(a.ctx, "file-progress", progress)  // 使用重命名后的包
            
        }(file)
    }

    wg.Wait()
    return atomic.LoadInt64(&processedCount), nil
}
func copyFiles(m map[string][]fs.FileInfo, path string) (int, error) {
	count := 0
	backupPath := filepath.Join(path, "备份")

	if err := ensureDir(backupPath); err != nil {
		return 0, fmt.Errorf("创建备份目录失败: %w", err)
	}

	for category, files := range m {
		categoryPath := filepath.Join(path, category)
		if err := ensureDir(categoryPath); err != nil {
			return 0, fmt.Errorf("创建分类目录失败: %w", err)
		}

		for _, fileInfo := range files {
			if err := moveFile(path, categoryPath, backupPath, fileInfo); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}

func moveFile(basePath, categoryPath, backupPath string, fileInfo fs.FileInfo) error {
	fileName := fileInfo.Name()
	oldPath := filepath.Join(basePath, fileName)

	// 先复制到分类目录
	if err := copyFileToDir(oldPath, categoryPath); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	// 再移动到备份目录
	if err := os.Rename(oldPath, filepath.Join(backupPath, fileName)); err != nil {
		return fmt.Errorf("移动到备份目录失败: %w", err)
	}

	return nil
}

func copyFileToDir(src, dstDir string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	dst := filepath.Join(dstDir, filepath.Base(src))
	return ioutil.WriteFile(dst, data, os.ModePerm)
}

func ensureDir(path string) error {
	if exists, err := pathExists(path); err != nil {
		return err
	} else if !exists {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func readFiles(path string) (map[string][]fs.FileInfo, error) {
	m := make(map[string][]fs.FileInfo)

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("读取目录失败: %w", err)
	}

	for _, fileInfo := range dir {
		if fileInfo.IsDir() {
			continue
		}

		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(fileInfo.Name()), "."))

		for category, extensions := range fileMap {
			if lo.Contains(extensions, ext) {
				m[category] = append(m[category], fileInfo)
				break
			}
		}
	}
	return m, nil
}
