// rlisten
package main

/*
import (
	"flag"
	"fmt"
	"net"
)

type Config struct {
	wrapUrl  string
	ssdbUrl  string
	wrapAddr *net.TCPAddr
	ssdbAddr *net.TCPAddr
	server   string
}

const (
	NL = "\r\n"
)

var (
	//	ErrInvalidSyntax = error.New("resp: invalid syntax")
	config = &Config{
		wrapUrl: "0.0.0.0:6380",
		ssdbUrl: "10.39.80.181:8888",
		server:  "",
	}
)

func init() {
	flag.StringVar(&config.ssdbUrl, "s", config.ssdbUrl, "ssdb ip:port")
	flag.StringVar(&config.wrapUrl, "l", config.wrapUrl, "listen ip:port")
}

func main() {
	flag.Parse()

	config.ssdbAddr, _ = net.ResolveTCPAddr("tcp", config.ssdbUrl)
	config.wrapAddr, _ = net.ResolveTCPAddr("tcp", config.wrapUrl)

	fmt.Println(config)

	wrap, err := net.DialTCP("tcp", nil, config.wrapAddr)
	if err != nil {
		fmt.Println("Dial wrap err: ", err.Error())
		return
	}
	defer wrap.Close()


		wrap.Write([]byte("*1" + NL + "$5" + NL + "MULTI" + NL + "*3" + NL + "$6" +
			NL + "INCRBY" + NL + "$32" + NL + "dtest01@tiscali.it/quota/storage" +
			NL + "$5" + NL + "-1622" + NL + "*3" + NL + "$6" + NL + "INCRBY" + NL + "$33" +
			NL + "dtest01@tiscali.it/quota/messages" + NL + "$2" + NL + "-2" + NL + "*1" +
			NL + "$4" + NL + "EXEC" + NL + "*2" + NL + "$3" + NL + "GET" + NL + "$32" +
			NL + "dtest01@tiscali.it/quota/storage" + NL))

	buf := make([]byte, 512)
	n, err := wrap.Read(buf)
	fmt.Printf("Reply: %+v\n%+v\n", buf[:n], string(buf[:n]))
}

*/
