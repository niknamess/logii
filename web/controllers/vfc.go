package controllers

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/shurcooL/httpfs/union"
	"github.com/spf13/afero"
)

var (
	fileName    string
	fullURLFile string
)

// RunHTTP run http api
func VFC(port string) string {
	//strconv.Itoa(port)
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
	return "running VFC" + " port: " + addr
}
