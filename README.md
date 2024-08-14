# go-youdao

Go 语言有道翻译模块，支持朗读。

## 获取模块

```shell
go get -u github.com/iuroc/go-youdao
```

## 获取翻译结果

`from` 和 `to` 是翻译语言组合，[取值说明](https://github.com/iuroc/apee-fanyi?tab=readme-ov-file#%E8%AF%AD%E8%A8%80%E5%88%97%E8%A1%A8)。

```go
input := "这是需要翻译的文本"
from := "AUTO"  // 翻译前的语言类型
to := "AUTO"    // 翻译后的语言类型
result, _ := youdao.GetTranslateResult(from, to, input)
fmt.Println("翻译结果:", result.Target)
```

### 翻译结果字段说明

- `Target`: 翻译结果文本
- `TargetPronounce`: 翻译结果的读音
- `To`: 翻译结果的语言类型
- `From`: 翻译前的语言类型
- `Src`: 翻译前的文本

## 获取朗读音频

```go
input := "这是需要朗读的文本"
language := "zh-CHS"  // 需要朗读的语言类型
data, err := youdao.GetPronounce(input, language)

// 将音频保存到文件
os.WriteFile("temp.mp3", data, 0600)
```

## CLI 工具

这里已经写好了一个 CLI 程序，可以将输入的文本进行翻译，并自动朗读翻译的结果。

你可以[点击这里查看源代码](./cmd/main.go)，也可以在本仓库的 [Releases](./releases) 下载编译后的版本。
