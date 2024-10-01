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

// downloadFile downloads a file from a given URL to a local path and shows a progress bar.
func downloadFile(url, filePath string, wg *sync.WaitGroup) {
	defer wg.Done() // 通知WaitGroup当前goroutine已完成
	fmt.Printf("正在下载 \t | %s\n", url)

	// 发起HTTP请求
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// 获取文件总大小
	totalSize := resp.ContentLength

	// 创建文件
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		return
	}
	defer out.Close()

	// 创建进度条
	progressBar := &ProgressBar{
		total: totalSize,
	}

	// 读取并写入文件
	buf := make([]byte, 1024*4) // 4KB的缓冲区
	written := 0
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}
		if n == 0 {
			break
		}

		// 写入文件
		wn, err := out.Write(buf[:n])
		if err != nil {
			fmt.Printf("Error writing to file %s: %v\n", filePath, err)
			return
		}
		written += wn

		// 更新进度条
		progressBar.Update(written)
	}

	fmt.Printf("\n已下载 \t\t | %s\n", url)
}

// ProgressBar is a simple struct to represent a progress bar
type ProgressBar struct {
	total int64
}

// Update updates the progress bar based on the number of bytes written
func (pb *ProgressBar) Update(written int) {
	percentage := float64(written) / float64(pb.total) * 100
	barLength := 50
	barFilled := int(percentage / 100 * float64(barLength))
	bar := strings.Repeat("=", barFilled) + strings.Repeat(" ", barLength-barFilled)
	fmt.Printf("\r[%-50s] %3d%%", bar, int(percentage))
}

func InstallFilesByUrl(urls []string, directory string, multithreaded bool) ([]string) {
	var wg sync.WaitGroup
	var files []string

	os.MkdirAll(directory, os.ModePerm) // 确保下载目录存在

	for _, url := range urls {
		fileName := strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
		filePath := filepath.Join(directory, fileName)
		files = append(files, filePath)
		wg.Add(1) // 增加WaitGroup的计数
		if multithreaded {
			go downloadFile(url, filePath, &wg) // 启动一个新的goroutine来下载文件
		} else {
			downloadFile(url, filePath, &wg) // 下载文件
		}
	}

	wg.Wait() // 等待所有goroutine完成
	fmt.Println("\n全部文件下载完成!")

	return files
}
