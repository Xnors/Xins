package installer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// downloadFile downloads a file from a given URL to a local path.
func downloadFile(url, filePath string, wg *sync.WaitGroup) {
	defer wg.Done() // 通知WaitGroup当前goroutine已完成

	// 发起HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		return
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error copying to file %s: %v\n", filePath, err)
	}
	fmt.Printf("Downloaded %s\n", url)
}

func InstallFilesByUrl(urls []string, directory string, multithreaded bool) {
	var wg sync.WaitGroup
	os.MkdirAll(directory, os.ModePerm) // 确保下载目录存在

	for _, url := range urls {
		fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
		filePath := filepath.Join(directory, fileName)
		wg.Add(1) // 增加WaitGroup的计数
		if multithreaded {
			go downloadFile(url, filePath, &wg) // 启动一个新的goroutine来下载文件
		} else {
			downloadFile(url, filePath, &wg) // 下载文件
		}
	}

	wg.Wait() // 等待所有goroutine完成
	fmt.Println("All files downloaded")
}

// func main() {
// 	urls := []string{
// 		"https://www.example.com/file1.txt",

// 	}
// 	InstallFilesByUrl(urls, "downloads", true)
// }
