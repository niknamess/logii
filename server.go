// Very basic socket server
// https://golangr.com/

package main

import (
	"bufio"
	"fmt"
	"net"

	"gitlab.topaz-atcs.com/tmcs/logi2/web"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
)

func Server() {
	fmt.Println("Start server...")
	fmt.Println("Pres \"VFC\" or \"WEB\"")

	// listen on port 8000
	ln, _ := net.Listen("tcp", ":8008")

	// accept connection
	conn, _ := ln.Accept()

	// run loop forever (or until ctrl-c)
	for {
		// get message, output
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print(message)
		if message == "VFC" {
			fmt.Print("VFC")
			controllers.VFC("10015")
		}
		if string(message) == "WEB" {
			web.ProcWeb("15000")
		}

		fmt.Print("Message Received:", string(message))
	}
}
