package replace

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf16"
)

var (
	ex, _    = os.Executable()
	currPath = filepath.Dir(ex) + "/docu/"
)

func RepAndWrite() {
	config, err := readConfig()
	if err != nil {
		return
	}
	jsonMap, err := readJson(config.JsonPath)
	if err != nil {
		return
	}
	txt, err := readTxt(config.TxtPath)
	if err != nil {
		return
	}
	start := time.Now()
	changeByteCount := 0
	beforeByteCount := len(txt)
	txt = strings.Replace(txt, "All", "所有", 1)
	ProcessBar.Max = float64(len(jsonMap))
	count := 0
	for en, zh := range jsonMap {
		changeByteCount += len(zh) - len(en)
		txt = strings.Replace(txt, "\r\n"+en+"\r\n", "\r\n"+zh+"\r\n", 1)
		count++
		ProcessBar.SetValue(float64(count))
	}
	Log(fmt.Sprintf("成功替换所有文本, 期望字节数变化:%d, 实际字节数变化:%d", changeByteCount, len(txt)-beforeByteCount))
	Log("开始写入文本到 " + config.NewTxtPath)
	//将替换后的文本写入到新文件中
	newFile, err := os.Create(currPath + config.NewTxtPath)
	if err != nil {
		LogErr("创建文件失败", err)
		return
	}
	defer newFile.Close()

	//以utf16写入
	utf16txt := utf16.Encode([]rune(txt))
	writer := bufio.NewWriter(newFile)
	err = binary.Write(writer, binary.LittleEndian, &utf16txt)
	if err != nil {
		LogErr("写入文件失败", err)
		return
	}
	cost := time.Since(start)
	Log("写入完毕!!! 本次运行耗时: " + cost.String())
}

type Config struct {
	JsonPath   string `yaml:"json-path"`
	TxtPath    string `yaml:"txt-path"`
	NewTxtPath string `yaml:"new-txt-path"`
}

func readConfig() (*Config, error) {
	file, err := os.Open(currPath + "config.yml")
	if err != nil {
		LogErr("打开 config.yml 失败", err)
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		LogErr("读取 config.yml 失败, 可能是格式问题", err)
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		LogErr("读取 config.yml 失败, 可能是格式问题", err)
		return nil, err
	}
	Log("成功读取配置 Json文件: " + config.JsonPath)
	Log("成功读取配置 英文文件: " + config.TxtPath)
	Log("成功读取配置 汉化文件: " + config.NewTxtPath)
	return &config, nil
}

func readJson(jsonPath string) (map[string]string, error) {
	jsonFile, err := os.Open(currPath + jsonPath)
	if err != nil {
		LogErr("文件打开失败", err)
		return nil, err
	}
	defer jsonFile.Close()
	byteSlice, _ := io.ReadAll(jsonFile)
	jsonMap := make(map[string]string, 400000)
	err = json.Unmarshal(byteSlice, &jsonMap)
	if err != nil {
		LogErr("json文件反序列化错误, 可能是格式问题", err)
		return nil, err
	}
	Log(fmt.Sprintf("成功读取文件 %v 有效键值对数: %d", jsonPath, len(jsonMap)))
	return jsonMap, nil
}

func readTxt(txtPath string) (string, error) {
	//打开要替换的txt文件
	file, err := os.Open(currPath + txtPath)
	if err != nil {
		LogErr("打开txt文件失败", err)
		return "", err
	}
	defer file.Close()
	//读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		LogErr("读取txt文件失败", err)
	}
	Log(fmt.Sprintf("成功读取文件 %v 有效字节数(英文数字1字节, 汉字3字节): %d", txtPath, len(content)))
	//将文件内容转换为字符串//utf16 -> utf8
	utf16Str := make([]uint16, len(content)/2)
	err = binary.Read(bytes.NewReader(content), binary.LittleEndian, &utf16Str)
	if err != nil {
		LogErr("utf16 -失败-> utf8", err)
	}
	return string(utf16.Decode(utf16Str)), nil
}
