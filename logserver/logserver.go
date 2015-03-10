package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net"
	//"os"
	"regexp"
	"strconv"
	"strings"
)

type Log struct {
	Aid           int64  //8
	Ip            int64  //10 字符串处理
	From          string //32
	File_name     string //128
	Crtime        int64  //10
	Line          int64  //10
	ContentLength int64  //8
	Content       string
}

func main() {

	l, err := net.Listen("tcp", "127.0.0.1:5000")
	checkError(err)
	defer l.Close()
	db := initMysql()
	defer db.Close()
	//debug
	db.Exec("TRUNCATE TABLE `127.0.0.1_yii`")
	db.Exec("TRUNCATE TABLE `127.0.0.1_apache`")

	var logs = make(chan Log, 30000)

	for {
		conn, err := l.Accept()
		checkError(err)

		go handleRequest(conn, db, logs)
	}
}

/**
 * 处理请求
 * conn socket对象
 * db 	数据库对象
 * logs 结构化日志管道
 *
 */
func handleRequest(conn net.Conn, db *sql.DB, logs chan Log) {
	defer conn.Close()
	go handleLog(logs, db)
	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 10240))
	// 数据缓冲
	databuf := make([]byte, 4096)
	// 消息长度
	length := 0
	// 消息长度uint64
	ulength := uint32(0)
	// 数据循环
	for {
		// 读取数据
		n, err := conn.Read(databuf)
		if err == io.EOF {
			fmt.Printf("Client exit: %s\n", conn.RemoteAddr())
			//fmt.Println(msgbuf.Len(), msgbuf.String())
		}
		if err != nil {
			fmt.Printf("Read error: %s\n", err)
			return
		}
		fmt.Println("databuf len:", n)

		// 数据添加到消息缓冲
		n, err = msgbuf.Write(databuf[:n])
		if err != nil {
			fmt.Printf("Buffer write error: %s\n", err)
			return
		}

		// 消息分割循环
		for {
			// 消息头
			if length == 0 && msgbuf.Len() >= 206 {
				binary.Read(msgbuf, binary.LittleEndian, &ulength)
				length = int(ulength)
				fmt.Println(length)
				//fmt.Println(msgbuf.String())
				//contentLength, _ := strconv.Atoi(string(msgbuf.String()[198:206]))
				//length = contentLength
				//fmt.Println(contentLength, msgbuf.Len())
				// 检查超长消息
				if length > 10240 {
					fmt.Printf("Message too length: %d\n", length)
					return
				}
			}
			//os.Exit(1)
			// 消息体
			if length > 0 && msgbuf.Len() >= length {
				//fmt.Printf("Client messge: %s\n", string(msgbuf.Next(length)))
				//fmt.Println(string(msgbuf.Next(length)))
				log := unpack(string(msgbuf.Next(length)))
				fmt.Printf("%+v", log)
				logs <- log
				length = 0
			} else {
				break
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initMysql() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/log_gather")
	checkError(err)
	//defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	checkError(err)
	return db
}

func unpack(str string) (logStruct Log) {
	//str = "0000000100000000000000000000000000apache000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000E:\\wamp\\logs\\apache_error_20150226.log142492272700000000110000000009line 11"
	var log Log
	//log.Aid = fmt.Sprintf("%d", str[0:7])
	id, _ := strconv.Atoi(str[0:8])
	ip, _ := strconv.Atoi(str[8:18])
	crtime, _ := strconv.Atoi(str[178:188])
	line, _ := strconv.Atoi(str[188:198])
	contentLength, _ := strconv.Atoi(str[198:206])
	log.Aid = int64(id)
	log.Ip = int64(ip)
	log.From = strings.TrimLeft(str[18:50], "0")
	log.File_name = strings.TrimLeft(str[50:178], "0")
	log.Crtime = int64(crtime)
	log.Line = int64(line)
	log.ContentLength = int64(contentLength)
	log.Content = str[206:]
	//fmt.Println(str[188:198], contentLength)
	//fmt.Printf("%+v", log)

	//debug
	//if str[8:18] != "2130706433" || (int64(contentLength) != int64(len(str[206:]))) {
	//	fmt.Println("\n\nerror Str", str)
	//}
	return log
}

func getTableName(log Log) (tableName string) {
	//strings.Replace(log.)
	tableName = "log_" + inet_ntoa(log.Ip).String() + "_" + log.From
	return tableName
}

func handleLog(logs chan Log, db *sql.DB) {
	for {
		select {
		case logStruct := <-logs:
			//do nothing
			writeLog(logStruct, db)
			//fmt.Println("c---->", logStruct)
		default:
			//warnning!
			//fmt.Errorf("TASK_CHANNEL is full!")
		}
	}
}

func writeLog(log Log, db *sql.DB) {
	tableName := getTableName(log)
	//fmt.Println(tableName)

	stmt, err := db.Prepare("insert into `" + tableName + "` values('',?,?,?,?,?,?)")
	defer stmt.Close()

	checkError(err)
	res, err := stmt.Exec(log.Aid, log.From, log.File_name, log.Crtime, log.Line, log.Content)
	checkError(err)
	if row, _ := res.LastInsertId(); row != 0 {
		stmt, err := db.Prepare("update app set last_line=last_line+?,last_time=? where ip=? and `from`=?")
		defer stmt.Close()

		checkError(err)
		re := regexp.MustCompile(`\n`)
		lines_arr := re.FindAllString(log.Content, -1)
		lines := int64(len(lines_arr))
		//fmt.Println("logline", lines)

		stmt.Exec(lines, log.Crtime, log.Ip, log.From)
	}
}

// Convert uint to net.IP http://www.sharejs.com
func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

// Convert net.IP to int64 ,  http://www.sharejs.com
func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}
