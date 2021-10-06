package web

import (
	//"crypto/tls"
	"fmt"
	"net/http"

	"github.com/alecthomas/kingpin"
	//"github.com/gorilla/csrf"

	"github.com/gorilla/mux"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var (
	dir  = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./logtest/test/").ExistingFilesOrDirs()
	port = kingpin.Flag("port", "Port number to host the server").Short('p').Default("15000").Int()
	cron = kingpin.Flag("cron", "configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]").Short('t').Default("0h").String()
	cert = kingpin.Flag("Test", "Test").Short('c').Default("").String()

	//search string
)

func ProcWeb(dir1 string) {
	kingpin.Parse()
	err := util.ParseConfig(*dir, *cron, *cert)

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/ws/{b64file}", Use(controllers.WSHandler)).Methods("GET")
	router.HandleFunc("/", Use(controllers.RootHandler)).Methods("GET")
	router.HandleFunc("/searchproject", controllers.SearchHandler)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))
	//router.PathPrefix("/").Handler(http.Handler())

	//router.HandleFunc("/", Use(controllers.RootHandler)).Methods("GET")
	//search := "32 "

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("/")
		http.ServeFile(w, r, "index.tmpl")
	})

	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", *port), Handler: router}

	panic(server.ListenAndServe())
} // else {

// Use - Stacking middlewares
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	//fmt.Println("zzz")
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}

//func searchHandler(w http.ResponseWriter, r *http.Request) {
//	search := r.URL.Query().Get("search_string")
//	println(search)
//}
