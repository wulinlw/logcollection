package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net"
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
	//fmt.Printf("%+v", db)

	for {
		conn, err := l.Accept()
		checkError(err)

		// Handle connections in a new goroutine.
		go handleRequest(conn, db)

	}
}
func handleRequest(conn net.Conn, db *sql.DB) {
	defer conn.Close()
	for {
		buf := make([]byte, 10240)
		reqLen, err := conn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Println(reqLen, string(buf[:reqLen]))
		logStruct := unpack(string(buf[:reqLen]))
		fmt.Println(logStruct)
		writeLog(logStruct, db)
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
	return log
}

func getTableName(log Log) (tableName string) {
	//strings.Replace(log.)
	tableName = inet_ntoa(log.Ip).String() + "_" + log.From
	return tableName
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
		stmt, err := db.Prepare("update app set last_line=?,last_time=?")
		defer stmt.Close()

		checkError(err)
		re := regexp.MustCompile(`\n`)
		lines_arr := re.FindAllString(log.Content, -1)
		lines := int64(len(lines_arr))
		//fmt.Println(log.Content, lines)
		var nextLine int64
		if lines > 1 {
			nextLine = log.Line + lines
		} else {
			nextLine = log.Line + 1
		}
		stmt.Exec(nextLine, log.Crtime)
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
