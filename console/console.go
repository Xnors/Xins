package console

import (
	"bufio"
	"fmt"
	"github.com/scylladb/termtables"
	"os"
	"os/exec"
	"strings"
	"xins/console/installer"
)

func isInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

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
	for {
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

		// fmt.Println("你输入了:", command, args)

		// 执行指令
		switch command {
		case "安装", "install", "i", "anzhuang", "az":
			if len(args) == 0 {
				fmt.Print("错误: 请输入你想要安装的软件名称!\n示例:\ti python\n")
			} else if len(args) == 1 {
				func_install(args[0])
			} else {
				fmt.Println("错误: 输入参数过多!")
			}

		case "q", "quit", "exit", "退出", "退出程序":
			fmt.Println("欢迎下次使用!")
			return
		default:
			fmt.Println("未知指令:", command)
		}
	}
}

func func_install(arg string) {
	var file string
	var urlToDownload string
	var OSBIT string
	var version string

	// 判断是否为64位系统
	if os.Getenv("ProgramFiles(x86)") != "" {
		OSBIT = "64bit"
	} else {
		OSBIT = "32bit"
	}

	// 读取镜像列表
	mirrors, err := ReadMirrors("mirrors.json")
	if err != nil {
		fmt.Println(err)
	}

	// 判断软件是否存在
	//// 创建一个切片来保存第一层的键
	var keys []string
	for key := range mirrors {
		keys = append(keys, key)
	}
	//// 判断软件是否存在
	if !isInSlice(keys, strings.ToLower(arg)) {
		fmt.Print("非常抱歉, 我们还没有收录这个软件!\n您可以前往 xnors-studio@outlook.com 联系我们添加软件!\n")
		return
	}

	fmt.Println("======[ 请输入版本号(必须输入完整的版本号哦) ]====== ")

	var versions []string
	for key := range mirrors[strings.ToLower(arg)] {
		fmt.Println("\t\t  ", key)
		versions = append(versions, key)
	}

	for {
		fmt.Print(">>> ")
		fmt.Scanf("%s\n", &version)
		if !isInSlice(versions, version) {
			fmt.Print("输入的版本号不存在, 请重新输入!\n")
			continue
		} else {
			break
		}
	}

	fmt.Printf("准备下载: %s 版本的 %s 安装包 ...\n", version, arg)

	// 获取下载链接
	urlToDownload = mirrors[strings.ToLower(arg)][version][OSBIT]["url"]

	fmt.Println(urlToDownload)

	fmt.Println("开始下载安装包...")
	file = installer.InstallOneFileByUrl(urlToDownload, "downloads")

	fmt.Print("\n======[ 下载完成! ]======\n\n")
	fmt.Print("==========[以下是安装教程]===========\n\n")
	fmt.Print(mirrors[strings.ToLower(arg)][version][OSBIT]["installDoc"])
	fmt.Print("==========[以上是安装教程]===========\n\n")
	// 运行安装包
	exec.Command(`.\` + file).Run()
}
