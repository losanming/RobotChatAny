package banner

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// 颜色常量
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	Bold   = "\033[1m"
)

var GenShin = "   _____             _____ _     _       _ \n  / ____|           / ____| |   (_)     | |\n | |  __  ___ _ __ | (___ | |__  _ _ __ | |\n | | |_ |/ _ \\ '_ \\ \\___ \\| '_ \\| | '_ \\| |\n | |__| |  __/ | | |____) | | | | | | | |_|\n  \\_____|\\___|_| |_|_____/|_| |_|_|_| |_(_)\n                                           \n                                           "

// PrintBanner 打印启动横幅
func PrintBanner(appName string, version string) {
	// 大型艺术字体
	banner := []string{
		"  _   _                 ____        _           _      ____        _  __",
		" | \\ | | _____      __ |  _ \\ _   _| |__   ___ | |_   / ___|  ___| |/ _|",
		" |  \\| |/ _ \\ \\ /\\ / / | |_) | | | | '_ \\ / _ \\| __| | |  _  / _ \\ | |_ ",
		" | |\\  |  __/\\ V  V /  |  _ <| |_| | |_) | (_) | |_  | |_| |  __/ |  _|",
		" |_| \\_|\\___| \\_/\\_/   |_| \\_\\\\__,_|_.__/ \\___/ \\__|  \\____|\\___|_|_|  ",
		"                                                                        ",
	}

	// 系统信息
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	goVersion := runtime.Version()
	osInfo := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
	cpuInfo := fmt.Sprintf("%d CPUs", runtime.NumCPU())

	// 打印横幅
	fmt.Println()
	for _, line := range banner {
		fmt.Printf("%s%s%s\n", Cyan, line, Reset)
	}
	fmt.Println()

	// 打印应用信息
	fmt.Printf("%s%s应用名称%s: %s%s%s\n", Bold, Yellow, Reset, Green, appName, Reset)
	fmt.Printf("%s%s版本号%s: %s%s%s\n", Bold, Yellow, Reset, Green, version, Reset)
	fmt.Printf("%s%s启动时间%s: %s%s%s\n", Bold, Yellow, Reset, Green, currentTime, Reset)
	fmt.Printf("%s%sGo版本%s: %s%s%s\n", Bold, Yellow, Reset, Green, goVersion, Reset)
	fmt.Printf("%s%s系统信息%s: %s%s%s\n", Bold, Yellow, Reset, Green, osInfo, Reset)
	fmt.Printf("%s%sCPU信息%s: %s%s%s\n", Bold, Yellow, Reset, Green, cpuInfo, Reset)
	fmt.Println()

	// 分隔线
	fmt.Printf("%s%s%s\n", Purple, strings.Repeat("=", 70), Reset)
	fmt.Println()
}

// PrintCustomBanner 打印自定义横幅
func PrintCustomBanner(text string, version string) {
	// 将文本转换为ASCII艺术字体
	asciiArt := generateASCIIArt(text)

	// 系统信息
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	goVersion := runtime.Version()
	osInfo := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	// 打印横幅
	fmt.Println()
	for _, line := range asciiArt {
		fmt.Printf("%s%s%s\n", Cyan, line, Reset)
	}
	fmt.Println()

	// 打印应用信息
	fmt.Printf("%s%s版本号%s: %s%s%s\n", Bold, Yellow, Reset, Green, version, Reset)
	fmt.Printf("%s%s启动时间%s: %s%s%s\n", Bold, Yellow, Reset, Green, currentTime, Reset)
	fmt.Printf("%s%sGo版本%s: %s%s%s\n", Bold, Yellow, Reset, Green, goVersion, Reset)
	fmt.Printf("%s%s系统信息%s: %s%s%s\n", Bold, Yellow, Reset, Green, osInfo, Reset)
	fmt.Println()

	// 分隔线
	fmt.Printf("%s%s%s\n", Purple, strings.Repeat("=", 70), Reset)
	fmt.Println()
}

// generateASCIIArt 生成简单的ASCII艺术字体
func generateASCIIArt(text string) []string {
	// 为中文和英文文本创建更好的显示效果

	// 预定义的大字体模板
	fonts := map[string][]string{
		"机": {
			"  ######  ",
			"     #    ",
			"     #    ",
			"##########",
			"     #    ",
			"     #    ",
			"##########",
		},
		"器": {
			"##########",
			"#   ##   #",
			"##########",
			"#   ##   #",
			"##########",
			"    ##    ",
			"##########",
		},
		"人": {
			"    ##    ",
			"    ##    ",
			"   #  #   ",
			"   #  #   ",
			"  #    #  ",
			" #      # ",
			"#        #",
		},
		"新": {
			"#      ###",
			" #    #   ",
			"  #  #    ",
			"#########",
			"    #     ",
			"   # #    ",
			"###   ####",
		},
		"自": {
			"##########",
			"    ##    ",
			"    ##    ",
			"##########",
			"    ##    ",
			"    ##    ",
			"    ##    ",
		},
	}

	// 初始化结果数组
	result := make([]string, 8)
	for i := range result {
		result[i] = ""
	}

	// 处理每个字符
	for _, char := range text {
		charStr := string(char)

		// 检查是否有预定义的字体
		if font, exists := fonts[charStr]; exists {
			// 使用预定义字体
			for i, line := range font {
				result[i] += line + " "
			}
		} else if char == ' ' {
			// 处理空格
			for i := range result {
				result[i] += "          "
			}
		} else {
			// 对于其他字符，创建简单的方块字体
			result[0] += "##########" + " "
			result[1] += "#        #" + " "
			result[2] += "#        #" + " "
			result[3] += fmt.Sprintf("#   %c    #", char) + " "
			result[4] += "#        #" + " "
			result[5] += "#        #" + " "
			result[6] += "##########" + " "
			result[7] += "          " + " "
		}
	}

	return result
}
