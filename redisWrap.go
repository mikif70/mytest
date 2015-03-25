// redisWrap

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Config struct {
	wrapUrl  string
	ssdbUrl  string
	wrapAddr *net.TCPAddr
	ssdbAddr *net.TCPAddr
	server   string
}

type Conn struct {
	conn     net.Conn
	bufRead  *bufio.Reader
	bufWrite *bufio.Writer
}

var (
	config = &Config{
		wrapUrl: "0.0.0.0:6380",
		ssdbUrl: "10.39.80.181:8888",
	}
	okReply   interface{} = "OK"
	pongReply interface{} = "PONG"
)

func (c *Conn) readLine() ([]byte, error) {
	p, err := c.bufRead.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		return nil, errors.New("long response line")
	}
	if err != nil {
		return nil, err
	}
	i := len(p) - 2
	fmt.Println(string(p), i)
	if i < 0 || p[i] != '\r' {
		return nil, errors.New("bad response line terminator")
	}
	return p[:i], nil
}

func (c *Conn) readReply() (interface{}, error) {
	line, err := c.readLine()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, errors.New("short response line")
	}
	switch line[0] {
	case '+':
		switch {
		case len(line) == 3 && line[1] == 'O' && line[2] == 'K':
			// Avoid allocation for frequent "+OK" response.
			return okReply, nil
		case len(line) == 5 && line[1] == 'P' && line[2] == 'O' && line[3] == 'N' && line[4] == 'G':
			// Avoid allocation in PING command benchmarks :)
			return pongReply, nil
		default:
			return string(line[1:]), nil
		}
	case '-':
		return fmt.Sprintf("server error: ", string(line[1:])), nil
	case ':':
		return parseInt(line[1:])
	case '$':
		n, err := parseLen(line[1:])
		if n < 0 || err != nil {
			return nil, err
		}
		p := make([]byte, n)
		_, err = io.ReadFull(c.bufRead, p)
		if err != nil {
			return nil, err
		}
		if line, err := c.readLine(); err != nil {
			return nil, err
		} else if len(line) != 0 {
			return nil, errors.New("bad bulk string format")
		}
		return string(p), nil
	case '*':
		n, err := parseLen(line[1:])
		if n < 0 || err != nil {
			return nil, err
		}
		r := make([]interface{}, n)
		for i := range r {
			r[i], err = c.readReply()
			if err != nil {
				return nil, err
			}
		}
		return r, nil
	}
	return nil, errors.New("unexpected response line")
}

// parseInt parses an integer reply.
func parseInt(p []byte) (interface{}, error) {
	if len(p) == 0 {
		return 0, errors.New("malformed integer")
	}

	var negate bool
	if p[0] == '-' {
		negate = true
		p = p[1:]
		if len(p) == 0 {
			return 0, errors.New("malformed integer")
		}
	}

	var n int64
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return 0, errors.New("illegal bytes in length")
		}
		n += int64(b - '0')
	}

	if negate {
		n = -n
	}
	return n, nil
}

// parseLen parses bulk string and array lengths.
func parseLen(p []byte) (int, error) {
	if len(p) == 0 {
		return -1, errors.New("malformed length")
	}

	if p[0] == '-' && len(p) == 2 && p[1] == '1' {
		// handle $-1 and $-1 null replies.
		return -1, nil
	}

	var n int
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return -1, errors.New("illegal bytes in length")
		}
		n += int(b - '0')
	}

	return n, nil
}

func init() {
	flag.StringVar(&config.ssdbUrl, "s", config.ssdbUrl, "ssdb ip:port")
	flag.StringVar(&config.wrapUrl, "l", config.wrapUrl, "listen ip:port")
}

func main() {
	flag.Parse()

	config.ssdbAddr, _ = net.ResolveTCPAddr("tcp", config.ssdbUrl)
	config.wrapAddr, _ = net.ResolveTCPAddr("tcp", config.wrapUrl)

	fmt.Println(config)

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

	ssdb, err := net.DialTCP("tcp", nil, config.ssdbAddr)
	if err != nil {
		fmt.Println("Dial err: ", err.Error())
		return
	}
	defer ssdb.Close()

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			fmt.Println("Accept err: ", err.Error())
			continue
		}

		go func() {
			defer conn.Close()

			c := &Conn{
				conn:     conn,
				bufRead:  bufio.NewReader(conn),
				bufWrite: bufio.NewWriter(conn),
			}

			start := time.Now()
			retval, _ := c.readReply()
			fmt.Println("parseCmd: ", time.Since(start))

			ret := retval.([]interface{})
			for i := range ret {
				fmt.Println(ret[i])
			}
		}()

	}
}
