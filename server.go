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
)

var str string = "VV"

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
			fmt.Print("VFC")
			go controllers.VFC("10015")
		}
		if s == "WEB" {
			terminal.CallClear()
			go web.ProcWeb("")
		}

		fmt.Print("Message Received:", string(message))
	}
}
