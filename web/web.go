package web

import (
	"fmt"
	"net/http"

	"github.com/alecthomas/kingpin"
	"github.com/gorilla/mux"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var (
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./view").ExistingFilesOrDirs()
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("/home/nik/projects/Course/logi2/repdata/").ExistingFilesOrDirs()
	dir  = kingpin.Arg("dir", "Directory path(s) to look for files").Default("/home/nik/projects/Course/tmcs-log-agent-storage/").ExistingFilesOrDirs()
	port = kingpin.Flag("port", "Port number to host the server").Short('p').Default("15000").Int()
	cron = kingpin.Flag("cron", "configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]").Short('t').Default("0h").String()
	cert = kingpin.Flag("Test", "Test").Short('c').Default("").String()
)

func ProcWeb(dir1 string) {
	kingpin.Parse()
	err := util.ParseConfig(*dir, *cron, *cert)

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/ws/", Use(controllers.BodyHandler)).Methods("POST")
	router.HandleFunc("/ws/{b64file}", Use(controllers.WSHandler)).Methods("GET")
	router.HandleFunc("/", Use(controllers.RootHandler)).Methods("GET")
	router.HandleFunc("/searchproject", controllers.SearchHandler)
	//router.HandleFunc( controllers.SearchHandler)
	//router.HandleFunc("/ws/", Use(controllers.BodyHandler)).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.tmpl")
	})

	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", *port), Handler: router}
	panic(server.ListenAndServe())
}

// Use - Stacking middlewares
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
