package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// 日志记录中间件
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s %s", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
	})
}

// 检查文件后缀是否在允许的列表中
func isAllowedFile(filename string) bool {
	allowedExt := map[string]bool{
		".html": true,
		".js":   true,
		".css":  true,
		".svg":  true,
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".json": true,
		".exe":  true,
		".dmg":  true,
		".apk":  true,
	}

	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExt[ext]
}

// 安全检查中间件
func securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 阻止目录遍历
		if strings.HasSuffix(r.URL.Path, "/") && r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// 检查文件后缀
		if filepath.Ext(r.URL.Path) != "" && !isAllowedFile(r.URL.Path) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 处理下载文件的处理器
func downloadHandler(dir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := path.Join(dir, strings.TrimPrefix(r.URL.Path, "/downloads/"))

		// 检查文件是否存在
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		// 设置适当的头部
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(filePath)))
		http.ServeFile(w, r, filePath)
	})
}

func main() {
	port := ":8080"

	// 确保目录存在
	os.MkdirAll("downloads", 0755)
	os.MkdirAll("images", 0755)

	// 静态文件服务
	staticFiles := http.FileServer(http.Dir("."))
	http.Handle("/", loggingMiddleware(securityMiddleware(staticFiles)))

	// 下载目录
	downloads := downloadHandler("downloads")
	http.Handle("/downloads/", loggingMiddleware(securityMiddleware(http.StripPrefix("/downloads/", downloads))))

	// 图片目录
	images := http.FileServer(http.Dir("images"))
	http.Handle("/images/", loggingMiddleware(securityMiddleware(http.StripPrefix("/images/", images))))

	// 启动服务器
	log.Printf("华创密信下载中心服务器启动在 http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("服务器启动失败: ", err)
	}
}
