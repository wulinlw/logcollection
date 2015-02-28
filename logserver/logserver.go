package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net"
	"os"
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
	ContentLength int64  //10
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
	buf := make([]byte, 10240)
	container := make([]byte, 0)
	var allBuf []byte
	for {

		reqLen, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("\n\n\nmsglength", reqLen)
		//fmt.Println(reqLen, string(buf[:reqLen]))
		//logStruct := unpack(string(buf[:reqLen]))
		//fmt.Println(logStruct)

		if reqLen <= 10240 {

			if len(container) != 0 {
				allBuf = append(container, buf[:reqLen]...)
				//fmt.Println("\n\n\n", string(allBuf[:208]), "\n\n\n")
			} else {
				allBuf = buf[:reqLen]
			}

			//var readedLen int = 0
			//fmt.Println("\nallBuf:", string(allBuf), "\n")
			for {

				if len(allBuf[:]) < 208 {
					container = allBuf[:]
					break
				} else if len(allBuf[:]) >= 208 {
					contentLength, _ := strconv.Atoi(string(allBuf[198:208]))
					fmt.Println(contentLength, len(allBuf[:]))
					if contentLength == 0 {
						//这里的调试信息 打印字符串出来，看下是否完整
						//fmt.Println("\n\n contentLength(0):", len(allBuf[:]), string(allBuf), "\n\n")
					}
					if 208+contentLength > len(allBuf[:]) {
						container = allBuf[:]
						//fmt.Println("\n\n lastbuf:", len(allBuf[:]), string(allBuf), "\n\n")
						break
					}
					//if 208+contentLength == len(allBuf[:]) {
					//	fmt.Println("\n\n =====:", len(allBuf[:]), string(allBuf), "\n\n")
					if false {
						os.Exit(0)
					}
					//}
					logStruct := unpack(string(allBuf[:208+contentLength]))
					//fmt.Println("\n\n\n", logStruct.Line)
					//debug
					//logStruct.Content = ""
					//fmt.Printf("%+v", logStruct)
					allBuf = allBuf[208+contentLength:]
					logs <- logStruct
				}
			}
		} else {
			log.Fatal("data too long")
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
	contentLength, _ := strconv.Atoi(str[198:208])
	log.Aid = int64(id)
	log.Ip = int64(ip)
	log.From = strings.TrimLeft(str[18:50], "0")
	log.File_name = strings.TrimLeft(str[50:178], "0")
	log.Crtime = int64(crtime)
	log.Line = int64(line)
	log.ContentLength = int64(contentLength)
	log.Content = str[208:]
	//fmt.Println(str[188:198], contentLength)
	//fmt.Printf("%+v", log)

	//debug
	if str[8:18] != "2130706433" || (int64(contentLength) != int64(len(str[208:]))) {
		fmt.Println("\n\nerror Str", str)
	}
	return log
}

func getTableName(log Log) (tableName string) {
	//strings.Replace(log.)
	tableName = inet_ntoa(log.Ip).String() + "_" + log.From
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
