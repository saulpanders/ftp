/*
	@saulpanders
	FTP server
	ftpServer.go -- basic message protocol design/testing

	inspired by "network programming with go"
*/

package main

import (
	"fmt"
	"net"
	"os"
)

const (
	DIR = "DIR"
	CD  = "CD"
	PWD = "PWD"
)

func main() {

	service := "0.0.0.0:1202"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	var buf [512]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			conn.Close()
			return
		}

		s := string(buf[0:n])

		//decode request
		if s[0:2] == CD {
			chdir(conn, s[3:])
		} else if s[0:3] == DIR {
			dirList(conn)
		} else if s[0:3] == PWD {
			pwd(conn)
		}
	}
}

func chdir(conn net.Conn, s string) {
	if os.Chdir(s) == nil {
		conn.Write([]byte("OK"))
	} else {
		conn.Write([]byte("ERROR"))
	}
}

func dirList(conn net.Conn) {
	defer conn.Write([]byte("\r\n"))

	dir, err := os.Open(".")
	if err != nil {
		return
	}

	names, err := dir.Readdirnames(-1)
	if err != nil {
		return
	}

	for _, nm := range names {
		conn.Write([]byte(nm + "\r\n"))
	}
}

func pwd(conn net.Conn) {
	s, err := os.Getwd()
	if err != nil {
		conn.Write([]byte(""))
		return
	}
	conn.Write([]byte(s))
}
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

/*
	This is a simple protocol. The most complicated data structure that we need to send is an array of strings for a
directory listing. In this case we don't need the heavy duty serialisation techniques of the last chapter. In this
case we can use a simple text format.
But even if we make the protocol simple, we still have to specify it in detail. We choose the following message
format:
All messages are in 7-bit US-ASCII
The messages are case-sensitive
Each message consists of a sequence of lines
The first word on the first line of each message describes the message type. All other words are message
data
All words are separated by exactly one space character
Each line is terminated by CR-LF
Some of the choices made above are weaker in real-life protocols. For example
Message types could be case-insensitive. This just requires mapping message type strings down to lowercase
before decoding
An arbitrary amount of white space could be left between words. This just adds a little more complication,
compressing white space
Continuation characters such as "\" can be used to break long lines over several lines. This starts to
make processing more complex
Just a "\n" could be used as line terminator, as well as "\r\n" . This makes recognising end of line a
bit harder
All of these variations exist in real protocols. Cumulatively, they make the string processing just more complex
than in our case.
*/
