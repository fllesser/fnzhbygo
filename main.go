package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	config := readConfig()
	replaceAndWrite(config)
}

func replaceAndWrite(config Config) {
	jsonMap := readJson(config.JsonPath)
	txt := readTxt(config.TxtPath)
	log.Println("读取完毕, 开始替换文本...")
	for old, _new := range jsonMap {
		txt = strings.Replace(txt, old, _new, -1)
	}

	log.Println("替换完毕, 开始写入文本到", config.NewTxtPath)
	//将替换后的文本写入到新文件中
	newFile, err := os.Create(config.NewTxtPath)
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	defer newFile.Close()

	//写入内容
	_, err = newFile.WriteString(txt)
	if err != nil {
		fmt.Println("写入文件失败:", err)
		return
	}
	log.Println("写入成功!!!")
}

type Config struct {
	JsonPath   string `yaml:"json-path"`
	TxtPath    string `yaml:"txt-path"`
	NewTxtPath string `yaml:"new-txt-path"`
}

func readConfig() Config {
	file, err := os.Open("config.yml")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	var config Config
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Println(err)
	}
	log.Printf("成功读取 config.yml:%+v", config)
	return config
}

func readJson(jsonPath string) map[string]string {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		log.Println("json文件读取错误", err)
		return nil
	}
	log.Println("成功读取", jsonPath)
	defer jsonFile.Close()
	bytes, _ := io.ReadAll(jsonFile)
	jsonMap := make(map[string]string, 400000)
	err = json.Unmarshal(bytes, &jsonMap)
	if err != nil {
		log.Println("json反序列化错误, 可能是格式问题, 具体错误信息: ", err)
	}
	return jsonMap
}

func readTxt(txtPath string) string {
	//打开要替换的txt文件
	file, err := os.Open(txtPath)
	if err != nil {
		log.Println("打开文件失败:", err)
		return ""
	}
	defer file.Close()
	log.Println("成功读取", txtPath)
	//读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		log.Println("读取文件失败:", err)
		return ""
	}
	//将文件内容转换为字符串
	return string(content)
}
