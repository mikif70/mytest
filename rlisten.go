// rlisten
package main

//import (
//	"bufio"
//	"bytes"
//	"flag"
//	"fmt"
//	"io"
//	"net"
//	"os"
//	"strconv"
//	"time"
//)

//type Config struct {
//	wrapUrl  string
//	ssdbUrl  string
//	wrapAddr *net.TCPAddr
//	ssdbAddr *net.TCPAddr
//	server   string
//}

//const (
//	SIMPLE_STRING = '+'
//	BULK_STRING   = '$'
//	INTEGER       = ':'
//	ARRAY         = '*'
//	ERROR         = '-'
//)

//var (
//	//	ErrInvalidSyntax = error.New("resp: invalid syntax")
//	config = &Config{
//		wrapUrl: "0.0.0.0:6380",
//		ssdbUrl: "10.39.80.181:8888",
//		server:  "",
//	}
//)

//type RESPReader struct {
//	*bufio.Reader
//}

//func NewReader(reader io.Reader) *RESPReader {
//	return &RESPReader{
//		Reader: bufio.NewReaderSize(reader, 32*1024),
//	}
//}

//func parseCmd(buf []byte) []byte {

//	row := bytes.Split(buf[:len(buf)-1], []byte("\n"))
//	ln := len(row)

//	retval := make([][][]byte, 0)

//	for i := 0; i < ln; i++ {
//		newAr := make([][]byte, 0)
//		switch row[i][0] {
//		case '*':
//			lenAr, _ := strconv.Atoi(string(row[i][1:]))
//			fmt.Println("Array: ", lenAr)
//			for a := 1; a <= lenAr*2; a++ {
//				newAr = append(newAr, row[i])
//				switch row[i+a][0] {
//				case '$':
//					newAr = append(newAr, row[i+a])
//					fmt.Println("len: ", string(row[i+a][1:]))
//				default:
//					if string(row[i+a]) == "MULTI" || string(row[i+a]) == "EXEC" {
//						newAr = make([][]byte, 0)
//						break
//					}
//					newAr = append(newAr, row[i+a])
//					fmt.Println("cmd: ", string(row[i+a]))
//				}
//			}
//			retval = append(retval, bytes.Join(newAr, []byte("\n")))
//			i += lenAr * 2
//		default:
//			fmt.Println("not found: ", string(row[i]))
//		}
//	}

//	fmt.Println(retval)

//	//	ret := bytes.Join(retval, []byte("\n"))

//	//	fmt.Println(ret)

//	return []byte{}
//}

//func init() {
//	flag.StringVar(&config.ssdbUrl, "s", config.ssdbUrl, "ssdb ip:port")
//	flag.StringVar(&config.wrapUrl, "l", config.wrapUrl, "listen ip:port")
//	flag.StringVar(&config.server, "c", config.server, "server")
//}

//func main() {
//	flag.Parse()

//	config.ssdbAddr, _ = net.ResolveTCPAddr("tcp", config.ssdbUrl)
//	config.wrapAddr, _ = net.ResolveTCPAddr("tcp", config.wrapUrl)

//	fmt.Println(config)

//	if config.server == "" {
//		fmt.Printf("Listen: %+v\n", config.wrapAddr)
//		ln, err := net.ListenTCP("tcp", config.wrapAddr)
//		if err != nil {
//			fmt.Println("Listen err: ", err.Error())
//			return
//		}
//		defer ln.Close()

//		file, err := os.Create("dovecot-redis.log")
//		if err != nil {
//			fmt.Println("File err: ", err.Error())
//			return
//		}
//		defer file.Close()

//		for {
//			conn, err := ln.AcceptTCP()
//			if err != nil {
//				fmt.Println("Accept err: ", err.Error())
//				continue
//			}

//			go func() {
//				defer conn.Close()
//				buf := make([]byte, 1024)
//				nr, err := conn.Read(buf)
//				if err != nil {
//					fmt.Println("Read1 err: ", err.Error())
//					return
//				}

//				bufRead := buf[:nr]

//				start := time.Now()
//				bufWrite := parseCmd(bufRead)
//				fmt.Println("parseCmd: ", time.Since(start))

//				file.Write(bufWrite)
//				file.Write([]byte("++++++++\n\r"))
//				fmt.Printf("Read %d bytes:\n%v++++++++\n", nr, string(buf[:nr]))

//				ssdb, err := net.DialTCP("tcp", nil, config.ssdbAddr)
//				if err != nil {
//					fmt.Println("Dial err: ", err.Error())
//					return
//				}
//				defer ssdb.Close()

//				ssdb.Write(buf[:nr])

//				nr, err = ssdb.Read(buf)
//				if err != nil {
//					fmt.Println("Read2 err: ", err.Error())
//					conn.Write([]byte("+OK\r\n"))
//					return
//				}

//				conn.Write(buf[:nr])
//				fmt.Printf("Reply %d bytes:\n%v--------\n\n", nr, string(buf[:nr]))
//				file.Write(buf[:nr])
//				file.Write([]byte("--------\n\r"))

//			}()
//		}
//	}

//	wrap, err := net.DialTCP("tcp", nil, config.wrapAddr)
//	if err != nil {
//		fmt.Println("Dial wrap err: ", err.Error())
//		return
//	}
//	defer wrap.Close()

//	wrap.Write([]byte(`*1
//$5
//MULTI
//*3
//$6
//INCRBY
//$32
//dtest01@tiscali.it/quota/storage
//$5
//-1622
//*3
//$6
//INCRBY
//$33
//dtest01@tiscali.it/quota/messages
//$2
//-2
//*1
//$4
//EXEC
//*2
//$3
//GET
//$32
//dtest01@tiscali.it/quota/storage
//`))

//}
