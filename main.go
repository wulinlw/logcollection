package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var taskUrl string = "http://www.test.com/log_collection/backend/task.php"

type Task struct {
	Id        json.Number `json:"id"`
	Ip        json.Number `json:"ip"`
	From      string      `json:"from"`
	Path      string      `json:"path"`
	Separator string      `json:"separator"`
	Last_line json.Number `json:"last_line"`
	Last_time json.Number `json:"last_time"`
	Describe  string      `json:"describe"`
}

func main() {
	task := getTask()
	fmt.Println(task)
	task = parseLogPath(&task)
	fmt.Println(task)

	for _, value := range task {
		file, err := os.Open(value.Path)
		checkError(err)
		defer file.Close()
		bfio := bufio.NewReader(file)
		Last_line, _ := value.Last_line.Int64()
		bfio, err = seek(bfio, int(Last_line))
		checkError(err)

		conn, err := net.Dial("tcp4", "127.0.0.1:5000")
		checkError(err)
		defer conn.Close()

		for i := 1; ; i++ {
			currentLine, err := readLine(bfio)
			if err != nil {
				break
			}

			sendLen, err := conn.Write([]byte(currentLine))
			checkError(err)
			fmt.Println(sendLen, currentLine)
		}

	}

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/**
 * 跳至指定行数
 */
func seek(r *bufio.Reader, lineNum int) (h *bufio.Reader, err error) {
	for i := 1; i < lineNum; i++ {
		_, err := r.ReadString(byte('\n'))
		if err != nil {
			log.Fatal(err)
			return r, err
		}
	}
	return r, nil
}

/**
 * 读取一行数据
 */
func readLine(r *bufio.Reader) (string, error) {
	currentLine, err := r.ReadString(byte('\n'))
	//checkError(err)
	//fmt.Println(currentLine)
	return currentLine, err
}

/**
 * 读取任务列表
 */
func getTask() (t []Task) {
	resp, err := http.Get(taskUrl)
	checkError(err)
	//fmt.Println(resp)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	//body = []byte(`[{"id":"1","ip":"2130706433","from":"apache"}]`)
	var task []Task
	err = json.Unmarshal([]byte(body), &task)
	checkError(err)
	//fmt.Println(task, task[0].Ip)
	return task
}

/**
 * 解析log路径中的时间字符串
 * 将   \logs\apache_error_{Ymd}.log
 * 解析为\logs\apache_error_20060102.log
 * Y	年 2006
 * m	月 01
 * d	日 02
 * H	时 15
 * i	分 04
 * s	秒 05
 */
func parseLogPath(task *[]Task) (t []Task) {
	for key, value := range *task {
		//fmt.Println(key, value)
		re := regexp.MustCompile(`{.*}`)
		dateFormatStr := re.FindAll([]byte(value.Path), 1)
		//fmt.Println(string(dateFormatStr[0]))//{Ymd}
		if dateFormatStr != nil {
			var dateStr string
			for _, v := range dateFormatStr[0][1 : len(dateFormatStr[0])-1] {
				//fmt.Println(k, v)
				switch string(v) {
				case "Y":
					dateStr += "2006"
				case "m":
					dateStr += "01"
				case "d":
					dateStr += "02"
				case "H":
					dateStr += "15"
				case "i":
					dateStr += "04"
				case "s":
					dateStr += "05"
				default:
					dateStr += string(v)
				}
			}
			dateStr = time.Now().Format(dateStr)
			(*task)[key].Path = strings.Replace((*task)[key].Path, string(dateFormatStr[0][:]), dateStr, 1)
			//fmt.Println((*task)[key].Path, dateStr)
		}
	}
	return *task
}
