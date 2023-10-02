package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf16"
)

func main() {
	start := time.Now()
	replaceAndWrite(readConfig())
	cost := time.Since(start)
	log.Println("本次运行耗时:", cost)
}

func replaceAndWrite(config Config) {
	jsonMap := readJson(config.JsonPath)
	txt := readTxt(config.TxtPath)
	log.Println("成功读取所有文件, 开始替换...")
	changeByteCount := 0
	beforeByteCount := len(txt)
	for old, _new := range jsonMap {
		changeByteCount += len(_new) - len(old)
		txt = strings.Replace(txt, old+"\r\n", _new+"\r\n", 1)
		//log.Println(old, "->", _new)
	}
	log.Println("成功替换所有文本", "预期字节数变化", changeByteCount,
		"实际字节数变化", len(txt)-beforeByteCount, ",开始写入文本到", config.NewTxtPath)
	//将替换后的文本写入到新文件中
	newFile, err := os.Create(config.NewTxtPath)
	if err != nil {
		log.Fatalf(red("创建文件失败, 错误信息: %v"), err)
	}
	defer newFile.Close()

	//以utf16写入
	utf16txt := utf16.Encode([]rune(txt))
	writer := bufio.NewWriter(newFile)
	err = binary.Write(writer, binary.LittleEndian, &utf16txt)
	//_, err = newFile.WriteString(txt)
	if err != nil {
		log.Fatalf(red("写入文件失败, 错误信息: %v"), err)
	}
	log.Println("成功写入!!!")
}

type Config struct {
	JsonPath   string `yaml:"json-path"`
	TxtPath    string `yaml:"txt-path"`
	NewTxtPath string `yaml:"new-txt-path"`
}

func readConfig() Config {
	file, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf(red("打开 config.yml 失败, 错误信息: %v"), err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf(red("读取 config.yml 失败, 可能是格式问题, 错误信息: %v"), err)
	}
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Fatalf(red("读取 config.yml 失败, 可能是格式问题, 错误信息: %v"), err)
	}
	log.Println("成功读取配置 Json文件:", config.JsonPath)
	log.Println("成功读取配置 英文文件:", config.TxtPath)
	log.Println("成功读取配置 汉化文件:", config.NewTxtPath)
	return config
}

func readJson(jsonPath string) map[string]string {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalf(red("文件 %v 打开失败, 错误信息: %v"), jsonPath, err)
	}
	defer jsonFile.Close()
	byteSlice, _ := io.ReadAll(jsonFile)
	jsonMap := make(map[string]string, 400000)
	err = json.Unmarshal(byteSlice, &jsonMap)
	if err != nil {
		log.Fatalf(red("文件 %v 反序列化错误, 可能是格式问题, 错误信息: %v"), jsonPath, err)
	}
	log.Println("成功读取文件", jsonPath, "有效键值对数:", len(jsonMap))
	return jsonMap
}

func readTxt(txtPath string) string {
	//打开要替换的txt文件
	file, err := os.Open(txtPath)
	if err != nil {
		log.Fatalf(red("打开txt文件失败, 错误信息: %v"), err)
	}
	defer file.Close()
	//读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf(red("读取txt文件失败, 错误信息: %v"), err)
	}
	log.Println("成功读取文件", txtPath, "有效字节数(英文数字1字节, 汉字3字节):", len(content))
	//将文件内容转换为字符串//utf16 -> utf8
	utf16Str := make([]uint16, len(content)/2)
	err = binary.Read(bytes.NewReader(content), binary.LittleEndian, &utf16Str)
	if err != nil {
		log.Fatalln(err)
	}
	return string(utf16.Decode(utf16Str))
}

func red(formatStr string) string {
	return "\033[31m" + formatStr + "\033[0m"
}
