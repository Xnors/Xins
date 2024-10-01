package console

import (
	"bufio"
	"fmt"
	"github.com/scylladb/termtables"
	"os"
	"strings"
	"xins/console/installer"
)

func Start() {
	// 读取配置文件
	cfgReader, err := ReadConfig("config.json")
	if err != nil {
		fmt.Println(err)
	}

	// 打印头部信息
	t := termtables.CreateTable()
	t.AddRow("Welcome to", "Xins Install Manager!")
	t.AddRow("版本号", cfgReader["version"])
	t.AddRow("作者/工作室", cfgReader["author"].(string)+"/"+cfgReader["studio"].(string))
	t.AddRow("邮箱", cfgReader["email"])
	fmt.Println(t.Render())
	// fmt.Println("-------------------------------------------------------------")
	// fmt.Print("|\tWelcome to \t|  Xins Install Manager!\t|\n")
	// fmt.Println("|\t版本号 \t\t| ", cfgReader["version"], "\t\t\t|")
	// fmt.Println("|\t作者/工作室 \t| ", cfgReader["author"], "/", cfgReader["studio"], "\t|")
	// fmt.Println("|\t邮箱 \t\t| ", cfgReader["email"], "\t|")
	// fmt.Println("-------------------------------------------------------------")

	core()
}

// 核心功能代码
func core() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(">>> ")

	// 读取一行输入
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("读取输入时发生错误:", err)
		return
	}

	// 去除换行符并分割输入
	input = strings.TrimSpace(input)
	parts := strings.Fields(input)

	// 存储指令和参数
	command := parts[0]
	args := parts[1:]

	fmt.Println("你输入了:", command, args)

	// 执行指令
	switch command {
	case "安装", "install", "i", "anzhuang", "az":
		if len(args) == 0 {
			fmt.Println("错误: 你想要安装什么?")
		}
		func_install(args)
		
	default:
		fmt.Println("错误: 未知指令:", command)
	}
}

func func_install(args []string) {
	fmt.Println("开始下载安装包...")

	if len(args) == 1 {
		installer.InstallFilesByUrl(args, "downloads", false)
	} else if len(args) > 1 {
		installer.InstallFilesByUrl(args, "downloads", true)
	}
}
