// rlisten
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp", "0.0.0.0:6380")
	fmt.Printf("Listen: %+v\n", addr)
	ln, err := net.ListenTCP("tcp", addr)
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

			file.Write(buf[:nr])
			file.Write([]byte("++++++++\n\r"))
			fmt.Printf("Read: %d - %v\n", nr, string(buf[:nr]))

			ssdbAddr, _ := net.ResolveTCPAddr("tcp", "10.39.80.181:8888")
			ssdb, err := net.DialTCP("tcp", nil, ssdbAddr)
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
			fmt.Printf("Reply: %d - %v\n", nr, string(buf[:nr]))
			file.Write(buf[:nr])
			file.Write([]byte("--------\n\r"))

		}()
	}

	fmt.Println("Hello World!")
}
