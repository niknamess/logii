// Very basic socket server
// https://golangr.com/

package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"gitlab.topaz-atcs.com/tmcs/logi2/terminal"
	"gitlab.topaz-atcs.com/tmcs/logi2/web"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var str string = "VFC"

func Server() {
	fmt.Println("Start server...")
	fmt.Println("Press \"VFC\" or \"WEB\" to start service")

	// listen on port 8000
	ln, _ := net.Listen("tcp", "127.0.0.1:8888")

	// accept connection
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// get message, output
		message, _ := bufio.NewReader(conn).ReadString('\n')
		//fmt.Print(message)
		s := strings.TrimSpace(message)
		if s == str {
			conn.Write([]byte("VFC: "))
			go controllers.VFC("10015")
			conn.Write([]byte("Выбрана служба vfc" + "\n"))
		}
		if s == "WEB" {
			terminal.CallClear()
			conn.Write([]byte("WEB: "))
			//conn.Write([]byte("Enter IP or enter \"stop\": " + "\n"))
			//отдeльная функция для   Web  отправки
			allip := enterIp(conn)

			go web.ProcWeb("-x", allip) //ща че нить придумаем
			//conn.Write([]byte("Выбрана функция web" + "\n"))
		}
		if s == "stopW" {
			terminal.CallClear()
			var stop []string

			go web.ProcWeb("-x", stop) //ща че нить придумаем
			//conn.Write([]byte("Выбрана функция web" + "\n"))
		}
		conn.Write([]byte(message + "\n"))
		//fmt.Print("Message Received:", string(message))
	}
}

func enterIp(conn net.Conn) []string {

	for {
		//fmt.Print("Enter IP:  ")
		//fmt.Scanln(&limit)
		conn.Write([]byte("Enter IP or enter \"stop\": " + "\n"))
		limit, _ := bufio.NewReader(conn).ReadString('\n')
		limit = strings.TrimSpace(limit)
		if limit == "stop" {
			break
		} else if util.CheckIPAddress(limit) {
			conn.Write([]byte("Valid: " + "\n"))
			ipaddr = append(ipaddr, limit)
			//limitSlice, _ := web.CheckConfig()
			//ipaddr = append(ipaddr, limitSlice...)
			ipaddr = removeDuplicateStr(ipaddr)
			//config := Config{DataBase: DatabaseConfig{Hostt: ipaddr, Port: "10015"}}
			//data, _ = toml.Marshal(&config)
		} else {
			conn.Write([]byte("Invalid: " + "\n"))
		}
	}
	return ipaddr
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
