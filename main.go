package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
var logMaxLine int = 30

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
type Log struct {
	Aid           int64  //8
	Ip            int64  //10
	From          string //32
	File_name     string //128
	Crtime        int64  //10
	Line          int64  //10
	ContentLength int64  //10
	Content       string
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
		//debug i <= 1
		for i := 1; i <= 3; i++ {
			chunkLog := readChunkLog(bfio, value.Separator)
			currentLine := int(Last_line) + i - 1
			str := pack(value, currentLine, chunkLog)
			//fmt.Println(str)
			if str == "" && str != "\n" {
				break
			}

			sendLen, err := conn.Write([]byte(str))
			checkError(err)
			fmt.Println(sendLen, str)
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
		//fmt.Println("--")
		_, err := r.ReadString(byte('\n'))
		if err != nil {
			log.Fatal(err)
			return r, err
		}
	}
	return r, nil
}

/**
 * 读取一条日志，可能是多行
 */
func readChunkLog(r *bufio.Reader, Separator string) string {
	var chunkLog string
	for i := 1; i <= logMaxLine; i++ {
		currentLine, err := r.ReadString(byte('\n'))
		chunkLog += currentLine
		if err == io.EOF {
			return chunkLog
		}
		if err != nil {
			return chunkLog
		}
		re := regexp.MustCompile(Separator)
		separatorExist := re.FindAllString(currentLine, 1)
		if separatorExist != nil {
			return chunkLog
		}
		//fmt.Println(currentLine)
	}
	return chunkLog
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

func pack(task Task, line int, chunkLog string) (str string) {
	var log Log
	log.Aid, _ = task.Id.Int64()
	log.Ip, _ = task.Ip.Int64()
	log.From = task.From
	log.File_name = task.Path
	log.Crtime = time.Now().Unix()
	log.Line = int64(line)
	log.ContentLength = int64(len(chunkLog))
	log.Content = chunkLog
	//fmt.Println(log)

	str = fmt.Sprintf("%08d", log.Aid)
	str += fmt.Sprintf("%010d", log.Ip)
	str += fmt.Sprintf("%032s", log.From)
	str += fmt.Sprintf("%0128s", log.File_name)
	str += fmt.Sprintf("%010d", log.Crtime)
	str += fmt.Sprintf("%010d", log.Line)
	str += fmt.Sprintf("%010d", log.ContentLength)
	str += log.Content

	//fmt.Println(str)
	return str

}
