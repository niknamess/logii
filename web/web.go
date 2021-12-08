package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gitlab.topaz-atcs.com/tmcs/logi2/generate_logs"

	"github.com/alecthomas/kingpin"
	"github.com/gorilla/mux"
	"gitlab.topaz-atcs.com/tmcs/logi2/logenc"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/controllers"
	"gitlab.topaz-atcs.com/tmcs/logi2/web/util"
)

var (
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./view").ExistingFilesOrDirs()
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("/home/nik/projects/Course/logi2/repdata/").ExistingFilesOrDirs()
	//dir = kingpin.Arg("dir", "Directory path(s) to look for files").Default("/home/nik/projects/Course/tmcs-log-agent-storage/").ExistingFilesOrDirs()
	dir            = kingpin.Arg("dir", "Directory path(s) to look for files").Default("./repdata").ExistingFilesOrDirs()
	port           = kingpin.Flag("port", "Port number to host the server").Short('p').Default("15000").Int()
	cron           = kingpin.Flag("cron", "configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]").Short('t').Default("0h").String()
	cert           = kingpin.Flag("Test", "Test").Short('c').Default("").String()
	missadr string = "nope"
)

func ProcWeb(dir1 string) {
	generate_logs.Remove("./testsave/", "gen_logs_coded")

	//util.GetFiles("localhost", "10015")
	logenc.CreateDir("./repdata/", "")
	logenc.CreateDir("./testsave/", "")

	kingpin.Parse()

	err := util.ParseConfig(*dir, *cron, *cert) //INDEXING FILE

	if err != nil {
		panic(err)
	}
	//go util.DeleteFile90("./repdata")
	go func() {
		time.Sleep(time.Second * 55)
		util.DiskInfo("./repdata")
	}()

	go Loop("192.168.0.193", "10015")
	//go Loop("192.168.0.214", "10015")
	//go Loop("192.168.0.213", "10015")

	router := mux.NewRouter()
	router.HandleFunc("/ws/{b64file}", Use(controllers.WSHandler)).Methods("GET")
	router.HandleFunc("/", Use(controllers.RootHandler)).Methods("GET")
	router.HandleFunc("/searchproject", controllers.SearchHandler)
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

func Loop(address string, port string) {

	for {
		if missadr != address {
			err := util.GetFiles(address, port)
			if err != nil {
				log.Println(err)
				fmt.Println(address)
				missadr = address
			}

			err = util.ParseConfig(*dir, *cron, *cert) //INDEXING FILE

			if err != nil {
				log.Println("LOOP", err)
				panic(err)
			}
		} else {
			fmt.Println(missadr)
			go reconect()
			time.Sleep(time.Second * 10)
			continue

		}

	}
}
func reconect() {
	for range time.Tick(time.Second * 1) {
		missadr = "nope"
		//time.Sleep(100 * time.Second) //time reconect
	}

}
