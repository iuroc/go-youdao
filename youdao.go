package youdao

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TranslateResult struct {
	// 翻译结果
	Target string `json:"tgt"`
	// 翻译结果的读音
	TargetPronounce string `json:"tgtPronounce"`
	// 翻译前的内容
	Src  string `json:"src"`
	From string
	To   string
}

// GetTranslateResult 获取翻译结果。
//
// from 和 to 取值说明：请查阅完整语言列表 https://github.com/iuroc/apee-fanyi?tab=readme-ov-file#%E8%AF%AD%E8%A8%80%E5%88%97%E8%A1%A8
func GetTranslateResult(from string, to string, input string) (*TranslateResult, error) {
	mysticTime := time.Now().UnixMilli()
	signStr := fmt.Sprintf("client=fanyideskweb&mysticTime=%d&product=webfanyi&key=fsdsogkndfokasodnaso", mysticTime)
	sign := md5Hex(signStr)
	body := url.Values{
		"i":          {input},
		"from":       {from},
		"to":         {to},
		"keyid":      {"webfanyi"},
		"sign":       {sign},
		"client":     {"fanyideskweb"},
		"product":    {"webfanyi"},
		"appVersion": {"1.0.0"},
		"vendor":     {"web"},
		"pointParam": {"client,mysticTime,product"},
		"mysticTime": {strconv.FormatInt(mysticTime, 10)},
		"keyfrom":    {"fanyi.web"},
	}
	request, err := http.NewRequest("POST", "https://dict.youdao.com/webtranslate", strings.NewReader(body.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Referer", "https://fanyi.youdao.com/")
	request.Header.Set("Cookie", "OUTFOX_SEARCH_USER_ID=0@0.0.0.0")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer response.Body.Close()
	bs, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	responseStr := string(bs)
	data, err := base64.URLEncoding.DecodeString(string(responseStr))
	if err != nil {
		return nil, fmt.Errorf("Base64 解码失败: %v", err)
	}
	result, err := decrypt(data)
	if err != nil {
		return nil, fmt.Errorf("AES 解密失败: %v", err)
	}
	res := &struct {
		// 状态码，正常为 0
		Code   int                 `json:"code"`
		Result [][]TranslateResult `json:"translateResult"`
		// 语言组合类型，完整列表参阅
		// https://github.com/iuroc/apee-fanyi?tab=readme-ov-file#%E8%AF%AD%E8%A8%80%E5%88%97%E8%A1%A8
		Type string `json:"type"`
	}{}
	err = json.Unmarshal(result, res)
	if err != nil || res.Code != 0 || len(res.Result) == 0 || len(res.Result[0]) == 0 {
		log.Fatalln("获取翻译结果失败:", err)
	}
	res.Result[0][0].From = strings.Split(res.Type, "2")[0]
	res.Result[0][0].To = strings.Split(res.Type, "2")[1]
	return &res.Result[0][0], nil
}

// GetPronounce 获取文本的朗读音频。
func GetPronounce(text string, language string) (data []byte, err error) {
	mysticTime := time.Now().UnixMilli()
	body := url.Values{
		"client":     {"web"},
		"keyfrom":    {"webfanyi"},
		"keyid":      {"voiceFanyiWeb"},
		"le":         {language},
		"mysticTime": {strconv.FormatInt(mysticTime, 10)},
		"product":    {"webfanyi"},
		"vendor":     {"web"},
		"word":       {text},
	}
	keys := make([]string, 0)
	for key := range body {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	body.Set("key", "qCG2vdP92hOXDcKa")
	keys = append(keys, "key")
	parts := make([]string, len(keys))
	for index, key := range keys {
		parts[index] = key + "=" + body.Get(key)
	}
	signStr := strings.Join(parts, "&")
	sign := md5Hex(signStr)
	body.Set("sign", sign)
	body.Set("pointParam", strings.Join(keys, ","))
	body.Del("key")
	response, err := http.Get("https://dict.youdao.com/pronounce/base?" + body.Encode())
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer response.Body.Close()
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	var jsonData struct {
		Status  int    `json:"status"`
		Message string `json:"msg"`
	}
	err = json.Unmarshal(data, &jsonData)
	if err == nil {
		return nil, fmt.Errorf("获取音频失败: %s", jsonData.Message)
	}
	return data, nil
}

// decrypt 实现 AES-128-CBC 解密。
func decrypt(ciphertext []byte) ([]byte, error) {
	key := md5.Sum([]byte("ydsecret://query/key/B*RGygVywfNBwpmBaZg*WT7SIOUP2T0C9WHMZN39j^DAdaZhAnxvGcCY6VYFwnHl"))
	iv := md5.Sum([]byte("ydsecret://query/iv/C@lZe2YzHtZ2CYgaXKSVfsb7Y4QWHjITPPZ0nQp87fBeJ!Iv6v^6fvi2WN@bYpJ4"))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("[aes.NewCipher] %v", err)
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("数据错误")
	}
	mode := cipher.NewCBCDecrypter(block, iv[:])
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	length := len(plaintext)
	paddingLen := int(plaintext[length-1])
	return plaintext[:length-paddingLen], nil
}

// md5Hex 生成 Hex 格式的 MD5。
func md5Hex(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
