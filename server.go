// Very basic socket server
// https://golangr.com/

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gitlab.topaz-atcs.com/tmcs/logi2/web"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var mail string = "Succes"

var ctxVFC, cancelVFC = context.WithCancel(context.Background())
var ctxWEB, cancelWEB = context.WithCancel(context.Background())

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
			MesToClient(c, "Выбрана служба vfc\n")
			go controllers.VFC("10015", ctxVFC)
		}
		if s == "WEB" {
			MesToClient(c, "Выбрана служба web\n")
			allip := enterIp(c)
			go web.ProcWeb("-x", allip, ctxWEB)
		}
		if s == "STOPWEB" {
			MesToClient(c, "Остановыка службы web\n")
			go func() {
				cancelWEB()
				fmt.Println("stop WEB")
				ctxWEB, cancelWEB = context.WithCancel(context.Background())
			}()

		}
		if s == "STOPVFC" {
			MesToClient(c, "Остановыка службы vfc\n")
			go func() {
				cancelVFC()
				fmt.Println("stop VFC")
				ctxVFC, cancelVFC = context.WithCancel(context.Background())
			}()
			//cancel()
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

func Server() string {
	fmt.Println("Server start")
	go func() {
		log.Println("App running, press CTRL + C to stop")
		select {}
	}()

	files, err := ioutil.ReadDir("/tmp/")
	if err != nil {

		log.Fatal(err)
	}

	for _, f := range files {
		if f.Name() == "echo.sock" {
			//e := os.Remove("/tmp/echo.sock")
			log.Fatal("FIND echo.sock ")

		}
	}

	l, err := net.Listen("unix", "/tmp/echo.sock")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func(c chan os.Signal) {
		// Wait for a SIGINT or SIGKILL:
		sig := <-c
		log.Printf("Caught signal %s: shutting down.", sig)
		// Stop listening (and unlink the socket if unix type):
		l.Close()
		// And we're done:
		os.Exit(0)
	}(sigc)
	for {
		fd, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		//shutdown.Listen()
		go echoServer(fd)

	}

}

func enterIp(c net.Conn) []string {

	data := []byte("Input ip address for running service:\n Enter \"stop\" to run service") //Send Client
	_, err := c.Write(data)
	if err != nil {
		log.Fatal("Write: ", err)
	}
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
			ipaddr = removeDuplicateStr(ipaddr)

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

func MesToClient(c net.Conn, message string) {
	data := []byte(message + "\n") //Send Client
	_, err := c.Write(data)
	if err != nil {
		log.Fatal("Write: ", err)
	}

}
