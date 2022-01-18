// Very basic socket server
// https://golangr.com/

package main

import (
	"log"
	"net"
	"strings"

	"gitlab.topaz-atcs.com/tmcs/logi2/web"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var mail string = "Succes"

func echoServer(c net.Conn) {
	for {
		buf := make([]byte, 512)
		nr, err := c.Read(buf)
		if err != nil {
			return
		}

		data := buf[0:nr]
		println("Server got:", string(data))
		s := strings.TrimSpace(string(data))
		if s == "VFC" {
			//conn.Write([]byte("VFC: "))
			data = []byte("Выбрана служба vfc") //Send Client
			_, err = c.Write(data)
			if err != nil {
				log.Fatal("Write: ", err)
			}
			go controllers.VFC("10015")
			//conn.Write([]byte("Выбрана служба vfc" + "\n"))
		}
		if s == "WEB" {
			data = []byte("Выбрана служба web") //Send Client
			_, err = c.Write(data)
			if err != nil {
				log.Fatal("Write: ", err)
			}
			//conn.Write([]byte("WEB: "))
			//conn.Write([]byte("Enter IP or enter \"stop\": " + "\n"))
			//отдeльная функция для   Web  отправки
			allip := enterIp(c)

			go web.ProcWeb("-x", allip) //ща че нить придумаем

			//conn.Write([]byte("Выбрана функция web" + "\n"))
		}
		//
		data = []byte(mail) //Send Client
		_, err = c.Write(data)
		if err != nil {
			log.Fatal("Write: ", err)
		}
		//
	}
}

func Server() {
	//logenc.DeleteOldsFiles()
	//e := os.Remove("/tmp/echo.sock")
	//if e != nil {
	//	log.Fatal(e)
	//s}
	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal("listen error:", err)
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}

		go echoServer(fd)
	}
}

/* var str string = "VFC"

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
} */

func enterIp(c net.Conn) []string {

	for {
		buf := make([]byte, 512)
		nr, _ := c.Read(buf)
		data := buf[0:nr]

		//data = []byte("В") //Send Client

		limit := string(data)
		limit = strings.TrimSpace(limit)
		if limit == "stop" {
			break
		} else if util.CheckIPAddress(limit) {
			data = []byte("Valid") //Send Client
			_, err := c.Write(data)
			if err != nil {
				log.Fatal("Write: ", err)
			}
			ipaddr = append(ipaddr, limit)
			//limitSlice, _ := web.CheckConfig()
			//ipaddr = append(ipaddr, limitSlice...)
			ipaddr = removeDuplicateStr(ipaddr)
			//config := Config{DataBase: DatabaseConfig{Hostt: ipaddr, Port: "10015"}}
			//data, _ = toml.Marshal(&config)
		} else {
			data = []byte("Invalid") //Send Client
			_, err := c.Write(data)
			if err != nil {
				log.Fatal("Write: ", err)
			}
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
