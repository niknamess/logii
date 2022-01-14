package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

var (
	limit  string
	ipaddr []string
)

func Client() {

	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", "127.0.0.1:8888")
	for {
		// Чтение входных данных от stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Now you can use only VFC(stable) and WEB(in work): ")
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// Отправляем в socket
		fmt.Fprintf(conn, text+"\n")
		// Прослушиваем ответ
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: " + message)
	}
}
