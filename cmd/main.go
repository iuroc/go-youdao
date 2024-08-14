package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/iuroc/go-youdao"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("🚩 请输入需要翻译的内容: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		fmt.Println()
		result, err := youdao.GetTranslateResult("AUTO", "AUTO", input)
		if err != nil {
			fmt.Println("❌", err)
			PrintLine()
			continue
		}
		fmt.Println("🎉 翻译结果:", result.Target)
		data, err := youdao.GetPronounce(result.Target, result.To)
		if err != nil {
			fmt.Println("❌", err)
			PrintLine()
			continue
		}
		filename := "temp.mp3"
		err = os.WriteFile(filename, data, 0600)
		if err != nil {
			fmt.Println("❌ 写入临时文件失败", err)
			PrintLine()
			continue
		}
		// Windows
		cmd := exec.Command("ffplay", "-nodisp", "-autoexit", filename)
		err = cmd.Start()
		if err != nil {
			// Termux
			cmd = exec.Command("play-audio", filename)
			err = cmd.Start()
			if err != nil {
				fmt.Println("❌ 播放音频失败", err)
				PrintLine()
			}
			continue
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println("❌ 播放音频失败", err)
			PrintLine()
			continue
		}
		err = os.Remove(filename)
		if err != nil {
			fmt.Println("❌ 删除临时文件失败", err)
			PrintLine()
			continue
		}
		PrintLine()
	}
}

func PrintLine() {
	fmt.Println()
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println()
}
