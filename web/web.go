package web

import (
	//"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/alecthomas/kingpin"
	//"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var (
	dir  = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./logtest/test/").ExistingFilesOrDirs()
	port = kingpin.Flag("port", "Port number to host the server").Short('p').Default("15000").Int()
	cron = kingpin.Flag("cron", "configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]").Short('t').Default("0h").String()
	cert = kingpin.Flag("Test", "Test").Short('c').Default("").String()
)

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%+v\n", r)
	fmt.Fprintln(w, "OK")
}

func ProcWeb(dir1 string) {
	kingpin.Parse()
	err := util.ParseConfig(*dir, *cron, *cert)

	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/ws/{b64file}", Use(controllers.WSHandler, controllers.GetContext)).Methods("GET")
	router.HandleFunc("/", Use(controllers.RootHandler, controllers.GetContext)).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static/")))
	search := "32 "

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("/")
		http.ServeFile(w, r, "index.tmpl")
	})
	router.HandleFunc("/searchproject", Use(searchHandler)).Methods("GET")

	print(search)
	//http.ListenAndServe(":8000", nil)

	csrfRouter := Use((router).ServeHTTP, controllers.CSRFExceptions)

	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", *port), Handler: handlers.CombinedLoggingHandler(os.Stdout, csrfRouter)}
	panic(server.ListenAndServe())
} // else {

// Use - Stacking middlewares
func Use(handler http.HandlerFunc, mid ...func(http.Handler) http.HandlerFunc) http.HandlerFunc {
	fmt.Println("zzz")
	for _, m := range mid {
		handler = m(handler)
	}
	return handler
}
