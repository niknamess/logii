package controllers

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/shurcooL/httpfs/union"
	"github.com/spf13/afero"
)

var (
	fileName    string
	fullURLFile string
)

var reader = bufio.NewReader(os.Stdin)

// RunHTTP run http api
func VFC(port string) error {

	fmt.Println("Start VFC")
	//strconv.Itoa(port)
	input := make(chan rune, 1)
	go stop(input)
	addr := ":" + port

	//dir := "/home/maxxant/Documents/log"
	dir := "./genrlogs./"
	//dir := "/home/nik/projects/Course/tmcs-log-agent-storage/"

	var listener net.Listener
	var err error
	listenErr := 0

	// wait for listening started
	for ok := false; !ok; {

		listener, err = net.Listen("tcp", addr)
		if err != nil {
			if listenErr == 0 {
				fmt.Println(err)
			}
			listenErr++
			time.Sleep(time.Second * 3)
		}
		ok = (err == nil)

		if ok {
			defer listener.Close()
		}
		//bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	//fmt.Println("listen ok: ", addr)

	fsbase := afero.NewBasePathFs(afero.NewOsFs(), dir)
	fsInput := afero.NewReadOnlyFs(fsbase)
	fsRoot := union.New(map[string]http.FileSystem{
		"/data": afero.NewHttpFs(fsInput),
	})

	router := mux.NewRouter()

	fileserver := http.FileServer(fsRoot)
	router.PathPrefix("/vfs/").Handler(http.StripPrefix("/vfs/", fileserver))
	//fmt.Println("running VFC" + " port: " + addr)
	//fmt.Println("Run new terminal for use service")

	srv := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.Serve(listener); err != nil {
		fmt.Println("Http serve error", err)
	}
	fmt.Printf("Input : %v\n", input)
	if input != nil {
		return srv.Shutdown(context.Background())
	}

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}
	return srv.Shutdown(context.Background())
}

func stop(input chan rune) {
	char, _, err := reader.ReadRune()
	if err != nil {
		log.Fatal(err)
	}
	input <- char
}

/*

ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	} */

/* 	go func() {
    for range time.Tick(time.Second) {
        select {
        case <- request.Context().Done():
            fmt.Println("request is outgoing")
            return
        default:
            fmt.Println("Current request is in progress")
        }
    }
}() */
