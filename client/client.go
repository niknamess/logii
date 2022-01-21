package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"gitlab.topaz-atcs.com/tmcs/logi2/terminal"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("From server:", string(buf[0:n])) //From server
	}
}

func Client() {

	terminal.CallClear()
	fmt.Println("Now you can use only VFC(stable) or WEB(stable) and stop this service")
	c, err := net.Dial("unix", "/tmp/echo.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	go reader(c)
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: ")

		text, _ := reader.ReadString('\n') //Send server
		//ui terminal

		_, err := c.Write([]byte(text)) //Send server
		if err != nil {
			log.Fatal("write error:", err)
			break
		}
		time.Sleep(1e9)
	}
}

/* func Client() {
	//reader := bufio.NewReader(os.Stdin)
	//ipaddress, _ := reader.ReadString('\n')
	//util.CheckIPAddress(ipaddress)
	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.1:8888")
	terminal.CallClear()
	fmt.Println("Now you can use only VFC(stable) and WEB(in work)")
	for {
		//message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print("Message from server: " + message)
		//terminal.CallClear()
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		//mt.Println("Now you can use only VFC(stable) and WEB(in work)")
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// Отправляем в socket
		fmt.Fprintf(conn, text+"\n")
		// Прослушиваем ответ
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
*/
