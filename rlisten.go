// rlisten
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

type Config struct {
	wrapUrl  string
	ssdbUrl  string
	wrapAddr *net.TCPAddr
	ssdbAddr *net.TCPAddr
}

var (
	config = &Config{
		wrapUrl: "0.0.0.0:6380",
		ssdbUrl: "10.39.80.181:8888",
	}
)

func parseCmd([]byte) {

}

func init() {
	flag.StringVar(&config.ssdbUrl, "s", config.ssdbUrl, "ssdb ip:port")
	flag.StringVar(&config.wrapUrl, "l", config.wrapUrl, "listen ip:port")
}

func main() {
	flag.Parse()

	config.ssdbAddr, _ = net.ResolveTCPAddr("tcp", config.ssdbUrl)
	config.wrapAddr, _ = net.ResolveTCPAddr("tcp", config.wrapUrl)
	fmt.Printf("Listen: %+v\n", config.wrapAddr)
	ln, err := net.ListenTCP("tcp", config.wrapAddr)
	if err != nil {
		fmt.Println("Listen err: ", err.Error())
		return
	}
	defer ln.Close()

	file, err := os.Create("dovecot-redis.log")
	if err != nil {
		fmt.Println("File err: ", err.Error())
		return
	}
	defer file.Close()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err: ", err.Error())
			continue
		}

		go func() {
			defer conn.Close()
			buf := make([]byte, 1024)
			nr, err := conn.Read(buf)
			if err != nil {
				fmt.Println("Read1 err: ", err.Error())
				return
			}

			parseCmd(buf[:nr])

			file.Write(buf[:nr])
			file.Write([]byte("++++++++\n\r"))
			fmt.Printf("Read %d bytes:\n%v++++++++\n", nr, string(buf[:nr]))

			ssdb, err := net.DialTCP("tcp", nil, config.ssdbAddr)
			if err != nil {
				fmt.Println("Dial err: ", err.Error())
				return
			}
			defer ssdb.Close()

			ssdb.Write(buf[:nr])

			nr, err = ssdb.Read(buf)
			if err != nil {
				fmt.Println("Read2 err: ", err.Error())
				conn.Write([]byte("+OK\r\n"))
				return
			}

			conn.Write(buf[:nr])
			fmt.Printf("Reply %d bytes:\n%v--------\n\n", nr, string(buf[:nr]))
			file.Write(buf[:nr])
			file.Write([]byte("--------\n\r"))

		}()
	}

	fmt.Println("Hello World!")
}
