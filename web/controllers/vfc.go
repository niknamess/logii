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

var reader = bufio.NewReader(os.Stdin)

// RunHTTP run http api
func VFC(port string, ctx context.Context) (err error) {

	fmt.Println("Start VFC")
	addr := ":" + port
	//dir := "/home/maxxant/Documents/log"
	dir := "./tmcs-log-agent-storage/"
	//dir := "/home/nik/projects/Course/tmcs-log-agent-storage/"

	var listener net.Listener
	//var err error
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
	go func() {
		if err := srv.Serve(listener); err != nil {
			fmt.Println("Http serve error", err)
		}
	}()
	<-ctx.Done()
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

/* func Tmain() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		cancel()
		fmt.Println("stop VFC")
	}()
	if err := VFC("10015", ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}

} */

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
