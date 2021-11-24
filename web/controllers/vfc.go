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

// RunHTTP run http api
func VFC(port string) {
	//strconv.Itoa(port)
	addr := ":" + port

	//dir := "/home/maxxant/Documents/log"
	dir := "./repdata"

	//fs http.FileSystem

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

	fmt.Println("listen ok: ", addr)

	fsbase := afero.NewBasePathFs(afero.NewOsFs(), dir)
	fsInput := afero.NewReadOnlyFs(fsbase)
	//fsInput := afero.NewRegexpFs(fsro, regexp.MustCompile(`\.ads$`))

	// fs0 := httpfs.New(mapfs.New(map[string]string{
	// 	"zzz-last-file.txt":   "It should be visited last.",
	// 	"a-file.txt":          "It` has stuff.",
	// 	"another-file.txt":    "Also stuff.",
	// 	"folderA/entry-A.txt": "Alpha.",
	// 	"folderA/entry-B.txt": "Beta.",
	// }))

	fsRoot := union.New(map[string]http.FileSystem{
		// "/fs0":   fs0,
		"/data": afero.NewHttpFs(fsInput),
	})

	router := mux.NewRouter()

	//fileserver := http.FileServer(fs.Dir("/"))
	fileserver := http.FileServer(fsRoot)
	router.PathPrefix("/vfs/").Handler(http.StripPrefix("/vfs/", fileserver))
	fmt.Println("running /vfs")

	// router.Handle("/debug/vars", http.DefaultServeMux)
	//router.Handle("/debug/metrics", exp.ExpHandler(metrics.DefaultRegistry))

	srv := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	if err := srv.Serve(listener); err != nil {
		fmt.Println("Http serve error", err)
	}
}
