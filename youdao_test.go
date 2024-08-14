package youdao

import (
	"os"
	"os/exec"
	"testing"
)

func TestGetTranslate(t *testing.T) {
	t.Log("🚩 正在翻译")
	result, err := GetTranslateResult("AUTO", "AUTO", "你好吗，亲爱的读者。")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("🚩 完成翻译")
	t.Logf("🚩 翻译结果: %s", result.Target)
}

func TestGetPronounce(t *testing.T) {
	t.Log("🚩 正在翻译")
	result, err := GetTranslateResult("AUTO", "AUTO", "hello world")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("🚩 完成翻译")
	t.Log("🚩 正在获取音频")
	data, err := GetPronounce(result.Target, "zh-CHS")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("🚩 完成获取音频")
	tempFileName := "temp.mp3"
	err = os.WriteFile(tempFileName, data, 0600)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("🚩 准备播放音频")
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit", tempFileName)
	err = cmd.Start()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("🚩 正在播放")
	err = cmd.Wait()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("🚩 播放完成")
	err = os.Remove(tempFileName)
	if err != nil {
		t.Fatal(err)
	}
}
