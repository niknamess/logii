package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"gitlab.topaz-atcs.com/tmcs/logi2/terminal"
)

var (
	limit  string
	ipaddr []string
)

func Client() {
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
