package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

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
		result, err := youdao.GetTranslateResult("AUTO", "AUTO", input)
		if err != nil {
			fmt.Println("❌", err)
			continue
		}
		fmt.Println()
		fmt.Println("🎉 翻译结果:", result.Target)
		fmt.Println()
		fmt.Println()
		data, err := youdao.GetPronounce(result.Target, result.To)
		if err != nil {
			fmt.Println("❌", err)
			continue
		}
		filename := "temp.mp3"
		err = os.WriteFile(filename, data, 0600)
		if err != nil {
			fmt.Println("❌ 写入临时文件失败", err)
			continue
		}
		cmd := exec.Command("ffplay", "-nodisp", "-autoexit", filename)
		err = cmd.Start()
		if err != nil {
			fmt.Println("❌ 播放音频失败", err)
			continue
		}
		err = cmd.Wait()
		if err != nil {
			fmt.Println("❌ 播放音频失败", err)
			continue
		}
		err = os.Remove(filename)
		if err != nil {
			fmt.Println("❌ 删除临时文件失败", err)
			continue
		}
	}
}
