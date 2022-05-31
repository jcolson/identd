package main

import (
	"flag"
	"log"
	"net"
	"os"
	"runtime/debug"
)

const (
	RESPONSE  = " : USERID : UNIX : "
	CONN_PORT = "113"
	CONN_TYPE = "tcp"
)

var helpFlag = flag.Bool("h", false, "Help")
var userFlag = flag.String("u", "root", "The userid to respond with in ident request.")

func main() {
	flag.Parse()
	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(1)
	}
	user := *userFlag
	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func(c net.Conn) {
			defer func() {
				if panicInfo := recover(); panicInfo != nil {
					log.Printf("%v, %s", panicInfo, string(debug.Stack()))
				}
			}()
			buf := make([]byte, 4096)
			i, err := c.Read(buf)
			if err != nil {
				log.Print(err)
				c.Close()
				return
			}
			for string(buf[i-1]) != "\n" {
				j, err := c.Read(buf)
				if err != nil {
					log.Print(err)
					c.Close()
					return
				}
				i += j
			}
			if string(buf[i-2]) == "\r" {
				i -= 1
			}
			c.Write(append(buf[:i-1], []byte(RESPONSE+user+"\r\n")...))
			c.Close()
		}(conn)
	}
}
