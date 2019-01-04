package bdlog

import (
	"bufio"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func AddLog(log *string, fileName string) error {
	fullName := beego.AppConfig.String("log_path") + fileName + ".txt"

	_, err := os.Stat(fullName)
	if os.IsNotExist(err) {
		file, err := os.Create(fullName)
		if err != nil {
			return err
		}
		fileStr := *log

		index := strings.Index(fileStr, "[PostBuildScript]")
		if index > 0 {
			logs := strings.Split(fileStr, "[PostBuildScript]")
			fileStr = logs[0]
		}
		_, err = io.WriteString(file, fileStr)
		if err != nil {
			return err
		}
		return nil
	}

	fmt.Println("exist")
	return err
}

func ReadLog(fileName string) (*string, error) {
	fullName := beego.AppConfig.String("log_path") + fileName + ".txt"
	if _, err := os.Stat(fullName); os.IsNotExist(err) {
		fmt.Println("not exist")
		return nil, err
	}
	b, err := ioutil.ReadFile(fullName)
	if err != nil {
		//fmt.Println(err.Error())
		return nil, err
	}
	log := string(b)
	return &log, nil
}

func ReadLogLine(fileName string, start, size int64) (*string, bool, error) {
	fullName := beego.AppConfig.String("log_path") + fileName + ".txt"
	if _, err := os.Stat(fullName); os.IsNotExist(err) {
		fmt.Println("not exist")
		return nil, false, err
	}
	file, err := os.OpenFile(fullName, os.O_RDONLY, 0660)
	if err != nil {
		fmt.Println(err.Error())
		return nil, false, err
	}
	r := bufio.NewReader(file)
	content := ""
	var i int64
	i = 1
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			fmt.Println(err.Error())
			return &content, false, err
		}
		if i == (start + size) {
			return &content, true, nil
		}
		if i >= start {
			content += line
		}
		if i == (start + size - 1) {
			i++
			continue
		}
		i++
	}

	return &content, true, nil
}
